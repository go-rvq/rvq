package schemaform

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"sort"
	"strings"
	"sync"
	"time"

	"github.com/gad-lang/gad"
)

// Request is a request of an address of the web a schema makes — the
// values of a field's `[options=…]` of a source of its own —:
//
//	Request("https://…/cities"; query={state: this.state}).then((r) => r.data)
//
// Request(url; method="GET", query={…}, headers={…}), and .then(fn): a
// request whose answer is fn of this one's — of the first, the response read
// as JSON. It is resolved by Go (Resolve): on the server, with its HTTP
// client; in the browser (WASM), the same Go — Go's HTTP there is the
// browser's fetch —, the page waiting on it as on a promise. A response is
// kept a while (RequestCacheTTL): the same address is asked once.
type Request struct {
	Method  string
	URL     string
	Headers map[string]string
	thens   []gad.Object
}

var (
	// RequestTimeout bounds a request.
	RequestTimeout = 5 * time.Second
	// RequestCacheTTL is how long a response is kept.
	RequestCacheTTL = 10 * time.Minute
	// RequestMaxBody is the most of a response read.
	RequestMaxBody int64 = 5 << 20
	// RequestClient is the HTTP client of the requests.
	RequestClient = &http.Client{}

	requestCache sync.Map // key → *cachedResponse
)

type cachedResponse struct {
	at   time.Time
	body any
}

// RequestType is the type of a Request; RequestFunc, what a schema calls it
// by (`Request(url; …)`, among the Builtins).
var (
	RequestType = gad.NewBuiltinObjType("Request")
	RequestFunc = gad.NewFunction("Request", newRequest)
)

func newRequest(c gad.Call) (gad.Object, error) {
	if err := c.Args.CheckLen(1); err != nil {
		return nil, err
	}
	if err := c.NamedArgs.CheckNames("method", "query", "headers"); err != nil {
		return nil, fmt.Errorf("Request: %w", err)
	}
	r := &Request{Method: http.MethodGet, URL: c.Args.Get(0).ToString()}
	if m := c.NamedArgs.GetValueOrNil("method"); m != nil && m != gad.Nil {
		r.Method = strings.ToUpper(m.ToString())
	}
	u, err := url.Parse(r.URL)
	if err != nil {
		return nil, fmt.Errorf("Request: %w", err)
	}
	if u.Scheme != "http" && u.Scheme != "https" {
		return nil, fmt.Errorf("Request: %q: only an address of http or https", r.URL)
	}
	if q := c.NamedArgs.GetValueOrNil("query"); q != nil && q != gad.Nil {
		values := u.Query()
		for _, kv := range pairsOf(q) {
			if kv.V != gad.Nil {
				values.Set(kv.K.ToString(), kv.V.ToString())
			}
		}
		u.RawQuery = values.Encode()
		r.URL = u.String()
	}
	if h := c.NamedArgs.GetValueOrNil("headers"); h != nil && h != gad.Nil {
		r.Headers = map[string]string{}
		for _, kv := range pairsOf(h) {
			r.Headers[kv.K.ToString()] = kv.V.ToString()
		}
	}
	return r, nil
}

// pairsOf is a dict or a key-value array as its pairs (a dict's by key).
func pairsOf(o gad.Object) gad.KeyValueArray {
	switch v := o.(type) {
	case gad.KeyValueArray:
		return v
	case gad.Dict:
		keys := make([]string, 0, len(v))
		for k := range v {
			keys = append(keys, k)
		}
		sort.Strings(keys)
		out := make(gad.KeyValueArray, len(keys))
		for i, k := range keys {
			out[i] = &gad.KeyValue{K: gad.Str(k), V: v[k]}
		}
		return out
	}
	return nil
}

func (r *Request) Type() gad.ObjectType { return RequestType }
func (r *Request) ToString() string     { return "Request(" + r.Method + " " + r.URL + ")" }
func (r *Request) IsFalsy() bool        { return false }
func (r *Request) Equal(right gad.Object) bool {
	o, ok := right.(*Request)
	return ok && o == r
}

// IndexGet is what a schema reads of a request: then(fn), its url, its
// method.
func (r *Request) IndexGet(_ *gad.VM, index gad.Object) (gad.Object, error) {
	switch index.ToString() {
	case "then":
		return gad.NewFunction("then", func(c gad.Call) (gad.Object, error) {
			if err := c.Args.CheckLen(1); err != nil {
				return nil, err
			}
			next := *r
			next.thens = append(append([]gad.Object{}, r.thens...), c.Args.Get(0))
			return &next, nil
		}), nil
	case "url":
		return gad.Str(r.URL), nil
	case "method":
		return gad.Str(r.Method), nil
	}
	return nil, gad.ErrInvalidIndex.NewError(index.ToString())
}

// Resolve is what the request answers: the response, as JSON, through its
// thens — a then that answers a request, that one resolved.
func (r *Request) Resolve(ctx context.Context, vm *gad.VM) (gad.Object, error) {
	body, err := r.fetch(ctx)
	if err != nil {
		return nil, err
	}
	obj := JSONObject(body)
	for _, fn := range r.thens {
		if obj, err = gad.NewInvoker(vm, fn).WithContext(ctx).Invoke(gad.Args{gad.Array{obj}}, nil); err != nil {
			return nil, fmt.Errorf("%s: then: %w", r.ToString(), err)
		}
		if next, ok := obj.(*Request); ok {
			if obj, err = next.Resolve(ctx, vm); err != nil {
				return nil, err
			}
		}
	}
	return obj, nil
}

func (r *Request) key() string {
	var b strings.Builder
	b.WriteString(r.Method + " " + r.URL)
	keys := make([]string, 0, len(r.Headers))
	for k := range r.Headers {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	for _, k := range keys {
		b.WriteString("\n" + k + ": " + r.Headers[k])
	}
	return b.String()
}

// fetch is the response of the request, as JSON: of the cache, while it
// holds it.
func (r *Request) fetch(ctx context.Context) (any, error) {
	key := r.key()
	if c, ok := requestCache.Load(key); ok && time.Since(c.(*cachedResponse).at) < RequestCacheTTL {
		return c.(*cachedResponse).body, nil
	}
	ctx, cancel := context.WithTimeout(ctx, RequestTimeout)
	defer cancel()
	req, err := http.NewRequestWithContext(ctx, r.Method, r.URL, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Accept", "application/json")
	for k, v := range r.Headers {
		req.Header.Set(k, v)
	}
	res, err := RequestClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", r.ToString(), err)
	}
	defer res.Body.Close()
	if res.StatusCode < 200 || res.StatusCode > 299 {
		return nil, fmt.Errorf("%s: %s", r.ToString(), res.Status)
	}
	var body any
	if err := json.NewDecoder(io.LimitReader(res.Body, RequestMaxBody)).Decode(&body); err != nil {
		return nil, fmt.Errorf("%s: the response is not JSON: %w", r.ToString(), err)
	}
	requestCache.Store(key, &cachedResponse{at: time.Now(), body: body})
	return body, nil
}

// JSONObject is v, as encoding/json reads it, as gad holds it: an object a
// dict, an array an array, a number an int when it is whole.
func JSONObject(v any) gad.Object {
	switch x := v.(type) {
	case nil:
		return gad.Nil
	case bool:
		return gad.Bool(x)
	case string:
		return gad.Str(x)
	case float64:
		if x == float64(int64(x)) {
			return gad.Int(int64(x))
		}
		return gad.Float(x)
	case []any:
		out := make(gad.Array, len(x))
		for i, e := range x {
			out[i] = JSONObject(e)
		}
		return out
	case map[string]any:
		out := make(gad.Dict, len(x))
		for k, e := range x {
			out[k] = JSONObject(e)
		}
		return out
	}
	return gad.Str(fmt.Sprint(v))
}
