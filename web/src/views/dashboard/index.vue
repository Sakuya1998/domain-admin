<template>
  <div class="page-container">
    <div class="dashboard-header">
      <h1 class="page-title">仪表盘</h1>
      <p class="page-description">欢迎使用 Domain Admin 域名管理系统</p>
    </div>
    
    <!-- 统计卡片 -->
    <el-row :gutter="20" class="stats-row">
      <el-col :xs="24" :sm="12" :md="6">
        <div class="stat-card">
          <div class="stat-icon user-icon">
            <el-icon><User /></el-icon>
          </div>
          <div class="stat-content">
            <div class="stat-number">{{ stats.userCount }}</div>
            <div class="stat-label">用户总数</div>
          </div>
        </div>
      </el-col>
      
      <el-col :xs="24" :sm="12" :md="6">
        <div class="stat-card">
          <div class="stat-icon role-icon">
            <el-icon><UserFilled /></el-icon>
          </div>
          <div class="stat-content">
            <div class="stat-number">{{ stats.roleCount }}</div>
            <div class="stat-label">角色总数</div>
          </div>
        </div>
      </el-col>
      
      <el-col :xs="24" :sm="12" :md="6">
        <div class="stat-card">
          <div class="stat-icon permission-icon">
            <el-icon><Lock /></el-icon>
          </div>
          <div class="stat-content">
            <div class="stat-number">{{ stats.permissionCount }}</div>
            <div class="stat-label">权限总数</div>
          </div>
        </div>
      </el-col>
      
      <el-col :xs="24" :sm="12" :md="6">
        <div class="stat-card">
          <div class="stat-icon online-icon">
            <el-icon><Connection /></el-icon>
          </div>
          <div class="stat-content">
            <div class="stat-number">{{ stats.onlineCount }}</div>
            <div class="stat-label">在线用户</div>
          </div>
        </div>
      </el-col>
    </el-row>
    
    <!-- 快速操作 -->
    <div class="content-card quick-actions">
      <h3 class="section-title">快速操作</h3>
      <el-row :gutter="20">
        <el-col :xs="24" :sm="12" :md="8">
          <el-card class="action-card" @click="$router.push('/system/users')">
            <div class="action-content">
              <el-icon class="action-icon"><UserFilled /></el-icon>
              <div class="action-text">
                <div class="action-title">用户管理</div>
                <div class="action-desc">管理系统用户信息</div>
              </div>
            </div>
          </el-card>
        </el-col>
        
        <el-col :xs="24" :sm="12" :md="8">
          <el-card class="action-card" @click="$router.push('/system/roles')">
            <div class="action-content">
              <el-icon class="action-icon"><Avatar /></el-icon>
              <div class="action-text">
                <div class="action-title">角色管理</div>
                <div class="action-desc">配置用户角色权限</div>
              </div>
            </div>
          </el-card>
        </el-col>
        
        <el-col :xs="24" :sm="12" :md="8">
          <el-card class="action-card" @click="$router.push('/system/permissions')">
            <div class="action-content">
              <el-icon class="action-icon"><Key /></el-icon>
              <div class="action-text">
                <div class="action-title">权限管理</div>
                <div class="action-desc">管理系统访问权限</div>
              </div>
            </div>
          </el-card>
        </el-col>
      </el-row>
    </div>
    
    <!-- 系统信息 -->
    <div class="content-card system-info">
      <h3 class="section-title">系统信息</h3>
      <el-descriptions :column="2" border>
        <el-descriptions-item label="系统名称">Domain Admin</el-descriptions-item>
        <el-descriptions-item label="系统版本">v1.0.0</el-descriptions-item>
        <el-descriptions-item label="当前用户">{{ userStore.userInfo?.nickname }}</el-descriptions-item>
        <el-descriptions-item label="用户角色">{{ userStore.userInfo?.role }}</el-descriptions-item>
        <el-descriptions-item label="登录时间">{{ loginTime }}</el-descriptions-item>
        <el-descriptions-item label="系统状态">
          <el-tag type="success">运行正常</el-tag>
        </el-descriptions-item>
      </el-descriptions>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { useUserStore } from '@/stores/user'
import { getDashboardStats, type DashboardStats } from '@/api/dashboard'
import { ElMessage } from 'element-plus'
import {
  User,
  UserFilled,
  Lock,
  Connection,
  Avatar,
  Key
} from '@element-plus/icons-vue'

const userStore = useUserStore()

const stats = ref<DashboardStats>({
  userCount: 0,
  roleCount: 0,
  permissionCount: 0,
  onlineCount: 0
})

const loginTime = ref('')
const loading = ref(false)

// 获取统计数据
const fetchStats = async () => {
  try {
    loading.value = true
    const data = await getDashboardStats()
    stats.value = data
  } catch (error) {
    console.error('获取统计数据失败:', error)
    ElMessage.error('获取统计数据失败')
  } finally {
    loading.value = false
  }
}

onMounted(() => {
  fetchStats()
  loginTime.value = new Date().toLocaleString()
})
</script>

<style scoped>
.page-container {
  padding: 24px;
  background-color: #f0f2f5;
  min-height: calc(100vh - 60px);
}

.content-card {
  background: #fff;
  border-radius: 8px;
  box-shadow: 0 2px 12px rgba(0, 0, 0, 0.1);
}

.dashboard-header {
  margin-bottom: 30px;
}

.page-title {
  font-size: 28px;
  font-weight: 600;
  color: #2c3e50;
  margin: 0 0 8px 0;
}

.page-description {
  font-size: 16px;
  color: #7f8c8d;
  margin: 0;
}

.stats-row {
  margin-bottom: 30px;
}

.stat-card {
  background: #fff;
  border-radius: 8px;
  padding: 24px;
  box-shadow: 0 2px 12px rgba(0, 0, 0, 0.1);
  display: flex;
  align-items: center;
  transition: transform 0.3s ease;
}

.stat-card:hover {
  transform: translateY(-2px);
}

.stat-icon {
  width: 60px;
  height: 60px;
  border-radius: 12px;
  display: flex;
  align-items: center;
  justify-content: center;
  margin-right: 16px;
  font-size: 24px;
  color: #fff;
}

.user-icon {
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
}

.role-icon {
  background: linear-gradient(135deg, #f093fb 0%, #f5576c 100%);
}

.permission-icon {
  background: linear-gradient(135deg, #4facfe 0%, #00f2fe 100%);
}

.online-icon {
  background: linear-gradient(135deg, #43e97b 0%, #38f9d7 100%);
}

.stat-content {
  flex: 1;
}

.stat-number {
  font-size: 32px;
  font-weight: 700;
  color: #2c3e50;
  line-height: 1;
  margin-bottom: 4px;
}

.stat-label {
  font-size: 14px;
  color: #7f8c8d;
}

.section-title {
  font-size: 20px;
  font-weight: 600;
  color: #2c3e50;
  margin: 0 0 20px 0;
}

.quick-actions {
  padding: 24px;
  margin-bottom: 30px;
}

.action-card {
  cursor: pointer;
  transition: all 0.3s ease;
  margin-bottom: 20px;
}

.action-card:hover {
  transform: translateY(-2px);
  box-shadow: 0 8px 25px rgba(0, 0, 0, 0.15);
}

.action-content {
  display: flex;
  align-items: center;
  padding: 8px;
}

.action-icon {
  font-size: 32px;
  color: #409eff;
  margin-right: 16px;
}

.action-text {
  flex: 1;
}

.action-title {
  font-size: 16px;
  font-weight: 600;
  color: #2c3e50;
  margin-bottom: 4px;
}

.action-desc {
  font-size: 14px;
  color: #7f8c8d;
}

.system-info {
  padding: 24px;
}

/* 响应式设计 */
@media (max-width: 768px) {
  .stat-card {
    margin-bottom: 16px;
  }
  
  .action-card {
    margin-bottom: 16px;
  }
  
  .page-title {
    font-size: 24px;
  }
}
</style>