package fileicons

import "testing"

func TestMDI(t *testing.T) {
	for ext, want := range map[string]string{"pdf": "mdi-file-pdf-box", ".XLSX": "mdi-file-excel", "zip": "mdi-folder-zip", "csv": "mdi-file-delimited", "": "mdi-file", "unknown": "mdi-file"} {
		if got := MDI(ext); got != want {
			t.Errorf("MDI(%q) = %q, want %q", ext, got, want)
		}
	}
}
