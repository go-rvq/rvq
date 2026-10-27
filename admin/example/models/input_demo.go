package models

import (
	"time"

	"github.com/go-rvq/rvq/admin/media/media_library"
	"github.com/lib/pq"
)

// NotifyChannels are the channels a notification goes by, a bit each: a
// field of bits (presets.BitsComponentFunc).
type NotifyChannels uint8

// The bits of NotifyChannels.
const (
	NotifyEmail NotifyChannels = 1 << iota
	NotifySMS
	NotifyPush
)

type InputDemo struct {
	ID               uint
	TextField1       string
	TextArea1        string
	Switch1          bool
	Slider1          int
	Select1          string
	RangeSlider1     pq.Int64Array `gorm:"type:integer[]"`
	Radio1           string
	FileInput1       string
	Combobox1        string
	Checkbox1        bool
	Autocomplete1    pq.StringArray `gorm:"type:text[]"`
	ButtonGroup1     string
	ChipGroup1       string
	ItemGroup1       string
	ListItemGroup1   string
	SlideGroup1      string
	ColorPicker1     string
	DatePicker1      string
	DatePickerMonth1 string
	TimePicker1      string
	MediaLibrary1    media_library.MediaBox `sql:"type:text;"`
	// Notifications1 is a field of bits: chips, one a channel
	Notifications1 NotifyChannels `gorm:"not null;default:0"`
	UpdatedAt      time.Time
	CreatedAt      time.Time
}
