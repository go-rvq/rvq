package examples_presets

import (
	"context"

	"github.com/go-rvq/rvq/admin/presets"
	"github.com/go-rvq/rvq/admin/presets/gorm2op"
	"github.com/go-rvq/rvq/x/i18n"
	"golang.org/x/text/language"
	"gorm.io/gorm"
)

// @snippet_begin(PresetsBitsFieldSample)

// Channels are the channels a subscriber is notified by: a bit each.
type Channels uint8

const (
	ChannelEmail Channels = 1 << iota
	ChannelSMS
	ChannelPush
)

type Subscriber struct {
	ID       uint
	Name     string
	Channels Channels `gorm:"not null;default:0"`
}

// The label and the description of each option are messages of the module:
// translated, never written in the code.
const BitsMessagesKey i18n.ModuleKey = "examples-bits"

type BitsMessages struct {
	ChannelEmail      string
	ChannelEmail_Desc string
	ChannelSMS        string
	ChannelSMS_Desc   string
	ChannelPush       string
	ChannelPush_Desc  string
}

var BitsMessages_en_US = &BitsMessages{
	ChannelEmail:      "Email",
	ChannelEmail_Desc: "Notified by email, at the address of the account.",
	ChannelSMS:        "SMS",
	ChannelSMS_Desc:   "Notified by a text message, at the phone of the account.",
	ChannelPush:       "Push",
	ChannelPush_Desc:  "Notified on the devices where the account is signed in.",
}

func bitsMessages(ctx context.Context) *BitsMessages {
	return i18n.MustGetModuleMessages(ctx, BitsMessagesKey, BitsMessages_en_US).(*BitsMessages)
}

func PresetsBitsField(b *presets.Builder, db *gorm.DB) (
	mb *presets.ModelBuilder,
	cl *presets.ListingBuilder,
	ce *presets.EditingBuilder,
	dp *presets.DetailingBuilder,
) {
	if err := db.AutoMigrate(&Subscriber{}); err != nil {
		panic(err)
	}
	b.DataOperator(gorm2op.DataOperator(db))
	b.I18n().RegisterForModule(language.English, BitsMessagesKey, BitsMessages_en_US)

	// one option a bit: its value, its label and its description
	option := func(bit Channels, label, desc func(*BitsMessages) string) presets.BitOption {
		return presets.BitOption{
			Bit:         uint64(bit),
			Label:       func(ctx context.Context) string { return label(bitsMessages(ctx)) },
			Description: func(ctx context.Context) string { return desc(bitsMessages(ctx)) },
		}
	}
	channels := presets.BitsComponentFunc(
		option(ChannelEmail, func(m *BitsMessages) string { return m.ChannelEmail }, func(m *BitsMessages) string { return m.ChannelEmail_Desc }),
		option(ChannelSMS, func(m *BitsMessages) string { return m.ChannelSMS }, func(m *BitsMessages) string { return m.ChannelSMS_Desc }),
		option(ChannelPush, func(m *BitsMessages) string { return m.ChannelPush }, func(m *BitsMessages) string { return m.ChannelPush_Desc }),
	)

	mb = b.Model(&Subscriber{})
	cl = mb.Listing("ID", "Name")
	ce = mb.Editing("Name", "Channels")
	dp = mb.Detailing("Name", "Channels")
	// the same component edits (the chips toggle; the value is the sum of the
	// bits marked) and shows (the set ones marked)
	ce.Field("Channels").ComponentFunc(channels)
	dp.Field("Channels").ComponentFunc(channels)
	return
}

// @snippet_end
