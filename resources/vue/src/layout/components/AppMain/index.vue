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

  </section>
</template>

<script setup lang="ts">
import { useSettingsStore, useTagsViewStore } from "@/store";
import variables from "@/styles/variables.module.scss";
import {RegisterMicroApps} from "@/utils/plugin";
import {start} from "qiankun";

const cachedViews = computed(() => useTagsViewStore().cachedViews); // 缓存页面集合
const minHeight = computed(() => {
  if (useSettingsStore().tagsView) {
    return `calc(100vh - ${variables["navbar-height"]} - ${variables["tags-view-height"]})`;
  } else {
    return `calc(100vh - ${variables["navbar-height"]})`;
  }
});

onMounted(()=>{
  if( !window["qiankunStarted"]){
    start({
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
