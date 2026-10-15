package js

import "embed"

// Portal is the bundle of the portals of a site (portal/dist:
// rvq-portal.umd.js): served by js/api/portal's Server.ScriptHandler.
//
//go:embed portal/dist
var Portal embed.FS
