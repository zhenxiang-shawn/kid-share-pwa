<script setup lang="ts">
import { computed, ref } from 'vue'
import { reqPostDiary } from '../../APP_interface/diary_posts'
import { diariesRequestData } from '../../APP_interface/diary_posts/type'
import useUserStore from '../../store/modules/user'
import { showNotify } from 'vant'

const onClickLeft = () => history.back()
const fileList = ref([
  { url: 'https://fastly.jsdelivr.net/npm/@vant/assets/leaf.jpeg' },
  // Uploader 根据文件后缀来判断是否为图片文件
  // 如果图片 URL 中不包含类型信息，可以添加 isImage 标记来声明
  // { url: 'https://cloud-image', isImage: true },
])

let userStore = useUserStore()
let post = ref<string>('')
let post_disabled = computed(() => {
  return post.value.length == 0 && fileList.value.length == 0
})
let imagePaths = computed(() => {
  return fileList.value.map((item) => item.url)
})

const sendPost = async () => {
  console.log('images:', fileList.value)
  const result = await reqPostDiary({
    username: userStore.username as string,
    content: post.value,
    image_paths: imagePaths.value,
  })

  if (result.code == 200) {
    post.value = ''
    fileList.value = []
    console.log('发布成功')
    showNotify({ type: 'success', message: '发布成功' })
  } else {
    console.log('发布失败')
    showNotify({ type: 'danger', message: '发布失败' })
  }
}
</script>

<template>
  <div>
    <van-nav-bar
      title="发布日记"
      left-text="返回"
      right-text="发布"
      left-arrow
      @click-left="onClickLeft"
      @click-right="sendPost"
      :right-disabled="post_disabled"
      v-model="post"
    />
    <van-cell-group class="post-group" inset>
      <van-field
        v-model="post"
        rows="3"
        autosize
        type="textarea"
        placeholder="请输入日记"
        class="post-input"
      />
      <van-uploader v-model="fileList" multiple></van-uploader>
    </van-cell-group>
  </div>
</template>

<style lang="scss" scoped>
.post-group {
  // border-top: 50px;
  margin: 10px 5px;
  van-uploader {
    display: flex;
    // 向左浮动显示
    float: left;
  }
}
</style>
