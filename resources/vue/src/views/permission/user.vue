<template>
  <div class="app-container">

    <div class="filter-container">
      <el-button
        size="mini"
        type="primary"
        icon="el-icon-plus"
        class="filter-item"
        @click="handleAddRole"
      >{{ $t('新建用户') }}
      </el-button>
    </div>
    <el-table

      :data="rolesList"
    >
      <el-table-column
        :label="$t('序号')"
        align="center"
        fixed
        width="50"
      >
        <template slot-scope="scope">
          {{ scope.$index + 1 }}
        </template>
      </el-table-column>
      <el-table-column align="center" label="id" width="220">
        <template slot-scope="scope">
          {{ scope.row.id }}
        </template>
      </el-table-column>
      <el-table-column align="center" :label="$t('用户名')" width="220">
        <template slot-scope="scope">
          {{ scope.row.username }}
        </template>
      </el-table-column>
      <el-table-column align="center" :label="$t('所属角色')" width="220">
        <template slot-scope="scope">
          {{ scope.row.role_name }}
        </template>
      </el-table-column>
      <el-table-column align="center" :label="$t('真实姓名')">
        <template slot-scope="scope">
          {{ scope.row.realname }}
        </template>
      </el-table-column>
      <el-table-column align="center" :label="$t('操作')" fixed="right" width="200">
        <template slot-scope="scope">
          <el-button

            type="primary"
            size="mini"
            icon="el-icon-edit"
            @click="handleEdit(scope)"
          >{{ $t('编辑') }}
          </el-button>
          <el-button

            type="danger"
            size="mini"
            icon="el-icon-delete"
            @click="handleDelete(scope)"
          >{{ $t('删除') }}
          </el-button>
        </template>
      </el-table-column>
    </el-table>

    <el-dialog
      :close-on-click-modal="false"
      :visible.sync="dialogVisible"
      :title="dialogType==='edit'?$t('编辑用户信息'):$t('新建用户信息')"
    >
      <el-form :model="role" label-width="100px" label-position="left">

        <el-form-item :label="$t('用户名')">
          <el-input v-model="role.username" :placeholder="$t('用户名')" />
        </el-form-item>
        <el-form-item :label="$t('密码')">
          <el-input v-model="role.password" :placeholder="$t('密码')" />
        </el-form-item>
        <el-form-item :label="$t('真实姓名')">
          <el-input v-model="role.realname" :placeholder="$t('真实姓名')" />
        </el-form-item>
        <el-form-item :label="$t('请选择角色')">
          <el-select v-model="role.role_id" :placeholder="$t('请选择角色')" clearable filterable>
            <el-option
              v-for="item in chanCfgList"
              :key="item.id"
              :label="item.name"
              :value="item.id"
            />
          </el-select>
        </el-form-item>

      </el-form>
      <div style="text-align:right;">
        <el-button
          size="mini"
          type="danger"
          icon="el-icon-close"
          @click="dialogVisible=false"
        >{{ $t('取消') }}
        </el-button>
        <el-button
          size="mini"
          type="primary"
          icon="el-icon-check"
          @click="confirmRole"
        >{{ $t('确认') }}
        </el-button>
      </div>
    </el-dialog>
  </div>
</template>

<script>
import { deepClone } from '@/utils'
import { DelUser, InsertUser, roleOption, UpdateUser, userList } from '@/api/user'

const defaultUser = {
  id: '',
  password: '',
  realname: '',
  role_id: '',
  username: '',
  role_name: ''
}

export default {
  data() {
    return {
      role: Object.assign({}, defaultUser),
      routes: [],
      rolesList: [],
      dialogVisible: false,
      dialogType: 'new',
      checkStrictly: false,
      defaultProps: {
        children: 'children',
        label: 'title'
      },
      chanCfgList: [],
      chanCfgMap: []
    }
  },
  computed: {},
  async created() {
    // Mock: get all routes and roles list from server
    await this.getRoleOpt()
    this.getUserList()
  },
  methods: {
    async getRoleOpt() {
      const res = await roleOption()
      for (var v of res.data) {
        this.chanCfgList.push(v)
        this.chanCfgMap[v['id']] = v['name']
      }
    },

    async getUserList() {
      const res = await userList()
      for (var k in res.data) {
        res.data[k]['role_name'] = this.chanCfgMap[res.data[k]['role_id']]
      }

      this.rolesList = res.data
    },

    handleAddRole() {
      this.role = Object.assign({}, defaultUser)
      this.dialogType = 'new'
      this.dialogVisible = true
    },
    handleEdit(scope) {
      this.dialogType = 'edit'
      this.dialogVisible = true
      this.checkStrictly = true
      this.role = deepClone(scope.row)
      this.role.password = ''
    },
    handleDelete({ $index, row }) {
      this.$confirm('确定删除该用户吗?', '警告', {
        confirmButtonText: '确认',
        cancelButtonText: '取消',
        type: 'warning'
      })
        .then(async() => {
          await DelUser({ id: row.id })
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
    async confirmRole() {
      const isEdit = this.dialogType === 'edit'
      if (this.role.password.length < 5 && this.role.password.trim() != '') {
        this.$message({
          message: '密码长度必须大于5',
          type: 'error'
        })
        return false
      }

      if (this.role.role_id == '' || this.role.role_id == 0) {
        this.$message({
          message: '角色不能为空',
          type: 'error'
        })
        return false
      }

      if (isEdit) {
        await UpdateUser(this.role)
        this.getUserList()
      } else {
        const { data } = await InsertUser(this.role)
        this.getUserList()
        /* this.role.id = data
          this.rolesList.push(this.role)*/
      }

      const { username, password, id, role_name, realname } = this.role
      this.dialogVisible = false
      this.$notify({
        title: 'Success',
        dangerouslyUseHTMLString: true,
        message: `
            <div>id: ${id}</div>
            <div>用户名: ${username}</div>
            <div>密码: ${password}</div>
            <div>角色: ${role_name}</div>
            <div>真实姓名: ${realname}</div>
          `,
        type: 'success'
      })
    }
  }
}
</script>

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
