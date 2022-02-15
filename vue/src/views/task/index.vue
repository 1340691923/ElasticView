<template>
  <div class="app-container">
    <el-card class="box-card">
      <div class="filter-container">
        <el-select
          v-model="input.actions"
          reserve-keyword
          collapse-tags
          placeholder="行为"
          clearable
          multiple
          clearable
          filterable
          @change="searchByFilter()"
        >
          <el-option v-for="(action, index) in actions" :key="index" :label="action" :value="action" />
        </el-select>
        <el-select
          v-model="input.task_id"
          reserve-keyword
          collapse-tags
          placeholder="任务ID"
          clearable
          multiple
          clearable
          filterable
          @change="searchByFilter()"
        >
          <el-option v-for="(taskID, index) in taskIdList" :key="index" :label="taskID" :value="taskID" />
        </el-select>
        <el-select
          v-model="input.node_id"
          reserve-keyword
          collapse-tags
          placeholder="节点ID"
          clearable
          multiple
          clearable
          filterable
          @change="searchByFilter()"
        >
          <el-option v-for="(node, index) in nodeIdList" :key="index" :label="node" :value="node" />
        </el-select>
        <el-select
          v-model="input.parent_task_id"
          placeholder="父任务ID"
          clearable
          clearable
          filterable
          @change="searchByFilter()"
        >
          <el-option v-for="(taskID, index) in parentTaskIdList" :key="index" :label="taskID" :value="taskID" />
        </el-select>
        <el-button type="primary" icon="el-icon-search" @click="search(true)">搜索</el-button>
        <el-button type="success" icon="el-icon-refresh" @click="refresh">刷新</el-button>
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
        <el-table-column align="center" label="行为" prop="action" width="200" />
        <el-table-column align="center" label="任务详细ID" prop="taskID" width="300" />
        <el-table-column align="center" label="任务ID" prop="id" width="100" />

        <el-table-column align="center" label="节点" prop="node" width="200" />

        <el-table-column align="center" label="父任务id" prop="parent_task_id" width="300">
          <template slot-scope="scope">
            <div>{{ scope.row.parent_task_id }}</div>
          </template>
        </el-table-column>
        <el-table-column align="center" label="开始详细时间" width="180">
          <template slot-scope="scope">
            <div>{{ timestampToTime(scope.row.start_time_in_millis) }}</div>
          </template>
        </el-table-column>
        <el-table-column align="center" label="状态" width="300">
          <template slot-scope="scope">

            <el-popover
                          placement="top-start"
                          title="描述"
                          width="600"
                          trigger="hover"
            >

              <div>{{ JSON.stringify(scope.row.status) }}</div>
                          <span slot="reference">{{ JSON.stringify(scope.row.status).substr(0, 50) + "..." }}</span>

              </span slot="reference">
            </el-popover>
          </template>
        </el-table-column>
        <el-table-column align="center" label="运行时间（秒）" width="180">
          <template slot-scope="scope">
            <div>{{ Number(scope.row.running_time_in_nanos/1000000000) }}</div>
          </template>
        </el-table-column>

        <el-table-column align="center" label="是否可撤销" prop="cancellable" width="100">
          <template slot-scope="scope">
            <div>{{ scope.row.cancellable }}</div>
          </template>
        </el-table-column>

        <el-table-column align="center" label="描述" prop="description" width="300">
          <template slot-scope="scope">

            <el-popover
                          placement="top-start"
                          title="描述"
                          width="600"
                          trigger="hover"
            >

              <div>{{ scope.row.description }}</div>
                          <span slot="reference">{{ scope.row.description.substr(0, 50) + "..." }}</span>
              </span slot="reference">
            </el-popover>
          </template>
        </el-table-column>
        <el-table-column align="center" label="开始时间" prop="start_time" width="100" />
        <el-table-column align="center" label="开始时间（毫秒）" prop="start_time_in_millis" width="150" />

        <el-table-column align="center" label="类型" prop="type" width="100" />

        <el-table-column align="center" label="操作" fixed="right" width="350">
          <template slot-scope="scope">
            <el-button-group>
              <el-button type="success" size="small" icon="el-icon-search" @click="look(scope.row)">查看</el-button>

              <el-button
                type="danger"
                size="small"
                icon="el-icon-delete"
                @click="CancelAction(scope.row.taskID)"
              >取消
              </el-button>

            </el-button-group>

          </template>
        </el-table-column>
      </el-table>

      <el-drawer
        ref="drawer"
        title="任务详细信息"
        :before-close="drawerHandleClose"
        :visible.sync="drawerShow"
        direction="rtl"
        close-on-press-escape
        destroy-on-close
        size="50%"
      >

        <json-editor
          v-model="taskDetail"
          styles="width: 100%"
          :read="true"
          title="任务详细信息"
        />
      </el-drawer>

    </el-card>
  </div>
</template>

<script>
import Moment from 'moment'
import { CancelAction, ListAction } from '@/api/es-task'

export default {
  name: 'ResTable',
  components: {
    'JsonEditor': () => import('@/components/JsonEditor/index'),
    'Add': () => import('@/views/back-up/components/addSnapshot'),
    'SnapshotRestore': () => import('@/views/back-up/components/snapshotRestore')
  },
  data() {
    return {
      taskDetail: '',
      loading: false,
      drawerShow: false,
      tableData: [],
      taskIdList: [],
      parentTaskIdList: [],
      nodeIdList: [],
      actions: [],
      input: {
        task_id: [],
        node_id: [],
        parent_task_id: '',
        actions: []
      }
    }
  },
  created() {
    this.search()
  },
  methods: {
    inArray(search, array) {
      for (var i in array) {
        if (array[i] == search) {
          return true
        }
      }
      return false
    },
    searchByFilter() {
      const tableData = []
      let flag = false
      if (this.input.parent_task_id != '') {
        flag = true
        for (const data of this.tableData) {
          if (data.parent_task_id == this.input.parent_task_id) {
            tableData.push(data)
          }
        }
      }
      if (this.input.actions.length > 0) {
        flag = true
        for (const data of this.tableData) {
          if (this.inArray(data.action, this.input.actions)) {
            tableData.push(data)
          }
        }
      }
      if (this.input.node_id.length > 0) {
        flag = true
        for (const data of this.tableData) {
          if (this.inArray(data.node, this.input.node_id)) {
            tableData.push(data)
          }
        }
      }
      if (this.input.task_id.length > 0) {
        flag = true
        for (const data of this.tableData) {
          if (this.inArray(data.taskID, this.input.task_id)) {
            tableData.push(data)
          }
        }
      }

      if (flag == true) {
        this.tableData = tableData
      } else {
        console.log(11)
        this.search()
      }
    },
    async CancelAction(taskid) {
      const input = {}
      input['es_connect'] = this.$store.state.baseData.EsConnectID
      input['task_id'] = taskid

      const { data, code, msg } = await CancelAction(input)
      if (code != 0) {
        this.$message({
          type: 'error',
          message: msg
        })
        return
      } else {
        this.search()
        this.$message({
          type: 'success',
          message: msg
        })
        return
      }
    },
    look(obj) {
      console.log(obj)
      this.taskDetail = JSON.stringify(obj, null, '\t')
      this.drawerShow = true
    },
    drawerHandleClose(done) {
      done()
    },
    timestampToTime(timeStamp) {
      const stamp = new Date(timeStamp)
      const time = Moment(stamp).format('YYYY-MM-DD HH:mm:ss')
      return time
    },
    async search(clear) {
      if (clear) {
        this.input = {
          task_id: [],
          node_id: [],
          parent_task_id: '',
          actions: []
        }
      }

      const input = {}

      input['es_connect'] = this.$store.state.baseData.EsConnectID

      const { data, code, msg } = await ListAction(input)

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
        this.tableData = []
        this.taskIdList = []
        this.parentTaskIdList = []
        this.nodeIdList = []

        const nodeSet = new Set()

        const actionsSet = new Set()
        for (const taskId in data) {
          const v = data[taskId]
          v['taskID'] = taskId
          this.tableData.push(v)
          this.taskIdList.push(taskId)
          if (v.parent_task_id != '') {
            this.parentTaskIdList.push(v.parent_task_id)
          }
          nodeSet.add(v.node)
          actionsSet.add(v.action)
        }
        this.nodeIdList = Array.from(nodeSet)
        this.actions = Array.from(actionsSet)
      }
    },
    async refresh() {
      await this.search()
      this.searchByFilter()
    }
  }
}
</script>

<style scoped>

</style>
