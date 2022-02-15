<template>
  <div class="app-container">
    <el-card class="box-card">
      <div class="filter-container">
        <el-tag class="filter-item">请输入关键词</el-tag>
        <el-input v-model="input" class="filter-item" style="width: 300px" clearable @input="search" />
        <el-button type="success" class="filter-item" @click="search">搜索</el-button>
      </div>

      <el-table
        v-loading="connectLoading"
        :data="list"
      >
        <el-table-column
          label="序号"
          align="center"
          fixed
          width="50"
        >
          <template slot-scope="scope">
            {{ scope.$index+1 }}
          </template>
        </el-table-column>

        <el-table-column v-for="(info,index) in tableInfo" :key="index" align="center" :label="tableInfo[index].desc" :width="info.width" :prop="info.data.toString()" :sortable="info.sort">
          <template slot-scope="scope">
            {{ scope.row[info.data.split('->').join('.') ] }}
          </template>
        </el-table-column>

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
        @current-change="handleCurrentChange"
      />
    </el-card>
  </div>
</template>

<script>
import { filterData } from '@/utils/table'
import { CatAction } from '@/api/es'

export default {
  name: 'Index',
  props: {
    catType: {
      type: String,
      default: ''
    },
    tableInfo: {
      type: Array,
      default: []
    }
  },
  data() {
    return {
      total: 0,
      connectLoading: false,
      page: 1,
      limit: 10,
      pageshow: true,
      list: [],
      input: ''
    }
  },
  mounted() {
    this.searchData()
  },
  methods: {
    search() {
      this.page = 1
      this.pageshow = false
      this.searchData()
      this.$nextTick(() => {
        this.pageshow = true
      })
    },
    filterData(list, input) {
      return filterData(list, input)
    },
    // 当每页数量改变
    handleSizeChange(val) {
      console.log(`每页 ${val} 条`)
      this.limit = val
      this.searchData()
    },
    // 当当前页改变
    handleCurrentChange(val) {
      console.log(`当前页: ${val}`)
      this.page = val
      this.searchData()
    },
    searchData() {
      this.connectLoading = true
      const form = {
        cat: this.catType,
        es_connect: this.$store.state.baseData.EsConnectID
      }
      CatAction(form).then(res => {
        if (res.code == 0) {
          let list = res.data
          if (list == null) {
            return
          }
          for (const index in list) {
            const obj = list[index]
            // 把 . 转成 ->
            for (const key in obj) {
              let value = Number(obj[key])
              if (isNaN(value)) {
                value = obj[key]
              }
              list[index][key.split('.').join('->')] = value
            }
          }

          list = filterData(list, this.input.trim())

          this.list = list.filter((item, index) =>
            index < this.page * this.limit && index >= this.limit * (this.page - 1)
          )
          this.total = list.length
        } else {
          this.$message({
            type: 'error',
            message: res.msg
          })
        }
        this.connectLoading = false
      }).catch(err => {
        console.log(err)

        this.connectLoading = false
      })
    }
  }
}
</script>

<style scoped>

</style>
