package login

import (
	"embed"
)

//go:embed assets
var assetsFS embed.FS

var assetsPathPrefix = "/auth/assets/"
var (
	StyleCSSURL = assetsPathPrefix + "style.css"
	ZxcvbnJSURL = assetsPathPrefix + "zxcvbn.js"
)

// setAssetsPrefix puts the login's files under prefix (Builder.URIPrefix):
// a process has one login, so the package's URLs follow it.
func setAssetsPrefix(prefix string) {
	assetsPathPrefix = prefix + "/auth/assets/"
	StyleCSSURL = assetsPathPrefix + "style.css"
	ZxcvbnJSURL = assetsPathPrefix + "zxcvbn.js"
}
