<template>
  <div class="filter-item">
    时间范围：

    <el-date-picker
      v-if="simple_format == false"
      v-model="date"
      type="daterange"
      align="right"

      unlink-panels
      range-separator="至"
      start-placeholder="开始日期"
      end-placeholder="结束日期"
      value-format="yyyy-MM-dd"

      :picker-options="pickerOptions"
      @change="searchData"
    />
    <el-date-picker
      v-else
      v-model="date"
      type="datetimerange"
      align="right"
      :default-time="['00:00:00', '23:59:59']"
      unlink-panels
      range-separator="至"
      start-placeholder="开始日期"
      end-placeholder="结束日期"
      format="yyyy-MM-dd HH:mm:ss"
      value-format="yyyy-MM-dd HH:mm:ss"
      :picker-options="pickerOptions"
      @change="searchData"
    />
  </div>
</template>

<script>
import { pickerOptions } from '@/utils/date'

import { dateFormat } from '@/utils/date'
const day = 3600 * 1000 * 24
export default {
  name: 'Index',
  props: [
    'dates',
    'simpleFormat'
  ],
  data() {
    return {
      date: this.dates,
      simple_format: this.simple_format,
      pickerOptions: {}
    }
  },
  mounted() {
    this.pickerOptions = pickerOptions
  },
  methods: {
    searchData() {
      if (this.date == null) {
        this.$emit('changeDate', [])
      } else {
        this.$emit('changeDate', this.date)
      }
    }
  }
}
</script>

<style scoped>
  .xwl>>>.el-input__inner{

    color:#000000!important;
    font-family: cursive!important;
    border-top:  1px red !important;
    border-left:  1px red !important;
    border-right:  1px red !important;
  }

</style>
