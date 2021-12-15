<template>
  <el-select
    :multiple="true"
    :value="value"
    :collapse-tags="true"
    @input="emit('input',$event)"
    placeholder="选择项目">
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
      url: "/api/targets",
    })
      .then((res) => {
        options.value = res
      })
      .catch((err) => {
        console.log(err);
      })
  })

</script>
