---
paths:
  - "backend/app/logic/auth/authenticator.go"
  - "backend/persistence/db/auth_oidc.go"
  - "backend/persistence/db/auth_oidc_gen.go"
---

# OIDC issuer は "https://" を除去した形で保存されている

## 背景 / 罠の説明

`backend/app/logic/auth/authenticator.go` の `Authenticator.VerifyAuthentication()` は、
ID トークンの claims を取得した直後に

```go
claims.Issuer = strings.Replace(claims.Issuer, "https://", "", -1)
```

で issuer を正規化してから呼び出し元へ返す。

そのため `auth_oidc` テーブル（`backend/persistence/db/auth_oidc.go` の `AuthOIDCRow.Issuer`）に
保存される issuer は、Google であれば `accounts.google.com`（`https://` なし）になる。
issuer を条件に検索・比較するコード（`FindPrincipalByOIDC` 等）を書くときは、
この正規化済みの文字列を前提にする。生の OIDC discovery document の issuer
（`https://accounts.google.com` 等）とそのまま比較しても一致しない。

`schema.sql` の `auth_oidc` テーブルは `UNIQUE uniq_issuer_subject (issuer, subject)` を持ち、
この正規化後の issuer と subject の組でユニーク制約が張られている。
