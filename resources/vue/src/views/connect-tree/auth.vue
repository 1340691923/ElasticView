<template>
  <div class="app-container">

    <div class="search-container">
      <el-form :inline="true">
        <el-form-item label="" >
          <el-button

            type="success"
            class="filter-item"
            @click="getList"
          >{{ $t('搜索') }}
          </el-button>
        </el-form-item>
        <el-form-item label="" >
          <el-button

            type="primary"
            class="filter-item"
            @click="handleAddEsCfg"
          >{{ $t('新建鉴权信息') }}
          </el-button>
        </el-form-item>
      </el-form>
    </div>
    <el-card shadow="never" class="table-container">
      <el-table
      v-if="refreshTable"
      v-loading="getListLoading"
      :data="list"

    >
      <el-table-column
        :label="$t('备注')"
        align="center"
      >
        <template #default="scope">
          {{ scope.row.remark }}
        </template>
      </el-table-column>
      <el-table-column align="center" :label="$t('已分配角色')" width="200">
        <template #default="scope">
          <template v-for="(v,index) in  scope.row.share_roles">
            <el-tag>{{roleMap[v]}}</el-tag>
          </template>
        </template>
      </el-table-column>

      <el-table-column align="center" :label="$t('用户名')" width="200">
        <template #default="scope">
          {{ scope.row.user }}
        </template>
      </el-table-column>
<!--      <el-table-column align="center" :label="$t('密码')" width="200">
        <template #default="scope">
          {{ scope.row.pwd }}
        </template>
      </el-table-column>-->
      <el-table-column align="center" :label="$t('root证书')" width="100" show-overflow-tooltip>
        <template #default="scope">
          {{ scope.row.rootpem }}
        </template>
      </el-table-column>
      <el-table-column align="center" :label="$t('cert证书')" width="100" show-overflow-tooltip>
        <template #default="scope">
          {{ scope.row.certpem }}
        </template>
      </el-table-column>

      <el-table-column align="center" :label="$t('key证书')" width="100" show-overflow-tooltip>
        <template #default="scope">
          {{ scope.row.keypem }}
        </template>
      </el-table-column>

      <el-table-column align="center" :label="$t('创建时间')" width="220">
        <template #default="scope">
          {{ scope.row.created }}
        </template>
      </el-table-column>
      <el-table-column align="center" :label="$t('修改时间')" width="220">
        <template #default="scope">
          {{ scope.row.updated }}
        </template>
      </el-table-column>
      <el-table-column align="center" :label="$t('操作')" fixed="right" width="200">
        <template #default="scope2">
          <el-button
            type="primary"
            @click="handleEditEsCfg(scope2)"
          >{{ $t('编辑') }}
          </el-button>
          <el-button
            type="danger"
            @click="deleteEsCfgAction(scope2)"
          >{{ $t('删除') }}
          </el-button>
        </template>
      </el-table-column>
    </el-table>
    </el-card>

    <el-drawer
      size="80%"
      v-model="dialogAddCfgVisible"
      :title="$t('新建鉴权配置')"
    >
      <el-form :model="addEsCfgForm" label-width="100px" label-position="left">
        <el-form-item :label="$t('备注')">
          <el-input clearable v-model="addEsCfgForm.remark" :placeholder="$t('备注')" />
        </el-form-item>
        <el-form-item :label="$t('可访问集群成员')">
          <el-transfer
            @change="changeShareRoles"
            v-if="addEsCfgForm"
            v-model="addEsCfgForm.share_roles"
            :titles="['全部角色', '当前集群角色']"
            :button-texts="['移除角色', '添加角色']"
            filterable
            :filter-method="filterMethod"
            filter-placeholder="请操作角色"
            :data="allRoleConfig"
          />
        </el-form-item>
        <el-form-item :label="$t('用户名')">
          <el-autocomplete
            v-model="addEsCfgForm.user"
            clearable
            :fetch-suggestions="querySearch"
            :placeholder="$t('用户名')"
          >
            <template #default="{ item }">
              <span>{{ item.value }}</span>
            </template>

          </el-autocomplete>

        </el-form-item>
        <el-form-item :label="$t('密码')">
          <el-input show-password type="password" clearable v-model="addEsCfgForm.pwd" :placeholder="$t('密码')" />
        </el-form-item>
        <el-form-item :label="$t('root证书')">
          <el-input clearable v-model="addEsCfgForm.rootpem" type="textarea" :placeholder="$t('root证书')" />
        </el-form-item>
        <el-form-item :label="$t('cert证书')">
          <el-input clearable v-model="addEsCfgForm.certpem" type="textarea" :placeholder="$t('cert证书')" />
        </el-form-item>
        <el-form-item :label="$t('key证书')">
          <el-input clearable v-model="addEsCfgForm.keypem" type="textarea" :placeholder="$t('key证书')" />
        </el-form-item>

        <el-form-item :label="$t('自定义请求头')">
          <es-header style="width: 100%" v-model="addEsCfgForm.header" ></es-header>
        </el-form-item>


      </el-form>
      <template #footer>
        <el-button

          type="danger"

          @click="dialogAddCfgVisible=false"
        >{{ $t('取消') }}
        </el-button>
        <el-button

          type="primary"

          @click="confirmEsCfg"
        >{{ $t('确认') }}
        </el-button>
      </template>
    </el-drawer>
    <el-drawer
      size="80%"

      v-model="dialogEditCfgVisible"
      :title="$t('编辑鉴权配置')"
    >
      <el-form :model="editEsCfgForm" label-width="200px" label-position="left">
        <el-form-item :label="$t('备注')">
          <el-input clearable v-model="editEsCfgForm.remark" :placeholder="$t('备注')" />
        </el-form-item>
        <el-form-item :label="$t('用户名')">
          <el-autocomplete
            v-model="editEsCfgForm.user"
            clearable
            :fetch-suggestions="querySearch"
            :placeholder="$t('用户名')"
          >
            <template #default="{ item }">
              <span>{{ item.value }}</span>
            </template>

          </el-autocomplete>

        </el-form-item>
        <el-form-item :label="$t('密码')">
          <el-input show-password type="password" clearable v-model="editEsCfgForm.pwd" :placeholder="$t('密码')" />
        </el-form-item>
        <el-form-item :label="$t('可访问集群成员')">
          <el-transfer
            @change="changeShareRoles"
            v-if="dialogEditCfgVisible"
            v-model="editEsCfgForm.share_roles"
            :titles="['全部角色', '当前集群角色']"
            :button-texts="['移除角色', '添加角色']"
            filterable
            :filter-method="filterMethod"
            filter-placeholder="请操作角色"
            :data="allRoleConfig"
          />

        </el-form-item>

        <el-form-item :label="$t('root证书')">
          <el-input clearable v-model="editEsCfgForm.rootpem" type="textarea" :placeholder="$t('root证书')" />
        </el-form-item>
        <el-form-item :label="$t('cert证书')">
          <el-input clearable v-model="editEsCfgForm.certpem" type="textarea" :placeholder="$t('cert证书')" />
        </el-form-item>
        <el-form-item :label="$t('key证书')">
          <el-input clearable v-model="editEsCfgForm.keypem" type="textarea" :placeholder="$t('key证书')" />
        </el-form-item>
        <el-form-item :label="$t('自定义请求头')">
          <es-header style="width: 100%" v-model="editEsCfgForm.header" ></es-header>
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button

          type="danger"

          @click="dialogEditCfgVisible=false"
        >{{ $t('取消') }}
        </el-button>
        <el-button

          type="primary"

          @click="confirmEditEsCfg"
        >{{ $t('确认') }}
        </el-button>

      </template>
    </el-drawer>


  </div>
</template>

<script>
import { deepClone } from '@/utils'
import { DeleteAction, InsertAction, GetEsCfgList, UpdateAction,
  UpdateEsCfgAction,InsertEsCfgAction,DeleteEsCfgAction,
} from '@/api/es-link'
import { PingAction } from '@/api/es'
import {roleOption} from "@/api/user";
import EsHeader from '@/views/connect-tree/EsHeader.vue'
import { View } from '@element-plus/icons-vue'

const defaultLink = {
  id: 0,
  ip: 'http://127.0.0.1:9200',
  remark: '',
  version: 6
}

export default {
  components:{
    EsHeader
  },
  data() {
    return {
      managerFormdialogVisible:false,
      refreshTable:true,
      default_expand_all:false,
      getListLoading:false,
      usernameWord:[
        { "value": "elastic" },
      ],
      testConnectLoading: false,
      connectLoading: false,
      link: Object.assign({}, defaultLink),
      list: [],
      dialogVisible: false,
      dialogType: 'new',
      dialogAddCfgVisible:false,
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
        header:[]
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
        share_roles:[],
        header:[]
      },
      selectEsLinkId:0,
      selectEsVersion:0,
      selectEsIp:'',
      roleMap:{},
      roleList:[],
      allRoleConfig:[]
    }
  },
  setup(){
    return {View}
  },
  computed: {

  },
  async created() {
    await this.initAllRoles()
    this.getList()
  },
  methods: {
    filterMethod(query, item) {
      return item.label.indexOf(query) > -1
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
    testConnectForm(form) {
      this.testConnectLoading = true
      form["ip"] = this.selectEsIp
      form["version"] = this.selectEsVersion
      PingAction(form).then(res => {
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
        this.testConnectLoading = false
      }).catch(err => {
        this.testConnectLoading = false
      })
    },
    testConnect(scope) {
      this.list[scope.$index].connectLoading = true

      PingAction({
        ip:scope.row.ip,
        user:scope.row.user,
        pwd:scope.row.pwd,
        version:scope.row.version,
        rootpem:scope.row.rootpem,
        certpem:scope.row.certpem,
        keypem:scope.row.keypem,
      }).then(res => {
        if (res.code == 0) {
          console.log('res', res)
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
        this.list[scope.$index].connectLoading = false
      }).catch(err => {
        this.list[scope.$index].connectLoading = false
      })
    },
    async getList() {
      this.getListLoading = true
      const res = await GetEsCfgList()
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

      this.list = res.data
      this.refreshTable = false
      this.$nextTick(() => {
        this.refreshTable = true
      })
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
    handleAddEsCfg() {
      this.addEsCfgForm.share_roles = []
      this.allRoleConfig = []
      for (let v of this.roleList) {

        const obj = {
          label: v.name,
          key: v.id.toString(),
          disabled: false
        }

        this.allRoleConfig.push(
          obj
        )
      }
      this.dialogAddCfgVisible = true
    },
    handleEditEsCfg(cfgScope) {
      this.editEsCfgForm.id = cfgScope.row.id
      this.editEsCfgForm.user = cfgScope.row.user
      this.editEsCfgForm.pwd = cfgScope.row.pwd
      this.editEsCfgForm.remark = cfgScope.row.remark
      this.editEsCfgForm.certpem = cfgScope.row.certpem
      this.editEsCfgForm.keypem = cfgScope.row.keypem
      this.editEsCfgForm.rootpem = cfgScope.row.rootpem
      this.editEsCfgForm.share_roles = cfgScope.row.share_roles
      this.editEsCfgForm.header = cfgScope.row.header
      this.allRoleConfig = []
      for (let v of this.roleList) {

        const obj = {
          label: v.name,
          key: v.id.toString(),
          disabled: false
        }


        this.allRoleConfig.push(
          obj
        )
      }

      this.dialogEditCfgVisible = true
    },
    async confirmEditEsCfg(){
      let form = JSON.parse(JSON.stringify(this.editEsCfgForm))
      let res = await UpdateEsCfgAction(form)
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
      this.dialogEditCfgVisible = false
    },
    async confirmEsCfg(){
      let form = JSON.parse(JSON.stringify(this.addEsCfgForm))
      let res = await InsertEsCfgAction(form)
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
      this.dialogAddCfgVisible = false
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
