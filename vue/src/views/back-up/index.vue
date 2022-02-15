<template>
  <div class="app-container">
    <el-card class="box-card">
      <div class="filter-container">
        <el-select
          v-model="snapshotNameList"
          clearable
          filterable
          multiple
          placeholder="请选择存储库"
          :loading="loading"
          @change="search"
        >
          <el-option
            v-for="(v,k,index) of resData"
            :key="index"
            :label="k"
            :value="k"
          />
        </el-select>
        <el-button type="warning" @click.native="openAddDialog = true">新建存储库</el-button>
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
        <el-table-column align="center" label="存储库名" prop="name" width="200" />
        <el-table-column align="center" label="存储库地址" prop="location" width="300" />
        <el-table-column align="center" label="类型" prop="type" width="200" />

        <el-table-column align="center" label="是否压缩" prop="compress" width="100" />
        <el-table-column align="center" label="分解块大小" prop="chunk_size" width="100" />
        <el-table-column align="center" label="是否只读(默认false)" prop="readonly" width="100" />

        <el-table-column align="center" label="制作快照的速度" width="100">
          <template slot-scope="scope">
            <div v-if="scope.row.max_snapshot_bytes_per_sec != ''">{{ scope.row.max_snapshot_bytes_per_sec }}/s</div>
            <div v-else>20mb/s</div>
          </template>
        </el-table-column>
        <el-table-column align="center" label="快照恢复的速度" width="100">
          <template slot-scope="scope">
            <div v-if="scope.row.max_restore_bytes_per_sec != ''">{{ scope.row.max_restore_bytes_per_sec }}/s</div>
            <div v-else>20mb/s</div>
          </template>
        </el-table-column>

        <el-table-column align="center" label="操作" fixed="right" width="350">
          <template slot-scope="scope">
            <el-button-group>
              <el-button type="success" size="small" icon="el-icon-search" @click="look(scope.row.name)">查看</el-button>
              <el-button type="warning" size="small" icon="el-icon-edit" @click="updateSnapshotData(scope.row)">修改
              </el-button>
              <el-button
                type="danger"
                size="small"
                icon="el-icon-delete"
                @click="SnapshotDeleteRepositoryAction(scope.row.name)"
              >删除
              </el-button>
              <el-button
                type="warning"
                size="small"
                icon="el-icon-delete"
                @click="CleanupeRepository(scope.row.name)"
              >清理
              </el-button>

            </el-button-group>

          </template>
        </el-table-column>
      </el-table>

      <el-drawer

        ref="drawer"
        title="JSON数据"
        :before-close="drawerHandleClose"
        :visible.sync="drawerShow"

        direction="rtl"
        close-on-press-escape
        destroy-on-close
        size="50%"
      >

        <json-editor
          v-if="drawerShow"
          v-model="JSON.stringify(resData[name],null, '\t')"
          class="res-body"
          styles="width: 100%"
          :read="true"
          title="JSON数据"
        />
      </el-drawer>
    </el-card>
    <add
      v-if="openAddDialog"
      :snapshot-data="snapshotData"
      :open="openAddDialog"
      @close="closeAddDialog"
    />
  </div>
</template>

<script>
import { CleanupeRepositoryAction, SnapshotDeleteRepositoryAction, SnapshotRepositoryListAction, SnapshotListAction } from '@/api/es-backup'

export default {
  name: 'ResTable',
  components: {
    'JsonEditor': () => import('@/components/JsonEditor/index'),
    'Add': () => import('@/views/back-up/components/addRepository')
  },
  data() {
    return {
      snapshotData: null,
      openAddDialog: false,
      loading: false,
      snapshotNameList: [],
      resData: {},
      name: '',
      drawerShow: false,
      tableData: [],
      index: 0,
      tableHeader: []
    }
  },
  created() {
    this.initSnapshotName()
    this.search()
  },

  methods: {
    async CleanupeRepository(name) {
      const input = {}
      input['es_connect'] = this.$store.state.baseData.EsConnectID
      input['name'] = name
      const { data, code, msg } = await CleanupeRepositoryAction(input)
      if (code != 0) {
        this.$message({
          type: 'error',
          message: msg
        })
        return
      } else {
        this.initSnapshotName()
        this.search()
        this.$message({
          type: 'success',
          message: msg
        })
        return
      }
    },
    async SnapshotDeleteRepositoryAction(name) {
      const input = {}
      input['es_connect'] = this.$store.state.baseData.EsConnectID
      input['name'] = name
      const { data, code, msg } = await SnapshotDeleteRepositoryAction(input)
      if (code != 0) {
        this.$message({
          type: 'error',
          message: msg
        })
        return
      } else {
        this.initSnapshotName()
        this.search()
        this.$message({
          type: 'success',
          message: msg
        })
        return
      }
    },
    updateSnapshotData(row) {
      this.snapshotData = row
      this.openAddDialog = true
    },
    closeAddDialog(addSuccess) {
      this.openAddDialog = false
      this.snapshotData = null
      if (addSuccess) {
        this.initSnapshotName()
        this.search()
      }
    },
    async initSnapshotName() {
      this.loading = true
      const input = {}
      input['es_connect'] = this.$store.state.baseData.EsConnectID
      const { data, code, msg } = await SnapshotRepositoryListAction(input)
      if (code == 0) {
        this.resData = data.res
        this.$notify({
          title: '成功',
          dangerouslyUseHTMLString: true,
          message: `
<strong> <b style="color: orange">您设置的path.repo为${data.pathRepo.join(',')}</b><br></strong>
          `,
          type: 'success'
        })
      } else if (code == 199999) {
        this.$notify({
          title: 'Error',
          dangerouslyUseHTMLString: true,
          message: `
<strong>
            <i style="color: orange">path.repo没有设置</i><br>
<i>在elasticsearch.yml 配置文件中配置仓库base目录</i><br>
<i>添加path.repo: /tmp/tmp (/tmp/tmp 为快照备份所在文件夹, <br><i style="color: orange">注意</i>首先要先创建这个文件夹)</i>

          `,
          type: 'error'
        })
        return
      } else {
        this.$message({
          type: 'error',
          message: msg
        })
        return
      }
      this.loading = false
    },
    async search() {
      const input = {}
      input['es_connect'] = this.$store.state.baseData.EsConnectID
      input['snapshot_info_list'] = this.snapshotNameList
      input['repository'] = 'test3'
      const { data, code, msg } = await SnapshotRepositoryListAction(input)

      if (code != 0) {
        this.$message({
          type: 'error',
          message: msg
        })
        return
      } else {
        this.tableData = data.list
      }
    },

    openDetail(row, index, tmp) {
      this.name = row.name
      this.drawerShow = true
    },
    look(index) {
      this.name = index
      this.drawerShow = true
    },
    drawerHandleClose(done) {
      done()
    }
  }
}
</script>

<style scoped>

</style>
