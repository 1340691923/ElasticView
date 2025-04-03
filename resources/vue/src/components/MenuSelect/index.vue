<template>
  <div class="menu-search-container">
    <!-- 搜索输入框 -->
    <el-input
      v-model="searchValue"
      :placeholder="$t('搜索菜单')"
      class="menu-search-input"
      :prefix-icon="Search"
      clearable
      @clear="handleClear"
    />
    
    <!-- 树形菜单 -->
    <div class="menu-tree-wrapper">
      <el-tree
        ref="treeRef"
        :data="addFieldToTree(theRoutes)"  
        :filter-node-method="filterNode"
        :expand-on-click-node="false"
        default-expand-all
        highlight-current
        @node-click="handleNodeClick"
      >
        <template #default="{ data }">
          <div class="custom-tree-node">
            <!-- 图标 -->
            <span class="tree-node-icon">
              <el-icon v-if="data.meta?.icon?.startsWith('el-icon')" class="sub-el-icon">
                <component :is="data.meta.icon.replace('el-icon-', '')" />
              </el-icon>
              <svg-icon v-else-if="data.meta?.icon" :icon-class="data.meta.icon" />
              <svg-icon v-else icon-class="menu" />
            </span>
            <!-- 标题 --> 
            <span class="tree-node-label">{{ translateRouteTitle(data.meta?.title) }}</span>
          </div>
        </template>
      </el-tree>
    </div>
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

// 定义emit
const emit = defineEmits(['select'])

const handleNodeClick = (data: any, node: any) => {
  if (!data.children || data.children.length === 0) {
    // 只有叶子节点才进行路由跳转
    router.push(data.path);
    emit('select');
  } else {
    // 非叶子节点切换展开/收起状态
    node.expanded = !node.expanded;
  }
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
.menu-search-container {
  @apply flex flex-col gap-4;

  // 搜索输入框样式
  .menu-search-input {
    :deep(.el-input__wrapper) {
      @apply bg-white/90 dark:bg-gray-800/90;
      @apply border border-gray-200 dark:border-gray-700;
      @apply rounded-lg;
      @apply shadow-sm;
      @apply transition-all duration-300;
      
      &:hover {
        @apply border-blue-300 dark:border-blue-600;
        @apply shadow;
      }
      
      &.is-focus {
        @apply border-blue-400 dark:border-blue-500;
        @apply shadow-md;
        @apply bg-white dark:bg-gray-800;
      }

      .el-input__prefix-icon {
        @apply text-gray-400 dark:text-gray-500;
        @apply transition-colors;
      }
    }
  }

  // 树形菜单容器
  .menu-tree-wrapper {
    @apply h-[400px] overflow-y-auto;
    @apply bg-white/80 dark:bg-gray-800/80;
    @apply rounded-lg;
    @apply border border-gray-200 dark:border-gray-700;
    @apply shadow-sm;

    :deep(.el-tree) {
      @apply bg-transparent p-3;
      
      // 树节点样式
      .el-tree-node__content {
        @apply h-10 rounded-lg mb-1 px-3;
        @apply transition-all duration-200;
        @apply hover:bg-blue-50 dark:hover:bg-blue-900/20;
        @apply hover:translate-x-1;
        
        &:hover .tree-node-label {
          @apply text-blue-500 dark:text-blue-400;
        }
      }

      // 选中状态
      .el-tree-node.is-current > .el-tree-node__content {
        @apply bg-blue-100 dark:bg-blue-900/30;
        @apply translate-x-2;
        
        .tree-node-label {
          @apply text-blue-600 dark:text-blue-400;
          @apply font-medium;
        }
        
        .tree-node-icon {
          @apply text-blue-500 dark:text-blue-400;
        }
      }
    }
  }

  // 自定义树节点
  .custom-tree-node {
    @apply flex items-center gap-3;
    @apply w-full;
    
    .tree-node-icon {
      @apply flex items-center justify-center;
      @apply w-5 h-5;
      @apply text-gray-400 dark:text-gray-500;
      @apply transition-colors;
    }

    .tree-node-label {
      @apply text-gray-700 dark:text-gray-200;
      @apply transition-all duration-200;
    }
  }
}

// 美化滚动条
.menu-tree-wrapper {
  &::-webkit-scrollbar {
    @apply w-1.5;
  }

  &::-webkit-scrollbar-track {
    @apply bg-transparent;
  }

  &::-webkit-scrollbar-thumb {
    @apply bg-gray-300 dark:bg-gray-600;
    @apply rounded-full;
    @apply hover:bg-gray-400 dark:hover:bg-gray-500;
    @apply transition-colors;
  }
}
</style>
