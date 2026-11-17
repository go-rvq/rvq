package notifications

import (
	"fmt"
	"github.com/go-rvq/rvq/admin/presets"
	"mime"
	"net/http"
	"path/filepath"
	"time"

	"github.com/go-rvq/rvq/admin/packages/comments"
	"github.com/go-rvq/rvq/admin/packages/user"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

// FileMaxSize is the largest attachment taken, in bytes.
const FileMaxSize = comments.FileMaxSize

// File is an attachment of a message (MessageID): of each copy of it, seen
// by who received it and by who sent it (FileHandler). Kept on disk (Path,
// under Builder.FilesDir); without a directory — or one of before —, in Data.
type File struct {
	ID          uuid.UUID `gorm:"type:uuid;primaryKey"`
	MessageID   uuid.UUID `gorm:"type:uuid;not null;index"`
	Name        string    `gorm:"size:255;not null"`
	ContentType string    `gorm:"size:255"`
	Size        int64
	// Path is under the directory of the files: "<message>/<id><ext>"
	Path      string `gorm:"size:255"`
	Data      []byte
	CreatedAt time.Time
}

func (File) TableName() string { return "user_notification_files" }

// AttachFiles keeps the files as the attachments of the message: in dir
// (comments.SaveFiles), or in the database (dir "").
func AttachFiles(db *gorm.DB, dir string, message uuid.UUID, files ...comments.Attachment) error {
	if dir != "" {
		stored, err := comments.SaveFiles(dir, message, files...)
		if err != nil {
			return err
		}
		for _, sf := range stored {
			f := File{ID: sf.ID, MessageID: message, Name: sf.Name, ContentType: sf.ContentType, Size: sf.Size, Path: sf.Path}
			if err := db.Create(&f).Error; err != nil {
				return err
			}
		}
		return nil
	}
	for _, a := range files {
		if len(a.Data) > FileMaxSize {
			return fmt.Errorf("%w: %s", comments.ErrFileTooLarge, a.Name)
		}
		ct := a.ContentType
		if ct == "" {
			ct = mime.TypeByExtension(filepath.Ext(a.Name))
		}
		f := File{ID: uuid.New(), MessageID: message, Name: a.Name, ContentType: ct, Size: int64(len(a.Data)), Data: a.Data}
		if err := db.Create(&f).Error; err != nil {
			return err
		}
	}
	return nil
}

// MoveFilesToDir writes the attachments kept in the database into dir, each
// in its place there, and empties their data: once, the ones of before a
// directory was given.
func MoveFilesToDir(db *gorm.DB, dir string) (moved int, err error) {
	var ids []uuid.UUID
	if err := db.Model(&File{}).Where("(path IS NULL OR path = '') AND data IS NOT NULL").Pluck("id", &ids).Error; err != nil {
		return 0, err
	}
	for _, id := range ids {
		var f File
		if err := db.First(&f, "id = ?", id).Error; err != nil {
			return moved, err
		}
		stored, err := comments.SaveFiles(dir, f.MessageID, comments.Attachment{Name: f.Name, ContentType: f.ContentType, Data: f.Data})
		if err != nil {
			return moved, err
		}
		if err := db.Model(&File{}).Where("id = ?", f.ID).
			Updates(map[string]any{"path": stored[0].Path, "data": nil}).Error; err != nil {
			return moved, err
		}
		moved++
	}
	return moved, nil
}

// FilesOf are the attachments of the messages (no data), by message.
func FilesOf(db *gorm.DB, messages ...uuid.UUID) (map[uuid.UUID][]File, error) {
	out := map[uuid.UUID][]File{}
	if len(messages) == 0 {
		return out, nil
	}
	var fs []File
	if err := db.Select("id", "message_id", "name", "content_type", "size", "path", "created_at").
		Where("message_id IN ?", messages).Order("created_at, name").Find(&fs).Error; err != nil {
		return nil, err
	}
	for _, f := range fs {
		out[f.MessageID] = append(out[f.MessageID], f)
	}
	return out, nil
}

// canSee says whether the user received or sent a copy of the message.
func canSee(db *gorm.DB, uid, message uuid.UUID) bool {
	var n int64
	db.Model(&Notification{}).Where("message_id = ? AND (user_id = ? OR source = ?)", message, uid, SourceUser(uid)).Count(&n)
	return n > 0
}

// filesPath is where the attachments are served, under the admin's prefix:
// the service notifications (presets.Service), no page of the menu.
const filesPath = "/" + presets.ServicePrefix + "notifications/files/"

// FileHandler serves the attachment of its id (the end of the path), of dir
// (or of the database), to who received or sent its message; to anyone
// else, not found (comments.ServeFile: what keeps it from running).
func FileHandler(db *gorm.DB, dir string) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		id, err := uuid.Parse(r.PathValue("id"))
		u := user.GetCurrentUser(r)
		if err != nil || u == nil || u.Anonymous() {
			http.NotFound(w, r)
			return
		}
		var f File
		if err := db.First(&f, "id = ?", id).Error; err != nil || !canSee(db, u.GetID(), f.MessageID) {
			http.NotFound(w, r)
			return
		}
		comments.ServeFile(w, r, dir, comments.StoredFile{ID: f.ID, Name: f.Name, ContentType: f.ContentType, Size: f.Size, Path: f.Path}, f.Data)
	})
}
