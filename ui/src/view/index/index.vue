<template>
  <div class='index container'>
    <el-row :gutter="22">
      <el-col :xs="24" :sm="24" :md="24" :lg="20" :xl="20">
        <el-row :gutter="22">
          <el-col>
            <el-card style="margin-bottom: 20px;margin-top:10px">
              <!--          <div flex="cross:center main:justify">-->

              <el-row :gutter="30">
                <!--              这里按照视图分层-->

                <el-col :xs="24" :sm="24" :md="5" :lg="5" :xl="5">
                  <selectTypes style="width: 100%" v-model:selectTypes="types"></selectTypes>
                </el-col>
                <el-col :xs="24" :sm="24" :md="8" :lg="8" :xl="8">
                  <selectLabels style="width: 100%" v-model:selectLabels="labels"></selectLabels>
                </el-col>
                <el-col :xs="24" :sm="24" :md={span:4,offset:7} :lg={span:4,offset:7} :xl={span:4,offset:7} >
                  <selectTimeRange class="time-range-container" v-model:timeRange="timeRange"
                                   style="width: 100%"></selectTimeRange>
                </el-col>

              </el-row>
              <!--          </div>-->

            </el-card>

            <div>
              <el-col :xs="24" :sm="24" :md="24" :lg="0" :xl="0" style="padding-right: 0px;padding-left: 0px; width: 100%">
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
                    :targets="targets"
                    :timeRange="timeRange"
                    :labels="labels"></chart>
              </el-card>
            </div>
          </el-col>
        </el-row>

      </el-col>
      <el-col :md="6" :lg="4" :xl="4" class="hidden-md-and-down">
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
import selectLabels from './components/selectLabels.vue'
import selectTimeRange from './components/selectTimeRange.vue'
import chart from './components/chart.vue'

const types = ref([])
const targets = ref([])
const labels = ref([])
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

.container{
  margin-left: 20px;
  margin-right: 20px;
}

@media only screen and (max-width: 1200px) {
  .hidden-md-and-down {
    display: none !important
  }

  .container{
    margin-left: 5px;
    margin-right: 5px;
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
