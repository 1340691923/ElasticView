<template>
  <div class="app-container">
    <back-to-top />
    <el-card class="box-card">
      <el-tabs v-model="activeName">

        <el-tab-pane label="查看索引信息" name="CatIndices">
          <cat v-if="activeName == 'CatIndices'" :cat-type="activeName" :table-info="catData[activeName]" />
        </el-tab-pane>
        <el-tab-pane label="显示别名,过滤器,路由信息" name="CatAliases">
          <cat v-if="activeName == 'CatAliases'":cat-type="activeName" :table-info="catData[activeName]" />
        </el-tab-pane>
        <el-tab-pane label="显示每个节点分片数量、占用空间" name="CatAllocation">
          <cat v-if="activeName == 'CatAllocation'" :cat-type="activeName" :table-info="catData[activeName]" />
        </el-tab-pane>
        <el-tab-pane label="显示索引文档的数量" name="CatCount">
          <cat v-if="activeName == 'CatCount'" :cat-type="activeName" :table-info="catData[activeName]" />
        </el-tab-pane>
        <el-tab-pane label="查看集群健康状况" name="CatHealth">
          <cat v-if="activeName == 'CatHealth'" :cat-type="activeName" :table-info="catData[activeName]" />
        </el-tab-pane>
        <el-tab-pane label="显示索引分片信息" name="CatShards">
          <cat v-if="activeName == 'CatShards'" :cat-type="activeName" :table-info="catData[activeName]" />
        </el-tab-pane>
      </el-tabs>
    </el-card>
  </div>
</template>

<script>
export default {
  components: {
    'BackToTop': () => import('@/components/BackToTop/index'),
    'Cat': () => import('@/components/Cat/index')
  },
  data() {
    return {
      activeName: 'CatIndices',
      catData: {
        'CatIndices': [
          {
            data: 'health',
            desc: '索引健康状态',
            sort: true,
            width: 100
          },
          {
            data: 'status',
            desc: '索引的开启状态',
            sort: true,
            width: 100
          },
          {
            data: 'index',
            desc: '索引名称',
            sort: true,
            width: 180
          },
          {
            data: 'uuid',
            desc: '索引uuid',
            sort: true,
            width: 300
          },
          {
            data: 'pri',
            desc: '索引主分片数',
            sort: true,
            width: 100
          },
          {
            data: 'rep',
            desc: '索引副本分片数量',
            sort: true,
            width: 100
          },
          {
            data: 'docs->count',
            desc: '索引中文档总数',
            sort: true,
            width: 100
          },
          {
            data: 'docs->deleted',
            desc: '索引中删除状态的文档',
            sort: true,
            width: 100
          },
          {
            data: 'store->size',
            desc: '主分片+副本分分片的大小',
            sort: true,
            width: 100
          },
          {
            data: 'pri->store->size',
            desc: '主分片的大小',
            sort: true,
            width: 200
          }
        ],
        'CatAliases': [
          {
            data: 'alias',
            desc: '别名',
            sort: true,
            width: 300
          },
          {
            data: 'index',
            desc: '索引别名指向',
            sort: true,
            width: 250
          },
          {
            data: 'filter',
            desc: '过滤器',
            sort: true,
            width: 220
          },
          {
            data: 'index',
            desc: '索引路由',
            sort: true,
            width: 220
          },
          {
            data: 'routing->search',
            desc: '搜索路由',
            sort: true,
            width: 220
          },
          {
            data: 'is_write_index',
            desc: '写索引',
            sort: true,
            width: 220
          }
        ],
        'CatAllocation': [
          {
            data: 'host',
            desc: '节点host',
            sort: true,
            width: 150
          },
          {
            data: 'ip',
            desc: '节点ip',
            sort: true,
            width: 150
          },
          {
            data: 'node',
            desc: '节点名称',
            sort: true,
            width: 170
          },
          {
            data: 'shards',
            desc: '节点承载分片数量',
            sort: true,
            width: 150
          },
          {
            data: 'disk->indices',
            desc: '索引占用空间大小',
            sort: true,
            width: 150
          },
          {
            data: 'disk->used',
            desc: '节点所在机器已使用的磁盘空间大小',
            sort: true,
            width: 150
          },
          {
            data: 'disk->avail',
            desc: '节点可用空间大小',
            sort: true,
            width: 170
          },
          {
            data: 'disk->total',
            desc: '节点总空间大小',
            sort: true,
            width: 170
          },
          {
            data: 'disk->percent',
            desc: '节点磁盘占用百分比',
            sort: true,
            width: 100
          }
        ],
        'CatCount': [
          {
            data: 'epoch',
            desc: '自标准时间（1970-01-01 00:00:00）以来的秒数',
            sort: true,
            width: 500
          },
          {
            data: 'timestamp',
            desc: '时分秒,utc时区',
            sort: true,
            width: 500
          },
          {
            data: 'count',
            desc: '文档总数',
            sort: true,
            width: 500
          }
        ],
        'CatHealth': [
          {
            data: 'cluster',
            desc: '集群名称',
            sort: true,
            width: 120
          },
          {
            data: 'status',
            desc: '集群状态',
            sort: true,
            width: 100
          },
          {
            data: 'node->total',
            desc: '节点总数',
            sort: true,
            width: 100
          },
          {
            data: 'node->data',
            desc: '数据节点总数',
            sort: true,
            width: 100
          },
          {
            data: 'shards',
            desc: '分片总数',
            sort: true,
            width: 100
          },
          {
            data: 'pri',
            desc: '主分片总数',
            sort: true,
            width: 100
          },
          {
            data: 'relo',
            desc: '复制节点总数',
            sort: true,
            width: 100
          },
          {
            data: 'init',
            desc: '初始化节点总数',
            sort: true,
            width: 100
          },
          {
            data: 'unassign',
            desc: '未分配分片总数',
            sort: true,
            width: 100
          },
          {
            data: 'pending_tasks',
            desc: '待定任务总数',
            sort: true,
            width: 100
          },
          {
            data: 'max_task_wait_time',
            desc: '等待最长任务的等待时间',
            sort: true,
            width: 100
          }, {
            data: 'active_shards_percent',
            desc: '活动分片百分比',
            sort: true,
            width: 100
          },
          {
            data: 'epoch',
            desc: '自标准时间（1970-01-01 00:00:00）以来的秒数',
            sort: true,
            width: 100
          },
          {
            data: 'timestamp',
            desc: '时分秒,utc时区',
            sort: true,
            width: 100
          }
        ],
        'CatShards': [
          {
            data: 'index',
            desc: '索引名称',
            sort: true,
            width: 220
          },
          {
            data: 'shard',
            desc: '分片序号',
            sort: true,
            width: 100
          },
          {
            data: 'prirep',
            desc: '分片类型(p为主分片，r为复制分片)',
            sort: true,
            width: 220
          },
          {
            data: 'state',
            desc: '分片状态',
            sort: true,
            width: 100
          },
          {
            data: 'docs',
            desc: '该分片存放的文档数量',
            sort: true,
            width: 140
          },
          {
            data: 'store',
            desc: '该分片占用的存储空间大小',
            sort: true,
            width: 200
          },
          {
            data: 'ip',
            desc: '该分片所在的服务器ip',
            sort: false,
            width: 200
          },
          {
            data: 'node',
            desc: '该分片所在的节点名称',
            sort: true,
            width: 200
          }
        ]
      }
    }
  },
  computed: {},
  created() {

  },
  methods: {}
}
</script>

<style lang="scss" scoped>
  .app-container {
    .roles-table {
      margin-top: 30px;
    }

    .permission-tree {
      margin-bottom: 30px;
    }
  }
</style>
