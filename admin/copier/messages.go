package copier

import (
	"context"

	"github.com/go-rvq/rvq/x/i18n"
)

const I18nCopierKey i18n.ModuleKey = "I18nCopierKey"

type Messages struct {
	ModuleDescription string `i18n:"hint='Says what the module is, for whoever translates its texts.'"`
	Copy              string `i18n:"label='Copy: action', hint='Action that makes a copy of a record.'"`
	CopyOf            string `i18n:"label='Copy: title', hint='The title of a copy; %s the title of the record copied.'"`
	CopyTitle         string `i18n:"label='Copy: title field', hint='Label of the field of the title of the copy.'"`
	Done              string `i18n:"label='Copy: done', hint='Shown once records were copied; %d how many.'"`
	ErrDenied         string `i18n:"label='Copy: denied', hint='Error when a record a copy takes may not be read by the user; %s its title.'"`
}

var Messages_en_US = &Messages{
	ModuleDescription: "The copies of records.",
	Copy:              "Make a copy",
	CopyOf:            "Copy of %s",
	CopyTitle:         "Title of the copy",
	Done:              "Copied: %d record(s)",
	ErrDenied:         "No permission to read “%s”.",
}

var Messages_pt_BR = &Messages{
	ModuleDescription: "As cópias de registros.",
	Copy:              "Fazer uma cópia",
	CopyOf:            "Cópia de %s",
	CopyTitle:         "Título da cópia",
	Done:              "Copiado: %d registro(s)",
	ErrDenied:         "Sem permissão para ler “%s”.",
}

func MustGetMessages(ctx context.Context) *Messages {
	return i18n.MustGetModuleMessages(ctx, I18nCopierKey, Messages_en_US).(*Messages)
}
