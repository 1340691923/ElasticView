<template>
  <div class="app-container">
    <div class="filter-container">
      <el-select
        id="index-health-status"
        v-model="status"
        class="filter-item width150"
        clearable
        filterable
        @change="search"
      >
        <el-option :label="$t('索引健康状态')" value=""/>
        <el-option label="green" value="green"/>
        <el-option label="yellow" value="yellow"/>
        <el-option label="red" value="red"/>
      </el-select>

      <el-tag class="filter-item">{{ $t('请输入关键词') }}</el-tag>

      <el-autocomplete
        class="filter-item width300"
        id="index-keyword"
        clearable
        :fetch-suggestions="querySearch"

        v-model="input"
        :placeholder="$t('请输入索引名')"
      >
        <i
          class="el-icon-edit el-input__icon"
          slot="suffix"
        >
        </i>
        <template slot-scope="{ item }">
          <span>{{ item.value }}</span>
        </template>

      </el-autocomplete>

      <el-button
        size="mini"
        id="index-search" type="primary" class="filter-item" icon="el-icon-search" @click="search">{{ $t('搜索') }}
      </el-button>
      <el-button
        size="mini"
        id="index-export" type="primary" class="filter-item" icon="el-icon-export" @click="exportXlsx">{{ $t('导出') }}
      </el-button>
      <el-button-group>

        <el-button
          size="mini"

          id="new-index"
          type="success"
          class="filter-item"
          icon="el-icon-plus"
          @click="openSettingDialog('','add')"
        >{{ $t('新建索引') }}
        </el-button>

        <el-button
          size="mini"

          :disabled="readOnlyAllowDeleteLoading"
          id="readOnlyAllowDelete"
          v-loading="readOnlyAllowDeleteLoading"
          type="warning"
          class="filter-item"
          icon="el-icon-sort"
          @click="readOnlyAllowDelete()"
        >
          {{ $t('将节点切换为可读写状态') }}
        </el-button>

        <el-button
          size="mini"

          :disabled="loadingGroup['_flush']"
          id="flushIndex"
          v-loading="loadingGroup['_flush']"
          type="info"
          class="filter-item"
          icon="el-icon-s-open"
          @click="runCommandByIndex('_flush','')"
        >{{ $t('将所有索引刷新到磁盘') }}
        </el-button>
      </el-button-group>

    </div>
    <div id="patch-operate" class="filter-container">
      <el-button-group>
        <el-button
          size="mini"

          :disabled="loadingGroup['close']"
          id="patchCloseIndex"
          v-loading="loadingGroup['close']"
          type="danger"
          icon="el-icon-circle-close"
          class="filter-item"
          @click="runCommandByIndex('close',selectIndexList.join(','))"
        >{{ $t('关闭') }}
        </el-button>

        <el-button
          size="mini"

          id="patchOpenIndex"
          v-loading="loadingGroup['open']"
          :disabled="loadingGroup['open']"
          type="success"

          icon="el-icon-success"
          class="filter-item"
          @click="runCommandByIndex('open',selectIndexList.join(','))"
        >{{ $t('打开') }}
        </el-button>
        <el-button
          size="mini"

          id="patchForcemergeIndex"
          v-loading="loadingGroup['_forcemerge']"
          :disabled="loadingGroup['_forcemerge']"
          icon="el-icon-connection"
          class="filter-item"
          @click="runCommandByIndex('_forcemerge',selectIndexList.join(','))"
        >{{ $t('强制合并索引') }}
        </el-button>

        <el-popover
          placement="top-start"
          :title="$t('提示')"
          width="200"
          trigger="hover"
          :content="$t('为了让最新的数据可以立即被搜索到')"
        >
          <el-button
            size="mini"

            id="patchRefreshIndex"
            slot="reference"
            v-loading="loadingGroup['_refresh']"
            :disabled="loadingGroup['_refresh']"
            class="filter-item"

            type="primary"
            icon="el-icon-refresh"
            @click="runCommandByIndex('_refresh',selectIndexList.join(','))"
          >{{ $t('刷新索引') }}
          </el-button>
        </el-popover>

        <el-popover
          placement="top-start"
          :title="$t('提示')"
          width="200"
          trigger="hover"
          :content="$t('让数据持久化到磁盘中')"
        >
          <el-button
            size="mini"

            id="patchFlushIndex"
            slot="reference"
            v-loading="loadingGroup['_flush']"
            :disabled="loadingGroup['_flush']"
            type="info"
            icon="el-icon-s-open"
            class="filter-item"
            @click="runCommandByIndex('_flush',selectIndexList.join(','))"
          >{{ $t('将索引刷新到磁盘') }}
          </el-button>
        </el-popover>

        <el-button
          size="mini"

          id="patchCacheClear"
          v-loading="loadingGroup['_cache/clear']"
          :disabled="loadingGroup['_cache/clear']"
          class="filter-item"

          type="warning"
          icon="el-icon-toilet-paper"
          @click="runCommandByIndex('_cache/clear',selectIndexList.join(','))"
        >{{ $t('清理缓存') }}
        </el-button>

        <el-button
          size="mini"

          id="patchDeleteIndex"
          v-loading="loadingGroup['deleteIndex']"
          :disabled="loadingGroup['deleteIndex']"
          class="filter-item"
          type="danger"
          icon="el-icon-delete"
          @click="deleteIndex(selectIndexList.join(','),'deleteIndex')"
        >{{ $t('删除索引') }}
        </el-button>
        <el-button
          size="mini"

          id="patchEmptyIndex"
          v-loading="loadingGroup['empty']"
          :disabled="loadingGroup['empty']"
          class="filter-item"
          type="danger"
          icon="el-icon-delete"
          @click="runCommandByIndex('empty',selectIndexList.join(','))"
        >{{ $t('清空索引') }}
        </el-button>
      </el-button-group>
    </div>
    <back-to-top/>

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
        :label="$t('序号')"
        align="center"
        fixed
        width="50"
      >
        <template slot-scope="scope">
          {{ scope.$index + 1 }}
        </template>
      </el-table-column>

      <el-table-column align="center" :label="$t('索引健康状态')" width="100">
        <template slot-scope="scope">
          <el-button
            size="mini"
            v-if="scope.row.health == 'green'" type="success" circle/>
          <el-button
            size="mini"
            v-if="scope.row.health == 'yellow'" type="warning" circle/>
          <el-button
            size="mini"
            v-if="scope.row.health == 'red'" type="danger" circle/>
        </template>
      </el-table-column>
      <el-table-column align="center" :label="$t('索引的开启状态')" width="100">
        <template slot-scope="scope">
          <el-button

            v-show="scope.row.status == 'open'"
            type="success"
            size="mini"
            icon="el-icon-success"
            @click="runCommandByIndex('close',scope.row.index)"
          >{{ $t('开启') }}
          </el-button>
          <el-button

            v-show="scope.row.status == 'close'"
            type="danger"
            size="mini"
            icon="el-icon-circle-close"
            @click="runCommandByIndex('open',scope.row.index)"
          >{{ $t('关闭') }}
          </el-button>
        </template>
      </el-table-column>
      <el-table-column align="center" :label="$t('索引名称')" width="180">
        <template slot-scope="scope">
          {{ scope.row.index }}
        </template>
      </el-table-column>
      <el-table-column align="center" :label="$t('索引uuid')" width="220">
        <template slot-scope="scope">
          {{ scope.row.uuid }}
        </template>
      </el-table-column>
      <el-table-column align="center" :label="$t('索引主分片数')" width="80" prop="pri" sortable>
        <template slot-scope="scope">
          {{ scope.row.pri }}
        </template>
      </el-table-column>
      <el-table-column align="center" :label="$t('索引副本分片数量')" width="80" prop="rep" sortable>
        <template slot-scope="scope">
          {{ scope.row.rep }}
        </template>
      </el-table-column>
      <el-table-column align="center" :label="$t('索引文档总数')" width="80" prop="docs->count" sortable>
        <template slot-scope="scope">
          {{ bigNumberTransform(scope.row["docs.count"]) }}
        </template>
      </el-table-column>
      <el-table-column align="center" :label="$t('索引中删除状态的文档')" width="80" prop="docs->deleted" sortable>
        <template slot-scope="scope">
          {{ scope.row["docs.deleted"] }}
        </template>
      </el-table-column>
      <el-table-column align="center" :label="$t('主分片+副本分分片的大小')" width="120" prop="store->size" sortable>
        <template slot-scope="scope">
          {{ scope.row["store.size"] }}
        </template>
      </el-table-column>
      <el-table-column align="center" :label="$t('主分片的大小')" width="150" prop="pri->store->size" sortable>
        <template slot-scope="scope">
          {{ scope.row["pri.store.size"] }}
        </template>
      </el-table-column>

      <el-table-column align="center" :label="$t('操作')" fixed="right" width="380">
        <template slot-scope="scope">

          <el-button-group>
            <!--  <el-button
size="mini"

                v-if="Object.keys(mappings[scope.row.index].mappings).length == 0"
                type="warning"
                size="mini"
                icon="el-icon-circle-plus-outline"
                @click="openMappingEditDialog(scope.row.index,false)"
              >新增映射结构
              </el-button>-->
            <el-button

              type="primary"
              size="mini"
              icon="el-icon-setting"
              @click="openSettingDialog(scope.row.index,'update')"
            >{{ $t('修改配置') }}
            </el-button>
            <el-button


              type="primary"
              size="mini"
              icon="el-icon-circle-plus-outline"
              @click="openMappingEditDialog(scope.row.index,false)"
            >{{ $t('修改映射') }}
            </el-button>
            <el-button


              icon="el-icon-more"
              type="primary"
              size="mini"
              @click="openDrawer(scope.row.index)"
            >{{ $t('更多') }}
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
        <el-tab-pane :label="$t('设置')" name="Settings">
          <json-editor
            v-if="activeName == 'Settings'"
            v-model="activeData"
            styles="width: 100%"
            :read="true"
            :title="$t('设置')"
          />

        </el-tab-pane>
        <el-tab-pane :label="$t('映射')" name="Mapping">
          <div class="filter-container operate">

            <el-tag type="warning" class="filter-item">{{ $t('切换为其它索引的映射') }}</el-tag>

            <index-select
              class="filter-item"
              :clearable="true"
              :placeholder="$t('请选择索引名')"
              @change="changeMapToAnotherIndex"
            />
            <el-tag type="primary" class="filter-item">{{ $t('操作') }}</el-tag>
            <el-button


              :disabled="loadingGroup['saveMappinng']"
              v-loading="loadingGroup['saveMappinng']"
              class="filter-item"
              size="mini"
              type="primary"
              icon="el-icon-check"
              @click="saveMappinng"
            >{{ $t('修改') }}
            </el-button>
            <el-link type="danger">{{ $t('【注意：只能新增映射字段不可修改映射字段类型】') }}</el-link>
          </div>
          <json-editor
            v-if="activeName == 'Mapping'"
            v-model="activeData"
            styles="width: 100%"
            :read="false"
            :title="$t('映射')"
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
        <el-tab-pane :label="$t('编辑索引配置')" name="editSettings">
          <el-form>
            <el-form-item :label="$t('编辑索引配置')">
              <el-button
                size="mini"
                type="primary" icon="el-icon-edit-outline" @click="submitSettings()">{{ $t('提交') }}
              </el-button>
              <el-button
                size="mini"
                icon="refresh" @click="resetSettings">{{ $t('重置') }}
              </el-button>
            </el-form-item>
            <el-form-item>
              <json-editor
                v-if="activeName == 'editSettings'"
                v-model="activeData"
                :point-out="pointOut"
                styles="width: 100%;"
                :read="false"
                :title="$t('编辑配置')"
                @getValue="getSettings"
              />
            </el-form-item>
          </el-form>
        </el-tab-pane>
        <el-tab-pane :label="$t('修改别名')" name="alias">
          <alias v-if="activeName == 'alias'" :index-name="indexName"/>
        </el-tab-pane>
        <div class="filter-container operate">
          <el-tag class="filter-item">{{ $t('操作') }}</el-tag>
          <el-button


            :disabled="loadingGroup['close']"
            v-loading="loadingGroup['close']"
            type="danger"
            size="mini"
            icon="el-icon-circle-close"
            class="filter-item"
            @click="runCommandByIndex('close',indexName)"
          >{{ $t('关闭') }}
          </el-button>

          <el-button

            :disabled="loadingGroup['open']"
            v-loading="loadingGroup['open']"
            type="success"
            size="mini"
            icon="el-icon-success"
            class="filter-item"
            @click="runCommandByIndex('open',indexName)"
          >{{ $t('打开') }}
          </el-button>
          <el-button

            :disabled="loadingGroup['_forcemerge']"
            v-loading="loadingGroup['_forcemerge']"
            size="mini"
            icon="el-icon-connection"
            class="filter-item"
            @click="runCommandByIndex('_forcemerge',indexName)"
          >{{ $t('强制合并索引') }}
          </el-button>
          <el-link type="danger">{{ $t('【forcemerge操作,手动释放磁盘空间】') }}</el-link>
          <el-popover
            placement="top-start"
            :title="$t('提示')"
            width="200"
            trigger="hover"
            :content="$t('为了让最新的数据可以立即被搜索到')"
          >
            <el-button

              slot="reference"
              :disabled="loadingGroup['_refresh']"
              v-loading="loadingGroup['_refresh']"
              class="filter-item"
              size="mini"
              type="primary"
              icon="el-icon-refresh"
              @click="runCommandByIndex('_refresh',indexName)"
            >{{ $t('刷新索引') }}
            </el-button>
          </el-popover>

          <el-popover
            placement="top-start"
            :title="$t('提示')"
            width="200"
            trigger="hover"
            :content="$t('让数据持久化到磁盘中')"
          >
            <el-button

              slot="reference"
              :disabled="loadingGroup['_flush']"
              v-loading="loadingGroup['_flush']"
              size="mini"
              type="info"
              icon="el-icon-s-open"
              class="filter-item"
              @click="runCommandByIndex('_flush',indexName)"
            >{{ $t('将索引刷新到磁盘') }}
            </el-button>
          </el-popover>

          <el-button

            :disabled="loadingGroup['_cache/clear']"
            v-loading="loadingGroup['_cache/clear']"
            class="filter-item"
            size="mini"
            type="warning"
            icon="el-icon-toilet-paper"
            @click="runCommandByIndex('_cache/clear',indexName)"
          >{{ $t('清理缓存') }}
          </el-button>

          <el-button

            :disabled="loadingGroup['deleteIndex']"
            v-loading="loadingGroup['deleteIndex']"
            class="filter-item"
            type="danger"
            size="mini"
            icon="el-icon-delete"
            @click="deleteIndex(indexName,'deleteIndex')"
          >{{ $t('删除索引') }}
          </el-button>
          <el-button
            :disabled="loadingGroup['empty']"
            v-loading="loadingGroup['empty']"
            class="filter-item"
            type="danger"
            size="mini"
            icon="el-icon-delete"
            @click="runCommandByIndex('empty',indexName)"
          >{{ $t('清空索引') }}
          </el-button>
        </div>
      </el-tabs>

    </el-drawer>
    <settings
      v-if="openSettings"
      :index-name="indexName"
      :settings-type="settingsType"
      @finished="search"
      :open="openSettings"
      @close="closeSettings"
    />
    <mappings
      v-if="openMappings"
      :index-name="indexName"
      :title="mappingTitle"
      :open="openMappings"
      @close="closeMappings"
    />

  </div>
</template>

<script>
import {Finish, IsFinish} from '@/api/guid'
import Driver from 'driver.js' // import driver.js
import 'driver.js/dist/driver.min.css' // import driver.js css
import steps from '@/views/indices/guide'
import {clone} from '@/utils/index'
import {filterData} from '@/utils/table'
import {CatAction, OptimizeAction, RecoverCanWrite} from '@/api/es'
import {bigNumberTransform} from '@/utils/format'
import {CreateAction, DeleteAction, GetSettingsAction, GetSettingsInfoAction, StatsAction} from '@/api/es-index'
import {esSettingsWords} from '@/utils/base-data'
import {ListAction, UpdateMappingAction} from '@/api/es-map'

import writeXlsxFile from "write-excel-file";

export default {
  name: 'indices',

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
      indexTishiList: [],
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
      total: 0,
      connectLoading: false,
      page: 1,
      limit: 10,
      pageshow: true,
      list: [],
      input: '',
      status: '',
      mappings: {},
      selectIndexList: [],
      openMappings: false,
      allList: [],
      max: 8,
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
    openMappingDialog(index) {
      this.indexName = index
      this.openMappings = true
    },
    async finishGuid() {
      const {data, code, msg} = await Finish({'guid_name': this.modName})
    },
    async startGuid() {
      const {data, code, msg} = await IsFinish({'guid_name': this.modName})

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
        this.activeData = JSON.stringify(res.data['list'][indexName].mappings, null, '\t')
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

      if (activeData[activeDataKeys[0]].toString() != 'false') {
        input['properties'] = activeData[activeDataKeys[0]]
        input['type_name'] = activeDataKeys[0]
      } else {
        input['properties'] = activeData
      }

      this.loadingGroup['saveMappinng'] = true
      const {data, code, msg} = await UpdateMappingAction(input)
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
    },
    openMappingEditDialog(indexName, haveMapping) {
      if (haveMapping) {
        this.mappingTitle = this.$t('新增字段')
      } else {
        this.mappingTitle = this.$t('新增映射结构')
      }
      this.indexName = indexName

      this.openMappings = true
    },
    async GetMapAction() {
      const input = {}
      input['es_connect'] = this.$store.state.baseData.EsConnectID
      const {code, data, msg} = await ListAction(input)
      if (code == 0) {
        this.mappings = data['list']
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
            this.activeData = JSON.stringify(res.data['list'][this.indexName].mappings, null, '\t')
            if (Object.keys(res.data['list'][this.indexName].mappings).length == 0) {
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

          const {data} = await GetSettingsAction(input)

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
    querySearch(queryString, cb) {

      let queryData = JSON.parse(JSON.stringify(this.indexTishiList))
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
    deleteIndex(indexName, loadingType) {
      this.$confirm('确定删除该索引吗?', '警告', {
        confirmButtonText: '确认',
        cancelButtonText: '取消',
        type: 'warning'
      })
        .then(async () => {
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
      this.pageLimit()
    },
    // 当当前页改变
    handleCurrentChange(val) {
      console.log(`当前页: ${val}`)
      this.page = val
      this.pageLimit()
    },
    pageLimit() {
      this.list = this.allList.filter((item, index) =>
        index < this.page * this.limit && index >= this.limit * (this.page - 1)
      )
    },
    searchData() {
      this.connectLoading = true
      const form = {
        cat: 'CatIndices',
        es_connect: this.$store.state.baseData.EsConnectID
      }
      this.indexTishiList = []
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
            this.indexTishiList.push({'value': obj["index"], 'data': obj["index"]})
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
          this.allList = tmpList
          this.total = tmpList.length
          this.pageLimit()
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
    },
    exportXlsx() {
      this.page = 1
      this.pageshow = false
      this.exportXlsxData()
      this.$nextTick(() => {
        this.pageshow = true
      })
    },
    exportXlsxData() {
      this.connectLoading = true
      const form = {
        cat: 'CatIndices',
        es_connect: this.$store.state.baseData.EsConnectID
      }
      this.indexTishiList = []
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
            this.indexTishiList.push({'value': obj["index"], 'data': obj["index"]})
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
          this.allList = tmpList
          this.total = tmpList.length
          this.pageLimit()
          writeXlsxFile(this.allList, {
             schema:[{
              column: this.$t("索引健康状态"),
              type: String,
              value: (row) => {return row.health},
             },{
              column: this.$t("索引名称"),
              type: String,
              value: (row) => {return row.index},
             },{
              column: this.$t("索引uuid"),
              type: String,
              value: (row) => {return row.uuid.toString()},
             },{
              column: this.$t("索引主分片数"),
              type: Number,
              value: (row) => {return row.pri},
             },{
              column: this.$t("索引副本分片数量"),
              type: Number,
              value: (row) => {return row.rep},
             },{
              column: this.$t("索引副本分片数量"),
              type: Number,
              value: (row) => {return row.rep},
             },{
              column: this.$t("索引文档总数"),
              type: String,
              value: (row) => {return row["docs.count"].toString()},
             },{
              column: this.$t("索引中删除状态的文档"),
              type: String,
              value: (row) => {return row["docs.deleted"].toString()},
             },{
              column: this.$t("主分片的大小"),
              type: String,
              value: (row) => {return row["pri.store.size"]},
             },{
              column: this.$t("主分片+副本分分片的大小"),
              type: String,
              value: (row) => {return row["store.size"]},
             }],
             fileName: "indexes.xlsx",
          }).then(res => {this.connectLoading = false}).catch(err => {this.connectLoading = false;this.$message({ type: 'error',message: err.toString()}); })
        } else {
          this.connectLoading = false
          this.$message({
            type: 'error',
            message: res.msg
          })
        }

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

::v-deep :focus {
  outline: 0;
}


::v-deep .el-table .sort-caret.descending{
  bottom: 0px;
}

</style>
