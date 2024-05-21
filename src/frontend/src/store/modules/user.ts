import { defineStore } from 'pinia'
import { GET_TOKEN, REMOVE_TOKEN, SET_TOKEN } from '../../utils/token'
import { reqLogin } from '../../api/user'
import { loginFormData, loginResponseData } from '../../api/user/type'

const useUserStore = defineStore('User', {
  state: () => {
    return {
      token: GET_TOKEN(),
      displayName: '',
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
        SET_TOKEN(result.data.token)
        this.username = result.data.username
        this.displayName = result.data.displayName
        this.avatar = result.data.avatar
        this.relation = result.data.relation
      }
    },
    // 重置用户信息
    resetUserInfo() {
      REMOVE_TOKEN()
      this.displayName = ''
      this.avatar = ''
      this.relation = ''
    },
  },
  getters: {},
})

export default useUserStore
