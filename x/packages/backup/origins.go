package backup

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"sort"
	"strconv"
	"strings"
)

// The origins of a state are its files not in its snapshot (another file
// system): their path in the state and where they are, a line each,
// quoted.

func writeOrigins(name string, origins map[string]string) error {
	f, err := os.OpenFile(name, os.O_CREATE|os.O_EXCL|os.O_WRONLY, 0o600)
	if err != nil {
		return err
	}
	keys := make([]string, 0, len(origins))
	for k := range origins {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	w := bufio.NewWriter(f)
	for _, k := range keys {
		fmt.Fprintf(w, "%s\t%s\n", strconv.Quote(k), strconv.Quote(origins[k]))
	}
	return errors.Join(w.Flush(), f.Close())
}

func readOrigins(name string) (map[string]string, error) {
	f, err := os.Open(name)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	out := map[string]string{}
	sc := bufio.NewScanner(f)
	sc.Buffer(make([]byte, 64<<10), 1<<20)
	for sc.Scan() {
		k, v, ok := strings.Cut(sc.Text(), "\t")
		if !ok {
			return nil, errors.New("backup: origins: a line with no tab")
		}
		var err error
		if k, err = strconv.Unquote(k); err == nil {
			v, err = strconv.Unquote(v)
		}
		if err != nil {
			return nil, fmt.Errorf("backup: origins: %w", err)
		}
		out[k] = v
	}
	return out, sc.Err()
}
