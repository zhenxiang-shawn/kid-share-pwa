<script lang="ts">
export default {
  name: 'Layout',
}
</script>
<script setup lang="ts">
import BottomTabBar from './bottom-tabbar/index.vue'
import MainDisplay from './main-display/index.vue'
</script>

<template>
  <div class="layout_container">
    <div class="main_display">
      <MainDisplay />
    </div>

    <div class="bottom_tab_bar">
      <BottomTabBar />
    </div>
  </div>
</template>

<style lang="scss" scoped>
.layout_container {
  width: 100%;
  height: 100vh;
  display: flex;
  // flex-direction: column;
  .bottom_tab_bar {
    height: 50px;
  }
  .main_display {
    width: 100%;
    height: calc(100vh - 50px);
    flex: 1;
    // overflow-y: auto;
  }
}
</style>
