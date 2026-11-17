package fs_tools

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/go-rvq/rvq/admin/packages/fs-tools/sandbox"
)

// termInfo is what the page is told of a session of the terminal, beside
// its state ({"info": …}): a chip (Label, Icon, Color) and, on a click, its
// details — in the user's language.
type termInfo struct {
	Label   string       `json:"label"`
	Icon    string       `json:"icon"`
	Color   string       `json:"color"`
	Details []termDetail `json:"details"`
}

// termDetail is a line of the details: what (Label) and its values, each a
// line; Warn, a value that is a risk.
type termDetail struct {
	Label  string      `json:"label"`
	Values []termValue `json:"values"`
}

type termValue struct {
	Text string `json:"text"`
	Warn bool   `json:"warn,omitempty"`
}

// info is how t keeps the session s apart, said by m.
func (t *terminal) info(m *Messages, s *session) termInfo {
	c := sandbox.CoverageOf(t.mode)
	in := termInfo{Label: m.TermNotIsolated, Icon: "mdi-shield-off-outline", Color: "warning"}
	isolation := termDetail{Label: m.TermIsolation}
	switch t.mode {
	case sandbox.ModeLandlock:
		in.Label, in.Icon, in.Color = fmt.Sprintf(m.TermIsolated, "Landlock"), "mdi-shield-lock-outline", "info"
		covers := []string{m.TermCovFiles}
		if c.TCP {
			covers = append(covers, "TCP")
		}
		if c.Signals {
			covers = append(covers, m.TermCovSignals)
		}
		isolation.Values = append(isolation.Values, termValue{Text: fmt.Sprintf(m.TermLandlock, c.ABI, strings.Join(covers, ", "))})
		if !c.Signals {
			isolation.Values = append(isolation.Values, termValue{Text: m.TermSignalsWarn, Warn: true})
		}
	case sandbox.ModeBwrap:
		in.Label, in.Icon, in.Color = fmt.Sprintf(m.TermIsolated, "bubblewrap"), "mdi-shield-lock-outline", "info"
		isolation.Values = append(isolation.Values, termValue{Text: m.TermBwrap})
	default:
		isolation.Values = append(isolation.Values, termValue{Text: m.TermOffWarn, Warn: true})
	}
	in.Details = append(in.Details, isolation)

	folders := termDetail{Label: m.TermFolders}
	for _, f := range s.folders {
		mode := m.ReadOnly
		if f.Write {
			mode = m.ReadWrite
		}
		folders.Values = append(folders.Values, termValue{Text: "$FILES_ROOT/" + f.At + " — " + mode})
	}
	in.Details = append(in.Details, folders)

	// the network: on (fstools.terminal.network of the app.yml), or off — the
	// server's administrator turns it on
	network := termValue{Text: m.TermNetOff}
	if t.net {
		network = termValue{Text: m.TermNetOn}
	}
	in.Details = append(in.Details, termDetail{Label: m.TermNetwork, Values: []termValue{network}})

	var limits []string
	if t.limits.CPU > 0 {
		limits = append(limits, "CPU "+fmtDuration(t.limits.CPU))
	}
	if t.limits.Memory > 0 {
		limits = append(limits, m.TermMemory+" "+fmtSize(t.limits.Memory))
	}
	if t.limits.FileSize > 0 {
		limits = append(limits, m.TermFileSize+" "+fmtSize(t.limits.FileSize))
	}
	if t.limits.Processes > 0 {
		limits = append(limits, m.TermProcesses+" "+strconv.Itoa(t.limits.Processes))
	}
	if len(limits) == 0 {
		limits = append(limits, m.TermNone)
	}
	in.Details = append(in.Details,
		termDetail{Label: m.TermLimits, Values: []termValue{{Text: strings.Join(limits, " · ")}}},
		termDetail{Label: m.TermIdle, Values: []termValue{{Text: fmt.Sprintf(m.TermIdleCloses, fmtDuration(t.idle))}}},
		termDetail{Label: m.TermSession, Values: []termValue{{Text: m.TermSessionTemp}}},
	)
	return in
}

// fmtDuration is d short: "45 s", "10 min", "1 h 30 min".
func fmtDuration(d time.Duration) string {
	switch {
	case d < time.Minute:
		return strconv.Itoa(int(d/time.Second)) + " s"
	case d < time.Hour:
		return strconv.Itoa(int(d/time.Minute)) + " min"
	}
	h, min := int(d/time.Hour), int(d%time.Hour/time.Minute)
	if min == 0 {
		return strconv.Itoa(h) + " h"
	}
	return strconv.Itoa(h) + " h " + strconv.Itoa(min) + " min"
}

// fmtSize is n bytes short, of 1024: "512 MiB", "2 GiB".
func fmtSize(n int64) string {
	units := []string{"B", "KiB", "MiB", "GiB", "TiB"}
	i, v := 0, float64(n)
	for v >= 1024 && i < len(units)-1 {
		v /= 1024
		i++
	}
	if v == float64(int64(v)) {
		return strconv.FormatInt(int64(v), 10) + " " + units[i]
	}
	return strconv.FormatFloat(v, 'f', 1, 64) + " " + units[i]
}
