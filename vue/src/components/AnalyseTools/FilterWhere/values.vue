<template>

  <small-select v-model="selectVal" :options="options" />

  <!-- <el-select
    filterable
    v-model="selectVal"
    style="width: 150px"
    size="mini"
    collapse-tags
    multiple
    placeholder="请选择"
    @change="changeValue"
  >
    <el-option v-for="(v,k,index) in options" :key="index" :label="v.label" :value="v.value" />
  </el-select>-->

</template>

<script>
import { GetValues } from '@/api/analysis'

export default {
  name: 'Values',
  components: {
    SmallSelect: () => import('@/components/AnalyseTools/FilterWhere/SmallSelect')
  },
  props: ['data', 'value', 'tableTyp'],
  data() {
    return {
      options: [],
      selectVal: this.value
    }
  },
  watch: {
    data(newV, oldV) {
      this.initValue(newV)
    },
    selectVal(newV, oldV) {
      this.$emit('input', this.selectVal)
    }
  },
  mounted() {
    this.initValue(this.data)
  },
  methods: {
    cleanValues() {
      this.selectVal = []
      this.$emit('input', this.selectVal)
    },
    changeValue() {
      this.$emit('input', this.selectVal)
    },
    initValue(data) {
      GetValues({ 'appid': this.$store.state.baseData.EsConnectID, table: this.tableTyp.toString(), col: data }).then(res => {
        const opt = []
        if (res.data.length > 0) {
          for (const v of res.data) {
            const obj = { label: v.value, value: v.value }
            opt.push(obj)
          }
        }
        this.options = opt
      }).catch(e => {
        console.warn(e)
      })
    }
  }
}
</script>

<style scoped>

</style>
