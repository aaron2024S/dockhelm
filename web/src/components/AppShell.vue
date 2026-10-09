<script setup lang="ts">
import { computed, onMounted, onUnmounted, ref, watch, type Component } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import {
  Archive,
  Box,
  CalendarClock,
  ChevronsLeft,
  ChevronsRight,
  Download,
  Gauge,
  Layers,
  LogOut,
  Menu,
  Network,
  RefreshCw,
  Rocket,
  Search,
  Settings as SettingsIcon,
} from 'lucide-vue-next'
import { useAppStore } from '@/stores/app'
import { useToastStore } from '@/stores/toast'
import { api } from '@/api/client'
import type { BackupStats, OverviewResponse, Schedule } from '@/api/types'

const app = useAppStore()
const toast = useToastStore()
const route = useRoute()
const router = useRouter()

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
 * 导航顺序、分组与数量角标照设计稿来（设计稿侧栏只有这 9 项）：
 * 五个常用页在最上面，其次「资源」，最后「系统」。
 * 「通知」是「设置」的子页 —— 进入后标签变成「设置 · 通知」，
 * 因此这里不要把通知单列成一项。「关于」不在设计稿侧栏里，
 * 入口放在设置页底部。
 */
const nav: NavItem[] = [
  { name: 'overview', label: '总览', icon: Gauge },
  { name: 'containers', label: '容器', icon: Box, count: 'containers' },
  { name: 'images', label: '镜像', icon: Layers, count: 'images' },
  { name: 'updates', label: '更新中心', icon: Download, count: 'updates' },
  { name: 'schedules', label: '计划任务', icon: CalendarClock, count: 'schedules' },
  { name: 'backup', label: '备份与恢复', icon: Archive, count: 'snapshots', group: '资源' },
  { name: 'networks', label: '网络与端口', icon: Network },
  { name: 'registries', label: '镜像加速源', icon: Rocket, group: '系统' },
  { name: 'settings', label: '设置', icon: SettingsIcon },
]

/** 侧栏条目在特定页面上要换一个更具体的名字（设计稿：设置 · 通知）。 */
const LABEL_OVERRIDE: Record<string, string> = { notify: '设置 · 通知' }

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
    case 'notify':
      return '设置 · 通知'
    default:
      return (route.meta.title as string) ?? 'Dockhelm'
  }
})

/** 侧栏 / 顶栏共用的「某一项当前是否高亮」。 */
function isActive(item: NavItem) {
  if (item.name === 'containers') {
    return route.name === 'containers' || route.name === 'container-detail'
  }
  // 通知是设置页的子页，进通知时「设置」这一项保持高亮（设计稿如此）
  if (item.name === 'settings') {
    return route.name === 'settings' || route.name === 'notify'
  }
  return route.name === item.name
}

/** 侧栏显示用的文字：通知页下「设置」要显示成「设置 · 通知」。 */
function labelOf(item: NavItem) {
  if (item.name === 'settings' && route.name === 'notify') return LABEL_OVERRIDE.notify as string
  return item.label
}

/** 顶栏搜索框的占位文案也随页面变（设计稿：搜索容器… / 搜索镜像… / 搜索快照…）。 */
const searchPlaceholder = computed(() => {
  switch (route.name) {
    case 'containers':
    case 'container-detail':
      return '搜索容器…'
    case 'images':
      return '搜索镜像…'
    case 'updates':
      return '搜索镜像…'
    case 'schedules':
      return '搜索任务…'
    case 'backup':
      return '搜索快照…'
    case 'overview':
      return '搜索容器、镜像、任务…'
    default:
      return '搜索设置…'
  }
})

function countOf(key?: CountKey) {
  if (!key) return 0
  return counts.value[key] ?? 0
}

/** 侧栏数量角标：更新中心恒为琥珀色（提醒），其余跟随是否高亮。 */
function countClass(item: NavItem, n: number) {
  if (item.count === 'updates' && n > 0) {
    return 'bg-[rgba(245,165,36,.15)] text-[#fbbf24]'
  }
  return isActive(item) ? 'bg-[#0c3230] text-[#5eead4]' : 'bg-ink-700 text-text-5'
}

/** 拉一次侧栏角标。三个接口并行，任何一个失败都不影响其它。 */
async function loadCounts() {
  const [ov, sch, bk] = await Promise.allSettled([
    api.get<OverviewResponse>('/api/overview'),
    api.get<{ schedules: Schedule[] }>('/api/schedules'),
    api.get<BackupStats>('/api/backups/stats'),
  ])
  if (ov.status === 'fulfilled') {
    counts.value.containers = ov.value.containers.total
    counts.value.images = ov.value.images.total
    // 角标数的是「有几个镜像待更新」而不是「几台容器」—— 与设计稿一致
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
async function quickCheck() {
  if (checking.value) return
  checking.value = true
  try {
    const res = await api.post<{ results: unknown[] }>('/api/updates/check', {})
    toast.success('巡检完成', `共检查 ${res.results?.length ?? 0} 个容器`)
    await loadCounts()
  } catch (e) {
    toast.error('巡检失败', e instanceof Error ? e.message : String(e))
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
        <svg viewBox="0 0 32 32" class="h-[15px] w-[15px]" fill="none" stroke="currentColor" stroke-width="2.6">
          <path d="M16 7l7 4v10l-7 4-7-4V11z" stroke-linejoin="round" />
          <path d="M16 15v10M9 11l7 4 7-4" stroke-linejoin="round" />
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
          <span class="h-[7px] w-[7px] rounded-full bg-[#f87171]" />Docker 离线
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

        <button
          type="button"
          class="dh-iconbtn"
          title="立即巡检更新（只读，不会动任何容器）"
          :disabled="checking"
          @click="quickCheck"
        >
          <RefreshCw class="h-[14px] w-[14px]" :class="checking ? 'dh-spin' : ''" />
        </button>

        <button
          type="button"
          class="dh-iconbtn max-md:hidden"
          :title="collapsed ? '展开侧栏' : '收起侧栏'"
          @click="collapsed = !collapsed"
        >
          <component :is="collapsed ? ChevronsRight : ChevronsLeft" class="h-[14px] w-[14px]" />
        </button>

        <button type="button" class="dh-avatar" title="账户设置" @click="go('settings')">A</button>

        <button type="button" class="dh-iconbtn" title="退出登录" @click="doLogout">
          <LogOut class="h-[14px] w-[14px]" />
        </button>
      </div>
    </header>

    <div class="flex min-h-0 flex-1">
      <!-- 侧边栏：纯导航，分组 + 数量角标 -->
      <aside
        class="z-40 flex flex-none flex-col gap-0.5 overflow-y-auto border-r border-line-2 bg-ink-875 px-2.5 py-3.5 transition-[width] duration-200 max-md:fixed max-md:inset-y-0 max-md:left-0 max-md:w-[172px] max-md:transition-transform"
        :class="[
          collapsed ? 'w-[62px]' : 'w-[172px]',
          mobileOpen ? 'max-md:translate-x-0' : 'max-md:-translate-x-full',
        ]"
      >
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
            :title="collapsed ? labelOf(item) : undefined"
            @click="go(item.name)"
          >
            <component :is="item.icon" class="h-[15px] w-[15px] flex-none" :stroke-width="2" />
            <span v-if="!collapsed" class="truncate">{{ labelOf(item) }}</span>
            <span
              v-if="!collapsed && item.count && countOf(item.count) > 0"
              class="ml-auto flex-none rounded-[6px] px-1.5 text-[11px] leading-[18px]"
              :class="countClass(item, countOf(item.count))"
            >
              {{ countOf(item.count) }}
            </span>
          </button>
        </template>
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
