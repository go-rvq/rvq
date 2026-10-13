package media

import (
	"errors"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/go-rvq/rvq/admin/media/base"
	"github.com/go-rvq/rvq/admin/media/media_library"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

// StartRecrop makes the crops of the files of the field key again in its
// sizes (Recrop), in the background; nothing when they are being made
// already. Its state is the use's (RecropState…).
func (b *Builder) StartRecrop(key string) bool {
	if !b.beginRecrop(key) {
		return false
	}
	go func() {
		if err := b.recrop(key); err != nil {
			log.Printf("ERROR: media: the crops of %s: %v", key, err)
		}
	}()
	return true
}

// Recrop makes the crops of the files of the field key again in its sizes —
// the files its uses show —: each size of the field set in the file, its area
// kept when its ratio is the same, else the largest of the ratio centered in
// the image; the boxes of the records then told the new sizes of their files.
// It waits for them to be made; StartRecrop does not.
func (b *Builder) Recrop(key string) error {
	if !b.beginRecrop(key) {
		return fmt.Errorf("media: the crops of %s are being made", key)
	}
	return b.recrop(key)
}

// beginRecrop marks the field's crops as being made; false when they are.
func (b *Builder) beginRecrop(key string) bool {
	now := time.Now()
	res := b.db.Model(&MediaBoxUsage{}).
		Where("usage_key = ? AND (recrop_state IS NULL OR recrop_state <> ?)", key, RecropRunning).
		Updates(map[string]any{"recrop_state": RecropRunning, "recrop_done": 0, "recrop_total": 0,
			"recrop_error": "", "recrop_at": &now})
	return res.Error == nil && res.RowsAffected == 1
}

func (b *Builder) recrop(key string) (err error) {
	u := b.usage(key)
	if u == nil {
		return fmt.Errorf("media: no field %s", key)
	}
	sizes := u.EffectiveSizes()
	finish := func(state string, errs []string) {
		up := map[string]any{"recrop_state": state, "recrop_error": strings.Join(errs, "\n")}
		if state == RecropDone {
			up["applied_sizes"] = encodeSizes(sizes)
		}
		b.db.Model(&MediaBoxUsage{}).Where("usage_key = ?", key).Updates(up)
	}

	var files []uuid.UUID
	if err = b.db.Model(&MediaLibraryUse{}).Where("record_table = ? AND field = ?", u.RecordTable, u.Field).
		Distinct().Pluck("media_library_id", &files).Error; err != nil {
		finish(RecropFailed, []string{err.Error()})
		return err
	}
	b.db.Model(&MediaBoxUsage{}).Where("usage_key = ?", key).Update("recrop_total", len(files))

	var errs []string
	for i, id := range files {
		if ferr := b.recropFile(u, id, sizes); ferr != nil {
			errs = append(errs, fmt.Sprintf("%s: %v", id, ferr))
		}
		b.db.Model(&MediaBoxUsage{}).Where("usage_key = ?", key).Update("recrop_done", i+1)
	}
	if len(errs) > 0 {
		finish(RecropFailed, errs)
		return errors.New(strings.Join(errs, "; "))
	}
	finish(RecropDone, nil)
	return nil
}

// recropFile makes the crops of the file id in sizes, and tells the boxes of
// the records of the field (u) that show it.
func (b *Builder) recropFile(u *MediaBoxUsage, id uuid.UUID, sizes []UsageSize) error {
	var m media_library.MediaLibrary
	if err := b.db.First(&m, "id = ?", id).Error; err != nil {
		return err
	}
	if !base.IsImageFormat(m.File.FileName) && m.SelectedType != media_library.ALLOW_TYPE_IMAGE {
		return nil // no image: no crops
	}
	opt := m.GetMediaOption()
	if opt.Sizes == nil {
		opt.Sizes = map[string]*base.Size{}
	}
	if opt.CropOptions == nil {
		opt.CropOptions = map[string]*base.CropOption{}
	}
	for name, s := range sizesOf(sizes, opt.Sizes) {
		opt.Sizes[name] = s
		if c := opt.CropOptions[name]; !sameRatio(c, s.Width, s.Height) {
			if cc := centeredCrop(m.File.Width, m.File.Height, s.Width, s.Height); cc != nil {
				opt.CropOptions[name] = cc
			} else {
				delete(opt.CropOptions, name)
			}
		}
	}
	opt.Crop = true
	if err := m.ScanMediaOptions(opt); err != nil {
		return err
	}
	if err := base.SaveUploadAndCropImage(&b.Config, b.db, &m); err != nil {
		return err
	}
	return b.tellBoxes(u, &m)
}

// tellBoxes writes the new sizes of the file m in the boxes of the records of
// the field (u) that show it.
func (b *Builder) tellBoxes(u *MediaBoxUsage, m *media_library.MediaLibrary) error {
	var uses []MediaLibraryUse
	if err := b.db.Where("record_table = ? AND field = ? AND media_library_id = ?", u.RecordTable, u.Field, m.ID).
		Find(&uses).Error; err != nil {
		return err
	}
	model := b.usageModel(u.Key)
	if model == nil || len(uses) == 0 {
		return nil
	}
	stmt := &gorm.Statement{DB: b.db}
	if err := stmt.Parse(model); err != nil {
		return err
	}
	f := stmt.Schema.LookUpField(u.Field)
	if f == nil || len(stmt.Schema.PrimaryFields) == 0 {
		return nil
	}
	for _, use := range uses {
		parts := strings.Split(use.RecordID, ",")
		if len(parts) != len(stmt.Schema.PrimaryFields) {
			continue
		}
		q := b.db.Table(stmt.Schema.Table)
		for i, pf := range stmt.Schema.PrimaryFields {
			q = q.Where(pf.DBName+" = ?", parts[i])
		}
		var raw string
		if err := q.Select(f.DBName).Row().Scan(&raw); err != nil {
			continue // a record gone meanwhile
		}
		var box media_library.MediaBox
		if box.Scan(raw) != nil || string(box.ID) != m.ID.String() {
			continue
		}
		box.Url = m.File.Url
		box.FileSizes = m.File.FileSizes
		upd := b.db.Table(stmt.Schema.Table)
		for i, pf := range stmt.Schema.PrimaryFields {
			upd = upd.Where(pf.DBName+" = ?", parts[i])
		}
		if err := upd.Update(f.DBName, box).Error; err != nil {
			return err
		}
	}
	return nil
}
