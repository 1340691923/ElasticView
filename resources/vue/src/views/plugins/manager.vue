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
    >
      <el-table-column align="center" :label="$t('插件名')" width="120">
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
      <el-table-column align="center" :label="$t('文件名')" width="200">
        <template #default="scope">
          {{ scope.row.plugin_file_name }}
        </template>
      </el-table-column>
      <el-table-column align="center" :label="$t('版本')" width="150">
        <template #default="scope">
          <template v-if="scope.row.has_update">
            <el-tag type="danger">{{ scope.row.version }}</el-tag>
            <el-icon><Right /></el-icon>
            <el-tag type="success">{{ scope.row.update_version }}</el-tag>
          </template>
          <template v-else>
            <el-tag  type="success">{{ scope.row.version }}</el-tag>
          </template>

        </template>
      </el-table-column>
      <el-table-column align="center" :label="$t('db存储路径')" width="300">
        <template #default="scope">
          {{ scope.row.store_path }}
        </template>
      </el-table-column>

      <el-table-column align="center" :label="$t('插件进程id')" width="100">
        <template #default="scope">
          {{ scope.row.pid }}
        </template>
      </el-table-column>

      <el-table-column align="center" :label="$t('cpu占用')" width="100">
        <template #default="scope">
          {{ scope.row.cpu_percent_str }}
        </template>
      </el-table-column>

      <el-table-column align="center" :label="$t('内存占用')" width="100">
        <template #default="scope">
          {{ scope.row.memory_percent_str }}
        </template>
      </el-table-column>

      <el-table-column align="center" :label="$t('是否调试')" width="100">
        <template #default="scope">
          {{ scope.row.backend_debug }}
        </template>
      </el-table-column>

      <el-table-column align="center" :label="$t('日志位置')" width="300">
        <template #default="scope">
            {{ scope.row.log_file_path }}
        </template>
      </el-table-column>

      <el-table-column align="center" :label="$t('开启时间')" width="200">
        <template #default="scope">
          <template v-if="!scope.row.is_exited">
            {{ scope.row.start_time }}
          </template>
        </template>
      </el-table-column>

      <el-table-column align="center" :label="$t('停止时间')" width="200">
        <template #default="scope">
          <template v-if="scope.row.is_exited">
            {{ scope.row.end_time }}
          </template>
        </template>
      </el-table-column>
      <el-table-column align="center" :label="$t('前端调试端口')" width="100">
        <template #default="scope">
          <template v-if="scope.row.frontend_debug">
            {{ scope.row.frontend_dev_port }}
          </template>
          <template v-else>
            未开启
          </template>
        </template>
      </el-table-column>


      <el-table-column align="center" :label="$t('操作')" fixed="right" :width="isMobile?100:250">
        <template #default="scope">
          <template v-if="!isMobile">
            <el-button v-if="scope.row.has_update" @click="installPlugin(scope.row.plugin_id,
            scope.row.update_version)" v-loading="installLoading" type="warning">
              更新
            </el-button>

            <el-button @click="gotoPProf(scope.row.plugin_id)" type="info">
              pprof
            </el-button>
            <el-button @click="unInstall(scope.row.plugin_id)" type="danger">
              卸载
            </el-button>
          </template>

          <template v-else >
            <el-dropdown trigger="click" >
              <el-button>管理</el-button>
              <template #dropdown>
                <el-dropdown-menu>
                  <el-dropdown-item  command="1">
                    <el-button v-if="scope.row.has_update" @click.stop="installPlugin(item.plugin_id,item.update_version)" v-loading="installLoading" type="warning">
                      更新
                    </el-button>
                  </el-dropdown-item>
                  <el-dropdown-item  command="1">
                    <el-button @click="gotoPProf(scope.row.plugin_id)" type="info">
                      pprof
                    </el-button>
                  </el-dropdown-item>
                  <el-dropdown-item  command="1">
                    <el-button @click="unInstall(scope.row.plugin_id)" type="danger">
                      卸载
                    </el-button>
                  </el-dropdown-item>

                </el-dropdown-menu>
              </template>
            </el-dropdown>
          </template>


        </template>
      </el-table-column>
    </el-table>

  </div>
</template>

<script setup lang="ts">

import {GetLocalPluginList, InstallPlugin, UnInstallPlugin} from "@/api/plugins";
import {useAppStore} from "@/store";
import {DeviceEnum} from "@/enums/DeviceEnum";
import {getBaseURL} from "@/utils/request";
import {getToken} from "@/utils/auth";
import {ElMessage, ElMessageBox, UploadProps, UploadUserFile} from "element-plus";
import {ref} from "vue";

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
      type: 'error',
      message: res.msg
    })
    return
  }
  getLocalPluginList()

  ElMessage.success({
    type: 'success',
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
      type: 'error',
      message: res.msg
    })
    return
  }
  getLocalPluginList()
  ElMessage.success({
    type: 'success',
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

onMounted(()=>{
  getLocalPluginList()
})

</script>

<style scoped>

</style>
