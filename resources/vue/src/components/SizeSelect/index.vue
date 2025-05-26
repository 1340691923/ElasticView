<template>
  <el-dropdown class="size-select" trigger="click" @command="handleSizeChange">
    <div>
      <svg-icon icon-class="size" />
    </div>
    <template #dropdown>
      <el-dropdown-menu>
        <el-dropdown-item
          v-for="item of sizeOptions"
          :key="item.value"
          :command="item.value"
        >
          {{ item.label }}
        </el-dropdown-item>
      </el-dropdown-menu>
    </template>
  </el-dropdown>
</template>

<script setup lang="ts">
import { SizeEnum } from "@/enums/SizeEnum";
import { useAppStore } from "@/store/modules/app";

const { t } = useI18n();
const sizeOptions = computed(() => {
  return [
    { label: t("正常"), value: SizeEnum.DEFAULT },
    { label: t("大"), value: SizeEnum.LARGE },
    { label: t("小"), value: SizeEnum.SMALL },
  ];
});

const appStore = useAppStore();
function handleSizeChange(size: string) {
  appStore.changeSize(size);
  ElMessage.success(t("文字大小已经修改"));
}
</script>
