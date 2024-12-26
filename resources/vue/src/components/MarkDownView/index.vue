<template>
  <MdPreview :theme="theme" :editorId="id" :modelValue="content" />
</template>

<script setup>

import { MdPreview } from 'md-editor-v3';
import 'md-editor-v3/lib/preview.css';
import {useSettingsStore} from "@/store";
import { ThemeEnum } from "@/enums/ThemeEnum";
const id = 'preview-only';
const settingsStore = useSettingsStore();

const props = defineProps({
  content: {
    type: String,
    default: "",
  },
});
const theme = computed(() => {
  return settingsStore.theme === ThemeEnum.DARK
    ? "dark"
    : "light";
});

</script>
