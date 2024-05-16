<script setup lang="ts" name="BottomTabBar">
import { ref } from 'vue'
const active = ref<number>(0)
const onChange = (index: number) => {
  // get
  console.log('index:', index)
}
</script>

<template>
  <div>
    <van-tabbar v-model="active" @change="onChange">
      <van-tabbar-item replace to="/home" icon="home-o">首页</van-tabbar-item>
      <!-- <van-tabbar-item replace to="/photos" icon="photo-o">
        照片
      </van-tabbar-item> -->
      <!-- <van-tabbar-item icon="friends-o">标签</van-tabbar-item> -->
      <van-tabbar-item replace to="/my" icon="setting-o">设置</van-tabbar-item>
    </van-tabbar>
  </div>
</template>

<style lang="scss" scoped></style>
