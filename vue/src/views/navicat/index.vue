<template>
  <div style="display:flex;justify-content:space-between">
    <div class="content_xwl">
      <div class="header_xwl" style="background: white">
        <div class="root_xwl">
          <div class="main_xwl">
            <span class="title_xwl" style="color: #202d3f">&nbsp;&nbsp;数据管理</span>
          </div>
          <div class="actions_xwl">

          </div>
        </div>
      </div>
      <split-pane :min-percent="4" :default-percent="20" split="vertical">
        <template slot="paneL">
          <div
            id="scollL"
            style="height: 95%;width: 100px;display: inline-block; height: 100%;vertical-align: top;width: 100%;background: white;"
          >
            <div style="width: 100%;height: calc(100% - 80px); overflow-x: hidden; overflow-y: auto;padding: 10px">
              <div
                style="width: 100%;height:90px;margin-bottom: 0px;z-index: 80;position: absolute;top:0px;left: 0px;padding: 20px;
                ;border-bottom: 1px solid #f0f2f5;background: white;display: flex;align-items: center;justify-content: center">
                <el-input v-model="filterStr" placeholder="请输入索引名" clearable></el-input>
              </div>

              <el-menu
                style="margin-top: 100px"
                active-text-color="rgb(64 158 255)"
                v-loading="loadingMenu">
                <div v-for="(v,index2) in getIndexList" :key="index2" :index="index2">

                  <el-menu-item
                    @click.native="clickItem(index2)"
                    :index="index2">
                    <el-dropdown>
                      <span class="el-dropdown-link">
                         <i v-if="v.health == 'red'" style="color: red" class="el-icon-s-grid"></i>
                         <i v-if="v.health == 'green'" style="color: #13ce66" class="el-icon-s-grid"></i>
                         <i v-if="v.health == 'yellow'" style="color: #ffba00" class="el-icon-s-grid"></i>
                      </span>
                      <el-dropdown-menu slot="dropdown">
                        <el-dropdown-item @click.native="deleteIndex(index2)" icon="el-icon-delete">删除
                        </el-dropdown-item>
                      </el-dropdown-menu>
                    </el-dropdown>
                    <span slot="title">{{v.index}}【{{v.storeSize}}】</span>
                  </el-menu-item>
                </div>
              </el-menu>
            </div>
          </div>
        </template>
        <template slot="paneR">
          <split-pane :default-percent="80" :min-percent="3" split="vertical">
            <template slot="paneL">
              <div
                style="width: 100%;height: calc(100% - 80px); overflow-x: hidden; overflow-y: auto;padding: 10px">

              <crud :attrMapProp="attrMap" :eventAttrOptionsProp="eventAttrOptions" :indexName="currentIndexName" v-if="refreshTab"></crud>
              </div>
            </template>
            <template slot="paneR">
              <div
                style="height: 95%;width: 100px;display: inline-block; height: 100%;vertical-align: top;width: 100%;background: white;"
              >
                <el-tabs v-loading="tabLoading" type="border-card" v-model="activeName">
                  <el-tab-pane label="索引设置" name="settings">
                    <json-editor
                      height="720"
                      v-if="refreshTab"
                      v-model="JSON.stringify(currentSettings,null, '\t')"
                      styles="width: 100%"
                      :read="true"
                      :title="currentIndexName"
                    />
                  </el-tab-pane>
                  <el-tab-pane label="映射结构" name="mappings">
                    <json-editor
                      height="720"
                      v-if="refreshTab"
                      v-model="JSON.stringify(currentMappings,null, '\t')"
                      styles="width: 100%"
                      :read="true"
                      :title="currentIndexName"
                    />
                  </el-tab-pane>
                </el-tabs>
              </div>
            </template>
          </split-pane>
        </template>
      </split-pane>
    </div>
  </div>
</template>

<script>

  import {CatAction} from '@/api/es'
  import {filterData} from '@/utils/table'
  import {DeleteAction, GetSettingsAction} from '@/api/es-index'
  import {ListAction} from '@/api/es-map'

  export default {
    name: 'navicat',
    components: {

      'JsonEditor': () => import('@/components/JsonEditor/index'),
      'Crud': () => import('@/views/navicat/crud'),
    },
    computed: {
      getIndexList() {
        if (this.filterStr == '') {
          return this.indexList.sort(this.compare("docsCount", true))
        }
        return filterData(this.indexList, this.filterStr).sort(this.compare("docsCount", true))
      }
    },
    data() {
      return {

        tabLoading: false,
        refreshTab: true,
        activeName: 'settings',
        loadingMenu: false,
        indexList: [],
        filterStr: '',
        currentIndexName: '',
        currentSettings: {},
        eventAttrOptions: [],
        attrMap:[]
      }
    },
    async beforeMount() {
      await this.init()
    },
    async mounted() {

    },
    methods: {
      reloadTab() {

        this.refreshTab = false // 先关闭，
        this.$nextTick(() => {
          this.refreshTab = true
        })
      },
      async clickItem(index) {
        this.currentIndexName = this.indexList[index].index
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
        let eventAttrOptions = [{"label": "筛选字段", "options": []}]
        let attrMap = { "2": []}
        let propertiesObj = {}
        switch (res.data.ver) {
          case 6:
            propertiesObj = this.currentMappings[this.currentIndexName].mappings[Object.keys(this.currentMappings[this.currentIndexName].mappings)[0]].properties
            break
          case 7:
            propertiesObj = this.currentMappings[this.currentIndexName].mappings.properties
            break
        }

        const Int = 1
        const Float = 2
        const String = 3
        //const DateTime = 4
        let propertiesObjKeys = Object.keys(propertiesObj)
        for (let k in propertiesObjKeys) {
          //其他类型暂时不支持 有人用我再继续写
          if (propertiesObj[propertiesObjKeys[k]].type) {
            eventAttrOptions[0].options.push({"value": propertiesObjKeys[k], "label": propertiesObjKeys[k]})
            let obj = {"attribute_name": propertiesObjKeys[k], "show_name": propertiesObjKeys[k]}
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


        /*for(let v of this.currentMappings[this.currentIndexName]){
          let tmp = { "value": "xwl_server_time", "label": "服务端入库时间" }
        }*/


        this.currentSettings = res2.data
        this.tabLoading = false
        this.reloadTab()
      },
      async deleteIndex(index) {

        let indexName = this.indexList[index].index
        this.$confirm('确定删除(' + indexName + ')该索引吗?', '警告', {
          confirmButtonText: '确认',
          cancelButtonText: '取消',
          type: 'warning'
        })
          .then(async () => {
            const $message = ELEMENT.Message({
              message: '  索引删除中...',
              customClass: 'theme-message',
              type: 'success',
              duration: 0,
              iconClass: 'el-icon-loading'
            })
            const input = {}
            input['es_connect'] = this.$store.state.baseData.EsConnectID
            input['index_name'] = indexName
            let res = await DeleteAction(input)
            $message.close()
            if (this.currentIndexName == indexName) {
              this.currentIndexName = ''
            }
            if (res.code == 0 || res.code == 200) {
              this.$message({
                type: 'success',
                message: indexName + '已删除'
              })
              this.indexList.splice(index, 1)
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
          var value1 = a[property];
          var value2 = b[property];
          if (sort) return value2 - value1;
          return value1 - value2;
        }
      },
      async init() {
        const form = {
          cat: 'CatIndices',
          es_connect: this.$store.state.baseData.EsConnectID
        }
        this.indexList = []
        this.loadingMenu = true
        const res = await CatAction(form)

        if (res.code != 0) {
          this.$message({
            type: 'error',
            message: res.msg
          })
          this.loadingMenu = false
          return
        }
        for (let v of res.data) {
          let obj = {health: v.health, index: v.index, storeSize: v["store.size"], docsCount: v['docs.count']}
          this.indexList.push(obj)
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
