package integration_test

import (
	"embed"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/go-rvq/rvq/admin/media/base"
	"github.com/go-rvq/rvq/admin/media/media_library"
	"github.com/go-rvq/rvq/admin/media/oss"
	"github.com/go-rvq/rvq/admin/media/storage"
	"github.com/go-rvq/rvq/admin/utils/uuidkey"
	"github.com/go-rvq/rvq/web/multipartestutils"
	"github.com/theplant/testenv"
	"gorm.io/gorm"
)

//go:embed *.png
var box embed.FS

var TestDB *gorm.DB

func TestMain(m *testing.M) {
	env, err := testenv.New().DBEnable(true).SetUp()
	if err != nil {
		panic(err)
	}
	defer env.TearDown()
	TestDB = env.DB
	m.Run()
}

func setup() (db *gorm.DB) {
	var err error
	db = TestDB

	// the records get UUID keys, as media.New arranges in an app
	uuidkey.MustRegister(db)

	if err = db.AutoMigrate(
		&media_library.MediaLibrary{},
	); err != nil {
		panic(err)
	}

	oss.Storage = storage.NewFileSystem("/tmp/media_test")

	return
}

func TestUpload(t *testing.T) {
	db := setup()
	f, err := box.ReadFile("testfile.png")
	if err != nil {
		panic(err)
	}

	fh := multipartestutils.CreateMultipartFileHeader("test.png", f)
	m := media_library.MediaLibrary{}

	err = m.File.Scan(fh)
	if err != nil {
		t.Fatal(err)
	}

	err = base.SaveUploadAndCropImage(&base.Config{}, db, &m)
	if err != nil {
		t.Fatal(err)
	}

	// stored under the key's two levels (uuidkey.ShortPath)
	dir := "/system/media_libraries/" + uuidkey.ShortPath(m.ID) + "/"
	if !strings.HasPrefix(m.File.Url, dir) {
		t.Fatalf("url %q is not under %q", m.File.Url, dir)
	}
	if _, err := os.Stat(filepath.Join("/tmp/media_test", m.File.Url)); err != nil {
		t.Errorf("the file is not where its url says: %v", err)
	}
}
