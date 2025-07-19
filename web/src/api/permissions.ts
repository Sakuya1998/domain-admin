import request from '@/utils/request'
import type { Permission, CreatePermissionForm, UpdatePermissionForm } from '@/types/permission'

// 获取权限列表
export function getPermissionList(params?: {
  page?: number
  limit?: number
  keyword?: string
}) {
  return request({
    url: '/permissions',
    method: 'get',
    params
  })
}

// 获取权限详情
export function getPermissionById(id: number) {
  return request({
    url: `/permissions/${id}`,
    method: 'get'
  })
}

// 创建权限
export function createPermission(data: CreatePermissionForm) {
  return request({
    url: '/permissions',
    method: 'post',
    data
  })
}

// 更新权限
export function updatePermission(id: number, data: UpdatePermissionForm) {
  return request({
    url: `/permissions/${id}`,
    method: 'put',
    data
  })
}

// 删除权限
export function deletePermission(id: number) {
  return request({
    url: `/permissions/${id}`,
    method: 'delete'
  })
}

// 更新权限状态
export function updatePermissionStatus(id: number, status: number) {
  return request({
    url: `/permissions/${id}/status`,
    method: 'put',
    data: { status }
  })
}