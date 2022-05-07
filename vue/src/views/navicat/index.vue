<template>
  <div style="display:flex;justify-content:space-between">
    <div class="content_xwl">
      <div class="header_xwl" style="background: white">
        <div class="root_xwl">
          <div class="main_xwl">
            <span class="title_xwl" style="color: #202d3f">&nbsp;&nbsp;{{$t('数据管理')}}</span>
          </div>
          <div class="actions_xwl" >
            <el-button type="primary" @click.native="showIndexSettings = true">{{$t('显示该索引结构')}}</el-button>
          </div>
        </div>
      </div>
      <split-pane :min-percent="4" :default-percent="15" split="vertical">
        <template slot="paneL">
          <div
            id="scollL"
            style="height: 95%;width: 100px;display: inline-block; height: 100%;vertical-align: top;width: 100%;background: white;"
          >
            <div style="width: 100%;height: calc(100% - 80px); overflow-x: hidden; overflow-y: auto;padding: 10px">
              <div
                style="width: 100%;height:90px;margin-bottom: 0px;z-index: 80;position: absolute;top:0px;left: 0px;padding: 20px;
                ;border-bottom: 1px solid #f0f2f5;background: white;display: flex;align-items: center;justify-content: center"
              >
                <el-autocomplete
                  style="width: 90%"
                  clearable
                  :fetch-suggestions="querySearch"

                  v-model="filterStr"
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

              </div>

              <el-menu
                v-loading="loadingMenu"
                style="margin-top: 100px"
                active-text-color="rgb(64 158 255)"
              >
                <div v-for="(v,index2) in getIndexList" :key="index2" :index="index2">

                  <el-menu-item
                    :index="index2.toString()"
                    @click.native="clickItem(v['index'])"
                  >
                    <el-dropdown>
                      <span class="el-dropdown-link">
                        <i v-if="v.health == 'red'" style="color: red" class="el-icon-s-grid" />
                        <i v-if="v.health == 'green'" style="color: #13ce66" class="el-icon-s-grid" />
                        <i v-if="v.health == 'yellow'" style="color: #ffba00" class="el-icon-s-grid" />
                      </span>
                      <el-dropdown-menu slot="dropdown">
                        <el-dropdown-item icon="el-icon-delete" @click.native="deleteIndex(getIndexList[index2].k,getIndexList[index2].index)">{{$t('删除')}}
                        </el-dropdown-item>
                      </el-dropdown-menu>
                    </el-dropdown>
                    <span slot="title">{{ v.index }}【{{ v.storeSize }}】</span>
                  </el-menu-item>
                </div>
              </el-menu>
            </div>
          </div>
        </template>
        <template  slot="paneR">
          <div
            style="width: 100%;height: calc(100% - 80px); overflow-x: hidden; overflow-y: auto;padding: 10px"
          >

            <crud v-if="refreshTab" :attr-map-prop="attrMap" :event-attr-options-prop="eventAttrOptions" :index-name="currentIndexName" />
          </div>
          <el-drawer
            ref="drawer"
            :title="$t('索引结构')"

            :visible.sync="showIndexSettings"

            direction="rtl"
            close-on-press-escape
            destroy-on-close
            size="50%"
          >
            <div
                 style="height: 95%;width: 100px;display: inline-block; height: 100%;vertical-align: top;width: 100%;background: white;"
            >
              <el-tabs v-model="activeName" v-loading="tabLoading" type="border-card">
                <el-tab-pane :label="$t('索引设置')" name="settings">
                  <json-editor
                    v-if="refreshTab"
                    v-model="JSON.stringify(currentSettings,null, '\t')"
                    height="720"
                    styles="width: 100%"
                    :read="true"
                    :title="currentIndexName"
                  />
                </el-tab-pane>
                <el-tab-pane :label="$t('映射结构')" name="mappings">
                  <json-editor
                    v-if="refreshTab"
                    v-model="JSON.stringify(currentMappings,null, '\t')"
                    height="720"
                    styles="width: 100%"
                    :read="true"
                    :title="currentIndexName"
                  />
                </el-tab-pane>
              </el-tabs>
            </div>

          </el-drawer>
<!--          <split-pane :default-percent="98" :min-percent="2" split="vertical">
            <template slot="paneL">

            </template>
            <template slot="paneR">

            </template>
          </split-pane>-->
        </template>
      </split-pane>
    </div>
  </div>
</template>

<script>

import { CatAction } from '@/api/es'
import { filterData } from '@/utils/table'
import { DeleteAction, GetSettingsAction } from '@/api/es-index'
import { ListAction } from '@/api/es-map'
import ElementUI from 'element-ui';

export default {
  name: 'Navicat',
  components: {
    'JsonEditor': () => import('@/components/JsonEditor/index'),
    'Crud': () => import('@/views/navicat/crud')
  },
  data() {
    return {
      indexTishiList:[],
      currentMappings:{},
      tabLoading: false,
      refreshTab: true,
      activeName: 'settings',
      loadingMenu: false,
      indexList: [],
      filterStr: '',
      currentIndexName: '',
      currentSettings: {},
      eventAttrOptions: [],
      showIndexSettings:false,
      attrMap: []
    }
  },
  computed: {
    getIndexList() {
      if (this.filterStr == '') {
        return this.indexList
      }
      return filterData(this.indexList, this.filterStr))
    }
  },
  async beforeMount() {
    await this.init()
  },
  async mounted() {
    this.toggleSideBar()
  },
  methods: {
    toggleSideBar() {
      this.$store.dispatch('app/toggleSideBar')
    },
    reloadTab() {
      this.refreshTab = false // 先关闭，
      this.$nextTick(() => {
        this.refreshTab = true
      })
    },
    async clickItem(index) {
      this.currentIndexName = index
      const input = {}
      input['es_connect'] = this.$store.state.baseData.EsConnectID
      input['index_name'] = this.currentIndexName
      this.tabLoading = true
      const res = await ListAction(input)

      if (res.code != 0) {
        this.tabLoading = false
        this.$message({
          type: 'error',
          message: res.msg
        })
      }

      const res2 = await GetSettingsAction(input)
      if (res2.code != 0) {
        this.tabLoading = false
        this.$message({
          type: 'error',
          message: res2.msg
        })
      }
      this.currentMappings = res.data.list
      const eventAttrOptions = [{ 'label': '筛选字段', 'options': [] }]
      const attrMap = { '2': [] }
      let propertiesObj = {}
      switch (res.data.ver) {
        case 6:
          propertiesObj = this.currentMappings[this.currentIndexName].mappings[Object.keys(this.currentMappings[this.currentIndexName].mappings)[0]].properties
          break
        case 7:
          propertiesObj = this.currentMappings[this.currentIndexName].mappings.properties
          break
        case 8:
          propertiesObj = this.currentMappings[this.currentIndexName].mappings.properties
          break
      }
      const Int = 1
      const Float = 2
      const String = 3
      // const DateTime = 4
      const propertiesObjKeys = Object.keys(propertiesObj)
      for (const k in propertiesObjKeys) {
        // 其他类型暂时不支持 有人用我再继续写
        if (propertiesObj[propertiesObjKeys[k]].type) {
          eventAttrOptions[0].options.push({ 'value': propertiesObjKeys[k], 'label': propertiesObjKeys[k] })
          const obj = { 'attribute_name': propertiesObjKeys[k], 'show_name': propertiesObjKeys[k] }
          switch (propertiesObj[propertiesObjKeys[k]].type) {
            case 'text':
            case 'keyword':
              obj['data_type'] = String

              break
            case 'byte':
            case 'short':
            case 'integer':
            case 'long':
              obj['data_type'] = Int

              break
            case 'float':
            case 'half_float':
            case 'scaled_float':
            case 'double':
              obj['data_type'] = Float

              break
          }
          attrMap['2'].push(obj)
        }
      }
      this.attrMap = attrMap
      this.eventAttrOptions = eventAttrOptions
      this.currentSettings = res2.data
      this.tabLoading = false
      this.reloadTab()
    },
    deleteListByK(k) {
      for (const index in this.indexList) {
        const v = this.indexList[index]
        if (v['k'].toString() == k.toString()) {
          this.indexList.splice(index, 1)
          break
        }
      }
    },
    async deleteIndex(key, indexName) {
      this.$confirm('确定删除(' + indexName + ')该索引吗?', '警告', {
        confirmButtonText: '确认',
        cancelButtonText: '取消',
        type: 'warning'
      })
        .then(async() => {
          const $message = ElementUI.Message({
            message: '  索引删除中...',
            customClass: 'theme-message',
            type: 'success',
            duration: 0,
            iconClass: 'el-icon-loading'
          })
          const input = {}
          input['es_connect'] = this.$store.state.baseData.EsConnectID
          input['index_name'] = indexName
          const res = await DeleteAction(input)
          $message.close()
          if (this.currentIndexName == indexName) {
            this.currentIndexName = ''
          }
          if (res.code == 0 || res.code == 200) {
            this.$message({
              type: 'success',
              message: indexName + '已删除'
            })
            this.deleteListByK(key)
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
    compare(property, sort) {
      return (a, b) => {
        var value1 = a[property]
        var value2 = b[property]
        if (sort) return value2 - value1
        return value1 - value2
      }
    },
    querySearch(queryString, cb) {

      let queryData = JSON.parse(JSON.stringify(this.indexTishiList))
      if(queryString == undefined)queryString = ""
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
    async init() {
      const form = {
        cat: 'CatIndices',
        es_connect: this.$store.state.baseData.EsConnectID
      }
      this.indexList = []
      this.loadingMenu = true
      this.indexTishiList = []
      const res = await CatAction(form)

      if (res.code != 0) {
        this.$message({
          type: 'error',
          message: res.msg
        })
        this.loadingMenu = false
        return
      }
      for (const k in res.data) {
        const v = res.data[k]
        const obj = { health: v.health, index: v.index, k: k, storeSize: v['store.size'], docsCount: v['docs.count'] }
        this.indexList.push(obj)
        this.indexTishiList.push( {'value': v.index, 'data': v.index})
      }
      this.loadingMenu = false
    }
  }
}
</script>

<style scoped src="@/styles/event.css"/>

<style>
  .eventNameDisplayInput .ant-input {
    resize: none;
    border: none;
  }

  .eventNameDisplayInput .ant-input:focus {
    border: none;
    box-shadow: none;
  }
</style>
