<template>
  <div class="app-container">
    <div class="mobile-filter-toggle" v-if="isMobile" @click="toggleFilter">
      <el-button type="primary" class="toggle-button">
        {{ isFilterVisible ? '收起筛选' : '展开筛选' }}
        <el-icon class="toggle-icon" :class="{ 'is-reverse': isFilterVisible }">
          <ArrowDown />
        </el-icon>
      </el-button>
    </div>

    <div
      class="search-container"
      :class="{
        'is-mobile': isMobile,
        'is-collapsed': isMobile && !isFilterVisible
      }"
    >
      <el-form :inline="true" class="search-form">
        <el-form-item label="">
          <el-input
            v-model="input.search_txt"
            clearable
            placeholder="插件名/描述:"

          >
            <template #prefix>
              <el-icon><Search /></el-icon>
            </template>
          </el-input>
        </el-form-item>

        <el-form-item label="排序:">
          <div class="sort-wrapper">
            <el-select v-model="input.order_by_col" class="custom-select">
              <el-option label="star次数" value="star_cnt" />
              <el-option label="下载人数" value="download_user_cnt" />
              <el-option label="下载次数" value="download_cnt" />
              <el-option label="最后更新时间" value="publish_time" />
          </el-select>
            <el-button
              class="sort-button"
              :type="input.order_by_desc ? 'default' : 'primary'"
              @click="toggleSort"
            >
              <el-icon>
                <component :is="input.order_by_desc ? SortDown : SortUp" />
              </el-icon>
            </el-button>
          </div>
        </el-form-item>

        <el-form-item label="安装状态:">
          <el-select v-model="input.has_download_type" class="custom-select">
            <el-option label="全部" :value="null" />
            <el-option label="未安装" :value="false" />
            <el-option label="已安装" :value="true" />
          </el-select>
        </el-form-item>

        <el-form-item>
              <el-button
            type="primary"
            :icon="Search"
            class="search-button"
                @click="getPluginMarket"
          >
            {{ $t('搜索') }}
              </el-button>
        </el-form-item>
      </el-form>
    </div>

    <div class="plugin-container" v-loading="pluginListLoading">
      <div class="plugin-list">
        <div
          v-for="(item, index) in pluginListData.list"
          :key="item.id"
          class="plugin-card-wrapper"
        >
          <el-card
            :class="['plugin-card', {'is-dark': settingsStore.theme === ThemeEnum.DARK}]"
            @click="lookPluginInfo(item)"
            :body-style="{ padding: '0' }"
          >
            <div class="plugin-header">
              <img :src="item.logo" class="plugin-logo" loading="lazy">
              <h3 class="plugin-title">{{ item.plugin_name }}</h3>
            </div>

            <div class="plugin-body">
              <p class="plugin-description">
                <el-tooltip :content="item.describe" placement="top">
                  {{ truncatedText(item.describe, 50) }}
              </el-tooltip>
              </p>

              <div class="plugin-tags">
                <el-tag v-if="item.has_download" type="success" effect="light">已安装</el-tag>
                <el-tag v-if="item.buy_coin_num > 0" type="primary" effect="light">
                  所需ev币: {{item.buy_coin_num}}个
                </el-tag>
                <el-tag v-else type="info" effect="light">免费</el-tag>
              </div>

              <div class="plugin-stats">
                <el-tooltip content="安装次数/人数" placement="top">
                  <div class="stat-item">
                    <el-icon><Download /></el-icon>
                    {{item.download_cnt}}/{{item.download_user_cnt}}
              </div>
                </el-tooltip>

                <el-tooltip content="Star数" placement="top">
                  <div class="stat-item">
                    <el-icon><Star /></el-icon>
                    {{item.star_cnt}}
              </div>
                </el-tooltip>
            </div>
          </div>

            <div class="plugin-footer">
              <span class="publish-time">{{item.publish_time}}</span>
              <el-button
                :type="item.star_state === 1 ? 'primary' : 'default'"
                :icon="item.star_state === 1 ? StarFilled : Star"
                circle
                @click.stop="starPlugin(item, index)"
              />
          </div>
          </el-card>
          </div>
      </div>

      <!-- 添加底部间距，为粘性分页提供空间 -->
      <div class="pagination-spacer"></div>
    </div>

    <!-- 粘性分页组件 -->
    <div class="pagination-container">
      <el-pagination
        v-model:current-page="input.page"
        v-model:page-size="input.limit"
        :total="pluginListData.count"
        :page-sizes="[10, 20, 30, 50]"
        background
        :layout="isMobile?'pager':'total, sizes, prev, pager, next, jumper'"
        @size-change="handlePluginListSizeChange"
        @current-change="handlePluginListPageChange"
      />
    </div>

    <el-drawer
      class="plugin-detail-drawer"
      v-model="dialog.visible"
      title="插件详情"
      :size="isMobile ? '100%' : '100%'"
      :append-to-body="true"
      :with-header="true"
    >
      <div class="plugin-detail">
        <el-card v-loading="installLoading" class="detail-card">
          <div class="detail-header">
            <div class="header-main">
              <div class="plugin-info" :class="{'plugin-info-mobile': isMobile}">
                <img :src="publishInput.pluginData.logo" class="plugin-logo">
                <div class="info-content">
                  <h1 class="plugin-title">{{publishInput.pluginData.plugin_name}}</h1>
                  <p class="plugin-desc">{{publishInput.pluginData.describe}}</p>
                </div>
                <div class="plugin-actions" :class="{'plugin-actions-mobile': isMobile}">
                  <el-button
                    type="primary"
                    v-if="!publishInput.pluginData.has_download"
                    @click="installPlugin(publishListData.list[0]?.version)"
                  >
                    {{isMobile ? '安装' : '安装最新版本'}}
                  </el-button>
                  <el-button
                    type="danger"
                    v-if="publishInput.pluginData.has_download"
                    @click="unInstall"
                  >
                    卸载
                  </el-button>
                </div>
              </div>

              <div class="tag-group" :class="{'tag-group-mobile': isMobile}">
                <el-tag v-if="publishInput.pluginData.has_download" type="success">已安装</el-tag>
                <el-tag>开发者: {{publishInput.pluginData.realname}}</el-tag>
                <el-tag :type="publishInput.pluginData.buy_coin_num > 0 ? 'primary' : 'info'">
                  {{publishInput.pluginData.buy_coin_num > 0 ? `所需ev币: ${publishInput.pluginData.buy_coin_num}个` : '免费'}}
                </el-tag>
                  </div>

              <div class="stats-group" :class="{'stats-group-mobile': isMobile}">
                <div class="stat-item">
                  <el-icon><Download /></el-icon>
                  <span>安装: {{publishInput.pluginData.download_cnt}}次/{{publishInput.pluginData.download_user_cnt}}人</span>
                  </div>
                <div class="stat-item">
                  <el-icon><Star /></el-icon>
                  <span>Star: {{publishInput.pluginData.star_cnt}}</span>
                </div>
                <div class="stat-item">
                  <el-icon><Clock /></el-icon>
                  <span>更新时间: {{publishInput.pluginData.publish_time}}</span>
              </div>
                    </div>
                  </div>
                </div>

          <el-tabs v-model="tabShowType" class="detail-tabs">
                <el-tab-pane label="介绍" name="readme">
              <div class="readme-content">
                  <mark-down-view :content="publishInput.pluginData.readme"></mark-down-view>
              </div>
                </el-tab-pane>

                <el-tab-pane label="版本列表" name="versions">
              <div class="version-list" :class="{'version-list-mobile': isMobile}">
                <el-timeline>
                  <el-timeline-item
                    v-for="(item,index) in limitedVersionList"
                    :key="index"
                    :timestamp="item.update_time"
                    placement="top"
                    class="version-timeline-item"
                  >
                    <el-card class="version-card">
                        <template #header>
                        <div class="version-header">
                          <span class="version-tag">v{{item.version}}</span>
                          <div class="version-actions">

                            <el-button
                              @click.stop="installPlugin(item.version)"
                              type="warning"
                              :icon="Download"
                              circle
                            />
                            </div>
                          </div>
                        </template>
                      <mark-down-view :content="item.changelog" class="changelog"></mark-down-view>

                      </el-card>
                    </el-timeline-item>
                  </el-timeline>
              </div>
                </el-tab-pane>

                <el-tab-pane label="评论" name="comments">
                  <div class="comments-list" v-loading="commentsLoading">
                    <div v-if="comments.length === 0" class="empty-comments">
                      <el-empty description="暂无评论" />
                    </div>
                    <div v-else class="comments-container">
                      <!-- 使用递归评论组件 -->
                      <comment-item
                        v-for="comment in comments"
                        :key="comment.id"
                        :comment="comment"
                        :plugin-id="publishInput.pluginData.id"
                        @refresh="getComments"
                      />
                    </div>

                                        <!-- 添加新评论 -->
                    <div class="add-comment-container">
      <el-input
        v-model="newCommentContent"
        type="textarea"
        :rows="3"
        placeholder="发表评论..."
        maxlength="200"
        show-word-limit
      />

      <div class="emoji-container">
        <el-button
          type="text"
          size="small"
          class="emoji-toggle"
          @click="showMainEmojiPicker = !showMainEmojiPicker"
        >
          😊 表情
        </el-button>
        <div v-if="showMainEmojiPicker" style="z-index:300000" class="emoji-picker">
          <div
            v-for="emoji in emojiList"
            :key="emoji"
            class="emoji-item"
            @click="insertMainEmoji(emoji)"
          >
            {{ emoji }}
          </div>
        </div>
      </div>

      <div class="comment-actions">
        <el-button type="primary" @click="submitComment">发表评论</el-button>
      </div>
    </div>

                    <!-- 评论区底部空间 -->
                    <div class="comments-spacer"></div>
                  </div>
                </el-tab-pane>
              </el-tabs>
        </el-card>
      </div>
    </el-drawer>

    <import-ev-key v-model:visible="importEvkeyDialogVisible" ></import-ev-key>
  </div>
</template>

<script lang="ts" setup>
import {useSettingsStore} from "@/store";
import {ThemeEnum} from "@/enums/ThemeEnum";
import {LikeComment,AddComment,PluginMarket, GetPluginInfo,InstallPlugin,StarPlugin,UnInstallPlugin,UploadPlugin,ListComments} from "@/api/plugins";
import MarkDownView from '@/components/MarkDownView/index.vue'
import ImportEvKey from '@/components/ImportEvKey/index.vue'
import CommentItem from '@/components/CommentItem.vue'

import {Star,StarFilled,Download, Search, SortUp, SortDown, Clock, ArrowDown, Pointer} from '@element-plus/icons-vue'
import {useAppStore} from "@/store";
import {DeviceEnum} from "@/enums/DeviceEnum";
import {getToken} from "@/utils/auth";
const appStore = useAppStore()



const isMobile = computed(() => appStore.device === DeviceEnum.MOBILE);

const dialog = reactive({
  visible: false,
})

const settingsStore = useSettingsStore();

const pluginListLoading = ref(false)

const publishLoading = ref(false)

const pluginListData = reactive({
  count: 0,
  list: []
})

const publishListData = reactive({
  count: 0,
  list: []
})

const tabShowType = ref("readme")

const pluginBoxClass = computed(() => {
  return {
    'plugin-box-black': settingsStore.theme === ThemeEnum.DARK,
    'plugin-box-light': settingsStore.theme !== ThemeEnum.DARK,
    // 你可以添加更多的条件
  }
})

const input = reactive({
  search_txt:'',
  order_by_col:'star_cnt',
  order_by_desc:true,
  has_download_type:null,
  page: 1,
  limit: 10,
})

const publishInput = reactive({
  pluginData: {
    id: 0,
    plugin_alias: "",
    plugin_name: "",
    user_id: 0,
    describe: "",
    plugin_lang: "",
    readme: "",
    create_time: "",
    update_time: "",
    state: 2,
    logo: "",
    msg: "",
    download_cnt: 0,
    star_cnt: 0,
    download_user_cnt: 0,
    buy_coin_num:0
  },
  publish_id: 0,
  page: 1,
  limit: 10
})

const importEvkeyDialogVisible = ref(false);

const openImportEvkeyDialogVisible = () => {
  importEvkeyDialogVisible.value = true;
};

const getBgColor = computed(()=>{
  return settingsStore.theme === ThemeEnum.DARK ? 'rgb(24, 27, 31)': ''
})

const unInstall = async ()=>{
  installLoading.value = true
  let res = await UnInstallPlugin({
    plugin_id:publishInput.pluginData.plugin_alias,
  })
  installLoading.value = false
  if (res.code != 0) {
    ElMessage.error({
      type: 'error',offset:60,
      message: res.msg
    })
    return
  }
  publishInput.pluginData.has_download = false
  ElMessage.success({
    type: 'success',offset:60,
    message: res.msg
  })
}

const getColor = computed(() => {
  return settingsStore.theme === ThemeEnum.DARK ? 'rgb(204, 204, 220)' : 'black'
})

const getColor2 = computed(() => {
  return settingsStore.theme === ThemeEnum.DARK ? 'rgb(122,122,133)' : 'gray'
})

const truncatedText = (text, maxLength) => {
  return text.length > maxLength ? text.slice(0, maxLength) + '...' : text;
};

const handlePluginListSizeChange = (v: number) => {
  input.limit = v
  getPluginMarket()
}

const handlePluginListPageChange = (v: number) => {
  input.page = v
  getPluginMarket()
}

const handlePublishSizeChange = (v: number) => {
  input.limit = v
  getPluginMarket()
}

const handlePublishPageChange = (v: number) => {
  input.page = v
  getPluginInfo()
}

const getPluginMarket = async () => {

  pluginListLoading.value = true
  const res = await PluginMarket(input)
  pluginListLoading.value = false
  if (res.code != 0) {
    ElMessage.error({
      type: 'error',offset:60,
      message: res.msg
    })
    return
  }

  pluginListData.count = res.data.count
  pluginListData.list = res.data.list
}

const lookPluginInfo = async (row) => {
  publishInput.pluginData = row
  dialog.visible = true

  await getPluginInfo()
  await getComments()
}

const installLoading = ref(false)

const confirmReloadPage = (pluginId) =>{
  ElMessageBox.confirm(pluginId+'插件安裝成功,是否立即刷新页面?', '是否刷新', {
    confirmButtonText: '确认',
    cancelButtonText: '取消',
    type: 'warning'
  })
    .then(async() => {
     window.location.reload()
    })
    .catch(err => {
      console.error(err)
    })


}

const installPlugin = async (version) => {
  installLoading.value = true

  let res = await InstallPlugin({
    plugin_id:publishInput.pluginData.plugin_alias,
    version:version
  })
  installLoading.value = false

  if (res.msg.indexOf('请前往')!==-1){
    openImportEvkeyDialogVisible()
    return
  }

  if (res.code != 0) {
    ElMessage.error({
      type: 'error',offset:60,
      message: res.msg
    })
    return
  }
  publishInput.pluginData.has_download = true

  if(res.msg.indexOf("刷新")!==-1){
    await confirmReloadPage(publishInput.pluginData.plugin_name)
  }else{
    ElMessage.success({
      type: 'success',offset:60,
      message: res.msg
    })
  }

}

const starPlugin = async (item,index) => {
  let res = await StarPlugin({
    plugin_id:item.id,
  })
  if (res.code != 0) {

    if (res.msg.indexOf('请前往')!==-1){

      openImportEvkeyDialogVisible()
      return
    }

    ElMessage.error({
      type: 'error',offset:60,
      message: res.msg
    })
    return
  }
  let msg = 'star成功'
  if(item.star_state == 1){
    msg = '取消star成功'
    pluginListData.list[index].star_state = 2
    pluginListData.list[index].star_cnt =  pluginListData.list[index].star_cnt - 1
  }else{
    pluginListData.list[index].star_state = 1
    pluginListData.list[index].star_cnt =  pluginListData.list[index].star_cnt + 1
  }

  ElMessage.success({
    type: 'success',offset:60,
    message: msg
  })
}

const getPluginInfo = async () => {
  publishLoading.value = true
  const res = await GetPluginInfo({
    page: publishInput.page,
    limit: publishInput.limit,
    plugin_id: publishInput.pluginData.id,
  })
  publishLoading.value = false

  if (res.code != 0) {
    ElMessage.error({
      type: 'error',offset:60,
      message: res.msg
    })
    return
  }
  if(res.data.list == null)res.data.list = []
  publishListData.count = res.data.count
  publishListData.list = res.data.list

}

const handleCloseDialog = () => {
  dialog.visible = false
}

const toggleSort = () => {
  input.order_by_desc = !input.order_by_desc;
  getPluginMarket(); // 切换后自动刷新列表
}

const isFilterVisible = ref(false)

const toggleFilter = () => {
  isFilterVisible.value = !isFilterVisible.value
}

const commentsLoading = ref(false)
const comments = ref([])
const newCommentContent = ref('')
const showMainEmojiPicker = ref(false)
const replyState = reactive({
  showReplyInput: false,
  replyToId: 0,
  replyContent: '',
  replyToName: ''
})

// emoji表情列表
const emojiList = [
  '😀', '😃', '😄', '😁', '😆', '😅', '😂', '🤣', '😊', '😇',
  '🙂', '🙃', '😉', '😌', '😍', '🥰', '😘', '😗', '😙', '😚',
  '😋', '😛', '😝', '😜', '🤪', '🤨', '🧐', '🤓', '😎', '🤩',
  '🥳', '😏', '😒', '😞', '😔', '😟', '😕', '🙁', '☹️', '😣',
  '😖', '😫', '😩', '🥺', '😢', '😭', '😤', '😠', '😡', '🤬',
  '🤯', '😳', '🥵', '🥶', '😱', '😨', '😰', '😥', '😓', '🤗',
  '🤔', '🤭', '🤫', '🤥', '😶', '😐', '😑', '😬', '🙄', '😯',
  '😦', '😧', '😮', '😲', '🥱', '😴', '🤤', '😪', '😵', '🤐',
  '🥴', '🤢', '🤮', '🤧', '😷', '🤒', '🤕', '🤑', '🤠', '👍',
  '👎', '👏', '🙏', '🤝', '💪', '❤️', '💔', '💯', '✨', '🔥'
]

// 向主评论框插入emoji
const insertMainEmoji = (emoji) => {
  newCommentContent.value += emoji
  showMainEmojiPicker.value = false // 选择后关闭emoji选择器
}

const getComments = async () => {
  commentsLoading.value = true
  try {
    const res = await ListComments({
      plugin_id: publishInput.pluginData.id
    })

    if (res.code !== 0) {
      ElMessage.error({
        type: 'error',
        offset: 60,
        message: res.msg
      })
      return
    }

    comments.value = res.data || []
  } catch (error) {
    console.error('获取评论失败', error)
    ElMessage.error({
      type: 'error',
      offset: 60,
      message: '获取评论失败'
    })
  } finally {
    commentsLoading.value = false
  }
}

const showReplyInput = (commentId) => {
  // 找到评论或回复的用户名
  let replyToName = '';

  // 查找主评论中是否有匹配的评论ID
  const mainComment = comments.value.find(c => c.id === commentId);
  if (mainComment) {
    replyToName = mainComment.realname;
  } else {
    // 如果在主评论中没找到，则在子评论中查找
    for (const comment of comments.value) {
      if (comment.children && comment.children.length > 0) {
        const childComment = comment.children.find(c => c.id === commentId);
        if (childComment) {
          replyToName = childComment.realname;
          break;
        }
      }
    }
  }

  replyState.showReplyInput = true;
  replyState.replyToId = commentId;
  replyState.replyContent = '';
  replyState.replyToName = replyToName;
}

const cancelReply = () => {
  replyState.showReplyInput = false
  replyState.replyToId = 0
  replyState.replyContent = ''
}

const submitComment = async () => {
  if (!newCommentContent.value.trim()) {
    ElMessage.warning({
      type: 'warning',
      offset: 60,
      message: '评论内容不能为空'
    })
    return
  }

  commentsLoading.value = true
  try {
    const res = await AddComment({
      plugin_id: publishInput.pluginData.id,
      content: newCommentContent.value,
      parent_id: 0
    })

    if (res.msg.indexOf('请前往')!==-1){
      openImportEvkeyDialogVisible()
      return
    }


    if (res.code !== 0) {
      ElMessage.error({
        type: 'error',
        offset: 60,
        message: res.msg
      })
      return
    }

    ElMessage.success({
      type: 'success',
      offset: 60,
      message: '评论成功'
    })

    newCommentContent.value = ''
    getComments() // 刷新评论列表
  } catch (error) {
    console.error('评论失败', error)
    ElMessage.error({
      type: 'error',
      offset: 60,
      message: '评论失败'
    })
  } finally {
    commentsLoading.value = false
  }
}

const limitedVersionList = computed(() => {
  return publishListData.list.slice(0, 100)
})

onMounted(() => {
  getPluginMarket()
})

</script>

<style lang="scss" scoped>
.mobile-filter-toggle {
  @apply mb-4;

  .toggle-button {
    @apply w-full;
    @apply flex items-center justify-center gap-2;

    .toggle-icon {
      @apply transition-transform duration-300;

      &.is-reverse {
        @apply transform rotate-180;
      }
    }
  }
}

.search-container {
  @apply mb-6 p-4;

  @apply backdrop-blur-sm;
  @apply rounded-lg;
  @apply shadow-sm;
  @apply transition-all duration-300;

  &.is-mobile {
    @apply overflow-hidden;

    &.is-collapsed {
      @apply h-0 p-0 mb-0;
      @apply opacity-0;
    }

    .search-form {
      @apply flex-col;

      .el-form-item {
        @apply w-full;

        :deep(.el-form-item__content) {
          @apply w-full;
          @apply flex;

          .custom-input,
          .custom-select,
          .sort-wrapper {
            @apply flex-1;
          }
        }
      }
    }
  }

  .search-form {
    @apply flex flex-wrap items-center gap-4;
  }
}

.custom-input {
  @apply w-64;
  :deep(.el-input__wrapper) {
    @apply shadow-sm;
    @apply transition-all duration-300;

    &:hover {
      @apply shadow;
    }
  }
}

.custom-select {
  @apply w-32;
  :deep(.el-input__wrapper) {
    @apply shadow-sm;
    @apply transition-all duration-300;

    &:hover {
      @apply shadow;
    }
  }
}

.search-button {
  @apply transition-all duration-300;
  @apply hover:scale-105;
}

.plugin-container {
  @apply my-6;
}

.plugin-list {
  @apply grid gap-6;
  grid-template-columns: repeat(auto-fill, minmax(300px, 1fr));
}

.pagination-spacer {
  @apply h-24; /* 提供分页组件空间 */
}

.plugin-card-wrapper {
  @apply transition-all duration-300;
  @apply hover:translate-y-[-4px];
  @apply cursor-pointer;
}

.plugin-card {
  @apply h-full;
  @apply cursor-pointer;
  @apply overflow-hidden;
  @apply transition-all duration-300;

  &:hover {
    @apply shadow-lg;
  }

  &.is-dark {
    @apply bg-transparent;
    @apply border-gray-700;

    :deep(.el-card__body) {
      @apply bg-transparent;
    }
  }
}

.plugin-header {
  @apply p-4;
  @apply flex items-center gap-4;
  @apply border-b border-gray-200 dark:border-gray-700;
  @apply cursor-pointer;

  .plugin-logo {
    @apply w-12 h-12;
    @apply object-contain;
    @apply rounded-lg;
  }

  .plugin-title {
    @apply m-0;
    @apply text-lg font-medium;
    @apply text-gray-900 dark:text-gray-100;
  }
}

.plugin-body {
  @apply p-4;
  @apply flex flex-col gap-4;
  @apply cursor-pointer;

  .plugin-description {
    @apply m-0;
    @apply text-gray-600 dark:text-gray-400;
    @apply line-clamp-2;
  }

  .plugin-tags {
    @apply flex flex-wrap gap-2;
  }

  .plugin-stats {
    @apply flex items-center gap-4;

    .stat-item {
      @apply flex items-center gap-1;
      @apply text-sm text-gray-500 dark:text-gray-400;
    }
  }
}

.plugin-footer {
  @apply p-4;
  @apply flex items-center justify-between;
  @apply border-t border-gray-200 dark:border-gray-700;
  @apply cursor-pointer;

  .publish-time {
    @apply text-sm text-gray-500 dark:text-gray-400;
  }
}


.pagination-container {
  @apply mt-4 flex justify-end;
  @apply overflow-x-auto;

:deep(.el-pagination) {

  @apply rounded-lg p-2;
  @apply shadow-sm;
  @apply min-w-fit;

@media (max-width: 768px) {
  @apply w-full;
  @apply flex justify-center;
  @apply text-sm;

  .el-pager {
    @apply flex-wrap;
  }

  .btn-prev,
  .btn-next {
    @apply min-w-[24px];
  }

  li {
    @apply min-w-[24px];
  }
}
}
}

// 添加暗色模式适配
.dark {
  .search-card {
    @apply bg-transparent;
  }

  .plugin-card {
    @apply bg-transparent;
    @apply border-gray-700;
  }

  :deep(.el-card) {
    @apply bg-transparent;

    .el-card__body {
      @apply bg-transparent;
    }
  }
}

.sort-wrapper {
  @apply flex items-center gap-2;
}

.sort-button {
  @apply flex items-center justify-center;
  @apply transition-all duration-300;

  :deep(.el-icon) {
    @apply text-sm;
  }
}

.plugin-detail {
  @apply p-4;
  @apply w-full;

  @media screen and (max-width: 768px) {
    @apply p-0;
  }

  .detail-card {
    @apply backdrop-blur-sm;
    @apply rounded-lg;
    @apply shadow-sm;
    @apply transition-all duration-300;
    @apply w-full;
    @apply p-6;

    :deep(.el-card__body) {
      .dark & {
        @apply bg-transparent;
      }
    }

    @media (max-width: 768px) {
      @apply p-0;
      @apply rounded-none;
      @apply shadow-none;

      :deep(.el-card__body) {
        padding: 0 !important;
        margin: 0 !important;
      }
    }
  }

  .detail-header {
    @apply mb-6;

    .header-main {
      @apply space-y-6;
    }

          .plugin-info {
        @apply flex items-start gap-6;
        @apply w-full;

        .plugin-logo {
          @apply w-20 h-20;
          @apply rounded-lg;
          @apply object-contain;
          @apply shadow-sm;
        }

        .info-content {
          @apply flex-1;

          .plugin-title {
            @apply text-2xl font-medium;
            @apply text-gray-900 dark:text-gray-100;
            @apply mb-2;
          }

          .plugin-desc {
            @apply text-gray-600 dark:text-gray-400;
            @apply text-sm;
          }
        }

        .plugin-actions {
          @apply flex flex-col gap-2;
          @apply ml-auto;
          @apply self-center;
        }

        &.plugin-info-mobile {
          @apply flex-col;

          .plugin-logo {
            @apply mx-auto mb-4;
            @apply w-24 h-24;
          }

          .info-content {
            @apply text-center mb-6;
            @apply w-full;

            .plugin-title {
              @apply text-xl;
              @apply mx-auto;
              @apply text-center;
            }

            .plugin-desc {
              @apply text-center mx-auto;
              @apply max-w-md;
            }
          }
        }

        .plugin-actions-mobile {
          @apply flex-row justify-center gap-4;
          @apply mt-2 mb-4;
          @apply w-full;
          @apply self-start;
        }
    }

    .tag-group {
      @apply flex flex-wrap gap-2;
      @apply mt-4;

      &.tag-group-mobile {
        @apply justify-center;
      }
    }

    .stats-group {
      @apply flex flex-wrap gap-6;
      @apply mt-4;
      @apply text-sm text-gray-500 dark:text-gray-400;

      .stat-item {
        @apply flex items-center gap-2;
      }

      &.stats-group-mobile {
        @apply justify-center flex-col items-center gap-3;

        .stat-item {
          @apply text-center;
        }
      }
    }
  }

  .detail-tabs {
    :deep(.el-tabs__content) {
      @apply mt-4;
    }
  }

  .version-list {
    @apply px-4;

    &.version-list-mobile {
      @apply px-1;
    }

    :deep(.el-timeline) {
      @apply space-y-6;
    }

    :deep(.el-timeline-item__node) {
      @apply bg-blue-500;
    }

    :deep(.el-timeline-item__timestamp) {
      @apply text-sm font-medium;
      @apply text-gray-500 dark:text-gray-400;
      @apply mb-2;
    }

    :deep(.version-timeline-item) {
      @apply mb-8;

      .el-timeline-item__wrapper {
        @apply py-2;
      }
    }

    .version-card {
      @apply bg-white/50 dark:bg-transparent;
      @apply backdrop-blur-sm;
      @apply border border-gray-200 dark:border-gray-700;
      @apply transition-all duration-300;
      @apply hover:shadow-md;
      @apply rounded-lg;
      @apply overflow-hidden;

      :deep(.el-card__header) {
        @apply p-0;
      }

      :deep(.el-card__body) {
        .dark & {
          @apply bg-transparent;
        }
      }

      .version-header {
        @apply flex items-center justify-between;
        @apply p-4 pb-3;
        @apply border-b border-gray-200 dark:border-gray-700;
        @apply flex-wrap gap-3;

        .version-tag {
          @apply text-lg font-medium;
          @apply text-gray-900 dark:text-gray-100;
          @apply flex items-center gap-2;
          @apply m-0;
          @apply py-1;

          &::before {
            content: '';
            @apply w-2 h-2;
            @apply rounded-full;
            @apply bg-green-500;
          }
        }

        .version-actions {
          @apply flex items-center gap-3;
          @apply flex-wrap;

          .version-support {
            @apply text-sm;
            @apply px-3 py-1;
            @apply rounded-full;
          }

          .el-button {
            @apply hover:scale-110;
            @apply transition-transform duration-300;
          }
        }
      }

      .changelog {
        @apply p-5;
        @apply text-gray-600 dark:text-gray-400;
        @apply text-sm;
        @apply leading-relaxed;

        :deep(h1, h2, h3, h4, h5, h6) {
          @apply font-medium;
          @apply mb-2;
        }

        :deep(p) {
          @apply mb-3;
        }

        :deep(ul, ol) {
          @apply pl-6;
          @apply mb-3;
        }

        :deep(code) {
          @apply px-1.5 py-0.5;
          @apply rounded;
          @apply bg-gray-100 dark:bg-gray-700;
          @apply text-sm;
        }

        :deep(pre) {
          @apply p-4;
          @apply rounded-lg;
          @apply bg-gray-100 dark:bg-gray-700;
          @apply overflow-x-auto;
          @apply mb-4;
        }
      }
    }
  }
}

.dialog-footer {
  @apply flex justify-end gap-2;
  @apply mt-4;
}

.plugin-detail-drawer {
  :deep(.el-drawer__body) {
    padding: 0 !important;
    margin: 0 !important;
    overflow-x: hidden;
  }

  :deep(.el-card) {
    @media (max-width: 768px) {
      box-shadow: none !important;
      border-radius: 0 !important;

      .el-card__body {
        padding: 0 !important;
        margin: 0 !important;
      }
    }
  }
}

/* 添加全局样式覆盖 */
:deep(.plugin-detail-drawer .el-drawer__body) {
  padding: 0 !important;
  margin: 0 !important;
}

// 优化标签样式
:deep(.el-tag) {
  @apply border-0;
  @apply shadow-sm;

  &.el-tag--warning {
    @apply bg-amber-100 text-amber-800;
    @apply dark:bg-amber-900/30 dark:text-amber-200;
  }

  &.el-tag--success {
    @apply bg-green-100 text-green-800;
    @apply dark:bg-green-900/30 dark:text-green-200;
  }
}

.comments-list {
  @apply mt-4;

  .empty-comments {
    @apply flex justify-center py-8;
  }

  .comments-container {
    @apply space-y-6 mb-6;
  }

      .add-comment-container {
      @apply mt-6 border-t border-gray-200 dark:border-gray-700 pt-4;
      @apply pb-4;

      .comment-actions {
        @apply flex justify-end mt-2;
      }
    }

    .comments-spacer {
      @apply h-20; /* 评论区底部间距 */
  }

  .comment-item {
    @apply p-4;
    @apply bg-white/50 dark:bg-gray-800/50;
    @apply backdrop-blur-sm;
    @apply border border-gray-200 dark:border-gray-700;
    @apply rounded-lg;
    @apply transition-all duration-300;
    @apply hover:shadow-md;
  }

  .comment-header {
    @apply flex items-center justify-between;
    @apply mb-2;

    .comment-user {
      @apply font-medium;
      @apply text-gray-900 dark:text-gray-100;
    }

    .comment-time {
      @apply text-sm;
      @apply text-gray-500 dark:text-gray-400;
    }
  }

  .comment-content {
    @apply py-2;
    @apply text-gray-700 dark:text-gray-300;
  }

  .comment-footer {
    @apply flex justify-between items-center;
    @apply mt-2;

    .comment-actions {
      @apply flex items-center gap-2;
    }

    .like-button, .reply-button {
      @apply flex items-center gap-1;
      @apply text-gray-500 dark:text-gray-400;
      @apply hover:text-blue-500 dark:hover:text-blue-400;
      @apply transition-colors duration-300;
    }
  }

  .reply-input-container {
    @apply mt-3 mb-3 pl-4 border-l-2 border-gray-200 dark:border-gray-700;

    .reply-to-tip {
      @apply text-sm font-medium text-blue-500 mb-2;
    }

    .reply-actions {
      @apply flex justify-end mt-2 gap-2;
    }
  }

  .comment-replies {
    @apply mt-4 ml-8;
    @apply space-y-4;
    @apply border-l-2 border-gray-200 dark:border-gray-700;
    @apply pl-4;
  }

  .emoji-container {
    @apply relative mt-2;

    .emoji-toggle {
      @apply text-blue-500;
    }

          .emoji-picker {
      @apply bg-white dark:bg-gray-800;
      @apply rounded-lg shadow-lg;
      @apply border border-gray-200 dark:border-gray-700;
      @apply flex flex-wrap;
      @apply overflow-y-auto;
      @apply absolute bottom-full left-0;
      @apply p-2;
      @apply mb-2;
      width: 280px;
      max-height: 180px;

      .emoji-item {
        @apply p-2 cursor-pointer text-xl;
        @apply hover:bg-gray-100 dark:hover:bg-gray-700;
        @apply rounded transition-colors;
        @apply flex items-center justify-center;
        width: 40px;
        height: 40px;
      }
    }
  }

  .reply-item {
    @apply p-3;
    @apply bg-gray-50 dark:bg-gray-800;
    @apply rounded-lg;
    @apply border border-gray-100 dark:border-gray-700;
  }

  .reply-header {
    @apply flex items-center justify-between;
    @apply mb-2;

    .reply-user {
      @apply font-medium;
      @apply text-gray-900 dark:text-gray-100;
    }

    .reply-time {
      @apply text-sm;
      @apply text-gray-500 dark:text-gray-400;
    }
  }

  .reply-content {
    @apply py-1;
    @apply text-gray-700 dark:text-gray-300;
  }

  .reply-footer {
    @apply flex justify-between items-center;
    @apply mt-2;

    .reply-actions {
      @apply flex items-center gap-2;
    }
  }
}

.version-mobile-info {
          @apply mt-4 pt-4;
          @apply border-t border-gray-100 dark:border-gray-700;
          @apply flex gap-2 flex-wrap justify-center;

          .el-tag {
            @apply mb-1;
          }
        }

@media screen and (max-width: 768px) {
  .version-header {
    @apply flex items-center justify-between;

    .version-tag {
      @apply m-0;
    }

    .version-actions {
      @apply self-center;
    }
  }

  .changelog {
    @apply px-4 py-4;
  }
}
</style>
