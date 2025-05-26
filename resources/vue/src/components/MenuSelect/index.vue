<template>
  <div class="menu-search-container">
    <el-input
      v-model="searchValue"
      :placeholder="$t('搜索菜单')"
      class="menu-search-input"
      :prefix-icon="Search"
      clearable
    />

    <div class="menu-tree-wrapper">
      <el-tree
        :data="filteredTreeData"
        :expand-on-click-node="false"
        default-expand-all
        highlight-current
        @node-click="handleNodeClick"
      >
        <template #default="{ data }">
          <div class="custom-tree-node">
            <span class="tree-node-icon">
              <el-icon v-if="data.meta?.icon?.startsWith('el-icon')" class="sub-el-icon">
                <component :is="data.meta.icon.replace('el-icon-', '')" />
              </el-icon>
              <svg-icon v-else-if="data.meta?.icon" :icon-class="data.meta.icon" />
              <svg-icon v-else icon-class="menu" />
            </span>
            <span class="tree-node-label">{{ data.label }}</span>
          </div>
        </template>
      </el-tree>
    </div>
  </div>
</template>

<script lang="ts" setup>
import { ref, computed } from 'vue'
import { useRouter } from 'vue-router'
import { Search } from '@element-plus/icons-vue'
import { usePermissionStore } from '@/store'
import { storeToRefs } from 'pinia'
import { translateRouteTitle } from '@/utils/i18n'
import path from 'path-browserify'

defineOptions({ name: 'MenuSelect' })

const searchValue = ref('')
const router = useRouter()
const emit = defineEmits(['select'])

const routesStore = usePermissionStore()
const { routes } = storeToRefs(routesStore)

// 给每个节点添加 label 和 value 字段
const addLabelToTree = (routes: any[]) => {
  routes.forEach((node: any) => {
    node.label = translateRouteTitle(node.meta?.title || '')
    node.value = node.name
    if (node.children && node.children.length > 0) {
      addLabelToTree(node.children)
    }
  })
  return routes
}

// 简化路由结构生成逻辑
const onlyOneShowingChildCall = (children = [], parent) => {
  const showingChildren = children.filter(item => !item.meta.hidden)
  if (showingChildren.length === 1) {
    const onlyOneChild = showingChildren[0]
    onlyOneChild.path = path.resolve(parent.path, onlyOneChild.path)
    return onlyOneChild
  }
  return false
}

const generateRoutes = (routes: any[], basePath = '/') => {
  const res = []
  for (let route of routes) {
    if (route.hasOwnProperty('children')) {
      const onlyOne = onlyOneShowingChildCall(route.children, route)
      if (onlyOne) route = onlyOne
    }
    if (!route.meta || route.meta.hidden) continue
    route.path = path.resolve(basePath, route.path)
    if (route.children) {
      route.children = generateRoutes(route.children, route.path)
    }
    res.push(route)
  }
  return res
}

const fullTree = computed(() => {
  const cloned = JSON.parse(JSON.stringify(routes.value))
  return addLabelToTree(generateRoutes(cloned))
})

// 根据关键词过滤菜单树
function filterTree(nodes: any[], keyword: string): any[] {
  const result = []
  for (const node of nodes) {
    const label = node.label?.toLowerCase() || ''
    const match = label.includes(keyword.toLowerCase())
    if (match) {
      result.push({ ...node, children: node.children ? filterTree(node.children, keyword) : [] })
    } else if (node.children && node.children.length > 0) {
      const filteredChildren = filterTree(node.children, keyword)
      if (filteredChildren.length > 0) {
        result.push({ ...node, children: filteredChildren })
      }
    }
  }
  return result
}

// 动态计算筛选后的树结构
const filteredTreeData = computed(() => {
  if (!searchValue.value) {
    return fullTree.value
  }
  return filterTree(fullTree.value, searchValue.value)
})

// 点击节点逻辑
const handleNodeClick = (data: any) => {
  if (!data.children || data.children.length === 0) {
    router.push(data.path)
    emit('select')
  }
}
</script>

<style lang="scss" scoped>
.menu-search-container {
  @apply flex flex-col gap-4;

  // 搜索输入框样式
  .menu-search-input {
    :deep(.el-input__wrapper) {

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

