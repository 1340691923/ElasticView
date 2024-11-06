<template>
  <el-drawer
    v-model="settingsVisible"
    size="300"
    :title="$t('ev后台配置')"
  >
    <el-divider>{{$t('ES连接')}}</el-divider>
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
        <size-select class="nav-action-item" />
      </el-tooltip>

      <el-tooltip
        :content="$t('切换全屏')"
        effect="dark"
        placement="bottom"
      >
        <!--全屏 -->
        <div class="nav-action-item" @click="toggle">
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


const route = useRoute();
const appStore = useAppStore();
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
.setting-item {
  @apply py-1 flex-x-between;
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
