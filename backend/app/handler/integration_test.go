package handler

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/acidlemon/guardmech/backend/app/config"
	"github.com/acidlemon/guardmech/backend/app/logic/auth"
	"github.com/acidlemon/guardmech/backend/app/logic/membership"
	"github.com/acidlemon/guardmech/backend/app/usecase"
	"github.com/acidlemon/guardmech/backend/app/usecase/payload"
	"github.com/acidlemon/guardmech/backend/db"
	"github.com/acidlemon/guardmech/backend/persistence"
	"github.com/gorilla/mux"
)

const (
	testIssuer  = "https://accounts.google.com"
	testSubject = "integration-test-subject"
	testEmail   = "owner@example.com"
)

// The db package builds its DSN in init(), so the GUARDMECH_DB_* env vars must be set
// before the test process starts (t.Setenv does not work for them).
func setupIntegrationDB(t *testing.T) {
	t.Helper()
	dbName := os.Getenv("GUARDMECH_DB_NAME")
	if dbName == "" {
		t.Skip("GUARDMECH_DB_NAME is not set; skipping MySQL integration test")
	}
	if !strings.HasSuffix(dbName, "_test") {
		t.Skipf("GUARDMECH_DB_NAME %q does not end with _test; refusing to drop tables on a non-test database", dbName)
	}
	t.Setenv("GUARDMECH_SIGNATURE_KEY", "integration-test-signature-key")
	t.Setenv("GUARDMECH_CRYPTO_KEY", "0123456789abcdef0123456789abcdef")

	ctx := t.Context()
	conn, err := db.GetConn(ctx)
	if err != nil {
		t.Fatalf("db.GetConn: %v", err)
	}
	defer func() {
		if err := conn.Close(); err != nil {
			t.Errorf("close conn: %v", err)
		}
	}()

	rows, err := conn.QueryContext(ctx, "SELECT table_name FROM information_schema.tables WHERE table_schema = DATABASE()")
	if err != nil {
		t.Fatalf("list tables: %v", err)
	}
	var tables []string
	for rows.Next() {
		var tbl string
		if err := rows.Scan(&tbl); err != nil {
			t.Fatalf("scan table name: %v", err)
		}
		tables = append(tables, tbl)
	}
	if err := rows.Err(); err != nil {
		t.Fatalf("iterate tables: %v", err)
	}
	if err := rows.Close(); err != nil {
		t.Fatalf("close rows: %v", err)
	}
	for _, tbl := range tables {
		if _, err := conn.ExecContext(ctx, "DROP TABLE IF EXISTS "+tbl); err != nil {
			t.Fatalf("drop table %s: %v", tbl, err)
		}
	}

	schema, err := os.ReadFile("../../../schema.sql")
	if err != nil {
		t.Fatalf("read schema.sql: %v", err)
	}
	for _, stmt := range strings.Split(string(schema), ";") {
		if strings.TrimSpace(stmt) == "" {
			continue
		}
		if _, err := conn.ExecContext(ctx, stmt); err != nil {
			t.Fatalf("apply schema statement %q: %v", stmt, err)
		}
	}
}

// seedOwnerPrincipal mirrors the first-user setup sequence in
// usecase.Authentication.VerifyAuth; keep in sync when that sequence changes.
func seedOwnerPrincipal(t *testing.T) {
	t.Helper()
	ctx := t.Context()
	conn, tx, err := db.GetTxConn(ctx)
	if err != nil {
		t.Fatalf("db.GetTxConn: %v", err)
	}
	defer func() {
		if err := conn.Close(); err != nil {
			t.Errorf("close conn: %v", err)
		}
	}()
	defer tx.AutoRollback()

	q := persistence.NewQuery(tx)
	cmd := persistence.NewCommand(tx)
	manager := membership.NewManager(q)

	token := &auth.OpenIDToken{
		Issuer:   testIssuer,
		Sub:      testSubject,
		Email:    testEmail,
		Verified: true,
		Name:     "Integration Owner",
	}
	pri, oidcAuth, err := manager.CreatePrincipalFromOpenID(ctx, token)
	if err != nil {
		t.Fatalf("CreatePrincipalFromOpenID: %v", err)
	}
	g, r, perm, err := manager.SetupPrincipalAsOwner(ctx, pri)
	if err != nil {
		t.Fatalf("SetupPrincipalAsOwner: %v", err)
	}
	roperm, err := manager.SetupSystemMembership(ctx)
	if err != nil {
		t.Fatalf("SetupSystemMembership: %v", err)
	}

	cmd.SavePermission(ctx, roperm)
	cmd.SavePermission(ctx, perm)
	cmd.SaveRole(ctx, r)
	cmd.SaveGroup(ctx, g)
	cmd.SavePrincipal(ctx, pri)
	cmd.SaveAuthOIDC(ctx, oidcAuth, pri)
	if err := cmd.Error(); err != nil {
		t.Fatalf("save owner principal: %v", err)
	}
	if err := tx.Commit(); err != nil {
		t.Fatalf("commit: %v", err)
	}
}

func ownerSessionCookie(t *testing.T, nextCheck time.Time) *http.Cookie {
	t.Helper()
	is := &usecase.IDSession{
		Issuer:  testIssuer,
		Subject: testSubject,
		Email:   testEmail,
		Membership: usecase.MembershipToken{
			NextCheck: nextCheck,
			Principal: &payload.SessionPrincipal{
				Email:       testEmail,
				Permissions: []string{membership.PermissionOwnerName},
			},
		},
	}
	sp := &SessionPayload{
		Data:     is,
		ExpireAt: time.Now().Add(config.SessionLifeTime),
	}
	value := sp.String()
	if value == "" {
		t.Fatal("failed to serialize session payload")
	}
	return &http.Cookie{Name: sessionKey, Value: value}
}

func TestIntegrationAuthRequest(t *testing.T) {
	setupIntegrationDB(t)
	seedOwnerPrincipal(t)

	r := mux.NewRouter()
	NewAuthMux().RegisterRoute(r)

	t.Run("有効なセッションCookieなら202とX-Guardmechヘッダを返す", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/auth/request", nil)
		req.AddCookie(ownerSessionCookie(t, time.Now().Add(-time.Minute)))
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		if w.Code != http.StatusAccepted {
			t.Fatalf("status = %d, want %d, body: %s", w.Code, http.StatusAccepted, w.Body.String())
		}
		if got := w.Header().Get("X-Guardmech-Email"); got != testEmail {
			t.Errorf("X-Guardmech-Email = %q, want %q", got, testEmail)
		}
		if got := w.Header().Get("X-Guardmech-Permissions"); !strings.Contains(got, membership.PermissionOwnerName) {
			t.Errorf("X-Guardmech-Permissions = %q, want contains %q", got, membership.PermissionOwnerName)
		}
	})

	t.Run("Cookieなしなら401を返す", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/auth/request", nil)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		if w.Code != http.StatusUnauthorized {
			t.Fatalf("status = %d, want %d", w.Code, http.StatusUnauthorized)
		}
	})
}

func TestIntegrationAdminAPIKeyFlow(t *testing.T) {
	setupIntegrationDB(t)
	seedOwnerPrincipal(t)

	r := mux.NewRouter()
	NewAdminMux().RegisterRoute(r)
	cookie := ownerSessionCookie(t, time.Now().Add(config.AuthorityValidationTimeout))

	postForm := func(t *testing.T, path string, form url.Values) *httptest.ResponseRecorder {
		t.Helper()
		req := httptest.NewRequest(http.MethodPost, path, strings.NewReader(form.Encode()))
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
		req.AddCookie(cookie)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)
		return w
	}

	var principalID string
	t.Run("Principalを作成できる", func(t *testing.T) {
		w := postForm(t, "/guardmech/api/principal", url.Values{
			"name":        {"api-client"},
			"description": {"integration test principal"},
		})
		if w.Code != http.StatusOK {
			t.Fatalf("status = %d, body: %s", w.Code, w.Body.String())
		}

		var res struct {
			Principal struct {
				Principal struct {
					ID string `json:"id"`
				} `json:"principal"`
			} `json:"principal"`
		}
		if err := json.Unmarshal(w.Body.Bytes(), &res); err != nil {
			t.Fatalf("unmarshal response: %v, body: %s", err, w.Body.String())
		}
		principalID = res.Principal.Principal.ID
		if principalID == "" {
			t.Fatalf("principal id is empty, body: %s", w.Body.String())
		}
	})

	var rawToken string
	t.Run("APIキーを発行でき生トークンとDTOだけが返る", func(t *testing.T) {
		if principalID == "" {
			t.Skip("principal 作成に失敗している")
		}
		w := postForm(t, "/guardmech/api/principal/"+principalID+"/new_key", url.Values{
			"name": {"integration-key"},
		})
		if w.Code != http.StatusOK {
			t.Fatalf("status = %d, body: %s", w.Code, w.Body.String())
		}

		var res struct {
			Token  string         `json:"token"`
			APIKey map[string]any `json:"api_key"`
		}
		if err := json.Unmarshal(w.Body.Bytes(), &res); err != nil {
			t.Fatalf("unmarshal response: %v, body: %s", err, w.Body.String())
		}
		rawToken = res.Token

		if !strings.HasPrefix(rawToken, "gmch-") || len(rawToken) != 48 {
			t.Fatalf("token = %q, want gmch- prefix and length 48", rawToken)
		}
		if _, ok := res.APIKey["hashed_token"]; ok {
			t.Errorf("api_key に hashed_token が露出している: %v", res.APIKey)
		}
		masked, _ := res.APIKey["masked_token"].(string)
		if want := strings.Repeat("*", 20) + rawToken[44:]; masked != want {
			t.Errorf("masked_token = %q, want %q", masked, want)
		}
	})

	t.Run("DBのhashed_tokenは生トークンのSHA-256と一致する", func(t *testing.T) {
		if rawToken == "" {
			t.Skip("APIキー発行に失敗している")
		}
		ctx := t.Context()
		conn, err := db.GetConn(ctx)
		if err != nil {
			t.Fatalf("db.GetConn: %v", err)
		}
		defer func() {
			if err := conn.Close(); err != nil {
				t.Errorf("close conn: %v", err)
			}
		}()

		var hashed string
		row := conn.QueryRowContext(ctx, "SELECT hashed_token FROM auth_apikey WHERE name = ?", "integration-key")
		if err := row.Scan(&hashed); err != nil {
			t.Fatalf("select hashed_token: %v", err)
		}
		if want := membership.HashAPIKeyToken(rawToken); hashed != want {
			t.Errorf("hashed_token = %q, want %q", hashed, want)
		}
	})
}
