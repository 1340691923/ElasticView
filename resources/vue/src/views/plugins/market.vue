<template>
  <div class="app-container">
    <div class="search-container">
      <el-form :inline="true">

        <el-form-item label="插件名/描述:">
          <el-input clearable v-model="input.search_txt" ></el-input>
        </el-form-item>
        <el-form-item label="排序:">
          <el-select style="width:100px" v-model="input.order_by_col">
            <el-option
              label="star次数"
              value="star_cnt"
            />
            <el-option
              label="下载人数"
              value="download_user_cnt"
            />
            <el-option
              label="下载次数"
              value="download_cnt"
            />
            <el-option
              label="最后更新时间"
              value="publish_time"
            />
          </el-select>
          <el-select style="width:100px" v-model="input.order_by_desc">
            <el-option
              label="正序"
              :value="false"
            />
            <el-option
              label="倒序"
              :value="true"
            />
          </el-select>
        </el-form-item>

        <el-form-item label="安装状态:">
          <el-select style="width:100px" v-model="input.has_download_type">
            <el-option
              label="全部"
              :value="null"
            />
            <el-option
              label="未安装"
              :value="false"
            />
            <el-option
              label="已安装"
              :value="true"
            />
          </el-select>
        </el-form-item>

        <el-form-item label="">

            <el-container style="margin-top: 0.7rem">
              <el-button
                type="success"
                class="filter-item"
                @click="getPluginMarket"
              >{{ $t('搜索') }}
              </el-button>

            </el-container>

        </el-form-item>

      </el-form>
    </div>

    <div class="plugin-container" v-loading="pluginListLoading">
      <div class="plugin-list">
        <a
          @click="lookPluginInfo(item)"
          v-for="(item, index) in pluginListData.list"
          :class="pluginBoxClass">
          <img
            :src="item.logo"
            class="plugin-img"
            loading="lazy" height="48px">
          <h2 :style="{'color':getColor}" class="plugin-name"> {{ item.plugin_name }}</h2>
          <div :style="{'color':getColor2}" class="plugin-content">
            <p>
              <el-tooltip
                :content="item.describe"
                placement="top"
              >
                {{
                  truncatedText(item.describe, 50)
                }}
              </el-tooltip>

            </p>
            <div class="css-f3blth-horizontal-group" style="width: 100%; height: auto;">
              <!--                      <div class="css-111p9b5-layoutChildrenWrapper">
                                      <el-tag>Signed</el-tag>
                                    </div>-->
              <div v-if="item.has_download" class="css-111p9b5-layoutChildrenWrapper">
                <el-tag type="success">已安装</el-tag>
              </div>
              <div class="css-111p9b5-layoutChildrenWrapper">
                <el-tag >安装:{{item.download_cnt}}次/{{item.download_user_cnt}}人</el-tag>
              </div>

              <div class="css-111p9b5-layoutChildrenWrapper">
                <el-tag type="warning">star:{{item.star_cnt}}人</el-tag>
              </div>
              <div class="css-111p9b5-layoutChildrenWrapper">
                <el-tag type="warning">最新发布时间:{{item.publish_time}}</el-tag>
              </div>
            </div>
          </div>

          <div class="css-wdk70v">
            <el-tooltip
              v-if="item.star_state == 1"
              content="取消star"
              placement="top"
            >
            <el-button @click.stop="starPlugin(item,index)"
                       type="success" :icon="StarFilled" circle ></el-button>
            </el-tooltip>
            <el-tooltip
              content="star"
              placement="top"
              v-else
            >
            <el-button @click.stop="starPlugin(item,index)"
                       type="warning" :icon="Star" circle ></el-button>
            </el-tooltip>
          </div>

        </a>

      </div>
    </div>

    <div style="margin-top: 16px;">
      <el-pagination
        background
        :current-page="input.page"
        :page-size="input.limit"
        :total="pluginListData.count"
        @current-change="handlePluginListPageChange"
        @size-change="handlePluginListSizeChange"
      />
    </div>

    <el-drawer

      v-model="dialog.visible"
      title="插件详情"
      :size="isMobile?'100%':'80%'"
    >
      <div :style="{background:getBgColor}" class="css-4hn7ji-page-inner">
        <el-card v-loading="installLoading">
          <div class="css-1hgwamg-page-header">
            <div class="css-1yhi3xa">
              <div class="css-bk8b94-title-info-container">
                <div class="css-swtwop"><img class="css-4kz4vr"
                                             :src="publishInput.pluginData.logo"
                >
                  <h1>{{publishInput.pluginData.plugin_name}}</h1></div>
                <div :style="{color:getColor}" class="css-1pdqlwl">

                  <div v-if="publishInput.pluginData.has_download" class="css-8exo7k">
                    <el-tag type="success">已安装</el-tag>
                  </div>

                  <div class="css-8exo7k">
                    <div><el-tag>开发者:{{publishInput.pluginData.realname}}</el-tag></div>
                  </div>


                  <div class="css-8exo7k">
                    <div><el-tag>安装:{{publishInput.pluginData.download_cnt}}次/{{publishInput.pluginData.download_user_cnt}}人</el-tag></div>
                  </div>
                  <div class="css-8exo7k">
                    <div><el-tag>star:{{publishInput.pluginData.star_cnt}}</el-tag></div>
                  </div>

                  <div class="css-8exo7k">
                    <div><el-tag>最新发布时间:{{publishInput.pluginData.publish_time}}</el-tag></div>
                  </div>
                </div>
              </div>
              <div class="css-5ax1kt">
                <div class="css-gjl87o-vertical-group" style="width: 100%; height: 100%;">
                  <div class="css-gxt817-layoutChildrenWrapper">
                    <div class="css-ffyaiw-horizontal-group" style="width: 100%; height: 100%;">
                      <div class="css-18qv8yz-layoutChildrenWrapper"></div>
                    </div>
                  </div>
                </div>
              </div>
            </div>
            <div :style="{color:getColor}" class="css-3oq5wu">
              <div class="css-b1ny57">
                <div>{{publishInput.pluginData.describe}}</div>
              </div>
            </div>
          </div>
          <div class="css-hshm0p">
            <div  class="css-3sx73p">
              <el-tabs v-model="tabShowType"  >
                <el-tab-pane label="介绍" name="readme">
                  <mark-down-view :content="publishInput.pluginData.readme"></mark-down-view>
                </el-tab-pane>
                <el-tab-pane label="版本列表" name="versions">

                  <el-timeline >
                    <el-timeline-item v-for="(item,index) in publishListData.list" :timestamp="item.update_time" placement="top">
                      <el-card>
                        <template #header>
                          <div class="card-header">
                            <span>{{`v${item.version}`}}</span>
                            <div style="float:right">

                              <el-tag v-if="!isMobile" type="warning">最大支持ev版本:{{item.gte_ev_dependency_ver}}</el-tag>
                              &nbsp;
                              <el-tag v-if="!isMobile" type="success">最小支持ev版本:{{item.lte_ev_dependency_ver}}</el-tag>
                              &nbsp;
                              <el-button @click.stop="installPlugin(item.version)"
                                         type="warning" :icon="Download" circle ></el-button>
                            </div>
                          </div>

                        </template>
                        <mark-down-view style="margin-top:5px" :content="item.changelog"></mark-down-view>
                      </el-card>
                    </el-timeline-item>

                  </el-timeline>


                  <el-pagination
                    background
                    :current-page="publishInput.page"
                    :page-size="publishInput.limit"
                    :total="publishListData.count"
                    @current-change="handlePublishPageChange"
                    @size-change="handlePublishSizeChange"
                  />

                </el-tab-pane>
              </el-tabs>
            </div>
          </div>
        </el-card>

      </div>

      <template #footer>
        <div class="dialog-footer">
          <el-button type="primary"  v-if="!publishInput.pluginData.has_download" @click="installPlugin(publishListData.list[0].version)">安装最新版本</el-button>
          <el-button type="danger" v-if="publishInput.pluginData.has_download" @click="unInstall">卸载</el-button>
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

import {Star,StarFilled,Download} from '@element-plus/icons-vue'
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
    download_user_cnt: 0
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
      type: 'error',
      message: res.msg
    })
    return
  }
  publishInput.pluginData.has_download = false
  ElMessage.success({
    type: 'success',
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
      type: 'error',
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

const installPlugin = async (version) => {
  installLoading.value = true

  let res = await InstallPlugin({
    plugin_id:publishInput.pluginData.plugin_alias,
    version:version
  })
  installLoading.value = false

  if (res.msg.indexOf('请先')!==-1){
    openImportEvkeyDialogVisible()
    return
  }

  if (res.code != 0) {
    ElMessage.error({
      type: 'error',
      message: res.msg
    })
    return
  }
  publishInput.pluginData.has_download = true
  ElMessage.success({
    type: 'success',
    message: res.msg
  })
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
      type: 'error',
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
    type: 'success',
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
      type: 'error',
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

import { ref } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'

onMounted(() => {
  getPluginMarket()
})

</script>

<style scoped>

.plugin-container {
  margin-top: 16px;
}

.plugin-list {
  display: grid;
  gap: 24px;
  grid-template-columns: repeat(auto-fill, minmax(272px, 1fr));

}


.plugin-box-black {
  display: grid;
  grid-template-columns: 48px 1fr 24px;
  grid-template-rows: auto;
  gap: 16px;
  grid-auto-flow: row;
  border-radius: 2px;
  padding: 24px;
  background: rgb(34, 37, 43);
  transition: background-color 250ms cubic-bezier(0.4, 0, 0.2, 1) 0ms, box-shadow 250ms cubic-bezier(0.4, 0, 0.2, 1) 0ms, border-color 250ms cubic-bezier(0.4, 0, 0.2, 1) 0ms, color 250ms cubic-bezier(0.4, 0, 0.2, 1) 0ms;
  box-shadow: 0 1px 1px var(--el-box-shadow-light);
}

.plugin-box-light {
  display: grid;
  grid-template-columns: 48px 1fr 24px;
  grid-template-rows: auto;
  gap: 16px;
  grid-auto-flow: row;
  border-radius: 2px;
  padding: 24px;
  background: #ffffff;
  transition: background-color 250ms cubic-bezier(0.4, 0, 0.2, 1) 0ms, box-shadow 250ms cubic-bezier(0.4, 0, 0.2, 1) 0ms, border-color 250ms cubic-bezier(0.4, 0, 0.2, 1) 0ms, color 250ms cubic-bezier(0.4, 0, 0.2, 1) 0ms;

  box-shadow: 0 1px 1px var(--el-box-shadow-light);
}

.css-111p9b5-layoutChildrenWrapper {
  margin-bottom: 8px;
  margin-right: 8px;
  display: flex;
  -webkit-box-align: center;
  align-items: center;
}


.css-111p9b5-layoutChildrenWrapper:last-child {
  margin-right: 0px;
}

.plugin-img {
  grid-area: 1 / 1 / 3 / 2;
  max-width: 100%;
  align-self: center;
  object-fit: contain;
}

.plugin-name {
  grid-area: 1 / 2 / 3 / 3;
  align-self: center;
  font-size: 1.28571rem;
  margin: 0px;
}

.plugin-content {
  grid-area: 3 / 1 / 4 / 3;
}

.css-wdk70v {
  grid-area: 1 / 3 / 2 / 4;
}

.css-f3blth-horizontal-group {
  display: flex;
  flex-flow: wrap;
  -webkit-box-pack: start;
  justify-content: flex-start;
  -webkit-box-align: center;
  align-items: center;
  height: 100%;
  max-width: 100%;
  margin-bottom: -8px;
}

.hover-red:hover {
  color: orangered;
}


.css-1hgwamg-page-header {
  display: flex;
  flex-direction: column;
  gap: 8px;
  margin-bottom: 16px;
}

.css-1yhi3xa {
  align-items: flex-start;
  display: flex;
  flex-flow: wrap;
  gap: 8px 24px;
}

.css-bk8b94-title-info-container {
  display: flex;
  flex: 1 1 0%;
  flex-wrap: wrap;
  gap: 8px 32px;
  -webkit-box-pack: justify;
  justify-content: space-between;
  max-width: 100%;
  min-width: 200px;
}


.css-swtwop {
  display: flex;
  flex-direction: row;
  max-width: 100%;
}

.css-4kz4vr {
  width: 32px;
  height: 32px;
  margin-right: 16px;
}

.css-1pdqlwl {
  display: flex;
  flex-direction: row;
  gap: 12px;
  overflow: auto;
}

.css-8exo7k {
  font-family: Inter, Helvetica, Arial, sans-serif;
  font-weight: 400;
  font-size: 0.857143rem;
  line-height: 1.5;
  letter-spacing: 0.0125em;
  display: flex;
  flex-direction: column;
  gap: 4px;
}

.css-eh43ok {
  border-left: 1px solid rgba(204, 204, 220, 0.12);
}

.css-8exo7k {
  font-family: Inter, Helvetica, Arial, sans-serif;
  font-weight: 400;
  font-size: 0.857143rem;
  line-height: 1.5;
  letter-spacing: 0.0125em;
  display: flex;
  flex-direction: column;
  gap: 4px;
}


.css-1cjwdfr {
  display: flex;
  flex-wrap: wrap;
  gap: 4px;
}

.css-ti9e4i {
  display: inline-flex;
  -webkit-box-align: center;
  align-items: center;
}

.css-3oq5wu {
  position: relative;
}

.css-3sx73p {
  border-bottom: 1px solid rgba(204, 204, 220, 0.12);
  overflow-x: auto;
}


.css-b1ny57 {
  display: flex;
  flex-direction: column;
  gap: 8px;
}
.css-swtwop h1 {
  display: flex;
  margin-bottom: 0px;
}
h1, .h1 {
  margin: 0px 0px 0.45em;
  font-size: 2rem;
  line-height: 1.14286;
  font-weight: 400;
  letter-spacing: -0.00893em;
  font-family: Inter, Helvetica, Arial, sans-serif;
}



</style>
