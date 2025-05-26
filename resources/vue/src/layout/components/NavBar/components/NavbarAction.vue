<template>
  <div class="flex">
    <template v-if="!isMobile">
      <div class="message nav-action-item-nohover">
        <el-container>
          <el-text type="primary">{{$t("数据源")}}:</el-text>
          <SelectLink style="margin-left: 1rem"></SelectLink>
        </el-container>
      </div>
    </template>

    <!-- 用户头像 -->
    <el-dropdown class="nav-action-item" trigger="click">
      <div class="flex-center h100% p10px">
        <span style="color:white" v-if="theme == 'dark'">{{ userStore.user.username }}</span>
        <span style="color:black" v-if="theme != 'dark'">{{ userStore.user.username }}</span>
      </div>
      <template #dropdown>
        <el-dropdown-menu>
          <el-dropdown-item><el-icon><Flag /></el-icon>版本:{{appVersion()}}</el-dropdown-item>
          <a target="_blank"  href="http://www.elastic-view.cn/" >
            <el-dropdown-item><el-icon><Postcard /></el-icon>{{ $t("官网") }}</el-dropdown-item>
          </a>
          <a target="_blank" href="https://txc.qq.com/products/666253" >
            <el-dropdown-item><el-icon><QuestionFilled /></el-icon>{{ $t("反馈") }}</el-dropdown-item>
          </a>
          <a target="_blank" href="https://raw.githubusercontent.com/1340691923/ElasticView/main/resources/show_img/weixin.jpg">
            <el-dropdown-item><el-icon><UserFilled /></el-icon>{{ $t("联系作者") }}</el-dropdown-item>
          </a>
          <el-dropdown-item divided @click="logout">
            <el-icon><SwitchButton /></el-icon>{{ $t("退出登录") }}
          </el-dropdown-item>
        </el-dropdown-menu>
      </template>
    </el-dropdown>

    <!-- 菜单定位按钮 -->
    <div
      class="nav-action-item flex items-center justify-center"
      :class="{ 'mobile-nav-action': isMobile }"
      @click="showMenuSearch = true"
    >
      <el-tooltip
        :content="$t('菜单定位')"
        placement="bottom"
      >
        <el-icon><Search /></el-icon>
      </el-tooltip>
    </div>

    <el-dialog
      v-model="showMenuSearch"
      :title="$t('菜单定位')"
      :width="isMobile ? '90%' : '60%'"
      height="100%"
      destroy-on-close
      class="menu-search-dialog"
      :modal-class="'menu-search-modal'"
    >
      <MenuSelect @select="handleMenuSelect"></MenuSelect>
    </el-dialog>

    <!-- 消息通知按钮 -->
    <div
      class="nav-action-item flex items-center justify-center"
      :class="{ 'mobile-nav-action': isMobile }"
    >
      <el-tooltip
        :content="$t('消息通知')"
        placement="bottom"
      >
        <div class="notice-icon-container">
          <NoticeDropdown></NoticeDropdown>
        </div>
      </el-tooltip>
    </div>

    <!-- 设置 -->
    <template v-if="defaultSettings.showSettings">
      <div
        class="nav-action-item flex items-center justify-center"
        :class="{ 'mobile-nav-action': isMobile }"
        @click="openSettings"
      >
        <svg-icon style="color:white" v-if="theme == 'dark'" icon-class="setting" class="setting-icon" />
        <svg-icon style="color:black" v-if="theme != 'dark'" icon-class="setting" class="setting-icon" />
      </div>
    </template>
  </div>
</template>
<script setup lang="ts">
import {
  useAppStore,
  useTagsViewStore,
  useUserStore,
  useSettingsStore,
} from "@/store";
import defaultSettings from "@/settings";
import { DeviceEnum } from "@/enums/DeviceEnum";
import { MessageTypeEnum, MessageTypeLabels } from "@/enums/MessageTypeEnum";
import NoticeDropdown from '@/components/NoticeDropdown/index.vue'
import {Search} from '@element-plus/icons-vue'

const appStore = useAppStore();
const tagsViewStore = useTagsViewStore();
const userStore = useUserStore();
const settingStore = useSettingsStore();
const searchMenuVisble = ref(false)
const route = useRoute();
const router = useRouter();

const isMobile = computed(() => appStore.device === DeviceEnum.MOBILE);

const { isFullscreen, toggle } = useFullscreen();

const activeTab = ref(MessageTypeEnum.MESSAGE);

const appVersion = ()=>{
  return window["appVersion"]
}

const showMenuSearch = ref(false)

const theme = computed(() => settingStore.theme);


const openSearchMenuVisble = ()=>{
  searchMenuVisble.value = true
}

/* 注销 */
function logout() {
  ElMessageBox.confirm("确定注销并退出系统吗？", "提示", {
    confirmButtonText: "确定",
    cancelButtonText: "取消",
    type: "warning",
    lockScroll: false,
  }).then(() => {
    userStore
      .logout()
      .then(() => {
        tagsViewStore.delAllViews();
      })
      .then(() => {
        router.push(`/login?redirect=${route.fullPath}`);
      });
  });
}

const handleMenuSelect = () => {
  showMenuSearch.value = false;
}

const openSettings = () => {
  settingStore.settingsVisible = true;
}
</script>
<style lang="scss" scoped>
.nav-action-item {
  display: inline-flex;
  align-items: center;
  justify-content: center;
  min-width: 40px;
  height: $navbar-height;
  cursor: pointer;

  &:hover {
    background: rgb(0 0 0 / 10%);
  }
}

.nav-action-item-nohover {
  display: inline-block;
  min-width: 40px;
  height: $navbar-height;
  line-height: $navbar-height;
  text-align: center;
  cursor: pointer;
}

:deep(.message .el-badge__content.is-fixed.is-dot) {
  top: 5px;
  right: 10px;
}

:deep(.el-divider--horizontal) {
  margin: 10px 0;
}

.dark .nav-action-item:hover {
  background: rgb(255 255 255 / 20%);
}

.layout-top .nav-action-item,
.layout-mix .nav-action-item {
  color: #fff;
}

.search-btn-wrapper {
  @apply flex items-center h-full;
}

.search-icon-btn {
  @apply flex items-center justify-center;
  @apply bg-transparent;
  @apply border border-gray-300 dark:border-gray-600;
  @apply text-gray-600 dark:text-gray-400;
  @apply hover:bg-gray-100 dark:hover:bg-gray-700;
  @apply transition-all duration-200;

  :deep(.el-icon) {
    @apply text-base;
  }

  &:hover {
    @apply border-blue-400 dark:border-blue-500;
    @apply text-blue-500 dark:text-blue-400;
    @apply shadow-sm;
  }
}

.square-icon-btn {
  @apply flex items-center justify-center;
  @apply bg-transparent;
  @apply border border-gray-300 dark:border-gray-600;
  @apply text-gray-600 dark:text-gray-400;
  @apply hover:bg-gray-100 dark:hover:bg-gray-700;
  @apply transition-all duration-200;
  @apply rounded-sm;
  @apply min-h-[32px] min-w-[32px];

  :deep(.el-icon) {
    @apply text-base;
  }

  &:hover {
    @apply border-blue-400 dark:border-blue-500;
    @apply text-blue-500 dark:text-blue-400;
    @apply shadow-sm;
  }
}

:deep(.menu-search-dialog) {
  z-index: 3000 !important;

  .el-dialog__body {
    @apply p-4;
  }

}

.setting-icon {
  @apply w-5 h-5;
}

:global(.menu-search-modal) {
  z-index: 2999 !important;
}

:global(.el-dialog__wrapper) {
  z-index: 3000 !important;
}

.flex-center {
  @apply flex items-center justify-center;
}

.notice-icon-container {
  @apply flex items-center justify-center;
  @apply h-[40px]; /* 设置固定高度与导航栏高度一致 */

  :deep(.notice) {
    @apply h-full flex items-center justify-center;

    :deep(.el-dropdown) {
      @apply h-full flex items-center justify-center;

      :deep(.el-badge) {
        @apply transform translate-y-0;
      }

      :deep(.el-icon) {
        @apply flex items-center justify-center;
        @apply text-base;
      }
    }
  }
}

.mobile-nav-action {
  @apply min-w-[36px]; /* 移动端下减小最小宽度 */
  @apply px-1; /* 减小内边距 */
}
</style>
