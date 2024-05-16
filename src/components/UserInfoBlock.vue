<script setup lang="ts" name="UserInfoBlock"></script>

<template>
  <div class="user-info-block">
    <!-- Personal Information -->
    <div class="user-info-area">
      <van-image round fit="cover" src="https://img.yzcdn.cn/vant/cat.jpeg"   width="5rem"
  height="5rem"></van-image>
  <span>UserName</span>
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
      <van-cell title="设置" is-link />
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
