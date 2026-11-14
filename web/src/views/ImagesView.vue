<script setup lang="ts">
import { computed, onMounted, ref } from 'vue'
import { Layers, RefreshCw, Search, Sparkles, Trash2 } from 'lucide-vue-next'
import { api } from '@/api/client'
import type { ImageView } from '@/api/types'
import { formatBytes, relativeTime } from '@/utils/format'
import { useToastStore } from '@/stores/toast'
import EmptyState from '@/components/EmptyState.vue'
import Modal from '@/components/Modal.vue'

const toast = useToastStore()
const images = ref<ImageView[]>([])
const totalSize = ref(0)
const loading = ref(false)
const pruning = ref(false)
const removing = ref(false)
const keyword = ref('')
const onlyUnused = ref(false)
const removeTarget = ref<ImageView | null>(null)
/** 「清理未使用镜像」的二次确认（存在带标签的未使用镜像时才弹）。 */
const pruneConfirm = ref(false)

const filtered = computed(() => {
  let list = images.value
  const k = keyword.value.trim().toLowerCase()
  if (k) {
    list = list.filter(
      (i) =>
        i.tags.some((t) => t.toLowerCase().includes(k)) ||
        (i.repo ?? '').toLowerCase().includes(k) ||
        i.shortId.includes(k),
    )
  }
  if (onlyUnused.value) list = list.filter((i) => i.containers === 0)
  return list
})

/**
 * 「未使用」的唯一口径：没有任何容器引用（containers === 0）。
 * 与表格「被引用」列的徽标同源 —— 之前这里筛的是 dangling（无 tag），
 * 于是列表里明明标着「未使用」的 ipv6sync:latest，勾上筛选却一个都不剩。
 */
const unusedImages = computed(() => images.value.filter((i) => i.containers === 0))
const unusedCount = computed(() => unusedImages.value.length)
const unusedSize = computed(() =>
  unusedImages.value.reduce((sum, i) => sum + (i.size || 0), 0),
)
/** 带标签的未使用镜像：删掉后本地就没了，用到它的容器下次启动要重新拉取。 */
const taggedUnused = computed(() => unusedImages.value.filter((i) => i.tags.length > 0))

async function load() {
  loading.value = true
  try {
    const res = await api.get<{ images: ImageView[]; totalSize: number }>('/api/images')
    images.value = res.images ?? []
    totalSize.value = res.totalSize ?? 0
  } catch (e) {
    toast.error('读取镜像失败', e instanceof Error ? e.message : String(e))
  } finally {
    loading.value = false
  }
}

async function prune(all: boolean) {
  pruning.value = true
  try {
    const res = await api.post<{ count: number; freedBytes: number }>(
      `/api/images/prune${all ? '?all=1' : ''}`,
    )
    // 数字要跟运行记录同一句：几个 + 多大（只报体积的话「清了几个」得靠猜）
    if (!res.count) toast.success('没有可清理的未使用镜像', '当前已是干净状态')
    else toast.success(`已清理 ${res.count} 个未使用镜像`, `释放 ${formatBytes(res.freedBytes)}`)
    pruneConfirm.value = false
    await load()
  } catch (e) {
    toast.error('清理失败', e instanceof Error ? e.message : String(e))
  } finally {
    pruning.value = false
  }
}

/** 有带标签的未使用镜像时先确认：它们删掉后要重新 pull，不能悄悄删。 */
function startPrune() {
  if (taggedUnused.value.length > 0) {
    pruneConfirm.value = true
    return
  }
  void prune(false)
}

async function confirmRemove() {
  const img = removeTarget.value
  if (!img || removing.value) return
  removing.value = true
  try {
    await api.del(`/api/images/${encodeURIComponent(img.id)}`, { force: true })
    toast.success('镜像已删除')
    removeTarget.value = null
    await load()
  } catch (e) {
    toast.error('删除失败', e instanceof Error ? e.message : String(e))
  } finally {
    removing.value = false
  }
}

onMounted(() => void load())
</script>

<template>
  <div class="flex flex-col gap-3.5 p-[18px]">
    <div class="dh-phead">
      <div class="dh-h1">镜像</div>
      <div class="dh-sub">
        {{ images.length }} 个 · 占用 {{ formatBytes(totalSize) }} ·
        未使用 {{ unusedCount }} 个（可回收 {{ formatBytes(unusedSize) }}）
      </div>
      <div class="ml-auto flex gap-2">
        <button class="dh-btn" :disabled="loading" @click="load">
          <RefreshCw class="h-3.5 w-3.5" :class="loading ? 'dh-spin' : ''" />刷新
        </button>
        <button class="dh-btn dh-btn-primary" :disabled="pruning || !unusedCount" @click="startPrune">
          <Sparkles class="h-3.5 w-3.5" :class="pruning ? 'dh-spin' : ''" />
          清理未使用镜像
        </button>
      </div>
    </div>

    <div class="flex flex-wrap items-center gap-2.5">
      <div class="relative min-w-[170px] flex-1 sm:max-w-[260px]">
        <Search class="pointer-events-none absolute left-2.5 top-1/2 h-3.5 w-3.5 -translate-y-1/2 text-text-5" />
        <input v-model="keyword" class="dh-input !pl-8" placeholder="按标签、仓库名或 ID 搜索" />
      </div>

      <label class="flex cursor-pointer items-center gap-2 text-[12px] text-text-3">
        <input v-model="onlyUnused" type="checkbox" class="h-[14px] w-[14px] accent-accent" />
        只看未使用镜像 ({{ unusedCount }})
      </label>
    </div>

    <div class="dh-card">
      <div class="dh-card-head">
        <Layers class="h-3.5 w-3.5 text-text-4" />
        <span>镜像列表</span>
        <span class="ml-auto text-[11.5px] font-normal text-text-5">{{ filtered.length }} 个</span>
      </div>
      <EmptyState
        v-if="!filtered.length"
        :icon="Layers"
        :title="loading ? '正在载入…' : '没有匹配的镜像'"
        description="镜像会随着容器更新不断积累，定期清理未使用镜像可以回收空间。"
      />
      <div v-else class="overflow-x-auto">
        <table class="dh-table">
          <thead>
            <tr>
              <th>标签</th>
              <th class="w-[130px]">镜像 ID</th>
              <th class="w-[100px]">大小</th>
              <th class="w-[130px]">创建时间</th>
              <th class="w-[90px]">被引用</th>
              <th class="w-[70px]" />
            </tr>
          </thead>
          <tbody>
            <tr v-for="img in filtered" :key="img.id">
              <td>
                <div v-if="img.tags.length" class="flex flex-wrap gap-1">
                  <span
                    v-for="t in img.tags"
                    :key="t"
                    class="rounded-md bg-ink-800 px-1.5 py-[2px] font-mono text-[11px] text-text-2"
                  >
                    {{ t }}
                  </span>
                </div>
                <!-- 无 tag 的「未使用镜像」也要给出可读名字：后端从 RepoDigests 推导出
                     它原来所属的仓库，绝不能让用户对着一串镜像 ID 猜这是什么。 -->
                <div v-else class="flex flex-wrap items-center gap-1.5">
                  <span class="font-mono text-[11.5px] text-text-2">{{ img.repo || img.shortId }}</span>
                  <span class="dh-badge dh-badge-plain">未使用镜像</span>
                </div>
              </td>
              <td class="font-mono text-[11.5px] text-text-4">{{ img.shortId }}</td>
              <td class="text-[12px] text-text-2">{{ formatBytes(img.size) }}</td>
              <td class="text-[11.5px] text-text-4">{{ relativeTime(img.created) }}</td>
              <td>
                <span v-if="img.containers" class="dh-badge dh-badge-accent">{{ img.containers }} 个容器</span>
                <span v-else class="dh-badge dh-badge-plain">未使用</span>
              </td>
              <td>
                <button
                  class="dh-btn dh-btn-sm dh-btn-danger"
                  :disabled="img.containers > 0"
                  :title="img.containers > 0 ? '还有容器在使用这个镜像，需要先删除容器' : '删除镜像'"
                  @click="removeTarget = img"
                >
                  <Trash2 class="h-3 w-3" />
                </button>
              </td>
            </tr>
          </tbody>
        </table>
      </div>
    </div>

    <Modal
      :open="!!removeTarget"
      title="删除镜像"
      :subtitle="removeTarget?.tags.join(', ') || removeTarget?.repo || removeTarget?.shortId"
      width="440px"
      :busy="removing"
      @close="removeTarget = null"
    >
      <div class="rounded-[10px] border border-line-err bg-soft-err px-3 py-2.5 text-[12px] leading-relaxed text-err-text">
        删除镜像本身不会删除容器，但如果这个镜像还在被容器使用，容器下次启动时会失败。
        如果只是想回收空间，用「清理未使用镜像」更安全。
      </div>
      <template #footer>
        <button class="dh-btn" :disabled="removing" @click="removeTarget = null">取消</button>
        <button class="dh-btn dh-btn-danger" :disabled="removing" @click="confirmRemove">
          {{ removing ? '删除中…' : '确认删除' }}
        </button>
      </template>
    </Modal>

    <!--
      清理确认：只有存在「带标签的未使用镜像」时才出现。
      这类镜像删掉后本地就没了，不像悬空镜像那样纯粹是残渣，必须先说清楚。
    -->
    <Modal
      :open="pruneConfirm"
      title="清理未使用镜像"
      width="480px"
      :busy="pruning"
      @close="pruneConfirm = false"
    >
      <div class="space-y-2.5 text-[12px] leading-relaxed text-text-3">
        <p>
          共 {{ unusedCount }} 个未使用镜像，合计 {{ formatBytes(unusedSize) }}。
        </p>
        <div v-if="taggedUnused.length">
          <p>
            其中 <b class="text-text-1">{{ taggedUnused.length }} 个带标签</b>：
          </p>
          <!-- 一行一个，别用「、」连成一段：仓库名很长，连起来会从词中间折断 -->
          <ul class="mt-1 space-y-0.5 font-mono text-[11.5px] text-text-2">
            <li v-for="t in taggedUnused.slice(0, 5)" :key="t.id" class="[overflow-wrap:anywhere]">
              {{ t.tags[0] }}
            </li>
            <li v-if="taggedUnused.length > 5" class="text-text-4">
              …等 {{ taggedUnused.length }} 个
            </li>
          </ul>
        </div>
        <div
          class="rounded-[10px] border border-line-err bg-soft-err px-3 py-2.5 text-[12px] leading-relaxed text-err-text"
        >
          带标签的镜像删掉后本地就没有了：用到它的容器下次启动会重新拉取，网络不通时会启动失败。
        </div>
      </div>
      <template #footer>
        <button class="dh-btn" :disabled="pruning" @click="pruneConfirm = false">取消</button>
        <button class="dh-btn dh-btn-primary" :disabled="pruning" @click="prune(true)">
          {{ pruning ? '清理中…' : '确认清理' }}
        </button>
      </template>
    </Modal>
  </div>
</template>
