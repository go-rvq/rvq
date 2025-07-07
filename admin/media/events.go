package media

import (
	"github.com/go-rvq/rvq/admin/media/media_library"
	"github.com/go-rvq/rvq/admin/presets"
	"github.com/go-rvq/rvq/web"
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

func registerEventFuncs(hub web.EventHandlerHub, p *presets.Builder, b *Builder) {
	mb := presets.NewModelBuilder(p, &media_library.MediaLibrary{})
	hub.RegisterEventHandler(openFileChooserEvent, fileChooser(p, b))
	hub.RegisterEventHandler(deleteFileEvent, deleteFileField())
	hub.RegisterEventHandler(cropImageEvent, cropImage(b))
	hub.RegisterEventHandler(loadImageCropperEvent, loadImageCropper(b))
	hub.RegisterEventHandler(imageSearchEvent, searchFile(b))
	hub.RegisterEventHandler(imageJumpPageEvent, jumpPage(b))
	hub.RegisterEventHandler(uploadFileEvent, uploadFile(b))
	hub.RegisterEventHandler(chooseFileEvent, chooseFile(b))
	hub.RegisterEventHandler(updateDescriptionEvent, updateDescription(b))
	hub.RegisterEventHandler(deleteConfirmationEvent, deleteConfirmation(mb))
	hub.RegisterEventHandler(doDeleteEvent, doDelete(b))
}
