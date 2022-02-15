<template>
  <div :style="styles" class="float">
    <span v-if="!simple" style="text-align: center;display:block;" class="font1">
      <el-button size="small" icon="el-icon-refresh" @click="format">美化</el-button>
      <el-button
        v-clipboard:copy="value"
        v-clipboard:success="onCopy"
        v-clipboard:error="onError"
        size="small"
        icon="el-icon-document-copy"
      >点击复制</el-button>
    </span>
    <editor
      v-model="value"
      lang="mysql"
      :options="editorOptions"
      theme="chrome"
      height="700"
      @init="editorInit"
    />
  </div>
</template>
<script>
import { format } from 'sql-formatter'
export default {
  name: 'Sql',
  props: {
    styles: {
      type: String,
      default: 'width: 50%'
    },
    read: {
      type: Boolean,
      default: false
    },
    title: {
      type: String,
      default: ''
    },
    value: {
      type: Object,
      default: {}
    },
    pointOut: {
      type: Array,
      default: []
    },
    simple: {
      type: Boolean,
      default: false
    }
  },
  data() {
    return {
      ed: null,
      editorOptions: {
        // 设置代码编辑器的样式
        enableBasicAutocompletion: true,
        enableSnippets: true,
        enableLiveAutocompletion: true,
        tabSize: 2,
        fontSize: 25,
        showPrintMargin: false // 去除编辑器里的竖线
      }
    }
  },
  watch: {
    value(v) {
      this.$emit('getValue', v)
    }
  },
  destroyed() {
    this.ed = null
  },
  created() {

  },
  methods: {
    format() {
      let sqlContent = ''
      sqlContent = this.ed.getValue()
      this.ed.setValue(format(sqlContent))
    },
    editorInit: function(ed) {
      require('brace/theme/chrome')
      require('brace/ext/language_tools')
      require('brace/mode/yaml')
      require('brace/ext/searchbox')
      require('brace/ext/emmet')
      require('brace/theme/monokai')
      require('brace/mode/mysql')
      require('brace/mode/less')
      require('brace/snippets/mysql')
      if (this.read) {
        ed.setReadOnly(true)
      }
      this.ed = ed
    },
    onCopy(e) { 　　 // 复制成功
      this.$message({
        message: '复制成功！',
        type: 'success'
      })
    },
    onError(e) {　　 // 复制失败
      this.$message({
        message: '复制失败！',
        type: 'error'
      })
    }
  }
}
</script>
<style scoped>
  .float {
    float: left;
  }

  .font1 {
    font-family: "Helvetica Neue", Helvetica, "PingFang SC", "Hiragino Sans GB", "Microsoft YaHei", "微软雅黑", Arial, sans-serif;
    color: green;
  }
</style>
