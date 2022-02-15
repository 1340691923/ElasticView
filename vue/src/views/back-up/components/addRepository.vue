<template>
  <div>
    <el-dialog :close-on-click-modal="false" :visible.sync="open" title="新增/修改快照存储库" @close="closeDialog">
      <el-card class="box-card">
        <el-form label-width="500px" label-position="left">
          <el-form-item label="存储库名">
            <el-input v-model="form.name" placeholder="存储库名" />
          </el-form-item>
          <el-form-item label="类型（type）">
            <el-select v-model="form.type" filterable>
              <el-option label="fs" value="fs" />
              <el-option label="url" value="url" />
            </el-select>
          </el-form-item>
          <el-form-item label="存储位置（location/url）">
            <el-input
              v-model="form.location"
              type="textarea"
              :autosize="{ minRows: 2, maxRows: 4}"
              placeholder="存储位置（location/url）"
            />
          </el-form-item>
          <el-form-item label="是否压缩 (compress)">
            <el-select v-model="form.compress" clearable filterable>
              <el-option label="是" value="true" />
              <el-option label="否" value="false" />
            </el-select>
          </el-form-item>
          <el-form-item label="大文件分解块大小">

            <el-input v-model="form.chunk_size" placeholder="1GB，10MB，5KB，500B。默认为null（无限制的块大小）" />
          </el-form-item>
          <el-form-item label="是否只读">
            <el-select v-model="form.readonly" filterable>
              <el-option label="是" value="true" />
              <el-option label="否" value="false" />
            </el-select>
          </el-form-item>
          <el-form-item label="节点恢复速率">
            <el-input v-model="form.max_restore_bytes_per_sec" placeholder="节点恢复速率" />
          </el-form-item>
          <el-form-item label="每个节点快照速率">
            <el-input v-model="form.max_snapshot_bytes_per_sec" placeholder="每个节点快照速率" />
          </el-form-item>
        </el-form>
        <div style="text-align:right;">
          <el-button type="danger" icon="el-icon-close" @click="closeDialog">取消</el-button>
          <el-button type="primary" icon="el-icon-check" @click="confirm">确认</el-button>
        </div>
      </el-card>
    </el-dialog>
  </div>
</template>

<script>
import { SnapshotCreateRepositoryAction } from '@/api/es-backup'

export default {
  name: 'Add',
  components: {},
  props: {
    open: {
      type: Boolean,
      default: false
    },
    snapshotData: {
      type: Object,
      default: null
    }
  },
  data() {
    return {
      isOpen: false,
      form: {
        name: '',
        type: 'fs',
        location: '',
        compress: 'false',
        max_restore_bytes_per_sec: '40mb',
        max_snapshot_bytes_per_sec: '40mb',
        chunk_size: '',
        readonly: 'false'
      }
    }
  },
  computed: {},

  created() {
    if (this.snapshotData != null) {
      this.form = Object.assign({}, this.snapshotData)
    }
  },
  methods: {
    closeDialog() {
      this.$emit('close', false)
    },
    async confirm() {
      const input = this.form
      input['es_connect'] = this.$store.state.baseData.EsConnectID
      const { code, data, msg } = await SnapshotCreateRepositoryAction(input)
      if (code == 0) {
        this.$emit('close', true)
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
    }
  }
}
</script>

<style scoped>

</style>
