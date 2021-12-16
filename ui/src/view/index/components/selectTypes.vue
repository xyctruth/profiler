<template>
  <el-select
    :collapse-tags="true"
    :multiple="true"
    :value="value"
    @input="emit('input',$event)"
    placeholder="选择类型"
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
  import {ref, onMounted, defineEmits, defineProps} from 'vue'
  import axios from "@/utils/request";

  const options = ref([])
  const props = defineProps({
    value: {
      type: Array,
      default: () => {
        return []
      }
    }
  })
  const emit = defineEmits(['input'])
  onMounted(() => {
    axios({
      url: "/api/group_sample_types",
    })
      .then((res) => {
        const types = ["heap","fgprof","profile","goroutine","allocs","block","threadcreate","mutex"]
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

</script>
