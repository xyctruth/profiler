<template>
  <div class='index container'>
    <el-row :gutter="20">
      <el-col :span="20">
        <el-card style="margin-bottom: 20px;">
          <div flex="cross:center main:justify">
            <div>
              <selectTypes style="margin-right: 20px;" v-model="types"></selectTypes>
              <selectProject style="margin-right: 20px;" v-model="projects"></selectProject>
            </div>
            <div>
              <selectTimeRange v-model:time="timeRange"></selectTimeRange>
            </div>
          </div>
        </el-card>
        <div>
          <el-card
            class="char-card"
            v-for="(item,index) in types"
            :key="index"
            style="margin-bottom: 30px;">
            <div :id="item"></div>
            <chart
              :index="index"
              :type="item"
              :projects="projects"
              :timeRange="timeRange"></chart>
          </el-card>
        </div>
      </el-col>
      <el-col :span="4">
        <el-affix :offset="20" class="guid">
          <el-card>
            <el-button
              @click="jump(item)"
              v-for="(item,index) in types" :key="index" type="text">
              {{item}}
            </el-button>
          </el-card>
        </el-affix>
      </el-col>
    </el-row>

  </div>
</template>

<script setup>
  import {ref} from 'vue'
  import selectTypes from './components/selectTypes.vue'
  import selectProject from './components/selectProject.vue'
  import selectTimeRange from './components/selectTimeRange.vue'
  import chart from './components/chart.vue'

  const types = ref([])
  const projects = ref([])
  const timeRange = ref([])
  const jump = (item) => {
    document.getElementById(item).scrollIntoView()
  }
</script>

<style lang="scss" scoped>
  .guid {
    .el-button {
      margin-left: 0;
      display: block;
    }
  }
</style>
<style lang="scss">
  .char-card {
    .el-card__body {
      padding: 20px 0;
    }
  }
</style>
