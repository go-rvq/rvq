package db_tools

import (
	"context"
	"io"
	"net/http"
	"time"
)

type Backuper interface {
	Open() (io.ReadCloser, error)
	GetCreatedAt() time.Time
	GetID() BackupID
	GetDbName() string
	GetMessage() string
	GetSize() int64
	IsAuto() bool
	DetailString(ctx context.Context) string
}

type ArchivedBackuper interface {
	Backuper
	OpenArchive() (io.ReadCloser, error)
}

type Backup struct {
	OpenFunc   func() (io.ReadCloser, error) `json:"-" gorm:"-"`
	DetailFunc func(ctx context.Context) string
	CreatedAt  time.Time
	DbName     string
	Message    string
	Size       int64
	Auto       bool
}

func (b *Backup) GetID() (id BackupID) {
	id.DbName = b.DbName
	id.Auto = b.Auto
	id.CreatedAt = b.CreatedAt
	return
}

func (b *Backup) Open() (io.ReadCloser, error) {
	return b.OpenFunc()
}

func (b *Backup) GetCreatedAt() time.Time {
	return b.CreatedAt
}

func (b *Backup) GetDbName() string {
	return b.DbName
}

func (b *Backup) GetMessage() string {
	return b.Message
}

func (b *Backup) GetSize() int64 {
	return b.Size
}

func (b *Backup) IsAuto() bool {
	return b.Auto
}

func (b *Backup) DetailString(ctx context.Context) string {
	return b.DetailFunc(ctx)
}

type BackupController interface {
	CurrentName(ctx context.Context) (name string, err error)
	Create(auto bool, message string) (bkp Backuper, err error)
	Remove(id BackupID) error
	RemoveOlder(auto bool, p *Persistence) (removed []Backuper, err error)
	Get(id BackupID) (Backuper, error)
	Download(aw http.ResponseWriter, r *http.Request, id BackupID) (err error)
	List(auto bool, cb func(bkp Backuper) error, filter ListFilter) (err error)
	Count(auto bool, filter ListFilter) (count int, err error)
}
