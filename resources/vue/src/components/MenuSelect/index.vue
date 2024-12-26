<template>
  <div>
    <el-tree-select
      placeholder="定位菜单"
      style="width:10rem"
      v-model="searchValue"
      :data="addFieldToTree(theRoutes)"
      default-expand-all
      filterable
      highlight-current
      :prefix-icon="Search"
      @node-click="handleSelect"
    >
      <template #default="{ data }">
        <!-- 根据 icon 类型决定使用的不同类型的图标组件 -->
        <template v-if="data.hasOwnProperty('meta')">
          <el-icon v-if="data.meta.hasOwnProperty('icon') && data.meta.icon.startsWith('el-icon')" class="sub-el-icon">
            <component :is="data.meta.icon.replace('el-icon-', '')" />
          </el-icon>
          <svg-icon v-else-if="data.meta.hasOwnProperty('icon')" :icon-class="data.meta.icon" />
          <svg-icon v-else icon-class="menu" />
          <span style="margin-left: 3px">{{ translateRouteTitle(data.meta.title) }}</span>
        </template>
      </template>
    </el-tree-select>
  </div>
</template>

<script lang="ts" setup>
import { Search } from '@element-plus/icons-vue'
import {usePermissionStore} from "@/store";
import { translateRouteTitle } from "@/utils/i18n";

import path from 'path-browserify'

defineOptions({
  name: 'MenuSelect',
})

const searchValue = ref<any>('')
const routesStore = usePermissionStore()
const { routes: routes } = storeToRefs(routesStore)

const theRoutes = computed(()=>{
  return generateRoutes(routes.value,'/')
})

const router = useRouter();
const addFieldToTree = (routes: any) => {
  routes.forEach((node: any) => {
    if(node.hasOwnProperty("meta")){
      node.value = node.name
      node.label = translateRouteTitle(node.meta.title)
      if (node.children && node.children.length > 0) addFieldToTree(node.children)
    }
  })
  return routes
}

const handleSelect = (item: any) => {
  router.push(item.path);
}

const filterHidden = (data: any) => {
  return data.reduce((acc: any, item: any) => {
    if (item.meta && item.meta.hidden) return acc
    const newItem = { ...item }
    if (item.children && item.children.length > 0) newItem.children = filterHidden(item.children)
    return [...acc, newItem]
  }, [])
}

const onlyOneShowingChildCall = (children = [], parent) => {
  let onlyOneChild = null
  const showingChildren = children.filter(item => !item.meta.hidden)

  if (showingChildren.length === 1) {
    onlyOneChild = showingChildren[0]
    onlyOneChild.path = path.resolve(parent.path, onlyOneChild.path)
    return onlyOneChild
  }

  return false
}

const generateRoutes = (routes, basePath = '/') => {
  const res = []

  for (let route of routes) {

    if(route.hasOwnProperty("children")) {
      let onlyOneShowingChild = onlyOneShowingChildCall(route.children, route)
      if(onlyOneShowingChild){
        route = onlyOneShowingChild
      }
    }

    if (!route.hasOwnProperty('meta')  ){
      console.log('route',route)
      continue
    }

    if (route.meta.hidden){
      continue
    }

    route.path = path.resolve(basePath, route.path)
    if (route.hasOwnProperty('children')) {
      route.children = generateRoutes(route.children,  route.path)
    }
    res.push(route)
  }
  return res
}

</script>

<style lang="scss" scoped>
.menu-search {
:deep() {
.el-input {
  width: 300px !important;
}
}
}
</style>
