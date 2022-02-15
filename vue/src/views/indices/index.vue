<template>
  <div class="app-container">
    <el-card class="box-card">
      <div class="filter-container">
        <el-select
          id="index-health-status"
          v-model="status"
          class="filter-item width150"
          clearable
          filterable
          @change="search"
        >
          <el-option label="索引健康状态" value="" />
          <el-option label="green" value="green" />
          <el-option label="yellow" value="yellow" />
          <el-option label="red" value="red" />
        </el-select>

        <el-tag class="filter-item">请输入关键词</el-tag>

        <el-input id="index-keyword" v-model="input" class="filter-item width300" clearable @input="search" />
        <el-button id="index-search" type="primary" class="filter-item" icon="el-icon-search" @click="search">搜索
        </el-button>
        <el-button-group>

          <el-button
            id="new-index"
            type="success"
            class="filter-item"
            icon="el-icon-plus"
            @click="openSettingDialog('','add')"
          >新建索引
          </el-button>

          <el-button
            id="readOnlyAllowDelete"
            v-loading="readOnlyAllowDeleteLoading"
            type="warning"
            class="filter-item"
            icon="el-icon-sort"
            @click="readOnlyAllowDelete()"
          >
            将节点切换为可读写状态
          </el-button>

          <el-button
            id="flushIndex"
            v-loading="loadingGroup['_flush']"
            type="info"
            class="filter-item"
            icon="el-icon-s-open"
            @click="runCommandByIndex('_flush','')"
          >将所有索引刷新到磁盘
          </el-button>
        </el-button-group>

      </div>
      <div id="patch-operate" class="filter-container">
        <el-button-group>
          <el-button
            id="patchCloseIndex"
            v-loading="loadingGroup['close']"
            type="danger"
            icon="el-icon-circle-close"
            class="filter-item"
            @click="runCommandByIndex('close',selectIndexList.join(','))"
          >关闭
          </el-button>

          <el-button
            id="patchOpenIndex"
            v-loading="loadingGroup['open']"
            type="success"

            icon="el-icon-success"
            class="filter-item"
            @click="runCommandByIndex('open',selectIndexList.join(','))"
          >打开
          </el-button>
          <el-button
            id="patchForcemergeIndex"
            v-loading="loadingGroup['_forcemerge']"
            icon="el-icon-connection"
            class="filter-item"
            @click="runCommandByIndex('_forcemerge',selectIndexList.join(','))"
          >强制合并索引
          </el-button>

          <el-popover
            placement="top-start"
            title="提示"
            width="200"
            trigger="hover"
            content="为了让最新的数据可以立即被搜索到"
          >
            <el-button
              id="patchRefreshIndex"
              slot="reference"
              v-loading="loadingGroup['_refresh']"
              class="filter-item"

              type="primary"
              icon="el-icon-refresh"
              @click="runCommandByIndex('_refresh',selectIndexList.join(','))"
            >刷新索引
            </el-button>
          </el-popover>

          <el-popover
            placement="top-start"
            title="提示"
            width="200"
            trigger="hover"
            content="让数据持久化到磁盘中"
          >
            <el-button
              id="patchFlushIndex"
              slot="reference"
              v-loading="loadingGroup['_flush']"

              type="info"
              icon="el-icon-s-open"
              class="filter-item"
              @click="runCommandByIndex('_flush',selectIndexList.join(','))"
            >将索引刷新到磁盘
            </el-button>
          </el-popover>

          <el-button
            id="patchCacheClear"
            v-loading="loadingGroup['_cache/clear']"
            class="filter-item"

            type="warning"
            icon="el-icon-toilet-paper"
            @click="runCommandByIndex('_cache/clear',selectIndexList.join(','))"
          >清理缓存
          </el-button>

          <el-button
            id="patchDeleteIndex"
            v-loading="loadingGroup['deleteIndex']"
            class="filter-item"
            type="danger"
            icon="el-icon-delete"
            @click="deleteIndex(selectIndexList.join(','),'deleteIndex')"
          >删除索引
          </el-button>
        </el-button-group>
      </div>
      <back-to-top />

      <el-table
        v-loading="connectLoading"
        :data="list"
        @selection-change="selectChange"
      >
        <el-table-column
          type="selection"
          width="55"
        />
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

        <el-table-column align="center" label="索引健康状态" width="100">
          <template slot-scope="scope">
            <el-button v-if="scope.row.health == 'green'" type="success" circle />
            <el-button v-if="scope.row.health == 'yellow'" type="warning" circle />
            <el-button v-if="scope.row.health == 'red'" type="danger" circle />
          </template>
        </el-table-column>
        <el-table-column align="center" label="索引的开启状态" width="100">
          <template slot-scope="scope">
            <el-button
              v-show="scope.row.status == 'open'"
              type="success"
              size="small"
              icon="el-icon-success"
              @click="runCommandByIndex('close',scope.row.index)"
            >开启
            </el-button>
            <el-button
              v-show="scope.row.status == 'close'"
              type="danger"
              size="small"
              icon="el-icon-circle-close"
              @click="runCommandByIndex('open',scope.row.index)"
            >关闭
            </el-button>
          </template>
        </el-table-column>
        <el-table-column align="center" label="索引名称" width="180">
          <template slot-scope="scope">
            {{ scope.row.index }}
          </template>
        </el-table-column>
        <el-table-column align="center" label="索引uuid" width="220">
          <template slot-scope="scope">
            {{ scope.row.uuid }}
          </template>
        </el-table-column>
        <el-table-column align="center" label="索引主分片数" width="80" prop="pri" sortable>
          <template slot-scope="scope">
            {{ scope.row.pri }}
          </template>
        </el-table-column>
        <el-table-column align="center" label="索引副本分片数量" width="80" prop="rep" sortable>
          <template slot-scope="scope">
            {{ scope.row.rep }}
          </template>
        </el-table-column>
        <el-table-column align="center" label="索引文档总数" width="80" prop="docs->count" sortable>
          <template slot-scope="scope">
            {{ bigNumberTransform(scope.row["docs.count"]) }}
          </template>
        </el-table-column>
        <el-table-column align="center" label="索引中删除状态的文档" width="80" prop="docs->deleted" sortable>
          <template slot-scope="scope">
            {{ scope.row["docs.deleted"] }}
          </template>
        </el-table-column>
        <el-table-column align="center" label="主分片+副本分分片的大小" width="120" prop="store->size" sortable>
          <template slot-scope="scope">
            {{ scope.row["store.size"] }}
          </template>
        </el-table-column>
        <el-table-column align="center" label="主分片的大小" width="150" prop="pri->store->size" sortable>
          <template slot-scope="scope">
            {{ scope.row["pri.store.size"] }}
          </template>
        </el-table-column>

        <el-table-column align="center" label="操作" fixed="right" width="400">
          <template slot-scope="scope">

            <el-button-group>
              <!--  <el-button
                  v-if="Object.keys(mappings[scope.row.index].mappings).length == 0"
                  type="warning"
                  size="small"
                  icon="el-icon-circle-plus-outline"
                  @click="openMappingEditDialog(scope.row.index,false)"
                >新增映射结构
                </el-button>-->
              <el-button
                type="primary"
                size="small"
                icon="el-icon-setting"
                @click="openSettingDialog(scope.row.index,'update')"
              >修改配置
              </el-button>

              <el-button
                icon="el-icon-more"
                type="primary"
                size="small"
                @click="openDrawer(scope.row.index)"
              >更多
              </el-button>
            </el-button-group>
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
      <el-drawer
        ref="drawer"
        :title="indexName"
        :before-close="drawerHandleClose"
        :visible.sync="drawerShow"
        direction="rtl"
        custom-class="demo-drawer"
        close-on-press-escape
        destroy-on-close
        size="50%"
      >
        <el-tabs v-model="activeName" class="margin-left-10" @tab-click="changeTab">
          <el-tab-pane label="设置" name="Settings">
            <json-editor
              v-if="activeName == 'Settings'"
              v-model="activeData"
              styles="width: 100%"
              :read="true"
              title="设置"
            />

          </el-tab-pane>
          <el-tab-pane label="映射" name="Mapping">
            <div class="filter-container operate">

              <el-tag type="warning" class="filter-item">切换为其它索引的映射</el-tag>

              <index-select
                class="filter-item"
                :clearable="true"
                placeholder="请选择索引名"
                @change="changeMapToAnotherIndex"
              />
              <el-tag type="primary" class="filter-item">操作</el-tag>
              <el-button
                v-loading="loadingGroup['saveMappinng']"
                class="filter-item"
                size="small"
                type="primary"
                icon="el-icon-check"
                @click="saveMappinng"
              >修改
              </el-button>
              <el-link type="danger">【注意：只能新增映射字段不可修改映射字段类型】</el-link>
            </div>
            <json-editor
              v-if="activeName == 'Mapping'"
              v-model="activeData"
              styles="width: 100%"
              :read="false"
              title="映射"
              @getValue="getMapping"
            />
          </el-tab-pane>
          <el-tab-pane label="Stats" name="Stats">
            <json-editor
              v-if="activeName == 'Stats'"
              v-model="activeData"

              styles="width: 100%;"
              :read="true"
              title="Stats"
            />
          </el-tab-pane>
          <el-tab-pane label="编辑索引配置" name="editSettings">
            <el-form>
              <el-form-item label="编辑索引配置">
                <el-button type="primary" icon="el-icon-edit-outline" @click="submitSettings()">提交</el-button>
                <el-button icon="refresh" @click="resetSettings">重置</el-button>
              </el-form-item>
              <el-form-item>
                <json-editor
                  v-if="activeName == 'editSettings'"
                  v-model="activeData"

                  :point-out="pointOut"
                  styles="width: 100%;"
                  :read="false"
                  title="编辑配置"
                  @getValue="getSettings"
                />
              </el-form-item>
            </el-form>
          </el-tab-pane>
          <el-tab-pane label="修改别名" name="alias">
            <alias v-if="activeName == 'alias'" :index-name="indexName" />
          </el-tab-pane>
          <div class="filter-container operate">
            <el-tag class="filter-item">操作</el-tag>
            <el-button
              v-loading="loadingGroup['close']"
              type="danger"
              size="small"
              icon="el-icon-circle-close"
              class="filter-item"
              @click="runCommandByIndex('close',indexName)"
            >关闭
            </el-button>

            <el-button
              v-loading="loadingGroup['open']"
              type="success"
              size="small"
              icon="el-icon-success"
              class="filter-item"
              @click="runCommandByIndex('open',indexName)"
            >打开
            </el-button>
            <el-button
              v-loading="loadingGroup['_forcemerge']"
              size="small"
              icon="el-icon-connection"
              class="filter-item"
              @click="runCommandByIndex('_forcemerge',indexName)"
            >强制合并索引
            </el-button>
            <el-link type="danger">【forcemerge操作,手动释放磁盘空间】</el-link>
            <el-popover
              placement="top-start"
              title="提示"
              width="200"
              trigger="hover"
              content="为了让最新的数据可以立即被搜索到"
            >
              <el-button
                slot="reference"
                v-loading="loadingGroup['_refresh']"
                class="filter-item"
                size="small"
                type="primary"
                icon="el-icon-refresh"
                @click="runCommandByIndex('_refresh',indexName)"
              >刷新索引
              </el-button>
            </el-popover>

            <el-popover
              placement="top-start"
              title="提示"
              width="200"
              trigger="hover"
              content="让数据持久化到磁盘中"
            >
              <el-button

                slot="reference"
                v-loading="loadingGroup['_flush']"
                size="small"
                type="info"
                icon="el-icon-s-open"
                class="filter-item"
                @click="runCommandByIndex('_flush',indexName)"
              >将索引刷新到磁盘
              </el-button>
            </el-popover>

            <el-button
              v-loading="loadingGroup['_cache/clear']"
              class="filter-item"
              size="small"
              type="warning"
              icon="el-icon-toilet-paper"
              @click="runCommandByIndex('_cache/clear',indexName)"
            >清理缓存
            </el-button>

            <el-button
              v-loading="loadingGroup['deleteIndex']"
              class="filter-item"
              type="danger"
              size="small"
              icon="el-icon-delete"
              @click="deleteIndex(indexName,'deleteIndex')"
            >删除索引
            </el-button>

          </div>
        </el-tabs>

      </el-drawer>
      <settings
        v-if="openSettings"
        :index-name="indexName"
        :settings-type="settingsType"
        :finished="search"
        :open="openSettings"
        @close="closeSettings"
      />
      <mappings
        v-if="openMappings"
        :index-name="indexName"
        :mappings="mappingInfo"
        :title="mappingTitle"
        :open="openMappings"
        @close="closeMappings"
      />

    </el-card>
  </div>
</template>

<script>
import { Finish, IsFinish } from '@/api/guid'
import Driver from 'driver.js' // import driver.js
import 'driver.js/dist/driver.min.css' // import driver.js css
import steps from '@/views/indices/guide'
import { clone } from '@/utils/index'
import { filterData } from '@/utils/table'
import { CatAction, OptimizeAction, RecoverCanWrite } from '@/api/es'
import { bigNumberTransform } from '@/utils/format'
import { CreateAction, DeleteAction, GetSettingsAction, GetSettingsInfoAction, StatsAction } from '@/api/es-index'
import { esSettingsWords } from '@/utils/base-data'
import { ListAction, UpdateMappingAction } from '@/api/es-map'

export default {
  name: 'CatIndices',

  components: {
    'Settings': () => import('@/views/indices/components/settings'),
    'Mappings': () => import('@/views/indices/components/mapping'),
    'BackToTop': () => import('@/components/BackToTop/index'),
    'JsonEditor': () => import('@/components/JsonEditor/index'),
    'Alias': () => import('@/views/indices/components/alias'),
    'IndexSelect': () => import('@/components/index/select')
  },

  data() {
    return {
      modName: '索引管理',
      aliasList: [],
      pointOut: esSettingsWords,
      settings: {},
      readOnlyAllowDeleteLoading: false,
      loadingGroup: {
        'close': false,
        'open': false,
        '_forcemerge': false,
        '_refresh': false,
        '_flush': false,
        '_cache/clear': false,
        'deleteIndex': false,
        '_all/_flush': false,
        'saveMappinng': false
      },

      forceMergeLoading: false,
      tabLoading: false,
      activeData: '{}',
      activeName: 'Settings',
      drawerShow: false,
      settingsType: 'add',
      mappingTitle: '',
      indexName: '',
      openSettings: false,
      openMappings: false,
      total: 0,
      connectLoading: false,
      page: 1,
      limit: 10,
      pageshow: true,
      list: [],
      input: '',
      status: '',
      mappingInfo: {},
      mappings: {},
      selectIndexList: []
    }
  },
  destroyed() {
    sessionStorage.setItem('CatIndices', this.input)
  },
  mounted() {
    const input = sessionStorage.getItem('CatIndices')
    if (input != null) {
      this.input = input
    }
    this.startGuid()
    this.GetMapAction()
    this.searchData()
  },
  methods: {
    async finishGuid() {
      const { data, code, msg } = await Finish({ 'guid_name': this.modName })
    },
    async startGuid() {
      const { data, code, msg } = await IsFinish({ 'guid_name': this.modName })

      if (!data) {
        console.log('开始新手引导')
        this.driver = new Driver({
          className: 'scoped-class',
          animate: true,
          opacity: 0.75,
          padding: 10,
          allowClose: false,
          overlayClickNext: false,
          doneBtnText: '完成',
          closeBtnText: '关闭',
          nextBtnText: '下一步',
          prevBtnText: '上一步',
          onNext: (Element) => {
            // console.log('Element', Element)
          }
        })
        setTimeout(() => {
          this.driver.defineSteps(steps)
          this.driver.start()
          this.finishGuid()
        }, 500)
      }
    },
    selectChange(row) {
      this.selectIndexList = []
      for (const v of row) {
        this.selectIndexList.push(v.index)
      }
    },
    async changeMapToAnotherIndex(indexName) {
      if (indexName == '') {
        indexName = this.indexName
      }
      const input = {}
      input['es_connect'] = this.$store.state.baseData.EsConnectID
      input['index_name'] = indexName

      const res = await ListAction(input)

      if (res.code == 0) {
        this.activeData = JSON.stringify(res.data[indexName].mappings, null, '\t')
      }
    },
    getMapping(v) {
      this.activeData = v
    },
    async saveMappinng() {
      let activeData = clone(this.activeData)
      try {
        activeData = JSON.parse(activeData)
      } catch (e) {
        this.$message({
          type: 'error',
          message: 'JSON格式不正确'
        })
        return
      }
      const activeDataKeys = Object.keys(activeData)
      if (activeDataKeys.length == 0) {
        this.$message({
          type: 'error',
          message: '请按格式写type名字'
        })
        return
      }

      const input = {}
      input['es_connect'] = this.$store.state.baseData.EsConnectID
      input['index_name'] = this.indexName
      input['type_name'] = activeDataKeys[0]
      input['properties'] = activeData[activeDataKeys[0]]
      this.loadingGroup['saveMappinng'] = true
      const { data, code, msg } = await UpdateMappingAction(input)
      this.loadingGroup['saveMappinng'] = false
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
      console.log('saveMappinng')
    },
    openMappingEditDialog(indexName, haveMapping) {
      if (haveMapping) {
        this.mappingInfo = this.mappings[indexName].mappings
        this.mappingTitle = '新增字段'
      } else {
        this.mappingTitle = '新增映射结构'
      }
      this.indexName = indexName

      this.openMappings = true
    },
    async GetMapAction() {
      const input = {}
      input['es_connect'] = this.$store.state.baseData.EsConnectID
      const { code, data, msg } = await ListAction(input)
      if (code == 0) {
        this.mappings = data
      } else {
        this.$message({
          type: 'error',
          message: msg
        })
        return
      }
    },
    submitSettings() {
      const input = {}
      input['es_connect'] = this.$store.state.baseData.EsConnectID
      try {
        input['settings'] = JSON.parse(this.activeData)
      } catch (e) {
        console.log(e)
        this.$message({
          type: 'error',
          message: 'JSON 解析异常'
        })
        return
      }
      input['index_name'] = this.indexName
      input['types'] = 'update'
      const loading = this.$loading({
        lock: true,
        text: 'Loading',
        spinner: 'el-icon-loading',
        background: 'rgba(0, 0, 0, 0.7)'
      })
      CreateAction(input).then(res => {
        if (res.code == 0 || res.code == 200) {
          this.$message({
            type: 'success',
            message: res.msg
          })
          this.search()
        } else {
          this.$message({
            type: 'error',
            message: res.msg
          })
        }
        loading.close()
      }).catch(err => {
        loading.close()
      })
    },
    getSettings(value) {
      this.activeData = value
    },
    resetSettings() {
      this.changeTab()
    },
    runCommandByIndex(command, indexName) {
      const input = {}
      input['es_connect'] = this.$store.state.baseData.EsConnectID
      input['index_name'] = indexName
      input['command'] = command

      this.loadingGroup[command] = true
      OptimizeAction(input).then(res => {
        if (res.code == 0 || res.code == 200) {
          this.$message({
            type: 'success',
            message: res.msg
          })
          this.search()
        } else {
          this.$message({
            type: 'error',
            message: res.msg
          })
        }
        this.loadingGroup[command] = false
      })
        .catch(err => {
          this.loadingGroup[command] = false
        })
    },
    async changeTab() {
      const input = {}
      let res = {}
      input['es_connect'] = this.$store.state.baseData.EsConnectID
      input['index_name'] = this.indexName
      switch (this.activeName) {
        case 'Settings':
          res = await GetSettingsInfoAction(input)

          if (res.code == 0) {
            this.activeData = JSON.stringify(res.data, null, '\t')
          }
          return

        case 'Mapping':
          res = await ListAction(input)

          if (res.code == 0) {
            this.activeData = JSON.stringify(res.data[this.indexName].mappings, null, '\t')
            if (Object.keys(res.data[this.indexName].mappings).length == 0) {
              this.$message({
                type: 'error',
                message: '您还没有设置映射结构'
              })
            }
          }
          return
        case 'Stats':

          res = await StatsAction(input)
          if (res.code == 0) {
            this.activeData = JSON.stringify(res.data, null, '\t')
          }
          return

        case 'editSettings':

          const { data } = await GetSettingsAction(input)

          const deleteKeyArr = [
            'creation_date', 'version', 'provided_name', 'uuid', 'format', 'number_of_shards'
          ]

          for (const key of deleteKeyArr) {
            if (data['index'].hasOwnProperty(key)) {
              delete data['index'][key]
            }
          }

          this.activeData = JSON.stringify(data, null, '\t')
          return
        case 'alias':
          return
        default:
          this.activeData = '{}'
          return
      }
    },
    openDrawer(indexName) {
      this.indexName = indexName
      this.changeTab()
      this.drawerShow = true
    },
    drawerHandleClose(done) {
      this.indexName = ''
      done()
    },
    bigNumberTransform(value) {
      return bigNumberTransform(value)
    },
    openSettingDialog(indexName, settingsType) {
      this.indexName = indexName
      this.settingsType = settingsType
      this.openSettings = true
    },
    readOnlyAllowDelete() {
      const input = {}
      input['es_connect'] = this.$store.state.baseData.EsConnectID
      this.readOnlyAllowDeleteLoading = true
      RecoverCanWrite(input).then(res => {
        if (res.code == 0 || res.code == 200) {
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
        this.readOnlyAllowDeleteLoading = false
      }).catch(err => {
        this.readOnlyAllowDeleteLoading = false
      })
    },
    deleteIndex(indexName, loadingType) {
      this.$confirm('确定删除该索引吗?', '警告', {
        confirmButtonText: '确认',
        cancelButtonText: '取消',
        type: 'warning'
      })
        .then(async() => {
          const input = {}
          input['es_connect'] = this.$store.state.baseData.EsConnectID
          input['index_name'] = indexName
          this.loadingGroup[loadingType] = true
          DeleteAction(input).then(res => {
            if (res.code == 0 || res.code == 200) {
              this.$message({
                type: 'success',
                message: indexName + '已删除'
              })
              this.searchData()
              this.loadingGroup[loadingType] = false
            } else {
              this.$message({
                type: 'error',
                message: res.msg
              })
              this.loadingGroup[loadingType] = false
            }
          }).catch(err => {

          })
        })
        .catch(err => {
          console.error(err)
        })
    },
    closeSettings() {
      this.indexName = ''
      this.settingsType = 'add'
      this.openSettings = false
    },
    closeMappings() {
      this.indexName = ''
      this.mappingTitle = 'add'
      this.openMappings = false
      this.mappingInfo = {}
    },
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
        cat: 'CatIndices',
        es_connect: this.$store.state.baseData.EsConnectID
      }
      CatAction(form).then(res => {
        if (res.code == 0) {
          const list = res.data

          for (const index in list) {
            const obj = list[index]
            // 把 . 转成 ->
            for (const key in obj) {
              let value = parseInt(obj[key])
              if (isNaN(value)) {
                value = obj[key]
              }
              list[index][key.split('.').join('->')] = value
            }
          }

          let tmpList = []
          if (this.status.trim() != '') {
            for (const v of list) {
              if (v['health'] == this.status.trim()) {
                tmpList.push(v)
              }
            }
          } else {
            tmpList = list
          }
          tmpList = filterData(tmpList, this.input.trim())
          this.list = tmpList.filter((item, index) =>
            index < this.page * this.limit && index >= this.limit * (this.page - 1)
          )
          this.total = tmpList.length
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
  .operate {

  }

  .aliasName {
    width: 400px;
  }

  .margin-left-10 {
    margin-left: 10px
  }

  .width300 {
    width: 300px;
  }

  .width150 {
    width: 150px;
  }

  /deep/ :focus {
    outline: 0;
  }
</style>
