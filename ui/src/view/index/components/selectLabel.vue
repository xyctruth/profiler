<template>
  <el-select
      v-model="value"
      :multiple="true"
      placeholder="Select Label"
      clearable="true"
      filterable>
    <el-option
        v-for="(item,index) in options"
        :key="index"
        :label="item"
        :value="item"
    >
    </el-option>
  </el-select>
</template>

<script setup>
import {ref, onMounted, defineEmits, watch} from 'vue'
import axios from "@/utils/request";
import router from "@/router/index.js";

const options = ref([])
const value = ref([])
const emit = defineEmits(['update:selectLabels'])

let query = router.currentRoute.value.query
if (query.labels){
  if (Array.isArray(query.labels)){
    value.value = query.labels
  }else{
    value.value = [query.labels]
  }
}

onMounted(() => {
  axios({
    url: "/api/labels",
  })
      .then((res) => {
        options.value = res
      })
      .catch((err) => {
        console.log(err);
      })
})

watch(value, (value) => {
  let query = router.currentRoute.value.query
  let path = router.currentRoute.value.path
  let newQuery = JSON.parse(JSON.stringify(query));
  newQuery.labels = value;
  router.push({path, query: newQuery});
  emit("update:selectLabels",  value)
}, {
  immediate: true
})

</script>
