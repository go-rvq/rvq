package models

import "github.com/go-rvq/rvq/admin/origin"

// Origin is where a revision was made from (origin.Origin): the author's
// address and browser and — when the application locates addresses
// (origin.SetFunc) — where the address is. Zero when nothing is known: a
// revision of the boot's seed, one recorded with no request.
type Origin = origin.Origin
