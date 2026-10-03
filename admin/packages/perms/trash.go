package perms

import (
	"context"
	"net/http"
	"reflect"

	"github.com/go-rvq/rvq/admin/model"
	"github.com/go-rvq/rvq/admin/presets"
	"github.com/go-rvq/rvq/admin/presets/gorm2op"
	"github.com/go-rvq/rvq/admin/softdelete"
	"github.com/go-rvq/rvq/web"
	"github.com/go-rvq/rvq/x/perm"
	"gorm.io/gorm"
)

// FilterTabTrash is the listing filter tab id that shows soft-deleted records.
const FilterTabTrash = "trash"

// TrashOptions configures SetupTrash.
type TrashOptions struct {
	// TableName of the soft-deletable model (required for the trash query and
	// the restore update).
	TableName string
	// OnRestore, when set, is called with the restored ids after a successful
	// restore (e.g. to write an audit record).
	OnRestore func(ctx *web.EventContext, ids []string) error
}

// SetupTrash turns a soft-delete listing into a trash-enabled one, gated by the
// PermTrash ('trash') listing permission: a "Trash" filter tab after the
// model's own tabs — after an "All" one, the default, when it has none —, the
// unscoped trash query and a restore bulk action are all only available to
// subjects allowed the trash verb (backend enforced). Deleted records use
// gorm's soft delete, so they are hidden from the default listing. The tab is
// added when the listing renders (presets.FilterTabsWrapper): the model's
// FilterTabsFunc, set before or after, keeps it.
func SetupTrash(mb *presets.ModelBuilder, db *gorm.DB, opts TrashOptions) {
	registerMessages(mb.Builder())
	lb := mb.Listing()

	lb.AppendFilterTabsWrapper(func(ctx *web.EventContext, tabs []*presets.FilterTab) []*presets.FilterTab {
		if mb.Permissioner().ReqListDo(ctx.R, PermTrash).Denied() {
			return tabs
		}
		m := msgs(ctx.Context())
		if len(tabs) == 0 {
			tabs = append(tabs, &presets.FilterTab{Label: m.TabAll, Default: true})
		}
		// on the right, its icon alone (its name a tooltip): apart from the
		// tabs of the records there are
		return append(tabs, &presets.FilterTab{ID: FilterTabTrash, Label: m.TabTrash,
			Icon: "mdi-trash-can-outline", End: true})
	})

	// who deletes a record, and from where (softdelete.Deletion)
	deletion := hasDeletion(mb.NewModel())
	if err := MigrateDeleted(db, opts.TableName, mb.NewModel()); err != nil {
		panic(err)
	}
	installDeletedColumns(mb, db, opts.TableName)

	mb.UpdateDataOperator(func(do presets.DataOperator) presets.DataOperator {
		return do.(*gorm2op.DataOperatorBuilder).WithDeleteCallbacks(func(cb *gorm2op.Callbacks[*gorm2op.DataOperatorBuilder]) {
			cb.Post(func(state *gorm2op.CallbackState) error {
				if !deletion {
					return nil
				}
				var r *http.Request
				if state.Ctx != nil {
					r = state.Ctx.R
				}
				return state.DB.Session(&gorm.Session{NewDB: true}).Unscoped().
					Model(state.Obj).UpdateColumns(DeletedBy(r)).Error
			})
		}).WrapPrepare(func(old gorm2op.Preparer) gorm2op.Preparer {
			return func(db *gorm.DB, mode gorm2op.Mode, obj interface{}, id model.ID, params *presets.SearchParams, ctx *web.EventContext) *gorm.DB {
				db = old(db, mode, obj, id, params, ctx)
				if mode == gorm2op.Search && params != nil &&
					params.Query.Get(presets.ActiveFilterTabQueryKey) == FilterTabTrash &&
					mb.Permissioner().ReqListDo(ctx.R, PermTrash).Allowed() {
					db = db.Unscoped().Where(opts.TableName + ".deleted_at IS NOT NULL")
				}
				// a record of the trash opens — its detail, its history, as
				// the parent of its children —: fetched deleted or not, by who
				// may see the trash
				if (mode.Has(gorm2op.Fetch) || mode.Has(gorm2op.FetchTitle)) && ctx != nil && ctx.R != nil &&
					mb.Permissioner().ReqListDo(ctx.R, PermTrash).Allowed() {
					db = db.Unscoped()
				}
				return db
			}
		})
	})

	lb.BulkAction("restore").
		SetI18nLabel(func(ctx context.Context) string { return msgs(ctx).Restore }).
		Icon("mdi-restore").
		UpdateFunc(func(selectedIds []string, ctx *web.EventContext, r *web.EventResponse) error {
			if mb.Permissioner().ReqListDo(ctx.R, PermTrash).Denied() {
				return perm.PermissionDenied
			}
			if err := db.Session(&gorm.Session{}).Unscoped().
				Table(opts.TableName).
				Where("id IN ?", selectedIds).
				UpdateColumns(restored(deletion)).Error; err != nil {
				return err
			}
			if opts.OnRestore != nil {
				return opts.OnRestore(ctx, selectedIds)
			}
			return nil
		})
}

// AutoTrash sets up the trash (SetupTrash) of every model of b made from now
// on whose records are soft-deleted — a gorm.DeletedAt field —, listed by
// the gorm data operator, not a singleton: the default of a soft-delete
// listing. Before the models.
func AutoTrash(b *presets.Builder, db *gorm.DB) {
	b.ModelConfigurators.Append(presets.ModelConfiguratorFunc(func(mb *presets.ModelBuilder) {
		if mb.GetSingleton() {
			return
		}
		if _, ok := mb.CurrentDataOperator().(*gorm2op.DataOperatorBuilder); !ok {
			return
		}
		if table, ok := softDeleteTable(db, mb.NewModel()); ok {
			SetupTrash(mb, db, TrashOptions{TableName: table})
		}
	}))
}

// softDeleteTable is the table of obj's model when its records are
// soft-deleted: a gorm.DeletedAt field.
func softDeleteTable(db *gorm.DB, obj any) (string, bool) {
	stmt := &gorm.Statement{DB: db}
	if err := stmt.Parse(obj); err != nil || stmt.Schema == nil {
		return "", false
	}
	for _, f := range stmt.Schema.Fields {
		if f.FieldType == reflect.TypeOf(gorm.DeletedAt{}) && f.DBName != "" {
			return stmt.Schema.Table, true
		}
	}
	return "", false
}

// restored are the columns of records restored: not deleted — by nobody, from
// nowhere, when the model keeps who deleted them.
func restored(deletion bool) map[string]any {
	if deletion {
		return softdelete.Restored()
	}
	return map[string]any{"deleted_at": nil}
}
