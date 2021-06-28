package membership

import (
	"context"
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

	roles  []*Role
	groups []*Group
	auth   *OIDCAuthorization
	// apiKeys []*
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

func (p *Principal) Permissions() []*Permission {
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

// write

func (p *Principal) CreateAPIKey(name string) (*AuthAPIKey, string, error) {
	newToken := "gmch-" + logic.GenerateRandomString(43)

	hashed, err := bcrypt.GenerateFromPassword([]byte(newToken), 12)
	if err != nil {
		log.Println("failed to run bcrypt. Maybe this is bug:", err)
		return nil, "", err
	}

	masked := strings.Repeat("*", 44) + newToken[44:]

	log.Println(newToken)
	log.Println(masked)

	key := &AuthAPIKey{
		AuthAPIKeyID: uuid.New(),
		Name:         name,
		HashedToken:  string(hashed),
		MaskedToken:  masked,
	}

	return key, newToken, nil
}

func (p *Principal) AttachNewRole(name, description string) (*Role, error) {
	// create New Role
	r, err := newRole(name, description)
	if err != nil {
		return nil, err
	}
	err = p.AttachRole(r)
	if err != nil {
		return nil, err
	}
	return r, nil
}
func (p *Principal) AttachRole(r *Role) error {
	p.roles = append(p.roles, r)

	return nil
}

/*
func (pr *Principal) AttachRole(ctx Context, tx *db.Tx, r *Role) error {
	_, err := tx.ExecContext(ctx, `INSERT INTO principal_role_map (principal_id, role_id) VALUES (?, ?)`, pr.SeqID, r.SeqID)
	return err
}

func (pr *Principal) FindRole(ctx Context, conn *sql.Conn) ([]*Role, error) {
	result := make([]*Role, 0, 32)

	// find role (direct attached)
	rows, err := conn.QueryContext(ctx,
		`SELECT r.seq_id, r.unique_id, r.name, r.description FROM principal_role_map AS m JOIN role_info AS r ON m.role_id = r.seq_id WHERE m.principal_id = ?`,
		pr.SeqID,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		var id int64
		var name, description string
		if err := rows.Scan(&id, &name, &description); err != nil {
			log.Println("scan error:", err)
			return nil, err
		}
		result = append(result, &Role{
			SeqID:       id,
			Name:        name,
			Description: description,
		})
	}

	// find role (attached via group)
	// TODO

	return result, nil
}
*/

//func (pr *Principal)
