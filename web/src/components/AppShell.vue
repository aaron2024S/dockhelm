<script setup lang="ts">
import { computed, onMounted, ref } from 'vue'
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
  RefreshCw,
  Rocket,
  Search,
  Settings as SettingsIcon,
  User,
} from 'lucide-vue-next'
import { useAppStore } from '@/stores/app'
import { useToastStore } from '@/stores/toast'
import { api } from '@/api/client'

const app = useAppStore()
const toast = useToastStore()
const route = useRoute()
const router = useRouter()

const nav = [
  { name: 'overview', label: '总览', icon: Gauge },
  { name: 'containers', label: '容器', icon: Box },
  { name: 'updates', label: '更新中心', icon: Download },
  { name: 'images', label: '镜像', icon: Layers },
  { name: 'schedules', label: '计划任务', icon: CalendarClock },
  { name: 'registries', label: '加速源', icon: Rocket },
  { name: 'backup', label: '备份与恢复', icon: Archive },
  { name: 'notify', label: '通知', icon: Bell },
  { name: 'settings', label: '设置', icon: SettingsIcon },
  { name: 'about', label: '关于', icon: Info },
]

const collapsed = ref(false)
const mobileOpen = ref(false)
const searchText = ref('')
const version = ref('')

onMounted(async () => {
  try {
    const res = await api.get<{ about: { version: string } }>('/api/about')
    version.value = res.about.version
  } catch {
    /* 忽略 */
  }
  void updateStatus()
})

/** 顶栏的 Docker 在线状态：靠 /api/health 轻量探活。 */
async function updateStatus() {
  try {
    const res = await api.get<{ docker: boolean }>('/api/health')
    app.dockerOnline = res.docker
  } catch {
    app.dockerOnline = false
  }
  window.setTimeout(updateStatus, 30000)
}

const currentTitle = computed(() => (route.meta.title as string) ?? 'Dockhelm')

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
  } catch (e) {
    toast.error('巡检失败', e instanceof Error ? e.message : String(e))
  } finally {
    checking.value = false
  }
}
</script>

<template>
  <div class="flex h-screen overflow-hidden bg-ink-950">
    <!-- 侧边栏 -->
    <aside
      class="z-40 flex flex-none flex-col border-r border-line-2 bg-ink-875 transition-[width] duration-200 max-md:fixed max-md:inset-y-0 max-md:left-0"
      :class="[
        collapsed ? 'w-[62px]' : 'w-[196px]',
        mobileOpen ? 'max-md:translate-x-0' : 'max-md:-translate-x-full',
        'max-md:w-[196px] max-md:transition-transform',
      ]"
    >
      <div class="flex h-[52px] flex-none items-center gap-2.5 px-3.5">
        <div class="grid h-[27px] w-[27px] flex-none place-items-center rounded-[9px] bg-accent text-accent-ink">
          <svg viewBox="0 0 32 32" class="h-[15px] w-[15px]" fill="none" stroke="currentColor" stroke-width="2.6">
            <path d="M16 7l7 4v10l-7 4-7-4V11z" stroke-linejoin="round" />
            <path d="M16 15v10M9 11l7 4 7-4" stroke-linejoin="round" />
          </svg>
        </div>
        <div v-if="!collapsed" class="min-w-0">
          <div class="truncate text-[14px] font-semibold tracking-[0.2px]">Dockhelm</div>
          <div class="truncate text-[11px] text-text-4">容器舵手{{ version ? ' · v' + version : '' }}</div>
        </div>
      </div>

      <nav class="flex min-h-0 flex-1 flex-col gap-0.5 overflow-y-auto px-2 pb-2 pt-1">
        <button
          v-for="item in nav"
          :key="item.name"
          type="button"
          class="flex w-full items-center gap-[9px] rounded-[9px] px-2.5 py-[7px] text-left text-[13px] transition-colors"
          :class="
            route.name === item.name ||
            (item.name === 'containers' && route.name === 'container-detail')
              ? 'bg-accent-soft text-accent'
              : 'text-text-3 hover:bg-ink-700 hover:text-text-1'
          "
          :title="collapsed ? item.label : undefined"
          @click="go(item.name)"
        >
          <component :is="item.icon" class="h-[15px] w-[15px] flex-none" :stroke-width="2" />
          <span v-if="!collapsed" class="truncate">{{ item.label }}</span>
        </button>
      </nav>

      <div class="flex-none border-t border-line-2 p-2">
        <button
          type="button"
          class="flex w-full items-center gap-[9px] rounded-[9px] px-2.5 py-[7px] text-[13px] text-text-3 hover:bg-ink-700 hover:text-text-1"
          @click="collapsed = !collapsed"
        >
          <component :is="collapsed ? ChevronsRight : ChevronsLeft" class="h-[15px] w-[15px] flex-none" />
          <span v-if="!collapsed">收起侧栏</span>
        </button>
      </div>
    </aside>

    <!-- 移动端遮罩 -->
    <div
      v-if="mobileOpen"
      class="fixed inset-0 z-30 bg-black/55 md:hidden"
      @click="mobileOpen = false"
    />

    <div class="flex min-w-0 flex-1 flex-col">
      <!-- 顶栏 -->
      <header
        class="flex h-[52px] flex-none items-center gap-3 border-b border-line-2 bg-ink-850 px-[18px]"
      >
        <button
          type="button"
          class="hidden h-[30px] w-[30px] place-items-center rounded-[9px] border border-line-3 text-text-4 max-md:grid"
          @click="mobileOpen = true"
        >
          <svg viewBox="0 0 24 24" class="h-4 w-4" fill="none" stroke="currentColor" stroke-width="2">
            <path d="M4 6h16M4 12h16M4 18h16" stroke-linecap="round" />
          </svg>
        </button>

        <div class="min-w-0 flex-none">
          <div class="truncate text-[15px] font-semibold">{{ currentTitle }}</div>
        </div>

        <div
          v-if="!app.dockerOnline"
          class="dh-badge dh-badge-err whitespace-nowrap"
          title="无法连接 Docker 守护进程"
        >
          <span class="h-[7px] w-[7px] rounded-full bg-[#f87171]" />Docker 离线
        </div>

        <div class="ml-auto flex items-center gap-2">
          <div class="relative max-md:hidden">
            <Search class="pointer-events-none absolute left-2.5 top-1/2 h-3.5 w-3.5 -translate-y-1/2 text-text-5" />
            <input
              v-model="searchText"
              class="dh-input !w-[200px] !py-[6px] !pl-8"
              placeholder="搜索容器…"
              @keyup.enter="submitSearch"
            />
          </div>

          <button
            type="button"
            class="grid h-[30px] w-[30px] place-items-center rounded-[9px] border border-line-3 text-text-4 transition-colors hover:border-line-4 hover:text-text-1 disabled:opacity-40"
            title="立即巡检更新（只读，不会动任何容器）"
            :disabled="checking"
            @click="quickCheck"
          >
            <RefreshCw class="h-[14px] w-[14px]" :class="checking ? 'dh-spin' : ''" />
          </button>

          <RouterLink
            to="/settings"
            class="grid h-[30px] w-[30px] place-items-center rounded-[9px] border border-line-3 text-text-4 transition-colors hover:border-line-4 hover:text-text-1"
            title="设置"
          >
            <User class="h-[14px] w-[14px]" />
          </RouterLink>

          <button
            type="button"
            class="grid h-[30px] w-[30px] place-items-center rounded-[9px] border border-line-3 text-text-4 transition-colors hover:border-line-4 hover:text-text-1"
            title="退出登录"
            @click="doLogout"
          >
            <LogOut class="h-[14px] w-[14px]" />
          </button>
        </div>
      </header>

      <!-- 内容 -->
      <main class="dh-scroll min-h-0 flex-1 bg-ink-900">
        <slot />
      </main>
    </div>
  </div>
</template>
