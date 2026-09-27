package media

import (
	"fmt"
	"strings"

	"github.com/google/uuid"
)

func mainPortalName(field string) string {
	return fmt.Sprintf("%s_portal", field)
}

func deleteConfirmPortalName(field string) string {
	return fmt.Sprintf("%s_deleteConfirm_portal", field)
}

func mediaBoxThumbnailsPortalName(field string) string {
	return fmt.Sprintf("%s_portal_thumbnails", field)
}

func cropperPortalName(field string) string {
	return fmt.Sprintf("%s_cropper_portal", field)
}

func dialogContentPortalName(field string) string {
	return fmt.Sprintf("%s_dialog_content", field)
}

func searchKeywordName(field string) string {
	return fmt.Sprintf("%s_file_chooser_search_keyword", field)
}

func currentPageName(field string) string {
	return fmt.Sprintf("%s_file_chooser_current_page", field)
}

// fileCroppingVarName is a JS name: the file's key without its dashes.
func fileCroppingVarName(id uuid.UUID) string {
	return "fileChooser" + strings.ReplaceAll(id.String(), "-", "") + "_cropping"
}
