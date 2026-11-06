package media

import (
	"bytes"
	"fmt"
	"image"
	"mime/multipart"
	"net/http"
	"path"
	"strings"

	"github.com/disintegration/imaging"
	"github.com/go-rvq/rvq/admin/media/base"
	"github.com/go-rvq/rvq/admin/media/media_library"
)

// The thumbnail of an image of a box (media_library.ThumbSizeName, of
// MediaBoxConfig.ThumbLimit) is made in the browser and sent with it
// (NewThumb<i> of an upload, Thumb of an edit): the server checks it — of
// the image's format, of its proportions, in the limit — and keeps it; only
// when none comes, or one that is not right, it makes it itself.

// postedThumb is the thumbnail the browser sent as key; nil, none.
func postedThumb(r *http.Request, key string) *multipart.FileHeader {
	if r.MultipartForm == nil {
		return nil
	}
	if fs := r.MultipartForm.File[key]; len(fs) > 0 {
		return fs[0]
	}
	return nil
}

// prepareThumb readies m, an image about to be saved, for its thumbnail of
// limit: made by the server as it is saved (its ThumbMax) when the browser
// sent none (posted), else kept from the browser after (keepThumb).
func prepareThumb(m *media_library.MediaLibrary, limit *base.Size, posted *multipart.FileHeader) {
	if m.SelectedType != media_library.ALLOW_TYPE_IMAGE || limit == nil {
		return
	}
	if posted != nil {
		m.File.ThumbMax = nil
		return
	}
	m.File.ThumbMax = limit
}

// keepThumb keeps posted, the thumbnail the browser made of the image of m
// (saved), when it is right (storeThumb); else the server makes it
// (makeThumb).
func (b *Builder) keepThumb(m *media_library.MediaLibrary, limit *base.Size, posted *multipart.FileHeader) error {
	if posted == nil || limit == nil || m.SelectedType != media_library.ALLOW_TYPE_IMAGE {
		return nil
	}
	m.File.ThumbMax = limit
	if !storeThumb(m, limit, posted) {
		if err := makeThumb(m, limit); err != nil {
			return err
		}
	}
	return b.db.Save(m).Error
}

// makeThumb makes the thumbnail of the image of m (its file, as saved) to
// fit limit.
func makeThumb(m *media_library.MediaLibrary, limit *base.Size) error {
	format, err := base.GetImageFormat(m.File.URL())
	if err != nil {
		return err
	}
	f, err := m.File.Retrieve(m.File.URL())
	if err != nil {
		return fmt.Errorf("media: the image of %s: %w", m.ID, err)
	}
	img, _, err := image.Decode(f)
	f.Close()
	if err != nil {
		return fmt.Errorf("media: the image of %s: %w", m.ID, err)
	}
	s := base.FitSize(img.Bounds().Dx(), img.Bounds().Dy(), *limit)
	var buf bytes.Buffer
	if err = imaging.Encode(&buf, imaging.Resize(img, s.Width, s.Height, imaging.CatmullRom), *format); err != nil {
		return err
	}
	n := buf.Len()
	if err = m.File.Store(m.File.URL(media_library.ThumbSizeName), nil, &buf); err != nil {
		return err
	}
	sizes := m.File.GetFileSizes()
	sizes[media_library.ThumbSizeName] = n
	m.File.FileSizes = sizes
	return nil
}

// storeThumb stores posted as the thumbnail of m when it is of the image's
// format, and of the size the image fits limit in (a pixel or two apart, of
// the rounding); false, not stored.
func storeThumb(m *media_library.MediaLibrary, limit *base.Size, posted *multipart.FileHeader) bool {
	f, err := posted.Open()
	if err != nil {
		return false
	}
	defer f.Close()
	c, format, err := image.DecodeConfig(f)
	if err != nil || format != imageFormatOf(m.File.URL()) || m.File.Width <= 0 || m.File.Height <= 0 {
		return false
	}
	want := base.FitSize(m.File.Width, m.File.Height, *limit)
	if abs(c.Width-want.Width) > 2 || abs(c.Height-want.Height) > 2 {
		return false
	}
	if _, err = f.Seek(0, 0); err != nil {
		return false
	}
	if err = m.File.Store(m.File.URL(media_library.ThumbSizeName), nil, f); err != nil {
		return false
	}
	sizes := m.File.GetFileSizes()
	sizes[media_library.ThumbSizeName] = int(posted.Size)
	m.File.FileSizes = sizes
	return true
}

// imageFormatOf is the format of the image of url (its extension), as
// image.DecodeConfig names it.
func imageFormatOf(url string) string {
	switch ext := strings.ToLower(path.Ext(url)); ext {
	case ".jpg", ".jpeg":
		return "jpeg"
	case "":
		return ""
	default:
		return ext[1:]
	}
}

func abs(n int) int {
	if n < 0 {
		return -n
	}
	return n
}
