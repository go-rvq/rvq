package dashboard

import (
	"context"

	"github.com/go-rvq/rvq/x/i18n"
)

// I18nDashboardKey is the module of the dashboard's words.
const I18nDashboardKey i18n.ModuleKey = "I18nDashboardKey"

// Messages are the dashboard's words (the widgets' are their modules').
type Messages struct {
	ModuleDescription string `i18n:"hint='Says what the module is, for whoever translates its texts.'"`
	Dashboard         string `i18n:"hint='Title of the dashboard page.'"`
	Period            string `i18n:"hint='Label of the period chooser at the top of the dashboard.'"`
	Period1h          string `i18n:"hint='Period chip: the last hour.'"`
	Period12h         string `i18n:"hint='Period chip: the last 12 hours.'"`
	Period24h         string `i18n:"hint='Period chip: the last 24 hours.'"`
	Period7d          string `i18n:"hint='Period chip: the last 7 days.'"`
	Period30d         string `i18n:"hint='Period chip: the last 30 days.'"`
	Period90d         string `i18n:"hint='Period chip: the last 90 days.'"`
	PeriodCustom      string `i18n:"hint='Period chip: a range of days chosen by the user.'"`
	From              string `i18n:"hint='First day of a custom period.'"`
	To                string `i18n:"hint='Last day of a custom period.'"`
	Apply             string `i18n:"hint='Button that applies a custom period.'"`
	Refresh           string `i18n:"hint='Button that loads a widget again, skipping its cache.'"`
	Open              string `i18n:"hint='Link of a widget to the listing behind it.'"`
	SeeAll            string `i18n:"hint='Link under a list of a widget to its whole listing.'"`
	Help              string `i18n:"hint='Menu item of a widget that opens its documentation.'"`
	Customize         string `i18n:"hint='Button that turns the customize mode on (hide, order, collapse).'"`
	Done              string `i18n:"hint='Button that leaves the customize mode.'"`
	Reset             string `i18n:"hint='Button that drops the customizations of the user.'"`
	MoveUp            string `i18n:"hint='Button that moves a widget before the previous one.'"`
	MoveDown          string `i18n:"hint='Button that moves a widget after the next one.'"`
	Collapse          string `i18n:"hint='Button that collapses a group of widgets.'"`
	Expand            string `i18n:"hint='Button that expands a collapsed group of widgets.'"`
	DefaultPeriod     string `i18n:"hint='Label of the period the dashboard opens with, for the user.'"`
	FollowsPeriod     string `i18n:"hint='Tooltip of a widget whose numbers follow the period chosen.'"`
	NoData            string `i18n:"hint='Said by a widget with nothing in the period.'"`
	NoWidgets         string `i18n:"hint='Said by a dashboard with no widget for the user.'"`
	New               string `i18n:"hint='Delta of a number that was zero in the previous period.'"`
	VsPrevious        string `i18n:"hint='Tooltip of the delta of a number.'"`
	LoadError         string `i18n:"hint='Said by a widget that failed to load.'"`
	Retry             string `i18n:"hint='Button that tries to load a failed widget again.'"`
	UpdatedAt         string `i18n:"hint='Footer of a widget: the time its numbers were read. {time} is the time.'"`
	ActionDone        string `i18n:"hint='Message after a button of a widget did its work.'"`
	ActionFailed      string `i18n:"hint='Message after a button of a widget failed.'"`
	Remove            string `i18n:"hint='Button that takes a widget out of the dashboard of the user (it may be added back).'"`
	Disable           string `i18n:"hint='Button of the customize mode that disables a widget: it stays in its place, shown only while customizing.'"`
	Enable            string `i18n:"hint='Button of the customize mode that enables a disabled widget again.'"`
	Disabled          string `i18n:"hint='Badge of a disabled widget, while customizing.'"`
	Drag              string `i18n:"hint='Tooltip of the handle that drags a widget to another place.'"`
	DragHint          string `i18n:"hint='Hint of the customize mode on how to move the widgets.'"`
	NewSection        string `i18n:"hint='Button and dialog title that makes a section of the user own.'"`
	EditSection       string `i18n:"hint='Button and dialog title that changes a section of the user own.'"`
	DeleteSection     string `i18n:"hint='Button that deletes a section of the user own (its widgets back to their sections).'"`
	SectionTitle      string `i18n:"hint='Field of the title of a section.'"`
	SectionIcon       string `i18n:"hint='Field of the icon of a section: a name of Material Design Icons.'"`
	Untitled          string `i18n:"hint='Shown in the place of the title of a section with none, while customizing.'"`
	AddWidgets        string `i18n:"hint='Button and dialog title that adds back the widgets removed.'"`
	AddSelected       string `i18n:"hint='Button that adds the widgets chosen back to their sections.'"`
	NoRemoved         string `i18n:"hint='Said by the dialog of the widgets to add when none was removed.'"`
	Save              string `i18n:"hint='Button that saves a dialog.'"`
	Cancel            string `i18n:"hint='Button that closes a dialog without saving.'"`
	ParentSection     string `i18n:"hint='Field of the section another section is in (a tree).'"`
	TopLevel          string `i18n:"hint='Choice of the field of the parent section: at the top, in no section.'"`
}

var Messages_en_US = &Messages{
	ModuleDescription: "The dashboard of the admin: its widgets, their period and the customizations of each user.",
	Dashboard:         "Dashboard",
	Period:            "Period",
	Period1h:          "Last hour",
	Period12h:         "Last 12 hours",
	Period24h:         "Last 24 hours",
	Period7d:          "7 days",
	Period30d:         "30 days",
	Period90d:         "90 days",
	PeriodCustom:      "Custom",
	From:              "From",
	To:                "To",
	Apply:             "Apply",
	Refresh:           "Refresh",
	Open:              "Open",
	SeeAll:            "See all",
	Help:              "About this widget",
	Customize:         "Customize",
	Done:              "Done",
	Reset:             "Restore defaults",
	MoveUp:            "Move up",
	MoveDown:          "Move down",
	Collapse:          "Collapse",
	Expand:            "Expand",
	DefaultPeriod:     "Default period",
	FollowsPeriod:     "Follows the period chosen",
	NoData:            "Nothing in this period.",
	NoWidgets:         "There is no widget for you yet.",
	New:               "new",
	VsPrevious:        "Compared with the previous period",
	LoadError:         "This widget could not be loaded.",
	Retry:             "Try again",
	UpdatedAt:         "Updated at {time}",
	ActionDone:        "Done.",
	ActionFailed:      "It could not be done: {error}",
	Remove:            "Remove from the dashboard",
	Disable:           "Disable",
	Enable:            "Enable",
	Disabled:          "Disabled",
	Drag:              "Drag to move",
	DragHint:          "Drag the widgets to place them, in a section or another.",
	NewSection:        "New section",
	EditSection:       "Edit the section",
	DeleteSection:     "Delete the section (its widgets go back to their sections)",
	SectionTitle:      "Title",
	SectionIcon:       "Icon (mdi-…)",
	Untitled:          "No title",
	AddWidgets:        "Add widgets",
	AddSelected:       "Add",
	NoRemoved:         "No widget was removed.",
	Save:              "Save",
	Cancel:            "Cancel",
	ParentSection:     "In the section",
	TopLevel:          "At the top (in no section)",
}

var Messages_pt_BR = &Messages{
	ModuleDescription: "O painel do admin: os widgets, o período deles e as personalizações de cada usuário.",
	Dashboard:         "Painel",
	Period:            "Período",
	Period1h:          "Última hora",
	Period12h:         "Últimas 12 horas",
	Period24h:         "Últimas 24 horas",
	Period7d:          "7 dias",
	Period30d:         "30 dias",
	Period90d:         "90 dias",
	PeriodCustom:      "Personalizado",
	From:              "De",
	To:                "Até",
	Apply:             "Aplicar",
	Refresh:           "Atualizar",
	Open:              "Abrir",
	SeeAll:            "Ver tudo",
	Help:              "Sobre este widget",
	Customize:         "Personalizar",
	Done:              "Concluir",
	Reset:             "Restaurar o padrão",
	MoveUp:            "Mover para cima",
	MoveDown:          "Mover para baixo",
	Collapse:          "Recolher",
	Expand:            "Expandir",
	DefaultPeriod:     "Período padrão",
	FollowsPeriod:     "Segue o período escolhido",
	NoData:            "Nada neste período.",
	NoWidgets:         "Ainda não há widget para você.",
	New:               "novo",
	VsPrevious:        "Comparado com o período anterior",
	LoadError:         "Não foi possível carregar este widget.",
	Retry:             "Tentar de novo",
	UpdatedAt:         "Atualizado às {time}",
	ActionDone:        "Feito.",
	ActionFailed:      "Não foi possível: {error}",
	Remove:            "Remover do painel",
	Disable:           "Desativar",
	Enable:            "Ativar",
	Disabled:          "Desativado",
	Drag:              "Arraste para mover",
	DragHint:          "Arraste os widgets para posicioná-los, na mesma seção ou em outra.",
	NewSection:        "Nova seção",
	EditSection:       "Editar a seção",
	DeleteSection:     "Excluir a seção (os widgets voltam às suas seções)",
	SectionTitle:      "Título",
	SectionIcon:       "Ícone (mdi-…)",
	Untitled:          "Sem título",
	AddWidgets:        "Adicionar widgets",
	AddSelected:       "Adicionar",
	NoRemoved:         "Nenhum widget foi removido.",
	Save:              "Salvar",
	Cancel:            "Cancelar",
	ParentSection:     "Dentro da seção",
	TopLevel:          "No topo (em nenhuma seção)",
}

// MessagesOf are the dashboard's words of ctx's language.
func MessagesOf(ctx context.Context) *Messages {
	return i18n.MustGetModuleMessages(ctx, I18nDashboardKey, Messages_en_US).(*Messages)
}

// PeriodLabel is the chip's word of the period key.
func (m *Messages) PeriodLabel(key string) string {
	switch key {
	case Period1h:
		return m.Period1h
	case Period12h:
		return m.Period12h
	case Period24h:
		return m.Period24h
	case Period7d:
		return m.Period7d
	case Period30d:
		return m.Period30d
	case Period90d:
		return m.Period90d
	case PeriodCustom:
		return m.PeriodCustom
	}
	return key
}
