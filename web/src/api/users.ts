import request from '@/utils/request'
import type { User, CreateUserForm, UpdateUserForm } from '@/types/user'

// 获取用户列表
export function getUserList(params?: {
  page?: number
  limit?: number
  keyword?: string
}) {
  return request({
    url: '/users',
    method: 'get',
    params
  })
}

// 获取用户详情
export function getUserById(id: number) {
  return request({
    url: `/users/${id}`,
    method: 'get'
  })
}

// 创建用户
export function createUser(data: CreateUserForm) {
  return request({
    url: '/users',
    method: 'post',
    data
  })
}

// 更新用户
export function updateUser(id: number, data: UpdateUserForm) {
  return request({
    url: `/users/${id}`,
    method: 'put',
    data
  })
}

// 删除用户
export function deleteUser(id: number) {
  return request({
    url: `/users/${id}`,
    method: 'delete'
  })
}

// 更新用户状态
export function updateUserStatus(id: number, status: number) {
  return request({
    url: `/users/${id}/status`,
    method: 'put',
    data: { status }
  })
}