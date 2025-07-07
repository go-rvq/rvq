package admin

import (
	"fmt"
	"net/http"
	"time"

	"github.com/qor5/admin/v3/example/models"
	"github.com/qor5/x/v3/login"
	"github.com/ua-parser/uap-go/uaparser"
	"gorm.io/gorm"
)

const (
	LoginTokenHashLen = 8 // The hash string length of the token stored in the DB.
)

func addSessionLogByUserID(db *gorm.DB, r *http.Request, userID uint) (err error) {
	token := login.GetSessionToken(loginBuilder, r)
	client := uaparser.NewFromSaved().Parse(r.Header.Get("User-Agent"))

	if err = db.Model(&models.LoginSession{}).Create(&models.LoginSession{
		UserID:    userID,
		Device:    fmt.Sprintf("%v - %v", client.UserAgent.Family, client.Os.Family),
		IP:        ip(r),
		TokenHash: getStringHash(token, LoginTokenHashLen),
		ExpiredAt: time.Now().Add(time.Duration(loginBuilder.GetSessionMaxAge()) * time.Second),
	}).Error; err != nil {
		return err
	}

	return nil
}

func updateCurrentSessionLog(db *gorm.DB, r *http.Request, userID uint, oldToken string) (err error) {
	token := login.GetSessionToken(loginBuilder, r)
	tokenHash := getStringHash(token, LoginTokenHashLen)
	oldTokenHash := getStringHash(oldToken, LoginTokenHashLen)
	if err = db.Model(&models.LoginSession{}).
		Where("user_id = ? and token_hash = ?", userID, oldTokenHash).
		Updates(map[string]interface{}{
			"token_hash": tokenHash,
			"expired_at": time.Now().Add(time.Duration(loginBuilder.GetSessionMaxAge()) * time.Second),
		}).Error; err != nil {
		return err
	}

	return nil
}

func expireCurrentSessionLog(db *gorm.DB, r *http.Request, userID uint) (err error) {
	token := login.GetSessionToken(loginBuilder, r)
	tokenHash := getStringHash(token, LoginTokenHashLen)
	if err = db.Model(&models.LoginSession{}).
		Where("user_id = ? and token_hash = ?", userID, tokenHash).
		Updates(map[string]interface{}{
			"expired_at": time.Now(),
		}).Error; err != nil {
		return err
	}

	return nil
}

func expireAllSessionLogs(db *gorm.DB, userID uint) (err error) {
	return db.Model(&models.LoginSession{}).
		Where("user_id = ?", userID).
		Updates(map[string]interface{}{
			"expired_at": time.Now(),
		}).Error
}

func expireOtherSessionLogs(db *gorm.DB, r *http.Request, userID uint) (err error) {
	token := login.GetSessionToken(loginBuilder, r)

	return db.Model(&models.LoginSession{}).
		Where("user_id = ? AND token_hash != ?", userID, getStringHash(token, LoginTokenHashLen)).
		Updates(map[string]interface{}{
			"expired_at": time.Now(),
		}).Error
}

func isTokenValid(v models.LoginSession) bool {
	return time.Now().Sub(v.ExpiredAt) > 0
}

func checkIsTokenValidFromRequest(db *gorm.DB, r *http.Request, userID uint) (valid bool, err error) {
	token := login.GetSessionToken(loginBuilder, r)
	if token == "" {
		return false, nil
	}
	sessionLog := models.LoginSession{}
	if err = db.Where("user_id = ? and token_hash = ?", userID, getStringHash(token, LoginTokenHashLen)).
		First(&sessionLog).
		Error; err != nil {
		if err != gorm.ErrRecordNotFound {
			return false, err
		}
		return false, nil
	}
	// IP check
	if sessionLog.IP != ip(r) {
		return false, nil
	}
	if isTokenValid(sessionLog) {
		return false, nil
	}
	return true, nil
}
