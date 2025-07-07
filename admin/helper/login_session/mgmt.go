package login_session

import (
	"fmt"
	"net/http"
	"sort"
	"time"

	"github.com/dustin/go-humanize"
	h "github.com/go-rvq/htmlgo"
	"github.com/go-rvq/rvq/web"
	"github.com/go-rvq/rvq/x/login"
	. "github.com/go-rvq/rvq/x/ui/vuetify"
	"github.com/ua-parser/uap-go/uaparser"
	"gorm.io/gorm"
)

const (
	LoginTokenHashLen = 8 // The hash string length of the token stored in the DB.
)

type Manager struct {
	lb *login.Builder
}

func NewManager(lb *login.Builder) *Manager {
	ConfigureMessages(lb.I18nBuilder())
	return &Manager{lb: lb}
}

func (m *Manager) AddSessionLogByUserID(db *gorm.DB, r *http.Request, userID uint) (err error) {
	token := login.GetSessionToken(m.lb, r)
	client := uaparser.NewFromSaved().Parse(r.Header.Get("User-Agent"))

	if err = db.Model(&LoginSession{}).Create(&LoginSession{
		UserID:    userID,
		Device:    fmt.Sprintf("%v - %v", client.UserAgent.Family, client.Os.Family),
		IP:        GetIP(r),
		TokenHash: GetStringHash(token, LoginTokenHashLen),
		ExpiredAt: time.Now().Add(time.Duration(m.lb.GetSessionMaxAge()) * time.Second),
	}).Error; err != nil {
		return err
	}

	return nil
}

func (m *Manager) UpdateCurrentSessionLog(db *gorm.DB, r *http.Request, userID uint, oldToken string) (err error) {
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

func (m *Manager) ExpireCurrentSessionLog(db *gorm.DB, r *http.Request, userID uint) (err error) {
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

func (m *Manager) ExpireAllSessionLogs(db *gorm.DB, userID uint) (err error) {
	return db.Model(&LoginSession{}).
		Where("user_id = ?", userID).
		Updates(map[string]interface{}{
			"expired_at": time.Now(),
		}).Error
}

func (m *Manager) ExpireOtherSessionLogs(db *gorm.DB, r *http.Request, userID uint) (err error) {
	token := login.GetSessionToken(m.lb, r)

	return db.Model(&LoginSession{}).
		Where("user_id = ? AND token_hash != ?", userID, GetStringHash(token, LoginTokenHashLen)).
		Updates(map[string]interface{}{
			"expired_at": time.Now(),
		}).Error
}

func (m *Manager) CheckIsTokenValidFromRequest(db *gorm.DB, r *http.Request, userID uint) (valid bool, err error) {
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

func (m *Manager) Sessions(db *gorm.DB, ctx *web.EventContext, userID uint) (comp h.HTMLComponent, err error) {
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
		{Title: msgr.Time, Key: "Time", Width: "25%"},
		{Title: msgr.Device, Key: "Device", Width: "25%"},
		{Title: msgr.IPAddress, Key: "IP", Width: "25%"},
		{Title: msgr.Status, Key: "Status", Width: "25%", Sortable: true},
	}

	comp = h.HTMLComponents{
		h.P(h.Text(msgr.LoginSessionsTips)),
		VDataTable().Headers(sessionTableHeaders).
			Items(items).
			ItemsPerPage(-1).HideDefaultFooter(true),
	}
	return
}
