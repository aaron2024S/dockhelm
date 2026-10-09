<script setup lang="ts">
import { computed, onMounted, onUnmounted, ref, watch } from 'vue'
import { RouterLink, useRoute } from 'vue-router'
import {
  Box,
  Download,
  MoreVertical,
  Play,
  RefreshCw,
  RotateCw,
  Search,
  Square,
  Trash2,
} from 'lucide-vue-next'
import { api, openStream } from '@/api/client'
import type { ContainerView } from '@/api/types'
import { containerStateLabel, relativeTime, shortImage } from '@/utils/format'
import { useToastStore } from '@/stores/toast'
import EmptyState from '@/components/EmptyState.vue'
import Modal from '@/components/Modal.vue'

const toast = useToastStore()
const route = useRoute()

const containers = ref<ContainerView[]>([])
const loading = ref(true)
const keyword = ref((route.query.q as string) ?? '')
const filter = ref<'all' | 'running' | 'stopped' | 'update'>('all')
const busyName = ref('')
const menuFor = ref('')
const removeTarget = ref<ContainerView | null>(null)
const removeVolumes = ref(false)
const removing = ref(false)
const selected = ref<Set<string>>(new Set())

const filtered = computed(() => {
  let list = containers.value
  const k = keyword.value.trim().toLowerCase()
  if (k) {
    list = list.filter(
      (c) =>
        c.name.toLowerCase().includes(k) ||
        c.image.toLowerCase().includes(k) ||
        c.project.toLowerCase().includes(k),
    )
  }
  if (filter.value === 'running') list = list.filter((c) => c.state === 'running')
  if (filter.value === 'stopped') list = list.filter((c) => c.state !== 'running')
  if (filter.value === 'update') list = list.filter((c) => c.hasUpdate)
  return list
})

const counts = computed(() => ({
  all: containers.value.length,
  running: containers.value.filter((c) => c.state === 'running').length,
  stopped: containers.value.filter((c) => c.state !== 'running').length,
  update: containers.value.filter((c) => c.hasUpdate).length,
}))

const updatableSelected = computed(() =>
  containers.value.filter((c) => selected.value.has(c.name) && !c.self && !c.excluded).map((c) => c.name),
)

async function load(silent = false) {
  if (!silent) loading.value = true
  try {
    const res = await api.get<{ containers: ContainerView[] }>('/api/containers')
    containers.value = res.containers ?? []
  } catch (e) {
    toast.error('读取容器失败', e instanceof Error ? e.message : String(e))
  } finally {
    loading.value = false
  }
}

async function act(c: ContainerView, action: string) {
  busyName.value = c.name
  menuFor.value = ''
  try {
    await api.post(`/api/containers/${encodeURIComponent(c.name)}/action`, { action })
    toast.success(`${c.name} 已${labelOf(action)}`)
    await load(true)
  } catch (e) {
    toast.error(`操作失败`, e instanceof Error ? e.message : String(e))
  } finally {
    busyName.value = ''
  }
}

function labelOf(action: string) {
  return { start: '启动', stop: '停止', restart: '重启', pause: '暂停', unpause: '恢复' }[action] ?? action
}

async function updateOne(c: ContainerView) {
  busyName.value = c.name
  menuFor.value = ''
  try {
    await api.post('/api/updates/apply', { names: [c.name] })
    toast.info(`已提交更新：${c.name}`, '如果镜像没有变化，容器不会被停止')
  } catch (e) {
    toast.error('提交失败', e instanceof Error ? e.message : String(e))
  } finally {
    busyName.value = ''
  }
}

async function updateSelected() {
  const names = updatableSelected.value
  if (!names.length) return
  try {
    await api.post('/api/updates/apply', { names })
    toast.info(`已提交 ${names.length} 个容器的更新`, '镜像未变化的容器会被自动跳过')
    selected.value = new Set()
  } catch (e) {
    toast.error('提交失败', e instanceof Error ? e.message : String(e))
  }
}

async function confirmRemove() {
  const c = removeTarget.value
  if (!c) return
  removing.value = true
  try {
    const res = await api.del<{ message: string }>(
      `/api/containers/${encodeURIComponent(c.name)}`,
      { force: true, volumes: removeVolumes.value },
    )
    toast.success(res.message || '已删除')
    removeTarget.value = null
    removeVolumes.value = false
    await load(true)
  } catch (e) {
    toast.error('删除失败', e instanceof Error ? e.message : String(e))
  } finally {
    removing.value = false
  }
}

function toggleSelect(name: string) {
  const next = new Set(selected.value)
  if (next.has(name)) next.delete(name)
  else next.add(name)
  selected.value = next
}

function toggleAll() {
  if (selected.value.size === filtered.value.length) {
    selected.value = new Set()
  } else {
    selected.value = new Set(filtered.value.map((c) => c.name))
  }
}

let inner: (() => void) | null = null

watch(
  () => route.query.q,
  (q) => {
    keyword.value = (q as string) ?? ''
  },
)

onMounted(() => {
  void load()
  // 更新批次或容器状态变化时自动刷新
  inner = openStream('/api/events/stream', (topic, ev) => {
    if (topic === 'update' && (ev.kind === 'batch_done' || ev.kind === 'container_status')) {
      void load(true)
    }
    if (topic === 'container' && ev.kind === 'action') {
      void load(true)
    }
  })
})

onUnmounted(() => inner?.())

function closeMenu() {
  menuFor.value = ''
}
</script>

<template>
  <div class="flex flex-col gap-3.5 p-[18px]" @click="closeMenu">
    <!-- 筛选条 -->
    <div class="flex flex-wrap items-center gap-2.5">
      <div class="dh-seg">
        <button :data-on="filter === 'all'" @click="filter = 'all'">全部 {{ counts.all }}</button>
        <button :data-on="filter === 'running'" @click="filter = 'running'">运行中 {{ counts.running }}</button>
        <button :data-on="filter === 'stopped'" @click="filter = 'stopped'">已停止 {{ counts.stopped }}</button>
        <button :data-on="filter === 'update'" @click="filter = 'update'">待更新 {{ counts.update }}</button>
      </div>

      <div class="relative min-w-[180px] flex-1 sm:max-w-[280px]">
        <Search class="pointer-events-none absolute left-2.5 top-1/2 h-3.5 w-3.5 -translate-y-1/2 text-text-5" />
        <input v-model="keyword" class="dh-input !pl-8" placeholder="按名称 / 镜像 / 项目筛选" />
      </div>

      <button class="dh-btn" :disabled="loading" @click="load()">
        <RefreshCw class="h-3.5 w-3.5" :class="loading ? 'dh-spin' : ''" />刷新
      </button>

      <template v-if="selected.size">
        <div class="ml-auto flex items-center gap-2">
          <span class="text-[12px] text-text-4">已选 {{ selected.size }} 个</span>
          <button class="dh-btn dh-btn-primary" :disabled="!updatableSelected.length" @click="updateSelected">
            <Download class="h-3.5 w-3.5" />更新选中的 {{ updatableSelected.length }} 个
          </button>
        </div>
      </template>
    </div>

    <div v-if="loading && !containers.length" class="dh-card grid h-[240px] place-items-center">
      <RefreshCw class="h-5 w-5 dh-spin text-text-5" />
    </div>

    <EmptyState
      v-else-if="!filtered.length"
      :icon="Box"
      title="没有匹配的容器"
      description="换个筛选条件，或者确认 Docker 守护进程是否正常。"
    />

    <!-- 卡片网格 -->
    <div v-else class="grid grid-cols-1 gap-3 md:grid-cols-2 xl:grid-cols-3 2xl:grid-cols-4">
      <div
        v-for="c in filtered"
        :key="c.id"
        class="flex flex-col gap-2.5 rounded-[14px] border bg-ink-700 p-3 transition-colors"
        :class="[
          c.hasUpdate ? 'border-[rgba(245,165,36,.42)]' : 'border-line-1 hover:border-line-4',
          selected.has(c.name) ? '!border-[rgba(45,212,191,.5)] bg-[#12201f]' : '',
        ]"
      >
        <div class="flex items-start gap-2.5">
          <label class="mt-[3px] flex-none cursor-pointer">
            <input
              type="checkbox"
              class="h-[14px] w-[14px] accent-[#2dd4bf]"
              :checked="selected.has(c.name)"
              @change="toggleSelect(c.name)"
            />
          </label>
          <div
            class="grid h-[34px] w-[34px] flex-none place-items-center rounded-[10px] bg-line-2 text-[13px] font-semibold text-[#5eead4]"
          >
            {{ c.name.slice(0, 2).toUpperCase() }}
          </div>
          <div class="min-w-0 flex-1">
            <RouterLink
              :to="`/containers/${encodeURIComponent(c.name)}`"
              class="block truncate text-[13px] font-semibold hover:text-accent"
            >
              {{ c.name }}
            </RouterLink>
            <div class="truncate font-mono text-[11.5px] text-text-5" :title="c.image">
              {{ shortImage(c.image) }}
            </div>
          </div>
          <div class="relative flex-none">
            <button
              class="grid h-6 w-6 place-items-center rounded-md text-text-5 hover:bg-ink-650 hover:text-text-1"
              @click.stop="menuFor = menuFor === c.name ? '' : c.name"
            >
              <MoreVertical class="h-3.5 w-3.5" />
            </button>
            <div
              v-if="menuFor === c.name"
              class="absolute right-0 top-7 z-20 w-[150px] overflow-hidden rounded-[10px] border border-line-3 bg-ink-750 py-1 shadow-[0_12px_30px_rgba(0,0,0,.5)]"
              @click.stop
            >
              <button class="block w-full px-3 py-1.5 text-left text-[12px] text-text-2 hover:bg-ink-650" @click="act(c, 'restart')">
                重启
              </button>
              <button class="block w-full px-3 py-1.5 text-left text-[12px] text-text-2 hover:bg-ink-650" @click="updateOne(c)">
                检查并更新
              </button>
              <button
                class="block w-full px-3 py-1.5 text-left text-[12px] text-[#fca5a5] hover:bg-ink-650"
                @click="removeTarget = c; menuFor = ''"
              >
                删除容器…
              </button>
            </div>
          </div>
        </div>

        <div class="flex flex-wrap items-center gap-1.5">
          <span v-if="c.self" class="dh-badge dh-badge-accent">Dockhelm 自身</span>
          <span v-else-if="c.excluded" class="dh-badge dh-badge-plain">已排除</span>
          <span
            class="dh-badge"
            :class="c.state === 'running' ? 'dh-badge-run' : 'dh-badge-stop'"
          >
            <span
              class="h-[7px] w-[7px] rounded-full"
              :class="c.state === 'running' ? 'bg-run' : 'bg-stop'"
            />
            {{ containerStateLabel(c.state) }}
          </span>
          <span v-if="c.hasUpdate" class="dh-badge dh-badge-warn">有新版本</span>
          <span v-if="c.project" class="dh-badge dh-badge-plain" :title="`compose 项目 ${c.project}`">
            {{ c.project }}
          </span>
        </div>

        <div v-if="c.ports.length" class="flex flex-wrap gap-1">
          <span
            v-for="p in c.ports.slice(0, 4)"
            :key="p"
            class="rounded-md bg-ink-800 px-1.5 py-[2px] font-mono text-[10.5px] text-text-4"
          >
            {{ p }}
          </span>
          <span v-if="c.ports.length > 4" class="px-1 text-[10.5px] text-text-5">+{{ c.ports.length - 4 }}</span>
        </div>

        <div class="mt-auto flex gap-1.5 border-t border-line-2 pt-2.5">
          <button
            v-if="c.state === 'running'"
            class="dh-btn dh-btn-sm flex-1"
            :disabled="busyName === c.name"
            @click="act(c, 'stop')"
          >
            <Square class="h-3 w-3" />停止
          </button>
          <button
            v-else
            class="dh-btn dh-btn-sm flex-1"
            :disabled="busyName === c.name"
            @click="act(c, 'start')"
          >
            <Play class="h-3 w-3" />启动
          </button>
          <button
            class="dh-btn dh-btn-sm flex-1"
            :disabled="busyName === c.name || c.self"
            @click="act(c, 'restart')"
          >
            <RotateCw class="h-3 w-3" :class="busyName === c.name ? 'dh-spin' : ''" />重启
          </button>
        </div>

        <div class="text-[10.5px] text-text-6">创建于 {{ relativeTime(c.created) }}</div>
      </div>
    </div>

    <!-- 批量选择条 -->
    <div v-if="filtered.length" class="flex items-center gap-2 text-[12px] text-text-5">
      <label class="flex cursor-pointer items-center gap-2">
        <input
          type="checkbox"
          class="h-[14px] w-[14px] accent-[#2dd4bf]"
          :checked="selected.size === filtered.length && filtered.length > 0"
          @change="toggleAll"
        />
        全选当前列表
      </label>
    </div>

    <!-- 删除确认 -->
    <Modal
      :open="!!removeTarget"
      title="删除容器"
      :subtitle="removeTarget ? `即将删除 ${removeTarget.name}` : ''"
      width="470px"
      @close="removeTarget = null; removeVolumes = false"
    >
      <div class="flex flex-col gap-3">
        <div class="rounded-[10px] border border-[rgba(248,113,113,.35)] bg-[rgba(248,113,113,.07)] px-3 py-2.5 text-[12px] leading-relaxed text-[#fca5a5]">
          此操作不可撤销。删除容器不会删除它的镜像，但容器自身的可写层数据会一并消失。
        </div>
        <label class="flex cursor-pointer items-start gap-2.5 text-[12.5px] text-text-2">
          <input v-model="removeVolumes" type="checkbox" class="mt-[3px] h-[14px] w-[14px] accent-[#f87171]" />
          <span>
            同时删除该容器的<b class="text-[#fca5a5]">匿名卷</b>
            <br />
            <span class="text-[11.5px] text-text-5">
              勾选后 docker 会连匿名卷一起删掉，卷里的数据将无法找回。命名卷不会被删除。
            </span>
          </span>
        </label>
      </div>
      <template #footer>
        <button class="dh-btn" @click="removeTarget = null; removeVolumes = false">取消</button>
        <button class="dh-btn dh-btn-danger" :disabled="removing" @click="confirmRemove">
          <Trash2 class="h-3.5 w-3.5" />{{ removing ? '删除中…' : '确认删除' }}
        </button>
      </template>
    </Modal>
  </div>
</template>
