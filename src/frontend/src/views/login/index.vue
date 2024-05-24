<script setup lang="ts">
import { ref } from 'vue'
import { showNotify } from 'vant'
import { useRoute, useRouter } from 'vue-router'
import useUserStore from '../../store/modules/user'
import pinia from '../../store'
// import useUserStore from '@/store/modules/user'

let username = ref<string>('test')
let password = ref<string>('test111')
let $router = useRouter()
let $route = useRoute()
let userStore = useUserStore(pinia)

const onSubmit = async () => {
  console.log('submit', username, password)
  // TODO(zhenxiang@) validate the form
  try {
    await userStore.login(username.value, password.value)
    //编程式导航跳转到展示数据首页
    //判断登录的时候,路由路径当中是否有query参数，如果有就往query参数挑战，没有跳转到首页
    let redirect: any = $route.query.redirect
    $router.push({ path: redirect || '/' })

    showNotify({ type: 'success', message: 'Login Success' })
  } catch (error) {
    showNotify({ type: 'danger', message: 'Login Failed' })
  }
}
</script>

<template>
  <van-form @submit="onSubmit">
    <van-cell-group inset>
      <van-field
        v-model="username"
        name="用户名"
        label="用户名"
        placeholder="用户名"
        :rules="[{ required: true, message: '请填写用户名' }]"
      />
      <van-field
        v-model="password"
        type="password"
        name="密码"
        label="密码"
        placeholder="密码"
        :rules="[{ required: true, message: '请填写密码' }]"
      />
    </van-cell-group>
    <div style="margin: 16px">
      <van-button round block type="primary" native-type="submit">
        提交
      </van-button>
    </div>
  </van-form>
</template>

<style lang="scss" scoped></style>
