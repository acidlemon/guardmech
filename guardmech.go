package guardmech

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net"
	"net/http"

	"github.com/acidlemon/guardmech/app/handler"
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

	childMux := http.NewServeMux()
	childMux.HandleFunc("/", func(w http.ResponseWriter, req *http.Request) {
		log.Println("defualt handler")
	})
	childMux.HandleFunc("/guardmech/api/", func(w http.ResponseWriter, req *http.Request) {
		log.Println("api request has come" + req.URL.Path)
		w.WriteHeader(200)
		io.WriteString(w, `{"message":"Hello World!"}`)
	})
	childMux.HandleFunc("/guardmech/api/users", func(w http.ResponseWriter, req *http.Request) {

	})
	childMux.Handle("/guardmech/admin/", http.FileServer(http.Dir("dist")))

	//authMux := auth.NewMux()
	authMux := handler.NewAuthMux()
	adminMux := handler.NewAdminMux()

	mux := http.NewServeMux()
	mux.Handle("/auth/", authMux.Mux())
	mux.HandleFunc("/guardmech/", func(w http.ResponseWriter, req *http.Request) {
		log.Println("guardmech request")
		childMux.ServeHTTP(w, req)
	})
	mux.Handle("/guardmech/api/", adminMux.Mux())
	mux.HandleFunc("/", func(w http.ResponseWriter, req *http.Request) {
		log.Println("catch all:", req.URL.Path)
	})

	return http.Serve(listener, mux)
}

func (g *GuardMech) ReverseProxy(w http.ResponseWriter, req *http.Request) {
}

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
