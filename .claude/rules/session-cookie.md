---
paths:
  - "backend/app/handler/session.go"
  - "backend/app/handler/cookie.go"
---

# セッション Cookie は独自形式でシリアライズされている

## 背景 / 罠の説明

セッションは JSON そのままではなく、`backend/app/handler/session.go` が独自の区切り文字で
組み立てた文字列を Cookie の値にする。

- 全体の形式: `data('-'*)signature`
- `data` 部の形式: `encVal(#'-')expireUnix`

（`SessionPayload.String()` で組み立て、`RestoreSessionPayload()` で
`strings.Split(fromCookie, "('-'*)")` → `strings.Split(data, "(#'-')")` の順に分解して復元する）

- 署名は HMAC-SHA256（`backend/app/handler/cookie.go` の `signator`、鍵は環境変数
  `GUARDMECH_SIGNATURE_KEY`）。
- 本文（`Data.String()` の結果）は AES-GCM で暗号化される（`cryptor`、鍵は環境変数
  `GUARDMECH_CRYPTO_KEY`）。`Encrypt()` は生成した nonce を暗号文の先頭に連結して base64 化し、
  `Decrypt()` は先頭 `gcm.NonceSize()` バイトを nonce として読み戻す。
- `GUARDMECH_CRYPTO_KEY` は AES の鍵長制約（16/24/32 バイト）を満たす必要がある。
  満たさない場合 `cryptor.Encrypt`/`Decrypt` 内の `aes.NewCipher` がエラーを返す。
- Cookie 名は用途で分かれている: CSRF 用セッションは `_guardmech_csrf`
  （`authSessionKey`）、ログイン後の本セッションは `_guardmech_session`
  （`sessionKey`。いずれも `backend/app/handler/auth.go` で定義）。
