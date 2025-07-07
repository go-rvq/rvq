package presets

import (
	"database/sql"
	"fmt"
	"mime/multipart"
	"path/filepath"
	"reflect"
	"time"

	"github.com/qor5/admin/v3/reflect_utils"
	"github.com/qor5/web/v3/datafield"
	"github.com/shopspring/decimal"
)

type FieldMode uint8

func (f FieldMode) Is(m ...FieldMode) bool {
	for _, other := range m {
		if f == other {
			return true
		}
	}
	return false
}

func (f FieldMode) IsNew() bool {
	return f.Is(NEW)
}

func (f FieldMode) IsEdit() bool {
	return f.Is(EDIT)
}

func (f FieldMode) IsDetail() bool {
	return f.Is(DETAIL)
}

func (f FieldMode) IsList() bool {
	return f.Is(LIST)
}

func (f FieldMode) IsWrite() bool {
	return f.Is(NEW, EDIT)
}

func (m FieldMode) Has(f FieldMode) bool {
	return m&f != 0
}
func (m FieldMode) HasAny(f ...FieldMode) bool {
	for _, f := range f {
		if m.Has(f) {
			return true
		}
	}
	return false
}

func (f FieldMode) Split() (modes []FieldMode) {
	if f.Has(NEW) {
		modes = append(modes, NEW)
	}
	if f.Has(EDIT) {
		modes = append(modes, EDIT)
	}
	if f.Has(DETAIL) {
		modes = append(modes, DETAIL)
	}
	if f.Has(LIST) {
		modes = append(modes, LIST)
	}
	return
}

const (
	NONE FieldMode = 1 << iota
	LIST
	DETAIL
	NEW
	EDIT

	WRITE = NEW | EDIT
	ALL   = LIST | DETAIL | WRITE
)

type FieldModeStack []FieldMode

func (s FieldModeStack) Global() FieldMode {
	return s[0]
}

func (s FieldModeStack) Dot() FieldMode {
	return s[len(s)-1]
}

func (s FieldModeStack) DotStack() FieldModeStack {
	return FieldModeStack{s.Dot()}
}

func (s FieldModeStack) Push(m FieldMode) FieldModeStack {
	return append(s, m)
}

func (f FieldModeStack) Is(mode ...FieldMode) bool {
	if len(f) == 0 {
		return false
	}
	return f.Dot().Is(mode...)
}

func (f FieldModeStack) IsNew() bool {
	return f.Is(NEW)
}

func (f FieldModeStack) IsEdit() bool {
	return f.Is(EDIT)
}

func (f FieldModeStack) IsDetail() bool {
	return f.Is(DETAIL)
}

func (f FieldModeStack) IsList() bool {
	return f.Is(LIST)
}

func (f FieldModeStack) IsWrite() bool {
	return f.Is(NEW, EDIT)
}

type FieldDefaultBuilder struct {
	valType    reflect.Type
	compFunc   FieldComponentFunc
	setterFunc FieldSetterFunc

	datafield.DataField[*FieldDefaultBuilder]
}

func NewFieldDefault(t reflect.Type) *FieldDefaultBuilder {
	return datafield.New(&FieldDefaultBuilder{valType: t})
}

func (b *FieldDefaultBuilder) ComponentFunc(v FieldComponentFunc) (r *FieldDefaultBuilder) {
	b.compFunc = v
	return b
}

func (b *FieldDefaultBuilder) SetterFunc(v FieldSetterFunc) (r *FieldDefaultBuilder) {
	b.setterFunc = v
	return b
}

var numberVals = []interface{}{
	int(0), int8(0), int16(0), int32(0), int64(0),
	uint(0), uint8(8), uint16(0), uint32(0), uint64(0),
	float32(0.0), float64(0.0), decimal.Decimal{},
	sql.NullInt32{}, sql.NullInt64{}, sql.NullFloat64{},
	sql.Null[uint]{}, sql.Null[uint8]{}, sql.Null[uint16]{}, sql.Null[uint32]{}, sql.Null[uint64]{},
}

var stringVals = []interface{}{
	"",
	[]rune(""),
	[]byte(""),
}

var timeVals = []interface{}{
	time.Now(),
	ptrTime(time.Now()),
}

type FieldDefaults struct {
	mode             FieldMode
	fieldTypes       []*FieldDefaultBuilder
	excludesPatterns []string
	disablesPatterns []string
}

func NewFieldDefaults(t FieldMode) (r *FieldDefaults) {
	r = &FieldDefaults{
		mode: t,
	}

	r.builtInFieldTypes()
	return
}

func (b *FieldDefaults) FieldType(v interface{}) (r *FieldDefaultBuilder) {
	return b.fieldTypeByTypeOrCreate(reflect.TypeOf(v))
}

func (b *FieldDefaults) Exclude(patterns ...string) (r *FieldDefaults) {
	b.excludesPatterns = patterns
	return b
}

func (b *FieldDefaults) Disable(patterns ...string) (r *FieldDefaults) {
	b.disablesPatterns = patterns
	return b
}

func (b *FieldDefaults) InspectFields(val interface{}, setup ...FieldSetuper) (r *FieldsBuilder) {
	_, fields := reflect_utils.UniqueFieldsOfReflectType(reflect.TypeOf(val))
	fieldBuilders := b.SetupFields(b.NewFieldBuilders(fields), setup...)

	return &FieldsBuilder{
		model:    val,
		fields:   fieldBuilders,
		defaults: b,
	}
}

func (b *FieldDefaults) NewFieldBuilders(fields reflect_utils.IndexableStructFields) (fbs FieldBuilders) {
	fbs = make(FieldBuilders, len(fields))
	var validCount int

	for _, f := range fields {
		if hasMatched(b.excludesPatterns, f.Name) {
			continue
		}
		fbs[validCount] = &FieldBuilder{
			rt: f.Type,
			NameLabel: NameLabel{
				name: f.Name,
			},
			mode:        ALL,
			structField: f,
			disabled:    hasMatched(b.disablesPatterns, f.Name),
		}
		validCount++
	}
	return fbs[:validCount]
}

func (b *FieldDefaults) SetupFields(fbs FieldBuilders, setupers ...FieldSetuper) (result FieldBuilders) {
	var (
		applyFt = func(f *FieldBuilder, ft *FieldDefaultBuilder) {
			if f.compFunc == nil {
				f.ComponentFunc(ft.compFunc)
			}

			if f.setterFunc == nil {
				f.SetterFunc(ft.setterFunc)
			}

			f.DataField.SetMapData(ft.Data().Clone())
		}
		withFt = func(f *FieldBuilder, rt reflect.Type) {
			if ft := b.fieldTypeByType(rt); ft != nil {
				applyFt(f, ft)
			}
		}
	)

	setuper := FieldSetupers(setupers)

	for _, f := range fbs {
		setuper.InitField(f)
		withFt(f, f.structField.Type)
		setuper.ConfigureField(f)
		result = append(result, f)
	}

	return result
}

func hasMatched(patterns []string, name string) bool {
	for _, p := range patterns {
		ok, err := filepath.Match(p, name)
		if err != nil {
			panic(err)
		}
		if ok {
			return true
		}
	}
	return false
}

func (b *FieldDefaults) String() string {
	var types []string
	for _, t := range b.fieldTypes {
		types = append(types, fmt.Sprint(t.valType))
	}
	return fmt.Sprintf("mode: %d, types %v", b.mode, types)
}

func (b *FieldDefaults) fieldTypeByType(tv reflect.Type) (r *FieldDefaultBuilder) {
	for _, ft := range b.fieldTypes {
		if ft.valType == tv {
			return ft
		}
	}
	return nil
}

func (b *FieldDefaults) fieldTypeByTypeOrCreate(tv reflect.Type) (r *FieldDefaultBuilder) {
	if r = b.fieldTypeByType(tv); r != nil {
		return
	}

	r = NewFieldDefault(tv)

	if b.mode == LIST {
		r.ComponentFunc(TDStringComponentFunc)
	} else {
		r.ComponentFunc(TextFieldComponentFunc)
	}
	b.fieldTypes = append(b.fieldTypes, r)
	return
}

func (b *FieldDefaults) builtInFieldTypes() {
	if b.mode == LIST {
		b.FieldType(true).
			ComponentFunc(TDReadonlyBoolComponentFunc)

		for _, v := range numberVals {
			b.FieldType(v).
				ComponentFunc(TDStringComponentFunc)
		}

		for _, v := range stringVals {
			b.FieldType(v).
				ComponentFunc(TDStringComponentFunc)
		}
		return
	}

	if b.mode == DETAIL {
		b.FieldType(true).
			ComponentFunc(CheckboxReadonlyComponentFunc)

		for _, v := range numberVals {
			b.FieldType(v).
				ComponentFunc(ReadonlyComponentFunc)
		}

		for _, v := range stringVals {
			b.FieldType(v).
				ComponentFunc(ReadonlyComponentFunc)
		}

		for _, v := range timeVals {
			b.FieldType(v).
				ComponentFunc(TimeReadonlyComponentFunc)
		}
		return
	}

	b.FieldType(true).
		ComponentFunc(CheckboxComponentFunc)

	for _, v := range numberVals {
		b.FieldType(v).
			ComponentFunc(NumberComponentFunc)
	}

	for _, v := range stringVals {
		b.FieldType(v).
			ComponentFunc(TextFieldComponentFunc)
	}

	for _, v := range timeVals {
		b.FieldType(v).
			ComponentFunc(TimeComponentFunc).
			SetterFunc(TimeComponentFuncSetter)
	}

	b.FieldType([]*multipart.FileHeader{}).
		ComponentFunc(FileFieldComponentFunc)
	b.FieldType(&multipart.FileHeader{}).
		ComponentFunc(FileFieldComponentFunc)

	b.FieldType(rune(0)).
		ComponentFunc(RuneFieldComponentFunc)

	b.Disable("ID")
	return
}
