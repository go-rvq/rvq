package integration_test

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"sync/atomic"
	"testing"

	"github.com/go-rvq/rvq/admin/presets"
	"github.com/go-rvq/rvq/admin/presets/actions"
	"github.com/go-rvq/rvq/admin/presets/gorm2op"
	"github.com/go-rvq/rvq/web"
	"github.com/go-rvq/rvq/web/multipartestutils"
	"github.com/go-rvq/rvq/x/i18n"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

// A bulk action that fails keeps its dialog open, saying why: a validation
// error as it is, any other error as its global error — never closed as done.
func TestBulkActionErrorKeepsTheDialog(t *testing.T) {
	db, err := gorm.Open(sqlite.Open(fmt.Sprintf("file:bulkerr_%d?mode=memory&cache=shared", atomic.AddInt64(&paSeq, 1))), &gorm.Config{})
	if err != nil {
		t.Fatal(err)
	}
	if err = db.AutoMigrate(&PAPost{}); err != nil {
		t.Fatal(err)
	}
	db.Create(&PAPost{ID: 1, Title: "old"})

	p := presets.New(i18n.New()).URIPrefix("/admin")
	p.DataOperator(gorm2op.DataOperator(db))
	l := p.Model(&PAPost{}, presets.ModelWithID("posts")).Listing("ID", "Title")
	l.BulkAction("Refused").UpdateFunc(func([]string, *web.EventContext, *web.EventResponse) error {
		var vErr web.ValidationErrors
		vErr.GlobalError("refused: translate the parent first")
		return &vErr
	})
	l.BulkAction("Failed").UpdateFunc(func([]string, *web.EventContext, *web.EventResponse) error {
		return errors.New("failed: no connection")
	})
	mux := http.NewServeMux()
	p.Build(mux)

	for name, want := range map[string]string{"Refused": "refused: translate the parent first", "Failed": "failed: no connection"} {
		r := multipartestutils.NewMultipartBuilder().PageURL("/admin/posts").EventFunc(actions.DoBulkAction).
			Query(presets.ParamBulkActionName, name).Query(presets.ParamSelectedIds, "1").BuildEventFuncRequest()
		w := httptest.NewRecorder()
		mux.ServeHTTP(w, r)
		body := w.Body.String()
		if !strings.Contains(body, want) {
			t.Errorf("%s: the error not shown: %.300s", name, body)
		}
		var res struct {
			RunScript string `json:"runScript"`
		}
		if err := json.Unmarshal(w.Body.Bytes(), &res); err != nil {
			t.Fatal(err)
		}
		if strings.Contains(res.RunScript, "closer.show = false") {
			t.Errorf("%s: the dialog closed: %s", name, res.RunScript)
		}
	}
}
