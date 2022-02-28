<template>
  <el-select
    v-model="value"
    :collapse-tags="true"
    :multiple="true"
    placeholder="Select Type"
    :clearable="true"
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
  const emit = defineEmits(['update:selectTypes'])

  let query = router.currentRoute.value.query
  if (query.types){
    if (Array.isArray(query.types)){
      value.value = query.types
    }else{
      value.value = [query.types]
    }
  }

  onMounted(() => {
    axios({
      url: "/api/group_sample_types",
    })
      .then((res) => {
        const types = ["goroutine","profile","heap","fgprof","allocs","block","threadcreate","mutex","trace"]
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
    newQuery.types = value;
    router.push({path, query: newQuery});
    emit("update:selectTypes",  value)
  }, {
    immediate: true
  })

</script>
