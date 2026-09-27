package integration_test

import (
	"embed"
	"testing"

	"github.com/go-rvq/rvq/admin/media/base"
	"github.com/go-rvq/rvq/admin/media/media_library"
	"github.com/go-rvq/rvq/admin/media/oss"
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

	db = db.Debug()
	// db.Logger = db.Logger.LogMode(logger.Info)

	if err = db.AutoMigrate(
		&media_library.MediaLibrary{},
	); err != nil {
		panic(err)
	}

	oss.Storage = filesystem.New("/tmp/media_test")

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

	err = base.SaveUploadAndCropImage(db, &m)
	if err != nil {
		t.Fatal(err)
	}
}
