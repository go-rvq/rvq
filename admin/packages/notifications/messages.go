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
	ModuleDescription        string `i18n:"hint='Says what the module is, for whoever translates its texts.'"`
	ServiceTitle             string `i18n:"hint='Title of the service of the notifications (its stable paths).'"`
	ServiceDesc              string `i18n:"hint='What the service of the notifications serves.'"`
	ServiceFileRoute         string `i18n:"hint='What the route of an attachment of a notification serves.'"`
	Notification             string `i18n:"hint='Name of a notification of the user (singular).'"`
	Notifications            string `i18n:"label='My notifications', hint='The notifications of the user: the listing (plural).'"`
	NotificationTitle        string `i18n:"label='Title', hint='Label of the title of a notification.'"`
	NotificationBody         string `i18n:"label='Text', hint='Label of the text of a notification.'"`
	NotificationLevel        string `i18n:"label='Importance', hint='Label of the level of importance of a notification (information, success, alert, critical).'"`
	NotificationSource       string `i18n:"label='Source', hint='Label of what made a notification.'"`
	NotificationLink         string `i18n:"label='Link', hint='Label of where a notification takes.'"`
	NotificationCreatedAt    string `i18n:"label='When', hint='Label of when a notification was made.'"`
	NotificationSender       string `i18n:"label='Sent by', hint='Label of who sent a notification: the login of the user, or the system.'"`
	NotificationConversation string `i18n:"label='Conversation', hint='Label of the messages of the conversation of a notification of a user, and its answer.'"`
	NotificationAttachments  string `i18n:"label='Attachments', hint='Label of the files sent with a notification.'"`
	NotificationReadAt       string `i18n:"label='Read at', hint='Label of when the user read a notification.'"`
	MarkAllRead              string `i18n:"hint='Button of the bell: marks every notification of the user read.'"`
	SeeAll                   string `i18n:"hint='Link of the bell to the listing of the notifications of the user.'"`
	Open                     string `i18n:"hint='Action of a notification that opens where it takes (marking it read).'"`
	Unread                   string `i18n:"hint='Tab of the notifications: the unread ones.'"`
	Read                     string `i18n:"hint='Tab of the notifications: the read ones.'"`
	All                      string `i18n:"hint='Tab of the notifications: all of them.'"`
	BellTitle                string `i18n:"hint='Title of the notifications in the bell: %d is how many are unread.'"`
	NotificationSystem       string `i18n:"label='System', hint='Label of whether a notification is the system\\'s (translated) and not a user\\'s.'"`
	Period                   string `i18n:"hint='Filter of the notifications: a period up to now.'"`
	PeriodToday              string `i18n:"hint='Option of the period: today.'"`
	Period7d                 string `i18n:"hint='Option of the period: the last 7 days.'"`
	Period30d                string `i18n:"hint='Option of the period: the last 30 days.'"`
	Period90d                string `i18n:"hint='Option of the period: the last 90 days.'"`
	Period365d               string `i18n:"hint='Option of the period: the last 365 days.'"`
	Origin                   string `i18n:"hint='Filter of the notifications: who made them.'"`
	OriginSystem             string `i18n:"hint='Option of the origin: the system.'"`
	OriginUsers              string `i18n:"hint='Option of the origin: the users (sent by one).'"`
	OriginOthers             string `i18n:"hint='Option of the origin: neither the system nor a user.'"`
	Send                     string `i18n:"hint='Action of the listing of the notifications: sends one to other users.'"`
	SendTo                   string `i18n:"hint='Field of the sending: the users it goes to.'"`
	SendRequired             string `i18n:"hint='Error of the sending with a field empty (every one is required).'"`
	SendDenied               string `i18n:"hint='Error of the sending by who may not list the users.'"`
	Sent                     string `i18n:"hint='Said after a notification was sent.'"`
	LevelInfo                string `i18n:"hint='Level of a notification: information.'"`
	LevelSuccess             string `i18n:"hint='Level of importance of a notification: something went well.'"`
	LevelWarning             string `i18n:"hint='Level of importance of a notification: an alert, to look at.'"`
	LevelError               string `i18n:"hint='Level of importance of a notification: critical, to act on.'"`
	CopyUserID               string `i18n:"hint='Tooltip of the icon that copies the id (UUID) of the user who sent a notification.'"`
	NoUnread                 string `i18n:"hint='Said by the widget of the notifications when none is unread.'"`
	MarkRead                 string `i18n:"hint='Button of a notification in the widget and in the listing: marks it read.'"`
	MarkUnread               string `i18n:"hint='Button of a read notification in the listing: marks it unread again.'"`
	MarkReadConfirm          string `i18n:"hint='Question of the dialog that marks a notification read.'"`
	MarkUnreadConfirm        string `i18n:"hint='Question of the dialog that marks a notification unread again.'"`
}

var Messages_en_US = &Messages{
	ServiceTitle:             "Notifications",
	ServiceDesc:              "The attachments of the messages of the notifications, to who received or sent them.",
	ServiceFileRoute:         "an attachment (by its id)",
	ModuleDescription:        "The notification center of the users.",
	Notification:             "Notification",
	Notifications:            "Notifications",
	NotificationTitle:        "Title",
	NotificationBody:         "Text",
	NotificationLevel:        "Importance",
	NotificationSource:       "Source",
	NotificationLink:         "Link",
	NotificationCreatedAt:    "When",
	NotificationSender:       "Sent by",
	NotificationConversation: "Conversation",
	NotificationAttachments:  "Attachments",
	NotificationReadAt:       "Read at",
	MarkAllRead:              "Mark all as read",
	SeeAll:                   "See all",
	Open:                     "Open",
	Unread:                   "Unread",
	Read:                     "Read",
	All:                      "All",
	BellTitle:                "Notifications (%d unread)",
	NotificationSystem:       "System",
	Period:                   "Period",
	PeriodToday:              "Today",
	Period7d:                 "Last 7 days",
	Period30d:                "Last 30 days",
	Period90d:                "Last 90 days",
	Period365d:               "Last 365 days",
	Origin:                   "Origin",
	OriginSystem:             "System",
	OriginUsers:              "Users",
	OriginOthers:             "Others",
	Send:                     "Send notification",
	SendTo:                   "To",
	SendRequired:             "Fill in every field: who it goes to, its title, its text and its importance.",
	SendDenied:               "You may not send notifications.",
	Sent:                     "Notification sent.",
	LevelInfo:                "Information",
	LevelSuccess:             "Success",
	LevelWarning:             "Alert",
	LevelError:               "Critical",
	CopyUserID:               "Copy the user's id (UUID)",
	NoUnread:                 "No unread notifications.",
	MarkRead:                 "Mark as read",
	MarkUnread:               "Mark as unread",
	MarkReadConfirm:          "Mark this notification as read?",
	MarkUnreadConfirm:        "Mark this notification as unread?",
}

var Messages_pt_BR = &Messages{
	ServiceTitle:             "Notificações",
	ServiceDesc:              "Os anexos das mensagens das notificações, a quem as recebeu ou enviou.",
	ServiceFileRoute:         "um anexo (pelo seu id)",
	ModuleDescription:        "A central de notificações dos usuários.",
	Notification:             "Notificação",
	Notifications:            "Notificações",
	NotificationTitle:        "Título",
	NotificationBody:         "Texto",
	NotificationLevel:        "Importância",
	NotificationSource:       "Fonte",
	NotificationLink:         "Link",
	NotificationCreatedAt:    "Quando",
	NotificationSender:       "Enviada por",
	NotificationConversation: "Conversa",
	NotificationAttachments:  "Anexos",
	NotificationReadAt:       "Lida em",
	MarkAllRead:              "Marcar todas como lidas",
	SeeAll:                   "Ver todas",
	Open:                     "Abrir",
	Unread:                   "Não lidas",
	Read:                     "Lidas",
	All:                      "Todas",
	BellTitle:                "Notificações (%d não lidas)",
	NotificationSystem:       "Sistema",
	Period:                   "Período",
	PeriodToday:              "Hoje",
	Period7d:                 "Últimos 7 dias",
	Period30d:                "Últimos 30 dias",
	Period90d:                "Últimos 90 dias",
	Period365d:               "Últimos 365 dias",
	Origin:                   "Origem",
	OriginSystem:             "Sistema",
	OriginUsers:              "Usuários",
	OriginOthers:             "Outras",
	Send:                     "Enviar notificação",
	SendTo:                   "Para",
	SendRequired:             "Preencha todos os campos: para quem vai, o título, o texto e a importância.",
	SendDenied:               "Você não pode enviar notificações.",
	Sent:                     "Notificação enviada.",
	LevelInfo:                "Informação",
	LevelSuccess:             "Sucesso",
	LevelWarning:             "Alerta",
	LevelError:               "Crítica",
	CopyUserID:               "Copiar o id (UUID) do usuário",
	NoUnread:                 "Nenhuma notificação não lida.",
	MarkRead:                 "Marcar como lida",
	MarkUnread:               "Marcar como não lida",
	MarkReadConfirm:          "Marcar esta notificação como lida?",
	MarkUnreadConfirm:        "Marcar esta notificação como não lida?",
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
