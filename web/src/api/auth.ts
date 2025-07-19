import request from '@/utils/request'
import type { LoginForm, UserInfo } from '@/types/user'

// 用户登录
export function login(data: LoginForm) {
  return request({
    url: '/auth/login',
    method: 'post',
    data
  })
}

// 获取用户信息
export function getUserProfile() {
  return request({
    url: '/auth/profile',
    method: 'get'
  })
}

// 更新用户信息
export function updateUserProfile(data: Partial<UserInfo>) {
  return request({
    url: '/auth/profile',
    method: 'put',
    data
  })
}