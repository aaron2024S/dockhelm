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
import PortChips from '@/components/PortChips.vue'

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
const checking = ref(false)
const bulkBusy = ref(false)
/** 提交更新任务中 —— 防连点（连点两次会发出两次批量更新请求）。 */
const applying = ref(false)

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

/**
 * 「更新 N 个容器」按钮真正要更新的那一批。
 *
 * 以前按钮文案数的是 `counts.update`（所有 hasUpdate 的容器），实际提交时却过滤掉
 * self/excluded —— 于是会出现「显示 3 个、点下去直接 return」这种计数与行为不一致。
 * 现在文案与提交用同一个来源。
 */
const updatableAll = computed(() => containers.value.filter((c) => c.hasUpdate && !c.self && !c.excluded))

/** 筛选后的列表是否已全选 —— 复选框的选中态要与之同步（以前恒为 false）。 */
const allSelected = computed(
  () => filtered.value.length > 0 && filtered.value.every((c) => selected.value.has(c.name)),
)

async function load(silent = false) {
  if (!silent) loading.value = true
  try {
    const res = await api.get<{ containers: ContainerView[] }>('/api/containers')
    containers.value = res.containers ?? []
  } catch (e) {
    // SSE 触发的静默刷新失败不弹提示：守护进程短暂不可达时会刷屏。
    // 手动点「刷新」失败才需要明确告诉用户。
    if (!silent) toast.error('读取容器失败', e instanceof Error ? e.message : String(e))
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
  if (!names.length || applying.value) return
  applying.value = true
  try {
    await api.post('/api/updates/apply', { names })
    toast.info(`已提交 ${names.length} 个容器的更新`, '镜像未变化的容器会被自动跳过')
    selected.value = new Set()
  } catch (e) {
    toast.error('提交失败', e instanceof Error ? e.message : String(e))
  } finally {
    applying.value = false
  }
}

async function checkAll() {
  checking.value = true
  try {
    const res = await api.post<{ results?: { status: string }[] }>('/api/updates/check', { deep: false })
    const n = (res.results ?? []).filter((r) => r.status === 'update_available').length
    toast.success('巡检完成', n ? `发现 ${n} 个有可用更新` : '所有容器都是最新的')
    await load(true)
  } catch (e) {
    toast.error('巡检失败', e instanceof Error ? e.message : String(e))
  } finally {
    checking.value = false
  }
}

async function updateAll() {
  const names = updatableAll.value.map((c) => c.name)
  if (!names.length || applying.value) return
  applying.value = true
  try {
    await api.post('/api/updates/apply', { names })
    toast.info(`已提交 ${names.length} 个容器的更新`, '镜像未变化的容器会被自动跳过')
  } catch (e) {
    toast.error('提交失败', e instanceof Error ? e.message : String(e))
  } finally {
    applying.value = false
  }
}

/**
 * 批量启停/重启：一次会动很多容器，必须二次确认。
 * 其余破坏性操作（删容器、删镜像、删快照、删计划）都有确认框，这几个以前点了就执行。
 */
async function bulkAct(action: 'start' | 'stop' | 'restart') {
  const names = containers.value
    .filter((c) => selected.value.has(c.name) && !c.self)
    .map((c) => c.name)
  if (!names.length) return
  bulkPending.value = { action, names }
}

const bulkPending = ref<{ action: 'start' | 'stop' | 'restart'; names: string[] } | null>(null)

async function confirmBulkAct() {
  const p = bulkPending.value
  if (!p) return
  bulkPending.value = null
  bulkBusy.value = true
  let ok = 0
  for (const name of p.names) {
    try {
      await api.post(`/api/containers/${encodeURIComponent(name)}/action`, { action: p.action })
      ok += 1
    } catch {
      // 单个失败不中断整批
    }
  }
  bulkBusy.value = false
  selected.value = new Set()
  toast.success(
    `已${labelOf(p.action)} ${ok}/${p.names.length} 个`,
    ok < p.names.length ? '部分容器操作失败，详见日志' : undefined,
  )
  await load(true)
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
  // 基于「筛选后的清单」逐项判断，别拿 selected.size 与 filtered.length 比数量：
  // 选中集合里可能残留着已被筛选隐藏的名字，两者相等时会误清空全选。
  if (allSelected.value) {
    const next = new Set(selected.value)
    for (const c of filtered.value) next.delete(c.name)
    selected.value = next
  } else {
    const next = new Set(selected.value)
    for (const c of filtered.value) next.add(c.name)
    selected.value = next
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
    <!-- 页头：标题 + 状态筛选 + 主操作 -->
    <div class="dh-phead">
      <div class="dh-h1">容器</div>
      <div class="flex flex-wrap gap-1.5">
        <button class="dh-chip" :data-on="filter === 'all'" @click="filter = 'all'">
          全部 {{ counts.all }}
        </button>
        <button class="dh-chip" :data-on="filter === 'running'" @click="filter = 'running'">
          运行中 {{ counts.running }}
        </button>
        <button class="dh-chip" :data-on="filter === 'stopped'" @click="filter = 'stopped'">
          已停止 {{ counts.stopped }}
        </button>
        <button class="dh-chip" :data-on="filter === 'update'" @click="filter = 'update'">
          有更新 {{ counts.update }}
        </button>
      </div>

      <div class="ml-auto flex flex-wrap gap-2">
        <button class="dh-btn" :disabled="loading" @click="load()">
          <RefreshCw class="h-3.5 w-3.5" :class="loading ? 'dh-spin' : ''" />刷新
        </button>
        <button class="dh-btn" :disabled="checking" @click="checkAll">
          <Download class="h-3.5 w-3.5" :class="checking ? 'dh-spin' : ''" />检测更新
        </button>
        <button
          class="dh-btn dh-btn-primary"
          :disabled="!updatableAll.length || applying"
          @click="updateAll"
        >
          {{ applying ? '提交中…' : `更新 ${updatableAll.length} 个容器` }}
        </button>
      </div>
    </div>

    <div class="flex flex-wrap items-center gap-2.5">
      <div class="relative min-w-[180px] flex-1 sm:max-w-[280px]">
        <Search class="pointer-events-none absolute left-2.5 top-1/2 h-3.5 w-3.5 -translate-y-1/2 text-text-5" />
        <input v-model="keyword" class="dh-input !pl-8" placeholder="按名称 / 镜像 / 项目筛选" />
      </div>

      <div v-if="selected.size" class="ml-auto flex flex-wrap items-center gap-2">
        <span class="text-[12px] text-text-4">已选 {{ selected.size }} 个</span>
        <button class="dh-btn dh-btn-sm" @click="toggleAll">
          {{ allSelected ? '取消全选' : '全选本页' }}
        </button>
        <button class="dh-btn dh-btn-sm" :disabled="bulkBusy" @click="bulkAct('restart')">重启</button>
        <button class="dh-btn dh-btn-sm" :disabled="bulkBusy" @click="bulkAct('stop')">停止</button>
        <button
          class="dh-btn dh-btn-sm dh-btn-primary"
          :disabled="!updatableSelected.length || applying"
          @click="updateSelected"
        >
          更新选中
        </button>
      </div>
    </div>

    <div class="dh-banner dh-banner-warn">
      <!-- 整句必须包在同一层里：.dh-banner 是 flex 容器，裸文本会被拆成独立格子，
           中间夹一个 <b> 就会排成三列（"文字 / 加粗 / 文字"并排）。 -->
      <span class="min-w-0 flex-1">
        批量更新前会先核对镜像摘要：<b>只有镜像真的变了才会重启容器</b>，未变化的容器会原样跳过。
      </span>
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
          c.hasUpdate ? 'border-line-warn' : 'border-line-1 hover:border-line-4',
          selected.has(c.name) ? '!border-accent-line bg-accent-soft' : '',
        ]"
      >
        <div class="flex items-start gap-2.5">
          <label class="mt-[3px] flex-none cursor-pointer">
            <input
              type="checkbox"
              class="h-[14px] w-[14px] accent-accent"
              :checked="selected.has(c.name)"
              @change="toggleSelect(c.name)"
            />
          </label>
          <div
            class="grid h-[34px] w-[34px] flex-none place-items-center rounded-[10px] bg-line-2 text-[13px] font-semibold text-accent-text"
          >
            {{ c.name.slice(0, 2).toUpperCase() }}
          </div>
          <div class="min-w-0 flex-1">
            <RouterLink
              :to="`/containers/${encodeURIComponent(c.name)}`"
              class="dh-tap-txt block min-w-0 text-[13px] font-semibold hover:text-accent"
              :title="c.name"
            >
              <span class="truncate">{{ c.name }}</span>
            </RouterLink>
            <div class="truncate font-mono text-[11.5px] text-text-5" :title="c.image">
              {{ shortImage(c.image) }}
            </div>
          </div>
          <div class="relative flex-none">
            <button
              class="dh-tap grid h-6 w-6 place-items-center rounded-md text-text-5 hover:bg-ink-650 hover:text-text-1"
              @click.stop="menuFor = menuFor === c.name ? '' : c.name"
            >
              <MoreVertical class="h-3.5 w-3.5" />
            </button>
            <div
              v-if="menuFor === c.name"
              class="absolute right-0 top-7 z-20 w-[150px] overflow-hidden rounded-[10px] border border-line-3 bg-ink-750 py-1 shadow-[var(--shadow-pop)]"
              @click.stop
            >
              <button class="block w-full px-3 py-1.5 text-left text-[12px] text-text-2 hover:bg-ink-650" @click="act(c, 'restart')">
                重启
              </button>
              <button class="block w-full px-3 py-1.5 text-left text-[12px] text-text-2 hover:bg-ink-650" @click="updateOne(c)">
                检查并更新
              </button>
              <button
                class="block w-full px-3 py-1.5 text-left text-[12px] text-err-text hover:bg-ink-650"
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

        <!-- 端口：发布到宿主机的用实底、加粗宿主端口；只在容器网络里的用虚线框区分。
             超过 3 条收成 +N，鼠标悬停看全部 —— 卡片本身不该被端口撑变形。 -->
        <PortChips :ports="c.portList ?? []" :max="3" />

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

    <!-- 批量选择条已并入页头，这里只保留一个「全选」快捷入口 -->
    <div v-if="filtered.length && !selected.size" class="flex items-center gap-2 text-[12px] text-text-5">
      <label class="flex cursor-pointer items-center gap-2">
        <input
          type="checkbox"
          class="h-[14px] w-[14px] accent-accent"
          :checked="allSelected"
          @change="toggleAll"
        />
        全选当前列表（{{ filtered.length }} 个）
      </label>
    </div>

    <!-- 批量启停/重启确认：一次会动多个容器，必须二次确认 -->
    <Modal
      :open="!!bulkPending"
      :title="`批量${bulkPending ? labelOf(bulkPending.action) : ''}容器`"
      :subtitle="
        bulkPending ? `即将对 ${bulkPending.names.length} 个容器执行${labelOf(bulkPending.action)}` : ''
      "
      width="470px"
      :busy="bulkBusy"
      @close="bulkPending = null"
    >
      <div class="flex flex-col gap-3">
        <div class="dh-banner dh-banner-warn">
          <span class="min-w-0 flex-1">
            批量操作会逐个执行，<b>单个容器失败不会中断其余容器</b>，结束后可在运行记录里查看结果。
          </span>
        </div>
        <div class="max-h-[200px] overflow-auto rounded-[10px] border border-line-2 bg-ink-800 p-2.5">
          <div
            v-for="n in bulkPending?.names ?? []"
            :key="n"
            class="font-mono text-[11.5px] leading-relaxed text-text-3"
          >
            {{ n }}
          </div>
        </div>
      </div>
      <template #footer>
        <button class="dh-btn" :disabled="bulkBusy" @click="bulkPending = null">取消</button>
        <button class="dh-btn dh-btn-primary" :disabled="bulkBusy" @click="confirmBulkAct">
          {{ bulkBusy ? '执行中…' : `确认${bulkPending ? labelOf(bulkPending.action) : ''}` }}
        </button>
      </template>
    </Modal>

    <!-- 删除确认 -->
    <Modal
      :open="!!removeTarget"
      title="删除容器"
      :subtitle="removeTarget ? `即将删除 ${removeTarget.name}` : ''"
      width="470px"
      :busy="removing"
      @close="removeTarget = null; removeVolumes = false"
    >
      <div class="flex flex-col gap-3">
        <div class="rounded-[10px] border border-line-err bg-soft-err px-3 py-2.5 text-[12px] leading-relaxed text-err-text">
          此操作不可撤销。删除容器不会删除它的镜像，但容器自身的可写层数据会一并消失。
        </div>
        <label class="flex cursor-pointer items-start gap-2.5 text-[12.5px] text-text-2">
          <input v-model="removeVolumes" type="checkbox" class="mt-[3px] h-[14px] w-[14px] accent-err" />
          <span>
            同时删除该容器的<b class="text-err-text">匿名卷</b>
            <br />
            <span class="text-[11.5px] text-text-5">
              勾选后 docker 会连匿名卷一起删掉，卷里的数据将无法找回。命名卷不会被删除。
            </span>
          </span>
        </label>
      </div>
      <template #footer>
        <button
          class="dh-btn"
          :disabled="removing"
          @click="removeTarget = null; removeVolumes = false"
        >
          取消
        </button>
        <button class="dh-btn dh-btn-danger" :disabled="removing" @click="confirmRemove">
          <Trash2 class="h-3.5 w-3.5" />{{ removing ? '删除中…' : '确认删除' }}
        </button>
      </template>
    </Modal>
  </div>
</template>
