<template>
  <div class="app-container">
    <div class="search-container">
      <el-form :inline="true">
        <el-form-item label="" >
          <el-button
            type="success"
            class="filter-item"
            @click="getList(1)"
          >{{$t('搜索')}}
          </el-button>
        </el-form-item>


        <el-form-item label="" >
          <el-button

            type="primary"
            class="filter-item"
            @click="handleAddRole"
          >{{ $t('新建数据源') }}
          </el-button>
        </el-form-item>

        <el-form-item :label=" $t('鉴权配置')" >

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
      <el-table-column width="50"  :label="$t('鉴权')" type="expand">
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
            <el-table-column align="center" :label="$t('分配权限组')" width="100">
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
      <el-table-column align="center"   :label="$t('数据源地址')" >
        <template #default="scope">
          {{ scope.row.ip }}
        </template>
      </el-table-column>
      <el-table-column align="center" :label="$t('数据源类型')" width="140">
        <template #default="scope">
          {{ scope.row.version }}
        </template>
      </el-table-column>
      <el-table-column align="center"  width="200" :label="$t('创建时间')" >
        <template #default="scope">
          {{ scope.row.created }}
        </template>
      </el-table-column>
      <el-table-column align="center" width="200"  :label="$t('修改时间')" >
        <template #default="scope">
          {{ scope.row.updated }}
        </template>
      </el-table-column>
      <el-table-column align="center"   fixed="right" :label="$t('操作')" width="130">
        <template #default="scope">

          <el-button
            type="primary"

            @click="handleEdit(scope)"
            :icon="Edit"
          >
          </el-button>
          <el-button
            type="danger"
            :icon="Delete"
            @click="handleDelete(scope)"
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
      :size="isMobile?'100%':'80%'"

      v-model="dialogVisible"
      :title="dialogType==='edit'?$t('编辑数据源'):$t('新建数据源')"
    >
      <el-form :model="link" label-width="100px" label-position="left">
        <el-form-item :label="$t('备注')">
          <el-input clearable v-model="link.remark" :placeholder="$t('备注')" />
        </el-form-item>
        <el-form-item :label="$t('数据源地址')">
          <el-input v-model="link.ip" :placeholder="$t('例如:http://127.0.0.1:9200')" />
        </el-form-item>
        <el-form-item :label="$t('数据源类型')">
          <el-select @change="changeVersion" v-model="link.version" :placeholder="$t('请选择数据源类型')" filterable>
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
import {Plus,Search,Edit,Delete} from "@element-plus/icons-vue";
import {useAppStore} from "@/store";
import {DeviceEnum} from "@/enums/DeviceEnum";

const defaultLink = {
  id: 0,
  ip: 'http://127.0.0.1:9200',
  remark: '',
  version: 'elasticsearch8.x',
  cfgIds:[],
}

export default {
  setup(){
    const appStore = useAppStore()
    const { device } = storeToRefs(appStore)
    return {Edit,Delete,device}
  },
  data() {
    return {
      input: {
        page: 1,
        limit: 10,
      },
      count: 0,
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
    this.getList(1)
  },
  computed: {
    isMobile(){
      return this.device === DeviceEnum.MOBILE
    },
  },
  methods: {
    changeVersion(){
      /***
       * <el-option label="elasticsearch6.x" value="elasticsearch6.x" />
       *             <el-option label="elasticsearch7.x" value="elasticsearch7.x" />
       *             <el-option label="elasticsearch8.x" value="elasticsearch8.x" />
       *             <el-option label="mysql" value="mysql" />
       *             <el-option label="redis" value="redis" />
       *             <el-option label="clickhouse" value="clickhouse" />
       *             <el-option label="postgres" value="postgres" />
       *             <el-option label="mongo" value="mongo" />
       */
      switch (this.link.version){
        case "elasticsearch6.x":
          this.link.ip = "http://127.0.0.1:9200"
          break
        case "elasticsearch7.x":
          this.link.ip = "http://127.0.0.1:9200"
          break
        case "elasticsearch8.x":
          this.link.ip = "https://127.0.0.1:9200"
          break
        case "mysql":
          this.link.ip = "127.0.0.1:3306"
          break
        case "clickhouse":
          this.link.ip = "127.0.0.1:9000"
          break
        case "postgres":
          this.link.ip = "127.0.0.1:5432"
          break
        case "mongo":
          this.link.ip = "127.0.0.1:27017"
          break
        case "redis":
          this.link.ip = "127.0.0.1:6379"
          break
      }
    },
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
      ElMessageBox.confirm('确定删除该配置吗?', '警告', {
        confirmButtonText: '确认',
        cancelButtonText: '取消',
        type: 'warning'
      })
        .then(async() => {
          const res = await DeleteEsCfgRelation({id:scope2.row.cfg_relation_id})
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
    filterMethod(query, item) {
      return item.label.indexOf(query) > -1
    },
    async getEsCfgOpt(){
      const res = await GetEsCfgOpt()
      if (res.code != 0) {
        ElMessage.error({
          type: 'error',offset:60,
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
            type: 'success',offset:60,
            message: `连接成功`
          })
        } else {
          ElMessage.error({
            type: 'error',offset:60,
            message: res.msg
          })
        }
      }).catch(err => {
        scope.row.connectLoading = false
        console.log("err",err)
      })

    },
    handleSizeChange(v) {
      this.input.limit = v
      this.getList(1)
    },
    async getList(page) {
      !page ? this.input.page = 1 :  this.input.page = page
      this.getListLoading = true
      const res = await ListAction({
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

      for (const k in res.data.list) {
        if(res.data.list[k].es_link_configs == null){
          res.data.list[k].es_link_configs = []
        }
        for (const k2 in res.data.list[k].es_link_configs) {
          res.data.list[k].es_link_configs[k2]['connectLoading'] = false
        }
      }
      this.count = res.data.count
      this.list = res.data.list
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
          type: 'error',offset:60,
          message: res.msg
        })
        return
      }
      ElMessage.success({
        type: 'success',offset:60,
        message: res.msg
      })
      this.getList()
    },
    handleDelete({ $index, row }) {
      ElMessageBox.confirm('确定删除该数据源吗?', '警告', {
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

