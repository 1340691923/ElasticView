<template>
  <div class="notice">

    <el-dropdown class="h-full items-center justify-center" trigger="click">
      <el-badge v-if="noticeList.length > 0" :offset="[0, 15]" :value="noticeList.length" :max="99">
        <el-icon>
          <Bell />
        </el-icon>
      </el-badge>

      <div v-else>
        <el-icon>
          <Bell />
        </el-icon>
      </div>

      <template #dropdown>
        <div class="p-5">
          <template v-if="noticeList.length > 0">
            <div v-for="(item, index) in noticeList.slice(0, 5)" :key="index" class="w500px py-3">
              <div class="flex-y-center">
                <el-tag :type="item.level">{{ item.type }}</el-tag>
                <el-text
                  size="small"
                  class="w200px cursor-pointer !ml-2 !flex-1"
                  truncated
                  @click="handleReadNotice(item)"
                >
                  {{ item.title }}
                </el-text>

                <div class="text-xs text-gray">
                  {{ formatISODateTime(item.created) }}
                </div>
              </div>
            </div>
            <el-divider />
            <div class="flex-x-between">
              <el-link type="primary" underline="never" @click="handleViewMoreNotice">
                <span class="text-xs">查看更多</span>
                <el-icon class="text-xs">
                  <ArrowRight />
                </el-icon>
              </el-link>

            </div>
          </template>
          <template v-else>
            <div class="flex-center h150px w350px">
              <el-empty :image-size="50" description="暂无消息" >
                <el-link type="primary" underline="never" @click="handleViewMoreNotice">
                  <span class="text-xs">查看历史消息</span>
                  <el-icon class="text-xs">
                    <ArrowRight />
                  </el-icon>
                </el-link>
              </el-empty>

            </div>
          </template>
        </div>
      </template>
    </el-dropdown>

    <el-drawer
      append-to-body
      v-model="noticeDialogVisible"
      :title="currentNotice?.title ?? '通知详情'"
      :size="isMobile ? '100%' : '50%'"
      direction="rtl"
      :with-header="true"
      :destroy-on-close="false"
      custom-class="notice-drawer"
      :class="{'mobile-drawer': isMobile}"
      :before-close="handleCloseDrawer"
      :modal-class="isMobile ? 'drawer-mobile-modal' : ''"
    >
      <div class="drawer-content" :class="{'mobile-drawer-content': isMobile}">
        <el-descriptions  :column="1" border>
          <el-descriptions-item width="10px"  label="标题" label-class-name="description-label">
            <span class="description-value">{{ currentNotice.title }}</span>
          </el-descriptions-item>

          <el-descriptions-item  label="类型" label-class-name="description-label">
            <el-tag
              :type="currentNotice.level"
              class="notice-tag"
            >
              {{ currentNotice.type }}
            </el-tag>
          </el-descriptions-item>

          <el-descriptions-item label="来源" label-class-name="description-label">
            <span class="description-value">{{ currentNotice.source || '无' }}</span>
          </el-descriptions-item>

          <el-descriptions-item label="通知影响范围" label-class-name="description-label">
            <span class="description-value">
              <template v-if="currentNotice.target_type == 'all'">全部用户</template>
               <template v-if="currentNotice.target_type == 'roles'">指定权限组</template>
               <template v-if="currentNotice.target_type == 'users'">指定用户</template>
            </span>
          </el-descriptions-item>

          <el-descriptions-item label="是否定时任务" label-class-name="description-label">
            <span  class="description-value">{{currentNotice.is_task==1?'是':'否'}}</span>
          </el-descriptions-item>


          <el-descriptions-item label="发布时间" label-class-name="description-label">
            <span class="description-value">{{ formatISODateTime(currentNotice.created) }}</span>
          </el-descriptions-item>

          <el-descriptions-item label="内容" label-class-name="description-label">
            <div class="notice-content" v-html="currentNotice.content" />
          </el-descriptions-item>
        </el-descriptions>

        <div class="drawer-footer">
          <el-button
            type="primary"
            @click="noticeDialogVisible  = false"
            class="close-button"
          >
            关闭
          </el-button>

          <el-button
            class="close-button"
            type="info"
            v-if="currentNotice.btn_jump_url"
            @click="handleJump(currentNotice.btn_jump_type,currentNotice.btn_jump_url)"
          >
            <template v-if="currentNotice.btn_desc != ''">
              {{currentNotice.btn_desc}}
            </template>
            <template v-else>
              跳转
            </template>
          </el-button>
        </div>
      </div>
    </el-drawer>

  </div>
</template>

<script setup lang="ts">
import router from "@/router";
import {DeviceEnum} from "@/enums/DeviceEnum";
import {useAppStore, useUserStore, useNoticeStore} from "@/store";

const noticeDialogVisible = ref(false);
const currentNotice = ref(null);
const appStore = useAppStore();
const isMobile = computed(() => appStore.device === DeviceEnum.MOBILE);

import {SubscribeToChannel, publish,unsubscribeFromChannel} from "@/utils/centrifuge";
import {ElMessage} from "element-plus";

const userStore = useUserStore();
const noticeStore = useNoticeStore();

// 使用全局store的通知列表
const noticeList = computed(() => noticeStore.noticeList);

/**
 * 获取我的通知公告
 */
async function handleQuery() {
  try {
    await noticeStore.fetchNoticeList({
      page:1,page_size:100,read_type:2,
    });
  } catch (error) {
    console.error('获取通知列表失败:', error);
    ElMessage.error('获取通知列表失败');
  }
}

// 阅读通知公告
async function handleReadNotice(data) {
  // 先设置数据，再显示抽屉
  currentNotice.value = data;

  // 确保DOM更新完成
  await nextTick();
  noticeDialogVisible.value = true;

  // 标记已读
  try {
    await noticeStore.markAsRead([data.id]);
  } catch (error) {
    console.error('标记已读失败:', error);
  }
}

// 查看更多
function handleViewMoreNotice() {
  router.push({ path: "/notice" });
}

function handleJump(typ: string, url: string) {
  if (typ === 'internal') {
    router.push({ path: url })
  } else {
    window.open(url, '_blank')
  }
}


function formatISODateTime(isoString) {
  const date = new Date(isoString);

  // 检查日期是否有效
  if (isNaN(date.getTime())) {
    return '无效日期';
  }

  const year = date.getFullYear();
  const month = String(date.getMonth() + 1).padStart(2, '0');
  const day = String(date.getDate()).padStart(2, '0');
  const hours = String(date.getHours()).padStart(2, '0');
  const minutes = String(date.getMinutes()).padStart(2, '0');
  const seconds = String(date.getSeconds()).padStart(2, '0');

  return `${year}-${month}-${day} ${hours}:${minutes}:${seconds}`;
}

onMounted(async() => {
  await handleQuery();
  let userId = userStore.user.userId
  let roles = userStore.user.roles

  let channelArr = ["EvAllMsgChannel",`EvUserMsgChannel:${userId}`]
  for(let v of roles){
    channelArr.push(`EvRoleMsgChannel:${v}`)
  }

  for (let ch of channelArr){
    SubscribeToChannel(ch, (data) => {
      console.log("收到通知消息：", data);
      // 使用store添加通知
      noticeStore.addNotice(data);

      ElNotification({
        title: "您收到一条新的通知消息！",
        message: data.title,
        type: "success",
        position: "top-right",
      });
    });
  }
});

onBeforeUnmount(() => {
  let userId = userStore.user.userId
  let roles = userStore.user.roles

  let channelArr = ["EvAllMsgChannel",`EvUserMsgChannel:${userId}`]
  for(let v of roles){
    channelArr.push(`EvRoleMsgChannel:${v}`)
  }
  for (let ch of channelArr){
    unsubscribeFromChannel(ch);
  }
});

const handleCloseDrawer = () => {
  noticeDialogVisible.value = false;
};
</script>

<style lang="scss" scoped>
.notice {
:deep(.el-badge) {
  display: flex;
  align-items: center;
  justify-content: center;
  width: 100%;
  height: 100%;
}

.notice-detail {
&__wrapper {
   padding: 0 20px;
 }

&__meta {
   display: flex;
   align-items: center;
   margin-bottom: 16px;
   font-size: 13px;
   color: var(--el-text-color-secondary);
 }

&__publisher {
   margin-right: 24px;

i {
  margin-right: 4px;
}
}

&__content {
   max-height: 60vh;
   padding-top: 16px;
   margin-bottom: 24px;
   overflow-y: auto;
   border-top: 1px solid var(--el-border-color);

&::-webkit-scrollbar {
   width: 6px;
 }
}
}
}
</style>

<style scoped>
.notice-drawer {
  --el-drawer-padding-primary: 20px;
  --el-drawer-bg-color: var(--el-bg-color);
  --el-drawer-title-font-size: 18px;
}

.drawer-content {
  padding: 0 20px;
  height: calc(100% - 60px);
  overflow-y: auto;
}

.description-label {
  width: 100px;
  font-weight: bold;
  color: var(--el-text-color-regular);
}

.description-value {
  color: var(--el-text-color-primary);
}

.notice-tag {
  font-size: 14px;
  padding: 0 10px;
  height: 28px;
  line-height: 28px;
}

.notice-content {
  padding: 15px;
  background: var(--el-bg-color-page);
  border-radius: 4px;
  line-height: 1.6;
  color: var(--el-text-color-primary);
  border: 1px solid var(--el-border-color-light);
}

.notice-content :deep(p) {
  margin: 0 0 10px 0;
}

.notice-content :deep(a) {
  color: var(--el-color-primary);
  text-decoration: none;
}

.notice-content :deep(a:hover) {
  text-decoration: underline;
}

.drawer-footer {
  margin-top: 20px;
  text-align: right;
  padding: 10px 0;
  border-top: 1px solid var(--el-border-color-light);
}

.close-button {
  width: 100px;
}

/* 暗夜模式特定样式 */
:root.dark .notice-content {
  background: var(--el-bg-color-overlay);
  border-color: var(--el-border-color);
}

:root.dark .drawer-footer {
  border-top-color: var(--el-border-color);
}

.mobile-drawer {
  :deep(.el-drawer__header) {
    margin-bottom: 10px;
    padding: 10px 20px !important;
  }

  :deep(.el-drawer__body) {
    padding: 0 !important;
    height: calc(100% - 50px);
    overflow: auto;
    transition: none !important;
  }
}

.mobile-drawer-content {
  padding: 0 10px;

  :deep(.el-descriptions__cell) {
    padding: 8px !important;
  }

  :deep(.el-descriptions__label) {
    width: 90px !important;
  }

  .notice-content {
    padding: 10px;
  }

  .drawer-footer {
    margin-top: 15px;

    .close-button {
      width: auto;
      min-width: 80px;
    }
  }
}

.drawer-mobile-modal {
  animation: none !important;
  transition: none !important;
}

/* 添加抽屉组件的全局样式 */
:deep(.el-overlay) {
  animation-duration: 0.15s !important;
}

:deep(.el-drawer) {
  animation-duration: 0.15s !important;
}

</style>
