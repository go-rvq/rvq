package db_tools

import (
	"context"
	"fmt"
	"strings"
	"text/template"

	h "github.com/go-rvq/htmlgo"
	"github.com/go-rvq/rvq/admin/worker"
	"github.com/go-rvq/rvq/x/i18n"
	db_tools "github.com/go-rvq/rvq/x/packages/db-tools"
	. "github.com/go-rvq/rvq/x/ui/vuetify"
	"golang.org/x/text/language"
)

const MessagesKey i18n.ModuleKey = "rqv-admin/db-tools"

func GetMessages(ctx context.Context) *Messages {
	return i18n.MustGetModuleMessages(ctx, MessagesKey, Messages_en_US).(*Messages)
}

type MessagesPersistence struct {
	Enabled     string
	Disabled    string
	Title       string
	Days        string
	Weeks       string
	Months      string
	Years       string
	NoOther     string
	OtherDays   string
	OtherWeeks  string
	OtherMonths string
	OtherYears  string
	Template    string
}

func (p *MessagesPersistence) Format(per *db_tools.Persistence) (s h.RawHTML, err error) {
	var t *template.Template
	if t, err = template.New("messages_persistence").
		Funcs(map[string]any{
			"join_and": func(sep, lastSep string, elem ...string) string {
				var valid []string
				for _, s := range elem {
					if len(s) > 0 {
						valid = append(valid, s)
					}
				}

				elem = valid
				switch len(elem) {
				case 0:
					return ""
				case 1:
					return elem[0]
				}

				s := strings.Join(elem[:len(elem)-1], sep)
				if len(elem) > 1 {
					return s + lastSep + elem[len(elem)-1]
				}
				return s
			},
		}).
		Parse(p.Template); err != nil {
		return
	}

	var (
		data = map[string]any{
			"enabled": p.Enabled,
			"valid":   !per.IsZero(),
			"days":    "",
			"weeks":   "",
			"months":  "",
			"years":   "",
		}
		other string
		w     strings.Builder
	)

	if !per.Enabled {
		data["enabled"] = p.Disabled
	}

	if per.Days > 0 {
		data["days"] = fmt.Sprintf(p.Days, per.Days)
	}

	if per.Weeks > 0 {
		data["weeks"] = fmt.Sprintf(p.Weeks, per.Weeks)
	}

	if per.Months > 0 {
		data["months"] = fmt.Sprintf(p.Months, per.Months)
	}

	if per.Years > 0 {
		data["years"] = fmt.Sprintf(p.Years, per.Years)
	}

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

	data["other"] = other

	err = t.Execute(&w, data)
	s = h.RawHTML(w.String())
	return
}

type Messages struct {
	AutoBackup                           string
	DatabaseAutoBackup                   string
	Database                             string
	CreateBackup                         string
	ConfigureBackupPersistence           string
	CreateDatabaseBackup                 string
	DatabaseOlderBackupsRemover          string
	BackupStartedInBackgroundJobTemplate string
	MessageFormMessage                   string
	Backups                              string
	CreatedAt                            string
	DbName                               string
	Size                                 string
	Message                              string
	Actions                              string
	BackupDetailTemplate                 string
	BackupRemovedTemplate                h.RawHTML
	BackupRemoveConfirmTemplate          h.RawHTML
	Auto                                 string
	Persistence                          MessagesPersistence
	PersistenceEnabled                   string
	PersistenceYears                     string
	PersistenceMonths                    string
	PersistenceWeeks                     string
	PersistenceDays                      string
	PersistenceOther                     string
}

var (
	Messages_en_US = &Messages{
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
			Title:       "Persistence",
			Days:        "%d days",
			Weeks:       "%d weeks",
			Months:      "%d months",
			Years:       "%d years",
			NoOther:     "Não definido",
			OtherDays:   "other days",
			OtherWeeks:  "other weeks",
			OtherMonths: "other months",
			OtherYears:  "other years",
			Template:    `{{if .valid}}Persists for {{- join_and " " ", " " e " .enabled .days .weeks .months .years .other}}{{else}}Not persists{{end}}.`,
		},
	}

	Messages_pt_BR = &Messages{
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
			Template:    `{{if .valid}}{{.enabled}}, manter por {{join_and ", " " e " .days .weeks .months .years .other}}{{else}}<span class='text-warning'>Não manter</span>{{end}}.`,
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
	b.RegisterForModules(language.English, MessagesKey, Messages_pt_BR).
		RegisterForModules(language.BrazilianPortuguese, MessagesKey, Messages_pt_BR)
}
