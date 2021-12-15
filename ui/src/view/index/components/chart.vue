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

  const canvas = ref(null)
  const canvasBox = ref(null)

  let chart = null
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
  const getCanvasWidthHeight = () => {
    if (canvasBox.value) {
      let {offsetWidth, offsetHeight} = canvasBox.value || {}
      return {
        width: offsetWidth || 800,
        height: offsetHeight || 400
      }
    }
  }

  const formatUnit = (value,unit)=>{
    if (unit === "bytes"){
      if (null == value || value == '') {
        return "0 Bytes";
      }
      var unitArr = new Array("Bytes", "KB", "MB", "GB", "TB", "PB", "EB", "ZB", "YB");
      var index = 0;
      var srcsize = parseFloat(value);
      index = Math.floor(Math.log(srcsize) / Math.log(1024));
      var size = srcsize / Math.pow(1024, index);
      size = size.toFixed(2);//保留的小数位数
      return size + unitArr[index];
    }else if (unit === "nanoseconds"){
      if (null == value || value == '') {
        return "0";
      }
      var unitArr = new Array("ns", "us", "ms", "s");
      var index = 0;
      var srcsize = parseFloat(value);
      index = Math.floor(Math.log(srcsize) / Math.log(1000));
      if (index > 3){
        index = 3
      }
      var size = srcsize / Math.pow(1000, index);
      size = size.toFixed(2);//保留的小数位数
      return size + unitArr[index];
    }
    else if (unit === "count"){
      return value
    }
    return value+unit
  }

  const setChart = ({data = [], title} = {}) => {
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
          top:80,
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
            formatter:(params)=>{
              return  formatUnit(params,unit)
            },
          }
        },
        tooltip: {
          show: true,
          trigger: 'item',
          axisPointer: {
            animation: false
          },
          formatter:(params)=>{
            console.log(params)
            var val = formatUnit(params.data.value[1],unit)

            return `<div style="margin: 0px 0 0;line-height:1;">
                        <div style="font-size:14px;color:#666;font-weight:400;line-height:1;">`+params.seriesName+`</div>
                        <div style="margin: 10px 0 0;line-height:1;"><div style="margin: 0px 0 0;line-height:1;">
                        <span style="display:inline-block;margin-right:4px;border-radius:10px;width:10px;height:10px;background-color:`+params.color+`;"></span>
                        <span style="float:right;margin-left:10px;font-size:14px;color:#666;font-weight:900">`+val+`</span>
                        <div style="clear:both"></div></div><div style="clear:both"></div>
                        </div>
                        <div style="clear:both"></div>
                     </div>`
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
          window.open(`${baseConfig.reqUrl}/web/profile/${params.data.sourceData.ProfileID}`)
        });
      }
      chart.setOption(Object.assign(chartOptions, {
        series: echartData,
      }), true);

    })
  }

  onMounted(() => {
    setChart()
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

