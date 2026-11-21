package activity

import (
	"context"
	"sync"
)

// actionLabels are the labels of the actions of other packages (a custom
// record: AddCustomizedRecord), by their name: RegisterActionLabel.
var actionLabels struct {
	sync.RWMutex
	m map[string]func(ctx context.Context) string
}

// RegisterActionLabel names the action of a custom record (its name, as
// AddCustomizedRecord got it) in the language of ctx, where a log shows it —
// the history of a record, the detail of a log —: without it, its name.
func RegisterActionLabel(action string, label func(ctx context.Context) string) {
	actionLabels.Lock()
	defer actionLabels.Unlock()
	if actionLabels.m == nil {
		actionLabels.m = map[string]func(ctx context.Context) string{}
	}
	actionLabels.m[action] = label
}

// ActionLabel is the action as a log shows it in the language of ctx: one of
// this package's (Create, Edit, Delete, View), one registered
// (RegisterActionLabel), else its name.
func ActionLabel(ctx context.Context, action string) string {
	return actionLabel(ctx, getMessages(ctx), action)
}

func actionLabel(ctx context.Context, msgr *Messages, action string) string {
	switch action {
	case ActivityCreate:
		return msgr.ActionCreate
	case ActivityEdit:
		return msgr.ActionEdit
	case ActivityDelete:
		return msgr.ActionDelete
	case ActivityView:
		return msgr.ActionView
	}
	actionLabels.RLock()
	label := actionLabels.m[action]
	actionLabels.RUnlock()
	if label != nil {
		if s := label(ctx); s != "" {
			return s
		}
	}
	return action
}
