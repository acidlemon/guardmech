package membership

import (
	"context"
	"crypto/sha512"
	"database/sql"
	"encoding/hex"
	"math/rand"
	"strings"

	"github.com/acidlemon/guardmech/db"
	"github.com/google/uuid"
)

type Service struct {
	repos Repository
}

func NewService(repos Repository) *Service {
	return &Service{
		repos: repos,
	}
}

func (s *Service) CreatePrincipal(ctx context.Context, tx *db.Tx, name, description string) (*Principal, error) {
	pri := &Principal{
		UniqueID:    uuid.New(),
		Name:        name,
		Description: description,
	}

	seqID, err := s.repos.SavePrincipal(ctx, tx, pri)
	if err != nil {
		return nil, err
	}
	pri.SeqID = seqID

	return pri, nil
}

func (s *Service) CreateAuth(ctx context.Context, tx *db.Tx, owner *Principal, idToken *OpenIDToken) (*Auth, error) {
	a := &Auth{
		UniqueID:  uuid.New(),
		Subject:   idToken.Sub,
		Issuer:    idToken.Issuer,
		Email:     idToken.Email,
		Principal: owner,
	}

	seqID, err := s.repos.SaveAuth(ctx, tx, a)
	if err != nil {
		return nil, err
	}
	a.SeqID = seqID

	return a, nil
}

var seedLetters = []byte("0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ_:")

// randomized letter string
func generateToken(length int) string {
	randbytes := make([]byte, ((length + 3) * 3 / 4))
	rand.Read(randbytes)

	b := strings.Builder{}
	limit := (length + 3) / 4
	for loop := 0; loop < limit; loop++ {
		idx := 3 * loop
		p0, p1, p2, p3 := 0x3F&randbytes[idx], (0xC0&randbytes[idx])>>2|(0x0F&randbytes[idx+1]), (0xF0&randbytes[idx+1])>>2|(0x03&randbytes[idx+2]), 0xFC&randbytes[idx+2]>>2

		b.WriteByte(seedLetters[p0])
		b.WriteByte(seedLetters[p1])
		b.WriteByte(seedLetters[p2])
		b.WriteByte(seedLetters[p3])
	}

	return b.String()[:length]
}

func hashToken(token, salt []byte, stretching int) []byte {
	input := make([]byte, 0, len(token)+len(salt))
	input = append(input, token...)
	input = append(input, salt...)
	inter := sha512.Sum512(input)
	for i := 0; i < stretching; i++ {
		inter = sha512.Sum512(inter[:])
	}

	return inter[:]
}

const stretchingCount = 4096 // like $6$12

func (s *Service) CreateAPIKey(ctx context.Context, tx *db.Tx, owner *Principal, name, description string) (*APIKey, string, error) {
	// generate token

	token := generateToken(36)
	salt := generateToken(28)

	maskedToken := token[:4] + strings.Repeat("*", len(token)-8) + token[len(token)-4:]
	hashedToken := hashToken([]byte(token), []byte(salt), stretchingCount)

	a := &APIKey{
		UniqueID:    uuid.New(),
		MaskedToken: maskedToken,
		HashedToken: hex.EncodeToString(hashedToken),
		Salt:        salt,
		Principal:   owner,
	}

	seqID, err := s.repos.SaveAPIKey(ctx, tx, a)
	if err != nil {
		return nil, "", err
	}
	a.SeqID = seqID

	return a, token, nil // return raw token once
}

func (s *Service) CreateFirstPrincipal(ctx context.Context, conn *sql.Conn, idToken *OpenIDToken) (*Principal, error) {
	tx, err := db.Begin(ctx, conn)
	defer tx.AutoRollback()

	// create principal
	name := idToken.Name
	if name == "" {
		name = idToken.Email
	}
	pri, err := s.CreatePrincipal(ctx, tx, name, "")
	if err != nil {
		return nil, err
	}

	// create auth
	_, err = s.CreateAuth(ctx, tx, pri, idToken)
	if err != nil {
		return nil, err
	}

	// create permission
	pe, err := s.CreatePermission(ctx, tx, PermissionOwner, "")
	if err != nil {
		return nil, err
	}

	//
	r, err := s.CreateRole(ctx, tx, RoleOwner, "")
	if err != nil {
		return nil, err
	}

	r.AttachPermission(ctx, tx, pe)
	pri.AttachRole(ctx, tx, r)

	tx.Commit()

	return pri, nil

}
func (s *Service) FindPrincipal(ctx context.Context, conn *sql.Conn, issuer, subject string) (*Principal, error) {
	return s.repos.FindPrincipal(ctx, conn, issuer, subject)
}

func (s *Service) FindPrincipalBySeqID(ctx context.Context, conn *sql.Conn, principalID int64) (*Principal, error) {
	return s.repos.FindPrincipalBySeqID(ctx, conn, principalID)
}

func (s *Service) FetchAllPrincipal(ctx context.Context, conn *sql.Conn) ([]*Principal, error) {
	return s.repos.FetchAllPrincipal(ctx, conn)
}
func (s *Service) FetchAllRole(ctx context.Context, conn *sql.Conn) ([]*Role, error) {
	return s.repos.FetchAllRole(ctx, conn)
}

func (s *Service) CreatePermission(ctx context.Context, tx *db.Tx, name, description string) (*Permission, error) {
	pe := &Permission{
		UniqueID:    uuid.New(),
		Name:        name,
		Description: description,
	}

	seqID, err := s.repos.SavePermission(ctx, tx, pe)
	if err != nil {
		return nil, err
	}
	pe.SeqID = seqID

	return pe, nil
}

func (s *Service) CreateRole(ctx context.Context, tx *db.Tx, name, description string) (*Role, error) {
	r := &Role{
		UniqueID:    uuid.New(),
		Name:        name,
		Description: description,
	}

	seqID, err := s.repos.SaveRole(ctx, tx, r)
	if err != nil {
		return nil, err
	}
	r.SeqID = seqID

	return r, nil
}
