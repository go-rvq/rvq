package microsite

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
)

type FileSystem struct {
	FileName string
	Url      string
}

func (fs FileSystem) Value() (driver.Value, error) {
	return json.Marshal(fs)
}

func (fs *FileSystem) Scan(value interface{}) error {
	switch v := value.(type) {
	case string:
		return json.Unmarshal([]byte(v), fs)
	case []byte:
		return json.Unmarshal(v, fs)
	default:
		return errors.New("not supported")
	}
}
