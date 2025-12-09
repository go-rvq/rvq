package corejs

import (
	"fmt"
	"strings"
)

const componentsRegistrator = `
window.__goplaidVueComponentRegisters = window.__goplaidVueComponentRegisters || [];
window.__goplaidVueComponentRegisters.push(%s)
`

func RegisterComponent(handlers ...string) string {
	return fmt.Sprintf(componentsRegistrator, strings.Join(handlers, ","))
}

func RegisterComponentWithDefaults(handlers ...string) string {
	return RegisterComponent(append([]string{
		`(app) => app.use(VueI18n.createI18n({}))`,
	}, handlers...)...)
}
