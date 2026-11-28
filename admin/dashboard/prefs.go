package dashboard

import (
	"encoding/json"
	"errors"
	"github.com/go-rvq/rvq/x/migrate"
	"slices"
	"strings"
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
	// Sizes are the widths the user gave the widgets (resized), by id
	Sizes map[string]Size `json:"sizes,omitempty"`
	// Added: of before (the opt-in widgets added), read and not used
	Added  []string `json:"added,omitempty"`
	Period string   `json:"period,omitempty"`
}

func (p *Prefs) collapsed(id string) bool { return slices.Contains(p.Collapsed, id) }

// MinSize is the narrowest a widget is resized to.
const MinSize Size = 2

// sizeOf is the width of w for the user: theirs (resized), else its own.
func (p *Prefs) sizeOf(w *Widget) Size {
	if n, ok := p.Sizes[w.ID]; ok && n >= MinSize && n <= Full {
		return n
	}
	return w.Size
}

// resize gives the widget id the width n (its own: n is forgotten).
func (p *Prefs) resize(w *Widget, n Size) {
	if n == w.Size {
		delete(p.Sizes, w.ID)
		return
	}
	if p.Sizes == nil {
		p.Sizes = map[string]Size{}
	}
	p.Sizes[w.ID] = n
}

func toggle(list []string, id string, on bool) []string {
	list = slices.DeleteFunc(list, func(s string) bool { return s == id })
	if on {
		list = append(list, id)
	}
	return list
}

// The choices of a user are kept a row (Pref), its Data the ones of each
// panel: {"panels": {"<panel>": Prefs}}. A row of before — one Prefs, of the
// dashboard alone — is the home's, its widgets' ids of before (a name alone)
// made "<group>/<name>" (normalize).

// storedPrefs are the choices of a user, by panel.
type storedPrefs struct {
	Panels map[string]Prefs `json:"panels"`
}

// loadAll are the choices of user, by panel (none: empty).
func (b *Builder) loadAll(user uuid.UUID) map[string]Prefs {
	all := map[string]Prefs{}
	if b.db == nil || user == uuid.Nil {
		return all
	}
	var row Pref
	if err := b.db.Where("user_id = ?", user).First(&row).Error; err != nil {
		return all
	}
	var raw map[string]json.RawMessage
	if json.Unmarshal([]byte(row.Data), &raw) != nil {
		return all
	}
	if _, ok := raw["panels"]; ok {
		var st storedPrefs
		if json.Unmarshal([]byte(row.Data), &st) == nil && st.Panels != nil {
			return st.Panels
		}
		return all
	}
	// of before: the home's
	var p Prefs
	if json.Unmarshal([]byte(row.Data), &p) == nil {
		all[PanelHome] = p
	}
	return all
}

// saveAll keeps the choices of user, by panel.
func (b *Builder) saveAll(user uuid.UUID, all map[string]Prefs) error {
	if b.db == nil || user == uuid.Nil {
		return errors.New("dashboard: no user")
	}
	data, err := json.Marshal(storedPrefs{Panels: all})
	if err != nil {
		return err
	}
	row := Pref{UserID: user, Data: string(data), UpdatedAt: time.Now()}
	return b.db.Save(&row).Error
}

// LoadPrefs are the choices of user of the panel (none: empty — everything
// shown).
func (b *Panel) LoadPrefs(user uuid.UUID) (p Prefs) {
	p = b.loadAll(user)[b.id]
	b.normalize(&p)
	return
}

// normalize makes the ids of before (a widget's name alone) of p the ids of
// the panel's widgets of that name ("<group>/<name>"), when one only.
func (b *Panel) normalize(p *Prefs) {
	byName := map[string][]string{}
	for _, w := range b.all() {
		byName[w.Name] = append(byName[w.Name], w.ID)
	}
	fix := func(ids []string) []string {
		for i, id := range ids {
			if !strings.Contains(id, "/") && len(byName[id]) == 1 {
				ids[i] = byName[id][0]
			}
		}
		return ids
	}
	for i := range p.Sections {
		p.Sections[i].Widgets = fix(p.Sections[i].Widgets)
	}
	p.Removed, p.Disabled, p.Hidden, p.Order = fix(p.Removed), fix(p.Disabled), fix(p.Hidden), fix(p.Order)
}

// SavePrefs keeps the choices of user of the panel.
func (b *Panel) SavePrefs(user uuid.UUID, p Prefs) error {
	all := b.loadAll(user)
	all[b.id] = p
	return b.saveAll(user, all)
}

// ResetPrefs drops the choices of user of the panel.
func (b *Panel) ResetPrefs(user uuid.UUID) error {
	if b.db == nil || user == uuid.Nil {
		return nil
	}
	all := b.loadAll(user)
	delete(all, b.id)
	if len(all) == 0 {
		return b.db.Where("user_id = ?", user).Delete(&Pref{}).Error
	}
	return b.saveAll(user, all)
}

// Migrate creates the table of the preferences.
func (b *Builder) Migrate() error {
	if b.db == nil {
		return nil
	}
	return migrate.Auto(b.db, Migrate)
}
