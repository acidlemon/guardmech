package guardmech

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"log"
	"net"
	"net/http"

	"github.com/acidlemon/guardmech/backend/app/handler"
	_ "github.com/k0kubun/pp/v3" // for development
)

type GuardMech struct {
}

func New() *GuardMech {
	gm := &GuardMech{}

	return gm
}

func (g *GuardMech) Run() error {
	listener, err := net.Listen("tcp", "0.0.0.0:2989")
	if err != nil {
		log.Println("failed to listen")
		return nil
	}

	adminWebMux := http.NewServeMux()
	// web assets
	adminWebMux.HandleFunc("/guardmech/", func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/guardmech/" {
			http.Redirect(w, r, "/guardmech/admin/", http.StatusPermanentRedirect)
			return
		}
		http.NotFound(w, r)
	})
	adminWebMux.HandleFunc("/guardmech/admin/", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "dist/index.html") // write SPA html
	})
	adminWebMux.Handle("/guardmech/admin/js/", http.StripPrefix("/guardmech/admin/js/", http.FileServer(http.Dir("dist/js"))))
	adminWebMux.Handle("/guardmech/admin/css/", http.StripPrefix("/guardmech/admin/css/", http.FileServer(http.Dir("dist/css"))))

	authMux := handler.NewAuthMux()
	adminAPIMux := handler.NewAdminMux()

	mux := http.NewServeMux()
	mux.Handle("/auth/", authMux.Mux())
	mux.Handle("/guardmech/", adminWebMux)
	mux.Handle("/guardmech/api/", adminAPIMux.Mux())
	mux.HandleFunc("/", func(w http.ResponseWriter, req *http.Request) {
		log.Println("catch all:", req.URL.Path)
	})

	return http.Serve(listener, mux)
}

// func (g *GuardMech) ReverseProxy(w http.ResponseWriter, req *http.Request) {
// }

func WrapServerError(res *http.Response, err error) {
	res.StatusCode = http.StatusInternalServerError
	log.Println("server error: ", err)

	html := `<!doctype html>
<html>
<head>
  <title>500 Internal Server Error</title>
</head>
<body>
<h1>Internal Server Error</h1>
<p>guardmech has panicked.</p>
<p>reason: %s</p>
</body>
</html>
`
	b := []byte(fmt.Sprintf(html, err.Error()))
	res.Body = ioutil.NopCloser(bytes.NewReader(b))
}
