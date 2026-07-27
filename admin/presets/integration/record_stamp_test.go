package integration_test

import (
	"net/http/httptest"
	"regexp"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/go-rvq/rvq/admin/presets"
	"github.com/go-rvq/rvq/admin/presets/actions"
	"github.com/go-rvq/rvq/admin/presets/gorm2op"
	"github.com/go-rvq/rvq/web"
	"github.com/go-rvq/rvq/web/multipartestutils"
	"github.com/go-rvq/rvq/x/i18n"
)

// A model that records WHEN it changed and by WHOM — the two fields the edit
// form's stale-record guard reads (see presets/record_stamp.go).
type StampedArticle struct {
	ID          uint
	Title       string
	UpdatedAt   time.Time
	UpdatedByID uint
}

// The same thing without the two fields: nothing to guard.
type UnstampedArticle struct {
	ID    uint
	Title string
}

type stampTestUser struct{ name, email string }

func (u stampTestUser) GetName() string  { return u.name }
func (u stampTestUser) GetEmail() string { return u.email }

type stampApp struct {
	handler *presets.Builder
	signer  *presets.HMACFormSigner
}

func newStampApp(t *testing.T, withUserFinder bool) *stampApp {
	t.Helper()

	db := TestDB
	if err := db.AutoMigrate(&StampedArticle{}); err != nil {
		t.Fatal(err)
	}
	db.Exec("DELETE FROM stamped_articles")

	// a fixed key, so the test can forge the stamps the server would produce
	signer := presets.NewHMACFormSigner([]byte("test key"))

	b := presets.New(i18n.New()).URIPrefix("/admin")
	b.DataOperator(gorm2op.DataOperator(db))
	b.SetFormSigner(signer)

	if withUserFinder {
		b.SetRecordUserFinder(func(id any, ctx *web.EventContext) (presets.RecordUser, error) {
			return stampTestUser{name: "Ana Souza", email: "ana@example.com"}, nil
		})
	}

	b.Model(&StampedArticle{}).URIName("articles")
	return &stampApp{handler: b, signer: signer}
}

func (a *stampApp) seed(t *testing.T) *StampedArticle {
	t.Helper()
	article := &StampedArticle{Title: "Original"}
	if err := TestDB.Create(article).Error; err != nil {
		t.Fatal(err)
	}
	return a.reload(t, article.ID)
}

func (a *stampApp) reload(t *testing.T, id uint) *StampedArticle {
	t.Helper()
	var got StampedArticle
	if err := TestDB.First(&got, id).Error; err != nil {
		t.Fatal(err)
	}
	return &got
}

// editForm renders the record's edit form and returns its body.
func (a *stampApp) editForm(t *testing.T, id uint) string {
	t.Helper()
	w := httptest.NewRecorder()
	r := multipartestutils.NewMultipartBuilder().
		PageURL("/admin/articles").
		EventFunc(actions.Edit).
		Query(presets.ParamID, itoa(id)).
		BuildEventFuncRequest()
	a.handler.ServeHTTP(w, r)
	return w.Body.String()
}

// update posts the edit form. An empty stamp means "do not send the field".
func (a *stampApp) update(t *testing.T, id uint, title, stamp string) string {
	t.Helper()
	b := multipartestutils.NewMultipartBuilder().
		PageURL("/admin/articles").
		EventFunc(actions.Update).
		Query(presets.ParamID, itoa(id)).
		AddField("Title", title)

	if stamp != "" {
		b = b.AddField(presets.RecordStampFormKey, stamp)
	}

	w := httptest.NewRecorder()
	a.handler.ServeHTTP(w, b.BuildEventFuncRequest())
	return w.Body.String()
}

func itoa(v uint) string {
	return strconv.FormatUint(uint64(v), 10)
}

// stampOf reads the stamp the rendered form carries. The body is the JSON of
// the event response, so its quotes come escaped.
var stampRe = regexp.MustCompile(`__formSign\\?":\s*\\?"([^"\\]+)`)

func stampOf(t *testing.T, body string) string {
	t.Helper()
	m := stampRe.FindStringSubmatch(body)
	if m == nil {
		t.Fatalf("the form carries no %s field", presets.RecordStampFormKey)
	}
	return m[1]
}

func TestRecordStampGuardsTheEditForm(t *testing.T) {
	app := newStampApp(t, false)
	article := app.seed(t)

	// 1. the form carries the stamp of the record it was rendered from
	stamp := stampOf(t, app.editForm(t, article.ID))
	value, err := app.signer.Unsign(stamp)
	if err != nil {
		t.Fatalf("the form's stamp is not signed by the app: %v", err)
	}
	if want := presets.RecordStampValue(article.UpdatedAt); value != want {
		t.Errorf("stamp = %q, want %q (the record's UpdatedAt)", value, want)
	}

	// 2. saving it back works — nothing changed underneath
	app.update(t, article.ID, "Edited", stamp)
	if got := app.reload(t, article.ID); got.Title != "Edited" {
		t.Errorf("title = %q, want %q", got.Title, "Edited")
	}

	// 3. the stamp the browser still holds is now stale: somebody else saved.
	//    The update must refuse it and keep the stored title.
	body := app.update(t, article.ID, "Written over", stamp)
	if got := app.reload(t, article.ID); got.Title != "Edited" {
		t.Errorf("a stale form overwrote the record: title = %q, want %q", got.Title, "Edited")
	}
	if !strings.Contains(body, "changed by someone else") {
		t.Errorf("the answer does not explain what happened:\n%s", firstLines(body, 3))
	}
	// with no author to name, it still says when it happened
	when := presets.Messages_en_US.FormatDateTime(app.reload(t, article.ID).UpdatedAt)
	if !strings.Contains(body, when) {
		t.Errorf("the answer does not say when the record changed (%s):\n%s", when, firstLines(body, 3))
	}
}

func TestRecordStampIsRequired(t *testing.T) {
	app := newStampApp(t, false)
	article := app.seed(t)

	// no stamp at all: a caller must not be able to skip the check by dropping
	// the field from the POST
	body := app.update(t, article.ID, "No stamp", "")
	if got := app.reload(t, article.ID); got.Title != "Original" {
		t.Errorf("an update with no stamp went through: title = %q", got.Title)
	}
	if !strings.Contains(body, "out of date") {
		t.Errorf("the answer does not explain what happened:\n%s", firstLines(body, 3))
	}
}

func TestRecordStampRejectsHandEditedValues(t *testing.T) {
	app := newStampApp(t, false)
	article := app.seed(t)

	stamp := stampOf(t, app.editForm(t, article.ID))

	// the value the record has now, with the signature of the older stamp: what
	// somebody would send after editing the field in the browser
	forged := presets.RecordStampValue(article.UpdatedAt.Add(time.Hour)) +
		stamp[strings.LastIndexByte(stamp, '.'):]

	app.update(t, article.ID, "Forged", forged)
	if got := app.reload(t, article.ID); got.Title != "Original" {
		t.Errorf("a hand-edited stamp went through: title = %q", got.Title)
	}

	// and a stamp signed with another key, as another deployment would produce
	other := presets.NewHMACFormSigner([]byte("another key")).
		Sign(presets.RecordStampValue(article.UpdatedAt))

	app.update(t, article.ID, "Other key", other)
	if got := app.reload(t, article.ID); got.Title != "Original" {
		t.Errorf("a stamp signed with another key went through: title = %q", got.Title)
	}
}

func TestRecordStampNamesTheAuthor(t *testing.T) {
	app := newStampApp(t, true)
	article := app.seed(t)

	stamp := stampOf(t, app.editForm(t, article.ID))

	// somebody else saves, recording themselves as the author
	TestDB.Model(&StampedArticle{}).Where("id = ?", article.ID).
		Updates(map[string]any{"title": "Theirs", "updated_by_id": 7, "updated_at": time.Now().Add(time.Second)})

	body := app.update(t, article.ID, "Mine", stamp)
	if !strings.Contains(body, "Ana Souza") || !strings.Contains(body, "ana@example.com") {
		t.Errorf("the message does not name who changed the record:\n%s", firstLines(body, 3))
	}

	// and WHEN they changed it — the stored UpdatedAt, the very value that did
	// not match the stamp
	when := presets.Messages_en_US.FormatDateTime(app.reload(t, article.ID).UpdatedAt)
	if !strings.Contains(body, when) {
		t.Errorf("the message does not say when the record changed (%s):\n%s", when, firstLines(body, 3))
	}
}

// A model with no UpdatedAt has nothing to compare: its form carries no stamp
// and its updates are not required to send one.
func TestRecordStampSkipsModelsWithoutUpdatedAt(t *testing.T) {
	if _, ok := presets.RecordUpdatedAt(&UnstampedArticle{}); ok {
		t.Fatal("this test needs a model with no UpdatedAt")
	}

	db := TestDB
	if err := db.AutoMigrate(&UnstampedArticle{}); err != nil {
		t.Fatal(err)
	}
	db.Exec("DELETE FROM unstamped_articles")

	product := &UnstampedArticle{Title: "Original"}
	if err := db.Create(product).Error; err != nil {
		t.Fatal(err)
	}

	b := presets.New(i18n.New()).URIPrefix("/admin")
	b.DataOperator(gorm2op.DataOperator(db))
	b.Model(&UnstampedArticle{}).URIName("unstamped-articles")

	w := httptest.NewRecorder()
	r := multipartestutils.NewMultipartBuilder().
		PageURL("/admin/unstamped-articles").
		EventFunc(actions.Edit).
		Query(presets.ParamID, itoa(product.ID)).
		BuildEventFuncRequest()
	b.ServeHTTP(w, r)

	if strings.Contains(w.Body.String(), presets.RecordStampFormKey) {
		t.Error("a model with no UpdatedAt must not carry a stamp")
	}

	// and the update goes through with no stamp at all
	w = httptest.NewRecorder()
	r = multipartestutils.NewMultipartBuilder().
		PageURL("/admin/unstamped-articles").
		EventFunc(actions.Update).
		Query(presets.ParamID, itoa(product.ID)).
		AddField("Title", "Edited").
		BuildEventFuncRequest()
	b.ServeHTTP(w, r)

	var got UnstampedArticle
	if err := db.First(&got, product.ID).Error; err != nil {
		t.Fatal(err)
	}
	if got.Title != "Edited" {
		t.Errorf("title = %q, want %q — an unguarded model must save normally", got.Title, "Edited")
	}
}

func firstLines(s string, n int) string {
	lines := strings.SplitN(s, "\n", n+1)
	if len(lines) > n {
		lines = lines[:n]
	}
	return strings.Join(lines, "\n")
}
