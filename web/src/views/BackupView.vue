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
  Loader2,
  RefreshCw,
  RotateCcw,
  Save,
  Trash2,
  Upload,
} from 'lucide-vue-next'
import { api } from '@/api/client'
import type { BackupStats, ContainerView, DiffEntry, ProjectInfo, RestoreResult, Settings, SnapshotItem } from '@/api/types'
import { formatBytes, relativeTime } from '@/utils/format'
import { useAppStore } from '@/stores/app'
import { useToastStore } from '@/stores/toast'
import EmptyState from '@/components/EmptyState.vue'
import Modal from '@/components/Modal.vue'
import SettingRow from '@/components/SettingRow.vue'
import ToggleSwitch from '@/components/ToggleSwitch.vue'

/** 快照里的单个挂载点记录（与后端 backup.VolumeEntry 一致）。 */
interface VolumeEntry {
  name: string
  type: string
  source: string
  destination: string
  packed: boolean
  bytes: number
  note: string
}

/** 快照的卷数据概况（与后端 backup.VolumesMeta 一致）。 */
interface VolumesMeta {
  included: boolean
  items: VolumeEntry[]
  warning?: string
}

const toast = useToastStore()
const app = useAppStore()
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
const snapshotWithVolumes = ref(false)
const removeTarget = ref<SnapshotItem | null>(null)

const restoreTarget = ref<SnapshotItem | null>(null)
const diff = ref<DiffEntry[]>([])
const diffLoading = ref(false)
const restoreResult = ref<RestoreResult | null>(null)
const restoring = ref(false)
const removing = ref(false)
/** 「清理过期快照」确认框 —— 这一步会真删文件，必须二次确认。 */
const confirmPrune = ref(false)
const showRestore = ref(false)
const volMeta = ref<VolumesMeta | null>(null)
const withVolumes = ref(false)

const showFile = ref<{ path: string; content: string } | null>(null)
const pruneKeep = ref(10)
const pruneDays = ref(30)

/** 保留策略：这些值会持久化到设置里，「每天自动清理」与手动清理用同一套。 */
const policy = ref<Settings | null>(null)
const savingPolicy = ref(false)

/**
 * 列表按快照时间倒序平铺成一张表 —— 一份快照恰好一行。
 * 时间戳本身是可字典序排序的（20060102-150405），直接比字符串即可。
 */
const rows = computed(() => [...backups.value].sort((a, b) => b.ts.localeCompare(a.ts)))

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
    const item = await api.post<SnapshotItem>('/api/backups/snapshot', {
      container: name,
      reason: 'manual',
      withVolumes: snapshotWithVolumes.value,
    })
    toast.success(
      `已备份 ${name}`,
      snapshotWithVolumes.value
        ? `快照 ${item.ts} · ${formatBytes(item.size)}${item.withData ? '（含卷数据）' : '（卷数据未打包，见快照说明）'}`
        : `快照 ${item.ts} · ${formatBytes(item.size)}`,
    )
    await load()
  } catch (e) {
    toast.error('备份失败', e instanceof Error ? e.message : String(e))
  } finally {
    busy.value = ''
  }
}

/**
 * 差异请求的序号。快速连点两行不同快照的「还原…」时，先发的请求可能后返回，
 * 把 A 快照的差异盖到 B 快照的弹窗上 —— 而用户正是看着这份差异去决定要不要还原的，
 * 看错代价太大。回调里只接受「当前这一次」的结果。
 */
let diffReq = 0

async function openRestore(item: SnapshotItem) {
  const req = ++diffReq
  restoreTarget.value = item
  restoreResult.value = null
  showRestore.value = true
  diffLoading.value = true
  diff.value = []
  volMeta.value = null
  withVolumes.value = false
  try {
    const res = await api.get<{ diff: DiffEntry[]; volumes?: VolumesMeta }>('/api/backups/diff', {
      container: item.container,
      ts: item.ts,
    })
    if (req !== diffReq) return // 已经有更新的目标了，这次结果作废
    diff.value = res.diff ?? []
    volMeta.value = res.volumes ?? null
  } catch (e) {
    if (req !== diffReq) return
    toast.error('无法比对差异', e instanceof Error ? e.message : String(e))
  } finally {
    if (req === diffReq) diffLoading.value = false
  }
}

/** 这份快照里真的被打包了的卷。 */
const packedVolumes = computed(() => (volMeta.value?.items ?? []).filter((v) => v.packed))
/** 知道路径、但没被打包进去的（bind 挂载 / 看不见的卷）——必须显式告诉用户。 */
const unpackedVolumes = computed(() => (volMeta.value?.items ?? []).filter((v) => !v.packed))

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
      withVolumes: withVolumes.value,
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

// ---------- 导入备份包 ----------

const showImport = ref(false)
const importContainer = ref('')
const importTS = ref('')
const importText = ref('')
const importing = ref(false)

/** 打开导入框时给一个合理的默认值：当前时间戳。 */
function openImport() {
  importContainer.value = ''
  importTS.value = new Date()
    .toLocaleString('sv-SE')
    .replace(/[-: ]/g, '')
    .slice(0, 15)
  importText.value = ''
  showImport.value = true
}

/** 从粘贴的 JSON 里猜出容器名，省得用户手打。 */
function guessContainer() {
  try {
    const doc = JSON.parse(importText.value)
    const name = String(doc?.inspect?.Name ?? '').replace(/^\//, '')
    if (name && !importContainer.value) importContainer.value = name
  } catch {
    /* 解析不了就让用户自己填 */
  }
}

async function doImport() {
  if (!importContainer.value.trim() || !importTS.value.trim() || !importText.value.trim()) return
  importing.value = true
  try {
    const item = await api.post<SnapshotItem>('/api/backups/import', {
      container: importContainer.value.trim(),
      ts: importTS.value.trim(),
      content: importText.value,
    })
    toast.success('快照已导入', `${item.container} · ${item.ts}`)
    showImport.value = false
    await load()
  } catch (e) {
    toast.error('导入失败', e instanceof Error ? e.message : String(e))
  } finally {
    importing.value = false
  }
}

async function confirmRemove() {
  const item = removeTarget.value
  if (!item || removing.value) return
  removing.value = true
  try {
    await api.del(`/api/backups/${encodeURIComponent(item.container)}/${encodeURIComponent(item.ts)}`)
    toast.success('快照已删除')
    removeTarget.value = null
    await load()
  } catch (e) {
    toast.error('删除失败', e instanceof Error ? e.message : String(e))
  } finally {
    removing.value = false
  }
}

async function doPrune() {
  confirmPrune.value = false
  busy.value = 'prune'
  try {
    // 不带参数 ⇒ 后端用设置页里的保留策略，与每天自动清理完全一致
    const res = await api.post<{ removed: number; freedBytes: number }>('/api/backups/prune', {})
    toast.success(`清理了 ${res.removed} 份快照`, `释放 ${formatBytes(res.freedBytes)}`)
    await load()
  } catch (e) {
    toast.error('清理失败', e instanceof Error ? e.message : String(e))
  } finally {
    busy.value = ''
  }
}

async function loadPolicy() {
  try {
    policy.value = await api.get<Settings>('/api/settings')
    if (policy.value) {
      pruneKeep.value = policy.value.backupKeepPerContainer
      pruneDays.value = policy.value.backupMaxAgeDays
    }
  } catch {
    policy.value = null
  }
}

async function savePolicy() {
  if (!policy.value) return
  savingPolicy.value = true
  try {
    // 只提交本页负责的四个字段（PATCH）。以前是把整份 Settings PUT 回去 ——
    // 那份快照是进页面时拉的，期间在设置页改过并发度/检测周期就会被悄悄还原。
    policy.value = await api.patch<Settings>('/api/settings', {
      backupKeepPerContainer: policy.value.backupKeepPerContainer,
      backupMaxAgeDays: policy.value.backupMaxAgeDays,
      backupMaxTotalMB: policy.value.backupMaxTotalMB,
      backupKeepPreUpdate: policy.value.backupKeepPreUpdate,
    })
    pruneKeep.value = policy.value.backupKeepPerContainer
    pruneDays.value = policy.value.backupMaxAgeDays
    // 同步全局副本，避免设置页还显示旧值
    void app.loadSettings()
    toast.success('保留策略已保存', '每天自动清理与「清理过期」都会按它执行')
  } catch (e) {
    toast.error('保存失败', e instanceof Error ? e.message : String(e))
  } finally {
    savingPolicy.value = false
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

/** 快照来源的中文名与配色。 */
function reasonLabel(reason: string) {
  switch (reason) {
    case 'pre_update':
      return '更新前'
    case 'scheduled':
      return '定时'
    case 'manual':
      return '手动'
    default:
      return '未知来源'
  }
}

function reasonTone(reason: string) {
  switch (reason) {
    case 'pre_update':
      return 'dh-badge-warn'
    case 'scheduled':
      return 'dh-badge-plain'
    default:
      return 'dh-badge-plain'
  }
}

onMounted(async () => {
  await load()
  await loadProjects()
  await loadPolicy()
})
</script>

<template>
  <div class="flex flex-col gap-3.5 p-[18px]">
    <div class="dh-phead">
      <div class="dh-h1">备份与恢复</div>
      <div class="dh-sub">
        快照 {{ stats?.snapshots ?? 0 }} 份 · 占用 {{ formatBytes(stats?.sizeBytes) }} · 覆盖
        {{ stats?.containers ?? 0 }} 个容器
      </div>
    </div>

    <!-- 能力边界说明 -->
    <div class="dh-banner dh-banner-info items-start">
      <HardDrive class="mt-[2px] h-4 w-4 flex-none" />
      <div class="min-w-0 flex-1 leading-relaxed">
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
        <label
          class="flex cursor-pointer items-center gap-1.5 text-[11.5px] text-text-4"
          :title="
            stats?.volumeRootMounted
              ? '连 named volume 的数据一起打包（体积可能很大）'
              : '宿主机的 /var/lib/docker/volumes 没有映射进来，卷数据打不了包'
          "
        >
          <input v-model="snapshotWithVolumes" type="checkbox" class="h-[13px] w-[13px] accent-warn" />
          含卷数据
        </label>
        <button class="dh-btn dh-btn-primary" :disabled="!snapshotTarget || busy === 'snapshot'" @click="doSnapshot">
          <Archive class="h-3.5 w-3.5" />立即备份
        </button>
      </template>

      <div class="ml-auto flex flex-wrap items-center gap-2">
        <button class="dh-btn dh-btn-sm" @click="openImport">
          <Upload class="h-3 w-3" />导入备份包
        </button>
        <button v-if="tab === 'snapshots'" class="dh-btn dh-btn-sm" :disabled="busy === 'prune'" @click="confirmPrune = true">
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
        <div class="mt-1 text-[11px] text-text-5" :title="String(stats?.dir ?? '')">
          <!-- 路径里没有空格，浏览器默认不换行也不断词 —— 会顶破这张卡、
               再连带把整个内容区撑成横向可滚（手机上就是「右边被推出去」）。 -->
          <span class="break-all">位于 {{ stats?.dir }}</span>
        </div>
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

    <!-- 保留策略 -->
    <div v-if="tab === 'snapshots' && policy" class="dh-card">
      <div class="dh-card-head">
        <Trash2 class="h-3.5 w-3.5 text-text-4" />
        <span>保留策略</span>
        <span class="ml-2 text-[11.5px] font-normal text-text-5">
          定时清理每天跑一次 · 也可随时点「清理过期」立即执行
        </span>
        <button class="dh-btn dh-btn-sm dh-btn-primary ml-auto" :disabled="savingPolicy" @click="savePolicy">
          <Loader2 v-if="savingPolicy" class="h-3 w-3 dh-spin" />
          <Save v-else class="h-3 w-3" />保存策略
        </button>
      </div>
      <div class="flex flex-col gap-3 p-3.5">
        <SettingRow title="每容器保留最近份数" sub="超出后按时间滚动覆盖（更新前快照不计入这个额度）">
          <select v-model.number="policy.backupKeepPerContainer" class="dh-select !w-[110px] !py-[5px] !text-[11.5px]">
            <option :value="0">不限制</option>
            <option :value="5">5 份</option>
            <option :value="10">10 份</option>
            <option :value="20">20 份</option>
            <option :value="50">50 份</option>
          </select>
        </SettingRow>
        <SettingRow title="快照保留期" sub="到期自动清理">
          <select v-model.number="policy.backupMaxAgeDays" class="dh-select !w-[110px] !py-[5px] !text-[11.5px]">
            <option :value="0">不限制</option>
            <option :value="7">7 天</option>
            <option :value="30">30 天</option>
            <option :value="90">90 天</option>
            <option :value="365">365 天</option>
          </select>
        </SettingRow>
        <SettingRow title="快照总容量上限" sub="达到上限后先清最旧的快照">
          <select v-model.number="policy.backupMaxTotalMB" class="dh-select !w-[110px] !py-[5px] !text-[11.5px]">
            <option :value="0">不限制</option>
            <option :value="512">512 MB</option>
            <option :value="2048">2 GB</option>
            <option :value="8192">8 GB</option>
            <option :value="20480">20 GB</option>
          </select>
        </SettingRow>
        <SettingRow title="更新前快照永不自动清理" sub="这是回滚的底牌 —— 自动更新出事后唯一能救回来的东西">
          <ToggleSwitch v-model="policy.backupKeepPreUpdate" label="更新前快照永不自动清理" />
        </SettingRow>
      </div>
    </div>

    <!-- 存储位置 -->
    <div v-if="tab === 'snapshots'" class="dh-card">
      <div class="dh-card-head">
        <HardDrive class="h-3.5 w-3.5 text-text-4" />
        <span>存储位置</span>
      </div>
      <div class="flex flex-col gap-3 p-3.5">
        <div class="rounded-[10px] border border-line-1 bg-ink-800 px-3 py-2.5">
          <div class="font-mono text-[12px] text-text-2">{{ stats?.dir || '/data/backups' }}</div>
          <div class="mt-1.5 text-[11.5px] leading-relaxed text-text-5">
            这是 Dockhelm <b class="text-text-4">容器内</b>的路径。建议把宿主目录映射进来，
            这样容器重建、换镜像都还在自己的快照。
          </div>
        </div>
        <div class="border-t border-line-1" />
        <div class="flex items-start gap-2.5">
          <span class="mt-[3px] h-[15px] w-[15px] flex-none rounded-[5px] bg-warn opacity-60" />
          <div class="text-[11.5px] leading-relaxed text-text-4">
            <b class="text-text-3">bind mount</b> 的数据在宿主机目录上，容器内默认看不见 ——
            Dockhelm 只能把路径记进快照，没法替你打包那份数据；
            <b class="text-text-3">named volume</b> 的内容才在 <code class="text-text-3">/var/lib/docker/volumes</code> 下，可以真正读到。
          </div>
        </div>
        <div v-if="pathMappings.length" class="border-t border-line-1 pt-3">
          <div class="mb-2 text-[12px] text-text-4">宿主机路径映射（启动时从自身容器自动识别）</div>
          <div class="flex flex-col gap-1.5">
            <div v-for="(m, i) in pathMappings" :key="i" class="flex flex-wrap items-center gap-2 text-[11.5px]">
              <span class="dh-badge" :class="m.source === 'auto' ? 'dh-badge-accent' : 'dh-badge-plain'">
                {{ m.source === 'auto' ? '自动识别' : '环境变量' }}
              </span>
              <code class="font-mono text-text-3">{{ m.host }}</code>
              <span class="text-text-6">→</span>
              <code class="font-mono text-text-3">{{ m.container }}</code>
              <span class="dh-badge" :class="m.visible ? 'dh-badge-run' : 'dh-badge-err'">
                {{ m.visible ? '可见' : '看不见' }}
              </span>
            </div>
          </div>
        </div>
      </div>
    </div>


    <!-- 快照列表 -->
    <div v-if="tab === 'snapshots'" class="dh-card">
      <div class="dh-card-head">
        <Archive class="h-3.5 w-3.5 text-text-4" />
        <span>配置快照</span>
        <span class="ml-auto text-[11.5px] font-normal text-text-5">
          {{ backups.length }} 份 · 覆盖 {{ stats?.containers ?? 0 }} 个容器
        </span>
      </div>

      <EmptyState
        v-if="!backups.length"
        :icon="Archive"
        :title="loading ? '正在载入…' : '还没有任何快照'"
        description="更新容器时 Dockhelm 会自动写一份快照（用于失败回滚），你也可以在这里手动备份。"
      />

      <!--
        平铺成一张表：一份快照一行。早前按容器名做分组小标题 + 数据行，
        结果每个容器都占掉两行（标题一行、快照一行），几十份快照扫起来很累，
        容器名也只出现在小标题上 —— 现在把容器名收进行内，一份快照就是一行。
      -->
      <div v-else class="overflow-x-auto">
        <table class="dh-table">
          <thead>
            <tr>
              <th>容器</th>
              <th class="w-[130px]">快照时间</th>
              <th class="w-[80px]">来源</th>
              <th class="w-[110px]">内容</th>
              <th class="w-[80px]">大小</th>
              <th class="w-[90px]">快照时状态</th>
              <th class="w-[130px]" />
            </tr>
          </thead>
          <tbody>
            <tr v-for="item in rows" :key="item.container + '/' + item.ts">
              <td>
                <div class="text-[12.5px] font-medium text-text-1">{{ item.container }}</div>
                <div class="max-w-[220px] truncate font-mono text-[10.5px] text-text-6" :title="item.image">
                  {{ item.image || '—' }}
                </div>
              </td>
              <td>
                <div class="text-[12px] text-text-3">{{ relativeTime(item.created) }}</div>
                <div class="font-mono text-[10.5px] text-text-6">{{ item.ts }}</div>
              </td>
              <td>
                <span class="dh-badge" :class="reasonTone(item.reason)">{{ reasonLabel(item.reason) }}</span>
              </td>
              <td>
                <span class="dh-badge" :class="item.withData ? 'dh-badge-accent' : 'dh-badge-plain'">
                  {{ item.withData ? '配置 + 卷数据' : '仅配置' }}
                </span>
              </td>
              <td class="font-mono text-[11.5px] text-text-4">{{ formatBytes(item.size) }}</td>
              <!--
                这个字段记录的是拍下这份快照那一刻容器的状态（数据来自快照里存的
                inspect），是历史值；容器**现在**是否在运行是另一回事 —— 本页还原
                弹窗的差异表里「当前值」那一列才是真当前状态。列头写「快照时状态」
                就是把这个限定词提到表头，格子里只留「运行中 / 已停止」。
              -->
              <td>
                <span
                  class="dh-badge whitespace-nowrap"
                  :class="item.running ? 'dh-badge-run' : 'dh-badge-stop'"
                  title="拍下这份快照时容器是运行中还是已停止。容器现在是否在运行请看容器页"
                >
                  {{ item.running ? '运行中' : '已停止' }}
                </span>
              </td>
              <td>
                <div class="flex justify-end gap-1.5">
                  <button class="dh-btn dh-btn-sm" @click="openRestore(item)">
                    <RotateCcw class="h-3 w-3" />还原…
                  </button>
                  <button class="dh-btn dh-btn-sm dh-btn-danger" @click="removeTarget = item">
                    <Trash2 class="h-3 w-3" />
                  </button>
                </div>
              </td>
            </tr>
          </tbody>
        </table>
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
            <span class="dh-badge dh-badge-plain">{{ p.containers?.length ?? 0 }} 个容器</span>
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

          <div v-if="p.readable?.length" class="mt-2.5 flex flex-col gap-1.5">
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
            v-if="p.unreadable?.length"
            class="mt-2.5 flex flex-col gap-1.5 rounded-[9px] border border-line-warn bg-soft-warn px-2.5 py-2"
          >
            <div class="flex items-center gap-2 text-[11.5px] text-warn-text">
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
      :busy="restoring"
      @close="showRestore = false; restoreResult = null"
    >
      <div v-if="restoreResult" class="flex flex-col gap-3">
        <div
          class="flex items-center gap-2 rounded-[10px] border px-3 py-2.5 text-[12.5px]"
          :class="
            restoreResult.ok
              ? 'border-line-ok bg-soft-ok text-run-text'
              : 'border-line-err bg-soft-err text-err-text'
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
        <div class="flex items-start gap-2.5 rounded-[10px] border border-line-warn bg-soft-warn px-3 py-2.5 text-[12px] leading-relaxed text-warn-text">
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
                  <td class="max-w-[220px] whitespace-pre-wrap break-all font-mono text-[11px] text-accent-text">
                    {{ d.snapshot || '（空）' }}
                  </td>
                  <td class="max-w-[220px] whitespace-pre-wrap break-all font-mono text-[11px] text-warn-text">
                    {{ d.current || '（空）' }}
                  </td>
                </tr>
              </tbody>
            </table>
          </div>
        </div>

      <!-- 卷数据 -->
        <div v-if="volMeta?.items?.length" class="border-t border-line-1 pt-3">
          <div class="mb-2 text-[12px] text-text-3">卷数据</div>

          <label
            class="flex cursor-pointer items-start gap-2.5 rounded-[9px] border px-3 py-2.5"
            :class="
              packedVolumes.length
                ? 'border-line-1 bg-ink-800'
                : 'cursor-not-allowed border-line-1 bg-ink-800 opacity-60'
            "
          >
            <input
              v-model="withVolumes"
              type="checkbox"
              class="mt-[3px] h-[14px] w-[14px] accent-warn"
              :disabled="!packedVolumes.length"
            />
            <span class="text-[12.5px]">
              <b :class="packedVolumes.length ? 'text-warn-text' : 'text-text-4'">
                含卷数据（覆盖现有文件）
              </b>
              <div class="mt-[2px] text-[11px] leading-relaxed text-text-5">
                <template v-if="packedVolumes.length">
                  这份快照打包了 {{ packedVolumes.length }} 个卷：
                  <span class="font-mono text-text-4">{{ packedVolumes.map((v) => v.name).join('、') }}</span>
                  （共 {{ formatBytes(packedVolumes.reduce((a, v) => a + v.bytes, 0)) }}）。
                  勾选后会把卷<b>当前的内容整个覆盖掉</b> —— 这是不可逆的。
                </template>
                <template v-else>
                  这份快照里没有任何被打包的卷数据，只能还原容器配置。
                </template>
              </div>
            </span>
          </label>

          <div v-if="unpackedVolumes.length" class="mt-2 flex flex-col gap-1">
            <div v-for="v in unpackedVolumes" :key="v.destination" class="flex items-start gap-2 text-[11.5px]">
              <span class="dh-badge flex-none" :class="v.type === 'bind' ? 'dh-badge-warn' : 'dh-badge-plain'">
                {{ v.type === 'bind' ? '绑定挂载' : v.type }}
              </span>
              <code class="flex-none font-mono text-text-4">{{ v.destination }}</code>
              <span class="text-text-6">{{ v.note }}</span>
            </div>
          </div>

          <div
            v-if="withVolumes && packedVolumes.length"
            class="mt-2.5 flex items-start gap-2 rounded-[9px] border border-line-err bg-soft-err px-3 py-2 text-[11.5px] leading-relaxed text-err-text"
          >
            <AlertTriangle class="mt-[1px] h-3.5 w-3.5 flex-none" />
            <span>已勾选「含卷数据」：还原过程中会把这些卷里的现有文件覆盖成快照里的版本，无法撤销。</span>
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

    <!-- 导入备份包 -->
    <Modal
      :open="showImport"
      title="导入备份包"
      subtitle="粘贴一份快照 JSON"
      width="640px"
      :busy="importing"
      @close="showImport = false"
    >
      <div class="flex flex-col gap-3">
        <div class="flex items-start gap-2.5 rounded-[10px] border border-line-1 bg-ink-800 px-3 py-2.5 text-[11.5px] leading-relaxed text-text-4">
          <Download class="mt-[2px] h-3.5 w-3.5 flex-none" />
          <span>
            接受 Dockhelm 写出的快照文件（内容形如 <code>{"_dockhelm":{…},"inspect":{…}}</code>）。
            导入只写文件、不动任何容器；导入后就能像本地快照一样在还原弹窗里使用。
            同名时间戳不会覆盖已有快照，会自动加后缀。
          </span>
        </div>
        <div class="grid grid-cols-2 gap-3">
          <div>
            <label class="mb-[5px] block text-[12px] text-text-4">容器名</label>
            <input
              v-model="importContainer"
              class="dh-input"
              placeholder="redis"
              @change="guessContainer"
            />
            <div class="mt-1 text-[11px] text-text-5">决定这份快照归到哪个容器的列表下。</div>
          </div>
          <div>
            <label class="mb-[5px] block text-[12px] text-text-4">快照时间戳</label>
            <input v-model="importTS" class="dh-input font-mono" placeholder="20261009-094200" />
            <div class="mt-1 text-[11px] text-text-5">格式 20060102-150405。</div>
          </div>
        </div>
        <div>
          <label class="mb-[5px] block text-[12px] text-text-4">快照 JSON</label>
          <textarea
            v-model="importText"
            rows="10"
            class="dh-input dh-scroll h-auto resize-y font-mono text-[11.5px] leading-relaxed"
            placeholder='{"_dockhelm":{...},"inspect":{...}}'
            @change="guessContainer"
          />
        </div>
      </div>
      <template #footer>
        <button class="dh-btn" @click="showImport = false">取消</button>
        <button
          class="dh-btn dh-btn-primary"
          :disabled="importing || !importContainer.trim() || !importTS.trim() || !importText.trim()"
          @click="doImport"
        >
          <Loader2 v-if="importing" class="h-3.5 w-3.5 dh-spin" />
          <Download v-else class="h-3.5 w-3.5" />导入
        </button>
      </template>
    </Modal>

    <!-- 删除快照 -->
    <Modal
      :open="!!removeTarget"
      title="删除快照"
      :subtitle="removeTarget ? `${removeTarget.container} · ${removeTarget.ts}` : ''"
      width="400px"
      :busy="removing"
      @close="removeTarget = null"
    >
      <div class="text-[12.5px] text-text-3">
        删除后这份配置快照就无法再用于还原。已经运行中的容器不受影响。
      </div>
      <template #footer>
        <button class="dh-btn" :disabled="removing" @click="removeTarget = null">取消</button>
        <button class="dh-btn dh-btn-danger" :disabled="removing" @click="confirmRemove">
          <Loader2 v-if="removing" class="h-3.5 w-3.5 dh-spin" />
          <Trash2 v-else class="h-3.5 w-3.5" />{{ removing ? '删除中…' : '确认删除' }}
        </button>
      </template>
    </Modal>

    <!-- 清理过期：会真删文件，必须二次确认（其余破坏性操作都有确认，这里以前是点了就删） -->
    <Modal
      :open="confirmPrune"
      title="清理过期快照"
      width="460px"
      :busy="busy === 'prune'"
      @close="confirmPrune = false"
    >
      <div class="flex flex-col gap-3">
        <div class="flex items-start gap-2.5 rounded-[10px] border border-line-warn bg-soft-warn px-3 py-2.5 text-[12.5px] leading-relaxed text-warn-text">
          <AlertTriangle class="mt-[2px] h-4 w-4 flex-none" />
          <div>
            会按下面的保留策略<b>直接删除快照文件</b>，连同打包进去的卷数据一起，<b>不可恢复</b>。
            「更新前快照永不自动清理」打开时，那部分不会被碰到。
          </div>
        </div>
        <div class="rounded-[10px] border border-line-1 bg-ink-800 px-3 py-2.5 text-[12px] leading-relaxed text-text-3">
          当前策略：每个容器保留最近
          <b>{{ policy?.backupKeepPerContainer ? policy.backupKeepPerContainer + ' 份' : '不限' }}</b>
          · 快照保留 {{ policy?.backupMaxAgeDays ? policy.backupMaxAgeDays + ' 天' : '不限' }}
          · 总容量上限 {{ policy?.backupMaxTotalMB ? policy.backupMaxTotalMB + ' MB' : '不限' }}
          <div class="mt-1 text-[11.5px] text-text-5">要改策略请在上面的「保留策略」卡里改并保存。</div>
        </div>
        <div class="text-[12px] text-text-4">当前共 {{ backups.length }} 份快照，占用 {{ formatBytes(stats?.sizeBytes) }}。</div>
      </div>
      <template #footer>
        <button class="dh-btn" :disabled="busy === 'prune'" @click="confirmPrune = false">取消</button>
        <button class="dh-btn dh-btn-danger" :disabled="busy === 'prune'" @click="doPrune">
          <Loader2 v-if="busy === 'prune'" class="h-3.5 w-3.5 dh-spin" />
          <Trash2 v-else class="h-3.5 w-3.5" />确认清理
        </button>
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
