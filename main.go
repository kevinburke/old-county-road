// go-html-boilerplate loads configuration from a file and starts a HTTP server
// that can render HTML templates and static assets.
//
// See config.yml for an explanation of the configuration options for the
// server, and the Makefile for various tasks you can run in coordination with
// the server (run tests, build assets, start the server).
package main

import (
	"bytes"
	"flag"
	"html/template"
	"io/ioutil"
	"net/http"
	"os"
	"regexp"
	"strconv"
	"strings"
	"time"

	log "github.com/inconshreveable/log15"
	"github.com/kevinburke/handlers"
	"github.com/kevinburke/nacl"
	"github.com/kevinburke/old-county-road/assets"
	"github.com/kevinburke/rest"
	yaml "gopkg.in/yaml.v2"
)

// DefaultPort is the listening port if no other port is specified.
var DefaultPort = 4965

// The server's Version.
const Version = "0.4"

var logger log.Logger

func init() {
	logger = handlers.Logger

	// Add more templates here.
}

// A HTTP server for static files. All assets are packaged up in the assets
// directory with the go-bindata binary. Run "make assets" to rerun the
// go-bindata binary.
type static struct {
	modTime time.Time
}

func (s *static) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path == "/favicon.ico" {
		r.URL.Path = "/static/favicon.ico"
	}
	if r.URL.Path == "/" {
		r.URL.Path = "/static/turnout.html"
	}
	bits, err := assets.Asset(strings.TrimPrefix(r.URL.Path, "/"))
	if err != nil {
		rest.NotFound(w, r)
		return
	}
	http.ServeContent(w, r, r.URL.Path, s.modTime, bytes.NewReader(bits))
}

// Render a template, or a server error.
func render(w http.ResponseWriter, r *http.Request, tpl *template.Template, name string, data interface{}) {
	buf := new(bytes.Buffer)
	if err := tpl.ExecuteTemplate(buf, name, data); err != nil {
		rest.ServerError(w, r, err)
		return
	}
	w.Write(buf.Bytes())
}

// NewServeMux returns a HTTP handler that covers all routes known to the
// server.
func NewServeMux() http.Handler {
	staticServer := &static{
		modTime: time.Now().UTC(),
	}

	r := new(handlers.Regexp)
	r.Handle(regexp.MustCompile(`(^/$|^/static|^/favicon.ico$)`), []string{"GET"}, handlers.GZip(staticServer))
	// Add more routes here. Routes not matched will get a 404 error page.
	// Call rest.RegisterHandler(404, http.HandlerFunc) to provide your own 404
	// page instead of the default.
	return r
}

// FileConfig represents the data in a config file.
type FileConfig struct {
	// SecretKey is used to encrypt sessions and other data before serving it to
	// the client. It should be a hex string that's exactly 64 bytes long. For
	// example:
	//
	//   d7211b215341871968869dontusethisc0ff1789fc88e0ac6e296ba36703edf8
	//
	// That key is invalid - you can generate a random key by running:
	//
	//   openssl rand -hex 32
	//
	// If no secret key is present, we'll generate one when the server starts.
	// However, this means that sessions may error when the server restarts.
	//
	// If a server key is present, but invalid, the server will not start.
	SecretKey string `yaml:"secret_key"`

	// Port to listen on. Set to 0 to choose a port at random. If unspecified,
	// defaults to 7065.
	Port *int `yaml:"port"`

	// Set to true to listen for HTTP traffic (instead of TLS traffic). Note
	// you need to terminate TLS to use HTTP server push.
	HTTPOnly bool `yaml:"http_only"`

	// For TLS configuration.
	CertFile string `yaml:"cert_file"`
	KeyFile  string `yaml:"key_file"`

	// Add other configuration settings here.
}

var cfg = flag.String("config", "config.yml", "Path to a config file")

func commonMain() (*FileConfig, http.Handler) {
	flag.Parse()
	data, err := ioutil.ReadFile(*cfg)
	c := new(FileConfig)
	if err == nil {
		if err := yaml.Unmarshal(data, c); err != nil {
			logger.Error("Couldn't parse config file", "err", err)
			os.Exit(2)
		}
	} else {
		logger.Error("Couldn't find config file", "err", err)
		os.Exit(2)
	}
	key, err := nacl.Load(c.SecretKey)
	if err != nil {
		logger.Error("Error getting secret key", "err", err)
		os.Exit(2)
	}
	// You can use the secret key with secretbox
	// (godoc.org/github.com/kevinburke/nacl/secretbox/) to generate cookies and
	// secrets. See flash.go and crypto.go for examples.
	_ = key

	if c.Port == nil {
		port, ok := os.LookupEnv("PORT")
		if ok {
			iPort, err := strconv.Atoi(port)
			if err != nil {
				logger.Error("Invalid port", "err", err, "port", port)
				os.Exit(2)
			}
			c.Port = &iPort
		} else {
			c.Port = &DefaultPort
		}
	}
	mux := NewServeMux()
	mux = handlers.UUID(mux)                               // add UUID header
	mux = handlers.Server(mux, "old-county-road/"+Version) // add Server header
	mux = handlers.Log(mux)                                // log requests/responses
	mux = handlers.Duration(mux)                           // add Duration header
	mux = handlers.STS(mux)                                // add STS header
	mux = handlers.RedirectProto(mux)                      // add STS header
	return c, mux
}
