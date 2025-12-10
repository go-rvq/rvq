package db_tools

import (
	"fmt"
	"strings"
	"time"
)

const BackupIDTimeFormat = "20060102T150405"

type BackupID struct {
	Auto      bool
	DbName    string
	CreatedAt time.Time
}

func (id BackupID) String() string {
	var a string
	if id.Auto {
		a = "A"
	} else {
		a = "M"
	}

	return a + ":" + id.DbName + ":" + id.CreatedAt.Format(BackupIDTimeFormat)
}

func (id *BackupID) Parse(v string) (err error) {
	if len(v) == 0 {
		return fmt.Errorf("empty id")
	}
	parts := strings.Split(v, ":")
	if len(parts) != 3 {
		return fmt.Errorf("malformed backup id: %s", v)
	}
	switch parts[0] {
	case "A":
		id.Auto = true
	case "M":
		id.Auto = false
	}
	id.DbName = parts[1]
	id.CreatedAt, err = time.Parse(BackupIDTimeFormat, parts[2])
	return
}

func (id BackupID) IsValid() bool {
	return id.DbName != "" && id.CreatedAt.IsZero()
}
