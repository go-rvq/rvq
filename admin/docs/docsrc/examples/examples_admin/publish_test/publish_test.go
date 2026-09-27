package publish_test

// The publish example driven as the admin drives it — the events its buttons
// fire — and checked by what it leaves in the database and what it tells the
// page to do next. The scenarios are the ones the recorded flows these tests
// replace covered: new, duplicate, publish and unpublish, schedule, the
// version list dialog, and deleting versions.

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"regexp"
	"strings"
	"testing"
	"time"

	"github.com/go-rvq/rvq/admin/docs/docsrc/examples/examples_admin"
	"github.com/go-rvq/rvq/admin/presets"
	"github.com/go-rvq/rvq/admin/presets/actions"
	"github.com/go-rvq/rvq/admin/publish"
	"github.com/go-rvq/rvq/web/multipartestutils"
)

const (
	productsURL = "/samples/publish-example/with-publish-products"
	versionsURL = productsURL + "-version-list-dialog"

	// the events of the version list dialog (publish/event.go), fired from its
	// buttons
	eventRenameVersion = "publish_eventRenameVersion"
	eventDeleteVersion = "publish_eventDeleteVersion"
	eventSelectVersion = "publish_eventSelectVersion"
	eventSchedule      = "publish_eventSchedulePublish"

	varCurrentDisplayID = "vars.publish_VarCurrentDisplayID"
)

var msgr = publish.Messages_en_US

// response is an event's answer, with everything it shows as text.
type response struct {
	multipartestutils.TestEventResponse
	Code int
}

// Shows reports whether the page is told s: in a script, a portal or the body.
func (r *response) Shows(s string) bool {
	if strings.Contains(r.RunScript, s) || strings.Contains(r.Body, s) {
		return true
	}
	for _, p := range r.UpdatePortals {
		if strings.Contains(p.Body, s) {
			return true
		}
	}
	return false
}

type event struct {
	url, name string
	query     [][2]string
	form      [][2]string
}

func (e event) Q(k, v string) event { e.query = append(e.query, [2]string{k, v}); return e }
func (e event) F(k, v string) event { e.form = append(e.form, [2]string{k, v}); return e }

func (e event) fire(t *testing.T) *response {
	t.Helper()
	b := multipartestutils.NewMultipartBuilder().PageURL(e.url).EventFunc(e.name)
	for _, q := range e.query {
		b = b.Query(q[0], q[1])
	}
	for _, f := range e.form {
		b = b.AddField(f[0], f[1])
	}
	w := httptest.NewRecorder()
	PresetsBuilder.ServeHTTP(w, b.BuildEventFuncRequest())
	r := &response{Code: w.Code}
	if w.Code != http.StatusOK {
		t.Fatalf("%s %s: %d\n%.500s", e.url, e.name, w.Code, w.Body.String())
	}
	if err := json.Unmarshal(w.Body.Bytes(), &r.TestEventResponse); err != nil {
		t.Fatalf("%s %s: %v\n%.500s", e.url, e.name, err, w.Body.String())
	}
	return r
}

func on(url, name string) event { return event{url: url, name: name} }

// reset empties the example's table: every test starts from nothing.
func reset(t *testing.T) {
	t.Helper()
	if err := DB.Exec(`TRUNCATE with_publish_products`).Error; err != nil {
		t.Fatal(err)
	}
}

type product = examples_admin.WithPublishProduct

func slug(p *product) string { return p.PrimarySlug() }

// create makes a product through the admin's New form and returns it.
func create(t *testing.T, name string, price int) *product {
	t.Helper()
	r := on(productsURL, actions.Create).F("Name", name).F("Price", fmt.Sprint(price)).fire(t)
	var p product
	if err := DB.Where("name = ?", name).Order("created_at DESC").First(&p).Error; err != nil {
		t.Fatalf("create %q made no record: %v\n%+v", name, err, r.TestEventResponse)
	}
	return &p
}

// versions are the product's versions not deleted, oldest first.
func versions(t *testing.T, id uint) []*product {
	t.Helper()
	var ps []*product
	if err := DB.Where("id = ?", id).Order("version").Find(&ps).Error; err != nil {
		t.Fatal(err)
	}
	return ps
}

func reload(t *testing.T, p *product) *product {
	t.Helper()
	var got product
	if err := DB.Unscoped().Where("id = ? AND version = ?", p.ID, p.Version.Version).First(&got).Error; err != nil {
		t.Fatal(err)
	}
	return &got
}

// duplicate makes a new version from p, as the version bar's Duplicate does.
func duplicate(t *testing.T, p *product) *product {
	t.Helper()
	before := versions(t, p.ID)
	r := on(productsURL, publish.EventDuplicateVersion).Q(presets.ParamID, slug(p)).fire(t)
	after := versions(t, p.ID)
	if len(after) != len(before)+1 {
		t.Fatalf("duplicate %s: %d versions, want %d\n%+v", slug(p), len(after), len(before)+1, r.TestEventResponse)
	}
	return after[len(after)-1]
}

func today() string { return DB.NowFunc().Format("2006-01-02") }

func TestNew(t *testing.T) {
	reset(t)
	p := create(t, "zz-new", 10)

	if want := today() + "-v01"; p.Version.Version != want {
		t.Errorf("version %q, want %q", p.Version.Version, want)
	}
	if p.Version.VersionName != p.Version.Version || p.Version.ParentVersion != "" {
		t.Errorf("a first version is named after itself and has no parent: %+v", p.Version)
	}
	if p.Status.Status != publish.StatusDraft || p.Price != 10 {
		t.Errorf("new product: status %q price %d", p.Status.Status, p.Price)
	}
}

func TestDuplicate(t *testing.T) {
	reset(t)
	v1 := create(t, "zz-dup", 20)

	// a published version is copied as a draft, out of any schedule
	on(productsURL, publish.EventPublish).Q(presets.ParamID, slug(v1)).fire(t)
	v2 := duplicate(t, reload(t, v1))
	if want := today() + "-v02"; v2.Version.Version != want || v2.Version.VersionName != want {
		t.Errorf("second version %+v, want %s", v2.Version, want)
	}
	if v2.Version.ParentVersion != v1.Version.Version {
		t.Errorf("parent %q, want %q", v2.Version.ParentVersion, v1.Version.Version)
	}
	if v2.Status.Status != publish.StatusDraft || v2.Status.OnlineUrl != "" || v2.Schedule.ActualStartAt != nil {
		t.Errorf("a copy is a draft, offline and unscheduled: %+v %+v", v2.Status, v2.Schedule)
	}
	if v2.Name != v1.Name || v2.Price != v1.Price {
		t.Errorf("a copy keeps the content: %q %d", v2.Name, v2.Price)
	}

	// and again, from the copy
	v3 := duplicate(t, v2)
	if v3.Version.ParentVersion != v2.Version.Version || v3.Version.Version != today()+"-v03" {
		t.Errorf("third version %+v", v3.Version)
	}
}

func TestPublishAndUnpublish(t *testing.T) {
	reset(t)
	v1 := create(t, "zz-pub", 30)

	on(productsURL, publish.EventPublish).Q(presets.ParamID, slug(v1)).fire(t)
	v1 = reload(t, v1)
	if v1.Status.Status != publish.StatusOnline || v1.Schedule.ActualStartAt == nil {
		t.Fatalf("published: %+v %+v", v1.Status, v1.Schedule)
	}

	// publishing another version takes the first offline: one is online
	v2 := duplicate(t, v1)
	on(productsURL, publish.EventPublish).Q(presets.ParamID, slug(v2)).fire(t)
	if s := reload(t, v2).Status.Status; s != publish.StatusOnline {
		t.Errorf("v2 %q, want online", s)
	}
	if s := reload(t, v1).Status.Status; s != publish.StatusOffline {
		t.Errorf("v1 %q, want offline once v2 is online", s)
	}

	on(productsURL, publish.EventUnpublish).Q(presets.ParamID, slug(v2)).fire(t)
	v2 = reload(t, v2)
	if v2.Status.Status != publish.StatusOffline || v2.Schedule.ActualEndAt == nil {
		t.Errorf("unpublished: %+v %+v", v2.Status, v2.Schedule)
	}
}

func TestSchedule(t *testing.T) {
	reset(t)
	p := create(t, "zz-sched", 40)
	at := func(d time.Duration) string {
		return publish.ScheduleTimeString(ptr(DB.NowFunc().Add(d)))
	}
	schedule := func(start, end string) *response {
		return on(productsURL, eventSchedule).Q(presets.ParamID, slug(p)).
			F("ScheduledStartAt", start).F("ScheduledEndAt", end).fire(t)
	}

	// now < start < end: saved
	start, end := at(time.Hour), at(2*time.Hour)
	schedule(start, end)
	got := reload(t, p)
	if publish.ScheduleTimeString(got.Schedule.ScheduledStartAt) != start || publish.ScheduleTimeString(got.Schedule.ScheduledEndAt) != end {
		t.Fatalf("schedule %v – %v, want %s – %s", got.Schedule.ScheduledStartAt, got.Schedule.ScheduledEndAt, start, end)
	}

	// refused, and the saved schedule kept
	for _, c := range []struct {
		name, start, end, message string
	}{
		{"start < end < now", at(-2 * time.Hour), at(-time.Hour), msgr.ScheduledStartAtShouldLaterThanNow},
		{"now < end < start", at(2 * time.Hour), at(time.Hour), msgr.ScheduledEndAtShouldLaterThanStartAt},
		{"end without start, on a draft", "", at(time.Hour), msgr.ScheduledStartAtShouldNotEmpty},
	} {
		t.Run(c.name, func(t *testing.T) {
			if r := schedule(c.start, c.end); !r.Shows(c.message) {
				t.Errorf("not told %q: %+v", c.message, r.TestEventResponse)
			}
			if s := publish.ScheduleTimeString(reload(t, p).Schedule.ScheduledStartAt); s != start {
				t.Errorf("a refused schedule changed the saved one: %s", s)
			}
		})
	}

	// both cleared: no schedule
	schedule("", "")
	if got := reload(t, p); got.Schedule.ScheduledStartAt != nil || got.Schedule.ScheduledEndAt != nil {
		t.Errorf("cleared schedule: %+v", got.Schedule)
	}
}

func ptr[T any](v T) *T { return &v }

// versionNames are the version names the version list dialog shows, in its
// order.
var versionCell = regexp.MustCompile(`<v-radio[^>]*></v-radio>\s*([^<\s][^<]*?)\s*<`)

func versionNames(r *response) []string {
	body := r.Body
	for _, p := range r.UpdatePortals {
		body += p.Body
	}
	var names []string
	for _, m := range versionCell.FindAllStringSubmatch(body, -1) {
		names = append(names, m[1])
	}
	return names
}

func TestVersionDialog(t *testing.T) {
	reset(t)
	v1 := create(t, "zz-dialog", 50)
	v2 := duplicate(t, v1)
	v3 := duplicate(t, v2)
	other := create(t, "zz-other", 1)

	// rename v2
	on(versionsURL, eventRenameVersion).Q(presets.ParamID, slug(v2)).F("VersionName", "zz-named").fire(t)
	if n := reload(t, v2).Version.VersionName; n != "zz-named" {
		t.Fatalf("renamed to %q", n)
	}

	list := func(query ...[2]string) []string {
		e := on(versionsURL, actions.OpenListingDialog).Q("select_id", slug(v3))
		for _, q := range query {
			e = e.Q(q[0], q[1])
		}
		return versionNames(e.fire(t))
	}

	// the record's versions only, newest first
	all := list()
	want := []string{v3.Version.VersionName, "zz-named", v1.Version.VersionName}
	if strings.Join(all, ",") != strings.Join(want, ",") {
		t.Errorf("versions %v, want %v (and none of %s's)", all, want, other.Name)
	}

	// the named versions tab (a tab's query goes with the "f_" of the listing
	// filters), and a search by name
	if got := list([2]string{"f_named_versions", "1"}); strings.Join(got, ",") != "zz-named" {
		t.Errorf("named versions %v", got)
	}
	if got := list([2]string{"keyword", "zz-named"}); strings.Join(got, ",") != "zz-named" {
		t.Errorf("search %v", got)
	}

	// choosing a version shows it
	r := on(productsURL, eventSelectVersion).Q("select_id", slug(v1)).fire(t)
	if !r.Shows(actions.Detailing) || !r.Shows(slug(v1)) {
		t.Errorf("select did not show %s: %s", slug(v1), r.RunScript)
	}
}

// deleteVersion deletes p from the version list dialog, with current the
// version the page shows.
func deleteVersion(t *testing.T, p, current *product) *response {
	t.Helper()
	return on(versionsURL, eventDeleteVersion).
		Q(presets.ParamID, slug(p)).
		Q("current_display_id", slug(current)).
		Q(presets.ParamListingQueries, "select_id="+slug(current)).
		fire(t)
}

func TestDeleteVersion(t *testing.T) {
	reset(t)
	v1 := create(t, "zz-del", 60)
	v2 := duplicate(t, v1)
	v3 := duplicate(t, v2)

	// not the one shown: gone, and the page keeps what it shows
	r := deleteVersion(t, v1, v3)
	if reload(t, v1).DeletedAt.Valid == false {
		t.Fatal("v1 was not deleted")
	}
	if strings.Contains(r.RunScript, varCurrentDisplayID) {
		t.Errorf("deleting another version changed the one shown: %s", r.RunScript)
	}

	// the one shown: the page moves to the next older version
	r = deleteVersion(t, v3, v3)
	if !r.Shows(fmt.Sprintf("%s = %q", varCurrentDisplayID, slug(v2))) {
		t.Errorf("after deleting the shown v3, v2 is not shown: %s", r.RunScript)
	}

	// the oldest left, shown: no older one, so the newest
	v4 := duplicate(t, v2)
	r = deleteVersion(t, v2, v2)
	if !r.Shows(fmt.Sprintf("%s = %q", varCurrentDisplayID, slug(v4))) {
		t.Errorf("after deleting the oldest, the newest v4 is not shown: %s", r.RunScript)
	}

	// the last one: back to the listing
	r = deleteVersion(t, v4, v4)
	if len(versions(t, v1.ID)) != 0 {
		t.Fatal("versions left")
	}
	if r.PushState == nil || !strings.HasSuffix(r.PushState.MyURL, "/with-publish-products") {
		t.Errorf("deleting the last version must go back to the listing: %+v", r.PushState)
	}
}
