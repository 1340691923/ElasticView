<template>
  <div>
    <el-form-item :label="isOne?'主表':'副表'">
      <el-select v-model="form.selectTable" filterable @change="GetTableColumns">
        <el-option v-for="(v,k,index) in tables" :key="index" :value="v" :label="v" />
      </el-select>
      <el-button
        size="mini"
        icon="el-icon-delete"
        type="danger"
        @click="deleteTable"
      >{{ $t('删除该表') }}{{ currentIndex }}
      </el-button>
    </el-form-item>
    <el-form-item v-if="!isOne" :label="$t('设置表间联系')">
      hello.
      <el-select v-model="form.selectTable" filterable @change="GetTableColumns">
        <el-option v-for="(v,k,index) in tables" :key="index" :value="v" :label="v" />
      </el-select>
      <el-tag>=</el-tag>
      test.
      <el-select v-model="form.selectTable" filterable @change="GetTableColumns">
        <el-option v-for="(v,k,index) in tables" :key="index" :value="v" :label="v" />
      </el-select>

    </el-form-item>
    <el-form-item label="表字段：">
      <div class="col-transfer">
        <el-transfer
          v-model="form.cols"
          :titles="['全部字段', '当前ES索引字段']"
          :button-texts="['移除字段', '添加字段']"
          filterable
          :filter-method="filterMethod"
          filter-placeholder="请操作成员"
          :data="allCols"
        />
      </div>
    </el-form-item>
  </div>
</template>

<script>
import { GetTableColumns, Tables } from '@/api/datax'

export default {
  name: 'Tables',
  props: {
    tables: {
      default: [],
      type: Array
    },
    selectType: {
      default: '',
      type: String
    },
    isOne: {
      default: false,
      type: Boolean
    }
  },
  data() {
    return {
      allCols: [],
      form: {
        selectTable: '',
        cols: []
      }
    }
  },
  computed: {},
  methods: {
    deleteTable() {
      this.$emit('deleteTable', this.currentIndex)
    },
    getSelectTypeObj() {
      return JSON.parse(this.selectType)
    },
    async GetTableColumns() {
      const res = await GetTableColumns({ id: this.getSelectTypeObj()['id'], table_name: this.form.selectTable })
      if (res.code != 0) {
        this.$message({
          type: 'error',
          message: res.msg
        })
        return
      }
      if (res.data == null) res.data = []
      this.allCols = []
      for (const v of res.data) {
        const obj = {
          key: v.Field,
          label: v.Comment == '' ? v.Field : `${v.Field}【${v.Comment}】`,
          disabled: false
        }
        this.allCols.push(obj)
      }
    },
    async changeTable() {
      await this.getTables()
      await this.GetTableColumns()
    },
    async GetTableColumns() {
      const res = await GetTableColumns({ id: this.getSelectTypeObj()['id'], table_name: this.form.selectTable })
      if (res.code != 0) {
        this.$message({
          type: 'error',
          message: res.msg
        })
        return
      }
      if (res.data == null) res.data = []
      this.allCols = []
      for (const v of res.data) {
        const obj = {
          key: v.Field,
          label: v.Comment == '' ? v.Field : `${v.Field}【${v.Comment}】`,
          disabled: false
        }
        this.allCols.push(obj)
      }
    },
    async getTables() {
      const res = await Tables({ id: this.getSelectTypeObj()['id'] })
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
    filterMethod(query, item) {
      return item.label.indexOf(query) > -1
    }
  }
}
</script>

<style scoped>
.col-transfer >>> .el-transfer-panel {
  width: 350px;
}
</style>
