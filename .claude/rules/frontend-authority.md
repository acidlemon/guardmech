---
paths:
  - "frontend/src/hooks/useUserAuthority.ts"
  - "frontend/src/components/AuthorityStatusBox.vue"
---

# useUserAuthority の AuthorityStatus 階層とセッション切れ処理

## 背景 / 罠の説明

`useUserAuthority`（`frontend/src/hooks/useUserAuthority.ts`）は `onMounted` で
`axios.get('/api/authority')` を呼び、レスポンスの `data.authority` を
`AuthorityStatus`（`'OWNER' | 'WRITABLE' | 'READONLY' | 'NONE' | 'ERROR' | ''`）として保持する。

- `canWrite` は `['OWNER', 'WRITABLE'].includes(...)`、
  `canRead` は `['OWNER', 'WRITABLE', 'READONLY'].includes(...)` という階層になっている
  （`OWNER` ⊃ `WRITABLE` ⊃ `READONLY` を暗黙に仮定した実装であり、型からは階層関係が読めない）。
  この階層は backend の `_GUARDMECH_OWNER`/`_GUARDMECH_WRITE`/`_GUARDMECH_READONLY`
  パーミッション（→ `.claude/rules/rbac-model.md`）に対応する。
- `axios.get('/api/authority')` が **catch に落ちた場合はセッション切れとみなし**
  `location.reload()` する（nginx の forward-auth を経由して sign_in にリダイレクトさせる想定）。
  レスポンスが返ったが `data.authority` が falsy な場合は `'ERROR'` として扱われ、
  reload はしない（この2つの失敗系統を混同しない）。
- `AuthorityStatusBox.vue` は `status` prop（`AuthorityStatus`）を `watch` し、
  `switch` で `message`/`variant` を出し分ける表示専用コンポーネント。
  `status` に対応する `case` が無い場合（例えば初期値の `''`）は `message`/`variant` が
  更新されず空のまま残る。
- axios のリクエストパスは `/api/...` の相対パスだが、`baseURL` は mount path 込みで
  設定されている（→ `.claude/rules/mount-path.md` の `frontend/src/main.ts` 参照）。
