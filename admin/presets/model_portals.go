package presets

type ModelBuilderPortals struct {
	id string
}

func NewModelBuilderPortals(id string) *ModelBuilderPortals {
	return &ModelBuilderPortals{id: id}
}

func (p *ModelBuilderPortals) ID() string {
	return p.id
}

func (p *ModelBuilderPortals) SetID(id string) *ModelBuilderPortals {
	p.id = id
	return p
}

func (p *ModelBuilderPortals) UID() string {
	return "presetsModel" + p.id
}

func (p *ModelBuilderPortals) New(name string) string {
	if name != "" {
		name = "--" + name
	}
	return p.UID() + name
}

func (p *ModelBuilderPortals) Action() *ListingPortals {
	return &ListingPortals{uid: p.New("action")}
}

func (p *ModelBuilderPortals) Temp() string {
	return p.New("Temp")
}
