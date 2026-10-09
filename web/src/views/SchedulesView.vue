<script setup lang="ts">
import { computed, onMounted, onUnmounted, ref } from 'vue'
import { Loader2, Play, Plus, RefreshCw, Trash2 } from 'lucide-vue-next'
import { api, openStream } from '@/api/client'
import type { ContainerView, RunLog, Schedule, ScheduleAction, Settings } from '@/api/types'
import { explainCron, formatDayTime } from '@/utils/format'
import { useToastStore } from '@/stores/toast'
import EmptyState from '@/components/EmptyState.vue'
import Modal from '@/components/Modal.vue'
import UiSwitch from '@/components/UiSwitch.vue'

const toast = useToastStore()
const schedules = ref<Schedule[]>([])
const actions = ref<ScheduleAction[]>([])
const containers = ref<ContainerView[]>([])
const excluded = ref<string[]>([])
const loading = ref(true)
const running = ref<number | null>(null)
const removeTarget = ref<Schedule | null>(null)
const saveError = ref('')
const saving = ref(false)
const removing = ref(false)
let closeStream: (() => void) | null = null

/** 表单当前是在「新建」还是「编辑」——两者共用同一张内联卡片。 */
const editingId = ref(0)

const form = ref({
  name: '',
  action: 'restart',
  targets: [] as string[],
  enabled: true,
  /** 重复方式 + 时间的组合，最终换算成 cron；repeat='custom' 时直接用 customCron。 */
  repeat: 'daily',
  time: '03:00',
  customCron: '0 3 * * *',
})

/** 空表单。 */
function blankForm() {
  return {
    name: '',
    action: 'restart',
    targets: [] as string[],
    enabled: true,
    repeat: 'daily',
    time: '03:00',
    customCron: '0 3 * * *',
  }
}

const repeatOptions = [
  { key: 'daily', label: '每天' },
  { key: 'hourly', label: '每小时' },
  { key: 'weekdays', label: '工作日（周一至周五）' },
  { key: 'weekly-1', label: '每周一' },
  { key: 'weekly-2', label: '每周二' },
  { key: 'weekly-3', label: '每周三' },
  { key: 'weekly-4', label: '每周四' },
  { key: 'weekly-5', label: '每周五' },
  { key: 'weekly-6', label: '每周六' },
  { key: 'weekly-0', label: '每周日' },
  { key: 'monthly', label: '每月 1 日' },
  { key: 'custom', label: '自定义 cron' },
]

/** 动作的短名（设计稿里的 chip 用短词）。 */
const SHORT_ACTION: Record<string, string> = {
  start: '启动',
  stop: '停止',
  restart: '重启',
  update: '更新镜像',
  backup: '备份快照',
  prune_images: '清理旧镜像',
}

/** 动作徽标的配色（对齐设计稿：重启=青、更新=琥珀、停止=灰、启动=绿）。 */
const ACTION_TONE: Record<string, string> = {
  start: 'dh-badge-run',
  stop: 'dh-badge-stop',
  restart: 'dh-badge-accent',
  update: 'dh-badge-warn',
  backup: 'dh-badge-plain',
  prune_images: 'dh-badge-plain',
}

function shortAction(key: string) {
  return SHORT_ACTION[key] ?? actions.value.find((a) => a.key === key)?.label ?? key
}

function actionTone(key: string) {
  return ACTION_TONE[key] ?? 'dh-badge-plain'
}

/** 由「重复 + 时间」拼出 cron（自定义时直接用输入的表达式）。 */
const composedCron = computed(() => {
  if (form.value.repeat === 'custom') return form.value.customCron.trim()
  const [h = '0', m = '0'] = form.value.time.split(':')
  const hh = String(Number(h))
  const mm = String(Number(m))
  switch (form.value.repeat) {
    case 'daily':
      return `${mm} ${hh} * * *`
    case 'hourly':
      return `${mm} * * * *`
    case 'weekdays':
      return `${mm} ${hh} * * 1-5`
    case 'monthly':
      return `${mm} ${hh} 1 * *`
    default: {
      const dow = form.value.repeat.split('-')[1] ?? '0'
      return `${mm} ${hh} * * ${dow}`
    }
  }
})

/** cron → 表单（编辑已有任务时把表达式翻译回「重复 + 时间」）。 */
function applyCron(cron: string) {
  const t = cron.trim()
  const hhmm = (h: string, m: string) => `${String(h).padStart(2, '0')}:${String(m).padStart(2, '0')}`
  let m = /^(\d+)\s+(\d+)\s+\*\s+\*\s+\*$/.exec(t)
  if (m) return { repeat: 'daily', time: hhmm(m[2] as string, m[1] as string) }
  m = /^(\d+)\s+(\d+)\s+\*\s+\*\s+1-5$/.exec(t)
  if (m) return { repeat: 'weekdays', time: hhmm(m[2] as string, m[1] as string) }
  m = /^(\d+)\s+(\d+)\s+\*\s+\*\s+([0-7])$/.exec(t)
  if (m) return { repeat: `weekly-${m[3] === '7' ? '0' : (m[3] as string)}`, time: hhmm(m[2] as string, m[1] as string) }
  m = /^(\d+)\s+(\d+)\s+1\s+\*\s+\*$/.exec(t)
  if (m) return { repeat: 'monthly', time: hhmm(m[2] as string, m[1] as string) }
  m = /^(\d+)\s+\*\s+\*\s+\*\s+\*$/.exec(t)
  if (m) return { repeat: 'hourly', time: hhmm('0', m[1] as string) }
  return { repeat: 'custom', time: '03:00' }
}

const actionMeta = computed(() => actions.value.find((a) => a.key === form.value.action))
/** 这个动作是否必须明确挑容器（清理悬空镜像这类全局动作不需要）。 */
const needsTargets = computed(() => actionMeta.value?.needsTargets !== false)

/**
 * 组装提交给后端的任务体。
 *
 * 后端 scheduleReq 开了 DisallowUnknownFields，**只**接受 name / cron / action / targets / enabled
 * 这五个键。早先这里把整条 Schedule（带 id、lastRun、lastStatus…）一起提交，后端直接回
 * `json: unknown field "id"` —— 新建和编辑都提交不了。
 */
function toPayload(enabled: boolean) {
  return {
    name: form.value.name.trim(),
    cron: composedCron.value,
    action: form.value.action,
    targets: needsTargets.value ? [...form.value.targets] : [],
    enabled,
  }
}

/** 这些容器永远不会被计划任务作用到，勾了也没用。 */
function blockedTarget(c: ContainerView) {
  return !!c.self || excluded.value.includes(c.name)
}

function blockedReason(c: ContainerView) {
  if (c.self) return '这是 Dockhelm 自己，任何计划任务都不会作用到它'
  if (excluded.value.includes(c.name)) return '在设置页的排除列表里，任何计划任务都会跳过它'
  return '选中这个容器'
}

async function load() {
  loading.value = true
  try {
    const [s, a, c, st] = await Promise.all([
      api.get<{ schedules: Schedule[] }>('/api/schedules'),
      api.get<{ actions: ScheduleAction[] }>('/api/schedules/actions'),
      api.get<{ containers: ContainerView[] }>('/api/containers'),
      api.get<Settings>('/api/settings').catch(() => null),
    ])
    schedules.value = s.schedules ?? []
    actions.value = a.actions ?? []
    containers.value = c.containers ?? []
    excluded.value = st?.exclude ?? []
  } catch (e) {
    toast.error('读取计划任务失败', e instanceof Error ? e.message : String(e))
  } finally {
    loading.value = false
  }
}

function resetForm() {
  form.value = blankForm()
  editingId.value = 0
  saveError.value = ''
}

/** 「新建任务」按钮：清空表单并把内联卡片滚进视野。 */
function focusForm() {
  resetForm()
  document.getElementById('schedule-form')?.scrollIntoView({ behavior: 'smooth', block: 'nearest' })
}

function startEdit(s: Schedule) {
  const mapped = applyCron(s.cron)
  form.value = {
    name: s.name,
    action: s.action,
    targets: [...(s.targets ?? [])],
    enabled: s.enabled,
    repeat: mapped.repeat,
    time: mapped.time,
    customCron: s.cron,
  }
  editingId.value = s.id
  saveError.value = ''
  document.getElementById('schedule-form')?.scrollIntoView({ behavior: 'smooth', block: 'nearest' })
}

async function save() {
  if (!form.value.name.trim()) {
    saveError.value = '请填写任务名称'
    return
  }
  if (needsTargets.value && !form.value.targets.length) {
    saveError.value = '请至少选择一个容器 ——「不选」不再等于全部容器'
    return
  }
  saveError.value = ''
  saving.value = true
  try {
    if (editingId.value) {
      await api.put(`/api/schedules/${editingId.value}`, toPayload(form.value.enabled))
      toast.success('任务已更新')
    } else {
      await api.post('/api/schedules', toPayload(form.value.enabled))
      toast.success('任务已创建')
    }
    resetForm()
    await load()
  } catch (e) {
    saveError.value = e instanceof Error ? e.message : String(e)
  } finally {
    saving.value = false
  }
}

async function toggleEnabled(s: Schedule) {
  try {
    // 单独的口子：启用状态与任务定义无关，不该把整份定义再提交一遍去过校验
    // （老的空目标任务那样会连停用都被挡住）。
    await api.post(`/api/schedules/${s.id}/enabled`, { enabled: !s.enabled })
    await load()
  } catch (e) {
    toast.error('更新失败', e instanceof Error ? e.message : String(e))
  }
}

async function runNow(s: Schedule) {
  running.value = s.id
  try {
    const res = await api.post<{ ok: boolean; message: string }>(`/api/schedules/${s.id}/run`)
    if (res.ok) toast.success('执行完成', res.message)
    else toast.error('执行失败', res.message)
    await load()
  } catch (e) {
    toast.error('执行失败', e instanceof Error ? e.message : String(e))
  } finally {
    running.value = null
  }
}

async function confirmRemove() {
  const s = removeTarget.value
  if (!s || removing.value) return
  removing.value = true
  try {
    await api.del(`/api/schedules/${s.id}`)
    toast.success('任务已删除')
    if (editingId.value === s.id) resetForm()
    removeTarget.value = null
    await load()
  } catch (e) {
    toast.error('删除失败', e instanceof Error ? e.message : String(e))
  } finally {
    removing.value = false
  }
}

function toggleTarget(name: string) {
  const list = form.value.targets
  if (list.includes(name)) form.value.targets = list.filter((t) => t !== name)
  else form.value.targets = [...list, name]
}

const statusBadge = (s: Schedule) => {
  if (!s.lastStatus) return { text: '未执行', cls: 'dh-badge-plain' }
  if (s.lastStatus === 'success') return { text: '成功', cls: 'dh-badge-run' }
  if (s.lastStatus === 'failed') return { text: '失败', cls: 'dh-badge-err' }
  return { text: s.lastStatus, cls: 'dh-badge-plain' }
}

/** 表头副标题：几个任务、几个今天还要跑。 */
const summary = computed(() => {
  const todayEnd = new Date()
  todayEnd.setHours(23, 59, 59, 999)
  const today = schedules.value.filter((s) => {
    if (!s.enabled || !s.nextRun) return false
    const t = Date.parse(s.nextRun)
    return !Number.isNaN(t) && t <= todayEnd.getTime()
  }).length
  return `${schedules.value.length} 个任务 · ${today} 个今天待执行`
})

/** 未来 24 小时内的执行安排，按时间排序。 */
const upcoming = computed(() => {
  const now = Date.now()
  return schedules.value
    .filter((s) => s.enabled && s.nextRun)
    .map((s) => ({ s, t: Date.parse(s.nextRun as string) }))
    .filter((x) => !Number.isNaN(x.t) && x.t <= now + 24 * 3600 * 1000)
    .sort((a, b) => a.t - b.t)
    .slice(0, 6)
})

/**
 * 目标列。
 *
 * 空目标**不再**等于「全部容器」——清理类动作本来就不需要目标，写「不适用」；
 * 其余动作的空目标是历史遗留（旧语义留下的），明确说清它不会执行。
 */
function targetsText(s: Schedule) {
  if (s.targets?.length) return s.targets
  if (s.action === 'prune_images') return ['不适用（全局动作）']
  return ['未指定 —— 任务不会执行']
}

// —— 执行历史 ——
const historyOpen = ref(false)
const historyLogs = ref<RunLog[]>([])
const historyLoading = ref(false)

async function openHistory() {
  historyOpen.value = true
  historyLoading.value = true
  try {
    const res = await api.get<{ logs: RunLog[] }>('/api/logs?limit=200')
    historyLogs.value = (res.logs ?? []).filter((l) => l.kind === 'schedule').slice(0, 60)
  } catch (e) {
    toast.error('读取执行历史失败', e instanceof Error ? e.message : String(e))
  } finally {
    historyLoading.value = false
  }
}

function historyStatus(l: RunLog) {
  if (l.status === 'success') return { text: '成功', cls: 'dh-badge-run' }
  if (l.status === 'failed') return { text: '失败', cls: 'dh-badge-err' }
  return { text: l.status || '—', cls: 'dh-badge-plain' }
}

onMounted(() => {
  void load()
  closeStream = openStream('/api/events/stream', (topic) => {
    if (topic === 'schedule') void load()
  })
})
onUnmounted(() => closeStream?.())
</script>

<template>
  <div class="flex flex-col gap-3.5 p-[18px]">
    <!-- 页头 -->
    <div class="dh-phead">
      <div class="dh-h1">计划任务</div>
      <div class="dh-sub">{{ summary }}</div>
      <div class="ml-auto flex gap-2">
        <button class="dh-btn" @click="openHistory">执行历史</button>
        <button class="dh-btn" :disabled="loading" @click="load">
          <RefreshCw class="h-3.5 w-3.5" :class="loading ? 'dh-spin' : ''" />刷新
        </button>
        <button class="dh-btn dh-btn-primary" @click="focusForm">
          <Plus class="h-3.5 w-3.5" />新建任务
        </button>
      </div>
    </div>

    <!-- 任务表 -->
    <div class="dh-card">
      <EmptyState
        v-if="!schedules.length"
        :icon="Plus"
        :title="loading ? '正在载入…' : '还没有计划任务'"
        description="在下面的「新建任务」里挑几个容器、选个时间和动作，就能定时启停或更新它们。"
      />
      <div v-else class="overflow-x-auto">
        <table class="dh-table">
          <thead>
            <tr>
              <th>任务名称</th>
              <th class="w-[150px]">目标</th>
              <th class="w-[110px]">动作</th>
              <th class="w-[170px]">计划</th>
              <th class="w-[150px]">上次执行</th>
              <th class="w-[150px]">下次执行</th>
              <th class="w-[90px] text-right">启用</th>
              <th class="w-[140px]" />
            </tr>
          </thead>
          <tbody>
            <tr v-for="s in schedules" :key="s.id" :class="s.enabled ? '' : 'opacity-55'">
              <td>
                <div class="text-[13px] font-semibold">{{ s.name }}</div>
                <div v-if="s.lastMessage" class="mt-0.5 max-w-[280px] truncate text-[11.5px] text-text-5" :title="s.lastMessage">
                  {{ s.lastMessage }}
                </div>
              </td>
              <td class="text-[11.5px] text-text-5">
                <div class="line-clamp-2" :class="!s.targets?.length && s.action !== 'prune_images' ? 'text-err-text' : ''">
                  {{ targetsText(s).join('、') }}
                </div>
              </td>
              <td>
                <span class="dh-badge" :class="actionTone(s.action)">{{ shortAction(s.action) }}</span>
              </td>
              <td class="text-[11.5px] text-text-5">
                <div>{{ explainCron(s.cron) }}</div>
                <div class="font-mono text-[11px]">{{ s.cron }}</div>
              </td>
              <td class="text-[11.5px] text-text-5">
                <div>{{ s.lastRun ? formatDayTime(s.lastRun) : '—' }}</div>
                <div v-if="s.lastStatus" :class="s.lastStatus === 'failed' ? 'text-err-text' : 'text-run-text'">
                  {{ statusBadge(s).text }}<template v-if="s.lastStatus === 'success' && s.targets?.length">
                    {{ ' ' + s.targets.length }}</template>
                </div>
              </td>
              <td class="text-[11.5px] text-text-5">
                {{ s.enabled ? (s.nextRun ? formatDayTime(s.nextRun) : '—') : '已停用' }}
              </td>
              <td class="text-right">
                <UiSwitch
                  :model-value="s.enabled"
                  :label="`启用 ${s.name}`"
                  @update:model-value="toggleEnabled(s)"
                />
              </td>
              <td>
                <div class="flex justify-end gap-1.5">
                  <button class="dh-btn dh-btn-sm" :disabled="running === s.id" @click="runNow(s)">
                    <Loader2 v-if="running === s.id" class="h-3 w-3 dh-spin" />
                    <Play v-else class="h-3 w-3" />运行
                  </button>
                  <button class="dh-btn dh-btn-sm" @click="startEdit(s)">编辑</button>
                  <button class="dh-btn dh-btn-sm dh-btn-danger" @click="removeTarget = s">
                    <Trash2 class="h-3 w-3" />
                  </button>
                </div>
              </td>
            </tr>
          </tbody>
        </table>
      </div>
    </div>

    <!-- 新建/编辑 + 时间轴 -->
    <div class="grid grid-cols-1 gap-3.5 xl:grid-cols-[1fr_.72fr]">
      <div id="schedule-form" class="dh-card">
        <div class="dh-card-head">
          <span>{{ editingId ? '编辑任务' : '新建任务' }}</span>
          <button v-if="editingId" class="ml-auto text-[11.5px] font-normal text-text-5 hover:text-accent" @click="resetForm">
            取消编辑
          </button>
        </div>
        <div class="dh-card-body flex flex-col gap-3.5">
          <div>
            <label class="dh-label">任务名称</label>
            <input v-model="form.name" class="dh-input" placeholder="例如：夜间重启下载器" />
          </div>

          <div v-if="needsTargets">
            <label class="dh-label">
              选择容器
              <span class="text-text-6">（必选；Dockhelm 自身与排除列表里的容器选不了）</span>
            </label>
            <div class="dh-scroll flex max-h-[132px] flex-wrap gap-1.5 overflow-auto rounded-[10px] border border-line-1 bg-ink-800 p-2">
              <button
                v-for="c in containers"
                :key="c.id"
                type="button"
                class="dh-chip"
                :data-on="form.targets.includes(c.name)"
                :disabled="blockedTarget(c)"
                :title="blockedReason(c)"
                @click="toggleTarget(c.name)"
              >
                {{ c.name }}
                <span v-if="c.self" class="text-[10px] text-text-5">自身</span>
                <span v-else-if="excluded.includes(c.name)" class="text-[10px] text-text-5">已排除</span>
              </button>
              <div v-if="!containers.length" class="p-1 text-[11.5px] text-text-5">读不到容器列表</div>
            </div>
            <div class="mt-1 text-[11px]" :class="form.targets.length ? 'text-text-5' : 'text-err-text'">
              <template v-if="form.targets.length">已选 {{ form.targets.length }} 个</template>
              <template v-else>还没选容器 —— 至少要选一个，任务才会执行</template>
            </div>
          </div>

          <div>
            <label class="dh-label">执行动作</label>
            <div class="flex flex-wrap gap-1.5">
              <button
                v-for="a in actions"
                :key="a.key"
                type="button"
                class="dh-chip"
                :data-on="form.action === a.key"
                :title="a.description"
                @click="form.action = a.key"
              >
                {{ shortAction(a.key) }}
              </button>
            </div>
          </div>

          <div class="grid grid-cols-[1.4fr_1fr] gap-2.5">
            <div>
              <label class="dh-label">重复</label>
              <select v-model="form.repeat" class="dh-select">
                <option v-for="r in repeatOptions" :key="r.key" :value="r.key">{{ r.label }}</option>
              </select>
            </div>
            <div>
              <label class="dh-label">{{ form.repeat === 'hourly' ? '每小时的第几分' : '时间' }}</label>
              <input
                v-if="form.repeat !== 'custom'"
                v-model="form.time"
                type="time"
                class="dh-input font-mono"
              />
              <input v-else v-model="form.customCron" class="dh-input font-mono" placeholder="0 3 * * *" />
            </div>
          </div>

          <div class="text-[11.5px] text-text-5">
            解析结果：<b class="font-mono text-text-3">{{ composedCron }}</b>
            <span v-if="explainCron(composedCron) && explainCron(composedCron) !== composedCron">
              · {{ explainCron(composedCron) }}
            </span>
          </div>

          <label class="flex cursor-pointer items-center gap-2.5 text-[12.5px] text-text-2">
            <UiSwitch v-model="form.enabled" />
            创建后立即启用
          </label>

          <div v-if="saveError" class="dh-banner dh-banner-err py-2 text-[12px]">{{ saveError }}</div>

          <button class="dh-btn dh-btn-primary" :disabled="saving" @click="save">
            <Loader2 v-if="saving" class="h-3.5 w-3.5 dh-spin" />
            {{ editingId ? '保存修改' : '创建任务' }}
          </button>
        </div>
      </div>

      <div class="dh-card">
        <div class="dh-card-head">时间轴 · 未来 24 小时</div>
        <div v-if="!upcoming.length" class="dh-card-body text-[12px] text-text-5">
          未来 24 小时内没有安排。
        </div>
        <div v-else class="dh-card-body flex flex-col gap-2.5">
          <div v-for="u in upcoming" :key="u.s.id" class="flex items-center gap-2.5">
            <span class="w-[74px] flex-none font-mono text-[11.5px] text-accent-text">
              {{ formatDayTime(u.t) }}
            </span>
            <span class="min-w-0 flex-1 truncate text-[12.5px]">{{ u.s.name }}</span>
            <span class="dh-badge flex-none" :class="actionTone(u.s.action)">
              {{ shortAction(u.s.action) }}
            </span>
          </div>
        </div>
      </div>
    </div>

    <!-- 删除确认 -->
    <Modal
      :open="!!removeTarget"
      title="删除计划任务"
      :subtitle="removeTarget?.name"
      width="400px"
      :busy="removing"
      @close="removeTarget = null"
    >
      <div class="text-[12.5px] text-text-3">删除后该任务不再自动执行，已有的执行记录仍然保留。</div>
      <template #footer>
        <button class="dh-btn" :disabled="removing" @click="removeTarget = null">取消</button>
        <button class="dh-btn dh-btn-danger" :disabled="removing" @click="confirmRemove">
          {{ removing ? '删除中…' : '确认删除' }}
        </button>
      </template>
    </Modal>

    <!-- 执行历史 -->
    <Modal :open="historyOpen" title="执行历史" subtitle="仅计划任务" width="720px" @close="historyOpen = false">
      <div v-if="historyLoading" class="grid h-24 place-items-center text-text-5">
        <Loader2 class="h-4 w-4 dh-spin" />
      </div>
      <div v-else-if="!historyLogs.length" class="py-8 text-center text-[12.5px] text-text-5">
        还没有计划任务的执行记录。
      </div>
      <div v-else class="overflow-x-auto">
        <table class="dh-table">
          <thead>
            <tr>
              <th class="w-[150px]">时间</th>
              <th class="w-[180px]">任务</th>
              <th>说明</th>
              <th class="w-[90px]">结果</th>
            </tr>
          </thead>
          <tbody>
            <tr v-for="l in historyLogs" :key="l.id">
              <td class="whitespace-nowrap text-[11.5px] text-text-5">{{ formatDayTime(l.ts) }}</td>
              <td class="max-w-[180px] truncate text-[12px] text-text-2">{{ l.ref || '—' }}</td>
              <td class="text-[12px] text-text-3">{{ l.message }}</td>
              <td><span class="dh-badge" :class="historyStatus(l).cls">{{ historyStatus(l).text }}</span></td>
            </tr>
          </tbody>
        </table>
      </div>
    </Modal>
  </div>
</template>
