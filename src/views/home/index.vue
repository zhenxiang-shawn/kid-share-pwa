<script setup lang="ts">
import { showToast } from 'vant'
import { ref } from 'vue'
import TimeLineItem from '../../components/TimeLineItem.vue'

let loading = ref<boolean>(false)
let finished = ref<boolean>(false)
let container = ref()
const fake_images1 = [
  'https://images.unsplash.com/photo-1568572933382-74d440642117',
  'https://images.unsplash.com/photo-1534377125276-8d48c3f9c3d3',
  'https://images.unsplash.com/photo-1551963831-b0d5d7de7c49',
  'https://images.unsplash.com/photo-1554629467-62245042f9a5',
  'https://images.unsplash.com/photo-1520880867055-1e30d1cb001c',
  'https://images.unsplash.com/photo-1513089184088-a4b54380a0b5',
  'https://images.unsplash.com/photo-1519687367256-4f5f4c5e3c6e',
  'https://images.unsplash.com/photo-1504198247232-7af2c1e6e8f2',
  'https://images.unsplash.com/photo-1518734712468-44c6f2a2c27b',
]

const fake_images2 = [
  'https://source.unsplash.com/960x640/?nature',
  'https://source.unsplash.com/960x640/?architecture',
  'https://source.unsplash.com/960x640/?technology',
  'https://source.unsplash.com/960x640/?people',
  'https://source.unsplash.com/960x640/?landscape',
]

const onRefresh = () => {
  setTimeout(() => {
    loading.value = false
    showToast('刷新成功')
  }, 1000)
}
</script>

<template>
  <div>
    <!-- 标题栏 -->
    <van-nav-bar title="Shy 娃时间线" />
    <!-- 配置下拉刷新 -->
    <van-pull-refresh v-model="loading" @refresh="onRefresh">
      <div class="main-content">
        contant area
        <van-list
          finished-text="没有更多了"
          :finished="finished"
          v-model="loading"
        >
          <TimeLineItem poster="Apple" date="2024-03-01" title="mama" />

          <TimeLineItem
            poster="二婶"
            date="2024-03-02"
            title="二婶"
            :images="fake_images1"
          />

          <TimeLineItem
            poster="小三姑"
            date="2024-03-03"
            title="三姑"
            :images="fake_images2"
          />
        </van-list>
      </div>
    </van-pull-refresh>
  </div>
</template>

<style lang="scss" scoped>
.main-content {
  overflow-y: auto; /* 显示垂直滚动条 */
  height: calc(100vh - var(--van-nav-bar-height) - var(--van-tabbar-height));
}
</style>
