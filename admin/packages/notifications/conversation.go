package notifications

import (
	"context"
	"log"
	"net/http"
	"strings"

	"github.com/go-rvq/rvq/admin/packages/comments"
	"github.com/go-rvq/rvq/web"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

// eventConversation is the event of the conversation of a notification (rvq
// comments).
const eventConversation = "notifications_conversation"

// Mail is a message of a user mailed to who received it (Builder.Mailer):
// From sent it to To; Copy is To's notification (its thread), Path its
// address in the admin; R the request it was sent in (nil: by mail).
type Mail struct {
	From, To uuid.UUID
	Copy     Notification
	Path     string
	Files    []comments.Attachment
	R        *http.Request
}

// Deliver sends the message n of the user from to the users to — rich text
// when n.HTML —, with files: a notification each, one message (MessageID),
// of its conversation (n.ThreadID; a new one, its own); mailed to each when
// the Mailer is set (a failure is logged: the notification stays).
func (b *Builder) Deliver(r *http.Request, from uuid.UUID, to []uuid.UUID, n Notification, files ...comments.Attachment) error {
	n.Source, n.System = SourceUser(from), false
	n.Module, n.TitleKey, n.BodyKey, n.TitleArgs, n.BodyArgs = "", "", "", nil, nil
	n.MessageID = uuid.New()
	if n.ThreadID == uuid.Nil {
		n.ThreadID = n.MessageID
	}
	if n.Level == "" {
		n.Level = LevelInfo
	}
	to = b.existing(to)
	err := b.db.Transaction(func(tx *gorm.DB) error {
		if err := Notify(tx, n, to...); err != nil {
			return err
		}
		return AttachFiles(tx, n.MessageID, files...)
	})
	if err != nil || b.Mailer == nil {
		return err
	}
	var copies []Notification
	b.db.Where("message_id = ?", n.MessageID).Find(&copies)
	ctx := context.Background()
	if r != nil {
		ctx = r.Context()
	}
	for _, c := range copies {
		if err := b.Mailer(ctx, &Mail{From: from, To: c.UserID, Copy: c, Path: b.listingHref() + "/" + c.ID.String(), Files: files, R: r}); err != nil {
			log.Printf("notifications: the mail of %s to %s: %v", c.ID, c.UserID, err)
		}
	}
	return nil
}

// Reply is an answer to a conversation (Answer): the message it answers, its
// text — rich text when HTML —, its files, where it came from out of the
// admin (ExternalID).
type Reply struct {
	// ParentID is the message it answers (none: the newest the user
	// received of the conversation)
	ParentID   uuid.UUID
	Body       string
	HTML       bool
	ExternalID string
	Files      []comments.Attachment
}

// Answer is the answer rep of the user from to the conversation thread: to
// everyone else of it (who received or sent one of its messages), under its
// title after "Re: ". Who is of none of it answers nothing
// (gorm.ErrRecordNotFound). An answer by mail too (the app's inbox).
func (b *Builder) Answer(r *http.Request, from, thread uuid.UUID, rep Reply) error {
	var first Notification
	if err := b.db.Where("thread_id = ? AND NOT system", thread).Order("created_at").First(&first).Error; err != nil {
		return err
	}
	if !b.inThread(from, thread) {
		return gorm.ErrRecordNotFound
	}
	to := b.participants(thread, from)
	if len(to) == 0 {
		return gorm.ErrRecordNotFound
	}
	title := first.Title
	if !strings.HasPrefix(strings.ToLower(title), "re:") {
		title = "Re: " + title
	}
	// the message it answers: of the conversation; none, the newest the
	// user received of it
	parent := rep.ParentID
	var n int64
	if parent != uuid.Nil {
		b.db.Model(&Notification{}).Where("thread_id = ? AND message_id = ?", thread, parent).Count(&n)
	}
	if n == 0 {
		var last Notification
		if err := b.db.Where("thread_id = ? AND user_id = ?", thread, from).Order("created_at DESC").First(&last).Error; err == nil {
			parent = last.MessageID
		} else {
			parent = first.MessageID
		}
	}
	return b.Deliver(r, from, to, Notification{Title: title, Body: rep.Body, HTML: rep.HTML, Level: first.Level,
		Link: first.Link, ThreadID: thread, ParentID: parent, ExternalID: rep.ExternalID}, rep.Files...)
}

// Answered says whether an answer came from externalID (a Message-ID).
func (b *Builder) Answered(externalID string) bool {
	var n int64
	b.db.Unscoped().Model(&Notification{}).Where("external_id = ?", externalID).Count(&n)
	return externalID != "" && n > 0
}

// inThread says whether the user received or sent a message of the thread.
func (b *Builder) inThread(uid, thread uuid.UUID) bool {
	var n int64
	b.db.Model(&Notification{}).Where("thread_id = ? AND (user_id = ? OR source = ?)", thread, uid, SourceUser(uid)).Count(&n)
	return n > 0
}

// participants are who received or sent a message of the thread, but but.
func (b *Builder) participants(thread, but uuid.UUID) (out []uuid.UUID) {
	var rows []Notification
	b.db.Unscoped().Select("user_id", "source").Where("thread_id = ?", thread).Find(&rows)
	seen := map[uuid.UUID]bool{but: true, uuid.Nil: true}
	add := func(id uuid.UUID) {
		if !seen[id] {
			seen[id] = true
			out = append(out, id)
		}
	}
	for _, r := range rows {
		add(r.UserID)
		if s, ok := strings.CutPrefix(r.Source, sourceUserPrefix); ok {
			if id, err := uuid.Parse(s); err == nil {
				add(id)
			}
		}
	}
	return
}

// configureConversation makes the conversation of a notification of a user
// (rvq comments): its messages the user received or sent, the oldest
// first — who sent each, its text, its files —, and the box that answers it.
func (b *Builder) configureConversation() *comments.Builder {
	cb := comments.New(b.p, eventConversation)
	cb.RichText = true
	// the files of an answer, where the app lets them (ReplyFiles)
	cb.Attachments = func(ctx *web.EventContext, _ string) bool {
		return b.ReplyFiles == nil || b.ReplyFiles(ctx)
	}
	// the key: the notification of the user (its thread)
	mine := func(ctx *web.EventContext, key string) (*Notification, bool) {
		uid := currentUser(ctx)
		var n Notification
		if uid == uuid.Nil || b.db.Where("id = ? AND user_id = ?", key, uid).First(&n).Error != nil {
			return nil, false
		}
		return &n, true
	}
	// the checks of the app's (ValidateFile, ValidateMessage), when it
	// gives them — after Install too
	cb.ValidateFile = func(ctx *web.EventContext, key string, f *comments.Attachment) error {
		if b.ValidateFile == nil {
			return nil
		}
		return b.ValidateFile(ctx, key, f)
	}
	cb.ValidateMessage = func(ctx *web.EventContext, w *comments.Write) error {
		if b.ValidateMessage == nil {
			return nil
		}
		return b.ValidateMessage(ctx, w)
	}
	cb.CanRead = func(ctx *web.EventContext, key string) bool {
		_, ok := mine(ctx, key)
		return ok
	}
	// a user's (a system's is answered by no one)
	cb.CanWrite = func(ctx *web.EventContext, key string) bool {
		n, ok := mine(ctx, key)
		return ok && !n.System
	}
	cb.Comments = func(ctx *web.EventContext, key string) ([]*comments.Comment, error) {
		n, ok := mine(ctx, key)
		if !ok {
			return nil, gorm.ErrRecordNotFound
		}
		uid := currentUser(ctx)
		var rows []Notification
		if err := b.db.Where("thread_id = ? AND (user_id = ? OR source = ?)", n.ThreadID, uid, SourceUser(uid)).
			Order("created_at").Find(&rows).Error; err != nil {
			return nil, err
		}
		// a message once (its copies, one each who received it)
		seen := map[uuid.UUID]bool{}
		var ids []uuid.UUID
		var msgs []Notification
		for _, r := range rows {
			if seen[r.MessageID] {
				continue
			}
			seen[r.MessageID] = true
			ids = append(ids, r.MessageID)
			msgs = append(msgs, r)
		}
		files, err := FilesOf(b.db, ids...)
		if err != nil {
			return nil, err
		}
		prefix := strings.TrimSuffix(b.p.GetURIPrefix(), "/")
		// a tree: each answer under the message it answers
		var out []*comments.Comment
		byMessage := map[uuid.UUID]*comments.Comment{}
		for _, m := range msgs {
			title, body := m.Text(ctx.Context())
			c := &comments.Comment{ID: m.MessageID.String(), Time: m.CreatedAt, Title: title, Body: body, HTML: m.HTML}
			if m.System {
				c.Author = GetMessages(ctx.Context()).OriginSystem
			} else if s, ok := strings.CutPrefix(m.Source, sourceUserPrefix); ok {
				c.Author = s
				if u := b.user(s); u != nil {
					c.Author, c.AuthorHint = u.GetAccountName(), u.GetName()
				}
			}
			for _, f := range files[m.MessageID] {
				c.Files = append(c.Files, comments.File{Name: f.Name, Size: f.Size, Type: f.ContentType, URL: prefix + filesPath + f.ID.String()})
			}
			byMessage[m.MessageID] = c
			if p := byMessage[m.ParentID]; p != nil && m.ParentID != m.MessageID {
				p.Replies = append(p.Replies, c)
			} else {
				out = append(out, c)
			}
		}
		return out, nil
	}
	cb.OnAdd = func(ctx *web.EventContext, w *comments.Write) error {
		n, ok := mine(ctx, w.Key)
		if !ok || n.System {
			return comments.ErrNotAllowed
		}
		// an answer to a message of it; the box on the top, to the one seen
		parent := n.MessageID
		if id, err := uuid.Parse(w.ParentID); err == nil {
			parent = id
		}
		return b.Answer(ctx.R, currentUser(ctx), n.ThreadID, Reply{ParentID: parent, Body: w.Body, HTML: true, Files: w.Files})
	}
	cb.Register(b.p.GetWebBuilder())
	return cb
}
