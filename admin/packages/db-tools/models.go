package db_tools

import (
	"github.com/go-rvq/rvq/thirdpart/gorm/datatypes"
	db_tools "github.com/go-rvq/rvq/x/packages/db-tools"
)

type DbBackupConfig struct {
	ID          uint `gorm:"primaryKey"`
	Persistence datatypes.NullJSONType[*db_tools.Persistence]
}
