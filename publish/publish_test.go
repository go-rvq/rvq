package publish_test

import (
	"context"
	"errors"
	"fmt"
	"io"
	"os"
	"sort"
	"testing"
	"time"

	"github.com/qor/oss"
	"github.com/qor5/admin/v3/presets"
	"github.com/qor5/admin/v3/publish"
	"github.com/theplant/sliceutils"
	"github.com/theplant/testenv"
	"github.com/theplant/testingutils"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"gorm.io/gorm/logger"
)

type Product struct {
	gorm.Model
	Name string
	Code string

	publish.Version
	publish.Schedule
	publish.Status
	publish.List
}

func (p *Product) getContent() string {
	return p.Code + p.Name + p.VersionName
}

func (p *Product) getUrl() string {
	return fmt.Sprintf("test/product/%s/index.html", p.Code)
}

func (p *Product) getListUrl() string {
	return "test/product/list/index.html"
}

func (p *Product) getListContent() string {
	return fmt.Sprintf("list page  %s", p.Code)
}

func (p *Product) GetPublishActions(mb *presets.ModelBuilder, db *gorm.DB, ctx context.Context, storage oss.StorageInterface) (objs []*publish.PublishAction, err error) {
	objs = append(objs, &publish.PublishAction{
		Url:      p.getUrl(),
		Content:  p.getContent(),
		IsDelete: false,
	})
	p.OnlineUrl = p.getUrl()

	var liveRecord Product
	db.Where("id = ? AND status = ?", p.ID, publish.StatusOnline).First(&liveRecord)
	if liveRecord.ID == 0 {
		return
	}

	if liveRecord.OnlineUrl != p.OnlineUrl {
		objs = append(objs, &publish.PublishAction{
			Url:      liveRecord.getUrl(),
			IsDelete: true,
		})
	}

	if val, ok := ctx.Value("skip_list").(bool); ok && val {
		return
	}
	objs = append(objs, &publish.PublishAction{
		Url:      p.getListUrl(),
		Content:  p.getListContent(),
		IsDelete: false,
	})
	return
}

func (p *Product) GetUnPublishActions(mb *presets.ModelBuilder, db *gorm.DB, ctx context.Context, storage oss.StorageInterface) (objs []*publish.PublishAction, err error) {
	objs = append(objs, &publish.PublishAction{
		Url:      p.OnlineUrl,
		IsDelete: true,
	})
	if val, ok := ctx.Value("skip_list").(bool); ok && val {
		return
	}
	objs = append(objs, &publish.PublishAction{
		Url:      p.getListUrl(),
		IsDelete: true,
	})
	return
}

type ProductWithoutVersion struct {
	gorm.Model
	Name string
	Code string

	publish.Status
	publish.List
}

func (p *ProductWithoutVersion) getContent() string {
	return p.Code + p.Name
}

func (p *ProductWithoutVersion) getUrl() string {
	return fmt.Sprintf("test/product_no_version/%s/index.html", p.Code)
}

func (p *ProductWithoutVersion) GetPublishActions(mb *presets.ModelBuilder, db *gorm.DB, ctx context.Context, storage oss.StorageInterface) (objs []*publish.PublishAction, err error) {
	objs = append(objs, &publish.PublishAction{
		Url:      p.getUrl(),
		Content:  p.getContent(),
		IsDelete: false,
	})

	if p.Status.Status == publish.StatusOnline && p.OnlineUrl != p.getUrl() {
		objs = append(objs, &publish.PublishAction{
			Url:      p.OnlineUrl,
			IsDelete: true,
		})
	}

	p.OnlineUrl = p.getUrl()
	return
}

func (p *ProductWithoutVersion) GetUnPublishActions(mb *presets.ModelBuilder, db *gorm.DB, ctx context.Context, storage oss.StorageInterface) (objs []*publish.PublishAction, err error) {
	objs = append(objs, &publish.PublishAction{
		Url:      p.OnlineUrl,
		IsDelete: true,
	})
	return
}

func (this ProductWithoutVersion) GetListUrl(pageNumber string) string {
	return fmt.Sprintf("/product_without_version/list/%v.html", pageNumber)
}

func (this ProductWithoutVersion) GetListContent(db *gorm.DB, onePageItems *publish.OnePageItems) string {
	pageNumber := onePageItems.PageNumber
	var result string
	for _, item := range onePageItems.Items {
		record := item.(*ProductWithoutVersion)
		result = result + fmt.Sprintf("product:%v ", record.Name)
	}
	result = result + fmt.Sprintf("pageNumber:%v", pageNumber)
	return result
}

func (this ProductWithoutVersion) Sort(array []interface{}) {
	var temp []*ProductWithoutVersion
	sliceutils.Unwrap(array, &temp)
	sort.Sort(SliceProductWithoutVersion(temp))
	for k, v := range temp {
		array[k] = v
	}
	return
}

type SliceProductWithoutVersion []*ProductWithoutVersion

func (x SliceProductWithoutVersion) Len() int           { return len(x) }
func (x SliceProductWithoutVersion) Less(i, j int) bool { return x[i].Name < x[j].Name }
func (x SliceProductWithoutVersion) Swap(i, j int)      { x[i], x[j] = x[j], x[i] }

type MockStorage struct {
	oss.StorageInterface
	Objects map[string]string
}

func (m *MockStorage) Get(path string) (f *os.File, err error) {
	content, exist := m.Objects[path]
	if !exist {
		err = fmt.Errorf("NoSuchKey: %s", path)
		return
	}

	pattern := fmt.Sprintf("s3*%d", time.Now().Unix())

	if f, err = os.CreateTemp("/tmp", pattern); err == nil {
		f.WriteString(content)
		f.Seek(0, 0)
	}

	return
}

func (m *MockStorage) Put(path string, r io.Reader) (*oss.Object, error) {
	fmt.Println("Calling mock s3 client - Put: ", path)
	b, err := io.ReadAll(r)
	if err != nil {
		panic(err)
	}
	if m.Objects == nil {
		m.Objects = make(map[string]string)
	}
	m.Objects[path] = string(b)

	return &oss.Object{}, nil
}

func (m *MockStorage) Delete(path string) error {
	fmt.Println("Calling mock s3 client - Delete: ", path)
	delete(m.Objects, path)
	return nil
}

var TestDB *gorm.DB

func TestMain(m *testing.M) {
	env, err := testenv.New().DBEnable(true).SetUp()
	if err != nil {
		panic(err)
	}
	defer env.TearDown()
	TestDB = env.DB
	TestDB.Logger = TestDB.Logger.LogMode(logger.Info)

	m.Run()
}

func TestPublishVersionContentToS3(t *testing.T) {
	db := TestDB
	db.AutoMigrate(&Product{})
	storage := &MockStorage{}

	productV1 := Product{
		Model:   gorm.Model{ID: 1},
		Code:    "0001",
		Name:    "coffee",
		Status:  publish.Status{Status: publish.StatusDraft},
		Version: publish.Version{Version: "v1"},
	}
	productV2 := Product{
		Model:   gorm.Model{ID: 1},
		Code:    "0002",
		Name:    "coffee",
		Status:  publish.Status{Status: publish.StatusDraft},
		Version: publish.Version{Version: "v2"},
	}
	db.Clauses(clause.OnConflict{UpdateAll: true}).Create(&productV1)
	db.Clauses(clause.OnConflict{UpdateAll: true}).Create(&productV2)

	p := publish.New(db, storage)
	// publish v1
	skipListTrueContext := context.WithValue(context.Background(), "skip_list", true)
	skipListFalseContext := context.WithValue(context.Background(), "skip_list", false)
	if err := p.Publish(&productV1, skipListTrueContext); err != nil {
		t.Error(err)
	}
	assertUpdateStatus(t, db, &productV1, publish.StatusOnline, productV1.getUrl())
	assertUploadFile(t, productV1.getContent(), productV1.getUrl(), storage)
	// assertUploadFile(t, productV1.getListContent(), productV1.getListUrl(), storage)

	// publish v2
	if err := p.Publish(&productV2, skipListFalseContext); err != nil {
		t.Error(err)
	}
	assertUpdateStatus(t, db, &productV2, publish.StatusOnline, productV2.getUrl())
	assertUploadFile(t, productV2.getContent(), productV2.getUrl(), storage)
	assertUploadFile(t, productV2.getListContent(), productV2.getListUrl(), storage)
	// if delete v1 file
	assertContentDeleted(t, productV1.getUrl(), storage)
	// if update v1 status to offline
	assertUpdateStatus(t, db, &productV1, publish.StatusOffline, productV1.getUrl())

	// unpublish v2
	if err := p.UnPublish(&productV2, skipListFalseContext); err != nil {
		t.Error(err)
	}
	assertUpdateStatus(t, db, &productV2, publish.StatusOffline, productV2.getUrl())
	assertContentDeleted(t, productV2.getUrl(), storage)
	assertContentDeleted(t, productV2.getListUrl(), storage)
}

func TestPublishList(t *testing.T) {
	db := TestDB
	db.AutoMigrate(&ProductWithoutVersion{})
	storage := &MockStorage{}

	productV1 := ProductWithoutVersion{
		Model:  gorm.Model{ID: 1},
		Code:   "1",
		Name:   "1",
		Status: publish.Status{Status: publish.StatusDraft},
	}
	productV2 := ProductWithoutVersion{
		Model:  gorm.Model{ID: 2},
		Code:   "2",
		Name:   "2",
		Status: publish.Status{Status: publish.StatusDraft},
	}

	productV3 := ProductWithoutVersion{
		Model:  gorm.Model{ID: 3},
		Code:   "3",
		Name:   "3",
		Status: publish.Status{Status: publish.StatusDraft},
	}

	db.Clauses(clause.OnConflict{UpdateAll: true}).Create(&productV1)
	db.Clauses(clause.OnConflict{UpdateAll: true}).Create(&productV2)
	db.Clauses(clause.OnConflict{UpdateAll: true}).Create(&productV3)

	publisher := publish.New(db, storage)
	listPublisher := publish.NewListPublishBuilder(db, storage)

	publisher.Publish(&productV1, context.Background())
	publisher.Publish(&productV3, context.Background())
	if err := listPublisher.Run(ProductWithoutVersion{}); err != nil {
		panic(err)
	}

	var expected string
	expected = "product:1 product:3 pageNumber:1"
	if storage.Objects["/product_without_version/list/1.html"] != expected {
		t.Error(errors.New(fmt.Sprintf(`
want: %v
get: %v
`, expected, storage.Objects["/product_without_version/list/1.html"])))
	}

	publisher.Publish(&productV2, context.Background())
	if err := listPublisher.Run(ProductWithoutVersion{}); err != nil {
		panic(err)
	}

	expected = "product:1 product:2 product:3 pageNumber:1"
	if storage.Objects["/product_without_version/list/1.html"] != expected {
		t.Error(errors.New(fmt.Sprintf(`
want: %v
get: %v
`, expected, storage.Objects["/product_without_version/list/1.html"])))
	}

	publisher.UnPublish(&productV2, context.Background())
	if err := listPublisher.Run(ProductWithoutVersion{}); err != nil {
		panic(err)
	}

	expected = "product:1 product:3 pageNumber:1"
	if storage.Objects["/product_without_version/list/1.html"] != expected {
		t.Error(errors.New(fmt.Sprintf(`
want: %v
get: %v
`, expected, storage.Objects["/product_without_version/list/1.html"])))
	}

	publisher.UnPublish(&productV3, context.Background())
	if err := listPublisher.Run(ProductWithoutVersion{}); err != nil {
		panic(err)
	}

	expected = "product:1 pageNumber:1"
	if storage.Objects["/product_without_version/list/1.html"] != expected {
		t.Error(errors.New(fmt.Sprintf(`
want: %v
get: %v
`, expected, storage.Objects["/product_without_version/list/1.html"])))
	}
}

func TestSchedulePublish(t *testing.T) {
	db := TestDB
	db.Migrator().DropTable(&Product{})
	db.AutoMigrate(&Product{})
	storage := &MockStorage{}

	productV1 := Product{
		Model:   gorm.Model{ID: 1},
		Version: publish.Version{Version: "2021-12-19-v01"},
		Code:    "1",
		Name:    "1",
		Status:  publish.Status{Status: publish.StatusDraft},
	}

	db.Clauses(clause.OnConflict{UpdateAll: true}).Create(&productV1)

	publisher := publish.New(db, storage)
	publisher.Publish(&productV1, context.Background())

	var expected string
	expected = "11"
	if storage.Objects["test/product/1/index.html"] != expected {
		t.Error(errors.New(fmt.Sprintf(`
	want: %v
	get: %v
	`, expected, storage.Objects["test/product/1/index.html"])))
	}

	productV1.Name = "2"
	startAt := db.NowFunc().Add(-24 * time.Hour)
	productV1.ScheduledStartAt = &startAt
	if err := db.Save(&productV1).Error; err != nil {
		panic(err)
	}
	schedulePublisher := publish.NewSchedulePublishBuilder(publisher)
	if err := schedulePublisher.Run(productV1); err != nil {
		panic(err)
	}
	expected = "12"
	if storage.Objects["test/product/1/index.html"] != expected {
		t.Error(errors.New(fmt.Sprintf(`
	want: %v
	get: %v
	`, expected, storage.Objects["test/product/1/index.html"])))
	}

	endAt := startAt.Add(time.Second * 2)
	productV1.ScheduledEndAt = &endAt
	if err := db.Save(&productV1).Error; err != nil {
		panic(err)
	}
	if err := schedulePublisher.Run(productV1); err != nil {
		panic(err)
	}
	expected = ""
	if storage.Objects["test/product/1/index.html"] != expected {
		t.Error(errors.New(fmt.Sprintf(`
	want: %v
	get: %v
	`, expected, storage.Objects["test/product/1/index.html"])))
	}
}

func TestPublishContentWithoutVersionToS3(t *testing.T) {
	db := TestDB
	db.AutoMigrate(&ProductWithoutVersion{})
	storage := &MockStorage{}

	product1 := ProductWithoutVersion{
		Model:  gorm.Model{ID: 1},
		Code:   "0001",
		Name:   "tea",
		Status: publish.Status{Status: publish.StatusDraft},
	}
	db.Clauses(clause.OnConflict{UpdateAll: true}).Create(&product1)
	ctx := context.Background()

	p := publish.New(db, storage)
	// publish product1
	if err := p.Publish(&product1, ctx); err != nil {
		t.Error(err)
	}
	assertNoVersionUpdateStatus(t, db, &product1, publish.StatusOnline, product1.getUrl())
	assertUploadFile(t, product1.getContent(), product1.getUrl(), storage)

	product1Clone := product1
	product1Clone.Code = "0002"

	// publish product1 again
	if err := p.Publish(&product1Clone, ctx); err != nil {
		t.Error(err)
	}
	assertNoVersionUpdateStatus(t, db, &product1Clone, publish.StatusOnline, product1Clone.getUrl())

	assertUploadFile(t, product1Clone.getContent(), product1Clone.getUrl(), storage)

	// if delete product1 old file
	assertContentDeleted(t, product1.getUrl(), storage)
	// unpublish product1
	if err := p.UnPublish(&product1Clone, ctx); err != nil {
		t.Error(err)
	}
	assertNoVersionUpdateStatus(t, db, &product1Clone, publish.StatusOffline, product1Clone.getUrl())
	// if delete product1 file
	assertContentDeleted(t, product1Clone.getUrl(), storage)
}

func assertUpdateStatus(t *testing.T, db *gorm.DB, p *Product, assertStatus string, asserOnlineUrl string) {
	t.Helper()

	var pindb Product
	err := db.Model(&Product{}).Where("id = ? AND version = ?", p.ID, p.Version.Version).First(&pindb).Error
	if err != nil {
		t.Fatal(err)
	}

	diff := testingutils.PrettyJsonDiff(publish.Status{
		Status:    assertStatus,
		OnlineUrl: asserOnlineUrl,
	}, pindb.Status)
	if diff != "" {
		t.Error(diff)
	}
	return
}

func assertContentDeleted(t *testing.T, url string, storage oss.StorageInterface) {
	t.Helper()

	t.Helper()
	_, err := storage.Get(url)
	if err == nil {
		t.Errorf("content for %s should be deleted", url)
	}
	return
}

func assertNoVersionUpdateStatus(t *testing.T, db *gorm.DB, p *ProductWithoutVersion, assertStatus string, asserOnlineUrl string) {
	var pindb ProductWithoutVersion
	err := db.Model(&ProductWithoutVersion{}).Where("id = ?", p.ID).First(&pindb).Error
	if err != nil {
		t.Fatal(err)
	}
	diff := testingutils.PrettyJsonDiff(publish.Status{
		Status:    assertStatus,
		OnlineUrl: asserOnlineUrl,
	}, pindb.Status)
	if diff != "" {
		t.Error(diff)
	}
	return
}

func assertUploadFile(t *testing.T, content string, url string, storage oss.StorageInterface) {
	t.Helper()
	f, err := storage.Get(url)
	if err != nil {
		t.Fatal(err)
	}
	c, err := io.ReadAll(f)
	if err != nil {
		t.Fatal(err)
	}
	diff := testingutils.PrettyJsonDiff(content, string(c))
	if diff != "" {
		t.Error(diff)
	}
}
