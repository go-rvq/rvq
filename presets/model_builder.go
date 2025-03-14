package presets

import "github.com/qor5/x/v3/i18n"

type ModelBuilderConfigAttributes struct {
	id           string
	label        string
	pluralLabel  string
	uriName      string
	singleton    bool
	dataOperator DataOperator
	moduleKey    i18n.ModuleKey
}

func (c *ModelBuilderConfigAttributes) ModuleKey() i18n.ModuleKey {
	return c.moduleKey
}

func (c *ModelBuilderConfigAttributes) SetModuleKey(moduleKey i18n.ModuleKey) {
	c.moduleKey = moduleKey
}

type ModelBuilderConfig struct {
	ModelBuilderConfigAttributes
}

func (c *ModelBuilderConfig) Apply(mb *ModelBuilder) {
	mb.ModelBuilderConfigAttributes = c.ModelBuilderConfigAttributes
}

func (c *ModelBuilderConfig) Id() string {
	return c.id
}

func (c *ModelBuilderConfig) SetId(id string) *ModelBuilderConfig {
	c.id = id
	return c
}

func (c *ModelBuilderConfig) Label() string {
	return c.label
}

func (c *ModelBuilderConfig) SetLabel(label string) *ModelBuilderConfig {
	c.label = label
	return c
}

func (c *ModelBuilderConfig) PluralLabel() string {
	return c.pluralLabel
}

func (c *ModelBuilderConfig) SetPluralLabel(pluralLabel string) *ModelBuilderConfig {
	c.pluralLabel = pluralLabel
	return c
}

func (c *ModelBuilderConfig) UriName() string {
	return c.uriName
}

func (c *ModelBuilderConfig) SetUriName(uriName string) *ModelBuilderConfig {
	c.uriName = uriName
	return c
}

func (c *ModelBuilderConfig) Singleton() bool {
	return c.singleton
}

func (c *ModelBuilderConfig) SetSingleton(singleton bool) *ModelBuilderConfig {
	c.singleton = singleton
	return c
}

func (c *ModelBuilderConfig) DataOperator() DataOperator {
	return c.dataOperator
}

func (c *ModelBuilderConfig) SetDataOperator(dataOperator DataOperator) *ModelBuilderConfig {
	c.dataOperator = dataOperator
	return c
}

func (c *ModelBuilderConfig) ModuleKey() i18n.ModuleKey {
	return c.moduleKey
}

func (c *ModelBuilderConfig) SetModuleKey(moduleKey i18n.ModuleKey) *ModelBuilderConfig {
	c.moduleKey = moduleKey
	return c
}

func ModelConfig() *ModelBuilderConfig {
	return &ModelBuilderConfig{}
}
