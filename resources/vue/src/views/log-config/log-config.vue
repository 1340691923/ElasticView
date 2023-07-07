<template>
  <div class="app-container">
    <div class="filter-container">
      <el-button
        size="mini"
        type="success"
        class="filter-item"
        @click="open = true"
        >{{ $t("新增配置") }}
      </el-button>
      <el-button
        size="mini"
        type="success"
        class="filter-item"
        @click="aliasDialog = true"
        >{{ $t("修改别名") }}
      </el-button>
      <el-button size="mini" type="success" class="filter-item" @click="search"
        >{{ $t("刷新") }}
      </el-button>
    </div>

    <el-table v-loading="tableLoading" :data="tableData">
      <el-table-column align="center" prop="indexName" label="索引" />
      <el-table-column align="center" prop="remark" label="备注" />
    </el-table>
    <div class="pagination-container">
      <el-pagination
        v-if="pageshow"
        background
        :current-page="input.page"
        :page-size="input.limit"
        layout="total, sizes, prev, pager, next, jumper"
        :total="count"
        @current-change="search"
        @size-change="handleSizeChange"
      />
    </div>

    <el-dialog
      :close-on-click-modal="false"
      :visible.sync="open"
      :title="$t('新增配置')"
      @close="closeDialog"
    >
      <el-form label-width="100px" label-position="left">
        <el-form-item label="索引名">
          <el-select
            v-model="form.indexName"
            v-loading="indexSelectLoading"
            filterable
            style="width:300px"
          >
            <el-option
              v-for="(v, k, index) in indexList"
              :key="index"
              :label="v"
              :value="v"
            />
          </el-select>
        </el-form-item>
        <el-form-item label="字段1">
          <el-select
            v-model="form.input_cols"
            filterable
            multiple
            style="width:300px"
          >
            <el-option
              v-for="(v, k, index) in searchInOptions"
              :key="index"
              :label="v.label"
              :value="v.value"
            />
          </el-select>
        </el-form-item>
        <el-form-item label="字段2">
          <el-select
            v-model="form.output_cols"
            class="filter-item"
            filterable
            multiple
            style="width:300px"
          >
            <el-option
              v-for="(v, k, index) in searchOutOptions"
              :key="index"
              :label="v.label"
              :value="v.value"
            />
          </el-select>
        </el-form-item>
        <el-form-item :label="$t('备注')">
          <el-input
            v-model="form.remark"
            style="width:300px"
            placeholder="备注"
          />
        </el-form-item>
      </el-form>
      <div style="text-align:right;">
        <el-button
          size="mini"
          type="danger"
          icon="el-icon-close"
          @click="closeDialog"
          >{{ $t("取消") }}
        </el-button>
        <el-button size="mini" type="primary" icon="el-icon-check" @click="add"
          >{{ $t("确认") }}
        </el-button>
      </div>
    </el-dialog>
    <el-dialog
      :close-on-click-modal="false"
      :visible.sync="aliasDialog"
      :title="$t('修改别名')"
      @close="closeAliasSetDialog"
    >
      <el-form label-width="100px" label-position="left">
        <el-form-item label="索引名">
          <el-select
            v-model="form.indexName"
            v-loading="indexSelectLoading"
            filterable
            style="width:300px"
          >
            <el-option
              v-for="(v, k, index) in indexList"
              :key="index"
              :label="v"
              :value="v"
            />
          </el-select>
          <div v-for="(v, index) in searchInOptions" :key="index">
            <span>字段名：{{ v.value }}</span>
            <span>别名：</span>
            <el-input
              v-model="v.label"
              style="width:300px"
              placeholder="别名"
            />
          </div>
        </el-form-item>
      </el-form>
      <div style="text-align:right;">
        <el-button
          size="mini"
          type="danger"
          icon="el-icon-close"
          @click="closeAliasSetDialog"
          >{{ $t("取消") }}
        </el-button>
        <el-button
          size="mini"
          type="primary"
          icon="el-icon-check"
          @click="onBtnConfirmOtherName"
          >{{ $t("确认") }}
        </el-button>
      </div>
    </el-dialog>
  </div>
</template>

<script>
import { IndexNamesAction } from "@/api/es-index";
import {
  getIndexCfg,
  setIndexCfg,
  SetMappingAlias,
  GetMappingAlias
} from "@/api/search";
import { ListAction } from "@/api/es-map";

const defaultForm = {
  indexName: "",
  remark: "",
  input_cols: [],
  output_cols: []
};

export default {
  name: "Link",
  data() {
    return {
      indexList: [],
      count: 0,
      pageshow: true,
      tableData: [],
      form: Object.assign({}, defaultForm),
      input: {
        page: 1,
        limit: 10,
        all: false
      },
      searchInOptions: [],
      searchOutOptions: [],
      tableLoading: false,
      indexSelectLoading: false,
      open: false,
      indexMap: {},
      aliasDialog: false
    };
  },
  watch: {
    "form.indexName": {
      deep: true,
      handler(newV, oldV) {
        if (newV) this.getListAction();
      }
    }
  },
  created() {
    this.getIndexList();
    this.search();
  },
  methods: {
    closeDialog() {
      this.open = false;
      this.form = Object.assign({}, defaultForm);
      this.searchInOptions = [];
      this.searchOutOptions = [];
    },
    closeAliasSetDialog() {
      this.aliasDialog = false;
      this.form = Object.assign({}, defaultForm);
      this.searchInOptions = [];
      this.searchOutOptions = [];
    },
    getIndexList() {
      this.indexSelectLoading = true;
      const input = {};
      input["es_connect"] = this.$store.state.baseData.EsConnectID;
      IndexNamesAction(input)
        .then(res => {
          this.indexSelectLoading = false;
          if (res.code == 0) {
            this.indexList = res.data;
          } else {
            this.$message({
              type: "error",
              message: res.msg
            });
          }
        })
        .catch(err => {
          this.indexSelectLoading = false;
          console.log(err);
        });
    },
    refreshPage() {
      this.pageshow = false;
      this.count = this.tableData.length;
      this.$nextTick(() => {
        this.pageshow = true;
      });
    },
    handleSizeChange(v) {
      this.input.limit = v;
      this.refreshPage();
    },
    async add() {
      const form = this.form;
      form["es_connect"] = this.$store.state.baseData.EsConnectID;

      const res = await setIndexCfg(form);
      if (res.code != 0) {
        this.$message({
          type: "error",
          message: res.msg
        });
        return;
      }
      this.$message({
        type: "success",
        message: res.msg
      });
      this.open = false;
      this.search();
    },

    async search() {
      this.tableLoading = true;
      const input = this.input;
      input["es_connect"] = this.$store.state.baseData.EsConnectID;
      const res = await getIndexCfg(input);
      if (res.code != 0) {
        this.$message({
          type: "error",
          message: res.msg
        });
        return;
      }
      this.tableData = res.data.list;
      this.count = res.data.count;
      if (this.tableData == null) {
        this.tableData = [];
      }

      this.$message({
        type: "success",
        message: res.msg
      });
      this.tableLoading = false;
    },
    async getListAction() {
      let res = await GetMappingAlias({
        es_connect: this.$store.state.baseData.EsConnectID,
        index_name: this.form.indexName
      });
      if (res.code != 0) {
        this.$message({
          type: "error",
          message: msg
        });
        return;
      }
      this.indexMap = res.data.res || {};
      const { data, code, msg } = await ListAction({
        es_connect: this.$store.state.baseData.EsConnectID,
        index_name: this.form.indexName
      });

      if (code != 0) {
        this.$message({
          type: "error",
          message: msg
        });
        return;
      }
      if (!(data && data.list && data.list[this.form.indexName])) return;
      let mappings = data.list[this.form.indexName].mappings;
      if (!mappings) return;
      let options = mappings.properties;
      let keys = Object.keys(options);
      this.searchOutOptions = keys.map(v => {
        return {
          value: v,
          label: this.aliasDialog
            ? this.indexMap[v] || ""
            : this.indexMap[v] || v
        };
      });
      this.searchInOptions = this.searchOutOptions.filter(
        v =>
          options[v.value].type == "text" || options[v.value].type == "keyword"
      );
      this.form.output_cols = keys;
    },
    async onBtnConfirmOtherName() {
      let indexName = this.form.indexName;
      let map = {};
      this.searchInOptions.forEach(v => {
        map[v.value] = v.label;
      });
      let res = await SetMappingAlias({
        es_connect: this.$store.state.baseData.EsConnectID,
        index_name: indexName,
        mapping_cfg: map
      });
      if (res.code != 0) {
        this.$message({
          type: "error",
          message: msg
        });
        return;
      }
      this.$message({
        type: "success",
        message: res.msg
      });
      this.closeAliasSetDialog();
    }
  }
};
</script>

<style scoped></style>
