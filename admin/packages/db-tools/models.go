package db_tools

import (
	"github.com/go-rvq/rvq/thirdpart/gorm/datatypes"
	db_tools "github.com/go-rvq/rvq/x/packages/db-tools"
	"github.com/google/uuid"
)

type DbBackupConfig struct {
	ID          uuid.UUID `gorm:"type:uuid;primaryKey"`
	Persistence datatypes.NullJSONType[*db_tools.Persistence]
	// Method is how the backups are made (BackupMethods; "" the default).
	Method string `gorm:"size:32;not null;default:''"`
	// AutoDisabled turns the automatic backups off (the manual ones go on).
	AutoDisabled bool `gorm:"not null;default:false"`
}
