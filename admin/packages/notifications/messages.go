package notifications

import (
	"context"

	"github.com/go-rvq/rvq/x/i18n"
	"golang.org/x/text/language"
)

// MessagesKey is the module of the words of the notifications.
const MessagesKey i18n.ModuleKey = "admin/packages/notifications"

// Messages are the words of the notification center.
type Messages struct {
	ModuleDescription     string `i18n:"hint='Says what the module is, for whoever translates its texts.'"`
	Notification          string `i18n:"hint='Name of a notification of the user (singular).'"`
	Notifications         string `i18n:"label='My notifications', hint='The notifications of the user: the listing (plural).'"`
	NotificationTitle     string `i18n:"label='Title', hint='Label of the title of a notification.'"`
	NotificationBody      string `i18n:"label='Text', hint='Label of the text of a notification.'"`
	NotificationLevel     string `i18n:"label='Level', hint='Label of the level of a notification (info, warning, error).'"`
	NotificationSource    string `i18n:"label='Source', hint='Label of what made a notification.'"`
	NotificationLink      string `i18n:"label='Link', hint='Label of where a notification takes.'"`
	NotificationCreatedAt string `i18n:"label='When', hint='Label of when a notification was made.'"`
	NotificationReadAt    string `i18n:"label='Read at', hint='Label of when the user read a notification.'"`
	MarkAllRead           string `i18n:"hint='Button of the bell: marks every notification of the user read.'"`
	SeeAll                string `i18n:"hint='Link of the bell to the listing of the notifications of the user.'"`
	Open                  string `i18n:"hint='Action of a notification that opens where it takes (marking it read).'"`
	Unread                string `i18n:"hint='Tab of the notifications: the unread ones.'"`
	All                   string `i18n:"hint='Tab of the notifications: all of them.'"`
	BellTitle             string `i18n:"hint='Title of the notifications in the bell: %d is how many are unread.'"`
}

var Messages_en_US = &Messages{
	ModuleDescription:     "The notification center of the users.",
	Notification:          "Notification",
	Notifications:         "My notifications",
	NotificationTitle:     "Title",
	NotificationBody:      "Text",
	NotificationLevel:     "Level",
	NotificationSource:    "Source",
	NotificationLink:      "Link",
	NotificationCreatedAt: "When",
	NotificationReadAt:    "Read at",
	MarkAllRead:           "Mark all as read",
	SeeAll:                "See all",
	Open:                  "Open",
	Unread:                "Unread",
	All:                   "All",
	BellTitle:             "Notifications (%d unread)",
}

var Messages_pt_BR = &Messages{
	ModuleDescription:     "A central de notificações dos usuários.",
	Notification:          "Notificação",
	Notifications:         "Minhas notificações",
	NotificationTitle:     "Título",
	NotificationBody:      "Texto",
	NotificationLevel:     "Nível",
	NotificationSource:    "Origem",
	NotificationLink:      "Link",
	NotificationCreatedAt: "Quando",
	NotificationReadAt:    "Lida em",
	MarkAllRead:           "Marcar todas como lidas",
	SeeAll:                "Ver todas",
	Open:                  "Abrir",
	Unread:                "Não lidas",
	All:                   "Todas",
	BellTitle:             "Notificações (%d não lidas)",
}

// ConfigureMessages registers the words.
func ConfigureMessages(b *i18n.Builder) {
	b.RegisterForModules(language.English, MessagesKey, Messages_en_US).
		RegisterForModules(language.BrazilianPortuguese, MessagesKey, Messages_pt_BR)
}

// GetMessages are the words in the language of ctx.
func GetMessages(ctx context.Context) *Messages {
	return i18n.MustGetModuleMessages(ctx, MessagesKey, Messages_en_US).(*Messages)
}
