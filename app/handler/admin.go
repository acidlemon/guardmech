package handler

import (
	"context"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"strconv"

	"github.com/acidlemon/guardmech/app/logic/membership"
	"github.com/acidlemon/guardmech/app/usecase"
	"github.com/acidlemon/guardmech/app/usecase/payload"
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
	r.HandleFunc("/guardmech/api/principal", a.CreatePrincipalHandler).Methods(http.MethodPost)
	r.HandleFunc("/guardmech/api/principals", a.ListPrincipalsHandler).Methods(http.MethodGet)
	r.HandleFunc("/guardmech/api/principal/{id:[0-9a-f-]+}", a.GetPrincipalHandler).Methods(http.MethodGet)
	r.HandleFunc("/guardmech/api/principal/{id:[0-9a-f-]+}", a.UpdatePrincipalHandler).Methods(http.MethodPost)
	r.HandleFunc("/guardmech/api/principal/{id:[0-9a-f-]+}", a.DeletePrincipalHandler).Methods(http.MethodDelete)
	r.HandleFunc("/guardmech/api/principal/{id:[0-9a-f-]+}/new_key", a.CreateAPIKeyHandler).Methods(http.MethodPost)
	r.HandleFunc("/guardmech/api/principal/{id:[0-9a-f-]+}/attach_role", a.AttachRoleToPrincipalHandler).Methods(http.MethodPost)
	r.HandleFunc("/guardmech/api/principal/{id:[0-9a-f-]+}/attach_group", a.AttachGroupToPrincipalHandler).Methods(http.MethodPost)
	r.HandleFunc("/guardmech/api/principal/{id:[0-9a-f-]+}/detach_role", a.DetachRoleToPrincipalHandler).Methods(http.MethodPost)
	r.HandleFunc("/guardmech/api/principal/{id:[0-9a-f-]+}/detach_group", a.DetachGroupToPrincipalHandler).Methods(http.MethodPost)

	r.HandleFunc("/guardmech/api/group", a.CreateGroupHandler).Methods(http.MethodPost)
	r.HandleFunc("/guardmech/api/groups", a.ListGroupsHandler)
	r.HandleFunc("/guardmech/api/group/{id:[0-9a-f-]+}", a.GetGroupHandler).Methods(http.MethodGet)
	r.HandleFunc("/guardmech/api/group/{id:[0-9a-f-]+}", a.UpdateGroupHandler).Methods(http.MethodPost)
	r.HandleFunc("/guardmech/api/group/{id:[0-9a-f-]+}", a.DeleteGroupHandler).Methods(http.MethodDelete)
	r.HandleFunc("/guardmech/api/group/{id:[0-9a-f-]+}/attach_role", a.AttachRoleToGroupHandler).Methods(http.MethodPost)
	r.HandleFunc("/guardmech/api/group/{id:[0-9a-f-]+}/detach_role", a.DetachRoleToGroupHandler).Methods(http.MethodPost)

	r.HandleFunc("/guardmech/api/role", a.CreateRoleHandler).Methods(http.MethodPost)
	r.HandleFunc("/guardmech/api/roles", a.ListRolesHandler)
	r.HandleFunc("/guardmech/api/role/{id:[0-9a-f-]+}", a.GetRoleHandler).Methods(http.MethodGet)
	r.HandleFunc("/guardmech/api/role/{id:[0-9a-f-]+}", a.UpdateRoleHandler).Methods(http.MethodPost)
	r.HandleFunc("/guardmech/api/role/{id:[0-9a-f-]+}", a.DeleteRoleHandler).Methods(http.MethodDelete)
	r.HandleFunc("/guardmech/api/role/{id:[0-9a-f-]+}/attach_permission", a.AttachPermissionToRoleHandler).Methods(http.MethodPost)
	r.HandleFunc("/guardmech/api/role/{id:[0-9a-f-]+}/detach_permission", a.DetachPermissionToRoleHandler).Methods(http.MethodPost)

	r.HandleFunc("/guardmech/api/permission", a.CreatePermissionHandler).Methods(http.MethodPost)
	r.HandleFunc("/guardmech/api/permissions", a.ListPermissionsHandler)
	r.HandleFunc("/guardmech/api/permission/{id:[0-9a-f-]+}", a.PermissionGetHandler).Methods(http.MethodGet)
	r.HandleFunc("/guardmech/api/permission/{id:[0-9a-f-]+}", a.UpdatePermissionHandler).Methods(http.MethodPost)
	r.HandleFunc("/guardmech/api/permission/{id:[0-9a-f-]+}", a.DeletePermissionHandler).Methods(http.MethodDelete)

	r.HandleFunc("/guardmech/api/mapping_rule", a.CreateMappingRuleHandler).Methods(http.MethodPost)
	r.HandleFunc("/guardmech/api/mapping_rules", a.ListMappingRulesHandler)
	r.HandleFunc("/guardmech/api/mapping_rule/{id:[0-9a-f-]+}", a.GetMappingRuleHandler).Methods(http.MethodGet)
	r.HandleFunc("/guardmech/api/mapping_rule/{id:[0-9a-f-]+}", a.UpdateMappingRuleHandler).Methods(http.MethodPost)

	r.Use(a.checkPermissionMiddleware)

	return r
}

func (a *AdminMux) ApiFallbackHandler(w http.ResponseWriter, req *http.Request) {
	log.Println("api request has come" + req.URL.Path)
	w.WriteHeader(200)
	io.WriteString(w, `{"message":"Hello World!"}`)
}

func (a *AdminMux) checkPermissionMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		c, err := req.Cookie(sessionKey)
		if err != nil {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		is := &usecase.IDSession{}
		_, err = RestoreSessionPayload(c.Value, is)
		if err != nil {
			httperr := NewHttpError(http.StatusUnauthorized, "failed to restore cookie", err)
			WriteHttpError(w, httperr)
			return
		}

		for _, perm := range is.Membership.Principal.Permissions {
			// TODO read only mode?
			if perm == membership.PermissionOwnerName {
				next.ServeHTTP(w, req)
				return
			}
		}

		http.Error(w, "Forbidden", http.StatusForbidden)
	})
}

// -- Principal

func (a *AdminMux) ListPrincipalsHandler(w http.ResponseWriter, req *http.Request) {
	list, err := a.u.ListPrincipals(req.Context())
	if err != nil {
		errorJSON(w, err)
		return
	}

	payloads := make([]*payload.PrincipalPayload, 0, len(list))
	for _, v := range list {
		payloads = append(payloads, payload.PrincipalPayloadFromEntity(v))
	}

	renderJSON(w, map[string]interface{}{
		"principals": payloads,
	})
}

func (a *AdminMux) GetPrincipalHandler(w http.ResponseWriter, req *http.Request) {
	vars := mux.Vars(req)
	id := vars["id"]
	pri, err := a.u.ShowPrincipal(req.Context(), id)
	if err != nil {
		errorJSON(w, err)
		return
	}

	p := payload.PrincipalPayloadFromEntity(pri)
	renderJSON(w, map[string]interface{}{
		"principal": p,
	})
}

func (a *AdminMux) UpdatePrincipalHandler(w http.ResponseWriter, req *http.Request) {
	// vars := mux.Vars(req)
	// id := vars["id"]

	// name := vars["name"]
	// description := vars["description"]

	// pri, err := a.u.UpdatePrincipal(req.Context(), name, description)
	// if err != nil {
	// 	errorJSON(w, err)
	// 	return
	// }

	// renderJSON(w, pri)

}

func (a *AdminMux) DeletePrincipalHandler(w http.ResponseWriter, req *http.Request) {
	vars := mux.Vars(req)
	id := vars["id"]

	err := a.u.DeletePrincipal(req.Context(), id)
	if err != nil {
		errorJSON(w, err)
		return
	}

	renderJSON(w, map[string]interface{}{
		"result": "ok",
	})
}

func (a *AdminMux) CreatePrincipalHandler(w http.ResponseWriter, req *http.Request) {
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

	p := payload.PrincipalPayloadFromEntity(pri)
	renderJSON(w, map[string]interface{}{
		"principal": p,
	})
}

func (a *AdminMux) CreateAPIKeyHandler(w http.ResponseWriter, req *http.Request) {
	vars := mux.Vars(req)
	id := vars["id"]
	// parameters
	err := req.ParseForm()
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

func (a *AdminMux) AttachGroupToPrincipalHandler(w http.ResponseWriter, req *http.Request) {
	vars := mux.Vars(req)
	principalID := vars["id"]

	err := req.ParseForm()
	if err != nil {
		errorJSON(w, err)
		return
	}
	groupID := req.Form.Get("group_id")

	pri, err := a.u.AttachGroupToPrincipal(req.Context(), principalID, groupID)
	if err != nil {
		errorJSON(w, err)
		return
	}

	p := payload.PrincipalPayloadFromEntity(pri)
	renderJSON(w, map[string]interface{}{
		"principal": p,
	})
}

func (a *AdminMux) AttachRoleToPrincipalHandler(w http.ResponseWriter, req *http.Request) {
	vars := mux.Vars(req)
	principalID := vars["id"]

	err := req.ParseForm()
	if err != nil {
		errorJSON(w, err)
		return
	}
	roleID := req.Form.Get("role_id")

	pri, err := a.u.AttachRoleToPrincipal(req.Context(), principalID, roleID)
	if err != nil {
		errorJSON(w, err)
		return
	}

	p := payload.PrincipalPayloadFromEntity(pri)
	renderJSON(w, map[string]interface{}{
		"principal": p,
	})
}

func (a *AdminMux) DetachGroupToPrincipalHandler(w http.ResponseWriter, req *http.Request) {
	vars := mux.Vars(req)
	principalID := vars["id"]

	err := req.ParseForm()
	if err != nil {
		errorJSON(w, err)
		return
	}
	groupID := req.Form.Get("group_id")

	pri, err := a.u.DetachGroupFromPrincipal(req.Context(), principalID, groupID)
	if err != nil {
		errorJSON(w, err)
		return
	}

	p := payload.PrincipalPayloadFromEntity(pri)
	renderJSON(w, map[string]interface{}{
		"principal": p,
	})
}

func (a *AdminMux) DetachRoleToPrincipalHandler(w http.ResponseWriter, req *http.Request) {
	vars := mux.Vars(req)
	principalID := vars["id"]

	err := req.ParseForm()
	if err != nil {
		errorJSON(w, err)
		return
	}
	roleID := req.Form.Get("role_id")

	pri, err := a.u.DetachRoleFromPrincipal(req.Context(), principalID, roleID)
	if err != nil {
		errorJSON(w, err)
		return
	}

	p := payload.PrincipalPayloadFromEntity(pri)
	renderJSON(w, map[string]interface{}{
		"principal": p,
	})
}

// -- Group

func (a *AdminMux) ListGroupsHandler(w http.ResponseWriter, req *http.Request) {
	list, err := a.u.ListGroups(req.Context())
	if err != nil {
		errorJSON(w, err)
		return
	}

	groups := make([]*payload.Group, 0, len(list))
	for _, v := range list {
		groups = append(groups, payload.GroupFromEntity(v))
	}

	renderJSON(w, map[string]interface{}{
		"groups": groups,
	})
}

func (a *AdminMux) CreateGroupHandler(w http.ResponseWriter, req *http.Request) {
	err := req.ParseForm()
	if err != nil {
		errorJSON(w, err)
		return
	}
	name := req.Form.Get("name")
	description := req.Form.Get("description")

	g, err := a.u.CreateGroup(req.Context(), name, description)
	if err != nil {
		errorJSON(w, err)
		return
	}

	group := payload.GroupFromEntity(g)
	renderJSON(w, map[string]interface{}{
		"group": group,
	})
}

func (a *AdminMux) GetGroupHandler(w http.ResponseWriter, req *http.Request) {
	vars := mux.Vars(req)
	id := vars["id"]

	g, err := a.u.FetchGroup(req.Context(), id)
	if err != nil {
		errorJSON(w, err)
		return
	}

	group := payload.GroupFromEntity(g)
	renderJSON(w, map[string]interface{}{
		"group": group,
	})
}

func (a *AdminMux) UpdateGroupHandler(w http.ResponseWriter, req *http.Request) {
	list, err := a.u.UpdateGroup(req.Context())
	if err != nil {
		errorJSON(w, err)
		return
	}

	renderJSON(w, map[string]interface{}{
		"groups": list,
	})
}

func (a *AdminMux) DeleteGroupHandler(w http.ResponseWriter, req *http.Request) {
	vars := mux.Vars(req)
	id := vars["id"]

	err := a.u.DeleteGroup(req.Context(), id)
	if err != nil {
		errorJSON(w, err)
		return
	}

	renderJSON(w, map[string]interface{}{
		"result": "ok",
	})
}

func (a *AdminMux) AttachRoleToGroupHandler(w http.ResponseWriter, req *http.Request) {
	vars := mux.Vars(req)
	groupID := vars["id"]

	err := req.ParseForm()
	if err != nil {
		errorJSON(w, err)
		return
	}
	roleID := req.Form.Get("role_id")

	g, err := a.u.AttachRoleToGroup(req.Context(), groupID, roleID)
	if err != nil {
		errorJSON(w, err)
		return
	}

	group := payload.GroupFromEntity(g)
	renderJSON(w, map[string]interface{}{
		"group": group,
	})
}

func (a *AdminMux) DetachRoleToGroupHandler(w http.ResponseWriter, req *http.Request) {
	vars := mux.Vars(req)
	groupID := vars["id"]

	err := req.ParseForm()
	if err != nil {
		errorJSON(w, err)
		return
	}
	roleID := req.Form.Get("role_id")

	g, err := a.u.DetachRoleFromGroup(req.Context(), groupID, roleID)
	if err != nil {
		errorJSON(w, err)
		return
	}

	group := payload.GroupFromEntity(g)
	renderJSON(w, map[string]interface{}{
		"group": group,
	})
}

// -- Role

func (a *AdminMux) ListRolesHandler(w http.ResponseWriter, req *http.Request) {
	list, err := a.u.ListRoles(req.Context())
	if err != nil {
		errorJSON(w, err)
		return
	}

	roles := make([]*payload.Role, 0, len(list))
	for _, v := range list {
		roles = append(roles, payload.RoleFromEntity(v))
	}

	renderJSON(w, map[string]interface{}{
		"roles": roles,
	})
}

func (a *AdminMux) CreateRoleHandler(w http.ResponseWriter, req *http.Request) {
	err := req.ParseForm()
	if err != nil {
		errorJSON(w, err)
		return
	}
	name := req.Form.Get("name")
	description := req.Form.Get("description")

	role, err := a.u.CreateRole(req.Context(), name, description)
	if err != nil {
		errorJSON(w, err)
		return
	}

	renderJSON(w, map[string]interface{}{
		"role": role,
	})
}

func (a *AdminMux) GetRoleHandler(w http.ResponseWriter, req *http.Request) {
	vars := mux.Vars(req)
	id := vars["id"]

	r, err := a.u.FetchRole(req.Context(), id)
	if err != nil {
		errorJSON(w, err)
		return
	}

	role := payload.RoleFromEntity(r)
	renderJSON(w, map[string]interface{}{
		"role": role,
	})
}

func (a *AdminMux) UpdateRoleHandler(w http.ResponseWriter, req *http.Request) {
	vars := mux.Vars(req)
	id := vars["id"]

	list, err := a.u.UpdateRole(req.Context(), id)
	if err != nil {
		errorJSON(w, err)
		return
	}

	renderJSON(w, map[string]interface{}{
		"roles": list,
	})
}

func (a *AdminMux) DeleteRoleHandler(w http.ResponseWriter, req *http.Request) {
	vars := mux.Vars(req)
	id := vars["id"]

	err := a.u.DeleteRole(req.Context(), id)
	if err != nil {
		errorJSON(w, err)
		return
	}

	renderJSON(w, map[string]interface{}{
		"result": "ok",
	})
}

func (a *AdminMux) AttachPermissionToRoleHandler(w http.ResponseWriter, req *http.Request) {
	vars := mux.Vars(req)
	roleID := vars["id"]

	err := req.ParseForm()
	if err != nil {
		errorJSON(w, err)
		return
	}
	permissionID := req.Form.Get("permission_id")

	r, err := a.u.AttachPermissionToRole(req.Context(), roleID, permissionID)
	if err != nil {
		errorJSON(w, err)
		return
	}

	role := payload.RoleFromEntity(r)
	renderJSON(w, map[string]interface{}{
		"role": role,
	})
}

func (a *AdminMux) DetachPermissionToRoleHandler(w http.ResponseWriter, req *http.Request) {
	vars := mux.Vars(req)
	roleID := vars["id"]

	err := req.ParseForm()
	if err != nil {
		errorJSON(w, err)
		return
	}
	permissionID := req.Form.Get("permission_id")

	r, err := a.u.DetachPermissionFromRole(req.Context(), roleID, permissionID)
	if err != nil {
		errorJSON(w, err)
		return
	}

	role := payload.RoleFromEntity(r)
	renderJSON(w, map[string]interface{}{
		"role": role,
	})
}

// -- Permission

func (a *AdminMux) ListPermissionsHandler(w http.ResponseWriter, req *http.Request) {
	list, err := a.u.ListPermissions(req.Context())
	if err != nil {
		errorJSON(w, err)
		return
	}

	permissions := make([]*payload.Permission, 0, len(list))
	for _, v := range list {
		permissions = append(permissions, payload.PermissionFromEntity(v))
	}

	renderJSON(w, map[string]interface{}{
		"permissions": permissions,
	})
}

func (a *AdminMux) CreatePermissionHandler(w http.ResponseWriter, req *http.Request) {
	err := req.ParseForm()
	if err != nil {
		errorJSON(w, err)
		return
	}
	name := req.Form.Get("name")
	description := req.Form.Get("description")

	perm, err := a.u.CreatePermission(req.Context(), name, description)
	if err != nil {
		errorJSON(w, err)
		return
	}

	permission := payload.PermissionFromEntity(perm)
	renderJSON(w, map[string]interface{}{
		"permission": permission,
	})
}

func (a *AdminMux) PermissionGetHandler(w http.ResponseWriter, req *http.Request) {
	vars := mux.Vars(req)
	id := vars["id"]

	perm, err := a.u.FetchPermission(req.Context(), id)
	if err != nil {
		errorJSON(w, err)
		return
	}

	permission := payload.PermissionFromEntity(perm)
	renderJSON(w, map[string]interface{}{
		"permission": permission,
	})
}

func (a *AdminMux) UpdatePermissionHandler(w http.ResponseWriter, req *http.Request) {
	list, err := a.u.UpdatePermission(req.Context())
	if err != nil {
		errorJSON(w, err)
		return
	}

	renderJSON(w, map[string]interface{}{
		"permissions": list,
	})
}

func (a *AdminMux) DeletePermissionHandler(w http.ResponseWriter, req *http.Request) {
	vars := mux.Vars(req)
	id := vars["id"]

	err := a.u.DeletePermission(req.Context(), id)
	if err != nil {
		errorJSON(w, err)
		return
	}

	renderJSON(w, map[string]interface{}{
		"result": "ok",
	})
}

// -- Mapping Rule

func (a *AdminMux) ListMappingRulesHandler(w http.ResponseWriter, req *http.Request) {
	list, err := a.u.ListMappingRules(req.Context())
	if err != nil {
		errorJSON(w, err)
		return
	}

	rules := make([]*payload.MappingRule, 0, len(list))
	for _, v := range list {
		rules = append(rules, payload.MappingRuleFromEntity(v))
	}

	renderJSON(w, map[string]interface{}{
		"mapping_rules": rules,
	})
}

func (a *AdminMux) CreateMappingRuleHandler(w http.ResponseWriter, req *http.Request) {
	err := req.ParseForm()
	if err != nil {
		errorJSON(w, err)
		return
	}
	name := req.Form.Get("name")
	description := req.Form.Get("description")
	ruleTypeStr := req.Form.Get("rule_type")
	detail := req.Form.Get("detail")
	associationType := req.Form.Get("association_type")
	associationID := req.Form.Get("association_id")
	ruleType, err := strconv.Atoi(ruleTypeStr)
	if err != nil {
		errorJSON(w, err)
		return
	}

	rule, err := a.u.CreateMappingRule(req.Context(), name, description, ruleType, detail, associationType, associationID)
	if err != nil {
		errorJSON(w, err)
		return
	}

	renderJSON(w, map[string]interface{}{
		"mapping_rule": rule,
	})
}

func (a *AdminMux) GetMappingRuleHandler(w http.ResponseWriter, req *http.Request) {
	vars := mux.Vars(req)
	id := vars["id"]

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
	vars := mux.Vars(req)
	id := vars["id"]

	rule, err := a.u.UpdateMappingRule(req.Context(), id)
	if err != nil {
		errorJSON(w, err)
		return
	}

	renderJSON(w, map[string]interface{}{
		"mapping_rule": rule,
	})
}

// -- other utilities

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
