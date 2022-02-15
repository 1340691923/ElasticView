<template>
  <div>
    <el-table
      :loading="connectLoading"
      :data="mappingList"
    >

      <el-table-column align="center" label="名称" min-width="10%">
        <template slot-scope="scope">
          {{ scope.row.name }}
        </template>
      </el-table-column>
      <el-table-column align="center" label="数据类型"min-width="10%">
        <template slot-scope="scope">
          <el-tag v-if="scope.row.hasOwnProperty('type')" type="warning"> {{ scope.row.type }}</el-tag>
          <el-tag v-else type="danger">多级</el-tag>
        </template>
      </el-table-column>
      <el-table-column align="center" label="dynamic" min-width="8%">
        <template slot-scope="scope">
          {{ scope.row.dynamic }}
        </template>
      </el-table-column>
      <el-table-column align="center" label="下级" type="expand">
        <template slot-scope="scope">
          <div v-if="scope.row.hasOwnProperty('properties')">
            <mapping-table :mappings="scope.row" />
          </div>
        </template>
      </el-table-column>
    </el-table>
  </div>
</template>

<script>
export default {
  name: 'MappingTable',
  components: {
    'MappingTable': () => import('@/views/indices/components/mappingTable')
  },
  props: {
    mappings: {
      type: Object,
      default: {}
    }
  },
  data() {
    return {
      mappingList: [],
      connectLoading: false
    }
  },
  mounted() {
    console.log(this.mappings, 'this.mappings', this.mappings.hasOwnProperty('properties'))
    if (this.mappings.hasOwnProperty('properties')) {
      const properties = this.mappings['properties']

      for (const name of Object.keys(properties)) {
        const propertiesValObj = properties[name]
        propertiesValObj['name'] = name
        this.mappingList.push(propertiesValObj)
      }
    }
  }

}
</script>

<style scoped>

</style>
