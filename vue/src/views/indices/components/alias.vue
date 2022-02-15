<template>
  <div>
    <el-form>
      <el-form-item>
        <el-button type="success" icon="el-icon-plus" @click="addAlias">新增别名</el-button>
        <el-button type="primary" icon="el-icon-check" @click="batchAdd">批量提交</el-button>
        <el-button type="info" icon="el-icon-delete" @click="getAlias">重置</el-button>
        <index-select :multiple="true" :clearable="true" placeholder="迁移别名到多个索引上" @change="changeAnotherIndex" />

      </el-form-item>
      <el-form-item
        v-for="(alias, index) in aliasList"
        :key="index"
        :label="'别名' + Number(index+1)"
      >
        <el-input v-model="aliasList[index].name" :readonly="aliasList[index].types !='new'" style="width:300px" />

        <el-button v-show="anotherIndex.length>0" icon="el-icon-right" type="success" @click="moveAliasToIndex(index)">迁移</el-button>

        <el-button
          v-clipboard:copy="aliasList[index].name"
          v-clipboard:success="onCopy"
          v-clipboard:error="onError"
          icon="el-icon-copy"
          type="success"
        >复制
        </el-button>
        <el-button
          v-show="aliasList[index].types =='new'"
          type="primary"
          icon="el-icon-check"
          @click="submitForm(index)"
        >提交
        </el-button>
        <el-button icon="el-icon-delete" type="danger" @click="removeAlias(index)">删除</el-button>

      </el-form-item>
    </el-form>
  </div>
</template>

<script>
import { GetAliasAction } from '@/api/es-index'
import { OperateAliasAction } from '../../../api/es-index'

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
      aliasList: [],
      anotherIndex: []
    }
  },
  mounted() {
    this.getAlias()
  },
  methods: {
    async moveAliasToIndex(index) {
      const alias = this.aliasList[index].name
      const input = {}
      input['types'] = 3
      input['new_index_list'] = this.anotherIndex
      input['alias_name'] = alias.trim()
      input['es_connect'] = this.$store.state.baseData.EsConnectID
      const { data, code, msg } = await OperateAliasAction(input)
      console.log(data, code, msg)
      if (code == 0) {
        this.$message({
          type: 'success',
          message: msg
        })
        return
      } else {
        this.$message({
          type: 'error',
          message: msg
        })
        return
      }
    },
    changeAnotherIndex(v) {
      this.anotherIndex = v
    },
    onCopy(e) { 　　 // 复制成功
      this.$message({
        message: '复制成功！',
        type: 'success'
      })
    },
    onError(e) {　　 // 复制失败
      this.$message({
        message: '复制失败！',
        type: 'error'
      })
    },
    async getAlias() {
      this.aliasList = []
      const input = {}
      input['es_connect'] = this.$store.state.baseData.EsConnectID
      input['index_name'] = this.indexName
      const { data, code, msg } = await GetAliasAction(input)
      console.log(data, code, msg)
      if (code == 0) {
        for (const k in data) {
          this.aliasList.push({ name: data[k].AliasName, types: 'old' })
        }
      } else {
        this.$message({
          message: msg,
          type: 'error'
        })
      }
    },
    async submitForm(index) {
      const alias = this.aliasList[index].name
      const input = {}
      input['types'] = 1
      input['index_name'] = this.indexName
      input['alias_name'] = alias.trim()
      input['es_connect'] = this.$store.state.baseData.EsConnectID
      const { data, code, msg } = await OperateAliasAction(input)
      console.log(data, code, msg)
      if (code == 0) {
        this.$message({
          type: 'success',
          message: msg
        })
        return
      } else {
        this.$message({
          type: 'error',
          message: msg
        })
        return
      }
    },
    async batchAdd() {
      const aliasList = []
      for (const alias of this.aliasList) {
        if (alias.types == 'new') {
          aliasList.push(alias.name)
        }
      }
      const input = {}
      input['types'] = 4
      input['index_name'] = this.indexName
      input['new_alias_name_list'] = aliasList
      input['es_connect'] = this.$store.state.baseData.EsConnectID
      const { data, code, msg } = await OperateAliasAction(input)
      if (code == 0) {
        this.$message({
          type: 'success',
          message: msg
        })
        return
      } else {
        this.$message({
          type: 'error',
          message: msg
        })
        return
      }
    },
    async removeAlias(index) {
      if (this.aliasList[index].types == 'new') {
        this.aliasList.splice(index, 1)
      } else if (this.aliasList[index].types == 'old') {
        const alias = this.aliasList[index].name
        const input = {}
        input['types'] = 2
        input['index_name'] = this.indexName
        input['alias_name'] = alias.trim()
        input['es_connect'] = this.$store.state.baseData.EsConnectID
        const { data, code, msg } = await OperateAliasAction(input)
        if (code == 0) {
          this.$message({
            type: 'success',
            message: msg
          })
          this.aliasList.splice(index, 1)
          return
        } else {
          this.$message({
            type: 'error',
            message: msg
          })
          return
        }
      }
    },
    addAlias() {
      this.aliasList.push({ name: '', types: 'new' })
    }

  }
}
</script>

<style scoped>

</style>
