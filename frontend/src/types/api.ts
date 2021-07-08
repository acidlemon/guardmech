export type PrincipalPayload = {
  principal: Principal
  auth_oidc: AuthOIDC
  auth_apikeys: AuthAPIKey[]
  groups: Group[]
  roles: Role[]
  permissions: Permission[]
}

export type AuthOIDC = {
  id: string
  issuer: string
  subject: string
  email: string
}

export type AuthAPIKey = {
  id: string
  name: string
  masked_token: string
}

export type Principal = {
  id: string
  name: string
  description: string
}

export type Group = {
  id: string
  name: string
  description: string
}

export type Role = {
  id: string
  name: string
  description: string
}

export type Permission = {
  id: string
  name: string
  description: string
}

export type MappingRule = {
  id: string
  name: string
  description: string
}
