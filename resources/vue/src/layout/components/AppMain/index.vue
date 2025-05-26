<template>
  <section class="app-main" :style="{ minHeight: minHeight }">
    <router-view>
      <template #default="{ Component, route }">
        <transition
          enter-active-class="animate__animated animate__fadeIn"
          mode="out-in"
        >
          <keep-alive :include="cachedViews">
            <component :is="Component" :key="route.path" />
          </keep-alive>
        </transition>
      </template>
    </router-view>
    <section id="Appmicro" />
    <float-btn @handlepaly="clickFloatBtn"></float-btn>

    <el-drawer
      v-model="showBigModel"
      direction="rtl"
      custom-class="drawer-style"
      :title="$t('百炼大模型')"
      :size="isMobile?'100%':'80%'"
    >
      <!-- 搜索框 -->

      <!-- 结果显示框 -->
      <el-scrollbar  v-loading="bigModelLoading" style="max-height: 56vh">
        <el-card v-if="bigModelResponse != ''">
              <mark-down-view :content="bigModelResponse "></mark-down-view>
        </el-card>
      </el-scrollbar>
      <el-divider />
      <el-input
        type="textarea"
        :autosize="{ minRows: 2, maxRows: 4}"
        v-model="bigModelInput"
        placeholder="AI大模型搜索，请输入问题"
        clearable
      >

      </el-input>
      <el-divider />
      <el-space justify="end">
        <el-button type="primary" v-loading="bigModelLoading" @click="sendBigModelQuery">
          发送
        </el-button>
        <el-button @click="clearBigModelResponse">
          清空
        </el-button>

      </el-space>
    </el-drawer>

  </section>
</template>

<script setup lang="ts">
import {useAppStore, useSettingsStore, useTagsViewStore} from "@/store";
import variables from "@/styles/variables.module.scss";
import {start} from "qiankun";
import FloatBtn from '@/components/FloatBtn.vue'
import {SearchBigMode} from "@/api/ai";
import {DeviceEnum} from "@/enums/DeviceEnum";
import MarkDownView from '@/components/MarkDownView/index.vue'

const cachedViews = computed(() => useTagsViewStore().cachedViews); // 缓存页面集合
const minHeight = computed(() => {
  if (useSettingsStore().tagsView) {
    return `calc(100vh - ${variables["navbar-height"]} - ${variables["tags-view-height"]})`;
  } else {
    return `calc(100vh - ${variables["navbar-height"]})`;
  }
});

const showBigModel = ref(false)

const bigModelInput = ref('')

const bigModelLoading = ref(false)

const bigModelResponse = ref('')

const clickFloatBtn = ()=>{
  showBigModel.value = !showBigModel.value
}

const clearBigModelResponse =()=>{
  bigModelInput.value = ''
  bigModelResponse.value = ''
}

const sendBigModelQuery = async ()=>{
  bigModelLoading.value = true
  let res =  await  SearchBigMode({content:bigModelInput.value})
  if(res.code != 0){
    bigModelLoading.value = false
    ElMessage.error(res.msg);
    return
  }
  bigModelLoading.value = false
  bigModelResponse.value = res.data

  return
}

const appStore = useAppStore()

const isMobile = computed(() => appStore.device === DeviceEnum.MOBILE);



onMounted(()=>{
  if( !window["qiankunStarted"]){
    start({
      prefetch:true,
      singular:true,
      sandbox: {
      // 开启严格的样式隔离模式。这种模式下 qiankun 会为每个微应用的容器包裹上一个 [shadow dom]节点，从而确保微应用的样式不会对全局造成影响。
      //strictStyleIsolation: false,
      // 设置实验性的样式隔离特性，即在子应用下面的样式都会包一个特殊的选择器规则来限定其影响范围
      //experimentalStyleIsolation: true
      }
    })
    window["qiankunStarted"] = true
  }
})

</script>

<style lang="scss" scoped>
.app-main {
  position: relative;
  background-color: var(--el-bg-color-page);
}
</style>
