import { createApp } from 'vue'
import { createPinia } from 'pinia' // 确保这一行不再报错
import App from './App.vue'
import router from './router'
import './styles/index.css' // 包含 Tailwind 的全局样式

const app = createApp(App)
const pinia = createPinia()

app.use(pinia)   // 注册 Pinia
app.use(router)  // 注册 路由
app.mount('#app')