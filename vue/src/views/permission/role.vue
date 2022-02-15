<template>
  <div class="app-container">
    <el-card class="box-card">
      <div class="filter-container">
        <el-button type="primary" class="filter-item" icon="el-icon-plus" @click="handleAddRole">新建角色</el-button>
      </div>
      <el-table
        :data="rolesList"
        @row-dblclick="handleEdit"
      >
        <el-table-column
          label="序号"
          align="center"
          fixed
          width="50"
        >
          <template slot-scope="scope">
            {{ scope.$index+1 }}
          </template>
        </el-table-column>
        <el-table-column align="center" label="角色id" width="220">
          <template slot-scope="scope">
            {{ scope.row.id }}
          </template>
        </el-table-column>
        <el-table-column align="center" label="角色名" width="220">
          <template slot-scope="scope">
            {{ scope.row.name }}
          </template>
        </el-table-column>
        <el-table-column align="header-center" label="角色详细信息">
          <template slot-scope="scope">
            {{ scope.row.description }}
          </template>
        </el-table-column>
        <el-table-column align="center" label="操作" width="220" fixed="right">
          <template slot-scope="scope">
            <el-button type="primary" size="small" icon="el-icon-edit" @click.stop="handleEdit(scope.row)">编辑
            </el-button>
            <el-button type="danger" size="small" icon="el-icon-delete" @click.stop="handleDelete(scope)">删除</el-button>
          </template>
        </el-table-column>
      </el-table>

      <el-dialog :close-on-click-modal="false" width="70%" :visible.sync="dialogVisible" :title="dialogType==='edit'?'修改角色':'新建角色'">
        <el-form :model="role" label-width="120px" label-position="left">
          <el-form-item label="角色名">
            <el-input v-model="role.name" placeholder="角色名" />
          </el-form-item>
          <el-form-item label="角色详情信息">
            <el-input
              v-model="role.description"
              :autosize="{ minRows: 2, maxRows: 4}"
              type="textarea"
              placeholder="角色详情信息"
            />
          </el-form-item>
          <el-form-item label="菜单栏">
            <el-input v-model="filterText" placeholder="输入关键字进行过滤" style="width: 300px" />
            <el-button icon="el-icon-check" @click="quanxuan">全选</el-button>
            <el-tree
              ref="tree"
              :filter-node-method="filterNode"
              :check-strictly="checkStrictly"
              :data="routesData"
              :props="defaultProps"
              show-checkbox
              node-key="path"
              class="permission-tree"
            />
          </el-form-item>
          <el-form-item v-if="role.id != 1" label="接口权限设置">
            <el-transfer
              v-if="dialogVisible"
              v-model="role.api"
              :titles="['全部接口权限', '角色拥有权限']"
              :button-texts="['移除权限', '添加权限']"
              filterable
              :filter-method="filterMethod"
              filter-placeholder="请选择接口名"
              :data="allApiConfig"
            />
          </el-form-item>
        </el-form>
        <div style="text-align:right;">
          <el-button type="danger" icon="el-icon-close" @click="dialogVisible=false">取消</el-button>
          <el-button type="primary" icon="el-icon-check" @click="confirmRole">确定</el-button>
        </div>
      </el-dialog>
    </el-card>
  </div>
</template>

<script>
import path from 'path'
import { deepClone } from '@/utils'
import { asyncRoutes } from '@/utils/router'
import { addRole, deleteRole, getRoles, updateRole } from '@/api/role'
import { UrlConfig } from '@/api/api-rbac'

const defaultRole = {
  id: '',
  name: '',
  description: '',
  routes: [],
  api: []
}

export default {
  data() {
    return {
      allApiConfig: [],
      value: [],
      filterMethod(query, item) {
        return item.label.indexOf(query) > -1
      },
      urlConfig: [],
      urlConfigMap: [],
      filterText: '',
      role: Object.assign({}, defaultRole),
      routes: [],
      rolesList: [],
      dialogVisible: false,
      dialogType: 'new',
      checkStrictly: false,
      defaultProps: {
        children: 'children',
        label: 'title'
      },
      chanCfgList: []
    }
  },
  computed: {
    routesData() {
      return this.routes
    }
  },
  watch: {
    filterText(val) {
      this.$refs['tree'].filter(val)
    }
  },
  created() {
    this.getUrlOpt()
    this.getRoutes()
    this.getRoles()
  },
  methods: {
    getUrlOpt() {
      UrlConfig().then(res => {
        if (res) {
          for (var v of res.data) {
            const url = v['url']
            const remark = v['remark']
            this.allApiConfig.push(
              {
                label: remark,
                key: url,
                disabled: false
              }
            )
          }
        }
      })
    },
    quanxuan() {
      this.$refs.tree.setCheckedNodes(this.routes)
    },
    async getRoutes() {
      var routers = asyncRoutes
      this.serviceRoutes = routers
      this.routes = this.generateRoutes(routers)
    },
    async getRoles() {
      const res = await getRoles()
      for (var k in res.data) {
        console.log(res.data[k])
        res.data[k]['routes'] = JSON.parse(res.data[k]['routes'])
      }
      this.rolesList = res.data
    },

    generateRoutes(routes, basePath = '/') {
      const res = []

      for (let route of routes) {
        if (route.hidden && !route.service) {
          continue
        }

        const onlyOneShowingChild = this.onlyOneShowingChild(route.children, route)

        if (route.children && onlyOneShowingChild && !route.alwaysShow) {
          route = onlyOneShowingChild
        }

        const data = {
          path: path.resolve(basePath, route.path),
          title: route.meta && route.meta.title
        }

        if (route.children) {
          data.children = this.generateRoutes(route.children, data.path)
        }
        res.push(data)
      }
      return res
    },
    generateArr(routes) {
      let data = []
      routes.forEach(route => {
        data.push(route)
        if (route.children) {
          const temp = this.generateArr(route.children)
          if (temp.length > 0) {
            data = [...data, ...temp]
          }
        }
      })
      return data
    },

    filterNode(value, data) {
      if (!value) return true
      return data.title.indexOf(value) !== -1
    },
    handleAddRole() {
      this.role = Object.assign({}, defaultRole)
      if (this.$refs.tree) {
        this.$refs.tree.setCheckedNodes([])
      }
      this.dialogType = 'new'
      this.dialogVisible = true
    },
    handleEdit(row) {
      console.log(row, 'row')
      this.dialogType = 'edit'
      this.dialogVisible = true
      this.checkStrictly = true
      this.role = deepClone(row)

      this.$nextTick(() => {
        const routes = this.generateRoutes(this.role.routes)
        this.$refs.tree.setCheckedNodes(this.generateArr(routes))
        this.checkStrictly = false
      })
    },
    handleDelete({ $index, row }) {
      this.$confirm('确定删除这个角色吗?', '警告', {
        confirmButtonText: '确定',
        cancelButtonText: '取消',
        type: 'warning'
      })
        .then(async() => {
          await deleteRole({ id: row.id })
          this.rolesList.splice($index, 1)
          this.$message({
            type: 'success',
            message: 'Delete succed!'
          })
        })
        .catch(err => {
          console.error(err)
        })
    },
    generateTree(routes, basePath = '/', checkedKeys) {
      const res = []

      for (const route of routes) {
        const routePath = path.resolve(basePath, route.path)

        if (route.children) {
          route.children = this.generateTree(route.children, routePath, checkedKeys)
        }

        if (checkedKeys.includes(routePath) || (route.children && route.children.length >= 1)) {
          res.push(route)
        }
      }
      return res
    },
    async confirmRole() {
      const isEdit = this.dialogType === 'edit'
      const checkedKeys = this.$refs.tree.getCheckedKeys()
      this.role.routes = this.generateTree(deepClone(this.serviceRoutes), '/', checkedKeys)
      var roleModel = this.role
      roleModel.routes = JSON.stringify(roleModel.routes)
      if (isEdit) {
        roleModel['id'] = this.role.id
        await updateRole(roleModel)
        for (let index = 0; index < this.rolesList.length; index++) {
          if (this.rolesList[index].id === this.role.id) {
            this.rolesList.splice(index, 1, Object.assign({}, this.role))
            break
          }
        }
      } else {
        roleModel.id = 0
        const { data } = await addRole(roleModel)
        this.role.id = data.id
        this.rolesList.push(this.role)
      }

      const { description, id, name } = this.role
      this.dialogVisible = false
      this.$notify({
        title: 'Success',
        dangerouslyUseHTMLString: true,
        message: `
            <div>id: ${id}</div>
            <div>Role Name: ${name}</div>
            <div>Description: ${description}</div>
          `,
        type: 'success'
      })
    },

    onlyOneShowingChild(children = [], parent) {
      let onlyOneChild = null
      const showingChildren = children.filter(item => !item.hidden)

      if (showingChildren.length === 1) {
        onlyOneChild = showingChildren[0]
        onlyOneChild.path = path.resolve(parent.path, onlyOneChild.path)
        return onlyOneChild
      }

      if (showingChildren.length === 0) {
        onlyOneChild = { ...parent, path: '', noShowingChildren: true }
        return onlyOneChild
      }

      return false
    }
  }
}
</script>

<style>
  .el-transfer-panel{
    width:  400px;
    height: 600px;
  }
  .el-transfer-panel__list.is-filterable{
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
