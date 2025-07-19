<template>
  <div class="page-container">
    <!-- 搜索表单 -->
    <div class="search-form">
      <el-form :model="searchForm" inline>
        <el-form-item label="关键词">
          <el-input
            v-model="searchForm.keyword"
            placeholder="请输入角色名称或描述"
            clearable
            style="width: 200px"
          />
        </el-form-item>
        <el-form-item>
          <el-button type="primary" @click="handleSearch">
            <el-icon><Search /></el-icon>
            搜索
          </el-button>
          <el-button @click="handleReset">
            <el-icon><Refresh /></el-icon>
            重置
          </el-button>
        </el-form-item>
      </el-form>
    </div>
    
    <!-- 操作按钮 -->
    <div class="action-buttons">
      <el-button type="primary" @click="handleCreate">
        <el-icon><Plus /></el-icon>
        新增角色
      </el-button>
    </div>
    
    <!-- 角色表格 -->
    <div class="content-card">
      <el-table
        v-loading="loading"
        :data="tableData"
        style="width: 100%"
        @selection-change="handleSelectionChange"
      >
        <el-table-column type="selection" width="55" />
        <el-table-column prop="id" label="ID" width="80" />
        <el-table-column prop="name" label="角色名称" min-width="120" />
        <el-table-column prop="description" label="描述" min-width="200" />
        <el-table-column prop="status" label="状态" width="100">
          <template #default="{ row }">
            <el-tag :type="row.status === 1 ? 'success' : 'danger'">
              {{ row.status === 1 ? '启用' : '禁用' }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="created_at" label="创建时间" width="180">
          <template #default="{ row }">
            {{ formatDate(row.created_at) }}
          </template>
        </el-table-column>
        <el-table-column label="操作" width="250" fixed="right">
          <template #default="{ row }">
            <el-button type="primary" size="small" @click="handleEdit(row)">
              编辑
            </el-button>
            <el-button type="info" size="small" @click="handlePermissions(row)">
              权限
            </el-button>
            <el-button
              :type="row.status === 1 ? 'warning' : 'success'"
              size="small"
              @click="handleToggleStatus(row)"
            >
              {{ row.status === 1 ? '禁用' : '启用' }}
            </el-button>
            <el-button type="danger" size="small" @click="handleDelete(row)">
              删除
            </el-button>
          </template>
        </el-table-column>
      </el-table>
      
      <!-- 分页 -->
      <div class="pagination-container">
        <el-pagination
          v-model:current-page="pagination.page"
          v-model:page-size="pagination.limit"
          :total="pagination.total"
          :page-sizes="[10, 20, 50, 100]"
          layout="total, sizes, prev, pager, next, jumper"
          @size-change="handleSizeChange"
          @current-change="handleCurrentChange"
        />
      </div>
    </div>
    
    <!-- 角色表单对话框 -->
    <el-dialog
      v-model="dialogVisible"
      :title="dialogTitle"
      width="600px"
      @close="handleDialogClose"
    >
      <el-form
        ref="formRef"
        :model="formData"
        :rules="formRules"
        label-width="80px"
      >
        <el-form-item label="角色名称" prop="name">
          <el-input v-model="formData.name" placeholder="请输入角色名称" />
        </el-form-item>
        <el-form-item label="描述" prop="description">
          <el-input
            v-model="formData.description"
            type="textarea"
            :rows="3"
            placeholder="请输入角色描述"
          />
        </el-form-item>
        <el-form-item label="状态" prop="status">
          <el-radio-group v-model="formData.status">
            <el-radio :label="1">启用</el-radio>
            <el-radio :label="0">禁用</el-radio>
          </el-radio-group>
        </el-form-item>
      </el-form>
      
      <template #footer>
        <el-button @click="dialogVisible = false">取消</el-button>
        <el-button type="primary" :loading="submitLoading" @click="handleSubmit">
          确定
        </el-button>
      </template>
    </el-dialog>
    
    <!-- 权限分配对话框 -->
    <el-dialog
      v-model="permissionDialogVisible"
      title="分配权限"
      width="800px"
      @close="handlePermissionDialogClose"
    >
      <div v-loading="permissionLoading">
        <el-tree
          ref="permissionTreeRef"
          :data="permissionTreeData"
          :props="treeProps"
          show-checkbox
          node-key="id"
          :default-checked-keys="checkedPermissions"
          :check-strictly="false"
        />
      </div>
      
      <template #footer>
        <el-button @click="permissionDialogVisible = false">取消</el-button>
        <el-button type="primary" :loading="permissionSubmitLoading" @click="handlePermissionSubmit">
          确定
        </el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, computed, onMounted } from 'vue'
import { ElMessage, ElMessageBox, type FormInstance, type FormRules } from 'element-plus'
import { Search, Refresh, Plus } from '@element-plus/icons-vue'
import {
  getRoleList,
  createRole,
  updateRole,
  deleteRole,
  updateRoleStatus,
  getRolePermissions,
  assignRolePermissions
} from '@/api/roles'
import { getPermissionList } from '@/api/permissions'
import type { Role, CreateRoleForm, UpdateRoleForm, Permission } from '@/types/role'

const loading = ref(false)
const submitLoading = ref(false)
const permissionLoading = ref(false)
const permissionSubmitLoading = ref(false)
const dialogVisible = ref(false)
const permissionDialogVisible = ref(false)
const isEdit = ref(false)
const formRef = ref<FormInstance>()
const permissionTreeRef = ref()

const searchForm = reactive({
  keyword: ''
})

const pagination = reactive({
  page: 1,
  limit: 20,
  total: 0
})

const tableData = ref<Role[]>([])
const selectedRows = ref<Role[]>([])
const currentRole = ref<Role | null>(null)
const permissionTreeData = ref<Permission[]>([])
const checkedPermissions = ref<number[]>([])

const formData = reactive<CreateRoleForm & { id?: number }>({
  name: '',
  description: '',
  status: 1
})

const formRules: FormRules = {
  name: [
    { required: true, message: '请输入角色名称', trigger: 'blur' },
    { min: 2, max: 50, message: '角色名称长度在 2 到 50 个字符', trigger: 'blur' }
  ],
  description: [
    { max: 200, message: '描述不能超过 200 个字符', trigger: 'blur' }
  ]
}

const treeProps = {
  children: 'children',
  label: 'name'
}

const dialogTitle = computed(() => isEdit.value ? '编辑角色' : '新增角色')

// 获取角色列表
const fetchRoleList = async () => {
  try {
    loading.value = true
    const params = {
      page: pagination.page,
      limit: pagination.limit,
      keyword: searchForm.keyword || undefined
    }
    
    const response = await getRoleList(params)
    tableData.value = response.data.items
    pagination.total = response.data.total || 0
  } catch (error) {
    console.error('获取角色列表失败:', error)
  } finally {
    loading.value = false
  }
}

// 获取权限列表
const fetchPermissionList = async () => {
  try {
    permissionLoading.value = true
    const response = await getPermissionList({ page: 1, limit: 1000 })
    permissionTreeData.value = buildPermissionTree(response.data.items || [])
  } catch (error) {
    console.error('获取权限列表失败:', error)
  } finally {
    permissionLoading.value = false
  }
}

// 构建权限树
const buildPermissionTree = (permissions: Permission[]): Permission[] => {
  const tree: Permission[] = []
  const map = new Map<number, Permission>()
  
  // 创建映射
  permissions.forEach(permission => {
    map.set(permission.id, { ...permission, children: [] })
  })
  
  // 构建树结构
  permissions.forEach(permission => {
    const node = map.get(permission.id)!
    if (permission.parent_id === 0) {
      tree.push(node)
    } else {
      const parent = map.get(permission.parent_id)
      if (parent) {
        parent.children = parent.children || []
        parent.children.push(node)
      }
    }
  })
  
  return tree
}

// 搜索
const handleSearch = () => {
  pagination.page = 1
  fetchRoleList()
}

// 重置
const handleReset = () => {
  searchForm.keyword = ''
  pagination.page = 1
  fetchRoleList()
}

// 新增角色
const handleCreate = () => {
  isEdit.value = false
  dialogVisible.value = true
  resetForm()
}

// 编辑角色
const handleEdit = (row: Role) => {
  isEdit.value = true
  dialogVisible.value = true
  Object.assign(formData, {
    id: row.id,
    name: row.name,
    description: row.description,
    status: row.status
  })
}

// 删除角色
const handleDelete = async (row: Role) => {
  try {
    await ElMessageBox.confirm(
      `确定要删除角色 "${row.name}" 吗？`,
      '提示',
      {
        confirmButtonText: '确定',
        cancelButtonText: '取消',
        type: 'warning'
      }
    )
    
    await deleteRole(row.id)
    ElMessage.success('删除成功')
    fetchRoleList()
  } catch (error: any) {
    if (error !== 'cancel') {
      console.error('删除角色失败:', error)
    }
  }
}

// 切换角色状态
const handleToggleStatus = async (row: Role) => {
  try {
    const newStatus = row.status === 1 ? 0 : 1
    await updateRoleStatus(row.id, newStatus)
    ElMessage.success('状态更新成功')
    fetchRoleList()
  } catch (error) {
    console.error('更新角色状态失败:', error)
  }
}

// 权限管理
const handlePermissions = async (row: Role) => {
  currentRole.value = row
  permissionDialogVisible.value = true
  
  // 获取权限列表
  await fetchPermissionList()
  
  // 获取角色已有权限
  try {
    const response = await getRolePermissions(row.id)
    checkedPermissions.value = response.data.map((p: Permission) => p.id)
  } catch (error) {
    console.error('获取角色权限失败:', error)
    checkedPermissions.value = []
  }
}

// 表单提交
const handleSubmit = async () => {
  if (!formRef.value) return
  
  try {
    await formRef.value.validate()
    
    submitLoading.value = true
    
    if (isEdit.value) {
      const updateData: UpdateRoleForm = {
        name: formData.name,
        description: formData.description,
        status: formData.status
      }
      await updateRole(formData.id!, updateData)
      ElMessage.success('更新成功')
    } else {
      await createRole(formData)
      ElMessage.success('创建成功')
    }
    
    dialogVisible.value = false
    fetchRoleList()
  } catch (error) {
    console.error('提交失败:', error)
  } finally {
    submitLoading.value = false
  }
}

// 权限提交
const handlePermissionSubmit = async () => {
  if (!currentRole.value) return
  
  try {
    permissionSubmitLoading.value = true
    
    const checkedKeys = permissionTreeRef.value.getCheckedKeys()
    const halfCheckedKeys = permissionTreeRef.value.getHalfCheckedKeys()
    const allCheckedKeys = [...checkedKeys, ...halfCheckedKeys]
    
    await assignRolePermissions(currentRole.value.id, allCheckedKeys)
    ElMessage.success('权限分配成功')
    permissionDialogVisible.value = false
  } catch (error) {
    console.error('权限分配失败:', error)
  } finally {
    permissionSubmitLoading.value = false
  }
}

// 对话框关闭
const handleDialogClose = () => {
  formRef.value?.resetFields()
  resetForm()
}

// 权限对话框关闭
const handlePermissionDialogClose = () => {
  currentRole.value = null
  checkedPermissions.value = []
  permissionTreeData.value = []
}

// 重置表单
const resetForm = () => {
  Object.assign(formData, {
    id: undefined,
    name: '',
    description: '',
    status: 1
  })
}

// 选择变化
const handleSelectionChange = (selection: Role[]) => {
  selectedRows.value = selection
}

// 分页大小变化
const handleSizeChange = (size: number) => {
  pagination.limit = size
  pagination.page = 1
  fetchRoleList()
}

// 当前页变化
const handleCurrentChange = (page: number) => {
  pagination.page = page
  fetchRoleList()
}

// 格式化日期
const formatDate = (date: string) => {
  return new Date(date).toLocaleString()
}

onMounted(() => {
  fetchRoleList()
})
</script>

<style scoped>
.pagination-container {
  display: flex;
  justify-content: center;
  margin-top: 20px;
}
</style>