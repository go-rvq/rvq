package dashboard

import (
	"encoding/json"
	"errors"
	"slices"
	"time"

	"github.com/google/uuid"
)

// Pref is what a user chose of the dashboard: its sections and the widgets
// of each (layout.go), the widgets removed and disabled, the groups
// collapsed, the period it opens with.
type Pref struct {
	UserID    uuid.UUID `gorm:"primaryKey;type:uuid"`
	Data      string    `gorm:"type:text"`
	UpdatedAt time.Time
}

// TableName is the table of the preferences.
func (Pref) TableName() string { return "dashboard_prefs" }

// Prefs are a user's choices (Pref.Data). Hidden and Order are of before
// the sections: a widget hidden is removed, the order the first layout's.
type Prefs struct {
	Sections  []Section `json:"sections,omitempty"`
	Removed   []string  `json:"removed,omitempty"`
	Disabled  []string  `json:"disabled,omitempty"`
	Hidden    []string  `json:"hidden,omitempty"`
	Order     []string  `json:"order,omitempty"`
	Collapsed []string  `json:"collapsed,omitempty"`
	Period    string    `json:"period,omitempty"`
}

func (p *Prefs) collapsed(id string) bool { return slices.Contains(p.Collapsed, id) }

func toggle(list []string, id string, on bool) []string {
	list = slices.DeleteFunc(list, func(s string) bool { return s == id })
	if on {
		list = append(list, id)
	}
	return list
}

// LoadPrefs are the choices of user (none: empty).
func (b *Builder) LoadPrefs(user uuid.UUID) (p Prefs) {
	if b.db == nil || user == uuid.Nil {
		return
	}
	var row Pref
	if err := b.db.Where("user_id = ?", user).First(&row).Error; err != nil {
		return
	}
	_ = json.Unmarshal([]byte(row.Data), &p)
	return
}

// SavePrefs keeps the choices of user.
func (b *Builder) SavePrefs(user uuid.UUID, p Prefs) error {
	if b.db == nil || user == uuid.Nil {
		return errors.New("dashboard: no user")
	}
	data, err := json.Marshal(p)
	if err != nil {
		return err
	}
	row := Pref{UserID: user, Data: string(data), UpdatedAt: time.Now()}
	return b.db.Save(&row).Error
}

// ResetPrefs drops the choices of user.
func (b *Builder) ResetPrefs(user uuid.UUID) error {
	if b.db == nil {
		return nil
	}
	return b.db.Where("user_id = ?", user).Delete(&Pref{}).Error
}

// Migrate creates the table of the preferences.
func (b *Builder) Migrate() error {
	if b.db == nil {
		return nil
	}
	return b.db.AutoMigrate(&Pref{})
}
