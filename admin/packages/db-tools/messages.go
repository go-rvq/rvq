package db_tools

import (
	"context"
	"fmt"
	"strings"

	"github.com/gad-lang/gad"
	h "github.com/go-rvq/htmlgo"
	"github.com/go-rvq/rvq/admin/worker"
	"github.com/go-rvq/rvq/x/i18n"
	"github.com/go-rvq/rvq/x/i18n/gadttpl"
	db_tools "github.com/go-rvq/rvq/x/packages/db-tools"
	. "github.com/go-rvq/rvq/x/ui/vuetify"
	"golang.org/x/text/language"
)

const MessagesKey i18n.ModuleKey = "rqv-admin/db-tools"

func GetMessages(ctx context.Context) *Messages {
	return i18n.MustGetModuleMessages(ctx, MessagesKey, Messages_en_US).(*Messages)
}

type MessagesPersistence struct {
	Enabled     string           `i18n:"type=html, hint='Shown when removing old backups is on (HTML).'"`
	Disabled    string           `i18n:"type=html, hint='Shown when removing old backups is off (HTML).'"`
	Title       string           `i18n:"hint='Title of the backup persistence.'"`
	Days        string           `i18n:"hint='How many days are kept.', fields=(;'%d'='the number')"`
	Weeks       string           `i18n:"hint='How many weeks are kept.', fields=(;'%d'='the number')"`
	Months      string           `i18n:"hint='How many months are kept.', fields=(;'%d'='the number')"`
	Years       string           `i18n:"hint='How many years are kept.', fields=(;'%d'='the number')"`
	NoOther     string           `i18n:"hint='Shown when no time is set for the other backups.'"`
	OtherDays   string           `i18n:"hint='The other backups kept by days.'"`
	OtherWeeks  string           `i18n:"hint='The other backups kept by weeks.'"`
	OtherMonths string           `i18n:"hint='The other backups kept by months.'"`
	OtherYears  string           `i18n:"hint='The other backups kept by years.'"`
	Template    gadttpl.Template `i18n:"type=gadt, hint='Summary of the persistence (HTML).', fields=(;valid='whether it is set', enabled='the Enabled or Disabled text', days='the days part', weeks='the weeks part', months='the months part', years='the years part', other='the other part', join_and='join_and(sep, lastSep, parts…): the parts not empty, sep between them and lastSep before the last')"`
}

func (p *MessagesPersistence) Format(per *db_tools.Persistence) (s h.RawHTML, err error) {
	data := gad.Dict{
		"enabled":  gad.Str(p.Enabled),
		"valid":    gad.Bool(!per.IsZero()),
		"days":     gad.Str(""),
		"weeks":    gad.Str(""),
		"months":   gad.Str(""),
		"years":    gad.Str(""),
		"join_and": gadttpl.JoinAnd,
	}

	if !per.Enabled {
		data["enabled"] = gad.Str(p.Disabled)
	}

	if per.Days > 0 {
		data["days"] = gad.Str(fmt.Sprintf(p.Days, per.Days))
	}

	if per.Weeks > 0 {
		data["weeks"] = gad.Str(fmt.Sprintf(p.Weeks, per.Weeks))
	}

	if per.Months > 0 {
		data["months"] = gad.Str(fmt.Sprintf(p.Months, per.Months))
	}

	if per.Years > 0 {
		data["years"] = gad.Str(fmt.Sprintf(p.Years, per.Years))
	}

	var other string
	switch per.Other {
	case db_tools.PersistenceOtherYears:
		other = p.OtherYears
	case db_tools.PersistenceOtherMonths:
		other = p.OtherMonths
	case db_tools.PersistenceOtherWeeks:
		other = p.OtherWeeks
	case db_tools.PersistenceOtherDays:
		other = p.OtherDays
	}

	data["other"] = gad.Str(other)

	var out string
	if out, err = p.Template.Render(data); err != nil {
		return
	}
	s = h.RawHTML(out)
	return
}

type Messages struct {
	ModuleDescription                    string              `i18n:"hint='Says what the module is, for whoever translates its texts.'"`
	AutoBackup                           string              `i18n:"hint='Name of the automatic backup settings.'"`
	DatabaseAutoBackup                   string              `i18n:"hint='Title of the automatic database backup job.'"`
	Database                             string              `i18n:"hint='Title of the database tools.'"`
	CreateBackup                         string              `i18n:"hint='Button that makes a database backup now.'"`
	ConfigureBackupPersistence           string              `i18n:"hint='Action that sets how long backups are kept.'"`
	CreateDatabaseBackup                 string              `i18n:"hint='Title of the job that makes a database backup.'"`
	DatabaseOlderBackupsRemover          string              `i18n:"label='Older backups remover', hint='Title of the job that removes the backups no longer kept.'"`
	BackupStartedInBackgroundJobTemplate string              `i18n:"label='Backup started', hint='Shown once a backup job was started.'"`
	MessageFormMessage                   string              `i18n:"label='Backup message', hint='Label of the note given to a backup.'"`
	Backups                              string              `i18n:"hint='Title of the list of backups.'"`
	CreatedAt                            string              `i18n:"hint='Column with when a backup was made.'"`
	DbName                               string              `i18n:"label='Database name', hint='Column with the database a backup is of.'"`
	Size                                 string              `i18n:"hint='Column with the size of a backup.'"`
	Message                              string              `i18n:"hint='Column with the note given to a backup.'"`
	Actions                              string              `i18n:"hint='Column with the actions of a backup.'"`
	BackupDetailTemplate                 string              `i18n:"label='Backup detail', hint='Line of the job log about a backup.', fields=(;'%s'='its detail')"`
	BackupRemovedTemplate                h.RawHTML           `i18n:"type=html, label='Backup removed', hint='Shown once a backup was removed (HTML).', fields=(;'%s'='the backup')"`
	BackupRemoveConfirmTemplate          h.RawHTML           `i18n:"type=html, label='Backup remove confirmation', hint='Asks to confirm the removal of a backup (HTML).', fields=(;'%s'='the backup')"`
	Auto                                 string              `i18n:"hint='Mark of a backup made automatically.'"`
	Persistence                          MessagesPersistence `i18n:"hint='The words of how long backups are kept.'"`
	PersistenceEnabled                   string              `i18n:"label='Persistence: enabled', hint='Label of whether old backups are removed.'"`
	PersistenceYears                     string              `i18n:"label='Persistence: years', hint='Label of for how many years a yearly backup is kept.'"`
	PersistenceMonths                    string              `i18n:"label='Persistence: months', hint='Label of for how many months a monthly backup is kept.'"`
	PersistenceWeeks                     string              `i18n:"label='Persistence: weeks', hint='Label of for how many weeks a weekly backup is kept.'"`
	PersistenceDays                      string              `i18n:"label='Persistence: days', hint='Label of for how many days a daily backup is kept.'"`
	PersistenceOther                     string              `i18n:"label='Persistence: other', hint='Label of how long the other backups are kept.'"`
}

var (
	Messages_en_US = &Messages{
		ModuleDescription:                    "The database tools: backups, their schedule and how long they are kept.",
		PersistenceEnabled:                   "Enabled",
		PersistenceYears:                     "Years",
		PersistenceMonths:                    "Months",
		PersistenceWeeks:                     "Weeks",
		PersistenceDays:                      "Days",
		PersistenceOther:                     "Other",
		Auto:                                 "Auto",
		AutoBackup:                           "Auto Backup",
		DatabaseAutoBackup:                   "Database Auto Backup",
		Database:                             "Data Base",
		CreateBackup:                         "Create Backup",
		CreateDatabaseBackup:                 "Create Database Backup",
		DatabaseOlderBackupsRemover:          "Database Older Backups Remover",
		BackupStartedInBackgroundJobTemplate: "Backup Started In Background",
		MessageFormMessage:                   "Message",
		Backups:                              "Backups",
		CreatedAt:                            "Created At",
		DbName:                               "DB Name",
		Size:                                 "Size",
		Message:                              "Message",
		Actions:                              "Actions",
		BackupDetailTemplate:                 "Backup Detail: %s",
		BackupRemovedTemplate:                "Backup Removed: %s",
		BackupRemoveConfirmTemplate:          "Backup Remove Confirm: %s",
		ConfigureBackupPersistence:           "Configure Backup Persistence",
		Persistence: MessagesPersistence{
			Enabled:     "<b class='text-primary'>ENABLED</b>",
			Disabled:    "<b class='text-warning'>NOT ENABLED</b>",
			Title:       "Persistence",
			Days:        "%d days",
			Weeks:       "%d weeks",
			Months:      "%d months",
			Years:       "%d years",
			NoOther:     "Not defined",
			OtherDays:   "other days",
			OtherWeeks:  "other weeks",
			OtherMonths: "other months",
			OtherYears:  "other years",
			Template:    `{ if valid begin }{= enabled }, keep for {= join_and(", ", " and ", days, weeks, months, years, other) }{ else }<span class='text-warning'>Do not keep</span>{ end }.`,
		},
	}

	Messages_pt_BR = &Messages{
		ModuleDescription:                    "As ferramentas do banco de dados: cópias de segurança, seu agendamento e por quanto tempo são mantidas.",
		AutoBackup:                           "Cópia de Segurança Automática",
		DatabaseAutoBackup:                   "Cópia de Segurança Automática do Banco de Dados",
		Database:                             "Banco de Dados",
		CreateBackup:                         "Criar Nova Cópia de Segurança",
		ConfigureBackupPersistence:           "Configurar armezagem da Cópia de Segurança",
		CreateDatabaseBackup:                 "Criar Cópia de Segurança do Banco de Dados",
		DatabaseOlderBackupsRemover:          "Remover Cópias de Segurança antigas do Banco de Dados",
		BackupStartedInBackgroundJobTemplate: `A Cópia de Segurança está sendo gerada em segundo plano. Acesse {link} para acompanhar.`,
		MessageFormMessage:                   "Mensagem",
		Backups:                              "Cópias de Segurança",
		CreatedAt:                            "Criado em",
		DbName:                               "Banco de Dados",
		Size:                                 "Tamanho",
		Message:                              "Mensagem",
		Actions:                              "Ações",
		BackupDetailTemplate:                 "Detalhes do Backup: %s",
		BackupRemovedTemplate:                "Cópia de Segurança <b>%s</b> EXCLUÍDA com sucesso",
		BackupRemoveConfirmTemplate:          "Tem certeza que deseja EXCLUIR a Cópia de Segurança <b>%s</b>?",
		Auto:                                 "Automático",
		Persistence: MessagesPersistence{
			Title:       "Armazenamento",
			Enabled:     "<b class='text-primary'>ATIVADO</b>",
			Disabled:    "<b class='text-warning'>NÃO ATIVADO</b>",
			Days:        "%d dias",
			Weeks:       "%d semanas",
			Months:      "%d meses",
			Years:       "%d anos",
			OtherDays:   "demais dias",
			OtherWeeks:  "demais semanas",
			OtherMonths: "demais meses",
			OtherYears:  "demais anos",
			Template:    `{ if valid begin }{= enabled }, manter por {= join_and(", ", " e ", days, weeks, months, years, other) }{ else }<span class='text-warning'>Não manter</span>{ end }.`,
		},
		PersistenceEnabled: "Ativado",
		PersistenceYears:   "Anos",
		PersistenceMonths:  "Meses",
		PersistenceWeeks:   "Semanas",
		PersistenceDays:    "Dias",
		PersistenceOther:   "Outros",
	}
)

func (m *Messages) BackupStartedInBackgroundJob(wb *worker.Builder, ctx context.Context) h.RawHTML {
	mb := wb.ModelBuilder()
	link, _ := h.Marshal(VBtn(mb.TTitlePlural(ctx)).
		Density(DensityCompact).
		Attr("href", mb.Info().ListingHref()).
		Tag("a"), ctx)
	return h.RawHTML(strings.ReplaceAll(m.BackupStartedInBackgroundJobTemplate, "{link}", string(link)))
}

func (m *Messages) BackupRemoved(detail string) h.RawHTML {
	return h.RawHTML(fmt.Sprintf(string(m.BackupRemovedTemplate), detail))
}

func (m *Messages) BackupRemoveConfirm(detail string) h.RawHTML {
	return h.RawHTML(fmt.Sprintf(string(m.BackupRemoveConfirmTemplate), detail))
}

func ConfigureMessages(b *i18n.Builder) {
	b.RegisterForModules(language.English, MessagesKey, Messages_en_US).
		RegisterForModules(language.BrazilianPortuguese, MessagesKey, Messages_pt_BR)
}
