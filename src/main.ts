// create App
import { createApp } from 'vue'
import App from './App.vue'
// import global styles
//引入模板的全局的样式
import './style.css'
// import Vant component
import Vant from 'vant'
import 'vant/lib/index.css'
// set pinia
import pinia from './store'
// set up router
import router from './router'

const app = createApp(App)

app.use(pinia)
app.use(router)
//引入路由鉴权文件
import './router/permission'
app.use(Vant)

app.mount('#app')
