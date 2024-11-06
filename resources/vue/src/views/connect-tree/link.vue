<template>
  <div class="app-container">
    <div class="search-container">
      <el-form :inline="true">
        <el-form-item label="" >
          <el-button
            type="success"
            class="filter-item"
            @click="getList"
          >{{$t('搜索')}}
          </el-button>
        </el-form-item>


        <el-form-item label="" >
          <el-button

            type="primary"
            class="filter-item"
            @click="handleAddRole"
          >{{ $t('新建连接信息') }}
          </el-button>
        </el-form-item>

        <el-form-item :label=" $t('ES鉴权配置')" >

          <el-switch
            @change="changeExpand"
            v-model="default_expand_all"
            active-text="展开"
            inactive-text="折叠">
          </el-switch>
        </el-form-item>
      </el-form>
    </div>
    <el-card shadow="never" class="table-container">

      <el-table
        :row-key="getRowKey"
        :expand-row-keys="expandRowKeys"
        @expand-change="handleExpandChange"
      v-loading="getListLoading"
      :data="list"
    >
      <el-table-column width="100"  :label="$t('鉴权配置')" type="expand">
        <template #default="props">
          <el-table style="width: 100%" :data="props.row.es_link_configs">
            <el-table-column
              :label="$t('备注')"
              align="center"

            >
              <template #default="scope">
                {{ scope.row.remark }}
              </template>
            </el-table-column>
            <el-table-column align="center" :label="$t('已分配角色')" width="400">
              <template #default="scope">
                <template v-for="(v,index) in  scope.row.share_roles">
                  <el-tag>{{roleMap[v]}}</el-tag>
                </template>
              </template>
            </el-table-column>

            <el-table-column align="center" :label="$t('用户名')" width="300">
              <template #default="scope">
                {{ scope.row.user }}
              </template>
            </el-table-column>

            <el-table-column align="center" :label="$t('root证书')" width="120" show-overflow-tooltip>
              <template #default="scope">
                {{ scope.row.rootpem }}
              </template>
            </el-table-column>
            <el-table-column align="center" :label="$t('cert证书')" width="120" show-overflow-tooltip>
              <template #default="scope">
                {{ scope.row.certpem }}
              </template>
            </el-table-column>
            <el-table-column align="center" :label="$t('key证书')" width="120" show-overflow-tooltip>
              <template #default="scope">
                {{ scope.row.keypem }}
              </template>
            </el-table-column>
            <el-table-column align="center" :label="$t('操作')" fixed="right" width="300">
              <template #default="scope2">

                <el-button
                  v-loading="scope2.row.connectLoading"
                  :disabled="scope2.row.connectLoading"
                  type="success"

                  @click="testConnect(scope2)"
                >{{ $t('ping') }}
                </el-button>

                <el-button
                  type="danger"

                  @click="deleteEsCfgRelation(scope2)"
                >{{ $t('删除') }}
                </el-button>
              </template>
            </el-table-column>
          </el-table>
        </template>
      </el-table-column>
      <el-table-column
        :label="$t('备注')"
        align="center"
        width="100"
      >
        <template #default="scope">
          {{ scope.row.remark }}
        </template>
      </el-table-column>
      <el-table-column align="center" :label="$t('连接地址')" >
        <template #default="scope">
          {{ scope.row.ip }}
        </template>
      </el-table-column>
      <el-table-column align="center" :label="$t('数据源类型')" width="100">
        <template #default="scope">
          {{ scope.row.version }}
        </template>
      </el-table-column>
      <el-table-column align="center" :label="$t('创建时间')" >
        <template #default="scope">
          {{ scope.row.created }}
        </template>
      </el-table-column>
      <el-table-column align="center" :label="$t('修改时间')" >
        <template #default="scope">
          {{ scope.row.updated }}
        </template>
      </el-table-column>
      <el-table-column align="center"   fixed="right" :label="$t('操作')" width="200">
        <template #default="scope">

          <el-button
            type="primary"

            @click="handleEdit(scope)"
          >{{ $t('编辑') }}
          </el-button>
          <el-button
            type="danger"

            @click="handleDelete(scope)"
          >{{ $t('删除') }}
          </el-button>
        </template>
      </el-table-column>
    </el-table>
    </el-card>
    <el-drawer
      size="80%"

      v-model="dialogVisible"
      :title="dialogType==='edit'?$t('编辑连接信息'):$t('新建连接信息')"
    >
      <el-form :model="link" label-width="200px" label-position="left">
        <el-form-item :label="$t('备注')">
          <el-input clearable v-model="link.remark" :placeholder="$t('备注')" />
        </el-form-item>
        <el-form-item :label="$t('连接地址')">
          <el-input v-model="link.ip" :placeholder="$t('例如:http://127.0.0.1:9200')" />
        </el-form-item>
        <el-form-item :label="$t('数据源类型')">
          <el-select v-model="link.version" :placeholder="$t('请选择数据源类型')" filterable>
            <el-option label="elasticsearch6.x" value="elasticsearch6.x" />
            <el-option label="elasticsearch7.x" value="elasticsearch7.x" />
            <el-option label="elasticsearch8.x" value="elasticsearch8.x" />


            <el-option label="mysql" value="mysql" />
            <el-option label="redis" value="redis" />
            <el-option label="clickhouse" value="clickhouse" />
            <el-option label="postgres" value="postgres" />
            <el-option label="mongo" value="mongo" />
          </el-select>
        </el-form-item>
        <el-form-item :label="$t('鉴权列表')">
          <el-select
            filterable
            v-model="link.cfgIds"
            style="width: 150px"
            collapse-tags
            multiple
            :placeholder="$t('请选择')"
          >
            <el-option v-for="(v,k,index) in esCfgOptlist" :key="index" :label="v.remark" :value="v.id" />
          </el-select>
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button
          type="danger"
          @click="dialogVisible=false"
        >{{ $t('取消') }}
        </el-button>
        <el-button
          type="primary"
          @click="confirm"
        >{{ $t('确认') }}
        </el-button>
      </template>
    </el-drawer>
  </div>
</template>

<script>
import { DeleteAction, InsertAction, ListAction, UpdateAction,
  GetEsCfgOpt,DeleteEsCfgRelation,DeleteEsCfgAction,
} from '@/api/es-link'
import { PingAction } from '@/api/es'

import {roleOption} from "@/api/user";

const defaultLink = {
  id: 0,
  ip: 'http://127.0.0.1:9200',
  remark: '',
  version: 'elasticsearch8.x',
  cfgIds:[],
}

export default {
  data() {
    return {
      esCfgOptlist:[],
      managerFormdialogVisible:false,
      default_expand_all:false,
      getListLoading:false,
      usernameWord:[
        { "value": "elastic" },
      ],
      testConnectLoading: false,
      connectLoading: false,
      loadingMap:{},
      link: Object.assign({}, defaultLink),
      list: [],
      dialogVisible: false,
      dialogType: 'new',
      dialogEditCfgVisible:false,
      addEsCfgForm:{
        es_link_id:0,
        user: '',
        pwd: '',
        remark: '',
        rootpem: '',
        certpem: '',
        keypem: '',
        share_roles:[],
      },
      editEsCfgForm:{
        linkId:0,
        id:0,
        user: '',
        pwd: '',
        remark: '',
        rootpem: '',
        certpem: '',
        keypem: '',
        share_roles:[]
      },
      selectEsLinkId:0,
      selectEsVersion:'elasticsearch6.x',
      selectEsIp:'',
      roleMap:{},
      roleList:[],
      allRoleConfig:[],
      expandRowKeys: [],
    }
  },
  async created() {
    await this.getEsCfgOpt()
    await this.initAllRoles()
    this.getList()
  },

  methods: {
    changeExpand(){
      if(this.default_expand_all) {
        this.expandRowKeys = this.list.map(item => item.id)
      }else{
        this.expandRowKeys = []
      }

    },
    //获取行的唯一标识符（key）
    getRowKey(row) {
      return row.id;
    },
    // 处理行展开事件
    handleExpandChange(row, expandedRows) {
      this.expandRowKeys = expandedRows.map(item => item.id)
      let expandRowKeys = []
      for(let i in this.expandRowKeys){
        expandRowKeys.push(this.expandRowKeys[i])
      }
      if(expandRowKeys.length  == this.list.map(item => item.id).length ){
        this.default_expand_all = true
      }
      if(expandRowKeys.length == 0 ){
        this.default_expand_all = false
      }

    },
    async deleteEsCfgRelation(scope2 ){
      const res = await DeleteEsCfgRelation({id:scope2.row.cfg_relation_id})
      if (res.code != 0) {
        ElMessage.error({
          type: 'error',
          message: res.msg
        })
        return
      }
      this.getList()
    },
    filterMethod(query, item) {
      return item.label.indexOf(query) > -1
    },
    async getEsCfgOpt(){
      const res = await GetEsCfgOpt()
      if (res.code != 0) {
        ElMessage.error({
          type: 'error',
          message: res.msg
        })
        return
      }
      if(res.data==null) res.data = []
      this.esCfgOptlist = res.data
    },
    async initAllRoles(){
      const res = await roleOption()
      if (res.code != 0) {
        ElMessage.error({
          type: 'error',
          message: res.msg
        })
        return
      }
      if(res.data==null) res.data = []
      this.roleMap = {}
      for(let v of res.data){
        this.roleMap[v.id] = v.name
      }
      this.roleList = res.data
    },
    testConnect(scope) {
      scope.row.connectLoading = true
      PingAction({
        id:scope.row.cfg_relation_id,
        ip:scope.row.ip,
        user:scope.row.user,
        pwd:scope.row.pwd,
        version:scope.row.version,
        rootpem:scope.row.rootpem,
        certpem:scope.row.certpem,
        keypem:scope.row.keypem,
        header:scope.row.header,
      }).then(res => {
        scope.row.connectLoading = false
        if (res.code == 0) {
          ElMessage.success({
            type: 'success',
            message: `连接成功`
          })
        } else {
          ElMessage.error({
            type: 'error',
            message: res.msg
          })
        }
      }).catch(err => {
        scope.row.connectLoading = false
        console.log("err",err)
      })

    },
    async getList() {
      this.getListLoading = true
      const res = await ListAction()
      this.getListLoading = false
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

      for (const k in res.data) {
        if(res.data[k].es_link_configs == null){
          res.data[k].es_link_configs = []
        }
        for (const k2 in res.data[k].es_link_configs) {
          res.data[k].es_link_configs[k2]['connectLoading'] = false
        }
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

      this.link.id = scope.row.id
      this.link.remark =  scope.row.remark
      this.link.version = scope.row.version
      this.link.ip = scope.row.ip
      this.link.cfgIds = []

      for(let v of scope.row.es_link_configs){
        this.link.cfgIds.push(v.id)
      }
    },
    async deleteEsCfgAction(scope){
      let res = await DeleteEsCfgAction({id:scope.row.id})
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
      this.getList()
    },
    handleDelete({ $index, row }) {
      ElMessageBox.confirm('确定删除该连接信息吗?', '警告', {
        confirmButtonText: '确认',
        cancelButtonText: '取消',
        type: 'warning'
      })
        .then(async() => {
          const res = await DeleteAction({ id: row.id })
          if (res.code != 0) {
            ElMessage.error({
              type: 'error',
              message: res.msg
            })
            return
          }
          this.getList()
        })
        .catch(err => {
          console.error(err)
        })
    },
    async confirm() {

      const isEdit = this.dialogType === 'edit'

      if (isEdit) {
        const res = await UpdateAction(this.link)
        if (res.code != 0) {
          ElMessage.error({
            type: 'error',
            message: res.msg
          })
          return
        }
        this.getList()
      } else {
        const res = await InsertAction(this.link)
        if (res.code != 0) {
          ElMessage.error({
            type: 'error',
            message: res.msg
          })
          return
        }
        this.getList()
      }

      this.dialogVisible = false
      ElMessage.success({
        type: 'success',
        message: isEdit ? '修改成功' : '创建成功'
      })
    },
  }
}
</script>

<style lang="scss" scoped>

</style>
