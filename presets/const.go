package presets

const (
	PermModule = "presets"
	PermList   = "presets:list"
	PermGet    = "presets:get"
	PermCreate = "presets:create"
	PermUpdate = "presets:update"
	PermDelete = "presets:delete"

	PermActions         = "actions"
	PermDoListingAction = "do_listing_action"
	PermBulkActions     = "bulk_actions"
)

var PermRead = []string{PermList, PermGet}

// params
const (
	ParamID                 = "id"
	ParamSelectedID         = "selected_id"
	ParamAction             = "action"
	ParamOverlay            = "overlay"
	ParamOverlayUpdateID    = "overlay_update_id"
	ParamBulkActionName     = "bulk_action"
	ParamListingActionName  = "listing_action"
	ParamSelectedIds        = "selected_ids"
	ParamListingQueries     = "presets_listing_queries"
	ParamAfterDeleteEvent   = "presets_after_delete_event"
	ParamPortalID           = "portal_id"
	ParamTargetPortal       = "target_portal"
	ParamPostChangeCallback = "presets_post_change_callback"
	ParamActionsDisabled    = "actions_disabled"
	ParamMustResult         = "must_result"

	// list editor
	ParamAddRowFormKey      = "listEditor_AddRowFormKey"
	ParamRemoveRowFormKey   = "listEditor_RemoveRowFormKey"
	ParamIsStartSort        = "listEditor_IsStartSort"
	ParamSortSectionFormKey = "listEditor_SortSectionFormKey"
	ParamSortResultFormKey  = "listEditor_SortResultFormKey"
)
