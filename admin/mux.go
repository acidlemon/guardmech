package admin

import (
	"context"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"strconv"

	"github.com/acidlemon/guardmech/infra"
	"github.com/acidlemon/guardmech/membership"
	"github.com/gorilla/mux"
)

type Context = context.Context

type Mux struct {
	usecase *Usecase
}

func NewMux() *Mux {
	repos := &infra.Membership{}
	am := &Mux{
		usecase: &Usecase{
			repos: repos,
			svc:   membership.NewService(repos),
		},
	}

	return am
}

func (a *Mux) Mux() http.Handler {
	r := mux.NewRouter()
	r.HandleFunc("/guardmech/api/", a.ApiFallbackHandler)
	r.HandleFunc("/guardmech/api/principals", a.ListPrincipalsHandler)
	r.HandleFunc("/guardmech/api/principal", a.CreatePrincipalHandler)
	r.HandleFunc("/guardmech/api/principal/{id:[0-9]+}", a.PrincipalGetHandler).Methods(http.MethodGet)
	r.HandleFunc("/guardmech/api/principal/{id:[0-9]+}", a.PrincipalPostHandler).Methods(http.MethodPost)
	r.HandleFunc("/guardmech/api/principal/{id:[0-9]+}/new_key", a.CreateAPIKeyHandler).Methods(http.MethodPost)
	r.HandleFunc("/guardmech/api/roles", a.ListRolesHandler)
	//	r.HandleFunc("/guardmech/api/role/{id:[0-9]+}", RoleHandler)
	//	r.HandleFunc("/guardmech/api/permissions", ListPermissionHandler)
	//	r.HandleFunc("/guardmech/api/permission/{id:[0-9]+}", PermissionHandler)

	return r
}

func (a *Mux) ApiFallbackHandler(w http.ResponseWriter, req *http.Request) {
	log.Println("api request has come" + req.URL.Path)
	w.WriteHeader(200)
	io.WriteString(w, `{"message":"Hello World!"}`)
}

func (a *Mux) ListPrincipalsHandler(w http.ResponseWriter, req *http.Request) {
	// TODO permission check
	list, err := a.usecase.ListPrincipals(req.Context())
	if err != nil {
		errorJSON(w, err)
		return
	}

	renderJSON(w, map[string]interface{}{
		"principals": list,
	})
}

func (a *Mux) PrincipalGetHandler(w http.ResponseWriter, req *http.Request) {
	vars := mux.Vars(req)
	idStr := vars["id"]
	var id int64
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		errorJSON(w, err)
		return
	}

	payload, err := a.usecase.ShowPrincipal(req.Context(), id)
	if err != nil {
		errorJSON(w, err)
		return
	}

	renderJSON(w, payload)
}

func (a *Mux) PrincipalPostHandler(w http.ResponseWriter, req *http.Request) {
	vars := mux.Vars(req)
	idStr := vars["id"]
	var id int64
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		errorJSON(w, err)
		return
	}

	// create
	//name := vars["name"]
	//description := vars["description"]
	log.Println("POST id=", id)
}

func (a *Mux) CreatePrincipalHandler(w http.ResponseWriter, req *http.Request) {
	err := req.ParseForm()
	if err != nil {
		errorJSON(w, err)
		return
	}
	name := req.Form.Get("name")
	description := req.Form.Get("description")

	pri, err := a.usecase.CreatePrincipal(req.Context(), name, description)
	if err != nil {
		errorJSON(w, err)
		return
	}

	renderJSON(w, pri)
}

func (a *Mux) CreateAPIKeyHandler(w http.ResponseWriter, req *http.Request) {
	vars := mux.Vars(req)
	idStr := vars["id"]
	var id int64
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		errorJSON(w, err)
		return
	}

	// parameters
	err = req.ParseForm()
	if err != nil {
		errorJSON(w, err)
		return
	}
	name := req.Form.Get("name")
	description := req.Form.Get("description")

	ap, rawToken, err := a.usecase.CreateAPIKey(req.Context(), id, name, description)
	if err != nil {
		errorJSON(w, err)
		return
	}

	result := map[string]interface{}{
		"token":   rawToken,
		"api_key": ap,
	}

	renderJSON(w, result)
}

func (a *Mux) ListRolesHandler(w http.ResponseWriter, req *http.Request) {
	// TODO permission check

	list, err := a.usecase.ListRoles(req.Context())
	if err != nil {
		errorJSON(w, err)
		return
	}

	renderJSON(w, map[string]interface{}{
		"roles": list,
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

	log.Println("500 Internal Server Error: ", string(b))

	w.WriteHeader(500)
	w.Write(b)
}
