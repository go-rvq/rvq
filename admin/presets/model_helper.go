package presets

import "github.com/go-rvq/rvq/web"

const (
	ToTitleRecordEncoderName = "presets__toTitle"
)

type ToTitleFactory[T any] func(ctx *web.EventContext) func(r T) string

func ToTitleRecordEncoderFactory[T any](l *ListingBuilder, ttf ...ToTitleFactory[T]) RecordEncoderFactory[any] {
	var tf ToTitleFactory[any]
	for _, tfT := range ttf {
		if tfT == nil {
			continue
		}

		tf = func(ctx *web.EventContext) func(r any) string {
			f := tfT(ctx)
			return func(r any) string {
				return f(r.(T))
			}
		}
	}
	return ToAnyTitleRecordEncoderFactory(l, tf)
}

func ToAnyTitleRecordEncoderFactory(l *ListingBuilder, ttf ...ToTitleFactory[any]) RecordEncoderFactory[any] {
	var tf ToTitleFactory[any]
	for _, tf = range ttf {
	}
	if tf == nil {
		tf = func(ctx *web.EventContext) func(r any) string {
			return func(r any) string {
				return l.mb.RecordTitle(r, ctx)
			}
		}
	}
	return func(ctx *web.EventContext) func(r any) any {
		toTitle := tf(ctx)
		return func(r any) any {
			return map[string]any{
				"ID":    l.mb.MustRecordID(r).String(),
				"Title": toTitle(r),
			}
		}
	}
}
