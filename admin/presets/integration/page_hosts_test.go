package integration_test

import (
	"net/http/httptest"
	"regexp"
	"strings"
	"testing"

	"github.com/go-rvq/rvq/admin/presets"
	"github.com/go-rvq/rvq/admin/presets/gorm2op"
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
	guards := []string{
		"v-if='vars." + presets.ListingNewScope + "?.show'",
		"v-if='" + presets.ListingItemDetailScope + "?.show'",
		"v-if='" + presets.ListingItemEditScope + "?.show'",
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

func firstMatch(s, pattern string) string {
	if m := regexp.MustCompile(pattern).FindString(s); m != "" {
		return m
	}
	return "(não encontrado)"
}
