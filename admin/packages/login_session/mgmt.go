package login_session

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"sort"
	"strings"
	"sync"
	"time"

	"github.com/dustin/go-humanize"
	h "github.com/go-rvq/htmlgo"
	"github.com/go-rvq/rvq/admin/utils/uuidkey"
	"github.com/go-rvq/rvq/web"
	"github.com/go-rvq/rvq/x/login"
	"github.com/go-rvq/rvq/x/place"
	. "github.com/go-rvq/rvq/x/ui/vuetify"
	"github.com/google/uuid"
	"github.com/ua-parser/uap-go/uaparser"
	"gorm.io/gorm"
)

const (
	LoginTokenHashLen = 8 // The hash string length of the token stored in the DB.
)

// AuthLocked is the Auth of an access that locked the account (a wrong
// password of HTTP Basic, the last one allowed).
const AuthLocked = "locked"

// DefaultAccessWindow is how long the accesses of one device, IP and way in
// are one row (Manager.AccessWindow).
const DefaultAccessWindow = time.Hour

// accessWriteEvery is how often an access's row is written: the requests in
// between are counted in memory.
const accessWriteEvery = time.Minute

type Manager struct {
	lb *login.Builder

	// PlaceFunc is where the IP of a request is (zero: not known): the
	// application's (a geolocation).
	PlaceFunc func(r *http.Request) place.Place
	// AccessWindow is how long the accesses of one device, IP and way in are
	// one row (DefaultAccessWindow when 0).
	AccessWindow time.Duration

	mu       sync.Mutex
	accesses map[string]*pendingAccess
}

// pendingAccess is an access's row and the requests not written to it yet.
type pendingAccess struct {
	id       uuid.UUID
	written  time.Time
	until    time.Time
	requests int
}

func NewManager(lb *login.Builder) *Manager {
	ConfigureMessages(lb.I18nBuilder())
	return &Manager{lb: lb, accesses: map[string]*pendingAccess{}}
}

func (m *Manager) place(r *http.Request) place.Place {
	if m.PlaceFunc == nil {
		return place.Place{}
	}
	return m.PlaceFunc(r)
}

func (m *Manager) window() time.Duration {
	if m.AccessWindow > 0 {
		return m.AccessWindow
	}
	return DefaultAccessWindow
}

func device(r *http.Request) string {
	client := uaparser.NewFromSaved().Parse(r.Header.Get("User-Agent"))
	return fmt.Sprintf("%v - %v", client.UserAgent.Family, client.Os.Family)
}

// Access is an access RecordAccess records: whose, its kind (KindWebDAV,
// KindGit, or the application's: "api"), how it was told who
// (login.AuthOf; AuthLocked when it locked the account), and what it was by
// when it says more than the browser ("Access key: deploy"; "" the
// browser's).
type Access struct {
	UserID uuid.UUID
	Kind   string
	Auth   string
	Device string
}

// RecordAccess records a request of a: the accesses of one device, IP and
// way in, in a window (AccessWindow), are one row — its requests counted,
// its last access written at most once a minute.
func (m *Manager) RecordAccess(db *gorm.DB, r *http.Request, a Access) error {
	uuidkey.MustRegister(db)
	ip, dev, now := GetIP(r), a.Device, time.Now()
	if dev == "" {
		dev = device(r)
	}
	userID, kind, auth := a.UserID, a.Kind, a.Auth
	key := strings.Join([]string{userID.String(), kind, auth, ip, dev}, "\x00")

	m.mu.Lock()
	p := m.accesses[key]
	if p != nil && now.Before(p.until) {
		p.requests++
		if now.Sub(p.written) < accessWriteEvery {
			m.mu.Unlock()
			return nil
		}
		n := p.requests
		p.requests, p.written = 0, now
		m.mu.Unlock()
		return db.Model(&LoginSession{}).Where("id = ?", p.id).Updates(map[string]any{
			"requests":       gorm.Expr("requests + ?", n),
			"last_access_at": now,
			"expired_at":     now.Add(m.window()),
		}).Error
	}
	m.mu.Unlock()

	// a new window — or one this process did not see (it restarted)
	var row LoginSession
	err := db.Where("user_id = ? AND kind = ? AND auth = ? AND ip = ? AND device = ? AND last_access_at > ?",
		userID, kind, auth, ip, dev, now.Add(-m.window())).Order("last_access_at DESC").First(&row).Error
	switch {
	case err == nil:
		err = db.Model(&row).Updates(map[string]any{
			"requests":       gorm.Expr("requests + 1"),
			"last_access_at": now,
			"expired_at":     now.Add(m.window()),
		}).Error
	case errors.Is(err, gorm.ErrRecordNotFound):
		row = LoginSession{
			UserID: userID, Device: dev, IP: ip, UserAgent: r.UserAgent(), Kind: kind, Auth: auth, Place: m.place(r),
			LastAccessAt: &now, Requests: 1, ExpiredAt: now.Add(m.window()),
		}
		err = db.Create(&row).Error
	}
	if err != nil {
		return err
	}
	m.mu.Lock()
	m.accesses[key] = &pendingAccess{id: row.ID, written: now, until: now.Add(m.window())}
	for k, a := range m.accesses {
		if now.After(a.until) {
			delete(m.accesses, k)
		}
	}
	m.mu.Unlock()
	return nil
}

// Purge deletes the rows older than retention — of a login ended, of an
// access whose last request is that old —: how many.
func (m *Manager) Purge(db *gorm.DB, retention time.Duration) (int64, error) {
	if retention <= 0 {
		return 0, nil
	}
	before := time.Now().Add(-retention)
	res := db.Where("expired_at < ? AND COALESCE(last_access_at, created_at) < ?", before, before).Delete(&LoginSession{})
	return res.RowsAffected, res.Error
}

func (m *Manager) AddSessionLogByUserID(db *gorm.DB, r *http.Request, userID uuid.UUID) (err error) {
	uuidkey.MustRegister(db) // its records have UUID keys
	token := login.GetSessionToken(m.lb, r)

	if err = db.Model(&LoginSession{}).Create(&LoginSession{
		UserID:    userID,
		Device:    device(r),
		IP:        GetIP(r),
		TokenHash: GetStringHash(token, LoginTokenHashLen),
		ExpiredAt: time.Now().Add(time.Duration(m.lb.GetSessionMaxAge()) * time.Second),
		Kind:      KindLogin,
		Auth:      login.AuthPassword,
		UserAgent: r.UserAgent(),
		Place:     m.place(r),
	}).Error; err != nil {
		return err
	}

	return nil
}

func (m *Manager) UpdateCurrentSessionLog(db *gorm.DB, r *http.Request, userID uuid.UUID, oldToken string) (err error) {
	token := login.GetSessionToken(m.lb, r)
	tokenHash := GetStringHash(token, LoginTokenHashLen)
	oldTokenHash := GetStringHash(oldToken, LoginTokenHashLen)
	if err = db.Model(&LoginSession{}).
		Where("user_id = ? and token_hash = ?", userID, oldTokenHash).
		Updates(map[string]interface{}{
			"token_hash": tokenHash,
			"expired_at": time.Now().Add(time.Duration(m.lb.GetSessionMaxAge()) * time.Second),
		}).Error; err != nil {
		return err
	}

	return nil
}

func (m *Manager) ExpireCurrentSessionLog(db *gorm.DB, r *http.Request, userID uuid.UUID) (err error) {
	token := login.GetSessionToken(m.lb, r)
	tokenHash := GetStringHash(token, LoginTokenHashLen)
	if err = db.Model(&LoginSession{}).
		Where("user_id = ? and token_hash = ?", userID, tokenHash).
		Updates(map[string]interface{}{
			"expired_at": time.Now(),
		}).Error; err != nil {
		return err
	}

	return nil
}

func (m *Manager) ExpireAllSessionLogs(db *gorm.DB, userID uuid.UUID) (err error) {
	return db.Model(&LoginSession{}).
		Where("user_id = ?", userID).
		Updates(map[string]interface{}{
			"expired_at": time.Now(),
		}).Error
}

func (m *Manager) ExpireOtherSessionLogs(db *gorm.DB, r *http.Request, userID uuid.UUID) (err error) {
	token := login.GetSessionToken(m.lb, r)

	return db.Model(&LoginSession{}).
		Where("user_id = ? AND token_hash != ?", userID, GetStringHash(token, LoginTokenHashLen)).
		Updates(map[string]interface{}{
			"expired_at": time.Now(),
		}).Error
}

func (m *Manager) CheckIsTokenValidFromRequest(db *gorm.DB, r *http.Request, userID uuid.UUID) (valid bool, err error) {
	token := login.GetSessionToken(m.lb, r)
	if token == "" {
		return false, nil
	}
	sessionLog := LoginSession{}
	if err = db.Where("user_id = ? and token_hash = ?", userID, GetStringHash(token, LoginTokenHashLen)).
		First(&sessionLog).
		Error; err != nil {
		if err != gorm.ErrRecordNotFound {
			return false, err
		}
		return false, nil
	}
	// IP check
	if sessionLog.IP != GetIP(r) {
		return false, nil
	}
	if IsTokenValid(sessionLog) {
		return false, nil
	}
	return true, nil
}

func (m *Manager) Sessions(db *gorm.DB, ctx *web.EventContext, userID uuid.UUID) (comp h.HTMLComponent, err error) {
	msgr := GetMessages(ctx.Context())
	var items []*LoginSession
	if err = db.Where("user_id = ?", userID).Find(&items).Error; err != nil {
		return
	}

	currentTokenHash := GetStringHash(login.GetSessionToken(m.lb, ctx.R), LoginTokenHashLen)

	var (
		expired        = msgr.Expired
		active         = msgr.Active
		currentSession = msgr.CurrentSession
		currentItem    *LoginSession
	)

	activeDevices := make(map[string]struct{})
	for _, item := range items {
		if IsTokenValid(*item) {
			item.Status = expired
		} else {
			item.Status = active
			activeDevices[fmt.Sprintf("%s#%s", item.Device, item.IP)] = struct{}{}
		}
		if item.TokenHash == currentTokenHash {
			item.Status = currentSession
			currentItem = item
		}

		item.Time = humanize.Time(item.CreatedAt)
		if item.LastAccessAt != nil {
			item.Time = humanize.Time(*item.LastAccessAt)
		}
		item.Access = accessLabel(msgr, item)
		item.PlaceText = item.Place.String()
	}

	{
		newItems := make([]*LoginSession, 0, len(items))

		for _, item := range items {
			if item == currentItem {
				continue
			}

			if item.Status == expired {
				_, ok := activeDevices[fmt.Sprintf("%s#%s", item.Device, item.IP)]
				if ok {
					continue
				}
			}
			newItems = append(newItems, item)
		}
		items = newItems
	}

	sort.Slice(items, func(i, j int) bool {
		if items[i].Status == expired &&
			items[j].Status == active {
			return false
		}
		if items[i].CreatedAt.Sub(items[j].CreatedAt) < 0 {
			return false
		}
		return true
	})

	if currentItem != nil {
		items = append([]*LoginSession{currentItem}, items...)
	}

	sessionTableHeaders := DataTableHeaderBasicSlice{
		{Title: msgr.Time, Key: "Time", Width: "15%"},
		{Title: msgr.Access, Key: "Access", Width: "20%"},
		{Title: msgr.Device, Key: "Device", Width: "20%"},
		{Title: msgr.IPAddress, Key: "IP", Width: "15%"},
		{Title: msgr.Place, Key: "PlaceText", Width: "15%"},
		{Title: msgr.Status, Key: "Status", Width: "15%", Sortable: true},
	}

	comp = h.HTMLComponents{
		h.P(h.Text(msgr.LoginSessionsTips)),
		VDataTable().Headers(sessionTableHeaders).
			Items(items).
			ItemsPerPage(-1).HideDefaultFooter(true),
	}
	return
}

// AccessLabel is the way of the session s, in the language of ctx: its kind,
// how it was told who, the requests of an access ("Git · access key · 12
// requests").
func AccessLabel(ctx context.Context, s *LoginSession) string {
	return accessLabel(GetMessages(ctx), s)
}

func accessLabel(msgr *Messages, s *LoginSession) string {
	kind := map[string]string{KindWebDAV: msgr.KindWebDAV, KindGit: msgr.KindGit}[s.Kind]
	if kind == "" {
		kind = msgr.KindLogin
	}
	auth := map[string]string{
		login.AuthPassword: msgr.AuthPassword, login.AuthAccessKey: msgr.AuthAccessKey,
		login.AuthSecureKey: msgr.AuthSecureKey, login.AuthSession: msgr.AuthSession, AuthLocked: msgr.AuthLocked,
	}[s.Auth]
	out := kind
	if auth != "" {
		out += " · " + auth
	}
	if s.Kind != KindLogin && s.Kind != "" && s.Requests > 0 {
		out += " · " + fmt.Sprintf(msgr.Requests, s.Requests)
	}
	return out
}
