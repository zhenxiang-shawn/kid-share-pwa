<script setup lang="ts">
import { showToast } from 'vant'
import { onMounted, ref, watch } from 'vue'
import TimeLineItem from '../../components/TimeLineItem.vue'
import useContentStore, {
  CLEAR_LOCAL_CONTENT,
} from '../../store/modules/content'
import { useRouter } from 'vue-router'

let pageNumber = 1
const pageLimit = 20
let loading = ref<boolean>(false)
let finished = ref<boolean>(false)
let useContent = useContentStore()
let $router = useRouter()

const getContent = async (initial: boolean = false) => {
  if (initial) {
    // initial value
    try {
      pageNumber = 1
      CLEAR_LOCAL_CONTENT()
      await useContent.updateContent(pageNumber++, pageLimit)
      finished.value = false // reset finished state when initial load
    } catch (e) {
      console.log(e)
    }
  } else {
    if (pageLimit * pageNumber >= useContent.contentsNum) {
      finished.value = true
      return 'ok'
    }
    try {
      await useContent.updateContent(pageNumber++, pageLimit)
      if (useContent.contents.length >= useContent.contentsNum) {
        finished.value = true
      }
    } catch (e) {
      console.log(e)
    }
  }
}

const addPost = () => {
  $router.push({
    path: '/post',
  })
}

const onRefresh = async () => {
  loading.value = true
  await getContent(true)
  loading.value = false
  showToast('刷新成功')
}

const onLoad = async () => {
  if (finished.value) {
    return
  }
  loading.value = true
  await getContent()
  loading.value = false
}

onMounted(() => {
  if (localStorage.getItem('content')) {
    getContent()
  } else {
    getContent(true)
  }
})

// Watch for changes in contents to update finished state
watch(
  () => useContent.contents,
  (newContents) => {
    if (newContents.length >= useContent.contentsNum) {
      finished.value = true
    }
  },
)
</script>

<template>
  <div>
    <!-- 标题栏 -->
    <van-nav-bar title="Shy 娃时间线" />
    <!-- 配置下拉刷新 -->
    <van-pull-refresh
      v-model="loading"
      @refresh="onRefresh"
      class="main-content"
    >
      contant area: {{ finished }}
      <van-list
        finished-text="没有更多了"
        :finished="finished"
        v-model="loading"
        @load="onLoad"
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
    </van-pull-refresh>
    <van-floating-bubble icon="add" @click="addPost" />
  </div>
</template>

<style lang="scss" scoped>
.main-content {
  overflow-y: auto; /* 显示垂直滚动条 */
  height: calc(100vh - var(--van-nav-bar-height) - var(--van-tabbar-height));
  van-list {
    height: 100%;
  }
}
</style>
