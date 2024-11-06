<template>
  <el-form>
    <el-table
      :data="localHeaders"
      style="width: 100%;"
      border
    >
      <el-table-column
        prop="key"
        :label="$t('键')"
      >
        <template #default="{ row, $index }">
          <el-input
            v-model="row.key"
            :placeholder="$t('键')"
            @input="updateHeaders"
          />
        </template>
      </el-table-column>

      <el-table-column
        prop="value"

        :label="$t('值')"
      >
        <template #default="{ row, $index }">
          <el-input
            v-model="row.value"
            :placeholder="$t('值')"
            @input="updateHeaders"
          />
        </template>
      </el-table-column>

      <el-table-column
        :label="$t('操作')"
        width="120"
      >
        <template #header>
          <el-button
            type="primary"
            @click="addHeader"
          >
            {{ $t('新增键值') }}
          </el-button>
        </template>
        <template #default="{ $index }">
          <el-button
            type="danger"
            @click="removeHeader($index)"
          > {{$t('删除')}}</el-button>
        </template>
      </el-table-column>
    </el-table>
  </el-form>
</template>

<script setup>
import { ref, watch } from 'vue'

const props = defineProps({
  modelValue: {
    type: Array,
    default: () => []
  }
})

const emit = defineEmits(['update:modelValue'])

const localHeaders = ref([...props.modelValue])

// Watch for changes in modelValue and update localHeaders
watch(
  () => props.modelValue,
  (newValue) => {
    localHeaders.value = [...newValue]
  },
  { deep: true }
)

const updateHeaders = () => {
  emit('update:modelValue', localHeaders.value)
}

const addHeader = () => {
  localHeaders.value.push({ key: '', value: '' })
  updateHeaders()
}

const removeHeader = (index) => {
  localHeaders.value.splice(index, 1)
  updateHeaders()
}
</script>

<style scoped>
.el-input {
  width: 100%;
}
</style>
