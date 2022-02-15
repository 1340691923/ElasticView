<template>
  <div :class="className" :style="{height:height,width:width}" />
</template>

<script>
// import echarts from 'echarts'
require('echarts/theme/macarons') // echarts theme
import resize from './mixins/resize'

export default {
  mixins: [resize],
  props: {
    className: {
      type: String,
      default: 'chart'
    },
    width: {
      type: String,
      default: '100%'
    },
    height: {
      type: String,
      default: '400px'
    },
    autoResize: {
      type: Boolean,
      default: true
    },
    legend: { // 设置区分（哪条线属于什么）
      type: Array,
      required: true
    },
    xAxisData: {
      type: Array,
      required: true
    },
    series: { // 描画数据
      type: Array,
      required: true
    },
    yAxisName: {
      type: String,
      default: ''
    },
    xAxisName: {
      type: String,
      default: ''
    },
    legendColor: {
      type: Array,
      required: true
    }
  },
  data() {
    return {
      chart: null
    }
  },
  watch: {
    series: {
      deep: true,
      handler() {
        this.setOptions()
      }
    }
  },
  mounted() {
    this.$nextTick(() => {
      this.initChart()
    })
  },
  beforeDestroy() {
    if (!this.chart) {
      return
    }
    this.chart.dispose()
    this.chart = null
  },
  methods: {
    initChart() {
      this.chart = echarts.init(this.$el, 'macarons')
      this.setOptions(this.chartData)
    },
    setOptions() {
      console.log('setOptions')
      const series = []
      for (var k in this.legend) {
        const tmpSeries = {
          name: this.legend[k],
          itemStyle: {
            normal: {
              color: this.legendColor[k],
              lineStyle: {
                color: this.legendColor[k],
                width: 2
              }
            }
          },
          smooth: true,
          type: 'line',
          data: this.series[k],
          animationDuration: 2800,
          animationEasing: 'quadraticOut'
        }
        if (k == 0) {
          tmpSeries.animationEasing = 'cubicInOut'
        }
        series.push(tmpSeries)
      }

      var option = {
        tooltip: { // 设置tip提示
          trigger: 'axis'
        },
        grid: {
          left: 30,
          right: 90,
          bottom: 60,
          top: 60,
          containLabel: true
        },
        legend: { // 设置区分（哪条线属于什么）
          data: this.legend
        },
        color: this.legendColor, // 设置区分（每条线是什么颜色，和 legend 一一对应）
        xAxis: { // 设置x轴
          type: 'category',
          boundaryGap: false, // 坐标轴两边不留白
          data: this.xAxisData,
          name: this.xAxisName // X轴 name
          /* nameTextStyle: {        //坐标轴名称的文字样式
            color: '#FA6F53',
            fontSize: 16,
            padding: [0, 0, 0, 20]
          },*/
          /* axisLine: {             //坐标轴轴线相关设置。
            lineStyle: {
              color: '#FA6F53',
            }
          }*/
        },
        yAxis: {
          name: this.yAxisName,
          /* nameTextStyle: {
            color: '#FA6F53',
            fontSize: 16,
            padding: [0, 0, 10, 0]
          },*/
          /* axisLine: {
            lineStyle: {
              color: '#FA6F53',
            }
          },*/
          type: 'value'
        },
        series: series
      }
      this.chart.setOption(
        option
      )
    }
  }
}
</script>

<!--
[
          {
          name: '点击次数2', itemStyle: {
            normal: {
              color: '#FF005A',
              lineStyle: {
                color: '#FF005A',
                width: 2
              }
            }
          },
          smooth: true,
          type: 'line',
          data: expectedData,
          animationDuration: 2800,
          animationEasing: 'cubicInOut'
        },
          {
          name: '点击次数',
          smooth: true,
          type: 'line',
          itemStyle: {
            normal: {
              color: '#3888fa',
              lineStyle: {
                color: '#3888fa',
                width: 2
              },
              areaStyle: {
                color: '#f3f8ff'
              }
            }
          },
          data: actualData,
          animationDuration: 2800,
          animationEasing: 'quadraticOut'
        }
        ]-->
