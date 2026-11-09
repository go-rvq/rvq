package notifications

import (
	"fmt"
	"mime"
	"net/http"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/go-rvq/rvq/admin/packages/comments"
	"github.com/go-rvq/rvq/admin/packages/user"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

// FileMaxSize is the largest attachment taken, in bytes.
const FileMaxSize = comments.FileMaxSize

// File is an attachment of a message (MessageID): of each copy of it, seen
// by who received it and by who sent it (FileHandler).
type File struct {
	ID          uuid.UUID `gorm:"type:uuid;primaryKey"`
	MessageID   uuid.UUID `gorm:"type:uuid;not null;index"`
	Name        string    `gorm:"size:255;not null"`
	ContentType string    `gorm:"size:255"`
	Size        int64
	Data        []byte `gorm:"not null"`
	CreatedAt   time.Time
}

func (File) TableName() string { return "user_notification_files" }

// AttachFiles keeps the files as the attachments of the message id.
func AttachFiles(db *gorm.DB, message uuid.UUID, files ...comments.Attachment) error {
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

// FilesOf are the attachments of the messages (no data), by message.
func FilesOf(db *gorm.DB, messages ...uuid.UUID) (map[uuid.UUID][]File, error) {
	out := map[uuid.UUID][]File{}
	if len(messages) == 0 {
		return out, nil
	}
	var fs []File
	if err := db.Select("id", "message_id", "name", "content_type", "size", "created_at").
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

// inlineTypes are the types an attachment is shown in the browser as.
var inlineTypes = map[string]bool{
	"image/png": true, "image/jpeg": true, "image/gif": true, "image/webp": true, "application/pdf": true,
}

// filesPath is where the attachments are served, under the admin's prefix.
const filesPath = "/notifications/files/"

// FileHandler serves the attachment of its id (the end of the path) to who
// received or sent its message; to anyone else, not found.
func FileHandler(db *gorm.DB) http.Handler {
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
		ct := f.ContentType
		if ct == "" {
			ct = "application/octet-stream"
		}
		// a user's file: shown in the browser only when it can run nothing
		// (an SVG can: downloaded)
		disposition := "attachment"
		if inlineTypes[strings.ToLower(strings.TrimSpace(strings.Split(ct, ";")[0]))] {
			disposition = "inline"
		} else {
			ct = "application/octet-stream"
		}
		w.Header().Set("Content-Type", ct)
		w.Header().Set("Content-Length", strconv.Itoa(len(f.Data)))
		w.Header().Set("Content-Disposition", mime.FormatMediaType(disposition, map[string]string{"filename": f.Name}))
		w.Header().Set("X-Content-Type-Options", "nosniff")
		w.Header().Set("Content-Security-Policy", "default-src 'none'; img-src 'self'; style-src 'unsafe-inline'; sandbox")
		w.Header().Set("Cache-Control", "private, max-age=3600")
		_, _ = w.Write(f.Data)
	})
}
