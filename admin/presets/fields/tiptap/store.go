package tiptap

import (
	"bytes"
	"context"

	"github.com/go-rvq/rvq/admin/media"
	"github.com/go-rvq/rvq/admin/media/base"
	"github.com/go-rvq/rvq/admin/media/media_library"
	"github.com/google/uuid"
	"github.com/iancoleman/strcase"
	"golang.org/x/net/html"
	"gorm.io/gorm"
)

type StoreImagesEntry struct {
	Format string
	Url    string
	Data   []byte
}

type StoreImagesResult struct {
	NewValue []byte
	Images   []*StoreImagesEntry
}

func (r *StoreImagesResult) Changed() bool {
	return len(r.Images) > 0
}

func (b *Builder) StoreImages(ctx context.Context, value []byte) (r *StoreImagesResult, err error) {
	var doc *html.Node
	if doc, err = html.Parse(bytes.NewReader(value)); err != nil {
		return
	}

	r = &StoreImagesResult{
		NewValue: value,
	}

	type image struct {
		n      *html.Node
		src    *html.Attribute
		format string
		data   []byte
		alt    string
	}

	var images []*image

	if err = traverseImg(doc, func(n *html.Node, src *html.Attribute, format string, data []byte) (err error) {
		img := &image{
			n:      n,
			format: format,
			src:    src,
			data:   data,
		}
		if altAttr := getAttr(n, "alt"); altAttr != nil {
			img.alt = altAttr.Val
		} else if n.Parent != nil && n.Parent.Data == "figure" {
			_ = traverseTag("figcaption", n.Parent, func(n *html.Node) error {
				s := text(n)
				if len(s) > 0 {
					img.alt = s
				}
				return nil
			})
		}

		images = append(images, img)

		return nil
	}); err != nil {
		return
	}

	for _, img := range images {
		var url string
		if url, err = b.store.Store(ctx, img.format, img.data, img.alt); err != nil {
			return
		}
		img.src.Val = url
		r.Images = append(r.Images, &StoreImagesEntry{
			Format: img.format,
			Url:    url,
			Data:   img.data,
		})
	}

	if r.Changed() {
		var b bytes.Buffer
		err = html.Render(&b, doc)
		if err != nil {
			return
		}
		r.NewValue = b.Bytes()
	}
	return
}

func MediaLibraryStorer(b *media.Builder, dbKey any) Storer {
	return StoreFunc(func(ctx context.Context, format string, data []byte, alt string) (url string, err error) {
		var (
			m = &media_library.MediaLibrary{
				SelectedType: media_library.ALLOW_TYPE_IMAGE,
				Hidden:       true,
			}
			fileName = strcase.ToSnake(alt)
			ext      = format
		)

		if len(fileName) == 0 {
			fileName = uuid.NewString()
		}

		switch ext {
		case "jpeg":
			ext = "jpg"
		}

		fileName += "." + ext

		if err = m.File.Scan(base.NewFileHeader(fileName, func() (base.Reader, error) {
			return bytes.NewReader(data), nil
		})); err != nil {
			return
		}

		db := ctx.Value(dbKey).(*gorm.DB)
		if err = base.SaveUploadAndCropImage(&b.Config, db, m); err != nil {
			return
		}

		url = m.File.URL("original")
		return
	})
}
