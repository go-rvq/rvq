package presets

import (
	"context"
	"errors"
	"io"
	"mime"
	"net"
	"net/http"
	"net/netip"
	"net/url"
	"strings"
	"syscall"
	"time"

	"github.com/go-rvq/rvq/web"
	"golang.org/x/net/html"
	"golang.org/x/net/html/charset"
)

// A page of another site read (LinkPreview): its title — for the text of a
// link to it left empty (rvq mdkit's link form: LinkTitleEvent) —, its
// description and its image — for its card (an external LinkCard) —. Read
// once — a small timeout, no second try —, of an address of the internet
// only: no scheme but http and https, no address of the machine, of its
// network or of the cloud's metadata — told by the address connected to, so
// a name that resolves to one, or a redirect to one, is refused too —; HTML
// only, its head only.

// LinkTitleEvent answers {title} of the page of the address posted (url);
// none known, "".
const LinkTitleEvent = "presets_LinkTitle"

// The limits of LinkTitleEvent.
var (
	// LinkTitleTimeout is how long it waits for the page, all of it.
	LinkTitleTimeout = 3 * time.Second
	// linkTitleMaxBytes is how much of the page is read.
	linkTitleMaxBytes int64 = 256 << 10
	// linkTitleMaxRedirects is how many redirects it follows.
	linkTitleMaxRedirects = 3
	// linkTitleMaxRunes is how long a title is kept.
	linkTitleMaxRunes = 300
)

// errLinkTitleAddress is an address LinkTitleEvent does not read.
var errLinkTitleAddress = errors.New("link title: not an address of the internet")

// publicAddr says ip is of the internet: not of the machine (loopback,
// unspecified), of its network (private, link-local — the clouds' metadata
// among them —, shared, unique local), nor multicast.
func publicAddr(ip netip.Addr) bool {
	ip = ip.Unmap()
	if !ip.IsValid() || ip.IsLoopback() || ip.IsUnspecified() || ip.IsPrivate() || ip.IsLinkLocalUnicast() ||
		ip.IsLinkLocalMulticast() || ip.IsInterfaceLocalMulticast() || ip.IsMulticast() {
		return false
	}
	// 100.64.0.0/10, the carriers' shared space
	if ip.Is4() && netip.MustParsePrefix("100.64.0.0/10").Contains(ip) {
		return false
	}
	return true
}

// linkTitleAllowed says LinkTitleEvent may connect to ip (publicAddr; a
// test's, its own server too).
var linkTitleAllowed = publicAddr

// linkTitleClient is the client of LinkTitleEvent: the address it connects
// to of the internet only (publicAddr), its redirects a few, of http and
// https.
var linkTitleClient = &http.Client{
	Transport: &http.Transport{
		Proxy: nil,
		DialContext: (&net.Dialer{
			Timeout: LinkTitleTimeout,
			// the address connected to, whatever name or redirect led to it
			Control: func(_, address string, _ syscall.RawConn) error {
				host, _, err := net.SplitHostPort(address)
				if err != nil {
					return err
				}
				ip, err := netip.ParseAddr(host)
				if err != nil || !linkTitleAllowed(ip) {
					return errLinkTitleAddress
				}
				return nil
			},
		}).DialContext,
		TLSHandshakeTimeout:   LinkTitleTimeout,
		ResponseHeaderTimeout: LinkTitleTimeout,
		MaxIdleConns:          4,
		IdleConnTimeout:       30 * time.Second,
	},
	CheckRedirect: func(req *http.Request, via []*http.Request) error {
		if len(via) > linkTitleMaxRedirects {
			return errors.New("link title: too many redirects")
		}
		if req.URL.Scheme != "http" && req.URL.Scheme != "https" {
			return errLinkTitleAddress
		}
		return nil
	},
}

// linkTitleURL is the address raw as LinkTitleEvent reads it: absolute, of
// http or https, of a host; else nil.
func linkTitleURL(raw string) *url.URL {
	u, err := url.Parse(strings.TrimSpace(raw))
	if err != nil || (u.Scheme != "http" && u.Scheme != "https") || u.Hostname() == "" || u.User != nil {
		return nil
	}
	return u
}

// ErrLinkPreview is a page of another site that could not be read: not of
// the internet, no answer in time, not HTML, no title.
var ErrLinkPreview = errors.New("link preview: the page could not be read")

// LinkPreview is the card of the page at raw — its address, its title, its
// summary and its cover (of https), what it says of itself to be shared
// first: of Open Graph (og:title, og:description, og:image), else of
// Twitter's card (twitter:…), else its <title> and its description —, read once within
// LinkTitleTimeout; ErrLinkPreview when it could not be read.
func LinkPreview(ctx context.Context, raw string) (*LinkCard, error) {
	u := linkTitleURL(raw)
	if u == nil {
		return nil, ErrLinkPreview
	}
	ctx, cancel := context.WithTimeout(ctx, LinkTitleTimeout)
	defer cancel()
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, u.String(), nil)
	if err != nil {
		return nil, ErrLinkPreview
	}
	req.Header.Set("Accept", "text/html,application/xhtml+xml")
	req.Header.Set("User-Agent", "Mozilla/5.0 (compatible; rvq-link-preview)")
	res, err := linkTitleClient.Do(req)
	if err != nil {
		return nil, ErrLinkPreview
	}
	defer res.Body.Close()
	ct := res.Header.Get("Content-Type")
	if mt, _, _ := mime.ParseMediaType(ct); res.StatusCode != http.StatusOK ||
		(mt != "" && mt != "text/html" && mt != "application/xhtml+xml") {
		return nil, ErrLinkPreview
	}
	body, err := charset.NewReader(io.LimitReader(res.Body, linkTitleMaxBytes), ct)
	if err != nil {
		return nil, ErrLinkPreview
	}
	h := headOf(body)
	// what the page says of itself to be shared first (Open Graph, then
	// Twitter's card), then its own title and description
	c := &LinkCard{Href: u.String(), Title: firstOf(h.ogTitle, h.twTitle, h.title),
		Summary: firstOf(h.ogDescription, h.twDescription, h.description)}
	for _, raw := range []string{h.ogImage, h.twImage} {
		// of the page's address (a redirect's); of https only
		if img, err := res.Request.URL.Parse(raw); raw != "" && err == nil && img.Scheme == "https" {
			c.Cover = img.String()
			break
		}
	}
	if c.Title == "" {
		return nil, ErrLinkPreview
	}
	return c, nil
}

// LinkTitle is the title of the page at raw (LinkPreview); "" when it could
// not be read.
func LinkTitle(ctx context.Context, raw string) string {
	if c, err := LinkPreview(ctx, raw); err == nil {
		return c.Title
	}
	return ""
}

// firstOf is the first of vs not empty.
func firstOf(vs ...string) string {
	for _, v := range vs {
		if v != "" {
			return v
		}
	}
	return ""
}

// pageHead is what the head of a page says of it.
type pageHead struct {
	title, ogTitle, twTitle, description, ogDescription, twDescription, ogImage, twImage string
}

// headOf is what the head of the HTML of r says: its first <title>, its
// metas description, og:title, og:description, og:image and twitter:title,
// twitter:description, twitter:image — each cleaned (cleanTitle) —; past the head (its body), nothing more.
func headOf(r io.Reader) (h pageHead) {
	z := html.NewTokenizer(r)
	in := false
	var title strings.Builder
	done := func() pageHead {
		h.title = cleanTitle(title.String())
		h.ogTitle, h.twTitle = cleanTitle(h.ogTitle), cleanTitle(h.twTitle)
		h.description, h.ogDescription, h.twDescription = cleanTitle(h.description), cleanTitle(h.ogDescription), cleanTitle(h.twDescription)
		h.ogImage, h.twImage = strings.TrimSpace(h.ogImage), strings.TrimSpace(h.twImage)
		return h
	}
	for {
		tt := z.Next()
		switch tt {
		case html.ErrorToken:
			return done()
		case html.StartTagToken, html.SelfClosingTagToken:
			name, hasAttr := z.TagName()
			switch string(name) {
			case "title":
				in = tt == html.StartTagToken && title.Len() == 0
			case "body":
				return done()
			case "meta":
				var key, content string
				for hasAttr {
					var k, v []byte
					k, v, hasAttr = z.TagAttr()
					switch string(k) {
					case "property", "name":
						key = strings.ToLower(string(v))
					case "content":
						content = string(v)
					}
				}
				switch key {
				case "og:title":
					h.ogTitle = firstOf(h.ogTitle, content)
				case "description":
					h.description = firstOf(h.description, content)
				case "og:description":
					h.ogDescription = firstOf(h.ogDescription, content)
				case "og:image", "og:image:secure_url":
					h.ogImage = firstOf(h.ogImage, content)
				case "twitter:title":
					h.twTitle = firstOf(h.twTitle, content)
				case "twitter:description":
					h.twDescription = firstOf(h.twDescription, content)
				case "twitter:image", "twitter:image:src":
					h.twImage = firstOf(h.twImage, content)
				}
			}
		case html.EndTagToken:
			if name, _ := z.TagName(); string(name) == "title" {
				in = false
			}
		case html.TextToken:
			if in {
				title.Write(z.Text())
			}
		}
	}
}

// titleOf is the <title> of the HTML of r (headOf); "" of none.
func titleOf(r io.Reader) string { return headOf(r).title }

// cleanTitle is t, its spaces one, cut at linkTitleMaxRunes.
func cleanTitle(t string) string {
	t = strings.Join(strings.Fields(t), " ")
	if r := []rune(t); len(r) > linkTitleMaxRunes {
		t = string(r[:linkTitleMaxRunes]) + "…"
	}
	return t
}

// linkTitle is LinkTitleEvent.
func (b *Builder) linkTitle(ctx *web.EventContext) (r web.EventResponse, err error) {
	r.Data = map[string]string{"title": LinkTitle(ctx.R.Context(), ctx.R.FormValue("url"))}
	return
}
