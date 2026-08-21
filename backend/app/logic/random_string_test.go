package logic

import (
	"strings"
	"testing"
)

func TestGenerateRandomString(t *testing.T) {
	t.Parallel()

	cases := []struct {
		name   string
		length int
	}{
		{name: "長さ0なら空文字列", length: 0},
		{name: "長さ1", length: 1},
		{name: "長さ43", length: 43},
	}

	allowed := "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			got := GenerateRandomString(tt.length)
			if len(got) != tt.length {
				t.Fatalf("len = %d, want %d", len(got), tt.length)
			}
			for _, r := range got {
				if !strings.ContainsRune(allowed, r) {
					t.Errorf("生成文字に許可外のルーン %q が含まれる: %q", r, got)
				}
			}
		})
	}
}

func TestGenerateRandomStringUniqueness(t *testing.T) {
	t.Parallel()

	t.Run("連続呼び出しで異なる値になる", func(t *testing.T) {
		t.Parallel()

		first := GenerateRandomString(43)
		second := GenerateRandomString(43)
		if first == second {
			t.Errorf("連続呼び出しで同じ値が返った: %q", first)
		}
	})
}
