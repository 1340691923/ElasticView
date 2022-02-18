<template>
  <div>
    <el-table
      :ref="tableRef"
      v-loading="connectLoading"
      :border="border"
      :span-method="spanMethod"
      :data="filterList"
      style="width: 100%"
      :max-height="maxHeight"
    >
      <template v-for="(col, index) in tableInfo" slot-scope="scope">
        <!-- 全部自定义 只需要这个分页 -->
        <slot :name="col.slot" :row="scope.row" />
      </template>
    </el-table>

    <el-pagination
      v-if="pageshow"
      class="pagination-container"
      :current-page="page"
      :page-sizes="[10, 20, 30, 50,100,150,200,500,1000]"
      :page-size="limit"
      layout="total, sizes, prev, pager, next, jumper"
      :total="total"
      @size-change="handleSizeChange"
      @current-change="handlePageChange"
    />
  </div>
</template>

<script>

import { filterData } from '@/utils/table'

export default {

  name: 'MetaAttr',
  props: {
    border: {
      type: Boolean,
      default: true
    },
    tableRef: {
      type: String,
      default: 'tableRef'
    },
    input: {
      type: String,
      default: ''
    },
    connectLoading: {
      type: Boolean,
      default: false
    },
    showTitle: {
      type: String,
      default: ''
    },
    limit: {
      type: Number,
      default: 20
    },
    maxHeight: {
      type: Number,
      default: 700
    },
    tableList: {
      type: Array,
      default: []
    },
    tableInfo: {
      type: Array,
      default: [{ slot: 'operate' }]
    },
    spanMethod: {
      type: Function
    }
  },
  data() {
    return {
      total: 0,
      page: 1,
      pageshow: true,
      list: [],
      trueList: []
    }
  },
  computed: {
    filterList() {
      var table = this.list.slice((this.page - 1) * this.limit, this.page * this.limit)
      this.refreshPage()
      return table
    }
  },
  watch: {
    input(newV, oldV) {
      let list = this.trueList
      list = filterData(list, newV.trim())
      this.total = list.length
      this.list = list
    }
  },
  mounted() {
    this.searchData()
  },
  methods: {
    handleSizeChange(v) {
      this.limit = v
      this.refreshPage()
    },
    handlePageChange(v) {
      this.page = v
      this.refreshPage()
    },
    refreshPage() {
      this.pageshow = false
      this.total = this.list.length
      this.$nextTick(() => {
        this.pageshow = true
      })
    },
    async searchData() {
      let list = this.tableList
      list = filterData(list, this.input.trim())
      this.total = list.length
      this.list = list
      this.trueList = list
    }
  }
}
</script>

<style scoped>

</style>
