<script setup lang="ts">
import { ref } from 'vue'
import { showNotify } from 'vant'
import { useRoute, useRouter } from 'vue-router'
import useUserStore from '../../store/modules/user'

let username = ref<string>('visitor')
let password = ref<string>('111111')
let $route = useRoute()
let $router = useRouter()
let useUser = useUserStore()

const onSubmit = () => {
  console.log('submit', username, password)

  // fake account
  if (username.value === 'visitor' && password.value === '111111') {
    useUser.login('visitor')
    showNotify({ type: 'success', message: 'Login Success' })
    //编程式导航跳转到展示数据首页
    //判断登录的时候,路由路径当中是否有query参数，如果有就往query参数挑战，没有跳转到首页
    let redirect: any = $route.query.redirect
    $router.push({ path: redirect || '/' })
  } else {
    showNotify({ type: 'danger', message: '账户密码错误' })
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
