export interface Role {
  id: number
  name: string
  display_name: string
  description: string
  status: number
  created_at: string
  updated_at: string
  permissions?: Permission[]
}

export interface Permission {
  id: number
  name: string
  display_name: string
  description: string
  resource: string
  action: string
  status: number
}

export interface CreateRoleForm {
  name: string
  display_name: string
  description?: string
  status?: number
}

export interface UpdateRoleForm {
  name?: string
  display_name?: string
  description?: string
  status?: number
}