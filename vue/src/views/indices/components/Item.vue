<template>
  <div>
    <el-row
      :style="'padding-left: ' + (newdeep - 1) * 15 + 'px; margin-top: 5px'"
    >
      <div class="item" v-if="rootFlag">
        <div
          :class="
            size == 'mini'
              ? 'item-control-mini'
              : size == 'small'
              ? 'item-control-small'
              : 'item-control-medium'
          "
          @click="hidden = !hidden"
        >
          <el-link :underline="false" v-if="item.childs != null">
            <i
              :class="
                hidden
                  ? 'el-icon-caret-right'
                  : 'el-icon-caret-bottom el-icon-blue'
              "
            ></i>
          </el-link>
          <i v-else style="margin-left: 14px"></i>
        </div>
        <div class="item-cell" @click="hidden = !hidden">
          <el-link :underline="false">
            <el-tag
              effect="plain"
              :class="
                size == 'mini'
                  ? 'item-name-mini el-tag'
                  : size == 'small'
                  ? 'item-name-small el-tag'
                  : 'item-name-medium el-tag'
              "
            >
              <span>{{
                getName(item.key) ? getName(item.key) : item.key
              }}</span>
            </el-tag>
          </el-link>
        </div>
        <div
          class="item-key item-cell"
          :style="'width: ' + (300 - (newdeep - 1) * 15) + 'px'"
        >
          <el-autocomplete
            clearable
            :fetch-suggestions="querySearch"

            v-model="item.key"
            :placeholder="parent.type == 'Array' ? '' : $t('请输入键')"
            :size="size"
            :disabled="item.isRoot || parent.type == 'Array' ? true : false"
            :style="'width: ' + (290 - (newdeep - 1) * 15) + 'px'"
          >
            <i
              class="el-icon-edit el-input__icon"
              slot="suffix"
            >
            </i>
            <template slot-scope="{ item }">
              <span>{{ item.value }}-{{ item.data }}</span>
            </template>

          </el-autocomplete>
        </div>
        <div class="item-type item-cell">
          <el-select
            v-model="item.type"
            :size="size"
            :placeholder="$t('请选择')"
            @change="changeSelect"
            class="select-body"
          >
            <el-option
              v-for="type in item.isRoot ? rootOptions : options"
              :key="type.value"
              :label="type.label"
              :value="type.value"
            >
            </el-option>
          </el-select>
        </div>
        <div class="item-value item-cell">
          <el-autocomplete
            :fetch-suggestions="querySearch2"
            clearable
            :size="size"
            v-model="item.value"
            v-if="item.type != 'Number' && item.type != 'Boolean'"
            :placeholder="
              item.type == 'Array' || item.type == 'Object' ? '' : $t('请输入值')
            "
            :disabled="
              item.type == 'Array' || item.type == 'Object' ? true : false
            "
            :class="
              size == 'mini'
                ? 'el-width-mini'
                : size == 'small'
                ? 'el-width-small'
                : 'el-width-medium'
            "
          >
            <i
              class="el-icon-edit el-input__icon"
              slot="suffix"
              >
            </i>
            <template slot-scope="{ item }">
              <span>{{ item.value }}-{{ item.data}}</span>
            </template>

          </el-autocomplete>
          <el-input-number
            v-model="item.value"
            v-else-if="item.type == 'Number'"
            :size="size"
            :class="
              size == 'mini'
                ? 'el-width-mini'
                : size == 'small'
                ? 'el-width-small'
                : 'el-width-medium'
            "
          ></el-input-number>
          <el-radio-group
            v-model="item.value"
            v-else
            :class="
              size == 'mini'
                ? 'el-width-mini'
                : size == 'small'
                ? 'el-width-small'
                : 'el-width-medium'
            "
          >
            <el-radio class="el-radio" :label="true">是</el-radio>
            <el-radio class="el-radio" :label="false">否</el-radio>
          </el-radio-group>
        </div>
        <div
          :class="
            size == 'mini'
              ? 'item-control-mini'
              : size == 'small'
              ? 'item-control-small'
              : 'item-control-medium'
          "
        >
          <el-tooltip
            class="item-control-cell"
            content="添加子元素"
            placement="top"
            v-if="item.type == 'Array' || item.type == 'Object'"
          >
            <el-link :underline="false" @click="addItem()"
            ><i class="el-icon-plus el-icon-blue"></i
            ></el-link>
          </el-tooltip>
          <el-popconfirm
            :title="$t('确定删除当前节点吗？')"
            @confirm="delItem"
            v-if="!item.isRoot"
          >
            <el-link
              slot="reference"
              :underline="false"
              class="item-control-cell"
            ><i class="el-icon-close el-icon-dim"></i
            ></el-link>
          </el-popconfirm>
        </div>
      </div>
    </el-row>
    <div :style="hidden ? 'display: none' : 'display: block'">
      <template v-if="item.childs && item.type == 'Object'">
        <span v-for="(child, index) in item.childs" :key="index">
          <VueJsonItem
            :size="size"
            :item="child"
            :parent="item"
            :names="names"
            :deep="newdeep"
            :openFlag="openFlag"
          />
        </span>
      </template>
      <template v-if="item.childs && item.type == 'Array'">
        <span v-for="(child, index) in item.childs" :key="index">
          <VueJsonItem
            :size="size"
            :item="child"
            :parent="item"
            :index="index"
            :names="names"
            :deep="newdeep"
            :openFlag="openFlag"
          />
        </span>
      </template>
    </div>
  </div>
</template>

<script>
  import {filterData} from '@/utils/table'

  export default {
    name: "VueJsonItem",
    props: {
      item: {
        type: Object,
        default: {
          key: "",
          value: "",
          type: "",
        },
      },
      names: {
        type: Array,
      },
      parent: {
        type: Object,
      },
      index: {
        type: Number,
        default: 0,
      },
      deep: {
        type: Number,
        default: 0,
      },
      size: {
        type: String,
        default: "small",
      },
      rootFlag: {
        type: Boolean,
        default: true,
      },
      openFlag: {
        type: Boolean,
        default: true,
      },
    },
    data() {
      return {
        hidden: !this.openFlag,
        newdeep: this.deep + 1,
        options: [
          {value: "String", labal: "String"},
          {value: "Object", labal: "Object"},
          {value: "Array", labal: "Array"},
          {value: "Number", labal: "Number"},
          {value: "Boolean", labal: "Boolean"},
        ],
        rootOptions: [
          {value: "Object", labal: "Object"},
          {value: "Array", labal: "Array"},
        ],
        max: 8,
        queryData:[
          {'value': "type", 'data': " 数据类型"},
          {'value': "format", 'data': " 时间格式化"},
          {'value': "analyzer", 'data': " 分词器"},
          {'value': "normalizer", 'data': " 分析器"},
          {'value': "boost", 'data': " 权重"},
          {'value': "coerce", 'data': " 强制类型转换"},

          {'value': "copy_to", 'data': " 合并参数"},
          {'value': "doc_values", 'data': " 文档值"},
          {'value': "dynamic", 'data': " 动态设置"},
          {'value': "enabled", 'data': " 是否开启字段"},
          {'value': "fielddata", 'data': " 字段数据"},
          {'value': "ignore_above", 'data': " 字段保存最大长度"},
          {'value': "ignore_malformed", 'data': " 忽略格式不对的数据"},
          {'value': "include_in_all", 'data': " _all 查询包含字段"},
          {'value': "index_options", 'data': " 索引设置"},
          {'value': "index", 'data': " 是否可以被搜索"},
          {'value': "fields", 'data': " 字段"},

          {'value': "norms", 'data': " 标准信息"},
          {'value': "null_value", 'data': " 空值"},
          {'value': "position_increment_gap", 'data': " 短语位置间隙"},
          {'value': "properties", 'data': " 属性"},
          {'value': "search_analyzer", 'data': " 搜索分析器"},
          {'value': "similarity", 'data': " 匹配算法"},
          {'value': "store", 'data': " 字段是否被存储"},
          {'value': "term_vector", 'data': " 词根信息"},
        ],
        queryData2:[
          { 'value': 'text', 'data': '字符串类型(可分词) ' },
          { 'value': 'keyword', 'data': '字符串类型(不可分词) ' },
          { 'value': 'long', 'data': '数字类型 ' },
          { 'value': 'integer', 'data': '数字类型 ' },
          { 'value': 'short', 'data': '数字类型 ' },
          { 'value': 'byte', 'data': '数字类型 ' },
          { 'value': 'double', 'data': '数字类型 ' },
          { 'value': 'float', 'data': '数字类型 ' },
          { 'value': 'half_float', 'data': '数字类型 ' },
          { 'value': 'scaled_float', 'data': '数字类型 ' },
          { 'value': 'date', 'data': '时间类型 ' },
          { 'value': 'boolean', 'data': '布尔类型 ' },
          { 'value': 'binary', 'data': '二进制类型 ' },
          { 'value': 'Array', 'data': '数组类型 ' },
          { 'value': 'object', 'data': '对象类型  ' },
          { 'value': 'nested', 'data': '嵌套类型 ' },
          { 'value': 'geo_point', 'data': '地理类型 ' },
          { 'value': 'geo_shape', 'data': '多边形类型 ' },
          { 'value': 'ip', 'data': 'ip类型 ' },
          { 'value': 'completion', 'data': '补全类型 ' },
          { 'value': 'token_count', 'data': '令牌计数类型 ' },
          { 'value': 'yyyy-MM-dd HH:mm:ss', 'data': '时间格式化 ' },
        ]
      };
    },
    methods: {
      querySearch(queryString, cb) {

        let queryData = JSON.parse(JSON.stringify(this.queryData))
        if(queryString == undefined)queryString = ""
        if (queryString.trim() == '') {
          if (queryData.length > this.max) {
            cb(queryData.slice(0, this.max))
          } else {
            cb(queryData)
          }
          return;
        }


        queryData = filterData(queryData, queryString.trim())

        if (queryData.length > this.max) {
          cb(queryData.slice(0, this.max))
        } else {
          cb(queryData)
        }
      },
      querySearch2(queryString, cb) {
        let queryData = JSON.parse(JSON.stringify(this.queryData2))
        if(queryString == undefined)queryString = ""
        if (queryString.trim() == '') {
          if (queryData.length > this.max) {
            cb(queryData.slice(0, this.max))
          } else {
            cb(queryData)
          }
          return;
        }

        queryData = filterData(queryData, queryString.trim())

        if (queryData.length > this.max) {
          cb(queryData.slice(0, this.max))
        } else {
          cb(queryData)
        }
      },
      /**增加元素 */
      addItem() {
        var childs = this.item.childs;
        var type = this.item.type;
        var count = 0;
        if (childs[childs.length - 1] != undefined) {
          count = parseInt(childs[childs.length - 1].key) + 1;
        }
        var additem = undefined;
        if (type != "Array") {
          additem = {type: "String", childs: null, value: null};
        } else {
          additem = {
            key: count,
            type: "String",
            childs: null,
            value: null,
          };
        }
        this.$set(childs, childs.length, additem);
      },
      /**删除元素 */
      delItem() {
        var childs = this.parent.childs;
        var item = this.item;
        for (var i in childs) {
          if (childs[i] == item) {
            childs.splice(i, 1);
          }
        }
      },
      /**更改类型 */
      changeSelect(option) {
        if (option == "Array" || option == "Object") {
          this.item.childs = [];
          this.item.value = null;
        } else if (option == "Number") {
          this.item.childs = null;
          this.item.value = 0;
        } else if (option == "Boolean") {
          this.item.childs = null;
          this.item.value = true;
        }
      },
      /**判断是否为空 */
      isNull(e) {
        var flag = false;
        if (e == null || e == "" || e == undefined) {
          flag = true;
        }
        if (e == 0) {
          flag = false;
        }
        return flag;
      },
      /**通过key获取中文 */
      getName(key) {
        var n = undefined;
        var names = this.names;
        for (var i in names) {
          var name = names[i];
          if (name.key == key) {
            n = name.name;
            break;
          }
        }
        if (n != undefined) {
          return n;
        } else {
          return "";
        }
      },
    },
  };
</script>

<style scoped>

  .item {
    display: flex;
  }

  .item-control-mini {
    height: 26px;
    line-height: 24px;
  }

  .item-control-small {
    height: 30px;
    line-height: 28px;
  }

  .item-control-medium {
    height: 32px;
    line-height: 30px;
  }

  .item-control-cell {
    margin-left: 10px;
    font-size: 16px;
  }

  .item-cell {
    padding-left: 10px;
  }

  .el-icon-blue {
    color: #409eff;
  }

  .el-icon-green {
    color: #67c23a;
  }

  .el-icon-dim {
    color: #909399;
  }

  .item-key {
    width: 150px;
  }

  .item-type {
    width: 100px;
  }

  .el-tag {
    text-align: center;
    width: 100px;
  }

  .el-width-medium {
    width: 208.67px;
    height: 34px;
    line-height: 34px;
    text-align: center;
  }

  .el-width-small {
    width: 193.33px;
    height: 30px;
    line-height: 30px;
    text-align: center;
  }

  .el-width-mini {
    width: 178px;
    height: 26px;
    line-height: 26px;
    text-align: center;
  }

  .el-radio {
    height: 32px;
    line-height: 32px;
  }

  .item-name-medium {
    height: 35px;
    line-height: 33px;
  }

  .item-name-small {
    height: 32px;
    line-height: 30px;
  }

  .item-name-mini {
    height: 28px;
    line-height: 26px;
  }

  @media screen and (max-width: 700px) {
    .select-body {
      width: 100px;
    }

    .input-body {
      width: 130px;
    }
  }
</style>
