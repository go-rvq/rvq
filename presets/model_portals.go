package presets

type ModelBuilderPortals struct {
	id string
}

func NewModelBuilderPortals(id string) *ModelBuilderPortals {
	return &ModelBuilderPortals{id: id}
}

func (b *ModelBuilderPortals) ID() string {
	return b.id
}

func (b *ModelBuilderPortals) SetID(id string) *ModelBuilderPortals {
	b.id = id
	return b
}

func (b *ModelBuilderPortals) UID() string {
	return "presetsModel" + b.id
}

func (b *ModelBuilderPortals) Listing() *ModelBuilderListingPortals {
	return &ModelBuilderListingPortals{uid: b.UID()}
}

type ModelBuilderListingPortals struct {
	uid string
}

func NewListingPortals(id string) *ModelBuilderListingPortals {
	return &ModelBuilderListingPortals{uid: id}
}

func (p *ModelBuilderListingPortals) ID() string {
	return p.uid
}

func (p *ModelBuilderListingPortals) SetID(id string) *ModelBuilderListingPortals {
	p.uid = id
	return p
}

func (p *ModelBuilderListingPortals) Name(name string) string {
	if name != "" {
		name = "--" + name
	}
	return p.uid + "--listing" + name
}

func (p *ModelBuilderListingPortals) Main() string {
	return p.Name("")
}

func (p *ModelBuilderListingPortals) DataTable() string {
	return p.Name("DataTable")
}

func (p *ModelBuilderListingPortals) DataTableAdditions() string {
	return p.DataTable() + "Additions"
}

func (p *ModelBuilderListingPortals) Record() string {
	return p.Name("Record")
}
