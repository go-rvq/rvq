package main

import (
	"github.com/go-rvq/rvq/admin/docs/docgo"
	"github.com/go-rvq/rvq/admin/docs/docsrc"
	"github.com/go-rvq/rvq/admin/docs/docsrc/assets"
)

func main() {
	docgo.New().
		Assets("/assets/", assets.Assets).
		MainPageTitle("RVQ Document").
		SitePrefix("/docs/").
		DocTree(docsrc.DocTree...).
		BuildStaticSite("../docs")
}
