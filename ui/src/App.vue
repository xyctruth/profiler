<script setup>
  import MyHeader from './components/MyHeader.vue'
</script>

<template>
  <div style=" width: 100%;overflow: hidden">
<!--    <MyHeader></MyHeader>-->
    <router-view/>
    <el-backtop />
  </div>
</template>

<style lang="scss">
  @import "@/assets/css/global.scss";

  #app {

  }
</style>
