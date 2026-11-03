import { createRouter, createWebHistory, type RouteRecordRaw } from 'vue-router'
import { useAppStore } from '@/stores/app'

/**
 * 用真实路由而不是 activeTab 状态机 —— 这是刻意与 dockerCopilot 拉开差距的一点：
 * 可以深链、可以前进后退、刷新后停在原页。
 */
const routes: RouteRecordRaw[] = [
  { path: '/', redirect: '/overview' },
  {
    path: '/login',
    name: 'login',
    component: () => import('@/views/LoginView.vue'),
    meta: { public: true, title: '登录' },
  },
  {
    path: '/overview',
    name: 'overview',
    component: () => import('@/views/OverviewView.vue'),
    meta: { title: '总览', icon: 'gauge' },
  },
  {
    path: '/containers',
    name: 'containers',
    component: () => import('@/views/ContainersView.vue'),
    meta: { title: '容器', icon: 'box' },
  },
  {
    path: '/containers/:name',
    name: 'container-detail',
    component: () => import('@/views/ContainerDetailView.vue'),
    meta: { title: '容器详情', hidden: true },
  },
  {
    path: '/images',
    name: 'images',
    component: () => import('@/views/ImagesView.vue'),
    meta: { title: '镜像', icon: 'layers' },
  },
  {
    path: '/schedules',
    name: 'schedules',
    component: () => import('@/views/SchedulesView.vue'),
    meta: { title: '计划任务', icon: 'calendar' },
  },
  {
    // 顺序启动：侧栏与页面标题同名，位置紧跟「计划任务」
    path: '/boot',
    name: 'boot',
    component: () => import('@/views/BootView.vue'),
    meta: { title: '顺序启动', icon: 'list-ordered' },
  },
  {
    path: '/registries',
    name: 'registries',
    component: () => import('@/views/RegistriesView.vue'),
    meta: { title: '加速源', icon: 'rocket' },
  },
  {
    path: '/backup',
    name: 'backup',
    component: () => import('@/views/BackupView.vue'),
    meta: { title: '备份与恢复', icon: 'archive' },
  },
  {
    path: '/notify',
    name: 'notify',
    component: () => import('@/views/NotifyView.vue'),
    meta: { title: '通知', icon: 'bell' },
  },
  {
    path: '/settings',
    name: 'settings',
    component: () => import('@/views/SettingsView.vue'),
    meta: { title: '设置', icon: 'settings' },
  },
  {
    path: '/about',
    name: 'about',
    component: () => import('@/views/AboutView.vue'),
    meta: { title: '关于', icon: 'info' },
  },
  { path: '/:pathMatch(.*)*', redirect: '/overview' },
]

const router = createRouter({
  history: createWebHistory(),
  routes,
  scrollBehavior: () => ({ top: 0 }),
})

/**
 * 登录后要回到的地址。
 *
 * 只接受站内绝对路径：`//evil.com`（协议相对）与 `http://…` 都是开放重定向，
 * 拿去拼「登录后跳转」就会被钓鱼用（受害者看到的是自家的登录页，跳走的是别人的站）。
 */
export function safeRedirect(raw: unknown): string | null {
  if (typeof raw !== 'string' || !raw) return null
  if (!raw.startsWith('/') || raw.startsWith('//')) return null
  return raw
}

/**
 * 路由守卫：**先判登录态，再渲染页面**。
 *
 * 之前没有守卫，未登录直接访问 /containers 时，组件会先渲染一遍（没有外壳、
 * 数据空白），请求 401 之后才被弹回登录页 —— 一次可见的闪烁，而且内部页面的
 * 骨架先露了出来。现在由 main.ts 在挂载前拿到登录态，守卫据此直接改道。
 */
router.beforeEach((to) => {
  const app = useAppStore()
  const isPublic = to.meta.public === true

  if (!app.loggedIn && !isPublic) {
    return {
      name: 'login',
      query: to.fullPath && to.fullPath !== '/' ? { redirect: to.fullPath } : {},
    }
  }
  // 已经登录还去登录页：直接送回去，别让它闪一下表单
  if (app.loggedIn && to.name === 'login') {
    const back = safeRedirect(to.query.redirect)
    return back ? { path: back } : { path: '/overview' }
  }
  return true
})

router.afterEach((to) => {
  const title = to.meta.title as string | undefined
  document.title = title ? `${title} · Dockhelm` : 'Dockhelm · 容器舵手'
})

export default router
