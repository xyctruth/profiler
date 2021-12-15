<template>
  <el-select v-model="value" placeholder="Select">
    <el-option
      v-for="item in options"
      :key="item.label"
      :label="item.label"
      :value="item.label"
    >
    </el-option>
  </el-select>
</template>

<script setup>
  import {ref, onMounted, defineEmits, defineProps, watch} from 'vue'
  import moment from 'moment'

  const value = ref("1小时内")
  const options = ref([
    {
      label: "1小时内",
      value: () => {
        return [moment().subtract(1, 'hours').toISOString(), moment().toISOString()]
      },
    },
    {
      label: "12小时内",
      value: () => {
        return [moment().subtract(12, 'hours').toISOString(), moment().toISOString()]
      },
    },
    {
      label: "一天内",
      value: () => {
        return [moment().subtract(1, 'days').toISOString(), moment().toISOString()]
      },
    },
    {
      label: "一周内",
      value: () => {
        return [moment().subtract(7, 'days').toISOString(), moment().toISOString()]
      },
    },
  ])
  const emit = defineEmits(['update:time'])
  const getCurrentValue = () => {
    let item = options.value.find((item) => {
      return item.label === value.value
    })
    return item ? item.value() : options[0].value()
  }
  watch(value, (value) => {
    emit("update:time", getCurrentValue())
  }, {
    immediate: true
  })


</script>

<style scope lang="scss">
  .selectTimeRange {

  }
</style>
