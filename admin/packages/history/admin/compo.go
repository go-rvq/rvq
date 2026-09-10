package admin

import (
	"encoding/hex"

	histmodels "github.com/go-rvq/rvq/admin/packages/history/models"
)

// shortHash is the abbreviated hex form of a revision hash, for listings.
func shortHash(hash histmodels.Hash) string {
	s := hex.EncodeToString(hash)
	if len(s) > 14 {
		return s[:14]
	}
	return s
}
