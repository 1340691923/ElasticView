<template>
  <el-dialog title="请填写 evKey" v-model="dialogVisible" @close="handleClose">
    <el-form :model="form" ref="formRef">
      <el-form-item label="evKey" prop="evKey">
        <el-input v-model="form.evKey" placeholder="请输入 evKey"></el-input>
        <div class="hint">请前往 http://dev.elastic-view.cn 注册账户后拷贝evKey至此并保存</div>
      </el-form-item>
    </el-form>
    <span slot="footer" class="dialog-footer">
      <el-button @click="handleClose">取消</el-button>
      <el-button type="primary" @click="submitForm">保存</el-button>
    </span>
  </el-dialog>
</template>

<script setup>
import { ref, watch } from 'vue';
import { ElMessage } from 'element-plus'; // 引入 ElMessage 进行提示
import {ImportEvKey, UnInstallPlugin} from "@/api/plugins";

const props = defineProps({
  visible: {
    type: Boolean,
    default: false
  },

});

const emit = defineEmits(['update:visible']);

const form = ref({
  evKey: ''
});
const dialogVisible = ref(false);

// 监听 props 的变化
watch(() => props.visible, (newValue) => {
  dialogVisible.value = newValue;
});

const handleClose = () => {
  dialogVisible.value = false;
  emit('update:visible', false); // 更新父组件状态
  form.value.evKey = ''; // 重置表单
};

const submitForm = async () => {
  const evKey = form.value.evKey;
  if (!evKey) {
    ElMessage.error('请填写 evKey');
    return;
  }
  let res = await ImportEvKey({
    ev_key:evKey,
  })

  if (res.code != 0) {
    ElMessage.error({
      type: 'error',
      message: res.msg
    })
    return
  }

  ElMessage.success({
    type: 'success',
    message: 'evKey导入成功，请重新进行想要进行的操作'
  })

  handleClose()
};
</script>

<style>
.hint {
  color: #888;
  font-size: 12px;
  margin-top: 5px;
}
</style>
