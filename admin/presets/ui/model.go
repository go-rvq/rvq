package ui

import (
	"fmt"

	"github.com/go-rvq/htmlgo"
	"github.com/go-rvq/rvq/admin/model"
	"github.com/go-rvq/rvq/admin/presets"
	"github.com/go-rvq/rvq/admin/presets/actions"
	"github.com/go-rvq/rvq/web"
)

type DetailOpenOptionsBuilder struct {
	ID        model.ID
	ParentsID model.IDSlice
	ModelInfo *presets.ModelInfo
	Dialog    *web.VueEventTagBuilder
}

func OpenDetailOptions(modelInfo *presets.ModelInfo, id model.ID, parentsID ...model.ID) *DetailOpenOptionsBuilder {
	return &DetailOpenOptionsBuilder{
		ID:        id,
		ParentsID: parentsID,
		ModelInfo: modelInfo,
		Dialog: web.Plaid().
			EventFunc(actions.Detailing).
			URL(modelInfo.ListingHref(parentsID...)).
			Query(presets.ParamID, id.String()),
	}
}

func (b *DetailOpenOptionsBuilder) ConfigureDialog(do func(p *web.VueEventTagBuilder)) *DetailOpenOptionsBuilder {
	do(b.Dialog)
	return b
}

func (b *DetailOpenOptionsBuilder) DialogPortal(name string) *DetailOpenOptionsBuilder {
	b.Dialog.Query(presets.ParamTargetPortal, name)
	return b
}

func (b *DetailOpenOptionsBuilder) DialogOverlay(overlayMode actions.OverlayMode) *DetailOpenOptionsBuilder {
	b.Dialog.Query(presets.ParamOverlay, overlayMode)
	return b
}

func (b *DetailOpenOptionsBuilder) DialogRenderBreadcrumbs() *DetailOpenOptionsBuilder {
	b.Dialog.Query(presets.ParamRenderBreadcrumbs, "true")
	return b
}

func (b *DetailOpenOptionsBuilder) BuildTargetLink() string {
	return b.ModelInfo.DetailingHref(b.ID, b.ParentsID...)
}

func (b *DetailOpenOptionsBuilder) BuildDialogEvent() string {
	return b.Dialog.Go()
}

func OpenDetail[T htmlgo.TagGetter](b *DetailOpenOptionsBuilder, tag T) T {
	t := tag.GetHTMLTagBuilder()
	t.Attr("@click", b.BuildDialogEvent())
	t.Attr("@click.middle",
		fmt.Sprintf(`(e) => e.view.window.open(%q, "_blank")`, b.BuildTargetLink()))
	return tag
}
