package handler

import "testing"

func TestBearerAPIKeyToken(t *testing.T) {
	t.Parallel()

	cases := []struct {
		name      string
		header    string
		wantToken string
		wantOK    bool
	}{
		{name: "Bearer と gmch- トークン", header: "Bearer gmch-abc", wantToken: "gmch-abc", wantOK: true},
		{name: "scheme は小文字でもよい", header: "bearer gmch-abc", wantToken: "gmch-abc", wantOK: true},
		{name: "scheme は大文字でもよい", header: "BEARER gmch-abc", wantToken: "gmch-abc", wantOK: true},
		{name: "トークン前後の空白は無視する", header: "Bearer   gmch-abc  ", wantToken: "gmch-abc", wantOK: true},
		{name: "gmch- でない Bearer はフォールスルー", header: "Bearer some-other-token", wantOK: false},
		{name: "Basic はフォールスルー", header: "Basic gmch-abc", wantOK: false},
		{name: "トークンが空", header: "Bearer ", wantOK: false},
		{name: "ヘッダが空", header: "", wantOK: false},
		{name: "scheme とトークンの区切りがない", header: "Bearergmch-abc", wantOK: false},
	}

	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			token, ok := bearerAPIKeyToken(tt.header)
			if ok != tt.wantOK {
				t.Fatalf("ok = %v, want %v", ok, tt.wantOK)
			}
			if ok && token != tt.wantToken {
				t.Errorf("token = %q, want %q", token, tt.wantToken)
			}
		})
	}
}
