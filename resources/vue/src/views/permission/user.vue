<template>
  <div class="app-container">

    <div class="search-container">
      <el-form :inline="true">
        <el-form-item label="">
          <el-button
            type="info"
            class="filter-item"
            @click="openAuthCfgDialog"
          >{{ $t('第三方认证设置') }}
          </el-button>
        </el-form-item>
        <el-form-item label="">
          <el-button
            type="warning"
            class="filter-item"
            @click="handleAddRole"
          >{{ $t('新建用户') }}
          </el-button>
        </el-form-item>
        <el-form-item label="">
          <el-button
            type="success"
            class="filter-item"
            @click="getUserList"
          >{{ $t('查询') }}
          </el-button>
        </el-form-item>
      </el-form>
    </div>
    <el-card shadow="never" class="table-container">
      <el-table
      :data="rolesList"
    >
      <el-table-column align="center" label="id" width="50">
        <template #default="scope">
          {{ scope.row.id }}
        </template>
      </el-table-column>
      <el-table-column align="center" :label="$t('用户名')" width="220">
        <template #default="scope">
          {{ scope.row.username }}
        </template>
      </el-table-column>
      <el-table-column align="center" :label="$t('所属角色')" width="220">
        <template #default="scope">

          <template v-for="(item) in scope.row.role_ids">
            <el-tag>{{chanCfgMap[item]}}</el-tag>
          </template>

        </template>
      </el-table-column>
      <el-table-column align="center" :label="$t('真实姓名')">
        <template #default="scope">
          {{ scope.row.realname }}
        </template>
      </el-table-column>
        <el-table-column align="center" :label="$t('电子邮箱')">
          <template #default="scope">
            {{ scope.row.email }}
          </template>
        </el-table-column>
        <el-table-column align="center" :label="$t('企业微信成员UserID')">
          <template #default="scope">
            {{ scope.row.work_wechat_uid }}
          </template>
        </el-table-column>
      <el-table-column align="center" :label="$t('操作')" fixed="right" width="340">
        <template #default="scope">

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
      </el-table-column>
    </el-table>
    </el-card>
    <el-drawer
      title="第三方认证设置"
      v-model="authDialogVisible"
      size="90%"
    >
      <oauth v-if="authDialogVisible"></oauth>

    </el-drawer>
    <el-drawer
      title="修改密码"
      v-model="pwdDialogVisible"
      size="30%"
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
      size="85%"
      :title="dialogType==='edit'?$t('编辑用户信息'):$t('新建用户信息')"
    >
      <el-form :model="role" label-width="200px" label-position="left">

        <el-form-item :label="$t('用户名')">
          <el-input v-model="role.username" :placeholder="$t('用户名')" />
        </el-form-item>
        <el-form-item v-if="dialogType!=='edit'" :label="$t('密码')">
          <el-input show-password type="password" v-model="role.password" :placeholder="$t('不填则为系统自动生成')" />
        </el-form-item>
        <el-form-item :label="$t('真实姓名')">
          <el-input v-model="role.realname" :placeholder="$t('真实姓名')" />
        </el-form-item>

        <el-form-item :label="$t('电子邮箱')">
          <el-input v-model="role.email" :placeholder="$t('电子邮箱')" />
        </el-form-item>
        <el-form-item :label="$t('企业微信成员UserID')">
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
  data() {
    return {
      updatePwdForm:{
        password:'',
        confirmPassword:'',
        id:0,
      },
      authDialogVisible:false,
      pwdDialogVisible:false,
      role: Object.assign({}, defaultUser),
      routes: [],
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
  computed: {},
  async created() {
    await this.getRoleOpt()
    this.getUserList()
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

    async getUserList() {
      const res = await userList()

      this.rolesList = res.data
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

    openAuthCfgDialog(){
      this.authDialogVisible = true
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
          this.getUserList()
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
          this.getUserList()
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
          this.getUserList()
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
        this.getUserList()
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
        this.getUserList()
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
