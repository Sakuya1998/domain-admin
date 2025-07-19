export interface User {
  id: number
  username: string
  email: string
  nickname: string
  role: string
  status: number
  created_at: string
  updated_at: string
}

export interface UserInfo {
  id: number
  username: string
  email: string
  nickname: string
  role: string
  status: number
}

export interface LoginForm {
  username: string
  password: string
}

export interface CreateUserForm {
  username: string
  email: string
  password: string
  nickname: string
  role: string
  status?: number
}

export interface UpdateUserForm {
  email?: string
  nickname?: string
  role?: string
  status?: number
}

export interface UpdateProfileForm {
  email?: string
  nickname?: string
}

export interface ChangePasswordForm {
  old_password: string
  new_password: string
  confirm_password: string
}