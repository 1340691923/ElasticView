<template>
  <div class="app-container">
      <div class="filter-container">
        <el-button @click="getList" class="filter-item" icon="el-icon-refresh" type="primary">{{ $t('刷新') }}</el-button>
        <el-button v-loading="openTaskLoading" :disabled="openTaskLoading" class="filter-item" icon="el-icon-plus"
                   @click="initForm" type="warning">
          {{ $t('新建数据抽取任务') }}
        </el-button>
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
          :label="$t('备注')"
          width="100">
        </el-table-column>
        <el-table-column
          align="center"
          prop="table_name"
          :label="$t('表名')"
          width="200">
        </el-table-column>
        <el-table-column
          align="center"
          prop="index_name"
          :label="$t('索引名')"
          width="200">
        </el-table-column>

        <el-table-column
          align="center"
          prop="dbcount"
          :label="$t('源数据条数')"
          width="100">
        </el-table-column>
        <el-table-column
          align="center"
          prop="escount"
          :label="$t('已导入数据条数')"
          width="100">
        </el-table-column>
        <el-table-column
          align="center"
          prop="crontab_spec"
          :label="$t('定时任务')"
          width="200">
        </el-table-column>

        <el-table-column
          :label="$t('状态')"
          align="center"
          width="200"
        >
          <template slot-scope="scope">
            <el-tag v-if="scope.row.status == '数据导入成功'" type="success">
              {{ $t(scope.row.status) }}
            </el-tag>
            <el-tag v-else-if="scope.row.status == '数据正在导入中...' ||scope.row.status ==  '正在运行中...'" type="warning">
              {{ $t(scope.row.status) }}
            </el-tag>

            <el-tag v-else type="danger">
              {{ $t(scope.row.status) }}
            </el-tag>
          </template>
        </el-table-column>

        <el-table-column
          align="center"
          prop="error_msg"
          :label="$t('附带信息')"
          width="200">
        </el-table-column>

        <el-table-column
          align="center"
          prop="created"
          :label="$t('创建时间')"
          width="150">
        </el-table-column>
        <el-table-column
          align="center"
          prop="updated"
          :label="$t('修改时间')"
          width="150">
        </el-table-column>
        <el-table-column align="center" :label="$t('操作')" fixed="right" width="300">
          <template slot-scope="scope">

            <el-button :disabled="scope.row.status != '正在运行中...' && scope.row.status != '数据正在导入中...'" type="danger"
                       size="small" icon="el-icon-close"
                       @click="cancel(scope.row.id)">{{ $t('取消') }}
            </el-button>
            <el-button type="danger" size="small" icon="el-icon-delete" @click="deleteById(scope.row.id)">{{ $t('删除') }}
            </el-button>
          </template>
        </el-table-column>
      </el-table>
      <el-dialog width="95%" :close-on-click-modal="false" @close="closeDialog" :visible.sync="open"
                 :title="$t('新建数据抽取任务')">
        <el-card class="box-card">
          <el-form label-width="400px" label-position="left">
            <el-form-item :label="$t('任务备注:')">
              <el-input v-model="form.remark" style="width: 300px"></el-input>
            </el-form-item>
            <el-form-item :label="$t('选择数据源:')">
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
            <el-form-item :label="$t('表字段：')">
              <div v-loading="transferLoading" class="col-transfer">
                <el-transfer
                  @change="changeTbCols"
                  v-model="form.cols.tableCols"
                  :titles="[$t('全部字段'), $t('当前所选表字段')]"
                  :button-texts="[$t('移除'), $t('添加')]"
                  filterable
                  :filter-method="filterMethod"
                  :filter-placeholder="$t('请操作字段')"
                  :data="allCols"
                />
              </div>
            </el-form-item>
            <el-form-item v-if="showAutoIncrementId" :label="$t('自增主键（注意：若存在连续自增的主键并该表非分区表则填)')">
              <el-select filterable clearable v-model="form.autoIncrementId">
                <el-option v-for="(v,k,index) in allCols" :key="index" :value="v.key" :label="v.label"/>
              </el-select>
            </el-form-item>
            <el-form-item :label="$t('索引名:')">
              <el-select @change="changeIndex" v-loading="indexSelectLoading" class="filter-item" filterable
                         v-model="form.indexName"
                         style="width: 350px">
                <el-option v-for="(v,k,index) in indexList" :key="index" :lable="v" :value="v"></el-option>
              </el-select>
              <el-button
                id="new-index"
                type="success"
                class="filter-item"
                icon="el-icon-refresh"
                @click="getIndexList"
              >{{ $t('刷新') }}
              </el-button>
              <el-button
                id="new-index"
                type="success"
                class="filter-item"
                icon="el-icon-plus"
                @click="openSettingDialog('','add')"
              >{{ $t('新建索引') }}
              </el-button>
              <el-button
                v-if="form.indexName != ''"
                type="primary"
                size="small"
                icon="el-icon-setting"
                @click="openSettingDialog(form.indexName,'update')"
              >{{ $t('修改配置') }}
              </el-button>
              <el-button
                v-if="form.indexName != ''"
                type="primary"
                size="small"
                icon="el-icon-circle-plus-outline"
                @click="openMappingEditDialog(form.indexName,false)"
              >{{ $t('修改映射') }}
              </el-button>
            </el-form-item>


            <el-form-item v-if="showMapping" :label="$t('字段映射:')">
              <div v-if="form.cols.tableCols.length == 0">
                <el-alert
                  :title="$t('请先选择表字段')"
                  type="warning"
                ></el-alert>
              </div>
              <div class="el-row" v-for="(v,k,index) in form.cols.tableCols">
                <mapping v-model="form.cols.esCols[k].col" :esMappingCol="esCols" :mysqlCol="v"></mapping>
              </div>
            </el-form-item>

            <el-form-item :label="$t('是否清空该索引重新导入数据:')">
              <el-radio v-model="form.reset" class="filter-item" :label="Boolean(true)">{{ $t('是') }}</el-radio>
              <el-radio v-model="form.reset" class="filter-item" :label="Boolean(false)">{{ $t('否') }}</el-radio>
            </el-form-item>
            <el-form-item :label="$t('协程数:')">
              <el-input type="number" v-model.number="form.goNum" style="width: 300px"></el-input>
            </el-form-item>
            <el-form-item :label="$t('源数据库每次limit条数:')">
              <el-input type="number" v-model.number="form.bufferSize" style="width: 300px"></el-input>
            </el-form-item>
            <el-form-item :label="$t('es入库批次数量:')">
              <el-input type="number" v-model.number="form.esBufferSize" style="width: 300px"></el-input>
            </el-form-item>
            <el-form-item :label="$t('es入库轮循间隔时间:')">
              <el-input type="number" v-model.number="form.esFlushInterval" style="width: 300px"></el-input>
            </el-form-item>

            <el-form-item :label="$t('数据库连接池最大打开的连接数(设置为0表示不限制):')">
              <el-input type="number" v-model.number="form.maxOpenConns" style="width: 300px"></el-input>
            </el-form-item>
            <el-form-item :label="$t('数据库连接池最大空闲的连接数:')">
              <el-input type="number" v-model.number="form.maxIdleConns" style="width: 300px"></el-input>
            </el-form-item>

            <el-form-item :label="$t('计划任务表达式:')">
              <el-autocomplete
                style="width: 300px"
                :placeholder="$t('计划任务表达式（若无需定时跑，则不填）')"
                clearable
                :fetch-suggestions="querySearch"

                v-model="form.crontab_spec"
              >
                <i
                  class="el-icon-edit el-input__icon"
                  slot="suffix"
                >
                </i>
                <template slot-scope="{ item }">
                  <span>{{ item.value }}-{{ item.data }}</span>
                </template>

              </el-autocomplete>

            </el-form-item>

          </el-form>
          <div style="text-align:right;">
            <el-button type="danger" icon="el-icon-close" @click="closeDialog">{{ $t('取消') }}</el-button>
            <el-button :disabled="addLoading" type="primary" icon="el-icon-check" v-loading="addLoading" @click="add">
              {{ $t('确认') }}
            </el-button>
          </div>
        </el-card>
      </el-dialog>

    <settings
      style="z-index: 99999"
      v-if="openSettings"
      :index-name="form.indexName"
      :settings-type="settingsType"
      @finished="getIndexList"
      :open="openSettings"
      @close="closeSettings"
    />
    <mappings
      @finished="changeIndex"
      style="z-index: 99999"
      v-if="openMappings"
      :index-name="form.indexName"
      :title="mappingTitle"
      :open="openMappings"
      @close="closeMappings"
    />
  </div>
</template>

<script>
import {CancelTaskById, GetTableColumns, GetTables, LinkSelectOpt, Transfer, TransferLogList} from "@/api/datax"
import {IndexNamesAction} from "@/api/es-index"
import {ListAction} from '@/api/es-map'
import {DeleteTaskById} from "../../api/datax";
import {filterData} from "@/utils/table";

const defaultForm = {
  crontab_spec: "",
  autoIncrementId: "",
  selectType: "{}",
  remark: "",
  selectTable: "",
  cols: {
    tableCols: [],
    esCols: []
  },
  indexName: "",
  reset: true,
  bufferSize: 100,
  esFlushInterval: 3,
  esBufferSize: 20000,
  maxOpenConns: 50,
  maxIdleConns: 50,
  goNum: 50,
  type_name: ""
}

export default {
  name: "list",
  data() {
    return {
      showAutoIncrementId: false,
      openSettings: false,
      openMappings: false,
      settingsType: 'add',
      mappingTitle: '',

      addLoading: false,
      openTaskLoading: false,
      connectLoading: false,
      test: "",
      showMapping: false,
      indexSelectLoading: false,
      open: false,
      form: Object.assign({}, defaultForm),
      linkSelectOpt: {},
      allCols: [],
      tables: [],
      indexList: [],
      esCols: [],
      ver: 6,
      tableList: [],
      indexName: "",
      crontabTishiList: [
        {'value': "0 30 2 * * *", 'data': "每天凌晨2：30跑一次"},
        {'value': "0 */5 * * * *", 'data': "每5分钟跑一次"},
      ],
      transferLoading: false,
      max: 20
    }
  },
  created() {
    this.getIndexList()
    this.getList()
  },
  components: {
    'Mapping': () => import("@/views/datax/components/mapping"),
    'Settings': () => import('@/views/indices/components/settings'),
    'Mappings': () => import('@/views/indices/components/mapping'),
    'BackToTop': () => import('@/components/BackToTop/index'),
    'JsonEditor': () => import('@/components/JsonEditor/index'),
    'Alias': () => import('@/views/indices/components/alias'),
    'IndexSelect': () => import('@/components/index/select')
  },
  computed: {
    getLabelName() {
      switch (this.getSelectTypeObj()["typ"]) {
        case "mysql":
        case "clickhouse":
          return this.$t("表名:")
        /*case "mongodb":
          return "集合名:"
          break*/
        default:
          return "未知:"
      }
    }
  },

  watch: {
    'form.cols.tableCols'(val, oldVal) {
      if (val != undefined) {
        this.form.cols.esCols = []
        for (let i in val) {
          this.form.cols.esCols.push({col: "", tbCol: val[i]})
        }
      }
    },
  },
  methods: {
    querySearch(queryString, cb) {

      let queryData = JSON.parse(JSON.stringify(this.crontabTishiList))
      if (queryString == undefined) queryString = ""
      if (queryString.trim() == '') {
        if (queryData.length > this.max) {
          cb(queryData.slice(0, this.max))
        } else {
          cb(queryData)
        }
        return;
      }

      queryData = filterData(queryData, queryString.trim())

      if (queryData.length > this.max) {
        cb(queryData.slice(0, this.max))
      } else {
        cb(queryData)
      }
    },
    openSettingDialog(indexName, settingsType) {
      this.form.indexName = indexName
      this.settingsType = settingsType
      this.openSettings = true
    },
    closeSettings() {
      this.settingsType = 'add'
      this.openSettings = false
    },
    closeMappings() {
      this.mappingTitle = 'add'
      this.openMappings = false
    },
    openMappingEditDialog(indexName, haveMapping) {
      if (haveMapping) {
        this.mappingTitle = this.$t('新增字段')
      } else {
        this.mappingTitle = this.$t('新增映射结构')
      }
      this.form.indexName = indexName

      this.openMappings = true
    },
    async cancel(id) {
      const {code, msg} = await CancelTaskById({id: id})
      if (code != 0) {
        this.$message({
          type: 'error',
          message: msg
        })
        return
      }
      this.$message({
        type: 'success',
        message: msg
      })
      this.getList()
      return
    },
    async deleteById(id) {
      const {code, msg} = await DeleteTaskById({id: id})
      if (code != 0) {
        this.$message({
          type: 'error',
          message: msg
        })
        return
      }
      this.$message({
        type: 'success',
        message: msg
      })
      this.getList()
      return
    },
    async getList() {
      this.connectLoading = true
      const res = await TransferLogList({es_connect: this.$store.state.baseData.EsConnectID})
      this.connectLoading = false
      if (res.code != 0) {
        this.$message({
          type: 'error',
          message: res.msg
        })
        return
      }
      this.tableList = res.data
    },
    changeTbCols(v) {
      this.refreshshowMapping()
    },
    refreshShowAutoIncrementId() {

      this.showAutoIncrementId = false
      this.$nextTick(() => {
        this.showAutoIncrementId = true
      })
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
              this.form.type_name = mappings[0]
              this.esCols = Object.keys(data.list[this.form.indexName].mappings[this.form.type_name].properties)
            }
            break
          case 7:
          case 8:
            this.esCols = Object.keys(data.list[this.form.indexName].mappings.properties)
            break
        }
      } catch (e) {
        this.esCols = []
      }
      this.refreshshowMapping()
    },
    async getIndexList() {
      console.log("test")
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

      form['es_connect'] = this.$store.state.baseData.EsConnectID
      this.addLoading = true
      const res = await Transfer(form)
      this.addLoading = false
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
      this.transferLoading = true
      const res = await GetTableColumns({id: this.getSelectTypeObj()['id'], table_name: this.form.selectTable})
      this.transferLoading = false
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
          label: v.Comment == '' ? `${v.Field}【${v.Type}】` : `${v.Field}【${v.Type}】【${v.Comment}】`,
          disabled: false
        }
        this.allCols.push(obj)
      }
      this.refreshShowAutoIncrementId()
      this.form.autoIncrementId = ""
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
      this.openTaskLoading = true

      const res = await LinkSelectOpt()
      if (res.code != 0) {
        this.openTaskLoading = false
        this.$message({
          type: 'error',
          message: res.msg
        })
        return
      }
      if (res.data == null) res.data = []

      if (res.data.length == 0) {
        this.openTaskLoading = false
        this.$message({
          type: 'error',
          message: '请先添加数据源'
        })
        return
      }

      this.linkSelectOpt = res.data

      if (this.linkSelectOpt.length > 0) {
        this.form.selectType = JSON.stringify(this.linkSelectOpt[0])
        await this.getTables()
        await this.GetTableColumns()
      }
      this.refreshShowAutoIncrementId()
      this.openTaskLoading = false
      this.open = true
      this.addLoading = false
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
