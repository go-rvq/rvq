package js

import "embed"

//go:embed vuetify/dist
var Vuetify embed.FS

//go:embed vuetify/runtime-dist
var VuetifyRuntime embed.FS
