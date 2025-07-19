import request from '@/utils/request'
import type { Role, CreateRoleForm, UpdateRoleForm } from '@/types/role'

// 获取角色列表
export function getRoleList(params?: {
  page?: number
  limit?: number
  keyword?: string
}) {
  return request({
    url: '/roles',
    method: 'get',
    params
  })
}

// 获取角色详情
export function getRoleById(id: number) {
  return request({
    url: `/roles/${id}`,
    method: 'get'
  })
}

// 创建角色
export function createRole(data: CreateRoleForm) {
  return request({
    url: '/roles',
    method: 'post',
    data
  })
}

// 更新角色
export function updateRole(id: number, data: UpdateRoleForm) {
  return request({
    url: `/roles/${id}`,
    method: 'put',
    data
  })
}

// 删除角色
export function deleteRole(id: number) {
  return request({
    url: `/roles/${id}`,
    method: 'delete'
  })
}

// 更新角色状态
export function updateRoleStatus(id: number, status: number) {
  return request({
    url: `/roles/${id}/status`,
    method: 'put',
    data: { status }
  })
}

// 获取角色权限
export function getRolePermissions(id: number) {
  return request({
    url: `/roles/${id}/permissions`,
    method: 'get'
  })
}

// 分配角色权限
export function assignRolePermissions(id: number, permissionIds: number[]) {
  return request({
    url: `/roles/${id}/permissions`,
    method: 'put',
    data: { permission_ids: permissionIds }
  })
}