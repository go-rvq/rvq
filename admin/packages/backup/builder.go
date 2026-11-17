// Package backup is the admin of the backup (x/packages/backup): the model
// of its states — asked by a form, made by a job of the worker (or by the
// schedule), shown with their status while they are made —, the permission
// of downloading them, and the routes of the client (ClientHandler: an
// access key, then a token of the session).
package backup

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"errors"
	"fmt"
	"io/fs"
	"log/slog"
	"net/http"
	"os"
	"slices"
	"strings"
	"sync"
	"time"

	"github.com/go-rvq/rvq/admin/model"
	"github.com/go-rvq/rvq/admin/presets"
	"github.com/go-rvq/rvq/admin/worker"
	"github.com/go-rvq/rvq/web"
	"github.com/go-rvq/rvq/x/login"
	"github.com/go-rvq/rvq/x/packages/backup"
	"gorm.io/gorm"
)

// The jobs of the backup.
const (
	JobBuildState = "backup_build_state"
	JobScheduled  = "backup_scheduled"
)

// ModelID is the model of the states: its URI and its menu's name.
const ModelID = "backup"

// Builder is the admin of the backup.
type Builder struct {
	db       *gorm.DB
	core     *backup.Builder
	wb       *worker.Builder
	p        *presets.Builder
	mb       *presets.ModelBuilder
	keys     login.KeyAuthenticator
	secret   []byte
	schedule string
	keep     int
	tokenTTL time.Duration
	log      *slog.Logger

	mu      sync.Mutex
	running map[string]context.CancelFunc

	manifests sync.Map
	usage     usage
	// allowDownload replaces the permission's check, clientBase the client's
	// path (tests)
	allowDownload func(r *http.Request) bool
	clientBase    string
}

// New is the admin of the states core makes; their jobs on wb.
func New(db *gorm.DB, core *backup.Builder, wb *worker.Builder) *Builder {
	return &Builder{db: db, core: core, wb: wb, schedule: "@daily", keep: 7, tokenTTL: time.Hour,
		log: web.NewLogger("backup"), running: map[string]context.CancelFunc{}}
}

// KeyAuth is how the client's access keys are authenticated (the
// accesskeys' Authenticator).
func (b *Builder) KeyAuth(a login.KeyAuthenticator) *Builder {
	b.keys = a
	return b
}

// Secret signs the tokens of the sessions (HMAC); none: one made at the
// start — the tokens die with the process (a client opens a session again).
func (b *Builder) Secret(s []byte) *Builder {
	b.secret = s
	return b
}

// Schedule is when a state is made by itself (a spec of the cron:
// "@daily"); "": never.
func (b *Builder) Schedule(spec string) *Builder {
	b.schedule = spec
	return b
}

// Keep is how many states ready the retention keeps (those kept aside).
func (b *Builder) Keep(n int) *Builder {
	b.keep = n
	return b
}

// TokenTTL is how long a token of a session is valid.
func (b *Builder) TokenTTL(d time.Duration) *Builder {
	b.tokenTTL = d
	return b
}

// Core is the maker of the states.
func (b *Builder) Core() *backup.Builder { return b.core }

// ModelBuilder is the model of the states.
func (b *Builder) ModelBuilder() *presets.ModelBuilder { return b.mb }

// AutoMigrate makes the tables of the states and of their sessions.
func (b *Builder) AutoMigrate() error {
	return b.db.AutoMigrate(&BackupState{}, &BackupSession{})
}

// Install registers the model, its jobs, the permission of downloading, the
// download of the manifest by the admin; the states a restart interrupted
// are failed.
func (b *Builder) Install(p *presets.Builder) error {
	b.p = p
	ConfigureMessages(p.I18n())
	if len(b.secret) == 0 {
		b.secret = make([]byte, 32)
		if _, err := rand.Read(b.secret); err != nil {
			return err
		}
	}
	if err := b.AutoMigrate(); err != nil {
		return err
	}
	if err := b.recover(); err != nil {
		return err
	}
	b.mb = p.Model(&BackupState{}, presets.ModelConfig().SetId(ModelID).SetModuleKey(MessagesKey)).
		Label("BackupState").SetPluralLabel("BackupStates").MenuIcon("mdi-backup-restore")
	b.configureModel(b.mb)
	b.configureJobs()
	p.MuxSetup(func(prefix string, mux *http.ServeMux) {
		// the manifest, to the admin (its session)
		mux.HandleFunc("GET "+strings.TrimSuffix(prefix, "/")+"/!backup-manifest/{id}", b.adminManifest)
	})
	go b.flushLoop()
	return nil
}

// recover fails the states being made when the application stopped (their
// folders beside removed by the core).
func (b *Builder) recover() error {
	if err := b.core.Clean(); err != nil {
		b.log.Warn("cleaning the states left", "err", err)
	}
	now := time.Now()
	err := b.db.Model(&BackupState{}).Where("status = ?", StatusBuilding).Updates(map[string]any{
		"status": StatusFailed, "error": Messages_en_US.Interrupted, "finished_at": now,
	}).Error
	if err != nil {
		return err
	}
	// ready, but their files gone (another machine, a folder removed): not
	// downloadable
	var ready []BackupState
	if err := b.db.Select("id").Where("status = ?", StatusReady).Find(&ready).Error; err != nil {
		return err
	}
	for _, s := range ready {
		if _, err := os.Stat(b.core.ManifestPath(s.ID)); errors.Is(err, fs.ErrNotExist) {
			if err := b.db.Model(&BackupState{}).Where("id = ?", s.ID).Updates(map[string]any{
				"status": StatusFailed, "error": Messages_en_US.Gone}).Error; err != nil {
				return err
			}
		}
	}
	return nil
}

func (b *Builder) flushLoop() {
	t := time.NewTicker(15 * time.Second)
	defer t.Stop()
	for range t.C {
		if err := b.usage.flush(b.db); err != nil {
			b.log.Warn("the usage of the sessions", "err", err)
		}
	}
}

// NewID is the id of a state made now: 20261117-153000-a1b2.
func NewID(now time.Time) string {
	var r [2]byte
	_, _ = rand.Read(r[:])
	return now.UTC().Format("20060102-150405") + "-" + hex.EncodeToString(r[:])
}

// jobArgs is the argument of the job that makes a state.
type jobArgs struct {
	StateID string
}

func (b *Builder) configureJobs() {
	b.wb.NewJob(JobBuildState).
		Title(func(ctx *web.EventContext) string { return GetMessages(ctx.Context()).JobBuildState }).
		Resource(&jobArgs{}).
		Handler(func(ctx context.Context, job worker.JobInterface) error {
			info, err := job.GetJobInfo()
			if err != nil {
				return err
			}
			a, ok := info.Argument.(*jobArgs)
			if !ok {
				return errors.New("backup: the job has no state")
			}
			return b.Build(ctx, a.StateID)
		})
	if b.schedule == "" {
		return
	}
	b.wb.NewJob(JobScheduled).
		Title(func(ctx *web.EventContext) string { return GetMessages(ctx.Context()).JobScheduled }).
		System(true).
		CronConfig(worker.JobCronConfig{Spec: b.schedule, Once: true}).
		Handler(func(ctx context.Context, job worker.JobInterface) error {
			s := &BackupState{ID: NewID(time.Now()), Status: StatusQueued, Origin: OriginScheduled}
			if err := b.db.Create(s).Error; err != nil {
				return err
			}
			job.AddLog("state " + s.ID)
			return b.Build(ctx, s.ID)
		})
}

// Build makes the state id, queued: building, then ready or failed (or
// canceled: Cancel); the retention after it.
func (b *Builder) Build(ctx context.Context, id string) error {
	var s BackupState
	if err := b.db.First(&s, "id = ?", id).Error; err != nil {
		return err
	}
	if s.Status != StatusQueued {
		return nil // canceled while queued, or made
	}
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()
	b.mu.Lock()
	b.running[id] = cancel
	b.mu.Unlock()
	defer func() {
		b.mu.Lock()
		delete(b.running, id)
		b.mu.Unlock()
	}()

	started := time.Now()
	if err := b.db.Model(&s).Updates(map[string]any{"status": StatusBuilding, "started_at": started}).Error; err != nil {
		return err
	}
	var last time.Time
	st := &backup.State{ID: s.ID, Created: s.CreatedAt, Progress: func(p backup.Progress) {
		// the progress, a second at most
		if time.Since(last) < time.Second {
			return
		}
		last = time.Now()
		_ = b.db.Model(&BackupState{}).Where("id = ?", s.ID).Updates(map[string]any{
			"progress_source": p.Source, "progress_files": p.Files, "progress_bytes": p.Bytes}).Error
	}}
	if len(s.Sources) > 0 {
		st.Sources = s.Sources
	}
	err := b.core.BuildState(ctx, st)
	finished := time.Now()
	up := map[string]any{"finished_at": finished, "progress_source": ""}
	switch {
	case err != nil && ctx.Err() != nil:
		up["status"], up["error"] = StatusCanceled, ""
	case err != nil:
		up["status"], up["error"] = StatusFailed, err.Error()
	default:
		r := st.Result
		up["status"] = StatusReady
		up["files"], up["dirs"], up["links"], up["bytes"], up["partial"] = r.Files, r.Dirs, r.Links, r.Bytes, r.Partial
		up["manifest_size"], up["manifest_sum"] = r.ManifestSize, r.ManifestSum.String()
		s.Results, s.Meta = r.Sources, r.Meta
	}
	if uerr := b.db.Model(&BackupState{}).Where("id = ?", s.ID).Updates(up).Error; uerr != nil {
		return uerr
	}
	if err == nil {
		// the json columns, by Save of their fields (Updates of a map would
		// not serialize them)
		if uerr := b.db.Model(&s).Select("Results", "Meta").Updates(&s).Error; uerr != nil {
			return uerr
		}
		if rerr := b.Retain(); rerr != nil {
			b.log.Warn("retention", "err", rerr)
		}
	}
	if ctx.Err() != nil {
		return nil // canceled: not the job's failure
	}
	return err
}

// Cancel cancels the state id: queued, it is not made; being made here, it
// stops (its folder removed).
func (b *Builder) Cancel(id string) error {
	b.mu.Lock()
	cancel := b.running[id]
	b.mu.Unlock()
	if cancel != nil {
		cancel()
		return nil
	}
	res := b.db.Model(&BackupState{}).Where("id = ? AND status = ?", id, StatusQueued).
		Updates(map[string]any{"status": StatusCanceled, "finished_at": time.Now()})
	if res.Error != nil {
		return res.Error
	}
	if res.RowsAffected == 0 {
		return fmt.Errorf("backup: state %s is not queued nor being made here", id)
	}
	return nil
}

// Remove removes the state id: its folder, its sessions, its record. One
// being made is not removed (ErrInProgress).
func (b *Builder) Remove(id string) error {
	var s BackupState
	if err := b.db.First(&s, "id = ?", id).Error; err != nil {
		return err
	}
	if s.Status == StatusBuilding {
		return ErrInProgress
	}
	if err := b.core.Remove(id); err != nil {
		return err
	}
	b.manifests.Delete(id)
	return b.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Where("state_id = ?", id).Delete(&BackupSession{}).Error; err != nil {
			return err
		}
		return tx.Delete(&BackupState{}, "id = ?", id).Error
	})
}

// ErrInProgress: a state being made is not removed.
var ErrInProgress = errors.New("backup: the state is being made")

// Retain removes the states beyond Keep: of those ready, the newest kept;
// of those failed or canceled, as many; the kept (Keep) never.
func (b *Builder) Retain() error {
	if b.keep <= 0 {
		return nil
	}
	var errs []error
	for _, statuses := range [][]string{{StatusReady}, {StatusFailed, StatusCanceled}} {
		var old []BackupState
		err := b.db.Where("status IN ? AND keep = ?", statuses, false).Order("created_at DESC").
			Offset(b.keep).Find(&old).Error
		if err != nil {
			return err
		}
		for _, s := range old {
			errs = append(errs, b.Remove(s.ID))
		}
	}
	return errors.Join(errs...)
}

// create makes the state asked by ctx's user and queues its job.
func (b *Builder) create(s *BackupState, ctx *web.EventContext) error {
	all := b.core.Sources()
	for _, src := range s.Sources {
		if !slices.Contains(all, src) {
			return fmt.Errorf("backup: no source %q", src)
		}
	}
	if len(s.Sources) == len(all) {
		s.Sources = nil // all
	}
	s.ID, s.Status, s.Origin = NewID(time.Now()), StatusQueued, OriginManual
	if id, name := userOf(login.CurrentUserFromRequest(ctx.R)); name != "" {
		s.CreatedByID, s.CreatedByName = model.UUID(id), name
	}
	if err := b.db.Create(s).Error; err != nil {
		return err
	}
	_, err := b.wb.CreateJob(ctx, JobBuildState, true, &jobArgs{StateID: s.ID})
	return err
}

// adminManifest serves the manifest of a state to the admin (its session),
// if it may download.
func (b *Builder) adminManifest(w http.ResponseWriter, r *http.Request) {
	if !b.mayDownload(r) {
		http.Error(w, "forbidden", http.StatusForbidden)
		return
	}
	s, ok := b.readyState(w, r.PathValue("id"))
	if !ok {
		return
	}
	b.sendManifest(w, r, s)
}
