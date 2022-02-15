<template>
  <div>
    <el-dialog :close-on-click-modal="false" :visible.sync="open" :title="title" @close="closeDialog">
      <el-card class="box-card">
        <el-form label-width="500px" label-position="left">
          <el-form-item label="索引名称">
            <el-input v-model="indexName" placeholder="索引名称" :disabled="settingsType != 'add'" />
          </el-form-item>
          <el-form-item label="number_of_shards (分片数)">
            <el-input v-model="form.number_of_shards" type="number" style="width: 300px" :disabled="isOpen" />
          </el-form-item>
          <el-form-item label="number_of_replicas (副本数)">
            <el-input v-model="form.number_of_replicas" type="number" style="width: 300px" />
          </el-form-item>
          <el-form-item label="refresh_interval (索引的刷新时间间隔)">
            <el-input v-model="form.refresh_interval" placeholder="索引的刷新时间间隔" />
          </el-form-item>
          <el-form-item label="translog.sync_interval (translog同步到磁盘的时间间隔)">
            <el-input v-model="form['index.translog.sync_interval']" :disabled="isOpen" />
          </el-form-item>

          <el-form-item label="translog.durability (tanslog同步到本地的方式)">
            <el-select v-model="form['index.translog.durability']" style="width: 300px" clearable filterable>
              <el-option label="request (直接写入到磁盘中)" value="request" />
              <el-option label="async (写入cache中，再写入磁盘)" value="async" />
            </el-select>
          </el-form-item>

          <el-form-item label="flush_threshold_size (满足translog同步的容量)">
            <el-input v-model="form['index.translog.flush_threshold_size']" />
          </el-form-item>
          <el-form-item label="merge.scheduler.max_thread_count (调高合并的最大线程)">
            <el-input v-model="form['index.merge.scheduler.max_thread_count']" />
          </el-form-item>
          <el-form-item label="merge.policy.max_merged_segment (最大分段大小)">
            <el-input v-model="form['index.merge.policy.max_merged_segment']" />
          </el-form-item>
          <el-form-item label="merge.policy.segments_per_tier (每层所允许的分段数)">
            <el-input v-model="form['index.merge.policy.segments_per_tier']" type="number" style="width: 300px" />
          </el-form-item>
          <el-form-item label="shard.check_on_startup (是否应在索引打开前检查分片是否损坏)">
            <el-select
              v-model="form['index.shard.check_on_startup']"
              clearable
              style="width: 300px"
              :disabled="isOpen"
              filterable
            >
              <el-option label="true(检查物理和逻辑损坏，这将消耗大量内存和CPU)" :value="true" />
              <el-option label="false(不检查物理和逻辑损坏)" :value="false" />
              <el-option label="checksum(检查物理损坏)" value="checksum" />
              <el-option label="fix(检查物理和逻辑损坏。有损坏的分片将被集群自动删除)" value="fix" />
            </el-select>
          </el-form-item>

          <el-form-item label="routing_partition_size (自定义路由值可以转发的目的分片数)">
            <el-input
              v-model="form['index.routing_partition_size']"
              :disabled="isOpen"
              type="number"
              style="width: 300px"
            />
          </el-form-item>

          <el-form-item label="index.codec (默认使用LZ4压缩方式存储数据)">
            <el-input v-model="form['index.codec']" :disabled="isOpen" />
          </el-form-item>

          <el-form-item label="auto_expand_replicas (基于可用节点的数量自动分配副本数量)">
            <el-input
              v-model="form['index.auto_expand_replicas']"
              style="width: 300px"
            />
          </el-form-item>
          <el-form-item label="max_result_window  (from+size 的最大值)">
            <el-input v-model="form['index.max_result_window']" type="number" style="width: 300px" />
          </el-form-item>
          <el-form-item label="index.max_rescore_window( rescore 的 window_size 的最大值)">
            <el-input v-model="form['index.max_rescore_window']" type="number" style="width: 300px" />
          </el-form-item>
          <el-form-item label="blocks.read_only (允许写入和元数据更改)">
            <el-select v-model="form['index.blocks.read_only']" clearable filterable>
              <el-option label="是" value="true" />
              <el-option label="否" value="false" />
            </el-select>
          </el-form-item>
          <el-form-item label="blocks.read_only_allow_delete (只允许读和删数据，不允许增和改数据)">
            <el-select v-model="form['index.blocks.read_only_allow_delete']" clearable filterable>
              <el-option label="是" value="true" />
              <el-option label="否" value="false" />
            </el-select>
          </el-form-item>
          <el-form-item label="blocks.read (禁用对索引的读取操作)">
            <el-select v-model="form['index.blocks.read']" clearable filterable>
              <el-option label="是" value="true" />
              <el-option label="否" value="false" />
            </el-select>
          </el-form-item>
          <el-form-item label="blocks.write (禁用对索引的写入操作)">
            <el-select v-model="form['index.blocks.write']" clearable filterable>
              <el-option label="是" value="true" />
              <el-option label="否" value="false" />
            </el-select>
          </el-form-item>
          <el-form-item label="blocks.metadata (禁用索引元数据的读取和写入)">
            <el-select v-model="form['index.blocks.metadata']" clearable filterable>
              <el-option label="是" value="true" />
              <el-option label="否" value="false" />
            </el-select>
          </el-form-item>
          <el-form-item label="max_refresh_listeners (索引的每个分片上可用的最大刷新侦听器数)">
            <el-input v-model="form['index.max_refresh_listeners']" type="number" style="width: 300px" />
          </el-form-item>
        </el-form>
        <div style="text-align:right;">
          <el-button type="danger" icon="el-icon-close" @click="closeDialog">取消</el-button>
          <el-button type="primary" icon="el-icon-check" @click="confirmSettings">确认</el-button>
        </div>
      </el-card>
    </el-dialog>
  </div>
</template>

<script>

import { CreateAction, GetSettingsAction, CatStatusAction } from '@/api/es-index'
import { clone } from '@/utils/index'

export default {
  name: 'Add',
  components: {},
  props: {
    indexName: {
      type: String,
      default: ''
    },
    settingsType: {
      type: String,
      default: 'add'
    },
    open: {
      type: Boolean,
      default: false
    }
  },
  data() {
    return {
      staticConfig: [
        'number_of_shards',
        'index.translog.sync_interval',
        'index.shard.check_on_startup',
        'index.routing_partition_size',
        'index.codec'
      ],
      isOpen: false,
      form: {
        'index.shard.check_on_startup': '',
        'index.routing_partition_size': '',
        'index.codec': '',
        'index.auto_expand_replicas': '',
        'index.max_result_window': '',
        'index.max_rescore_window': '',
        'index.blocks.read_only': '',
        'index.blocks.read': '',
        'index.blocks.write': '',
        'index.blocks.metadata': '',
        'index.blocks.read_only_allow_delete': '',
        'index.max_refresh_listeners': '',
        'number_of_shards': '',
        'number_of_replicas': '',
        'refresh_interval': '',
        'index.translog.sync_interval': '',
        'index.translog.durability': '',
        'index.translog.flush_threshold_size': '',
        'index.merge.scheduler.max_thread_count': '',
        'index.merge.policy.max_merged_segment': '',
        'index.merge.policy.segments_per_tier': ''
      },
      dialogVisible: true
    }
  },
  computed: {
    isStatic() {
      return (this.settingsType == 'add') || (this.open = false)
    },
    title() {
      if (this.settingsType == 'add') {
        return '新增索引配置'
      } else {
        return '修改索引配置'
      }
    }
  },

  created() {
    if (this.indexName != '') {
      this.catIndexStatus()
      this.catIndexSettings()
    }
  },
  methods: {
    confirmSettings() {
      const form = this.form
      if (form['index.routing_partition_size'] >= form['index.number_of_shards']) {
        this.$message({
          type: 'error',
          message: '自定义路由值可以转发的目的分片数 必须小于 主分片数'
        })
        return
      }

      if (this.indexName == '') {
        this.$message({
          type: 'error',
          message: '索引名不能为空'
        })
        return
      }
      const input = {}
      input['es_connect'] = this.$store.state.baseData.EsConnectID
      input['settings'] = clone(this.form)
      input['index_name'] = this.indexName
      input['types'] = this.settingsType
      if (this.settingsType == 'update') {
        for (const config of this.staticConfig) {
          delete input['settings'][config]
        }
      }

      for (const settingsKey in input['settings']) {
        const settingsVal = input['settings'][settingsKey]
        if (settingsVal == '') {
          delete input['settings'][settingsKey]
        }
      }

      const loading = this.$loading({
        lock: true,
        text: 'Loading',
        spinner: 'el-icon-loading',
        background: 'rgba(0, 0, 0, 0.7)'
      })
      CreateAction(input).then(res => {
        console.log(input['settings'])
        if (res.code == 0 || res.code == 200) {
          this.$message({
            type: 'success',
            message: res.msg
          })
          this.settingsType = 'update'
          this.$emit('finished')
        } else {
          this.$message({
            type: 'error',
            message: res.msg
          })
        }
        loading.close()
      }).catch(err => {
        loading.close()
      })
    },
    catIndexSettings() {
      const input = {}
      input['es_connect'] = this.$store.state.baseData.EsConnectID
      input['index_name'] = this.indexName
      GetSettingsAction(input).then(res => {
        if (res.code == 0 || res.code == 200) {
          const index = res.data['index']
          this.form.number_of_replicas = index['number_of_replicas']
          this.form.number_of_shards = index['number_of_shards']
          if (index.hasOwnProperty('auto_expand_replicas')) {
            this.form['index.auto_expand_replicas'] = index['auto_expand_replicas']
          }
          if (index.hasOwnProperty('blocks')) {
            const blocks = index['blocks']
            if (blocks.hasOwnProperty('metadata')) {
              this.form['index.blocks.metadata'] = blocks['metadata']
            }
            if (blocks.hasOwnProperty('read')) {
              this.form['index.blocks.read'] = blocks['read']
            }
            if (blocks.hasOwnProperty('read_only')) {
              this.form['index.blocks.read_only'] = blocks['read_only']
            }
            if (blocks.hasOwnProperty('read_only_allow_delete')) {
              this.form['index.blocks.read_only_allow_delete'] = blocks['read_only_allow_delete']
            }
            if (blocks.hasOwnProperty('write')) {
              this.form['index.blocks.write'] = blocks['write']
            }
          }
          if (index.hasOwnProperty('merge')) {
            const merge = index['merge']
            if (merge.hasOwnProperty('policy')) {
              const policy = merge['policy']
              if (policy.hasOwnProperty('max_merged_segment')) {
                this.form['index.merge.policy.max_merged_segment'] = policy['max_merged_segment']
              }
              if (policy.hasOwnProperty('max_merged_segment')) {
                this.form['index.merge.policy.segments_per_tier'] = policy['segments_per_tier']
              }
            }
            if (merge.hasOwnProperty('scheduler')) {
              const scheduler = merge['scheduler']
              if (scheduler.hasOwnProperty('max_thread_count')) {
                this.form['index.merge.scheduler.max_thread_count'] = scheduler['max_thread_count']
              }
            }
          }

          if (index.hasOwnProperty('max_refresh_listeners')) {
            this.form['index.max_refresh_listeners'] = index['max_refresh_listeners']
          }
          if (index.hasOwnProperty('max_rescore_window')) {
            this.form['index.max_rescore_window'] = index['max_rescore_window']
          }
          if (index.hasOwnProperty('max_result_window')) {
            this.form['index.max_result_window'] = index['max_result_window']
          }

          if (index.hasOwnProperty('refresh_interval')) {
            this.form['refresh_interval'] = index['refresh_interval']
          }
          if (index.hasOwnProperty('translog')) {
            const translog = index['translog']
            if (translog.hasOwnProperty('durability')) {
              this.form['index.translog.durability'] = translog['durability']
            }
            if (translog.hasOwnProperty('flush_threshold_size')) {
              this.form['index.translog.flush_threshold_size'] = translog['flush_threshold_size']
            }
          }
          console.log(this.form, 'form')
        } else {
          this.$message({
            type: 'error',
            message: res.msg
          })
        }
      }).catch(err => {
        console.log(err)
      })
    },
    catIndexStatus() {
      const input = {}
      input['es_connect'] = this.$store.state.baseData.EsConnectID
      input['index_name'] = this.indexName

      CatStatusAction(input).then(res => {
        if (res.code == 0 || res.code == 200) {
          console.log(res)
          this.isOpen = (res.data[0].status == 'open')
        } else {
          this.$message({
            type: 'error',
            message: res.msg
          })
        }
      }).catch(err => {
        console.error(err)
      })
    },
    closeDialog() {
      this.$emit('close')
    }
  }
}
</script>

<style scoped>

</style>
