package membership

import (
	"context"
	"fmt"
	"log"
	"strings"

	"github.com/acidlemon/guardmech/app/logic"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type Context = context.Context

type Principal struct {
	PrincipalID uuid.UUID
	Name        string
	Description string

	auth    *OIDCAuthorization
	apiKeys []*AuthAPIKey
	roles   []*Role
	groups  []*Group
}

func (p *Principal) AttachedRoles() []*Role {
	if p.roles == nil {
		return []*Role{}
	}

	return p.roles
}

func (p *Principal) Roles() []*Role {
	// direct roles
	tmp := []*Role{}
	if len(p.roles) != 0 {
		tmp = append(tmp, p.roles...)
	}

	if len(p.groups) == 0 {
		return tmp
	}

	// indirect roles
	for _, g := range p.groups {
		tmp = append(tmp, g.Roles()...)
	}

	// uniq list
	existsMap := map[string]bool{}
	result := make([]*Role, 0, len(tmp))
	for _, r := range tmp {
		uidstr := r.RoleID.String()
		if _, exist := existsMap[uidstr]; !exist {
			result = append(result, r)
			existsMap[uidstr] = true
		}
	}

	return result
}

func (p *Principal) Groups() []*Group {
	if p.groups == nil {
		return []*Group{}
	}

	return p.groups
}

func (p *Principal) HavingPermissions() []*Permission {
	roles := p.Roles()

	tmp := []*Permission{}
	if len(roles) == 0 {
		return tmp
	}

	for _, r := range roles {
		tmp = append(tmp, r.Permissions()...)
	}

	// uniq list
	existsMap := map[string]bool{}
	result := make([]*Permission, 0, len(tmp))
	for _, p := range tmp {
		uidstr := p.PermissionID.String()
		if _, exist := existsMap[uidstr]; !exist {
			result = append(result, p)
			existsMap[uidstr] = true
		}
	}

	return result
}

// return OpenID Connect Authorization info. May be nil.
func (p *Principal) OIDCAuthorization() *OIDCAuthorization {
	return p.auth
}

func (p *Principal) APIKeys() []*AuthAPIKey {
	if len(p.apiKeys) == 0 {
		return []*AuthAPIKey{}
	}

	return p.apiKeys
}

//-- write

// Add New APIKey
func (p *Principal) CreateAPIKey(name string) (*AuthAPIKey, string, error) {
	if name == "" {
		return nil, "", fmt.Errorf("CreateAPIKey: name is required")
	}

	newToken := "gmch-" + logic.GenerateRandomString(43)

	hashed, err := bcrypt.GenerateFromPassword([]byte(newToken), 12)
	if err != nil {
		log.Println("failed to run bcrypt. Maybe this is bug:", err)
		return nil, "", err
	}

	masked := strings.Repeat("*", 20) + newToken[44:]

	key := &AuthAPIKey{
		AuthAPIKeyID: uuid.New(),
		Name:         name,
		HashedToken:  string(hashed),
		MaskedToken:  masked,
	}

	return key, newToken, nil
}

func (p *Principal) AttachNewGroup(name, description string) (*Group, error) {
	if name == "" {
		return nil, fmt.Errorf("AttachNewGroup: name is required")
	}

	// create New Group
	g := newGroup(name, description)
	err := p.AttachGroup(g)
	if err != nil {
		return nil, err
	}
	return g, nil
}

func (p *Principal) AttachGroup(g *Group) error {
	for _, v := range p.groups {
		if v.GroupID == g.GroupID {
			// already exists
			return nil // TODO error?
		}
	}

	p.groups = append(p.groups, g)
	return nil
}

func (p *Principal) DetachGroup(g *Group) error {
	for i, v := range p.groups {
		if v.GroupID == g.GroupID {
			p.groups = append(p.groups[:i], p.groups[i+1:]...)
			return nil
		}
	}

	return nil
}

func (p *Principal) AttachNewRole(name, description string) (*Role, error) {
	if name == "" {
		return nil, fmt.Errorf("AttachNewRole: name is required")
	}

	// create New Role
	r := newRole(name, description)
	err := p.AttachRole(r)
	if err != nil {
		return nil, err
	}
	return r, nil
}
func (p *Principal) AttachRole(r *Role) error {
	for _, v := range p.roles {
		if v.RoleID == r.RoleID {
			// already exists
			return nil // TODO error?
		}
	}

	p.roles = append(p.roles, r)
	return nil
}

func (p *Principal) DetachRole(r *Role) error {
	for i, v := range p.roles {
		if v.RoleID == r.RoleID {
			p.roles = append(p.roles[:i], p.roles[i+1:]...)
			return nil
		}
	}
	return nil // Not Found, but no error (idempotence)
}
