package presets

import (
	"net/http"
	"net/url"
	"strconv"
	"testing"

	"github.com/go-rvq/rvq/web"
	"github.com/go-rvq/rvq/x/i18n"
)

func TestIsMenuItemActive(t *testing.T) {
	cases := []struct {
		// path means current url path
		path string
		// link means menu item link
		link string

		excepted bool
	}{
		{"", "/", true},
		{"/", "/", true},
		{"/", "/order", false},
		{"/order", "/order", true},
		{"/order/1", "/order", true},
		{"/order#", "/order", true},
		{"/product", "/order", false},
		{"/product", "/", false},
	}

	type io struct {
		ctx      *web.EventContext
		m        *ModelBuilder
		excepted bool
	}

	var toIO []io
	b := New(i18n.New())
	for _, c := range cases {
		// the menu item follows the model's link when it has one
		m := NewModelBuilder(b, &struct{}{})
		m.link = c.link

		toIO = append(toIO, io{
			ctx: &web.EventContext{
				R: &http.Request{
					URL: &url.URL{
						Path: c.path,
					},
				},
			},
			m:        m,
			excepted: c.excepted,
		})
	}

	for i, io := range toIO {
		t.Run(strconv.Itoa(i), func(t *testing.T) {
			if got := io.m.isMenuItemActive(io.ctx); got != io.excepted {
				t.Errorf("isMenuItemActive() = %v, excepted %v", got, io.excepted)
			}
		})
	}
}
