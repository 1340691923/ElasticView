<template>
  <el-date-picker
    type="datetimerange"
    range-separator="To"
    :clearable="false"
    :editable="false"
    :model-value="modelValue"
    popper-class="vab-date-picker"
    :shortcuts="shortcuts"
    @calendar-change="handleChange" />
</template>

<script setup lang="ts">
defineProps<{
  modelValue: any
}>()

import dayjs from "dayjs";

const emits = defineEmits(['update:modelValue'])

const handleChange = (data: any) => {
  data[0] && data[1] && emits('update:modelValue', data)
}

const shortcuts =   [
  {
    text: '今天',
    value: () => {

      let date = [dayjs().format('YYYY-MM-DD 00:00:00'), dayjs().format('YYYY-MM-DD 23:59:59')]
      handleChange(date)
      return date
    },
  },
  {
    text: '昨天',
    value: () => {
      let date = [dayjs().subtract(1, 'day').format('YYYY-MM-DD 00:00:00'), dayjs().subtract(1, 'day').format('YYYY-MM-DD 23:59:59')]

      handleChange(date)
      return date
    },
  },
  {
    text: '一周',
    value: () => {
      let date = [dayjs().subtract(7, 'day').format('YYYY-MM-DD 00:00:00'), dayjs().format('YYYY-MM-DD 23:59:59')]

      handleChange(date)
      return date
    },
  },
  {
    text: '两周',
    value: () => {
      let date = [dayjs().subtract(14, 'day').format('YYYY-MM-DD 00:00:00'), dayjs().format('YYYY-MM-DD 23:59:59')]

      handleChange(date)
      return date
    },
  },
  {
    text: '一月',
    value: () => {
      let date = [dayjs().subtract(30, 'day').format('YYYY-MM-DD 00:00:00'), dayjs().format('YYYY-MM-DD 23:59:59')]

      handleChange(date)
      return date
    },
  },
  {
    text: '三个月',
    value: () => {
      let date = [dayjs().subtract(90, 'day').format('YYYY-MM-DD 00:00:00'), dayjs().format('YYYY-MM-DD 23:59:59')]

      handleChange(date)
      return date
    }
  },
  {
    text: '半年',
    value: () => {
      let date = [dayjs().subtract(180, 'day').format('YYYY-MM-DD 00:00:00'), dayjs().format('YYYY-MM-DD 23:59:59')]

      handleChange(date)
      return date
    },
  },
  {
    text: '一年',
    value: () => {
      let date = [dayjs().subtract(360, 'day').format('YYYY-MM-DD 00:00:00'), dayjs().format('YYYY-MM-DD 23:59:59')]

      handleChange(date)
      return date
    },

  },
]
</script>

<style lang="scss">
@media (max-width: 768px) {
$panel-width: 70px;

  .vab-date-picker {
    width: 77%;

  .el-date-range-picker {
    width: 100% !important;

  .el-picker-panel__body {
    min-width: 76%;
    margin-left: $panel-width;
  }
}

.el-picker-panel__shortcut {
  text-align: center;
}

.el-date-range-picker__content {
  width: 100% !important;
  padding: 5px;
  margin: 0px;

.el-date-table {
th {
  padding: 0px;
}

td {
  padding: 0px;
}
}
}

.is-left {
  padding-bottom: 0px;
  border-bottom: 1px solid var(--el-datepicker-inner-border-color);
}

.is-right {
  padding-top: 0px;
}

.el-picker-panel__sidebar {
  width: $panel-width;
}

td {
  height: 24px;

.el-date-table-cell {
  height: 24px;

.el-date-table-cell__text {
  width: 18px;
  height: 18px;
  line-height: 18px;
}
}
}

.el-date-range-picker__header div {
  font-size: 12px;
}
}
}
</style>
