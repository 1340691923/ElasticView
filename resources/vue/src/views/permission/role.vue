<template>
  <div class="app-container">

    <div class="search-container">
      <el-form :inline="true">
        <el-form-item label="角色名:" prop="keywords">

          <el-input clearable v-model="searchCfg.role_name" style="width:160px"></el-input>

          <el-button
            style="margin-left: 10px"
            type="success"
            class="filter-item"
            @click="GetRoles(1)"
          >
            {{ $t('搜索') }}
          </el-button>
        </el-form-item>
        <el-form-item label="" prop="keywords">
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
        <el-form-item v-if="_this.role.id != 1" :label="$t('菜单栏')">
          <el-input v-model="_this.filterText" :placeholder="$t('输入关键字进行过滤')" style="width: 300px" />
          <el-button
            @click="selectAll"
          >{{ $t('全选') }}
          </el-button>

        </el-form-item>
        <el-form-item  v-if="_this.role.id != 1" >
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
        </el-form-item>
        <el-form-item v-if="_this.role.id != 1" :label="$t('接口权限')">
          <el-input v-model="_this.filterText4Api" :placeholder="$t('输入关键字进行过滤')" style="width: 300px" />
          <el-button

            @click="selectAll4Api"
          >{{ $t('全选') }}
          </el-button>

        </el-form-item>
        <el-form-item v-if="_this.role.id != 1" >
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
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button

          type="danger"

          @click="_this.dialogVisible=false"
        >{{ $t('取消') }}
        </el-button>
        <el-button

          type="primary"

          @click="confirmRole"
        >{{ $t('确定') }}
        </el-button>
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
import {Plus,Search,Edit,Delete} from "@element-plus/icons-vue";
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
          type: 'error',
          message: res.msg
        })
        return
      }
      await GetRoles(1)
      if (res.code != 0) {
        ElMessage.error({
          type: 'error',
          message: res.msg
        })
        return
      }
      ElMessage.success({
        type: 'success',
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
        type: 'error',
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
        type: 'error',
        message: res.msg
      })
      return
    }
    GetRoles(1)
  }
  _this.dialogVisible = false
  ElMessage.success({
    type: 'success',
    message: "操作成功"
  })
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

.roles-table {
  margin-top: 30px;
}

.permission-tree {
  margin-bottom: 30px;
}

}
</style>
