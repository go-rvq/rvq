package worker

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"net/url"
	"reflect"
	"strings"
	"time"

	. "github.com/go-rvq/htmlgo"
	"github.com/go-rvq/rvq/admin/activity"
	"github.com/go-rvq/rvq/admin/model"
	"github.com/go-rvq/rvq/admin/presets"
	"github.com/go-rvq/rvq/admin/utils/uuidkey"
	"github.com/go-rvq/rvq/web"
	"github.com/go-rvq/rvq/x/i18n"
	"github.com/go-rvq/rvq/x/perm"
	. "github.com/go-rvq/rvq/x/ui/vuetify"
	"github.com/go-rvq/rvq/x/ui/vuetifyx"
	"github.com/google/uuid"
	rcron "github.com/robfig/cron/v3"
	"gorm.io/gorm"
)

type NewOption func(b *Builder)

func WithQueue(q Queue) NewOption {
	return func(b *Builder) {
		b.q = q
	}
}

func WithCron(cron *rcron.Cron) NewOption {
	return func(b *Builder) {
		b.cron = cron
	}
}

func WithCronLogger(log *slog.Logger) NewOption {
	return func(b *Builder) {
		b.cronLogger = log
	}
}

type Builder struct {
	db                   *gorm.DB
	q                    Queue
	jpb                  *presets.Builder // for render job form
	pb                   *presets.Builder
	jbs                  []*JobBuilder
	mb                   *presets.ModelBuilder
	getCurrentUserIDFunc func(r *http.Request) string
	ab                   *activity.Builder
	cron                 *rcron.Cron
	cronContext          *web.EventContext
	cronLogger           *slog.Logger
}

func New(i18nB *i18n.Builder, db *gorm.DB, option ...NewOption) *Builder {
	uuidkey.MustRegister(db) // its records have UUID keys
	if db == nil {
		panic("db can not be nil")
	}

	if err := renameLegacy(db); err != nil {
		panic(err)
	}

	err := db.AutoMigrate(&Job{}, &JobInstance{}, &JobLog{}, &GoQueError{})
	if err != nil {
		panic(err)
	}

	r := &Builder{
		db:  db,
		jpb: presets.New(i18nB),
	}

	for _, opt := range option {
		opt(r)
	}

	if r.q == nil {
		r.q = NewGoQueQueue(db)
	}

	if r.cronLogger == nil {
		r.cronLogger = web.NewLogger("worker").WithGroup("cron")
	}

	req := (&http.Request{
		Header: http.Header{},
	}).WithContext(context.Background())

	r.cronContext = &web.EventContext{
		R: req,
	}

	return r
}

// default queue is go-que queue
func (b *Builder) Queue(q Queue) *Builder {
	b.q = q
	return b
}

func (b *Builder) GetCurrentUserIDFunc(f func(r *http.Request) string) *Builder {
	b.getCurrentUserIDFunc = f
	return b
}

// Activity sets Activity Builder to log activities
func (b *Builder) Activity(ab *activity.Builder) *Builder {
	b.ab = ab
	return b
}

func (b *Builder) NewJob(name string) *JobBuilder {
	for _, jb := range b.jbs {
		if jb.name == name {
			panic(fmt.Sprintf("worker %s already exists", name))
		}
	}

	j := newJob(b, name)
	b.jbs = append(b.jbs, j)

	return j
}

func (b *Builder) getJobBuilder(name string) *JobBuilder {
	for _, jb := range b.jbs {
		if jb.name == name {
			return jb
		}
	}

	return nil
}

func (b *Builder) mustGetJobBuilder(name string) *JobBuilder {
	jb := b.getJobBuilder(name)

	if jb == nil {
		panic(fmt.Sprintf("no job %s", name))
	}

	return jb
}

func (b *Builder) getJobBuilderByJobID(id uuid.UUID) (*JobBuilder, error) {
	j := Job{}
	err := b.db.Where("id = ?", id).First(&j).Error
	if err != nil {
		return nil, err
	}

	return b.getJobBuilder(j.Job), nil
}

func (b *Builder) setStatus(id uuid.UUID, status string) error {
	return b.db.Model(&Job{}).Where("id = ?", id).
		Updates(map[string]interface{}{
			"status": status,
		}).
		Error
}

func (b *Builder) ModelBuilder() *presets.ModelBuilder {
	return b.mb
}

var permVerifier *perm.Verifier

func (b *Builder) URI() string {
	return b.mb.Info().ListingHref()
}

func (b *Builder) Install(pb *presets.Builder) error {
	b.pb = pb
	permVerifier = perm.NewVerifier("workers", pb.GetPermission())

	ConfigureMessages(pb.I18n())

	mb := pb.Model(&Job{}, presets.ModelConfig().SetModuleKey(MessagesKey)).
		Label("Workers").
		URIName("workers").
		MenuIcon("mdi-briefcase")

	b.mb = mb
	mb.RegisterEventFunc(EventSelectJob, b.eventSelectJob)
	mb.RegisterEventFunc(EventAbortJob, b.eventAbortJob)
	mb.RegisterEventFunc(EventRerunJob, b.eventRerunJob)
	mb.RegisterEventFunc(EventUpdateJob, b.eventUpdateJob)
	mb.RegisterEventFunc(EventUpdateJobProgressing, b.eventUpdateJobProgressing)
	mb.RegisterEventFunc(EventLoadHiddenLogs, b.eventLoadHiddenLogs)
	mb.RegisterEventFunc(ActionJobInputParams, b.eventActionJobInputParams)
	mb.RegisterEventFunc(ActionJobCreate, b.eventActionJobCreate)
	mb.RegisterEventFunc(ActionJobResponse, b.eventActionJobResponse)
	mb.RegisterEventFunc(ActionJobClose, b.eventActionJobClose)
	mb.RegisterEventFunc(ActionJobProgressing, b.eventActionJobProgressing)

	lb := mb.Listing("ID", "Job", "Status", "CreatedAt")
	lb.RowMenu().Empty()

	lb.FilterDataFunc(func(ctx *web.EventContext) vuetifyx.FilterData {
		var (
			i18m = i18n.MustGetModuleMessages(ctx.Context(), MessagesKey, Messages_en_US).(*Messages)
			m    = GetMessages(ctx.Context())
			jobs []*vuetifyx.SelectItem
		)

		for _, jb := range b.jbs {
			jobs = append(jobs, &vuetifyx.SelectItem{Value: jb.name, Text: jb.GetTitle(ctx)})
		}

		return []*vuetifyx.FilterItem{
			{
				Key:          "status",
				Label:        "Status",
				ItemType:     vuetifyx.ItemTypeMultipleSelect,
				SQLCondition: `status %s ?`,
				Options: []*vuetifyx.SelectItem{
					{Text: i18m.StatusNew, Value: JobStatusNew},
					{Text: i18m.StatusScheduled, Value: JobStatusScheduled},
					{Text: i18m.StatusRunning, Value: JobStatusRunning},
					{Text: i18m.StatusCancelled, Value: JobStatusCancelled},
					{Text: i18m.StatusDone, Value: JobStatusDone},
					{Text: i18m.StatusException, Value: JobStatusException},
					{Text: i18m.StatusKilled, Value: JobStatusKilled},
				},
			},
			{
				Key:          "job",
				Label:        m.Job,
				ItemType:     vuetifyx.ItemTypeMultipleSelect,
				SQLCondition: `job %s ?`,
				Options:      jobs,
				Invisible:    len(jobs) == 0,
			},
		}
	})
	lb.FilterTabsFunc(func(ctx *web.EventContext) []*presets.FilterTab {
		msgr := i18n.MustGetModuleMessages(ctx.Context(), MessagesKey, Messages_en_US).(*Messages)
		return []*presets.FilterTab{
			{
				Label: msgr.FilterTabAll,
				Query: url.Values{"all": []string{"1"}},
			},
			{
				Label: msgr.FilterTabRunning,
				Query: url.Values{"status": []string{JobStatusRunning}},
			},
			{
				Label: msgr.FilterTabScheduled,
				Query: url.Values{"status": []string{JobStatusScheduled}},
			},
			{
				Label: msgr.FilterTabDone,
				Query: url.Values{"status": []string{JobStatusDone}},
			},
			{
				Label: msgr.FilterTabErrors,
				Query: url.Values{"status": []string{JobStatusException}},
			},
		}
	})
	lb.Field("Job").ComponentFunc(func(field *presets.FieldContext, ctx *web.EventContext) HTMLComponent {
		job := field.Obj.(*Job)
		name := job.Job
		if b := b.getJobBuilder(name); b != nil {
			name = b.GetTitle(field.EventContext)
		}
		return Td(Text(name))
	})
	lb.Field("Status").ComponentFunc(func(field *presets.FieldContext, ctx *web.EventContext) HTMLComponent {
		msgr := i18n.MustGetModuleMessages(ctx.Context(), MessagesKey, Messages_en_US).(*Messages)
		job := field.Obj.(*Job)
		return Td(Text(msgr.GetStatus(job.Status)))
	})

	eb := mb.Editing("Job", "Args")

	eb.Validators.AppendFunc(func(obj interface{}, mode presets.FieldModeStack, ctx *web.EventContext) (err web.ValidationErrors) {
		msgr := i18n.MustGetModuleMessages(ctx.Context(), MessagesKey, Messages_en_US).(*Messages)
		job := obj.(*Job)
		if job.Job == "" {
			err.FieldError("Job", msgr.PleaseSelectJob)
		}

		return err
	})

	eb.Field("Job").ComponentFunc(func(field *presets.FieldContext, ctx *web.EventContext) HTMLComponent {
		job := field.Obj.(*Job)
		return web.Portal(b.jobSelectList(ctx, job.Job)).Name("worker_jobSelectList")
	})

	eb.Field("Args").
		ComponentFunc(func(field *presets.FieldContext, ctx *web.EventContext) HTMLComponent {
			var vErr web.ValidationErrors
			if ve, ok := ctx.Flash.(*web.ValidationErrors); ok {
				vErr = *ve
				if fvErr := vErr.GetFieldErrors(field.Name); len(fvErr) > 0 {
					errM := make(map[string][]string)
					if err := json.Unmarshal([]byte(fvErr[0]), &errM); err == nil {
						for f, es := range errM {
							for _, e := range es {
								ve.FieldError(f, e)
							}
						}
					}
				}
			}

			job := field.Obj.(*Job)
			return web.Portal(b.jobEditingContent(ctx, job.Job, job.Args)).Name("worker_jobEditingContent")
		})

	eb.SaveFunc(func(obj interface{}, id model.ID, ctx *web.EventContext) (err error) {
		job := obj.(*Job)
		if job.Job == "" {
			return errors.New("job is required")
		}
		j, err := b.createJob(ctx, job)
		if err != nil {
			return err
		}
		if b.ab != nil {
			b.ab.AddRecords(activity.ActivityCreate, ctx.R.Context(), j)
		}
		return
	})

	eb.CreateFunc(func(obj interface{}, ctx *web.EventContext) (err error) {
		job := obj.(*Job)
		if job.Job == "" {
			return errors.New("job is required")
		}
		j, err := b.createJob(ctx, job)
		if err != nil {
			return err
		}
		if b.ab != nil {
			b.ab.AddRecords(activity.ActivityCreate, ctx.R.Context(), j)
		}
		return
	})

	mb.Detailing("DetailingPage").Field("DetailingPage").ComponentFunc(func(field *presets.FieldContext, ctx *web.EventContext) HTMLComponent {
		msgr := i18n.MustGetModuleMessages(ctx.Context(), MessagesKey, Messages_en_US).(*Messages)

		job := field.Obj.(*Job)
		inst, err := getModelJobInstance(b.db, job.ID)
		if err != nil {
			return Text(err.Error())
		}

		var scheduledJobDetailing []HTMLComponent
		eURL := b.URI()
		if inst.Status == JobStatusScheduled {
			jb := b.getJobBuilder(job.Job)
			if jb != nil && jb.r != nil {
				args := jb.newResourceObject()
				err := json.Unmarshal([]byte(inst.Args), &args)
				if err != nil {
					return Text(err.Error())
				}
				body := jb.rmb.Editing().ToComponent(&presets.ToComponentOptions{}, args, field.Mode.DotStack(), ctx)
				scheduledJobDetailing = []HTMLComponent{
					body,
					If(editIsAllowed(ctx.R, job.Job) == nil,
						Div().Class("d-flex mt-3").Children(
							VSpacer(),
							VBtn(msgr.ActionCancelJob).Color("error").Class("mr-2").
								Attr("@click", web.Plaid().
									URL(eURL).
									EventFunc(EventAbortJob).
									Query("jobID", fmt.Sprintf("%d", job.ID)).
									Query("job", job.Job).
									Go()),
							VBtn(msgr.ActionUpdateJob).Color("primary").
								Attr("@click", web.Plaid().
									URL(eURL).
									EventFunc(EventUpdateJob).
									Query("jobID", fmt.Sprintf("%d", job.ID)).
									Query("job", job.Job).
									Go()),
						),
					),
				}
			} else {
				scheduledJobDetailing = []HTMLComponent{
					VAlert().Density(DensityCompact).Type("warning").Children(
						Text(msgr.NoticeJobWontBeExecuted),
					),
					Div(Text("args: " + inst.Args)),
				}
			}
		}

		return Div(
			Div(Text(getTJob(ctx.Context(), job.Job))).Class("mb-3 text-h6 font-weight-regular"),
			If(inst.Status == JobStatusScheduled,
				scheduledJobDetailing...,
			).Else(
				web.Scope(
					web.Portal().
						Loader(web.Plaid().EventFunc(EventUpdateJobProgressing).
							URL(eURL).
							Query("jobID", fmt.Sprintf("%d", job.ID)).
							Query("job", job.Job),
						).
						AutoReloadInterval("locals.worker_updateJobProgressingInterval"),
				).LocalsInit("{worker_updateJobProgressingInterval: 2000}"),
			),
			web.Portal().Name("worker_snackbar"),
		)
	})

	if b.ab != nil {
		b.ab.RegisterModel(mb).SkipCreate().SkipUpdate().SkipDelete().
			AddTypeHanders(time.Time{}, func(old, now interface{}, prefixField string) []activity.Diff {
				fm := "2006-01-02 15:04:05"
				oldString := old.(time.Time).Format(fm)
				nowString := now.(time.Time).Format(fm)
				if oldString != nowString {
					return []activity.Diff{
						{Field: prefixField, Old: oldString, Now: nowString},
					}
				}
				return []activity.Diff{}
			}).
			AddTypeHanders(Schedule{}, func(old, now interface{}, prefixField string) []activity.Diff {
				fm := "2006-01-02 15:04:05"
				oldString := old.(Schedule).ScheduleTime.Format(fm)
				nowString := now.(Schedule).ScheduleTime.Format(fm)
				if oldString != nowString {
					return []activity.Diff{
						{Field: prefixField, Old: oldString, Now: nowString},
					}
				}
				return []activity.Diff{}
			})
	}

	return nil
}

func (b *Builder) Listen() {
	var (
		jds     []*JobDefinition
		crons   []*JobBuilder
		newCron = b.cron == nil
	)

	for _, jb := range b.jbs {
		jds = append(jds, &JobDefinition{
			Name:    jb.name,
			Handler: jb.h,
		})

		if jb.cronConfig.Valid() {
			crons = append(crons, jb)
		}
	}

	err := b.q.Listen(jds, func(jobID uuid.UUID) (QueJobInterface, error) {
		jb, err := b.getJobBuilderByJobID(jobID)
		if err != nil {
			return nil, err
		}
		if jb == nil {
			return nil, errors.New("failed to find job (job name modified?)")
		}

		return jb.getJobInstance(jobID)
	})

	if err != nil {
		panic(err)
	}

	if newCron {
		b.cron = rcron.New()
		defer b.cron.Start()
	}

	for _, jb := range crons {
		cfg := jb.cronConfig
		jl := b.cronLogger.With("spec", cfg.Spec, "once", cfg.Once, "job", jb.name)
		jl.Info("add")

		if _, err = b.cron.AddJob(cfg.Spec, rcron.FuncJob(func() {
			if _, err := b.CreateJob(b.cronContext, jb.name, cfg.Once, cfg.Spec); err != nil {
				jl.Error("new instance failed: " + err.Error())
			} else {
				jl.Info("new instance created")
			}
		})); err != nil {
			jl.Error(err.Error())
		}
	}
}

// Shutdown stops the cron — no job is created any more — and the queue's
// workers, which hold connections of their own while they run.
func (b *Builder) Shutdown(ctx context.Context) error {
	if b.cron != nil {
		select {
		case <-b.cron.Stop().Done():
		case <-ctx.Done():
		}
	}
	return b.q.Shutdown(ctx)
}

func (b *Builder) createJob(ctx *web.EventContext, job *Job) (j *Job, err error) {
	if err = editIsAllowed(ctx.R, job.Job); err != nil {
		return
	}

	jb := b.mustGetJobBuilder(job.Job)

	// encode args
	args, vErr := jb.unmarshalForm(ctx)
	if vErr.HaveErrors() {
		errM := make(map[string][]string)
		argsT := reflect.TypeOf(jb.r).Elem()
		for i := 0; i < argsT.NumField(); i++ {
			fName := argsT.Field(i).Name
			errM[fName] = vErr.GetFieldErrors(fName)
		}
		bErrM, _ := json.Marshal(errM)
		err = errors.New(string(bErrM))
		return
	}

	// encode context
	context := make(map[string]interface{})
	for key, v := range DefaultOriginalPageContextHandler(ctx) {
		context[key] = v
	}

	if jb.contextHandler != nil {
		for key, v := range jb.contextHandler(ctx) {
			context[key] = v
		}
	}

	err = b.db.Transaction(func(tx *gorm.DB) error {
		j = &Job{
			Job:    job.Job,
			Status: JobStatusNew,
		}
		err = b.db.Create(j).Error
		if err != nil {
			return err
		}
		var inst *JobInstance
		inst, err = jb.newJobInstance(ctx.R, j.ID, job.Job, job.Once, args, context)
		if err != nil {
			return err
		}
		return b.q.Add(ctx.R.Context(), inst)
	})
	return
}

func (b *Builder) CreateJob(ctx *web.EventContext, name string, once bool, args any) (j *Job, err error) {
	jb := b.mustGetJobBuilder(name)
	if jb == nil {
		err = errors.New("failed to find job (job name modified?)")
		return
	}

	// encode context
	context := make(map[string]interface{})
	for key, v := range DefaultOriginalPageContextHandler(ctx) {
		context[key] = v
	}

	if jb.contextHandler != nil {
		for key, v := range jb.contextHandler(ctx) {
			context[key] = v
		}
	}

	err = b.db.Transaction(func(tx *gorm.DB) error {
		j = &Job{
			Job:    name,
			Status: JobStatusNew,
			Once:   once,
		}
		err = b.db.Create(j).Error
		if err != nil {
			return err
		}
		var inst *JobInstance
		inst, err = jb.newJobInstance(ctx.R, j.ID, name, once, args, context)
		if err != nil {
			return err
		}
		return b.q.Add(ctx.R.Context(), inst)
	})
	return
}

func (b *Builder) eventSelectJob(ctx *web.EventContext) (er web.EventResponse, err error) {
	job := ctx.R.FormValue("jobName")
	er.
		UpdatePortal(
			"worker_jobEditingContent",
			b.jobEditingContent(ctx, job, nil),
		).
		UpdatePortal(
			"worker_jobSelectList",
			b.jobSelectList(ctx, job),
		)

	return
}

func (b *Builder) eventAbortJob(ctx *web.EventContext) (er web.EventResponse, err error) {
	msgr := i18n.MustGetModuleMessages(ctx.Context(), MessagesKey, Messages_en_US).(*Messages)

	jobID := paramJobID(ctx)
	jobName := ctx.R.FormValue("job")

	if pErr := editIsAllowed(ctx.R, jobName); pErr != nil {
		return er, pErr
	}

	jb := b.mustGetJobBuilder(jobName)
	inst, err := jb.getJobInstance(jobID)
	if err != nil {
		return er, err
	}
	isScheduled := inst.Status == JobStatusScheduled

	err = b.doAbortJob(ctx.R.Context(), inst)
	if err != nil {
		_, ok := err.(*cannotAbortError)
		if !ok {
			return er, err
		}
		er.UpdatePortal("worker_snackbar",
			VSnackbar().ModelValue(true).Timeout(3000).Color("warning").Children(
				Text(msgr.NoticeJobCannotBeAborted),
			))
	}

	er.Reload = true
	er.RunScript = "vars.worker_updateJobProgressingInterval = 2000"

	if b.ab != nil {
		action := "Abort"
		if isScheduled {
			action = "Cancel"
		}
		b.ab.AddCustomizedRecord(action, false, ctx.R.Context(), &Job{
			Model: uuidkey.Model{
				ID: inst.JobID,
			},
		})
	}

	return er, nil
}

type cannotAbortError struct {
	err error
}

func (e *cannotAbortError) Error() string {
	return e.err.Error()
}

func (b *Builder) doAbortJob(ctx context.Context, inst *JobInstance) (err error) {
	switch inst.Status {
	case JobStatusRunning:
		return b.q.Kill(ctx, inst)
	case JobStatusNew, JobStatusScheduled:
		return b.q.Remove(ctx, inst)
	default:
		return &cannotAbortError{
			err: fmt.Errorf("job status is %s, cannot be aborted/canceled", inst.Status),
		}
	}
}

func (b *Builder) eventRerunJob(ctx *web.EventContext) (er web.EventResponse, err error) {
	jobID := paramJobID(ctx)
	jobName := ctx.R.FormValue("job")

	if pErr := editIsAllowed(ctx.R, jobName); pErr != nil {
		return er, pErr
	}

	jb := b.mustGetJobBuilder(jobName)
	old, err := jb.getJobInstance(jobID)
	if err != nil {
		return er, err
	}

	if old.Once {
		return er, errors.New("job is not done")
	}

	if old.Status != JobStatusDone {
		return er, errors.New("job is not done")
	}

	inst, err := jb.newJobInstance(ctx.R, jobID, jobName, old.Once, old.Args, old.Context)
	if err != nil {
		return er, err
	}
	err = b.setStatus(jobID, JobStatusNew)
	if err != nil {
		return er, err
	}
	err = b.q.Add(ctx.R.Context(), inst)
	if err != nil {
		return er, err
	}

	er.Reload = true
	er.RunScript = "vars.worker_updateJobProgressingInterval = 2000"

	if b.ab != nil {
		b.ab.AddCustomizedRecord("Rerun", false, ctx.R.Context(), &Job{
			Model: uuidkey.Model{
				ID: inst.JobID,
			},
		})
	}
	return
}

func (b *Builder) eventUpdateJob(ctx *web.EventContext) (er web.EventResponse, err error) {
	msgr := i18n.MustGetModuleMessages(ctx.Context(), MessagesKey, Messages_en_US).(*Messages)

	jobID := paramJobID(ctx)
	jobName := ctx.R.FormValue("job")

	if pErr := editIsAllowed(ctx.R, jobName); pErr != nil {
		return er, pErr
	}

	jb := b.mustGetJobBuilder(jobName)
	newArgs, argsVErr := jb.unmarshalForm(ctx)
	if argsVErr.HaveErrors() {
		return er, errors.New("invalid arguments")
	}

	contexts := make(map[string]interface{})
	for key, v := range DefaultOriginalPageContextHandler(ctx) {
		contexts[key] = v
	}
	if jb.contextHandler != nil {
		for key, v := range jb.contextHandler(ctx) {
			contexts[key] = v
		}
	}

	old, err := jb.getJobInstance(jobID)
	if err != nil {
		return er, err
	}
	oldArgs, _ := jb.parseArgs(old.Args)
	err = b.doAbortJob(ctx.R.Context(), old)
	if err != nil {
		_, ok := err.(*cannotAbortError)
		if !ok {
			return er, err
		}
		er.UpdatePortal("worker_snackbar",
			VSnackbar().ModelValue(true).Timeout(3000).Color("warning").Children(
				Text(msgr.NoticeJobCannotBeAborted),
			))
		er.Reload = true
		return er, nil
	}

	newInst, err := jb.newJobInstance(ctx.R, jobID, jobName, old.Once, newArgs, contexts)
	if err != nil {
		return er, err
	}
	err = b.q.Add(ctx.R.Context(), newInst)
	if err != nil {
		return er, err
	}

	er.Reload = true
	er.RunScript = "vars.worker_updateJobProgressingInterval = 2000"
	if b.ab != nil {
		b.ab.AddEditRecordWithOldAndContext(
			ctx.R.Context(),
			&Job{
				Model: uuidkey.Model{
					ID: newInst.JobID,
				},
				Args: oldArgs,
			},
			&Job{
				Model: uuidkey.Model{
					ID: newInst.JobID,
				},
				Args: newArgs,
			},
		)
	}
	return er, nil
}

func (b *Builder) eventUpdateJobProgressing(ctx *web.EventContext) (er web.EventResponse, err error) {
	msgr := i18n.MustGetModuleMessages(ctx.Context(), MessagesKey, Messages_en_US).(*Messages)

	jobID := paramJobID(ctx)
	jobName := ctx.R.FormValue("job")

	inst, err := getModelJobInstance(b.db, jobID)
	if err != nil {
		return er, err
	}

	canEdit := editIsAllowed(ctx.R, jobName) == nil
	logs := make([]string, 0, 100)
	hasMoreLogs := false
	{
		var count int64
		err = b.db.Model(&JobLog{}).
			Where("job_instance_id = ?", inst.ID).
			Count(&count).
			Error
		if err != nil {
			return er, err
		}
		if count > 100 {
			hasMoreLogs = true
		}
		if count > 0 {
			var mLogs []*JobLog
			err = b.db.Where("job_instance_id = ?", inst.ID).
				Order("created_at desc").
				Limit(100).
				Find(&mLogs).
				Error
			if err != nil {
				return er, err
			}
			for i := len(mLogs) - 1; i >= 0; i-- {
				logs = append(logs, mLogs[i].Log)
			}
		}
	}
	er.Body = b.jobProgressing(canEdit, msgr, jobID, jobName, inst.Status, inst.Progress, logs, hasMoreLogs, inst.ProgressText, presets.ParentsModelID(ctx.R))
	if inst.Status != JobStatusNew && inst.Status != JobStatusRunning && inst.Status != JobStatusKilled {
		er.RunScript = "vars.worker_updateJobProgressingInterval = 0"
	} else {
		er.RunScript = "vars.worker_updateJobProgressingInterval = 2000"
	}
	return er, nil
}

func (b *Builder) eventLoadHiddenLogs(ctx *web.EventContext) (er web.EventResponse, err error) {
	jobID := paramJobID(ctx)
	currentCount := ctx.ParamAsInt("currentCount")

	inst, err := getModelJobInstance(b.db, jobID)
	if err != nil {
		return er, err
	}

	var logs []*JobLog
	err = b.db.Where("job_instance_id = ?", inst.ID).
		Order("created_at desc").
		Offset(currentCount).
		Find(&logs).
		Error
	if err != nil {
		return er, err
	}
	logLines := make([]HTMLComponent, 0, len(logs))
	for i := len(logs) - 1; i >= 0; i-- {
		logLines = append(logLines, P().Style(`
    margin: 0;
    margin-bottom: 4px;`).Children(Text(logs[i].Log)))
	}
	er.UpdatePortal("worker_hiddenLogs", Div(logLines...))
	return er, nil
}

func (b *Builder) jobProgressing(
	canEdit bool,
	msgr *Messages,
	id uuid.UUID,
	job string,
	status string,
	progress uint,
	logs []string,
	hasMoreLogs bool,
	progressText string,
	parentsID presets.IDSlice,
) HTMLComponent {
	logLines := make([]HTMLComponent, 0, len(logs)+1)
	if hasMoreLogs {
		logLines = append(logLines, web.Portal(
			VBtn("Load hidden logs").Attr("@click", web.Plaid().EventFunc(EventLoadHiddenLogs).
				Query("jobID", id).
				Query("currentCount", len(logs)).Go()).
				Size(SizeSmall).
				Variant(VariantFlat).
				Class("mb-3"),
		).Name("worker_hiddenLogs"))
	}
	for _, l := range logs {
		logLines = append(logLines, P().Style(`
    margin: 0;
    margin-bottom: 4px;`).Children(Text(l)))
	}
	// https://stackoverflow.com/a/44051405/10150757
	var reverseStyle string
	if len(logs) > 18 {
		reverseStyle = "display: flex;flex-direction: column-reverse;"
		for i, j := 0, len(logLines)-1; i < j; i, j = i+1, j-1 {
			logLines[i], logLines[j] = logLines[j], logLines[i]
		}
	}
	inRefresh := status == JobStatusNew || status == JobStatusRunning
	eURL := b.URI()
	return Div(
		Div(Text(msgr.DetailTitleStatus)).Class("text-caption"),
		Div().Class("d-flex align-center mb-5").Children(
			Div().Style("width: 120px").Children(
				Text(fmt.Sprintf("%s (%d%%)", msgr.GetStatus(status), progress)),
			),
			VProgressLinear().ModelValue(int(progress)),
		),

		Div(Text(msgr.DetailTitleLog)).Class("text-caption"),
		Div().Class("mb-3").Style(fmt.Sprintf(`
		background-color: #222;
		color: #fff;
		font-family: menlo,Roboto,Helvetica,Arial,sans-serif;
		height: 300px;
		padding: 8px;
		overflow: auto;
		box-sizing: border-box;
		font-size: 12px;
		line-height: 1;
		%s
		`, reverseStyle)).Children(
			logLines...,
		),

		If(progressText != "",
			Div().Class("mb-3").Children(
				RawHTML(progressText),
			),
		),

		If(canEdit,
			Div().Class("d-flex mt-3").Children(
				VSpacer(),
				If(inRefresh,
					VBtn(msgr.ActionAbortJob).Color("error").
						Attr("@click", web.Plaid().
							URL(eURL).
							EventFunc(EventAbortJob).
							Query("jobID", fmt.Sprintf("%d", id)).
							Query("job", job).
							Go()),
				),
				If(status == JobStatusDone,
					VBtn(msgr.ActionRerunJob).Color("primary").
						Attr("@click", web.Plaid().
							URL(eURL).
							EventFunc("worker_rerunJob").
							Query("jobID", fmt.Sprintf("%d", id)).
							Query("job", job).
							Go()),
				),
			),
		),
	)
}

func (b *Builder) jobSelectList(
	ctx *web.EventContext,
	job string,
) HTMLComponent {
	var vErr web.ValidationErrors
	if ve, ok := ctx.Flash.(*web.ValidationErrors); ok {
		vErr = *ve
	}
	var alert HTMLComponent
	if v := vErr.GetFieldErrors("Job"); len(v) > 0 {
		alert = VAlert(Text(strings.Join(v, ","))).Type("error")
	}
	var (
		currentTitle string
		items        = make([]HTMLComponent, 0, len(b.jbs))
	)

	for _, jb := range b.jbs {
		if !jb.global || jb.system {
			continue
		}

		if editIsAllowed(ctx.R, jb.name) == nil {
			title := jb.GetTitle(ctx)

			if jb.name == job {
				currentTitle = title
			}

			items = append(items,
				VListItem(
					VListItemTitle(
						A(Text(title)).Attr("@click",
							web.Plaid().
								URL(b.URI()).
								EventFunc(EventSelectJob).
								Query("jobName", jb.name).
								Go(),
						),
					)),
			)
		}
	}

	return Div(
		Input("").Type("hidden").Attr(web.VField("Job", job)...),
		If(job == "",
			alert,
			VList(items...).Nav(true).Density(DensityCompact),
		).Else(
			Div(
				VIcon("arrow_back").Attr("@click",
					web.Plaid().EventFunc(EventSelectJob).
						Query("jobName", "").
						Go(),
				),
			).Class("mb-3"),
			Div(Text(currentTitle)).Class("mb-3 text-h6").Style("font-weight: inherit"),
		),
	)
}

func (b *Builder) jobEditingContent(
	ctx *web.EventContext,
	job string,
	args interface{},
) HTMLComponent {
	if job == "" {
		return Template()
	}

	jb := b.mustGetJobBuilder(job)
	var argsObj interface{}
	if args != nil {
		argsObj = args
	} else {
		argsObj = jb.r
	}

	if jb.rmb == nil {
		return Template()
	}
	return jb.rmb.Editing().ToComponent(&presets.ToComponentOptions{}, argsObj, presets.FieldModeStack{presets.NEW}, ctx)
}

// renameLegacy moves the tables and the foreign keys off the old "qor_" prefix,
// which the models carried until they were renamed. It runs before AutoMigrate,
// which would otherwise create the new tables empty and add the new columns
// beside the old ones, leaving the rows behind either way.
func renameLegacy(db *gorm.DB) error {
	m := db.Migrator()

	for _, t := range [][2]string{
		{"qor_jobs", "jobs"},
		{"qor_job_instances", "job_instances"},
		{"qor_job_logs", "job_logs"},
	} {
		if m.HasTable(t[0]) && !m.HasTable(t[1]) {
			if err := m.RenameTable(t[0], t[1]); err != nil {
				return err
			}
		}
	}

	for _, c := range []struct {
		model    any
		old, new string
	}{
		{&JobInstance{}, "qor_job_id", "job_id"},
		{&JobLog{}, "qor_job_instance_id", "job_instance_id"},
	} {
		if m.HasTable(c.model) && m.HasColumn(c.model, c.old) && !m.HasColumn(c.model, c.new) {
			if err := m.RenameColumn(c.model, c.old, c.new); err != nil {
				return err
			}
		}
	}

	return nil
}
