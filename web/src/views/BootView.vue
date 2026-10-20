<script setup lang="ts">
/**
 * 顺序启动（/boot）。
 *
 * 一句话：本页负责将启动策略非 unless-stopped 且非 always 的容器按顺序启动。
 * 页面不做任何策略修改 —— 哪些容器归本页管，完全由 Docker 里那个容器的
 * restartPolicy 实时决定（no / on-failure 归本页，always / unless-stopped 归 Docker）。
 */
import { computed, onMounted, onUnmounted, ref } from 'vue'
import { useRouter } from 'vue-router'
import { GripVertical, History, Loader2, Play } from 'lucide-vue-next'
import { api } from '@/api/client'
import type { BootContainer, BootLastRun, BootPlan, BootResponse, BootStatus } from '@/api/types'
import { containerStateLabel, formatDayTime } from '@/utils/format'
import { useToastStore } from '@/stores/toast'
import Modal from '@/components/Modal.vue'
import UiSwitch from '@/components/UiSwitch.vue'

const toast = useToastStore()
const router = useRouter()

const loading = ref(true)
const saving = ref(false)
const confirming = ref(false)
const plan = ref<BootPlan>({
  enabled: true,
  initialDelaySec: 30,
  waitTimeoutSec: 60,
  defaultDelaySec: 3,
  items: [],
})
const rows = ref<BootContainer[]>([])
const self = ref({ name: '', restartPolicy: '' })
const lastRun = ref<BootLastRun | null>(null)
const status = ref<BootStatus | null>(null)

/** 只有「归本页管」的才可编辑；另一组只读。 */
const managedRows = computed(() => rows.value.filter((r) => r.managed))
const dockerRows = computed(() => rows.value.filter((r) => !r.managed))

/** 本次会被真正启动的（开关开着 + 当前没在跑）—— 二次确认里报这个数。 */
const pendingStart = computed(() => managedRows.value.filter((r) => r.auto && r.state !== 'running'))

const running = computed(() => !!status.value?.running)

/**
 * 有 total 才算「有进度」。
 * 等 Docker 就绪、初始等待这两段还没数过容器，total 是 0 —— 这时候进度条只能是
 * 一条 0% 的空槽，说明文字放页头那一格里更省地方（用户 2026-10-20 提的）。
 */
const hasProgress = computed(() => !!status.value?.total)

const fmtSecs = (s: number) => (s >= 60 ? `${Math.floor(s / 60)} 分 ${s % 60} 秒` : `${s} 秒`)

/**
 * 「照当前设置从头跑一遍」一共要多久。
 * 口径固定：初始等待 + 每个容器约 8 秒（启动并等到 running）+ 各自延时，只数开着
 * 「开机启动」的行（关掉的不归 dockhelm 拉）。
 * 不按运行状态过滤 —— 按实时状态算出来的数会随执行进度一直缩水，而且全都在跑时
 * 会退化成「预计 全部已在运行」这种说不通的句子。已在运行的跳过属于"实际更短"，
 * 放进 title 里解释。
 */
const estimate = computed(() => {
  const rows = managedRows.value.filter((r) => r.auto)
  if (!rows.length) return null
  const init = plan.value.initialDelaySec
  const startup = rows.length * 8
  const delays = rows.reduce((n, r) => n + r.delaySec, 0)
  const secs = init + startup + delays
  return {
    text: fmtSecs(secs),
    title:
      `按当前设置估算：初始等待 ${init} 秒 + ${rows.length} 个容器 × 约 8 秒（启动并等到 running）` +
      ` + 各容器延时合计 ${delays} 秒 = ${fmtSecs(secs)}。已在运行的容器会被跳过，实际用时更短。`,
  }
})

const lastRunText = computed(() => {
  const l = lastRun.value
  if (!l) return '上次执行：从未执行'
  const when = formatDayTime(l.ts)
  if (!l.total) return `上次执行：${when} · ${l.message}`
  const parts = [`成功 ${l.success ?? 0}`, `跳过 ${l.skipped ?? 0}`]
  if (l.failed) parts.push(`失败 ${l.failed}`)
  return `上次执行：${when} · ${parts.join('、')}`
})

const lastFailed = computed(() => lastRun.value?.status === 'failed' || (lastRun.value?.failed ?? 0) > 0)

const progressText = computed(() => {
  const s = status.value
  if (!s?.running) return ''
  if (s.phase === 'waiting-docker') return '等待 Docker 就绪…'
  if (s.phase === 'initial-delay') return `初始等待 ${plan.value.initialDelaySec} 秒后开始…`
  if (!s.total) return '正在读取容器…'
  const tail = s.phase === 'starting' && s.current ? `正在启动 ${s.current}` : '准备中'
  const counts = [`成功 ${s.success}`, `跳过 ${s.skipped}`]
  if (s.failed) counts.push(`失败 ${s.failed}`)
  return `第 ${s.index}/${s.total} · ${tail} · ${counts.join('、')}`
})

const progressPercent = computed(() => {
  const s = status.value
  if (!s?.total) return 0
  return Math.min(100, Math.round((s.index / s.total) * 100))
})

/**
 * Dockhelm 自己必须是 unless-stopped —— 它要是也归「开机不自启」，
 * 面板起不来，这一页就永远不会执行（方案 §七 风险 1）。
 */
const selfWarn = computed(() => {
  const p = self.value.restartPolicy
  if (!self.value.name || !p) return ''
  if (p !== 'no' && p !== 'on-failure') return ''
  return `Dockhelm 自身容器（${self.value.name}）的重启策略是 ${p}：开机时它不会自己起来，面板起不来，这一页也就不会执行。建议在 NAS 面板把它改回 unless-stopped。`
})

async function load() {
  loading.value = true
  try {
    const d = await api.get<BootResponse>('/api/boot')
    plan.value = d.plan
    rows.value = d.containers ?? []
    self.value = d.self ?? { name: '', restartPolicy: '' }
    lastRun.value = d.lastRun ?? null
    status.value = d.status
    if (d.status?.running) startPolling()
  } catch (e) {
    toast.error('读取顺序启动设置失败', e instanceof Error ? e.message : String(e))
  } finally {
    loading.value = false
  }
}

async function save() {
  saving.value = true
  try {
    const r = await api.put<{ plan: BootPlan }>('/api/boot', plan.value)
    plan.value = r.plan
    applyPlanOrder()
  } catch (e) {
    toast.error('保存失败', e instanceof Error ? e.message : String(e))
    await load()
  } finally {
    saving.value = false
  }
}

/**
 * 按清单顺序重排本地行。
 *
 * 列表渲染的是 server 返回的 rows（它们带着容器的实时状态），不是 plan.items ——
 * 所以只改 plan.items 再保存，界面上还是旧顺序，看起来就像「拖拽无效」
 * （刷新一下又是新顺序，因为服务端确实存对了）。保存成功后按 plan 重排一次。
 */
function applyPlanOrder() {
  const rank = new Map(plan.value.items.map((it, i) => [it.name, i]))
  const pos = (r: BootContainer) => rank.get(r.name) ?? Number.MAX_SAFE_INTEGER
  // Array.prototype.sort 是稳定排序，清单外的行（归 Docker 管的）保持原有相对顺序
  rows.value = [...rows.value].sort((a, b) => pos(a) - pos(b))
}

// ---------- 行内编辑 ----------

function setAuto(row: BootContainer, v: boolean) {
  const it = plan.value.items.find((x) => x.name === row.name)
  if (!it) return
  it.auto = v
  row.auto = v
  void save()
}

function setDelay(row: BootContainer, ev: Event) {
  const el = ev.target as HTMLInputElement
  const raw = Number(el.value)
  const v = Number.isFinite(raw) ? Math.min(300, Math.max(0, Math.round(raw))) : plan.value.defaultDelaySec
  el.value = String(v)
  const it = plan.value.items.find((x) => x.name === row.name)
  if (!it || it.delaySec === v) return
  it.delaySec = v
  row.delaySec = v
  void save()
}

function applyAllDelay() {
  const v = plan.value.defaultDelaySec
  const n = plan.value.items.length
  if (!n) return
  for (const it of plan.value.items) it.delaySec = v
  for (const r of rows.value) r.delaySec = v
  void save()
  // 一次性改动十几行，给个明确回执，否则看不出到底改了几个
  toast.success('已应用到全部行', `${n} 个容器的延时都设为 ${v} 秒`)
}

function setParam(key: 'initialDelaySec' | 'waitTimeoutSec' | 'defaultDelaySec', ev: Event) {
  const el = ev.target as HTMLInputElement
  const raw = Number(el.value)
  if (!Number.isFinite(raw)) return
  plan.value[key] = Math.max(0, Math.round(raw))
  void save()
}

// ---------- 拖拽排序（Pointer Events，鼠标与触屏共用一套） ----------
//
// 为什么不用原生 HTML5 拖拽：`dragstart/dragover/drop` 在触屏上根本不触发
// （iOS Safari 至今不支持，Android 上也不可靠），手机上那一列手柄按下去毫无反应。
// Pointer Events 是鼠标 / 触屏 / 触控笔统一的一套事件，`setPointerCapture` 之后
// 指针移出元素也照样收得到 move，不会中途丢失。
//
// 交互约定：
//   - 鼠标：按住整行任意空白处拖（点输入框、开关、链接除外）；
//   - 触屏：只能按住左侧手柄拖 —— 行内还有输入框和开关，整行可拖会抢掉滚动和点按。
//     手柄上 `touch-action: none`，浏览器不会把它当成滚动手势，也就不会发 pointercancel。
//   - 移动不足 4px 视为点击，不触发重排；
//   - 拖动过程中不移动 DOM（只画一条落点线），松手才重排并保存 —— 拖动中改 DOM
//     会让缓存的行位置失效，50 行容器里很容易错位。

const listRef = ref<HTMLElement | null>(null)
const dragFrom = ref(-1)
const dropSlot = ref(-1) // 落点：插到第几行之前（= 行数 表示插到末尾）
const dragActive = ref(false) // 已越过阈值，进入真正的拖拽

const DRAG_THRESHOLD = 4 // px，小于这个距离算点击
const EDGE_ZONE = 64 // 距滚动容器上下边缘多少像素开始自动滚
const EDGE_MAX_SPEED = 18 // 自动滚的最大速度（px / 帧）

let dragPointerId = -1
let dragStartY = 0
let lastPointerY = 0
let dragRects: { top: number; bottom: number }[] = []
let scrollEl: HTMLElement | null = null
let scrollSpeed = 0
let rafId = 0
let captureEl: HTMLElement | null = null

function onRowPointerDown(i: number, e: PointerEvent) {
  const target = e.target as HTMLElement
  if (e.pointerType === 'mouse') {
    if (e.button !== 0) return
    // 行里的输入框 / 开关 / 链接得照常工作，不能一按就被拖走
    if (target.closest('input, button, a, label, select, textarea')) return
  } else if (!target.closest('[data-grip]')) {
    // 触屏只有手柄可拖
    return
  }
  if (managedRows.value.length < 2) return

  dragFrom.value = i
  dropSlot.value = i
  dragPointerId = e.pointerId
  dragStartY = e.clientY
  lastPointerY = e.clientY
  dragActive.value = false
  // 拖动期间布局不变（只画落点线），量一次就够了；每次 move 都量 50 行太浪费
  dragRects = measureRows()
  scrollEl = findScrollParent(e.currentTarget as HTMLElement | null)
  scrollSpeed = 0

  window.addEventListener('pointermove', onDragMove, { passive: false })
  window.addEventListener('pointerup', onDragEnd)
  window.addEventListener('pointercancel', onDragCancel)
  rafId = requestAnimationFrame(autoScrollTick)
  // 捕获指针：指针划出这一行、甚至划出窗口，move / up 仍然发给它，
  // 不会出现「拖到一半事件没了、行一直半透明」。
  captureEl = (e.currentTarget as HTMLElement) ?? null
  try {
    captureEl?.setPointerCapture(e.pointerId)
  } catch {
    // 某些环境（合成事件、已失活的指针）不支持，退回纯 window 监听
  }
  document.body.classList.add('dh-drag-lock')
  e.preventDefault()
}

function measureRows(): { top: number; bottom: number }[] {
  const host = listRef.value
  if (!host) return []
  return Array.from(host.querySelectorAll<HTMLElement>('[data-boot-name]')).map((el) => {
    const r = el.getBoundingClientRect()
    return { top: r.top, bottom: r.bottom }
  })
}

/** 从行往上找真正在滚的那个祖先（本页是 AppShell 的 <main>）。 */
function findScrollParent(from: HTMLElement | null): HTMLElement | null {
  let p = from?.parentElement ?? null
  while (p) {
    const st = getComputedStyle(p)
    if (/(auto|scroll)/.test(st.overflowY) && p.scrollHeight > p.clientHeight + 1) return p
    p = p.parentElement
  }
  return null
}

/** 指针纵向位置 → 插入槽位（0..n）。落在上半部算插到这一行之前，下半部算之后。 */
function slotAt(y: number): number {
  const n = dragRects.length
  for (let i = 0; i < n; i++) {
    const r = dragRects[i]!
    if (y < r.bottom) return y < (r.top + r.bottom) / 2 ? i : i + 1
  }
  return n
}

function onDragMove(e: PointerEvent) {
  if (e.pointerId !== dragPointerId) return
  lastPointerY = e.clientY
  if (!dragActive.value) {
    if (Math.abs(e.clientY - dragStartY) < DRAG_THRESHOLD) return
    dragActive.value = true
  }
  e.preventDefault()
  dropSlot.value = slotAt(e.clientY)
  updateScrollSpeed(e.clientY)
}

/**
 * 自动滚动：拖到容器上下边缘附近时按距离加速滚。
 *
 * 必须用 rAF 循环而不是「每次 pointermove 滚一点」—— 指针停在边缘不动时
 * 是不产生 move 事件的，那种写法会滚一下就卡住。
 */
function updateScrollSpeed(y: number) {
  if (!scrollEl) {
    scrollSpeed = 0
    return
  }
  const rect = scrollEl.getBoundingClientRect()
  const top = rect.top + 72 // 让开吸顶工具条，否则落点线藏在它后面
  const bottom = rect.bottom - 20
  if (y < top + EDGE_ZONE) {
    scrollSpeed = -Math.ceil(((top + EDGE_ZONE - y) / EDGE_ZONE) * EDGE_MAX_SPEED)
  } else if (y > bottom - EDGE_ZONE) {
    scrollSpeed = Math.ceil(((y - (bottom - EDGE_ZONE)) / EDGE_ZONE) * EDGE_MAX_SPEED)
  } else {
    scrollSpeed = 0
  }
  // 越靠近边缘越快，但别让它一帧跳出去半屏
  scrollSpeed = Math.max(-EDGE_MAX_SPEED, Math.min(EDGE_MAX_SPEED, scrollSpeed))
}

function autoScrollTick() {
  if (dragFrom.value < 0) return
  if (scrollEl && scrollSpeed !== 0) {
    const before = scrollEl.scrollTop
    scrollEl.scrollTop = before + scrollSpeed
    const moved = scrollEl.scrollTop - before
    if (moved !== 0) {
      // 内容滚了多少，缓存的行位置就整体平移多少 —— 比重新量一遍便宜且不闪
      for (const r of dragRects) {
        r.top -= moved
        r.bottom -= moved
      }
      dropSlot.value = slotAt(lastPointerY)
    }
  }
  rafId = requestAnimationFrame(autoScrollTick)
}

function onDragEnd(e: PointerEvent) {
  if (e.pointerId !== dragPointerId) return
  const from = dragFrom.value
  const slot = dropSlot.value
  const active = dragActive.value
  cleanupDrag()
  if (active && from >= 0 && slot >= 0) commitReorder(from, slot)
}

function onDragCancel(e: PointerEvent) {
  if (e.pointerId !== dragPointerId) return
  cleanupDrag()
}

function cleanupDrag() {
  window.removeEventListener('pointermove', onDragMove)
  window.removeEventListener('pointerup', onDragEnd)
  window.removeEventListener('pointercancel', onDragCancel)
  if (rafId) cancelAnimationFrame(rafId)
  rafId = 0
  if (captureEl && dragPointerId >= 0) {
    try {
      if (captureEl.hasPointerCapture(dragPointerId)) captureEl.releasePointerCapture(dragPointerId)
    } catch {
      // 元素已被卸载，捕获随之失效，不用管
    }
  }
  captureEl = null
  document.body.classList.remove('dh-drag-lock')
  dragFrom.value = -1
  dropSlot.value = -1
  dragActive.value = false
  dragRects = []
  scrollEl = null
  scrollSpeed = 0
  dragPointerId = -1
}

/**
 * 把「从 from 拖到 slot」落实到清单上。
 *
 * ⚠ plan.items 里是**全部**容器（含 always / unless-stopped 的），页面上被过滤成两组。
 * 所以不能整体重排 items —— 那会把只读组的位置也搅乱。做法是：只在受管行占据的
 * 那些槽位上，按新顺序重新填一遍；不受管的条目原地不动。
 */
function commitReorder(from: number, slot: number) {
  const names = managedRows.value.map((r) => r.name)
  if (slot === from || slot === from + 1) return // 原地放下，什么都没变
  const moved = names.splice(from, 1)[0]
  if (!moved) return
  names.splice(slot > from ? slot - 1 : slot, 0, moved)

  const byName = new Map(plan.value.items.map((it) => [it.name, it]))
  const managedSet = new Set(names)
  const queue: BootPlan['items'] = []
  for (const n of names) {
    const it = byName.get(n)
    if (it) queue.push(it)
  }
  let k = 0
  plan.value.items = plan.value.items.map((it) => {
    if (!managedSet.has(it.name)) return it
    return queue[k++] ?? it
  })
  // 先本地重排，界面立刻跟手；再异步保存（等接口回来才动的话会有明显延迟感）
  applyPlanOrder()
  void save()
}

/** 落点线：插在某行之前画在它顶部，插在末尾画在最后一行底部。 */
function dropClass(i: number) {
  if (!dragActive.value || dropSlot.value < 0) return ''
  if (dropSlot.value === i) return 'dh-boot-drop-top'
  if (dropSlot.value === dragRects.length && i === dragRects.length - 1) return 'dh-boot-drop-bottom'
  return ''
}

// ---------- 执行与轮询 ----------

let timer: number | undefined

function startPolling() {
  if (timer) return
  timer = window.setInterval(() => void refreshStatus(), 2000)
}

function stopPolling() {
  if (!timer) return
  window.clearInterval(timer)
  timer = undefined
}

async function refreshStatus() {
  try {
    const d = await api.get<{ status: BootStatus }>('/api/boot/status')
    const wasRunning = running.value
    status.value = d.status
    if (wasRunning && !d.status.running) {
      stopPolling()
      await load()
      toast.success('顺序启动已结束', '结果已写入运行记录')
    }
  } catch {
    stopPolling()
  }
}

async function doRun() {
  confirming.value = false
  try {
    await api.post('/api/boot/run', { skipDelay: true })
    startPolling()
    await refreshStatus()
  } catch (e) {
    toast.error('无法开始', e instanceof Error ? e.message : String(e))
  }
}

/**
 * 跳到设置页的「运行记录」并定位过去。
 *
 * 不用 hash + router.scrollBehavior：页面滚动发生在 AppShell 的 <main> 里，
 * 而 vue-router 的 el 定位是相对 window 的，滚不动那个容器。
 */
function gotoRunLogs() {
  void router.push('/settings').then(() => {
    window.setTimeout(() => {
      document.getElementById('runlogs')?.scrollIntoView({ behavior: 'smooth', block: 'start' })
    }, 120)
  })
}

function openRunLogs() {
  if (!lastFailed.value) return
  gotoRunLogs()
}

onMounted(() => void load())
onUnmounted(() => {
  stopPolling()
  cleanupDrag()
})
</script>

<template>
  <div class="flex flex-col">
    <!-- 吸顶工具条：标题 + 立即执行 + 上次执行 + 进度 -->
    <div class="sticky top-0 z-20 border-b border-line-1 bg-ink-950 px-[18px] pb-2.5 pt-[18px]">
      <div class="flex flex-wrap items-center gap-x-2.5 gap-y-2">
        <div class="dh-h1">顺序启动</div>
        <span v-if="saving" class="dh-sub">保存中…</span>

        <!-- 手机上按钮留在标题行右侧、摘要整行换行，避免工具条吃掉 1/6 屏高 -->
        <button class="dh-btn order-1 ml-auto" :disabled="running || loading" @click="confirming = true">
          <Loader2 v-if="running" class="h-3.5 w-3.5 dh-spin" />
          <Play v-else class="h-3.5 w-3.5" />
          {{ running ? '执行中…' : '立即执行一次' }}
        </button>

        <div class="order-2 flex min-w-0 flex-wrap items-center gap-x-3 gap-y-1 max-md:w-full md:order-3">
          <!-- 还没到「第 x/y 个」之前的阶段（等 Docker / 初始等待）没有百分比可画，
               就把这句话放在这一格里，不再单独占一行空进度条 -->
          <span v-if="running && !hasProgress" class="flex items-center gap-1.5 text-[11.5px] text-text-5">
            <Loader2 class="h-3 w-3 dh-spin" />{{ progressText }}
          </span>
          <span
            v-else
            class="text-[11.5px]"
            :class="lastFailed ? 'cursor-pointer text-err-text hover:underline' : 'text-text-5'"
            :title="lastRun ? lastRun.message : '还没有执行过'"
            @click="openRunLogs"
          >
            {{ lastRunText }}
          </span>

          <button
            class="dh-btn dh-btn-ghost dh-btn-sm !gap-1.5 text-[11.5px] max-md:ml-auto max-md:!px-2"
            title="打开设置页的运行记录"
            @click="gotoRunLogs"
          >
            <History class="h-3.5 w-3.5" /><span class="max-md:hidden">运行记录</span>
          </button>
        </div>
      </div>

      <div v-if="running && hasProgress" class="mt-2.5 flex items-center gap-2.5">
        <div class="h-[3px] flex-1 overflow-hidden rounded-full bg-line-1">
          <div class="h-full rounded-full bg-accent transition-all" :style="{ width: progressPercent + '%' }" />
        </div>
        <span class="whitespace-nowrap text-[11.5px] text-text-5">{{ progressText }}</span>
      </div>
    </div>

    <div class="flex flex-col gap-3.5 px-[18px] pb-[18px] pt-3">
      <div class="dh-sub">
        本页负责将启动策略非 unless-stopped 且非 always 的容器按顺序启动。
      </div>

      <div v-if="selfWarn" class="dh-banner-warn">
        {{ selfWarn }}
      </div>

      <!-- 执行参数 -->
      <div class="dh-card">
        <div class="dh-card-head items-baseline">
          <span class="shrink-0">执行参数</span>
          <span class="text-[11px] font-normal leading-[1.55] text-text-6">
            延时 = 容器启动成功后等 N 秒再启动下一个。新加入的容器默认取这个值；「应用到全部行」只改每行的延时，不动初始等待与容器启动超时。
          </span>
        </div>
        <div class="dh-card-body">
          <div class="flex flex-wrap items-center gap-x-7 gap-y-3">
            <label class="flex items-center gap-2 text-[12.5px] text-text-4">
              初始等待
              <input
                class="dh-input !w-[76px] text-center font-mono"
                type="number"
                min="0"
                max="1800"
                :value="plan.initialDelaySec"
                @change="setParam('initialDelaySec', $event)"
              />
              秒
            </label>
            <label
              class="flex items-center gap-2 text-[12.5px] text-text-4"
              title="发出 start 后最多等多久它进入 running；到点还没起来就记这个容器失败，接着起下一个"
            >
              容器启动超时
              <input
                class="dh-input !w-[76px] text-center font-mono"
                type="number"
                min="5"
                max="600"
                :value="plan.waitTimeoutSec"
                @change="setParam('waitTimeoutSec', $event)"
              />
              秒
            </label>

            <!-- 延时是「每行一个值 + 一个批量」成对的东西，用一个小框把它俩框在一起，
                 和上面两项并列时不会被读成「三项参数一起批量应用」 -->
            <div
              class="flex items-center gap-2.5 rounded-[10px] border border-line-4 bg-ink-650 py-1.5 pl-3 pr-2"
            >
              <label class="flex items-center gap-2 text-[12.5px] text-text-4">
                默认延时
                <input
                  class="dh-input !w-[76px] text-center font-mono"
                  type="number"
                  min="0"
                  max="300"
                  :value="plan.defaultDelaySec"
                  @change="setParam('defaultDelaySec', $event)"
                />
                秒
              </label>
              <button
                class="dh-btn dh-btn-sm text-[11.5px]"
                :disabled="!managedRows.length || saving"
                :title="`把全部 ${managedRows.length} 个容器的延时都设成 ${plan.defaultDelaySec} 秒`"
                @click="applyAllDelay"
              >
                应用到全部行
              </button>
            </div>
          </div>
        </div>
      </div>

      <!-- 由 dockhelm 按顺序启动 -->
      <div class="dh-card">
        <div class="dh-card-head items-baseline">
          <span class="shrink-0">由 dockhelm 按顺序启动 · {{ managedRows.length }} 个</span>
          <span v-if="managedRows.length" class="ml-1 text-[11px] font-normal leading-[1.55] text-text-6">
            每次 Dockhelm 启动都会执行一次，已在运行的跳过，只有开机启动已打开的未运行容器会被拉起。
          </span>
          <span
            v-if="estimate"
            class="ml-auto shrink-0 pl-3 text-[11.5px] font-normal text-text-5"
            :title="estimate.title"
          >
            预计全部启动完成需 {{ estimate.text }}
          </span>
        </div>

        <div v-if="!managedRows.length" class="dh-card-body text-[12.5px] text-text-4">
          {{ loading ? '正在载入…' : '还没有归本页管的容器。把容器的重启策略改成 on-failure（或 no），刷新后就会出现在这里。' }}
        </div>

        <div v-else class="overflow-x-auto">
          <div ref="listRef" class="min-w-0 px-3 pb-2 pt-2">
            <div class="dh-boot-grid dh-boot-head px-2 pb-1.5 text-[11px] text-text-6">
              <span />
              <span>#</span>
              <span>容器</span>
              <span>延时</span>
              <span class="whitespace-nowrap">开机启动</span>
              <span>状态</span>
              <span>策略</span>
            </div>

            <div
              v-for="(row, i) in managedRows"
              :key="row.name"
              :data-boot-name="row.name"
              class="dh-boot-row dh-boot-grid rounded-lg px-2 py-1.5 text-[12.5px]"
              :class="[
                dropClass(i),
                dragActive && dragFrom === i ? 'dh-boot-dragging' : '',
                row.state === 'running' ? 'dh-boot-run' : '',
              ]"
              @pointerdown="onRowPointerDown(i, $event)"
            >
              <span class="dh-c-grip" data-grip title="按住拖动手柄排序（触屏也可以）">
                <GripVertical class="h-3.5 w-3.5" />
              </span>
              <span class="dh-c-idx text-text-5">{{ i + 1 }}</span>
              <span class="dh-c-name truncate font-medium" :title="row.name">{{ row.name }}</span>
              <span class="dh-c-delay">
                <span class="dh-delay-label">延时</span>
                <!-- 手机上宽度收到 44px：高度保持 38px 不动，那是全站触屏
                     下的最小点按尺寸（style.css 的 min-height），压小反而难点 -->
                <input
                  class="dh-input !w-[52px] !px-1 text-center font-mono text-[11.5px] max-md:!w-[44px] max-md:!py-0"
                  type="number"
                  min="0"
                  max="300"
                  :value="row.delaySec"
                  :title="`${row.name} 启动成功后等多少秒再起下一个`"
                  @change="setDelay(row, $event)"
                />
                <span class="dh-delay-label">秒</span>
              </span>
              <div class="dh-c-switch flex justify-start max-md:justify-end">
                <UiSwitch
                  :model-value="row.auto"
                  :label="`${row.name} 开机启动`"
                  @update:model-value="(v: boolean) => setAuto(row, v)"
                />
              </div>
              <span class="dh-c-meta">
                <span :class="row.state === 'running' ? 'text-run-text' : 'text-text-5'">
                  {{ containerStateLabel(row.state) }}
                </span>
                <span class="font-mono text-[11px] text-text-4" :title="row.restartPolicy">
                  {{ row.restartPolicy || '—' }}
                </span>
              </span>
            </div>
          </div>
        </div>
      </div>

      <!-- 由 Docker 自行启动 -->
      <div class="dh-card">
        <div class="dh-card-head items-baseline">
          <span class="shrink-0">由 Docker 自行启动 · {{ dockerRows.length }} 个</span>
          <span v-if="dockerRows.length" class="ml-1 text-[11px] font-normal leading-[1.55] text-text-6">
            这些容器由 Docker 在开机时自动启动，不受本程序控制。想让它按顺序启动，去 NAS 面板把重启策略改成
            on-failure，刷新后会自动加入到本程序的顺序启动管理列表。
          </span>
        </div>

        <div v-if="!dockerRows.length" class="dh-card-body text-[12.5px] text-text-4">
          {{ loading ? '正在载入…' : '没有 always / unless-stopped 的容器 —— 全部受本页控制。' }}
        </div>

        <div v-else class="overflow-x-auto">
          <div class="min-w-0 px-3 pb-2 pt-2">
            <div class="dh-boot-grid dh-boot-head px-2 pb-1.5 text-[11px] text-text-6">
              <span />
              <span />
              <span>容器</span>
              <span />
              <span />
              <span>状态</span>
              <span>策略</span>
            </div>
            <div
              v-for="row in dockerRows"
              :key="row.name"
              class="dh-boot-grid px-2 py-1.5 text-[12.5px] text-text-5"
            >
              <span class="dh-c-grip" />
              <span class="dh-c-idx">—</span>
              <span class="dh-c-name truncate" :title="row.name">{{ row.name }}</span>
              <span class="dh-c-delay" />
              <span class="dh-c-filler" />
              <span class="dh-c-meta">
                <span :class="row.state === 'running' ? '' : 'text-text-6'">
                  {{ containerStateLabel(row.state) }}
                </span>
                <span class="font-mono text-[11px]">{{ row.restartPolicy || '—' }}</span>
              </span>
            </div>
          </div>
        </div>
      </div>
    </div>

    <Modal
      :open="confirming"
      title="立即执行一次顺序启动？"
      :subtitle="`将按当前顺序启动 ${pendingStart.length} 个已停止的容器，全程串行，随时可在页头看到进度。`"
      width="460px"
      @close="confirming = false"
    >
      <div class="text-[12.5px] leading-relaxed text-text-4">
        已在运行的、开关关掉的、以及归 Docker 管的容器都会被跳过。本次会跳过开头的初始等待，立刻开始。
      </div>
      <template #footer>
        <button class="dh-btn" @click="confirming = false">取消</button>
        <button class="dh-btn dh-btn-primary" @click="doRun">
          <Play class="h-3.5 w-3.5" />执行
        </button>
      </template>
    </Modal>
  </div>
</template>

<style scoped>
/* 运行中的行：左侧一条淡绿竖线（方案 §2.1）。
   用伪元素而不是 inset box-shadow —— 后者会顺着行的圆角弯成「(」形。 */
.dh-boot-row {
  position: relative;
}
.dh-boot-run::before {
  content: '';
  position: absolute;
  left: 0;
  top: 5px;
  bottom: 5px;
  width: 2px;
  border-radius: 2px;
  background: var(--color-run);
}
/* 拖拽的落点线：用内阴影画一条 2px 的边。
   不能用 border-t-2 —— 那会把这行撑高 2px，整张表在拖动过程中会上下跳，
   落点跟着偏，很多时候「拖了没反应」就是它造成的。 */
.dh-boot-drop-top {
  box-shadow: inset 0 2px 0 0 var(--color-accent);
}
.dh-boot-drop-bottom {
  box-shadow: inset 0 -2px 0 0 var(--color-accent);
}

/* 被拖动的那一行：压暗 + 底色，长列表里一眼能找回自己在拖哪一行 */
.dh-boot-dragging {
  background: var(--color-accent-row);
  opacity: 0.5;
}

/* 手柄：触屏上必须是拖拽的唯一入口，所以这里 touch-action: none 让浏览器
   别把它识别成滚动手势（否则会中途发 pointercancel，拖到一半就断）。
   桌面鼠标同样走 Pointer Events，手势光标靠 cursor 表达。 */
.dh-c-grip {
  display: flex;
  align-items: center;
  justify-content: center;
  cursor: grab;
  color: var(--color-text-6);
  touch-action: none;
  -webkit-user-select: none;
  user-select: none;
  -webkit-touch-callout: none;
}
.dh-boot-dragging .dh-c-grip {
  cursor: grabbing;
}

/* ------------------------------------------------------------------
   列表网格：桌面 7 列平铺；手机（<768px）折成两行，
   避免 660px 的表格在 375px 屏上只能看到左侧一半、
   延时输入框和开关被挤到屏幕外（实测要横向滚 300px 才够得着）。

   桌面列序：拖动手柄 / # / 容器 / 延时 / 开机启动 / 状态 / 策略
   —— 「延时」和「开机启动」紧跟在容器名后面，是这一行能改的两个值，
   眼睛不用横跨整个屏幕去对。容器名封顶 300px（大约 20 个汉字或 45 个
   拉丁字符，超过的截断 + title 悬浮看全名），剩下的宽度全给最右的
   「策略」：它是只读文本，被拉宽不影响"谁是谁"的对应关系。
   ------------------------------------------------------------------ */
.dh-boot-grid {
  display: grid;
  grid-template-columns: 18px 22px minmax(0, 300px) 54px 52px 60px minmax(92px, 1fr);
  align-items: center;
  column-gap: 8px;
}
/* 状态 + 策略桌面分属两列：用 contents 让它们脱离这一层，直接参与父网格 */
.dh-c-meta {
  display: contents;
}
/* 「延时」「秒」这两个字只在手机布局里出现（桌面靠表头） */
.dh-delay-label {
  display: none;
}

@media (max-width: 767px) {
  /* 手机上不表头：每行自带标签，列头反而占高度 */
  .dh-boot-head {
    display: none;
  }
  .dh-boot-grid {
    grid-template-columns: 24px 16px minmax(0, 1fr) auto;
    grid-template-areas:
      'grip idx name sw'
      'grip idx meta delay';
    column-gap: 7px;
    row-gap: 2px;
    padding-top: 6px;
    padding-bottom: 6px;
  }
  /* 手柄竖跨两行：24px 宽 × 整行高，手指按得住 */
  .dh-c-grip {
    grid-area: grip;
    align-self: stretch;
  }
  .dh-c-idx {
    grid-area: idx;
  }
  .dh-c-name {
    grid-area: name;
    min-width: 0;
  }
  .dh-c-meta {
    grid-area: meta;
    display: flex;
    gap: 8px;
    align-items: baseline;
    min-width: 0;
  }
  .dh-c-delay {
    grid-area: delay;
    display: flex;
    align-items: center;
    gap: 4px;
  }
  .dh-c-switch {
    grid-area: sw;
  }
  .dh-delay-label {
    display: inline;
    font-size: 11px;
    color: var(--color-text-6);
  }
  /* 只读组没有第 7 格，不让它生成隐式行 */
  .dh-c-filler {
    display: none;
  }
}

/* 拖拽期间锁住选中：否则鼠标划过大段文字会一路把行名、策略都选蓝，
   松手后选区还留在页面上。必须在全局作用域（body 在组件之外）。 */
:global(body.dh-drag-lock) {
  user-select: none;
  -webkit-user-select: none;
  cursor: grabbing;
}
</style>
