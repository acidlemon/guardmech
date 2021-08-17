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

export type MappingRuleType = number
export const MappingRuleTypeSpecificDomain: MappingRuleType = 1
export const MappingRuleTypeWholeDomain: MappingRuleType = 2
export const MappingRuleTypeGroupMember: MappingRuleType = 3
export const MappingRuleTypeSpecificAddress: MappingRuleType = 4

export type MappingRule = {
  id: string
  name: string
  description: string
  rule_type: MappingRuleType
  detail: string
  priority: number
  association_type: string
  association_id: string
}
