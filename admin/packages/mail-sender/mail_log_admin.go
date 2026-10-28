package mail_sender

import (
	"context"
	"errors"
	"fmt"
	"net/url"
	"strings"

	h "github.com/go-rvq/htmlgo"
	"github.com/go-rvq/rvq/admin/presets"
	"github.com/go-rvq/rvq/admin/worker"
	"github.com/go-rvq/rvq/web"
	v "github.com/go-rvq/rvq/x/ui/vuetify"
	vx "github.com/go-rvq/rvq/x/ui/vuetifyx"
	"gorm.io/gorm"
)

// MailLogModelID is the id of the model of the central of the mails sent:
// its permissions (admin:mail_logs:…), its documents, its item of the menu.
const MailLogModelID = "mail_logs"

// JobMailLogsPurge deletes, every day, the MailLogs older than the
// retention of the settings (MailSender.MailLogRetentionDays).
const JobMailLogsPurge = "MailLogs:Purge"

// installMailLogs is the central of the mails sent, a child of the settings
// of the sending (sender: a singleton, so no id of it in its address —
// <sender>/mail-logs —, its node in the menu under the sender's): its
// listing, the newest first, all or the failed ones, by kind, period,
// recipient and status; the detail of each one — its body (redacted) in a
// sandbox, the link to the record it is of, the other messages of its
// conversation —; sending again one that had nothing hidden. Nobody edits,
// creates or deletes them.
func (b *Builder) installMailLogs(p *presets.Builder, sender *presets.ModelBuilder) {
	mb := NewModel(p, &MailLog{}, presets.ModelConfig().SetId(MailLogModelID)).
		URIName("mail-logs").
		MenuIcon("mdi-email-search-outline").
		GlobalSearch(false)
	mb.SetEditingDisabled(true).SetCreatingDisabled(true).SetDeletingDisabled(true)
	sender.AddChild(mb)
	b.logs = mb
	// in the menu, under the sender's node: to who may list them
	sender.MenuChildren(func(ctx *web.EventContext) ([]*presets.MenuNode, string) {
		if mb.Permissioner().Lister(ctx.R).Denied() {
			return nil, ""
		}
		href := mb.Info().ListingHref()
		node := &presets.MenuNode{Title: mb.TTitlePlural(ctx.Context()), Value: href,
			Props: map[string]any{"href": href, "prependIcon": "mdi-email-search-outline"}}
		active := ""
		if strings.HasPrefix(strings.TrimSuffix(ctx.R.URL.Path, "/"), href) {
			active = href
		}
		return []*presets.MenuNode{node}, active
	})

	l := mb.Listing("CreatedAt", "Status", "Kind", "Recipients", "Subject", "Method").
		SearchColumns("subject", "recipients", "record_path", "message_id").
		OrderBy("created_at DESC")
	l.Field("CreatedAt").ComponentFunc(func(field *presets.FieldContext, ctx *web.EventContext) h.HTMLComponent {
		return h.Td(h.Text(field.Obj.(*MailLog).CreatedAt.Format("2006-01-02 15:04:05")))
	})
	l.Field("Status").ComponentFunc(func(field *presets.FieldContext, ctx *web.EventContext) h.HTMLComponent {
		return h.Td(statusChip(field.Obj.(*MailLog), ctx))
	})

	l.FilterTabsFunc(func(ctx *web.EventContext) []*presets.FilterTab {
		msgr := GetMessages(ctx.Context())
		return []*presets.FilterTab{
			{Label: msgr.MailLogTabAll, Query: url.Values{"status": []string{}}},
			{Label: msgr.MailLogTabFailed, Query: url.Values{"status": []string{MailLogFailed}}},
		}
	})
	l.FilterDataFunc(func(ctx *web.EventContext) vx.FilterData {
		msgr := GetMessages(ctx.Context())
		var kinds []string
		b.db.Model(&MailLog{}).Distinct("kind").Where("kind <> ''").Order("kind").Pluck("kind", &kinds)
		var kindOptions []*vx.SelectItem
		for _, k := range kinds {
			kindOptions = append(kindOptions, &vx.SelectItem{Text: k, Value: k})
		}
		return []*vx.FilterItem{
			{Key: "kind", Label: msgr.MailLogKind, ItemType: vx.ItemTypeSelect, SQLCondition: `kind %s ?`, Options: kindOptions},
			{Key: "created", Label: msgr.MailLogCreatedAt, ItemType: vx.ItemTypeDatetimeRange, SQLCondition: `created_at %s ?`},
			{Key: "recipients", Label: msgr.MailLogRecipients, ItemType: vx.ItemTypeString, SQLCondition: `recipients %s ?`},
			{Key: "status", Label: msgr.MailLogStatus, ItemType: vx.ItemTypeSelect, SQLCondition: `status %s ?`,
				Options: []*vx.SelectItem{
					{Text: msgr.MailLogStatusSent, Value: MailLogSent},
					{Text: msgr.MailLogStatusFailed, Value: MailLogFailed},
				}},
		}
	})

	d := mb.Detailing("Status", "CreatedAt", "Kind", "Origin", "Recipients", "ReplyTo", "Subject", "Method", "Error",
		"MessageID", "InReplyTo", "UserID", "RemoteIP", "UserAgent", "LocaleCode", "Body", "Conversation")
	d.Field("Status").ComponentFunc(func(field *presets.FieldContext, ctx *web.EventContext) h.HTMLComponent {
		return vx.VXReadonlyField(statusChip(field.Obj.(*MailLog), ctx)).Label(field.Label)
	})
	d.Field("Error").ComponentFunc(func(field *presets.FieldContext, ctx *web.EventContext) h.HTMLComponent {
		if e := field.Obj.(*MailLog).Error; e != "" {
			return v.VAlert(h.Text(e)).Type("error").Variant(v.VariantTonal).Density(v.DensityCompact).Class("mb-4")
		}
		return nil
	})
	// the record it is of, by its record path
	d.Field("Origin").ComponentFunc(func(field *presets.FieldContext, ctx *web.EventContext) h.HTMLComponent {
		rp := field.Obj.(*MailLog).RecordPath
		if rp == "" {
			return nil
		}
		var link h.HTMLComponent = h.Text(rp)
		if href, err := p.RecordPathHref(rp); err == nil {
			link = v.VBtn(GetMessages(ctx.Context()).MailLogOpenOrigin).Href(href).
				PrependIcon("mdi-open-in-new").Variant(v.VariantTonal).Size(v.SizeSmall)
		}
		return vx.VXReadonlyField(link).Label(field.Label)
	})
	// its body, as it went (redacted), out of the page: a sandbox
	d.Field("Body").ComponentFunc(func(field *presets.FieldContext, ctx *web.EventContext) h.HTMLComponent {
		l := field.Obj.(*MailLog)
		msgr := GetMessages(ctx.Context())
		if l.BodyHidden {
			return vx.VXReadonlyField(h.I(msgr.MailLogBodyHidden)).Label(field.Label)
		}
		return vx.VXReadonlyField(
			h.Tag("iframe").Attr("sandbox", "").Attr("srcdoc", l.Body).
				Style("width:100%;min-height:360px;border:1px solid rgba(128,128,128,.3);border-radius:4px;background:#fff"),
		).Label(field.Label)
	})
	// the other messages of its conversation
	d.Field("Conversation").ComponentFunc(func(field *presets.FieldContext, ctx *web.EventContext) h.HTMLComponent {
		l := field.Obj.(*MailLog)
		id := l.ConversationID()
		if id == "" {
			return nil
		}
		var others []*MailLog
		b.db.Where("(message_id = ? OR in_reply_to = ?) AND id <> ?", id, id, l.ID).Order("created_at").Find(&others)
		if len(others) == 0 {
			return nil
		}
		var rows h.HTMLComponents
		for _, o := range others {
			rows = append(rows, h.Tr(
				h.Td(h.A(h.Text(o.CreatedAt.Format("2006-01-02 15:04:05"))).Href(mb.Info().DetailingHref(o.ID.String()))),
				h.Td(statusChip(o, ctx)),
				h.Td(h.Text(o.Kind)),
				h.Td(h.Text(o.Recipients)),
				h.Td(h.Text(o.Subject)),
			))
		}
		return vx.VXReadonlyField(v.VTable(h.Tbody(rows...)).Density(v.DensityCompact)).Label(field.Label)
	})

	// sent again, as it was: only one with nothing hidden
	d.Action("Resend").
		SetEnabledObj(func(obj any, id string, ctx *web.EventContext) (bool, error) {
			return !obj.(*MailLog).Redacted, nil
		}).
		UpdateFunc(func(id string, ctx *web.EventContext) error {
			var l MailLog
			if err := b.db.First(&l, "id = ?", id).Error; err != nil {
				return err
			}
			var settings MailSender
			if err := b.db.First(&settings).Error; err != nil {
				return fmt.Errorf("mail sender: no settings: %w", err)
			}
			if err := settings.Resend(&l); err != nil {
				return err
			}
			ctx.Flash = GetMessages(ctx.Context()).SendMailSuccessfully
			return nil
		})
}

// statusChip is the status of l: sent, or failed.
func statusChip(l *MailLog, ctx *web.EventContext) h.HTMLComponent {
	msgr := GetMessages(ctx.Context())
	if l.Status == MailLogFailed {
		return v.VChip(h.Text(msgr.MailLogStatusFailed)).Color("error").Size(v.SizeSmall).Variant(v.VariantTonal)
	}
	return v.VChip(h.Text(msgr.MailLogStatusSent)).Color("success").Size(v.SizeSmall).Variant(v.VariantTonal)
}

// PurgeJob is the job that deletes, every day, the MailLogs older than the
// retention of the settings (JobMailLogsPurge).
func (b *Builder) PurgeJob(w *worker.Builder) *Builder {
	w.NewJob(JobMailLogsPurge).
		CronConfig(worker.JobCronConfig{Spec: "@daily", Once: true}).
		Handler(func(ctx context.Context, job worker.JobInterface) error {
			var settings MailSender
			if err := b.db.First(&settings).Error; err != nil {
				if errors.Is(err, gorm.ErrRecordNotFound) {
					return nil
				}
				return err
			}
			n, err := PurgeMailLogs(b.db, settings.MailLogRetentionDays)
			if err == nil {
				job.AddLog(fmt.Sprintf("%d messages older than the retention deleted", n))
			}
			return err
		})
	return b
}
