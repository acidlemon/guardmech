package handler

import (
	"context"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"strconv"

	"github.com/acidlemon/guardmech/app/usecase"
	"github.com/gorilla/mux"
)

type Context = context.Context

type AdminMux struct {
	u *usecase.Administration
}

func NewAdminMux() *AdminMux {
	u := usecase.NewAdministration()

	am := &AdminMux{
		u: u,
	}

	return am
}

func (a *AdminMux) Mux() http.Handler {
	r := mux.NewRouter()
	r.HandleFunc("/guardmech/api/", a.ApiFallbackHandler)
	// r.HandleFunc("/guardmech/api/principals", a.ListPrincipalsHandler)
	r.HandleFunc("/guardmech/api/principal", a.CreatePrincipalHandler)
	// r.HandleFunc("/guardmech/api/principal/{id:[0-9]+}", a.GetPrincipalHandler).Methods(http.MethodGet)
	//	r.HandleFunc("/guardmech/api/principal/{id:[0-9]+}", a.UpdatePrincipalHandler).Methods(http.MethodPost)
	r.HandleFunc("/guardmech/api/principal/{id:[0-9]+}/new_key", a.CreateAPIKeyHandler).Methods(http.MethodPost)

	// r.HandleFunc("/guardmech/api/roles", a.ListRolesHandler)
	r.HandleFunc("/guardmech/api/role/new", a.CreateRoleHandler).Methods(http.MethodPost)
	r.HandleFunc("/guardmech/api/role/{id:[0-9]+}", a.GetRoleHandler).Methods(http.MethodGet)
	r.HandleFunc("/guardmech/api/role/{id:[0-9]+}", a.UpdateRoleHandler).Methods(http.MethodPost)

	// r.HandleFunc("/guardmech/api/rules", a.ListRulesHandler)
	// r.HandleFunc("/guardmech/api/rule/new", a.CreateRuleHandler).Methods(http.MethodPost)
	// r.HandleFunc("/guardmech/api/rule/{id:[0-9]+}", a.GetRuleHandler).Methods(http.MethodGet)
	// r.HandleFunc("/guardmech/api/rule/{id:[0-9]+}", a.UpdateRuleHandler).Methods(http.MethodPost)

	r.HandleFunc("/guardmech/api/groups", a.ListGroupsHandler)
	r.HandleFunc("/guardmech/api/group/new", a.CreateGroupHandler).Methods(http.MethodPost)
	r.HandleFunc("/guardmech/api/group/{id:[0-9]+}", a.GetGroupHandler).Methods(http.MethodGet)
	r.HandleFunc("/guardmech/api/group/{id:[0-9]+}", a.UpdateGroupHandler).Methods(http.MethodPost)

	r.HandleFunc("/guardmech/api/permissions", a.ListPermissionsHandler)
	r.HandleFunc("/guardmech/api/permission/new", a.CreatePermissionHandler).Methods(http.MethodPost)
	r.HandleFunc("/guardmech/api/permission/{id:[0-9]+}", a.PermissionGetHandler).Methods(http.MethodGet)
	r.HandleFunc("/guardmech/api/permission/{id:[0-9]+}", a.PermissionPostHandler).Methods(http.MethodPost)

	return r
}

func (a *AdminMux) ApiFallbackHandler(w http.ResponseWriter, req *http.Request) {
	log.Println("api request has come" + req.URL.Path)
	w.WriteHeader(200)
	io.WriteString(w, `{"message":"Hello World!"}`)
}

/*
func (a *AdminMux) ListPrincipalsHandler(w http.ResponseWriter, req *http.Request) {
	// TODO permission check
	list, err := a.u.ListPrincipals(req.Context())
	if err != nil {
		errorJSON(w, err)
		return
	}

	renderJSON(w, map[string]interface{}{
		"principals": list,
	})
}

func (a *AdminMux) GetPrincipalHandler(w http.ResponseWriter, req *http.Request) {
	// TODO permission check

	vars := mux.Vars(req)
	idStr := vars["id"]
	var id int64
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		errorJSON(w, err)
		return
	}

	payload, err := a.u.ShowPrincipal(req.Context(), id)
	if err != nil {
		errorJSON(w, err)
		return
	}

	renderJSON(w, payload)
}
*/

// func (a *Mux) UpdatePrincipalHandler(w http.ResponseWriter, req *http.Request) {
// 	// TODO permission check

// 	vars := mux.Vars(req)
// 	idStr := vars["id"]
// 	var id int64
// 	id, err := strconv.ParseInt(idStr, 10, 64)
// 	if err != nil {
// 		errorJSON(w, err)
// 		return
// 	}

// 	name := vars["name"]
// 	description := vars["description"]
// 	log.Println("POST id=", id)

// 	pri, err := a.u.UpdatePrincipal(req.Context(), name, description)
// 	if err != nil {
// 		errorJSON(w, err)
// 		return
// 	}

// 	renderJSON(w, pri)
// }

func (a *AdminMux) CreatePrincipalHandler(w http.ResponseWriter, req *http.Request) {
	// TODO permission check

	err := req.ParseForm()
	if err != nil {
		errorJSON(w, err)
		return
	}
	name := req.Form.Get("name")
	description := req.Form.Get("description")

	pri, err := a.u.CreatePrincipal(req.Context(), name, description)
	if err != nil {
		errorJSON(w, err)
		return
	}

	renderJSON(w, pri)
}

func (a *AdminMux) CreateAPIKeyHandler(w http.ResponseWriter, req *http.Request) {
	// TODO permission check

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

	ap, rawToken, err := a.u.CreateAPIKey(req.Context(), id, name)
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

/*
func (a *AdminMux) ListRolesHandler(w http.ResponseWriter, req *http.Request) {
	// TODO permission check

	list, err := a.u.ListRoles(req.Context())
	if err != nil {
		errorJSON(w, err)
		return
	}

	renderJSON(w, map[string]interface{}{
		"roles": list,
	})
}
*/
func (a *AdminMux) CreateRoleHandler(w http.ResponseWriter, req *http.Request) {
	// TODO permission check

	role, err := a.u.CreateRole(req.Context())
	if err != nil {
		errorJSON(w, err)
		return
	}

	renderJSON(w, map[string]interface{}{
		"role": role,
	})
}

func (a *AdminMux) GetRoleHandler(w http.ResponseWriter, req *http.Request) {
	// TODO permission check

	vars := mux.Vars(req)
	idStr := vars["id"]
	var id int64
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		errorJSON(w, err)
		return
	}

	role, err := a.u.FetchRole(req.Context(), id)
	if err != nil {
		errorJSON(w, err)
		return
	}

	renderJSON(w, map[string]interface{}{
		"role": role,
	})
}

func (a *AdminMux) UpdateRoleHandler(w http.ResponseWriter, req *http.Request) {
	// TODO permission check

	vars := mux.Vars(req)
	idStr := vars["id"]
	var id int64
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		errorJSON(w, err)
		return
	}

	list, err := a.u.UpdateRole(req.Context(), id)
	if err != nil {
		errorJSON(w, err)
		return
	}

	renderJSON(w, map[string]interface{}{
		"roles": list,
	})
}

func (a *AdminMux) ListMappingRulesHandler(w http.ResponseWriter, req *http.Request) {
	// TODO permission check

	list, err := a.u.ListMappingRules(req.Context())
	if err != nil {
		errorJSON(w, err)
		return
	}

	renderJSON(w, map[string]interface{}{
		"mapping_rules": list,
	})
}

func (a *AdminMux) CreateMappingRuleHandler(w http.ResponseWriter, req *http.Request) {
	// TODO permission check

	rule, err := a.u.CreateMappingRule(req.Context())
	if err != nil {
		errorJSON(w, err)
		return
	}

	renderJSON(w, map[string]interface{}{
		"mapping_rule": rule,
	})
}

func (a *AdminMux) GetMappingRuleHandler(w http.ResponseWriter, req *http.Request) {
	// TODO permission check

	vars := mux.Vars(req)
	idStr := vars["id"]
	var id int64
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		errorJSON(w, err)
		return
	}

	rule, err := a.u.FetchMappingRule(req.Context(), id)
	if err != nil {
		errorJSON(w, err)
		return
	}

	renderJSON(w, map[string]interface{}{
		"mapping_rule": rule,
	})
}

func (a *AdminMux) UpdateMappingRuleHandler(w http.ResponseWriter, req *http.Request) {
	// TODO permission check

	vars := mux.Vars(req)
	idStr := vars["id"]
	var id int64
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		errorJSON(w, err)
		return
	}

	rule, err := a.u.UpdateMappingRule(req.Context(), id)
	if err != nil {
		errorJSON(w, err)
		return
	}

	renderJSON(w, map[string]interface{}{
		"mapping_rule": rule,
	})
}

func (a *AdminMux) ListGroupsHandler(w http.ResponseWriter, req *http.Request) {
	// TODO permission check

	list, err := a.u.ListGroups(req.Context())
	if err != nil {
		errorJSON(w, err)
		return
	}

	renderJSON(w, map[string]interface{}{
		"roles": list,
	})
}

func (a *AdminMux) CreateGroupHandler(w http.ResponseWriter, req *http.Request) {
	// TODO permission check

	list, err := a.u.CreateGroup(req.Context())
	if err != nil {
		errorJSON(w, err)
		return
	}

	renderJSON(w, map[string]interface{}{
		"roles": list,
	})
}

func (a *AdminMux) GetGroupHandler(w http.ResponseWriter, req *http.Request) {
	// TODO permission check

	list, err := a.u.FetchGroup(req.Context())
	if err != nil {
		errorJSON(w, err)
		return
	}

	renderJSON(w, map[string]interface{}{
		"roles": list,
	})
}

func (a *AdminMux) UpdateGroupHandler(w http.ResponseWriter, req *http.Request) {
	// TODO permission check

	list, err := a.u.UpdateGroup(req.Context())
	if err != nil {
		errorJSON(w, err)
		return
	}

	renderJSON(w, map[string]interface{}{
		"roles": list,
	})
}

func (a *AdminMux) ListPermissionsHandler(w http.ResponseWriter, req *http.Request) {
	// TODO permission check

	list, err := a.u.ListPermissions(req.Context())
	if err != nil {
		errorJSON(w, err)
		return
	}

	renderJSON(w, map[string]interface{}{
		"roles": list,
	})
}

func (a *AdminMux) CreatePermissionHandler(w http.ResponseWriter, req *http.Request) {
	// TODO permission check

	list, err := a.u.CreatePermission(req.Context())
	if err != nil {
		errorJSON(w, err)
		return
	}

	renderJSON(w, map[string]interface{}{
		"permission": list,
	})
}

func (a *AdminMux) PermissionGetHandler(w http.ResponseWriter, req *http.Request) {
	// TODO permission check

	list, err := a.u.FetchPermission(req.Context())
	if err != nil {
		errorJSON(w, err)
		return
	}

	renderJSON(w, map[string]interface{}{
		"roles": list,
	})
}

func (a *AdminMux) PermissionPostHandler(w http.ResponseWriter, req *http.Request) {
	// TODO permission check

	list, err := a.u.UpdatePermission(req.Context())
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
