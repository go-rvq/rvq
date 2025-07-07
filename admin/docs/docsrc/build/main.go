package main

import (
	"github.com/go-rvq/rvq/admin/docs/docsrc"
	"github.com/go-rvq/rvq/admin/docs/docsrc/assets"
	"github.com/theplant/docgo"
)

func main() {
	docgo.New().
		Assets("/assets/", assets.Assets).
		MainPageTitle("QOR5 Document").
		SitePrefix("/docs/").
		DocTree(docsrc.DocTree...).
		BuildStaticSite("../docs")
}
