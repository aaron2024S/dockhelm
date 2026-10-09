import { createApp } from 'vue'
import { createPinia } from 'pinia'
import App from './App.vue'
import router from './router'
import { useThemeStore } from './stores/theme'
import './style.css'

const pinia = createPinia()
const app = createApp(App)
app.use(pinia)
app.use(router)

// 主题对齐：<html> 上的 dark/light 已由 index.html 的内联脚本设好（防首屏闪烁），
// 这里把 store 状态与它同步，并注册「跟随系统」时对系统主题变化的监听。
useThemeStore(pinia).init()

app.mount('#app')
