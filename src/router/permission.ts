// 路由守卫:
// 路由鉴权
import useUserStore from '../store/modules/user';
import router from '../router';

const userStore = useUserStore()
const username = userStore.username
router.beforeEach(async (to:any, from: any, next:any) => {
    // update page title
    const token = userStore
    if (token) {
        if (to.path === '/login') {
            next({path: '/'})
        } else {
            if (username) {
                next()
            }
        }
    } else {
        if (to.path == '/login') {
            next()
          } else {
            next({ path: '/login', query: { redirect: to.path } })
          }
    }
})
// 后置守卫
// router.afterEach(( to: any, from: any ) => {})