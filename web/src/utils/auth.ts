import Cookies from 'js-cookie'

const TOKEN_KEY = 'domain_admin_token'

export function getToken(): string | undefined {
  return Cookies.get(TOKEN_KEY)
}

export function setToken(token: string): void {
  Cookies.set(TOKEN_KEY, token, { expires: 7 }) // 7天过期
}

export function removeToken(): void {
  Cookies.remove(TOKEN_KEY)
}