package microsite

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"path"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/go-rvq/rvq/admin/media/storage"
	"github.com/go-rvq/rvq/admin/microsite/utils"
	"github.com/go-rvq/rvq/admin/presets"
	"github.com/go-rvq/rvq/admin/publish"
	"github.com/hashicorp/go-multierror"
	"github.com/mholt/archiver/v4"
	"gorm.io/gorm"
)

type MicroSite struct {
	gorm.Model
	publish.Status
	publish.Schedule
	publish.Version
	Name        string
	Description string
	PrePath     string
	Package     FileSystem `gorm:"type:text"`
	FilesList   string     `gorm:"type:text"`
	UnixKey     string
}

func (ms *MicroSite) PermissionRN() []string {
	return []string{"microsite_models", strconv.Itoa(int(ms.ID)), ms.Version.Version}
}

func (ms *MicroSite) PrimarySlug() string {
	return fmt.Sprintf("%v_%v", ms.ID, ms.Version.Version)
}

func (ms *MicroSite) PrimaryColumnValuesBySlug(slug string) map[string]string {
	segs := strings.Split(slug, "_")
	if len(segs) != 2 {
		panic("wrong slug")
	}

	return map[string]string{
		"id":      segs[0],
		"version": segs[1],
	}
}

func (ms MicroSite) GetID() uint {
	return ms.ID
}

func (ms MicroSite) GetUnixKey() string {
	return ms.UnixKey
}

func (ms *MicroSite) SetUnixKey() {
	ms.UnixKey = strconv.FormatInt(time.Now().UnixMilli(), 10)
}

func (ms MicroSite) GetPublishedPath(fileName string) string {
	return path.Join(strings.TrimPrefix(strings.TrimSuffix(ms.PrePath, "/"), "/"), fileName)
}

func (ms MicroSite) GetPublishedUrl(domain, fileName string) string {
	return strings.TrimSuffix(domain, "/") + "/" + ms.GetPublishedPath(fileName)
}

func (ms MicroSite) GetPackageUrl(domain string) string {
	return strings.TrimSuffix(domain, "/") + "/" + strings.TrimPrefix(ms.Package.Url, "/")
}

func (ms MicroSite) GetFileList() (arr []string) {
	json.Unmarshal([]byte(ms.FilesList), &arr)
	return
}

func (ms *MicroSite) SetFilesList(filesList []string) {
	list, err := json.Marshal(filesList)
	if err != nil {
		return
	}
	ms.FilesList = string(list)
}

func (ms *MicroSite) GetPackage() FileSystem {
	return ms.Package
}

func (ms *MicroSite) SetPackage(fileName, url string) {
	ms.Package.FileName = fileName
	ms.Package.Url = url
}

type contextKeyType int

const contextKey contextKeyType = iota

func (b *Builder) ContextValueProvider(in context.Context) context.Context {
	return context.WithValue(in, contextKey, b)
}

func builderFromContext(c context.Context) (b *Builder, ok bool) {
	b, ok = c.Value(contextKey).(*Builder)
	return
}

func (ms *MicroSite) GetPublishActions(mb *presets.ModelBuilder, db *gorm.DB, ctx context.Context, storage storage.Storage) (objs []*publish.PublishAction, err error) {
	if len(ms.GetFileList()) == 0 {
		return
	}
	mib, ok := builderFromContext(ctx)
	if !ok {
		panic("use publisher.ContextValueFuncs(micrositeBuilder.ContextValueFunc) to set up microsite builder into context")
	}

	var previewPaths []string

	wg := sync.WaitGroup{}
	var copyError error
	var mutex sync.Mutex
	for _, v := range ms.GetFileList() {
		wg.Add(1)
		copySemaphore <- struct{}{}
		go func(v string) {
			defer func() {
				wg.Done()
				<-copySemaphore
			}()
			err = utils.Copy(storage, getPreviewPath(ms, v, mib), ms.GetPublishedPath(v))
			if err != nil {
				mutex.Lock()
				copyError = multierror.Append(copyError, err).ErrorOrNil()
				mutex.Unlock()
				return
			}
			mutex.Lock()
			previewPaths = append(previewPaths, getPreviewPath(ms, v, mib))
			mutex.Unlock()
		}(v)
	}

	wg.Wait()

	if len(previewPaths) > 0 {
		err = utils.DeleteObjects(storage, previewPaths)
	}
	err = multierror.Append(err, copyError).ErrorOrNil()

	return
}

func (ms *MicroSite) GetUnPublishActions(mb *presets.ModelBuilder, db *gorm.DB, ctx context.Context, storage storage.Storage) (objs []*publish.PublishAction, err error) {
	var paths []string
	for _, v := range ms.GetFileList() {
		paths = append(paths, ms.GetPublishedPath(v))
	}
	err = utils.DeleteObjects(storage, paths)
	if err != nil {
		return
	}
	return
}

func (ms *MicroSite) UnArchiveAndPublish(getPath func(string) string, fileName string, f io.Reader, storage storage.Storage) (filesList []string, err error) {
	format, reader, err := archiver.Identify(fileName, f)
	if err != nil {
		if err == archiver.ErrNoMatch {
			err = utils.Upload(storage, getPath(fileName), f)
			return
		}
		return
	}

	wg := sync.WaitGroup{}
	var putError error
	var mutex sync.Mutex

	err = format.(archiver.Extractor).Extract(context.Background(), reader, nil, func(ctx context.Context, f archiver.File) (err error) {
		if f.IsDir() {
			return
		}

		rc, err := f.Open()
		if err != nil {
			return
		}
		defer rc.Close()

		filesList = append(filesList, f.NameInArchive)

		publishedPath := getPath(f.NameInArchive)
		wg.Add(1)
		putSemaphore <- struct{}{}
		go func() {
			defer func() {
				<-putSemaphore
				wg.Done()
			}()
			err2 := utils.Upload(storage, publishedPath, rc)
			if err2 != nil {
				mutex.Lock()
				putError = multierror.Append(putError, err2).ErrorOrNil()
				mutex.Unlock()
			}
		}()

		return
	})
	wg.Wait()
	err = multierror.Append(err, putError).ErrorOrNil()
	return
}
