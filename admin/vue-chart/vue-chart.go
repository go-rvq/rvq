package vue_chart

import (
	"github.com/go-rvq/rvq/admin/presets"
	vue_chart "github.com/go-rvq/rvq/x/ui/vue-chart"
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
