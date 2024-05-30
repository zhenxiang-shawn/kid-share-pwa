<script setup lang="ts" name="UserInfoBlock">
import { showNotify } from 'vant'
import useUserStore from '../store/modules/user'
import { computed } from 'vue'

let userStore = useUserStore()
let props = defineProps(['image_url'])
let image_url = computed(() => {
  return props.image_url == ''
    ? props.image_url
    : 'https://fastly.jsdelivr.net/npm/@vant/assets/cat.jpeg'
})
const logout = () => {
  // 需要向服务器发请求[退出登录接口]******

  // 仓库当中关于用于相关的数据清空[token|username|avatar]
  userStore.logout()
  window.location.reload() // TODO(dir to login page)
  // notification
  showNotify({ type: 'warning', message: '已登出' })
}
</script>

<template>
  <div class="user-info-block">
    <!-- Personal Information -->
    <div class="user-info-area" @click="$router.push('/my/info')">
      <van-image
        round
        fit="cover"
        :src="image_url"
        width="5rem"
        height="5rem"
      ></van-image>
      <span>{{ userStore.displayName }}</span>
    </div>
    <!-- avator -->
    <!-- user name -->
    <!-- Post Information -->
    <van-grid :column-num="3">
      <van-grid-item icon="star-o" text="收藏" />
      <van-grid-item icon="clock-o" text="历史" />
      <van-grid-item icon="orders-o" text="作品" />
    </van-grid>
    <!-- Function Area -->
    <van-cell-group>
      <van-cell title="消息通知" is-link />
      <van-cell title="用户反馈" is-link />
      <van-cell title="登出" @click="logout" />
    </van-cell-group>
  </div>
</template>

<style lang="scss" scoped>
.user-info-block {
  height: 100%;
  display: flex;
  flex-direction: column;
  .user-info-area {
    display: flex;
    flex-direction: column;
    align-items: center;
    justify-content: center;
    height: 10rem;
    background-color: #47d4b1;
  }
}
</style>
