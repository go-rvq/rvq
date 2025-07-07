package web

import (
	"fmt"
	"strings"
)

type ChangeEvent string

const (
	Changed ChangeEvent = "changed"
	Created ChangeEvent = "created"
	Updated ChangeEvent = "updated"
	Deleted ChangeEvent = "deleted"

	ChangedEventKey = "$CHANGED_EVENT$"
	IdKey           = "$CHANGED_ID$"
)

func ApplyChangeEvent(s string, event ChangeEvent, id string) string {
	return strings.NewReplacer(
		IdKey, id,
		ChangedEventKey, fmt.Sprintf(string(event)),
	).Replace(s)
}

func ParseChangedEventValue(v string) (event ChangeEvent, value string) {
	parts := strings.SplitN(v, ":", 1)
	if len(parts) == 2 {
		switch ChangeEvent(parts[0]) {
		case Changed, Updated, Deleted:
			return ChangeEvent(parts[0]), parts[1]
		}
	}
	return
}

func (b *VueEventTagBuilder) ChangeEvent(event ChangeEvent) (r *VueEventTagBuilder) {
	return b.Query("presets_change_event", fmt.Sprintf("%s:%s", event, IdKey))
}

func (b *VueEventTagBuilder) CreatedEvent() (r *VueEventTagBuilder) {
	return b.ChangeEvent(Created)
}

func (b *VueEventTagBuilder) UpdatedEvent() (r *VueEventTagBuilder) {
	return b.ChangeEvent(Updated)
}

func (b *VueEventTagBuilder) DeletedEvent() (r *VueEventTagBuilder) {
	return b.ChangeEvent(Deleted)
}

func (b *VueEventTagBuilder) AnyEvent() (r *VueEventTagBuilder) {
	return b.Query("presets_change_event", fmt.Sprintf("%s:%s", ChangedEventKey, IdKey))
}
