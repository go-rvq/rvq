package models

import (
	"context"
	"fmt"
	"sort"
	"strconv"
	"strings"

	"github.com/go-rvq/rvq/admin/media/storage"
	"github.com/go-rvq/rvq/admin/presets"
	"github.com/go-rvq/rvq/admin/publish"
	"github.com/theplant/sliceutils"
	"gorm.io/gorm"
)

type ListModel struct {
	gorm.Model

	Title string

	publish.Status
	publish.Schedule
	publish.List
	publish.Version
}

func (lm *ListModel) PrimarySlug() string {
	return fmt.Sprintf("%v_%v", lm.ID, lm.Version.Version)
}

func (lm *ListModel) PermissionRN() []string {
	return []string{"list_models", strconv.Itoa(int(lm.ID)), lm.Version.Version}
}

func (lm *ListModel) PrimaryColumnValuesBySlug(slug string) map[string]string {
	segs := strings.Split(slug, "_")
	if len(segs) != 2 {
		panic("wrong slug")
	}

	return map[string]string{
		"id":      segs[0],
		"version": segs[1],
	}
}

func (lm *ListModel) GetPublishActions(mb *presets.ModelBuilder, db *gorm.DB, ctx context.Context, _ storage.Storage) (objs []*publish.PublishAction, err error) {
	objs = append(objs, &publish.PublishAction{
		Url:      lm.getPublishUrl(),
		Content:  lm.getPublishContent(),
		IsDelete: false,
	})

	if lm.Status.Status == publish.StatusOnline && lm.OnlineUrl != lm.getPublishUrl() {
		objs = append(objs, &publish.PublishAction{
			Url:      lm.OnlineUrl,
			IsDelete: true,
		})
	}

	lm.OnlineUrl = lm.getPublishUrl()
	return
}

func (lm *ListModel) GetUnPublishActions(mb *presets.ModelBuilder, db *gorm.DB, ctx context.Context, _ storage.Storage) (objs []*publish.PublishAction, err error) {
	objs = append(objs, &publish.PublishAction{
		Url:      lm.OnlineUrl,
		IsDelete: true,
	})
	return
}

func (lm ListModel) getPublishUrl() string {
	return fmt.Sprintf("/list_model/%v/index.html", lm.ID)
}

func (lm ListModel) getPublishContent() string {
	return fmt.Sprintf("id: %v, title: %v", lm.ID, lm.Title)
}

func (lm ListModel) GetListUrl(pageNumber string) string {
	return fmt.Sprintf("/list_model/list/%v.html", pageNumber)
}

func (lm ListModel) GetListContent(db *gorm.DB, onePageItems *publish.OnePageItems) string {
	pageNumber := onePageItems.PageNumber
	var result string
	for _, item := range onePageItems.Items {
		result = result + fmt.Sprintf("%v</br>", item)
	}
	result = result + fmt.Sprintf("</br>pageNumber: %v</br>", pageNumber)
	return result
}

func (lm ListModel) Sort(array []interface{}) {
	var temp []*ListModel
	sliceutils.Unwrap(array, &temp)
	sort.Sort(SliceListModel(temp))
	for k, v := range temp {
		array[k] = v
	}
}

type SliceListModel []*ListModel

func (x SliceListModel) Len() int           { return len(x) }
func (x SliceListModel) Less(i, j int) bool { return x[i].Title < x[j].Title }
func (x SliceListModel) Swap(i, j int)      { x[i], x[j] = x[j], x[i] }
