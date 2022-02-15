<template>
  <div>
    <el-dialog :close-on-click-modal="false" width="95%" :visible.sync="dialogVisible" title="查询结果" @close="close">
      <div class="filter-container">
        <el-tag class="filter-item">请输入关键词</el-tag>
        <el-input v-model="input" class="filter-item" style="width: 300px" clearable @input="search" />
        <el-button type="success" class="filter-item" @click="search">搜索</el-button>
        <el-button v-if="ISDoc" type="primary" class="filter-item" @click.native="openAddDialog = true">添加文档</el-button>
      </div>

      <el-table
        v-if="showTable"
        :data="getTableData"
        use-virtual
        :row-height="30"
      >
        <el-table-column
          label="ID"
          align="center"
          fixed
          min-width="100"
        >
          <template slot-scope="scope">
            <template v-if="ISDoc">
              {{ scope.row["_id"] }}
            </template>
            <template v-else>
              {{ scope.$index+1 }}
            </template>
          </template>
        </el-table-column>
        <el-table-column
          v-for="(val,key) in tableHeader"
          v-if="val != 'xwl_index' && val != '_id'"
          :key="key"
          :prop="val"
          :sortable="true"
          align="center"
          :label="val"
        >
          <template slot-scope="scope">
            <el-popover
              v-if=" scope.row[val] !=undefined"
                          placement="top-start"
                          trigger="hover"
            >
              <div>{{ scope.row[val].toString() }}</div>
               <span v-if="scope.row[val].toString().length>=20"  slot="reference">{{ scope.row[val].toString().substr(0, 20) + "..." }}
              </span>
               <span v-else slot="reference">{{ scope.row[val].toString() }}
              </span>
            </el-popover>
          </template>
        </el-table-column>
        <el-table-column align="center" label="操作" fixed="right" width="200">
          <template slot-scope="scope">
            <el-button-group>
              <el-button
                v-if="ISDoc"
                type="primary"
                size="small"
                icon="el-icon-edit"
                @click="look(scope.$index)"
              >
                编辑
              </el-button>
              <el-button v-else type="success" size="small" icon="el-icon-search" @click="look(scope.$index)">
                查看
              </el-button>
              <el-button
                v-if="ISDoc"
                type="danger"
                size="small"
                icon="el-icon-delete"
                @click="deleteByID(scope.row,scope.$index)"
              >删除
              </el-button>
            </el-button-group>
          </template>
        </el-table-column>

      </el-table>

    </el-dialog>
    <el-drawer
      ref="drawer"
      title="详细数据"
      :before-close="drawerHandleClose"
      :visible.sync="drawerShow"

      direction="rtl"
      close-on-press-escape
      destroy-on-close
      size="50%"
    >
      <div v-if="ISDoc" class="filter-container">
        <el-tag class="filter-item">操作</el-tag>
        <el-button type="primary" class="filter-item" @click="updateByID">修改</el-button>
      </div>
      <json-editor
        v-if="isArray"
        v-model="JSON.stringify(jsonData[index],null, '\t')"
        height="900"
        class="res-body"
        styles="width: 100%"
        :read="true"
        title="详细数据"
      />

      <json-editor
        v-if="!isArray"
        v-model="JSON.stringify(jsonData['hits']['hits'][index],null, '\t')"
        height="900"
        class="res-body"
        styles="width: 100%"
        :read="!ISDoc"
        title="详细数据"
        @getValue="getEditDoc"
      />

    </el-drawer>
    <el-drawer
      ref="drawer"
      title="新增文档"
      :before-close="drawerHandleClose"
      :visible.sync="openAddDialog"

      direction="rtl"
      close-on-press-escape
      destroy-on-close
      size="50%"
    >
      <div class="filter-container">
        <el-tag class="filter-item">操作</el-tag>
        <el-button type="primary" icon="el-icon-edit" class="filter-item" @click="add">提交</el-button>
      </div>

      <json-editor
        v-if="openAddDialog"
        v-model="JSON.stringify(properties,null, '\t')"
        height="900"
        class="res-body"
        styles="width: 100%"
        title="新增文档"
        @getValue="getNewDoc"
      />

    </el-drawer>
  </div>
</template>

<script>
import { filterData } from '@/utils/table'
import { clone } from '@/utils/index'
import { DeleteRowByIDAction, InsertAction, UpdateByIDAction } from '@/api/es-doc'
import { ListAction } from '@/api/es-map'

export default {
  name: 'ResTable',
  components: {
    'JsonEditor': () => import('@/components/JsonEditor/index')
  },
  props: {

    dialogVisible: {
      type: Boolean,
      default: false
    },
    jsonData: {
      type: Array,
      default: []
    },
    searchPath: {
      type: String,
      default: ''
    }
  },
  data() {
    return {
      openAddDialog: false,
      showTable: true,
      drawerShow: false,
      tableData: [],
      index: 0,
      isArray: false,
      tableHeader: [],
      ISDoc: false,
      input: '',
      properties: {},
      newDoc: {}
    }
  },
  computed: {
    getTableData() {
      if (this.input == '') {
        return this.tableData
      } else {
        return filterData(this.tableData, this.input.trim())
      }
    }
  },
  created() {
    this.initTableData()
  },
  methods: {
    getNewDoc(doc) {
      try {
        this.newDoc = JSON.parse(doc)
      } catch (e) {

      }
    },
    async add() {
      const editData = this.jsonData['hits']['hits'][this.index]
      const doc = this.newDoc

      const input = {}
      input['es_connect'] = this.$store.state.baseData.EsConnectID
      input['index'] = editData['_index']
      input['type_name'] = editData['_type']
      input['json'] = doc

      const res = await InsertAction(input)
      if (res.code == 0) {
        this.$message({
          type: 'success',
          message: res.msg + '(_id为:' + res.data._id + ')',
          duration: 20000
        })
      } else {
        this.$message({
          type: 'error',
          message: res.msg
        })
      }
    },
    search() {
      this.showTable = false
      this.$nextTick(() => {
        this.showTable = true
      })
    },
    async initTableData() {
      const resData = clone(this.jsonData)
      if (Array.isArray(resData)) {
        if (resData.length > 500) {
          this.$message({
            type: 'error',
            message: '请减少数据条数'
          })
          this.$emit('close', false)
          return
        }
        this.tableData = resData

        for (const index in this.tableData) {
          const obj = this.tableData[index]

          for (const key in obj) {
            let value = this.strToNum(obj[key])
            if (value == false) {
              value = obj[key]
            }
            if (key.indexOf('.') != -1) {
              this.tableData[index][key.split('.').join('->')] = value
              delete this.tableData[index][key]
            } else {
              this.tableData[index][key] = value
            }

            this.tableData[index]['xwl_index'] = index
          }
        }

        this.tableHeader = Object.keys(this.tableData[0])
        this.isArray = true
      } else {
        if (resData.hasOwnProperty('hits')) {
          if (resData['hits']['hits'].length > 0) {
            this.ISDoc = true
            if (resData['hits']['hits'].length > 500) {
              this.$message({
                type: 'error',
                message: '请减少查詢的数据条数'
              })
              this.$emit('close', false)
              return
            }

            const sourceArr = []

            for (const index in resData['hits']['hits']) {
              const _source = resData['hits']['hits'][index]['_source']
              _source['_id'] = resData['hits']['hits'][index]['_id']
              _source['_score'] = resData['hits']['hits'][index]['_score']
              _source['xwl_index'] = index
              sourceArr.push(_source)
            }

            this.isArray = false
            this.tableData = sourceArr

            for (const index in this.tableData) {
              for (const key in this.tableData[index]) {
                if (typeof this.tableData[index][key] === 'object' || typeof this.tableData[index][key] === 'array') {
                  this.tableData[index][key] = JSON.stringify(this.tableData[index][key])
                }
              }
            }
          }
        }

        const input = {}
        input['es_connect'] = this.$store.state.baseData.EsConnectID
        input['index_name'] = this.jsonData['hits']['hits'][0]['_index']

        const { data, code, msg } = await ListAction(input)

        if (code == 0) {
          const mappings = Object.keys(data[input['index_name']].mappings)
          const tableHeader = Object.keys(data[input['index_name']].mappings[mappings[0]].properties)
          this.properties = data[input['index_name']].mappings[mappings[0]].properties
          const tmpTableHeader = []
          tmpTableHeader.push('_id')
          for (const i in tableHeader) {
            if (tableHeader[i] != '_id') {
              tmpTableHeader.push(tableHeader[i])
            }
          }
          this.tableHeader = tmpTableHeader
        } else {
          this.$message({
            type: 'error',
            message: msg
          })
        }
      }
    },
    getEditDoc(v) {
      try {
        const editDoc = JSON.parse(v)
        this.jsonData['hits']['hits'][this.index] = editDoc
      } catch (e) {

      }
    },
    async updateByID() {
      const editData = this.jsonData['hits']['hits'][this.index]
      const doc = editData['_source']

      const input = {}
      input['es_connect'] = this.$store.state.baseData.EsConnectID
      input['index'] = editData['_index']
      input['type_name'] = editData['_type']
      input['id'] = editData['_id']
      input['json'] = doc

      const res = await UpdateByIDAction(input)
      if (res.code == 0) {
        const _source = clone(this.jsonData['hits']['hits'][this.index]['_source'])
        _source['_id'] = this.jsonData['hits']['hits'][this.index]['_id']
        _source['_score'] = this.jsonData['hits']['hits'][this.index]['_score']
        _source['xwl_index'] = this.index
        for (const key in _source) {
          if (typeof _source[key] === 'object' || _source[key] === 'array') {
            _source[key] = JSON.stringify(_source[key])
          }
        }
        this.tableData[this.index] = _source
        this.showTable = false
        this.$nextTick(() => {
          this.showTable = true
        })
        this.$message({
          type: 'success',
          message: res.msg
        })
      } else {
        this.$message({
          type: 'error',
          message: res.msg
        })
      }
    },
    async deleteByID(row, index) {
      this.$confirm('确定删除该条文档吗?', '警告', {
        confirmButtonText: '确认',
        cancelButtonText: '取消',
        type: 'warning'
      })
        .then(async() => {
          const input = {}
          input['es_connect'] = this.$store.state.baseData.EsConnectID
          input['index_name'] = this.jsonData['hits']['hits'][0]['_index']
          input['type'] = this.jsonData['hits']['hits'][0]['_type']
          input['id'] = row['_id']

          console.log(input)
          const res = await DeleteRowByIDAction(input)
          if (res.code == 0) {
            this.tableData.splice(index, 1)
            this.$message({
              type: 'success',
              message: res.msg
            })
          } else {
            this.$message({
              type: 'error',
              message: res.msg
            })
          }
        })
        .catch(err => {
          console.error(err)
        })
    },
    strToNum(str) {
      var convertNum = Number(str) // 将字符串强制转换为数字
      if (str === '') { // 排除空字符串
        return false
      } else {
        if (str.includes(' ')) { // 排除空格
          return false
        } else {
          if (isNaN(convertNum)) {
            return false
          } else {
            return convertNum
          }
        }
      }
    },
    look(index) {
      this.index = index
      this.drawerShow = true
    },
    drawerHandleClose(done) {
      done()
    },
    close() {
      this.$emit('close', false)
    }

  }
}
</script>

<style scoped>

</style>
