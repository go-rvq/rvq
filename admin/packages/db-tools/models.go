package db_tools

import (
	"github.com/go-rvq/rvq/thirdpart/gorm/datatypes"
	db_tools "github.com/go-rvq/rvq/x/packages/db-tools"
	"github.com/google/uuid"
)

type DbBackupConfig struct {
	ID          uuid.UUID `gorm:"type:uuid;primaryKey"`
	Persistence datatypes.NullJSONType[*db_tools.Persistence]
}
