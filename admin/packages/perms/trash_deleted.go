package perms

import (
	"context"
	"fmt"
	"net/http"
	"reflect"
	"time"

	h "github.com/go-rvq/htmlgo"
	"github.com/go-rvq/rvq/admin/origin"
	"github.com/go-rvq/rvq/admin/origin/originui"
	"github.com/go-rvq/rvq/admin/packages/user"
	"github.com/go-rvq/rvq/admin/presets"
	"github.com/go-rvq/rvq/admin/presets/actions"
	"github.com/go-rvq/rvq/admin/softdelete"
	"github.com/go-rvq/rvq/web"
	v "github.com/go-rvq/rvq/x/ui/vuetify"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

// UsersTable is the table of the users the deleted_by_id of a trash
// references, and UserNameColumn the column of their names.
var (
	UsersTable     = "users"
	UserNameColumn = "name"
)

// ActionDeletedOrigin is the action of a deleted record — of its detail, in
// the admin's dialog — that shows where it was deleted from: the address, the
// browser, the place and, its coordinates known, the place on a map. The
// trash's "Deleted from" column opens it.
const ActionDeletedOrigin = "deleted_origin"

// hasDeletion is whether obj's model keeps who deleted a record and from where
// (softdelete.Deletion).
func hasDeletion(obj any) bool {
	t := reflect.Indirect(reflect.ValueOf(obj)).Type()
	_, ok := t.FieldByName("DeletedOrigin")
	_, ok2 := t.FieldByName("DeletedByID")
	return ok && ok2
}

// MigrateDeleted adds to table — of obj's model, which keeps who deleted a
// record (softdelete.Deletion) — the foreign key of its deleted_by_id to the
// users (UsersTable): nil when the user is deleted. Nothing before the table
// is.
func MigrateDeleted(db *gorm.DB, table string, obj any) error {
	db = db.Session(&gorm.Session{NewDB: true})
	m := db.Migrator()
	if !hasDeletion(obj) || !m.HasTable(table) || !m.HasColumn(table, "deleted_by_id") || !m.HasTable(UsersTable) {
		return nil
	}
	name := "fk_" + table + "_deleted_by"
	if m.HasConstraint(table, name) {
		return nil
	}
	if err := db.Exec(fmt.Sprintf(`ALTER TABLE %q ADD CONSTRAINT %q FOREIGN KEY (deleted_by_id) REFERENCES %q (id) ON DELETE SET NULL`,
		table, name, UsersTable)).Error; err != nil {
		return fmt.Errorf("the deleted_by_id of %s: %w", table, err)
	}
	return nil
}

// DeletedBy are the values of the columns of who deleted a record — the
// user of r — and from where (softdelete.Columns), to set with its
// deleted_at.
func DeletedBy(r *http.Request) map[string]any {
	if r == nil {
		return softdelete.Columns(nil, origin.Origin{})
	}
	var by *uuid.UUID
	if u := user.GetCurrentUser(r); u != nil {
		id := u.GetID()
		by = &id
	}
	return softdelete.Columns(by, origin.Of(r))
}

// deletionOf is the deletion obj keeps; nil when it keeps none.
func deletionOf(obj any) *softdelete.Deletion {
	v := reflect.Indirect(reflect.ValueOf(obj))
	f := v.FieldByName("Deletion")
	if f.IsValid() {
		if d, ok := f.Addr().Interface().(*softdelete.Deletion); ok {
			return d
		}
	}
	return nil
}

// deletedAtOf is when obj was deleted; zero while it is not (or its model has
// no gorm.DeletedAt).
func deletedAtOf(obj any) time.Time {
	f := reflect.Indirect(reflect.ValueOf(obj)).FieldByName("DeletedAt")
	if !f.IsValid() {
		return time.Time{}
	}
	if d, ok := f.Interface().(gorm.DeletedAt); ok && d.Valid {
		return d.Time
	}
	return time.Time{}
}

// userName is the name of the user id; "" when not found.
func userName(db *gorm.DB, id *uuid.UUID) (name string) {
	if id != nil {
		db.Session(&gorm.Session{NewDB: true}).Table(UsersTable).
			Select(UserNameColumn).Where("id = ?", *id).Limit(1).Scan(&name)
	}
	return
}

// installDeletedColumns adds to the trash of mb — in its tab only — the
// columns of when, by whom and from where each record was deleted, and the
// action that shows where from with a map.
func installDeletedColumns(mb *presets.ModelBuilder, db *gorm.DB, table string) {
	inTrash := func(ctx *presets.FieldContext) bool {
		return ctx.Mode.IsTrash()
	}
	label := func(f func(*Messages) string) func(context.Context) string {
		return func(c context.Context) string { return f(msgs(c)) }
	}
	lb := mb.Listing()
	lb.Field("DeletedAt").
		SetI18nLabel(label(func(m *Messages) string { return m.DeletedAt })).
		SetEnabled(inTrash).
		ComponentFunc(func(field *presets.FieldContext, ctx *web.EventContext) h.HTMLComponent {
			var s string
			if d := deletionOf(field.Obj); d != nil && d.DeletedAt.Valid {
				s = d.DeletedAt.Time.Local().Format(time.DateTime)
			}
			return h.Td(h.Text(s)).Style("white-space:nowrap")
		})
	installDeletedNotice(mb, db)
	if !hasDeletion(mb.NewModel()) {
		lb.AppendTrailingFields("DeletedAt")
		return
	}
	lb.Field("DeletedByID").
		SetI18nLabel(label(func(m *Messages) string { return m.DeletedBy })).
		SetEnabled(inTrash).
		ComponentFunc(func(field *presets.FieldContext, ctx *web.EventContext) h.HTMLComponent {
			var name string
			if d := deletionOf(field.Obj); d != nil {
				name = userName(db, d.DeletedByID)
			}
			return h.Td(h.Text(name))
		})
	lb.Field("DeletedOrigin").
		SetI18nLabel(label(func(m *Messages) string { return m.DeletedOrigin })).
		SetEnabled(inTrash).
		ComponentFunc(func(field *presets.FieldContext, ctx *web.EventContext) h.HTMLComponent {
			d := deletionOf(field.Obj)
			if d == nil {
				return h.Td()
			}
			return h.Td(deletedOriginLink(mb, field.Obj, d, ctx))
		})
	lb.AppendTrailingFields("DeletedAt", "DeletedByID", "DeletedOrigin")

	mb.Detailing().Action(ActionDeletedOrigin).
		SetI18nLabel(label(func(m *Messages) string { return m.DeletedOrigin })).
		Icon("mdi-map-marker-remove-outline").
		SetTrash(presets.TrashOnly).
		// of a deleted record only: the trash's column opens it
		SetEnabledObj(func(obj any, _ string, _ *web.EventContext) (bool, error) {
			d := deletionOf(obj)
			return d != nil && d.DeletedAt.Valid, nil
		}).
		ComponentFunc(func(id string, ctx *web.EventContext) (h.HTMLComponent, error) {
			if mb.Permissioner().ReqListDo(ctx.R, PermTrash).Denied() {
				return nil, presets.ErrActionNotAllowed
			}
			mid, err := mb.ParseRecordID(id)
			if err != nil {
				return nil, err
			}
			obj := mb.NewModel()
			mid.SetTo(obj)
			// deleted: out of the scoped queries
			if err = db.Session(&gorm.Session{NewDB: true}).Unscoped().Table(table).Where(obj).First(obj).Error; err != nil {
				return nil, err
			}
			m := msgs(ctx.Context())
			return originui.Body(originui.Labels{
				IP: m.OriginIP, Browser: m.OriginBrowser, Place: m.OriginPlace,
				Coordinates: m.OriginCoordinates, NoMap: m.OriginNoMap,
			}, deletionOf(obj).DeletedOrigin), nil
		})
}

// deletedOriginLink is where obj — d its deletion — was deleted from, on one
// line, opening the deleted_origin action (of its detail, in the admin's
// dialog).
func deletedOriginLink(mb *presets.ModelBuilder, obj any, d *softdelete.Deletion, ctx *web.EventContext) h.HTMLComponent {
	id := mb.MustRecordID(obj).String()
	open := web.Plaid().
		EventFunc(actions.Action).
		Query(presets.ParamID, id).
		Query(presets.ParamAction, ActionDeletedOrigin).
		Query(presets.ParamOverlay, actions.Dialog).
		URL(mb.Info().ListingHrefCtx(ctx)).
		Go()
	return originui.Link(d.DeletedOrigin, open, "data-deleted-origin", id)
}

// installDeletedNotice warns, at the top of the detail of a deleted record,
// that it is deleted — and when, by whom and from where, when its model keeps
// them (softdelete.Deletion).
func installDeletedNotice(mb *presets.ModelBuilder, db *gorm.DB) {
	mb.Detailing().AppendNoticeFunc(func(obj any, ctx *web.EventContext) h.HTMLComponent {
		at := deletedAtOf(obj)
		if at.IsZero() {
			return nil
		}
		m := msgs(ctx.Context())
		row := func(label string, value h.HTMLComponent) h.HTMLComponent {
			return h.Div(h.Strong(label+": "), value)
		}
		body := h.HTMLComponents{row(m.DeletedAt, h.Text(at.Local().Format(time.DateTime)))}
		if d := deletionOf(obj); d != nil {
			if name := userName(db, d.DeletedByID); name != "" {
				body = append(body, row(m.DeletedBy, h.Text(name)))
			}
			if d.DeletedOrigin.String() != "" {
				body = append(body, row(m.DeletedOrigin, deletedOriginLink(mb, obj, d, ctx)))
			}
		}
		return v.VAlert(body...).
			Type("warning").
			Variant(v.VariantTonal).
			Title(m.DeletedNotice).
			Icon("mdi-delete-alert-outline").
			Class("mb-4").
			Attr("data-deleted-notice", "true")
	})
}
