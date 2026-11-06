package js

import "embed"

// ImageEditor is vx-image-editor (js/image-editor): the editor of the images
// of the media library, and window.rvqImageEditor.
//
//go:embed image-editor/dist
var ImageEditor embed.FS
