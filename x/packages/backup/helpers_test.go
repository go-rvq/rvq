package backup

import (
	"compress/gzip"
	"io"
	"os"
)

func gzipWriter(w io.Writer) *gzip.Writer { return gzip.NewWriter(w) }

// writeOriginsAppend adds an origin to the file of a state's origins.
func writeOriginsAppend(name, p, src string) error {
	o, err := readOrigins(name)
	if err != nil {
		return err
	}
	o[p] = src
	if err := os.Remove(name); err != nil {
		return err
	}
	return writeOrigins(name, o)
}
