package guardmech

import (
	"encoding/json"
	"io"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

func ApiMux() http.Handler {
	r := mux.NewRouter()
	r.HandleFunc("/guardmech/api/", ApiFallbackHandler)
	r.HandleFunc("/guardmech/api/principals", ListPrincipalsHandler)
	r.HandleFunc("/guardmech/api/principal/{id:[0-9]+}", PrincipalHandler)
	r.HandleFunc("/guardmech/api/roles", ListRolesHandler)
//	r.HandleFunc("/guardmech/api/role/{id:[0-9]+}", RoleHandler)
//	r.HandleFunc("/guardmech/api/permissions", ListPermissionHandler)
//	r.HandleFunc("/guardmech/api/permission/{id:[0-9]+}", PermissionHandler)

	return r
}

func ApiFallbackHandler(w http.ResponseWriter, req *http.Request) {
	log.Println("api request has come" + req.URL.Path)
	w.WriteHeader(200)
	io.WriteString(w, `{"message":"Hello World!"}`)
}

func ListPrincipalsHandler(w http.ResponseWriter, req *http.Request) {
	// TODO permission check

	ctx := req.Context()
	conn, err := db.Conn(ctx)
	if err != nil {
		errorJSON(w, err)
		return
	}
	defer conn.Close()

	prs, err := FetchAllPrincipal(ctx, conn)
	if err != nil {
		errorJSON(w, err)
		return
	}

	renderJSON(w, map[string]interface{}{
		"principals": prs,
	})
}

func PrincipalHandler(w http.ResponseWriter, req *http.Request) {
	ctx := req.Context()
	conn, err := db.Conn(ctx)
	if err != nil {
		errorJSON(w, err)
		return
	}
	defer conn.Close()

	vars := mux.Vars(req)
	id := vars["id"]

	switch req.Method {
	case http.MethodPost:
		// create
		//name := vars["name"]
		//description := vars["description"]

		break

	default:
		// read
		payload, err := FindPrincipalByID(ctx, conn, id)
		if err != nil {
			errorJSON(w, err)
			return
		}

		renderJSON(w, payload)
		break
	}
}

func ListRolesHandler(w http.ResponseWriter, req *http.Request) {
	// TODO permission check

	ctx := req.Context()
	conn, err := db.Conn(ctx)
	if err != nil {
		errorJSON(w, err)
		return
	}
	defer conn.Close()

	roles, err := FetchAllRole(ctx, conn)
	if err != nil {
		errorJSON(w, err)
		return
	}

	renderJSON(w, map[string]interface{}{
		"roles": roles,
	})
}

func renderJSON(w http.ResponseWriter, data interface{}) {
	b, err := json.Marshal(data)
	code := 200
	if err != nil {
		b = []byte(`{"error":"failed to marshal json"}`)
		code = 500
	}

	w.WriteHeader(code)
	w.Write(b)
}

func errorJSON(w http.ResponseWriter, err error) {
	b, err := json.Marshal(map[string]string{
		"error": err.Error(),
	})
	if err != nil {
		b = []byte(`{"error":"error occurred, additionally failed to marshal json"}`)
	}

	w.WriteHeader(500)
	w.Write(b)
}
