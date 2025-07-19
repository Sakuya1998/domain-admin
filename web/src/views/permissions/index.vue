<template>
  <div class="page-container">
    <!-- 搜索表单 -->
    <div class="search-form">
      <el-form :model="searchForm" inline>
        <el-form-item label="关键词">
          <el-input
            v-model="searchForm.keyword"
            placeholder="请输入权限名称或路径"
            clearable
            style="width: 200px"
          />
        </el-form-item>
        <el-form-item label="类型">
          <el-select v-model="searchForm.type" placeholder="请选择类型" clearable style="width: 120px">
            <el-option label="菜单" value="menu" />
            <el-option label="按钮" value="button" />
            <el-option label="API" value="api" />
          </el-select>
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
        新增权限
      </el-button>
      <el-button @click="handleExpandAll">
        <el-icon><ArrowDown /></el-icon>
        展开全部
      </el-button>
      <el-button @click="handleCollapseAll">
        <el-icon><ArrowUp /></el-icon>
        收起全部
      </el-button>
    </div>
    
    <!-- 权限表格 -->
    <div class="content-card">
      <!-- 调试信息 -->
      <div v-if="tableData.length === 0" style="padding: 20px; text-align: center; color: #999;">
        <p>表格数据为空，原始权限数据数量: {{ allPermissions.length }}</p>
        <p v-if="allPermissions.length > 0">有原始数据但表格为空，可能是树结构构建问题</p>
        <el-button @click="forceRefresh" type="primary">强制刷新</el-button>
      </div>
      <el-table
        ref="tableRef"
        v-loading="loading"
        :data="tableData"
        style="width: 100%"
        row-key="id"
        :tree-props="{ children: 'children', hasChildren: 'hasChildren' }"
        :default-expand-all="false"
      >
        <el-table-column prop="display_name" label="权限名称" min-width="200">
          <template #default="{ row }">
            <el-icon v-if="row.type === 'menu'" style="margin-right: 5px"><Menu /></el-icon>
            <el-icon v-else-if="row.type === 'button'" style="margin-right: 5px"><Operation /></el-icon>
            <el-icon v-else style="margin-right: 5px"><Link /></el-icon>
            {{ row.display_name }}
          </template>
        </el-table-column>
        <el-table-column prop="resource" label="路径/标识" min-width="150" />
        <el-table-column prop="type" label="类型" width="100">
          <template #default="{ row }">
            <el-tag :type="getTypeTagType(row.type)">{{ getTypeText(row.type) }}</el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="action" label="请求方式" width="100">
          <template #default="{ row }">
            <el-tag v-if="row.action" :type="getMethodTagType(row.action)" size="small">
              {{ row.action }}
            </el-tag>
            <span v-else>-</span>
          </template>
        </el-table-column>
        <el-table-column prop="id" label="排序" width="80" />
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
        <el-table-column label="操作" width="200" fixed="right">
          <template #default="{ row }">
            <el-button type="primary" size="small" @click="handleEdit(row)">
              编辑
            </el-button>
            <el-button type="success" size="small" @click="handleCreateChild(row)">
              新增子权限
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
    </div>
    
    <!-- 权限表单对话框 -->
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
        label-width="100px"
      >
        <el-form-item label="上级权限" prop="parent_id">
          <el-tree-select
            v-model="formData.parent_id"
            :data="parentOptions"
            :props="{ value: 'id', label: 'name', children: 'children' }"
            placeholder="请选择上级权限（不选择则为顶级权限）"
            clearable
            check-strictly
            style="width: 100%"
          />
        </el-form-item>
        <el-form-item label="权限标识" prop="name">
          <el-input v-model="formData.name" placeholder="请输入权限标识（英文）" />
        </el-form-item>
        <el-form-item label="权限名称" prop="display_name">
          <el-input v-model="formData.display_name" placeholder="请输入权限显示名称" />
        </el-form-item>
        <el-form-item label="权限类型" prop="type">
          <el-radio-group v-model="formData.type">
            <el-radio label="menu">菜单</el-radio>
            <el-radio label="button">按钮</el-radio>
            <el-radio label="api">API</el-radio>
          </el-radio-group>
        </el-form-item>
        <el-form-item label="路径/标识" prop="resource">
          <el-input v-model="formData.resource" placeholder="请输入路径或标识" />
        </el-form-item>
        <el-form-item v-if="formData.type === 'api'" label="请求方式" prop="action">
          <el-select v-model="formData.action" placeholder="请选择请求方式">
            <el-option label="GET" value="GET" />
            <el-option label="POST" value="POST" />
            <el-option label="PUT" value="PUT" />
            <el-option label="DELETE" value="DELETE" />
            <el-option label="PATCH" value="PATCH" />
          </el-select>
        </el-form-item>
        <el-form-item label="描述" prop="description">
          <el-input
            v-model="formData.description"
            type="textarea"
            :rows="3"
            placeholder="请输入权限描述"
          />
        </el-form-item>
        <el-form-item label="排序" prop="sort">
          <el-input-number v-model="formData.sort" :min="0" :max="9999" />
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
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, computed, onMounted, watch } from 'vue'
import { ElMessage, ElMessageBox, type FormInstance, type FormRules } from 'element-plus'
import {
  Search,
  Refresh,
  Plus,
  ArrowDown,
  ArrowUp,
  Menu,
  Operation,
  Link
} from '@element-plus/icons-vue'
import {
  getPermissionList,
  createPermission,
  updatePermission,
  deletePermission,
  updatePermissionStatus
} from '@/api/permissions'
import type { Permission, CreatePermissionForm, UpdatePermissionForm } from '@/types/permission'

const loading = ref(false)
const submitLoading = ref(false)
const dialogVisible = ref(false)
const isEdit = ref(false)
const formRef = ref<FormInstance>()
const tableRef = ref()

const searchForm = reactive({
  keyword: '',
  type: ''
})

const tableData = ref<Permission[]>([])
const allPermissions = ref<Permission[]>([])
const parentOptions = ref<Permission[]>([])

const formData = reactive<CreatePermissionForm & { id?: number }>({
  parent_id: 0,
  name: '',
  display_name: '',
  resource: '',
  type: 'menu',
  action: '',
  description: '',
  sort: 0,
  status: 1
})

const formRules: FormRules = {
  name: [
    { required: true, message: '请输入权限标识', trigger: 'blur' },
    { min: 2, max: 50, message: '权限标识长度在 2 到 50 个字符', trigger: 'blur' }
  ],
  display_name: [
    { required: true, message: '请输入权限名称', trigger: 'blur' },
    { min: 2, max: 50, message: '权限名称长度在 2 到 50 个字符', trigger: 'blur' }
  ],
  resource: [
    { required: true, message: '请输入路径或标识', trigger: 'blur' }
  ],
  type: [
    { required: true, message: '请选择权限类型', trigger: 'change' }
  ],
  action: [
     {
       validator: (rule, value, callback) => {
         if (formData.type === 'api' && !value) {
           callback(new Error('请选择请求方式'))
         } else {
           callback()
         }
       },
       trigger: 'change'
     }
   ]
}

const dialogTitle = computed(() => isEdit.value ? '编辑权限' : '新增权限')

// 监听权限类型变化
watch(() => formData.type, (newType) => {
  if (newType !== 'api') {
    formData.action = ''
  }
})

// 获取权限列表
const fetchPermissionList = async () => {
  try {
    loading.value = true
    console.log('开始获取权限列表...')
    const response = await getPermissionList({ page: 1, limit: 1000 })
    console.log('权限列表响应:', response)
    console.log('响应类型:', typeof response)
    console.log('响应数据结构:', JSON.stringify(response, null, 2))
    
    // 处理后端返回的数据结构 {code: 200, data: {items: [...], total: ...}}
    if (response && response.code === 200 && response.data) {
      console.log('权限数据:', response.data)
      console.log('权限项目:', response.data.items)
      allPermissions.value = response.data.items || []
      console.log('设置后的allPermissions:', allPermissions.value)
      console.log('权限数量:', allPermissions.value.length)
    } else {
      allPermissions.value = []
      console.warn('权限列表数据格式异常:', response)
    }
    filterAndBuildTree()
    console.log('过滤后的tableData:', tableData.value)
  } catch (error) {
    console.error('获取权限列表失败:', error)
    allPermissions.value = []
  } finally {
    loading.value = false
  }
}

// 过滤并构建树
const filterAndBuildTree = () => {
  console.log('开始过滤和构建树，原始数据:', allPermissions.value)
  let filteredPermissions = [...allPermissions.value]
  
  // 关键词过滤
  if (searchForm.keyword) {
    filteredPermissions = filteredPermissions.filter(permission =>
      permission.display_name.includes(searchForm.keyword) ||
      permission.resource.includes(searchForm.keyword)
    )
  }
  
  // 类型过滤
  if (searchForm.type) {
    filteredPermissions = filteredPermissions.filter(permission =>
      permission.type === searchForm.type
    )
  }
  
  console.log('过滤后的权限:', filteredPermissions)
  tableData.value = buildPermissionTree(filteredPermissions)
  console.log('构建树后的tableData:', tableData.value)
}

// 构建权限树
const buildPermissionTree = (permissions: Permission[]): Permission[] => {
  console.log('构建权限树，输入数据:', permissions)
  
  if (permissions.length === 0) {
    console.log('权限数据为空，返回空数组')
    return []
  }
  
  const tree: Permission[] = []
  const map = new Map<number, Permission>()
  
  // 创建映射
  permissions.forEach(permission => {
    map.set(permission.id, { ...permission, children: [] })
  })
  
  console.log('权限映射:', map)
  
  // 构建树结构
  permissions.forEach(permission => {
    const node = map.get(permission.id)!
    if (permission.parent_id === 0) {
      tree.push(node)
      console.log('添加根节点:', node.name)
    } else {
      const parent = map.get(permission.parent_id)
      if (parent) {
        parent.children = parent.children || []
        parent.children.push(node)
        console.log('添加子节点:', node.name, '到父节点:', parent.name)
      } else {
        // 如果父节点不存在，将其作为根节点处理
        console.log('父节点不存在，将', node.name, '作为根节点处理')
        tree.push(node)
      }
    }
  })
  
  console.log('构建完成的树结构:', tree)
  return tree
}

// 构建父级选项
const buildParentOptions = (excludeId?: number) => {
  const options = allPermissions.value
    .filter(p => p.id !== excludeId)
    .map(p => ({ ...p, children: [] }))
  
  parentOptions.value = buildPermissionTree(options)
}

// 搜索
const handleSearch = () => {
  filterAndBuildTree()
}

// 重置
const handleReset = () => {
  searchForm.keyword = ''
  searchForm.type = ''
  filterAndBuildTree()
}

// 展开全部
const handleExpandAll = () => {
  if (tableRef.value) {
    tableRef.value.store.states.defaultExpandAll.value = true
    tableRef.value.store.updateTableScrollY()
  }
}

// 收起全部
const handleCollapseAll = () => {
  if (tableRef.value) {
    tableRef.value.store.states.defaultExpandAll.value = false
    tableRef.value.store.updateTableScrollY()
  }
}

// 新增权限
const handleCreate = () => {
  isEdit.value = false
  dialogVisible.value = true
  resetForm()
  buildParentOptions()
}

// 新增子权限
const handleCreateChild = (row: Permission) => {
  isEdit.value = false
  dialogVisible.value = true
  resetForm()
  formData.parent_id = row.id
  buildParentOptions()
}

// 编辑权限
const handleEdit = (row: Permission) => {
  isEdit.value = true
  dialogVisible.value = true
  Object.assign(formData, {
    id: row.id,
    parent_id: row.parent_id,
    name: row.name,
    display_name: row.display_name,
    resource: row.resource,
    type: row.type,
    action: row.action || '',
    description: row.description,
    sort: row.sort,
    status: row.status
  })
  buildParentOptions(row.id)
}

// 删除权限
const handleDelete = async (row: Permission) => {
  try {
    await ElMessageBox.confirm(
      `确定要删除权限 "${row.name}" 吗？删除后其子权限也会被删除！`,
      '提示',
      {
        confirmButtonText: '确定',
        cancelButtonText: '取消',
        type: 'warning'
      }
    )
    
    await deletePermission(row.id)
    ElMessage.success('删除成功')
    fetchPermissionList()
  } catch (error: any) {
    if (error !== 'cancel') {
      console.error('删除权限失败:', error)
    }
  }
}

// 切换权限状态
const handleToggleStatus = async (row: Permission) => {
  try {
    const newStatus = row.status === 1 ? 0 : 1
    await updatePermissionStatus(row.id, newStatus)
    ElMessage.success('状态更新成功')
    fetchPermissionList()
  } catch (error) {
    console.error('更新权限状态失败:', error)
  }
}

// 表单提交
const handleSubmit = async () => {
  if (!formRef.value) return
  
  try {
    await formRef.value.validate()
    
    submitLoading.value = true
    
    if (isEdit.value) {
      const updateData: UpdatePermissionForm = {
        parent_id: formData.parent_id,
        name: formData.name,
        display_name: formData.display_name,
        resource: formData.resource,
        type: formData.type,
        action: formData.action || undefined,
        description: formData.description,
        sort: formData.sort,
        status: formData.status
      }
      await updatePermission(formData.id!, updateData)
      ElMessage.success('更新成功')
    } else {
      await createPermission(formData)
      ElMessage.success('创建成功')
    }
    
    dialogVisible.value = false
    fetchPermissionList()
  } catch (error) {
    console.error('提交失败:', error)
  } finally {
    submitLoading.value = false
  }
}

// 对话框关闭
const handleDialogClose = () => {
  formRef.value?.resetFields()
  resetForm()
}

// 重置表单
const resetForm = () => {
  Object.assign(formData, {
    id: undefined,
    parent_id: 0,
    name: '',
    display_name: '',
    resource: '',
    type: 'menu',
    action: '',
    description: '',
    sort: 0,
    status: 1
  })
}

// 格式化日期
const formatDate = (date: string) => {
  return new Date(date).toLocaleString()
}

// 获取类型标签类型
const getTypeTagType = (type: string) => {
  const typeMap: Record<string, string> = {
    menu: 'primary',
    button: 'success',
    api: 'warning'
  }
  return typeMap[type] || 'info'
}

// 获取类型文本
const getTypeText = (type: string) => {
  const textMap: Record<string, string> = {
    menu: '菜单',
    button: '按钮',
    api: 'API'
  }
  return textMap[type] || type
}

// 获取请求方式标签类型
const getMethodTagType = (method: string) => {
  const methodMap: Record<string, string> = {
    GET: 'success',
    POST: 'primary',
    PUT: 'warning',
    DELETE: 'danger',
    PATCH: 'info'
  }
  return methodMap[method] || 'info'
}

// 强制刷新
const forceRefresh = () => {
  console.log('强制刷新，清空数据重新获取')
  allPermissions.value = []
  tableData.value = []
  fetchPermissionList()
}

onMounted(() => {
  fetchPermissionList()
})
</script>

<style scoped>
.pagination-container {
  display: flex;
  justify-content: center;
  margin-top: 20px;
}
</style>