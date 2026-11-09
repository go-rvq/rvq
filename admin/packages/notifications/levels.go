package notifications

import (
	"context"
	"fmt"
	"slices"

	h "github.com/go-rvq/htmlgo"
	"github.com/go-rvq/rvq/admin/dashboard"
	v "github.com/go-rvq/rvq/x/ui/vuetify"
)

// LevelColor is the color of the level in the palette of the theme (info,
// success, warning, error); an unknown one, info's.
func LevelColor(level string) string {
	if slices.Contains(Levels, level) {
		return level
	}
	return LevelInfo
}

// LevelIcon is the icon of the level.
func LevelIcon(level string) string {
	switch level {
	case LevelSuccess:
		return "mdi-check-circle-outline"
	case LevelWarning:
		return "mdi-alert-outline"
	case LevelError:
		return "mdi-alert-octagon-outline"
	}
	return "mdi-information-outline"
}

// LevelLabel is the name of the level in the language of ctx.
func LevelLabel(ctx context.Context, level string) string {
	m := GetMessages(ctx)
	switch level {
	case LevelSuccess:
		return m.LevelSuccess
	case LevelWarning:
		return m.LevelWarning
	case LevelError:
		return m.LevelError
	}
	return m.LevelInfo
}

// levelChip is the level as a chip of its color, with its icon.
func levelChip(ctx context.Context, level string) h.HTMLComponent {
	return v.VChip(h.Text(LevelLabel(ctx, level))).Color(LevelColor(level)).Size(v.SizeSmall).
		Variant(v.VariantTonal).PrependIcon(LevelIcon(level)).Attr("data-notification-level", LevelColor(level))
}

// levelIcon is the icon of the level in its color, its name the tooltip.
func levelIcon(ctx context.Context, level string) h.HTMLComponent {
	label := LevelLabel(ctx, level)
	return v.VIcon(LevelIcon(level)).Color(LevelColor(level)).Attr("aria-label", label).
		Children(v.VTooltip().Activator("parent").Location("top").Text(label))
}

// levelRowStyle is the background of a row of the level: its color of the
// theme, light.
func levelRowStyle(level string) string {
	return fmt.Sprintf("background-color: rgba(var(--v-theme-%s), 0.08)", LevelColor(level))
}

// dashboardLevel is the level as the dashboard says it.
func dashboardLevel(level string) dashboard.Level {
	return dashboard.Level(LevelColor(level))
}
