package sitemap

import (
	"net/http"
)

type contextKey string

const hostWithSchemeKey contextKey = "HostWithScheme"

func (site *SiteMapBuilder) MountTo(mux *http.ServeMux) {
	mux.Handle(site.pathName, site)
}

func (index *SiteMapIndexBuilder) MountTo(mux *http.ServeMux) {
	mux.Handle(index.pathName, index)
	for _, site := range index.siteMaps {
		mux.Handle(site.pathName, site)
	}
}

func (robot *RobotsBuilder) MountTo(mux *http.ServeMux) {
	mux.Handle("/robots.txt", robot)
}

func (site *SiteMapBuilder) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/xml;charset=UTF-8")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(EncodeToXmlByRequest(r, site)))
}

func (index *SiteMapIndexBuilder) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/xml;charset=UTF-8")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(EncodeToXmlByRequest(r, index)))
}

func (robot *RobotsBuilder) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/plain;charset=UTF-8")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(robot.ToTxt()))
}

func EncodeToXmlByRequest(r *http.Request, encoder EncodeToXmlInterface) string {
	var host string
	if r.URL.Host != "" {
		host = r.URL.Scheme + "://" + r.URL.Host
	}

	return encoder.EncodeToXml(WithHost(host, r.Context()))
}
