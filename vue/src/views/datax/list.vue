<template>
  <div class="app-container">
    <el-card class="box-card">
      <div class="filter-container">
        <el-button class="filter-item" @click="initForm" type="primary">新建数据抽取任务</el-button>
      </div>
      <el-dialog width="95%" :close-on-click-modal="false" @close="closeDialog" :visible.sync="open" title="新建数据抽取任务">
        <el-card class="box-card">
          <el-form label-width="200px" label-position="left">
            <el-form-item label="任务备注:">
              <el-input v-model="form.remark" style="width: 300px"></el-input>
            </el-form-item>
            <el-form-item label="选择数据源:">
              <el-select @change="changeTable"  filterable v-model="selectType">
                <el-option v-for="(v,k,index) in linkSelectOpt" :key="index" :value="JSON.stringify(v)" :label="v.remark.concat(`(${v.typ})`)" />
              </el-select>
            </el-form-item>
            <el-form-item :label="getLabelName">
              <el-select @change="GetTableColumns" filterable v-model="form.selectTable">
                <el-option v-for="(v,k,index) in tables" :key="index" :value="v" :label="v" />
              </el-select>
            </el-form-item>
            <el-form-item label="表字段：">
              <div class="col-transfer">
                <el-transfer
                  v-model="form.cols"
                  :titles="['全部字段', '当前ES索引字段']"
                  :button-texts="['移除字段', '添加字段']"
                  filterable
                  :filter-method="filterMethod"
                  filter-placeholder="请操作成员"
                  :data="allCols"
                />
              </div>
            </el-form-item>
            <el-form-item label="索引名:">
              <el-input v-model="form.indexName" style="width: 300px"></el-input>
            </el-form-item>
            <el-form-item label="是否建立新索引:">
              <el-radio  v-model="form.reset" class="filter-item" label="yes" >是</el-radio>
              <el-radio  v-model="form.reset" class="filter-item" label="no" >否</el-radio>
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
  import {LinkSelectOpt,Tables,GetTableColumns} from "@/api/datax"

  const defaultForm = {
    remark: "",
    selectTable:"",
    cols:[],
    indexName:"",
    reset:"yes",
    bufferSize:5000,
    flushInterval:5000,
  }

  export default {
    name: "list",
    data() {
      return {
        open: false,
        selectType:"{}",
        form: Object.assign({}, defaultForm),
        linkSelectOpt:{

        },
        allCols:[],
        tables:[],
      }
    },
    mounted() {

    },
    computed:{
      getLabelName(){
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
    watch: {
      'form.selectTable'(newV, oldV) {
        this.form.indexName = newV
      }
    },
    methods: {
      add(){
        console.log("form",this.form)
      },
      filterMethod(query, item) {
        return item.label.indexOf(query) > -1
      },
      async changeTable(){
        await this.getTables()
        await this.GetTableColumns()
      },
      async GetTableColumns(){
        const res =  await GetTableColumns({id:this.getSelectTypeObj()['id'],table_name:this.form.selectTable})
        if(res.code != 0){
          this.$message({
            type: 'error',
            message: res.msg
          })
          return
        }
        if(res.data == null)res.data = []
        this.allCols = []
        for(let v of res.data){
          const obj = {
            key: v.Field,
            label:v.Comment==''?v.Field:`${v.Field}【${v.Comment}】`,
            disabled: false
          }
          this.allCols.push(obj)
        }
      },
      async getTables(){
        const res =  await Tables({id:this.getSelectTypeObj()['id']})
        if(res.code != 0){
          this.$message({
            type: 'error',
            message: res.msg
          })
          return
        }
        this.tables = res.data

        if(this.tables == null) this.tables = []
        if(this.tables.length >0){
          this.form.selectTable = this.tables[0]
        }

      },
      getSelectTypeObj(){
        return JSON.parse(this.selectType)
      },
      closeDialog() {
        this.open = false
      },
      async initForm() {
        this.form.remark = ""
        const res =  await LinkSelectOpt()
        if(res.code != 0){
          return
        }
        if(res.data == null) res.data = []
        this.linkSelectOpt = res.data

        if(this.linkSelectOpt.length > 0 ){
          this.selectType = JSON.stringify(this.linkSelectOpt[0])
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
  width:350px;
}
</style>
