<script src="../../store/modules/user.js"></script>
<template>
  <div class="login-container">
    <div class="Login_loginWrapper">
      <div>
        <!--        <div class="TextBox_productName">轻便的ElasticSearch可视化管理工具</div>-->
        <div class="TextBox_productDesc">ElasticView</div>
        <div class="TextBox_productSubDesc">您的ElasticSearch小管家😘😘😘</div>
      </div>
      <div class="LoginBox_boxWrap">
        <div>
          <div class="LoginBox_headerText">{{ getTitle() }}</div>
        </div>
        <el-form
          ref="loginForm"
          :model="loginForm"
          :rules="loginRules"
          class="login-form"
          autocomplete="on"
          label-position="left"
        >

          <el-form-item prop="username">

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

              <el-input
                :key="passwordType"
                ref="password"
                v-model="loginForm.password"
                :type="passwordType"
                placeholder="Password"
                name="password"
                tabindex="2"
                autocomplete="on"
                show-password
                @keyup.native="checkCapslock"
                @blur="capsTooltip = false"
                @keyup.enter.native="handleLogin"
              />

            </el-form-item>
          </el-tooltip>

          <el-button
            size="mini"

            icon="el-icon-switch-button"
            :loading="loading"
            type="primary"
            style="width:100%;margin-bottom:30px;"
            @click.native.prevent="handleLogin"
          >登录
          </el-button>

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
    </div>

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

      rememberme: false,
      loginForm: {
        username: '',
        password: ''
      },
      loginRules: {
        username: [{required: true, trigger: 'blur'}],
        password: [{required: true, trigger: 'blur', validator: validatePassword}]
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
      handler: function (route) {
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
    links(row) {
      window.location.href = row[1]
    },
    getTitle() {
      return process.env.VUE_APP_BASE_TITLE
    },
    checkCapslock({shiftKey, key} = {}) {
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
    },
    setLocalData(user, pwd) {
      let title = this.getTitle()
      let userInfo = {
        username: user,
        password: pwd
      }
      localStorage.setItem(title, JSON.stringify(userInfo))
    },
    getLocalData() {
      let title = this.getTitle()
      let info = localStorage.getItem(title)
      if (info != null) {
        let obj = JSON.parse(info)
        if (obj.username) {
          this.loginForm.username = obj.username
        }

        if (obj.password) {
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
              if (this.rememberme) {
                this.setLocalData(this.loginForm.username, this.loginForm.password)
                //this.setCookie(this.loginForm.username,this.loginForm.password,7)
              } else {
                this.deleteCookie()
              }
              //this.$websocket.dispatch('WEBSOCKET_INIT',"ws://127.0.0.1:11122/api/Chat?token="+getToken())
              this.$router.push({path: this.redirect || '/', query: this.otherQuery})
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
    setCookie(name, pass, days) {
      let expire = new Date()
      expire.setDate(expire.getDate() + days)
      document.cookie = `C-username=${name};expires=${expire}`
      document.cookie = `C-password=${pass};expires=${expire}`
    },
    getCookie() {
      if (document.cookie.length) {
        let arr = document.cookie.split('; ')
        for (let i = 0; i < arr.length; i++) {
          let arr2 = arr[i].split('=')
          if (arr2[0] === 'C-username') {
            this.loginForm.username = arr2[1]
          } else if (arr2[0] === 'C-password') {
            this.loginForm.password = arr2[1]
            this.rememberme = true
          }
        }
      }
    },
    deleteCookie() {
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

<style scoped src="@/styles/login.css"/>
