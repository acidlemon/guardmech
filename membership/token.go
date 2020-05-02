package membership

type OpenIDToken struct {
	Issuer   string `json:"iss"`
	Sub      string `json:"sub"`
	Email    string `json:"email"`
	Verified bool   `json:"email_verified"`
	Name     string `json:"name"`
}
