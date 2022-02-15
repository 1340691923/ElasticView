<template>
  <el-select v-model="indexName" :reserve-keyword="multiple" :collapse-tags="multiple" :disabled="disabled" :placeholder="placeholder" clearable :multiple="multiple" :clearable="clearable" filterable @change="change()">

    <el-option v-if="haveAll == true" label="全部" value="*" />
    <el-option v-for="(indexName, index) in indexList" :key="index" :label="indexName" :value="indexName" />
  </el-select>
</template>

<script>
import { IndexNamesAction } from '@/api/es-index.js'
export default {
  name: 'Index',
  props: {
    indexName: {
      type: String | Array,
      default: ''
    },
    multiple: {
      type: Boolean,
      default: false
    },
    clearable: {
      type: Boolean,
      default: false
    },
    disabled: {
      type: Boolean,
      default: false
    },
    placeholder: {
      type: String,
      default: ''
    },
    haveAll: {
      type: Boolean,
      default: false
    }
  },
  data() {
    return {
      indexList: []
    }
  },
  computed: {

  },
  mounted() {
    this.getIndexList()
  },
  methods: {
    change() {
      this.$emit('change', this.indexName)
    },
    getIndexList() {
      const input = {}
      input['es_connect'] = this.$store.state.baseData.EsConnectID
      IndexNamesAction(input).then(res => {
        if (res.code == 0) {
          this.indexList = res.data
        } else {
          this.$message({
            type: 'error',
            message: res.msg
          })
        }
      }).catch(err => {
        console.log(err)
      })
    }
  }
}
</script>

<style scoped>

</style>
