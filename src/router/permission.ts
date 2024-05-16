// 路由守卫:
// 路由鉴权
import useUserStore from '../store/modules/user'
import router from '../router'
import pinia from '../store'

const userStore = useUserStore(pinia)

router.beforeEach((to: any, from: any, next: any) => {
  console.log('before each')
  // update page title

  // get token
  const user_token = userStore.token
  const username = userStore.username
  if (user_token) {
    if (to.path == '/login') {
      next({ path: '/' })
    } else {
      if (username) {
        next()
      } else {
        // fetch user info
        // 如果 token 过期,登出,并重新定向到登录页面
        next()
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
