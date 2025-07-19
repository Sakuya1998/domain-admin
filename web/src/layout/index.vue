<template>
  <el-container class="layout-container">
    <!-- 侧边栏 -->
    <el-aside :width="isCollapse ? '64px' : '200px'" class="sidebar">
      <div class="logo">
        <img v-if="!isCollapse" src="/vite.svg" alt="Logo" class="logo-img" />
        <span v-if="!isCollapse" class="logo-text">Domain Admin</span>
        <img v-else src="/vite.svg" alt="Logo" class="logo-img-mini" />
      </div>
      
      <el-menu
        :default-active="$route.path"
        :collapse="isCollapse"
        :unique-opened="true"
        router
        class="sidebar-menu"
      >
        <el-menu-item index="/dashboard">
          <el-icon><Odometer /></el-icon>
          <template #title>仪表盘</template>
        </el-menu-item>
        
        <el-sub-menu index="/system">
          <template #title>
            <el-icon><Setting /></el-icon>
            <span>系统管理</span>
          </template>
          <el-menu-item index="/system/users">
            <el-icon><User /></el-icon>
            <template #title>用户管理</template>
          </el-menu-item>
          <el-menu-item index="/system/roles">
            <el-icon><Avatar /></el-icon>
            <template #title>角色管理</template>
          </el-menu-item>
          <el-menu-item index="/system/permissions">
            <el-icon><Lock /></el-icon>
            <template #title>权限管理</template>
          </el-menu-item>
        </el-sub-menu>
      </el-menu>
    </el-aside>
    
    <!-- 主内容区 -->
    <el-container>
      <!-- 顶部导航 -->
      <el-header class="header">
        <div class="header-left">
          <el-button
            type="text"
            @click="toggleCollapse"
            class="collapse-btn"
          >
            <el-icon><Expand v-if="isCollapse" /><Fold v-else /></el-icon>
          </el-button>
        </div>
        
        <div class="header-right">
          <el-dropdown @command="handleCommand">
            <span class="user-info">
              <el-avatar :size="32" class="user-avatar">
                {{ userStore.userInfo?.nickname?.charAt(0) || 'U' }}
              </el-avatar>
              <span class="username">{{ userStore.userInfo?.nickname || '用户' }}</span>
              <el-icon class="el-icon--right"><ArrowDown /></el-icon>
            </span>
            <template #dropdown>
              <el-dropdown-menu>
                <el-dropdown-item command="profile">个人中心</el-dropdown-item>
                <el-dropdown-item command="logout" divided>退出登录</el-dropdown-item>
              </el-dropdown-menu>
            </template>
          </el-dropdown>
        </div>
      </el-header>
      
      <!-- 主内容 -->
      <el-main class="main-content">
        <router-view />
      </el-main>
    </el-container>
  </el-container>
</template>

<script setup lang="ts">
import { ref } from 'vue'
import { useRouter } from 'vue-router'
import { useUserStore } from '@/stores/user'
import {
  Odometer,
  User,
  Avatar,
  Lock,
  Setting,
  Expand,
  Fold,
  ArrowDown
} from '@element-plus/icons-vue'
import { ElMessageBox } from 'element-plus'

const router = useRouter()
const userStore = useUserStore()

const isCollapse = ref(false)

const toggleCollapse = () => {
  isCollapse.value = !isCollapse.value
}

const handleCommand = async (command: string) => {
  switch (command) {
    case 'profile':
      router.push('/profile')
      break
    case 'logout':
      try {
        await ElMessageBox.confirm('确定要退出登录吗？', '提示', {
          confirmButtonText: '确定',
          cancelButtonText: '取消',
          type: 'warning'
        })
        
        userStore.logout()
        router.push('/login')
      } catch {
        // 用户取消
      }
      break
  }
}
</script>

<style scoped>
.layout-container {
  height: 100vh;
}

.sidebar {
  background: linear-gradient(180deg, #667eea 0%, #764ba2 100%);
  transition: width 0.3s;
  overflow: hidden;
  box-shadow: 2px 0 8px rgba(0, 0, 0, 0.1);
}

.logo {
  height: 60px;
  display: flex;
  align-items: center;
  justify-content: center;
  background: rgba(255, 255, 255, 0.1);
  backdrop-filter: blur(10px);
  color: #fff;
  font-size: 16px;
  font-weight: bold;
  border-bottom: 1px solid rgba(255, 255, 255, 0.1);
}

.logo-img {
  width: 32px;
  height: 32px;
  margin-right: 8px;
  filter: brightness(1.2);
}

.logo-img-mini {
  width: 32px;
  height: 32px;
  filter: brightness(1.2);
}

.logo-text {
  white-space: nowrap;
  text-shadow: 0 1px 2px rgba(0, 0, 0, 0.3);
}

.sidebar-menu {
  border: none;
  background: transparent;
}

.sidebar-menu .el-menu-item {
  color: rgba(10, 92, 87, 0.85);
  border-radius: 8px;
  margin: 4px 8px;
  transition: all 0.3s ease;
  background: transparent;
  font-weight: 500;
}

.sidebar-menu .el-menu-item:hover {
  background: rgba(255, 255, 255, 0.15);
  color: #fff;
  transform: translateX(4px);
}

.sidebar-menu .el-menu-item.is-active {
  background: rgba(255, 255, 255, 0.2);
  color: #fff;
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.2);
}

.sidebar-menu .el-sub-menu {
  /* 移除 overflow: hidden 以确保子菜单可见 */
}

.sidebar-menu .el-sub-menu .el-sub-menu__title {
  color: rgba(255, 255, 255, 0.45);
  border-radius: 8px;
  margin: 4px 8px;
  transition: all 0.3s ease;
  position: relative;
  background: transparent;
  font-weight: 500;
}

.sidebar-menu .el-sub-menu .el-sub-menu__title:hover {
  background: rgba(255, 255, 255, 0.15);
  color: #fff;
  transform: translateX(4px);
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.1);
}

.sidebar-menu .el-sub-menu.is-opened .el-sub-menu__title {
  background: rgba(255, 255, 255, 0.2);
  color: #fff;
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.15);
}

.sidebar-menu .el-sub-menu .el-sub-menu__title .el-sub-menu__icon-arrow {
  transition: transform 0.3s ease;
}

.sidebar-menu .el-sub-menu.is-opened .el-sub-menu__title .el-sub-menu__icon-arrow {
  transform: rotateZ(180deg);
}

.sidebar-menu .el-sub-menu .el-menu {
  background: #8780c4 !important;
  border-radius: 0 0 12px 12px;
  margin: 0 8px 8px 8px;
  padding: 8px 0;
  box-shadow: inset 0 2px 4px rgba(0, 0, 0, 0.15);
  backdrop-filter: blur(8px);
  border: 1px solid rgba(255, 255, 255, 0.1);
}

.sidebar-menu .el-sub-menu ul.el-menu {
  background-color: #8780c4 !important;
  background: #8780c4 !important;
}

.sidebar-menu .el-sub-menu .el-menu.el-menu--inline {
  background: #8780c4 !important;
  background-color: #8780c4 !important;
}

ul.el-menu.el-menu--inline[role="menu"] {
  background: #8780c4 !important;
  background-color: #8780c4 !important;
}

.sidebar-menu ul.el-menu.el-menu--inline[role="menu"] {
  background: #8780c4 !important;
  background-color: #8780c4 !important;
}

.el-sub-menu ul.el-menu.el-menu--inline[role="menu"] {
  background: #8780c4 !important;
  background-color: #8780c4 !important;
}

.sidebar-menu .el-sub-menu ul.el-menu.el-menu--inline[role="menu"] {
  background: #8780c4 !important;
  background-color: #8780c4 !important;
}

.sidebar-menu .el-sub-menu .el-menu-item {
  background: transparent;
  color: #333333 !important;
  border-radius: 6px;
  margin: 2px 16px;
  transition: all 0.3s ease;
  position: relative;
  overflow: hidden;
}

.sidebar-menu .el-sub-menu .el-menu-item span {
  color: #333333 !important;
  font-weight: 500;
}

.sidebar-menu .el-sub-menu .el-menu-item .el-icon {
  color: #333333 !important;
}

.sidebar-menu .el-sub-menu .el-menu-item::before {
  content: '';
  position: absolute;
  left: 0;
  top: 0;
  width: 3px;
  height: 100%;
  background: linear-gradient(180deg, #667eea 0%, #764ba2 100%);
  transform: scaleY(0);
  transition: transform 0.3s ease;
  border-radius: 0 3px 3px 0;
}

.sidebar-menu .el-sub-menu .el-menu-item:hover {
  background: rgba(255, 255, 255, 0.15);
  color: #fff;
  transform: translateX(6px);
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.15);
}

.sidebar-menu .el-sub-menu .el-menu-item:hover::before {
  transform: scaleY(1);
}

.sidebar-menu .el-sub-menu .el-menu-item.is-active {
  background: rgba(255, 255, 255, 0.25);
  color: #fff;
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.2);
  transform: translateX(6px);
}

.sidebar-menu .el-sub-menu .el-menu-item.is-active::before {
  transform: scaleY(1);
}

.header {
  background-color: #fff;
  border-bottom: 1px solid #e4e7ed;
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 0 20px;
  box-shadow: 0 1px 4px rgba(0, 21, 41, 0.08);
}

.header-left {
  display: flex;
  align-items: center;
}

.collapse-btn {
  font-size: 18px;
  color: #606266;
}

.header-right {
  display: flex;
  align-items: center;
}

.user-info {
  display: flex;
  align-items: center;
  cursor: pointer;
  padding: 0 12px;
  border-radius: 4px;
  transition: background-color 0.3s;
}

.user-info:hover {
  background-color: #f5f7fa;
}

.user-avatar {
  margin-right: 8px;
  background-color: #409eff;
}

.username {
  margin-right: 4px;
  color: #606266;
  font-size: 14px;
}

.main-content {
  background-color:rgb(240, 244, 245);
  padding: 0;
  overflow-y: auto;
}
</style>