package guardmech

import (
	"bytes"
	"context"
	"database/sql"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"net/http/httputil"
	"net/url"
	"strings"
)

type GuardMech struct {
	rp *httputil.ReverseProxy
}

func New() *GuardMech {
	u, err := url.Parse("http://oauth2_proxy:4180")
	if err != nil {
		log.Println(err)
		return nil
	}
	rp := httputil.NewSingleHostReverseProxy(u)
	gm := &GuardMech{
		rp: rp,
	}

	modifier := func(res *http.Response) error {
		if res.StatusCode == http.StatusAccepted {
			// pass authentication
			gm.Authenticate(res.Header.Get("X-Auth-Request-Email"), res)
		}
		return nil
	}
	rp.ModifyResponse = modifier

	return gm
}

func (g *GuardMech) Run() error {
	listener, err := net.Listen("tcp", "0.0.0.0:2989")
	if err != nil {
		log.Println("failed to listen")
		return nil
	}

	mux := http.NewServeMux()
	mux.HandleFunc("/oauth2/", func(w http.ResponseWriter, req *http.Request) {
		log.Println("request has come " + req.URL.Path)
		g.ReverseProxy(w, req)
	})
	mux.HandleFunc("/guardmech/", func(w http.ResponseWriter, req *http.Request) {
		log.Println("request has come " + req.URL.Path)
		w.WriteHeader(200)
		io.WriteString(w, "Hello World!")
	})

	return http.Serve(listener, mux)
}

func (g *GuardMech) ReverseProxy(w http.ResponseWriter, req *http.Request) {
	g.rp.ServeHTTP(w, req)
}

func (g *GuardMech) Authenticate(account string, res *http.Response) {
	ctx := res.Request.Context()
	conn, err := db.Conn(ctx)
	if err != nil {
		WrapServerError(res, err)
		return
	}
	ok, err := HasPrincipal(ctx, conn)
	if err != nil {
		WrapServerError(res, err)
		return
	}

	// first cut
	if !ok {
		log.Println("No User!! Entering setup mode.")
		err = g.Setup(ctx, conn, account)
		if err != nil {
			WrapServerError(res, err)
			return
		}
	}

	log.Println(account)

	// output headers
	pr, err := FindPrincipal(ctx, conn, account)
	if err != nil {
		WrapServerError(res, err)
		return
	}
	roles, err := pr.FindRole(ctx, conn)
	if err != nil {
		WrapServerError(res, err)
		return
	}

	names := make([]string, 0, len(roles))
	for _, r := range roles {
		names = append(names, r.Name)
	}
	res.Header.Add("X-Auth-Role", strings.Join(names, " "))

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

func (g *GuardMech) Setup(ctx context.Context, conn *sql.Conn, account string) error {
	// create owner user
	err := CreateFirstPrincipal(ctx, conn, account)
	if err != nil {
		log.Println("failed to setup:", err.Error())
		return err
	}

	return nil
}
