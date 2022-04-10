<template>
  <div class="app-container">
    <el-card class="box-card">
      <div class="filter-container">
        <el-select v-model="input.status"  class="filter-item" clearable filterable placeholder="请选择状态" @change="searchData(1)">
          <el-option
            v-for="(v,k,index) in statusMap"
            :key="index"
            :label="v"
            :value="Number(k)"
          />
        </el-select>
        <el-button icon="el-icon-search" @click="searchData(1)" class="filter-item">刷新</el-button>
      </div>

      <el-table
        :data="tableData"
        v-loading="tableLoading"
        :cell-class-name="rowClass"
        style="width: 100%">
        <el-table-column label="ID" sortable width="80" prop="id" align="center"/>
        <el-table-column label="任务ID" sortable width="100" prop="task_id" align="center"/>
        <el-table-column align="center" label="执行命令" width="230">
          <template slot-scope="scope">
            {{ scope.row.action }}
          </template>
        </el-table-column>
        <el-table-column align="center" label="执行时间" width="230">
          <template slot-scope="scope">
            {{ scope.row.exec_time_format }}
          </template>
        </el-table-column>
        <el-table-column align="center" label="状态" width="130">
          <template slot-scope="scope">
            <el-tag v-if="scope.row.status == 0">{{ statusMap[scope.row.status] }}</el-tag>
            <el-tag v-if="scope.row.status == 1" type="info">{{ statusMap[scope.row.status] }}</el-tag>
            <el-tag v-if="scope.row.status == 2" type="success">{{ statusMap[scope.row.status] }}</el-tag>
            <el-tag v-if="scope.row.status == 3" type="danger">{{ statusMap[scope.row.status] }}</el-tag>
            <el-tag v-if="scope.row.status == 4" type="warning">{{ statusMap[scope.row.status] }}</el-tag>
          </template>
        </el-table-column>
        <el-table-column align="center" label="任务信息" width="230">
          <template slot-scope="scope">
            {{ scope.row.msg }}
          </template>
        </el-table-column>
        <el-table-column align="center" label="创建时间" width="230">
          <template slot-scope="scope">
            {{ scope.row.create_time }}
          </template>
        </el-table-column>
        <el-table-column align="center" label="修改时间" width="230">
          <template slot-scope="scope">
            {{ scope.row.update_time }}
          </template>
        </el-table-column>
        <el-table-column  label="操作"  width="100" fixed="right" align="center"  >
          <template slot-scope="scope">
            <el-button
              v-if="scope.row.status == 0 || scope.row.status == 1"
              size="mini"
              type="primary"
              @click="cancelAction(scope.row.task_id)">
              取消
            </el-button>
          </template>
        </el-table-column>
      </el-table>
      <div class="pagination-container">
        <el-pagination
          v-if="pageshow"
          background
          :current-page="input.page"
          :page-size="input.limit"
          layout="total, sizes, prev, pager, next, jumper"
          :total="count"
          @current-change="searchData"
          @size-change="handleSizeChange"
        />
      </div>
    </el-card>
  </div>
</template>

<script>
import { deepClone } from '@/utils'
import { ListAction,CancelAction} from '@/api/timing'
export default {
  name: 'day',
  data() {
    return {

      statusMap: {
        0: '等待执行',
        1: '正在执行',
        2: '成功',
        3: '失败',
        4: '取消',
      },
      count: 0,
      pageshow: true,
      input: {
        status: null,
        action: null,
        limit: 10,
        date: [],
        page: 1
      },
      taskMap: {},
      tableData: [],
      tableLoading: false,
      downloadLoading: false,
    }
  },
  created() {
    this.searchData()
  },
  computed: {

  },
  components: {

  },
  methods: {
    async cancelAction(id){
      const {msg,code} = await CancelAction({taskId:id})

      if(code != 0){
        this.$message({
          type: 'error',
          message:msg
        })
      }else{
        this.$message({
          type: 'success',
          message:msg
        })
        this.searchData()
      }
    },
    handlePageChange(v) {
      this.input.page = v
      this.refreshPage()
    },
    refreshPage() {
      this.pageshow = false
      this.count = this.tableData.length
      this.$nextTick(() => {
        this.pageshow = true
      })
    },
    handleSizeChange(v) {
      this.input.limit = v
      this.refreshPage()
    },

    rowClass({ row, column, rowIndex, columnIndex }) {
      return 'aaa'
    },

    searchData(page) {
      !page ? this.input.page = 1 : this.input.page = page
      this.tableLoading = true
      let input = deepClone(this.input)

      if(input.status === "") input.status = null
      if(input.action === "") input.action = null

      ListAction(input).then(res => {
        if (res.code == 0) {
          this.tableData = res.data.data
          this.count = res.data.count
          if (this.tableData == null) {
            this.tableData = []
          }
          this.$message({
            type: 'success',
            message: res.msg
          })
          this.tableLoading = false
        } else {
          this.$message({
            type: 'error',
            message: res.msg
          })
        }
        this.tableLoading = false
      }).catch(err => {
        this.tableLoading = false
        console.log(err)
      })
    }
  }
}
</script>

<style>
.bgColor {
  background: #6959CD;
  color: white
}
</style>
