<template>
  <div class="app-container">
    <el-card class="box-card">
      <div class="filter-container">
        <el-select
          v-model="input.operater_action"
          reserve-keyword
          collapse-tags
          placeholder="接口名"
          style="width: 250px;"
          class="filter-item"
          clearable
          filterable
          @change="getList(1)"
        >

          <el-option
            v-for="item in urlConfig"
            :key="item.url"
            :label="item.remark"
            :value="item.url"
          />
        </el-select>
        <el-select
          v-model="input.operater_role_id"
          placeholder="请选择角色"
          clearable
          class="filter-item"
          filterable
          @change="getList(1)"
        >
          <el-option label="请选择角色" :value="Number(0)" />
          <el-option
            v-for="item in roleConfig"
            :key="item.id"
            :label="item.name"
            :value="item.id"
          />
        </el-select>
        <el-select
          v-model="input.operater_id"
          placeholder="请选择用户"
          clearable
          class="filter-item"
          filterable
          @change="getList(1)"
        >
          <el-option label="请选择用户" :value="Number(0)" />
          <el-option
            v-for="item in userConfig"
            :key="item.id"
            :label="item.username"
            :value="item.id"
          />
        </el-select>
        <date :dates="input.date" class="filter-item" @changeDate="changeDate" />
        <el-button icon="el-icon-search" type="primary" class="filter-item" @click="getList(1)">查询</el-button>
      </div>

      <el-table
        v-loading="tableLoading"
        :data="logList"
        style="width: 100%;margin-top:30px;"
        element-loading-text="请给我点时间！"
      >
        <el-table-column
          label="序号"
          align="center"
          fixed
          min-width="50"
        >
          <template slot-scope="scope">
            {{ scope.$index + 1 }}
          </template>
        </el-table-column>

        <el-table-column align="center" label="ID" min-width="100">
          <template slot-scope="scope">
            {{ scope.row.id }}
          </template>
        </el-table-column>

        <el-table-column align="center" label="所属角色" min-width="100">
          <template slot-scope="scope">
            {{ roleConfigMap[scope.row.operater_id] }}
          </template>
        </el-table-column>

        <el-table-column align="center" label="操作用户" min-width="100">
          <template slot-scope="scope">
            {{ scope.row.operater_name }}
          </template>
        </el-table-column>

        <el-table-column align="center" label="操作接口" min-width="220">
          <template slot-scope="scope">
            {{ urlConfigMap[scope.row.operater_action] }}
          </template>
        </el-table-column>

        <el-table-column align="center" label="接口请求参数" min-width="270">
          <template slot-scope="scope">
            <div v-if="!scope.row.isFormatData">
              <span v-html="scope.row.body_str" />
              <a-tooltip placement="right" style="cursor: pointer">
                <template slot="title">
                  <span>格式化数据</span>
                </template>
                <a-button
                  v-show="scope.row.body_str!=''"
                  type="link"
                  icon="eye"
                  @click="lookData(scope.row.index,true)"
                />
              </a-tooltip>
            </div>
            <div v-else>
              <json-editor
                v-model="scope.row.dataFormat"
                font-size="15"
                height="400"
                class="req-body"
                styles="width: 100%"
                :read="true"
                title="上报数据"
              />
              <a-tooltip placement="right" style="cursor: pointer">
                <template slot="title">
                  <span>还原数据</span>
                </template>
                <a-button style="color: red" type="link" icon="eye-invisible" @click="lookData(scope.row.index,false)" />
              </a-tooltip>
            </div>
          </template>
        </el-table-column>

        <el-table-column align="center" label="创建时间" min-width="200">
          <template slot-scope="scope">
            {{ scope.row.created }}
          </template>
        </el-table-column>

      </el-table>
      <div class="pagination-container">
        <el-pagination
          background
          :small="device === 'mobile'"
          :current-page="input.page"
          :page-size="input.limit"
          :layout="getLayout"
          :total="count"
          @current-change="getList"
          @size-change="handleSizeChange"
        />
      </div>
    </el-card>
  </div>
</template>

<script>
import { UrlConfig } from '@/api/api-rbac'
import { getList } from '@/api/operate.js'
import { roleOption, userList } from '@/api/user'
import Vuex from 'vuex'

export default {
  components: {
    Date: () => import('@/components/Date'),
    JsonEditor: () => import('@/components/JsonEditor/index')
  },
  data() {
    return {
      input: {
        page: 1,
        limit: 10,
        operater_action: '',
        operater_id: 0,
        operater_role_id: 0,
        date: []
      },
      count: 0,
      logList: [],
      tableLoading: false,
      urlConfig: [],
      urlConfigMap: [],
      roleConfig: [],
      roleConfigMap: [],
      userConfig: []
    }
  },
  computed: {
    ...Vuex.mapGetters([
      'device',
    ]),
    getLayout(){
      return this.device !== 'mobile' ? 'total, sizes, prev, pager, next, jumper' : 'prev, pager, next'
    }
  },
  created() {
    this.initList()
  },

  methods: {
    lookData(index, typ) {
      for (const i in this.logList) {
        if (this.logList[i].index == index) {
          this.logList[i].isFormatData = typ
        }
      }
    },
    changeDate(v) {
      this.input.date = v
      this.getList(1)
    },
    async initList() {
      const roleOptionRes = await roleOption()
      if (roleOptionRes) {
        for (var v of roleOptionRes.data) {
          this.roleConfig.push(v)
          this.roleConfigMap[v['id']] = v['name']
        }
      }
      const urlConfigRes = await UrlConfig()
      if (urlConfigRes) {
        for (var v of urlConfigRes.data) {
          this.urlConfig.push(v)
          this.urlConfigMap[v['url']] = v['remark']
        }
      }
      const userListRes = await userList()
      if (userListRes) {
        for (var v of userListRes.data) {
          this.userConfig.push(v)
        }
      }
      this.getList()
    },
    handleSizeChange(v) {
      this.input.limit = v
      this.getList(1)
    },
    getUrlOpt() {
      UrlConfig().then(res => {
        if (res) {
          for (var v of res.data) {
            this.urlConfig.push(v)
            this.urlConfigMap[v['url']] = v['remark']
          }
        }
      })
    },
    getRoleOpt() {
      roleOption().then(res => {
        if (res) {
          for (var v of res.data) {
            this.roleConfig.push(v)
            this.roleConfigMap[v['id']] = v['name']
          }
        }
      })
    },
    getUserOpt() {
      userList().then(res => {
        if (res) {
          for (var v of res.data) {
            this.userConfig.push(v)
          }
        }
      })
    },

    getList(page) {
      !page ? this.input.page = 1 : this.input.page = page
      this.tableLoading = true

      if (this.input.operater_id == '') {
        this.input.operater_id = 0
      }

      if (this.input.operater_role_id == '') {
        this.input.operater_role_id = 0
      }

      getList(this.input).then(res => {
        if (res.code == 0) {
          if (res.data.list == null) {
            res.data.list = []
          }
          const list = []
          let index = 0
          for (const v of res.data.list) {
            if (v['body_str'] != '') {
              v['dataFormat'] = JSON.stringify(JSON.parse(v['body_str']), null, '\t')
            }
            v['isFormatData'] = false
            v['index'] = index
            list.push(v)
            index++
          }
          this.logList = list

          this.tableLoading = false
          this.count = Number(res.data.count)
        } else {
          this.$message({
            offset: 60,

            type: 'error',
            message: res.msg
          })
        }
        this.tableLoading = false
      }).catch(err => {
        this.tableLoading = false
        console.log('err', err)
      })
    }
  }
}
</script>

<style lang="scss" scoped>

</style>

