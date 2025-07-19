import { createRouter, createWebHistory } from 'vue-router'
import type { RouteRecordRaw } from 'vue-router'
import { useUserStore } from '@/stores/user'
import NProgress from 'nprogress'
import 'nprogress/nprogress.css'

NProgress.configure({ showSpinner: false })

const routes: RouteRecordRaw[] = [
  {
    path: '/login',
    name: 'Login',
    component: () => import('@/views/login/index.vue'),
    meta: {
      title: '登录',
      requiresAuth: false
    }
  },
  {
    path: '/',
    name: 'Layout',
    component: () => import('@/layout/index.vue'),
    redirect: '/dashboard',
    meta: {
      requiresAuth: true
    },
    children: [
      {
        path: '/dashboard',
        name: 'Dashboard',
        component: () => import('@/views/dashboard/index.vue'),
        meta: {
          title: '仪表盘',
          icon: 'Dashboard'
        }
      },
      {
        path: '/system',
        name: 'System',
        redirect: '/system/users',
        meta: {
          title: '系统管理',
          icon: 'Setting'
        },
        children: [
          {
            path: '/system/users',
            name: 'Users',
            component: () => import('@/views/users/index.vue'),
            meta: {
              title: '用户管理',
              icon: 'User'
            }
          },
          {
            path: '/system/roles',
            name: 'Roles',
            component: () => import('@/views/roles/index.vue'),
            meta: {
              title: '角色管理',
              icon: 'UserFilled'
            }
          },
          {
            path: '/system/permissions',
            name: 'Permissions',
            component: () => import('@/views/permissions/index.vue'),
            meta: {
              title: '权限管理',
              icon: 'Lock'
            }
          }
        ]
      },
      {
        path: '/profile',
        name: 'Profile',
        component: () => import('@/views/profile/index.vue'),
        meta: {
          title: '个人中心',
          icon: 'User'
        }
      }
    ]
  },
  {
    path: '/:pathMatch(.*)*',
    name: 'NotFound',
    component: () => import('@/views/404.vue')
  }
]

const router = createRouter({
  history: createWebHistory(),
  routes
})

// 路由守卫
router.beforeEach(async (to, from, next) => {
  NProgress.start()
  
  const userStore = useUserStore()
  const token = userStore.token
  
  if (to.meta.requiresAuth !== false) {
    if (!token) {
      next('/login')
      return
    }
    
    // 如果有token但没有用户信息，获取用户信息
    if (!userStore.userInfo) {
      try {
        await userStore.getUserInfo()
      } catch (error) {
        // 静默处理token失效，不显示错误消息
        userStore.logout()
        next('/login')
        return
      }
    }
  }
  
  if (to.path === '/login' && token) {
    next('/')
    return
  }
  
  next()
})

router.afterEach(() => {
  NProgress.done()
})

export default router