<template>
  <div>
    <div class="filter-container operate">
      <el-tag class="filter-item">操作</el-tag>
      <el-button
        v-loading="loadingGroup['saveMappinng']"
        class="filter-item"
        size="small"
        type="warning"
        icon="el-icon-toilet-paper"
        @click="saveMappinng"
      >保存/修改映射
      </el-button>
      <index-select :multiple="true" :clearable="true" placeholder="迁移别名到多个索引上" @change="changeAnotherIndex" />
    </div>
    <json-editor
      v-model="mapping"
      v-loading="tabLoading"
      styles="width: 100%"
      :read="true"
      title="映射"
      @getValue="getMapping"
    />

  </div>
</template>

<script>

export default {
  name: 'Alias',
  components: {
    'IndexSelect': () => import('@/components/index/select')
  },
  props: {
    indexName: {
      type: String,
      default: ''
    }
  },
  data() {
    return {
      loadingGroup: {
        'saveMappinng': false
      },
      anotherIndex: '',
      mapping: ''
    }
  },
  mounted() {

  },
  methods: {
    getMapping(v) {
      this.mapping = v
    }

  }
}
</script>

<style scoped>

</style>
