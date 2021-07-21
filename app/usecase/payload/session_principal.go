package payload

import entity "github.com/acidlemon/guardmech/app/logic/membership"

// Small Payload for Session Cookie
type SessionPrincipal struct {
	Email       string   `json:"email"`
	Groups      []string `json:"groups,omitempty"`
	Roles       []string `json:"roles,omitempty"`
	Permissions []string `json:"permissions,omitempty"`
}

func SessionPrincipalFromEntity(pri *entity.Principal) *SessionPrincipal {
	gs := pri.Groups()
	rs := pri.Roles()
	ps := pri.HavingPermissions()

	groups := make([]string, 0, len(gs))
	roles := make([]string, 0, len(rs))
	perms := make([]string, 0, len(ps))

	for _, v := range gs {
		groups = append(groups, v.Name)
	}
	for _, v := range rs {
		roles = append(roles, v.Name)
	}
	for _, v := range ps {
		perms = append(perms, v.Name)
	}

	return &SessionPrincipal{
		Email:       pri.OIDCAuthorization().Email,
		Groups:      groups,
		Roles:       roles,
		Permissions: perms,
	}
}
