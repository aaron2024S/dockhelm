import { createRouter, createWebHistory, type RouteRecordRaw } from 'vue-router'

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
    path: '/updates',
    name: 'updates',
    component: () => import('@/views/UpdatesView.vue'),
    meta: { title: '更新中心', icon: 'download' },
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
    path: '/networks',
    name: 'networks',
    component: () => import('@/views/NetworksView.vue'),
    meta: { title: '网络与端口', icon: 'network' },
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

router.afterEach((to) => {
  const title = to.meta.title as string | undefined
  document.title = title ? `${title} · Dockhelm` : 'Dockhelm · 容器舵手'
})

export default router
