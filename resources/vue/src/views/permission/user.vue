<template>
  <div class="app-container">

    <div class="search-container" :class="{ 'is-collapsed': isCollapsed && isMobile }">
      <div class="search-header" v-if="isMobile" @click="toggleCollapse">
        <span>{{ $t('搜索条件') }}</span>
        <el-icon :class="{ 'is-collapsed': isCollapsed }">
          <ArrowDown />
        </el-icon>
      </div>
      <el-form :inline="true">

        <el-form-item label="用户名:">
          <el-input clearable v-model="input.user_name" style="width: 200px" ></el-input>
        </el-form-item>
        <el-form-item label="姓名:">
          <el-input clearable v-model="input.real_name" style="width: 200px" ></el-input>
        </el-form-item>
        <el-form-item label="是否封禁:">
          <el-switch
            @change="getUserList(1)"
            v-model="input.is_ban"
            active-text="是"
            inactive-text="否">
          </el-switch>
        </el-form-item>
        <el-form-item label="">
          <el-button
            type="success"
            class="filter-item"
            @click="getUserList(1)"
          >{{ $t('查询') }}
          </el-button>
          <el-button
            type="warning"
            class="filter-item"
            @click="handleAddRole"
          >{{ $t('新增') }}
          </el-button>

        </el-form-item>



      </el-form>
    </div>
    <el-card shadow="never" class="table-container">
      <el-table
      :data="rolesList"
      v-loading="loading"

    >
      <el-table-column align="center" sortable label="id" width="50">
        <template #default="scope">
          {{ scope.row.id }}
        </template>
      </el-table-column>
      <el-table-column align="center"  :label="$t('用户名')" width="100">
        <template #default="scope">
          {{ scope.row.username }}
        </template>
      </el-table-column>
      <el-table-column align="center" width="160" :label="$t('角色')">
        <template #default="scope">
          <div class="role-tags">
            <template v-if="scope.row.role_ids.length <= 2">
              <el-tag v-for="item in scope.row.role_ids" :key="item">
                {{chanCfgMap[item]}}
              </el-tag>
            </template>
            <template v-else>
              <el-tag v-for="item in scope.row.role_ids.slice(0, 1)" :key="item">
                {{chanCfgMap[item]}}
              </el-tag>
              <el-popover
                placement="bottom"
                trigger="hover"
                :width="200"
                popper-class="role-popover"
              >
                <template #reference>
                  <el-tag class="more-tag">
                    +{{ scope.row.role_ids.length - 1 }}
                  </el-tag>
                </template>
                <div class="popover-tags">
                  <el-tag
                    v-for="item in scope.row.role_ids.slice(1)"
                    :key="item"
                    style="margin-left:1rem"
                    size="small"
                  >
                    {{chanCfgMap[item]}}
                  </el-tag>
                </div>
              </el-popover>
          </template>
          </div>
        </template>
      </el-table-column>
      <el-table-column align="center" width="100"  :label="$t('姓名')">
        <template #default="scope">
          {{ scope.row.realname }}
        </template>
      </el-table-column>
        <el-table-column align="center"  width="200" :label="$t('电子邮箱')">
          <template #default="scope">
            {{ scope.row.email }}
          </template>
        </el-table-column>

        <el-table-column align="center"   width="100" :label="$t('企微成员UserID')">
          <template #default="scope">
            {{ scope.row.work_wechat_uid }}
          </template>
        </el-table-column>

        <el-table-column align="center" sortable  width="180" :label="$t('创建时间')">
          <template #default="scope">
            {{ scope.row.create_time }}
          </template>
        </el-table-column>

        <el-table-column align="center"  sortable width="180" :label="$t('修改时间')">
          <template #default="scope">
            {{ scope.row.update_time }}
          </template>
        </el-table-column>

        <el-table-column align="center" sortable  width="180" :label="$t('最后登录时间')">
          <template #default="scope">
            {{ scope.row.last_login_time }}
          </template>
        </el-table-column>

      <el-table-column align="center" :label="$t('操作')" fixed="right" :width="isMobile?80:400">
        <template #default="scope">
          <template v-if="!isMobile">
            <el-button
              type="primary"
              @click="handleEdit(scope)"
            >{{ $t('编辑') }}
            </el-button>
            <el-button
              type="info"
              @click="openPwdDialog(scope.row)"
            >{{ $t('重置密码') }}
            </el-button>

            <el-button
              v-if="scope.row.is_ban == 1"
              type="primary"
              @click="UnSealUser(scope)"
            >{{ $t('解封') }}
            </el-button>
            <el-button
              v-if="scope.row.is_ban == 0"
              type="danger"
              @click="SealUser(scope)"
            >{{ $t('封禁') }}
            </el-button>
            <el-button
              type="danger"
              @click="handleDelete(scope)"
            >{{ $t('删除') }}
            </el-button>
          </template>
          <template v-else >
            <el-dropdown trigger="click" >
              <el-button>管理</el-button>
              <template #dropdown>
                <el-dropdown-menu>
                  <el-dropdown-item  command="1">
                    <el-button
                      type="primary"
                      @click.stop="handleEdit(scope)"
                    >{{ $t('编辑') }}
                    </el-button>
                  </el-dropdown-item>
                  <el-dropdown-item  >
                    <el-button
                      type="info"
                      @click="openPwdDialog(scope.row)"
                    >{{ $t('重置密码') }}
                    </el-button>
                  </el-dropdown-item>
                  <el-dropdown-item  >
                    <el-button
                      v-if="scope.row.is_ban == 1"
                      type="primary"
                      @click="UnSealUser(scope)"
                    >{{ $t('解封') }}
                    </el-button>
                    <el-button
                      v-if="scope.row.is_ban == 0"
                      type="danger"
                      @click="SealUser(scope)"
                    >{{ $t('封禁') }}
                    </el-button>
                  </el-dropdown-item>

                  <el-dropdown-item  >
                    <el-button
                      type="danger"
                      @click="handleDelete(scope)"
                    >{{ $t('删除') }}
                    </el-button>
                  </el-dropdown-item>
                </el-dropdown-menu>
              </template>
            </el-dropdown>
          </template>


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
        @current-change="getUserList"
        @size-change="handleSizeChange"
      />
    </div>


    <el-drawer
      title="修改密码"
      v-model="pwdDialogVisible"
      :size="isMobile?'100%':'50%'"
    >
      <el-form label-width="100px" >
        <el-form-item label="新密码" required>
          <el-input  v-model="updatePwdForm.password" show-password type="password" placeholder="请输入新密码" />
        </el-form-item>
        <el-form-item label="确认密码" required>
          <el-input v-model="updatePwdForm.confirmPassword" show-password type="password" placeholder="请再次输入新密码" />
        </el-form-item>
      </el-form>

      <template #footer>
        <el-button @click="pwdDialogVisible = false">取消</el-button>
        <el-button type="primary" @click="submitPwdForm">修改</el-button>
      </template>
    </el-drawer>
    <el-drawer

      v-model="dialogVisible"
      :size="isMobile?'100%':'50%'"
      :title="dialogType==='edit'?$t('编辑用户信息'):$t('新建用户信息')"
    >
      <el-form :model="role" label-width="100px" label-position="left">

        <el-form-item :label="$t('用户名')">
          <el-input v-model="role.username" :placeholder="$t('用户名')" />
        </el-form-item>
        <el-form-item v-if="dialogType!=='edit'" :label="$t('密码')">
          <el-input show-password type="password" v-model="role.password" :placeholder="$t('不填则为系统自动生成')" />
        </el-form-item>
        <el-form-item :label="$t('姓名')">
          <el-input v-model="role.realname" :placeholder="$t('姓名')" />
        </el-form-item>

        <el-form-item :label="$t('电子邮箱')">
          <el-input v-model="role.email" :placeholder="$t('电子邮箱')" />
        </el-form-item>
        <el-form-item :label="$t('企微UserID')">
          <el-input v-model="role.work_wechat_uid" :placeholder="$t('企业微信成员UserID。对应管理端的账号，企业内必须唯一')" />
        </el-form-item>

        <el-form-item :label="$t('请选择角色')">
          <el-select v-model="role.role_ids" multiple :placeholder="$t('请选择角色')" clearable filterable>
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
          type="danger"
          @click="dialogVisible=false"
        >{{ $t('取消') }}
        </el-button>
        <el-button
          type="primary"
          @click="confirmRole"
        >{{ $t('确认') }}
        </el-button>
      </div>
    </el-drawer>
  </div>
</template>

<script>
import { deepClone } from '@/utils'
import {
  DelUser, InsertUser, roleOption, UpdateUser, userList, ModifyPwdByUserId,
  UnSealUserAction,SealUserAction
} from '@/api/user'
import Oauth from '../permission/oauth.vue'
import {ElMessage} from "element-plus";

import {useAppStore} from "@/store";
import {DeviceEnum} from "@/enums/DeviceEnum";
import { ArrowDown } from '@element-plus/icons-vue'



const defaultUser = {
  id: '',
  password: '',
  realname: '',
  email:'',
  work_wechat_uid:'',
  role_ids: [],
  username: '',
  role_name: ''
}

export default {
  setup(){
    const appStore = useAppStore()
    const { device } = storeToRefs(appStore)
    return {device}
  },
  data() {
    return {
      input: {
        page: 1,
        limit: 10,
        user_name:'',
        real_name:'',
        is_ban:false
      },
      count: 0,
      updatePwdForm:{
        password:'',
        confirmPassword:'',
        id:0,
      },
      pwdDialogVisible:false,
      role: Object.assign({}, defaultUser),
      routes: [],
      loading:false,
      rolesList: [],
      dialogVisible: false,
      dialogType: 'new',
      checkStrictly: false,
      defaultProps: {
        children: 'children',
        label: 'title'
      },
      chanCfgList: [],
      chanCfgMap: [],
      isCollapsed: true,
    }
  },
  computed: {
    isMobile(){
      return this.device === DeviceEnum.MOBILE
    },
  },
  async created() {
    await this.getRoleOpt()
    this.getUserList(1)
  },
  components:{
    Oauth
  },
  methods: {
    async getRoleOpt() {
      const res = await roleOption()
      for (var v of res.data) {
        this.chanCfgList.push(v)
        this.chanCfgMap[v['id']] = v['name']
      }
    },
    handleSizeChange(v) {
      this.input.limit = v
      this.getUserList(1)
    },
    async getUserList(page) {
      !page ? this.input.page = 1 :  this.input.page = page
      this.loading = true
      const res = await userList({
        page:page,
        page_size: this.input.limit,
        user_name:this.input.user_name,
        real_name:this.input.real_name,
        is_ban:this.input.is_ban,
      })
      this.loading = false
      this.count = res.data.count

      this.rolesList = res.data.list
    },
    async submitPwdForm(){

      if(this.updatePwdForm.password == ''){
        ElMessage.error({
          message: "密码不能为空",
          type: 'error'
        })
        return
      }
      if(this.updatePwdForm.confirmPassword != this.updatePwdForm.password){
        ElMessage.error({
          message: "两次密码输入不一致",
          type: 'error'
        })
        return
      }
      const res = await ModifyPwdByUserId(this.updatePwdForm)
      if(res.code != 0 ){
        ElMessage.error({
          message: res.msg,
          type: 'error'
        })
        return
      }
      ElMessage.success({
        message: res.msg,
        type: 'success'
      })
      this.updatePwdForm.password = ''
      this.updatePwdForm.confirmPassword = ''
      this.pwdDialogVisible = false
    },

    openPwdDialog(item){
      this.updatePwdForm.id = item.id
      this.pwdDialogVisible = true
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
      ElMessageBox.confirm("确认删除该用户吗?", "警告", {
        confirmButtonText: "确定",
        cancelButtonText: "取消",
        type: "warning",
      })
        .then(async() => {
          const res = await DelUser({ id: row.id })
          if(res.code != 0){
            ElMessage({
              type: 'error',offset:60,
              message: res.msg
            })
            return
          }
          this.getUserList(1)
          ElMessage.success({
            type: 'success',offset:60,
            message: res.msg
          })
        })
        .catch(err => {
          console.error(err)
        })
    },

    SealUser({ $index, row }) {
      ElMessageBox.confirm("确认封禁该用户吗?", "警告", {
        confirmButtonText: "确定",
        cancelButtonText: "取消",
        type: "warning",
      })
        .then(async() => {
          const res = await SealUserAction({ id: row.id })
          if(res.code != 0){
            ElMessage({
              type: 'error',offset:60,
              message: res.msg
            })
            return
          }
          this.getUserList(1)
          ElMessage.success({
            type: 'success',offset:60,
            message: res.msg
          })
        })
        .catch(err => {
          console.error(err)
        })
    },

    UnSealUser({ $index, row }) {
      ElMessageBox.confirm("确认解封该用户吗?", "警告", {
        confirmButtonText: "确定",
        cancelButtonText: "取消",
        type: "warning",
      })
        .then(async() => {
          const res = await UnSealUserAction({ id: row.id })
          if(res.code != 0){
            ElMessage({
              type: 'error',offset:60,
              message: res.msg
            })
            return
          }
          this.getUserList(1)
          ElMessage.success({
            type: 'success',offset:60,
            message: res.msg
          })
        })
        .catch(err => {
          console.error(err)
        })
    },

    async confirmRole() {
      const isEdit = this.dialogType === 'edit'


      if (this.role.role_ids.length == 0) {
        ElMessage.error({
          message: '角色不能为空',
          type: 'error'
        })
        return false
      }

      if (isEdit) {
        const res = await UpdateUser(this.role)
        if(res.code != 0 ){
          ElMessage.error({
            message: res.msg,
            type: 'error'
          })
          return
        }
        ElMessage.success({
          message: res.msg,
          type: 'success'
        })
        this.getUserList(1)
      } else {

        if (this.role.password.length < 5 && this.role.password.trim() != '') {
          ElMessage.error({
            message: '密码长度必须大于5',
            type: 'error'
          })
          return false
        }

        const res = await InsertUser(this.role)
        if(res.code != 0 ){
          ElMessage.error({
            message: res.msg,
            type: 'error'
          })
          return
        }
        ElMessage.success({
          message: res.msg,
          type: 'success'
        })
        this.getUserList(1)
      }

      this.dialogVisible = false

    },
    toggleCollapse() {
      this.isCollapsed = !this.isCollapsed
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
