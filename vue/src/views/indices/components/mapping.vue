<template>
  <div>
    <el-dialog :close-on-click-modal="false" :visible.sync="open" :title="title.concat(`【${indexName}】`)" width="60%" @close="closeDialog">
        <div class="app-container">
          <div class="filter-container">

            <el-tag class="filter-item"  v-if="showTypeName">type名：</el-tag>
            <el-input v-if="showTypeName" :readonly="typeReadonly"  style="width: 200px" class="filter-item" v-model="type_name" ></el-input>
            <el-tag class="filter-item" >dynamic：</el-tag>
            <el-select v-model="dynamic" class="filter-item" size="small">
              <el-option label="动态映射" value="true" />
              <el-option label="静态映射" value="false" />
              <el-option label="严格映射" value="strict" />
            </el-select>

            <el-button v-loading="loading" icon="el-icon-check" type="success" size="small" @click="saveMappinng" class="filter-item">保存/修改映射</el-button>
          </div>
          <VueJsonHelper
            v-if="showVueJsonHelper"
            :size="size"
            :names="names"
            :json-str="jsonStr"
            :root-flag="rootFlag"
            :open-flag="openFlag"
            :back-top-flag="backTopFlag"
            @jsonListener="jsonListener"
            :shadow-flag="false"
            :border-flag="false"
          />
        </div>
    </el-dialog>
  </div>
</template>

<script>
  import {ListAction,UpdateMappingAction} from '@/api/es-map'

  export default {
    name: 'Mapping',

    props: {
      indexName: {
        type: String,
        default: ''
      },
      open: {
        type: Boolean,
        default: false
      },
      title: {
        type: String,
        default: '新增映射结构'
      }
    },

    data() {
      return {
        loading:false,
        dynamic:"false",
        drawerShow: false,
        connectLoading: false,
        size: 'small',
        names: [
          {key: "Root", name: "properties"},
          {key: "type", name: "数据类型"},
          {key: "format", name: "时间格式化"},
          {key: "flag", name: "标示"},
          {key: "EduexPerience", name: "教育经历"},
          {key: "year", name: "年份"},
          {key: "education", name: "学历"},
        ],
        rootFlag: true,
        openFlag: true,
        backTopFlag: false,
        jsonStr: "{}",
        showVueJsonHelper:false,
        type_name:"",
        ver:6,
        showTypeName:false,
        typeReadonly:false
      }
    },
    mounted() {
      this.init()
    },
    methods: {
      async saveMappinng() {
        let properties = {}
        try {
          properties = JSON.parse(this.jsonStr)
        } catch (e) {
          this.$message({
            type: 'error',
            message: 'JSON格式不正确'
          })
          return
        }

        const input = {}
        input['es_connect'] = this.$store.state.baseData.EsConnectID
        input['index_name'] = this.indexName
        let activeData = {}
        activeData["properties"] = properties
        activeData["dynamic"] = this.dynamic

        switch (this.ver) {
          case 6:
            input['properties'] = activeData
            input['type_name'] = this.type_name
            break
          case 7:
          case 8:
            input['properties'] = activeData
            break
        }

        this.loading = true
        const { data, code, msg } = await UpdateMappingAction(input)
        this.loading = false
        if (code == 0) {
          this.$message({
            type: 'success',
            message: msg
          })
        } else {
          this.$message({
            type: 'error',
            message: msg
          })
        }
      },
      refreshVueJsonHelper() {
        console.log("234")
        this.showVueJsonHelper = false
        this.$nextTick(() => {
          this.showVueJsonHelper = true
        })
      },
      async init() {

        const input = {}
        input['es_connect'] = this.$store.state.baseData.EsConnectID

        input['index_name'] = this.indexName

        const {data, code, msg} = await ListAction(input)

        if (code != 0) {
          this.$message({
            type: 'error',
            message: res.msg
          })
          return
        }

        this.ver = data.ver

        switch (this.ver) {
          case 6:
            const mappings = Object.keys(data.list[this.indexName].mappings)

            this.showTypeName = true
            if (mappings.length == 0) {
              this.typeReadonly = false
            }else{

              this.type_name = mappings[0]
              this.typeReadonly = true
              this.dynamic = data.list[this.indexName].mappings[this.type_name].hasOwnProperty("dynamic")?data.list[this.indexName].mappings[this.type_name]["dynamic"]:"false"
              this.jsonStr = data.list[this.indexName].mappings[this.type_name].hasOwnProperty("properties")?JSON.stringify(data.list[this.indexName].mappings[this.type_name].properties):"{}"
              this.jsonStr = JSON.stringify(data.list[this.indexName].mappings[this.type_name].properties)
            }
            break
          case 7:
          case 8:
            this.dynamic = data.list[this.indexName].mappings.hasOwnProperty("dynamic")?data.list[this.indexName].mappings["dynamic"]:"false"
            this.jsonStr = data.list[this.indexName].mappings.hasOwnProperty("properties")?JSON.stringify(data.list[this.indexName].mappings.properties):"{}"
            break
        }
        this.refreshVueJsonHelper()

      },
      jsonListener(data) {
        this.jsonStr = JSON.stringify(data)
      },
      closeDialog() {
        this.$emit('close')
      }
    }
  }
</script>

<style scoped>

</style>
