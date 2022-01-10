<template>
  <div ref="canvasBox" style="width:100%;height:400px;">
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
    targets: {
      type: Array,
      default: () => {
        return {}
      }
    },
    labels: {
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
    nextTick(() => {
      var unit = ""
      const baseSetting = {
        type: "scatter",
        showSymbol: false,
        sampling: 'lttb',
        showAllSymbol: false,
        symbolSize: 10,
        emphasis: {
          width: 3,
          focus: 'series',
        },
      }

      const chartOptions = {
        animation: false, // 关闭加载动画
        legend: {
          type: 'scroll',
          orient: 'vertical',
          right: 10,
          top: 20,
          bottom: 20,
          selectedMode: 'multiple',
        },
        grid: {
          left: 100,
          right: 300,
          top: 80,
        },
        dataZoom: {
          start: 0,
          end: 100,
          type: "slider",
          show: true,
          realtime: false, // 是否实时刷新
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
          confine: true,
          trigger: 'item',
          axisPointer: {
            animation: false
          },
          formatter: (params) => {
            return formatTooltip(params, unit)
          }
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
            if (!unit){
              unit = p.SampleTypeUnit
            }
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
          height,
          useDirtyRect: true,
        });
        chart.on('click', function (params) {
          if (title === "trace") {
            window.open(`${baseConfig.reqUrl}/api/trace/ui/${params.data.sourceData.ProfileID}`)
          }else{
            window.open(`${baseConfig.reqUrl}/api/pprof/ui/${params.data.sourceData.ProfileID}?si=${title}`)

          }
        });
      }

      if (title ==="trace"){
        chartOptions.legend.selectedMode = 'single'
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

  watch(() => [props.type, props.targets,props.labels, props.timeRange], () => {
    if (props.type) {
      let sampleType = props.type
      let targetList = props.targets
      let labelList = props.labels
      let startTime = props.timeRange[0] ? `&start_time=${props.timeRange[0]}` : ""
      let endTime = props.timeRange[1] ? `&end_time=${props.timeRange[1]}` : ""

      var labels = labelList.filter(function (item) {
        if (item) {
          return true;
        } else {
          return false
        }
      }).map(function (item) {
          return {Key:item[0],Value:item[1]}
      })


      let loading = ElLoading.service({
        text: "加载中..."
      })

      axios({
        url: `/api/profile_meta/${sampleType}?${startTime}${endTime}`,
        params: {
          "labels": labels,
          "targets":targetList,
        }
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

