package publish

import (
	"fmt"
	"reflect"
	"time"

	"github.com/qor5/admin/v3/model"
	"github.com/qor5/admin/v3/reflect_utils"
	"github.com/qor5/admin/v3/utils/db_utils"
	"gorm.io/gorm"
)

func (version *Version) GetNextVersion(t *time.Time) string {
	if t == nil {
		return ""
	}
	date := t.Format("2006-01-02")
	return fmt.Sprintf("%s-v%02v", date, 1)
}

func (version *Version) CreateVersion(db *gorm.DB, mid model.ID, obj interface{}) (string, error) {
	date := db.NowFunc().Format("2006-01-02")
	var count int64
	if err := db_utils.ModelIdWhere(db.Unscoped(), obj, mid, "Version").
		Where("version like ?", date+"%").
		Order("version DESC").
		Count(&count).Error; err != nil {
		return "", err
	}

	versionName := fmt.Sprintf("%s-v%02v", date, count+1)
	version.Version = versionName
	version.VersionName = versionName
	return version.Version, nil
}

func IsVersion(obj interface{}) (IsVersion bool) {
	_, IsVersion = reflect_utils.GetStruct(reflect.TypeOf(obj)).(VersionInterface)
	return
}
