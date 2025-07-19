export interface Permission {
  id: number
  parent_id: number
  name: string
  display_name: string
  resource: string
  type: string
  action?: string
  description: string
  sort: number
  status: number
  created_at: string
  updated_at: string
  children?: Permission[]
  hasChildren?: boolean
}

export interface CreatePermissionForm {
  parent_id: number
  name: string
  display_name: string
  resource: string
  type: string
  action?: string
  description?: string
  sort?: number
  status?: number
}

export interface UpdatePermissionForm {
  parent_id?: number
  name?: string
  display_name?: string
  resource?: string
  type?: string
  action?: string
  description?: string
  sort?: number
  status?: number
}