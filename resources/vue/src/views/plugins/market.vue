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
        <el-form-item label="插件名/描述:">
          <el-input
            v-model="input.search_txt"
            clearable
            placeholder="搜索插件..."
            class="custom-input"
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
              size="default"
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
                size="small"
                @click.stop="starPlugin(item, index)"
              />
          </div>
          </el-card>
          </div>
      </div>
    </div>

    <div class="pagination-container">
      <el-pagination
        v-model:current-page="input.page"
        v-model:page-size="input.limit"
        :total="pluginListData.count"
        :page-sizes="[10, 20, 30, 50]"
        background
        :layout="isMobile ? 'prev, pager, next' : 'total, sizes, prev, pager, next, jumper'"
        @size-change="handlePluginListSizeChange"
        @current-change="handlePluginListPageChange"
      />
    </div>

    <el-drawer
      v-model="dialog.visible"
      title="插件详情"
      :size="isMobile ? '100%' : '80%'"
    >
      <div class="plugin-detail">
        <el-card v-loading="installLoading" class="detail-card">
          <div class="detail-header">
            <div class="header-main">
              <div class="plugin-info">
                <img :src="publishInput.pluginData.logo" class="plugin-logo">
                <div class="info-content">
                  <h1 class="plugin-title">{{publishInput.pluginData.plugin_name}}</h1>
                  <p class="plugin-desc">{{publishInput.pluginData.describe}}</p>
                </div>
                  </div>

              <div class="tag-group">
                <el-tag v-if="publishInput.pluginData.has_download" type="success">已安装</el-tag>
                <el-tag>开发者: {{publishInput.pluginData.realname}}</el-tag>
                <el-tag :type="publishInput.pluginData.buy_coin_num > 0 ? 'primary' : 'info'">
                  {{publishInput.pluginData.buy_coin_num > 0 ? `所需ev币: ${publishInput.pluginData.buy_coin_num}个` : '免费'}}
                </el-tag>
                  </div>

              <div class="stats-group">
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
              <div class="version-list">
                <el-timeline>
                  <el-timeline-item
                    v-for="(item,index) in publishListData.list"
                    :key="index"
                    :timestamp="item.update_time"
                    placement="top"
                  >
                    <el-card class="version-card">
                        <template #header>
                        <div class="version-header">
                          <span class="version-tag">v{{item.version}}</span>
                          <div class="version-actions">
                            <el-tag v-if="!isMobile" type="warning" class="version-support">
                              最大支持ev版本: {{item.gte_ev_dependency_ver}}
                            </el-tag>
                            <el-tag v-if="!isMobile" type="success" class="version-support">
                              最小支持ev版本: {{item.lte_ev_dependency_ver}}
                            </el-tag>
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
              </el-tabs>
        </el-card>
      </div>

      <template #footer>
        <div class="dialog-footer">
          <el-button
            type="primary"
            v-if="!publishInput.pluginData.has_download"
            @click="installPlugin(publishListData.list[0].version)"
          >
            安装最新版本
          </el-button>
          <el-button
            type="danger"
            v-if="publishInput.pluginData.has_download"
            @click="unInstall"
          >
            卸载
          </el-button>
          <el-button @click="handleCloseDialog">取消</el-button>
        </div>
      </template>
    </el-drawer>

    <import-ev-key v-model:visible="importEvkeyDialogVisible" ></import-ev-key>
  </div>
</template>

<script lang="ts" setup>
import {useSettingsStore} from "@/store";
import {ThemeEnum} from "@/enums/ThemeEnum";
import {PluginMarket, GetPluginInfo,InstallPlugin,StarPlugin,UnInstallPlugin,UploadPlugin} from "@/api/plugins";
import MarkDownView from '@/components/MarkDownView/index.vue'
import ImportEvKey from '@/components/ImportEvKey/index.vue'

import {Star,StarFilled,Download, Search, SortUp, SortDown, Clock, ArrowDown} from '@element-plus/icons-vue'
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

import { ref } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'

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
  @apply bg-white/90 dark:bg-gray-800/90;
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

.plugin-card-wrapper {
  @apply transition-all duration-300;
  @apply hover:translate-y-[-4px];
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
    @apply bg-gray-800;
    @apply border-gray-700;
  }
}

.plugin-header {
  @apply p-4;
  @apply flex items-center gap-4;
  @apply border-b border-gray-200 dark:border-gray-700;

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

  .publish-time {
    @apply text-sm text-gray-500 dark:text-gray-400;
  }
}

.pagination-container {
  @apply mt-6;
  @apply flex justify-center;
  @apply overflow-x-hidden;

  :deep(.el-pagination) {
    @apply flex flex-wrap justify-center;
    @apply gap-2;

    @media (max-width: 768px) {
      @apply w-full;
      @apply px-4;

      .el-pagination__sizes {
        @apply hidden;
      }

      .btn-prev,
      .btn-next {
        @apply min-w-[32px];
      }

      .el-pager {
        @apply flex-wrap justify-center;
      }
    }
  }
}

// 添加暗色模式适配
.dark {
  .search-card {
    @apply bg-gray-800/90;
  }

  .plugin-card {
    @apply bg-gray-800;
    @apply border-gray-700;
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

  .detail-card {
    @apply bg-white/90 dark:bg-gray-800/90;
    @apply backdrop-blur-sm;
    @apply rounded-lg;
    @apply shadow-sm;
    @apply transition-all duration-300;
    @apply w-full;

    @apply p-6;

    @media (max-width: 768px) {
      @apply p-4;
    }
  }

  .detail-header {
    @apply mb-6;

    .header-main {
      @apply space-y-6;
    }

    .plugin-info {
      @apply flex items-start gap-6;

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
    }

    .tag-group {
      @apply flex flex-wrap gap-2;
      @apply mt-4;
    }

    .stats-group {
      @apply flex flex-wrap gap-6;
      @apply mt-4;
      @apply text-sm text-gray-500 dark:text-gray-400;

      .stat-item {
        @apply flex items-center gap-2;
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

    .version-card {
      @apply bg-white/50 dark:bg-gray-800/50;
      @apply backdrop-blur-sm;
      @apply border border-gray-200 dark:border-gray-700;
      @apply transition-all duration-300;
      @apply hover:shadow-md;

      .version-header {
        @apply flex items-center justify-between;
        @apply p-4 pb-3;
        @apply border-b border-gray-200 dark:border-gray-700;

        .version-tag {
          @apply text-lg font-medium;
          @apply text-gray-900 dark:text-gray-100;
          @apply flex items-center gap-2;

          &::before {
            content: '';
            @apply w-2 h-2;
            @apply rounded-full;
            @apply bg-green-500;
          }
        }

        .version-actions {
          @apply flex items-center gap-3;

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
        @apply p-4;
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

:deep(.el-drawer__body) {
  @apply p-0;
  @apply overflow-x-hidden;
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
</style>
