package presets

func (p *ModelBuilderPortals) Listing() *ListingPortals {
	return &ListingPortals{model: p, uid: p.New("listing")}
}

type ListingPortals struct {
	model *ModelBuilderPortals
	uid   string
}

func (p *ListingPortals) Model() *ModelBuilderPortals {
	return p.model
}

func NewListingPortals(uid string) *ListingPortals {
	return &ListingPortals{uid: uid}
}

func (b *ListingBuilder) Portals(portalID string) *ListingPortals {
	return NewModelBuilderPortals(portalID).Listing()
}

func (p *ListingPortals) UID() string {
	return p.uid
}

func (p *ListingPortals) SetUID(v string) *ListingPortals {
	p.uid = v
	return p
}

func (p *ListingPortals) New(name string) string {
	if name != "" {
		name = "--" + name
	}
	return p.uid + name
}

func (p *ListingPortals) Main() string {
	return p.New("")
}

func (p *ListingPortals) DataTable() string {
	return p.New("DataTable")
}

func (p *ListingPortals) DataTableAdditions() string {
	return p.DataTable() + "Additions"
}

func (p *ListingPortals) Temp() string {
	return p.New("Temp")
}
