package integration_test

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/go-rvq/rvq/admin/presets"
	"github.com/go-rvq/rvq/admin/presets/actions"
	. "github.com/go-rvq/rvq/admin/presets/integration"
	"github.com/go-rvq/rvq/web"
	"github.com/go-rvq/rvq/web/multipartestutils"
)

type afHost struct {
	ID   string `gorm:"primarykey"`
	Text string
}

type afLinha struct {
	Valor string
	Obs   string
}

type afPedido struct {
	Nome   string
	Linhas []*afLinha
}

// The list editor of the form of an action (go-rvq/rvq#17): its + sends the
// action again — the row added, the dialog drawn again with what was typed,
// the action not done —; the items given as JSON added; done, the action
// gets the rows kept (not the removed), in their order.
func TestActionFormListEditor(t *testing.T) {
	db, err := NewNestedModelsDB()
	if err != nil {
		t.Fatal(err)
	}
	if err = db.AutoMigrate(&afHost{}); err != nil {
		t.Fatal(err)
	}
	p := NewNestedModelsApp(db)
	host := p.Model(&afHost{}, presets.ModelConfig().SetId("af-hosts"))
	host.Listing("ID", "Text")

	lmb := presets.NewModelBuilder(p, &afLinha{})
	led := lmb.Editing("Valor", "Obs")
	pmb := presets.NewModelBuilder(p, &afPedido{})
	ped := pmb.Editing("Nome", "Linhas")
	ped.Field("Linhas").Nested(presets.NestedSlice(lmb, &led.FieldsBuilder).
		SetComponentBuilder(presets.NewListEditorTableBuilder()))

	var got *afPedido
	presets.ActionForm[*afPedido](host.Detailing("Text").Action("Pedir"), ped,
		func(ctx *presets.ActionFormContext[*afPedido]) error {
			got = ctx.Form
			return nil
		}).Build()
	p.Build()
	db.Create(&afHost{ID: "h1", Text: "H"})

	serve := func(mp *multipartestutils.Builder) string {
		t.Helper()
		r := mp.BuildEventFuncRequest()
		r = r.WithContext(web.WithSkipFormSign(r.Context()))
		w := httptest.NewRecorder()
		p.ServeHTTP(w, r)
		if w.Code != http.StatusOK {
			t.Fatalf("%s: %d %s", r.URL, w.Code, w.Body.String())
		}
		return w.Body.String()
	}
	acao := func() *multipartestutils.Builder {
		return multipartestutils.NewMultipartBuilder().PageURL("/admin/af-hosts").EventFunc(actions.DoAction).
			Query(presets.ParamAction, "Pedir").Query(presets.ParamID, "h1").
			Query(presets.ParamOverlay, string(actions.Dialog)).
			AddField("Nome", "pedido-1").AddField("Linhas."+presets.ListEditorPresentField, "1").
			AddField("Linhas[0].Valor", "10").AddField("Linhas[0]."+presets.ListEditorPositionField, "0")
	}

	// the dialog: its + sends the action again
	dialogo := serve(multipartestutils.NewMultipartBuilder().PageURL("/admin/af-hosts").EventFunc(actions.Action).
		Query(presets.ParamAction, "Pedir").Query(presets.ParamID, "h1"))
	if !strings.Contains(dialogo, presets.ParamActionFormAddRow) {
		t.Errorf("the + of the list of the dialog does not send the action: %.400s", dialogo)
	}

	// +: a row more, what was typed kept, not done
	body := serve(acao().Query(presets.ParamActionFormAddRow, "Linhas"))
	if got != nil {
		t.Fatal("the + did the action")
	}
	if !strings.Contains(body, "Linhas[1]") || !strings.Contains(body, "pedido-1") {
		t.Errorf("the row not added, or the form not kept: %.600s", body)
	}

	// the items given as JSON
	body = serve(acao().Query(presets.ParamActionFormAddRow, "Linhas").
		AddField(presets.ListEditorAddRowParamJsonItems, `[{"Valor":"77"},{"Valor":"88"}]`))
	if got != nil || !strings.Contains(body, "77") || !strings.Contains(body, "88") || !strings.Contains(body, "Linhas[2]") {
		t.Errorf("the items given not added: %.600s", body)
	}

	// done: the rows kept, in their order
	serve(acao().
		AddField("Linhas[1].Valor", "20").AddField("Linhas[1]."+presets.ListEditorPositionField, "1").
		AddField("Linhas[1]."+presets.ListEditorDeletedField, "true").
		AddField("Linhas[2].Valor", "30").AddField("Linhas[2]."+presets.ListEditorPositionField, "2").
		AddField("Linhas[2]."+presets.ListEditorNewField, "true"))
	if got == nil {
		t.Fatal("not done")
	}
	var valores []string
	for _, l := range got.Linhas {
		valores = append(valores, l.Valor)
	}
	if strings.Join(valores, ",") != "10,30" || got.Nome != "pedido-1" {
		t.Errorf("the form done: %q %v", got.Nome, valores)
	}
}
