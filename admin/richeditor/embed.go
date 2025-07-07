package richeditor

import (
	"bytes"
	"embed"
	"fmt"

	"github.com/go-rvq/rvq/web"
	"github.com/go-rvq/rvq/x/ui/redactor"
)

//go:embed redactor
var box embed.FS

func JSComponentsPack() web.ComponentsPack {
	var js [][]byte
	js = append(js, []byte(redactor.JSComponentsPack()))
	for _, p := range Plugins {
		pj, err := box.ReadFile(fmt.Sprintf("redactor/%s.min.js", p))
		if err != nil {
			continue
		}
		js = append(js, pj)
	}
	if len(PluginsJS) > 0 {
		js = append(js, PluginsJS...)
	}

	return web.ComponentsPack(bytes.Join(js, []byte("\n\n")))
}

func CSSComponentsPack() web.ComponentsPack {
	var css [][]byte
	css = append(css, []byte(redactor.CSSComponentsPack()))
	custom, err := box.ReadFile("redactor/redactor.custom.css")
	if err != nil {
		panic(err)
	}
	css = append(css, custom)
	if len(PluginsCSS) > 0 {
		css = append(css, PluginsCSS...)
	}
	return web.ComponentsPack(bytes.Join(css, []byte("\n\n")))
}
