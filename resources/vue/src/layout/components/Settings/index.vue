<template>
  <el-drawer
    v-model="settingsVisible"
    :size="isMobile?'100%':'30%'"
    :title="$t('ev后台配置')"
  >
    <el-divider>{{$t('数据源')}}</el-divider>
    <div class="flex-center">
      <SelectLink></SelectLink>
    </div>

    <el-divider>{{ $t("主题") }}</el-divider>

    <div class="flex-center">
      <el-tooltip
        :content="$t('主题设置')"
        effect="dark"
        placement="bottom"
      >
      <el-switch
        v-model="isDark"
        active-icon="Moon"
        inactive-icon="Sunny"
        @change="changeTheme"
      />
      </el-tooltip>
    </div> 

    <el-divider>{{$t('界面工具')}}</el-divider>
    <div class="flex-center" >

      <el-tooltip
        :content="$t('布局大小')"
        effect="dark"
        placement="bottom"
      >
        <size-select />
      </el-tooltip>

      <el-tooltip

        :content="$t('切换全屏')"
        effect="dark"
        placement="bottom"
      >
        <!--全屏 -->
        <div style="margin-left: 3rem" @click="toggle">
          <svg-icon
            :icon-class="isFullscreen ? 'fullscreen-exit' : 'fullscreen'"
          />
        </div>
      </el-tooltip>
    </div>

    <el-divider>{{ $t("界面设置") }}</el-divider>

    <div class="setting-item">
      <span class="text-xs">{{ $t("主题颜色") }}</span>
      <ThemeColorPicker
        v-model="settingsStore.themeColor"
        @update:model-value="changeThemeColor"
      />
    </div>

    <div class="setting-item">
      <span class="text-xs">{{ $t("开启 Tags-View") }}</span>
      <el-switch v-model="settingsStore.tagsView" />
    </div>

    <div class="setting-item">
      <span class="text-xs">{{ $t("固定 Header") }}</span>
      <el-switch v-model="settingsStore.fixedHeader" />
    </div>

    <div class="setting-item">
      <span class="text-xs">{{ $t("侧边栏 Logo") }}</span>
      <el-switch v-model="settingsStore.sidebarLogo" />
    </div>

<!--    <div class="setting-item">
      <span class="text-xs">{{ $t("settings.watermark") }}</span>
      <el-switch v-model="settingsStore.watermarkEnabled" /> 
    </div>-->

    <el-divider>{{ $t("导航设置") }}</el-divider>

    <LayoutSelect
      v-model="settingsStore.layout"
      @update:model-value="changeLayout"
    />
  </el-drawer>
</template>

<script setup lang="ts">
import { useSettingsStore, usePermissionStore, useAppStore } from "@/store";
import { LayoutEnum } from "@/enums/LayoutEnum";
import { ThemeEnum } from "@/enums/ThemeEnum";
import {DeviceEnum} from "@/enums/DeviceEnum";

const appStore = useAppStore()

const isMobile = computed(() => appStore.device === DeviceEnum.MOBILE);


const route = useRoute();

const settingsStore = useSettingsStore();
const permissionStore = usePermissionStore();

const { isFullscreen, toggle } = useFullscreen();


const settingsVisible = computed({
  get() {
    return settingsStore.settingsVisible;
  },
  set() {
    settingsStore.settingsVisible = false;
  },
});

/** 切换主题颜色 */
function changeThemeColor(color: string) {
  settingsStore.changeThemeColor(color);
}

/** 切换主题 */
const isDark = ref<boolean>(settingsStore.theme === ThemeEnum.DARK);
const changeTheme = (val: any) => {
  isDark.value = val;
  settingsStore.changeTheme(isDark.value ? ThemeEnum.DARK : ThemeEnum.LIGHT);
};

/** 切换布局 */
function changeLayout(layout: string) {
  settingsStore.changeLayout(layout);
  if (layout === LayoutEnum.MIX) {
    route.name && againActiveTop(route.name as string);
  }
  window.location.reload()
}

/** 重新激活顶部菜单 */
function againActiveTop(newVal: string) {
  const parent = findOutermostParent(permissionStore.routes, newVal);
  if (appStore.activeTopMenu !== parent.path) {
    appStore.activeTopMenu(parent.path);
  }
}

/** 递归查找最外层父节点 */
function findOutermostParent(tree: any[], findName: string) {
  let parentMap: any = {};

  function buildParentMap(node: any, parent: any) {
    parentMap[node.name] = parent;

    if (node.children) {
      for (let i = 0; i < node.children.length; i++) {
        buildParentMap(node.children[i], node);
      }
    }
  }

  for (let i = 0; i < tree.length; i++) {
    buildParentMap(tree[i], null);
  }

  let currentNode = parentMap[findName];
  while (currentNode) {
    if (!parentMap[currentNode.name]) {
      return currentNode;
    }
    currentNode = parentMap[currentNode.name];
  }

  return null;
}
</script>

<style lang="scss" scoped>
// 设置项通用样式
.setting-item {
  @apply py-3 flex items-center justify-between;
  padding: 0.75rem 1rem;
  border-radius: 0.5rem;
  transition: all 0.3s ease;
  margin: 0.5rem 0;
  
  &:hover {
    background-color: rgba(0, 0, 0, 0.02);
  }

  span {
    font-size: 0.875rem;
    font-weight: 500;
  }
}

// 分隔线样式
:deep(.el-divider) {
  margin: 1.5rem 0 1rem;
  
  .el-divider__text {
    font-size: 0.875rem;
    font-weight: 500;
    background-color: var(--el-bg-color);
    padding: 0 1rem;
  }
}

// 居中布局容器
.flex-center {
  display: flex;
  align-items: center;
  justify-content: center;
  gap: 2rem;
  padding: 1rem;
  
  .el-tooltip {
    cursor: pointer;
    transition: all 0.3s ease;
    
    &:hover {
      transform: translateY(-2px);
    }
  }
}

// 主题切换开关样式
:deep(.el-switch) {
  --el-switch-on-color: var(--el-color-primary);
  
  &.is-checked .el-switch__core {
    box-shadow: 0 0 0.5rem var(--el-color-primary-light-5);
  }
}

// 抽屉样式
:deep(.el-drawer) {
  .el-drawer__header {
    margin-bottom: 0;
    padding: 1rem 1.5rem;
    border-bottom: 1px solid var(--el-border-color-light);
    
    .el-drawer__title {
      font-size: 1rem;
      font-weight: 600;
    }
  }
  
  .el-drawer__body {
    padding: 1rem 1.5rem;
  }
}

// SVG 图标样式
.svg-icon {
  font-size: 1.25rem;
  cursor: pointer;
  transition: all 0.3s ease;
  padding: 0.5rem;
  border-radius: 0.375rem;
  
  &:hover {
    background-color: rgba(0, 0, 0, 0.05);
    transform: translateY(-2px);
  }
}

// 深色模式适配
:deep(.dark) {
  .setting-item {
    &:hover {
      background-color: rgba(255, 255, 255, 0.05);
    }
  }
  
  .el-divider__text {
    background-color: var(--el-bg-color);
  }
  
  .svg-icon:hover {
    background-color: rgba(255, 255, 255, 0.1);
  }
}

// 主题色选择器容器
:deep(.theme-picker) {
  display: flex;
  gap: 0.5rem;
  
  .theme-item {
    width: 1.5rem;
    height: 1.5rem;
    border-radius: 0.375rem;
    cursor: pointer;
    transition: all 0.3s ease;
    
    &:hover {
      transform: scale(1.1);
    }
    
    &.active {
      box-shadow: 0 0 0 2px var(--el-color-white),
                 0 0 0 4px var(--el-color-primary);
    }
  }
}

// 布局选择器样式
:deep(.layout-select) {
  display: grid;
  grid-template-columns: repeat(3, 1fr);
  gap: 1rem;
  padding: 1rem 0;
  
  .layout-item {
    border: 2px solid var(--el-border-color);
    border-radius: 0.5rem;
    padding: 0.5rem;
    cursor: pointer;
    transition: all 0.3s ease;
    
    &:hover {
      transform: translateY(-2px);
    }
    
    &.active {
      border-color: var(--el-color-primary);
      box-shadow: 0 0 0.5rem var(--el-color-primary-light-5);
    }
    
    img {
      width: 100%;
      border-radius: 0.25rem;
    }
    
    .layout-title {
      text-align: center;
      margin-top: 0.5rem;
      font-size: 0.875rem;
    }
  }
}
</style>
<style lang="scss" scoped>
.nav-action-item {
  display: inline-block;
  min-width: 40px;
  height: $navbar-height;
  line-height: $navbar-height;
  color: var(--el-text-color);
  text-align: center;
  cursor: pointer;

&:hover {
   background: rgb(0 0 0 / 10%);
 }
}

:deep(.message .el-badge__content.is-fixed.is-dot) {
  top: 5px;
  right: 10px;
}

:deep(.el-divider--horizontal) {
  margin: 10px 0;
}

.dark .nav-action-item:hover {
  background: rgb(255 255 255 / 20%);
}

.layout-top .nav-action-item,
.layout-mix .nav-action-item {
  color: #fff;
}
</style>
