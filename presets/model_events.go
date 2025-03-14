package presets

import (
	"github.com/qor5/admin/v3/presets/actions"
	"github.com/qor5/web/v3"
)

func (mb *ModelBuilder) registerDefaultEventFuncs() {
	mb.RegisterEventFunc(actions.New, mb.editing.formNew)
	mb.RegisterEventFunc(actions.Edit, mb.editing.formEdit)
	mb.RegisterEventFunc(actions.DeleteConfirmation, mb.listing.deleteConfirmation)
	mb.RegisterEventFunc(actions.Update, mb.editing.defaultUpdate)
	mb.RegisterEventFunc(actions.Create, mb.editing.defaultCreate)
	mb.RegisterEventFunc(actions.DoDelete, mb.listing.doDelete)
	mb.RegisterEventFunc(actions.DoBulkAction, mb.listing.doBulkAction)
	mb.RegisterEventFunc(actions.DoListingAction, mb.listing.doListingAction)
	mb.RegisterEventFunc(actions.OpenBulkActionDialog, mb.listing.openBulkActionDialog)
	mb.RegisterEventFunc(actions.OpenActionDialog, mb.listing.openActionDialog)

	mb.RegisterEventFunc(actions.Action, mb.detailing.formAction)
	mb.RegisterEventFunc(actions.DoAction, mb.detailing.doAction)
	mb.RegisterEventFunc(actions.Detailing, mb.detailing.detailingEvent)
	// mb.RegisterEventFunc(actions.DetailingContent, mb.detailing.detailingContent)
	mb.RegisterEventFunc(actions.DoSaveDetailingField, mb.detailing.SaveDetailField)
	mb.RegisterEventFunc(actions.DoEditDetailingField, mb.detailing.EditDetailField)
	mb.RegisterEventFunc(actions.DoEditDetailingListField, mb.detailing.EditDetailListField)
	mb.RegisterEventFunc(actions.DoSaveDetailingListField, mb.detailing.SaveDetailListField)
	mb.RegisterEventFunc(actions.DoDeleteDetailingListField, mb.detailing.DeleteDetailListField)
	mb.RegisterEventFunc(actions.DoCreateDetailingListField, mb.detailing.CreateDetailListField)

	mb.RegisterEventFunc(actions.ListData, mb.listing.records)
	mb.RegisterEventFunc(actions.ReloadList, mb.listing.reloadList)
	mb.RegisterEventFunc(actions.OpenListingDialog, mb.listing.openListingDialog)
	mb.RegisterEventFunc(actions.OpenListingDialogForSelection, mb.listing.openListingDialogForSelection)
	mb.RegisterEventFunc(actions.UpdateListingDialog, mb.listing.updateListingDialog)

	// list editor
	mb.RegisterEventFunc(actions.AddRowEvent, addListItemRow(mb))
	mb.RegisterEventFunc(actions.RemoveRowEvent, removeListItemRow(mb))
	mb.RegisterEventFunc(actions.SortEvent, sortListItems(mb))
}

func (mb *ModelBuilder) RegisterEventFunc(id string, f web.EventFunc) string {
	return mb.RegisterEventHandler(id, web.EventFunc(func(ctx *web.EventContext) (r web.EventResponse, err error) {
		WithModel(ctx, mb)
		if r, err = f(ctx); err != nil {
			r.UpdatePortal(FlashPortalName, RenderFlash(err, ""))
			err = nil
		}
		return
	}))
}

func (mb *ModelBuilder) RegisterEventHandler(eventFuncId string, ef web.EventHandler) (key string) {
	return mb.EventsHub.RegisterEventHandler(eventFuncId, ef)
}
