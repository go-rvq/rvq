package base

import (
	"bytes"
	"database/sql/driver"
	"fmt"
	"image"
	"io"
	"strings"

	"gorm.io/gorm"
	"gorm.io/gorm/schema"
)

// Media is an interface including methods that needs for a media library storage
type Media interface {
	Scan(value interface{}) error
	Value() (driver.Value, error)

	GetURLTemplate(*Option) string
	GetURL(option *Option, db *gorm.DB, field *schema.Field, templater URLTemplater) string

	GetFileHeader() FileHeader
	GetFileName() string

	GetSizes() map[string]*Size
	NeedCrop() bool
	Cropped(values ...bool) bool
	GetCropOption(name string) *image.Rectangle
	GetFileSizes() map[string]int

	Store(url string, option *Option, reader io.Reader) error
	Retrieve(url string) (FileInterface, error)

	IsImage() bool

	URL(style ...string) string
	Ext() string
	String() string
}

// MediaSymlinker is an Media interface including methods to create symbolic link.
type MediaSymlinker interface {
	// Symlink create symbolic link from name to target.
	// If symbolic link creation isn't supported, return ErrSymlinkNotSupported.
	Symlink(target string, name string, option *Option) (err error)
}

// FileInterface media file interface
type FileInterface interface {
	io.ReadSeeker
	io.Closer
}

type FileBytes struct {
	*bytes.Reader
}

func (fb *FileBytes) Close() error {
	return nil
}

// Size is a struct, used for `GetSizes` method, it will return a slice of Size, media library will crop images automatically based on it
type Size struct {
	Width   int
	Height  int
	Padding bool
	// v-col sm
	Sm int
	// v-col col
	Cols int
}

func NewSize(width int, height int, fix ...int) *Size {
	if len(fix) > 0 {
		width, height = FixDimension(fix[0], width, height)
	}
	return &Size{Width: width, Height: height}
}

func (s *Size) FixDimension(max int) *Size {
	s.Width, s.Height = FixDimension(max, s.Width, s.Height)
	return s
}

func (s *Size) RescaleToMax(maxWidth, maxHeight int) (s2 *Size) {
	s2 = NewSize(s.Width, s.Height, maxWidth)
	if s.Height > maxHeight {
		ratio := float64(maxHeight) / float64(s.Height)
		s2.Width = int(float64(s.Width) * ratio)
		s2.Height = int(float64(s.Height) * ratio)
	}
	return
}

func (s *Size) String() string {
	return fmt.Sprintf("%dx%d", s.Width, s.Height)
}

// URLTemplater is a interface to return url template
type URLTemplater interface {
	GetURLTemplate(*Option) string
}

// Option media library option
type Option map[string]string

// get option with name
func (option Option) Get(key string) string {
	return option[strings.ToUpper(key)]
}

// set option
func (option Option) Set(key string, val string) {
	option[strings.ToUpper(key)] = val
	return
}
