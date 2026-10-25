package profile

import (
	h "github.com/go-rvq/htmlgo"
	"github.com/go-rvq/rvq/admin/media"
	"github.com/go-rvq/rvq/admin/media/base"
	"github.com/go-rvq/rvq/admin/media/media_library"
	"github.com/go-rvq/rvq/admin/packages/user"
	"github.com/go-rvq/rvq/admin/presets"
	"github.com/go-rvq/rvq/web"
	"gorm.io/gorm"
)

// AvatarUser is a user with a picture of their own: the profile edits it
// (Avatar, an image of the media library, cropped square), the top bar shows
// it in place of the initials. A user model without it: the initials.
type AvatarUser interface {
	GetAvatar() media_library.MediaBox
	SetAvatar(v media_library.MediaBox)
}

// AvatarSize is the size of the picture of the user (its crop, square).
const AvatarSize = "avatar"

// avatarURL is the picture of u ("" when none).
func avatarURL(u any) string {
	au, ok := u.(AvatarUser)
	if !ok {
		return ""
	}
	box := au.GetAvatar()
	if box.Url == "" {
		return ""
	}
	return box.URL(AvatarSize)
}

// configureAvatarField: the picture in the profile, for a user whose model
// has one (AvatarUser); none otherwise.
func configureAvatarField(db *gorm.DB, e *presets.EditingBuilder, d *presets.DetailingBuilder) {
	cfg := &media_library.MediaBoxConfig{
		AllowType: "image",
		Sizes:     map[string]*base.Size{AvatarSize: {Width: 256, Height: 256}},
	}
	only := func(readonly bool) presets.FieldComponentFunc {
		return func(field *presets.FieldContext, ctx *web.EventContext) h.HTMLComponent {
			if _, ok := user.GetCurrentUser(ctx.R).(AvatarUser); !ok {
				return nil
			}
			return media.MediaBoxComponentFunc(db, readonly)(field, ctx)
		}
	}
	e.Field("Avatar").WithContextValue(media.MediaBoxConfig, cfg).ComponentFunc(only(false))
	d.Field("Avatar").WithContextValue(media.MediaBoxConfig, cfg).ComponentFunc(only(true))
}
