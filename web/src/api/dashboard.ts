import request from '@/utils/request'

export interface DashboardStats {
  userCount: number
  roleCount: number
  permissionCount: number
  onlineCount: number
}

// 获取仪表盘统计数据
export function getDashboardStats(): Promise<DashboardStats> {
  return request({
    url: '/dashboard/stats',
    method: 'get'
  }).then((response: any) => response.data)
}