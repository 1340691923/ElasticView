<template>
  <el-row :class="borderFlag?'el-row-border':''" :style="shadowFlag?'box-shadow: 0 2px 12px 0 rgba(0, 0, 0, 0.1)':''">
      <el-col v-if="isJson"  class="item-body">
        <Item
          :item="item"
          :deep="deep"
          :parent="item"
          :names="names"
          :size="size"
          :index="0"
          :rootFlag="rootFlag"
          :openFlag="openFlag"
        />
      </el-col>
      <el-col v-else>
        <i>{{ errorResult }}</i>
      </el-col>
      <el-backtop v-if="backTopFlag"></el-backtop>
    </el-row>
</template>

<script>
import Item from "./Item";
export default {
  name: "VueJsonHelper",
  components: {
    "Item": Item,
  },
  data() {
    return {
      deep: 0,
      isJson: false,
      jsonData: undefined,
      errorResult: undefined,
      item: {
        key: "Root",
        value: null,
        type: null,
        isRoot: true,
        childs: [],
      },
    };
  },
  props: {
    names: {
      type: Array,
    },
    size: {
      type: String,
      default: 'small'
    },
    jsonStr: {
      type: String,
    },
    rootFlag: {
      type: Boolean,
      default: true
    },
    openFlag: {
      type: Boolean,
      default: true
    },
    borderFlag: {
      type: Boolean,
      default: true
    },
    shadowFlag: {
      type: Boolean,
      default: false
    },
    backTopFlag: {
      type: Boolean,
      default: false
    }
  },
  watch: {
    item: {
      handler(newVal, oldVal) {
        if (newVal == oldVal) {
          var json = this.handleJsonData(newVal);
          if (json != undefined) {
            this.$emit("jsonListener", json);
          }
        }
      },
      immediate: true,
      deep: true,
    },
  },
  created() {
    this.isJson = this.judgeJson();
    this.item.childs = this.handleJson(this.jsonData);
    this.item.type = this.handleType();
  },
  methods: {
    /**判断是否为json */
    judgeJson() {
      var flag = false;
      try {
        this.jsonData = JSON.parse(this.jsonStr);
        flag = true;
      } catch (e) {
        this.errorResult = e;
        flag = false;
      }
      return flag;
    },
    /**处理JSONData数据 */
    handleJsonData(jsonData) {
      var objs = {};
      var arr = [];
      var type = jsonData.type;
      var childs = jsonData.childs;
      for (var i in childs) {
        if (childs[i].type != "Object" && childs[i].type != "Array") {
          if (type == "Object") {
            objs[childs[i].key] = childs[i].value;
          } else if (type == "Array") {
            arr.push(childs[i].value);
          }
        } else {
          if (type == "Object") {
            objs[childs[i].key] = this.handleJsonData(childs[i]);
          } else if (type == "Array") {
            arr.push(this.handleJsonData(childs[i]));
          }
        }
      }
      if (type == "Object") {
        return objs;
      } else if (type == "Array") {
        return arr;
      }
    },
    /**处理JSON数据 */
    handleJson(json) {
      var jsonData = [];
      for (var i in json) {
        var type = this.judgeType(json[i]);
        if (type == "Object" || type == "Array") {
          var item = {
            key: i,
            value: null,
            type: type,
            childs: this.handleJson(json[i]),
          };
          jsonData.push(item);
        } else {
          var item = {
            key: i,
            value: json[i],
            type: type,
          };
          jsonData.push(item);
        }
      }
      return jsonData;
    },
    /**判断数据类型 */
    judgeType(data) {
      var type = Object.prototype.toString.call(data);
      if (type === "[object String]") {
        type = "String";
      } else if (type === "[object Number]") {
        type = "Number";
      } else if (type === "[object Boolean]") {
        type = "Boolean";
      } else if (type === "[object Array]") {
        type = "Array";
      } else if (type === "[object Object]") {
        type = "Object";
      } else {
        type = null;
      }
      return type;
    },
    /**处理根节点数据类型 */
    handleType() {
      return this.judgeType(this.jsonData);
    },
  },
};
</script>

<style scoped>

@media screen and (max-width: 700px) {
  .item-body {
    white-space: nowrap;
    overflow-x: scroll;
    overflow-y: hidden;
  }
}
.content {
  color: #888;
}
.el-row-border{
  padding: 10px;
  border: 1px solid #DCDFE6 ;
  border-radius: 10px;
}
</style>
