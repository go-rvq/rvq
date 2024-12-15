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
		InitFields(fieldBuilders *FieldBuilders)
		FieldSetuper
	}

	ModelSetupers []ModelSetuper
)

func (s ModelSetupers) Init() {
	for _, m := range s {
		m.Init()
	}
}

func (s ModelSetupers) InitFields(fieldBuilders *FieldBuilders) {
	for _, m := range s {
		m.InitFields(fieldBuilders)
	}
}

func (s ModelSetupers) InitField(f *FieldBuilder) {
	for _, m := range s {
		m.InitField(f)
	}
}

func (s ModelSetupers) ConfigureField(f *FieldBuilder) {
	for _, m := range s {
		m.ConfigureField(f)
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

func (defaultModelSetup) InitFields(fb *FieldBuilders) {
}

func (defaultModelSetup) InitField(f *FieldBuilder) {
	switch f.name {
	case "CreatedAt", "UpdatedAt", "DeletedAt":
		f.SetMode(DETAIL)
		f.SetAudited(true)
	}
}

func (defaultModelSetup) ConfigureField(f *FieldBuilder) {
	switch f.name {
	case "CreatedAt", "UpdatedAt", "DeletedAt":
		f.DisableZeroRender()
	}
}
