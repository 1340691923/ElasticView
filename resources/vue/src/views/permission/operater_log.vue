<template>
  <div class="app-container">
    <div class="search-container">
      <el-form :inline="true">
        <el-form-item :label="$t('接口名')" >
          <el-select
            v-model="input.operater_action"
            reserve-keyword
            collapse-tags
            :placeholder="$t('接口名')"
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
        </el-form-item>

        <el-form-item :label=" $t('用户')" >
          <el-select
            v-model="input.operater_id"
            :placeholder="  $t('请选择用户')"
            clearable
            style="width: 120px"
            filterable
            @change="getList(1)"
          >
            <el-option :label="$t('请选择用户')" :value="Number(0)" />
            <el-option
              v-for="item in userConfig"
              :key="item.id"
              :label="item.username"
              :value="item.id"
            />
          </el-select>
        </el-form-item>
        <el-form-item  :label="$t('时间')" >
          <date v-model="input.date" class="filter-item"  />
        </el-form-item>
        <el-form-item label="" >
          <el-button type="primary"  @click="getList(1)"> {{$t('查询')}}</el-button>
        </el-form-item>
      </el-form>

    </div>
    <el-card shadow="never" class="table-container">
      <el-table
      v-loading="tableLoading"
      :data="logList"
      max-height="800"
      :element-loading-text="`${$t('请给我点时间')}！`"
    >

        <el-table-column type="expand">
          <template #default="props">

            <json-viewer
              boxed
              copyable
              :expanded="false"
              theme="my-json-theme"
              :value="props.row.dataFormat"
            ></json-viewer>
          </template>
        </el-table-column>


      <el-table-column align="center" :label="$t('操作用户')" min-width="100">
        <template #default="scope">
          {{ scope.row.operater_name }}
        </template>
      </el-table-column>

      <el-table-column align="center" :label="$t('操作接口')" min-width="220">
        <template #default="scope">
          {{ urlConfigMap[scope.row.operater_action] }}
        </template>
      </el-table-column>

        <el-table-column align="center" :label="$t('所花时间')" min-width="200">
          <template #default="scope">
            {{ scope.row.cost_time }}
          </template>


        </el-table-column>
        <el-table-column align="center" :label="$t('状态')" min-width="200">
          <template #default="scope">
            {{ scope.row.status }}
          </template>


        </el-table-column>
      <el-table-column align="center" :label="$t('创建时间')" min-width="200">
        <template #default="scope">
          {{ scope.row.created }}
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
        @current-change="getList"
        @size-change="handleSizeChange"
      />
    </div>
  </div>
</template>

<script>
import { UrlConfig } from '@/api/api-rbac'
import { getList } from '@/api/operate'
import {  userList } from '@/api/user'
import dayjs from "dayjs";

export default {
  components: {
    Date: () => import('@/components/Date/index.vue'),
  },
  data() {
    return {
      input: {
        page: 1,
        limit: 10,
        operater_action: '',
        operater_id: 0,

        date: [
          dayjs().format('YYYY-MM-DD 00:00:00'), dayjs().format('YYYY-MM-DD 23:59:59')
        ]
      },
      count: 0,
      logList: [],
      tableLoading: false,
      urlConfig: [],
      urlConfigMap: [],

      userConfig: []
    }
  },
  computed: {


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

      const urlConfigRes = await UrlConfig()
      if (urlConfigRes) {
        for (var v of urlConfigRes.data.cfg) {
          this.urlConfig.push(v)
          this.urlConfigMap[v['url']] = v['remark']
        }
      }
      const userListRes = await userList()
      if (userListRes) {
        for (var v of userListRes.data.list) {
          this.userConfig.push(v)
        }
      }
      this.getList(1)
    },
    handleSizeChange(v) {
      this.input.limit = v
      this.getList(1)
    },
    getUrlOpt() {
      UrlConfig().then(res => {
        if (res) {
          for (var v of res.data.cfg) {
            this.urlConfig.push(v)
            this.urlConfigMap[v['url']] = v['remark']
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


      this.input.date[0] = dayjs(this.input.date[0]).format('YYYY-MM-DD HH:mm:ss')
      this.input.date[1] = dayjs(this.input.date[1]).format('YYYY-MM-DD HH:mm:ss')

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
          ElMessage.error({
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

