package membership

import (
	"context"
	"database/sql"
	"log"

	"github.com/acidlemon/guardmech/db"
	"github.com/google/uuid"
	"github.com/k0kubun/pp"
)

type Auth struct {
	SeqID     int64
	UniqueID  uuid.UUID  `json:"uuid"`
	Issuer    string     `json:"issuer"`
	Subject   string     `json:"subject"`
	Email     string     `json:"email"`
	Principal *Principal `json:"-"`
}

type APIKey struct {
	SeqID     int64
	UniqueID  uuid.UUID  `json:"uuid"`
	Token     string     `json:"token"`
	Principal *Principal `json:"-"`
}

type Principal struct {
	SeqID       int64
	UniqueID    uuid.UUID `json:"uuid"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
}

type PrincipalPayload struct {
	Principal   *Principal    `json:"principal"`
	Auths       []*Auth       `json:"auths"`
	APIKeys     []*APIKey     `json:"api_keys"`
	Groups      []*Group      `json:"groups"`
	Roles       []*Role       `json:"roles"`
	Permissions []*Permission `json:"permissions"`
}

func (pr *Principal) AttachRole(ctx context.Context, tx *db.Tx, r *Role) error {
	_, err := tx.ExecContext(ctx, `INSERT INTO principal_role_map (principal_id, role_id) VALUES (?, ?)`, pr.SeqID, r.SeqID)
	return err
}

func (pr *Principal) FindRole(ctx context.Context, conn *sql.Conn) ([]*Role, error) {
	result := make([]*Role, 0, 32)

	pp.Print(pr)

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

//func (pr *Principal)
