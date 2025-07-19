<template>
  <div class="error-page">
    <div class="error-container">
      <div class="error-content">
        <div class="error-code">404</div>
        <div class="error-title">页面不存在</div>
        <div class="error-description">
          抱歉，您访问的页面不存在或已被删除
        </div>
        <div class="error-actions">
          <el-button type="primary" @click="goHome">
            <el-icon><HomeFilled /></el-icon>
            返回首页
          </el-button>
          <el-button @click="goBack">
            <el-icon><ArrowLeft /></el-icon>
            返回上页
          </el-button>
        </div>
      </div>
      <div class="error-image">
        <svg viewBox="0 0 404 300" xmlns="http://www.w3.org/2000/svg">
          <!-- 背景 -->
          <rect width="404" height="300" fill="#f8f9fa"/>
          
          <!-- 云朵 -->
          <ellipse cx="100" cy="80" rx="25" ry="15" fill="#e9ecef"/>
          <ellipse cx="120" cy="75" rx="30" ry="18" fill="#e9ecef"/>
          <ellipse cx="140" cy="80" rx="25" ry="15" fill="#e9ecef"/>
          
          <ellipse cx="300" cy="60" rx="20" ry="12" fill="#e9ecef"/>
          <ellipse cx="315" cy="55" rx="25" ry="15" fill="#e9ecef"/>
          <ellipse cx="330" cy="60" rx="20" ry="12" fill="#e9ecef"/>
          
          <!-- 山峰 -->
          <polygon points="50,250 150,150 250,250" fill="#dee2e6"/>
          <polygon points="200,250 300,120 400,250" fill="#ced4da"/>
          
          <!-- 太阳 -->
          <circle cx="350" cy="90" r="25" fill="#ffd43b"/>
          
          <!-- 404文字 -->
          <text x="202" y="200" font-family="Arial, sans-serif" font-size="48" font-weight="bold" text-anchor="middle" fill="#6c757d">404</text>
          
          <!-- 小鸟 -->
          <path d="M80 120 Q85 115 90 120 Q85 125 80 120" fill="#495057"/>
          <path d="M85 118 Q90 113 95 118" fill="none" stroke="#495057" stroke-width="1"/>
          
          <path d="M320 100 Q325 95 330 100 Q325 105 320 100" fill="#495057"/>
          <path d="M325 98 Q330 93 335 98" fill="none" stroke="#495057" stroke-width="1"/>
        </svg>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { useRouter } from 'vue-router'
import { HomeFilled, ArrowLeft } from '@element-plus/icons-vue'

const router = useRouter()

const goHome = () => {
  router.push('/')
}

const goBack = () => {
  router.go(-1)
}
</script>

<style scoped>
.error-page {
  min-height: 100vh;
  display: flex;
  align-items: center;
  justify-content: center;
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  padding: 20px;
}

.error-container {
  display: flex;
  align-items: center;
  gap: 60px;
  max-width: 1000px;
  width: 100%;
}

.error-content {
  flex: 1;
  text-align: center;
  color: white;
}

.error-code {
  font-size: 120px;
  font-weight: bold;
  line-height: 1;
  margin-bottom: 20px;
  background: linear-gradient(45deg, #ff6b6b, #feca57);
  -webkit-background-clip: text;
  -webkit-text-fill-color: transparent;
  background-clip: text;
}

.error-title {
  font-size: 32px;
  font-weight: 600;
  margin-bottom: 16px;
}

.error-description {
  font-size: 18px;
  margin-bottom: 40px;
  opacity: 0.9;
  line-height: 1.6;
}

.error-actions {
  display: flex;
  gap: 16px;
  justify-content: center;
  flex-wrap: wrap;
}

.error-actions .el-button {
  padding: 12px 24px;
  font-size: 16px;
}

.error-image {
  flex: 1;
  max-width: 400px;
}

.error-image svg {
  width: 100%;
  height: auto;
  border-radius: 12px;
  box-shadow: 0 10px 30px rgba(0, 0, 0, 0.2);
}

@media (max-width: 768px) {
  .error-container {
    flex-direction: column;
    gap: 40px;
    text-align: center;
  }
  
  .error-code {
    font-size: 80px;
  }
  
  .error-title {
    font-size: 24px;
  }
  
  .error-description {
    font-size: 16px;
  }
  
  .error-actions {
    flex-direction: column;
    align-items: center;
  }
  
  .error-actions .el-button {
    width: 200px;
  }
}
</style>