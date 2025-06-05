<template>

    <!-- 根据 icon 类型决定使用的不同类型的图标组件 -->
    <el-icon v-if="icon && icon.startsWith('el-icon')" class="sub-el-icon">
      <component :is="icon.replace('el-icon-', '')" />
    </el-icon>
    <img  v-else-if="icon.indexOf('call_plugin')!==-1" class="img-left w-4 h-4 hover:rotate-180 transition-all duration-500"  :src="getIconUrl(icon)" ></img>
    <svg-icon v-else-if="icon" :icon-class="icon" />
    <svg-icon v-else icon-class="menu" />

    <span v-if="titles" class="ml-1">{{ translateRouteTitle(titles) }}</span>

</template>

<script setup lang="ts">
import { translateRouteTitle } from "@/utils/i18n";

defineProps({
  icon: {
    type: String,
    default: "",
  },
  titles: {
    type: String,
    default: "",
  },
});

const getIconUrl = (path)=>{
  if(!import.meta.env.PROD){
    return import.meta.env.VITE_APP_API_URL+path
  }
  return path
}
</script>

<style lang="scss" scoped>
.sub-el-icon {
  width: 14px !important;
  margin-right: 0 !important;
  color: currentcolor;
}

.hideSidebar {
  .el-sub-menu,
  .el-menu-item {
    .svg-icon,
    .sub-el-icon {
      margin-left: 20px;
    }
  }


  .img-left {
    margin-left: 20px;
  }
}
</style>
