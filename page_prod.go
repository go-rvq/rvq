//go:build prod

package web

import (
	"fmt"
	"os"
)

func init() {
	if exe, _ := os.Executable(); exe != "" {
		if s, _ := os.Stat(exe); s != nil {
			ExeMTime = fmt.Sprint(s.ModTime().Unix())
		}
	}
}
