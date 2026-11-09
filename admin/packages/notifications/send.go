package notifications

import (
	"context"
	"errors"
	"fmt"
	"reflect"
	"slices"
	"strings"

	h "github.com/go-rvq/htmlgo"
	"github.com/go-rvq/rvq/admin/packages/comments"
	"github.com/go-rvq/rvq/admin/packages/filecards"
	"github.com/go-rvq/rvq/admin/packages/user"
	"github.com/go-rvq/rvq/admin/presets"
	"github.com/go-rvq/rvq/web"
	v "github.com/go-rvq/rvq/x/ui/vuetify"
	vx "github.com/go-rvq/rvq/x/ui/vuetifyx"
	"github.com/google/uuid"
)

// ActionSend is the action of the listing that sends a notification to other
// users (Builder.Users).
const ActionSend = "send_notification"

// sendLimit is how many users the dialog of the sending offers.
const sendLimit = 1000

// The fields of the dialog of the sending.
const (
	sendTo    = "Recipients"
	sendTitle = "Title"
	sendBody  = "Body"
	sendLevel = "Level"
	sendFiles = "Attachments"
)

// configureSend puts in the listing the action that sends a notification to
// other users: who may list the users (b.Users) chooses them; the text as
// written, never translated (SourceUser of who sends).
func (b *Builder) configureSend(lb *presets.ListingBuilder) {
	users := b.Users
	lb.Action(ActionSend).Icon("mdi-send-outline").
		SetI18nLabel(func(ctx context.Context) string { return GetMessages(ctx).Send }).
		SetEnabled(func(_ string, ctx *web.EventContext) (bool, error) {
			return currentUser(ctx) != uuid.Nil && users.Permissioner().ReqLister(ctx.R).IsAllowed() == nil, nil
		}).
		ComponentFunc(func(_ string, ctx *web.EventContext) (h.HTMLComponent, error) {
			m := GetMessages(ctx.Context())
			type item struct {
				Title string `json:"title"`
				Value string `json:"value"`
			}
			var items []item
			for _, u := range b.users() {
				title := u.GetName()
				if acc := u.GetAccountName(); acc != "" && acc != title {
					title += " (" + acc + ")"
				}
				items = append(items, item{Title: title, Value: u.GetID().String()})
			}
			cm := comments.GetMessages(ctx.Context())
			var levelChips []h.HTMLComponent
			for _, l := range Levels {
				levelChips = append(levelChips, v.VChip(h.Text(LevelLabel(ctx.Context(), l))).Value(l).
					// its color chosen and not (a chip of a group: base-color)
					Color(LevelColor(l)).Attr("base-color", LevelColor(l)).PrependIcon(LevelIcon(l)).
					Attr(":variant", fmt.Sprintf(`form[%q] === %q ? "flat" : "outlined"`, sendLevel, l)))
			}
			return h.Div(
				v.VAutocomplete().Label(m.SendTo+" *").Items(items).ItemTitle("title").ItemValue("value").
					Multiple(true).Chips(true).ClosableChips(true).
					Attr(web.VField(sendTo, []string{})...).Attr("data-send-to", true),
				v.VTextField().Label(m.NotificationTitle+" *").Attr(web.VField(sendTitle, "")...),
				// rich text, as the body of a page; files
				h.Div(vx.VXTipTapEditor().Label(m.NotificationBody+" *").Output("html").
					Attr(web.VField(sendBody, "")...)).Class("mb-2").Attr("data-send-body", true),
				h.Div(h.Text(m.NotificationAttachments)).Class("text-caption text-medium-emphasis"),
				h.Div(filecards.Drop(fmt.Sprintf("form[%q]", sendFiles), cm.DropFiles, cm.EditFile, cm.RemoveFile)).
					Attr("data-send-files", true),
				// the importance: a chip a level, its icon and its color, one
				// chosen (filled)
				h.Div(h.Text(m.NotificationLevel+" *")).Class("text-caption text-medium-emphasis mt-2"),
				v.VChipGroup(levelChips...).Mandatory(true).Attr(web.VField(sendLevel, LevelInfo)...).
					Attr("data-send-level", true),
			), nil
		}).
		UpdateFunc(func(_ string, ctx *web.EventContext) error {
			m := GetMessages(ctx.Context())
			from := currentUser(ctx)
			if from == uuid.Nil || users.Permissioner().ReqLister(ctx.R).IsAllowed() != nil {
				return errors.New(m.SendDenied)
			}
			// the form read (multipart) before its lists
			_ = ctx.R.FormValue(sendTitle)
			// the users: a list (Recipients[0], Recipients[1]…)
			var to []uuid.UUID
			for k, vs := range ctx.R.Form {
				if k != sendTo && !strings.HasPrefix(k, sendTo+"[") {
					continue
				}
				for _, s := range vs {
					if id, err := uuid.Parse(strings.TrimSpace(s)); err == nil {
						to = append(to, id)
					}
				}
			}
			// every field required (the level, one of them)
			title, body := strings.TrimSpace(ctx.R.FormValue(sendTitle)), strings.TrimSpace(ctx.R.FormValue(sendBody))
			level := ctx.R.FormValue(sendLevel)
			if len(to) == 0 || title == "" || comments.PlainText(body) == "" || !slices.Contains(Levels, level) {
				return errors.New(m.SendRequired)
			}
			files, err := comments.Attachments(ctx.R, sendFiles)
			if err != nil {
				return err
			}
			if err := b.conv.Validate(ctx, &comments.Write{Body: body, Files: files}); err != nil {
				return err
			}
			if err := b.Deliver(ctx.R, from, to, Notification{Title: title, Body: body, HTML: true, Level: level}, files...); err != nil {
				return err
			}
			ctx.Flash = m.Sent
			return nil
		})
}

// users are the users the sending offers, by name.
func (b *Builder) users() (out []user.User) {
	slice := b.Users.NewModelSlice()
	if err := b.db.Order("name").Limit(sendLimit).Find(slice).Error; err != nil {
		return nil
	}
	rv := reflect.ValueOf(slice).Elem()
	for i := 0; i < rv.Len(); i++ {
		e := rv.Index(i)
		if e.Kind() != reflect.Pointer {
			e = e.Addr()
		}
		if u, ok := e.Interface().(user.User); ok {
			out = append(out, u)
		}
	}
	return
}

// existing are the ids of ids that are users (no Users: ids).
func (b *Builder) existing(ids []uuid.UUID) []uuid.UUID {
	if b.Users == nil {
		return ids
	}
	var out []uuid.UUID
	b.db.Model(b.Users.NewModel()).Where("id IN ?", ids).Pluck("id", &out)
	return out
}
