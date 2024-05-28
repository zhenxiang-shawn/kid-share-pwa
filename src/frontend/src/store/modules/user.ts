import { defineStore } from 'pinia'
import { GET_TOKEN, REMOVE_TOKEN, SET_TOKEN } from '../../utils/token'
import { reqLogin } from '../../APP_interface/user'
import { loginFormData, loginResponseData } from '../../APP_interface/user/type'
import {
  GET_USER_DISPLAY_NAME,
  GET_USER_NAME,
  SET_USER_INFO,
} from '../../utils/user'

const useUserStore = defineStore('User', {
  state: () => {
    return {
      token: GET_TOKEN(),
      username: GET_USER_NAME(),
      displayName: GET_USER_DISPLAY_NAME(),
      avatar: '',
      relation: '',
    }
  },
  actions: {
    // 用户登陆
    async login(username, password) {
      let loginForm: loginFormData = {
        username: username,
        password: password,
      }
      const result: loginResponseData = await reqLogin(loginForm)
      if (result.code === 200) {
        SET_TOKEN(result.data.token as string)
        SET_USER_INFO(
          result.data.username as string,
          result.data.display_name as string,
          result.data.avatar as string,
          result.data.relation as string,
        )
        this.username = result.data.username
        this.displayName = result.data.display_name
        this.avatar = result.data.avatar
        this.relation = result.data.relation
        //能保证当前async函数返回一个成功的promise
        return 'ok'
      } else {
        return Promise.reject(new Error(result.message))
      }
    },
    // 登出
    logout() {
      // 重置用户信息
      REMOVE_TOKEN()
      this.username = ''
      this.displayName = ''
      this.avatar = ''
      this.relation = ''
    },
  },
})

export default useUserStore
