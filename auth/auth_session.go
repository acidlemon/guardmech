package auth

import "strings"

type AuthSession struct {
	CSRFToken string
	Path      string
}

func (as *AuthSession) String() string {
	b := strings.Builder{}
	b.WriteString(as.CSRFToken)
	b.WriteString("|")
	b.WriteString(as.Path)

	return b.String()
}

func (as *AuthSession) Restore(data string) error {
	ss := strings.Split(data, "|")
	as.CSRFToken = ss[0]
	as.Path = ss[1]
	if as.Path == "" {
		as.Path = "/"
	}

	return nil
}
