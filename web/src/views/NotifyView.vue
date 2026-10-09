<script setup lang="ts">
import { computed, onMounted, ref } from 'vue'
import {
  Bell,
  BellOff,
  CheckCircle2,
  History,
  Info,
  Loader2,
  Plus,
  Save,
  Send,
  Settings2,
  Trash2,
  XCircle,
} from 'lucide-vue-next'
import { api } from '@/api/client'
import type {
  Channel,
  ChannelPreset,
  EventDef,
  NotifyEventRow,
  NotifyRecord,
  NotifySettings,
} from '@/api/types'
import { relativeTime } from '@/utils/format'
import { useToastStore } from '@/stores/toast'
import EmptyState from '@/components/EmptyState.vue'
import Modal from '@/components/Modal.vue'
import UiSwitch from '@/components/UiSwitch.vue'

const toast = useToastStore()
const tab = ref<'channels' | 'events' | 'history'>('channels')

const presets = ref<ChannelPreset[]>([])
const templateVars = ref<{ key: string; label: string }[]>([])
const channels = ref<Channel[]>([])
const catalog = ref<EventDef[]>([])
const eventRows = ref<Record<string, NotifyEventRow>>({})
const settings = ref<NotifySettings | null>(null)
const history = ref<NotifyRecord[]>([])
const loading = ref(true)
const saving = ref(false)
const testingId = ref<number | null>(null)

const showEditor = ref(false)
const editingChannel = ref<Channel | null>(null)
const removeTarget = ref<Channel | null>(null)
const showVars = ref(false)

const groups = computed(() => {
  const out: { name: string; items: EventDef[] }[] = []
  for (const e of catalog.value) {
    let g = out.find((x) => x.name === e.group)
    if (!g) {
      g = { name: e.group, items: [] }
      out.push(g)
    }
    g.items.push(e)
  }
  return out
})

const enabledCount = computed(() =>
  catalog.value.filter((e) => eventRows.value[e.event]?.enabled).length,
)

const currentPreset = computed(() =>
  presets.value.find((p) => p.type === editingChannel.value?.type),
)

async function load() {
  loading.value = true
  try {
    const [p, e, s, ch, h] = await Promise.all([
      api.get<{ presets: ChannelPreset[]; vars: { key: string; label: string }[] }>('/api/notify/presets'),
      api.get<{ catalog: EventDef[]; current: NotifyEventRow[] }>('/api/notify/events'),
      api.get<{ settings: NotifySettings }>('/api/notify/settings'),
      api.get<{ channels: Channel[] }>('/api/notify/channels'),
      api.get<{ history: NotifyRecord[] }>('/api/notify/history', { limit: 80 }),
    ])
    presets.value = p.presets ?? []
    templateVars.value = p.vars ?? []
    catalog.value = e.catalog ?? []
    const map: Record<string, NotifyEventRow> = {}
    for (const def of e.catalog ?? []) {
      const cur = (e.current ?? []).find((c) => c.event === def.event)
      map[def.event] = cur ?? { event: def.event, enabled: def.default, level: def.level }
    }
    eventRows.value = map
    settings.value = s.settings
    channels.value = ch.channels ?? []
    history.value = h.history ?? []
  } catch (err) {
    toast.error('读取通知配置失败', err instanceof Error ? err.message : String(err))
  } finally {
    loading.value = false
  }
}

function openCreate() {
  const def = presets.value[0]
  editingChannel.value = {
    id: 0,
    name: def?.label ?? '新渠道',
    type: def?.type ?? 'webhook',
    enabled: true,
    config: {},
    createdAt: '',
  }
  showEditor.value = true
}

function openEdit(c: Channel) {
  editingChannel.value = { ...c, config: { ...(c.config ?? {}) } }
  showEditor.value = true
}

function needField(key: string) {
  const p = currentPreset.value
  const v = editingChannel.value?.config?.[key]
  return p?.fields.find((f) => f.key === key)?.required && !String(v ?? '').trim()
}

async function saveChannel() {
  const ch = editingChannel.value
  if (!ch) return
  if (!ch.name.trim()) {
    toast.error('请填写渠道名称')
    return
  }
  try {
    if (ch.id) {
      await api.put(`/api/notify/channels/${ch.id}`, ch)
      toast.success('渠道已更新')
    } else {
      await api.post('/api/notify/channels', ch)
      toast.success('渠道已创建')
    }
    showEditor.value = false
    await load()
  } catch (e) {
    toast.error('保存失败', e instanceof Error ? e.message : String(e))
  }
}

/** 列表里的启用开关：先乐观更新，失败再退回。 */
async function toggleChannel(c: Channel, v: boolean) {
  const prev = c.enabled
  c.enabled = v
  try {
    await api.put(`/api/notify/channels/${c.id}`, c)
  } catch (e) {
    c.enabled = prev
    toast.error('保存失败', e instanceof Error ? e.message : String(e))
  }
}

async function testChannel(c: Channel) {
  testingId.value = c.id
  try {
    const res = await api.post<{ ok: boolean; message: string }>(`/api/notify/channels/${c.id}/test`)
    if (res.ok) toast.success('测试消息已发送', `${c.name} 配置可用`)
  } catch (e) {
    toast.error('测试失败', e instanceof Error ? e.message : String(e))
  } finally {
    testingId.value = null
    void loadHistory()
  }
}

async function confirmRemove() {
  const c = removeTarget.value
  if (!c) return
  try {
    await api.del(`/api/notify/channels/${c.id}`)
    toast.success('渠道已删除')
    removeTarget.value = null
    await load()
  } catch (e) {
    toast.error('删除失败', e instanceof Error ? e.message : String(e))
  }
}

async function toggleEvent(def: EventDef) {
  const cur = eventRows.value[def.event]
  if (!cur) return
  cur.enabled = !cur.enabled
  try {
    await api.put('/api/notify/events', { events: [cur] })
  } catch (e) {
    cur.enabled = !cur.enabled
    toast.error('保存失败', e instanceof Error ? e.message : String(e))
  }
}

async function bulkEvents(preset: 'recommended' | 'all' | 'none') {
  const rows: NotifyEventRow[] = catalog.value.map((def) => ({
    event: def.event,
    enabled: preset === 'all' ? true : preset === 'none' ? false : def.default,
    level: eventRows.value[def.event]?.level ?? def.level,
  }))
  try {
    await api.put('/api/notify/events', { events: rows })
    for (const r of rows) eventRows.value[r.event] = r
    toast.success('已应用')
  } catch (e) {
    toast.error('保存失败', e instanceof Error ? e.message : String(e))
  }
}

async function saveSettings() {
  if (!settings.value) return
  saving.value = true
  try {
    const res = await api.put<{ settings: NotifySettings }>('/api/notify/settings', settings.value)
    settings.value = res.settings
    toast.success('通知设置已保存')
  } catch (e) {
    toast.error('保存失败', e instanceof Error ? e.message : String(e))
  } finally {
    saving.value = false
  }
}

async function loadHistory() {
  try {
    const res = await api.get<{ history: NotifyRecord[] }>('/api/notify/history', { limit: 80 })
    history.value = res.history ?? []
  } catch {
    /* 忽略 */
  }
}

async function clearHistory() {
  try {
    await api.del('/api/notify/history')
    history.value = []
    toast.success('推送历史已清空')
  } catch (e) {
    toast.error('清空失败', e instanceof Error ? e.message : String(e))
  }
}

const levelBadge = (lv: string) =>
  lv === 'urgent' ? 'dh-badge-err' : 'dh-badge-plain'

/**
 * 在模板里安全地渲染「双花括号」变量写法。
 * 不能直接在模板里写字面量 —— Vue 的分隔符会在字符串闭合前抢先结束插值。
 */
const tplVar = (name: string) => `{{${name}}}`

onMounted(() => void load())
</script>

<template>
  <div class="flex flex-col gap-3.5 p-[18px]">
    <div class="dh-phead">
      <div class="dh-h1">通知</div>
      <div class="dh-sub">
        {{ channels.filter((c) => c.enabled).length }} 个渠道已启用 · 已订阅 {{ enabledCount }} 个事件
      </div>
    </div>

    <div class="flex flex-wrap items-center gap-2.5">
      <div class="dh-seg">
        <button :data-on="tab === 'channels'" @click="tab = 'channels'">
          渠道 {{ channels.length }}
        </button>
        <button :data-on="tab === 'events'" @click="tab = 'events'">
          事件订阅 {{ enabledCount }}
        </button>
        <button :data-on="tab === 'history'" @click="tab = 'history'">推送历史</button>
      </div>

      <div class="ml-auto flex items-center gap-2">
        <button v-if="tab === 'channels'" class="dh-btn dh-btn-primary" @click="openCreate">
          <Plus class="h-3.5 w-3.5" />添加渠道
        </button>
        <template v-if="tab === 'events'">
          <button class="dh-btn dh-btn-sm" @click="bulkEvents('recommended')">恢复推荐</button>
          <button class="dh-btn dh-btn-sm" @click="bulkEvents('none')">全部关闭</button>
          <button class="dh-btn dh-btn-sm" @click="showVars = true">
            <Info class="h-3 w-3" />模板变量
          </button>
        </template>
        <template v-if="tab === 'history'">
          <button class="dh-btn dh-btn-sm" @click="loadHistory">刷新</button>
          <button class="dh-btn dh-btn-sm dh-btn-danger" :disabled="!history.length" @click="clearHistory">
            <Trash2 class="h-3 w-3" />清空
          </button>
        </template>
      </div>
    </div>

    <!-- 渠道 -->
    <template v-if="tab === 'channels'">
      <EmptyState
        v-if="!channels.length"
        :icon="BellOff"
        :title="loading ? '正在载入…' : '还没有配置通知渠道'"
        description="支持 Telegram、Bark、ntfy、企业微信、钉钉、飞书、Server 酱、PushPlus、邮件 SMTP，以及任意自定义 Webhook。"
        action-label="添加第一个渠道"
        @action="openCreate"
      />
      <div v-else class="dh-card">
        <div class="dh-card-head">
          <Bell class="h-3.5 w-3.5 text-text-4" />
          <span>通知渠道</span>
          <span class="ml-auto text-[11.5px] font-normal text-text-5">
            {{ channels.filter((c) => c.enabled).length }} / {{ channels.length }} 已启用
          </span>
        </div>
        <div class="flex flex-col">
          <div
            v-for="c in channels"
            :key="c.id"
            class="flex flex-wrap items-center gap-2.5 border-b border-[#171f2a] px-3.5 py-3 last:border-b-0"
          >
            <UiSwitch
              :model-value="c.enabled"
              :label="`启用 ${c.name}`"
              @update:model-value="(v: boolean) => toggleChannel(c, v)"
            />
            <div class="min-w-[150px] flex-1">
              <div class="text-[12.5px] font-medium">{{ c.name }}</div>
              <div class="text-[11px] text-text-5">
                {{ presets.find((p) => p.type === c.type)?.label ?? c.type }}
              </div>
            </div>
            <div class="flex gap-1.5">
              <button class="dh-btn dh-btn-sm" :disabled="testingId === c.id" @click="testChannel(c)">
                <Loader2 v-if="testingId === c.id" class="h-3 w-3 dh-spin" />
                <Send v-else class="h-3 w-3" />测试
              </button>
              <button class="dh-btn dh-btn-sm" @click="openEdit(c)">编辑</button>
              <button class="dh-btn dh-btn-sm dh-btn-danger" @click="removeTarget = c">
                <Trash2 class="h-3 w-3" />
              </button>
            </div>
          </div>
        </div>
      </div>

      <!-- 全局设置 -->
      <div v-if="settings" class="dh-card">
        <div class="dh-card-head">
          <Settings2 class="h-3.5 w-3.5 text-text-4" />
          <span>投递策略</span>
          <button class="dh-btn dh-btn-sm dh-btn-primary ml-auto" :disabled="saving" @click="saveSettings">
            <Save class="h-3 w-3" />保存
          </button>
        </div>
        <div class="grid grid-cols-1 gap-x-8 gap-y-4 p-3.5 lg:grid-cols-2">
          <label class="flex items-start gap-2.5">
            <UiSwitch v-model="settings.enabled" />
            <span class="text-[12.5px]">
              启用通知
              <div class="text-[11px] text-text-5">总开关。关闭后不发送任何通知，也不写推送历史。</div>
            </span>
          </label>

          <div>
            <label class="dh-label">面板地址（用于模板变量 <code>{{ tplVar('url') }}</code>）</label>
            <input v-model="settings.panelURL" class="dh-input" placeholder="http://192.168.1.10:5923" />
          </div>

          <label class="flex items-start gap-2.5">
            <UiSwitch v-model="settings.quietEnabled" />
            <span class="text-[12.5px]">
              启用静默时段
              <div class="text-[11px] text-text-5">支持跨天，例如 23:00–07:00。</div>
            </span>
          </label>

          <div class="flex items-end gap-2">
            <div class="flex-1">
              <label class="dh-label">开始</label>
              <input v-model="settings.quietStart" type="time" class="dh-input" />
            </div>
            <div class="flex-1">
              <label class="dh-label">结束</label>
              <input v-model="settings.quietEnd" type="time" class="dh-input" />
            </div>
          </div>

          <div>
            <label class="dh-label">静默期的普通事件</label>
            <select v-model="settings.quietNormalMode" class="dh-select">
              <option value="digest">攒着，时段结束后汇总发一条</option>
              <option value="drop">直接丢弃</option>
            </select>
          </div>

          <label class="flex items-start gap-2.5">
            <UiSwitch v-model="settings.quietUrgentSend" />
            <span class="text-[12.5px]">
              紧急事件在静默期照常发送
              <div class="text-[11px] text-text-5">更新失败、容器意外退出、登录告警属于紧急事件。</div>
            </span>
          </label>

          <div>
            <label class="dh-label">同容器同事件去重窗口（分钟）</label>
            <input v-model.number="settings.dedupeWindow" type="number" min="0" class="dh-input" />
            <div class="mt-1 text-[11px] text-text-5">
              崩溃循环的容器一分钟能产生几十条事件，靠这个窗口压住。
            </div>
          </div>

          <div>
            <label class="dh-label">每日推送上限（条）</label>
            <input v-model.number="settings.dailyLimit" type="number" min="0" class="dh-input" />
            <div class="mt-1 text-[11px] text-text-5">今日已发送 {{ settings.sentToday }} 条，超出上限后丢弃。</div>
          </div>

          <div class="lg:col-span-2">
            <div class="rounded-[9px] border border-line-1 bg-ink-800 px-3 py-2.5 text-[11.5px] leading-relaxed text-text-5">
              投递失败时会自动重试 3 次并按指数退避，失败只记入推送历史，
              <b class="text-text-3">绝不影响更新主流程</b>。
              当前状态：{{ settings.inQuietHours ? '处于静默时段' : '非静默时段' }}。
            </div>
          </div>
        </div>
      </div>
    </template>

    <!-- 事件订阅 -->
    <template v-else-if="tab === 'events'">
      <div class="flex items-start gap-3 rounded-[14px] border border-line-1 bg-ink-750 px-4 py-3 text-[12px] leading-relaxed text-text-4">
        <Info class="mt-[2px] h-4 w-4 flex-none text-accent" />
        <div>
          订阅是按<b>事件</b>而不是按渠道设置的：某个事件一旦开启，会同时发往所有已启用的渠道。
          标注「紧急」的事件可以在静默时段照常发送。登录成功与登录失败是两个独立事件，可以分别开关。
        </div>
      </div>

      <div v-for="g in groups" :key="g.name" class="dh-card">
        <div class="dh-card-head">
          <Bell class="h-3.5 w-3.5 text-text-4" />
          <span>{{ g.name }}</span>
          <span class="ml-auto text-[11.5px] font-normal text-text-5">
            {{ g.items.filter((i) => eventRows[i.event]?.enabled).length }} / {{ g.items.length }} 已开启
          </span>
        </div>
        <div class="flex flex-col">
          <div
            v-for="e in g.items"
            :key="e.event"
            class="flex flex-wrap items-center gap-3 border-b border-[#171f2a] px-3.5 py-2.5 last:border-b-0"
          >
            <UiSwitch
              :model-value="eventRows[e.event]?.enabled ?? false"
              :label="e.label"
              @update:model-value="toggleEvent(e)"
            />
            <div class="min-w-[200px] flex-1">
              <div class="flex items-center gap-2">
                <span class="text-[12.5px] font-medium">{{ e.label }}</span>
                <span class="dh-badge" :class="levelBadge(e.level)">
                  {{ e.level === 'urgent' ? '紧急' : '普通' }}
                </span>
                <code class="text-[10.5px] text-text-6">{{ e.event }}</code>
              </div>
              <div class="mt-0.5 text-[11px] leading-relaxed text-text-5">{{ e.description }}</div>
            </div>
          </div>
        </div>
      </div>
    </template>

    <!-- 推送历史 -->
    <template v-else>
      <div class="dh-card">
        <div class="dh-card-head">
          <History class="h-3.5 w-3.5 text-text-4" />
          <span>推送历史</span>
          <span class="ml-auto text-[11.5px] font-normal text-text-5">{{ history.length }} 条</span>
        </div>
        <EmptyState
          v-if="!history.length"
          :icon="History"
          title="还没有推送记录"
          description="通知发出后（无论成功或失败）都会在这里留下一条记录。"
        />
        <div v-else class="overflow-x-auto">
          <table class="dh-table">
            <thead>
              <tr>
                <th class="w-[140px]">时间</th>
                <th class="w-[150px]">渠道</th>
                <th class="w-[150px]">事件</th>
                <th class="w-[70px]">结果</th>
                <th>内容</th>
              </tr>
            </thead>
            <tbody>
              <tr v-for="h in history" :key="h.id">
                <td class="whitespace-nowrap text-[11.5px] text-text-5">{{ relativeTime(h.ts) }}</td>
                <td class="text-[12px] text-text-3">{{ h.channel || '—' }}</td>
                <td>
                  <span class="dh-badge" :class="levelBadge(h.level)">
                    {{ h.event }}
                  </span>
                </td>
                <td>
                  <span class="dh-badge" :class="h.ok ? 'dh-badge-run' : 'dh-badge-err'">
                    <CheckCircle2 v-if="h.ok" class="h-3 w-3" />
                    <XCircle v-else class="h-3 w-3" />
                    {{ h.ok ? '成功' : '失败' }}
                  </span>
                </td>
                <td class="max-w-[420px]">
                  <div class="truncate text-[12px] text-text-2" :title="h.title">{{ h.title }}</div>
                  <div v-if="!h.ok && h.errmsg" class="mt-0.5 truncate text-[11px] text-[#fca5a5]" :title="h.errmsg">
                    {{ h.errmsg }}
                  </div>
                </td>
              </tr>
            </tbody>
          </table>
        </div>
      </div>
    </template>

    <!-- 渠道编辑器 -->
    <Modal
      :open="showEditor"
      :title="editingChannel?.id ? '编辑渠道' : '添加通知渠道'"
      width="620px"
      @close="showEditor = false"
    >
      <div v-if="editingChannel" class="flex flex-col gap-3.5">
        <div class="grid grid-cols-2 gap-3">
          <div>
            <label class="dh-label">渠道名称</label>
            <input v-model="editingChannel.name" class="dh-input" placeholder="例如：家庭群通知" />
          </div>
          <div>
            <label class="dh-label">类型</label>
            <select
              v-model="editingChannel.type"
              class="dh-select"
              :disabled="!!editingChannel.id"
              @change="editingChannel.config = {}"
            >
              <option v-for="p in presets" :key="p.type" :value="p.type">{{ p.label }}</option>
            </select>
          </div>
        </div>

        <div v-if="currentPreset" class="text-[11.5px] leading-relaxed text-text-5">
          {{ currentPreset.description }}
        </div>

        <div class="flex flex-col gap-3">
          <div v-for="f in currentPreset?.fields ?? []" :key="f.key">
            <label class="dh-label">
              {{ f.label }}
              <span v-if="f.required" class="text-[#fca5a5]">*</span>
            </label>
            <textarea
              v-if="f.type === 'textarea'"
              :value="String(editingChannel.config[f.key] ?? '')"
              class="dh-textarea font-mono text-[11.5px]"
              :placeholder="f.placeholder"
              @input="editingChannel.config[f.key] = ($event.target as HTMLTextAreaElement).value"
            />
            <input
              v-else
              :value="String(editingChannel.config[f.key] ?? '')"
              :type="f.type === 'password' ? 'password' : 'text'"
              class="dh-input"
              :class="needField(f.key) ? '!border-[rgba(248,113,113,.5)]' : ''"
              :placeholder="f.placeholder"
              @input="editingChannel.config[f.key] = ($event.target as HTMLInputElement).value"
            />
            <div v-if="f.help" class="mt-1 text-[11px] text-text-5">{{ f.help }}</div>
          </div>
        </div>

        <label class="flex cursor-pointer items-center gap-2.5 text-[12.5px] text-text-2">
          <UiSwitch v-model="editingChannel.enabled" />
          启用这个渠道
        </label>
      </div>
      <template #footer>
        <button class="dh-btn" @click="showEditor = false">取消</button>
        <button class="dh-btn dh-btn-primary" @click="saveChannel">
          {{ editingChannel?.id ? '保存修改' : '创建渠道' }}
        </button>
      </template>
    </Modal>

    <!-- 变量说明 -->
    <Modal :open="showVars" title="模板变量" width="560px" @close="showVars = false">
      <div class="text-[12px] leading-relaxed text-text-4">
        自定义 Webhook 的 URL、请求头、请求体里都可以使用这些变量，写法是
        <code class="text-text-2">{{ tplVar('变量名') }}</code>。URL 中的变量会自动做 URL 转义。
      </div>
      <div class="mt-3 overflow-hidden rounded-[10px] border border-line-1">
        <table class="dh-table">
          <thead>
            <tr>
              <th class="w-[140px]">变量</th>
              <th>说明</th>
            </tr>
          </thead>
          <tbody>
            <tr v-for="v in templateVars" :key="v.key">
              <td class="font-mono text-[11.5px] text-accent">{{ tplVar(v.key) }}</td>
              <td class="text-[12px] text-text-3">{{ v.label }}</td>
            </tr>
          </tbody>
        </table>
      </div>
      <template #footer>
        <button class="dh-btn dh-btn-primary" @click="showVars = false">关闭</button>
      </template>
    </Modal>

    <!-- 删除渠道 -->
    <Modal
      :open="!!removeTarget"
      title="删除渠道"
      :subtitle="removeTarget?.name"
      width="400px"
      @close="removeTarget = null"
    >
      <div class="text-[12.5px] text-text-3">删除后该渠道不再接收任何通知。</div>
      <template #footer>
        <button class="dh-btn" @click="removeTarget = null">取消</button>
        <button class="dh-btn dh-btn-danger" @click="confirmRemove">确认删除</button>
      </template>
    </Modal>
  </div>
</template>
