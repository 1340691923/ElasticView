<template>
  <div class="app-container">

    <div class="filter-container">

      <el-button
        size="mini"
        type="success"
        class="filter-item"
        @click="open = true"
      >{{ $t('新增配置') }}
      </el-button>

    </div>

    <el-table v-loading="tableLoading" :data="tableData">

      <el-table-column
        align="center"
        prop="config"
        label="索引"

      />
      <el-table-column
        align="center"
        prop="remark"
        label="备注"

      />

    </el-table>
    <div class="pagination-container">
      <el-pagination
        v-if="pageshow"
        background
        :current-page="input.page"
        :page-size="input.limit"
        layout="total, sizes, prev, pager, next, jumper"
        :total="count"
        @current-change="search"
        @size-change="handleSizeChange"
      />
    </div>

    <el-dialog :close-on-click-modal="false" :visible.sync="open" :title="$t('新增配置')" @close="closeDialog">
      <el-card class="box-card">
        <el-form label-width="300px" label-position="left">
          <el-form-item label="索引名">
            <el-select
              v-model="form.indexName"
              v-loading="indexSelectLoading"
              class="filter-item"
              filterable
              style="width: 180px"
            >
              <el-option v-for="(v,k,index) in indexList" :key="index" :lable="v" :value="v" />
            </el-select>
          </el-form-item>
          <el-form-item :label="$t('备注')">
            <el-input v-model="form.remark" placeholder="备注" />
          </el-form-item>
        </el-form>
        <div style="text-align:right;">
          <el-button
            size="mini"
            type="danger"
            icon="el-icon-close"
            @click="closeDialog"
          >{{ $t('取消') }}
          </el-button>
          <el-button
            size="mini"
            type="primary"
            icon="el-icon-check"
            @click="add"
          >{{ $t('确认') }}
          </el-button>
        </div>
      </el-card>
    </el-dialog>

  </div>
</template>

<script>

import {  InsertLink, LinkInfoList } from '@/api/datax'
import {IndexNamesAction} from "@/api/es-index";
import {getIndexCfg, setIndexCfg} from "@/api/search";

const defaultForm = {
  indexName: '',
  remark: "",
}

export default {
  name: 'Link',
  data() {
    return {
      indexList:[],
      count: 0,
      pageshow: true,
      tableData: [],
      form: Object.assign({}, defaultForm),
      input: {
        page: 1,
        limit: 10
      },
      tableLoading: false,
      indexSelectLoading:false,
      open: false
    }
  },
  created() {
    this.getIndexList()
    this.search()
  },
  methods: {
    closeDialog() {
      this.open = false
      this.form = Object.assign({}, defaultForm)
    },
    getIndexList(){
      this.indexSelectLoading = true
      const input = {}
      input['es_connect'] = this.$store.state.baseData.EsConnectID
      IndexNamesAction(input).then(res => {
        this.indexSelectLoading = false
        if (res.code == 0) {
          this.indexList = res.data
        } else {
          this.$message({
            type: 'error',
            message: res.msg
          })
        }
      }).catch(err => {
        this.indexSelectLoading = false
        console.log(err)
      })
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
    async add() {
      console.log(123)
      const res = await setIndexCfg(this.form)
      if (res.code != 0) {
        this.$message({
          type: 'error',
          message: res.msg
        })
        return
      }
      this.$message({
        type: 'success',
        message: res.msg
      })
      this.open = false
      this.search()
    },

    async search() {
      this.tableLoading = true
      const input = this.input
      input['es_connect'] = this.$store.state.baseData.EsConnectID
      const res = await getIndexCfg(input)
      if (res.code != 0) {
        this.$message({
          type: 'error',
          message: res.msg
        })
        return
      }
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

    }
    },

}
</script>

<style scoped>

</style>
