package i18n

import (
	"context"
	"time"
)

type TimeLayoutSpec struct {
	Default TimeLayout
	Short   TimeLayout
	Full    TimeLayout
}

func (s *TimeLayoutSpec) FromLevel(level TimeFormatLevel) TimeLayout {
	switch level {
	case TimeFormatShort:
		return s.Short
	case TimeFormatFull:
		return s.Full
	default:
		return s.Default
	}
}

func (s *TimeLayoutSpec) Level(level TimeFormatLevel) TimeLayoutFormatter {
	return TimeFormats.Get(s.FromLevel(level), s.Default)
}

type TimeFormatLevel uint8

type TimeFormatHandler func(ctx context.Context) func(t time.Time) string

type TimeFormatLeveled struct {
	m      *DefaultMessages
	s      *TimeLayoutSpec
	defaul TimeLayoutString
}

func (h *TimeFormatLeveled) Default() TimeFormatHandler {
	return h.Of(TimeFormatDefault)
}

func (h *TimeFormatLeveled) Short() TimeFormatHandler {
	return h.Of(TimeFormatShort)
}

func (h *TimeFormatLeveled) Full() TimeFormatHandler {
	return h.Of(TimeFormatFull)
}

func (h *TimeFormatLeveled) Of(level TimeFormatLevel) TimeFormatHandler {
	var f TimeLayoutFormatter
	if f = h.s.Level(level); f == nil {
		f = h.defaul
	}
	return func(ctx context.Context) func(t time.Time) string {
		return f.Formatter(h.m, ctx)
	}
}

const (
	TimeFormatDefault TimeFormatLevel = iota
	TimeFormatShort
	TimeFormatFull

	DefaultTimeFormatter     TimeLayoutString = "15:04:05"
	DefaultDateFormatter     TimeLayoutString = "2006-01-02"
	DefaultDateTimeFormatter TimeLayoutString = "2006-01-02 15:04:05"
)

type TimeLayout string

type TimeLayoutFormatter interface {
	Formatter(m *DefaultMessages, ctx context.Context) func(t time.Time) string
}

type TimeLayoutString string

func (l TimeLayoutString) Formatter(*DefaultMessages, context.Context) func(t time.Time) string {
	return func(t time.Time) string {
		return t.Format(string(l))
	}
}

type TimeLayoutFunc func(m *DefaultMessages, ctx context.Context) func(t time.Time) string

func (f TimeLayoutFunc) Formatter(m *DefaultMessages, ctx context.Context) func(t time.Time) string {
	return f(m, ctx)
}

type TimeLayoutFormat struct {
	Name      string
	Formatter TimeLayoutFormatter
}

type TimeFormatsMap map[TimeLayout]*TimeLayoutFormat

var TimeFormats = TimeFormatsMap{}

func (m TimeFormatsMap) Add(name TimeLayout, formatter TimeLayoutFormatter) {
	m[name] = &TimeLayoutFormat{string(name), formatter}
}

func (m TimeFormatsMap) AddLayout(name TimeLayout, layout string) {
	m.Add(name, TimeLayoutString(layout))
}

func (m TimeFormatsMap) AddFunc(name TimeLayout, formatter TimeLayoutFunc) {
	m.Add(name, formatter)
}

func (m TimeFormatsMap) AddFuncs(f func(add func(name TimeLayout, formatter TimeLayoutFunc))) {
	f(m.AddFunc)
}

func (m TimeFormatsMap) AddLayouts(f func(add func(name TimeLayout, layout string))) {
	f(m.AddLayout)
}

func (m TimeFormatsMap) Get(layout TimeLayout, fallback ...TimeLayout) (f TimeLayoutFormatter) {
	if f := m[layout]; f != nil {
		return f.Formatter
	}
	for _, fb := range fallback {
		if f := m[fb]; f != nil {
			return f.Formatter
		}
	}
	return
}

func RegisterTimeLayoutFormatterAlias(from TimeLayout, to ...TimeLayout) {
	f := TimeFormats[from]
	for _, s := range to {
		TimeFormats[s] = f
	}
}

func (m *DefaultMessages) DateFormatter() *TimeFormatLeveled {
	return &TimeFormatLeveled{m: m, s: &m.DateLayout, defaul: DefaultDateFormatter}
}

func (m *DefaultMessages) TimeFormatter() *TimeFormatLeveled {
	return &TimeFormatLeveled{m: m, s: &m.TimeLayout, defaul: DefaultTimeFormatter}
}

func (m *DefaultMessages) DateTimeFormatter() *TimeFormatLeveled {
	return &TimeFormatLeveled{m: m, s: &m.DateTimeLayout, defaul: DefaultDateTimeFormatter}
}
