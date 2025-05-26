<template>
  <div class="login-container">
    <!-- 顶部导航栏美化 -->
    <div class="flex-x-between absolute-lt w-full p-4 z-10">
      <div class="flex-center">
        <el-image
          :src="logo"
          class="w-8 h-8 hover:rotate-180 transition-all duration-500"
        />
        <span class="text-2xl font-bold bg-gradient-to-r from-blue-500 to-teal-500 text-transparent bg-clip-text mx-2">
          {{ defaultSettings.title }}
        </span>
        <el-tag type="success" class="animate-pulse">{{ appVersion() }}</el-tag>
      </div>
      <div class="flex-center">
        <div class="top-bar">
          <el-switch
            v-model="isDark"
            inline-prompt
            active-icon="Moon"
            inactive-icon="Sunny"
            class="hover:scale-105 transition-all"
            @change="toggleTheme"
          />
          <lang-select class="ml-3 cursor-pointer hover:scale-105 transition-all"/>
        </div>
      </div>
    </div>

    <!-- 登录区域美化 -->
    <div class="login-content backdrop-blur-sm bg-white/30 dark:bg-gray-800/30">
      <div v-if="width > 992" class="login-image">
        <el-image v-if="isDark" :src="bg_dark" class="w-full h-full object-cover hover:scale-105 transition-all duration-500"/>
        <el-image v-else :src="bg_light" class="w-full h-full object-cover hover:scale-105 transition-all duration-500"/>

      </div>
      <div class="login-box">
        <h2 class="text-2xl font-bold text-center mb-8 bg-gradient-to-r from-blue-500 to-teal-500 text-transparent bg-clip-text">
          欢迎登录
        </h2>
        <el-form
          ref="loginFormRef"
          :model="loginData"
          :rules="loginRules"
          class="login-form"
        >
          <!-- 用户名输入框美化 -->
          <el-form-item prop="username">
            <el-input
              ref="username"
              v-model="loginData.username"
              :placeholder="$t('用户名')"
              name="username"
              size="large"
              class="custom-input"
            >
            </el-input>
          </el-form-item>

          <!-- 密码输入框美化 -->
          <el-tooltip
            :visible="isCapslock"
            :content="$t('大写锁定已打开')"
            placement="right"
          >
            <el-form-item prop="password">
              <el-input
                v-model="loginData.password"
                :placeholder="$t('密码')"
                name="password"
                size="large"
                class="custom-input"
                :type="passwordVisible ? 'text' : 'password'"
                @keyup="checkCapslock"
                @keyup.enter="handleLoginSubmit"
              >
                <template #suffix>
                  <el-icon
                    class="cursor-pointer text-gray-400 hover:text-blue-500 transition-colors"
                    @click="passwordVisible = !passwordVisible"
                  >
                    <component :is="passwordVisible ? View : Hide" />
                  </el-icon>
                </template>
              </el-input>
            </el-form-item>
          </el-tooltip>

          <!-- 登录按钮美化 -->
          <el-button
            :loading="loading"
            type="primary"
            size="large"
            class="w-full h-12 text-lg font-semibold hover:transform hover:scale-105 transition-all duration-300 bg-gradient-to-r from-blue-500 to-teal-500 border-none"
            @click.prevent="handleLoginSubmit"
          >
            {{ $t("登录") }}
          </el-button>

          <!-- 记住密码区域美化 -->
          <div class="flex-x-between w-full py-4">
            <el-checkbox
              v-model="rememberme"
              class="rememberme text-gray-600 hover:text-blue-500 transition-colors"
            >
              记住密码
            </el-checkbox>
          </div>

          <!-- 其他登录方式美化 -->
          <el-divider class="my-6">
            <span class="text-gray-400">其他</span>
          </el-divider>

          <el-form-item class="flex justify-center">
            <span class="flex gap-4">
             <template v-for="(v, index) in data.oauthConfigs">
  <div
    :title="v.enable ? v.name : `${v.name}（未开启认证）`"
    class="inline-block transition-all duration-300"
    :class="v.enable ? 'hover:scale-110 cursor-pointer' : 'opacity-50 cursor-not-allowed'"
  >
    <a
      v-if="v.enable"
      :href="v.oauthUrl"
    >
      <el-image
        v-if="v.img === 'work_wechat'"
        :src="work_wechat"
        :alt="v.name"
        class="provider-img"
      />
    </a>
    <template v-else>
      <el-image
        v-if="v.img === 'work_wechat'"
        :src="work_wechat"
        :alt="v.name"
        class="provider-img"
      />
    </template>
  </div>
</template>

            </span>
          </el-form-item>
        </el-form>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref } from 'vue'
import { View, Hide } from '@element-plus/icons-vue'
import {GetOAuthList} from "@/api/user";

const logo = ref(new URL(`@/assets/logo.png`, import.meta.url).href);
//图片列表
const work_wechat = ref(new URL(`@/assets/images/work_wechat.png`, import.meta.url).href);
//


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
import {initCentrifuge} from '@/utils/centrifuge.js'
// 使用导入的依赖和库
const userStore = useUserStore();
const settingsStore = useSettingsStore();
const route = useRoute();
// 窗口高度
const {height} = useWindowSize();
// 国际化 Internationalization
const {t} = useI18n();


const bg_light = ref(new URL(`@/assets/images/login_bg.png`, import.meta.url).href)

const bg_dark = ref(new URL(`@/assets/images/login_dark_bg2.png`, import.meta.url).href)


const login_bg = computed(()=>{
  if(isDark)return bg_dark
  return bg_light
});



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

const passwordVisible = ref(false)

const appVersion = ()=>{
  return window["appVersion"]
}

const loginRules = computed(() => {
  return {
    username: [
      {
        required: true,
        trigger: "blur",
        message: t("请输入用户名"),
      },
    ],
    password: [
      {
        required: true,
        trigger: "blur",
        message: t("请输入密码"),
      },
      {
        min: 5,
        message: "密码长度不能少于5个字符",
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
          initCentrifuge()
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
        initCentrifuge()
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

<style lang="scss" scoped>
.login-container {
  @apply min-h-screen flex items-center justify-center bg-gradient-to-br from-blue-50 to-teal-50 dark:from-gray-900 dark:to-gray-800;

  .login-content {
    @apply rounded-xl shadow-2xl transition-all duration-500;

    .login-box {
      @apply bg-white/80 dark:bg-gray-800/80 backdrop-blur-lg;
    }
  }

  .login-input :deep(.el-input__wrapper) {
    @apply rounded-lg border-0 shadow-sm;
  }

  .provider-img {
    @apply w-10 h-10  shadow-md;
  }
}

.custom-input {
  :deep(.el-input__wrapper) {
    @apply bg-white/80 dark:bg-gray-800/80;
    @apply border border-gray-200 dark:border-gray-700;
    @apply shadow-sm;
    @apply rounded-lg;
    @apply transition-all duration-300;
    @apply backdrop-blur-sm;
    @apply px-4;

    // 输入框获得焦点时的样式
    &.is-focus {
      @apply border-blue-400 dark:border-blue-500;
      @apply shadow-md;
      @apply bg-white dark:bg-gray-800;
    }

    // 悬停效果
    &:hover {
      @apply border-blue-300 dark:border-blue-600;
      @apply shadow;
    }

    // 输入文本样式
    .el-input__inner {
      @apply text-gray-700 dark:text-gray-200;
      @apply placeholder:text-gray-400 dark:placeholder:text-gray-500;
      @apply h-11;
      @apply text-base;
    }

    // 密码可见性图标样式
    .el-input__suffix {
      @apply text-gray-400;
    }
  }
}

// 错误状态样式
:deep(.el-form-item.is-error) .custom-input {
  .el-input__wrapper {
    @apply border-red-300 dark:border-red-500;

    &:hover, &.is-focus {
      @apply border-red-400 dark:border-red-600;
    }

    .el-input__prefix-icon {
      @apply text-red-400 dark:text-red-500;
    }
  }
}

// 表单项间距
.el-form-item {
  @apply mb-6;
}
</style>

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
