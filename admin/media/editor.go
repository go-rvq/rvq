package media

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"mime"
	"mime/multipart"
	"path"
	"strings"
	"time"

	h "github.com/go-rvq/htmlgo"
	"github.com/go-rvq/rvq/admin/media/base"
	"github.com/go-rvq/rvq/admin/media/media_library"
	"github.com/go-rvq/rvq/admin/presets"
	rvqjs "github.com/go-rvq/rvq/js"
	"github.com/go-rvq/rvq/web"
	"github.com/go-rvq/rvq/x/ui/cropper"
)

// The image editor (vx-image-editor, js/image-editor): the click on a
// thumbnail opens it on the original — "Edit image": Crop (fixed,
// perspective, free), Colors, Annotate —, the browser draws the result and
// sends it (saveImageEdit). The original never changes: kept apart
// (MediaLibraryStorage.Edit.Source) the first time, the edit beside it to
// open again; the result is the file, its sizes made from it. The fixed crop
// is still the size's (CropOptions), over the image edited.

const (
	editImageEvent     = "mediaLibrary_EditImage"
	saveImageEditEvent = "mediaLibrary_SaveImageEdit"
)

// CameraZIndex is the z-index of the dialog of the camera: over the dialogs
// it opens from (a VDialog's default, 2400, would tie with them); the image
// editor (3100) and the lightbox (3200) are over it.
const CameraZIndex = 3000

// sourceStyle is the style of the URL the original is kept at.
const sourceStyle = "source"

// imageEditorPack is a file of the editor's build (js/image-editor/dist).
func imageEditorPack(name string) web.ComponentsPack {
	v, err := rvqjs.ImageEditor.ReadFile("image-editor/dist/" + name)
	if err != nil {
		panic(err)
	}
	return web.ComponentsPack(v)
}

// editImageScript opens the editor on the file id of the media box of field,
// from its size thumb, and saves what it makes.
func editImageScript(field, id, thumb string, cfg *media_library.MediaBoxConfig) string {
	save := web.Plaid().
		EventFunc(saveImageEditEvent).
		Form(web.Var("{}")).
		Query("field", field).
		Query("id", id).
		Query("thumb", thumb).
		FieldValue("cfg", h.JSONString(cfg)).
		FieldValue(field+".Values", web.Var(fmt.Sprintf("form[%q]", field+".Values"))).
		FieldValue("Edit", web.Var("JSON.stringify(s.edit)")).
		FieldValue("CropOption", web.Var("s.cropOption ? JSON.stringify(s.cropOption) : ''")).
		FieldValue("Result", web.Var("s.file ? [s.file] : []")).
		FieldValue("Restore", web.Var("!s.file && d.edited ? '1' : ''")).
		Go()
	return web.Plaid().
		EventFunc(editImageEvent).
		Form(web.Var("{}")).
		Query("field", field).
		Query("id", id).
		Query("thumb", thumb).
		FieldValue("cfg", h.JSONString(cfg)).
		ThenScript(`const d = r && r.data; if (!d || !window.rvqImageEditor) return;
return window.rvqImageEditor.open(d).then(function(res) { if (res.action !== 'save') return; const s = res.save; ` + save + ` })`).
		Go()
}

// editImage answers what the editor opens with (r.Data): the original, the
// edit saved (its fixed crop: the size's), the ratio of the size.
func editImage(mb *Builder) web.EventFunc {
	return func(ctx *web.EventContext) (r web.EventResponse, err error) {
		_, id, thumb, cfg := getParams(mb, ctx)
		var m media_library.MediaLibrary
		if err = mb.db.First(&m, "id = ?", id).Error; err != nil {
			return
		}
		if !base.IsImageFormat(m.File.FileName) && !base.IsImageFormat(m.File.Url) {
			return
		}
		moption := m.GetMediaOption()
		var aspect float64
		if size := moption.Sizes[thumb]; size != nil && size.Width > 0 && size.Height > 0 {
			aspect = float64(size.Width) / float64(size.Height)
		}

		src, name := m.File.URL("original"), m.File.GetFileName()
		state := map[string]any{}
		edited := m.File.Edit != nil && m.File.Edit.Source != ""
		if edited {
			src = m.File.Edit.Source
			if m.File.Edit.SourceName != "" {
				name = m.File.Edit.SourceName
			}
			if len(m.File.Edit.State) > 0 {
				_ = json.Unmarshal(m.File.Edit.State, &state)
			}
		}
		// the fixed crop is the size's
		delete(state, "fixo")
		if co := moption.CropOptions[thumb]; co != nil && aspect > 0 {
			state["fixo"] = map[string]int{"x": co.X, "y": co.Y, "width": co.Width, "height": co.Height}
			if _, ok := state["mode"]; !ok {
				state["mode"] = "fixo"
			}
		}

		data := map[string]any{
			"src":    fmt.Sprintf("%s?%d", src, time.Now().UnixNano()),
			"name":   name,
			"type":   mime.TypeByExtension(strings.ToLower(path.Ext(src))),
			"aspect": aspect,
			"edited": edited,
		}
		if len(state) > 0 {
			data["edit"] = state
		}
		if cfg.EditTab != "" {
			data["tab"] = cfg.EditTab
		}
		if cfg.EditMode != "" {
			data["mode"] = cfg.EditMode
		}
		r.Data = data
		return
	}
}

// saveImageEdit saves what the editor made: the result (Result), the edit
// (Edit), the fixed crop of the size (CropOption); Restore, the original back
// (the edit undone: "Back to the original").
func saveImageEdit(b *Builder) web.EventFunc {
	return func(ctx *web.EventContext) (r web.EventResponse, err error) {
		db := b.db
		field, id, thumb, cfg := getParams(b, ctx)

		box := &media_library.MediaBox{}
		if err = box.Scan(ctx.R.FormValue(field + ".Values")); err != nil {
			return
		}
		var m media_library.MediaLibrary
		if err = db.First(&m, "id = ?", id).Error; err != nil {
			return
		}
		if err = b.updateDescIsAllowed(ctx.R, &m); err != nil {
			return
		}

		state := json.RawMessage(ctx.R.FormValue("Edit"))
		if len(state) > 0 && !json.Valid(state) {
			return r, fmt.Errorf("media: the edit is not JSON")
		}
		var result *multipart.FileHeader
		if ctx.R.MultipartForm != nil {
			if fs := ctx.R.MultipartForm.File["Result"]; len(fs) > 0 {
				result = fs[0]
			}
		}

		switch {
		case result != nil:
			err = b.saveEdited(&m, result, state)
		case ctx.R.FormValue("Restore") != "":
			err = b.restoreOriginal(&m)
		case m.File.Edit != nil && len(state) > 0:
			// the same image (only the size's crop): the edit kept for the
			// next opening
			m.File.Edit.State = state
			err = db.Save(&m).Error
		}
		if err != nil {
			presets.ShowMessage(&r, err.Error(), "error")
			return r, nil
		}

		if co := ctx.R.FormValue("CropOption"); co != "" && thumb != base.DefaultSizeKey {
			var v cropper.Value
			if err = json.Unmarshal([]byte(co), &v); err != nil {
				return
			}
			moption := m.GetMediaOption()
			if moption.CropOptions == nil {
				moption.CropOptions = map[string]*base.CropOption{}
			}
			moption.CropOptions[thumb] = &base.CropOption{X: int(v.X), Y: int(v.Y), Width: int(v.Width), Height: int(v.Height)}
			moption.Crop = true
			if err = m.ScanMediaOptions(moption); err != nil {
				return
			}
			if err = base.SaveUploadAndCropImage(&b.Config, db, &m); err != nil {
				presets.ShowMessage(&r, err.Error(), "error")
				return r, nil
			}
		}

		box.Url = m.File.Url
		box.FileName = m.File.FileName
		box.FileSizes = m.File.FileSizes
		box.Width = m.File.Width
		box.Height = m.File.Height
		r.UpdatePortal(
			mediaBoxThumbnailsPortalName(field),
			mediaBoxThumbnails(ctx, box, field, cfg, false, false),
		)
		return
	}
}

// keepSource keeps the original of m apart, once: the first edit.
func (b *Builder) keepSource(m *media_library.MediaLibrary) error {
	if m.File.Edit != nil && m.File.Edit.Source != "" {
		return nil
	}
	f, err := m.File.Retrieve(m.File.URL("original"))
	if err != nil {
		return fmt.Errorf("media: the original of %s: %w", m.ID, err)
	}
	defer f.Close()
	dst := m.File.URL(sourceStyle)
	if err = m.File.Store(dst, nil, f); err != nil {
		return fmt.Errorf("media: keeping the original of %s: %w", m.ID, err)
	}
	m.File.Edit = &media_library.ImageEdit{Source: dst, SourceName: m.File.GetFileName()}
	return nil
}

// saveEdited makes result the file of m — its sizes made from it, the crops
// of the sizes started over —, the original kept apart and state the edit.
func (b *Builder) saveEdited(m *media_library.MediaLibrary, result *multipart.FileHeader, state json.RawMessage) error {
	if err := b.keepSource(m); err != nil {
		return err
	}
	edit := m.File.Edit
	edit.State = state
	if err := m.File.Scan(result); err != nil {
		return err
	}
	m.File.Edit = edit
	m.File.CropOptions = map[string]*base.CropOption{}
	return base.SaveUploadAndCropImage(&b.Config, b.db, m)
}

// restoreOriginal makes the original (kept apart) the file of m again: the
// edit gone.
func (b *Builder) restoreOriginal(m *media_library.MediaLibrary) error {
	if m.File.Edit == nil || m.File.Edit.Source == "" {
		return nil
	}
	f, err := m.File.Retrieve(m.File.Edit.Source)
	if err != nil {
		return fmt.Errorf("media: the original of %s: %w", m.ID, err)
	}
	buf, err := io.ReadAll(f)
	f.Close()
	if err != nil {
		return err
	}
	name := m.File.Edit.SourceName
	if name == "" {
		name = path.Base(m.File.Edit.Source)
	}
	if err = m.File.Scan(base.NewFileHeader(name, func() (base.Reader, error) { return bytes.NewReader(buf), nil })); err != nil {
		return err
	}
	m.File.Edit = nil
	m.File.CropOptions = map[string]*base.CropOption{}
	return base.SaveUploadAndCropImage(&b.Config, b.db, m)
}

// keepUploadEdit keeps, with the result of an upload edited in the browser
// (m), its original and the edit (MediaBoxConfig.UploadEdit): reopened, the
// edit goes on from them.
func (b *Builder) keepUploadEdit(m *media_library.MediaLibrary, original *multipart.FileHeader, state json.RawMessage) error {
	f, err := original.Open()
	if err != nil {
		return err
	}
	defer f.Close()
	dst := m.File.URL(sourceStyle)
	if ext := path.Ext(original.Filename); ext != "" {
		dst = strings.TrimSuffix(dst, path.Ext(dst)) + strings.ToLower(ext)
	}
	if err = m.File.Store(dst, nil, f); err != nil {
		return err
	}
	if len(state) > 0 && !json.Valid(state) {
		state = nil
	}
	m.File.Edit = &media_library.ImageEdit{Source: dst, SourceName: original.Filename, State: state}
	return b.db.Save(m).Error
}

// largestSize is the key and the ratio of the largest size of cfg — the
// fixed crop of an upload —; "", 0 when it has none.
func largestSize(cfg *media_library.MediaBoxConfig) (key string, ratio float64) {
	var area int
	for k, s := range cfg.Sizes {
		if s == nil || s.Width <= 0 || s.Height <= 0 {
			continue
		}
		if a := s.Width * s.Height; a > area || (a == area && k < key) {
			area, key, ratio = a, k, float64(s.Width)/float64(s.Height)
		}
	}
	return
}

// uploadThroughEditor is the script that uploads files (a JS expression: an
// array of File) by send — a script reading list, [{file, original, edit,
// cropOption}] —: through the editor first when cfg.EditOnUpload, a photo of
// the camera, or always (a JS expression: the switch "Edit before
// uploading"). keep (cfg.UploadEdit) posts, with each result, its original,
// the edit and the fixed crop (NewOriginal<i>, NewEdit<i>, NewCrop<i>).
func uploadThroughEditor(cfg *media_library.MediaBoxConfig, files, send string, camera bool, always string) string {
	_, ratio := largestSize(cfg)
	opts := map[string]any{"aspect": ratio, "camera": camera}
	if cfg.EditTab != "" {
		opts["tab"] = cfg.EditTab
	}
	if cfg.EditMode != "" {
		opts["mode"] = cfg.EditMode
	}
	if cfg.EditMaxPixels > 0 {
		opts["maxPixels"] = cfg.EditMaxPixels
	}
	if always == "" {
		always = "false"
	}
	return fmt.Sprintf(`(function(files) {
const send = function(list) { %s };
const edit = window.rvqImageEditor && files.some(window.rvqImageEditor.isImage) && (%t || %t || (%s));
if (!edit) return send(files.map(function(f) { return {file: f} }));
return window.rvqImageEditor.editUploads(files, %s).then(function(list) {
if (list === 'retake') { if (typeof locals !== 'undefined') locals.fileChooserCamera = true; return }
if (list) return send(list);
});
})(%s)`, send, cfg.EditOnUpload, camera, always, h.JSONString(opts), files)
}

// uploadFormScript puts list (the editor's) in form to be posted: the files
// (NewFiles) and, when cfg keeps the edits, each one's original, edit and
// fixed crop.
func uploadFormScript(cfg *media_library.MediaBoxConfig, form string) string {
	s := fmt.Sprintf(`Object.keys(%[1]s).forEach(function(k) { if (/^New(Original|Edit|Crop)\d+$/.test(k)) delete %[1]s[k] });
%[1]s.NewFiles = list.map(function(u) { return u.file });
list.forEach(function(u, i) { if (u.cropOption) %[1]s['NewCrop' + i] = JSON.stringify(u.cropOption) });`, form)
	if cfg.UploadEdit {
		s += fmt.Sprintf(`
list.forEach(function(u, i) { if (u.original && u.edit) { %[1]s['NewOriginal' + i] = [u.original]; %[1]s['NewEdit' + i] = JSON.stringify(u.edit) } });`, form)
	}
	return s
}

// afterUpload keeps, of the i-th file uploaded (m), what the editor made with
// it: the fixed crop of the largest size (NewCrop<i>), and its original and
// edit (NewOriginal<i>, NewEdit<i>) when cfg keeps them.
func (b *Builder) afterUpload(ctx *web.EventContext, i int, m *media_library.MediaLibrary, cfg *media_library.MediaBoxConfig) error {
	form := ctx.R.MultipartForm
	if form == nil || m.SelectedType != media_library.ALLOW_TYPE_IMAGE {
		return nil
	}
	if cfg != nil && cfg.UploadEdit {
		if fs := form.File[fmt.Sprintf("NewOriginal%d", i)]; len(fs) > 0 {
			var state json.RawMessage
			if v := form.Value[fmt.Sprintf("NewEdit%d", i)]; len(v) > 0 {
				state = json.RawMessage(v[0])
			}
			if err := b.keepUploadEdit(m, fs[0], state); err != nil {
				return err
			}
		}
	}
	key, _ := largestSize(cfg)
	if v := form.Value[fmt.Sprintf("NewCrop%d", i)]; key != "" && len(v) > 0 && v[0] != "" {
		var c cropper.Value
		if err := json.Unmarshal([]byte(v[0]), &c); err != nil {
			return err
		}
		moption := m.GetMediaOption()
		if moption.CropOptions == nil {
			moption.CropOptions = map[string]*base.CropOption{}
		}
		moption.CropOptions[key] = &base.CropOption{X: int(c.X), Y: int(c.Y), Width: int(c.Width), Height: int(c.Height)}
		moption.Crop = true
		if err := m.ScanMediaOptions(moption); err != nil {
			return err
		}
		return base.SaveUploadAndCropImage(&b.Config, b.db, m)
	}
	return nil
}
