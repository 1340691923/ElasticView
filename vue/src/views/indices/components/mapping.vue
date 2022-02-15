<template>
  <div>
    <el-dialog :close-on-click-modal="false" :visible.sync="open" :title="title" width="95%" @close="closeDialog">
      <el-card class="box-card">
        <div class="filter-container">
          <el-tag class="filter-item">类型名</el-tag>
          <el-input
            v-model="typeName"
            placeholder="类型名"
            style="width: 500px"
            class="filter-item"
            clearable
          />
          <el-tag class="filter-item">dynamic</el-tag>

          <el-select v-model="dynamic" class="filter-item" clearable filterable>
            <el-option label="默认" value="" />
            <el-option label="true" value="true" />
            <el-option label="false" value="false" />
            <el-option label="strict" value="strict" />
          </el-select>
          <el-button class="filter-item" type="primary" @click.native="drawerShow = true">新增字段</el-button>
          <el-button class="filter-item" type="success" @click="addMapping">确认提交</el-button>
        </div>
        <el-table
          :data="tableData"
        >
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
          <el-table-column align="center" label="字段名" width="600">
            <template slot-scope="scope">
              {{ scope.row.name }}
            </template>
          </el-table-column>
          <el-table-column align="center" label="类型" width="600">
            <template slot-scope="scope">
              {{ scope.row.typ }}
            </template>
          </el-table-column>
          <el-table-column align="center" label="操作" fixed="right" width="400">
            <template slot-scope="scope">
              <el-button-group>
                <el-button
                  icon="el-icon-edit"
                  type="primary"
                  size="small"
                  @click="openDrawer(scope.row.index)"
                >编辑
                </el-button>
                <el-button
                  icon="el-icon-delete"
                  type="danger"
                  size="small"
                  @click.native="tableData.splice(scope.$index, 1)"
                >删除
                </el-button>
              </el-button-group>
            </template>
          </el-table-column>
        </el-table>

      </el-card>
    </el-dialog>
    <el-drawer
      ref="drawer"
      title="新增映射属性"
      :before-close="drawerHandleClose"
      :visible.sync="drawerShow"

      direction="rtl"
      close-on-press-escape
      destroy-on-close
      size="50%"
    >
      <el-form :model="properties" label-width="180px" label-position="right">
        <el-form-item label="字段名">
          <el-input v-model="properties.name" placeholder="字段名" />
        </el-form-item>
        <el-form-item label="类型">
          <el-select v-model="properties.typ" placeholder="类型" clearable filterable>
            <el-option label="binary" value="binary" />
            <el-option label="text" value="text" />
            <el-option label="boolean" value="boolean" />
            <el-option label="keyword" value="keyword" />
            <el-option label="constant_keyword" value="constant_keyword" />
            <el-option label="wildcard" value="wildcard" />
            <el-option label="long" value="long" />
            <el-option label="integer" value="integer" />
            <el-option label="short" value="short" />
            <el-option label="byte" value="byte" />
            <el-option label="double" value="double" />
            <el-option label="float" value="float" />
            <el-option label="half_float" value="half_float" />
            <el-option label="scaled_float" value="scaled_float" />
            <el-option label="unsigned_long" value="unsigned_long" />
            <el-option label="text" value="text" />
            <el-option label="text" value="text" />
            <el-option label="text" value="text" />
            <el-option label="text" value="text" />
            <el-option label="text" value="text" />
            <el-option label="text" value="text" />
            <el-option label="text" value="text" />
            <el-option label="text" value="text" />
            <el-option label="text" value="text" />
          </el-select>
        </el-form-item>
        <el-form-item label="analyzer">
          <el-input v-model="properties.analyzer" placeholder="analyzer" />
        </el-form-item>
        <el-form-item label="search_analyzer">
          <el-input v-model="properties.search_analyzer" placeholder="search_analyzer" />
        </el-form-item>
        <el-form-item label="search_quote_analyzer">
          <el-input v-model="properties.search_quote_analyzer" placeholder="search_quote_analyzer" />
        </el-form-item>
        <el-form-item label="doc_values">
          <el-input v-model="properties.doc_values" placeholder="doc_values" />
        </el-form-item>
        <el-form-item label="store">
          <el-input v-model="properties.store" placeholder="store" />
        </el-form-item>
        <el-form-item label="boost">
          <el-input v-model="properties.boost" placeholder="boost" />
        </el-form-item>
        <el-form-item label="index">
          <el-input v-model="properties.index" placeholder="index" />
        </el-form-item>
        <el-form-item label="null_value">
          <el-input v-model="properties.null_value" placeholder="null_value" />
        </el-form-item>
        <el-form-item label="meta">
          <el-input v-model="properties.meta" placeholder="meta" />
        </el-form-item>
        <el-form-item label="index">
          <el-input v-model="properties.index" placeholder="index" />
        </el-form-item>
        <el-form-item label="eager_global_ordinals">
          <el-input v-model="properties.eager_global_ordinals" placeholder="eager_global_ordinals" />
        </el-form-item>
        <el-form-item label="ignore_above">
          <el-input v-model="properties.ignore_above" placeholder="ignore_above" />
        </el-form-item>
        <el-form-item label="index_options">
          <el-input v-model="properties.index_options" placeholder="index_options" />
        </el-form-item>
        <el-form-item label="norms">
          <el-input v-model="properties.norms" placeholder="norms" />
        </el-form-item>
        <el-form-item label="similarity">
          <el-input v-model="properties.similarity" placeholder="similarity" />
        </el-form-item>
        <el-form-item label="normalizer">
          <el-input v-model="properties.normalizer" placeholder="normalizer" />
        </el-form-item>
        <el-form-item label="split_queries_on_whitespace">
          <el-input v-model="properties.split_queries_on_whitespace" placeholder="split_queries_on_whitespace" />
        </el-form-item>
        <el-form-item label="meta">
          <el-input v-model="properties.meta" placeholder="meta" />
        </el-form-item>
        <el-form-item label="coerce">
          <el-input v-model="properties.coerce" placeholder="coerce" />
        </el-form-item>
        <el-form-item label="similarity">
          <el-input v-model="properties.similaritysimilarity" placeholder="similarity" />
        </el-form-item>
        <el-form-item label="similarity">
          <el-input v-model="properties.similaritysimilarity" placeholder="similarity" />
        </el-form-item>
        <el-form-item label="similarity">
          <el-input v-model="properties.similaritysimilarity" placeholder="similarity" />
        </el-form-item>
        <el-form-item label="similarity">
          <el-input v-model="properties.similaritysimilarity" placeholder="similarity" />
        </el-form-item>
        <el-form-item label="similarity">
          <el-input v-model="properties.similaritysimilarity" placeholder="similarity" />
        </el-form-item>
      </el-form>
    </el-drawer>
  </div>
</template>

<script>
export default {
  name: 'Mapping',

  props: {
    indexName: {
      type: String,
      default: ''
    },
    mappings: {
      type: Object,
      default: {}
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
      drawerShow: false,
      connectLoading: false,
      mappingList: [],
      typeName: '',
      dynamic: '',
      tableData: [{ name: 1 }],
      properties: {
        typ: ''
      }
    }
  },
  mounted() {
    console.log(this.indexName, 'this.indexName')
    this.typeName = this.indexName
  },
  methods: {
    openDrawer() {

    },
    drawerHandleClose(done) {
      done()
    },
    addColumn() {

    },
    addMapping() {
      console.log('新增映射')
    },
    closeDialog() {
      this.$emit('close')
    }
  }
}
</script>

<style scoped>

</style>
