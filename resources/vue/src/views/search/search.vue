<template>
  <div class="app-container">
    <div>
      <el-checkbox v-model="indexIsInput">模糊搜索</el-checkbox>
      <el-select
        multiple
        v-model="indexValue"
        filterable
        placeholder="请选择索引"
        v-if="!indexIsInput"
      >
        <el-option
          v-for="(v, index) in indexOptions"
          :key="index"
          :label="v.remark"
          :value="v.indexName"
        ></el-option>
      </el-select>
      <el-input
        v-model="inputIndex"
        style="width: 200px;"
        placeholder="请输入索引"
        v-else
      ></el-input>
      <el-tag>字段名</el-tag>
      <el-input v-model="otherInput" style="width: 100px;"></el-input>
      <el-input
        v-model="searchInput"
        style="width: 300px;"
        placeholder="请输入内容"
      ></el-input>
      <el-button type="primary" @click="onBtnSearch">搜索</el-button>
    </div>
    <div>
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
    <div>
      <el-pagination
        background
        layout="prev, pager, next"
        :total="total"
        :current-page="currentPage"
        :page-size=pageLimit
        @current-change="changePage"
      >
      </el-pagination>
    </div>
  </div>
</template>

<script>
import Moment from "moment";
import { getIndexCfg, SearchLog } from "@/api/search";
export default {
  name: "search",
  components: {
    JsonEditor: () => import("@/components/JsonEditor/index")
  },
  data() {
    return {
      tableData: [],
      total: 0,
      indexOptions: [],
      indexValue: [],
      inputIndex: "",
      searchInput: "",
      otherInput: "ip",
      currentPage: 1,
      pageLimit: 10,
      userTime: 0,
      indexIsInput: false
    };
  },
  computed: {},
  mounted() {
    this.getIndexCfg();
  },
  methods: {
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
      this.indexOptions = res.data.list;
      if (this.tableData == null) {
        this.indexOptions = [];
      }
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
      const res = await SearchLog({
        index_names: indexNames,
        search_col: this.otherInput,
        search_text: this.searchInput,
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
      let data = res.data || {};
      this.total = data.count ? data.count.value : 0;
      this.userTime = (Date.now() - time) / 1000;
      this.tableData = data.list.hits.hits;

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
