<template>
  <div class="app-container">
    <el-card class="box-card">
      <div class="filter-container">
        <el-button class="filter-item" @click="initForm" type="primary">新建数据抽取任务</el-button>
      </div>
      <el-dialog width="95%" :close-on-click-modal="false" @close="closeDialog" :visible.sync="open" title="新建数据抽取任务">
        <el-card class="box-card" >

          <el-form label-width="200px" label-position="left">
            <el-form-item label="任务备注:">
              <el-input v-model="form.remark" style="width: 300px"></el-input>
            </el-form-item>
            <el-form-item label="选择数据源:">
              <el-select @change="changeDataSource" filterable v-model="selectType">
                <el-option v-for="(v,k,index) in linkSelectOpt" :key="index" :value="JSON.stringify(v)" :label="v.remark.concat(`(${v.typ})`)" />
              </el-select>
              <el-button icon="el-icon-plus" type="primary" @click="addTable">添加表</el-button>
            </el-form-item>
            <template v-for="(v,k,index) in form.tables">
             --- {{index}}{{k}}{{v}}
              <tables :tables="tables" :currentIndex="index" :selectType="selectType" @deleteTable="deleteTable" :key="index"></tables>
            </template>

            <el-form-item label="索引名:">
              <el-input v-model="form.indexName" style="width: 300px"></el-input>
            </el-form-item>
            <el-form-item label="是否清空该索引重新导入:">
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
  import {LinkSelectOpt,GetTables} from "@/api/datax"

  const defaultForm = {
    remark: "",
    tables:[],
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
        tables:[]
      }
    },
    mounted() {

    },
    components:{
      "Tables":()=>import("@/views/datax/components/Tables")
    },
    computed:{

    },

    methods: {
      deleteTable(index){
        console.log("index",index)
        this.form.tables.splice(index,1)
      },
      getSelectTypeObj(){
        return JSON.parse(this.selectType)
      },
      async getTables(){
        const res =  await GetTables({id:this.getSelectTypeObj()['id']})
        if(res.code != 0){
          this.$message({
            type: 'error',
            message: res.msg
          })
          return
        }
        this.tables = res.data
        if(this.tables == null) this.tables = []
        this.form.tables = []
      },
      async changeDataSource(){
        await this.getTables()
      },
      addTable(){
        this.form.tables.push({
          selectTable:"",
          cols:[],
        })
      },
      add(){
        console.log("form",this.form)
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
          await  this.getTables()
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
