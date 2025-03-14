package base

type Configor interface {
	getConfig() *Config
}

type Config struct {
	OriginalMaxWidth  int
	OriginalMaxHeight int
}

type WithConfigField struct {
	Config Config
}

func (c *WithConfigField) getConfig() *Config {
	return &c.Config
}
