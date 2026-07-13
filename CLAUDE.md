# CLAUDE.md

## 1. guardmech とは

nginx の `auth_request`（forward-auth）と組み合わせて使う認証認可ゲートウェイ。

- 認証: OpenID Connect（`coreos/go-oidc/v3`）と API キー。
- 認可: Principal / Group / Role / Permission による RBAC と、ログイン時の MappingRule によるロール/グループ自動付与。
- Go 製バックエンド + Vue3（TypeScript, vue-cli）の SPA 管理画面 + MySQL。

## 2. アーキテクチャ

- nginx が保護対象への各リクエストで `auth_request /auth/request` を backend（`:2989`）に投げる。forward-auth の契約として、`/auth/request` は `202 Accepted`（許可）/ `401 Unauthorized`（未ログイン → `error_page 401 = /auth/sign_in`）/ `403 Forbidden`（`error_page 403 = /auth/unauthorized`）のいずれかのみを返す（`nginx.default.template` の `location = /auth/request` コメント参照）。
- 認可 OK 時、backend は `X-Guardmech-Email` / `X-Guardmech-Groups` / `X-Guardmech-Roles` / `X-Guardmech-Permissions` ヘッダを返し、nginx が `auth_request_set` で受けて下流へ伝播する（`nginx.default.template` 参照）。
- ローカルは `docker-compose.yml`（`guardmech` / `frontend` / `nginx` / `mysql` / `dumper_app`）で確認できる。`dumper_app`（`scalify/http-dump` イメージ）は保護対象アプリのデモ。
- 実行バイナリのエントリは `backend/cmd/guardmech/main.go` の `main()` → `backend/guardmech.go` の `GuardMech.Run()`。HTTP は `0.0.0.0:2989` で listen する。

## 3. ディレクトリ構成（クリーンアーキテクチャ的分離）

- `backend/app/logic/` … ドメイン層。`membership`（Principal/Group/Role/Permission/MappingRule/APIKey）と `auth`（OIDC 認証・GroupInquirer）。
- `backend/app/usecase/` … アプリケーション層（`authentication.go`, `admin.go` 等）。`payload/` は API レスポンス DTO。
- `backend/app/handler/` … HTTP ハンドラ（`auth.go` 認証フロー, `admin.go` 管理 API, `session.go`/`cookie.go` セッション）。
- `backend/app/config/` … タイムアウト等の定数。
- `backend/db/` … MySQL コネクションプールとトランザクション（`GetConn`/`GetTxConn`）。
- `backend/persistence/` … 永続化。`query.go`/`command.go` がドメインモデルと DB 表現を変換し、`db/` に seacle 生成の Row マッパを含む（詳細は `.claude/rules/seacle-generated.md` 参照）。
- `backend/oidconnect/` … OIDC プロバイダ抽象（`gsuite`, `generic`）。
- `backend/gen/` … seacle コード生成のエントリ（`main.go`）。
- `frontend/` … Vue3 SPA。`src/pages/`, `src/components/`, `src/hooks/`, `src/router/`。

## 4. ビルド・ローカル開発

- モジュールパスは `github.com/acidlemon/guardmech/backend`（`backend/go.mod`）。
- `Dockerfile` のバックエンドビルドステージは `golang:1.16` を使い、`GOPATH=/stash` の下の `/stash/src/github.com/acidlemon/guardmech/backend` にソースを配置してから `go get && go build -o guardmech cmd/guardmech/main.go` を実行する（`backend/Dockerfile.local` も同様の配置）。`go.sum` はコミット対象外（`.gitignore` 参照）。
- 単体でビルドする場合: `cd backend && go build -o guardmech cmd/guardmech/main.go`。
- SQL マッパの再生成: `cd backend && go generate ./...`（`backend/gen/main.go` が `prepareGenset()` の対応表に従って `persistence/db/*_gen.go` を生成する。ORM は `github.com/acidlemon/seacle`）。`docker-compose.yml` はローカル開発時に `../seacle` を `guardmech` コンテナへ隣接マウントする。
- frontend: `cd frontend && npm install && npm run serve`（dev）/ `npm run build`（本番）。`vue.config.js` の `publicPath` は `/guardmech/admin/` 固定。
- DB スキーマ: リポジトリルート `schema.sql`。
- 起動スクリプト `backend/run_guardmech.sh` は `GUARDMECH_MOUNT_PATH` 設定時に `dist/index.html` を書き換える（→ `.claude/rules/mount-path.md` 参照）。

## 5. テスト・CI

- backend: `cd backend && go test ./...`。
- frontend: `npm run type-check` / `npm run lint` / `npm run build`。
- CI:
  - `.github/workflows/go.yml`: go 1.21 で `go get -v -t -d ./...` の後 `go build -v .` のみを実行する（`go test` ステップは無い）。
  - `.github/workflows/frontend.yml`: node 16.x で `npm run type-check && npm run lint && npm run build` を実行する。
  - `.github/workflows/release.yml`: master / `feature/build-*` push とタグ push で Docker イメージをビルドし `ghcr.io` へ push する。

## 6. 主要な環境変数（すべて `GUARDMECH_` プレフィックス。値の実体は各自の環境で設定）

- OIDC: `GUARDMECH_CLIENT_ID` / `GUARDMECH_CLIENT_SECRET` / `GUARDMECH_REDIRECT_URL`
- セッション: `GUARDMECH_SIGNATURE_KEY`（HMAC 署名鍵） / `GUARDMECH_CRYPTO_KEY`（AES 暗号鍵）
- DB: `GUARDMECH_DB_HOST` / `GUARDMECH_DB_PORT` / `GUARDMECH_DB_USER` / `GUARDMECH_DB_PASSWORD` / `GUARDMECH_DB_NAME`
- マウント: `GUARDMECH_MOUNT_PATH`（サブパス配信時のプレフィックス）
- GSuite グループ照会に Google サービスアカウントを使う場合 `GOOGLE_APPLICATION_CREDENTIALS`（`backend/run_guardmech.sh` は `GOOGLE_APPLICATION_CREDENTIALS_BASE64` を base64 デコードして展開する）
