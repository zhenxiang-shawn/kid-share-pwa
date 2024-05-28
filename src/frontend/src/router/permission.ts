// 路由守卫:
// 路由鉴权
import useUserStore from '../store/modules/user'
import router from '.'
import pinia from '../store'
import { reqUserInfo } from '../APP_interface/user'
import type { userInfoReponseData } from '../APP_interface/user/type'

const userStore = useUserStore(pinia)

router.beforeEach(async (to: any, from: any, next: any) => {
  // update page title

  // get token
  const user_token = userStore.token as string
  // console.log('grab token: ', user_token)
  const username = userStore.username as string
  // console.log('grab username: ', username)
  if (user_token) {
    if (to.path == '/login') {
      next({ path: '/' })
    } else {
      if (username) {
        next()
      } else {
        // console.log('No user info found, pass next !!!!')
        // fetch user info
        let result = await reqUserInfo(username)
        if (result.code == 200) {
          console.log('Router_resluting user info:', result.data)
          next({ ...to })
        } else if (result.code == 401) {
          // 如果 token 过期,登出,并重新定向到登录页面
          console.log('Token expired, logout')
          // token 过期,登出
          userStore.logout()
          // 用户未登录,重定向到登录页面
          next({ path: '/login', query: { redirect: to.path } })
        } else {
          console.log('No user info found, But with error:', result.code)
        }

        // next()
      }
    }
  } else {
    console.log('No token found, redirect to login')
    // 用户未登录,重定向到登录页面
    if (to.path == '/login') {
      next()
    } else {
      next({ path: '/login', query: { redirect: to.path } })
    }
  }
})
// 后置守卫
// router.afterEach(( to: any, from: any ) => {})
