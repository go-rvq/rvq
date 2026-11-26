package media

import (
	"context"
	"fmt"
	"path/filepath"
	"reflect"
	"strings"

	h "github.com/go-rvq/htmlgo"
	"github.com/go-rvq/rvq/admin/media/base"
	"github.com/go-rvq/rvq/admin/media/media_library"
	"github.com/go-rvq/rvq/admin/presets"
	"github.com/go-rvq/rvq/web"
	"github.com/go-rvq/rvq/x/ui/fileicons"
	v "github.com/go-rvq/rvq/x/ui/vuetify"
	"gorm.io/gorm"
)

// The attachments of a record: its images, or its files, each with a
// description, in their order — a list of records of a model of their own
// (one a media: a media box, a description, a position), nested in the
// record's form and detail, all of them alike wherever they are:
//
//   - the images: cards (ImagesCards) — the lightbox through them, the
//     description of each edited in its media box, kept in its record —,
//     and a place to drop them;
//   - the files: a table — the icon of the extension, the name (cut, whole by
//     a click), the description, the size — and a place to drop them.
//
//	imgs := media.ImageAttachments(b, db, &Image{})
//	imgs.Install(&ed.FieldsBuilder, "Images")
//	files := media.FileAttachments(b, db, &File{})
//	files.Install(&ed.FieldsBuilder, "Files")
//
// The record of an attachment has the fields of the media box, the
// description and the position (Fields: Media, Description, Position by
// default) and an ID. Saved with its owner, in its order: Save.

// AttachmentsBuilder makes a list of attachments (ImageAttachments,
// FileAttachments).
type AttachmentsBuilder struct {
	b           *presets.Builder
	db          *gorm.DB
	model       any
	images      bool
	media       string
	description string
	position    string
	label       func(ctx context.Context) string
	placeholder func(ctx context.Context) string
	mb          *presets.ModelBuilder
	ed          *presets.EditingBuilder
}

// ImageAttachments are the images of a record: model, the record of an image.
func ImageAttachments(b *presets.Builder, db *gorm.DB, model any) *AttachmentsBuilder {
	return &AttachmentsBuilder{b: b, db: db, model: model, images: true, media: "Media", description: "Description", position: "Position"}
}

// FileAttachments are the files of a record: model, the record of a file.
func FileAttachments(b *presets.Builder, db *gorm.DB, model any) *AttachmentsBuilder {
	return &AttachmentsBuilder{b: b, db: db, model: model, media: "Media", description: "Description", position: "Position"}
}

// Fields are the names of the fields of the record of an attachment: its media
// box, its description and its position.
func (a *AttachmentsBuilder) Fields(media, description, position string) *AttachmentsBuilder {
	a.media, a.description, a.position = media, description, position
	return a
}

// Label is the label of the list (the media library's Images or Files when
// none).
func (a *AttachmentsBuilder) Label(f func(ctx context.Context) string) *AttachmentsBuilder {
	a.label = f
	return a
}

// Placeholder is the words of the place to drop the files on (the media
// library's when none).
func (a *AttachmentsBuilder) Placeholder(f func(ctx context.Context) string) *AttachmentsBuilder {
	a.placeholder = f
	return a
}

// config is the media box of an attachment: an image's, its description in
// its record (DescriptionField); a file's, of any type.
func (a *AttachmentsBuilder) config() *media_library.MediaBoxConfig {
	if a.images {
		return &media_library.MediaBoxConfig{AllowType: "image", DescriptionField: a.description}
	}
	return &media_library.MediaBoxConfig{DisableDescription: true}
}

// Model is the model of the record of an attachment (made once).
func (a *AttachmentsBuilder) Model() *presets.ModelBuilder {
	if a.mb != nil {
		return a.mb
	}
	a.mb = presets.NewModelBuilder(a.b, a.model)
	if a.images {
		a.ed = a.mb.Editing(a.media)
		a.ed.HiddenField("ID")
		a.ed.Field(a.media).WithContextValue(MediaBoxConfig, a.config()).
			SetterFunc(MediaBoxSetterFunc(a.db)).
			SetI18nLabel(func(ctx context.Context) string { return GetMessages(ctx).AttachmentImage }).
			ComponentFunc(ImageCardComponentFunc(a.db))
		return a.mb
	}
	a.ed = a.mb.Editing("Icon", a.media, a.description, "Size")
	a.ed.HiddenField("ID")
	nothing := func(interface{}, *presets.FieldContext, *web.EventContext) error { return nil }
	boxOf := func(field *presets.FieldContext) media_library.MediaBox {
		rv := reflect.Indirect(reflect.ValueOf(field.Obj))
		if rv.Kind() == reflect.Struct {
			if f := rv.FieldByName(a.media); f.IsValid() {
				if mbx, ok := f.Interface().(media_library.MediaBox); ok {
					return mbx
				}
			}
		}
		return media_library.MediaBox{}
	}
	a.ed.Field("Icon").SetI18nLabel(func(context.Context) string { return "" }).
		ComponentFunc(func(field *presets.FieldContext, _ *web.EventContext) h.HTMLComponent {
			return v.VIcon(FileIcon(boxOf(field).FileName)).Color("primary")
		}).SetterFunc(nothing)
	// the name: the file chosen — cut, whole by a click, its value kept in the
	// form —; none yet, the chooser
	box := MediaBoxComponentFunc(a.db, false)
	a.ed.Field(a.media).WithContextValue(MediaBoxConfig, a.config()).
		SetterFunc(MediaBoxSetterFunc(a.db)).
		SetI18nLabel(func(ctx context.Context) string { return GetMessages(ctx).AttachmentName }).
		ComponentFunc(func(field *presets.FieldContext, ctx *web.EventContext) h.HTMLComponent {
			mbx := boxOf(field)
			if mbx.FileName == "" {
				return box(field, ctx)
			}
			value, _ := mbx.Value()
			return web.Scope(
				h.Input("").Type("hidden").Attr(web.VField(field.FormKey+".Values", fmt.Sprint(value))...),
				h.Span(mbx.FileName).
					Attr(":class", "locals.whole ? 'text-break' : 'd-inline-block text-truncate'").
					Style("max-width: 220px; cursor: pointer; vertical-align: middle").
					Attr("title", mbx.FileName).
					Attr("@click", "locals.whole = !locals.whole").
					Attr("data-attachment-name", true),
			).Slot("{ locals }").LocalsInit("{whole: false}")
		})
	a.ed.Field(a.description).
		SetI18nLabel(func(ctx context.Context) string { return GetMessages(ctx).AttachmentDescription }).
		ComponentFunc(func(field *presets.FieldContext, _ *web.EventContext) h.HTMLComponent {
			if field.ReadOnly {
				return h.Span(field.StringValue()).Class("text-pre-wrap")
			}
			return v.VTextarea().Variant(v.FieldVariantUnderlined).Rows(1).AutoGrow(true).
				Attr(web.VField(field.FormKey, field.StringValue())...).Label(field.InputLabel())
		})
	a.ed.Field("Size").SetI18nLabel(func(ctx context.Context) string { return GetMessages(ctx).AttachmentSize }).
		ComponentFunc(func(field *presets.FieldContext, _ *web.EventContext) h.HTMLComponent {
			t := "-"
			if n := boxOf(field).FileSizes["original"]; n > 0 {
				t = base.ByteCountSI(n)
			}
			return h.Span(t).Class("text-no-wrap text-medium-emphasis")
		}).SetterFunc(nothing)
	return a.mb
}

// Nested is the nested slice of the attachments: the cards of the images, or
// the table of the files; each with the place to drop them.
func (a *AttachmentsBuilder) Nested() *presets.NestedSliceBuilder {
	a.Model()
	words := func(ctx *presets.ListEditorContainerContext) string {
		if a.placeholder != nil {
			return a.placeholder(ctx.FieldContext.EventContext.Context())
		}
		m := GetMessages(ctx.FieldContext.EventContext.Context())
		if a.images {
			return m.DropImages
		}
		return m.DropFiles
	}
	if a.images {
		return presets.NestedSlice(a.mb, &a.ed.FieldsBuilder).
			SetComponentBuilder(presets.NewListEditorCardsBuilder().
				FloatingDelete(true).
				LayoutFunc(func(*presets.FieldsBuilder) presets.FieldsLayout { return presets.FieldsLayout{a.media} }).
				AfterFunc(func(ctx *presets.ListEditorContainerContext) h.HTMLComponent {
					return a.drop(ctx, words(ctx))
				}))
	}
	return presets.NestedSlice(a.mb, &a.ed.FieldsBuilder).
		SetComponentBuilder(presets.NewListEditorTableBuilder().
			LayoutFunc(func(*presets.FieldsBuilder) presets.FieldsLayout {
				return presets.FieldsLayout{"Icon", a.media, a.description, "Size"}
			}).
			AfterFunc(func(ctx *presets.ListEditorContainerContext) h.HTMLComponent {
				return a.drop(ctx, words(ctx))
			}))
}

// drop is, in the form, the place to drop the files of the list — each
// uploaded, an attachment of its own —; nothing when read only.
func (a *AttachmentsBuilder) drop(ctx *presets.ListEditorContainerContext, words string) h.HTMLComponent {
	if ctx.FieldContext.ReadOnly {
		return nil
	}
	return FileDrop(a.config(), words, fmt.Sprintf("$listEditorAddItems(boxes.map(m => ({%s: m})))", a.media))
}

// Install puts the attachments on the field name of fb: its label, the nested
// slice, the list marked (data-attachments: images or files).
func (a *AttachmentsBuilder) Install(fb *presets.FieldsBuilder, name string) *presets.FieldBuilder {
	kind := "files"
	if a.images {
		kind = "images"
	}
	return fb.Field(name).
		SetI18nLabel(func(ctx context.Context) string {
			if a.label != nil {
				return a.label(ctx)
			}
			if a.images {
				return GetMessages(ctx).AttachmentImages
			}
			return GetMessages(ctx).AttachmentFiles
		}).
		Nested(a.Nested()).
		WrapComponentFunc(func(old presets.FieldComponentFunc) presets.FieldComponentFunc {
			return func(field *presets.FieldContext, ctx *web.EventContext) h.HTMLComponent {
				return h.Div(old(field, ctx)).Attr("data-attachments", kind)
			}
		})
}

// Kept are the attachments of the list key the form keeps (not removed), in
// their order, of current (the owner's as loaded, a slice of pointers of the
// model); ok false when the list was not in the form.
func (a *AttachmentsBuilder) Kept(ctx *web.EventContext, key string, current any) (kept []any, ok bool) {
	part := presets.PartitionListEditorItems(presets.ListEditorFormValues(ctx), key, reflect.ValueOf(current))
	if part == nil {
		return nil, false
	}
	for _, rv := range part.Kept {
		kept = append(kept, rv.Interface())
	}
	return kept, true
}

// Save makes the attachments of the owner ownerID (its field ownerField in
// the record of an attachment) kept (SaveOf).
func (a *AttachmentsBuilder) Save(tx *gorm.DB, ownerField string, ownerID any, kept []any, set func(obj any)) error {
	return a.SaveOf(tx, map[string]any{ownerField: ownerID}, kept, set)
}

// SaveOf makes the attachments of the owner — its fields in the record of an
// attachment and their values (one: its id; two of a polymorphic owner: its
// id and its kind) — kept (Kept): those of before removed, each of kept
// created in its order — its position 1, 2, 3…, the fields of the owner —,
// and set (set, nil none: the organization of a record, say) before.
func (a *AttachmentsBuilder) SaveOf(tx *gorm.DB, owner map[string]any, kept []any, set func(obj any)) error {
	q := tx.Session(&gorm.Session{NewDB: true})
	for f, v := range owner {
		q = q.Where(tx.NamingStrategy.ColumnName("", f)+" = ?", v)
	}
	if err := q.Delete(reflect.New(reflect.TypeOf(a.model).Elem()).Interface()).Error; err != nil {
		return err
	}
	for i, obj := range kept {
		rv := reflect.Indirect(reflect.ValueOf(obj))
		for f, v := range owner {
			if fv := rv.FieldByName(f); fv.IsValid() && fv.CanSet() {
				fv.Set(reflect.ValueOf(v).Convert(fv.Type()))
			}
		}
		if f := rv.FieldByName(a.position); f.IsValid() && f.CanSet() {
			f.SetInt(int64(i + 1))
		}
		if set != nil {
			set(obj)
		}
		if err := tx.Session(&gorm.Session{NewDB: true}).Create(obj).Error; err != nil {
			return err
		}
	}
	return nil
}

// The kinds of a change of the attachments of an owner (AttachmentChange).
const (
	AttachmentAdded     = "added"
	AttachmentRemoved   = "removed"
	AttachmentDescribed = "described"
)

// AttachmentChange is a change of the attachments of an owner
// (SaveOfChanges): one added (Now its description), removed (Old its
// description) or described again (Old, Now) — Name its file's name.
type AttachmentChange struct {
	Kind, Name, Old, Now string
}

// SaveOfChanges is SaveOf, and the changes it made to the attachments of the
// owner: of the ones stored before — an attachment by its ID —, those that
// left, those that came, those of another description; in their order.
func (a *AttachmentsBuilder) SaveOfChanges(tx *gorm.DB, owner map[string]any, kept []any, set func(obj any)) ([]AttachmentChange, error) {
	q := tx.Session(&gorm.Session{NewDB: true})
	for f, v := range owner {
		q = q.Where(tx.NamingStrategy.ColumnName("", f)+" = ?", v)
	}
	olds := reflect.New(reflect.SliceOf(reflect.TypeOf(a.model)))
	if err := q.Order(tx.NamingStrategy.ColumnName("", a.position)).Find(olds.Interface()).Error; err != nil {
		return nil, err
	}
	field := func(obj any, name string) reflect.Value {
		return reflect.Indirect(reflect.ValueOf(obj)).FieldByName(name)
	}
	idOf := func(obj any) string {
		if f := field(obj, "ID"); f.IsValid() && !f.IsZero() {
			return fmt.Sprint(f.Interface())
		}
		return ""
	}
	text := func(obj any, name string) string {
		if f := field(obj, name); f.IsValid() {
			if mbx, ok := f.Interface().(media_library.MediaBox); ok {
				return mbx.FileName
			}
			return fmt.Sprint(f.Interface())
		}
		return ""
	}
	before := map[string]any{}
	var order []string
	for i := 0; i < olds.Elem().Len(); i++ {
		o := olds.Elem().Index(i).Interface()
		before[idOf(o)] = o
		order = append(order, idOf(o))
	}
	var changes []AttachmentChange
	seen := map[string]bool{}
	for _, k := range kept {
		id := idOf(k)
		o, ok := before[id]
		if id == "" || !ok {
			changes = append(changes, AttachmentChange{Kind: AttachmentAdded, Name: text(k, a.media), Now: text(k, a.description)})
			continue
		}
		seen[id] = true
		if od, nd := text(o, a.description), text(k, a.description); od != nd {
			changes = append(changes, AttachmentChange{Kind: AttachmentDescribed, Name: text(k, a.media), Old: od, Now: nd})
		}
	}
	for _, id := range order {
		if !seen[id] {
			o := before[id]
			changes = append(changes, AttachmentChange{Kind: AttachmentRemoved, Name: text(o, a.media), Old: text(o, a.description)})
		}
	}
	return changes, a.SaveOf(tx, owner, kept, set)
}

// FileIcon is the icon (MDI) of a file by its name's extension.
func FileIcon(name string) string {
	ext := strings.ToLower(strings.TrimPrefix(filepath.Ext(name), "."))
	if ext == "" {
		return "mdi-file-outline"
	}
	return fileicons.MDI(ext)
}
