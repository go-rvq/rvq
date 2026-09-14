package cli

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"html"
	"io"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strconv"
	"strings"

	"github.com/go-rvq/rvq/cli"
	"github.com/go-rvq/rvq/web"
	"github.com/go-rvq/rvq/x/login"
)

// Config wires the application-specific parts the http_api command needs.
// Handler builds the app's HTTP handler on the given mux (the request is served
// against it in-process); FindUser resolves an account name to the app's user
// object (e.g. login.Builder.FindUserByAccount).
type Config struct {
	Handler  func(mux *http.ServeMux) http.Handler
	FindUser func(account string) (any, error)
}

// HttpApiCommand runs an HTTP request against the application handler in the same
// process, for use by an operator on the server host.
//
// The request is tagged with web.RequestSourceCLI so handlers can tell it apart
// from network traffic. With --user, it runs as that account through
// login.WithTrustedUser, the in-process equivalent of an authenticated session.
//
// A JSON body is flattened into dot/bracket form fields
// (e.g. {"Config":{"TypeID":2}} -> Config.TypeID=2) and sent as multipart form
// data, since the admin event handlers read form values rather than JSON. A key
// of the form "file:KEY" attaches the file at its string value (a path on disk)
// as the multipart file field KEY, read directly from disk and left in place
// (not treated as a temporary to remove after the request).
//
// In the response, the flash portal (see flashPortalName) is handled specially:
// its snackbar text is stripped of HTML and reported under "flash", and the
// portal is dropped from the encoded response.
//
// A spec may carry "expectedResponse" to assert the result (status code, body
// equal/contains/starts/ends, and required JSON keys); a failed assertion stops
// the run and exits non-zero before the next request.
func HttpApiCommand(cfg Config) *cli.Command {
	var (
		method      string
		userAccount string
		dataArg     string
		contentType string
		rawOut      bool
	)

	return &cli.Command{
		Name:        "http_api",
		Usage:       "URI | [FILE.json...]",
		Description: "runs one or more HTTP requests against the app in the same process",
		New: func(ctx *cli.CommandContext) error {
			fs := ctx.Flags()
			fs.StringVar(&method, "http-method", http.MethodGet, "HTTP method (default for JSON specs)")
			fs.StringVar(&userAccount, "user", "", "account name to run the request as (default for JSON specs)")
			fs.StringVar(&dataArg, "data", "", "request body; '@FILE' reads a file, '@-' or omitted reads piped stdin")
			fs.StringVar(&contentType, "content-type", "application/json", "request Content-Type (default for JSON specs)")
			fs.BoolVar(&rawOut, "raw", false, "print the raw response body without extracting the flash portal")
			return nil
		},
		Run: func(ctx *cli.CommandContext) error {
			o := httpApiOptions{
				method:      strings.ToUpper(method),
				userAccount: userAccount,
				data:        dataArg,
				contentType: contentType,
				raw:         rawOut,
			}
			args := ctx.Args
			// A single positional that is not an existing file is the URI: the
			// simple form `http_api [--http-method M] [--user U] URI`.
			if len(args) == 1 && !fileExists(args[0]) {
				o.uri = args[0]
				return runSingle(ctx, cfg, o)
			}
			// Otherwise requests come as JSON from FILE(s) or stdin: a JSON array
			// runs many requests, a JSON object runs one.
			return runFromJSON(ctx, cfg, o, args)
		},
	}
}

type httpApiOptions struct {
	method      string
	userAccount string
	data        string
	contentType string
	uri         string
	raw         bool
}

// requestSpec is one request described in JSON. It carries everything a request
// needs — the user login, the HTTP method and the URI — so no command flag is
// required; any omitted field falls back to the corresponding flag default.
type requestSpec struct {
	Method           string            `json:"method"`
	URI              string            `json:"uri"`
	Login            string            `json:"login"`
	User             string            `json:"user"` // alias for login
	ContentType      string            `json:"contentType"`
	Body             json.RawMessage   `json:"body"`
	ExpectedResponse *expectedResponse `json:"expectedResponse"`
}

// account returns the user login for this spec: "login", or the "user" alias.
func (s requestSpec) account() string {
	return coalesce(s.Login, s.User)
}

// expectedResponse validates a served response. Every field set must hold, or
// the command fails (non-zero exit).
type expectedResponse struct {
	Status *int       `json:"status"` // exact HTTP status code
	Body   *bodyMatch `json:"body"`   // string checks on the raw response body
	Keys   []string   `json:"keys"`   // dot-paths that must exist in the JSON response
}

// bodyMatch holds the string checks for the response body; each set field must
// hold.
type bodyMatch struct {
	Equal    *string `json:"equal"`
	Contains *string `json:"contains"`
	Starts   *string `json:"starts"`
	Ends     *string `json:"ends"`
}

// dispatchResult is what one served request reports back.
type dispatchResult struct {
	Status   int      `json:"status"`
	URI      string   `json:"uri,omitempty"`
	Flash    []string `json:"flash"`
	Response any      `json:"response,omitempty"`
	Raw      string   `json:"raw,omitempty"`
}

func runSingle(ctx *cli.CommandContext, cfg Config, o httpApiOptions) error {
	if cfg.Handler == nil {
		return fmt.Errorf("http_api: no in-process handler configured")
	}

	rawBody, err := loadHTTPBody(o.data)
	if err != nil {
		return fmt.Errorf("read body: %w", err)
	}

	res, err := serve(ctx.Context, cfg, o.method, o.uri, o.userAccount, o.contentType, rawBody, o.raw, nil)
	if err != nil {
		return err
	}
	res.URI = "" // single mode: the URI is obvious from the invocation
	return writeJSON(ctx.Out, res)
}

func runFromJSON(ctx *cli.CommandContext, cfg Config, o httpApiOptions, files []string) error {
	if cfg.Handler == nil {
		return fmt.Errorf("http_api: no in-process handler configured")
	}

	specs, array, err := gatherSpecs(files)
	if err != nil {
		return err
	}
	if len(specs) == 0 {
		return fmt.Errorf("no request given (provide a URI, a JSON file, or JSON on stdin)")
	}

	results := make([]*dispatchResult, 0, len(specs))
	for i, spec := range specs {
		method := strings.ToUpper(coalesce(spec.Method, o.method))
		user := coalesce(spec.account(), o.userAccount)
		ct := coalesce(spec.ContentType, o.contentType)
		if spec.URI == "" {
			return fmt.Errorf("request %d: missing \"uri\"", i)
		}

		var body []byte
		if len(spec.Body) > 0 && !bytes.Equal(bytes.TrimSpace(spec.Body), []byte("null")) {
			body = spec.Body
		}

		res, err := serve(ctx.Context, cfg, method, spec.URI, user, ct, body, o.raw, spec.ExpectedResponse)
		if err != nil {
			return fmt.Errorf("request %d (%s): %w", i, spec.URI, err)
		}
		if array {
			res.URI = spec.URI
		}
		results = append(results, res)
	}

	// A JSON array of requests reports a JSON array of results; a single JSON
	// object reports a single result object.
	if array {
		return writeJSON(ctx.Out, results)
	}
	return writeJSON(ctx.Out, results[0])
}

// serve builds one request, dispatches it into the app handler in-process, and
// returns the processed result. rawBody is the request body before wire
// encoding: when it is JSON it is flattened into multipart form fields.
func serve(baseCtx context.Context, cfg Config, method, uri, user, contentType string, rawBody []byte, raw bool, expected *expectedResponse) (*dispatchResult, error) {
	wireBody, wireContentType, err := prepareBody(contentType, rawBody)
	if err != nil {
		return nil, err
	}

	c := baseCtx
	if c == nil {
		c = context.Background()
	}
	c = web.WithRequestSource(c, web.RequestSourceCLI)

	if user != "" {
		if cfg.FindUser == nil {
			return nil, fmt.Errorf("http_api: --user given but no user finder configured")
		}
		u, err := cfg.FindUser(user)
		if err != nil {
			return nil, fmt.Errorf("user %q: %w", user, err)
		}
		c = login.WithTrustedUser(c, u)
	}

	var bodyReader io.Reader
	if len(wireBody) > 0 {
		bodyReader = bytes.NewReader(wireBody)
	}

	req := httptest.NewRequest(method, uri, bodyReader).WithContext(c)
	if wireContentType != "" && len(wireBody) > 0 {
		req.Header.Set("Content-Type", wireContentType)
	}

	rec := httptest.NewRecorder()
	cfg.Handler(http.NewServeMux()).ServeHTTP(rec, req)

	respBody := rec.Body.Bytes()
	respType := rec.Result().Header.Get("Content-Type")

	// Validate before continuing: a failed expectation ends the run with an
	// error (non-zero exit).
	if err := validateResponse(expected, rec.Code, respBody); err != nil {
		return nil, err
	}

	if raw || !isJSONContentType(respType) {
		return &dispatchResult{Status: rec.Code, Raw: string(respBody)}, nil
	}
	return processResponse(rec.Code, respBody), nil
}

// validateResponse checks a served response against expected. It reports the
// first failed check; nil expected always passes.
func validateResponse(expected *expectedResponse, status int, respBody []byte) error {
	if expected == nil {
		return nil
	}

	if expected.Status != nil && status != *expected.Status {
		return fmt.Errorf("status = %d, expected %d", status, *expected.Status)
	}

	if b := expected.Body; b != nil {
		s := string(respBody)
		if b.Equal != nil && s != *b.Equal {
			return fmt.Errorf("body != expected (equal)")
		}
		if b.Contains != nil && !strings.Contains(s, *b.Contains) {
			return fmt.Errorf("body does not contain %q", *b.Contains)
		}
		if b.Starts != nil && !strings.HasPrefix(s, *b.Starts) {
			return fmt.Errorf("body does not start with %q", *b.Starts)
		}
		if b.Ends != nil && !strings.HasSuffix(s, *b.Ends) {
			return fmt.Errorf("body does not end with %q", *b.Ends)
		}
	}

	if len(expected.Keys) > 0 {
		var root any
		if err := json.Unmarshal(respBody, &root); err != nil {
			return fmt.Errorf("cannot check keys: response is not JSON: %w", err)
		}
		for _, k := range expected.Keys {
			if !keyExists(root, k) {
				return fmt.Errorf("expected key %q not found in response", k)
			}
		}
	}

	return nil
}

// keyExists reports whether the dot-separated path exists in a decoded JSON
// value, walking object members segment by segment.
func keyExists(root any, path string) bool {
	cur := root
	for _, seg := range strings.Split(path, ".") {
		m, ok := cur.(map[string]any)
		if !ok {
			return false
		}
		v, ok := m[seg]
		if !ok {
			return false
		}
		cur = v
	}
	return true
}

// prepareBody turns a request body into its wire form. A JSON body is flattened
// into dot/bracket multipart form fields (the admin reads form values, not
// JSON); anything else is sent as-is.
func prepareBody(contentType string, rawBody []byte) (body []byte, wireContentType string, err error) {
	if len(rawBody) == 0 || !isJSONContentType(contentType) {
		return rawBody, contentType, nil
	}
	var decoded any
	if err = json.Unmarshal(rawBody, &decoded); err != nil {
		return nil, "", fmt.Errorf("parse JSON body: %w", err)
	}
	values := url.Values{}
	var files []formFile
	flattenJSON("", decoded, values, &files)
	return encodeMultipart(values, files)
}

// processResponse decodes the JSON EventResponse, lifts the flash portal text out
// (HTML stripped) into Flash and removes that portal from the response.
func processResponse(status int, respBody []byte) *dispatchResult {
	res := &dispatchResult{Status: status, Flash: []string{}}

	var resp map[string]any
	if err := json.Unmarshal(respBody, &resp); err != nil {
		res.Raw = string(respBody)
		return res
	}

	if ups, ok := resp["updatePortals"].([]any); ok {
		kept := make([]any, 0, len(ups))
		for _, up := range ups {
			m, _ := up.(map[string]any)
			if m != nil {
				if name, _ := m["name"].(string); name == flashPortalName {
					if b, _ := m["body"].(string); b != "" {
						if text := stripHTML(b); text != "" {
							res.Flash = append(res.Flash, text)
						}
					}
					continue // remove the flash portal from the response
				}
			}
			kept = append(kept, up)
		}
		if len(kept) == 0 {
			delete(resp, "updatePortals")
		} else {
			resp["updatePortals"] = kept
		}
	}

	res.Response = resp
	return res
}

// flashPortalName mirrors presets.FlashPortalName ("flash"). It is duplicated
// here to keep this command free of an import cycle with the presets package.
const flashPortalName = "flash"

var htmlTagRe = regexp.MustCompile(`(?s)<[^>]*>`)

// stripHTML removes HTML tags and unescapes entities, collapsing whitespace, so
// a rendered flash snackbar becomes its plain message text.
func stripHTML(s string) string {
	s = htmlTagRe.ReplaceAllString(s, " ")
	s = html.UnescapeString(s)
	return strings.Join(strings.Fields(s), " ")
}

func isJSONContentType(ct string) bool {
	return strings.Contains(strings.ToLower(ct), "json")
}

// loadHTTPBody resolves the --data flag: "@FILE" reads a file, "@-" reads stdin,
// an empty value reads stdin when it is piped (not a terminal), otherwise the
// literal string is the body.
func loadHTTPBody(data string) ([]byte, error) {
	switch {
	case data == "@-":
		return io.ReadAll(os.Stdin)
	case strings.HasPrefix(data, "@"):
		return os.ReadFile(data[1:])
	case data != "":
		return []byte(data), nil
	default:
		if fi, err := os.Stdin.Stat(); err == nil && fi.Mode()&os.ModeCharDevice == 0 {
			return io.ReadAll(os.Stdin)
		}
		return nil, nil
	}
}

// formFile is a multipart file part: the flattened field name and the path on
// disk whose contents are uploaded under it.
type formFile struct {
	Field string
	Path  string
}

const fileKeyPrefix = "file:"

// flattenJSON walks a decoded JSON value into form fields: object members join
// with '.', array elements with '[i]', matching the admin's form field naming.
// A key "file:KEY" with a string value is collected into files as a file part
// (field KEY) instead of a plain form value.
func flattenJSON(prefix string, v any, out url.Values, files *[]formFile) {
	switch t := v.(type) {
	case map[string]any:
		for k, val := range t {
			name := k
			isFile := false
			if strings.HasPrefix(k, fileKeyPrefix) {
				name = k[len(fileKeyPrefix):]
				isFile = true
			}
			key := name
			if prefix != "" {
				key = prefix + "." + name
			}
			if isFile {
				if p, ok := val.(string); ok {
					*files = append(*files, formFile{Field: key, Path: p})
				} else {
					out.Add(key, fmt.Sprintf("%v", val))
				}
				continue
			}
			flattenJSON(key, val, out, files)
		}
	case []any:
		for i, val := range t {
			flattenJSON(prefix+"["+strconv.Itoa(i)+"]", val, out, files)
		}
	case nil:
		out.Add(prefix, "")
	case bool:
		out.Add(prefix, strconv.FormatBool(t))
	case float64:
		out.Add(prefix, strconv.FormatFloat(t, 'f', -1, 64))
	case string:
		out.Add(prefix, t)
	default:
		out.Add(prefix, fmt.Sprintf("%v", t))
	}
}

func encodeMultipart(values url.Values, files []formFile) (body []byte, contentType string, err error) {
	var buf bytes.Buffer
	mw := multipart.NewWriter(&buf)

	keys := make([]string, 0, len(values))
	for k := range values {
		keys = append(keys, k)
	}
	sort.Strings(keys) // stable output

	for _, k := range keys {
		for _, v := range values[k] {
			if err = mw.WriteField(k, v); err != nil {
				return nil, "", err
			}
		}
	}

	// Read each "file:" entry straight from disk into its multipart part; the
	// source path is left untouched (it is not a temporary upload to clean up).
	for _, f := range files {
		if err = writeFilePart(mw, f); err != nil {
			return nil, "", err
		}
	}

	if err = mw.Close(); err != nil {
		return nil, "", err
	}
	return buf.Bytes(), mw.FormDataContentType(), nil
}

func writeFilePart(mw *multipart.Writer, f formFile) error {
	src, err := os.Open(f.Path)
	if err != nil {
		return fmt.Errorf("file field %q: %w", f.Field, err)
	}
	defer src.Close()

	part, err := mw.CreateFormFile(f.Field, filepath.Base(f.Path))
	if err != nil {
		return fmt.Errorf("file field %q: %w", f.Field, err)
	}
	if _, err = io.Copy(part, src); err != nil {
		return fmt.Errorf("file field %q: %w", f.Field, err)
	}
	return nil
}

// gatherSpecs collects request specs from the given JSON files, or from stdin
// when none is given. It also reports whether the requests should be presented
// as an array: a top-level JSON array, or more than one input file, means many.
func gatherSpecs(files []string) (specs []requestSpec, array bool, err error) {
	switch len(files) {
	case 0:
		data, err := io.ReadAll(os.Stdin)
		if err != nil {
			return nil, false, fmt.Errorf("read stdin: %w", err)
		}
		return parseSpecsMode(data)
	case 1:
		data, err := os.ReadFile(files[0])
		if err != nil {
			return nil, false, fmt.Errorf("read %s: %w", files[0], err)
		}
		specs, array, err := parseSpecsMode(data)
		if err != nil {
			return nil, false, fmt.Errorf("%s: %w", files[0], err)
		}
		return specs, array, nil
	default:
		// Several files are inherently several requests.
		for _, f := range files {
			data, err := os.ReadFile(f)
			if err != nil {
				return nil, false, fmt.Errorf("read %s: %w", f, err)
			}
			s, _, err := parseSpecsMode(data)
			if err != nil {
				return nil, false, fmt.Errorf("%s: %w", f, err)
			}
			specs = append(specs, s...)
		}
		return specs, true, nil
	}
}

// parseSpecsMode decodes a JSON array of request specs (array=true) or a single
// spec object (array=false).
func parseSpecsMode(data []byte) (specs []requestSpec, array bool, err error) {
	trimmed := bytes.TrimSpace(data)
	if len(trimmed) == 0 {
		return nil, false, nil
	}
	if trimmed[0] == '[' {
		if err := json.Unmarshal(trimmed, &specs); err != nil {
			return nil, false, fmt.Errorf("parse JSON array of requests: %w", err)
		}
		return specs, true, nil
	}
	var spec requestSpec
	if err := json.Unmarshal(trimmed, &spec); err != nil {
		return nil, false, fmt.Errorf("parse JSON request: %w", err)
	}
	return []requestSpec{spec}, false, nil
}

func fileExists(name string) bool {
	fi, err := os.Stat(name)
	return err == nil && !fi.IsDir()
}

func coalesce(v, fallback string) string {
	if v != "" {
		return v
	}
	return fallback
}

func writeJSON(out io.Writer, v any) error {
	enc := json.NewEncoder(out)
	enc.SetIndent("", "  ")
	return enc.Encode(v)
}
