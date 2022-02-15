<template>
  <div>
    <el-dialog :close-on-click-modal="false" :visible.sync="open" title="创建快照" @close="closeDialog">
      <el-card class="box-card">
        <el-form label-width="500px" label-position="left">
          <el-form-item label="仓库名">
            <el-select
              v-model="form.repositoryName"
              clearable
              filterable
              placeholder="请选择存储库"
            >
              <el-option
                v-for="(v,k,index) of repositoryNameList"
                :key="index"
                :label="k"
                :value="k"
              />
            </el-select>
          </el-form-item>
          <el-form-item label="快照名">
            <el-input v-model="form.snapshotName" placeholder="快照名" />
          </el-form-item>
          <el-form-item label="需要备份的索引">
            <index-select :multiple="true" :have-all="true" :clearable="true" placeholder="迁移别名到多个索引上" @change="changeIndex" />
          </el-form-item>

          <el-form-item label="ignore_unavailable   【把这个选项设置为 true 的时候在创建快照的过程中会忽略不存在的索引,如果没有设置ignore_unavailable，在索引不存在的情况下快照请求将会失败。】">
            <el-select v-model="form.ignore_unavailable" filterable>
              <el-option label="不设置" :value="null" />
              <el-option label="是" :value="true" />
              <el-option label="否" :value="false" />
            </el-select>
          </el-form-item>
          <el-form-item label="include_global_state 【通过设置 include_global_state 为false 能够防止 集群的全局状态被作为快照的一部分存储起来】">
            <el-select v-model="form.include_global_state" filterable>
              <el-option label="不设置" :value="null" />
              <el-option label="是" :value="true" />
              <el-option label="否" :value="false" />
            </el-select>
          </el-form-item>
          <el-form-item label="partial  【如果快照中的1个或多个索引不是全部主分片都可用，就会导致整个创建快照的过程失败。 通过设置 partial 为 true 可以改变这个行为】">
            <el-select v-model="form.partial" clearable filterable>
              <el-option label="不设置" :value="null" />
              <el-option label="是" :value="true" />
              <el-option label="否" :value="false" />
            </el-select>
          </el-form-item>
          <el-form-item label="创建方式">
            <el-select v-model="form.wait" clearable filterable>
              <el-option label="不设置" :value="null" />
              <el-option label="异步创建" :value="true" />
              <el-option label="同步创建" :value="false" />
            </el-select>
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
import { CreateSnapshotAction, SnapshotRepositoryListAction } from '@/api/es-backup'

export default {
  name: 'Add',
  components: {},
  components: {
    'IndexSelect': () => import('@/components/index/select')
  },
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
        snapshotName: '',
        repositoryName: '',
        indexList: [],
        ignore_unavailable: null,
        include_global_state: null,
        partial: null,
        wait: null
      },
      repositoryNameList: {}
    }
  },
  computed: {},

  created() {
    this.initRepositoryName()
  },
  methods: {
    changeIndex(index) {
      console.log(index)
      this.form.indexList = []
      this.form.indexList = index
    },
    async initRepositoryName() {
      this.loading = true
      const input = {}
      input['es_connect'] = this.$store.state.baseData.EsConnectID
      const { data, code, msg } = await SnapshotRepositoryListAction(input)
      if (code == 0) {
        this.repositoryNameList = data.res
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
    closeDialog() {
      this.$emit('close', false)
    },
    async confirm() {
      const input = this.form
      input['es_connect'] = this.$store.state.baseData.EsConnectID

      const { code, data, msg } = await CreateSnapshotAction(input)
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
