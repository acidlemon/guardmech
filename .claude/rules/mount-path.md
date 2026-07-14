---
paths:
  - "backend/guardmech.go"
  - "backend/app/handler/auth.go"
  - "backend/app/handler/common.go"
  - "backend/run_guardmech.sh"
  - "frontend/vue.config.js"
  - "frontend/src/main.ts"
  - "frontend/src/router/index.ts"
---

# GUARDMECH_MOUNT_PATH 対応は複数箇所に分散している

## 背景 / 罠の説明

`GUARDMECH_MOUNT_PATH`（サブパス配信用のプレフィックス）への対応は
バックエンドとフロントエンドの複数箇所に分かれており、片方だけ直すと壊れる。

- `backend/guardmech.go`: `Run()` 内で `mount := os.Getenv("GUARDMECH_MOUNT_PATH")` を取得し、
  空でなければ `root.PathPrefix(mount).Subrouter()` でルータを切り替える。
  静的アセットは `http.StripPrefix(mount+"/guardmech/admin/js", ...)` のように
  `mount` を手動で前置した文字列で `StripPrefix` している。
- `backend/app/handler/auth.go`: `SignIn()` がリダイレクト先や埋め込み HTML のフォーム action に
  `os.Getenv("GUARDMECH_MOUNT_PATH")` を都度前置する（ソース中に「TODO とても行儀が悪いので直す」の
  コメントがある）。
- `backend/app/handler/common.go`: `WriteHttpError()` がエラーページの戻り先リンクにも
  同様に `GUARDMECH_MOUNT_PATH` を前置する。
- `backend/run_guardmech.sh`: コンテナ起動時、`GUARDMECH_MOUNT_PATH` が設定されていれば
  `perl -pi -e` で `dist/index.html` の `href="/guardmech/admin...` / `src="/guardmech/admin...` を
  書き換える。**ビルド時ではなく起動時パッチ**なので、`dist/index.html` を直接見ても
  ローカルビルド直後は mount path が反映されていない。
- `frontend/src/main.ts`: `document.currentScript.src` の URL から
  `/admin/js/` の手前までを切り出して `Axios.defaults.baseURL` に設定する。
- `frontend/src/router/index.ts`: 同様に `document.currentScript.src` から
  `/js/` の手前までを `createWebHistory()` の base に使う。
  **main.ts のアンカー文字列は `/admin/js/`、router 側は `/js/` で異なる**ため、
  片方だけ書き換えると片方が反応しなくなる。
- `frontend/vue.config.js`: `publicPath: '/guardmech/admin/'` で固定。ここは mount path 非依存。

## チェックリスト

- mount path 関連を触るときは、上記のうち関係する箇所を横断で洗い出してから直す
- `main.ts` と `router/index.ts` のアンカー文字列（`/admin/js/` と `/js/`）を
  それぞれ独立に確認する（片方を直したつもりでもう片方が直っていない、という事故が起きやすい）
- `backend/run_guardmech.sh` のパッチはビルド成果物に対する起動時書き換えである点を踏まえ、
  動作確認は実際にコンテナを起動した状態で行う
