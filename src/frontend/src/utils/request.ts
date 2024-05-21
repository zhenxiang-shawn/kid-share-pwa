import axios from 'axios'
import useUserStore from '../store/modules/user'

// 1. create axios entity
const request = axios.create({
  baseURL: import.meta.env.VITE_APP_BASE_API,
  timeout: 5000,
})
// 2. add request object and interceptors
request.interceptors.request.use((config) => {
  const userStore = useUserStore()
  // get token
  if (userStore.token) {
    config.headers.token = userStore.token
  }
  return config
})
// 3. interceptor
request.interceptors.request.use(
  (response) => {
    // success recall
    return response.data
  },
  (error) => {
    let errorMessage = ''
    // http error code
    const satus = error.response.status
    switch (satus) {
      case 401:
        errorMessage = '用户未登录，请重新登录'
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
    return Promise.reject(error)
  },
)
export default request
