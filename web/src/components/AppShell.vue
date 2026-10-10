<script setup lang="ts">
import { computed, onMounted, onUnmounted, ref, watch, type Component } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import {
  Archive,
  Bell,
  Box,
  CalendarClock,
  ChevronsLeft,
  ChevronsRight,
  Download,
  Gauge,
  Info,
  Layers,
  LogOut,
  Menu,
  Monitor,
  Moon,
  RefreshCw,
  Rocket,
  Search,
  Settings as SettingsIcon,
  Sun,
} from 'lucide-vue-next'
import { useAppStore } from '@/stores/app'
import { useThemeStore, type ThemeMode } from '@/stores/theme'
import { useToastStore } from '@/stores/toast'
import { api } from '@/api/client'
import type { BackupStats, OverviewResponse, Schedule } from '@/api/types'

const app = useAppStore()
const theme = useThemeStore()
const toast = useToastStore()
const route = useRoute()
const router = useRouter()

/** 主题按钮：图标即当前模式，点一下切下一档。 */
const THEME_ICON: Record<ThemeMode, Component> = { dark: Moon, light: Sun, system: Monitor }
const themeIcon = computed(() => THEME_ICON[theme.mode])
const themeTip = computed(() => {
  const now = theme.mode === 'system' ? `跟随系统（当前${theme.effective === 'dark' ? '深色' : '浅色'}）` : theme.label
  return `主题：${now} · 点击切换为「${theme.nextLabel}」`
})

/** 数量角标的取值来源。 */
type CountKey = 'containers' | 'images' | 'updates' | 'schedules' | 'snapshots'

interface NavItem {
  name: string
  label: string
  icon: Component
  count?: CountKey
  /** 非空表示这一项前面要插一个分组标题（设计稿里的「资源 / 系统」）。 */
  group?: string
}

/**
 * 导航顺序、分组与数量角标照设计稿来，共 10 项。
 *
 * 与设计稿的两处刻意差别：
 *
 * ① **没有「网络与端口」这一项**。设计稿里它只是一个侧栏条目、没有对应屏；
 *    而端口信息在容器页逐个展示更有用（每个容器的端口就摆在它自己卡片上，
 *    不用跨页对照）。因此并进了容器页。
 *
 * ② **「通知」与「关于」从设置页底部提到侧栏单列**。设计稿把它们当成设置的子页，
 *    唯一的入口是设置页底部一行 11.5px 的灰字 —— 字号那么小、又贴着页面最底，
 *    实际等于藏起来。它们是两个完整的独立功能面，单列之后能直接点到，
 *    也省掉了「进通知页就把『设置』改名成『设置 · 通知』」那套花招。
 */
const nav: NavItem[] = [
  { name: 'overview', label: '总览', icon: Gauge },
  { name: 'containers', label: '容器', icon: Box, count: 'containers' },
  { name: 'images', label: '镜像', icon: Layers, count: 'images' },
  { name: 'updates', label: '更新中心', icon: Download, count: 'updates' },
  { name: 'schedules', label: '计划任务', icon: CalendarClock, count: 'schedules' },
  { name: 'backup', label: '备份与恢复', icon: Archive, count: 'snapshots', group: '资源' },
  { name: 'registries', label: '镜像加速源', icon: Rocket, group: '系统' },
  { name: 'notify', label: '通知', icon: Bell },
  { name: 'settings', label: '设置', icon: SettingsIcon },
  { name: 'about', label: '关于', icon: Info },
]

const collapsed = ref(localStorage.getItem('dockhelm.side') === 'collapsed')
const mobileOpen = ref(false)
const searchText = ref('')
const hostName = ref('')
const counts = ref<Record<CountKey, number>>({
  containers: 0,
  images: 0,
  updates: 0,
  schedules: 0,
  snapshots: 0,
})

watch(collapsed, (v) => localStorage.setItem('dockhelm.side', v ? 'collapsed' : 'open'))

/** 顶栏中间那行小字（设计稿每屏都不一样：设置类页面统一带「设置 · 」前缀）。 */
const contextLabel = computed(() => {
  switch (route.name) {
    case 'overview':
      return hostName.value ? `Docker 控制台 · ${hostName.value}` : 'Docker 控制台'
    case 'container-detail':
      return `容器 · ${String(route.params.name ?? '')}`
    case 'registries':
      return '设置 · 镜像加速源'
    default:
      return (route.meta.title as string) ?? 'Dockhelm'
  }
})

/** 侧栏 / 顶栏共用的「某一项当前是否高亮」。 */
function isActive(item: NavItem) {
  if (item.name === 'containers') {
    return route.name === 'containers' || route.name === 'container-detail'
  }
  return route.name === item.name
}

/**
 * 顶栏搜索框的占位文案。
 *
 * 以前它随页面变（「搜索镜像…」「搜索快照…」「搜索任务…」），但回车**一律**跳到
 * 容器列表 —— 在镜像页输入「nginx」回车会跑到容器列表去按容器名筛，文案与行为对不上。
 * 现在文案只说它真正能做的事，细节写在 title 里。
 */
const searchPlaceholder = '搜索容器…'

function countOf(key?: CountKey) {
  if (!key) return 0
  return counts.value[key] ?? 0
}

/** 侧栏数量角标：更新中心恒为琥珀色（提醒），其余跟随是否高亮。 */
function countClass(item: NavItem, n: number) {
  if (item.count === 'updates' && n > 0) {
    return 'bg-soft-warn text-warn-text'
  }
  return isActive(item) ? 'bg-accent-soft text-accent-text' : 'bg-ink-700 text-text-5'
}

/**
 * 拉一次侧栏角标。三个接口并行，任何一个失败都不影响其它。
 *
 * 带序号：onMounted、每 60 秒的定时器、以及每次换页都会触发它，
 * 三者可能并发。不带序号时先发的响应会盖掉后发的，角标数字凭空回退。
 */
let countsSeq = 0

async function loadCounts() {
  const seq = ++countsSeq
  const [ov, sch, bk] = await Promise.allSettled([
    api.get<OverviewResponse>('/api/overview'),
    api.get<{ schedules: Schedule[] }>('/api/schedules'),
    api.get<BackupStats>('/api/backups/stats'),
  ])
  if (seq !== countsSeq) return // 已有更新的一次在跑，这次结果作废
  if (ov.status === 'fulfilled') {
    counts.value.containers = ov.value.containers.total
    counts.value.images = ov.value.images.total
    // 角标数的是「有几个镜像待更新」而不是「几个容器」—— 原设计稿写的是「几台」，
    // 2026-10-09 起全站统一用「个」计量容器
    counts.value.updates = new Set((ov.value.updates.items ?? []).map((i) => i.image)).size
    hostName.value = ov.value.docker?.name ?? ''
    app.dockerOnline = !ov.value.dockerError
  }
  if (sch.status === 'fulfilled') counts.value.schedules = (sch.value.schedules ?? []).length
  if (bk.status === 'fulfilled') counts.value.snapshots = bk.value.snapshots ?? 0
}

/** 顶栏的 Docker 在线状态：靠 /api/health 轻量探活。 */
async function updateStatus() {
  try {
    const res = await api.get<{ docker: boolean }>('/api/health')
    app.dockerOnline = res.docker
  } catch {
    app.dockerOnline = false
  }
}

let statusTimer = 0
let countTimer = 0

onMounted(() => {
  void loadCounts()
  void updateStatus()
  statusTimer = window.setInterval(updateStatus, 30000)
  countTimer = window.setInterval(loadCounts, 60000)
})

onUnmounted(() => {
  window.clearInterval(statusTimer)
  window.clearInterval(countTimer)
})

// 换页时刷一次角标 —— 刚建完任务、删完快照，回到列表就能看到数字变了。
watch(() => route.fullPath, () => void loadCounts())

/** 顶栏搜索：回车直接跳到容器列表并带上关键字。 */
function submitSearch() {
  const q = searchText.value.trim()
  router.push({ name: 'containers', query: q ? { q } : {} })
  searchText.value = '' // 提交后清空，免得下次回车又搜同一个词
  mobileOpen.value = false
}

async function doLogout() {
  await app.logout()
  toast.info('已退出登录')
  router.push({ name: 'login' })
}

function go(name: string) {
  router.push({ name })
  mobileOpen.value = false
}

const checking = ref(false)

/**
 * 顶栏「检查更新」：**只读检测，永不更新容器**。
 *
 * 这是全站唯一的手动检测入口 —— 容器页、更新中心页头、总览空态原先各有一颗
 * 同功能按钮（三种叫法），已全部删除，避免同一个动作散在四处。
 * 真正会动容器的只有更新中心的「立即执行自动更新」。
 *
 * 检测完成后：① 刷侧栏角标；② 自增 app.checkTick，让当前页重载自己的检测结果。
 */
async function quickCheck() {
  if (checking.value) return
  checking.value = true
  try {
    const res = await api.post<{ results?: { status: string }[] }>('/api/updates/check', {})
    const n = (res.results ?? []).filter((r) => r.status === 'update_available').length
    toast.success('检查完成', n ? `发现 ${n} 个容器有可用更新` : '所有容器都是最新的')
    await loadCounts()
    app.checkTick++
  } catch (e) {
    toast.error('检查失败', e instanceof Error ? e.message : String(e))
  } finally {
    checking.value = false
  }
}
</script>

<template>
  <div class="flex h-screen flex-col overflow-hidden bg-ink-950">
    <!-- 顶栏：全宽，左起 logo / 品牌 / 当前页，右侧搜索与账户 -->
    <header
      class="flex h-[52px] flex-none items-center gap-2.5 border-b border-line-2 bg-ink-850 px-[18px]"
    >
      <button
        type="button"
        class="dh-iconbtn hidden max-md:grid"
        title="打开导航"
        @click="mobileOpen = true"
      >
        <Menu class="h-4 w-4" />
      </button>

      <div
        class="grid h-[27px] w-[27px] flex-none place-items-center rounded-[9px] bg-accent text-accent-ink"
      >
        <svg viewBox="0 0 32 32" class="h-[15px] w-[15px]">
          <circle cx="16" cy="16" r="12.3" fill="none" stroke="currentColor" stroke-width="2.6" />
          <path
            fill-rule="evenodd"
            fill="currentColor"
            d="M16 6.6 18.5 13.5 25.4 16 18.5 18.5 16 25.4 13.5 18.5 6.6 16 13.5 13.5ZM18 16A2 2 0 1 0 14 16A2 2 0 1 0 18 16Z"
          />
        </svg>
      </div>
      <div class="flex-none text-[14px] font-semibold tracking-[0.2px]">Dockhelm</div>
      <div class="min-w-0 truncate text-[12px] text-text-4">{{ contextLabel }}</div>

      <div class="ml-auto flex items-center gap-2">
        <span
          v-if="!app.dockerOnline"
          class="dh-badge dh-badge-err whitespace-nowrap"
          title="无法连接 Docker 守护进程"
        >
          <span class="h-[7px] w-[7px] rounded-full bg-err" />Docker 离线
        </span>

        <div class="relative max-md:hidden">
          <Search class="pointer-events-none absolute left-2.5 top-1/2 h-3.5 w-3.5 -translate-y-1/2 text-text-5" />
          <input
            v-model="searchText"
            class="dh-input !w-[190px] !py-[6px] !pl-8"
            :placeholder="searchPlaceholder"
            title="回车跳到容器列表，按关键字筛选"
            @keyup.enter="submitSearch"
          />
        </div>

        <!-- 主题：深色 / 浅色 / 跟随系统 三档循环，图标与文字都表示当前档 -->
        <button type="button" class="dh-topbtn" :title="themeTip" @click="theme.cycle()">
          <component :is="themeIcon" class="h-[14px] w-[14px] flex-none" />
          <span class="max-md:hidden">{{ theme.label }}</span>
        </button>

        <!-- 全站唯一的手动检测入口：只检查有没有新版本，绝不会更新容器 -->
        <button
          type="button"
          class="dh-topbtn"
          title="检查有没有新版本（只读，不会更新任何容器）"
          :disabled="checking"
          @click="quickCheck"
        >
          <RefreshCw class="h-[14px] w-[14px] flex-none" :class="checking ? 'dh-spin' : ''" />
          <span class="max-md:hidden">检查更新</span>
        </button>

        <button type="button" class="dh-topbtn" title="退出登录" @click="doLogout">
          <LogOut class="h-[14px] w-[14px] flex-none" />
          <span class="max-md:hidden">退出</span>
        </button>
      </div>
    </header>

    <div class="flex min-h-0 flex-1">
      <!-- 侧边栏：纯导航，分组 + 数量角标；折叠开关钉在最底下 -->
      <aside
        class="z-40 flex flex-none flex-col border-r border-line-2 bg-ink-875 px-2.5 py-3.5 transition-[width] duration-200 max-md:fixed max-md:inset-y-0 max-md:left-0 max-md:w-[172px] max-md:transition-transform"
        :class="[
          collapsed ? 'w-[62px]' : 'w-[172px]',
          mobileOpen ? 'max-md:translate-x-0' : 'max-md:-translate-x-full',
        ]"
      >
        <nav class="dh-scroll flex min-h-0 flex-1 flex-col gap-0.5 overflow-y-auto">
          <template v-for="item in nav" :key="item.name">
            <div v-if="item.group && !collapsed" class="dh-nav-group">{{ item.group }}</div>
            <div v-else-if="item.group" class="my-1 h-px flex-none bg-line-2" />

            <button
              type="button"
              class="flex w-full items-center gap-[9px] rounded-[9px] px-2.5 py-2 text-left text-[13px] transition-colors"
              :class="
                isActive(item)
                  ? 'bg-accent-soft text-accent'
                  : 'text-text-3 hover:bg-ink-700 hover:text-text-1'
              "
              :title="collapsed ? item.label : undefined"
              @click="go(item.name)"
            >
              <component :is="item.icon" class="h-[15px] w-[15px] flex-none" :stroke-width="2" />
              <span v-if="!collapsed" class="truncate">{{ item.label }}</span>
              <span
                v-if="!collapsed && item.count && countOf(item.count) > 0"
                class="ml-auto flex-none rounded-[6px] px-1.5 text-[11px] leading-[18px]"
                :class="countClass(item, countOf(item.count))"
              >
                {{ countOf(item.count) }}
              </span>
            </button>
          </template>
        </nav>

        <!--
          折叠开关：钉在侧栏最底部并配文字。
          原先它在顶栏、只有一对箭头图标（‹‹），既不好找也看不出是干什么的；
          侧栏底部是「导航的收尾位置」，摆在这里语义自明。
          收窄态只剩图标（宽度 62px 放不下字），靠 title 兜底。
          移动端不显示：那时侧栏是抽屉，开关由顶栏的汉堡按钮负责。
        -->
        <div class="mt-1 hidden flex-none border-t border-line-2 pt-1.5 md:block">
          <button
            type="button"
            class="flex w-full items-center gap-[9px] rounded-[9px] px-2.5 py-2 text-left text-[12.5px] text-text-4 transition-colors hover:bg-ink-700 hover:text-text-1"
            :class="collapsed ? 'justify-center' : ''"
            :title="collapsed ? '展开侧栏' : '收起侧栏'"
            @click="collapsed = !collapsed"
          >
            <component
              :is="collapsed ? ChevronsRight : ChevronsLeft"
              class="h-[15px] w-[15px] flex-none"
              :stroke-width="2"
            />
            <span v-if="!collapsed" class="truncate">收起侧栏</span>
          </button>
        </div>
      </aside>

      <!-- 移动端遮罩 -->
      <div
        v-if="mobileOpen"
        class="fixed inset-0 z-30 bg-black/55 md:hidden"
        @click="mobileOpen = false"
      />

      <!-- 内容 -->
      <main class="dh-scroll min-h-0 flex-1 bg-ink-900">
        <slot />
      </main>
    </div>
  </div>
</template>
