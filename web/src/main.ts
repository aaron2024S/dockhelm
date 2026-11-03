import { createApp } from 'vue'
import { createPinia } from 'pinia'
import App from './App.vue'
import router from './router'
import { setUnauthorizedHandler } from './api/client'
import { useAppStore } from './stores/app'
import { useThemeStore } from './stores/theme'
import './style.css'

const pinia = createPinia()
const app = createApp(App)
app.use(pinia)

// 主题对齐：<html> 上的 dark/light 已由 /theme-init.js 设好（防首屏闪烁），
// 这里把 store 状态与它同步，并注册「跟随系统」时对系统主题变化的监听。
useThemeStore(pinia).init()

/**
 * 会话失效时的统一处理：当然回不到原页面，但要记住「本来想去哪」。
 *
 * 走路由而不是 location.replace：后者会整页重载、丢掉已加载的状态，
 * 而且登录完只能落到总览 —— 从容器详情页被踢出去，回来要重新点三层。
 */
setUnauthorizedHandler(() => {
  const store = useAppStore(pinia)
  store.loggedIn = false
  const cur = router.currentRoute.value
  if (cur.name === 'login') return
  void router.replace({
    name: 'login',
    query: cur.fullPath && cur.fullPath !== '/' ? { redirect: cur.fullPath } : {},
  })
})

/**
 * 先把登录态拿到手，再装路由、再挂载。
 *
 * 顺序很关键：路由守卫要判断「有没有登录」，而登录态来自 /api/session。
 * 以前是挂载之后才在 onMounted 里取，于是直接访问 /containers 时会先按
 * 「已登录」把内部页面渲染出来（没有外壳、数据空白），随后才被 401 弹回去 ——
 * 用户看到的是一次明显的闪烁，而且页面上不该出现的骨架先露了出来。
 */
async function start() {
  const store = useAppStore(pinia)
  await store.bootstrap()
  app.use(router)
  await router.isReady()
  app.mount('#app')
}

void start()
