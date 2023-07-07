<template>
  <div>
    <div>
      <!-- <el-checkbox v-model="indexIsInput" style="margin-right: 10px;"
        >模糊搜索</el-checkbox
      > -->
      <el-select
        multiple
        v-model="indexValue"
        filterable
        placeholder="请选择索引"
        style="width: 300px;"
        v-if="!indexIsInput"
      >
        <el-option
          v-for="(v, index) in indexOptions"
          :key="index"
          :label="v.remark || v.indexName"
          :value="v.indexName"
        ></el-option>
      </el-select>
      <el-input
        v-model="inputIndex"
        style="width: 300px;"
        placeholder="请输入索引"
        v-else
      ></el-input>
      <el-button type="primary" @click="onBtnSearch" style="margin-left: 10px;"
      >搜索</el-button
      >
      <div
        v-for="(v, index) in SearchLogFilter"
        :key="index"
        style="margin-top: 10px;"
      >
        <!-- <el-tag>字段名</el-tag> -->
        <el-input
          v-model="v.search_col"
          style="width: 100px;"
          v-if="indexIsInput"
        ></el-input>
        <el-input
          :value="indexMap[v.search_col] || v.search_col"
          style="width: 100px;"
          :disabled="true"
          v-else
        ></el-input>
        <el-input
          v-model="v.search_text"
          style="width: 300px;margin-left: 10px;"
          placeholder="请输入内容"
        ></el-input>
      </div>
    </div>
    <div style="margin-top: 10px;" v-if="!isFirstSearch">
      找到与{{ !indexIsInput ? indexValue.join(",") : inputIndex }}相关的结果{{
        total
      }}个，用时{{ userTime }}s
    </div>
    <div style="max-height: 500px;margin: 10px 0px;overflow: auto;">
      <div
        v-for="(v, index) in tableData"
        :key="index"
        class="content-item"
        style="align-items: center;padding: 20px;background-color: #f4f4f5;margin-top: 10px;"
      >
        <div>{{ index + 1 }}.</div>
        <div style="flex:1 1;margin-left: 20px;" class="">
          <div class="content-item">
            <div class="text-green">From：</div>
            <div>{{ v._index }}</div>
          </div>
          <div class="content-item">
            <div class="text-green">Content：</div>
            <div>
              <div>{</div>
              <div v-for="(v, index, i) in v._source" :key="i">
                <div>
                  &nbsp;&nbsp;&nbsp;&nbsp;"{{ index }}" : {{ v ? `"${v}"` : v }}
                </div>
              </div>
              <div>}</div>
            </div>
          </div>
        </div>
      </div>
    </div>
    
    <div v-if="tableData != undefined && tableData.length > 0">
      <el-pagination
        background
        layout="prev, pager, next"
        :total="total"
        :current-page="currentPage"
        :page-size="pageLimit"
        @current-change="changePage"
      >
      </el-pagination>
    </div>
  </div>
</template>

<script>
import Moment from "moment";
import { getIndexCfg, SearchLog, GetMappingAlias } from "@/api/search";

export default {
  name: "LogSearch",
  components: {},
  props: {
    value: {
      type: String,
      default: ""
    }
  },
  data() {
    return {
      tableData: [],
      total: 0,
      indexOptions: [],
      indexValue: [],
      inputIndex: "",
      SearchLogFilter: [],
      currentPage: 1,
      pageLimit: 10,
      userTime: 0,
      indexIsInput: false,
      indexMap: {},
      isFirstSearch: true
    };
  },
  computed: {},
  watch: {
    indexValue: {
      handler(newV, oldV) {
        this.setSearchLogFilter(newV);
      },
      deep: true
    },
    indexIsInput(newV, oldV) {
      if (newV) {
        this.SearchLogFilter = [{ search_col: "ip", search_text: "" }];
        return;
      }
      this.setSearchLogFilter(this.indexValue);
    }
  },
  mounted() {
    this.getIndexCfg();
  },
  methods: {
    setSearchLogFilter(value) {
      let col = [];
      for (const key of value) {
        let item = this.indexOptions.find(
          v => v.remark == key || v.indexName == key
        );
        if (item) {
          col.push(...item.input_cols);
        }
      }
      col = [...new Set(col)];
      this.SearchLogFilter = col.map(v => {
        return { search_col: v, search_text: "" };
      });
    },
    onBtnSearch() {
      this.SearchLog();
    },
    async getIndexCfg() {
      const res = await getIndexCfg({
        page: 1,
        limit: 10,
        all: true,
        es_connect: this.$store.state.baseData.EsConnectID
      });
      if (res.code != 0) {
        this.$message({
          type: "error",
          message: res.msg
        });
        return;
      }
      let indexOptions = res.data.list;
      if (indexOptions == null) {
        indexOptions = [];
      }
      let filterIndexOptions = [];
      //获取别名
      let map = {};
      for (const v of indexOptions) {
        // console.log("进入到了这里啊", this.value, v.indexName);
        if (v.indexName.indexOf(this.value) >= 0) {
          filterIndexOptions.push(v);
          let ret = await GetMappingAlias({
            es_connect: this.$store.state.baseData.EsConnectID,
            index_name: v.indexName
          });
          if (ret && ret.code == 0) {
            let tempMap = ret.data.res || {};
            for (const key in tempMap) {
              if (!map[key]) map[key] = tempMap[key];
            }
          }
        }
      }
      this.indexOptions = filterIndexOptions;
      this.indexMap = map;
      // console.log("这个东西啊", this.indexMap, filterIndexOptions);
      if (this.indexOptions.length > 0) {
        this.indexValue = [this.indexOptions[0].indexName];
      }
      this.$message({
        type: "success",
        message: res.msg
      });
    },
    async SearchLog() {
      let time = Date.now();
      let indexIsInput = this.indexIsInput;
      let indexNames = [];
      indexNames = indexIsInput ? [this.inputIndex] : this.indexValue;
      this.tableData = []

      const res = await SearchLog({
        index_names: indexNames,
        search_log_filter: this.SearchLogFilter,
        mode: indexIsInput ? 1 : 0,
        page: this.currentPage,
        limit: this.pageLimit,
        es_connect: this.$store.state.baseData.EsConnectID
      });
      if (res.code != 0) {
        this.$message({
          type: "error",
          message: res.msg
        });
        return;
      }
      this.isFirstSearch = false
      let data = res.data || {};
      this.total = data.count ? data.count.value : 0;
      this.userTime = (Date.now() - time) / 1000;
      if(Object.keys(data).length > 0){
        this.tableData = data.list.hits.hits;
      }
      this.$message({
        type: "success",
        message: res.msg
      });
    },
    changePage(v) {
      this.currentPage = v;
      this.SearchLog();
    }
  }
};
</script>

<style>
.content-item {
  display: flex;
  justify-content: flex-start;
}
.text-green {
  color: blue;
}
</style>
