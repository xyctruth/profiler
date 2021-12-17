import { createApp } from 'vue'
import ElementPlus from 'element-plus'
import zhCn from 'element-plus/es/locale/lang/zh-cn'
import 'element-plus/dist/index.css'
import App from './App.vue'
import router from "./router";
import vuex from "./store";
import axios from './utils/request.js'

const app = createApp(App)

app.use(ElementPlus, {
    locale: zhCn,
})
app.use(router)
app.use(vuex)
app.config.globalProperties.$axios = axios
app.mount('#app')
