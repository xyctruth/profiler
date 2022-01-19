<template>
  <el-cascader
      v-model="value"
      placeholder="Select Labels"
      :options="options"
      :props="props"
      :clearable="true"
      separator="="
      filterable
      :collapse-tags="true">
  </el-cascader>
</template>

<script setup>
  import {ref, onMounted, defineEmits, watch} from 'vue'
  import axios from "@/utils/request";
  import router from "@/router/index.js";

  const options = ref([])
  const value = ref([])
  const emit = defineEmits(['update:selectLabels'])

  const props = {
    multiple: true,
  }

  let query = router.currentRoute.value.query
  if (query.labels){
    if (!Array.isArray(query.labels)){
      value.value = [query.labels.split(',')]
    }else{
      for (const p of query.labels) {
        value.value.push(p.split(','))
      }
    }
  }

  onMounted(() => {
    axios({
      url: "/api/group_labels",
    }).then((res) => {
      var data = []
      for (var key in res) {
        var l = {
          label: key,
          value:key,
          children: [],
        }

        for (const p of res[key]) {
          l.children.push({value: p.Value, label: p.Value})
        }
        data.push(l)
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
