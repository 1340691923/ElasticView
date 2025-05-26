<template>
  <div class="app-container">

    <div class="search-container">
      <el-form :inline="true">
        <el-form-item label="" >
          <el-button

            type="success"
            class="filter-item"
            @click="getList(1)"
          >{{ $t('搜索') }}
          </el-button>
        </el-form-item>
        <el-form-item label="" >
          <el-button

            type="primary"
            class="filter-item"
            @click="handleAddEsCfg"
          >{{ $t('新建数据源鉴权信息') }}
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
      <el-table-column align="center" :label="$t('分配权限组')" width="140">
        <template #default="scope">
          <div class="role-tags">
            <template v-if="scope.row.share_roles.length <= 2">
              <el-tag v-for="item in scope.row.share_roles" :key="item">
                {{roleMap[item]}}
              </el-tag>
            </template>
            <template v-else>
              <el-tag v-for="item in scope.row.share_roles.slice(0, 1)" :key="item">
                {{roleMap[item]}}
              </el-tag>
              <el-popover
                placement="bottom"
                trigger="hover"
                :width="200"
                popper-class="role-popover"
              >
                <template #reference>
                  <el-tag class="more-tag">
                    +{{ scope.row.share_roles.length - 1 }}
                  </el-tag>
                </template>
                <div class="popover-tags">
                  <el-tag
                    v-for="item in scope.row.share_roles.slice(1)"
                    :key="item"
                    style="margin-left:1rem"
                    size="small"
                  >
                    {{roleMap[item]}}
                  </el-tag>
                </div>
              </el-popover>
            </template>
          </div>

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
      <el-table-column align="center" :label="$t('操作')" fixed="right" width="130">
        <template #default="scope2">
          <el-button
            type="primary"
            @click="handleEditEsCfg(scope2)"
            :icon="Edit"
          >
          </el-button>
          <el-button
            type="danger"
            @click="deleteEsCfgAction(scope2)"
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
        :layout="isMobile?'pager':'total, sizes, prev, pager, next, jumper'"
        @current-change="getList"
        @size-change="handleSizeChange"
      />
    </div>
    <el-drawer
      :size="isMobile?'100%':'50%'"
      v-model="dialogAddCfgVisible"
      :title="$t('新建鉴权配置')"
    >
      <el-form :model="addEsCfgForm" label-width="100px" label-position="left">
        <el-form-item :label="$t('备注')">
          <el-input clearable v-model="addEsCfgForm.remark" :placeholder="$t('备注')" />
        </el-form-item>
        <el-form-item :label="$t('可访问权限组')">

          <el-select
            multiple
            v-model="addEsCfgForm.share_roles"
            reserve-keyword
            collapse-tags
            :placeholder="$t('可访问权限组')"

            class="filter-item"
            filterable
          >

            <el-option
              v-for="item in allRoleConfig"
              :key="item.key"
              :label="item.label"
              :value="item.key"
            />
          </el-select>

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
      :size="isMobile?'100%':'50%'"

      v-model="dialogEditCfgVisible"
      :title="$t('编辑鉴权配置')"
    >
      <el-form :model="editEsCfgForm" label-width="100px" label-position="left">
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
        <el-form-item :label="$t('可访问权限组')">

          <el-select
            multiple
            v-model="editEsCfgForm.share_roles"
            reserve-keyword
            collapse-tags
            :placeholder="$t('可访问权限组')"

            class="filter-item"
            filterable
          >

            <el-option
              v-for="item in allRoleConfig"
              :key="item.key"
              :label="item.label"
              :value="item.key"
            />
          </el-select>
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
import {Delete, Edit, View} from '@element-plus/icons-vue'
import {useAppStore} from "@/store";
import {DeviceEnum} from "@/enums/DeviceEnum";
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
  setup(){
    const appStore = useAppStore()
    const { device } = storeToRefs(appStore)
    return {Edit,Delete,device,View}
  },
  computed: {
    isMobile(){
      return this.device === DeviceEnum.MOBILE
    },
  },
  data() {
    return {
      input: {
        page: 1,
        limit: 10,
      },
      count: 0,
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

  async created() {
    await this.initAllRoles()
    this.getList(1)
  },
  methods: {
    filterMethod(query, item) {
      return item.label.indexOf(query) > -1
    },
    async initAllRoles(){
      const res = await roleOption()
      if (res.code != 0) {
        ElMessage.error({
          type: 'error',offset:60,
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
            type: 'success',offset:60,
            message: `连接成功`
          })
        } else {
          ElMessage.error({
            type: 'error',offset:60,
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
            type: 'success',offset:60,
            message: `连接成功`
          })
        } else {
          ElMessage.error({
            type: 'error',offset:60,
            message: res.msg
          })
        }
        this.list[scope.$index].connectLoading = false
      }).catch(err => {
        this.list[scope.$index].connectLoading = false
      })
    },
    handleSizeChange(v) {
      this.input.limit = v
      this.getList(1)
    },
    async getList(page) {
      !page ? this.input.page = 1 :  this.input.page = page
      this.getListLoading = true
      const res = await GetEsCfgList({
        page:page,
        page_size: this.input.limit
      })
      this.getListLoading = false
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
      this.count = res.data.count
      this.list = res.data.list
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
          type: 'error',offset:60,
          message: res.msg
        })
        return
      }
      ElMessage.success({
        type: 'success',offset:60,
        message: res.msg
      })
      this.getList(1)
      this.dialogEditCfgVisible = false
    },
    async confirmEsCfg(){
      let form = JSON.parse(JSON.stringify(this.addEsCfgForm))
      let res = await InsertEsCfgAction(form)
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

      this.getList(1)
      this.dialogAddCfgVisible = false
    },
    async deleteEsCfgAction(scope){

      ElMessageBox.confirm('确定删除该鉴权配置吗?', '警告', {
        confirmButtonText: '确认',
        cancelButtonText: '取消',
        type: 'warning'
      })
        .then(async() => {
          let res = await DeleteEsCfgAction({id:scope.row.id})
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
          this.getList(1)
        })
        .catch(err => {
          console.error(err)
        })


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
              type: 'error',offset:60,
              message: res.msg
            })
            return
          }
          this.getList(1)
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
            type: 'error',offset:60,
            message: res.msg
          })
          return
        }
        this.getList(1)
      } else {
        const res = await InsertAction(this.link)
        if (res.code != 0) {
          ElMessage.error({
            type: 'error',offset:60,
            message: res.msg
          })
          return
        }
        this.getList(1)
      }

      this.dialogVisible = false
      ElMessage.success({
        type: 'success',offset:60,
        message: isEdit ? '修改成功' : '创建成功'
      })
    },
  }
}
</script>


<style lang="scss" scoped>
.app-container {
.search-container {
  margin-bottom: 1.5rem;
  padding: 1rem;
  backdrop-filter: blur(8px);
  border-radius: 0.5rem;
  transition: all 0.3s;

.search-header {
  display: none;
}

@media (max-width: 768px) {
  .search-header {
    display: flex;
    justify-content: space-between;
    align-items: center;
    padding: 0.5rem 0;
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

&:last-child {
   display: flex;
   gap: 0.5rem;

.el-button {
  flex: 1;
  margin: 0 !important;
}
}
}
}
}

:deep(.el-form) {
  display: flex;
  flex-wrap: wrap;
  align-items: flex-end;
  gap: 1rem;

.el-form-item {
  margin-bottom: 0;

&__label {
   color: #4b5563;
 }
}

.el-input__wrapper,
.el-select .el-input__wrapper {
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

.el-button {
  box-shadow: 0 1px 2px rgba(0, 0, 0, 0.05);
  transition: all 0.3s;

&:hover {
   transform: scale(1.05);
   box-shadow: 0 4px 6px rgba(0, 0, 0, 0.1);
 }

&.filter-item {
   min-width: 80px;
 }
}
}
}

.table-container {
  border-radius: 0.5rem;
  overflow: hidden;
  backdrop-filter: blur(8px);
  transition: all 0.3s;
  box-shadow: 0 1px 2px rgba(0, 0, 0, 0.05);

:deep(.el-table) {
  background-color: transparent;

.el-table__row {
  transition: all 0.3s;

&:hover {
   transform: translateY(-1px);
   box-shadow: 0 2px 4px rgba(0, 0, 0, 0.05);
 }
}

th {
  color: #4b5563;
  font-weight: 500;
  border-bottom: 1px solid #e5e7eb;
  padding: 12px 0;

.cell {
  font-size: 0.95rem;
}
}

td {
  border-bottom: 1px solid #f3f4f6;
  padding: 16px 0;

.cell {
  line-height: 1.6;
}
}

.el-button {
  transition: all 0.3s;
  margin: 0 4px;
  padding: 6px 12px;

&:not(.el-button--link):hover {
   transform: scale(1.05);
   box-shadow: 0 2px 4px rgba(0, 0, 0, 0.1);
 }

&--primary {
   background: linear-gradient(135deg, #60a5fa, #3b82f6);
   border: none;

&:hover {
   background: linear-gradient(135deg, #3b82f6, #2563eb);
 }
}

&--danger {
   background: linear-gradient(135deg, #f87171, #ef4444);
   border: none;

&:hover {
   background: linear-gradient(135deg, #ef4444, #dc2626);
 }
}

&--info {
   background: linear-gradient(135deg, #93c5fd, #60a5fa);
   border: none;

&:hover {
   background: linear-gradient(135deg, #60a5fa, #3b82f6);
 }
}
}

.el-tag {
  transition: all 0.3s;
  margin: 0 4px;
  border-radius: 4px;
  padding: 4px 8px;
  font-size: 0.9rem;
  border: none;
  box-shadow: 0 1px 2px rgba(0, 0, 0, 0.05);

&:hover {
   transform: translateY(-1px);
   box-shadow: 0 2px 4px rgba(0, 0, 0, 0.1);
 }
}

.el-table__empty-block {
  background-color: transparent;

.el-table__empty-text {
  color: #6b7280;
}
}

.role-tags {
  display: flex;
  flex-wrap: nowrap;
  align-items: center;
  gap: 4px;

.el-tag {
  transition: all 0.3s;
  margin: 0;
  white-space: nowrap;
  max-width: 120px;
  overflow: hidden;
  text-overflow: ellipsis;

&.more-tag {
   background: linear-gradient(135deg, #93c5fd, #60a5fa);
   color: white;
   cursor: pointer;
   padding: 0 8px;

&:hover {
   background: linear-gradient(135deg, #60a5fa, #3b82f6);
   transform: translateY(-1px);
 }
}
}
}

.role-popover {
.popover-tags {
  display: flex;
  flex-wrap: wrap;
  gap: 4px;
  max-height: 200px;
  overflow-y: auto;
  padding: 4px;

&::-webkit-scrollbar {
   width: 4px;
 }

&::-webkit-scrollbar-thumb {
   background-color: #cbd5e1;
   border-radius: 2px;
 }

.el-tag {
  margin: 0;
}
}
}
}
}

.pagination-container {
  margin-top: 1rem;
  display: flex;
  justify-content: flex-end;
  overflow-x: auto;

:deep(.el-pagination) {
  border-radius: 0.5rem;
  padding: 0.5rem;
  min-width: fit-content;

@media (max-width: 768px) {
  width: 100%;
  display: flex;
  justify-content: center;
  font-size: 0.875rem;

  .el-pager {
    flex-wrap: wrap;
  }

  .btn-prev,
  .btn-next {
    min-width: 24px;
  }

  li {
    min-width: 24px;
  }
}

.el-pagination__total,
.el-pagination__jump {
  color: #6b7280;
}

.el-pager li {
  background-color: transparent;
  color: #4b5563;
  border: 1px solid transparent;
  transition: all 0.3s;

&:hover {
   background-color: #f3f4f6;
 }

&.is-active {
   background-color: #3b82f6;
   color: white;
 }
}
}
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

  :deep(.el-form) {
  .el-form-item__label {
    color: #d1d5db;
  }

  .el-input__wrapper,
  .el-select .el-input__wrapper {
    background-color: #374151;
    border-color: #4b5563;

&:hover {
   border-color: #60a5fa;
 }

&.is-focus {
   border-color: #60a5fa;
 }
}
}
}

.table-container {
  background-color: rgba(31, 41, 55, 0.9);

:deep(.el-table) {
.el-table__row:hover {
  background-color: rgba(55, 65, 81, 0.7) !important;
}

th {
  background-color: rgba(55, 65, 81, 0.9) !important;
  color: #d1d5db;
  border-bottom-color: #4b5563;
}

td {
  border-bottom-color: #374151;
}

.el-button {
&--primary {
   background: linear-gradient(135deg, #3b82f6, #2563eb);
&:hover {
   background: linear-gradient(135deg, #2563eb, #1d4ed8);
 }
}

&--danger {
   background: linear-gradient(135deg, #ef4444, #dc2626);
&:hover {
   background: linear-gradient(135deg, #dc2626, #b91c1c);
 }
}

&--info {
   background: linear-gradient(135deg, #60a5fa, #3b82f6);
&:hover {
   background: linear-gradient(135deg, #3b82f6, #2563eb);
 }
}
}

.el-table__empty-block {
.el-table__empty-text {
  color: #9ca3af;
}
}

.role-tags {
.el-tag {
  background-color: #374151;
  color: #e5e7eb;
  border: 1px solid #4b5563;

&:hover {
   border-color: #60a5fa;
   background-color: #3b4252;
 }

&.more-tag {
   background: linear-gradient(135deg, #3b82f6, #2563eb);
   border: none;
   color: white;

&:hover {
   background: linear-gradient(135deg, #2563eb, #1d4ed8);
 }
}
}
}

.role-popover {
.popover-tags::-webkit-scrollbar-thumb {
  background-color: #4b5563;
}
}
}
}

.pagination-container {
:deep(.el-pagination) {
  background-color: rgba(31, 41, 55, 0.9);

.el-pagination__total,
.el-pagination__jump {
  color: #9ca3af;
}

.el-pager li {
  color: #d1d5db;

&:hover {
   background-color: #374151;
 }

&.is-active {
   background-color: #2563eb;
 }
}
}
}
}
}

// 添加弹出框的深色模式样式（注意：这个需要放在全局样式中，因为 popper 是挂载在 body 下的）
   :deep(.el-popper.role-popover) {
&.is-dark {
   background-color: #1f2937;
   border: 1px solid #374151;

.popover-tags {
.el-tag {
  background-color: #374151;
  color: #e5e7eb;
  border: 1px solid #4b5563;

&:hover {
   border-color: #60a5fa;
   background-color: #3b4252;
 }
}

&::-webkit-scrollbar {
   width: 4px;
 }

&::-webkit-scrollbar-thumb {
   background-color: #4b5563;
 }

&::-webkit-scrollbar-track {
   background-color: #1f2937;
 }
}
}

.el-popper__arrow::before {
  background-color: #1f2937;
  border-color: #374151;
}
}
</style>

