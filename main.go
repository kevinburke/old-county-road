package main

import (
	"bytes"
	"html/template"
	"log"
	"net"
	"net/http"
	"os"
	"regexp"
	"strings"
	"time"

	"github.com/kevinburke/go-html-template/assets"
	"github.com/kevinburke/handlers"
	"github.com/kevinburke/rest"
)

var homepageTpl *template.Template

func init() {
	homepageHTML := assets.MustAssetString("templates/index.html")
	homepageTpl = template.Must(template.New("homepage").Parse(homepageHTML))

	// Add more templates here.
}

// Static file HTTP server; all assets are packaged up in the assets directory
// with go-bindata.
type static struct {
	modTime time.Time
}

func (s *static) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path == "/favicon.ico" {
		r.URL.Path = "/static/favicon.ico"
	}
	bits, err := assets.Asset(strings.TrimPrefix(r.URL.Path, "/"))
	if err != nil {
		rest.NotFound(w, r)
		return
	}
	http.ServeContent(w, r, r.URL.Path, s.modTime, bytes.NewReader(bits))
}

func render(w http.ResponseWriter, tpl *template.Template, name string, data interface{}) {
	buf := new(bytes.Buffer)
	if err := tpl.ExecuteTemplate(w, name, data); err != nil {
		http.Error(w, err.Error(), 500)
	}
	w.Write(buf.Bytes())
}

func NewServeMux() http.Handler {
	staticServer := &static{
		modTime: time.Now().UTC(),
	}

	r := new(handlers.Regexp)
	r.Handle(regexp.MustCompile(`(^/static|^/favicon.ico$)`), []string{"GET"}, handlers.GZip(staticServer))
	r.HandleFunc(regexp.MustCompile(`^/$`), []string{"GET"}, func(w http.ResponseWriter, r *http.Request) {
		render(w, homepageTpl, "homepage", nil)
	})
	// Add more routes here.
	return r
}

func main() {
	mux := NewServeMux()
	port, ok := os.LookupEnv("PORT")
	if !ok {
		port = "7065"
	}
	ln, err := net.Listen("tcp", ":"+port)
	if err != nil {
		log.Fatal(err)
	}
	log.Println("Listening on port", port)
	http.Serve(ln, handlers.Duration(handlers.Log(handlers.Debug(handlers.UUID(mux)))))
}
