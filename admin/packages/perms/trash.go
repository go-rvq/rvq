package perms

import (
	"context"

	"github.com/go-rvq/rvq/admin/model"
	"github.com/go-rvq/rvq/admin/presets"
	"github.com/go-rvq/rvq/admin/presets/gorm2op"
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
// PermTrash ('trash') listing permission: an "All"/"Trash" filter tab pair, the
// unscoped trash query and a restore bulk action are all only available to
// subjects allowed the trash verb (backend enforced). Deleted records use
// gorm's soft delete, so they are hidden from the default listing.
func SetupTrash(mb *presets.ModelBuilder, db *gorm.DB, opts TrashOptions) {
	registerMessages(mb.Builder())
	lb := mb.Listing()

	lb.FilterTabsFunc(func(ctx *web.EventContext) []*presets.FilterTab {
		m := msgs(ctx.Context())
		tabs := []*presets.FilterTab{{Label: m.TabAll, Default: true}}
		if mb.Permissioner().ReqListDo(ctx.R, PermTrash).Allowed() {
			tabs = append(tabs, &presets.FilterTab{ID: FilterTabTrash, Label: m.TabTrash})
		}
		return tabs
	})

	mb.UpdateDataOperator(func(do presets.DataOperator) presets.DataOperator {
		return do.(*gorm2op.DataOperatorBuilder).WrapPrepare(func(old gorm2op.Preparer) gorm2op.Preparer {
			return func(db *gorm.DB, mode gorm2op.Mode, obj interface{}, id model.ID, params *presets.SearchParams, ctx *web.EventContext) *gorm.DB {
				db = old(db, mode, obj, id, params, ctx)
				if mode == gorm2op.Search && params != nil &&
					params.Query.Get(presets.ActiveFilterTabQueryKey) == FilterTabTrash &&
					mb.Permissioner().ReqListDo(ctx.R, PermTrash).Allowed() {
					db = db.Unscoped().Where(opts.TableName + ".deleted_at IS NOT NULL")
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
				Update("deleted_at", nil).Error; err != nil {
				return err
			}
			if opts.OnRestore != nil {
				return opts.OnRestore(ctx, selectedIds)
			}
			return nil
		})
}
