<template>
  <div class="MyHeader container">
    <div flex="cross:center main:justify" class="main-header">
      <!--左边内容-->
      <div>
        <img class="logo" src="http://cdnabc.gytycn.com/image/2020-8-14-144653877-f05d74f6ddf911eab93300163e1299a2.jpg" alt="">
      </div>
      <!--右边内容-->
      <div>
        <el-button type="text">github</el-button>
        <el-button type="text">Help</el-button>
      </div>

    </div>
  </div>
</template>

<script setup>

</script>

<style scope lang="scss">
  .MyHeader {
    .main-header {
      height:60px;
      .logo{
        width:60px;
      }
    }
  }
</style>
