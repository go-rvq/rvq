package base

type Configor interface {
	getConfig() *Config
}

type Config struct {
	// OriginalMaxWidth × OriginalMaxHeight is the largest an image's
	// original is kept in (0: no limit on that side); turned with the image —
	// a portrait one fits it on its side.
	OriginalMaxWidth  int
	OriginalMaxHeight int
}

// originalMax is the limit of the original of an image w × h: the config's,
// turned when the image is turned the other way.
func (c *Config) originalMax(w, h int) (mw, mh int) {
	mw, mh = c.OriginalMaxWidth, c.OriginalMaxHeight
	if mw > 0 && mh > 0 && (w > h) != (mw > mh) && w != h && mw != mh {
		mw, mh = mh, mw
	}
	return
}

type WithConfigField struct {
	Config Config
}

func (c *WithConfigField) getConfig() *Config {
	return &c.Config
}
