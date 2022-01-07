<template>
  <el-select
      v-model="value"
      :multiple="true"
      placeholder="Select Label"
      clearable="true"
      filterable>
    <el-option-group
        v-for="group in options"
        :key="group.label"
        :label="group.label"
    >
      <el-option
          v-for="(item,index) in group.options"
          :key="index"
          :label="item"
          :value="item"
      >
      </el-option>
    </el-option-group>
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
    url: "/api/group_labels",
  })
      .then((res) => {
        const types = ["custom","generate"]
        var data = []
        for (const key of types) {
          if (res[key]) {
            data.push({
              label: key,
              options: res[key]
            })
          }
        }
        options.value = data
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
