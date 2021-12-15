<template>
  <div ref="canvasBox" style="width:100%;height:600px;">
    <div
      class="canvas"
      :id="index+'canvas'"
      ref="canvas"></div>
  </div>
</template>

<script setup>
  import {defineProps, ref, onMounted, watch, nextTick} from 'vue'
  import axios from '@/utils/request'
  import baseConfig from '@/config/config'
  import moment from 'moment'
  import * as echarts from 'echarts';
  import {ElLoading} from 'element-plus'
  import {formatTooltip,formatUnit} from '@/utils/chart'
  const canvas = ref(null)
  const canvasBox = ref(null)
  let chart = null

  const getCanvasWidthHeight = () => {
    if (canvasBox.value) {
      let {offsetWidth, offsetHeight} = canvasBox.value || {}
      return {
        width: offsetWidth || 800,
        height: offsetHeight || 400
      }
    }
  }

  const props = defineProps({
    index: {
      type: [String, Number],
    },
    type: {
      type: String,
      default: () => {
        return {}
      }
    },
    projects: {
      type: Array,
      default: () => {
        return {}
      }
    },
    timeRange: {
      type: Array,
      default: () => {
        return []
      }
    }
  })

  const setChart = ({data = [], title} = {}) => {
    console.log("set chart")
    nextTick(() => {
      var unit = ""
      const baseSetting = {
        type: "line",
        showSymbol: true,
        lineStyle: {
          width: 3,
        },
        symbolSize: 10
      }
      const chartOptions = {
        legend: {
          type: 'scroll',
          orient: 'vertical',
          right: 10,
          top: 20,
          bottom: 20,
        },
        grid: {
          left: 100,
          right: 300,
          top: 80,
        },
        dataZoom: {
          start: 0,
          type: "slider",
          show: true,
        },
        title: {
          text: title,
          x: 'center'
        },
        xAxis: {
          type: "time",
          axisLabel: {
            formatter: (v) => {
              return moment(v).format('MM-DD HH:mm')
            }
          }
        },
        yAxis: {
          axisLabel: {
            formatter: (params) => {
              return formatUnit(params, unit)
            },
          }
        },
        tooltip: {
          show: true,
          trigger: 'item',
          axisPointer: {
            animation: false
          },
          formatter: (params) => {
            return formatTooltip(params, unit)
          },
        },
        series: []
      }

      const echartData = []
      if (data) {
        let count = 0
        for (const meta of data) {
          count++
          const item = {
            ...baseSetting,
            name: meta.TargetName,
            data: [],
          }
          for (const p of meta.ProfileMetas) {
            unit = p.SampleTypeUnit
            let label = moment(p.Timestamp).format('YYYY-MM-DD HH:mm:ss')
            if (label) {
              item.data.push({
                sourceData: p,
                value: [label, p.Value]
              })
            }
          }
          echartData.push(item)
        }
      }

      if (!chart) {
        let {width, height} = getCanvasWidthHeight()
        chart = echarts.init(canvas.value, null, {
          devicePixelRatio: window.devicePixelRatio || 2,
          width,
          height
        });
        chart.on('click', function (params) {
          window.open(`${baseConfig.reqUrl}/pprof/register/${params.data.sourceData.ProfileID}?si=${title}`)
        });
      }
      chart.setOption(Object.assign(chartOptions, {
        series: echartData,
      }), true);

    })
  }

  onMounted(() => {
    window.addEventListener('resize', () => {
      if (chart) {
        let {width, height} = getCanvasWidthHeight()
        chart.resize({
          width,
          height
        });
      }
    })
  })

  watch(() => [props.type, props.projects, props.timeRange], () => {
    if (props.type) {
      let sampleType = props.type
      let targetList = props.projects
      let startTime = props.timeRange[0] ? `&start_time=${props.timeRange[0]}` : ""
      let endTime = props.timeRange[1] ? `&end_time=${props.timeRange[1]}` : ""
      let targetFilter = targetList.reduce((result, item) => {
        return result + `&targets=${item}`
      }, "")

      let loading = ElLoading.service({
        text: "加载中..."
      })
      axios({
        url: `/api/profile_meta/${sampleType}?${startTime}${endTime}${targetFilter}`
      })
        .then(targets => {
          setChart({data: targets, title: sampleType})
        })
        .finally(() => {
          loading.close()
        })
    }
  }, {
    deep: true,
    immediate: true,
  })

</script>

<style lang="scss" scoped>
  .canvas {
    width: 100%;
    height: 100%;
  }
</style>

