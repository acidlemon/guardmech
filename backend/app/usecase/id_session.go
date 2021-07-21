package usecase

import (
	"fmt"
	"strings"
	"time"

	"encoding/json"

	"github.com/acidlemon/guardmech/backend/app/usecase/payload"
)

type IDSession struct {
	Issuer     string
	Subject    string
	Email      string
	Membership MembershipToken
}

type MembershipToken struct {
	NextCheck time.Time                 `json:"next_check"`
	Principal *payload.SessionPrincipal `json:"principal"`
}

const IDSeparator string = "('-'o)"

func (is *IDSession) String() string {
	b := strings.Builder{}
	b.WriteString(is.Issuer)
	b.WriteString(IDSeparator)
	b.WriteString(is.Subject)
	b.WriteString(IDSeparator)
	b.WriteString(is.Email)
	b.WriteString(IDSeparator)
	enc := json.NewEncoder(&b)
	enc.Encode(is.Membership)

	return b.String()
}

func (is *IDSession) Restore(data string) error {
	ss := strings.Split(data, IDSeparator)
	if len(ss) < 4 {
		return fmt.Errorf("not enough session data: len=%d", len(ss))
	}

	is.Issuer = ss[0]
	is.Subject = ss[1]
	is.Email = ss[2]
	err := json.Unmarshal([]byte(ss[3]), &is.Membership)
	if err != nil {
		return err
	}

	return nil
}
