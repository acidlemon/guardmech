package guardmech

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"os"

	"github.com/acidlemon/guardmech/backend/app/handler"
	"github.com/gorilla/mux"
	_ "github.com/k0kubun/pp/v3" // for development
)

type GuardMech struct {
}

func New() *GuardMech {
	gm := &GuardMech{}

	return gm
}

const spaIndexFile = "dist/index.html"

func (g *GuardMech) Run() error {
	listener, err := net.Listen("tcp", "0.0.0.0:2989")
	if err != nil {
		log.Println("failed to listen")
		return nil
	}

	root := mux.NewRouter()
	var r *mux.Router
	mount := os.Getenv("GUARDMECH_MOUNT_PATH")
	if mount != "" {
		r = root.PathPrefix(mount).Subrouter()
	} else {
		r = root
	}

	authMux := handler.NewAuthMux()
	adminAPIMux := handler.NewAdminMux()

	authMux.RegisterRoute(r)
	adminAPIMux.RegisterRoute(r)

	// web assets
	r.HandleFunc("/guardmech/", func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == mount+"/guardmech/" {
			http.Redirect(w, r, mount+"/guardmech/admin/", http.StatusPermanentRedirect)
			return
		}
		http.NotFound(w, r)
	})
	// r.HandleFunc("/guardmech/admin/", func(w http.ResponseWriter, r *http.Request) {
	// 	http.ServeFile(w, r, spaIndexFile) // write SPA html
	// })
	r.PathPrefix("/guardmech/admin/js/").Handler(http.StripPrefix(mount+"/guardmech/admin/js", http.FileServer(http.Dir("dist/js"))))
	r.PathPrefix("/guardmech/admin/css/").Handler(http.StripPrefix(mount+"/guardmech/admin/css", http.FileServer(http.Dir("dist/css"))))
	r.PathPrefix("/guardmech/admin/").HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, spaIndexFile) // write SPA html
	})

	root.HandleFunc("/", func(w http.ResponseWriter, req *http.Request) {
		log.Println("catch all:", req.URL.Path)
	})

	return http.Serve(listener, root)
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
