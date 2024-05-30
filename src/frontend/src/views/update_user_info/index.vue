<script setup lang="ts">
import { onMounted, ref } from 'vue'
import useUserStore from '../../store/modules/user'

const onClickLeft = () => history.back()
let userStore = useUserStore()
let username = ref<string>(userStore.username!)
let display_name = ref<string>(userStore.displayName!)
let relationship = ref<string>(userStore.relation)
let password = ref<string>('666666')
let showPicker = ref<boolean>(false)

let titles = [
  { text: '妈妈', value: '妈妈' },
  { text: '爸爸', value: '爸爸' },
  { text: '爷爷', value: '爷爷' },
  { text: '奶奶', value: '奶奶' },
  { text: '外公', value: '外公' },
  { text: '外婆', value: '外婆' },
  { text: '叔叔', value: '叔叔' },
  { text: '小姨', value: '小姨' },
  { text: '哥哥', value: '哥哥' },
  { text: '姐姐', value: '姐姐' },
  { text: '弟弟', value: '弟弟' },
  { text: '妹妹', value: '妹妹' },
  { text: '姨妈', value: '姨妈' },
  { text: '舅舅', value: '舅舅' },
]

let result = ref<string>('')

const onConfirm = ({ selectedOptions }) => {
  result.value = selectedOptions[0]?.text
  showPicker.value = false
}

onMounted(() => {
  console.log('mounted')
})
</script>

<template>
  <van-nav-bar
    title="个人信息"
    left-text="返回"
    left-arrow
    @click-left="onClickLeft"
  />

  <!-- update avatar -->

  <!-- update display name -->
  <van-cell-group inset class="user_info">
    <!-- 输入任意文本 -->
    <van-field v-model="username" label="用户名" disabled />
    <!-- <van-field v-model="password" type="password" label="密码" /> -->
    <van-field v-model="display_name" label="昵称" disabled />
    <van-field
      v-model="result"
      is-link
      name="picker"
      label="头衔"
      :placeholder="relationship"
      @click="showPicker = true"
    />
    <van-popup v-model:show="showPicker" position="bottom">
      <van-picker
        :columns="titles"
        @confirm="onConfirm"
        @cancel="showPicker = false"
      />
    </van-popup>
    <!-- 输入密码 -->
  </van-cell-group>

  <!-- update relationship -->
</template>

<style lang="scss" scoped>
.user_info {
  margin-top: 10px;
}
</style>
