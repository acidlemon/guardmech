# guardmech

guardmech is an authentication and authorization gateway designed to be used
together with nginx `auth_request` (forward-auth). It sits in front of your
services and decides, on every request, whether the caller is allowed in.

## Authentication and authorization

- **Authentication** supports OpenID Connect (via `github.com/coreos/go-oidc/v3`)
  and API keys.
- **Authorization** is role-based (RBAC), built around Principals, Groups,
  Roles, and Permissions. Mapping rules can automatically assign roles and
  groups to a Principal at sign-in time, based on things like the
  authenticated domain or existing group membership.

## How it works

For each incoming request, nginx issues a subrequest to
`auth_request /auth/request`, which guardmech serves on port `2989`.
guardmech responds with:

- `202 Accepted` if the request is authenticated and authorized
- `401 Unauthorized` if the request is not authenticated
- `403 Forbidden` if the request is authenticated but not authorized

On a successful response, guardmech returns the resolved identity as
`X-Guardmech-Email`, `X-Guardmech-Groups`, `X-Guardmech-Roles`, and
`X-Guardmech-Permissions` headers, which nginx propagates to the
downstream service via `auth_request_set`.

## Components

- **Backend**: a Go application (entry point at `backend/cmd/guardmech/main.go`)
  that implements authentication, authorization, and an admin API.
- **Frontend**: a Vue 3 (TypeScript) single-page application for administering
  Principals, Groups, Roles, Permissions, and mapping rules.
- **MySQL**: stores accounts and authorization data (see `schema.sql`).

For local development, all of the above are run together via
`docker-compose.yml`.

## License

MIT License. See [LICENSE](LICENSE) for details.
