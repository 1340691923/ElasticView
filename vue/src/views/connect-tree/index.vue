<template>
  <div class="app-container">
    <el-card class="box-card">
      <div class="filter-container">
        <el-button type="primary" icon="el-icon-plus" class="filter-item" @click="handleAddRole">新建连接信息</el-button>
      </div>
      <back-to-top />
      <el-table

        :data="list"
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

        <el-table-column align="center" label="HOST" width="220">
          <template slot-scope="scope">
            {{ scope.row.ip }}
          </template>
        </el-table-column>
        <el-table-column align="center" label="用户名" width="300">
          <template slot-scope="scope">
            {{ scope.row.user }}
          </template>
        </el-table-column>
        <el-table-column align="center" label="密码" width="300">
          <template slot-scope="scope">
            {{ scope.row.pwd }}
          </template>
        </el-table-column>
        <el-table-column align="center" label="备注" width="220">
          <template slot-scope="scope">
            {{ scope.row.remark }}
          </template>
        </el-table-column>
        <el-table-column align="center" label="版本" width="100">
          <template slot-scope="scope">
            {{ scope.row.version }}
          </template>
        </el-table-column>
        <el-table-column align="center" label="创建时间" width="220">
          <template slot-scope="scope">
            {{ scope.row.created }}
          </template>
        </el-table-column>
        <el-table-column align="center" label="修改时间" width="220">
          <template slot-scope="scope">
            {{ scope.row.updated }}
          </template>
        </el-table-column>
        <el-table-column align="center" label="操作" fixed="right" width="300">
          <template slot-scope="scope">
            <el-button
              v-loading="scope.row.connectLoading"
              type="success"
              size="small"
              icon="el-icon-link"
              @click="testConnect(scope)"
            >测试连接
            </el-button>
            <el-button type="primary" size="small" icon="el-icon-edit" @click="handleEdit(scope)">编辑</el-button>
            <el-button type="danger" size="small" icon="el-icon-delete" @click="handleDelete(scope)">删除</el-button>
          </template>
        </el-table-column>
      </el-table>

      <el-dialog :close-on-click-modal="false" :visible.sync="dialogVisible" :title="dialogType==='edit'?'编辑用户信息':'新建用户信息'">
        <el-form :model="link" label-width="100px" label-position="left">
          <el-form-item label="IP">
            <el-input v-model="link.ip" placeholder="例如:http://127.0.0.1:9200" />
          </el-form-item>
          <el-form-item label="用户名">
            <el-input v-model="link.user" placeholder="用户名" />
          </el-form-item>
          <el-form-item label="密码">
            <el-input v-model="link.pwd" placeholder="密码" />
          </el-form-item>
          <el-form-item label="备注">
            <el-input v-model="link.remark" placeholder="备注" />
          </el-form-item>
          <el-form-item label="版本">
            <el-select v-model="link.version" placeholder="请选择版本" filterable>
              <el-option label="6" :value="Number(6)" />
              <el-option label="7" :value="Number(7)" />
            </el-select>
          </el-form-item>
        </el-form>
        <div style="text-align:right;">
          <el-button
            v-loading="testConnectLoading"
            type="success"
            icon="el-icon-link"
            @click="testConnectForm"
          >测试连接
          </el-button>
          <el-button type="danger" icon="el-icon-close" @click="dialogVisible=false">取消</el-button>
          <el-button type="primary" icon="el-icon-check" @click="confirm">确认</el-button>
        </div>
      </el-dialog>
    </el-card>
  </div>
</template>

<script>
import { deepClone } from '@/utils'
import { DeleteAction, InsertAction, ListAction, UpdateAction } from '@/api/es-link'
import { PingAction } from '@/api/es'

const defaultLink = {
  created: '',
  id: 0,
  ip: 'http://127.0.0.1:9200',
  pwd: '',
  remark: '',
  updated: '',
  user: '',
  version: 6
}

export default {
  components: {
    'BackToTop': () => import('@/components/BackToTop/index')
  },
  data() {
    return {
      testConnectLoading: false,
      connectLoading: false,
      link: Object.assign({}, defaultLink),
      list: [],
      dialogVisible: false,
      dialogType: 'new',
      checkStrictly: false,
      defaultProps: {
        children: 'children',
        label: 'title'
      }
    }
  },
  computed: {},
  created() {
    this.getList()
  },
  methods: {
    testConnectForm() {
      this.testConnectLoading = true
      PingAction(this.link).then(res => {
        if (res.code == 0) {
          this.$message({
            type: 'success',
            message: `连接成功,ES版本为 :${res.data.version.number}`
          })
        } else {
          this.$message({
            type: 'error',
            message: res.msg
          })
        }
        this.testConnectLoading = false
      }).catch(err => {
        this.testConnectLoading = false
      })
    },
    testConnect(scope) {
      this.list[scope.$index].connectLoading = true
      PingAction(scope.row).then(res => {
        if (res.code == 0) {
          this.$message({
            type: 'success',
            message: `连接成功,ES版本为 :${res.data.version.number}`
          })
        } else {
          this.$message({
            type: 'error',
            message: res.msg
          })
        }
        this.list[scope.$index].connectLoading = false
      }).catch(err => {
        this.list[scope.$index].connectLoading = false
      })
    },
    async getList() {
      const res = await ListAction()
      for (const k in res.data) {
        res.data[k]['connectLoading'] = false
      }
      this.list = res.data
    },

    handleAddRole() {
      this.link = Object.assign({}, defaultLink)
      this.dialogType = 'new'
      this.dialogVisible = true
    },
    handleEdit(scope) {
      this.dialogType = 'edit'
      this.dialogVisible = true
      this.checkStrictly = true
      this.link = deepClone(scope.row)
    },
    handleDelete({ $index, row }) {
      this.$confirm('确定删除该连接信息吗?', '警告', {
        confirmButtonText: '确认',
        cancelButtonText: '取消',
        type: 'warning'
      })
        .then(async() => {
          const res = await DeleteAction({ id: row.id })
          if (res.code != 0) {
            this.$message({
              type: 'error',
              message: res.msg
            })
            return
          }
          this.list.splice($index, 1)
          this.$message({
            type: 'success',
            message: 'Delete succed!'
          })
        })
        .catch(err => {
          console.error(err)
        })
    },
    async confirm() {
      if (this.link.remark == '') {
        this.$message({
          type: 'error',
          message: '请填写备注'
        })
        return
      }
      const isEdit = this.dialogType === 'edit'

      if (isEdit) {
        const res = await UpdateAction(this.link)
        if (res.code != 0) {
          this.$message({
            type: 'error',
            message: res.msg
          })
          return
        }
        this.getList()
      } else {
        const res = await InsertAction(this.link)
        if (res.code != 0) {
          this.$message({
            type: 'error',
            message: res.msg
          })
          return
        }
        this.getList()
      }

      const { created, ip, id, pwd, remark, updated, user, version } = this.link
      this.dialogVisible = false
      this.$notify({
        title: 'Success',
        dangerouslyUseHTMLString: true,
        message: `
            <div>id: ${id}</div>
            <div>IP: ${ip}</div>
            <div>用户名: ${user}</div>
            <div>密码: ${pwd}</div>
            <div>备注: ${remark}</div>
            <div>版本: ${version}</div>
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
