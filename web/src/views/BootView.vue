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
import { GripVertical, Loader2, Play } from 'lucide-vue-next'
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

const estimateText = computed(() => {
  const starts = pendingStart.value.length
  if (!starts) return managedRows.value.length ? '全部已在运行' : '—'
  // 8 秒是「启动并等它 running」的平均耗时，只用于这一行文案
  const secs = plan.value.initialDelaySec + pendingStart.value.reduce((n, r) => n + 8 + r.delaySec, 0)
  return secs >= 60 ? `${Math.floor(secs / 60)} 分 ${secs % 60} 秒` : `${secs} 秒`
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
  for (const it of plan.value.items) it.delaySec = v
  for (const r of rows.value) r.delaySec = v
  void save()
}

function setParam(key: 'initialDelaySec' | 'waitTimeoutSec' | 'defaultDelaySec', ev: Event) {
  const el = ev.target as HTMLInputElement
  const raw = Number(el.value)
  if (!Number.isFinite(raw)) return
  plan.value[key] = Math.max(0, Math.round(raw))
  void save()
}

// ---------- 拖拽排序（原生 HTML5，不引第三方库） ----------

const dragFrom = ref(-1)
const dragOver = ref(-1)

function onDragStart(i: number, e: DragEvent) {
  dragFrom.value = i
  if (e.dataTransfer) {
    e.dataTransfer.effectAllowed = 'move'
    e.dataTransfer.setData('text/plain', String(i))
  }
}

function onDragOver(i: number, e: DragEvent) {
  if (dragFrom.value < 0) return
  e.preventDefault()
  if (e.dataTransfer) e.dataTransfer.dropEffect = 'move'
  dragOver.value = i
}

function resetDrag() {
  dragFrom.value = -1
  dragOver.value = -1
}

function onDrop(i: number) {
  const from = dragFrom.value
  resetDrag()
  if (from < 0 || from === i) return
  const fromName = managedRows.value[from]?.name
  const toName = managedRows.value[i]?.name
  if (!fromName || !toName) return
  const arr = [...plan.value.items]
  const fi = arr.findIndex((x) => x.name === fromName)
  const ti = arr.findIndex((x) => x.name === toName)
  if (fi < 0 || ti < 0) return
  const [moved] = arr.splice(fi, 1)
  if (!moved) return
  arr.splice(ti, 0, moved)
  plan.value.items = arr
  void save()
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
onUnmounted(() => stopPolling())
</script>

<template>
  <div class="flex flex-col">
    <!-- 吸顶工具条：标题 + 立即执行 + 上次执行 + 进度 -->
    <div class="sticky top-0 z-20 border-b border-line-1 bg-ink-950 px-[18px] pb-2.5 pt-[18px]">
      <div class="flex flex-wrap items-center gap-2.5">
        <div class="dh-h1">顺序启动</div>
        <span v-if="saving" class="dh-sub">保存中…</span>

        <div class="ml-auto flex flex-wrap items-center gap-3">
          <button class="dh-btn" :disabled="running || loading" @click="confirming = true">
            <Loader2 v-if="running" class="h-3.5 w-3.5 dh-spin" />
            <Play v-else class="h-3.5 w-3.5" />
            {{ running ? '执行中…' : '立即执行一次' }}
          </button>
          <span
            class="text-[11.5px]"
            :class="lastFailed ? 'cursor-pointer text-err-text hover:underline' : 'text-text-5'"
            :title="lastRun ? lastRun.message : '还没有执行过'"
            @click="openRunLogs"
          >
            {{ lastRunText }}
          </span>
        </div>
      </div>

      <div v-if="running" class="mt-2.5 flex items-center gap-2.5">
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
        <div class="dh-card-head">
          <span>执行参数</span>
          <button class="ml-auto text-[11.5px] font-normal text-text-5 hover:text-accent" @click="applyAllDelay">
            批量应用到全部
          </button>
        </div>
        <div class="dh-card-body flex flex-wrap items-center gap-x-7 gap-y-3">
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
          <label class="flex items-center gap-2 text-[12.5px] text-text-4">
            单容器上限
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
          <label class="flex items-center gap-2 text-[12.5px] text-text-4">
            新行默认延时
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
          <span class="text-[11.5px] text-text-6">延时 = 这一行启动成功后等 N 秒再起下一个</span>
        </div>
      </div>

      <!-- 由 dockhelm 按顺序启动 -->
      <div class="dh-card">
        <div class="dh-card-head items-start">
          <span class="shrink-0 pt-[2px]">由 dockhelm 按顺序启动 · {{ managedRows.length }} 个</span>
          <span v-if="managedRows.length" class="ml-1 text-[11px] font-normal leading-[1.55] text-text-6">
            每次 Dockhelm 启动都会执行一次，已在运行的跳过，只有开机启动已打开的未运行容器会被拉起。
          </span>
          <span class="ml-auto shrink-0 pl-3 pt-[1px] text-[11.5px] font-normal text-text-5">
            预计 {{ estimateText }}
          </span>
        </div>

        <div v-if="!managedRows.length" class="dh-card-body text-[12.5px] text-text-4">
          {{ loading ? '正在载入…' : '还没有归本页管的容器。把容器的重启策略改成 on-failure（或 no），刷新后就会出现在这里。' }}
        </div>

        <div v-else class="overflow-x-auto">
          <div class="min-w-[660px] px-3 pb-2 pt-2">
            <div class="grid grid-cols-[18px_22px_minmax(0,1fr)_72px_92px_54px_56px] items-center gap-2 px-2 pb-1.5 text-[11px] text-text-6">
              <span />
              <span>#</span>
              <span>容器</span>
              <span>状态</span>
              <span>策略</span>
              <span>延时</span>
              <span class="text-right whitespace-nowrap">开机启动</span>
            </div>

            <div
              v-for="(row, i) in managedRows"
              :key="row.name"
              class="dh-boot-row grid grid-cols-[18px_22px_minmax(0,1fr)_72px_92px_54px_56px] items-center gap-2 rounded-lg px-2 py-1.5 text-[12.5px]"
              :class="[
                dragOver === i && dragFrom >= 0 && dragFrom !== i ? 'dh-boot-drop' : '',
                dragFrom === i ? 'opacity-45' : '',
                row.state === 'running' ? 'dh-boot-run' : '',
              ]"
              draggable="true"
              @dragstart="onDragStart(i, $event)"
              @dragover="onDragOver(i, $event)"
              @drop.prevent="onDrop(i)"
              @dragend="resetDrag"
            >
              <span class="cursor-grab text-text-6 active:cursor-grabbing" title="按住拖动排序">
                <GripVertical class="h-3.5 w-3.5" />
              </span>
              <span class="text-text-5">{{ i + 1 }}</span>
              <span class="truncate font-medium" :title="row.name">{{ row.name }}</span>
              <span :class="row.state === 'running' ? 'text-run-text' : 'text-text-5'">
                {{ containerStateLabel(row.state) }}
              </span>
              <span class="font-mono text-[11px] text-text-4" :title="row.restartPolicy">
                {{ row.restartPolicy || '—' }}
              </span>
              <input
                class="dh-input !w-[52px] !px-1 text-center font-mono text-[11.5px]"
                type="number"
                min="0"
                max="300"
                :value="row.delaySec"
                :title="`${row.name} 启动成功后等多少秒再起下一个`"
                @change="setDelay(row, $event)"
              />
              <div class="flex justify-end">
                <UiSwitch
                  :model-value="row.auto"
                  :label="`${row.name} 开机启动`"
                  @update:model-value="(v: boolean) => setAuto(row, v)"
                />
              </div>
            </div>
          </div>
        </div>
      </div>

      <!-- 由 Docker 自行启动 -->
      <div class="dh-card">
        <div class="dh-card-head items-start">
          <span class="shrink-0 pt-[2px]">由 Docker 自行启动 · {{ dockerRows.length }} 个</span>
          <span v-if="dockerRows.length" class="ml-1 text-[11px] font-normal leading-[1.55] text-text-6">
            这些容器由 Docker 在开机时自动启动，不受本程序控制。想让它按顺序启动，去 NAS 面板把重启策略改成
            on-failure，刷新后会自动加入到本程序的顺序启动管理列表。
          </span>
        </div>

        <div v-if="!dockerRows.length" class="dh-card-body text-[12.5px] text-text-4">
          {{ loading ? '正在载入…' : '没有 always / unless-stopped 的容器 —— 全部受本页控制。' }}
        </div>

        <div v-else class="overflow-x-auto">
          <div class="min-w-[660px] px-3 pb-2 pt-2">
            <div class="grid grid-cols-[18px_22px_minmax(0,1fr)_72px_92px_54px_56px] items-center gap-2 px-2 pb-1.5 text-[11px] text-text-6">
              <span />
              <span />
              <span>容器</span>
              <span>状态</span>
              <span>策略</span>
              <span />
              <span />
            </div>
            <div
              v-for="row in dockerRows"
              :key="row.name"
              class="grid grid-cols-[18px_22px_minmax(0,1fr)_72px_92px_54px_56px] items-center gap-2 px-2 py-1.5 text-[12.5px] text-text-5"
            >
              <span />
              <span>—</span>
              <span class="truncate" :title="row.name">{{ row.name }}</span>
              <span :class="row.state === 'running' ? '' : 'text-text-6'">
                {{ containerStateLabel(row.state) }}
              </span>
              <span class="font-mono text-[11px]">{{ row.restartPolicy || '—' }}</span>
              <span />
              <span />
            </div>
          </div>
        </div>
      </div>

      <div class="flex items-center text-[11.5px] text-text-6">
        <button class="ml-auto text-text-5 hover:text-accent" @click="gotoRunLogs">查看运行记录</button>
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
/* 拖拽的落点提示：用内阴影画一条 2px 顶线。
   不能用 border-t-2 —— 那会把这行撑高 2px，整张表在拖动过程中会上下跳，
   落点跟着偏，很多时候「拖了没反应」就是它造成的。 */
.dh-boot-drop {
  box-shadow: inset 0 2px 0 0 var(--color-accent);
}
</style>
