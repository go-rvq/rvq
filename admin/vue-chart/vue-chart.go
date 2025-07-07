package vue_chart

import (
	"github.com/qor5/admin/v3/presets"
	vue_chart "github.com/qor5/x/v3/ui/vue-chart"
)

func New() *Buider {
	return &Buider{}
}

type Buider struct {
}

func (b *Buider) Install(pb *presets.Builder) error {
	pb.ExtraAsset("/vue-chart.js", "text/javascript", vue_chart.JSComponentsPack())
	return nil
}
