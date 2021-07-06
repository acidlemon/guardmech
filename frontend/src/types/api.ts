export type PrincipalPayload = {
  principal: Principal
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
