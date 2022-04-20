<template>
  <div class="app-container">
    <el-card class="box-card">
      <div class="filter-container">
        <el-button class="filter-item" @click="initForm" type="primary">新建数据抽取任务</el-button>
      </div>

      <el-table
        v-loading="connectLoading"
        :data="tableList"
      >

        <el-table-column
          align="center"
          prop="id"
          label="id"
          width="80">
        </el-table-column>
        <el-table-column
          align="center"
          prop="id"
          label="备注"
          width="150">
        </el-table-column>
        <el-table-column
          align="center"
          prop="table_name"
          label="表名"
          width="200">
        </el-table-column>
        <el-table-column
          align="center"
          prop="index_name"
          label="索引名"
          width="250">
        </el-table-column>
        <el-table-column
          align="center"
          prop="status"
          label="状态"
          width="200">
        </el-table-column>
        <el-table-column
          align="center"
          prop="error_msg"
          label="错误信息"
          width="300">
        </el-table-column>

        <el-table-column
          align="center"
          prop="created"
          label="创建时间"
          width="150">
        </el-table-column>
        <el-table-column
          align="center"
          prop="updated"
          label="修改时间"
          width="150">
        </el-table-column>
        <el-table-column align="center" label="操作" fixed="right" width="300">
          <template slot-scope="scope">
           <el-button type="danger" size="small" icon="el-icon-close" @click="cancel(scope.row.id)">{{$t('取消')}}</el-button>
          </template>
        </el-table-column>
      </el-table>
      <el-dialog width="95%" :close-on-click-modal="false" @close="closeDialog" :visible.sync="open" title="新建数据抽取任务">
        <el-card class="box-card">
          <el-form label-width="200px" label-position="left">
            <el-form-item label="任务备注:">
              <el-input v-model="form.remark" style="width: 300px"></el-input>
            </el-form-item>
            <el-form-item label="选择数据源:">
              <el-select @change="changeTable" filterable v-model="form.selectType">
                <el-option v-for="(v,k,index) in linkSelectOpt" :key="index" :value="JSON.stringify(v)"
                           :label="v.remark.concat(`(${v.typ})`)"/>
              </el-select>
            </el-form-item>
            <el-form-item :label="getLabelName">
              <el-select @change="GetTableColumns" filterable v-model="form.selectTable">
                <el-option v-for="(v,k,index) in tables" :key="index" :value="v" :label="v"/>
              </el-select>
            </el-form-item>
            <el-form-item label="表字段：">
              <div class="col-transfer">
                <el-transfer
                  @change="changeTbCols"
                  v-model="form.cols.tableCols"
                  :titles="['全部字段', '当前ES索引字段']"
                  :button-texts="['移除', '添加']"
                  filterable
                  :filter-method="filterMethod"
                  filter-placeholder="请操作字段"
                  :data="allCols"
                />
              </div>
            </el-form-item>
            <el-form-item label="索引名:">

              <el-select @change="changeIndex" v-loading="indexSelectLoading" class="filter-item" filterable v-model="form.indexName"
                         style="width: 350px">
                <el-option v-for="(v,k,index) in indexList" :key="index" :lable="v" :value="v"></el-option>
              </el-select>
            </el-form-item>

            <el-form-item v-if="showMapping" label="字段映射:">
              <div v-if="form.cols.tableCols.length == 0">
                <el-alert
                  title="请先选择表字段"
                  type="warning"
                  ></el-alert>
              </div>
              <div class="el-row" v-for="(v,k,index) in form.cols.tableCols">
                <mapping  v-model="form.cols.esCols[k].col" :esMappingCol="esCols" :mysqlCol="v"></mapping>
              </div>
            </el-form-item>

            <el-form-item label="是否清空该索引重新导入数据:">
              <el-radio v-model="form.reset" class="filter-item" :label="Boolean(true)">是</el-radio>
              <el-radio v-model="form.reset" class="filter-item" :label="Boolean(false)">否</el-radio>
            </el-form-item>
            <el-form-item label="入库批次数量:">
              <el-input type="number" v-model="form.bufferSize" style="width: 300px"></el-input>
            </el-form-item>
            <el-form-item label="入库轮循间隔时间(单位秒):">
              <el-input type="number" v-model="form.flushInterval" style="width: 300px"></el-input>
            </el-form-item>
          </el-form>
          <div style="text-align:right;">
            <el-button type="danger" icon="el-icon-close" @click="closeDialog">取消</el-button>
            <el-button type="primary" icon="el-icon-check" @click="add">确认</el-button>
          </div>
        </el-card>
      </el-dialog>
    </el-card>
  </div>
</template>

<script>
  import {GetTableColumns, GetTables, LinkSelectOpt,Transfer,TransferLogList} from "@/api/datax"
  import {IndexNamesAction} from "@/api/es-index"
  import {ListAction} from '@/api/es-map'
  const defaultForm = {
    selectType:"{}",
    remark: "",
    selectTable: "",
    cols: {
      tableCols:[],
      esCols:[]
    },
    indexName: "",
    reset: true,
    bufferSize: 5000,
    flushInterval: 5000
  }

  export default {
    name: "list",
    data() {
      return {
        connectLoading:false,
        test:"",
        showMapping:false,
        indexSelectLoading: false,
        open: false,
        form: Object.assign({}, defaultForm),
        linkSelectOpt: {},
        allCols: [],
        tables: [],
        indexList: [],
        esCols:[],
        ver:6,
        tableList:[],
      }
    },
    async mounted() {
      this.getIndexList()
      this.getList()
    },
    computed: {
      getLabelName() {
        switch (this.getSelectTypeObj()["typ"]) {
          case "mysql":
          case "clickhouse":
            return "表名:"
            break
          /*case "mongodb":
            return "集合名:"
            break*/
          default:
            return "未知:"
        }
      }
    },
    components:{
      Mapping:()=>import("@/views/datax/components/mapping")
    },
    watch: {
      'form.cols.tableCols'(val, oldVal) {
        if(val !=undefined){
          this.form.cols.esCols=[]
          for(let i in  val){
            this.form.cols.esCols.push({col:"",tbCol:val[i]})
          }
        }
      },
    },
    methods: {
      cancel(id){
        alert(id)
      },
      async getList(){
        const res = await TransferLogList()
        if (res.code != 0) {
          this.$message({
            type: 'error',
            message: res.msg
          })
          return
        }
        this.tableList = res.data
      },
      changeTbCols(v){
        this.refreshshowMapping()
      },
      refreshshowMapping() {
        this.showMapping = false
        this.$nextTick(() => {
          this.showMapping = true
        })
      },
      async changeIndex() {

        const input = {}
        input['es_connect'] = this.$store.state.baseData.EsConnectID

        input['index_name'] = this.form.indexName

        const {data, code, msg} = await ListAction(input)

        if (code != 0) {
          this.$message({
            type: 'error',
            message: msg
          })
          return
        }

        this.ver = data.ver

        try {
          switch (this.ver) {
            case 6:
              const mappings = Object.keys(data.list[this.form.indexName].mappings)
              if (mappings.length > 0) {
                let type_name = mappings[0]
                this.esCols = Object.keys(data.list[this.form.indexName].mappings[type_name].properties)
              }
              break
            case 7:
            case 8:
              this.esCols = Object.keys(data.list[this.form.indexName].mappings.properties)
              break
          }
        }catch (e) {
          this.esCols = []
        }
        this.refreshshowMapping()

      },
      async getIndexList() {
        const input = {}
        input['es_connect'] = this.$store.state.baseData.EsConnectID
        this.indexSelectLoading = true
        const res = await IndexNamesAction(input)
        this.indexSelectLoading = false
        if (res.code != 0) {
          this.$message({
            type: 'error',
            message: res.msg
          })
          return
        }
        this.indexList = res.data
      },
      async add() {
        let form = this.form
        console.log("this.form",this.form)
        form['es_connect'] = this.$store.state.baseData.EsConnectID
        const res =  await Transfer(form)
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
        this.form = Object.assign({}, defaultForm)
        this.getList()
        this.closeDialog()
      },
      filterMethod(query, item) {
        return item.label.indexOf(query) > -1
      },
      async changeTable() {
        await this.getTables()
        await this.GetTableColumns()

      },
      async GetTableColumns() {
        const res = await GetTableColumns({id: this.getSelectTypeObj()['id'], table_name: this.form.selectTable})
        if (res.code != 0) {
          this.$message({
            type: 'error',
            message: res.msg
          })
          return
        }
        if (res.data == null) res.data = []
        this.allCols = []
        for (let v of res.data) {
          const obj = {
            key: v.Field,
            label: v.Comment == '' ? v.Field : `${v.Field}【${v.Comment}】`,
            disabled: false
          }
          this.allCols.push(obj)
        }
        this.form.cols.tableCols = []
        this.form.cols.esCols = []
      },
      async getTables() {
        const res = await GetTables({id: this.getSelectTypeObj()['id']})
        if (res.code != 0) {
          this.$message({
            type: 'error',
            message: res.msg
          })
          return
        }
        this.tables = res.data

        if (this.tables == null) this.tables = []
        if (this.tables.length > 0) {
          this.form.selectTable = this.tables[0]
        }

      },
      getSelectTypeObj() {
        return JSON.parse(this.form.selectType)
      },
      closeDialog() {
        this.open = false
      },
      async initForm() {
        this.form.remark = ""
        const res = await LinkSelectOpt()
        if (res.code != 0) {
          return
        }
        if (res.data == null) res.data = []
        this.linkSelectOpt = res.data

        if (this.linkSelectOpt.length > 0) {
          this.form.selectType = JSON.stringify(this.linkSelectOpt[0])
          await this.getTables()
          await this.GetTableColumns()
        }
        this.open = true

      }
    }
  }
</script>

<style scoped>
  .col-transfer >>> .el-transfer-panel {
    width: 35%;
  }

  .el-row {
    margin-bottom: 20px;

  &
  :last-child {
    margin-bottom: 0;
  }

  }



</style>
