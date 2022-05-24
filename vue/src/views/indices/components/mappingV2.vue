<template>
  <div>
    <div class="filter-container operate">
      <el-tag class="filter-item">{{ $t('操作') }}</el-tag>
      <el-button

        v-loading="loadingGroup['saveMappinng']"
        class="filter-item"
        size="mini"
        type="warning"
        icon="el-icon-toilet-paper"
        @click="saveMappinng"
      >{{ $t('保存/修改映射') }}
      </el-button>
      <index-select :multiple="true" :clearable="true" :placeholder="$t('迁移别名到多个索引上')" @change="changeAnotherIndex" />
    </div>
    <json-editor
      v-model="mapping"
      v-loading="tabLoading"
      styles="width: 100%"
      :read="true"
      :title="$t('映射')"
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
