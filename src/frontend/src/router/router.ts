import { RouteRecordRaw } from 'vue-router'
export const constRoute: Array<RouteRecordRaw> = [
  // login
  {
    path: '/login',
    // component: () => import('@/views/login/index.vue')
    component: () => import('../views/login/index.vue'),
    name: 'login',
  },
  {
    path: '/',
    component: () => import('../layout/index.vue'),
    name: 'layout',
    redirect: '/home',
    children: [
      {
        path: '/home',
        component: () => import('../views/home/index.vue'),
        name: 'home',
      },
      {
        path: '/post',
        component: () => import('../views/posts/index.vue'),
        name: 'post',
      },
      {
        path: '/my',
        component: () => import('../views/my/index.vue'),
        name: 'my',
      },
      {
        path: '/404',
        component: () => import('../views/404/index.vue'),
        name: '404',
      },
    ],
  },
]
