<template>
  <el-select
    v-model="currValue"
    size="mini"
    multiple
    allow-create
    placeholder="请选择"
    style="width: 150px"
    class="yt-select"
    filterable
    v-bind="$attrs"
    :filter-method="userFilter"
    :disabled="disabled"
    :clearable="clearable"
    @change="change"
  >
    <el-option
      v-for="option in renderOption"
      :key="option.value"
      :value="option.value"
      :label="option.label"
    >{{ option.label }}
    </el-option>
  </el-select>

</template>

<script>

export default {
  name: 'SmallSelect',
  props: {
    value: {
      type: [String, Number],
      default: ''
    },
    max: {
      type: Number,
      default: 30
    },
    disabled: {
      type: Boolean,
      default: false
    },
    clearable: {
      type: Boolean,
      default: true
    },
    options: {
      type: Array,
      default: () => []
    }
  },
  data() {
    return {
      renderOption: []
    }
  },
  computed: {
    currValue: {
      get() {
        return this.value || ''
      },
      set(value) {
        this.$emit('input', value)
      }
    }
  },
  watch: {
    value() {
      this.addValueOptions()
    },
    options: {
      handler(V) {
        this.init()
      },
      deep: true
    }
  },
  created() {
    this.init()
  },
  methods: {
    async init() {
      this.userFilter()
      this.addValueOptions()
    },
    addValueOptions() {
      if (this.currValue) {
        const target = this.options.find((item) => { // 从大option中找到当前条
          return item.value === this.currValue
        })
        if (target) { // 将当前条与小option比对，没有则加入
          if (this.renderOption.every(item => item.value !== target.value)) {
            this.renderOption.unshift(target)
          }
        }
      }
    },
    addFilterOptions(label) {
      // 每次查找输入时，若有精确匹配的条目，保证该条目一定在renderOption内
      const target = this.options.find((item) => { // 从大option中找到当前条
        return item.label === label
      })
      if (target) { // 将当前条与小option比对，没有则加入
        if (this.renderOption.every(item => item.label !== target.label)) {
          this.renderOption.unshift(target)
        }
      }
    },
    userFilter(query = '') {
      let arr = []
      if (query != '') {
        arr = this.options.filter((item) => {
          return item.label.includes(query) || item.value.includes(query)
        })
      } else {
        arr = JSON.parse(JSON.stringify(this.options))
      }

      if (arr.length > this.max) {
        this.renderOption = arr.slice(0, this.max)
        this.addFilterOptions(query)
      } else {
        this.renderOption = arr
      }
    },
    change(value) {
      this.$emit('change', value)
      if (!value) { // 单选清空-optons初始化下
        this.userFilter()
      }
    }
  }
}
</script>
