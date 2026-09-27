package seo

import (
	"database/sql/driver"
	"encoding/json"

	"github.com/go-rvq/rvq/admin/utils/uuidkey"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

// SEOConfigGoogleMapsKey is the Data key holding the Google Maps browser API key.
const SEOConfigGoogleMapsKey = "google_maps_api_key"

// SEOConfigData is the SEO configuration's free-form JSON payload. Well-known
// keys have typed accessors (e.g. the Google Maps API key); anything else is
// carried through untouched.
type SEOConfigData map[string]any

// Scan reads the JSON payload from the database.
func (d *SEOConfigData) Scan(value interface{}) error {
	if *d == nil {
		*d = SEOConfigData{}
	}
	switch v := value.(type) {
	case nil:
		return nil
	case []byte:
		if len(v) == 0 {
			return nil
		}
		return json.Unmarshal(v, d)
	case string:
		if v == "" {
			return nil
		}
		return json.Unmarshal([]byte(v), d)
	}
	return nil
}

// Value writes the JSON payload to the database.
func (d SEOConfigData) Value() (driver.Value, error) {
	if d == nil {
		return "{}", nil
	}
	b, err := json.Marshal(d)
	return string(b), err
}

// GetString returns a string value from the payload (empty when absent).
func (d SEOConfigData) GetString(key string) string {
	if d == nil {
		return ""
	}
	s, _ := d[key].(string)
	return s
}

// SetString sets a string value on the payload.
func (d *SEOConfigData) SetString(key, val string) {
	if *d == nil {
		*d = SEOConfigData{}
	}
	(*d)[key] = val
}

// SEOConfig is the SEO settings singleton (one row): a free-form JSON payload
// (Data) with typed accessors for known values such as the Google Maps API key
// (used by the built-in zipcodes variable's map picker). Table seo_config.
type SEOConfig struct {
	ID   uuid.UUID     `gorm:"type:uuid;primaryKey"`
	Data SEOConfigData `sql:"type:text"`
}

// TableName is seo_config.
func (SEOConfig) TableName() string { return "seo_config" }

// GoogleMapsAPIKey is the configured Maps API key.
func (c *SEOConfig) GoogleMapsAPIKey() string { return c.Data.GetString(SEOConfigGoogleMapsKey) }

// SetGoogleMapsAPIKey sets the Maps API key.
func (c *SEOConfig) SetGoogleMapsAPIKey(key string) { c.Data.SetString(SEOConfigGoogleMapsKey, key) }

// SEOConfigID is the key of the singleton SEO config: fixed, and derived from
// its name (a UUID v5), so every database — and the migration from the integer
// key 1 — agrees on it.
var SEOConfigID = uuid.NewSHA1(uuid.NameSpaceURL, []byte("rvq:seo_config"))

// LoadSEOConfig returns the singleton SEO config (SEOConfigID), creating an
// empty row on first use. Best-effort: a DB error yields a zero config so
// callers never fail over configuration.
func LoadSEOConfig(db *gorm.DB) *SEOConfig {
	uuidkey.MustRegister(db) // its records have UUID keys
	c := &SEOConfig{ID: SEOConfigID}
	if db == nil {
		return c
	}
	db.Session(&gorm.Session{}).FirstOrCreate(c, "id = ?", SEOConfigID)
	return c
}
