<script setup lang="ts">
import { computed, onMounted, onUnmounted, ref } from 'vue'
import {
  CalendarClock,
  CheckCircle2,
  Loader2,
  Play,
  Plus,
  RefreshCw,
  Trash2,
  XCircle,
} from 'lucide-vue-next'
import { api, openStream } from '@/api/client'
import type { ContainerView, Schedule, ScheduleAction } from '@/api/types'
import { explainCron, relativeTime } from '@/utils/format'
import { useToastStore } from '@/stores/toast'
import EmptyState from '@/components/EmptyState.vue'
import Modal from '@/components/Modal.vue'
import UiSwitch from '@/components/UiSwitch.vue'

const toast = useToastStore()
const schedules = ref<Schedule[]>([])
const actions = ref<ScheduleAction[]>([])
const containers = ref<ContainerView[]>([])
const loading = ref(true)
const running = ref<number | null>(null)
const showEditor = ref(false)
const removeTarget = ref<Schedule | null>(null)
const saveError = ref('')
let closeStream: (() => void) | null = null

const cronPresets = [
  { label: '每天 03:00', cron: '0 3 * * *' },
  { label: '每天 04:00', cron: '0 4 * * *' },
  { label: '每小时', cron: '0 * * * *' },
  { label: '每 6 小时', cron: '0 */6 * * *' },
  { label: '每 30 分钟', cron: '*/30 * * * *' },
  { label: '每周一 03:00', cron: '0 3 * * 1' },
  { label: '每月 1 号 03:00', cron: '0 3 1 * *' },
]

const form = ref<Schedule>({
  id: 0,
  name: '',
  cron: '0 3 * * *',
  action: 'restart',
  targets: [],
  enabled: true,
  lastRun: '',
  lastStatus: '',
  lastMessage: '',
  createdAt: '',
})

const actionMeta = computed(() => actions.value.find((a) => a.key === form.value.action))

async function load() {
  loading.value = true
  try {
    const [s, a, c] = await Promise.all([
      api.get<{ schedules: Schedule[] }>('/api/schedules'),
      api.get<{ actions: ScheduleAction[] }>('/api/schedules/actions'),
      api.get<{ containers: ContainerView[] }>('/api/containers'),
    ])
    schedules.value = s.schedules ?? []
    actions.value = a.actions ?? []
    containers.value = c.containers ?? []
  } catch (e) {
    toast.error('读取计划任务失败', e instanceof Error ? e.message : String(e))
  } finally {
    loading.value = false
  }
}

function openCreate() {
  form.value = {
    id: 0,
    name: '',
    cron: '0 3 * * *',
    action: 'restart',
    targets: [],
    enabled: true,
    lastRun: '',
    lastStatus: '',
    lastMessage: '',
    createdAt: '',
  }
  saveError.value = ''
  showEditor.value = true
}

function openEdit(s: Schedule) {
  form.value = { ...s, targets: [...(s.targets ?? [])] }
  saveError.value = ''
  showEditor.value = true
}

async function save() {
  saveError.value = ''
  try {
    if (form.value.id) {
      await api.put(`/api/schedules/${form.value.id}`, form.value)
      toast.success('任务已更新')
    } else {
      await api.post('/api/schedules', form.value)
      toast.success('任务已创建')
    }
    showEditor.value = false
    await load()
  } catch (e) {
    saveError.value = e instanceof Error ? e.message : String(e)
  }
}

async function toggleEnabled(s: Schedule) {
  try {
    await api.put(`/api/schedules/${s.id}`, { ...s, enabled: !s.enabled })
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
  if (!s) return
  try {
    await api.del(`/api/schedules/${s.id}`)
    toast.success('任务已删除')
    removeTarget.value = null
    await load()
  } catch (e) {
    toast.error('删除失败', e instanceof Error ? e.message : String(e))
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

const actionLabel = (key: string) => actions.value.find((a) => a.key === key)?.label ?? key

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
    <div class="flex flex-wrap items-center gap-2.5">
      <div class="text-[12.5px] text-text-4">
        用标准 5 段 cron 表达式（分 时 日 月 周）定时执行容器的启动 / 停止 / 重启 / 更新 / 备份。
      </div>
      <div class="ml-auto flex gap-2">
        <button class="dh-btn" :disabled="loading" @click="load">
          <RefreshCw class="h-3.5 w-3.5" :class="loading ? 'dh-spin' : ''" />刷新
        </button>
        <button class="dh-btn dh-btn-primary" @click="openCreate">
          <Plus class="h-3.5 w-3.5" />新建任务
        </button>
      </div>
    </div>

    <div class="dh-card">
      <div class="dh-card-head">
        <CalendarClock class="h-3.5 w-3.5 text-text-4" />
        <span>计划任务</span>
        <span class="ml-auto text-[11.5px] font-normal text-text-5">{{ schedules.length }} 个</span>
      </div>

      <EmptyState
        v-if="!schedules.length"
        :icon="CalendarClock"
        :title="loading ? '正在载入…' : '还没有计划任务'"
        description="例如：每天凌晨 3 点自动检查并更新所有容器；或者每晚 23:30 停止某个占资源的容器。"
        action-label="新建第一个任务"
        @action="openCreate"
      />

      <div v-else class="overflow-x-auto">
        <table class="dh-table">
          <thead>
            <tr>
              <th class="w-[60px]">启用</th>
              <th>名称</th>
              <th class="w-[150px]">计划</th>
              <th class="w-[120px]">动作</th>
              <th class="w-[140px]">目标</th>
              <th class="w-[150px]">上次执行</th>
              <th class="w-[120px]">下次执行</th>
              <th class="w-[130px]" />
            </tr>
          </thead>
          <tbody>
            <tr v-for="s in schedules" :key="s.id">
              <td>
                <UiSwitch :model-value="s.enabled" :label="`启用 ${s.name}`" @update:model-value="toggleEnabled(s)" />
              </td>
              <td>
                <div class="text-[12.5px] font-medium">{{ s.name }}</div>
                <div v-if="s.lastMessage" class="mt-0.5 max-w-[280px] truncate text-[11px] text-text-5" :title="s.lastMessage">
                  {{ s.lastMessage }}
                </div>
              </td>
              <td>
                <div class="font-mono text-[11.5px] text-text-2">{{ s.cron }}</div>
                <div class="text-[11px] text-text-5">{{ explainCron(s.cron) }}</div>
              </td>
              <td><span class="dh-badge dh-badge-plain">{{ actionLabel(s.action) }}</span></td>
              <td>
                <span v-if="!s.targets.length" class="dh-badge dh-badge-accent">全部容器</span>
                <span v-else class="dh-badge dh-badge-plain" :title="s.targets.join(', ')">
                  {{ s.targets.length }} 个
                </span>
              </td>
              <td>
                <div class="flex items-center gap-1.5">
                  <span class="dh-badge" :class="statusBadge(s).cls">
                    <CheckCircle2 v-if="s.lastStatus === 'success'" class="h-3 w-3" />
                    <XCircle v-else-if="s.lastStatus === 'failed'" class="h-3 w-3" />
                    {{ statusBadge(s).text }}
                  </span>
                </div>
                <div class="mt-1 text-[11px] text-text-5">{{ s.lastRun ? relativeTime(s.lastRun) : '—' }}</div>
              </td>
              <td class="text-[11.5px] text-text-4">
                {{ s.enabled ? (s.nextRun ? relativeTime(s.nextRun) : '—') : '已停用' }}
              </td>
              <td>
                <div class="flex gap-1.5">
                  <button class="dh-btn dh-btn-sm" :disabled="running === s.id" @click="runNow(s)">
                    <Loader2 v-if="running === s.id" class="h-3 w-3 dh-spin" />
                    <Play v-else class="h-3 w-3" />运行
                  </button>
                  <button class="dh-btn dh-btn-sm" @click="openEdit(s)">编辑</button>
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

    <!-- 编辑器 -->
    <Modal
      :open="showEditor"
      :title="form.id ? '编辑计划任务' : '新建计划任务'"
      width="600px"
      @close="showEditor = false"
    >
      <div class="flex flex-col gap-3.5">
        <div>
          <label class="dh-label">任务名称</label>
          <input v-model="form.name" class="dh-input" placeholder="例如：每天凌晨更新全部容器" />
        </div>

        <div>
          <label class="dh-label">执行动作</label>
          <div class="grid grid-cols-1 gap-2 sm:grid-cols-2">
            <button
              v-for="a in actions"
              :key="a.key"
              type="button"
              class="rounded-[10px] border px-3 py-2 text-left transition-colors"
              :class="
                form.action === a.key
                  ? 'border-[rgba(45,212,191,.5)] bg-[#12201f]'
                  : 'border-line-1 bg-ink-800 hover:border-line-4'
              "
              @click="form.action = a.key"
            >
              <div class="text-[12.5px] font-medium" :class="form.action === a.key ? 'text-accent' : 'text-text-2'">
                {{ a.label }}
              </div>
              <div class="mt-0.5 text-[11px] leading-relaxed text-text-5">{{ a.description }}</div>
            </button>
          </div>
        </div>

        <div>
          <label class="dh-label">计划（cron 表达式）</label>
          <input v-model="form.cron" class="dh-input font-mono" placeholder="0 3 * * *" />
          <div class="mt-1.5 flex flex-wrap gap-1.5">
            <button
              v-for="p in cronPresets"
              :key="p.cron"
              type="button"
              class="rounded-full border px-2.5 py-[3px] text-[11.5px] transition-colors"
              :class="
                form.cron === p.cron
                  ? 'border-[rgba(45,212,191,.5)] bg-[#10231f] text-accent'
                  : 'border-line-3 text-text-3 hover:border-line-4'
              "
              @click="form.cron = p.cron"
            >
              {{ p.label }}
            </button>
          </div>
          <div class="mt-1.5 text-[11.5px] text-text-5">
            解析结果：<b class="text-text-3">{{ explainCron(form.cron) || '无法解析' }}</b>
          </div>
        </div>

        <div v-if="actionMeta?.needsTargets">
          <label class="dh-label">
            目标容器
            <span class="text-text-6">（不选 = 全部容器，会自动排除 Dockhelm 自身与排除列表）</span>
          </label>
          <div class="dh-scroll max-h-[190px] overflow-auto rounded-[10px] border border-line-1 bg-ink-800 p-2">
            <label
              v-for="c in containers"
              :key="c.id"
              class="flex cursor-pointer items-center gap-2.5 rounded-md px-2 py-1.5 hover:bg-ink-750"
            >
              <input
                type="checkbox"
                class="h-[14px] w-[14px] accent-[#2dd4bf]"
                :checked="form.targets.includes(c.name)"
                @change="toggleTarget(c.name)"
              />
              <span class="min-w-0 flex-1 truncate text-[12px] text-text-2">{{ c.name }}</span>
              <span v-if="c.self" class="dh-badge dh-badge-accent">自身</span>
              <span v-else-if="c.excluded" class="dh-badge dh-badge-plain">已排除</span>
              <span class="dh-badge" :class="c.state === 'running' ? 'dh-badge-run' : 'dh-badge-stop'">
                {{ c.state === 'running' ? '运行' : '停止' }}
              </span>
            </label>
          </div>
          <div class="mt-1 text-[11px] text-text-5">已选 {{ form.targets.length }} 个</div>
        </div>

        <label class="flex cursor-pointer items-center gap-2.5 text-[12.5px] text-text-2">
          <UiSwitch v-model="form.enabled" />
          创建后立即启用
        </label>

        <div
          v-if="saveError"
          class="rounded-[9px] border border-[rgba(248,113,113,.35)] bg-[rgba(248,113,113,.08)] px-3 py-2 text-[12px] text-[#fca5a5]"
        >
          {{ saveError }}
        </div>
      </div>
      <template #footer>
        <button class="dh-btn" @click="showEditor = false">取消</button>
        <button class="dh-btn dh-btn-primary" :disabled="!form.name.trim()" @click="save">
          {{ form.id ? '保存修改' : '创建任务' }}
        </button>
      </template>
    </Modal>

    <Modal
      :open="!!removeTarget"
      title="删除计划任务"
      :subtitle="removeTarget?.name"
      width="400px"
      @close="removeTarget = null"
    >
      <div class="text-[12.5px] text-text-3">删除后该任务不再自动执行，已有的执行记录仍然保留。</div>
      <template #footer>
        <button class="dh-btn" @click="removeTarget = null">取消</button>
        <button class="dh-btn dh-btn-danger" @click="confirmRemove">确认删除</button>
      </template>
    </Modal>
  </div>
</template>
