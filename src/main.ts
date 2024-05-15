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

app.use(router)
app.use(Vant)
app.use(pinia)
app.mount('#app')
