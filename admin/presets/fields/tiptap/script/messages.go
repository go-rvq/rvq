package script

import (
	"context"
	"fmt"

	"github.com/gad-lang/gad/parser/source"
	h "github.com/go-rvq/htmlgo"
	"github.com/go-rvq/rvq/x/i18n"
	"golang.org/x/text/language"
)

const MessagesKey i18n.ModuleKey = "rqv-admin/tiptap"

func GetMessages(ctx context.Context) *Messages {
	return i18n.MustGetModuleMessages(ctx, MessagesKey, Messages_en_US).(*Messages)
}

type Messages struct {
	ErrScriptFailure     i18n.ErrorString
	ParseErrorTemplate   string
	CompileErrorTemplate string
	RunErrorTemplate     string
	EditorUsage          h.HTMLComponent
}

func (m *Messages) FormateTypeError(errType ScriptErrorType, pos source.FilePos, message string) string {
	var template *string
	switch errType {
	case ScriptErrorTypeParse:
		template = &m.ParseErrorTemplate
	case ScriptErrorTypeCompile:
		template = &m.CompileErrorTemplate
	case ScriptErrorTypeRun:
		template = &m.RunErrorTemplate
	default:
		return ""
	}
	return fmt.Sprintf(*template, pos.Line, pos.Column, message)
}

var (
	Messages_en_US = &Messages{
		ErrScriptFailure:     i18n.ErrorString(ErrScriptFailure.Error()),
		ParseErrorTemplate:   "Parse ERROR at [%d:%d]: %s",
		CompileErrorTemplate: "Compile ERROR at [%d:%d]: %s",
		RunErrorTemplate:     "Run ERROR at [%d:%d]: %s",
	}

	Messages_pt_BR = &Messages{
		ErrScriptFailure:     "Falha de script",
		ParseErrorTemplate:   "ERRO de interpretação em [%d:%d]: %s",
		CompileErrorTemplate: "ERRO de compilação em [%d:%d]: %s",
		RunErrorTemplate:     "ERRO de execução em [%d:%d]: %s",
	}
)

func ConfigureMessages(b *i18n.Builder) {
	b.RegisterForModules(language.English, MessagesKey, Messages_pt_BR).
		RegisterForModules(language.BrazilianPortuguese, MessagesKey, Messages_pt_BR)
}
