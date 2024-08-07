<template>
  <div class="right-menu-item">
    <el-select
      v-model="Esi18n"
      style="width:60px"
      filterable
      default-first-option
      placeholder="I18N"
      @change="changeI18n"
    >
      <el-option value="zh" label="中文" />
      <el-option value="en" label="English" />
    </el-select>

    <el-select
      v-model="linkID"
      style="width:100px"
      filterable
      default-first-option
      :placeholder="$t('请选择ES连接')"
      @change="change"
    >
      <!-- <el-option :value="Number(0)" :label="$t('请选择ES连接')" /> -->
      <el-option v-for="item in opt" :key="item.id" :value="Number(item.id)" :label="item.remark" />
    </el-select>
    <el-button
      size="mini"
      type="primary"

      @click="refresh"
    > {{ $t('刷新') }}
    </el-button>
    <slot name="avatar" />
  </div>
</template>

<script>
import { OptAction } from '@/api/es-link'
import { setI18nLanguage } from '../../utils/lang'

export default {
  inject: ['reload'],
  name: 'SelectLink',
  data() {
    return {
      opt: [],
      linkID: '',
      Esi18n: localStorage.getItem('lang') || 'zh',
      time: null,
      timeSecend: 60
    }
  },
  computed: {},
  watch: {},
  mounted() {
    const obj = this.$store.state.baseData.EsConnectID
    this.linkID = Number(obj)
    this.getEsOpt()
    this.startLoop()
  },

  beforeDestroy() {
    // 清除定时器
    clearInterval(this.time)
    this.time = null
    console.log('beforeDestroy')
  },
  methods: {
    startLoop() {
      this.time = setInterval(() => {
        this.getEsOpt()
      }, this.timeSecend * 1000)
    },
    async getEsOpt() {
      const res = await OptAction({ 'getByLocal': 1 })
      if (res.data == null) res.data = []
      this.opt = res.data

    },
    refresh() {
      this.getEsOpt()
      this.$message({
        type: 'success',
        message: '刷新ES连接信息成功'
      })
    },
    change(link) {
      localStorage.setItem('EsConnect', link)
      this.$store.dispatch('baseData/SetEsConnect', link)
      this.reload()
    },
    changeI18n(i18n) {
      setI18nLanguage(i18n)
      this.reload()
    }
  }

}
</script>

<style lang="scss" scoped>
.header-search {
  font-size: 0 !important;

.search-icon {
  cursor: pointer;
  font-size: 18px;
  vertical-align: middle;
}

.header-search-select {
  font-size: 18px;
  transition: width 0.2s;
  width: 0;
  overflow: hidden;
  background: transparent;
  border-radius: 0;
  display: inline-block;
  vertical-align: middle;

  ::v-deep .el-input__inner {
  border-radius: 0;
  border: 0;
  padding-left: 0;
  padding-right: 0;
  box-shadow: none !important;
  border-bottom: 1px solid #d9d9d9;
  vertical-align: middle;
}

}

&
.show {

.header-search-select {
  width: 210px;
  margin-left: 10px;
}

}
}
</style>
