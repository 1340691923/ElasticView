<template>
  <div>

    <div class="xwl_main" style="margin-top: 20px">
      <div class="relation-editor globalFilters_xwl">

        <div class="relation-relation">
          <em class="relation-relation-line" />
          <div
            v-if="user_filter.filts.length > 1"
            class="relation-relation-value"
            @click="changeRelationLine"
          >{{ user_filter.relation }}
          </div>
        </div>
        <div class="relation-main">

          <div v-for="(v,index) in user_filter.filts" :key="index" class="relation-row">

            <div class="ta-multa-filter-condition">

              <div v-if="v.filterType == 'SIMPLE'" class="action-row row___xwl">

                <action-row :key="index" v-model="user_filter.filts[index]" v-if="changeActionRow" :options="options" :table-typ="tableTyp" :data-type-map="dataTypeMap" class="action-left" />
                <div class="action-right">
                  <a-button-group>

                    <a-button
                      type="link"
                      class="actions_xwl_btn"
                      icon="filter"
                      @click="addRelationSimple(index)"
                    />
                    <a-button
                      type="link"
                      class="actions_xwl_btn"
                      icon="close-circle"
                      @click="deleteRelationSimple(index)"
                    />

                  </a-button-group>
                </div>
              </div>
              <div v-else class="relation-editor">
                <div class="relation-relation">
                  <div class="relation-relation-sub">◆</div>
                  <em class="relation-relation-line-sub" />
                  <div class="relation-relation-value" @click="changeRelationLine1(index)">
                    {{ v.relation }}
                  </div>
                </div>
                <div class="relation-main">
                  <div class="relation-row ">
                    <div v-for="(v2,index2) in v.filts" :key="index2" class="action-row row___xwl">
                      <div class="action-left">
                        <action-row :key="(index+'_'+index2).toString()" v-model="v.filts[index2]" :options="options" :table-typ="tableTyp" :data-type-map="dataTypeMap" class="action-left" />
                      </div>
                      <div class="action-right">
                        <a-button-group>

                          <a-button
                            type="link"
                            class="actions_xwl_btn"
                            icon="filter"
                            @click="addRelationCOMPOUND(index)"
                          />

                          <a-button
                            type="link"
                            class="actions_xwl_btn"
                            icon="close-circle"
                            @click="deleteRelationCOMPOUND(index,index2)"
                          />

                        </a-button-group>
                      </div>
                    </div>
                  </div>
                </div>
              </div>
            </div>
          </div>
        </div>
      </div>
    </div>
    <div style="padding: 0 12px;">
      <span :style="{color:fontColor}" class="footadd___2D4YB" @click="addRelation1">
        <a-icon type="filter" />
        增加条件
      </span>
    </div>
  </div>
</template>

<script>

export default {
  name: 'Index',
  components: {
    ActionRow: () => import('@/components/AnalyseTools/FilterWhere/ActionRow')
  },
  props: {
    value: {
      type: Object,
      default: {
        relationLine: '且',
        relationArr: [],
        ftv: []
      }
    },
    fontColor: {
      type: String,
      default: '#3d90ff'
    },
    tableTyp: {
      type: Number,
      default: 0
    },
    dataTypeMap: {
      type: Array,
      default: []
    },
    options: {
      type: Array,
      default: []
    }
  },
  data() {
    return {
      changeActionRow:true,
      user_filter: this.value
    }
  },
  methods: {
    refreshChangeActionRow() {
      this.changeActionRow = false
      this.$nextTick(() => {
        this.changeActionRow = true
      })
    },
    addRelationSimple(index) {
      const obj = this.user_filter.filts[index]
      const COMPOUNDObj = {
        filterType: 'COMPOUND',
        filts: [],
        relation: '且'
      }

      COMPOUNDObj.filts.push(obj)

      COMPOUNDObj.filts.push({
        columnName: this.options[0].options[0].value,
        comparator: '=',
        filterType: 'SIMPLE',
        ftv: ''

      })

      Vue.set(this.user_filter.filts, index, COMPOUNDObj)
      this.$emit('input', this.user_filter)
    },
    addRelationCOMPOUND(index) {
      this.user_filter.filts[index].filts.push(
        {
          columnName: this.options[0].options[0].value,
          comparator: '=',
          filterType: 'SIMPLE',
          ftv: ''

        })
      this.$emit('input', this.user_filter)
    },
    deleteRelationCOMPOUND(index1, index2) {
      if (this.user_filter.filts[index1].filts.length > 2) {
        this.user_filter.filts[index1].filts.splice(index2, 1)
      } else {
        const arr = JSON.parse(JSON.stringify(this.user_filter.filts[index1].filts))
        arr.splice(index2, 1)
        const obj = arr[0]
        Vue.set(this.user_filter.filts, index1, obj)
      }

      this.$emit('input', this.user_filter)
    },

    addRelation1() {
      if (this.options[0].options.length == 0) {
        this.$message({
          offset: 60,

          message: '该应用没有用户上报数据或者属性被隐藏',
          type: 'error'
        })
        return
      }
      this.user_filter.filts.push(
        {
          columnName: this.options[0].options[0].value,
          comparator: '=',
          filterType: 'SIMPLE',
          ftv: ''
        })
      this.$emit('input', this.user_filter)
    },
    deleteRelationSimple(i) {

      this.user_filter.filts.splice(i, 1)

      this.$emit('input', this.user_filter)
      this.refreshChangeActionRow()
    },

    changeRelationLine() {
      this.user_filter.relation == '或' ? this.user_filter.relation = '且' : this.user_filter.relation = '或'
      this.$emit('input', this.user_filter)
    },
    changeRelationLine1(index) {
      this.user_filter.filts[index].relation == '或' ? this.user_filter.filts[index].relation = '且' : this.user_filter.filts[index].relation = '或'
      this.$emit('input', this.user_filter)
    }
  }
}
</script>

<style lang="scss" scoped>
  .sider_xwl {
    display: inline-block;
    width: 270px;
    height: calc(100vh - 50px);
    padding: 0;
    overflow: hidden;
    vertical-align: top;
  }

  .top_xwl {
    display: flex;
    align-items: center;
    justify-content: flex-start;
    width: 100%;
    height: 57px;
    border-bottom: 1px solid #f0f0f0;
    border-right: 1px solid #f0f0f0;
    padding-right: 20px;
  }

  .body_xwl {
    border-bottom: 1px solid #f0f0f0;
    border-right: 1px solid #f0f0f0;
    height: calc(100vh - 120px);
    overflow-y: auto;
    background-color: white;
  }

  .content_xwl {
    display: inline-block;
    width: 100%;
    height: calc(100vh - 50px);
    padding: 0;
    overflow: hidden;
    vertical-align: top;
    background: #f0f2f5;
    border-bottom: 2px solid #e6e6e6;
  }

  .header_xwl {
    position: relative;
    z-index: 1000;
    width: 100%;
    height: 56px;
    padding: 0;
    line-height: 56px;

  }

  .main_xwl {
    display: flex;
    align-items: center;
  }

  .title_xwl {
    display: inline-block;
    margin-right: 16px;
    color: black;
    font-weight: 500;
    font-size: 18px;
  }

  .dashbordName_xwl {
    display: flex;
    align-items: center;
    justify-content: space-between;
    height: 40px;

    border-radius: 2px;
    cursor: pointer;
  }

  .root_xwl {
    display: flex;
    align-items: center;
    justify-content: space-between;
    width: 100%;
    height: 56px;
    padding: 0 16px;
    background-color: #fff;
    border-left: 1px solid #f0f0f0;
  }

  .actions_xwl {
    display: flex;
    align-items: center;
  }

  .actions_xwl_btn:hover {
    color: orangered;
  }

  .actions_xwl_btn {
    color: #67729d;
  }

  .echartBox {
    width: 48%;
    height: 400px;
    padding: 5px;
    margin-bottom: 5px;
  }

  .echartBox:hover {
    cursor: pointer;
  }

  .echartBox_title {
    cursor: pointer;
  }

  .echartBox_title:hover {
    color: #3d90ff;
  }

  .footadd___2D4YB {
    display: inline-block;
    margin-right: 16px;
    padding: 4px;
    color: #3d90ff;
    font-size: 13px;
    line-height: 20px;
    border-radius: 2px;
    cursor: pointer;
    transition: all .3s;
    transition-property: all;
    transition-duration: 0.3s;
    transition-timing-function: ease;
    transition-delay: 0s;
  }

  .right_res {
    position: relative;
    width: 100%;
    height: calc(100% - 32px);

    overflow-x: hidden;
    overflow-y: auto;
    white-space: normal;
    transition: width .3s;
  }

  .row___xwl {
    min-height: 40px;
    padding: 0 4px 0 8px;
    transition-property: all;
    transition-duration: 0.3s;
    transition-timing-function: ease;
    transition-delay: 0s;
    padding: 10px;
  }

  .row___xwl:hover {
    box-shadow: 0 0 3px 0 #1890ff;
    transition-property: all;
    transition-duration: 0.3s;
    transition-timing-function: ease;
    transition-delay: 0s;
  }

  .action-row {
    position: relative;
    display: flex;
    align-items: center;
    justify-content: space-between;
    width: 100%;
    min-height: 24px;
  }

  .relation-relation-line {
    position: absolute;
    top: 0;
    left: 12px;
    width: 1px;
    height: 100%;
    background-color: #d9dfe6;
    transition-property: all;
    transition-duration: 0.3s;
    transition-timing-function: ease;
    transition-delay: 0s;
  }

  .relation-relation-value {
    position: absolute;
    top: 50%;
    left: 0;
    width: 24px;
    height: 24px;
    margin-top: -12px;
    color: #3d90ff;
    font-size: 12px;
    line-height: 22px;
    text-align: center;
    background-color: #fff;
    border: 1px solid #d9dfe6;
    border-radius: 2px;
    cursor: pointer;
    transition-property: all;
    transition-duration: 0.3s;
    transition-timing-function: ease;
    transition-delay: 0s;
  }

  .xwl_main {

  }

  .globalFilters_xwl {
    padding: 0 8px 0 12px;
  }

  .relation-editor {
    position: relative;
    display: flex;
    width: 100%;
  }

  .relation-editor .relation-relation {
    position: relative;
    display: flex;
    flex-direction: column;
    align-items: center;
    justify-content: flex-start;
    width: 24px;
    margin-right: 4px;
  }

  .relation-relation-line {
    position: absolute;
    top: 0;
    left: 12px;
    width: 1px;
    height: 100%;
    background-color: #d9dfe6;
    transition-property: all;
    transition-duration: 0.3s;
    transition-timing-function: ease;
    transition-delay: 0s;
  }

  .relation-editor .relation-relation .relation-relation-value {
    position: absolute;
    top: 50%;
    left: 0;
    width: 24px;
    height: 24px;
    margin-top: -12px;
    color: #3d90ff;
    font-size: 12px;
    line-height: 22px;
    text-align: center;
    background-color: #fff;
    border: 1px solid #d9dfe6;
    border-radius: 2px;
    cursor: pointer;
    transition-property: all;
    transition-duration: 0.3s;
    transition-timing-function: ease;
    transition-delay: 0s;
  }

  .relation-editor .relation-main {
    flex-grow: 1;
    flex-shrink: 1;
    flex-basis: 0%;
  }

  .relation-editor .relation-relation .relation-relation-sub {
    position: relative;
    top: -7px;
    left: .5px;
    z-index: 100;
    color: #a3acc5;
    font-size: 11px;
  }

  .relation-editor .relation-relation .relation-relation-line-sub {
    position: absolute;
    top: 13px;
    left: 12px;
    width: 1px;
    height: calc(100% - 10px);
    background-color: #d9dfe6;
    transition-property: all;
    transition-duration: 0.3s;
    transition-timing-function: ease;
    transition-delay: 0s;
  }

  .ta-multa-filter-condition {
    padding: 0;
  }

  .relation-row .action-row {
    align-items: center;
    justify-content: flex-start;
  }

  .row___xwl {
    min-height: 32px;
    padding: 2px 4px;
  }

  .action-row .action-left {
    display: flex;
    align-items: center;
  }

  .ta-filter-condition {
    display: inline-block;
    min-height: 32px;
    padding-bottom: 2px;
    white-space: normal;
  }

  .action-row .action-right {
    display: flex;
    align-items: center;
  }

  .relation-relation:hover {
    cursor: pointer;
  }

  .relation-relation:hover .relation-relation-line {
    background-color: #3d90ff;
  }

  .relation-relation:hover .relation-relation-value {
    background-color: #3d90ff;
    color: white;
  }

  .relation-relation:hover .relation-relation-line-sub {
    background-color: #3d90ff;
  }

  .relation-relation:hover .relation-relation-sub {
    color: #3d90ff;
  }

  ::-webkit-scrollbar {
    width: 8px; /*对垂直流动条有效*/
    height: 10px; /*对水平流动条有效*/
  }

  /*定义滚动条的轨道颜色、内阴影及圆角*/
  ::-webkit-scrollbar-track {
    border-radius: 3px;
  }

  /*定义滑块颜色、内阴影及圆角*/
  ::-webkit-scrollbar-thumb {
    border-radius: 7px;
    background-color: #e6e6e6;
  }

  /*定义两端按钮的样式*/
  ::-webkit-scrollbar-button {
  }

  /*定义右下角汇合处的样式*/
  ::-webkit-scrollbar-corner {
    background: khaki;
  }

  .eventRow_xwl {
    position: relative;
  }

  .rename_xwl {
    width: auto;
    min-width: 50px;
    max-width: 260px;
    height: 24px;
    margin: 8px 0 0;
    padding: 0;
    line-height: 24px;
    background: inherit;
  }

  .rename_xwl span {
    max-width: 260px;
    font-weight: 500;
    font-size: 13px;
  }

  .rename_xwl input {
    font-weight: 500;
    font-size: 13px;
  }

  .eventRow_xwl .eventItem_xwl {
    max-width: 400px;
    padding: 2px 0 6px;
    overflow: hidden;
    white-space: normal;
  }

  .filters_xwl {
    padding-left: 32px;
  }

  .formula_xwl {
    padding-right: 8px;
    padding-left: 32px;
  }

</style>
