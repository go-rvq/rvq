package web

import (
	"fmt"
)

type idEventFunc struct {
	id string
	ef EventHandler
}

type EventsHub struct {
	eventFuncs []*idEventFunc
	wraper     func(ef EventHandler) EventHandler
}

func (eh *EventsHub) Wraper(f func(ef EventHandler) EventHandler) *EventsHub {
	eh.wraper = f
	return eh
}

func (eh *EventsHub) Wrap(f func(ef EventHandler) EventHandler) *EventsHub {
	for _, ef := range eh.eventFuncs {
		ef.ef = f(ef.ef)
	}
	return eh
}

func (eh *EventsHub) String() string {
	var rs []string
	for _, ne := range eh.eventFuncs {
		rs = append(rs, ne.id)
	}
	return fmt.Sprintf("%#+v", rs)
}

func (eh *EventsHub) RegisterEventHandler(eventFuncId string, ef EventHandler) (key string) {
	key = eventFuncId
	if eh.eventHandleById(eventFuncId) != nil {
		return
	}

	if eh.wraper != nil {
		ef = eh.wraper(ef)
	}

	eh.eventFuncs = append(eh.eventFuncs, &idEventFunc{eventFuncId, ef})
	return
}

func (eh *EventsHub) RegisterEventFunc(eventFuncId string, ef EventFunc) (key string) {
	return eh.RegisterEventHandler(eventFuncId, ef)
}

func (eh *EventsHub) addMultipleEventFuncs(vs ...interface{}) (key string) {
	if len(vs)%2 != 0 {
		panic("id and func not paired")
	}
	for i := 0; i < len(vs); i = i + 2 {
		eh.RegisterEventHandler(vs[i].(string), EventFunc(vs[i+1].(func(ctx *EventContext) (r EventResponse, err error))))
	}
	return
}

func (eh *EventsHub) eventHandleById(id string) (r EventHandler) {
	for _, ne := range eh.eventFuncs {
		if ne.id == id {
			r = ne.ef
			return
		}
	}
	return
}

func (eh *EventsHub) Merge(hub *EventsHub) {
	eh.eventFuncs = append(eh.eventFuncs, hub.eventFuncs...)
}
