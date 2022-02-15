<script src="../../store/modules/user.js"></script>
<template>
  <div class="login-container">

    <el-form ref="loginForm" :model="loginForm" :rules="loginRules" class="login-form" autocomplete="on" label-position="left">

      <div class="title-container">
        <h3 class="title">{{ getTitle() }}</h3>
      </div>

      <el-form-item prop="username">
        <span class="svg-container">
          <i class="el-icon-s-custom" />
        </span>
        <el-input
          ref="username"
          v-model="loginForm.username"
          placeholder="Username"
          name="username"
          type="text"
          tabindex="1"
          autocomplete="on"
        />
      </el-form-item>

      <el-tooltip v-model="capsTooltip" content="Caps lock is On" placement="right" manual>
        <el-form-item prop="password">
          <span class="svg-container">
            <i class="el-icon-lock" />
          </span>
          <el-input
            :key="passwordType"
            ref="password"
            v-model="loginForm.password"
            :type="passwordType"
            placeholder="Password"
            name="password"
            tabindex="2"
            autocomplete="on"
            @keyup.native="checkCapslock"
            @blur="capsTooltip = false"
            @keyup.enter.native="handleLogin"
          />
          <span class="show-pwd" @click="showPwd">
            <i :class="passwordType === 'password' ? 'el-icon-view' : 'el-icon-view'" />
          </span>
        </el-form-item>
      </el-tooltip>
      <!--<el-cascader
        placeholder="试试搜索：火柴人"
        :options="options"
        filterable
        @change="links"
      />-->
      <el-button icon="el-icon-switch-button" :loading="loading" type="primary" style="width:100%;margin-bottom:30px;" @click.native.prevent="handleLogin">登录</el-button>

      <div style="position:relative">
        <div class="tips">
          <span>
            <el-checkbox
              v-model="rememberme"
              class="rememberme"
              @change="changeRemember"
            >
              记住密码
            </el-checkbox>
          </span>
        </div>
        <div class="tips">
          <span style="margin-right:18px;" />
          <span />
        </div>
      </div>
    </el-form>
  </div>
</template>

<script>

  export default {
    name: 'Login',
    data() {
      const validatePassword = (rule, value, callback) => {
        if (value.length < 5) {
          callback(new Error('密码位数必须大于5'))
        } else {
          callback()
        }
      }

      return {
        options: [],

        rememberme:false,
          loginForm: {
          username: '',
          password: ''
        },
        loginRules: {
          username: [{ required: true, trigger: 'blur' }],
          password: [{ required: true, trigger: 'blur', validator: validatePassword }]
        },
        passwordType: 'password',
        capsTooltip: false,
        loading: false,
        showDialog: false,
        redirect: undefined,
        otherQuery: {},

      }
    },
    watch: {
      $route: {
        handler: function(route) {
          const query = route.query
          if (query) {
            this.redirect = query.redirect
            this.otherQuery = this.getOtherQuery(query)
          }
        },
        immediate: true
      }
    },
    mounted() {

      this.getLocalData()
      if (this.loginForm.username === '') {
        this.$refs.username.focus()
      } else if (this.loginForm.password === '') {
        this.$refs.password.focus()
      }
    },
    methods: {
      links(row){
        window.location.href = row[1]
      },
      getTitle(){
        return process.env.VUE_APP_BASE_TITLE
      },
      checkCapslock({ shiftKey, key } = {}) {
        if (key && key.length === 1) {
          if (shiftKey && (key >= 'a' && key <= 'z') || !shiftKey && (key >= 'A' && key <= 'Z')) {
            this.capsTooltip = true
          } else {
            this.capsTooltip = false
          }
        }
        if (key === 'CapsLock' && this.capsTooltip === true) {
          this.capsTooltip = false
        }
      },
      showPwd() {
        if (this.passwordType === 'password') {
          this.passwordType = ''
        } else {
          this.passwordType = 'password'
        }
        this.$nextTick(() => {
          this.$refs.password.focus()
        })
      },
      changeRemember() {
        console.log(this.rememberme);
      },
      setLocalData(user,pwd){
        let title = this.getTitle()
        let userInfo = {
          username:user,
          password:pwd
        }
        localStorage.setItem(title,JSON.stringify(userInfo))
      },
      getLocalData(){
        let title = this.getTitle()
        let info = localStorage.getItem(title)
        if(info != null){
          let obj = JSON.parse(info)
          if (obj.username){
            this.loginForm.username= obj.username
          }

          if (obj.password){
            this.loginForm.password = obj.password
          }

          this.rememberme = true
        }
      },
      handleLogin() {
        this.$refs.loginForm.validate(valid => {
          if (valid) {
            this.loading = true
            this.$store.dispatch('user/login', this.loginForm)
              .then(() => {
                if(this.rememberme){
                  this.setLocalData(this.loginForm.username,this.loginForm.password)
                  //this.setCookie(this.loginForm.username,this.loginForm.password,7)
                }else{
                  this.deleteCookie()
                }
                //this.$websocket.dispatch('WEBSOCKET_INIT',"ws://127.0.0.1:11122/api/Chat?token="+getToken())
                this.$router.push({ path: this.redirect || '/', query: this.otherQuery })
                this.loading = false
              })
              .catch(() => {
                this.loading = false
              })
          } else {
            console.log('error submit!!')
            return false
          }
        })
      },
      setCookie(name, pass, days){
        let expire = new Date()
        expire.setDate(expire.getDate() + days)
        document.cookie = `C-username=${name};expires=${expire}`
        document.cookie = `C-password=${pass};expires=${expire}`
      },
      getCookie(){
        if(document.cookie.length){
          let arr = document.cookie.split('; ')
          for(let i=0; i<arr.length; i++){
            let arr2 = arr[i].split('=')
            if(arr2[0] === 'C-username'){
              this.loginForm.username = arr2[1]
            }else if(arr2[0] === 'C-password'){
              this.loginForm.password = arr2[1]
              this.rememberme = true
            }
          }
        }
      },
      deleteCookie(){
        localStorage.removeItem(this.getTitle())
        //this.setCookie('', '', -1);
      },
      getOtherQuery(query) {
        return Object.keys(query).reduce((acc, cur) => {
          if (cur !== 'redirect') {
            acc[cur] = query[cur]
          }
          return acc
        }, {})
      }
    }
  }
</script>

<style lang="scss">

  $bg:#283443;
  $light_gray:#fff;
  $cursor: #fff;

  @supports (-webkit-mask: none) and (not (cater-color: $cursor)) {
    .login-container .el-input input {
      color: $cursor;
    }
  }

  /* reset element-ui css */
  .login-container {
    .el-input {
      display: inline-block;
      height: 47px;
      width: 85%;

      input {
        background: transparent;
        border: 0px;
        -webkit-appearance: none;
        border-radius: 0px;
        padding: 12px 5px 12px 15px;
        color: $light_gray;
        height: 47px;
        caret-color: $cursor;

        &:-webkit-autofill {
          box-shadow: 0 0 0px 1000px $bg inset !important;
          -webkit-text-fill-color: $cursor !important;
        }
      }
    }

    .el-form-item {
      border: 1px solid rgba(255, 255, 255, 0.1);
      background: rgba(0, 0, 0, 0.1);
      border-radius: 5px;
      color: #454545;
    }
  }
</style>

<style lang="scss" scoped>
  $bg:#2d3a4b;
  $dark_gray:#889aa4;
  $light_gray:#eee;

  .login-container {
    min-height: 100%;
    width: 100%;
    background-color: $bg;
    overflow: hidden;

    .login-form {
      position: relative;
      width: 520px;
      max-width: 100%;
      padding: 160px 35px 0;
      margin: 0 auto;
      overflow: hidden;
    }

    .tips {
      font-size: 14px;
      color: #fff;
      margin-bottom: 10px;

      span {
        &:first-of-type {
          margin-right: 16px;
        }
      }
    }

    .svg-container {
      padding: 6px 5px 6px 15px;
      color: $dark_gray;
      vertical-align: middle;
      width: 30px;
      display: inline-block;
    }

    .title-container {
      position: relative;

      .title {
        font-size: 26px;
        color: $light_gray;
        margin: 0px auto 40px auto;
        text-align: center;
        font-weight: bold;
      }
    }

    .show-pwd {
      position: absolute;
      right: 10px;
      top: 7px;
      font-size: 16px;
      color: $dark_gray;
      cursor: pointer;
      user-select: none;
    }

    .thirdparty-button {
      position: absolute;
      right: 0;
      bottom: 6px;
    }

    @media only screen and (max-width: 470px) {
      .thirdparty-button {
        display: none;
      }
    }
  }
</style>
