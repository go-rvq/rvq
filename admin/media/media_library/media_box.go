package media_library

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"
	"path"
	"strings"

	"github.com/go-rvq/rvq/admin/media/base"
)

const (
	ALLOW_TYPE_FILE  = "file"
	ALLOW_TYPE_IMAGE = "image"
	ALLOW_TYPE_VIDEO = "video"
)

type MediaBox struct {
	ID          MediaID
	Url         string
	VideoLink   string
	FileName    string
	Description string
	FileSizes   map[string]int `json:",omitempty"`
	// for default image
	Width  int `json:",omitempty"`
	Height int `json:",omitempty"`
}

// MediaBoxConfig configure MediaBox metas
type MediaBoxConfig struct {
	Sizes     map[string]*base.Size
	Max       uint
	AllowType string
	// the background color of MediaBox
	BackgroundColor string
	// disable crop
	DisableCrop bool
	// DisableDescription hides the description for accessibility of the
	// field (a photo of a document, say, which has no other reader)
	DisableDescription bool
	// DescriptionField is the field of the record — beside the media box, by
	// its name — that holds the description, out of the media box: the box
	// shows it, and saves it there (its Description left empty). Empty, the
	// box's own Description.
	DescriptionField string `json:",omitempty"`
	// Card draws the box as a card of a list of images (media.ImagesCards):
	// the image itself — a click opens it in the lightbox, among the list's
	// images —, its dimensions and size over it, its name and description
	// under it; none yet, only the button that chooses it. Every event of the
	// box (choose, delete, edit) draws it so again.
	Card bool `json:",omitempty"`
	// EditOnUpload opens the image editor on the images chosen to upload,
	// in the browser, before they are sent (else only by the "Edit" of the
	// upload); a photo of the camera always opens it.
	EditOnUpload bool `json:",omitempty"`
	// UploadEdit sends, with the result of an image edited on the upload,
	// its original and the edit: the edit stays undoable (reopened on the
	// original). Else only the result.
	UploadEdit bool `json:",omitempty"`
	// EditTab and EditMode are where the editor opens: "recortar", "cores",
	// "anotar"; "fixo", "perspectiva", "livre" (a check's photo:
	// recortar, perspectiva).
	EditTab  string `json:",omitempty"`
	EditMode string `json:",omitempty"`
	// EditMaxPixels: a photo larger is made smaller before it is edited and
	// sent (a camera's of 48 MP); 0, never.
	EditMaxPixels int `json:",omitempty"`
	// OriginalMax is the largest an image's original is kept in: a larger
	// one, uploaded, is made smaller to fit it, its proportions kept (and
	// the limit turned with the image: a portrait one fits it on its side).
	// nil, a box with no sizes has DefaultOriginalMax (Full HD); one with
	// sizes, the media library's (base.Config).
	OriginalMax *base.Size `json:",omitempty"`
	// ThumbMax is the largest the thumbnail of an image of the box
	// (ThumbSizeName) is made in, its proportions kept: always made, not a
	// size of the field (none of its crops, nor shown in its form). nil, a
	// box with no sizes has DefaultThumbMax; one with sizes, none.
	ThumbMax *base.Size `json:",omitempty"`
	// Usage is the use of the field — "<table>.<Field>", a MediaBoxUsage —
	// whose sizes (the admin's, else the code's) the server reads: the
	// sizes the browser sends back are not the ones it trusts.
	Usage string `json:",omitempty"`
}

// DefaultOriginalMax is the largest an image's original is kept in, in a box
// with no sizes and no OriginalMax: Full HD.
var DefaultOriginalMax = base.Size{Width: 1920, Height: 1080}

// DefaultThumbMax is the largest the thumbnail of an image of a box with no
// sizes and no ThumbMax is made in.
var DefaultThumbMax = base.Size{Width: 320, Height: 320}

// ThumbLimit is the largest the thumbnail of an image uploaded to the box is
// made in (ThumbMax, else DefaultThumbMax with no sizes); nil, none. A nil
// config is a box's with no sizes.
func (c *MediaBoxConfig) ThumbLimit() *base.Size {
	if c != nil && c.ThumbMax != nil && c.ThumbMax.Width > 0 && c.ThumbMax.Height > 0 {
		s := *c.ThumbMax
		return &s
	}
	if c == nil || len(c.Sizes) == 0 {
		s := DefaultThumbMax
		return &s
	}
	return nil
}

// OriginalLimit is the largest the original of an image uploaded to the box
// is kept in (OriginalMax, else DefaultOriginalMax with no sizes); nil, the
// media library's. A nil config is a box's with no sizes.
func (c *MediaBoxConfig) OriginalLimit() *base.Size {
	if c != nil && c.OriginalMax != nil && c.OriginalMax.Width > 0 && c.OriginalMax.Height > 0 {
		s := *c.OriginalMax
		return &s
	}
	if c == nil || len(c.Sizes) == 0 {
		s := DefaultOriginalMax
		return &s
	}
	return nil
}

func (mediaBox *MediaBox) Scan(data interface{}) (err error) {
	switch values := data.(type) {
	case []byte:
		if len(values) > 0 {
			return json.Unmarshal(values, mediaBox)
		}
	case string:
		return mediaBox.Scan([]byte(values))
	}
	return nil
}

func (mediaBox MediaBox) Value() (driver.Value, error) {
	if mediaBox.ID.IsZero() {
		return nil, nil
	}
	results, err := json.Marshal(mediaBox)
	return string(results), err
}

// IsImage return if it is an image
func (mediaBox *MediaBox) IsImage() bool {
	return base.IsImageFormat(mediaBox.Url)
}

func (mediaBox *MediaBox) IsVideo() bool {
	return base.IsVideoFormat(mediaBox.Url)
}

func (mediaBox *MediaBox) IsSVG() bool {
	return base.IsSVGFormat(mediaBox.Url)
}

func (mediaBox *MediaBox) URL(styles ...string) string {
	if mediaBox.Url != "" && len(styles) > 0 {
		ext := path.Ext(mediaBox.Url)
		return fmt.Sprintf("%v.%v%v", strings.TrimSuffix(mediaBox.Url, ext), styles[0], ext)
	}
	return mediaBox.Url
}

func (mediaBox MediaBox) WebpURL(styles ...string) string {
	url := mediaBox.URL(styles...)
	ext := path.Ext(url)
	extArr := strings.Split(ext, "?")
	i := strings.LastIndex(url, ext)
	return url[:i] + strings.Replace(url[i:], extArr[0], ".webp", 1)
}
