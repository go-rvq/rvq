package backup

import (
	"context"
	"fmt"
	"strconv"
	"strings"
	"time"

	h "github.com/go-rvq/htmlgo"
	"github.com/go-rvq/rvq/admin/model"
	"github.com/go-rvq/rvq/admin/presets"
	"github.com/go-rvq/rvq/admin/presets/gorm2op"
	"github.com/go-rvq/rvq/web"
	v "github.com/go-rvq/rvq/x/ui/vuetify"
	vx "github.com/go-rvq/rvq/x/ui/vuetifyx"
)

// fill replaces the {name}s of s by kv (name, value, …).
func fill(s string, kv ...string) string {
	for i := 0; i+1 < len(kv); i += 2 {
		s = strings.ReplaceAll(s, "{"+kv[i]+"}", kv[i+1])
	}
	return s
}

// Size is n bytes, short (1000s: kB, MB, GB).
func Size(n int64) string {
	if n < 1000 {
		return strconv.FormatInt(n, 10) + " B"
	}
	u := []string{"kB", "MB", "GB", "TB"}
	f, i := float64(n)/1000, 0
	for f >= 1000 && i < len(u)-1 {
		f /= 1000
		i++
	}
	if f < 10 {
		return fmt.Sprintf("%.1f %s", f, u[i])
	}
	return fmt.Sprintf("%.0f %s", f, u[i])
}

// statusText is the word of a status.
func statusText(m *Messages, status string) string {
	switch status {
	case StatusQueued:
		return m.StatusQueued
	case StatusBuilding:
		return m.StatusBuilding
	case StatusReady:
		return m.StatusReady
	case StatusFailed:
		return m.StatusFailed
	case StatusCanceled:
		return m.StatusCanceled
	}
	return status
}

// statusChip is the status of s: a chip (its color, its icon); being made,
// its progress — and the page reloads itself while it is so (poller).
func statusChip(s *BackupState, ctx *web.EventContext) h.HTMLComponent {
	m := GetMessages(ctx.Context())
	color, icon := "grey", "mdi-help-circle-outline"
	switch s.Status {
	case StatusQueued:
		color, icon = "info", "mdi-timer-sand"
	case StatusBuilding:
		color, icon = "primary", "mdi-progress-upload"
	case StatusReady:
		color, icon = "success", "mdi-check-circle"
	case StatusFailed:
		color, icon = "error", "mdi-alert-circle"
	case StatusCanceled:
		color, icon = "grey", "mdi-cancel"
	}
	chips := h.HTMLComponents{v.VChip(h.Text(statusText(m, s.Status))).Color(color).PrependIcon(icon).
		Size(v.SizeSmall).Variant(v.VariantTonal).Attr("data-backup-status", s.Status)}
	if s.Partial {
		chips = append(chips, v.VChip(h.Text(m.PartialChip)).Color("warning").Size(v.SizeSmall).Variant(v.VariantOutlined).Class("ms-1"))
	}
	if s.Status == StatusBuilding && s.ProgressSource != "" {
		chips = append(chips, h.Span(fill(m.ProgressText, "source", s.ProgressSource,
			"files", strconv.Itoa(s.ProgressFiles), "size", Size(s.ProgressBytes))).Class("text-caption text-medium-emphasis ms-2"))
	}
	if s.InProgress() {
		chips = append(chips, poller())
	}
	return h.Span("").Children(chips...).Class("d-inline-flex align-center")
}

// poller reloads the page every 5 s while it is shown (a state queued or
// being made); one a page, whatever the rows that draw it; the page not
// seen, or a login up, not then.
func poller() h.HTMLComponent {
	return h.Tag("user-component").
		Attr(":setup", `[({inject, window}) => {
	if (window.__rvqBackupPoll) return;
	const plaid = inject('plaid');
	window.__rvqBackupPoll = window.setInterval(() => {
		if (window.document.hidden || (window.__goplaid && window.__goplaid.loginPending)) return;
		plaid().reload().go();
	}, 5000);
}]`).
		Attr("@before-unmount", "() => { window.clearInterval(window.__rvqBackupPoll); window.__rvqBackupPoll = 0 }").
		Attr("data-backup-poller", true)
}

func (b *Builder) configureModel(mb *presets.ModelBuilder) {
	l := mb.Listing("CreatedAt", "Description", "Status", "Origin", "CreatedByName", "Files", "Bytes", "Keep").
		OrderBy("created_at DESC").SearchColumns("description", "id")
	l.Field("Status").ComponentFunc(func(field *presets.FieldContext, ctx *web.EventContext) h.HTMLComponent {
		return h.Td(statusChip(field.Obj.(*BackupState), ctx))
	})
	l.Field("Origin").ComponentFunc(func(field *presets.FieldContext, ctx *web.EventContext) h.HTMLComponent {
		return h.Td(h.Text(originText(GetMessages(ctx.Context()), field.Obj.(*BackupState).Origin)))
	})
	l.Field("Bytes").ComponentFunc(func(field *presets.FieldContext, ctx *web.EventContext) h.HTMLComponent {
		s := field.Obj.(*BackupState)
		if s.Status != StatusReady {
			return h.Td()
		}
		return h.Td(h.Text(Size(s.Bytes)))
	})
	l.Field("Files").ComponentFunc(func(field *presets.FieldContext, ctx *web.EventContext) h.HTMLComponent {
		s := field.Obj.(*BackupState)
		if s.Status != StatusReady {
			return h.Td()
		}
		return h.Td(h.Text(strconv.Itoa(s.Files)))
	})
	l.Field("Keep").ComponentFunc(func(field *presets.FieldContext, ctx *web.EventContext) h.HTMLComponent {
		if !field.Obj.(*BackupState).Keep {
			return h.Td()
		}
		return h.Td(v.VIcon("mdi-pin").Size(v.SizeSmall).Attr("title", GetMessages(ctx.Context()).KeptChip))
	})

	b.configureForm(mb)
	b.configureDetailing(mb)

	// the permission of downloading the states: the listing's
	mb.AppendListPermActions(func() []*presets.PermNodeAction {
		return []*presets.PermNodeAction{{Name: PermDownload,
			Title:       func(ctx context.Context) string { return GetMessages(ctx).PermDownloadTitle },
			Description: func(ctx context.Context) string { return GetMessages(ctx).PermDownloadDesc }}}
	})

	// canceled: queued or being made
	cancel := mb.Detailing().Action("cancel").Icon("mdi-cancel").
		SetEnabledObj(func(obj any, id string, ctx *web.EventContext) (bool, error) {
			return obj.(*BackupState).InProgress(), nil
		}).
		UpdateFunc(func(id string, ctx *web.EventContext) error {
			if err := b.Cancel(id); err != nil {
				return err
			}
			ctx.Flash = GetMessages(ctx.Context()).Canceled
			return nil
		})
	mb.Listing().AppendItemAction(cancel)

	mb.Listing().DeleteFunc(func(obj any, id model.ID, cascade bool, ctx *web.EventContext) error {
		if err := b.Remove(id.String()); err != nil {
			if err == ErrInProgress {
				return fmt.Errorf("%s", GetMessages(ctx.Context()).ErrInProgress)
			}
			return err
		}
		return nil
	})

	// the sessions of a state: who downloaded it, from where, what
	sessions := presets.NewModelBuilder(b.p, &BackupSession{}, presets.ModelConfig().SetId("backup_sessions").SetModuleKey(MessagesKey))
	sessions.Label("BackupSession").SetPluralLabel("BackupSessions")
	mb.AddChild(sessions)
	b.configureSessions(sessions)
}

func originText(m *Messages, o string) string {
	if o == OriginScheduled {
		return m.OriginScheduled
	}
	return m.OriginManual
}

// configureForm: a new state asks its description, its sources (all by
// default), whether it is kept; then, the description and kept alone.
func (b *Builder) configureForm(mb *presets.ModelBuilder) {
	ed := mb.Editing("Description", "Sources", "Keep")
	ed.Field("Description").ComponentFunc(func(field *presets.FieldContext, ctx *web.EventContext) h.HTMLComponent {
		return v.VTextarea().Label(field.Label).Rows(2).AutoGrow(true).
			Hint(GetMessages(ctx.Context()).BackupStateDescription_Hint).PersistentHint(true).
			Attr(web.VField(field.FormKey, field.Obj.(*BackupState).Description)...).ErrorMessages(field.Errors...)
	})
	ed.Field("Sources").ComponentFunc(func(field *presets.FieldContext, ctx *web.EventContext) h.HTMLComponent {
		s := field.Obj.(*BackupState)
		all := b.core.Sources()
		values := s.Sources
		if len(values) == 0 {
			values = all
		}
		sel := v.VSelect().Label(field.Label).Multiple(true).Chips(true).Items(all).
			Hint(GetMessages(ctx.Context()).BackupStateSources_Hint).PersistentHint(true).
			Attr(web.VField(field.FormKey, values)...).ErrorMessages(field.Errors...).Attr("data-backup-sources", true)
		if s.ID != "" {
			sel.Disabled(true) // asked once
		}
		return sel
	}).SetterFunc(func(obj any, field *presets.FieldContext, ctx *web.EventContext) error {
		if s := obj.(*BackupState); s.ID == "" {
			s.Sources = ctx.FormSliceValues(field.FormKey)
		}
		return nil
	})
	ed.Field("Keep").ComponentFunc(func(field *presets.FieldContext, ctx *web.EventContext) h.HTMLComponent {
		return v.VCheckbox().Label(field.Label).Hint(GetMessages(ctx.Context()).BackupStateKeep_Hint).PersistentHint(true).
			Attr(web.VField(field.FormKey, field.Obj.(*BackupState).Keep)...)
	})
	ed.Validators.AppendFunc(func(obj any, _ presets.FieldModeStack, ctx *web.EventContext) (errs web.ValidationErrors) {
		if s := obj.(*BackupState); s.ID == "" && len(s.Sources) == 0 {
			errs.FieldError("Sources", GetMessages(ctx.Context()).ErrNoSource)
		}
		return
	})
	ed.CreateFunc(func(obj any, ctx *web.EventContext) error { return b.create(obj.(*BackupState), ctx) })
	ed.SaveFunc(func(obj any, id model.ID, ctx *web.EventContext) error {
		s := obj.(*BackupState)
		return b.db.Model(&BackupState{}).Where("id = ?", id.String()).
			Updates(map[string]any{"description": s.Description, "keep": s.Keep}).Error
	})
}

func (b *Builder) configureDetailing(mb *presets.ModelBuilder) {
	d := mb.Detailing("Status", "Description", "CreatedByName", "Origin", "CreatedAt", "StartedAt", "FinishedAt",
		"Error", "Files", "Bytes", "ManifestSum", "Results", "Download")
	ro := func(field *presets.FieldContext, c h.HTMLComponent) h.HTMLComponent {
		return vx.VXReadonlyField(c).Label(field.Label)
	}
	d.Field("Status").ComponentFunc(func(field *presets.FieldContext, ctx *web.EventContext) h.HTMLComponent {
		s := field.Obj.(*BackupState)
		c := h.Div(statusChip(s, ctx))
		if s.Partial {
			c.AppendChildren(h.Div(h.Text(GetMessages(ctx.Context()).PartialWarning)).Class("text-caption text-warning mt-1"))
		}
		return ro(field, c)
	})
	d.Field("Origin").ComponentFunc(func(field *presets.FieldContext, ctx *web.EventContext) h.HTMLComponent {
		return ro(field, h.Text(originText(GetMessages(ctx.Context()), field.Obj.(*BackupState).Origin)))
	})
	d.Field("Error").ComponentFunc(func(field *presets.FieldContext, ctx *web.EventContext) h.HTMLComponent {
		if e := field.Obj.(*BackupState).Error; e != "" {
			return v.VAlert(h.Text(e)).Type("error").Variant(v.VariantTonal).Density(v.DensityCompact).Class("mb-4")
		}
		return nil
	})
	for _, f := range []string{"CreatedAt", "StartedAt", "FinishedAt"} {
		name := f
		d.Field(name).ComponentFunc(func(field *presets.FieldContext, ctx *web.EventContext) h.HTMLComponent {
			s := field.Obj.(*BackupState)
			t := s.StartedAt
			switch name {
			case "CreatedAt":
				t = &s.CreatedAt
			case "FinishedAt":
				t = s.FinishedAt
			}
			if t == nil {
				return nil
			}
			return ro(field, h.Text(presets.InTimeLocation(ctx, *t).Format("2006-01-02 15:04:05")))
		})
	}
	d.Field("Bytes").ComponentFunc(func(field *presets.FieldContext, ctx *web.EventContext) h.HTMLComponent {
		s := field.Obj.(*BackupState)
		if s.Status != StatusReady {
			return nil
		}
		return ro(field, h.Text(Size(s.Bytes)))
	})
	d.Field("Files").ComponentFunc(func(field *presets.FieldContext, ctx *web.EventContext) h.HTMLComponent {
		s := field.Obj.(*BackupState)
		if s.Status != StatusReady {
			return nil
		}
		return ro(field, h.Text(fmt.Sprintf("%d (+ %d / %d)", s.Files, s.Dirs, s.Links)))
	})
	d.Field("ManifestSum").ComponentFunc(func(field *presets.FieldContext, ctx *web.EventContext) h.HTMLComponent {
		s := field.Obj.(*BackupState)
		if s.ManifestSum == "" {
			return nil
		}
		return ro(field, h.Code(s.ManifestSum))
	})
	d.Field("Results").ComponentFunc(func(field *presets.FieldContext, ctx *web.EventContext) h.HTMLComponent {
		s := field.Obj.(*BackupState)
		if len(s.Results) == 0 {
			return nil
		}
		m := GetMessages(ctx.Context())
		rows := h.HTMLComponents{}
		for _, r := range s.Results {
			var warns h.HTMLComponent
			if len(r.Warnings) > 0 {
				items := h.HTMLComponents{}
				for _, w := range r.Warnings {
					items = append(items, h.Li(h.Text(w)))
				}
				warns = h.Ul(items...).Class("ps-4 text-warning")
			}
			rows = append(rows, h.Tr(h.Td(h.Code(r.Name)), h.Td(h.Text(strconv.Itoa(r.Files))), h.Td(h.Text(Size(r.Bytes))),
				h.Td(h.Text(r.Duration.Round(time.Millisecond).String())), h.Td(warns)))
		}
		return ro(field, v.VTable(
			h.Thead(h.Tr(h.Th(m.ColSource), h.Th(m.ColFiles), h.Th(m.ColSize), h.Th(m.ColDuration), h.Th(m.ColWarnings))),
			h.Tbody(rows...),
		).Density(v.DensityCompact).Attr("data-backup-results", true))
	})
	d.Field("Download").ComponentFunc(func(field *presets.FieldContext, ctx *web.EventContext) h.HTMLComponent {
		return ro(field, b.downloadSection(field.Obj.(*BackupState), ctx))
	})
}

// downloadSection is how a state is downloaded: ready, the manifest and the
// client's address; else why not.
func (b *Builder) downloadSection(s *BackupState, ctx *web.EventContext) h.HTMLComponent {
	m := GetMessages(ctx.Context())
	switch {
	case s.InProgress():
		return h.Div(h.Text(fill(m.NotReady, "status", statusText(m, s.Status)))).Class("text-medium-emphasis")
	case s.Status != StatusReady:
		return h.Div(h.Text(m.NotDownloadable)).Class("text-medium-emphasis")
	}
	if b.mb.Permissioner().ReqListDo(ctx.R, PermDownload).Denied() {
		return nil
	}
	manifest := strings.TrimSuffix(b.p.GetURIPrefix(), "/") + "/!backup-manifest/" + s.ID
	return h.Div(
		v.VBtn(m.ManifestLink).Href(manifest).PrependIcon("mdi-download").Variant(v.VariantTonal).Size(v.SizeSmall).
			Attr("download", true).Attr("data-backup-manifest", true),
		h.Div(h.Strong(m.StateIDLabel+": "), h.Code(s.ID)).Class("mt-3"),
		h.Div(h.Strong(m.ClientURL+": "), h.Code(clientBase(ctx.R)+b.ClientPath())).Class("mt-1"),
		h.Div(h.Text(m.KeyNeeded)).Class("mt-2 text-caption"),
	).Attr("data-backup-download", true)
}

// configureSessions: the history of the downloads of a state, read only.
func (b *Builder) configureSessions(mb *presets.ModelBuilder) {
	// of the state above: its id the last of the parents
	mb.UpdateDataOperator(func(do presets.DataOperator) presets.DataOperator {
		return do.(*gorm2op.DataOperatorBuilder).
			WithReadCallbacks(func(cb *gorm2op.Callbacks[*gorm2op.DataOperatorBuilder]) {
				cb.Pre(func(state *gorm2op.CallbackState) error {
					ids := presets.ParentsModelID(state.Ctx.R)
					if len(ids) == 0 {
						state.DB = state.DB.Where("1 = 0")
						return nil
					}
					state.DB = state.DB.Where("state_id = ?", fmt.Sprint(ids.Last().Value()))
					return nil
				})
			})
	})
	mb.Listing("CreatedAt", "UserName", "IP", "UserAgent", "ExpiresAt", "Files", "Bytes").OrderBy("created_at DESC")
	mb.Listing().Field("Bytes").ComponentFunc(func(field *presets.FieldContext, ctx *web.EventContext) h.HTMLComponent {
		return h.Td(h.Text(Size(field.Obj.(*BackupSession).Bytes)))
	})
	mb.Detailing("CreatedAt", "UserName", "IP", "UserAgent", "ExpiresAt", "Files", "Bytes")
	mb.SetCreatingDisabled(true).SetEditingDisabled(true).SetDeletingDisabled(true)
}
