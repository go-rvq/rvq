package presets

type (
	ModelSetupFactory   func(mb *ModelBuilder) ModelSetuper
	ModelSetupFactories []ModelSetupFactory
)

func (f *ModelSetupFactories) Append(factory ...ModelSetupFactory) {
	*f = append(*f, factory...)
}

func (f ModelSetupFactories) Of(mb *ModelBuilder) (s ModelSetupers) {
	s = make(ModelSetupers, len(f))
	for i, factory := range f {
		s[i] = factory(mb)
	}
	return
}

var DefaultModelSetupFactories ModelSetupFactories

type (
	ModelSetuper interface {
		Init()
		InitFields(mode FieldMode, fieldBuilders *FieldBuilders)
		InitField(mode FieldMode, f *FieldBuilder)
		ConfigureField(mode FieldMode, f *FieldBuilder)
	}

	ModelSetupers []ModelSetuper
)

func (s ModelSetupers) Init() {
	for _, m := range s {
		m.Init()
	}
}

func (s ModelSetupers) InitFields(mode FieldMode, fieldBuilders *FieldBuilders) {
	for _, m := range s {
		m.InitFields(mode, fieldBuilders)
	}
}

func (s ModelSetupers) InitField(mode FieldMode, f *FieldBuilder) {
	for _, m := range s {
		m.InitField(mode, f)
	}
}

func (s ModelSetupers) ConfigureField(mode FieldMode, f *FieldBuilder) {
	for _, m := range s {
		m.ConfigureField(mode, f)
	}
}

func init() {
	DefaultModelSetupFactories.Append(func(*ModelBuilder) ModelSetuper {
		return defaultModelSetup{}
	})
}

type defaultModelSetup struct {
}

func (defaultModelSetup) Init() {}

func (defaultModelSetup) InitFields(mode FieldMode, fieldBuilders *FieldBuilders) {
}

func (defaultModelSetup) InitField(mode FieldMode, f *FieldBuilder) {
	switch f.name {
	case "CreatedAt", "UpdatedAt", "DeletedAt":
		f.SetMode(DETAIL)
		f.SetAudited(true)
	}
}

func (defaultModelSetup) ConfigureField(mode FieldMode, f *FieldBuilder) {
	switch f.name {
	case "CreatedAt", "UpdatedAt", "DeletedAt":
		f.DisableZeroRender()
	}
}
