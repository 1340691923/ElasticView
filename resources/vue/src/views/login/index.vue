<template>
  <div class="login-container">

    <div class="flex-x-between absolute-lt w-full p-2">
      <div class="flex-center">

        <el-image
          :src="logo"
          style="width: 30px; height: 30px;"
        />

        <span class="text-2xl font-bold bg-gradient-to-r from-blue-500 to-teal-500 text-transparent bg-clip-text mx-1">
          {{ defaultSettings.title }}
        </span>

        <el-tag type="success">{{ appVersion() }}</el-tag>
      </div>
      <div class="flex-center">
        <!-- 顶部工具栏 -->
        <div class="top-bar">
          <el-switch
            v-model="isDark"
            inline-prompt
            active-icon="Moon"
            inactive-icon="Sunny"
            @change="toggleTheme"
          />
          <lang-select class="ml-2 cursor-pointer"/>
        </div>

      </div>
    </div>


    <!-- 登录表单 -->

    <div class="login-content">

      <div v-if="width > 992" class="login-image">
        <el-image :src="login_bg"/>
      </div>
      <div class="login-box">
        <h2 class="text-xl font-medium text-center flex-center relative">
          登 录
        </h2>
        <el-form
          ref="loginFormRef"
          :model="loginData"
          :rules="loginRules"
          class="login-form"
        >
          <el-form-item prop="username">
            <el-input
              ref="username"
              v-model="loginData.username"
              :placeholder="$t('用户名')"
              name="username"
              size="large"
              class="h-[48px] input-with-select"

            >
              <template #prepend>
                <i-ep-user class="mx-2"/>
              </template>
            </el-input>
          </el-form-item>

          <el-tooltip
            :visible="isCapslock"
            :content="$t('大写锁定已打开')"
            placement="right"
          >
            <el-form-item prop="password">
              <el-input
                v-model="loginData.password"
                :placeholder="$t('密码')"
                type="password"
                name="password"
                @keyup="checkCapslock"
                @keyup.enter="handleLoginSubmit"
                size="large"
                class="h-[48px] pr-2  input-with-select"
                show-password
              >
                <template #prepend>
                  <i-ep-lock class="mx-2"/>
                </template>
              </el-input>
            </el-form-item>
          </el-tooltip>

          <el-button
            :loading="loading"
            type="primary"
            size="large"
            class="w-full"
            @click.prevent="handleLoginSubmit"
          >{{ $t("登录") }}
          </el-button>

          <div class="flex-x-between w-full py-1">

            <el-checkbox
              v-model="rememberme"
              class="rememberme"
            >
              记住密码
            </el-checkbox>


          </div>
          <el-divider>
            其他登录方式
          </el-divider>

          <el-form-item style="border:none" >

              <span>
                <template
                v-for="(v,index) in data.oauthConfigs"
                >
                  <a v-if="v.enable" :title="v.name" :href="v.oauthUrl">
                    <el-image :alt="v.name" v-if="v.img == 'work_wechat'" :src="work_wechat" class="provider-img"> </el-image>
                  </a>
                </template>
              </span>
          </el-form-item>


        </el-form>
      </div>
    </div>


    <!--    <el-card class="login-card">
          <div class="text-center relative">
            <h2>{{ defaultSettings.title }}</h2>
            <el-tag class="ml-2 absolute-rt">{{ defaultSettings.version }}</el-tag>
          </div>


        </el-card>-->
  </div>
</template>

<script setup lang="ts">

import {GetOAuthList} from "@/api/user";

const logo = ref(new URL(`@/assets/logo.png`, import.meta.url).href);
//图片列表
const work_wechat = ref(new URL(`@/assets/images/work_wechat.png`, import.meta.url).href);


//
const login_bg = ref(new URL(`@/assets/images/login_bg.png`, import.meta.url).href);

const appStore = useAppStore();

const width = useWindowSize().width;

// 外部库和依赖
import {LocationQuery, useRoute} from "vue-router";

const rememberme = ref(false)

// 内部依赖
import {useAppStore, usePermissionStore, useSettingsStore, useUserStore} from "@/store";
import {LoginData} from "@/api/auth";
import router from "@/router";
import defaultSettings from "@/settings";
import {ThemeEnum} from "@/enums/ThemeEnum";

// 类型定义
import type {FormInstance} from "element-plus";

// 导入 login.scss 文件
import "@/styles/login.scss";
import {setToken} from "@/utils/auth";

// 使用导入的依赖和库
const userStore = useUserStore();
const settingsStore = useSettingsStore();
const route = useRoute();
// 窗口高度
const {height} = useWindowSize();
// 国际化 Internationalization
const {t} = useI18n();

// 是否暗黑模式
const isDark = ref(settingsStore.theme === ThemeEnum.DARK);
// 是否显示 ICP 备案信息
const icpVisible = ref(true);
// 按钮 loading 状态
const loading = ref(false);
// 是否大写锁定
const isCapslock = ref(false);
// 登录表单ref
const loginFormRef = ref<FormInstance>();

const loginData = ref<LoginData>({
  username: "",
  password: "",
} as LoginData);

const appVersion = ()=>{
  return window["appVersion"]
}

const loginRules = computed(() => {
  return {
    username: [
      {
        required: true,
        trigger: "blur",
        message: t("login.message.username.required"),
      },
    ],
    password: [
      {
        required: true,
        trigger: "blur",
        message: t("login.message.password.required"),
      },
      {
        min: 5,
        message: t("login.message.password.min"),
        trigger: "blur",
      },
    ],
  };
});

/** 登录表单提交 */
function handleLoginSubmit() {
  loginFormRef.value?.validate(async (valid: boolean) => {
    if (valid) {
      loading.value = true;
      userStore
        .login(loginData.value)
        .then((res) => {
          console.log("res", res)

          if (res.code != 0) {
            ElMessage.error(res.msg)
            return
          }
          if (rememberme.value == true) {
            setLocalData(loginData.value.username, loginData.value.password)
          } else {
            localStorage.removeItem("ev-login_info")
          }

          setToken(res.data.token)
          ElMessage.success(res.msg)
          const {path, queryParams} = parseRedirect();

          router.push({path: path, query: queryParams});
        })
        .catch(() => {

        })
        .finally(() => {
          loading.value = false;
        });
    }
  });
}

/** 解析 redirect 字符串 为 path 和  queryParams */
function parseRedirect(): {
  path: string;
  queryParams: Record<string, string>;
} {
  const query: LocationQuery = route.query;
  const redirect = (query.redirect as string) ?? "/";

  const url = new URL(redirect, window.location.origin);
  const path = url.pathname;
  const queryParams: Record<string, string> = {};

  url.searchParams.forEach((value, key) => {
    if (key != "code" && key != "state") {
      queryParams[key] = value;
    }
  });

  return {path, queryParams};
}



/** 主题切换 */
const toggleTheme = () => {
  const newTheme = settingsStore.theme === ThemeEnum.DARK ? ThemeEnum.LIGHT : ThemeEnum.DARK;
  settingsStore.changeTheme(newTheme);
};

/** 根据屏幕宽度切换设备模式 */
watchEffect(() => {
  if (height.value < 600) {
    icpVisible.value = false;
  } else {
    icpVisible.value = true;
  }
});

/** 检查输入大小写 */
function checkCapslock(event: KeyboardEvent) {
  // 防止浏览器密码自动填充时报错
  if (event instanceof KeyboardEvent) {
    isCapslock.value = event.getModifierState("CapsLock");
  }
}

const getLocalData = () => {

  let info = localStorage.getItem("ev-login_info")

  if (info != null) {
    let obj = JSON.parse(info)
    if (obj.username) {
      loginData.value.username = obj.username
    }

    if (obj.password) {
      loginData.value.password = obj.password
    }

    rememberme.value = true
  }
}

const setLocalData = (user, pwd) => {

  let userInfo = {
    username: user,
    password: pwd
  }
  console.log(userInfo, user, pwd)
  localStorage.setItem("ev-login_info", JSON.stringify(userInfo))

}

const GetCurrentURL = ()=>{
  const currentUrl = window.location.href;
  const url = new URL(currentUrl);

  let loginCallback = `${url.origin}`

  return loginCallback
}

const data = reactive({
  oauthConfigs:[]
})

const InitOAuthConfis = async ()=>{

  const res = await GetOAuthList({
    call_back:GetCurrentURL()
  })

  if(res.code != 0){
    return
  }

  if(res.data == null) res.data = []

  data.oauthConfigs = res.data
}

onMounted(async () => {

  //todo..
  InitOAuthConfis()

  if (route.query.hasOwnProperty('code') && route.query.hasOwnProperty('state')) {
    let code = route.query.code
    let state = route.query.state
    userStore
      .login({
        oauth_code:code,
        state,
      })
      .then((res) => {
        console.log("res", res)

        if (res.code != 0) {
          ElMessage.error(res.msg)
          return
        }

        setToken(res.data.token)
        ElMessage.success(res.msg)
        const {path, queryParams} = parseRedirect();

        router.push({path: path, query: queryParams});
      })
      .catch(() => {

      })
      .finally(() => {
        loading.value = false;
      });
  }

  getLocalData()
})

</script>

<style lang="scss" scoped></style>

<style scoped>
.login-container .login-content {
  display: flex;
  width: 100%;
  min-width: 400px;
  max-width: 850px;
  overflow: hidden;
  border-radius: 5px;
  box-shadow: var(--el-box-shadow-light);
}

.login-container .login-content .login-image {
  display: flex;
  flex: 3;
  align-items: center;
  justify-content: center;
  background: linear-gradient(60deg, #f6fbfd, #f7fbfe);
}

.login-container .login-content .login-box {
  display: flex;
  flex: 2;
  flex-direction: column;
  justify-content: center;
  min-width: 400px;
  padding: 30px;
}


@media (width <= 768px

) {
  .login-container .login-content .login-box {
    width: 100%
  }
}


@media (width <= 768px

) {
  .login-container .login-content {
    flex-direction: column;
    max-width: 100%;
    height: 100vh;
    border-radius: 0;
    box-shadow: none
  }
}
.provider-img {
  width: 30px;
  margin: 5px;
}

</style>

<style>

</style>
