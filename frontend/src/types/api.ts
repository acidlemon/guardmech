export type PrincipalPayload = {
  principal: Principal
  auth_oidc: AuthOIDC
  auth_apikeys: AuthAPIKey[]
  groups: Group[]
  attached_roles: Role[]
  having_roles: Role[]
  having_permissions: Permission[]
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
  attached_roles: Role[]
  having_permissions: Permission[]
}

export type Role = {
  id: string
  name: string
  description: string
  attached_permissions: Permission[]
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
  rule_type: number
  detail: string
  priority: number
  association_type: string
  association_id: string
}
