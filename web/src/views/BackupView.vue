<script setup lang="ts">
import { computed, onMounted, ref } from 'vue'
import {
  Archive,
  AlertTriangle,
  CheckCircle2,
  Eye,
  FileCode2,
  FolderTree,
  HardDrive,
  RefreshCw,
  RotateCcw,
  Trash2,
} from 'lucide-vue-next'
import { api } from '@/api/client'
import type { BackupStats, ContainerView, DiffEntry, ProjectInfo, RestoreResult, SnapshotItem } from '@/api/types'
import { formatBytes, relativeTime } from '@/utils/format'
import { useToastStore } from '@/stores/toast'
import EmptyState from '@/components/EmptyState.vue'
import Modal from '@/components/Modal.vue'

const toast = useToastStore()
const tab = ref<'snapshots' | 'projects'>('snapshots')
const backups = ref<SnapshotItem[]>([])
const stats = ref<BackupStats | null>(null)
/** 当前生效的宿主路径映射（启动时自动识别 + DOCKHELM_HOST_ROOTS） */
const pathMappings = computed(() => stats.value?.pathMappings ?? [])
const containers = ref<ContainerView[]>([])
const projects = ref<ProjectInfo[]>([])
const loading = ref(true)
const busy = ref('')
const snapshotTarget = ref('')
const removeTarget = ref<SnapshotItem | null>(null)

const restoreTarget = ref<SnapshotItem | null>(null)
const diff = ref<DiffEntry[]>([])
const diffLoading = ref(false)
const restoreResult = ref<RestoreResult | null>(null)
const restoring = ref(false)
const showRestore = ref(false)

const showFile = ref<{ path: string; content: string } | null>(null)
const pruneKeep = ref(10)
const pruneDays = ref(30)
const pruneDaysArr = [7, 30, 90, 365]

const grouped = computed(() => {
  const map = new Map<string, SnapshotItem[]>()
  for (const b of backups.value) {
    const list = map.get(b.container) ?? []
    list.push(b)
    map.set(b.container, list)
  }
  return [...map.entries()]
})

async function load() {
  loading.value = true
  // 用 allSettled：容器列表依赖 Docker 守护进程，它一时连不上不该把备份页
  // 其它信息（尤其是路径映射）一起清空 —— 那些数据本地就有。
  const [b, s, c] = await Promise.allSettled([
    api.get<{ backups: SnapshotItem[] }>('/api/backups'),
    api.get<BackupStats>('/api/backups/stats'),
    api.get<{ containers: ContainerView[] }>('/api/containers'),
  ])
  const errText = (e: unknown) => (e instanceof Error ? e.message : String(e))
  if (b.status === 'fulfilled') backups.value = b.value.backups ?? []
  else toast.error('读取快照失败', errText(b.reason))
  if (s.status === 'fulfilled') stats.value = s.value
  else toast.error('读取备份统计失败', errText(s.reason))
  if (c.status === 'fulfilled') containers.value = c.value.containers ?? []
  else toast.error('读取容器列表失败', errText(c.reason))
  loading.value = false
}

async function loadProjects() {
  try {
    const res = await api.get<{ projects: ProjectInfo[] }>('/api/backups/projects')
    projects.value = res.projects ?? []
  } catch (e) {
    toast.error('读取项目失败', e instanceof Error ? e.message : String(e))
  }
}

async function doSnapshot() {
  const name = snapshotTarget.value
  if (!name) return
  busy.value = 'snapshot'
  try {
    await api.post('/api/backups/snapshot', { container: name, reason: 'manual' })
    toast.success(`已备份 ${name} 的配置快照`)
    await load()
  } catch (e) {
    toast.error('备份失败', e instanceof Error ? e.message : String(e))
  } finally {
    busy.value = ''
  }
}

async function openRestore(item: SnapshotItem) {
  restoreTarget.value = item
  restoreResult.value = null
  showRestore.value = true
  diffLoading.value = true
  diff.value = []
  try {
    const res = await api.get<{ diff: DiffEntry[] }>('/api/backups/diff', {
      container: item.container,
      ts: item.ts,
    })
    diff.value = res.diff ?? []
  } catch (e) {
    toast.error('无法比对差异', e instanceof Error ? e.message : String(e))
  } finally {
    diffLoading.value = false
  }
}

async function doRestore() {
  const item = restoreTarget.value
  if (!item) return
  restoring.value = true
  try {
    const res = await api.post<RestoreResult>('/api/backups/restore', {
      container: item.container,
      ts: item.ts,
      preSnapshot: true,
      keepBackupContainer: true,
    })
    restoreResult.value = res
    if (res.ok) toast.success('还原完成')
    else toast.error('还原未成功', res.message)
    await load()
  } catch (e) {
    toast.error('还原失败', e instanceof Error ? e.message : String(e))
  } finally {
    restoring.value = false
  }
}

async function confirmRemove() {
  const item = removeTarget.value
  if (!item) return
  try {
    await api.del(`/api/backups/${encodeURIComponent(item.container)}/${encodeURIComponent(item.ts)}`)
    toast.success('快照已删除')
    removeTarget.value = null
    await load()
  } catch (e) {
    toast.error('删除失败', e instanceof Error ? e.message : String(e))
  }
}

async function doPrune() {
  busy.value = 'prune'
  try {
    const res = await api.post<{ removed: number; freedBytes: number }>('/api/backups/prune', {
      keepPerContainer: pruneKeep.value,
      maxAgeDays: pruneDays.value,
      maxTotalMB: 2048,
    })
    toast.success(`清理了 ${res.removed} 份快照`, `释放 ${formatBytes(res.freedBytes)}`)
    await load()
  } catch (e) {
    toast.error('清理失败', e instanceof Error ? e.message : String(e))
  } finally {
    busy.value = ''
  }
}

async function openFile(f: { hostPath: string }) {
  try {
    const res = await api.get<{ path: string; content: string }>('/api/backups/file', { path: f.hostPath })
    showFile.value = { path: res.path, content: res.content }
  } catch (e) {
    toast.error('无法读取文件', e instanceof Error ? e.message : String(e))
  }
}

onMounted(async () => {
  await load()
  await loadProjects()
})
</script>

<template>
  <div class="flex flex-col gap-3.5 p-[18px]">
    <!-- 能力边界说明 -->
    <div class="flex items-start gap-3 rounded-[14px] border border-[rgba(45,212,191,.28)] bg-[rgba(45,212,191,.07)] px-4 py-3 text-[12.5px] leading-relaxed text-[#99f6e4]">
      <HardDrive class="mt-[2px] h-4 w-4 flex-none" />
      <div class="min-w-0 flex-1">
        Dockhelm 备份的是<b>容器配置快照</b>（<code>docker inspect</code> 的结果），可以一键还原成同名容器。
        <b>绑定挂载的数据在宿主机目录上</b>，Dockhelm 容器默认看不见 —— 要备份那份数据，得把宿主目录也挂进来。
        compose 项目真正的来源是那份 <code>yaml</code>，用「compose 项目」标签页查看。
      </div>
    </div>

    <div class="flex flex-wrap items-center gap-2.5">
      <div class="dh-seg">
        <button :data-on="tab === 'snapshots'" @click="tab = 'snapshots'">容器配置快照</button>
        <button :data-on="tab === 'projects'" @click="tab = 'projects'">compose 项目</button>
      </div>

      <template v-if="tab === 'snapshots'">
        <select v-model="snapshotTarget" class="dh-select !w-auto !min-w-[190px]">
          <option value="">选择要备份的容器…</option>
          <option v-for="c in containers" :key="c.id" :value="c.name">{{ c.name }}</option>
        </select>
        <button class="dh-btn dh-btn-primary" :disabled="!snapshotTarget || busy === 'snapshot'" @click="doSnapshot">
          <Archive class="h-3.5 w-3.5" />立即备份
        </button>
      </template>

      <div class="ml-auto flex flex-wrap items-center gap-2">
        <select v-model.number="pruneKeep" class="dh-select !w-auto !py-[6px] !text-[11.5px]">
          <option :value="5">每容器保留 5 份</option>
          <option :value="10">每容器保留 10 份</option>
          <option :value="20">每容器保留 20 份</option>
          <option :value="50">每容器保留 50 份</option>
        </select>
        <select v-model.number="pruneDays" class="dh-select !w-auto !py-[6px] !text-[11.5px]">
          <option v-for="d in pruneDaysArr" :key="d" :value="d">保留 {{ d }} 天</option>
        </select>
        <button class="dh-btn dh-btn-sm" :disabled="busy === 'prune'" @click="doPrune">
          <Trash2 class="h-3 w-3" />清理过期
        </button>
        <button class="dh-btn dh-btn-sm" :disabled="loading" @click="load">
          <RefreshCw class="h-3 w-3" :class="loading ? 'dh-spin' : ''" />刷新
        </button>
      </div>
    </div>

    <!-- 统计 -->
    <div class="grid grid-cols-2 gap-3 lg:grid-cols-4">
      <div class="dh-card p-3.5">
        <div class="text-[12px] text-text-4">快照数量</div>
        <div class="mt-1.5 text-[20px] font-semibold leading-none">{{ stats?.snapshots ?? 0 }}</div>
        <div class="mt-1 text-[11px] text-text-5">覆盖 {{ stats?.containers ?? 0 }} 个容器</div>
      </div>
      <div class="dh-card p-3.5">
        <div class="text-[12px] text-text-4">占用空间</div>
        <div class="mt-1.5 text-[20px] font-semibold leading-none">{{ formatBytes(stats?.sizeBytes) }}</div>
        <div class="mt-1 text-[11px] text-text-5">位于 {{ stats?.dir }}</div>
      </div>
      <div class="dh-card p-3.5">
        <div class="text-[12px] text-text-4">Docker 数据根目录</div>
        <div class="mt-1.5">
          <span class="dh-badge" :class="stats?.dockerRootVisible ? 'dh-badge-run' : 'dh-badge-warn'">
            {{ stats?.dockerRootVisible ? '可见' : '不可见' }}
          </span>
        </div>
        <div class="mt-1 truncate font-mono text-[11px] text-text-5">{{ stats?.dockerRoot || '—' }}</div>
      </div>
      <div class="dh-card p-3.5">
        <div class="text-[12px] text-text-4">命名卷目录</div>
        <div class="mt-1.5">
          <span class="dh-badge" :class="stats?.volumeRootMounted ? 'dh-badge-run' : 'dh-badge-warn'">
            {{ stats?.volumeRootMounted ? '已挂载' : '未挂载' }}
          </span>
        </div>
        <div class="mt-1 text-[11px] leading-relaxed text-text-5">
          {{ stats?.volumeRootMounted ? '可以读取卷数据' : '想备份卷数据需只读挂载 /var/lib/docker/volumes' }}
        </div>
      </div>
    </div>

    <!-- 路径映射：把「Dockhelm 到底看见了什么」摊开给人看 -->
    <div class="dh-card">
      <div class="dh-card-head">
        <FolderTree class="h-3.5 w-3.5 text-text-4" />
        <span>宿主路径映射</span>
        <span class="ml-auto text-[11.5px] font-normal text-text-5">
          冒号两边不必写一样，右边叫什么都可以
        </span>
      </div>
      <div class="dh-card-body">
        <div v-if="!pathMappings.length" class="text-[12px] leading-relaxed text-text-4">
          没有识别到任何挂载映射 —— 说明 Dockhelm 看不到宿主机的 docker 目录，读不到你的 compose 文件。
          在 compose 里挂一行即可，例如
          <code class="font-mono text-accent">/volume1/docker:/host/docker</code>
          （右边叫什么名字都行，启动时会自动识别）。
        </div>
        <div v-else class="space-y-1.5">
          <div
            v-for="m in pathMappings"
            :key="m.host + '>' + m.container"
            class="flex flex-wrap items-center gap-x-2 gap-y-1"
          >
            <span class="dh-badge" :class="m.visible ? 'dh-badge-run' : 'dh-badge-warn'">
              {{ m.visible ? '可见' : '不可见' }}
            </span>
            <span class="font-mono text-[11.5px] text-text-3">{{ m.host }}</span>
            <span class="text-text-5">→</span>
            <span class="font-mono text-[11.5px] text-accent">{{ m.container }}</span>
            <span v-if="m.host === m.container" class="text-[11px] text-text-5">两边一致</span>
            <span v-if="m.source === 'env'" class="text-[11px] text-text-5">来自 DOCKHELM_HOST_ROOTS</span>
          </div>
        </div>
      </div>
    </div>

    <!-- 快照列表 -->
    <div v-if="tab === 'snapshots'" class="dh-card">
      <div class="dh-card-head">
        <Archive class="h-3.5 w-3.5 text-text-4" />
        <span>配置快照</span>
        <span class="ml-auto text-[11.5px] font-normal text-text-5">{{ backups.length }} 份</span>
      </div>

      <EmptyState
        v-if="!backups.length"
        :icon="Archive"
        :title="loading ? '正在载入…' : '还没有任何快照'"
        description="更新容器时 Dockhelm 会自动写一份快照（用于失败回滚），你也可以在这里手动备份。"
      />

      <div v-else class="flex flex-col">
        <template v-for="[name, list] in grouped" :key="name">
          <div class="flex items-center gap-2 border-b border-line-1 bg-ink-750 px-3.5 py-2">
            <span class="text-[12.5px] font-medium">{{ name }}</span>
            <span class="dh-badge dh-badge-plain">{{ list.length }} 份</span>
          </div>
          <div
            v-for="item in list"
            :key="item.ts"
            class="flex flex-wrap items-center gap-2.5 border-b border-[#171f2a] px-3.5 py-2.5 last:border-b-0"
          >
            <div class="min-w-[150px] flex-1">
              <div class="font-mono text-[11.5px] text-text-2">{{ item.ts }}</div>
              <div class="text-[11px] text-text-5">
                {{ relativeTime(item.created) }} · {{ formatBytes(item.size) }}
              </div>
            </div>
            <div class="min-w-[160px] flex-1 truncate font-mono text-[11px] text-text-4" :title="item.image">
              {{ item.image }}
            </div>
            <span class="dh-badge" :class="item.running ? 'dh-badge-run' : 'dh-badge-stop'">
              {{ item.running ? '当时运行中' : '当时已停止' }}
            </span>
            <div class="flex gap-1.5">
              <button class="dh-btn dh-btn-sm" @click="openRestore(item)">
                <RotateCcw class="h-3 w-3" />还原…
              </button>
              <button class="dh-btn dh-btn-sm dh-btn-danger" @click="removeTarget = item">
                <Trash2 class="h-3 w-3" />
              </button>
            </div>
          </div>
        </template>
      </div>
    </div>

    <!-- compose 项目 -->
    <div v-else class="dh-card">
      <div class="dh-card-head">
        <FolderTree class="h-3.5 w-3.5 text-text-4" />
        <span>compose 项目</span>
        <span class="ml-auto text-[11.5px] font-normal text-text-5">
          来自容器标签 com.docker.compose.project.config_files
        </span>
      </div>

      <EmptyState
        v-if="!projects.length"
        :icon="FolderTree"
        :title="loading ? '正在载入…' : '没有发现 compose 项目'"
        description="只有由 docker compose 创建的容器才会带上项目标签。"
      />

      <div v-else class="flex flex-col">
        <div v-for="p in projects" :key="p.project" class="border-b border-line-1 p-3.5 last:border-b-0">
          <div class="flex flex-wrap items-center gap-2">
            <span class="text-[13px] font-semibold">{{ p.project }}</span>
            <span class="dh-badge dh-badge-plain">{{ p.containers.length }} 个容器</span>
            <span v-if="p.workDir" class="truncate font-mono text-[11px] text-text-6" :title="p.workDir">
              {{ p.workDir }}
            </span>
          </div>
          <div class="mt-1.5 flex flex-wrap gap-1">
            <span
              v-for="c in p.containers"
              :key="c"
              class="rounded-md bg-ink-800 px-1.5 py-[2px] font-mono text-[10.5px] text-text-4"
            >
              {{ c }}
            </span>
          </div>

          <div v-if="p.readable.length" class="mt-2.5 flex flex-col gap-1.5">
            <div
              v-for="f in p.readable"
              :key="f.path"
              class="flex items-center gap-2 rounded-[9px] border border-line-1 bg-ink-800 px-2.5 py-1.5"
            >
              <FileCode2 class="h-3.5 w-3.5 flex-none text-accent" />
              <div class="min-w-0 flex-1">
                <div class="text-[12px] text-text-2">{{ f.name }}</div>
                <div class="truncate font-mono text-[10.5px] text-text-6">{{ f.hostPath }}</div>
              </div>
              <span class="text-[11px] text-text-5">{{ formatBytes(f.size) }}</span>
              <button class="dh-btn dh-btn-sm" @click="openFile(f)">
                <Eye class="h-3 w-3" />查看
              </button>
            </div>
          </div>

          <div
            v-if="p.unreadable.length"
            class="mt-2.5 flex flex-col gap-1.5 rounded-[9px] border border-[rgba(245,165,36,.3)] bg-[rgba(245,165,36,.06)] px-2.5 py-2"
          >
            <div class="flex items-center gap-2 text-[11.5px] text-[#fcd34d]">
              <AlertTriangle class="h-3.5 w-3.5" />
              以下文件<b>知道路径但在 Dockhelm 容器里看不见</b>，请自行备份：
            </div>
            <div
              v-for="u in p.unreadable"
              :key="u"
              class="truncate font-mono text-[10.5px] text-text-5"
              :title="u"
            >
              {{ u }}
            </div>
            <div class="text-[11px] leading-relaxed text-text-5">
              解决办法：在 Dockhelm 的 compose 里把这个目录也挂进来即可 ——
              <b class="text-text-3">右边叫什么名字都行</b>（例如
              <code>- /volume1/docker:/host/docker</code>），启动时会自动识别，
              容器标签里的宿主路径就能换算过去。
            </div>
          </div>
        </div>
      </div>
    </div>

    <!-- 还原 -->
    <Modal
      :open="showRestore"
      title="还原容器配置"
      :subtitle="restoreTarget ? `${restoreTarget.container} · ${restoreTarget.ts}` : ''"
      width="660px"
      @close="showRestore = false; restoreResult = null"
    >
      <div v-if="restoreResult" class="flex flex-col gap-3">
        <div
          class="flex items-center gap-2 rounded-[10px] border px-3 py-2.5 text-[12.5px]"
          :class="
            restoreResult.ok
              ? 'border-[rgba(52,211,153,.4)] bg-[rgba(52,211,153,.08)] text-[#4ade80]'
              : 'border-[rgba(248,113,113,.35)] bg-[rgba(248,113,113,.08)] text-[#fca5a5]'
          "
        >
          <component :is="restoreResult.ok ? CheckCircle2 : AlertTriangle" class="h-4 w-4 flex-none" />
          {{ restoreResult.message }}
        </div>
        <div class="dh-scroll max-h-[260px] overflow-auto rounded-[10px] border border-line-1 bg-ink-800 p-3">
          <div v-for="(s, i) in restoreResult.steps" :key="i" class="py-[3px] font-mono text-[11.5px] text-text-3">
            {{ s }}
          </div>
        </div>
      </div>

      <div v-else class="flex flex-col gap-3">
        <div class="flex items-start gap-2.5 rounded-[10px] border border-[rgba(245,165,36,.3)] bg-[rgba(245,165,36,.07)] px-3 py-2.5 text-[12px] leading-relaxed text-[#fcd34d]">
          <AlertTriangle class="mt-[2px] h-4 w-4 flex-none" />
          <div>
            还原会用快照里的配置创建一个<b>新的同名容器</b>，当前容器会被停止并改名成
            <code>&lt;名字&gt;__restorebak_&lt;时间&gt;</code> 保留下来（不会删除）。
            还原前会自动为<b>当前状态</b>再存一份快照，改错了还能退回来。
          </div>
        </div>

        <div>
          <div class="mb-2 flex items-center gap-2 text-[12px] text-text-3">
            <span>与当前状态的差异</span>
            <RefreshCw v-if="diffLoading" class="h-3 w-3 dh-spin text-text-5" />
          </div>
          <div v-if="diffLoading" class="text-[12px] text-text-5">正在比对…</div>
          <div
            v-else-if="!diff.length"
            class="rounded-[9px] border border-line-1 bg-ink-800 px-3 py-2.5 text-[12px] text-text-4"
          >
            快照与当前配置没有差异。
          </div>
          <div v-else class="dh-scroll max-h-[240px] overflow-auto rounded-[10px] border border-line-1">
            <table class="dh-table">
              <thead>
                <tr>
                  <th class="w-[100px]">字段</th>
                  <th>快照里的值</th>
                  <th>当前值</th>
                </tr>
              </thead>
              <tbody>
                <tr v-for="d in diff" :key="d.field">
                  <td class="text-[12px] text-text-3">{{ d.field }}</td>
                  <td class="max-w-[220px] whitespace-pre-wrap break-all font-mono text-[11px] text-[#5eead4]">
                    {{ d.snapshot || '（空）' }}
                  </td>
                  <td class="max-w-[220px] whitespace-pre-wrap break-all font-mono text-[11px] text-[#fbbf24]">
                    {{ d.current || '（空）' }}
                  </td>
                </tr>
              </tbody>
            </table>
          </div>
        </div>
      </div>

      <template #footer>
        <button class="dh-btn" @click="showRestore = false; restoreResult = null">
          {{ restoreResult ? '关闭' : '取消' }}
        </button>
        <button v-if="!restoreResult" class="dh-btn dh-btn-danger" :disabled="restoring" @click="doRestore">
          <RotateCcw class="h-3.5 w-3.5" />{{ restoring ? '还原中…' : '确认还原' }}
        </button>
      </template>
    </Modal>

    <!-- 删除快照 -->
    <Modal
      :open="!!removeTarget"
      title="删除快照"
      :subtitle="removeTarget ? `${removeTarget.container} · ${removeTarget.ts}` : ''"
      width="400px"
      @close="removeTarget = null"
    >
      <div class="text-[12.5px] text-text-3">
        删除后这份配置快照就无法再用于还原。已经运行中的容器不受影响。
      </div>
      <template #footer>
        <button class="dh-btn" @click="removeTarget = null">取消</button>
        <button class="dh-btn dh-btn-danger" @click="confirmRemove">确认删除</button>
      </template>
    </Modal>

    <!-- 查看文件 -->
    <Modal
      :open="!!showFile"
      title="文件内容"
      :subtitle="showFile?.path"
      width="720px"
      @close="showFile = null"
    >
      <pre class="dh-scroll max-h-[520px] overflow-auto whitespace-pre-wrap break-all rounded-[10px] border border-line-1 bg-ink-800 p-3 font-mono text-[11.5px] leading-[1.7] text-text-2">{{ showFile?.content }}</pre>
    </Modal>
  </div>
</template>
