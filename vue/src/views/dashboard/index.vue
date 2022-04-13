<template>
  <div>
    <el-card class="box-card">
      <div class="dashboard-editor-container">
        <github-corner class="github-corner" />

        <panel-group  />
        <!---todo..-->
        <!--<el-row style="background:#fff;padding:16px 16px 0;margin-bottom:32px;">
          <line-chart :chart-data="lineChartData" />
        </el-row>-->

        <!--<el-row :gutter="32">
          <el-col :xs="24" :sm="24" :lg="8">
            <div class="chart-wrapper">
              <raddar-chart />
            </div>
          </el-col>
          <el-col :xs="24" :sm="24" :lg="8">
            <div class="chart-wrapper">
              <pie-chart />
            </div>
          </el-col>
          <el-col :xs="24" :sm="24" :lg="8">
            <div class="chart-wrapper">
              <bar-chart />
            </div>
          </el-col>
        </el-row>-->

        <!--<el-row :gutter="8">
          <el-col
            :xs="{span: 24}"
            :sm="{span: 24}"
            :md="{span: 24}"
            :lg="{span: 12}"
            :xl="{span: 12}"
            style="padding-right:8px;margin-bottom:30px;"
          >
            <transaction-table />
          </el-col>-->
        <!--  <el-col
            :xs="{span: 24}"
            :sm="{span: 12}"
            :md="{span: 12}"
            :lg="{span: 6}"
            :xl="{span: 6}"
            style="margin-bottom:30px;"
          >
            <todo-list />
          </el-col>-->
        <!--  <el-col
            :xs="{span: 24}"
            :sm="{span: 12}"
            :md="{span: 12}"
            :lg="{span: 6}"
            :xl="{span: 6}"
            style="margin-bottom:30px;"
          >
            <box-card />
          </el-col>
        </el-row>-->
      </div>
    </el-card>
  </div>
</template>

<script>
import GithubCorner from '@/components/GithubCorner'
import PanelGroup from './components/PanelGroup'


export default {
  name: 'DashboardAdmin',
  components: {
    GithubCorner,
    PanelGroup,
  },
  data() {
    return {

    }
  },
  methods: {

  }
}
</script>

<style lang="scss" scoped>
  .dashboard-editor-container {
    padding: 32px;
    background-color: rgb(240, 242, 245);
    position: relative;

    .github-corner {
      position: absolute;
      top: 0px;
      border: 0;
      right: 0;
    }

    .chart-wrapper {
      background: #fff;
      padding: 16px 16px 0;
      margin-bottom: 32px;
    }
  }

  @media (max-width: 1024px) {
    .chart-wrapper {
      padding: 8px;
    }
  }
</style>
