<template>
  <div class="app-container">

    <div class="filter-container">
      <el-button
        size="mini"
        type="primary"
        icon="el-icon-plus"
        class="filter-item"
        @click="handleAddRole"
      >{{ $t('新建连接信息') }}
      </el-button>
    </div>
    <back-to-top />
    <el-table
      :data="list"
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

      <el-table-column align="center" :label="$t('HOST')" width="220">
        <template slot-scope="scope">
          {{ scope.row.ip }}
        </template>
      </el-table-column>
      <el-table-column align="center" :label="$t('用户名')" width="300">
        <template slot-scope="scope">
          {{ scope.row.user }}
        </template>
      </el-table-column>

      <el-table-column align="center" :label="$t('备注')" width="220">
        <template slot-scope="scope">
          {{ scope.row.remark }}
        </template>
      </el-table-column>
      <el-table-column align="center" :label="$t('版本')" width="100">
        <template slot-scope="scope">
          {{ scope.row.version }}
        </template>
      </el-table-column>
      <el-table-column align="center" :label="$t('root证书')" width="100" show-overflow-tooltip>
        <template slot-scope="scope">
          {{ scope.row.rootpem }}
        </template>
      </el-table-column>
      <el-table-column align="center" :label="$t('cert证书')" width="100" show-overflow-tooltip>
        <template slot-scope="scope">
          {{ scope.row.certpem }}
        </template>
      </el-table-column>
      <el-table-column align="center" :label="$t('key证书')" width="100" show-overflow-tooltip>
        <template slot-scope="scope">
          {{ scope.row.keypem }}
        </template>
      </el-table-column>
      <el-table-column align="center" :label="$t('创建时间')" width="220">
        <template slot-scope="scope">
          {{ scope.row.created }}
        </template>
      </el-table-column>
      <el-table-column align="center" :label="$t('修改时间')" width="220">
        <template slot-scope="scope">
          {{ scope.row.updated }}
        </template>
      </el-table-column>
      <el-table-column align="center" :label="$t('操作')" fixed="right" width="300">
        <template slot-scope="scope">
          <el-button

            v-loading="scope.row.connectLoading"
            :disabled="scope.row.connectLoading"
            type="success"
            size="mini"
            icon="el-icon-link"
            @click="testConnect(scope)"
          >{{ $t('测试连接') }}
          </el-button>
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
      :title="dialogType==='edit'?$t('编辑连接信息'):$t('新建连接信息')"
    >
      <el-form :model="link" label-width="100px" label-position="left">
        <el-form-item :label="$t('IP')">
          <el-input v-model="link.ip" :placeholder="$t('例如:http://127.0.0.1:9200')" />
        </el-form-item>
        <el-form-item :label="$t('用户名')">
          <el-autocomplete
            v-model="link.user"
            clearable
            :fetch-suggestions="querySearch"
            :placeholder="$t('用户名')"
          >
            <i
              slot="suffix"
              class="el-icon-user"
            />
            <template slot-scope="{ item }">
              <span>{{ item.value }}</span>
            </template>

          </el-autocomplete>

        </el-form-item>
        <el-form-item :label="$t('密码')">
          <el-input clearable v-model="link.pwd" :placeholder="$t('密码')" />
        </el-form-item>
        <el-form-item :label="$t('root证书')">
          <el-input clearable v-model="link.rootpem" type="textarea" :placeholder="$t('root证书')" />
        </el-form-item>
        <el-form-item :label="$t('cert证书')">
          <el-input clearable v-model="link.certpem" type="textarea" :placeholder="$t('cert证书')" />
        </el-form-item>
        <el-form-item :label="$t('key证书')">
          <el-input clearable v-model="link.keypem" type="textarea" :placeholder="$t('key证书')" />
        </el-form-item>
        <el-form-item :label="$t('备注')">
          <el-input clearable v-model="link.remark" :placeholder="$t('备注')" />
        </el-form-item>
        <el-form-item :label="$t('版本')">
          <el-select v-model="link.version" :placeholder="$t('请选择版本')" filterable>
            <el-option label="6" :value="Number(6)" />
            <el-option label="7" :value="Number(7)" />
            <el-option label="8" :value="Number(8)" />
          </el-select>
        </el-form-item>
      </el-form>
      <div style="text-align:right;">
        <el-button
          v-loading="testConnectLoading"

          size="mini"
          :disabled="testConnectLoading"
          type="success"
          icon="el-icon-link"
          @click="testConnectForm"
        >{{ $t('测试连接') }}
        </el-button>
        <el-button
          size="mini"
          type="primary"
          icon="el-icon-check"
          @click="confirm"
        >{{ $t('确认') }}
        </el-button>
        <el-button
          size="mini"
          type="danger"
          icon="el-icon-close"
          @click="dialogVisible=false"
        >{{ $t('取消') }}
        </el-button>
      </div>
    </el-dialog>

  </div>
</template>

<script>
import { deepClone } from '@/utils'
import { DeleteAction, InsertAction, ListAction, UpdateAction } from '@/api/es-link'
import { PingAction } from '@/api/es'
import {filterData} from "@/utils/table";

const defaultLink = {
  created: '',
  id: 0,
  ip: 'http://127.0.0.1:9200',
  pwd: '',
  remark: '',
  updated: '',
  user: '',
  rootpem: '',
  certpem: '',
  keypem: '',
  version: 6
}

export default {
  components: {
    'BackToTop': () => import('@/components/BackToTop/index')
  },
  data() {
    return {
      usernameWord:[
        { "value": "elastic" },
      ],
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

    querySearch(queryString, cb) {
      var usernameWord = this.usernameWord;
      var results = queryString ? usernameWord.filter(this.createFilter(queryString)) : usernameWord;
      // 调用 callback 返回建议列表的数据
      cb(results);
    },
    createFilter(queryString) {
      return (usernameWord) => {
        return (usernameWord.value.toLowerCase().indexOf(queryString.toLowerCase()) === 0);
      };
    },
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
          console.log('res', res)
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
      this.link.pwd = ''
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

      this.dialogVisible = false
      this.$message({
        type: 'success',
        message: isEdit ? '修改成功' : '创建成功'
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
