package actions

const (
	New                        = "presets_New"
	Edit                       = "presets_Edit"
	Action                     = "presets_Action"
	DeleteConfirmation         = "presets_DeleteConfirmation"
	Update                     = "presets_Update"
	Create                     = "presets_Create"
	DoAction                   = "presets_DoAction"
	DoDelete                   = "presets_DoDelete"
	DoBulkAction               = "presets_DoBulkAction"
	DoListingAction            = "presets_DoListingAction"
	OpenBulkActionDialog       = "presets_OpenBulkActionDialog"
	OpenActionDialog           = "presets_OpenActionDialog"
	NotificationCenter         = "presets_NotificationCenter"
	Detailing                  = "presets_Detailing"
	DetailingContent           = "presets_DetailingContent"
	DoSaveDetailingField       = "presets_Detailing_Field_Save"
	DoEditDetailingField       = "presets_Detailing_Field_Edit"
	DoEditDetailingListField   = "presets_Detailing_List_Field_Edit"
	DoSaveDetailingListField   = "presets_Detailing_List_Field_Save"
	DoDeleteDetailingListField = "presets_Detailing_List_Field_Delete"
	DoCreateDetailingListField = "presets_Detailing_List_Field_Create"

	ListData                      = "presets_ListData"
	ReloadList                    = "presets_ReloadList"
	OpenListingDialog             = "presets_OpenListingDialog"
	OpenListingDialogForSelection = "presets_OpenListingDialogForSelection"
	UpdateListingDialog           = "presets_UpdateListingDialog"

	// list editor
	AddRowEvent    = "listEditor_addRowEvent"
	RemoveRowEvent = "listEditor_removeRowEvent"
	SortEvent      = "listEditor_sortEvent"
)

type OverlayMode string

const (
	Dialog       OverlayMode = "Dialog"
	StartDrawer  OverlayMode = "StartDrawer"
	EndDrawer    OverlayMode = "EndDrawer"
	LeftDrawer   OverlayMode = "LeftDrawer"
	RightDrawer  OverlayMode = "RightDrawer"
	TopDrawer    OverlayMode = "TopDrawer"
	BottomDrawer OverlayMode = "BottomDrawer"
	Content      OverlayMode = "Content"
)

func (m OverlayMode) PortalName() string {
	return "presets_Overlay" + string(m)
}

func (m OverlayMode) ContentPortalName() string {
	return "presets_Overlay" + string(m) + "Content"
}

func (m OverlayMode) CloseScript() string {
	if m == "" {
		return ""
	}
	if m == Dialog {
		return "locals.presets" + string(m) + " = false"
	}
	return "vars.presets" + string(m) + " = false"
}

func (m OverlayMode) String() string {
	return string(m)
}

func (m OverlayMode) Is(o ...OverlayMode) bool {
	for _, i := range o {
		if i == m {
			return true
		}
	}
	return false
}

func (m OverlayMode) IsDrawer() bool {
	return m.Is(LeftDrawer, RightDrawer, TopDrawer, BottomDrawer, StartDrawer, EndDrawer)
}

func (m OverlayMode) IsDialog() bool {
	return m == Dialog
}

func (m OverlayMode) Overlayed() bool {
	return m.IsDrawer() || m.IsDialog()
}

func (m OverlayMode) Up() OverlayMode {
	if m.Overlayed() {
		return Dialog
	}
	return RightDrawer
}
