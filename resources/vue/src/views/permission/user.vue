<template>
  <div class="app-container">

    <div class="search-container">
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
      <el-table-column align="center" :label="$t('所属角色')" >
        <template #default="scope">

          <template v-for="(item) in scope.row.role_ids">
            <el-tag>{{chanCfgMap[item]}}</el-tag>
          </template>

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
      :size="isMobile?'100%':'30%'"
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
      chanCfgMap: []
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
              type: 'error',
              message: res.msg
            })
            return
          }
          this.getUserList(1)
          ElMessage.success({
            type: 'success',
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
              type: 'error',
              message: res.msg
            })
            return
          }
          this.getUserList(1)
          ElMessage.success({
            type: 'success',
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
              type: 'error',
              message: res.msg
            })
            return
          }
          this.getUserList(1)
          ElMessage.success({
            type: 'success',
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
