<template>
  <div class="page-container">
    <div class="profile-header">
      <h1 class="page-title">个人中心</h1>
      <p class="page-description">管理您的个人信息和账户设置</p>
    </div>
    
    <el-row :gutter="20">
      <el-col :span="8">
        <!-- 用户信息卡片 -->
        <el-card class="profile-card">
          <template #header>
            <div class="card-header">
              <span>基本信息</span>
            </div>
          </template>
          
          <div class="profile-info">
            <div class="avatar-section">
              <el-avatar :size="80" class="profile-avatar">
                {{ userStore.userInfo?.nickname?.charAt(0) || 'U' }}
              </el-avatar>
              <div class="user-details">
                <h3>{{ userStore.userInfo?.nickname || '未设置' }}</h3>
                <p class="user-role">{{ userStore.userInfo?.username || '未知用户' }}</p>
              </div>
            </div>
            
            <el-divider />
            
            <div class="info-item">
              <span class="label">用户名：</span>
              <span class="value">{{ userStore.userInfo?.username || '未设置' }}</span>
            </div>
            
            <div class="info-item">
              <span class="label">昵称：</span>
              <span class="value">{{ userStore.userInfo?.nickname || '未设置' }}</span>
            </div>
            
            <div class="info-item">
              <span class="label">邮箱：</span>
              <span class="value">{{ userStore.userInfo?.email || '未设置' }}</span>
            </div>
            
            <div class="info-item">
              <span class="label">状态：</span>
              <el-tag :type="userStore.userInfo?.status === 1 ? 'success' : 'danger'" size="small">
                {{ userStore.userInfo?.status === 1 ? '正常' : '禁用' }}
              </el-tag>
            </div>
            
            <div class="info-item">
              <span class="label">创建时间：</span>
              <span class="value">{{ formatDate(userStore.userInfo?.createdAt) }}</span>
            </div>
          </div>
        </el-card>
      </el-col>
      
      <el-col :span="16">
        <!-- 编辑信息表单 -->
        <el-card class="edit-card">
          <template #header>
            <div class="card-header">
              <span>编辑信息</span>
            </div>
          </template>
          
          <el-form
            ref="formRef"
            :model="form"
            :rules="rules"
            label-width="100px"
            class="profile-form"
          >
            <el-form-item label="昵称" prop="nickname">
              <el-input v-model="form.nickname" placeholder="请输入昵称" />
            </el-form-item>
            
            <el-form-item label="邮箱" prop="email">
              <el-input v-model="form.email" placeholder="请输入邮箱" type="email" />
            </el-form-item>
            
            <el-form-item>
              <el-button type="primary" @click="updateProfile" :loading="loading">
                保存修改
              </el-button>
              <el-button @click="resetForm">重置</el-button>
            </el-form-item>
          </el-form>
        </el-card>
        
        <!-- 修改密码 -->
        <el-card class="password-card" style="margin-top: 20px;">
          <template #header>
            <div class="card-header">
              <span>修改密码</span>
            </div>
          </template>
          
          <el-form
            ref="passwordFormRef"
            :model="passwordForm"
            :rules="passwordRules"
            label-width="100px"
            class="password-form"
          >
            <el-form-item label="当前密码" prop="oldPassword">
              <el-input
                v-model="passwordForm.oldPassword"
                type="password"
                placeholder="请输入当前密码"
                show-password
              />
            </el-form-item>
            
            <el-form-item label="新密码" prop="newPassword">
              <el-input
                v-model="passwordForm.newPassword"
                type="password"
                placeholder="请输入新密码"
                show-password
              />
            </el-form-item>
            
            <el-form-item label="确认密码" prop="confirmPassword">
              <el-input
                v-model="passwordForm.confirmPassword"
                type="password"
                placeholder="请再次输入新密码"
                show-password
              />
            </el-form-item>
            
            <el-form-item>
              <el-button type="primary" @click="changePassword" :loading="passwordLoading">
                修改密码
              </el-button>
              <el-button @click="resetPasswordForm">重置</el-button>
            </el-form-item>
          </el-form>
        </el-card>
      </el-col>
    </el-row>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, onMounted } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import type { FormInstance, FormRules } from 'element-plus'
import { useUserStore } from '@/stores/user'
import { updateUserProfile as updateProfileApi } from '@/api/auth'

const userStore = useUserStore()
const formRef = ref<FormInstance>()
const passwordFormRef = ref<FormInstance>()
const loading = ref(false)
const passwordLoading = ref(false)

// 基本信息表单
const form = reactive({
  nickname: '',
  email: ''
})

// 密码表单
const passwordForm = reactive({
  oldPassword: '',
  newPassword: '',
  confirmPassword: ''
})

// 表单验证规则
const rules: FormRules = {
  nickname: [
    { required: true, message: '请输入昵称', trigger: 'blur' },
    { min: 2, max: 20, message: '昵称长度在 2 到 20 个字符', trigger: 'blur' }
  ],
  email: [
    { type: 'email', message: '请输入正确的邮箱地址', trigger: ['blur', 'change'] }
  ]
}

const passwordRules: FormRules = {
  oldPassword: [
    { required: true, message: '请输入当前密码', trigger: 'blur' }
  ],
  newPassword: [
    { required: true, message: '请输入新密码', trigger: 'blur' },
    { min: 6, max: 20, message: '密码长度在 6 到 20 个字符', trigger: 'blur' }
  ],
  confirmPassword: [
    { required: true, message: '请再次输入新密码', trigger: 'blur' },
    {
      validator: (rule, value, callback) => {
        if (value !== passwordForm.newPassword) {
          callback(new Error('两次输入密码不一致'))
        } else {
          callback()
        }
      },
      trigger: 'blur'
    }
  ]
}

// 格式化日期
const formatDate = (dateString?: string) => {
  if (!dateString) return '未知'
  return new Date(dateString).toLocaleString('zh-CN')
}

// 初始化表单数据
const initForm = () => {
  if (userStore.userInfo) {
    form.nickname = userStore.userInfo.nickname || ''
    form.email = userStore.userInfo.email || ''
  }
}

// 更新个人信息
const updateProfile = async () => {
  if (!formRef.value) return
  
  try {
    await formRef.value.validate()
    loading.value = true
    
    await updateProfileApi(form)
    await userStore.getUserInfo() // 重新获取用户信息
    
    ElMessage.success('个人信息更新成功')
  } catch (error) {
    console.error('更新个人信息失败:', error)
    ElMessage.error('更新个人信息失败')
  } finally {
    loading.value = false
  }
}

// 修改密码
const changePassword = async () => {
  if (!passwordFormRef.value) return
  
  try {
    await passwordFormRef.value.validate()
    passwordLoading.value = true
    
    // TODO: 实现修改密码API
    await new Promise(resolve => setTimeout(resolve, 1000)) // 模拟API调用
    
    ElMessage.success('密码修改成功')
    resetPasswordForm()
  } catch (error) {
    console.error('修改密码失败:', error)
    ElMessage.error('修改密码失败')
  } finally {
    passwordLoading.value = false
  }
}

// 重置基本信息表单
const resetForm = () => {
  initForm()
}

// 重置密码表单
const resetPasswordForm = () => {
  passwordForm.oldPassword = ''
  passwordForm.newPassword = ''
  passwordForm.confirmPassword = ''
  passwordFormRef.value?.clearValidate()
}

onMounted(() => {
  initForm()
})
</script>

<style scoped>
.page-container {
  padding: 24px;
  background-color: #f0f2f5;
  min-height: calc(100vh - 60px);
}

.profile-header {
  margin-bottom: 24px;
}

.page-title {
  font-size: 24px;
  font-weight: 600;
  color: #1f2937;
  margin: 0 0 8px 0;
}

.page-description {
  color: #6b7280;
  margin: 0;
  font-size: 14px;
}

.profile-card,
.edit-card,
.password-card {
  border-radius: 8px;
  box-shadow: 0 1px 3px rgba(0, 0, 0, 0.1);
}

.card-header {
  font-weight: 600;
  color: #1f2937;
}

.profile-info {
  padding: 0;
}

.avatar-section {
  display: flex;
  align-items: center;
  margin-bottom: 20px;
}

.profile-avatar {
  background-color: #409eff;
  margin-right: 16px;
}

.user-details h3 {
  margin: 0 0 4px 0;
  font-size: 18px;
  font-weight: 600;
  color: #1f2937;
}

.user-role {
  margin: 0;
  color: #6b7280;
  font-size: 14px;
}

.info-item {
  display: flex;
  align-items: center;
  margin-bottom: 12px;
  padding: 8px 0;
}

.info-item:last-child {
  margin-bottom: 0;
}

.label {
  font-weight: 500;
  color: #374151;
  min-width: 80px;
}

.value {
  color: #6b7280;
  flex: 1;
}

.profile-form,
.password-form {
  max-width: 500px;
}

:deep(.el-form-item__label) {
  font-weight: 500;
  color: #374151;
}

:deep(.el-card__header) {
  background-color: #f9fafb;
  border-bottom: 1px solid #e5e7eb;
}
</style>