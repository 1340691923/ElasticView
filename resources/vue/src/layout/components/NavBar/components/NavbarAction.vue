<template>
  <div class="flex">
    <template v-if="!isMobile">

      <div class="message nav-action-item-nohover">
        <el-container>
          <el-text type="primary" >数据源:</el-text><SelectLink style="margin-left: 1rem" ></SelectLink>
        </el-container>
      </div>

      <div style="margin-left: 1rem" class="message nav-action-item-nohover" >
        <MenuSelect  class="message nav-action-item-nohover"></MenuSelect>
<!--        <el-button @click="openSearchMenuVisble" :icon="Search" circle />-->
      </div>
<!--      <div class="message nav-action-item-nohover" >
        <el-popover
          placement="bottom"
          title="请输入您需要快速到达的功能"
          :width="200"
          trigger="click"
        >
          <template #reference>
            <el-button  :icon="Search" circle />
          </template>
          <MenuSelect  class="message nav-action-item-nohover"></MenuSelect>

        </el-popover>
      </div>-->
      <el-drawer
        size="30%"
        v-model="searchMenuVisble"
        title="请选择您需要快速到达的功能"
      >
        <MenuSelect  class="message nav-action-item-nohover"></MenuSelect>
      </el-drawer>
<!--      <MenuSelect  class="message nav-action-item-nohover"></MenuSelect>-->

<!--      <SelectLink style="margin-left: 1rem" class="message nav-action-item-nohover"></SelectLink>-->
      <!-- 消息通知 -->
<!--      <el-dropdown class="message nav-action-item" trigger="click">
        <el-badge is-dot>
          <div class="flex-center h100% p10px">
            <i-ep-bell />
          </div>
        </el-badge>
        <template #dropdown>
          <div class="px-5 py-2">
            <el-tabs v-model="activeTab">
              <el-tab-pane
                v-for="(label, key) in MessageTypeLabels"
                :label="label"
                :name="key"
                :key="key"
              >
                <div
                  class="w-[380px] py-2"
                  v-for="message in getFilteredMessages(key)"
                  :key="message.id"
                >
                  <el-link type="primary">
                    <el-text class="w-350px" size="default" truncated>
                      {{ message.title }}
                    </el-text>
                  </el-link>
                </div>
              </el-tab-pane>
            </el-tabs>
            <el-divider />
            <div class="flex-x-between">
              <el-link type="primary" :underline="false">
                <span class="text-xs">查看更多</span>
                <el-icon class="text-xs"><ArrowRight /></el-icon>
              </el-link>
              <el-link type="primary" :underline="false">
                <span class="text-xs">全部已读</span>
              </el-link>
            </div>
          </div>
        </template>
      </el-dropdown>-->

    </template>

    <!-- 用户头像 -->
    <el-dropdown class="nav-action-item" trigger="click">
      <div class="flex-center h100% p10px">

        <!-- todo... -->
<!--        <img
          :src="userStore.user.avatar + '?imageView2/1/w/80/h/80'"
          class="rounded-full mr-10px w24px w24px"
        />-->
        <span>{{ userStore.user.username }}</span>
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

<!--          <el-dropdown-item><el-icon><UserFilled /></el-icon>{{ $t("个人信息") }}</el-dropdown-item>-->
          <a target="_blank" href="https://raw.githubusercontent.com/1340691923/ElasticView/main/resources/show_img/weixin.jpg">
            <el-dropdown-item><el-icon><UserFilled /></el-icon>{{ $t("联系作者") }}</el-dropdown-item>
          </a>


          <el-dropdown-item divided @click="logout">
            <el-icon><SwitchButton /></el-icon>{{ $t("退出登录") }}
          </el-dropdown-item>
        </el-dropdown-menu>
      </template>
    </el-dropdown>

    <!-- 设置 -->
    <template v-if="defaultSettings.showSettings">
      <div class="nav-action-item" @click="settingStore.settingsVisible = true">
        <svg-icon icon-class="setting" />
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

const showSearchMenuVisble = ref(false)

const showSearchMenu = ()=>{
  showSearchMenuVisble.value = true
}

const messages = ref([
 /* {
    id: 1,
    title: "系统升级通知：服务器将于今晚12点进行升级维护，请提前保存工作内容。",
    type: MessageTypeEnum.MESSAGE,
  },
  {
    id: 2,
    title: "新功能发布：我们的应用程序现在支持多语言功能。",
    type: MessageTypeEnum.MESSAGE,
  },
  {
    id: 3,
    title: "重要提醒：请定期更改您的密码以保证账户安全。",
    type: MessageTypeEnum.MESSAGE,
  },
  {
    id: 4,
    title: "通知：您有一条未读的系统消息，请及时查看。",
    type: MessageTypeEnum.NOTICE,
  },
  {
    id: 5,
    title: "新订单通知：您有一笔新的订单需要处理。",
    type: MessageTypeEnum.NOTICE,
  },
  {
    id: 6,
    title: "审核提醒：您的审核请求已被批准。",
    type: MessageTypeEnum.NOTICE,
  },
  { id: 7, title: "待办事项：完成用户权限设置。", type: MessageTypeEnum.TODO },
  { id: 8, title: "待办事项：更新产品列表。", type: MessageTypeEnum.TODO },
  { id: 9, title: "待办事项：备份数据库。", type: MessageTypeEnum.TODO },*/
]);

const getFilteredMessages = (type: MessageTypeEnum) => {
  return messages.value.filter((message) => message.type === type);
};

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
</script>
<style lang="scss" scoped>
.nav-action-item {
  display: inline-block;
  min-width: 40px;
  height: $navbar-height;
  line-height: $navbar-height;
  color: var(--el-text-color);
  text-align: center;
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
  color: var(--el-text-color);
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
</style>
