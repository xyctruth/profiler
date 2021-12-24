<template>
  <div class='index container'>
    <el-row :gutter="22">
      <el-col :xs="24" :sm="24" :md="24" :lg="18" :xl="18">
        <el-row :gutter="22">
          <el-col>
            <el-card style="margin-bottom: 20px;margin-top:10px">
              <!--          <div flex="cross:center main:justify">-->

              <el-row :gutter="30">
                <!--              这里按照视图分层-->

                <el-col :xs="24" :sm="24" :md="8" :lg="9" :xl="4">
                  <selectTypes style="margin-right: 20px; width: 90%" v-model:selectTypes="types"></selectTypes>
                </el-col>
                <el-col :xs="24" :sm="24" :md="8" :lg="9" :xl="4">
                  <selectProject style="margin-right: 20px;width: 90%"
                                 v-model:selectProjects="projects"></selectProject>
                </el-col>
                <el-col :xs="24" :sm="24" :md="8" :lg="9" :xl={span:6,offset:10} >
                  <selectTimeRange class="time-range-container" v-model:timeRange="timeRange"
                                   style="width: 90%"></selectTimeRange>
                </el-col>

              </el-row>
              <!--          </div>-->

            </el-card>

            <div>
              <el-col :xs="24" :sm="24" :md="0" :lg="0" :xl="0" style="padding-right: 0px;padding-left: 0px; width: 100%">
                <el-affix :offset="10" class="guid">
                  <el-card>
                    <el-button
                        @click="jump(item)"
                        v-for="(item,index) in types" :key="index" type="text"
                        style="padding-right: 20px;"
                    >
                      {{ item }}
                    </el-button>
                  </el-card>
                </el-affix>
              </el-col>
            </div>
            <div>
              <el-card
                  class="char-card"
                  v-for="(item,index) in types"
                  :key="index"
                  style="margin-bottom: 30px;overflow: auto;overflow: scroll  ">
                <div :id="item"></div>
                <chart
                    style="min-width: 1200px;"
                    :index="index"
                    :type="item"
                    :projects="projects"
                    :timeRange="timeRange"></chart>
              </el-card>
            </div>
          </el-col>
        </el-row>

      </el-col>
      <el-col :md="6" :lg="6" :xl="6" class="hidden-md-and-down">
        <el-row>
          <el-col :xs="24" :sm="24" :md="24" :lg="24" :xl="24">
            <el-affix :offset="10" class="guid">
              <el-card>
                <el-button
                    @click="jump(item)"
                    v-for="(item,index) in types" :key="index" type="text"
                    style="padding-right: 20px;"
                >
                  {{ item }}
                </el-button>
              </el-card>
            </el-affix>
          </el-col>
        </el-row>
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

console.log('local')
const types = ref([])
const projects = ref([])
const timeRange = ref([])
const jump = (item) => {
  document.getElementById(item).scrollIntoView()
}
</script>

<style lang="scss" scoped>
.guid {
  height: auto !important;
  width: auto !important;

  .el-button {
    margin-left: 0;
    display: block;
  }
}
</style>
<style lang="scss">

@media only screen and (max-width: 1200px) {
  .hidden-md-and-down {
    display: none !important
  }
}

.el-col {
  min-height: 1px
}

.char-card {
  .el-card__body {
    padding: 20px 0;
  }
}

.el-col {
  margin-bottom: 10px;
}

.time-range-container {
  .el-input {
    width: 100%;
  }
}


</style>
