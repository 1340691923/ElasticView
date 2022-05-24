<template>
  <div>
    <el-dialog
      :close-on-click-modal="false"
      :visible.sync="open"
      :title="title.concat(`【${indexName}】`)"
      width="60%"
      @close="closeDialog"
    >
      <div class="app-container">
        <div class="filter-container">

          <el-tag v-if="showTypeName" class="filter-item">type名：</el-tag>
          <el-input
            v-if="showTypeName"
            v-model="type_name"
            :readonly="typeReadonly"
            style="width: 200px"
            class="filter-item"
          />
          <el-tag class="filter-item">dynamic：</el-tag>
          <el-select v-model="dynamic" class="filter-item" size="mini">
            <el-option :label="$t('动态映射')" value="true" />
            <el-option :label="$t('静态映射')" value="false" />
            <el-option :label="$t('严格映射')" value="strict" />
          </el-select>

          <el-button
            v-loading="loading"
            :disabled="loading"
            icon="el-icon-check"
            type="success"
            size="mini"
            class="filter-item"
            @click="saveMappinng"
          >{{ $t('保存/修改映射') }}
          </el-button>
        </div>
        <vue-json-helper
          v-if="showVueJsonHelper"
          :size="size"
          :names="names"
          :json-str="jsonStr"
          :root-flag="rootFlag"
          :open-flag="openFlag"
          :back-top-flag="backTopFlag"
          :shadow-flag="false"
          :border-flag="false"
          @jsonListener="jsonListener"
        />
      </div>
    </el-dialog>
  </div>
</template>

<script>
import { ListAction, UpdateMappingAction } from '@/api/es-map'

export default {
  name: 'Mapping',
  components: {
    'VueJsonHelper': () => import('@/views/indices/components/Helper')
  },
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
      loading: false,
      dynamic: 'false',
      drawerShow: false,
      connectLoading: false,
      size: 'small',
      names: [
        { key: 'Root', name: 'properties' },
        { key: 'type', name: '数据类型' },
        { key: 'format', name: '时间格式化' },
        { key: 'analyzer', name: '分词器' },
        { key: 'normalizer', name: '分析器' },
        { key: 'boost', name: '权重' },
        { key: 'coerce', name: '强制类型转换' },

        { key: 'copy_to', name: '合并参数' },
        { key: 'doc_values', name: '文档值' },
        { key: 'dynamic', name: '动态设置' },
        { key: 'enabled', name: '是否开启字段' },
        { key: 'fielddata', name: '字段数据' },
        { key: 'ignore_above', name: '字段保存最大长度' },
        { key: 'ignore_malformed', name: '忽略格式不对的数据' },
        { key: 'include_in_all', name: '_all 查询包含字段' },
        { key: 'index_options', name: '索引设置' },
        { key: 'index', name: '是否可以被搜索' },
        { key: 'fields', name: '字段' },

        { key: 'norms', name: '标准信息' },
        { key: 'null_value', name: '空值' },
        { key: 'position_increment_gap', name: '短语位置间隙' },
        { key: 'properties', name: '属性' },
        { key: 'search_analyzer', name: '搜索分析器' },
        { key: 'similarity', name: '匹配算法' },
        { key: 'store', name: '字段是否被存储' },
        { key: 'term_vector', name: '词根信息' }
      ],
      rootFlag: true,
      openFlag: true,
      backTopFlag: false,
      jsonStr: '{}',
      showVueJsonHelper: false,
      type_name: '',
      ver: 6,
      showTypeName: false,
      typeReadonly: false
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
      const activeData = {}
      activeData['properties'] = properties
      activeData['dynamic'] = this.dynamic

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
        this.$emit('finished')
      } else {
        this.$message({
          type: 'error',
          message: msg
        })
      }
    },
    refreshVueJsonHelper() {
      this.showVueJsonHelper = false
      this.$nextTick(() => {
        this.showVueJsonHelper = true
      })
    },
    async init() {
      const input = {}
      input['es_connect'] = this.$store.state.baseData.EsConnectID

      input['index_name'] = this.indexName

      const { data, code, msg } = await ListAction(input)

      if (code != 0) {
        this.$message({
          type: 'error',
          message: msg
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
          } else {
            this.type_name = mappings[0]
            this.typeReadonly = true
            this.dynamic = data.list[this.indexName].mappings[this.type_name].hasOwnProperty('dynamic') ? data.list[this.indexName].mappings[this.type_name]['dynamic'] : 'false'
            this.jsonStr = data.list[this.indexName].mappings[this.type_name].hasOwnProperty('properties') ? JSON.stringify(data.list[this.indexName].mappings[this.type_name].properties) : '{}'
            this.jsonStr = JSON.stringify(data.list[this.indexName].mappings[this.type_name].properties)
          }
          break
        case 7:
        case 8:
          this.dynamic = data.list[this.indexName].mappings.hasOwnProperty('dynamic') ? data.list[this.indexName].mappings['dynamic'] : 'false'
          this.jsonStr = data.list[this.indexName].mappings.hasOwnProperty('properties') ? JSON.stringify(data.list[this.indexName].mappings.properties) : '{}'
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
