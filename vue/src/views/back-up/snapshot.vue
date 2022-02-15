<template>
  <div class="app-container">
    <el-card class="box-card">
      <div class="filter-container">
        <el-select
          v-model="repositoryName"
          clearable
          filterable
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
        <el-button type="warning" @click.native="openAddDialog = true">新建快照</el-button>
        <el-button type="success" @click="search">刷新</el-button>
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

        <el-table-column align="center" label="快照名" prop="id" width="100" />
        <el-table-column align="center" label="备份索引数" prop="indices" width="50" />

        <el-table-column align="center" label="状态" prop="status" width="100" />

        <el-table-column align="center" label="状态" width="100">
          <template slot-scope="scope">
            <el-tag v-if="scope.row.status == 'SUCCESS'" type="success">成功</el-tag>
            <el-tag v-else-if="scope.row.status == 'IN_PROGRESS'" type="warning">还在进行中</el-tag>
            <el-tag v-else>{{ scope.row.status }}</el-tag>
          </template>
        </el-table-column>

        <el-table-column align="center" label="开始时间" prop="start_time" width="100" />

        <el-table-column align="center" label="结束时间" prop="end_time" width="100" />

        <el-table-column align="center" label="开始详细时间" width="180">
          <template slot-scope="scope">
            <div>{{ timestampToTime(scope.row.start_epoch) }}</div>
          </template>
        </el-table-column>

        <el-table-column align="center" label="结束详细时间" width="180">
          <template slot-scope="scope">
            <div>{{ timestampToTime(scope.row.end_epoch) }}</div>
          </template>
        </el-table-column>

        <el-table-column align="center" label="耗费时长" prop="duration" width="90" />
        <el-table-column align="center" label="分片总数" prop="total_shards" width="90" />
        <el-table-column align="center" label="成功分片数" prop="successful_shards" width="90" />
        <el-table-column align="center" label="失败分片数" prop="failed_shards" width="90" />

        <el-table-column align="center" label="操作" fixed="right" width="350">
          <template slot-scope="scope">
            <el-button-group>
              <el-button type="success" size="small" icon="el-icon-search" @click="look(scope.row.id)">查看</el-button>

              <el-button
                type="danger"
                size="small"
                icon="el-icon-delete"
                @click="SnapshotDeleteAction(scope.row.id)"
              >删除
              </el-button>
              <el-button
                type="warning"
                size="small"
                icon="el-icon-refresh"
                @click="openRestore(scope.row.id)"
              >恢复
              </el-button>
              <el-button
                type="primary"
                size="small"
                icon="el-icon-refresh"
                @click="status(scope.row.id)"
              >状态
              </el-button>

            </el-button-group>

          </template>
        </el-table-column>
      </el-table>

      <el-drawer
        ref="drawer"
        title="快照详细信息"
        :before-close="drawerHandleClose"
        :visible.sync="drawerShow"
        direction="rtl"
        close-on-press-escape
        destroy-on-close
        size="50%"
      >

        <json-editor

          v-model="JSON.stringify(snapshotDetail,null, '\t')"
          class="res-body"
          styles="width: 100%"
          :read="true"
          title="快照详细信息"
        />
      </el-drawer>
    </el-card>
    <add
      v-if="openAddDialog"
      :open="openAddDialog"
      @close="closeAddDialog"
    />
    <snapshot-restore
      v-if="openRestoreDialog"
      :repository="repositoryName"
      :snapshot="name"
      :open="openRestoreDialog"
      @close="closeoRestoreDialog"
    />

  </div>
</template>

<script>
import { SnapshotListAction, SnapshotRepositoryListAction, SnapshotDeleteAction, SnapshotDetailAction, SnapshotStatusAction } from '@/api/es-backup'

import { timestampToTime } from '@/utils/time'

export default {
  name: 'ResTable',
  components: {
    'JsonEditor': () => import('@/components/JsonEditor/index'),
    'Add': () => import('@/views/back-up/components/addSnapshot'),
    'SnapshotRestore': () => import('@/views/back-up/components/snapshotRestore')

  },
  data() {
    return {
      openRestoreDialog: false,
      openAddDialog: false,
      loading: false,
      repositoryName: '',
      resData: {},
      name: '',
      drawerShow: false,
      tableData: [],
      index: 0,

      tableHeader: [],
      snapshotDetail: ''
    }
  },
  created() {
    this.initRepositoryName()
  },
  methods: {
    async status(name) {
      const input = {}
      input['es_connect'] = this.$store.state.baseData.EsConnectID
      input['snapshot'] = name
      input['repository'] = this.repositoryName

      const { data, code, msg } = await SnapshotStatusAction(input)
      if (code != 0) {
        this.$message({
          type: 'error',
          message: msg
        })
        return
      } else {
        this.snapshotDetail = data.snapshots
        this.drawerShow = true
        this.$message({
          type: 'success',
          message: msg
        })
        return
      }
    },
    openRestore(name) {
      this.name = name
      this.openRestoreDialog = true
    },
    timestampToTime(timestamp) {
      return timestampToTime(timestamp)
    },

    async SnapshotDeleteAction(name) {
      const input = {}
      input['es_connect'] = this.$store.state.baseData.EsConnectID
      input['snapshot'] = name
      input['repository'] = this.repositoryName

      const { data, code, msg } = await SnapshotDeleteAction(input)
      if (code != 0) {
        this.$message({
          type: 'error',
          message: msg
        })
        return
      } else {
        this.initRepositoryName()
        this.search()
        this.$message({
          type: 'success',
          message: msg
        })
        return
      }
    },
    closeoRestoreDialog(addSuccess) {
      this.openRestoreDialog = false
      if (addSuccess) {
        this.initRepositoryName()
        this.search()
      }
    },
    closeAddDialog(addSuccess) {
      this.openAddDialog = false
      if (addSuccess) {
        this.initRepositoryName()
        this.search()
      }
    },
    async initRepositoryName() {
      this.loading = true
      const input = {}
      input['es_connect'] = this.$store.state.baseData.EsConnectID
      const { data, code, msg } = await SnapshotRepositoryListAction(input)
      if (code == 0) {
        this.resData = data.res

        if (Object.keys(this.resData).length > 0) {
          this.repositoryName = Object.keys(this.resData)[0]
        }

        this.search()
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
      input['repository'] = this.repositoryName
      const { data, code, msg } = await SnapshotListAction(input)

      if (code != 0) {
        this.$message({
          type: 'error',
          message: msg
        })
        return
      } else {
        this.$message({
          type: 'success',
          message: msg
        })
        this.tableData = data
      }
    },

    openDetail(row, index, tmp) {
      this.name = row.name
      this.drawerShow = true
    },
    async look(index) {
      this.name = index
      const input = {}
      input['es_connect'] = this.$store.state.baseData.EsConnectID
      input['repository'] = this.repositoryName
      input['snapshot'] = index

      const { data, code, msg } = await SnapshotDetailAction(input)
      if (code != 0) {
        this.$message({
          type: 'error',
          message: msg
        })
        return
      } else {
        this.snapshotDetail = data.snapshots
        this.drawerShow = true
      }
    },
    drawerHandleClose(done) {
      done()
    }
  }
}
</script>

<style scoped>

</style>
