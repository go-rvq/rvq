package media

import (
	"github.com/qor5/admin/v3/media/media_library"
	"github.com/qor5/admin/v3/presets"
	"github.com/qor5/web/v3"
)

const (
	openFileChooserEvent    = "mediaLibrary_OpenFileChooserEvent"
	deleteFileEvent         = "mediaLibrary_DeleteFileEvent"
	cropImageEvent          = "mediaLibrary_CropImageEvent"
	loadImageCropperEvent   = "mediaLibrary_LoadImageCropperEvent"
	imageSearchEvent        = "mediaLibrary_ImageSearchEvent"
	imageJumpPageEvent      = "mediaLibrary_ImageJumpPageEvent"
	uploadFileEvent         = "mediaLibrary_UploadFileEvent"
	chooseFileEvent         = "mediaLibrary_ChooseFileEvent"
	updateDescriptionEvent  = "mediaLibrary_UpdateDescriptionEvent"
	deleteConfirmationEvent = "mediaLibrary_DeleteConfirmationEvent"
	doDeleteEvent           = "mediaLibrary_DoDelete"
)

func registerEventFuncs(hub web.EventFuncHub, p *presets.Builder, b *Builder) {
	mb := presets.NewModelBuilder(p, &media_library.MediaLibrary{})
	hub.RegisterEventFunc(openFileChooserEvent, fileChooser(b))
	hub.RegisterEventFunc(deleteFileEvent, deleteFileField())
	hub.RegisterEventFunc(cropImageEvent, cropImage(b))
	hub.RegisterEventFunc(loadImageCropperEvent, loadImageCropper(b))
	hub.RegisterEventFunc(imageSearchEvent, searchFile(b))
	hub.RegisterEventFunc(imageJumpPageEvent, jumpPage(b))
	hub.RegisterEventFunc(uploadFileEvent, uploadFile(b))
	hub.RegisterEventFunc(chooseFileEvent, chooseFile(b))
	hub.RegisterEventFunc(updateDescriptionEvent, updateDescription(b))
	hub.RegisterEventFunc(deleteConfirmationEvent, deleteConfirmation(mb))
	hub.RegisterEventFunc(doDeleteEvent, doDelete(b))
}
