package integration_test

import (
	"fmt"
	"net/http/httptest"
	"regexp"
	"strings"
	"testing"

	"github.com/go-rvq/rvq/admin/presets"
	"github.com/go-rvq/rvq/admin/presets/actions"
	"github.com/go-rvq/rvq/admin/presets/gorm2op"
	"github.com/go-rvq/rvq/web"
	"github.com/go-rvq/rvq/web/multipartestutils"
	"github.com/go-rvq/rvq/x/i18n"
)

type PageHostProduct struct {
	ID   uint
	Name string
}

// Numa listagem em PÁGINA o botão "+" é renderizado pelo layout na barra de
// ações, longe da listagem. O host de criação guarda o estado em `vars` para que
// os dois lados falem do mesmo objeto — e o bloco guardado dele fica JUNTO com os
// de detalhe e edição, dentro do conteúdo, que é onde o overlay tem de abrir.
func TestPageListingHostsAreTogether(t *testing.T) {
	db := TestDB
	if err := db.AutoMigrate(&PageHostProduct{}); err != nil {
		t.Fatal(err)
	}
	db.Exec("DELETE FROM page_host_products")
	db.Create(&PageHostProduct{Name: "P1"})

	b := presets.New(i18n.New()).URIPrefix("/admin")
	b.DataOperator(gorm2op.DataOperator(db))
	mb := b.Model(&PageHostProduct{}).URIName("produtos")
	mb.Detailing("Name") // para os três hosts existirem

	w := httptest.NewRecorder()
	b.ServeHTTP(w, httptest.NewRequest("GET", "/admin/produtos", nil))
	body := w.Body.String()
	for _, r := range [][2]string{{`<`, "<"}, {`>`, ">"}, {`\n`, "\n"}, {`\t`, "\t"}, {`\"`, `"`}, {"&#39;", "'"}} {
		body = strings.ReplaceAll(body, r[0], r[1])
	}

	// o estado do host de criação vive em vars…
	if !strings.Contains(body, `"`+presets.ListingNewScope+`": $closer(`) {
		t.Errorf("o host de criação não foi declarado em vars:\n%s", firstMatch(body, `assign=[^>]{0,200}`))
	}

	// …e é assim que o botão da barra de ações o abre
	if !strings.Contains(body, "vars."+presets.ListingNewScope+".show = true") {
		t.Errorf("o botão + não abre o host em vars:\n%s", firstMatch(body, `data-event='new'[^>]{0,120}`))
	}

	// os três blocos guardados ficam lado a lado, no conteúdo
	// numa página os três vivem em vars: é o que lhes dá endereço global, sem o
	// qual o overlay teria de ser renderizado dentro do portal do host
	guards := []string{
		"v-if='vars." + presets.ListingNewScope + "?.show'",
		"v-if='vars." + presets.ListingItemDetailScope + "?.show'",
		"v-if='vars." + presets.ListingItemEditScope + "?.show'",
	}
	positions := make([]int, len(guards))
	for i, g := range guards {
		if positions[i] = strings.Index(body, g); positions[i] < 0 {
			t.Fatalf("bloco ausente: %s", g)
		}
	}

	// todos depois da tabela, e nenhum dentro da barra de ações (que vem antes)
	table := strings.Index(body, "<v-table")
	if table < 0 {
		table = strings.Index(body, "<table")
	}
	if table < 0 {
		table = strings.Index(body, "DataTable")
	}
	if table < 0 {
		t.Fatal("a listagem não renderizou a tabela")
	}
	for i, p := range positions {
		if p < table {
			t.Errorf("o bloco %q ficou ANTES da tabela (na barra de ações?)", guards[i])
		}
	}
}

// Um drawer só se dimensiona contra a janela quando está no portal do LAYOUT;
// dentro do portal de quem o abriu ele se registra na caixa que o contém e abre
// abaixo do app bar. O host cujo estado está em `vars` manda o endereço do
// closer junto (ParamCloserRef), e é isso que permite responder lá.
func TestPageListingNewOpensInTheLayoutDrawer(t *testing.T) {
	db := TestDB
	if err := db.AutoMigrate(&PageHostProduct{}); err != nil {
		t.Fatal(err)
	}

	b := presets.New(i18n.New()).URIPrefix("/admin")
	b.DataOperator(gorm2op.DataOperator(db))
	b.Model(&PageHostProduct{}).URIName("produtos")

	ref := "vars." + presets.ListingNewScope
	w := httptest.NewRecorder()
	b.ServeHTTP(w, multipartestutils.NewMultipartBuilder().
		PageURL("/admin/produtos").
		EventFunc("presets_New").
		Query(presets.ParamTargetPortal, "hostPortal").
		Query(presets.ParamOverlay, "RightDrawer").
		Query(presets.ParamCloserProvided, "true").
		Query(presets.ParamCloserRef, ref).
		BuildEventFuncRequest())

	body := strings.ReplaceAll(w.Body.String(), `"`, `"`)

	if !strings.Contains(body, `"name":"`+actions.RightDrawer.PortalName()+`"`) {
		t.Errorf("a resposta não foi para o portal do layout:\n%s", firstMatch(body, `"name":"[^"]*"`))
	}
	if !strings.Contains(body, `:closer='`+ref+`'`) {
		t.Errorf("o drawer não ligou no closer do host:\n%s", firstMatch(body, `:closer='[^']*'`))
	}

	// e sem a referência continua respondendo no portal de quem pediu
	w = httptest.NewRecorder()
	b.ServeHTTP(w, multipartestutils.NewMultipartBuilder().
		PageURL("/admin/produtos").
		EventFunc("presets_New").
		Query(presets.ParamTargetPortal, "hostPortal").
		Query(presets.ParamOverlay, "RightDrawer").
		Query(presets.ParamCloserProvided, "true").
		BuildEventFuncRequest())

	if !strings.Contains(w.Body.String(), `"name":"hostPortal"`) {
		t.Errorf("sem referência, a resposta deveria ficar no portal de quem pediu:\n%s",
			firstMatch(w.Body.String(), `"name":"[^"]*"`))
	}
}

// O formulário devolvido com os erros da validação volta para onde foi aberto:
// o portal do layout, quando o host manda o endereço do closer — no portal de
// quem pediu ele seria um segundo drawer, atrás do aberto, e os erros não
// apareceriam.
func TestEditErrorsInTheLayoutDrawer(t *testing.T) {
	db := TestDB
	if err := db.AutoMigrate(&PageHostProduct{}); err != nil {
		t.Fatal(err)
	}
	p := &PageHostProduct{Name: "zz-edit-errors"}
	if err := db.Create(p).Error; err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { db.Delete(p) })

	b := presets.New(i18n.New()).URIPrefix("/admin")
	b.DataOperator(gorm2op.DataOperator(db))
	m := b.Model(&PageHostProduct{}).URIName("produtos")
	m.Editing("Name").Validators.AppendFunc(func(obj any, _ presets.FieldModeStack, _ *web.EventContext) (errs web.ValidationErrors) {
		if obj.(*PageHostProduct).Name == "" {
			errs.FieldError("Name", "zz-name-required")
		}
		return
	})

	update := func(ref string) string {
		mb := multipartestutils.NewMultipartBuilder().
			PageURL("/admin/produtos").
			EventFunc("presets_Update").
			Query(presets.ParamID, fmt.Sprint(p.ID)).
			Query(presets.ParamTargetPortal, "hostPortal").
			Query(presets.ParamOverlay, "RightDrawer").
			Query(presets.ParamCloserProvided, "true").
			AddField("Name", "")
		if ref != "" {
			mb = mb.Query(presets.ParamCloserRef, ref)
		}
		// the stamp of the record is not what this is about
		r := mb.BuildEventFuncRequest()
		r = r.WithContext(web.WithSkipFormSign(r.Context()))
		w := httptest.NewRecorder()
		b.ServeHTTP(w, r)
		return w.Body.String()
	}

	body := update("vars." + presets.DetailingEditScope)
	if !strings.Contains(body, "zz-name-required") {
		t.Fatalf("sem o erro da validação:\n%s", firstMatch(body, `(?i)error[^"]{0,300}`))
	}
	if !strings.Contains(body, `"name":"`+actions.RightDrawer.PortalName()+`"`) {
		t.Errorf("os erros não voltaram ao portal do layout:\n%s", firstMatch(body, `"name":"[^"]*"`))
	}
	// sem a referência, no portal de quem pediu
	if body = update(""); !strings.Contains(body, `"name":"hostPortal"`) {
		t.Errorf("sem referência, os erros deveriam ficar no portal de quem pediu:\n%s", firstMatch(body, `"name":"[^"]*"`))
	}
}

func firstMatch(s, pattern string) string {
	if m := regexp.MustCompile(pattern).FindString(s); m != "" {
		return m
	}
	return "(não encontrado)"
}
