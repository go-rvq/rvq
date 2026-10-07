package vuetifyx

import (
	"encoding/json"
	"strings"
	"unicode/utf8"

	h "github.com/go-rvq/htmlgo"
	"github.com/go-rvq/rvq/web"
	"github.com/sergi/go-diff/diffmatchpatch"
)

// LineChanges compares two texts as IntelliJ does: the 1-based lines removed
// from old and added to new, and, for a line changed (a removal paired with
// the addition after it), the exact spans changed in it — on the old side
// (delRanges, keyed by old line) and the new (insRanges, keyed by new line),
// each [start, end, hunk] in runes, hunk -1. What <vx-code>'s diff marks
// (MarkLines, ChangeRanges) take.
func LineChanges(oldStr, newStr string) (removed, added []int, delRanges, insRanges map[int][][3]int) {
	delRanges = map[int][][3]int{}
	insRanges = map[int][][3]int{}
	// a last line with no "\n" is the same line as one with it
	if oldStr != "" && !strings.HasSuffix(oldStr, "\n") {
		oldStr += "\n"
	}
	if newStr != "" && !strings.HasSuffix(newStr, "\n") {
		newStr += "\n"
	}
	splitLines := func(s string) []string {
		s = strings.TrimSuffix(s, "\n")
		if s == "" {
			return nil
		}
		return strings.Split(s, "\n")
	}
	diffs := LineDiff(oldStr, newStr)
	oldLine, newLine := 0, 0
	for i := 0; i < len(diffs); i++ {
		d := diffs[i]
		switch d.Type {
		case diffmatchpatch.DiffEqual:
			n := countLines(d.Text)
			oldLine += n
			newLine += n
		case diffmatchpatch.DiffDelete:
			dels := splitLines(d.Text)
			var ins []string
			if i+1 < len(diffs) && diffs[i+1].Type == diffmatchpatch.DiffInsert {
				ins = splitLines(diffs[i+1].Text)
				i++ // the insert paired with it
			}
			for k, ln := range dels {
				oldLine++
				removed = append(removed, oldLine)
				if k < len(ins) {
					if dr, _ := charRanges(ln, ins[k]); len(dr) > 0 {
						delRanges[oldLine] = dr
					}
				}
			}
			for k, ln := range ins {
				newLine++
				added = append(added, newLine)
				if k < len(dels) {
					if _, ir := charRanges(dels[k], ln); len(ir) > 0 {
						insRanges[newLine] = ir
					}
				}
			}
		case diffmatchpatch.DiffInsert:
			for range splitLines(d.Text) {
				newLine++
				added = append(added, newLine)
			}
		}
	}
	return
}

// LineDiff is the diff of two texts by their lines — each Diff's Text whole
// lines, with their "\n" —: every line both have is paired (Myers' exact
// diff, no timeout). diffmatchpatch's own DiffLinesToChars writes the lines
// as numbers in a string, which its character diff then compares digit by
// digit: lines both texts have came out changed.
func LineDiff(a, b string) []diffmatchpatch.Diff {
	index := map[string]rune{}
	var lines []string
	encode := func(s string) []rune {
		var out []rune
		for len(s) > 0 {
			i := strings.IndexByte(s, '\n')
			line := s
			if i >= 0 {
				line = s[:i+1]
			}
			s = s[len(line):]
			r, ok := index[line]
			if !ok {
				// a rune per line, past the surrogates
				r = rune(0x10000 + len(lines))
				index[line] = r
				lines = append(lines, line)
			}
			out = append(out, r)
		}
		return out
	}
	ra, rb := encode(a), encode(b)
	dmp := diffmatchpatch.New()
	dmp.DiffTimeout = 0
	diffs := dmp.DiffMainRunes(ra, rb, false)
	for i, d := range diffs {
		var text strings.Builder
		for _, r := range d.Text {
			text.WriteString(lines[r-0x10000])
		}
		diffs[i].Text = text.String()
	}
	return diffs
}

// charRanges char-diffs two changed lines: the spans deleted (in the old
// line's runes) and inserted (in the new line's), each [start, end, -1].
func charRanges(oldLn, newLn string) (del, ins [][3]int) {
	dmp := diffmatchpatch.New()
	diffs := dmp.DiffMain(oldLn, newLn, false)
	dmp.DiffCleanupSemantic(diffs)
	oldOff, newOff := 0, 0
	for _, d := range diffs {
		n := utf8.RuneCountInString(d.Text)
		switch d.Type {
		case diffmatchpatch.DiffEqual:
			oldOff += n
			newOff += n
		case diffmatchpatch.DiffDelete:
			del = append(del, [3]int{oldOff, oldOff + n, -1})
			oldOff += n
		case diffmatchpatch.DiffInsert:
			ins = append(ins, [3]int{newOff, newOff + n, -1})
			newOff += n
		}
	}
	return
}

// countLines counts the lines of a chunk of a line diff (the last possibly
// with no "\n").
func countLines(s string) int {
	if s == "" {
		return 0
	}
	n := strings.Count(s, "\n")
	if !strings.HasSuffix(s, "\n") {
		n++
	}
	return n
}

// DiffFile is a file of a <vx-diff-browser>: its path, its status (git's:
// M, A, D, R, ?), its language (Prism's), its two versions and how they
// differ (NewDiffFile computes it).
type DiffFile struct {
	Path string `json:"path"`
	// From is a renamed (moved) file's path before: Old is its content
	// there, New here.
	From     string `json:"from,omitempty"`
	Status   string `json:"status,omitempty"`
	Language string `json:"language,omitempty"`
	Old      string `json:"old,omitempty"`
	New      string `json:"new,omitempty"`
	Binary   bool   `json:"binary,omitempty"`
	// ReadOnly: not to be edited in the browser, though it saves (Save)
	ReadOnly  bool             `json:"readOnly,omitempty"`
	Removed   []int            `json:"removed,omitempty"`
	Added     []int            `json:"added,omitempty"`
	DelRanges map[int][][3]int `json:"delRanges,omitempty"`
	InsRanges map[int][][3]int `json:"insRanges,omitempty"`
}

// NewDiffFile is the file path as it was (old) and is (new), compared
// (LineChanges). Content with a NUL byte is binary: not compared, not sent.
// A renamed file's path before is From (set after).
func NewDiffFile(path, status, language, old, new string) DiffFile {
	f := DiffFile{Path: path, Status: status, Language: language}
	if strings.IndexByte(old, 0) >= 0 || strings.IndexByte(new, 0) >= 0 {
		f.Binary = true
		return f
	}
	f.Old, f.New = old, new
	f.Removed, f.Added, f.DelRanges, f.InsRanges = LineChanges(old, new)
	return f
}

// DiffSummary is a file of a <vx-diff-browser> as its tree has it — its path,
// status, path before (a renamed one's) and language —, its contents asked
// for when its tab opens (VXDiffBrowserBuilder.Load).
func DiffSummary(path, status, from, language string) DiffFile {
	return DiffFile{Path: path, Status: status, From: from, Language: language}
}

// DiffContentScript is the script an event of Load answers with: the file's
// diff (NewDiffFile) set on the tab's content.
func DiffContentScript(f DiffFile) string {
	b, err := json.Marshal(f)
	if err != nil {
		return DiffContentErrorScript(err.Error())
	}
	return "content.value = " + string(b)
}

// DiffSavedScript is the script an event of Save answers with when it saved
// the file.
func DiffSavedScript() string { return "content.saved = true" }

// DiffContentErrorScript is the script an event of Load answers with when it
// has no diff to give: what the tab says.
func DiffContentErrorScript(message string) string {
	b, _ := json.Marshal(message)
	return "content.error = " + string(b)
}

// VXDiffBrowserBuilder renders <vx-diff-browser>: the changes of a set of
// files, browsed — a tree of the files on the left (their folders, their
// paths), a tab on the right for each file opened from it, the old and the
// current side by side, compared as IntelliJ does. Its panels are dockview's.
type VXDiffBrowserBuilder struct {
	tag *h.HTMLTagBuilder
}

// VXDiffBrowser browses files: complete (NewDiffFile), or summaries
// (DiffSummary) whose contents Load asks for.
func VXDiffBrowser(files ...DiffFile) *VXDiffBrowserBuilder {
	if files == nil {
		files = []DiffFile{}
	}
	return &VXDiffBrowserBuilder{tag: h.Tag("vx-diff-browser").Attr(":files", files)}
}

// Load is how a file's diff is asked for, when its tab opens: the event,
// called with the tab's content in its scope ({content}) and the file's
// path, from (a renamed one's path before) and status in its query; it
// answers with DiffContentScript (or DiffContentErrorScript) as its
// RunScript, which sets content.value. The browser keeps a file's diff only
// while its tab is open.
func (b *VXDiffBrowserBuilder) Load(event *web.VueEventTagBuilder) *VXDiffBrowserBuilder {
	call := event.Scope(web.Var("{content}")).
		Query("path", web.Var("content.file.path")).
		Query("from", web.Var("content.file.from || ''")).
		Query("status", web.Var("content.file.status || ''")).
		Go()
	b.tag.Attr(":load", "(content) => "+call)
	return b
}

// Save makes the current sides editable — the difference computed again as
// it is edited, a change of the old taken back by its button, undo and redo,
// a save button (and Ctrl+S) —, and is how a file is saved: the event, called
// with the save's state in its scope ({content}: content.value the text), the
// file's path in its query (path) and the text in its form (content); it
// answers with DiffSavedScript, or DiffContentErrorScript, as its RunScript.
// A file ReadOnly is not edited.
func (b *VXDiffBrowserBuilder) Save(event *web.VueEventTagBuilder) *VXDiffBrowserBuilder {
	call := event.Scope(web.Var("{content}")).
		Query("path", web.Var("path")).
		FieldValue("content", web.Var("content.value")).
		Go()
	b.tag.Attr(":save", "(path, content) => "+call)
	return b
}

// EditTexts are the words of the editing: undo, redo, save, saving, saved,
// not saved, and the title of a revert button.
func (b *VXDiffBrowserBuilder) EditTexts(undo, redo, save, saving, saved, unsaved, revert string) *VXDiffBrowserBuilder {
	b.tag.Attr("undo-text", undo, "redo-text", redo, "save-text", save, "saving-text", saving,
		"saved-text", saved, "unsaved-text", unsaved, "revert-text", revert)
	return b
}

// NavTexts are the titles of the buttons that go to the previous and the
// next change.
func (b *VXDiffBrowserBuilder) NavTexts(prev, next string) *VXDiffBrowserBuilder {
	b.tag.Attr("prev-text", prev, "next-text", next)
	return b
}

// LoadTexts are what a tab says while its diff is asked for, and when it
// could not be.
func (b *VXDiffBrowserBuilder) LoadTexts(loading, failed string) *VXDiffBrowserBuilder {
	b.tag.Attr("loading-text", loading, "load-error-text", failed)
	return b
}

// VModel binds the summary (the files of the tree) to expr, in place of
// the files given.
func (b *VXDiffBrowserBuilder) VModel(expr string) *VXDiffBrowserBuilder {
	b.tag.Attr("v-model", expr)
	return b
}

// Height is how tall it is (default "70vh"); its panels fill it.
func (b *VXDiffBrowserBuilder) Height(v string) *VXDiffBrowserBuilder {
	b.tag.Attr("height", v)
	return b
}

// Labels are its words: the tree's title, the two sides', what a binary
// file says, and what it says when there is no file.
func (b *VXDiffBrowserBuilder) Labels(files, old, new, binary, empty string) *VXDiffBrowserBuilder {
	b.tag.Attr("files-title", files, "old-label", old, "new-label", new, "binary-text", binary, "empty-text", empty)
	return b
}

// RenameLabels are the words of a renamed (moved) file: what precedes its
// new path, and what its tab says when its content did not change.
func (b *VXDiffBrowserBuilder) RenameLabels(renamedTo, unchanged string) *VXDiffBrowserBuilder {
	b.tag.Attr("renamed-text", renamedTo, "unchanged-text", unchanged)
	return b
}

// ActionsSlot is drawn by each file of the tree, its scope { file } — the
// DiffFile (file.path, file.status): a discard button, say.
func (b *VXDiffBrowserBuilder) ActionsSlot(children ...h.HTMLComponent) *VXDiffBrowserBuilder {
	b.tag.Children(web.Slot(children...).Name("actions").Scope("{ file }"))
	return b
}

func (b *VXDiffBrowserBuilder) Attr(vs ...interface{}) *VXDiffBrowserBuilder {
	b.tag.Attr(vs...)
	return b
}

func (b *VXDiffBrowserBuilder) Write(ctx *h.Context) error {
	return b.tag.Write(ctx)
}
