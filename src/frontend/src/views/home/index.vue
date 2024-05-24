<script setup lang="ts">
import { showToast } from 'vant'
import { onMounted, ref } from 'vue'
import TimeLineItem from '../../components/TimeLineItem.vue'
import useContentStore from '../../store/modules/content'
import { useRouter } from 'vue-router'

let loading = ref<boolean>(false)
let finished = ref<boolean>(false)
let useContent = useContentStore()
let $router = useRouter()

const fake_images1 = [
  'https://images.unsplash.com/photo-1568572933382-74d440642117',
  'https://images.unsplash.com/photo-1534377125276-8d48c3f9c3d3',
  'https://images.unsplash.com/photo-1551963831-b0d5d7de7c49',
  'https://images.unsplash.com/photo-1554629467-62245042f9a5',
  'https://images.unsplash.com/photo-1520880867055-1e30d1cb001c',
  'https://images.unsplash.com/photo-1513089184088-a4b54380a0b5',
  'https://images.unsplash.com/photo-1519687367256-4f5f4c5e3c6e',
  'https://images.unsplash.com/photo-1504198247232-7af2c1e6e8f2',
  'https://1000logos.net/wp-content/uploads/2016/10/Apple-Logo.png',
]

const fake_images2 = [
  'https://source.unsplash.com/960x640/?nature',
  'https://source.unsplash.com/960x640/?architecture',
  'https://source.unsplash.com/960x640/?technology',
  'https://source.unsplash.com/960x640/?people',
  'https://source.unsplash.com/960x640/?landscape',
]

const getContent = () => {
  // load data from localStorage

  // if there is no localstorage, request from server
  try {
    useContent.getContent()
  } catch (e) {
    console.log(e)
  }
}

const addPost = () => {
  $router.push({
    path: '/post',
  })
}

const onRefresh = () => {
  getContent()
  setTimeout(() => {
    loading.value = false
    showToast('刷新成功')
  }, 1000)
}

onMounted(() => {
  getContent()
})
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
          <TimeLineItem
            v-for="(item, index) in useContent.contents"
            :key="item.id"
            :content="item.content"
            :poster="item.username"
            :timestamp="item.timestamp"
            :images="item.images"
          />
        </van-list>
      </div>
    </van-pull-refresh>
    <van-floating-bubble icon="add" @click="addPost" />
  </div>
</template>

<style lang="scss" scoped>
.main-content {
  overflow-y: auto; /* 显示垂直滚动条 */
  height: calc(100vh - var(--van-nav-bar-height) - var(--van-tabbar-height));
}
</style>
