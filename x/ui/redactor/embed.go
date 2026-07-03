package redactor

import (
	"bytes"

	rvqjs "github.com/go-rvq/rvq/js"
	"github.com/go-rvq/rvq/web"
)

func JSComponentsPack() web.ComponentsPack {
	var js [][]byte
	j1, err := rvqjs.Redactor.ReadFile("redactor/dist/redactorjs.umd.cjs")
	if err != nil {
		panic(err)
	}
	js = append(js, j1)
	return web.ComponentsPack(bytes.Join(js, []byte("\n\n")))
}

func CSSComponentsPack() web.ComponentsPack {
	c, err := rvqjs.Redactor.ReadFile("redactor/dist/redactor.css")
	if err != nil {
		panic(err)
	}
	return web.ComponentsPack(c)
}
