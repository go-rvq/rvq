package publish

import (
	"errors"
	"fmt"
	"reflect"
	"sync"

	h "github.com/go-rvq/htmlgo"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"

	"github.com/go-rvq/rvq/admin/presets"
	"github.com/go-rvq/rvq/web"
	. "github.com/go-rvq/rvq/x/ui/vuetify"
)

func DraftCountComponentFunc(db *gorm.DB) presets.FieldComponentFunc {
	return func(field *presets.FieldContext, ctx *web.EventContext) h.HTMLComponent {
		var (
			count int64
			obj   = field.Obj
		)
		modelSchema, err := schema.Parse(obj, &sync.Map{}, db.NamingStrategy)
		if err != nil {
			return h.Td(h.Text("0"))
		}
		setPrimaryKeysConditionWithoutVersion(db.Model(reflect.New(modelSchema.ModelType).Interface()), obj, modelSchema).
			Where("status = ?", StatusDraft).Count(&count)

		return h.Td(h.Text(fmt.Sprint(count)))
	}
}

func LiveComponentFunc(db *gorm.DB, b *LiveChipsBuilder) presets.FieldComponentFunc {
	return presets.ReadOnlyFieldComponentFuncWrapper(LiveFieldComponentFunc(db, b))
}

func LiveFieldComponentFunc(db *gorm.DB, b *LiveChipsBuilder) presets.FieldComponentFunc {
	return LiveFieldComponentFuncBuilder(db, func(status string, toStatus string, msgr *Messages) h.HTMLComponent {
		return b.Auto(status, toStatus, msgr)
	})
}

func LiveFieldComponentFuncBuilder(db *gorm.DB, build func(status string, toStatus string, msgr *Messages) h.HTMLComponent) presets.FieldComponentFunc {
	return func(field *presets.FieldContext, ctx *web.EventContext) (comp h.HTMLComponent) {
		msgr := GetMessages(ctx.Context())

		var (
			ok            bool
			err           error
			modelSchema   *schema.Schema
			scheduleStart Schedule

			obj = field.Obj
		)

		defer func() {
			if err != nil {
				comp = h.Text("-")
				return
			}
		}()
		if modelSchema, err = schema.Parse(obj, &sync.Map{}, db.NamingStrategy); err != nil {
			return
		}

		var (
			g = func() *gorm.DB {
				return setPrimaryKeysConditionWithoutFields(db.Model(reflect.New(modelSchema.ModelType).Interface()), obj, modelSchema, "Version", "LocaleCode")
			}
			nowTime = db.NowFunc()
		)
		st, ok := obj.(StatusInterface)
		if !ok {
			err = errors.New("ErrorModel")
			return
		}

		sc, ok := obj.(ScheduleInterface)
		if !ok {
			return build(st.EmbedStatus().Status, "", msgr)
		}

		var (
			statusFieldName = modelSchema.FieldsByName["Status"].DBName
			startFieldName  = modelSchema.FieldsByName["ScheduledStartAt"].DBName
		)

		var toStatus string
		if st.EmbedStatus().Status != StatusOnline {
			if sc.EmbedSchedule().ScheduledStartAt != nil {
				toStatus = StatusOnline
			}
		} else {
			err := g().Select(startFieldName).Where(fmt.Sprintf("%s <> ? AND %s > ?", statusFieldName, startFieldName), StatusOnline, nowTime).Order(startFieldName).Limit(1).Scan(&scheduleStart).Error
			if err != nil {
				return
			}
			currentEndAt := sc.EmbedSchedule().ScheduledEndAt
			if scheduleStart.ScheduledStartAt != nil && (currentEndAt == nil || !scheduleStart.ScheduledStartAt.After(*currentEndAt)) {
				toStatus = "+1"
			} else if currentEndAt != nil && !currentEndAt.Before(nowTime) {
				toStatus = StatusOffline
			}
		}

		return build(st.EmbedStatus().Status, toStatus, msgr)
	}
}

func StatusListFunc(b *LiveChipsBuilder) presets.FieldComponentFunc {
	return func(field *presets.FieldContext, ctx *web.EventContext) h.HTMLComponent {
		msgr := GetMessages(ctx.Context())

		if s, ok := field.Obj.(StatusInterface); ok {
			return h.Td(b.Status(s.EmbedStatus().Status, msgr))
		}
		return nil
	}
}

type LiveChipsBuilder struct {
	WithLabel bool
	WithBg    bool
}

var (
	LiveChipsFormBuilder = LiveChipsBuilder{WithLabel: true, WithBg: true}
	LiveChipsListBuilder LiveChipsBuilder
)

func (b *LiveChipsBuilder) LiveIcons(status string, isScheduled bool, msgr *Messages) (comps h.HTMLComponents) {
	var (
		label, color = GetStatusLabelColor(status, msgr)
		i            = VIcon("mdi-radiobox-marked").Title(label).Color(color)
	)

	comps = h.HTMLComponents{i}

	if isScheduled {
		i := VIcon("mdi-menu-right")
		if b.WithLabel {
			i.Class("ml-1")
		}
		if b.WithBg {
			i.Size(SizeSmall)
		}
		comps = append(comps, i)
	}

	return
}

func (b *LiveChipsBuilder) Live(status string, isScheduled bool, msgr *Messages) *VChipBuilder {
	var (
		label, color = GetStatusLabelColor(status, msgr)
		comps        h.HTMLComponents
	)

	if status == StatusOnline {
		i := VIcon("mdi-radiobox-marked")
		if b.WithBg {
			i.Size(SizeSmall)
		}
		if b.WithLabel {
			i.Class("mr-1")
		}
		comps = append(comps, i)
	}

	if b.WithLabel {
		comps = append(comps, h.Span(label))
	}

	if isScheduled {
		i := VIcon("mdi-menu-right")
		if b.WithLabel {
			i.Class("ml-1")
		}
		if b.WithBg {
			i.Size(SizeSmall)
		}
		comps = append(comps, i)
	}

	chip := VChip(comps...).Color(color).Density(DensityCompact).Tile(true).Class("px-1")

	if !isScheduled {
		return chip
	}

	return chip.Class("rounded-s-lg")
}

func (b *LiveChipsBuilder) Status(status string, msgr *Messages) *VChipBuilder {
	return b.Live(status, false, msgr).Class("rounded")
}

func (b *LiveChipsBuilder) Lives(status string, toStatus string, msgr *Messages) h.HTMLComponent {
	if toStatus != "" {
		return h.Components(
			b.Live(status, true, msgr).Class("rounded-s"),
			b.Live(toStatus, false, msgr).Class("rounded-e"),
		)
	}
	return b.Status(status, msgr).Class("rounded")
}

func (b *LiveChipsBuilder) LivesIcon(status string, toStatus string, msgr *Messages) h.HTMLComponents {
	if toStatus != "" {
		return append(
			b.LiveIcons(status, true, msgr),
			b.LiveIcons(toStatus, false, msgr)...,
		)
	}
	return b.LiveIcons(status, false, msgr)
}

func (b *LiveChipsBuilder) Auto(status string, toStatus string, msgr *Messages) h.HTMLComponent {
	if b.WithBg {
		return b.Lives(status, toStatus, msgr)
	}
	return b.LivesIcon(status, toStatus, msgr)
}

func GetStatusLabelColor(status string, msgr *Messages) (label, color string) {
	switch status {
	case StatusOnline:
		return msgr.StatusOnline, ColorSuccess
	case StatusOffline:
		return msgr.StatusOffline, "#9E9E9E"
	case StatusDraft:
		return msgr.StatusDraft, ColorWarning
	default:
		return status, ColorSuccess
	}
}
