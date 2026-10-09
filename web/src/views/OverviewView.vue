<script setup lang="ts">
import { computed, onMounted, onUnmounted, ref } from 'vue'
import { RouterLink } from 'vue-router'
import {
  Activity,
  AlertTriangle,
  ArrowUpRight,
  Box,
  CircleDot,
  Download,
  HardDrive,
  Layers,
  Loader2,
  RefreshCw,
  Server,
} from 'lucide-vue-next'
import { api, openStream } from '@/api/client'
import type { OverviewResponse, RunLog } from '@/api/types'
import { formatBytes, relativeTime, shortImage } from '@/utils/format'
import { useToastStore } from '@/stores/toast'
import EmptyState from '@/components/EmptyState.vue'

const toast = useToastStore()
const data = ref<OverviewResponse | null>(null)
const loading = ref(true)
const errorMsg = ref('')
const live = ref<{ text: string; status: string }[]>([])
let closeStream: (() => void) | null = null

async function load() {
  loading.value = true
  errorMsg.value = ''
  try {
    data.value = await api.get<OverviewResponse>('/api/overview')
  } catch (e) {
    errorMsg.value = e instanceof Error ? e.message : String(e)
  } finally {
    loading.value = false
  }
}

async function checkNow() {
  try {
    await api.post('/api/updates/check', {})
    toast.success('巡检已开始', '只读检查，不会停止任何容器')
    window.setTimeout(load, 2500)
  } catch (e) {
    toast.error('巡检失败', e instanceof Error ? e.message : String(e))
  }
}

const diskPercent = computed(() => {
  const d = data.value?.disk
  if (!d || !d.total) return 0
  return ((d.total - d.free) / d.total) * 100
})

const statusOfLog = (l: RunLog) => {
  if (l.status === 'success' || l.status === 'up_to_date') return 'accent'
  if (l.status === 'failed') return 'err'
  return 'plain'
}

onMounted(() => {
  void load()
  closeStream = openStream('/api/events/stream', (topic, ev) => {
    if (topic !== 'update' && topic !== 'schedule') return
    const message = typeof ev.data?.message === 'string' ? ev.data.message : ev.kind
    live.value.unshift({ text: String(message), status: String(ev.status ?? 'info') })
    if (live.value.length > 6) live.value.pop()
    // 批次或巡检结束时刷新总览
    if (ev.kind === 'batch_done' || ev.kind === 'check_done') {
      window.setTimeout(load, 800)
    }
  })
})

onUnmounted(() => closeStream?.())
</script>

<template>
  <div class="flex flex-col gap-3.5 p-[18px]">
    <!-- 顶部提示条 -->
    <div
      v-if="data?.updates.available"
      class="flex items-center gap-3 rounded-[14px] border border-[rgba(245,165,36,.3)] bg-[rgba(245,165,36,.07)] px-4 py-3 text-[13px] text-[#fcd34d]"
    >
      <AlertTriangle class="h-4 w-4 flex-none" />
      <div class="min-w-0 flex-1">
        有 <b>{{ data.updates.available }}</b> 个容器存在可用更新。Dockhelm 只会更新镜像真正变化的容器，
        拉取后若镜像 ID 未变会直接跳过，不会把它们停掉。
      </div>
      <RouterLink to="/updates" class="dh-btn dh-btn-sm flex-none !border-[rgba(245,165,36,.45)] !text-[#fcd34d]">
        去更新中心
        <ArrowUpRight class="h-3 w-3" />
      </RouterLink>
    </div>

    <div
      v-if="errorMsg"
      class="flex items-center gap-3 rounded-[14px] border border-[rgba(248,113,113,.35)] bg-[rgba(248,113,113,.07)] px-4 py-3 text-[13px] text-[#fca5a5]"
    >
      <AlertTriangle class="h-4 w-4 flex-none" />
      <span class="min-w-0 flex-1">{{ errorMsg }}</span>
      <button class="dh-btn dh-btn-sm" @click="load">重试</button>
    </div>

    <div
      v-if="data?.dockerError"
      class="flex items-center gap-3 rounded-[14px] border border-[rgba(248,113,113,.35)] bg-[rgba(248,113,113,.07)] px-4 py-3 text-[13px] text-[#fca5a5]"
    >
      <AlertTriangle class="h-4 w-4 flex-none" />
      <span class="min-w-0 flex-1">无法连接 Docker 守护进程：{{ data.dockerError }}</span>
    </div>

    <!-- 指标卡 -->
    <div class="grid grid-cols-2 gap-3 lg:grid-cols-4">
      <div class="dh-card p-3.5">
        <div class="flex items-center gap-2 text-[12px] text-text-4">
          <Box class="h-3.5 w-3.5" />容器
        </div>
        <div class="mt-1.5 flex items-baseline gap-1.5">
          <span class="text-[22px] font-semibold leading-none tracking-[-0.3px]">
            {{ data?.containers.running ?? '—' }}
          </span>
          <span class="text-[12px] text-text-5">/ {{ data?.containers.total ?? '—' }} 运行中</span>
        </div>
        <div class="mt-2 flex flex-wrap gap-1.5">
          <span class="dh-badge dh-badge-run"><span class="h-[7px] w-[7px] rounded-full bg-run" />运行 {{ data?.containers.running ?? 0 }}</span>
          <span class="dh-badge dh-badge-stop">停止 {{ data?.containers.stopped ?? 0 }}</span>
          <span v-if="data?.containers.unhealthy" class="dh-badge dh-badge-err">异常 {{ data.containers.unhealthy }}</span>
        </div>
      </div>

      <div class="dh-card p-3.5">
        <div class="flex items-center gap-2 text-[12px] text-text-4">
          <Download class="h-3.5 w-3.5" />可用更新
        </div>
        <div class="mt-1.5 flex items-baseline gap-1.5">
          <span
            class="text-[22px] font-semibold leading-none tracking-[-0.3px]"
            :class="data?.updates.available ? 'text-[#fbbf24]' : ''"
          >
            {{ data?.updates.available ?? 0 }}
          </span>
          <span class="text-[12px] text-text-5">
            {{ data?.updates.checkedAt ? relativeTime(data.updates.checkedAt) + '检查' : '尚未检查' }}
          </span>
        </div>
        <div class="mt-2 flex items-center gap-2">
          <button class="dh-btn dh-btn-sm" :disabled="loading" @click="checkNow">
            <RefreshCw class="h-3 w-3" />立即巡检
          </button>
          <span v-if="data?.updates.unknown" class="dh-badge dh-badge-plain">
            {{ data.updates.unknown }} 个无法判定
          </span>
        </div>
      </div>

      <div class="dh-card p-3.5">
        <div class="flex items-center gap-2 text-[12px] text-text-4">
          <Layers class="h-3.5 w-3.5" />镜像
        </div>
        <div class="mt-1.5 flex items-baseline gap-1.5">
          <span class="text-[22px] font-semibold leading-none tracking-[-0.3px]">{{ data?.images.total ?? '—' }}</span>
          <span class="text-[12px] text-text-5">共 {{ formatBytes(data?.images.sizeBytes) }}</span>
        </div>
        <div class="mt-2">
          <RouterLink to="/images" class="dh-btn dh-btn-sm">
            管理镜像
            <ArrowUpRight class="h-3 w-3" />
          </RouterLink>
        </div>
      </div>

      <div class="dh-card p-3.5">
        <div class="flex items-center gap-2 text-[12px] text-text-4">
          <HardDrive class="h-3.5 w-3.5" />数据盘
        </div>
        <div class="mt-1.5 flex items-baseline gap-1.5">
          <span class="text-[22px] font-semibold leading-none tracking-[-0.3px]">
            {{ data?.disk.total ? formatBytes(data.disk.free) : '—' }}
          </span>
          <span class="text-[12px] text-text-5">可用</span>
        </div>
        <div class="mt-2">
          <div class="dh-bar"><i :style="{ width: `${Math.min(diskPercent, 100)}%` }" /></div>
          <div class="mt-1 text-[11px] text-text-5">
            共 {{ data?.disk.total ? formatBytes(data.disk.total) : '—' }}
          </div>
        </div>
      </div>
    </div>

    <div class="grid grid-cols-1 gap-3.5 xl:grid-cols-[1.35fr_1fr]">
      <!-- 有更新的容器 -->
      <div class="dh-card">
        <div class="dh-card-head">
          <Download class="h-3.5 w-3.5 text-text-4" />
          <span>待更新容器</span>
          <span class="ml-auto text-[11.5px] font-normal text-text-5">
            {{ data?.updates.items.length ?? 0 }} 个
          </span>
        </div>
        <div v-if="!data?.updates.checkedAt" class="dh-card-body text-[12.5px] text-text-4">
          还没有做过巡检。点上面的「立即巡检」可以只读地检查一遍所有容器。
        </div>
        <EmptyState
          v-else-if="!data.updates.items.length"
          :icon="CircleDot"
          title="没有待更新的容器"
          description="所有可比对的容器都与仓库摘要一致。"
        />
        <div v-else class="divide-y divide-line-1">
          <RouterLink
            v-for="item in data.updates.items"
            :key="item.container"
            :to="`/containers/${encodeURIComponent(item.container)}`"
            class="flex items-center gap-3 px-3.5 py-2.5 transition-colors hover:bg-ink-750"
          >
            <div class="grid h-[30px] w-[30px] flex-none place-items-center rounded-[9px] bg-line-2 text-[11px] font-semibold text-[#5eead4]">
              {{ item.container.slice(0, 2).toUpperCase() }}
            </div>
            <div class="min-w-0 flex-1">
              <div class="truncate text-[12.5px] font-medium">{{ item.container }}</div>
              <div class="truncate font-mono text-[11px] text-text-5">{{ shortImage(item.image) }}</div>
            </div>
            <span class="dh-badge dh-badge-warn flex-none">有新版本</span>
          </RouterLink>
        </div>
      </div>

      <!-- 环境 + 实时活动 -->
      <div class="flex flex-col gap-3.5">
        <div class="dh-card">
          <div class="dh-card-head">
            <Server class="h-3.5 w-3.5 text-text-4" />
            <span>Docker 环境</span>
          </div>
          <div class="dh-card-body grid grid-cols-2 gap-x-4 gap-y-2 text-[12px]">
            <template v-if="data?.docker">
              <div class="text-text-5">引擎版本</div>
              <div class="truncate text-text-2">{{ data.docker.version }}</div>
              <div class="text-text-5">系统</div>
              <div class="truncate text-text-2">{{ data.docker.os }}</div>
              <div class="text-text-5">架构</div>
              <div class="truncate text-text-2">{{ data.docker.arch }}</div>
              <div class="text-text-5">内核</div>
              <div class="truncate text-text-2">{{ data.docker.kernel }}</div>
              <div class="text-text-5">CPU / 内存</div>
              <div class="truncate text-text-2">
                {{ data.docker.cpus }} 核 / {{ formatBytes(data.docker.memTotal) }}
              </div>
              <div class="text-text-5">数据根目录</div>
              <div class="col-span-1 truncate font-mono text-[11px] text-text-2">{{ data.docker.rootDir }}</div>
              <div class="text-text-5">加速源</div>
              <div class="truncate text-text-2">
                <template v-if="data.docker.mirrors.length">
                  {{ data.docker.mirrors.length }} 个（守护进程配置）
                </template>
                <template v-else>
                  <span class="text-text-5">未配置</span>
                </template>
              </div>
            </template>
            <template v-else>
              <div class="col-span-2 text-text-5">无法读取 Docker 信息</div>
            </template>
          </div>
        </div>

        <div class="dh-card min-h-[160px]">
          <div class="dh-card-head">
            <Activity class="h-3.5 w-3.5 text-text-4" />
            <span>实时活动</span>
            <Loader2 v-if="loading" class="ml-auto h-3 w-3 dh-spin text-text-5" />
          </div>
          <div v-if="!live.length" class="dh-card-body text-[12.5px] text-text-4">
            正在监听更新与计划任务的进度，触发后这里会实时出现。
          </div>
          <div v-else class="flex flex-col">
            <div
              v-for="(item, i) in live"
              :key="i"
              class="flex items-center gap-2 border-b border-[#171f2a] px-3.5 py-2 text-[12px] last:border-b-0"
            >
              <span
                class="h-[6px] w-[6px] flex-none rounded-full"
                :class="item.status === 'failed' ? 'bg-err' : item.status === 'success' ? 'bg-run' : 'bg-accent'"
              />
              <span class="min-w-0 flex-1 truncate text-text-3">{{ item.text }}</span>
            </div>
          </div>
        </div>
      </div>
    </div>

    <!-- 最近操作 -->
    <div class="dh-card">
      <div class="dh-card-head">
        <Activity class="h-3.5 w-3.5 text-text-4" />
        <span>最近操作记录</span>
        <RouterLink to="/settings" class="ml-auto text-[11.5px] font-normal text-text-5 hover:text-accent">
          查看全部
        </RouterLink>
      </div>
      <div v-if="!data?.recent?.length" class="dh-card-body text-[12.5px] text-text-4">暂无记录</div>
      <div v-else class="overflow-x-auto">
        <table class="dh-table">
          <thead>
            <tr>
              <th class="w-[140px]">时间</th>
              <th class="w-[80px]">类型</th>
              <th class="w-[160px]">对象</th>
              <th>说明</th>
            </tr>
          </thead>
          <tbody>
            <tr v-for="l in data.recent" :key="l.id">
              <td class="whitespace-nowrap text-[11.5px] text-text-5">{{ relativeTime(l.ts) }}</td>
              <td>
                <span class="dh-badge" :class="`dh-badge-${statusOfLog(l)}`">{{ l.kind }}</span>
              </td>
              <td class="max-w-[160px] truncate font-mono text-[11.5px] text-text-3">{{ l.ref || '—' }}</td>
              <td class="text-[12px] text-text-2">{{ l.message }}</td>
            </tr>
          </tbody>
        </table>
      </div>
    </div>
  </div>
</template>
