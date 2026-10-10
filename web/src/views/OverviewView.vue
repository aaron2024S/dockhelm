<script setup lang="ts">
import { computed, onMounted, onUnmounted, ref, watch } from 'vue'
import { useRouter } from 'vue-router'
import { AlertTriangle, ArrowRight } from 'lucide-vue-next'
import { api, openStream } from '@/api/client'
import type { OverviewResponse, RunLog, Schedule } from '@/api/types'
import { formatBytes, formatDayTime, runKindLabel } from '@/utils/format'
import { useAppStore } from '@/stores/app'

const app = useAppStore()
const router = useRouter()
const data = ref<OverviewResponse | null>(null)
const schedules = ref<Schedule[]>([])
const loading = ref(true)
const errorMsg = ref('')
let closeStream: (() => void) | null = null

/**
 * 延时刷新用的定时器。
 * 巡检/批次结束后会延后几秒再拉一次数据；如果这几秒内用户切走了页面，
 * 原实现仍会执行 —— 会往已卸载的组件里写 state，且失败的请求会静默堆积。
 * 这里统一登记，卸载时全部清掉。
 */
let timers: number[] = []
function later(fn: () => void, ms: number) {
  timers.push(window.setTimeout(fn, ms))
}

async function load() {
  loading.value = true
  errorMsg.value = ''
  try {
    const [ov, sch] = await Promise.all([
      api.get<OverviewResponse>('/api/overview'),
      api.get<{ schedules: Schedule[] }>('/api/schedules'),
    ])
    data.value = ov
    schedules.value = sch.schedules ?? []
  } catch (e) {
    errorMsg.value = e instanceof Error ? e.message : String(e)
  } finally {
    loading.value = false
  }
}

// 本页原先的空态「立即巡检」按钮已删：手动检测入口统一到顶栏「检查更新」，
// 检测完成后顶栏自增 app.checkTick，本页 watch 它重新拉取总览。

/**
 * 「容器状态」卡的轮询。
 *
 * CPU / 内存 / 磁盘原来只在 load() 时拿一次，卡片却挂着「实时」徽标 —— 用户看到的
 * 是进度条永远不动的快照。现在每 10 秒打一次轻量的 /api/usage（后端 CPU/内存恰好
 * 有 10 秒缓存），只合并这三个数，不整页刷新、不触发 loading。
 * 页面切到后台时跳过，别在看不见的地方白白打接口。
 */
let usageTimer: number | undefined
async function pollUsage() {
  if (document.hidden) return
  try {
    const u = await api.get<{
      cpuPercent: number
      memUsed: number
      memTotal: number
      diskFree: number
      diskTotal: number
    }>('/api/usage')
    if (!data.value) return
    data.value.usage = { cpuPercent: u.cpuPercent, memUsed: u.memUsed, memTotal: u.memTotal }
    if (u.diskTotal) data.value.disk = { free: u.diskFree, total: u.diskTotal }
  } catch {
    /* 静默：一张卡的数字，不值得为它弹错误条，下一轮再试 */
  }
}

/** 把 "68.4 GB" 拆成数值和单位，好让单位用小字号跟在后面（设计稿的 .v small）。 */
function sizeParts(n: number | undefined | null) {
  const [v = '0', u = 'B'] = formatBytes(n).split(' ')
  return { v, u }
}

const donut = computed(() => {
  const total = data.value?.containers?.total ?? 0
  const running = data.value?.containers?.running ?? 0
  const C = 2 * Math.PI * 46
  const runArc = total ? (running / total) * C : 0
  const stopArc = total ? ((total - running) / total) * C : 0
  return { C, runArc, stopArc }
})

const memTotal = computed(
  () => data.value?.usage?.memTotal || data.value?.docker?.memTotal || 0,
)
const memUsed = computed(() => data.value?.usage?.memUsed ?? 0)
const diskUsed = computed(() => {
  const d = data.value?.disk
  if (!d?.total) return 0
  return d.total - d.free
})
const pct = (used: number, total: number) =>
  total > 0 ? Math.min((used / total) * 100, 100) : 0

const cpuPercent = computed(() => Math.min(data.value?.usage?.cpuPercent ?? 0, 100))

/** 待更新的容器数（概览卡一行 = 一个容器，头部徽标与行数保持一致）。 */
const pendingCount = computed(() => data.value?.updates?.items?.length ?? 0)

/** 待更新镜像：按镜像归并，只给顶部统计卡用（「待更新镜像」数的是镜像不是容器）。 */
const updateGroups = computed(() => {
  const map = new Map<
    string,
    { image: string; containers: string[]; localDigest: string; remoteDigest: string }
  >()
  for (const it of data.value?.updates?.items ?? []) {
    const entry = map.get(it.image)
    if (entry) {
      entry.containers.push(it.container)
    } else {
      map.set(it.image, {
        image: it.image,
        containers: [it.container],
        localDigest: it.localDigest,
        remoteDigest: it.remoteDigest,
      })
    }
  }
  return [...map.values()]
})

/** 计划任务里最近的一次「下次执行」。 */
const nextRun = computed(() => {
  const times = schedules.value
    .filter((s) => s.enabled && s.nextRun)
    .map((s) => Date.parse(s.nextRun as string))
    .filter((t) => !Number.isNaN(t))
    .sort((a, b) => a - b)
  return times.length ? times[0] : 0
})

const recentRows = computed(() => (data.value?.recent ?? []).slice(0, 8))

/** 图标：取名字首字母。 */
function initial(name: string) {
  return (name[0] ?? '?').toUpperCase()
}

/** 待更新列表里的镜像名：只留最后一段仓库名、去掉 tag（设计稿写的就是 qbittorrent / jellyfin / nginx）。 */
function imageShort(image: string) {
  const noTag = (image.split('@')[0] ?? image).split('/').pop() ?? image
  const colon = noTag.indexOf(':')
  return colon > 0 ? noTag.slice(0, colon) : noTag
}

function kindClass(kind: string, status: string) {
  if (status === 'failed') return 'dh-badge-err'
  if (kind === 'schedule') return 'dh-badge-accent'
  return 'dh-badge-plain'
}

function resultBadge(l: RunLog) {
  if (l.status === 'success') return { text: '成功', cls: 'dh-badge-run' }
  if (l.status === 'up_to_date') return { text: '跳过', cls: 'dh-badge-accent' }
  if (l.status === 'failed') return { text: '失败', cls: 'dh-badge-err' }
  if (l.status === 'done') return { text: '完成', cls: 'dh-badge-run' }
  return { text: l.status || '—', cls: 'dh-badge-plain' }
}

/** 时间列：今天只出 HH:mm，昨天带前缀，更早出日期。 */
function logTime(ts: string) {
  const t = Date.parse(ts)
  if (Number.isNaN(t)) return '—'
  const d = new Date(t)
  const p = (n: number) => String(n).padStart(2, '0')
  const hm = `${p(d.getHours())}:${p(d.getMinutes())}`
  const today0 = new Date()
  today0.setHours(0, 0, 0, 0)
  const day0 = new Date(t)
  day0.setHours(0, 0, 0, 0)
  const diff = Math.round((day0.getTime() - today0.getTime()) / 86400000)
  if (diff === 0) return hm
  if (diff === -1) return `昨天 ${hm}`
  return formatDayTime(ts, 'date')
}

// 顶栏「检查更新」跑完一轮后重新拉总览（含待更新清单）——本页不再有自己的检测按钮。
watch(
  () => app.checkTick,
  () => void load(),
)

onMounted(() => {
  void load()
  // 容器状态卡：每 10 秒刷新一次 CPU / 内存 / 磁盘（与后端缓存的 10 秒 TTL 对齐）
  usageTimer = window.setInterval(() => void pollUsage(), 10000)
  // 订阅事件只为「批次跑完 / 巡检结束」时自动刷新总览，界面本身不再展示活动流。
  closeStream = openStream('/api/events/stream', (topic, ev) => {
    if (topic !== 'update' && topic !== 'schedule') return
    if (ev.kind === 'batch_done' || ev.kind === 'check_done') {
      later(() => void load(), 800)
    }
  })
})

onUnmounted(() => {
  closeStream?.()
  if (usageTimer !== undefined) window.clearInterval(usageTimer)
  timers.forEach((t) => window.clearTimeout(t))
  timers = []
})
</script>

<template>
  <div class="flex flex-col gap-3.5 p-[18px]">
    <div v-if="errorMsg" class="dh-banner dh-banner-err">
      <AlertTriangle class="h-4 w-4 flex-none" />
      <span class="min-w-0 flex-1">{{ errorMsg }}</span>
      <button class="dh-btn dh-btn-sm" @click="load">重试</button>
    </div>

    <div v-if="data?.dockerError" class="dh-banner dh-banner-err">
      <AlertTriangle class="h-4 w-4 flex-none" />
      <span class="min-w-0 flex-1">无法连接 Docker 守护进程：{{ data.dockerError }}</span>
    </div>

    <!-- 四个指标 -->
    <div class="grid grid-cols-2 gap-3 lg:grid-cols-4">
      <div class="dh-metric">
        <div class="k">容器</div>
        <div class="v">
          {{ data?.containers?.total ?? '—' }} <small>个</small>
        </div>
        <div class="s">
          运行 {{ data?.containers?.running ?? 0 }} · 停止 {{ data?.containers?.stopped ?? 0 }}
          <template v-if="data?.containers?.unhealthy"> · 异常 {{ data.containers.unhealthy }}</template>
        </div>
      </div>

      <div class="dh-metric">
        <div class="k">待更新镜像</div>
        <div class="v" :class="updateGroups.length ? 'text-warn-text' : ''">
          {{ updateGroups.length }} <small>个</small>
        </div>
        <div class="s">
          {{ data?.updates?.checkedAt ? `影响 ${data.updates.items?.length ?? 0} 个容器` : '还没做过巡检' }}
        </div>
      </div>

      <div class="dh-metric">
        <div class="k">计划任务</div>
        <div class="v">
          {{ schedules.length }} <small>个</small>
        </div>
        <div class="s">下次 {{ nextRun ? formatDayTime(nextRun) : '暂无启用中的任务' }}</div>
      </div>

      <div class="dh-metric">
        <div class="k">镜像占用</div>
        <div class="v">
          {{ sizeParts(data?.images?.sizeBytes).v }} <small>{{ sizeParts(data?.images?.sizeBytes).u }}</small>
        </div>
        <div class="s">可回收 {{ formatBytes(data?.images?.reclaimable) }}</div>
      </div>
    </div>

    <!-- 容器状态 + 待更新镜像 -->
    <div class="grid grid-cols-1 gap-3.5 xl:grid-cols-2">
      <div class="dh-card">
        <div class="dh-card-head">
          <span>容器状态</span>
          <span class="ml-auto dh-badge dh-badge-plain">实时</span>
        </div>
        <div class="dh-card-body flex items-center gap-[22px]">
          <!-- SVG 的呈现属性（stroke="..."）里不能写 var()，变量只在 style 属性/样式表里生效，
               所以下面这些颜色统一走 style，跟随主题切换。 -->
          <svg viewBox="0 0 120 120" class="h-[118px] w-[118px] flex-none">
            <circle
              cx="60"
              cy="60"
              r="46"
              fill="none"
              style="stroke: var(--color-line-1)"
              stroke-width="13"
            />
            <circle
              cx="60"
              cy="60"
              r="46"
              fill="none"
              style="stroke: var(--color-run)"
              stroke-width="13"
              stroke-linecap="round"
              :stroke-dasharray="`${donut.runArc} ${donut.C}`"
              transform="rotate(-90 60 60)"
            />
            <circle
              v-if="donut.stopArc > 0"
              cx="60"
              cy="60"
              r="46"
              fill="none"
              style="stroke: var(--color-stop)"
              stroke-width="13"
              :stroke-dasharray="`${donut.stopArc} ${donut.C}`"
              :stroke-dashoffset="-donut.runArc"
              transform="rotate(-90 60 60)"
            />
            <text
              x="60"
              y="56"
              text-anchor="middle"
              style="fill: var(--color-text-1)"
              font-size="21"
              font-weight="600"
            >
              {{ data?.containers?.running ?? 0 }}
            </text>
            <text x="60" y="74" text-anchor="middle" style="fill: var(--color-text-4)" font-size="11">
              运行中
            </text>
          </svg>

          <div class="flex min-w-0 flex-1 flex-col gap-[11px]">
            <div class="flex gap-4 text-[12px] text-text-4">
              <span class="flex items-center gap-1.5">
                <i class="h-[7px] w-[7px] rounded-full bg-run" />运行中 {{ data?.containers?.running ?? 0 }}
              </span>
              <span class="flex items-center gap-1.5">
                <i class="h-[7px] w-[7px] rounded-full bg-stop" />已停止 {{ data?.containers?.stopped ?? 0 }}
              </span>
            </div>

            <div class="flex flex-col gap-1.5">
              <div class="flex justify-between text-[11.5px] text-text-4">
                <span>CPU 总占用</span><span>{{ cpuPercent.toFixed(0) }}%</span>
              </div>
              <div class="dh-bar"><i :style="{ width: `${cpuPercent}%` }" /></div>
            </div>

            <div class="flex flex-col gap-1.5">
              <div class="flex justify-between text-[11.5px] text-text-4">
                <span>内存 {{ formatBytes(memUsed) }} / {{ formatBytes(memTotal) }}</span>
                <span>{{ pct(memUsed, memTotal).toFixed(0) }}%</span>
              </div>
              <div class="dh-bar">
                <i :style="{ width: `${pct(memUsed, memTotal)}%`, background: 'var(--color-accent-text)' }" />
              </div>
            </div>

            <!-- 磁盘：拿不到宿主机文件系统信息时（例如后端跑在 Windows 上）整行隐藏 -->
            <div v-if="data?.disk?.total" class="flex flex-col gap-1.5">
              <div class="flex justify-between text-[11.5px] text-text-4">
                <span>磁盘 {{ formatBytes(diskUsed) }} / {{ formatBytes(data?.disk?.total) }}</span>
                <span>{{ pct(diskUsed, data?.disk?.total ?? 0).toFixed(0) }}%</span>
              </div>
              <div class="dh-bar">
                <i :style="{ width: `${pct(diskUsed, data?.disk?.total ?? 0)}%`, background: 'var(--color-chart-indigo)' }" />
              </div>
            </div>
          </div>
        </div>
      </div>

      <div class="dh-card flex flex-col">
        <div class="dh-card-head">
          <span>待更新镜像</span>
          <span class="ml-auto dh-badge" :class="pendingCount ? 'dh-badge-warn' : 'dh-badge-plain'">
            {{ pendingCount }} 个
          </span>
        </div>

        <div v-if="!pendingCount" class="dh-card-body flex flex-1 flex-col items-center justify-center gap-2 py-6 text-center">
          <div class="text-[12.5px] text-text-3">
            {{ data?.updates?.checkedAt ? '所有镜像都是最新的' : '还没有做过巡检' }}
          </div>
          <div class="text-[11.5px] leading-relaxed text-text-5">
            {{
              data?.updates?.checkedAt
                ? '没有任何容器需要更新。'
                : '点顶栏右上角的「检查更新」可以只读地检查一遍所有容器。'
            }}
          </div>
        </div>
        <div v-else class="dh-card-body flex flex-col gap-[11px]">
          <!-- 一行 = 一个待更新的容器。以前按镜像简称分组：ghcr.io/music-assistant/server
               只剩 "server"、redis:alpine 剩 "redis"，用户根本对不上是哪个容器。
               摘要哈希（3b39059 → 130a28b）与「N 个」徽标对用户没有信息量，一并去掉。 -->
          <div
            v-for="it in (data?.updates?.items ?? []).slice(0, 4)"
            :key="it.container"
            class="flex items-center gap-2.5"
          >
            <div
              class="grid h-[28px] w-[28px] flex-none place-items-center rounded-[9px] bg-line-2 text-[11px] font-semibold text-accent-text"
            >
              {{ initial(it.container) }}
            </div>
            <div class="min-w-0 flex-1">
              <div class="truncate text-[12.5px] font-semibold" :title="it.container">{{ it.container }}</div>
              <div class="truncate text-[11.5px] text-text-5" :title="it.image">{{ imageShort(it.image) }}</div>
            </div>
          </div>

          <button class="dh-btn dh-btn-primary mt-0.5" @click="router.push('/containers')">
            前往容器页面
            <ArrowRight class="h-3.5 w-3.5" />
          </button>
        </div>
      </div>
    </div>

    <!-- 最近执行记录 -->
    <div class="dh-card">
      <div class="dh-card-head">
        <span>最近执行记录</span>
        <span class="ml-auto text-[12px] font-normal text-text-4">最近 15 条</span>
      </div>
      <div v-if="!recentRows.length" class="dh-card-body text-[12.5px] text-text-4">暂无记录</div>
      <div v-else class="dh-card-body flex flex-col py-1">
        <div v-for="l in recentRows" :key="l.id" class="dh-tl">
          <div class="w-[86px] flex-none text-[11.5px] text-text-5">{{ logTime(l.ts) }}</div>
          <span class="dh-badge flex-none" :class="kindClass(l.kind, l.status)">{{ runKindLabel(l.kind) }}</span>
          <span class="min-w-0 flex-1 truncate text-[12.5px] text-text-2" :title="l.message">
            {{ l.message }}
          </span>
          <span class="dh-badge flex-none" :class="resultBadge(l).cls">{{ resultBadge(l).text }}</span>
        </div>
      </div>
    </div>
  </div>
</template>
