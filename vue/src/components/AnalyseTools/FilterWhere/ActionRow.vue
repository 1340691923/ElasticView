<template>
  <div>
    <div class="ta-filter-condition ">

      <el-select
        v-model="v.columnName"
        filterable
        size="mini"
        style="width: 300px"
        placeholder="请选择"
        @change="changeColumnNameSelect()"
      >
        <el-option-group
          v-for="group in options"
          :key="group.label"
          :label="group.label"
        >
          <el-option
            v-for="item in group.options"
            :key="item.value"
            :label="item.label"
            :value="item.value"
          >
            <span style="float: left">{{ item.label }}</span>
            <span style="float: right; color: #8492a6; font-size: 13px">{{ item.value }}</span>
          </el-option>
        </el-option-group>
      </el-select>
      <el-select v-model="v.comparator" size="mini" style="width: 80px">
        <el-option v-for="(v,k,index) in getDataTypeCalcuSymbol(v.columnName)" :key="index" :label="v" :value="k" />
      </el-select>

      <template v-if="noValueSymbolArr.indexOf(v.comparator)!==-1" />
      <template v-else-if="rangeSymbolArr.indexOf(v.comparator)!==-1">
        <el-input v-model="v.ftv[0]" clearable size="mini" type="number" style="width: 300px" />~
        <el-input v-model="v.ftv[1]" clearable size="mini" type="number" style="width: 300px" />
      </template>

      <el-input v-else-if="inputSymbolArr.indexOf(v.comparator)!==-1" v-model="v.ftv" clearable size="mini" style="width: 300px" />
      <template v-else-if="rangeTimeSymbolArr.indexOf(v.comparator)!==-1">
        <el-date-picker
          v-model="v.ftv"

          size="mini"
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
        />
      </template>

      <template v-else>
        <keep-alive>
          <select-values ref="values" v-model="v.ftv" :table-typ="tableTyp" style="width: 300px" :data="v.columnName" />
        </keep-alive>
      </template>

    </div>
  </div>

</template>

<script>
import { pickerOptions } from '@/utils/date'
import { dataTypeCalcuSymbol, noValueSymbolArr, inputSymbolArr, rangeSymbolArr, rangeTimeSymbolArr } from '@/utils/base-data'

export default {
  name: 'ActionRow',
  components: {
    SelectValues: () => import('@/components/AnalyseTools/FilterWhere/values')
  },
  props: {
    value: {
      type: Object,
      default: {}
    },
    datas: {
      type: Object,
      default: {}
    },
    options: {
      type: Array,
      default: []
    },
    dataTypeMap: {
      type: Array,
      default: []
    },
    tableTyp: {
      type: Number,
      default: 0
    }
  },
  data() {
    return {
      pickerOptions: pickerOptions,
      v: this.value,
      inputSymbolArr: inputSymbolArr,
      noValueSymbolArr: noValueSymbolArr,
      rangeSymbolArr: rangeSymbolArr,
      rangeTimeSymbolArr: rangeTimeSymbolArr
    }
  },
  watch: {
    'v.columnName'(val, oldVal) {
      this.$emit('input', this.v)
    },
    'v.comparator'(val, oldVal) {
      this.$emit('input', this.v)
    },
    'v.ftv'(val, oldVal) {
      this.$emit('input', this.v)
    }
  },
  mounted() {
    this.pickerOptions = pickerOptions
  },
  methods: {
    getDataTypeCalcuSymbol(data) {
      let typ = 0

      if (this.dataTypeMap.hasOwnProperty(this.tableTyp.toString())) {
        for (const v of this.dataTypeMap[this.tableTyp.toString()]) {
          if (data == v.attribute_name) {
            typ = v.data_type
          }
        }
      }

      return dataTypeCalcuSymbol[typ]
    },
    changeColumnNameSelect() {
      const tmp = this.getDataTypeCalcuSymbol(this.v.columnName)

      for (const i in tmp) {
        this.v.comparator = i
        break
      }

      this.$refs['values'].cleanValues()
      this.$emit('input', this.v)
    }

  }
}
</script>

<style scoped>

</style>
