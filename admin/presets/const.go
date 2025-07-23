package presets

import "net/http"

const (
	PermModule = "presets"
	PermList   = "presets:list"
	PermGet    = "presets:get"
	PermCreate = "presets:create"
	PermUpdate = "presets:update"
	PermDelete = "presets:delete"

	PermActions         = "action"
	PermDoListingAction = "do_listing_action"
	PermBulkActions     = "bulk_action"
)

var PermRead = []string{PermList, PermGet}

// params
const (
	ParamID                        = "id"
	ParamSelectedID                = "selected_id"
	ParamAction                    = "action"
	ParamOverlay                   = "overlay"
	ParamOverlayUpdateID           = "overlay_update_id"
	ParamBulkActionName            = "bulk_action"
	ParamSelectedIds               = "selected_ids"
	ParamListingQueries            = "presets_listing_queries"
	ParamAfterDeleteEvent          = "presets_after_delete_event"
	ParamPortalID                  = "portal_id"
	ParamTargetPortal              = "target_portal"
	ParamEditFormUnscoped          = "presets_edit_form_unscoped"
	ParamPostChangeCallback        = "presets_post_change_callback"
	ParamPostDeleteCallback        = "presets_post_delete_callback"
	ParamPostExecuteActionCallback = "presets_post_execute_action_callback"
	ParamActionsDisabled           = "actions_disabled"
	ParamMustResult                = "must_result"
	ParamListingEncoder            = "presets_listingEncoder"

	// list editor
	ParamAddRowFormKey      = "listEditor_AddRowFormKey"
	ParamRemoveRowFormKey   = "listEditor_RemoveRowFormKey"
	ParamIsStartSort        = "listEditor_IsStartSort"
	ParamSortSectionFormKey = "listEditor_SortSectionFormKey"
	ParamSortResultFormKey  = "listEditor_SortResultFormKey"
)

func PermFromRequest(r *http.Request) string {
	method := r.FormValue("_method")
	if method == "" {
		method = r.Method
	}
	return PermFromHttpMethod(method)
}

func PermFromHttpMethod(method string) string {
	switch method {
	case "POST":
		return PermCreate
	case "PUT":
		return PermUpdate
	case "DELETE":
		return PermDelete
	default:
		return PermGet
	}
}
