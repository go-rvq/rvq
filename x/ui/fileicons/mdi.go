package fileicons

import "strings"

// The icon of a file by its extension, of Material Design Icons (mdi-…): the
// one of a link, a button, a card of the file. One not known, mdi-file.

var mdiByExt = map[string]string{}

func init() {
	for icon, exts := range map[string][]string{
		"mdi-file-pdf-box":           {"pdf"},
		"mdi-file-word":              {"doc", "docx", "odt", "rtf"},
		"mdi-file-excel":             {"xls", "xlsx", "ods"},
		"mdi-file-delimited":         {"csv", "tsv"},
		"mdi-file-powerpoint":        {"ppt", "pptx", "odp"},
		"mdi-folder-zip":             {"zip", "rar", "7z", "gz", "tgz", "bz2", "xz", "tar"},
		"mdi-file-image":             {"jpg", "jpeg", "png", "gif", "webp", "svg", "bmp", "tif", "tiff", "heic", "avif", "ico"},
		"mdi-file-video":             {"mp4", "webm", "mov", "avi", "mkv", "wmv", "m4v", "ogv"},
		"mdi-file-music":             {"mp3", "wav", "ogg", "oga", "flac", "m4a", "aac", "opus"},
		"mdi-file-document-outline":  {"txt", "md", "log"},
		"mdi-code-json":              {"json", "xml", "yaml", "yml"},
		"mdi-file-code":              {"html", "htm", "css", "js", "ts", "go", "py", "sql"},
		"mdi-book-open-page-variant": {"epub", "mobi"},
	} {
		for _, ext := range exts {
			mdiByExt[ext] = icon
		}
	}
}

// MDI is the icon (mdi-…) of a file of the extension ext ("pdf", ".PDF").
func MDI(ext string) string {
	if icon, ok := mdiByExt[strings.ToLower(strings.TrimPrefix(ext, "."))]; ok {
		return icon
	}
	return "mdi-file"
}
