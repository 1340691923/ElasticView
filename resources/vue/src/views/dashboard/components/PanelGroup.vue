<template>
  <el-row :gutter="40" class="panel-group">
    <el-col :xs="12" :sm="12" :lg="8" class="card-panel-col">
      <div class="card-panel" @click="handleSetLineChartData('shared')">
        <div class="card-panel-icon-wrapper icon-people">
          <svg-icon icon-class="form" class-name="card-panel-icon" />
        </div>
        <div class="card-panel-description">
          <div class="card-panel-text">
            {{ $t('总分片数') }}
          </div>
          <span v-loading="loading4" class="card-panel-num">
            {{ count.allShared }}
          </span>
        </div>
      </div>
    </el-col>
    <el-col :xs="12" :sm="12" :lg="8" class="card-panel-col">
      <div class="card-panel" @click="handleSetLineChartData('shared')">
        <div class="card-panel-icon-wrapper icon-message">
          <svg-icon icon-class="edit" class-name="card-panel-icon" />
        </div>
        <div class="card-panel-description">
          <div class="card-panel-text">
            {{ $t('成功的分片数') }}
          </div>
          <span v-loading="loading4" class="card-panel-num">
            {{ count.successfulShared }}
          </span>
        </div>
      </div>
    </el-col>
    <el-col :xs="12" :sm="12" :lg="8" class="card-panel-col">
      <div class="card-panel" @click="handleSetLineChartData('shoppings')">
        <div class="card-panel-icon-wrapper icon-shopping">
          <svg-icon style="color: red" icon-class="edit" class-name="card-panel-icon" />
        </div>
        <div class="card-panel-description">
          <div class="card-panel-text">
            {{ $t('失败的分片数') }}
          </div>
          <span v-loading="loading4" class="card-panel-num">
            {{ count.failedShared }}
          </span>
        </div>
      </div>
    </el-col>
    <el-col :xs="12" :sm="12" :lg="8" class="card-panel-col">
      <div class="card-panel" @click="handleSetLineChartData('purchases')">
        <div class="card-panel-icon-wrapper icon-message">
          <svg-icon icon-class="tree" class-name="card-panel-icon" />
        </div>
        <div class="card-panel-description">
          <div class="card-panel-text">
            {{ $t('索引数') }}
          </div>
          <span v-loading="loading1" class="card-panel-num">
            {{ count.index }}
          </span>
        </div>

      </div>
    </el-col>

    <el-col :xs="12" :sm="12" :lg="8" class="card-panel-col">
      <div class="card-panel" @click="handleSetLineChartData('shoppings')">
        <div class="card-panel-icon-wrapper icon-shopping">
          <svg-icon icon-class="tree-table" class-name="card-panel-icon" />
        </div>
        <div class="card-panel-description">
          <div class="card-panel-text">
            {{ $t('文档数') }}
          </div>
          <span v-loading="loading3" class="card-panel-num">
            {{ count.document }}
          </span>
        </div>
      </div>
    </el-col>
    <el-col :xs="12" :sm="12" :lg="8" class="card-panel-col">
      <div class="card-panel" @click="handleSetLineChartData('shoppings')">
        <div class="card-panel-icon-wrapper icon-shopping">
          <svg-icon icon-class="dashboard" class-name="card-panel-icon" />
        </div>
        <div class="card-panel-description">
          <div class="card-panel-text">
            {{ $t('索引所占空间大小') }}
          </div>
          <span v-loading="loading2" class="card-panel-num">
            {{ count.size }}
          </span>

        </div>
      </div>
    </el-col>
  </el-row>
</template>

<script>
import { CatAction } from '@/api/es'
import { IndexsCountAction } from '@/api/es-index'
export default {
  components: {
    'CountTo': () => import('vue-count-to')
  },
  data() {
    return {
      addUserCount: 0,
      count: {
        allShared: 0,
        successfulShared: 0,
        index: 0,
        failedShared: 0,
        document: 0,
        size: 0
      },
      loading1: false,
      loading2: false,
      loading3: false,
      loading4: false
    }
  },
  mounted() {
    this.getIndexCount()
    this.getSegments()
    this.catStats()
    this.catAllocation()
  },
  methods: {
    handleSetLineChartData(type) {
      // this.$emit('handleSetLineChartData', type)
    },

    async getIndexCount() {
      const form = {
        // cat: 'CatIndices',
        es_connect: this.$store.state.baseData.EsConnectID
        // index_bytes_format: 'kb'
      }
      this.loading1 = true
      const { data, code, msg } = await IndexsCountAction(form)
      this.loading1 = false
      if (code == 0) {
        this.count.index = data
      }
    },
    async catAllocation() {
      const form = {
        cat: 'CatAllocation',
        es_connect: this.$store.state.baseData.EsConnectID
      }
      this.loading2 = true
      const { data, code, msg } = await CatAction(form)
      this.loading2 = false
      if (code == 0) {
        this.count.size = data[0]['disk.indices']
      }
    },
    async catStats() {
      const form = {
        cat: 'CatStats',
        es_connect: this.$store.state.baseData.EsConnectID
      }
      this.loading3 = true
      const { data, code, msg } = await CatAction(form)
      this.loading3 = false
      if (code == 0) {
        this.count.document = data.indices.docs.count
      }
    },
    async getSegments() {
      const form = {
        cat: 'CatSegments',
        es_connect: this.$store.state.baseData.EsConnectID
      }
      this.loading4 = true
      const { data, code, msg } = await CatAction(form)
      this.loading4 = false
      if (code == 0) {
        this.count.failedShared = Number(data._shards.failed)
        this.count.allShared = Number(data._shards.total)
        this.count.successfulShared = Number(data._shards.successful)
      }
    }
  }
}
</script>

<style lang="scss" scoped>
.panel-group {
  margin-top: 18px;

.card-panel-col {
  margin-bottom: 32px;
}

.card-panel {
  height: 108px;
  cursor: pointer;
  font-size: 12px;
  position: relative;
  overflow: hidden;
  color: #666;
  background: #fff;
  box-shadow: 4px 4px 40px rgba(0, 0, 0, .05);
  border-color: rgba(0, 0, 0, .05);

&
:hover {

.card-panel-icon-wrapper {
  color: #fff;
}

.icon-people {
  background: #40c9c6;
}

.icon-message {
  background: #36a3f7;
}

.icon-money {
  background: #f4516c;
}

.icon-shopping {
  background: #34bfa3
}

}

.icon-people {
  color: #40c9c6;
}

.icon-message {
  color: #36a3f7;
}

.icon-money {
  color: #f4516c;
}

.icon-shopping {
  color: #34bfa3
}

.card-panel-icon-wrapper {
  float: left;
  margin: 14px 0 0 14px;
  padding: 16px;
  transition: all 0.38s ease-out;
  border-radius: 6px;
}

.card-panel-icon {
  float: left;
  font-size: 48px;
}

.card-panel-description {
  float: right;
  font-weight: bold;
  margin: 26px;
  margin-left: 0px;

.card-panel-text {
  line-height: 18px;
  color: rgba(0, 0, 0, 0.45);
  font-size: 16px;
  margin-bottom: 12px;
}

.card-panel-num {
  font-size: 20px;
}

}
}
}

@media (max-width: 550px) {
  .card-panel-description {
    display: none;
  }

  .card-panel-icon-wrapper {
    float: none !important;
    width: 100%;
    height: 100%;
    margin: 0 !important;

  .svg-icon {
    display: block;
    margin: 14px auto !important;
    float: none !important;
  }
}

}
</style>
