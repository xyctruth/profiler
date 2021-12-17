<template>
  <div>
    <el-date-picker v-model="value" type="date" :disabled-date="disabledDate">
    </el-date-picker>
  </div>
</template>

<script setup>
  import { reactive, toRefs, defineEmits, watch} from 'vue'
  import moment from 'moment'

  const state = reactive({
    disabledDate(time) {
      return time.getTime() > Date.now()
    },
    value: new Date(),
  })
  const {disabledDate, value} = toRefs(state)

  const emit = defineEmits(['update:timeRange'])
  const getCurrentValue = () => {

    var startDate = moment(value.value).startOf('day');
    var endDate = moment(value.value).add(1,'days').startOf('day')

    return [startDate.toISOString(),endDate.toISOString()]
  }

  watch(value, (value) => {
    emit("update:timeRange", getCurrentValue())
  }, {
    immediate: true
  })

</script>

