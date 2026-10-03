package presets

import "net/http"

const (
	PermModule = "presets"
	// The permissions of a resource — the last part of the resource asked,
	// "presets:site/:seo/:seo_config:<7>:@edit" —: an "@", then its name; an
	// action is a "!" and its name ("…:<7>:!publish", ActionPerm).
	PermList              = "@list"
	PermGet               = "@get"
	PermCreate            = "@create"
	PermUpdate            = "@edit"
	PermDelete            = "@delete"
	PermDeleteWithRelated = "@delete_with_related"

	PermActions         = "action"
	PermDoListingAction = "do_listing_action"
	PermBulkActions     = "bulk_action"
)

var PermRead = []string{PermList, PermGet}

// params
const (
	ParamID               = "id"
	ParamSelectedID       = "selected_id"
	ParamAction           = "action"
	ParamOverlay          = "overlay"
	ParamOverlayUpdateID  = "overlay_update_id"
	ParamBulkActionName   = "bulk_action"
	ParamSelectedIds      = "selected_ids"
	ParamListingQueries   = "presets_listing_queries"
	ParamAfterDeleteEvent = "presets_after_delete_event"
	ParamPortalID         = "portal_id"
	ParamTargetPortal     = "target_portal"
	// ParamCloserProvided tells the responder that the caller already provides the
	// overlay's `closer` (a form host owns it, see FormHost), so the response must
	// NOT wrap its content in a new closer scope — otherwise it would create a
	// child closer and closing it would no longer destroy the host's form.
	ParamCloserProvided = "presets_closer_provided"

	// ParamCloserRef carries HOW to address the caller's closer from anywhere —
	// `vars.$presetsCreating`, say. A host whose state is on `vars` can send it,
	// and then the overlay does not have to be rendered inside the host's own
	// portal to reach the closer: a drawer can go to the layout's portal, which
	// is the only place it sizes itself against the whole window.
	ParamCloserRef                 = "presets_closer_ref"
	ParamPostChangeCallback        = "presets_post_change_callback"
	ParamPostDeleteCallback        = "presets_post_delete_callback"
	ParamPostExecuteActionCallback = "presets_post_execute_action_callback"
	ParamActionsDisabled           = "actions_disabled"
	ParamMustResult                = "must_result"
	ParamListingEncoder            = "presets_listingEncoder"
	ParamRenderBreadcrumbs         = "presets_renderBreadcrumbs"
	ParamsItemTextKey              = "presets_itemTextKey"

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
