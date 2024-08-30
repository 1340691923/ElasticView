<template>
  <div>
    <el-dialog :close-on-click-modal="false" :visible.sync="open" :title="$t('恢复备份')" @close="closeDialog">
      <el-card class="box-card">
        <el-form label-width="500px" label-position="left">
          <el-form-item :label="$t('仓库名')">
            <el-input v-model="repository" readonly />
          </el-form-item>
          <el-form-item :label="$t('快照名')">
            <el-input v-model="snapshot" readonly />
          </el-form-item>
          <el-form-item :label="$t('需要恢复的索引')">
            <index-select
              :multiple="true"
              :have-all="true"
              :clearable="true"
              :placeholder="$t('选择需要恢复的索引')"
              @change="changeIndex"
              @mount="initIndexSelected"
              ref="indexSelect"
            />
          </el-form-item>
          <el-form-item :label="$t('rename_pattern 【正则表达式】')">
            <el-input v-model="form.rename_pattern" />
          </el-form-item>
          <el-form-item :label="$t('rename_replacement 【正则表达式】')">
            <el-input v-model="form.rename_replacement" />
          </el-form-item>

          <!--          <el-form-item
            :label="$t('ignore_unavailable   【把这个选项设置为 true 的时候在创建快照的过程中会忽略不存在的索引,如果没有设置ignore_unavailable，在索引不存在的情况下快照请求将会失败。】')"
          >
            <el-select v-model="form.ignore_unavailable" filterable>
              <el-option :label="$t('不设置')" :value="null" />
              <el-option :label="$t('是')" :value="true" />
              <el-option :label="$t('否')" :value="false" />
            </el-select>
          </el-form-item>-->
          <el-form-item
            :label="$t('include_global_state 【通过设置 include_global_state 为false 能够防止 集群的全局状态被作为快照的一部分存储起来】')"
          >
            <el-select v-model="form.include_global_state" filterable>
              <el-option :label="$t('不设置')" :value="null" />
              <el-option :label="$t('是')" :value="true" />
              <el-option :label="$t('否')" :value="false" />
            </el-select>
          </el-form-item>
          <el-form-item :label="$t('partial  【如果快照中的1个或多个索引不是全部主分片都可用，就会导致整个创建快照的过程失败。 通过设置 partial 为 true 可以改变这个行为】')">
            <el-select v-model="form.partial" clearable filterable>
              <el-option :label="$t('不设置')" :value="null" />
              <el-option :label="$t('是')" :value="true" />
              <el-option :label="$t('否')" :value="false" />
            </el-select>
          </el-form-item>
          <el-form-item :label="$t('创建方式')">
            <el-select v-model="form.wait" clearable filterable>
              <el-option :label="$t('不设置')" :value="null" />
              <el-option :label="$t('异步创建')" :value="true" />
              <el-option :label="$t('同步创建')" :value="false" />
            </el-select>
          </el-form-item>
        </el-form>
        <div style="text-align:right;">
          <el-button
            size="mini"
            type="danger"
            icon="el-icon-close"
            @click="closeDialog"
          >{{ $t('取消') }}
          </el-button>
          <el-button
            size="mini"
            type="primary"
            icon="el-icon-check"
            @click="confirm"
          >{{ $t('确认') }}
          </el-button>
        </div>
      </el-card>
    </el-dialog>
  </div>
</template>

<script>
import { SnapshotRestoreAction,SnapshotDetailAction } from '@/api/es-backup'
import { OptimizeAction } from '@/api/es'

export default {
  name: 'Add',
  components: {
    'IndexSelect': () => import('@/components/index/select')
  },
  props: {
    open: {
      type: Boolean,
      default: false
    },
    snapshot: {
      type: String,
      default: ''
    },
    repository: {
      type: String,
      default: ''
    }
  },
  data() {
    return {
      isOpen: false,
      form: {
        snapshotName: this.snapshot,
        repositoryName: this.repository,
        ignore_unavailable: null,
        include_global_state: null,
        partial: null,
        wait: null,
        indexList: [],
        rename_pattern: '',
        rename_replacement: ''
      }
    }
  },
  computed: {},
  created() {

  },

  methods: {
    changeIndex(index) {
      this.form.indexList = []
      this.form.indexList = index
    },
    closeDialog() {
      this.$emit('close', false)
    },
    async confirm() {
      const input = this.form
      input['es_connect'] = this.$store.state.baseData.EsConnectID
      const { code, data, msg } = await SnapshotRestoreAction(input)
      if (code == 0) {
        this.$emit('close', true)
        this.$message({
          type: 'success',
          message: msg
        })
        return
      }
      if (msg.indexOf(' with same name already exists ') !== -1) {
        this.$confirm('错误信息:【' + msg + '】,是否先关闭所有索引再进行恢复快照?', '警告', {
          confirmButtonText: '确认',
          cancelButtonText: '取消',
          type: 'warning'
        })
          .then(async() => {
            const input = {}
            input['es_connect'] = this.$store.state.baseData.EsConnectID
            input['index_name'] = '*'
            input['command'] = 'close'
            const res = await OptimizeAction(input)
            if (res.code != 0) {
              this.$message({
                type: 'error',
                message: msg
              })
              return
            }
            this.confirm()
          })
          .catch(err => {
            console.error(err)
          })
      } else {
        this.$message({
          type: 'error',
          message: msg
        })
      }
    },
    //自动加载快照中的索引列表
    async initIndexSelected(){
      const input = {}
      input['es_connect'] = this.$store.state.baseData.EsConnectID
      input['repository'] = this.form.repositoryName
      input['snapshot'] = this.form.snapshotName

      const { data, code, msg } = await SnapshotDetailAction(input)
      if (code != 0) {
        this.$message({
          type: 'error',
          message: msg
        })
        return
      } else {
        this.form.indexList = []
        if(data.snapshots && data.snapshots.length>0){
          this.$nextTick(() => {
           this.$refs.indexSelect.setIndexList(  data.snapshots[0].indices )
          })
        }
      }
    }
  }
}
</script>

<style scoped>

</style>
