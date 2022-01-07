<template>
  <el-select
    v-model="value"
    :multiple="true"
    :collapse-tags="true"
    placeholder="Select Target"
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
  const emit = defineEmits(['update:selectTargets'])

  let query = router.currentRoute.value.query
  if (query.targets){
    if (Array.isArray(query.targets)){
      value.value = query.targets
    }else{
      value.value = [query.targets]
    }
  }

  onMounted(() => {
    axios({
      url: "/api/targets",
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
  newQuery.targets = value;
  router.push({path, query: newQuery});
  emit("update:selectTargets",  value)
}, {
  immediate: true
})

</script>
