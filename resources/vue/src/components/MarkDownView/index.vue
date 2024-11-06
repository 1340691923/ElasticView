<template>
  <MdPreview :editorId="id" :modelValue="content" />

</template>

<script setup>

import { MdPreview } from 'md-editor-v3';
import 'md-editor-v3/lib/preview.css';

const id = 'preview-only';

const props = defineProps({
  content: {
    type: String,
    default: "",
  },
});


</script>
