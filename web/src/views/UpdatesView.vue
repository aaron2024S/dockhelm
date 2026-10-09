<script setup lang="ts">
import { computed, onMounted, onUnmounted, ref } from 'vue'
import { RouterLink } from 'vue-router'
import {
  AlertTriangle,
  CheckCircle2,
  Clock,
  Download,
  HelpCircle,
  Info,
  Loader2,
  Play,
  RefreshCw,
  ShieldQuestion,
  ShieldAlert,
  Zap,
} from 'lucide-vue-next'
import { api, openStream } from '@/api/client'
import type { AutoUpdateInfo, CheckResult, Settings, UpdatesResponse } from '@/api/types'
import { checkLabel, formatDateTime, relativeTime, shortImage } from '@/utils/format'
import { useToastStore } from '@/stores/toast'
import EmptyState from '@/components/EmptyState.vue'
import Modal from '@/components/Modal.vue'
import SettingRow from '@/components/SettingRow.vue'
import ToggleSwitch from '@/components/ToggleSwitch.vue'

const toast = useToastStore()
const results = ref<CheckResult[]>([])
const checkedAt = ref('')
const selfName = ref('')
const loading = ref(false)
const deepLoading = ref(false)
const applying = ref(false)
const selected = ref<Set<string>>(new Set())
const progress = ref<{ name: string; message: string; status: string }[]>([])
const showConfirm = ref(false)
const forceUpdate = ref(false)
const showInfo = ref(false)
let closeStream: (() => void) | null = null

// —— 自动更新 ——
const auto = ref<AutoUpdateInfo | null>(null)
const autoBusy = ref(false)
const policy = ref<Settings | null>(null)
const savingPolicy = ref(false)

const available = computed(() => results.value.filter((r) => r.status === 'update_available'))
const other = computed(() => results.value.filter((r) => r.status !== 'update_available'))
const selectedNames = computed(() => [...selected.value])

/** 会被自动更新的那几台（预览）。 */
const willUpdate = computed(() => (auto.value?.candidates ?? []).filter((c) => c.willUpdate))
/** 被跳过的（含排除列表、自身、以及未检测到更新的）。 */
const skipped = computed(() => (auto.value?.candidates ?? []).filter((c) => !c.willUpdate && (c.protected || c.excluded)))

function toggle(name: string) {
  const s = new Set(selected.value)
  if (s.has(name)) s.delete(name)
  else s.add(name)
  selected.value = s
}

function selectAllAvailable() {
  if (selected.value.size === available.value.length) selected.value = new Set()
  else selected.value = new Set(available.value.map((r) => r.container))
}

async function load() {
  try {
    const res = await api.get<UpdatesResponse>('/api/updates')
    results.value = res.results ?? []
    checkedAt.value = res.checkedAt
    selfName.value = res.selfName
  } catch (e) {
    toast.error('读取巡检结果失败', e instanceof Error ? e.message : String(e))
  }
}

async function check(deep: boolean) {
  if (deep) deepLoading.value = true
  else loading.value = true
  try {
    const res = await api.post<{ results: CheckResult[] }>('/api/updates/check', { deep })
    results.value = res.results ?? []
    checkedAt.value = new Date().toISOString()
    const n = (res.results ?? []).filter((r) => r.status === 'update_available').length
    toast.success(deep ? '深度检测完成' : '巡检完成', `发现 ${n} 个有可用更新`)
  } catch (e) {
    toast.error('巡检失败', e instanceof Error ? e.message : String(e))
  } finally {
    deepLoading.value = false
    loading.value = false
  }
}

async function deepCheckOne(name: string) {
  try {
    await api.post('/api/updates/deep-check', { name })
    await load()
    toast.success(`${name} 深度检测完成`)
  } catch (e) {
    toast.error('深度检测失败', e instanceof Error ? e.message : String(e))
  }
}

async function apply() {
  const names = selectedNames.value
  if (!names.length) return
  applying.value = true
  progress.value = []
  try {
    await api.post('/api/updates/apply', { names, force: forceUpdate.value })
    showConfirm.value = false
    toast.info(`已提交 ${names.length} 个容器的更新任务`, '正在后台执行，进度见下方')
  } catch (e) {
    toast.error('提交失败', e instanceof Error ? e.message : String(e))
    applying.value = false
  }
}

const toneOf = (status: string) => checkLabel(status).tone

// ---------- 自动更新 ----------

async function loadAuto() {
  try {
    auto.value = await api.get<AutoUpdateInfo>('/api/updates/auto')
  } catch {
    // 自动更新信息读不到不影响主流程（例如守护进程暂时不可达）
    auto.value = null
  }
}

async function loadPolicy() {
  try {
    policy.value = await api.get<Settings>('/api/settings')
  } catch {
    policy.value = null
  }
}

async function savePolicy() {
  if (!policy.value) return
  savingPolicy.value = true
  try {
    policy.value = await api.put<Settings>('/api/settings', policy.value)
    await loadAuto()
    toast.success('更新策略已保存', '下一轮自动更新就会按新策略执行')
  } catch (e) {
    toast.error('保存失败', e instanceof Error ? e.message : String(e))
  } finally {
    savingPolicy.value = false
  }
}

/** 切换自动更新总开关（立即生效，不必再点保存）。 */
async function toggleAuto(on: boolean) {
  if (!policy.value) return
  policy.value.autoApply = on
  await savePolicy()
}

/** 立即跑一轮；dryRun 为真时只巡检不动容器。 */
async function runAutoCycle(dryRun: boolean) {
  autoBusy.value = true
  try {
    await api.post('/api/updates/auto-run', { dryRun })
    toast.info(dryRun ? '已开始巡检（不动容器）' : '已开始自动更新', '完成后本页会自动刷新')
    // 后端在后台跑，轮询几次把结果拉回来
    for (let i = 0; i < 40; i++) {
      await new Promise((r) => window.setTimeout(r, 1500))
      const info = await api.get<AutoUpdateInfo>('/api/updates/auto').catch(() => null)
      if (info) auto.value = info
      if (info && !info.running) break
    }
    await Promise.all([load(), loadAuto()])
    const last = auto.value?.lastRun
    if (last && !last.dryRun) {
      toast.success('自动更新完成', `更新 ${last.updated} 台，失败 ${last.failed} 台`)
    }
  } catch (e) {
    toast.error('启动失败', e instanceof Error ? e.message : String(e))
  } finally {
    autoBusy.value = false
  }
}

onMounted(() => {
  void load()
  void loadAuto()
  void loadPolicy()
  closeStream = openStream('/api/events/stream', (topic, ev) => {
    if (topic !== 'update') return
    const d = ev.data ?? {}
    const name = String(d.container ?? '')
    if (ev.kind === 'step' || ev.kind === 'container_status' || ev.kind === 'pull_progress') {
      const msg = String(d.message ?? d.status ?? '')
      if (!msg) return
      const idx = progress.value.findIndex((p) => p.name === name)
      const entry = { name, message: msg, status: String(ev.status ?? 'running') }
      if (idx >= 0) progress.value[idx] = entry
      else progress.value.unshift(entry)
      if (progress.value.length > 12) progress.value.pop()
    }
    if (ev.kind === 'container_status') {
      const st = String(d.status ?? '')
      if (['updated', 'failed', 'broken', 'up_to_date'].includes(st)) {
        const entry = progress.value.find((p) => p.name === name)
        if (entry) entry.status = st === 'updated' || st === 'up_to_date' ? 'success' : 'failed'
      }
    }
    if (ev.kind === 'batch_done') {
      applying.value = false
      selected.value = new Set()
      void load()
      const msg = `更新 ${d.updated ?? 0} 个，已是最新/跳过 ${d.skipped ?? 0} 个，失败 ${d.failed ?? 0} 个`
      toast.success('批量更新完成', msg)
    }
    // 自动更新：巡检结束 / 整轮结束都要把这块刷新一下
    if (ev.kind === 'auto_check_done' || ev.kind === 'auto_done') {
      void loadAuto()
      if (ev.kind === 'auto_done') void load()
    }
  })
})

onUnmounted(() => closeStream?.())
</script>

<template>
  <div class="flex flex-col gap-3.5 p-[18px]">
    <div class="dh-phead">
      <div class="dh-h1">更新中心</div>
      <div class="dh-sub">
        {{ checkedAt ? `上次检测 ${relativeTime(checkedAt)} · 共比对 ${results.length} 个容器` : '还没有检测过' }}
      </div>
      <div class="ml-auto flex gap-2">
        <button class="dh-btn" :disabled="loading || deepLoading" @click="check(true)">
          <Zap class="h-3.5 w-3.5" :class="deepLoading ? 'dh-spin' : ''" />深度检测
        </button>
        <button class="dh-btn" :disabled="loading || deepLoading" @click="check(false)">
          <RefreshCw class="h-3.5 w-3.5" :class="loading ? 'dh-spin' : ''" />重新检测
        </button>
        <button
          class="dh-btn dh-btn-primary"
          :disabled="!available.length || applying"
          @click="showConfirm = true"
        >
          <Download class="h-3.5 w-3.5" />更新 {{ available.length }} 个镜像
        </button>
      </div>
    </div>

    <!-- 说明条：把「为什么不会误停容器」讲清楚 -->
    <div class="dh-banner dh-banner-info items-start">
      <Info class="mt-[2px] h-4 w-4 flex-none" />
      <div class="min-w-0 flex-1 leading-relaxed">
        更新流程是「<b>先拉取、再比对镜像 ID</b>」：镜像一个字节没变就直接结束，
        <b>容器不会被停止、不会被重建</b>。检测也走 Docker 守护进程自己解析的仓库端点，
        因此「检测到的新版本」与「真正拉到的镜像」永远同源。
        <button class="ml-1 underline decoration-dotted" @click="showInfo = true">了解详情</button>
      </div>
    </div>

    <!-- 汇总 -->
    <div class="grid grid-cols-2 gap-3 lg:grid-cols-4">
      <div class="dh-card p-3.5">
        <div class="flex items-center gap-2 text-[12px] text-text-4"><Download class="h-3.5 w-3.5" />有可用更新</div>
        <div class="mt-1.5 text-[22px] font-semibold leading-none" :class="available.length ? 'text-[#fbbf24]' : ''">
          {{ available.length }}
        </div>
      </div>
      <div class="dh-card p-3.5">
        <div class="flex items-center gap-2 text-[12px] text-text-4"><CheckCircle2 class="h-3.5 w-3.5" />已是最新</div>
        <div class="mt-1.5 text-[22px] font-semibold leading-none">
          {{ results.filter((r) => r.status === 'up_to_date').length }}
        </div>
      </div>
      <div class="dh-card p-3.5">
        <div class="flex items-center gap-2 text-[12px] text-text-4"><ShieldQuestion class="h-3.5 w-3.5" />无法判定</div>
        <div class="mt-1.5 text-[22px] font-semibold leading-none">
          {{ results.filter((r) => r.status === 'unknown').length }}
        </div>
        <div class="mt-1 text-[11px] text-text-6">网络/认证问题导致，绝不当作有更新</div>
      </div>
      <div class="dh-card p-3.5">
        <div class="flex items-center gap-2 text-[12px] text-text-4"><HelpCircle class="h-3.5 w-3.5" />本地镜像</div>
        <div class="mt-1.5 text-[22px] font-semibold leading-none">
          {{ results.filter((r) => r.status === 'no_upstream').length }}
        </div>
        <div class="mt-1 text-[11px] text-text-6">本地构建，无远端可比对</div>
      </div>
    </div>

    <!-- 进度 -->
    <div v-if="progress.length" class="dh-card">
      <div class="dh-card-head">
        <Loader2 v-if="applying" class="h-3.5 w-3.5 dh-spin text-accent" />
        <Zap v-else class="h-3.5 w-3.5 text-text-4" />
        <span>{{ applying ? '更新进行中' : '最近一次更新进度' }}</span>
        <span class="ml-auto text-[11.5px] font-normal text-text-5">{{ progress.length }} 条</span>
      </div>
      <div class="flex flex-col">
        <div
          v-for="(p, i) in progress"
          :key="i"
          class="flex items-center gap-2.5 border-b border-[#171f2a] px-3.5 py-2 text-[12px] last:border-b-0"
        >
          <span
            class="h-[6px] w-[6px] flex-none rounded-full"
            :class="p.status === 'failed' ? 'bg-err' : p.status === 'success' ? 'bg-run' : 'bg-accent'"
          />
          <span class="w-[170px] flex-none truncate font-mono text-[11.5px] text-text-3">{{ p.name }}</span>
          <span class="min-w-0 flex-1 truncate text-text-2">{{ p.message }}</span>
        </div>
      </div>
    </div>

    <!-- 待更新 -->
    <div class="dh-card">
      <div class="dh-card-head">
        <Download class="h-3.5 w-3.5 text-text-4" />
        <span>有可用更新的容器</span>
        <span class="text-[11.5px] font-normal text-text-5">
          {{ checkedAt ? relativeTime(checkedAt) + '检查' : '尚未巡检' }}
        </span>
        <div v-if="available.length" class="ml-auto flex items-center gap-2">
          <label class="flex cursor-pointer items-center gap-2 text-[11.5px] font-normal text-text-4">
            <input type="checkbox" class="h-[13px] w-[13px] accent-[#2dd4bf]" :checked="selected.size === available.length && available.length > 0" @change="selectAllAvailable" />
            全选
          </label>
          <button class="dh-btn dh-btn-sm dh-btn-primary" :disabled="!selected.size" @click="showConfirm = true">
            更新选中的 {{ selected.size }} 个
          </button>
        </div>
      </div>

      <EmptyState
        v-if="!available.length"
        :icon="CheckCircle2"
        :title="checkedAt ? '没有可用更新' : '还没有巡检过'"
        :description="checkedAt ? '所有可比对的容器都与仓库摘要一致。' : '点击右上角「巡检」只读检查一遍（不会动任何容器）。'"
      />
      <div v-else class="overflow-x-auto">
        <table class="dh-table">
          <thead>
            <tr>
              <th class="w-[34px]" />
              <th>容器</th>
              <th>镜像</th>
              <th class="w-[150px]">本地摘要</th>
              <th class="w-[150px]">仓库摘要</th>
              <th class="w-[110px]">判定</th>
              <th class="w-[110px]">操作</th>
            </tr>
          </thead>
          <tbody>
            <tr v-for="r in available" :key="r.container">
              <td>
                <input
                  type="checkbox"
                  class="h-[14px] w-[14px] accent-[#2dd4bf]"
                  :checked="selected.has(r.container)"
                  @change="toggle(r.container)"
                />
              </td>
              <td>
                <RouterLink :to="`/containers/${encodeURIComponent(r.container)}`" class="text-[12.5px] font-medium hover:text-accent">
                  {{ r.container }}
                </RouterLink>
              </td>
              <td class="max-w-[220px] truncate font-mono text-[11.5px] text-text-3">{{ shortImage(r.image) }}</td>
              <td class="font-mono text-[11px] text-text-5">{{ (r.localDigest || '—').slice(0, 19) }}</td>
              <td class="font-mono text-[11px] text-text-5">{{ (r.remoteDigest || '—').slice(0, 19) }}</td>
              <td><span class="dh-badge dh-badge-warn">有新版本</span></td>
              <td>
                <button class="dh-btn dh-btn-sm" @click="deepCheckOne(r.container)" title="真的拉一次镜像来确认（权威判定）">
                  <Zap class="h-3 w-3" />确认真实性
                </button>
              </td>
            </tr>
          </tbody>
        </table>
      </div>
    </div>

    <!-- 自动更新 -->
    <div class="dh-card">
      <div class="dh-card-head">
        <Zap class="h-3.5 w-3.5" :class="auto?.enabled ? 'text-[#fbbf24]' : 'text-text-4'" />
        <span>自动更新</span>
        <span class="ml-2 text-[11.5px] font-normal" :class="auto?.enabled ? 'text-[#fcd34d]' : 'text-text-5'">
          {{ auto?.enabled ? '检测到新版本会自动更新' : '只检测，不会自动动容器' }}
        </span>
        <div class="ml-auto flex items-center gap-2">
          <button class="dh-btn dh-btn-sm" :disabled="autoBusy" @click="runAutoCycle(true)">
            <RefreshCw class="h-3 w-3" :class="autoBusy ? 'dh-spin' : ''" />立即巡检一轮
          </button>
          <button class="dh-btn dh-btn-sm dh-btn-primary" :disabled="autoBusy || !auto?.enabled" @click="runAutoCycle(false)">
            <Play class="h-3 w-3" />立即执行自动更新
          </button>
        </div>
      </div>

      <div class="flex flex-col gap-3 p-3.5">
        <!-- 总开关 + 下次巡检时间 -->
        <div class="grid grid-cols-1 gap-3 lg:grid-cols-3">
          <div class="rounded-[10px] border border-line-1 bg-ink-800 p-3 lg:col-span-2">
            <SettingRow title="检测到新版本后自动更新" sub="每轮巡检结束后，把有更新的容器（排除列表与自己除外）自动重建到新镜像">
              <ToggleSwitch
                v-if="policy"
                :model-value="policy.autoApply"
                :disabled="savingPolicy"
                label="自动更新"
                @update:model-value="toggleAuto"
              />
              <span v-else class="text-[11.5px] text-text-5">…</span>
            </SettingRow>
          </div>
          <div class="rounded-[10px] border border-line-1 bg-ink-800 p-3">
            <div class="flex items-center gap-1.5 text-[12px] text-text-4"><Clock class="h-3.5 w-3.5" />下次巡检</div>
            <div class="mt-1.5 text-[14px] font-semibold text-text-1">
              {{ auto?.nextCheckAt ? formatDateTime(auto.nextCheckAt) : '未开启周期巡检' }}
            </div>
            <div class="mt-1 text-[11px] text-text-5">
              {{ auto?.lastCheckAt ? `上次 ${relativeTime(auto.lastCheckAt)}` : '还没有巡检记录' }}
              · 每 {{ auto?.checkIntervalHours ?? 0 }} 小时一次
            </div>
          </div>
        </div>

        <!-- 本轮会发生什么 -->
        <div class="rounded-[10px] border border-line-1 bg-ink-800">
          <div class="flex flex-wrap items-center gap-2 border-b border-line-1 px-3 py-2">
            <span class="text-[12px] font-medium text-text-3">本轮会发生什么</span>
            <span class="dh-badge dh-badge-warn">{{ willUpdate.length }} 台将被更新</span>
            <span class="dh-badge dh-badge-plain">{{ skipped.length }} 台被保护/排除</span>
            <RouterLink to="/settings" class="ml-auto text-[11.5px] text-text-5 hover:text-accent">
              去设置里调整检测频率与排除列表 →
            </RouterLink>
          </div>
          <div v-if="!auto?.candidates?.length" class="px-3 py-3 text-[12px] text-text-5">
            还没有巡检结果，点右上角「立即巡检一轮」先跑一次。
          </div>
          <div v-else class="overflow-x-auto">
            <table class="dh-table">
              <thead>
                <tr>
                  <th class="w-[200px]">容器</th>
                  <th class="w-[220px]">镜像</th>
                  <th class="w-[110px]">状态</th>
                  <th>去向</th>
                </tr>
              </thead>
              <tbody>
                <tr v-for="c in auto.candidates" :key="c.name">
                  <td class="font-mono text-[11.5px] text-text-3">{{ c.name }}</td>
                  <td class="max-w-[220px] truncate font-mono text-[11px] text-text-5">{{ shortImage(c.image) }}</td>
                  <td>
                    <span class="dh-badge" :class="c.running ? 'dh-badge-run' : 'dh-badge-stop'">
                      {{ c.running ? '运行中' : '已停止' }}
                    </span>
                  </td>
                  <td class="text-[11.5px]">
                    <span v-if="c.willUpdate" class="dh-badge dh-badge-warn">将更新</span>
                    <span v-else-if="c.protected" class="dh-badge dh-badge-accent">受保护</span>
                    <span v-else-if="c.excluded" class="dh-badge dh-badge-plain">已排除</span>
                    <span v-else class="dh-badge dh-badge-plain">无需更新</span>
                    <span v-if="c.reason" class="ml-2 text-text-5">{{ c.reason }}</span>
                  </td>
                </tr>
              </tbody>
            </table>
          </div>
        </div>

        <!-- 上一轮结果 -->
        <div v-if="auto?.lastRun" class="rounded-[10px] border border-line-1 bg-ink-800 px-3 py-2.5">
          <div class="flex flex-wrap items-center gap-2 text-[12px]">
            <span class="font-medium text-text-3">上一轮（{{ auto.lastRun.trigger === 'manual' ? '手动触发' : '定时触发' }}）</span>
            <span class="dh-badge dh-badge-plain">{{ relativeTime(auto.lastRun.startedAt) }}</span>
            <span v-if="auto.lastRun.dryRun" class="dh-badge dh-badge-plain">仅巡检</span>
            <span v-else class="dh-badge dh-badge-run">更新 {{ auto.lastRun.updated }} 台</span>
            <span v-if="auto.lastRun.failed" class="dh-badge dh-badge-err">失败 {{ auto.lastRun.failed }} 台</span>
            <span class="ml-auto text-[11.5px] text-text-5">
              巡检 {{ auto.lastRun.checked }} 台 · 发现 {{ auto.lastRun.available }} 台有更新 · 耗时 {{ (auto.lastRun.durationMs / 1000).toFixed(1) }}s
            </span>
          </div>
        </div>

        <div
          v-if="auto && !auto.enabled"
          class="flex items-start gap-2 rounded-[8px] border border-line-1 px-2.5 py-2 text-[11.5px] leading-relaxed text-text-4"
        >
          <Info class="mt-[1px] h-3.5 w-3.5 flex-none" />
          <span>
            自动更新默认关闭：检测到新版本只会打上「有新版本」标记，要不要更新由你决定。
            打开开关后，每轮巡检结束就会按上面的清单自动重建容器 —— 开启前建议先看清清单。
          </span>
        </div>
      </div>
    </div>

    <!-- 更新策略 -->
    <div v-if="policy" class="dh-card">
      <div class="dh-card-head">
        <ShieldAlert class="h-3.5 w-3.5 text-text-4" />
        <span>更新策略</span>
        <button class="dh-btn dh-btn-sm dh-btn-primary ml-auto" :disabled="savingPolicy" @click="savePolicy">
          <Loader2 v-if="savingPolicy" class="h-3 w-3 dh-spin" />
          <Download v-else class="h-3 w-3" />保存策略
        </button>
      </div>
      <div class="flex flex-col gap-3 p-3.5">
        <SettingRow title="仅在镜像真正变化时重启容器" sub="比对更新前后的镜像 ID，一致就完全不动它">
          <span class="dh-badge dh-badge-plain">始终开启</span>
        </SettingRow>
        <SettingRow title="同一镜像的多个容器只拉取一次" sub="4 台共用 nginx:alpine 时，下载 1 次、重建 4 台">
          <ToggleSwitch v-model="policy.pullOnce" label="同一镜像只拉取一次" />
        </SettingRow>
        <SettingRow title="更新前自动备份容器配置" sub="失败可一键回滚到更新前">
          <ToggleSwitch v-model="policy.backupBefore" label="更新前自动备份容器配置" />
        </SettingRow>
        <SettingRow title="更新后清理旧镜像" sub="确认无引用后才删除">
          <ToggleSwitch v-model="policy.cleanupAfter" label="更新后清理旧镜像" />
        </SettingRow>
      </div>
    </div>

    <!-- 其它容器 -->
    <div v-if="other.length" class="dh-card">
      <div class="dh-card-head">
        <CheckCircle2 class="h-3.5 w-3.5 text-text-4" />
        <span>其余容器</span>
        <span class="ml-auto text-[11.5px] font-normal text-text-5">{{ other.length }} 个</span>
      </div>
      <div class="overflow-x-auto">
        <table class="dh-table">
          <thead>
            <tr>
              <th>容器</th>
              <th>镜像</th>
              <th class="w-[110px]">判定</th>
              <th>说明</th>
            </tr>
          </thead>
          <tbody>
            <tr v-for="r in other" :key="r.container">
              <td>
                <RouterLink :to="`/containers/${encodeURIComponent(r.container)}`" class="text-[12.5px] font-medium hover:text-accent">
                  {{ r.container }}
                </RouterLink>
                <span v-if="r.container === selfName" class="ml-1.5 dh-badge dh-badge-accent">自身</span>
              </td>
              <td class="max-w-[240px] truncate font-mono text-[11.5px] text-text-3">{{ shortImage(r.image) }}</td>
              <td>
                <span class="dh-badge" :class="`dh-badge-${toneOf(r.status)}`">{{ checkLabel(r.status).text }}</span>
              </td>
              <td class="text-[11.5px] text-text-4">{{ r.reason }}</td>
            </tr>
          </tbody>
        </table>
      </div>
    </div>

    <!-- 确认 -->
    <Modal
      :open="showConfirm"
      title="确认执行更新"
      :subtitle="`共 ${selected.size} 个容器`"
      @close="showConfirm = false"
    >
      <div class="flex flex-col gap-3">
        <div class="rounded-[10px] border border-line-1 bg-ink-800 px-3 py-2.5 text-[12px] leading-relaxed text-text-3">
          Docker 会按顺序对每个容器执行：拉取镜像 → 比对镜像 ID。
          <b class="text-text-1">ID 没变就完全跳过</b>，容器不会被停止或重建。
          只有镜像真的变化时才会走「停旧 → 改名保留 → 建新 → 健康检查（失败自动回滚）」。
        </div>
        <div class="dh-scroll max-h-[200px] overflow-auto rounded-[10px] border border-line-1 bg-ink-800 p-2.5">
          <div v-for="n in selectedNames" :key="n" class="px-1 py-[3px] font-mono text-[11.5px] text-text-3">
            {{ n }}
          </div>
        </div>
        <label class="flex cursor-pointer items-start gap-2.5 text-[12px] text-text-2">
          <input v-model="forceUpdate" type="checkbox" class="mt-[3px] h-[14px] w-[14px] accent-[#f5a524]" />
          <span>
            强制重建（即使镜像 ID 未变化也重建容器）
            <br />
            <span class="text-[11.5px] text-text-5">
              一般不需要。勾选后连「已是最新」的容器也会被停掉重建 —— 这正是 dockerCopilot 曾经的行为。
            </span>
          </span>
        </label>
      </div>
      <template #footer>
        <button class="dh-btn" @click="showConfirm = false">取消</button>
        <button class="dh-btn dh-btn-primary" :disabled="applying" @click="apply">
          <Loader2 v-if="applying" class="h-3.5 w-3.5 dh-spin" />
          <Download v-else class="h-3.5 w-3.5" />
          开始更新
        </button>
      </template>
    </Modal>

    <!-- 原理说明 -->
    <Modal :open="showInfo" title="为什么 Dockhelm 不会误停容器" width="620px" @close="showInfo = false">
      <div class="flex flex-col gap-3 text-[12.5px] leading-relaxed text-text-3">
        <div class="flex items-start gap-2.5">
          <AlertTriangle class="mt-[2px] h-4 w-4 flex-none text-[#fbbf24]" />
          <div>
            <b class="text-text-1">dockerCopilot 的两个缺陷。</b>
            一是它的检测路径自拼 registry 请求、用一份硬编码的加速站列表，从不读取守护进程
            <code class="text-text-2">daemon.json</code> 里真正生效的 <code class="text-text-2">registry-mirrors</code>，
            于是「检测说有新版本、拉取说已是最新」会永久互相矛盾。
            二是它的更新流程无条件执行 <code class="text-text-2">pull → stop → rename → create → start</code>，
            <b class="text-text-1">算出了新镜像 ID 却只用来决定要不要删旧镜像，从不用于决定要不要停容器</b>。
            两者叠加，就出现了「一次更新十几台，连没更新的容器也全被停掉」。
          </div>
        </div>
        <div class="flex items-start gap-2.5">
          <CheckCircle2 class="mt-[2px] h-4 w-4 flex-none text-[#4ade80]" />
          <div>
            <b class="text-text-1">Dockhelm 的做法。</b>
            检测走守护进程的 <code class="text-text-2">/distribution</code> 接口（与 pull 同一套仓库端点解析），
            检测失败一律标记为「无法判定」而不是「有新版本」。
            更新时先拉取，再比对容器使用的镜像 ID：<b class="text-text-1">一致就直接返回，容器一个字节都不碰</b>；
            不一致才停止旧容器、改名保留（<code class="text-text-2">&lt;名字&gt;__bak_&lt;时间&gt;</code>）、
            用原配置创建同名新容器、启动并做健康检查，失败自动把旧容器改回原名并启动。
          </div>
        </div>
        <div class="flex items-start gap-2.5">
          <Info class="mt-[2px] h-4 w-4 flex-none text-accent" />
          <div>
            <b class="text-text-1">「深度检测」是什么。</b>
            它真的拉一次镜像再比对镜像 ID —— 这是唯一 100% 同源的判定。
            镜像已最新时守护进程只会下载 manifest 与 config（几 KB），不会下载层文件。
            如果你怀疑某个容器被误报，用它确认即可。
          </div>
        </div>
      </div>
      <template #footer>
        <button class="dh-btn dh-btn-primary" @click="showInfo = false">明白了</button>
      </template>
    </Modal>
  </div>
</template>
