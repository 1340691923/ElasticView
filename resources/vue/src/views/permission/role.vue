<template>
  <div class="app-container">

    <div class="search-container" :class="{ 'is-collapsed': isCollapsed && isMobile }">
      <div class="search-header" v-if="isMobile" @click="toggleCollapse">
        <span>搜索条件</span>
        <el-icon :class="{ 'is-collapsed': isCollapsed }">
          <ArrowDown />
        </el-icon>
      </div>
      <el-form :inline="true">
        <el-form-item label="角色名:" prop="keywords">
          <el-input clearable v-model="searchCfg.role_name" style="width:160px"></el-input>
        </el-form-item>
        <el-form-item class="button-group">
          <el-button
            type="success"
            class="filter-item"
            @click="GetRoles(1)"
          >
            {{ $t('搜索') }}
          </el-button>
          <el-button
            type="warning"
            class="filter-item"
            @click="handleAddRole"
          >
            {{ $t('新增') }}
          </el-button>
        </el-form-item>
      </el-form>
    </div>
    <el-card shadow="never" class="table-container">
      <el-table
      v-loading="_this.loading"
      :data="_this.rolesList"
      @row-dblclick="handleEdit"
    >
      <el-table-column align="center" :label="$t('id')" width="40">
        <template #default="scope">
          {{ scope.row.id }}
        </template>
      </el-table-column>
      <el-table-column align="center" :label="$t('角色名')" width="200">
        <template #default="scope">
          {{ scope.row.name }}
        </template>
      </el-table-column>
      <el-table-column align="header-center" :label="$t('角色详细信息')">
        <template #default="scope">
          {{ scope.row.description }}
        </template>
      </el-table-column>
      <el-table-column align="center" :label="$t('操作')" width="140" fixed="right">
        <template #default="scope">
          <el-button
            type="primary"
            @click.stop="handleEdit(scope.row)"
            :icon="Edit"
          >
          </el-button>
          <el-button
            type="danger"
            @click.stop="handleDelete(scope)"
            :icon="Delete"
          >
          </el-button>
        </template>
      </el-table-column>
    </el-table>
    </el-card>
    <div class="pagination-container">
      <el-pagination
        background
        :current-page="input.page"
        :page-size="input.limit"
        :total="count"
        @current-change="GetRoles"
        @size-change="handleSizeChange"
      />
    </div>

    <el-drawer
      :size="isMobile?'100%':'70%'"
      v-model="_this.dialogVisible"
      :title="_this.dialogType==='edit'?$t('修改角色'): $t('新建角色')"
    >
      <el-tabs v-model="activeTab">
        <!-- 基本信息 Tab -->
        <el-tab-pane :label="$t('基本信息')" name="basic">
          <el-form :model="_this.role" label-width="80px" label-position="left">
            <el-form-item :label="$t('角色名')">
              <el-input v-model="_this.role.name" :placeholder="$t('角色名')" />
            </el-form-item>
            <el-form-item :label="$t('角色详情信')">
              <el-input
                v-model="_this.role.description"
                :autosize="{ minRows: 2, maxRows: 4}"
                type="textarea"
                :placeholder="$t('角色详情信息')"
              />
            </el-form-item>
          </el-form>
        </el-tab-pane>

        <!-- 菜单权限 Tab -->
        <el-tab-pane 
          :label="$t('菜单权限')" 
          name="menu"
          v-if="_this.role.id != 1"
        >
          <div class="permission-header">
            <el-input 
              v-model="_this.filterText" 
              :placeholder="$t('输入关键字进行过滤')" 
              class="filter-input"
            />
            <el-button @click="selectAll">{{ $t('全选') }}</el-button>
          </div>
          <el-tree
            default-expand-all
            check-on-click-node
            ref="menu_tree"
            :filter-node-method="filterNode"
            :check-strictly="_this.checkStrictly"
            :data="_this.routes"
            :props="_this.defaultProps"
            show-checkbox
            node-key="path"
            class="permission-tree"
          />
        </el-tab-pane>

        <!-- 接口权限 Tab -->
        <el-tab-pane 
          :label="$t('接口权限')" 
          name="api"
          v-if="_this.role.id != 1"
        >
          <div class="permission-header">
            <el-input 
              v-model="_this.filterText4Api" 
              :placeholder="$t('输入关键字进行过滤')" 
              class="filter-input"
            />
            <el-button @click="selectAll4Api">{{ $t('全选') }}</el-button>
          </div>
          <el-tree
            default-expand-all
            check-on-click-node
            ref="api_tree"
            :filter-node-method="filterNode4Api"
            :check-strictly="_this.checkStrictly"
            :data="_this.apiCfg"
            :props="_this.defaultProps4Api"
            show-checkbox
            node-key="value"
            class="permission-tree"
          />
        </el-tab-pane>
      </el-tabs>

      <template #footer>
        <div class="drawer-footer">
          <el-button type="danger" @click="_this.dialogVisible=false">
            {{ $t('取消') }}
          </el-button>
          <el-button type="primary" @click="confirmRole">
            {{ $t('确定') }}
          </el-button>
        </div>
      </template>
    </el-drawer>

  </div>
</template>

<script lang="ts" setup>
import path from 'path-browserify'

import { deepClone } from '@/utils'
import { asyncRoutes } from '@/utils/asyncRoutes'
import { addRole, deleteRole, getRoles, updateRole } from '@/api/role'
import { UrlConfig } from '@/api/api-rbac'
import {GetRoutesConfig} from "@/api/user";
import {Plus,Search,Edit,Delete,ArrowDown} from "@element-plus/icons-vue";
import {useAppStore} from "@/store";
import {DeviceEnum} from "@/enums/DeviceEnum";

const appStore = useAppStore()

const isMobile = computed(() => appStore.device === DeviceEnum.MOBILE);

const count = ref(0)

const handleSizeChange = (v) => {
  input.limit = v
  GetRoles(1)
}

const defaultRole = {
  id: '',
  name: '',
  description: '',
  routes: [],
  api: []
}

const searchCfg = reactive({
  role_name:"",
})

const _this = reactive({

  asyncRoutes,
  allAsyncRoutes:[],
  editLinkRoleId:0,
  loading:false,
  allApiConfig: [],
  value: [],
  filterMethod(query, item) {
    return item.label.indexOf(query) > -1
  },
  urlConfig: [],
  urlConfigMap: [],
  filterText: '',
  filterText4Api: '',
  role: Object.assign({}, defaultRole),
  routes: [],
  rolesList: [],
  dialogVisible: false,
  esLinkDialogVisible:false,
  dialogType: 'new',
  checkStrictly: false,
  defaultProps: {
    children: 'children',
    label: 'title'
  },
  defaultProps4Api: {
    children: 'options',
    label: 'label'
  },
  chanCfgList: [],
  apiCfg:[],
  selectAllApiCfg:[],
  apiMap:{},
  data: [],
})

const menu_tree = ref()
const api_tree = ref()

watch(_this.filterText, (val) => {
  menu_tree.value.filter(val)
});

watch(_this.filterText4Api, (val) => {
  api_tree.value.filter(val)
});

onMounted(async ()=>{
  await GetUrlOpt()
  GetRoutes()
  GetRoles(1)
})

//获取接口权限配置
const  GetUrlOpt = async () => {
  let res = await UrlConfig({need_auth:true})

  if (res.code == 0) {
    let cfg_with_module = JSON.parse(JSON.stringify(res.data.cfg_with_module))
    let i = cfg_with_module.length;
    while(i--) {
      let router = cfg_with_module[i]
      let j = router.options.length;
      let options = router.options
      while(j--) {
        if(!options[j].needAuth){
          options.splice(j, 1);
        }
        if(options.length == 0){
          cfg_with_module.splice(i,1)
        }
      }
    }
    _this.apiMap = {}
    _this.selectAllApiCfg = []
    for(let group of cfg_with_module){

      for(let router of group.options){
        _this.selectAllApiCfg.push(router)
        _this.apiMap[router["value"]] = router
      }
    }

    _this.apiCfg = cfg_with_module

  }

}
//菜单权限全选
const selectAll= () => {
  menu_tree.value.setCheckedNodes(_this.routes)
}
const selectAll4Api = () => {
  api_tree.value.setCheckedNodes(_this.selectAllApiCfg)
}

const GetRoutes = async () =>{

  let res = await GetRoutesConfig({
    "routers":asyncRoutes,
  })

  _this.allAsyncRoutes = res.data

  var routers = res.data
  _this.routes = generateRoutes(routers)
}

const onlyOneShowingChildFn = (children = [], parent) => {
  let onlyOneChild = null
  const showingChildren = children.filter(item => !item.hidden)

  if (showingChildren.length === 1) {
    onlyOneChild = showingChildren[0]
    onlyOneChild.path =  path.resolve(parent.path, onlyOneChild.path)
    return onlyOneChild
  }

  if (showingChildren.length === 0) {
    onlyOneChild = { ...parent, path: '', noShowingChildren: true }
    return onlyOneChild
  }

  return false
}

const generateRoutes = (routes, basePath = '/')=> {
  const res = []

  for (let route of routes) {

    if (route.hidden && !route.service) {
      continue
    }

    const onlyOneShowingChild = onlyOneShowingChildFn(route.children, route)

    if (route.children && route.children.length > 0  && onlyOneShowingChild && !route.alwaysShow) {
      route = onlyOneShowingChild
    }

    const data = {
      path: path.resolve(basePath, route.path),
      title: route.meta && route.meta.title
    }

    if (route.children && route.children.length >0 ) {
      data.children = generateRoutes(route.children, data.path)
    }
    res.push(data)
  }
  return res
}
//配置tree所需要的配置
const generateArr = (routes) => {
  let data = []
  routes.forEach(route => {
    data.push(route)
    if (route.children && route.children.length >0 ) {
      const temp = generateArr(route.children)
      if (temp.length > 0) {
        data = [...data, ...temp]
      }
    }
  })
  return data
}
//生成菜单
const generateTree = (routes, basePath = '/', checkedKeys) => {
  const res = []
  for (const route of routes) {
    if (route.children && route.children.length >0 ) {
      route.children = generateTree(route.children, path.resolve(basePath, route.path), checkedKeys)
    }

    if (checkedKeys.includes(path.resolve(basePath, route.path)) || (route.children && route.children.length >= 1)) {
      res.push(route)
    }
  }
  return res
}
const filterNode = (value, data) => {
  if (!value) return true
  return data.title.indexOf(value) !== -1
}
const filterNode4Api = (value, data) => {
  if (!value) return true
  return data.label.indexOf(value) !== -1
}
const handleAddRole = () => {
  _this.role = Object.assign({}, defaultRole)
  if (menu_tree.value) {
    menu_tree.value.setCheckedNodes([])
  }
  if (api_tree.value) {
    api_tree.value.setCheckedNodes([])
  }
  _this.dialogType = 'new'
  _this.dialogVisible = true
}
const handleEdit = (row) => {
  _this.dialogType = 'edit'
  _this.dialogVisible = true
  _this.checkStrictly = true
  _this.role = deepClone(row)

  nextTick(() => {
    const routes = generateRoutes(_this.role.routes)

    menu_tree.value.setCheckedNodes(generateArr(routes))
    //todo...
    let apiList = []
    for(let v of _this.role.api){
      apiList.push(_this.apiMap[v])
    }
    if(api_tree.value){
      api_tree.value.setCheckedNodes(apiList)
    }

    _this.checkStrictly = false
  })
}

const input =reactive({
  page:1, 
  limit:10,
})

const GetRoles = async (page) => {
  !page ? input.page = 1 : input.page = page
  _this.loading=true
  const res = await getRoles({
    role_name:searchCfg.role_name,
    page:page,
    page_size:input.limit
  })
  _this.loading=false

  for (var k in res.data.list) {
    res.data.list[k]['routes'] = JSON.parse(res.data.list[k]['routes'])
  }
  count.value = res.data.count
  _this.rolesList = res.data.list
}
const handleDelete = ({ $index, row }) => {
  ElMessageBox.confirm('确定删除这个角色吗?', '警告', {
    confirmButtonText: '确定',
    cancelButtonText: '取消',
    type: 'warning'
  })
    .then(async() => {
      let res = await deleteRole({ id: row.id })
      if (res.code != 0) {
        ElMessage.error({
          type: 'error',offset:60,
          message: res.msg
        })
        return
      }
      await GetRoles(1)
      if (res.code != 0) {
        ElMessage.error({
          type: 'error',offset:60,
          message: res.msg
        })
        return
      }
      ElMessage.success({
        type: 'success',offset:60,
        message: res.msg
      })
      return

    })
    .catch(err => {
      console.error(err)
    })
}

const confirmRole = async () => {
  if (api_tree.value){
    const checkedKeys4Api = api_tree.value.getCheckedNodes(true,false)
    let api = []

    for(let v of checkedKeys4Api){
      api.push(v.value)
    }

    _this.role.api = api
  }


  const isEdit = _this.dialogType === 'edit'
  const checkedKeys = menu_tree.value.getCheckedKeys()
  _this.role.routes = generateTree(deepClone(_this.allAsyncRoutes), '/', checkedKeys)
  var roleModel = _this.role
  roleModel.routes = JSON.stringify(roleModel.routes)
  if (isEdit) {
    roleModel['id'] = _this.role.id
    let res = await updateRole(roleModel)
    if (res.code != 0) {
      ElMessage.error({
        type: 'error',offset:60,
        message: res.msg
      })
      return
    }
    GetRoles(1)
  } else {
    roleModel.id = 0
    const res = await addRole(roleModel)
    if (res.code != 0) {
      ElMessage.error({
        type: 'error',offset:60,
        message: res.msg
      })
      return
    }
    GetRoles(1)
  }
  _this.dialogVisible = false
  ElMessage.success({
    type: 'success',offset:60,
    message: "操作成功"
  })
}

// 添加 activeTab 响应式变量
const activeTab = ref('basic')

// 添加折叠状态控制
const isCollapsed = ref(true)

// 切换折叠状态
const toggleCollapse = () => {
  isCollapsed.value = !isCollapsed.value
}

</script>

<style>
.el-transfer-panel {
  width: 400px;
  height: 600px;
}

.el-transfer-panel__list.is-filterable {
  height: 500px;
}

</style>

<style lang="scss" scoped>
.app-container {
  // 搜索区域样式
  .search-container {
    margin-bottom: 1.5rem;
    padding: 1rem;
    backdrop-filter: blur(8px);
    border-radius: 0.5rem;
    box-shadow: 0 1px 2px rgba(0, 0, 0, 0.05);
    transition: all 0.3s;

    :deep(.el-form) {
      display: flex;
      flex-wrap: wrap;
      align-items: center;
      gap: 1rem;

      .el-form-item {
        margin-bottom: 0;
        margin-right: 0;

        &__label {
          font-weight: 500;
        }
      }

      .el-input {
        .el-input__wrapper {
          box-shadow: 0 1px 2px rgba(0, 0, 0, 0.05);
          border: 1px solid #e5e7eb;
          transition: all 0.3s;

          &:hover {
            border-color: #60a5fa;
          }

          &.is-focus {
            border-color: #3b82f6;
            box-shadow: 0 0 0 2px rgba(59, 130, 246, 0.2);
          }
        }
      }

      .el-button {
        padding: 8px 16px;
        transition: all 0.3s;
        
        &:hover {
          transform: translateY(-1px);
          box-shadow: 0 4px 6px rgba(0, 0, 0, 0.1);
        }

        &.filter-item {
          min-width: 80px;
        }
      }
    }

    // 添加移动端折叠相关样式
    .search-header {
      display: none; // 默认隐藏
    }

    @media (max-width: 768px) {
      .search-header {
        display: flex;
        justify-content: space-between;
        align-items: center;
        padding: 0.5rem 1rem;
        cursor: pointer;
        user-select: none;
        
        span {
          font-weight: 500;
        }

        .el-icon {
          transition: transform 0.3s ease;
          
          &.is-collapsed {
            transform: rotate(-180deg);
          }
        }
      }

      &.is-collapsed {
        :deep(.el-form) {
          display: none;
        }
      }

      :deep(.el-form) {
        animation: slideDown 0.3s ease;
        
        .el-form-item {
          margin-right: 0;
          width: 100%;
          
          .el-input {
            width: 100% !important;
          }
          
          // 按钮组样式
          &.button-group {
            display: flex;
            gap: 0.5rem;
            margin-top: 0.5rem;
            
            .el-button {
              flex: 1;
              margin: 0 !important;
              width: auto;
            }
          }
        }
      }
    }
  }

  // 表格容器样式
  .table-container {
    border-radius: 0.5rem;
    overflow: hidden;
    backdrop-filter: blur(8px);
    transition: all 0.3s;
    box-shadow: 0 1px 2px rgba(0, 0, 0, 0.05);

    :deep(.el-table) {
      background-color: transparent;
      
      // 表格行样式
      .el-table__row {
        transition: all 0.3s;
        
        &:hover {
          transform: translateY(-1px);
          box-shadow: 0 2px 4px rgba(0, 0, 0, 0.05);
        }
      }
      
      // 表头样式
      th {
        font-weight: 500;
        border-bottom: 1px solid #e5e7eb;
        padding: 12px 0;
        
        .cell {
          font-size: 0.95rem;
        }
      }
      
      // 单元格样式
      td {
        border-bottom: 1px solid #f3f4f6;
        padding: 16px 0;
        
        .cell {
          line-height: 1.6;
        }
      }

      // 操作按钮样式
      .el-button {
        padding: 6px 12px;
        transition: all 0.3s;
        border-radius: 6px;
        
        &:not(:last-child) {
          margin-right: 8px;
        }
        
        &:hover {
          transform: translateY(-1px);
          box-shadow: 0 2px 4px rgba(0, 0, 0, 0.1);
        }

        &.el-button--primary {
          background: linear-gradient(135deg, #60a5fa, #3b82f6);
          border: none;
          
          &:hover {
            background: linear-gradient(135deg, #3b82f6, #2563eb);
          }
        }

        &.el-button--danger {
          background: linear-gradient(135deg, #f87171, #ef4444);
          border: none;
          
          &:hover {
            background: linear-gradient(135deg, #ef4444, #dc2626);
          }
        }
      }
    }
  }

  // 分页容器样式
  .pagination-container {
    margin-top: 1rem;
    display: flex;
    justify-content: flex-end;
    
    :deep(.el-pagination) {
      padding: 0.75rem;
      border-radius: 0.5rem;
      box-shadow: 0 1px 2px rgba(0, 0, 0, 0.05);

      .el-pager li {
        border-radius: 4px;
        transition: all 0.3s;
      
        &.is-active {
          background-color: #3b82f6;
        }
      }
    }
  }

  // 抽屉内容样式
  :deep(.el-drawer__body) {
    padding: 0;
    display: flex;
    flex-direction: column;

    .el-tabs {
      flex: 1;
      display: flex;
      flex-direction: column;
      
      .el-tabs__header {
        margin: 0;
        padding: 0 20px;
        border-bottom: 1px solid #e5e7eb;

        .el-tabs__nav-wrap {
          &::after {
            display: none;
          }

          .el-tabs__nav {
            border: none;
          }

          .el-tabs__item {
            height: 48px;
            line-height: 48px;
            padding: 0 24px;
            font-size: 14px;
            transition: all 0.3s;

            &.is-active {
              font-weight: 500;
            }
          }
        }
      }
      
      .el-tabs__content {
        flex: 1;
        padding: 20px;
        overflow-y: auto;
      }
    }

    // 权限树样式
    .permission-tree {
      border: 1px solid #e5e7eb;
      border-radius: 6px;
      padding: 12px;
      margin-top: 12px;

      .el-tree-node__content {
        height: 36px;
        transition: all 0.3s;
      }

      .el-checkbox__input {
        .el-checkbox__inner {
          border-radius: 4px;
        }
      }
    }
  }

  .permission-header {
    display: flex;
    gap: 12px;
    margin-bottom: 16px;
    align-items: center;

    .filter-input {
      width: 300px;
    }
  }

  .drawer-footer {
    padding: 16px 20px;
    border-top: 1px solid #e5e7eb;
    display: flex;
    justify-content: flex-end;
    gap: 12px;
  }
}

@keyframes slideDown {
  from {
    opacity: 0;
    transform: translateY(-10px);
  }
  to {
    opacity: 1;
    transform: translateY(0);
  }
}

// 深色模式
@media (prefers-color-scheme: dark) {
  .app-container {
    .search-container {
      background-color: rgba(31, 41, 55, 0.9);

      .search-header {
        color: #e5e7eb;
      }
    }

    .table-container {
      background-color: rgba(31, 41, 55, 0.9);
    }

    :deep(.el-drawer__body) {
      background-color: #1f2937;

      .el-tabs {
        .el-tabs__header {
          border-bottom-color: #374151;
        }
      }

      .permission-tree {
        border-color: #4b5563;
        background-color: rgba(55, 65, 81, 0.3);

        .el-tree-node__content:hover {
          background-color: rgba(75, 85, 99, 0.3);
        }
      }
    }

    .drawer-footer {
      border-top-color: #4b5563;
      background-color: #1f2937;
    }
  }
}
</style>
