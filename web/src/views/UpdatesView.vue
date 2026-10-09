<script setup lang="ts">
import { computed, onMounted, onUnmounted, ref } from 'vue'
import { RouterLink } from 'vue-router'
import {
  AlertTriangle,
  CheckCircle2,
  Download,
  HelpCircle,
  Info,
  Loader2,
  RefreshCw,
  ShieldQuestion,
  Zap,
} from 'lucide-vue-next'
import { api, openStream } from '@/api/client'
import type { CheckResult, UpdatesResponse } from '@/api/types'
import { checkLabel, relativeTime, shortImage } from '@/utils/format'
import { useToastStore } from '@/stores/toast'
import EmptyState from '@/components/EmptyState.vue'
import Modal from '@/components/Modal.vue'

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

const available = computed(() => results.value.filter((r) => r.status === 'update_available'))
const other = computed(() => results.value.filter((r) => r.status !== 'update_available'))
const selectedNames = computed(() => [...selected.value])

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

onMounted(() => {
  void load()
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
