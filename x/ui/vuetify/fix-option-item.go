package vuetify

type DefaultOptionItem struct {
	Text     string               `json:"text"`
	Value    string               `json:"value"`
	Children []*DefaultOptionItem `json:"children"`
}
