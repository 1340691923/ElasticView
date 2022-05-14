<template>
  <div>
    <el-input style="width: 30%" v-model="mysqlCol" class="filter-item" readonly></el-input>
    <el-tag>{{ $t('等于') }}</el-tag>
    <el-select class="filter-item" filterable v-model="v" style="width: 30%">
      <el-option v-for="(v,k,index) of esMappingCol" :key="index" :lable="k" :value="v"></el-option>
    </el-select>
  </div>
</template>

<script>
export default {
  name: "mapping",
  props: {
    mysqlCol: {
      type: String,
      default: ""
    },
    esMappingCol: {
      type: Array,
      default: []
    },
    value: {
      type: String,
      default: ""
    }
  },
  data() {
    return {
      v: this.value,
    }
  },
  watch: {
    v(newV, oldV) {
      this.$emit('input', this.v)
    }
  },
}
</script>

<style scoped>

</style>
