package media_library

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"
	"path"
	"strings"

	"github.com/go-rvq/rvq/admin/media/base"

	"github.com/go-rvq/rvq/admin/media/oss"
	"gorm.io/gorm"
)

var (
	QorPreviewSizeName = "@qor_preview"
	QorPreviewMaxSize  = 200
)

type MediaLibrary struct {
	gorm.Model
	SelectedType string
	File         MediaLibraryStorage `sql:"size:4294967295;" mediaLibrary:"url:/system/{{class}}/{{primary_key}}/{{column}}.{{extension}}"`
}

type MediaOption struct {
	Video          string                      `json:",omitempty"`
	FileName       string                      `json:",omitempty"`
	URL            string                      `json:",omitempty"`
	OriginalURL    string                      `json:",omitempty"`
	OriginalStored bool                        `json:",omitempty"`
	CropOptions    map[string]*base.CropOption `json:",omitempty"`
	Sizes          map[string]*base.Size       `json:",omitempty"`
	SelectedType   string                      `json:",omitempty"`
	Description    string                      `json:",omitempty"`
	Crop           bool
}

func (mediaLibrary *MediaLibrary) ScanMediaOptions(mediaOption MediaOption) error {
	bytes, err := json.Marshal(mediaOption)
	if err == nil {
		return mediaLibrary.File.Scan(bytes)
	}
	return err
}

func (mediaLibrary *MediaLibrary) GetMediaOption() MediaOption {
	return MediaOption{
		Video:        mediaLibrary.File.Video,
		FileName:     mediaLibrary.File.GetFileName(),
		URL:          mediaLibrary.File.URL(),
		OriginalURL:  mediaLibrary.File.URL("original"),
		CropOptions:  mediaLibrary.File.CropOptions,
		Sizes:        mediaLibrary.File.GetSizes(),
		SelectedType: mediaLibrary.File.SelectedType,
		Description:  mediaLibrary.File.Description,
	}
}

func (mediaLibrary *MediaLibrary) SetSelectedType(typ string) {
	mediaLibrary.SelectedType = typ
}

func (mediaLibrary *MediaLibrary) GetSelectedType() string {
	return mediaLibrary.SelectedType
}

type MediaLibraryStorage struct {
	oss.OSS
	Sizes        map[string]*base.Size `json:",omitempty"`
	Video        string
	SelectedType string
	Description  string
}

func (mediaLibraryStorage MediaLibraryStorage) GetSizes() map[string]*base.Size {
	if len(mediaLibraryStorage.Sizes) == 0 && !(mediaLibraryStorage.GetFileHeader() != nil || mediaLibraryStorage.Crop) {
		return map[string]*base.Size{}
	}

	sizes := map[string]*base.Size{
		QorPreviewSizeName: base.NewSize(mediaLibraryStorage.Width, mediaLibraryStorage.Height, QorPreviewMaxSize),
	}

	for key, value := range mediaLibraryStorage.Sizes {
		sizes[key] = value
	}
	return sizes
}

func (mediaLibraryStorage *MediaLibraryStorage) Scan(data interface{}) (err error) {
	switch values := data.(type) {
	case []byte:
		if mediaLibraryStorage.Sizes == nil {
			mediaLibraryStorage.Sizes = map[string]*base.Size{}
		}
		// cropOptions := mediaLibraryStorage.CropOptions
		sizeOptions := mediaLibraryStorage.Sizes

		if string(values) != "" {
			mediaLibraryStorage.Base.Scan(values)
			if err = json.Unmarshal(values, mediaLibraryStorage); err == nil {
				if mediaLibraryStorage.CropOptions == nil {
					mediaLibraryStorage.CropOptions = map[string]*base.CropOption{}
				}

				// for key, value := range cropOptions {
				// 	if _, ok := mediaLibraryStorage.CropOptions[key]; !ok {
				// 		mediaLibraryStorage.CropOptions[key] = value
				// 	}
				// }

				for key, value := range sizeOptions {
					if key == base.DefaultSizeKey {
						continue
					}
					if _, ok := mediaLibraryStorage.Sizes[key]; !ok {
						mediaLibraryStorage.Sizes[key] = value
					}

				}

				for key, value := range mediaLibraryStorage.CropOptions {
					if key == base.DefaultSizeKey {
						continue
					}
					if _, ok := mediaLibraryStorage.Sizes[key]; !ok {
						mediaLibraryStorage.Sizes[key] = &base.Size{Width: value.Width, Height: value.Height}
					}

				}
			}
		}
	case string:
		err = mediaLibraryStorage.Scan([]byte(values))
	case []string:
		for _, str := range values {
			if err = mediaLibraryStorage.Scan(str); err != nil {
				return err
			}
		}
	default:
		return mediaLibraryStorage.Base.Scan(data)
	}
	return nil
}

func (mediaLibraryStorage MediaLibraryStorage) Value() (driver.Value, error) {
	results, err := json.Marshal(mediaLibraryStorage)
	return string(results), err
}

func (mediaLibraryStorage MediaLibraryStorage) URL(styles ...string) string {
	if mediaLibraryStorage.Url != "" && len(styles) > 0 {
		ext := path.Ext(mediaLibraryStorage.Url)
		return fmt.Sprintf("%v.%v%v", strings.TrimSuffix(mediaLibraryStorage.Url, ext), styles[0], ext)
	}
	return mediaLibraryStorage.Url
}
