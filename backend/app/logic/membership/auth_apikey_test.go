package membership

import (
	"encoding/hex"
	"strings"
	"testing"
)

func TestCreateAPIKey(t *testing.T) {
	t.Parallel()

	t.Run("name空文字ならエラー", func(t *testing.T) {
		t.Parallel()
		p := &Principal{}
		if _, _, err := p.CreateAPIKey(""); err == nil {
			t.Error("name空文字でエラーが返らなかった")
		}
	})

	t.Run("正常系", func(t *testing.T) {
		t.Parallel()
		p := &Principal{}
		key, rawToken, err := p.CreateAPIKey("test-key")
		if err != nil {
			t.Fatalf("CreateAPIKey: %v", err)
		}

		if !strings.HasPrefix(rawToken, "gmch-") {
			t.Errorf("rawToken が gmch- 始まりでない: %q", rawToken)
		}
		if len(rawToken) != 48 {
			t.Errorf("rawToken の長さ = %d, want 48", len(rawToken))
		}
		if key.HashedToken != HashAPIKeyToken(rawToken) {
			t.Errorf("HashedToken = %q, want HashAPIKeyToken(rawToken) = %q", key.HashedToken, HashAPIKeyToken(rawToken))
		}
		if want := strings.Repeat("*", 20) + rawToken[44:]; key.MaskedToken != want {
			t.Errorf("MaskedToken = %q, want %q", key.MaskedToken, want)
		}
	})

	t.Run("連続呼び出しでトークンとIDが異なる", func(t *testing.T) {
		t.Parallel()
		p := &Principal{}
		key1, token1, err := p.CreateAPIKey("key1")
		if err != nil {
			t.Fatalf("CreateAPIKey(1回目): %v", err)
		}
		key2, token2, err := p.CreateAPIKey("key2")
		if err != nil {
			t.Fatalf("CreateAPIKey(2回目): %v", err)
		}
		if token1 == token2 {
			t.Errorf("連続呼び出しで同じトークンが返った: %q", token1)
		}
		if key1.AuthAPIKeyID == key2.AuthAPIKeyID {
			t.Errorf("連続呼び出しで同じAuthAPIKeyIDが返った: %s", key1.AuthAPIKeyID)
		}
	})
}

func TestHashAPIKeyToken(t *testing.T) {
	t.Parallel()

	t.Run("同一入力なら同一出力", func(t *testing.T) {
		t.Parallel()
		if HashAPIKeyToken("gmch-abc") != HashAPIKeyToken("gmch-abc") {
			t.Error("同一入力で異なるハッシュが返った")
		}
	})

	t.Run("既知入力のSHA-256ダイジェストと一致する", func(t *testing.T) {
		t.Parallel()
		want := "dcfa279c42d658397b282e52daa0c7b75a7ace438a090e821bcb3d8b7152e506"
		if got := HashAPIKeyToken("gmch-abc"); got != want {
			t.Errorf("HashAPIKeyToken(\"gmch-abc\") = %q, want %q", got, want)
		}
	})

	t.Run("出力は64文字のhex", func(t *testing.T) {
		t.Parallel()
		got := HashAPIKeyToken("gmch-abc")
		if len(got) != 64 {
			t.Errorf("len = %d, want 64", len(got))
		}
		if _, err := hex.DecodeString(got); err != nil {
			t.Errorf("hex として decode できない: %q", got)
		}
	})

	t.Run("異なる入力なら異なる出力", func(t *testing.T) {
		t.Parallel()
		if HashAPIKeyToken("gmch-abc") == HashAPIKeyToken("gmch-abd") {
			t.Error("異なる入力で同じハッシュが返った")
		}
	})
}
