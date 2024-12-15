package presets

import "github.com/qor5/web/v3"

func GetFlashMessages(ctx *web.EventContext) (v FlashMessages) {
	v, _ = ctx.ContextValue(ctxFlashMessages).(FlashMessages)
	return v
}

func AddFlashMessage(ctx *web.EventContext, message ...*FlashMessage) {
	if old := web.GetContextValuer(ctx.R.Context(), ctxFlashMessages); old != nil {
		old.Set(append(old.Get().(FlashMessages), message...))
	} else {
		ctx.WithContextValue(ctxFlashMessages, FlashMessages(message))
	}
}
