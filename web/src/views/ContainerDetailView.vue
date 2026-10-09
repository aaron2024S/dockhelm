<script setup lang="ts">
import { computed, onMounted, onUnmounted, ref } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import {
  ArrowLeft,
  Box,
  Download,
  Eye,
  EyeOff,
  HardDrive,
  Network,
  Play,
  RefreshCw,
  RotateCw,
  ScrollText,
  Square,
  Terminal,
} from 'lucide-vue-next'
import { api, openStream } from '@/api/client'
import { relativeTime } from '@/utils/format'
import { useToastStore } from '@/stores/toast'

const route = useRoute()
const router = useRouter()
const toast = useToastStore()

const name = computed(() => decodeURIComponent(String(route.params.name ?? '')))
const detail = ref<Record<string, any> | null>(null)
const summary = ref<Record<string, any> | null>(null)
const stats = ref<Record<string, number> | null>(null)
const logs = ref('')
const tab = ref<'logs' | 'config' | 'mounts' | 'networks'>('logs')
const loading = ref(true)
const busy = ref(false)
const showSensitive = ref(false)
const tail = ref(200)
let timer: number | undefined
let closeStream: (() => void) | null = null

const container = computed(() => summary.value ?? {})
const envList = computed<{ key: string; value: string; sensitive?: string }[]>(() => {
  const list = (summary.value?.env as { key: string; value: string; sensitive?: string }[]) ?? []
  return list
})

const mounts = computed<Record<string, any>[]>(() => (summary.value?.mounts as Record<string, any>[]) ?? [])
const networks = computed<Record<string, any>[]>(() => (summary.value?.networks as Record<string, any>[]) ?? [])

async function load() {
  loading.value = true
  try {
    const res = await api.get<{ inspect: Record<string, unknown>; summary: Record<string, unknown> }>(
      `/api/containers/${encodeURIComponent(name.value)}`,
    )
    detail.value = res.inspect
    summary.value = res.summary
  } catch (e) {
    toast.error('读取容器失败', e instanceof Error ? e.message : String(e))
  } finally {
    loading.value = false
  }
}

async function loadLogs() {
  try {
    logs.value = await api.text(`/api/containers/${encodeURIComponent(name.value)}/logs`, { tail: tail.value })
  } catch (e) {
    logs.value = '读取日志失败：' + (e instanceof Error ? e.message : String(e))
  }
}

async function loadStats() {
  try {
    const res = await api.get<{ summary: Record<string, number> }>(
      `/api/containers/${encodeURIComponent(name.value)}/stats`,
    )
    stats.value = res.summary
  } catch {
    stats.value = null
  }
}

async function act(action: string) {
  busy.value = true
  try {
    await api.post(`/api/containers/${encodeURIComponent(name.value)}/action`, { action })
    toast.success('操作已执行')
    await load()
    await loadStats()
  } catch (e) {
    toast.error('操作失败', e instanceof Error ? e.message : String(e))
  } finally {
    busy.value = false
  }
}

async function snapshot() {
  try {
    await api.post('/api/backups/snapshot', { container: name.value, reason: 'manual' })
    toast.success('已保存配置快照')
  } catch (e) {
    toast.error('备份失败', e instanceof Error ? e.message : String(e))
  }
}

async function updateNow() {
  try {
    await api.post('/api/updates/apply', { names: [name.value] })
    toast.info('已提交更新任务', '镜像没变化的话容器不会被停止')
  } catch (e) {
    toast.error('提交失败', e instanceof Error ? e.message : String(e))
  }
}

const healthTone = computed(() => {
  const h = container.value.health as string
  if (h === 'healthy') return 'dh-badge-run'
  if (h === 'unhealthy') return 'dh-badge-err'
  return 'dh-badge-plain'
})

onMounted(async () => {
  await load()
  await loadLogs()
  await loadStats()
  timer = window.setInterval(() => {
    if (tab.value === 'logs') void loadLogs()
    void loadStats()
  }, 8000)
  closeStream = openStream('/api/events/stream', (topic, ev) => {
    if (topic === 'container' && String(ev.data?.name ?? '') === name.value) void load()
  })
})

onUnmounted(() => {
  if (timer) window.clearInterval(timer)
  closeStream?.()
})

const running = computed(() => container.value.health !== undefined && (detail.value?.['State'] as any)?.Running)
</script>

<template>
  <div class="flex flex-col gap-3.5 p-[18px]">
    <div class="dh-phead">
      <button class="dh-iconbtn" title="返回列表" @click="router.back()">
        <ArrowLeft class="h-4 w-4" />
      </button>
      <div class="grid h-[30px] w-[30px] flex-none place-items-center rounded-[9px] bg-line-2 text-[12px] font-semibold text-[#5eead4]">
        {{ name.slice(0, 2).toUpperCase() }}
      </div>
      <div class="min-w-0">
        <div class="dh-h1 truncate">{{ name }}</div>
        <div class="dh-sub truncate font-mono">{{ container.id }}</div>
      </div>

      <div class="ml-auto flex flex-wrap items-center gap-2">
        <button v-if="running" class="dh-btn dh-btn-sm" :disabled="busy" @click="act('stop')">
          <Square class="h-3 w-3" />停止
        </button>
        <button v-else class="dh-btn dh-btn-sm" :disabled="busy" @click="act('start')">
          <Play class="h-3 w-3" />启动
        </button>
        <button class="dh-btn dh-btn-sm" :disabled="busy" @click="act('restart')">
          <RotateCw class="h-3 w-3" :class="busy ? 'dh-spin' : ''" />重启
        </button>
        <button class="dh-btn dh-btn-sm" @click="updateNow">
          <Download class="h-3 w-3" />检查并更新
        </button>
        <button class="dh-btn dh-btn-sm" @click="snapshot">
          <HardDrive class="h-3 w-3" />备份配置
        </button>
      </div>
    </div>

    <!-- 概览条 -->
    <div class="grid grid-cols-2 gap-3 lg:grid-cols-4">
      <div class="dh-card p-3.5">
        <div class="text-[12px] text-text-4">状态</div>
        <div class="mt-1.5 flex flex-wrap items-center gap-1.5">
          <span
            class="dh-badge"
            :class="detail?.['State'] && (detail['State'] as any).Running ? 'dh-badge-run' : 'dh-badge-stop'"
          >
            {{ detail?.['State'] && (detail['State'] as any).Running ? '运行中' : '已停止' }}
          </span>
          <span v-if="container.health" class="dh-badge" :class="healthTone">{{ container.health }}</span>
        </div>
        <div class="mt-1.5 text-[11px] text-text-5">
          退出码 {{ container.exitCode ?? '—' }} · 重启 {{ container.restarts ?? 0 }} 次
        </div>
      </div>
      <div class="dh-card p-3.5">
        <div class="text-[12px] text-text-4">CPU</div>
        <div class="mt-1.5 text-[20px] font-semibold leading-none">
          {{ stats ? (stats.cpuPercent ?? 0).toFixed(1) : '—' }}<span class="ml-0.5 text-[12px] text-text-5">%</span>
        </div>
        <div class="mt-1.5 text-[11px] text-text-5">{{ stats?.onlineCpus ?? '—' }} 核可用</div>
      </div>
      <div class="dh-card p-3.5">
        <div class="text-[12px] text-text-4">内存</div>
        <div class="mt-1.5 text-[20px] font-semibold leading-none">
          {{ stats ? (stats.memPercent ?? 0).toFixed(1) : '—' }}<span class="ml-0.5 text-[12px] text-text-5">%</span>
        </div>
        <div class="mt-1.5 text-[11px] text-text-5">
          {{ stats ? Math.round((stats.memUsage ?? 0) / 1048576) + ' MB' : '—' }} /
          {{ stats ? Math.round((stats.memLimit ?? 0) / 1048576) + ' MB' : '—' }}
        </div>
      </div>
      <div class="dh-card p-3.5">
        <div class="text-[12px] text-text-4">网络</div>
        <div class="mt-1.5 text-[16px] font-semibold leading-tight">
          ↓ {{ stats ? Math.round((stats.netRx ?? 0) / 1048576) + ' MB' : '—' }}
        </div>
        <div class="mt-[3px] text-[16px] font-semibold leading-tight text-text-3">
          ↑ {{ stats ? Math.round((stats.netTx ?? 0) / 1048576) + ' MB' : '—' }}
        </div>
      </div>
    </div>

    <!-- 详情 -->
    <div class="dh-card">
      <div class="dh-card-head">
        <div class="dh-seg">
          <button :data-on="tab === 'logs'" @click="tab = 'logs'">日志</button>
          <button :data-on="tab === 'config'" @click="tab = 'config'">配置</button>
          <button :data-on="tab === 'mounts'" @click="tab = 'mounts'">挂载 {{ mounts.length }}</button>
          <button :data-on="tab === 'networks'" @click="tab = 'networks'">网络 {{ networks.length }}</button>
        </div>
        <div v-if="tab === 'logs'" class="ml-auto flex items-center gap-2">
          <select v-model.number="tail" class="dh-select !w-[110px] !py-[5px] !text-[11.5px]" @change="loadLogs">
            <option :value="100">最近 100 行</option>
            <option :value="200">最近 200 行</option>
            <option :value="500">最近 500 行</option>
            <option :value="2000">最近 2000 行</option>
          </select>
          <button class="dh-btn dh-btn-sm" @click="loadLogs"><RefreshCw class="h-3 w-3" />刷新</button>
        </div>
      </div>

      <!-- 日志 -->
      <div v-if="tab === 'logs'">
        <div v-if="loading" class="grid h-[200px] place-items-center">
          <RefreshCw class="h-5 w-5 dh-spin text-text-5" />
        </div>
        <pre
          v-else
          class="dh-scroll max-h-[520px] overflow-auto whitespace-pre-wrap break-all px-3.5 py-3 font-mono text-[11.5px] leading-[1.65] text-text-2"
          >{{ logs || '（没有日志输出）' }}</pre
        >
      </div>

      <!-- 配置 -->
      <div v-else-if="tab === 'config'" class="grid grid-cols-1 gap-0 lg:grid-cols-2">
        <div class="border-b border-line-1 p-3.5 lg:border-b-0 lg:border-r">
          <div class="mb-2 text-[12px] font-medium text-text-3">基本信息</div>
          <div class="grid grid-cols-[92px_1fr] gap-x-3 gap-y-2 text-[12px]">
            <div class="text-text-5">镜像</div>
            <div class="break-all font-mono text-[11.5px] text-text-2">{{ container.image }}</div>
            <div class="text-text-5">镜像 ID</div>
            <div class="font-mono text-[11.5px] text-text-2">{{ container.imageId }}</div>
            <div class="text-text-5">启动命令</div>
            <div class="break-all font-mono text-[11.5px] text-text-2">
              {{ [...(container.entrypoint ?? []), ...(container.cmd ?? [])].join(' ') || '（镜像默认）' }}
            </div>
            <div class="text-text-5">工作目录</div>
            <div class="font-mono text-[11.5px] text-text-2">{{ container.workingDir || '—' }}</div>
            <div class="text-text-5">运行用户</div>
            <div class="font-mono text-[11.5px] text-text-2">{{ container.user || 'root' }}</div>
            <div class="text-text-5">重启策略</div>
            <div class="text-text-2">{{ container.restart }}</div>
            <div class="text-text-5">网络模式</div>
            <div class="text-text-2">{{ container.networkMode }}</div>
            <div class="text-text-5">特权模式</div>
            <div class="text-text-2">{{ container.privileged ? '是' : '否' }}</div>
            <div class="text-text-5">启动时间</div>
            <div class="text-text-2">{{ relativeTime(container.startedAt) }}</div>
          </div>
        </div>
        <div class="p-3.5">
          <div class="mb-2 flex items-center gap-2">
            <span class="text-[12px] font-medium text-text-3">环境变量</span>
            <button class="dh-btn dh-btn-ghost !p-1" @click="showSensitive = !showSensitive">
              <component :is="showSensitive ? EyeOff : Eye" class="h-3.5 w-3.5" />
            </button>
            <span class="text-[11px] text-text-6">
              {{ showSensitive ? '敏感值已显示' : '敏感值已打码' }}
            </span>
          </div>
          <div class="dh-scroll max-h-[360px] overflow-auto">
            <table class="w-full">
              <tbody>
                <tr v-for="e in envList" :key="e.key" class="border-b border-[#171f2a] last:border-b-0">
                  <td class="py-1.5 pr-3 align-top font-mono text-[11px] text-text-4">{{ e.key }}</td>
                  <td class="break-all py-1.5 font-mono text-[11px] text-text-2">
                    <template v-if="e.sensitive === 'true' && !showSensitive">
                      <span class="text-text-6">••••••••</span>
                    </template>
                    <template v-else>{{ e.value }}</template>
                  </td>
                </tr>
              </tbody>
            </table>
          </div>
        </div>
      </div>

      <!-- 挂载 -->
      <div v-else-if="tab === 'mounts'" class="p-3.5">
        <div
          v-if="!mounts.length"
          class="text-[12.5px] text-text-4"
        >
          这个容器没有任何挂载。
        </div>
        <div v-else class="flex flex-col gap-2.5">
          <div class="flex items-start gap-2 rounded-[10px] border border-line-1 bg-ink-800 px-3 py-2 text-[11.5px] leading-relaxed text-text-4">
            <Box class="mt-[1px] h-3.5 w-3.5 flex-none" />
            <span>
              「绑定挂载」的数据在<b>宿主机目录</b>上，Dockhelm 容器默认看不见 —— 要备份这份数据，
              需要把对应宿主目录也挂进 Dockhelm（冒号右边叫什么名字都可以，启动时会自动识别）。「命名卷」的数据在
              <code class="text-text-3">/var/lib/docker/volumes</code> 下，只读挂载该目录即可备份。
            </span>
          </div>
          <div
            v-for="(m, i) in mounts"
            :key="i"
            class="rounded-[10px] border border-line-1 bg-ink-800 px-3 py-2.5"
          >
            <div class="flex flex-wrap items-center gap-2">
              <span class="dh-badge" :class="m.type === 'volume' ? 'dh-badge-accent' : 'dh-badge-plain'">
                {{ m.type === 'volume' ? '命名卷/匿名卷' : m.type }}
              </span>
              <span class="dh-badge" :class="m.rw ? 'dh-badge-run' : 'dh-badge-warn'">
                {{ m.rw ? '读写' : '只读' }}
              </span>
              <code class="ml-auto font-mono text-[11px] text-text-3">{{ m.destination }}</code>
            </div>
            <div class="mt-1.5 break-all font-mono text-[11px] text-text-2">
              {{ m.name || m.source }}
            </div>
            <div class="mt-1 text-[11px] text-text-6">{{ m.note }}</div>
          </div>
        </div>
      </div>

      <!-- 网络 -->
      <div v-else class="p-3.5">
        <div v-if="!networks.length" class="text-[12.5px] text-text-4">没有网络信息</div>
        <div v-else class="grid grid-cols-1 gap-2.5 md:grid-cols-2">
          <div
            v-for="n in networks"
            :key="n.name"
            class="rounded-[10px] border border-line-1 bg-ink-800 px-3 py-2.5"
          >
            <div class="flex items-center gap-2">
              <Network class="h-3.5 w-3.5 text-text-3" />
              <span class="text-[12.5px] font-medium">{{ n.name }}</span>
            </div>
            <div class="mt-1.5 grid grid-cols-[62px_1fr] gap-x-3 gap-y-1 text-[11.5px]">
              <div class="text-text-5">容器 IP</div>
              <div class="font-mono text-text-2">{{ n.ip || '—' }}</div>
              <div class="text-text-5">网关</div>
              <div class="font-mono text-text-2">{{ n.gateway || '—' }}</div>
              <div class="text-text-5">别名</div>
              <div class="break-all font-mono text-text-2">{{ (n.aliases ?? []).join(', ') || '—' }}</div>
            </div>
          </div>
        </div>
      </div>
    </div>

    <div class="flex items-center gap-2 text-[11.5px] text-text-6">
      <ScrollText class="h-3.5 w-3.5" />
      日志每 8 秒自动刷新一次
      <span class="ml-auto inline-flex items-center gap-1">
        <Terminal class="h-3.5 w-3.5" />容器名 <code class="text-text-4">{{ name }}</code>
      </span>
    </div>
  </div>
</template>
