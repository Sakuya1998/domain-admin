import { defineStore } from 'pinia'
import { ref } from 'vue'
import { login, getUserProfile } from '@/api/auth'
import { removeToken, setToken, getToken } from '@/utils/auth'
import type { LoginForm, UserInfo } from '@/types/user'

export const useUserStore = defineStore('user', () => {
  const token = ref<string>(getToken() || '')
  const userInfo = ref<UserInfo | null>(null)
  
  // 登录
  const userLogin = async (loginForm: LoginForm) => {
    try {
      const response = await login(loginForm)
      const { data } = response
      const { token: userToken } = data
      
      token.value = userToken
      setToken(userToken)
      
      return response
    } catch (error) {
      throw error
    }
  }
  
  // 获取用户信息
  const getUserInfo = async () => {
    try {
      const response = await getUserProfile()
      const { data } = response
      userInfo.value = data
      return response
    } catch (error) {
      throw error
    }
  }
  
  // 登出
  const logout = () => {
    token.value = ''
    userInfo.value = null
    removeToken()
  }
  
  return {
    token,
    userInfo,
    userLogin,
    getUserInfo,
    logout
  }
})