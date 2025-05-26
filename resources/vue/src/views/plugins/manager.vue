<template>
  <div class="app-container">
    <div class="search-container">
      <el-form :inline="true">



        <el-form-item>
         <el-container>
           <el-upload
             ref="uploadRef"
             :headers="uploadHeader"
             v-model:file-list="fileList"

             :action="uploadFileAction()"
             :on-preview="handlePreview"
             :on-remove="handleRemove"
             :before-remove="beforeRemove"
             :limit="1"
             :on-exceed="handleExceed"
             :on-success="uploadSucc"
           >
             <el-button type="primary">离线上传</el-button>

           </el-upload>
           <el-button   style="margin-left: 1rem" @click="getLocalPluginList" type="primary">
             刷新
           </el-button>
         </el-container>

        </el-form-item>
      </el-form>
    </div>
    <el-table
      :data="data.pluginList"
      v-loading="loading"
      show-overflow-tooltip
      class="plugin-table"
    >
      <el-table-column align="center" :label="$t('插件名')" >
        <template #default="scope">
          <el-tag :type="scope.row.is_exited?'danger':'success'">
            {{ scope.row.plugin_name }}
          </el-tag>
        </template>
      </el-table-column>
      <el-table-column align="center" :label="$t('开发者')" width="120">
        <template #default="scope">
          {{ scope.row.developer }}
        </template>
      </el-table-column>
      <el-table-column align="center" :label="$t('别名')" width="200">
        <template #default="scope">
          {{ scope.row.plugin_alias }}
        </template>
      </el-table-column>
      <el-table-column align="center" :label="$t('文件名')" width="250">
        <template #default="scope">
          {{ scope.row.plugin_file_name }}
        </template>
      </el-table-column>
      <el-table-column align="center" :label="$t('版本')" width="170">
        <template #default="scope">
          <template v-if="scope.row.has_update">
            <el-tag type="danger">{{ scope.row.version }}</el-tag>
            <el-icon class="mx-1"><Right /></el-icon>
            <el-tag type="success">{{ scope.row.update_version }}</el-tag>
          </template>
          <template v-else>
            <el-tag type="success">{{ scope.row.version }}</el-tag>
          </template>
        </template>
      </el-table-column>
      <el-table-column align="center" :label="$t('性能')" width="240">
        <template #default="scope">
          <div class="performance-info">
            <el-tooltip effect="dark" placement="top">
              <template #content>
                <div class="tooltip-content">
                  <div class="tooltip-item">
                    <span class="label">进程ID:</span>
                    <span class="value">{{ scope.row.pid }}</span>
                  </div>
                  <div class="tooltip-item">
                    <span class="label">CPU:</span>
                    <span class="value">{{ scope.row.cpu_percent_str }}</span>
                  </div>
                  <div class="tooltip-item">
                    <span class="label">内存:</span>
                    <span class="value">{{ scope.row.memory_percent_str }}</span>
                  </div>
                </div>
              </template>
              <div class="performance-indicators">
                <el-icon><Monitor /></el-icon>
                <span class="ml-1">{{ scope.row.pid }}</span>
                <el-progress
                  :percentage="parseFloat(scope.row.cpu_percent_str)"
                  :stroke-width="4"
                  class="performance-progress"
                  :color="getCpuColor"
                />
                <el-progress
                  :percentage="parseFloat(scope.row.memory_percent_str)"
                  :stroke-width="4"
                  class="performance-progress"
                  :color="getMemoryColor"
                />
              </div>
            </el-tooltip>
          </div>
        </template>
      </el-table-column>
      <el-table-column align="center" fixed="right" :label="$t('操作')" width="150">
        <template #default="scope">
          <div class="action-buttons">
            <el-button v-if="scope.row.has_update" @click="installPlugin(scope.row.plugin_id, scope.row.update_version)" type="warning" >
              更新
            </el-button>
            <el-button @click="showDetails(scope.row)" type="info" >
              详情
            </el-button>
          </div>
        </template>
      </el-table-column>
    </el-table>

    <!-- 详情抽屉 -->
    <el-drawer
      v-model="detailsVisible"
      :title="currentPlugin?.plugin_name"
      direction="rtl"
      :size="isMobile?'100%':'80%'"
    >
      <div class="plugin-details">
        <div class="detail-card">
          <div class="card-title">基本信息</div>
          <div class="card-content">
            <div class="detail-item">
              <span class="detail-label">插件名称</span>
              <span class="detail-value">{{ currentPlugin?.plugin_name }}</span>
            </div>
            <div class="detail-item">
              <span class="detail-label">开发者</span>
              <span class="detail-value">{{ currentPlugin?.developer }}</span>
            </div>
            <div class="detail-item">
              <span class="detail-label">别名</span>
              <span class="detail-value">{{ currentPlugin?.plugin_alias }}</span>
            </div>
            <div class="detail-item">
              <span class="detail-label">文件名</span>
              <span class="detail-value">{{ currentPlugin?.plugin_file_name }}</span>
            </div>
          </div>
        </div>

        <div class="detail-card">
          <div class="card-title">版本信息</div>
          <div class="card-content">
            <div class="detail-item">
              <span class="detail-label">当前版本</span>
              <span class="detail-value">
                <el-tag :type="currentPlugin?.has_update ? 'danger' : 'success'">
                  {{ currentPlugin?.version }}
                </el-tag>
              </span>
            </div>
            <div class="detail-item" v-if="currentPlugin?.has_update">
              <span class="detail-label">可更新版本</span>
              <span class="detail-value">
                <el-tag type="success">{{ currentPlugin?.update_version }}</el-tag>
              </span>
            </div>
          </div>
        </div>

        <div class="detail-card">
          <div class="card-title">性能监控</div>
          <div class="card-content">
            <div class="detail-item">
              <span class="detail-label">进程ID</span>
              <span class="detail-value">{{ currentPlugin?.pid }}</span>
            </div>
            <div class="detail-item">
              <span class="detail-label">CPU占用</span>
              <div class="detail-progress">
                <el-progress
                  :percentage="parseFloat(currentPlugin?.cpu_percent_str)"
                  :color="getCpuColor"
                />
              </div>
            </div>
            <div class="detail-item">
              <span class="detail-label">内存占用</span>
              <div class="detail-progress">
                <el-progress
                  :percentage="parseFloat(currentPlugin?.memory_percent_str)"
                  :color="getMemoryColor"
                />
              </div>
            </div>
          </div>
        </div>

        <div class="detail-card">
          <div class="card-title">调试信息</div>
          <div class="card-content">
            <div class="detail-item">
              <span class="detail-label">调试模式</span>
              <span class="detail-value">
                <el-tag :type="currentPlugin?.backend_debug ? 'warning' : 'info'">
                  {{ currentPlugin?.backend_debug ? '已开启' : '未开启' }}
                </el-tag>
              </span>
            </div>
            <div class="detail-item">
              <span class="detail-label">前端调试端口</span>
              <span class="detail-value">
                {{ currentPlugin?.frontend_debug ? currentPlugin?.frontend_dev_port : '未开启' }}
              </span>
            </div>
          </div>
        </div>

        <div class="detail-card">
          <div class="card-title">路径信息</div>
          <div class="card-content">
            <div class="detail-item">
              <span class="detail-label">存储路径</span>
              <span class="detail-value break-all">{{ currentPlugin?.store_path }}</span>
            </div>
            <div class="detail-item">
              <span class="detail-label">日志位置</span>
              <span class="detail-value break-all">{{ currentPlugin?.log_file_path }}</span>
            </div>
          </div>
        </div>

        <div class="detail-card">
          <div class="card-title">运行状态</div>
          <div class="card-content">
            <div class="detail-item">
              <span class="detail-label">运行状态</span>
              <span class="detail-value">
                <el-tag :type="currentPlugin?.is_exited ? 'danger' : 'success'">
                  {{ currentPlugin?.is_exited ? '已停止' : '运行中' }}
                </el-tag>
              </span>
            </div>
            <div class="detail-item" v-if="!currentPlugin?.is_exited">
              <span class="detail-label">开启时间</span>
              <span class="detail-value">{{ currentPlugin?.start_time }}</span>
            </div>
            <div class="detail-item" v-if="currentPlugin?.is_exited">
              <span class="detail-label">停止时间</span>
              <span class="detail-value">{{ currentPlugin?.end_time }}</span>
            </div>
          </div>
        </div>

        <div class="detail-actions">
          <el-button @click="gotoPProf(currentPlugin?.plugin_id)" type="primary">
            性能分析
          </el-button>
          <el-button @click="unInstall(currentPlugin?.plugin_id)" type="danger">
            卸载插件
          </el-button>
        </div>
      </div>
    </el-drawer>
  </div>
</template>

<script setup lang="ts">

import {GetLocalPluginList, InstallPlugin, UnInstallPlugin} from "@/api/plugins";
import {useAppStore} from "@/store";
import {DeviceEnum} from "@/enums/DeviceEnum";
import {getBaseURL} from "@/utils/request";
import {getToken} from "@/utils/auth";
import {ElMessage, ElMessageBox, UploadProps, UploadUserFile} from "element-plus";
import {ref, computed} from "vue";
import { Monitor, Right } from '@element-plus/icons-vue'

const appStore = useAppStore()

const isMobile = computed(() => appStore.device === DeviceEnum.MOBILE);


const data = reactive({
  pluginList:[]
})

const uninstallLoading = ref(false)

const unInstall = async (pluginId)=>{
  uninstallLoading.value = true
  let res = await UnInstallPlugin({
    plugin_id:pluginId,
  })
  uninstallLoading.value = false
  if (res.code != 0) {
    ElMessage.error({
      type: 'error',offset:60,
      message: res.msg
    })
    return
  }
  getLocalPluginList()

  ElMessage.success({
    type: 'success',offset:60,
    message: res.msg
  })
}

const loading = ref(false)

const getLocalPluginList = async ()=>{
  loading.value = false
  let res =  await GetLocalPluginList({})
  loading.value = false
  if(res.code != 0){
    ElMessage.error(res.msg);
    return
  }

  data.pluginList = res.data

  return
}

const installLoading = ref(false)

const gotoPProf = (pluginId)=>{
  const url = `${getBaseURL()}api/call_plugin/${pluginId}/debug/pprof/`; // 替换为目标链接
  const target = "_blank"; // 在新窗口打开
  window.open(url, target);
}

const installPlugin = async (pluginId,version) => {
  installLoading.value = true
  let res = await InstallPlugin({
    plugin_id:pluginId,
    version:version
  })
  installLoading.value = false
  if (res.code != 0) {
    ElMessage.error({
      type: 'error',offset:60,
      message: res.msg
    })
    return
  }
  getLocalPluginList()
  ElMessage.success({
    type: 'success',offset:60,
    message: res.msg
  })
}

const uploadHeader = {
  "X-Token":getToken(),
  "X-Version":window["appVersion"]
}

const uploadFileAction = ()=>{
  return getBaseURL() + "api/plugins/UploadPlugin"
}

const uploadRef = ref(null);

const uploadSucc = (res)=>{
  if(res.code != 0){
    ElMessage.error(res.msg);
    return
  }
  getLocalPluginList()
  uploadRef.value.clearFiles();
  ElMessage.success(res.msg);
}

const fileList = ref<UploadUserFile[]>([
])

const handleRemove: UploadProps['onRemove'] = (file, uploadFiles) => {
  console.log(file, uploadFiles)
}

const handlePreview: UploadProps['onPreview'] = (uploadFile) => {
  console.log(uploadFile)
}

const handleExceed: UploadProps['onExceed'] = (files, uploadFiles) => {
  ElMessage.warning(
    `The limit is 1, you selected ${files.length} files this time, add up to ${
      files.length + uploadFiles.length
    } totally`
  )
}

const beforeRemove: UploadProps['beforeRemove'] = (uploadFile, uploadFiles) => {
  return ElMessageBox.confirm(
    `Cancel the transfer of ${uploadFile.name} ?`
  ).then(
    () => true,
    () => false
  )
}

const detailsVisible = ref(false)
const currentPlugin = ref(null)

const showDetails = (plugin) => {
  currentPlugin.value = plugin
  detailsVisible.value = true
}

const getCpuColor = (percentage) => {
  if (percentage < 30) return '#67C23A'
  if (percentage < 70) return '#E6A23C'
  return '#F56C6C'
}

const getMemoryColor = (percentage) => {
  if (percentage < 50) return '#67C23A'
  if (percentage < 80) return '#E6A23C'
  return '#F56C6C'
}

onMounted(()=>{
  getLocalPluginList()
})

</script>

<style lang="scss" scoped>
.plugin-table {
  @apply rounded-lg overflow-hidden;
  @apply shadow-sm;

  @apply backdrop-blur-sm;
  @apply transition-all duration-300;

  :deep(.el-table__header) {
    @apply bg-gray-50/90 dark:bg-gray-700/90;

    th {
      @apply bg-transparent;
      @apply text-gray-600 dark:text-gray-300;
      @apply font-medium;
      @apply border-b border-gray-200 dark:border-gray-600;
      @apply transition-colors duration-300;
    }
  }

  :deep(.el-table__body) {
    tr {
      @apply transition-colors duration-300;

      &:hover > td {
        @apply bg-gray-50/80 dark:bg-gray-700/50;
      }

      td {
        @apply border-b border-gray-100 dark:border-gray-700;
        @apply text-gray-600 dark:text-gray-300;
      }
    }
  }

  // 优化标签样式
  :deep(.el-tag) {
    @apply border-0;
    @apply shadow-sm;
    @apply transition-all duration-300;

    &.el-tag--success {
      @apply bg-green-100 text-green-800;
      @apply dark:bg-green-900/30 dark:text-green-200;
    }

    &.el-tag--danger {
      @apply bg-red-100 text-red-800;
      @apply dark:bg-red-900/30 dark:text-red-200;
    }
  }

  // 优化按钮样式
  :deep(.el-button) {
    @apply transition-all duration-300;

    &:not(.is-disabled):hover {
      @apply transform scale-105;
      @apply shadow-md;
    }

    &.el-button--warning {
      @apply bg-amber-500 border-amber-500;
      @apply hover:bg-amber-600 hover:border-amber-600;
    }

    &.el-button--info {
      @apply bg-blue-500 border-blue-500;
      @apply hover:bg-blue-600 hover:border-blue-600;
    }

    &.el-button--danger {
      @apply bg-red-500 border-red-500;
      @apply hover:bg-red-600 hover:border-red-600;
    }
  }

  // 优化加载状态
  :deep(.el-loading-mask) {
    @apply backdrop-blur-sm;
    @apply bg-white/50 dark:bg-gray-800/50;
  }
}

.performance-info {
  @apply inline-flex items-center;

  .performance-indicators {
    @apply flex items-center gap-2;

    .performance-progress {
      @apply w-16;
    }
  }
}

.plugin-details {
  @apply space-y-6 overflow-y-auto;
  height: calc(100vh - 120px);

  .detail-card {
    @apply bg-white/80 dark:bg-gray-800/80;
    @apply rounded-lg shadow-sm;
    @apply backdrop-blur-sm;
    @apply overflow-hidden;
    @apply transition-all duration-300;

    .card-title {
      @apply px-4 py-3;
      @apply bg-gray-50/80 dark:bg-gray-700/80;
      @apply text-gray-700 dark:text-gray-200;
      @apply font-medium text-sm;
      @apply border-b border-gray-100 dark:border-gray-600;
    }

    .card-content {
      @apply p-4 space-y-3;
    }
  }

  .detail-item {
    @apply flex items-center gap-2;

    .detail-label {
      @apply text-gray-500 dark:text-gray-400;
      @apply text-sm;
      @apply w-24 flex-shrink-0;
    }

    .detail-value {
      @apply text-gray-800 dark:text-gray-200;
      @apply flex-1;
      @apply text-sm;
    }

    .detail-progress {
      @apply flex-1;
    }
  }

  .detail-actions {
    @apply sticky bottom-0 left-0 right-0;
    @apply bg-white/80 dark:bg-gray-800;
    @apply backdrop-blur-sm dark:backdrop-blur-none;
    @apply p-4 flex gap-4;
    @apply border-t border-gray-100 dark:border-gray-700;
    @apply shadow-sm;
  }
}

.tooltip-content {
  @apply space-y-1;

  .tooltip-item {
    @apply flex items-center gap-2;

    .label {
      @apply text-gray-300;
    }

    .value {
      @apply text-white font-medium;
    }
  }
}

.break-all {
  @apply break-all;
}
</style>
