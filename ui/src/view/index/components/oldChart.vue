<template>
  <div>
    <canvas
      :id="index+'canvas'"
      style="cursor: pointer"
      ref="canvas"></canvas>
  </div>
</template>

<script setup>
  import {defineProps, ref, onMounted, watch, nextTick} from 'vue'
  import Chart from '@/utils/chart'
  import axios from '@/utils/request'
  import baseConfig from '@/config/config'
  const canvas = ref(null)
  import moment from 'moment'
  import autocolors from 'chartjs-plugin-autocolors/dist/chartjs-plugin-autocolors.esm';

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
  watch(() => [props.type, props.projects, props.timeRange], () => {
    let sampleType = props.type
    let targetList = props.projects
    let startTime = props.timeRange[0] ? `&start_time=${props.timeRange[0]}` : ""
    let endTime = props.timeRange[1] ? `&end_time=${props.timeRange[1]}` : ""
    let targetFilter = targetList.reduce((result, item) => {
      return result + `&targets=${item}`
    }, "")
    axios({
      url: `/api/samples/${sampleType}?${startTime}${endTime}${targetFilter}`
    })
      .then(targets => {
        const data = {
          labels: [],
          datasets: []
        };

        for (const target of targets) {
          const item = {
            label: target.Name,
            style: {
              cursor: "pointer",
            },
            data: [],
          }
          for (const sample of target.Samples) {
            const label = moment(sample.Key).format("YYYY-MM-DD HH:mm:ss")
            data.labels.push(label)
            item.data.push({x: label, y: sample.CountValue, profileID: sample.ProfileID})
          }
          data.datasets.push(item)
        }

        data.labels = [...new Set(data.labels)]
        data.labels = data.labels.sort((a, b) => {
          return new Date(a) - new Date(b);
        })

        console.log(data)
        if (chart) {
          chart.data = data
          chart.update();
        } else {
          const config = {
            type: 'line',
            data: data,
            options: {
              interaction: {
                intersect: false,
              },
              plugins: {
                title: {
                  display: true,
                  text: sampleType,
                  font: {
                    size: 25
                  }
                },
              },
              onClick: (e, el, chart) => {
                if (!el[0]) return;
                const profileID = el[0].element.$context.raw.profileID
                window.open(`https://www.speedscope.app/#profileURL=${baseConfig.reqUrl}/api/profile/` + profileID);
              },
              elements: {
                bar: {
                  pointStyle: "cross"
                }

              }
            },
          };
          chart = new Chart(
            canvas.value.getContext("2d"),
            config
          );
        }
      })
  }, {
    deep: true,
    immediate: true,
  })

</script>

