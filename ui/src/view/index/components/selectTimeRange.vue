<template>
  <div>
    <el-date-picker v-model="value" type="date" :disabled-date="disabledDate">
    </el-date-picker>
  </div>
</template>

<script setup>
  import { reactive, toRefs, defineEmits, watch} from 'vue'
  import router from "@/router/index.js";
  import moment from 'moment'

  var date = moment().startOf('day')
  var query = router.currentRoute.value.query
  if (query.date){
    date = moment(query.date,moment.ISO_8601)
  }

  const state = reactive({
    disabledDate(time) {
      return time.getTime() > moment().toDate()
    },
    value:date.toDate(),
  })
  const {disabledDate, value} = toRefs(state)

  const emit = defineEmits(['update:timeRange'])
  const getCurrentValue = (value) => {
    var startDate = moment(value).startOf('day');
    var endDate = moment(value).add(1,'days').startOf('day')
    return [startDate.toISOString(),endDate.toISOString()]
  }

  watch(value, (value) => {
    let query = router.currentRoute.value.query
    let path = router.currentRoute.value.path
    let newQuery = JSON.parse(JSON.stringify(query));
    newQuery.date = moment(value).toISOString();
    router.push({path, query: newQuery});
    emit("update:timeRange",  getCurrentValue(value))
  }, {
    immediate: true
  })

</script>

