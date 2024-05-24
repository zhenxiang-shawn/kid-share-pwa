import axios from 'axios'
import useUserStore from '../store/modules/user'
import { showNotify } from 'vant'

// 1. create axios entity
const request = axios.create({
  // baseURL: import.meta.env.VITE_APP_BASE_API,
  baseURL: '/api',
  timeout: 5000,
})
// 2. add request interceptors
request.interceptors.request.use((config) => {
  const userStore = useUserStore()
  // get token
  if (userStore.token) {
    config.headers.Authorization = `Bearer ${userStore.token}`
  }
  console.log('Request config:', config)
  console.log('Request URL:', config.url)
  return config
})
// 3. add response interceptor
request.interceptors.response.use(
  (response) => {
    console.log('Response:', response)
    // success recall
    return response.data
  },
  (error) => {
    const userStore = useUserStore()
    console.log('Error:', error)
    let errorMessage = ''
    // http error code
    const satus = error.response.status
    switch (satus) {
      case 401:
        errorMessage = '401: 用户未登录，请重新登录'
        userStore.logout()
        break
      case 403:
        errorMessage = '用户没有权限'
        break
      case 404:
        errorMessage = '请求地址错误'
        break
      default:
        errorMessage = '网络异常，请稍后重试'
        break
    }
    showNotify({ type: 'danger', message: errorMessage })
    return Promise.reject(error)
  },
)
export default request
