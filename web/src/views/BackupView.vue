<script setup lang="ts">
import { computed, onMounted, reactive, ref } from 'vue'
import {
  Archive,
  AlertTriangle,
  CheckCircle2,
  ChevronDown,
  ChevronRight,
  Copy,
  Download,
  Eye,
  FileCode2,
  FolderTree,
  HardDrive,
  Layers,
  Loader2,
  RefreshCw,
  RotateCcw,
  Save,
  Trash2,
  Upload,
} from 'lucide-vue-next'
import { api } from '@/api/client'
import type {
  BackupStats,
  BatchResult,
  ContainerView,
  DiffEntry,
  ProjectBackupItem,
  ProjectInfo,
  ProjectRestoreResult,
  RestoreResult,
  Settings,
  SnapshotItem,
} from '@/api/types'
import { formatBytes, relativeTime } from '@/utils/format'
import { useAppStore } from '@/stores/app'
import { useToastStore } from '@/stores/toast'
import EmptyState from '@/components/EmptyState.vue'
import Modal from '@/components/Modal.vue'
import SettingRow from '@/components/SettingRow.vue'
import ToggleSwitch from '@/components/ToggleSwitch.vue'

const toast = useToastStore()
const app = useAppStore()

const tab = ref<'snapshots' | 'projects'>('snapshots')
const backups = ref<SnapshotItem[]>([])
const stats = ref<BackupStats | null>(null)
/** 当前生效的宿主路径映射（启动时自动识别 + DOCKHELM_HOST_ROOTS） */
const pathMappings = computed(() => stats.value?.pathMappings ?? [])
const containers = ref<ContainerView[]>([])
const projects = ref<ProjectInfo[]>([])
/** 「一个项目的备份里有什么」那张卡里给几个「看一眼」的按钮。 */
const p0Readable = computed(() => projects.value.find((p) => p.readable?.length)?.readable ?? [])
const loading = ref(true)
const busy = ref('')
const snapshotTarget = ref('')

/** 快照列表的两种看法：按批次（一次备份动作一行）/ 按容器（一份快照一行）。 */
const view = ref<'batch' | 'container'>('batch')
const filter = ref('')
/** 展开了明细的批次（键是批次时间戳）。默认全部折叠。 */
const expanded = reactive<Record<string, boolean>>({})
/** 批次里超过 6 个容器时，「展开全部」才把剩下的列出来。 */
const batchFull = reactive<Record<string, boolean>>({})

const removeTarget = ref<SnapshotItem | null>(null)
const removeBatch = ref<Batch | null>(null)
const removing = ref(false)

const restoreTarget = ref<SnapshotItem | null>(null)
const diff = ref<DiffEntry[]>([])
const diffLoading = ref(false)
const restoreResult = ref<RestoreResult | null>(null)
const restoring = ref(false)
const showRestore = ref(false)

const batchTarget = ref<Batch | null>(null)
const batchRunning = ref(false)
const batchProgress = ref(0)
const batchResults = ref<{ container: string; ok: boolean; message: string }[]>([])

const projectRestoreTarget = ref<ProjectInfo | null>(null)
const projectChosen = ref('')
const projectRestoring = ref(false)
const projectResult = ref<ProjectRestoreResult | null>(null)
const mountHelpTarget = ref<ProjectInfo | null>(null)

const showFile = ref<{ path: string; content: string } | null>(null)

/** 保留策略：这些值会持久化到设置里，「每天自动清理」与手动清理用同一套。 */
const policy = ref<Settings | null>(null)
const savingPolicy = ref(false)
/** 「清理过期快照」确认框 —— 这一步会真删文件，必须二次确认。 */
const confirmPrune = ref(false)

// ---------- 批次聚合 ----------

/** 一次备份动作 = 一批共用同一个时间戳的快照。 */
interface Batch {
  ts: string
  items: SnapshotItem[]
  size: number
  reason: string
  created: string
}

/**
 * 列表按**批**聚合，而不是一份快照一行。
 *
 * 「立即备份全部容器」会给这一批的每份快照写同一个时间戳，所以按 ts 分组就还原出了
 * 「一次备份动作」。单个容器的备份天然各占一批（只有一份），于是显示成「单个」。
 */
const batches = computed<Batch[]>(() => {
  const map = new Map<string, SnapshotItem[]>()
  for (const it of backups.value) {
    const arr = map.get(it.ts)
    if (arr) arr.push(it)
    else map.set(it.ts, [it])
  }
  const out: Batch[] = []
  for (const [ts, items] of map) {
    const sorted = [...items].sort((a, b) => a.container.localeCompare(b.container))
    const head = sorted[0]
    out.push({
      ts,
      items: sorted,
      size: sorted.reduce((a, i) => a + i.size, 0),
      reason: head?.reason ?? 'manual',
      created: sorted.reduce((a, i) => (i.created > a ? i.created : a), head?.created ?? ''),
    })
  }
  return out.sort((a, b) => b.ts.localeCompare(a.ts))
})

/** 批次里的第一份快照。批次至少含一份，这里给个安全兜底免得模板里到处判空。 */
function firstOf(b: Batch): SnapshotItem {
  const first = b.items[0]
  if (first) return first
  return {
    container: '',
    ts: b.ts,
    path: '',
    size: 0,
    image: '',
    running: false,
    created: b.created,
    reason: b.reason,
  }
}

/** 项目最近一次备份（没有就返回 null）。 */
function latestOf(p: ProjectInfo): ProjectBackupItem | null {
  return p.backups?.[0] ?? null
}

const flatRows = computed(() => {
  const kw = filter.value.trim().toLowerCase()
  const all = [...backups.value].sort((a, b) => b.ts.localeCompare(a.ts))
  return kw ? all.filter((i) => i.container.toLowerCase().includes(kw)) : all
})

const pageSub = computed(() => {
  if (tab.value === 'projects') {
    const visible = projects.value.filter((p) => (p.readable?.length ?? 0) > 0).length
    const latest = projects.value
      .map((p) => p.backups?.[0]?.ts ?? '')
      .sort()
      .pop()
    return `${projects.value.length} 个 compose 项目 · ${visible} 个可备份${
      latest ? ' · 最近备份 ' + formatStamp(latest) : ''
    }`
  }
  const latest = batches.value[0]?.ts
  return `快照 ${stats.value?.snapshots ?? 0} 份 · 覆盖 ${stats.value?.containers ?? 0} 个容器 · 占用 ${formatBytes(
    stats.value?.sizeBytes,
  )}${latest ? ' · 最近一次 ' + formatStamp(latest) : ''}`
})

/** 时间戳（20060102-150405）→ 「今天 12:03」这种给人看的形式。 */
function formatStamp(ts: string): string {
  const m = /^(\d{4})(\d{2})(\d{2})-(\d{2})(\d{2})(\d{2})/.exec(ts)
  if (!m) return ts
  const t = new Date(Number(m[1]), Number(m[2]) - 1, Number(m[3]), Number(m[4]), Number(m[5]), Number(m[6]))
  const now = new Date()
  const hm = `${m[4]}:${m[5]}`
  const sameDay = (a: Date, b: Date) =>
    a.getFullYear() === b.getFullYear() && a.getMonth() === b.getMonth() && a.getDate() === b.getDate()
  if (sameDay(t, now)) return `今天 ${hm}`
  const y = new Date(now.getTime() - 86400000)
  if (sameDay(t, y)) return `昨天 ${hm}`
  return `${m[2]}-${m[3]} ${hm}`
}

function isFullBatch(b: Batch) {
  return b.items.length > 1
}

function toggleBatch(ts: string) {
  expanded[ts] = !expanded[ts]
}

/** 批次明细：默认只列前 6 个，多的收进「展开全部」。 */
function batchItems(b: Batch): SnapshotItem[] {
  if (batchFull[b.ts] || b.items.length <= 6) return b.items
  return b.items.slice(0, 6)
}

// ---------- 数据加载 ----------

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

async function loadPolicy() {
  try {
    policy.value = await api.get<Settings>('/api/settings')
  } catch {
    policy.value = null
  }
}

// ---------- 容器配置快照 ----------

async function snapshotOne() {
  const name = snapshotTarget.value
  if (!name) return
  snapshotTarget.value = ''
  busy.value = 'snapshot'
  try {
    const item = await api.post<SnapshotItem>('/api/backups/snapshot', { container: name, reason: 'manual' })
    toast.success(`已备份 ${name}`, `快照 ${item.ts} · ${formatBytes(item.size)}`)
    await load()
  } catch (e) {
    toast.error('备份失败', e instanceof Error ? e.message : String(e))
  } finally {
    busy.value = ''
  }
}

async function snapshotAll() {
  busy.value = 'snapshot-all'
  try {
    const res = await api.post<BatchResult>('/api/backups/snapshot-all', { reason: 'manual' })
    const detail = `${res.items.length} 个容器成功 · ${formatBytes(
      res.items.reduce((a, i) => a + i.size, 0),
    )}`
    if (res.failed?.length) {
      toast.info(`已备份 ${res.items.length}/${res.total} 个容器`, `${detail}；${res.failed.length} 个失败`)
    } else {
      toast.success(`已备份全部 ${res.items.length} 个容器`, detail)
    }
    expanded[res.ts] = true // 刚拍的那一批直接展开，省得用户再点一下
    await load()
  } catch (e) {
    toast.error('全量备份失败', e instanceof Error ? e.message : String(e))
  } finally {
    busy.value = ''
  }
}

/** 还原接口失败（HTTP 5xx）时，后端把完整结果放在响应体里 —— 别丢掉真正的失败原因。 */
function restoreFromError(e: unknown): RestoreResult | null {
  const p = (e as { payload?: unknown }).payload
  if (p && typeof p === 'object' && 'container' in p && 'steps' in p) return p as RestoreResult
  return null
}

/**
 * 差异请求的序号。快速连点两行不同快照的「还原…」时，先发的请求可能后返回，
 * 把 A 快照的差异盖到 B 快照的弹窗上 —— 而用户正是看着这份差异去决定要不要还原的。
 */
let diffReq = 0

async function openRestore(item: SnapshotItem) {
  const req = ++diffReq
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
    if (req !== diffReq) return // 已经有更新的目标了，这次结果作废
    diff.value = res.diff ?? []
  } catch (e) {
    if (req !== diffReq) return
    toast.error('无法比对差异', e instanceof Error ? e.message : String(e))
  } finally {
    if (req === diffReq) diffLoading.value = false
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
    const res = restoreFromError(e)
    if (res) {
      restoreResult.value = res
      toast.error('还原未成功', res.message)
    } else {
      toast.error('还原失败', e instanceof Error ? e.message : String(e))
    }
    await load()
  } finally {
    restoring.value = false
  }
}

// ---------- 整批还原 ----------

function openBatchRestore(b: Batch) {
  batchTarget.value = b
  batchResults.value = []
  batchProgress.value = 0
}

/**
 * 整批还原：一个容器接一个，绝不并发。
 *
 * 并发重建 24 个容器会把守护进程和 NAS 磁盘一起打满，而且失败时根本分不清
 * 是哪一步出的问题。串行慢一点，但每一步的结果都能如实列给用户。
 */
async function runBatchRestore() {
  const b = batchTarget.value
  if (!b || batchRunning.value) return
  batchRunning.value = true
  batchResults.value = b.items.map((i) => ({ container: i.container, ok: false, message: '待还原' }))
  batchProgress.value = 0
  for (let i = 0; i < b.items.length; i++) {
    const it = b.items[i]
    if (!it) continue
    batchResults.value[i] = { container: it.container, ok: false, message: '还原中…' }
    try {
      const res = await api.post<RestoreResult>('/api/backups/restore', {
        container: it.container,
        ts: it.ts,
        preSnapshot: true,
        keepBackupContainer: true,
      })
      batchResults.value[i] = { container: it.container, ok: res.ok, message: res.message }
    } catch (e) {
      const res = restoreFromError(e)
      batchResults.value[i] = {
        container: it.container,
        ok: false,
        message: res?.message ?? (e instanceof Error ? e.message : String(e)),
      }
    }
    batchProgress.value = i + 1
  }
  batchRunning.value = false
  await load()
  const bad = batchResults.value.filter((r) => !r.ok).length
  if (bad) toast.error(`本批有 ${bad} 个容器没还原成功`, '原因见列表')
  else toast.success(`本批 ${batchProgress.value} 个容器已全部还原`)
}

// ---------- 导出 / 删除 / 清理 ----------

async function downloadFile(url: string, filename: string) {
  const res = await fetch(url, { credentials: 'same-origin' })
  if (!res.ok) {
    let msg = `下载失败（HTTP ${res.status}）`
    try {
      const j = (await res.json()) as { error?: string }
      if (j?.error) msg = j.error
    } catch {
      /* 不是 json 就用默认文案 */
    }
    throw new Error(msg)
  }
  const blob = await res.blob()
  const a = document.createElement('a')
  a.href = URL.createObjectURL(blob)
  a.download = filename
  document.body.appendChild(a)
  a.click()
  a.remove()
  URL.revokeObjectURL(a.href)
}

async function exportSnapshot(item: SnapshotItem) {
  try {
    await downloadFile(
      `/api/backups/export?container=${encodeURIComponent(item.container)}&ts=${encodeURIComponent(item.ts)}`,
      `${item.container}-${item.ts}.json`,
    )
  } catch (e) {
    toast.error('导出失败', e instanceof Error ? e.message : String(e))
  }
}

/** 导出：单个批次就是那份 json；整批导出打包成一个 zip（一次动作一个文件，别让浏览器弹 N 次下载）。 */
async function batchExport(b: Batch) {
  if (!isFullBatch(b)) {
    await exportSnapshot(firstOf(b))
    return
  }
  try {
    await downloadFile(`/api/backups/export-batch?ts=${encodeURIComponent(b.ts)}`, `snapshots-${b.ts}.zip`)
    toast.success('已导出整批快照', `共 ${b.items.length} 个容器`)
  } catch (e) {
    toast.error('导出失败', e instanceof Error ? e.message : String(e))
  }
}

async function confirmRemoveBatch() {
  const b = removeBatch.value
  if (!b || removing.value) return
  removing.value = true
  let done = 0
  try {
    for (const it of b.items) {
      try {
        await api.del(`/api/backups/${encodeURIComponent(it.container)}/${encodeURIComponent(it.ts)}`)
        done++
      } catch {
        /* 单个失败不停：剩下的照样删，最后按实际数量汇报 */
      }
    }
    toast.success(`已删除本批 ${done}/${b.items.length} 份快照`)
    removeBatch.value = null
    await load()
  } finally {
    removing.value = false
  }
}

/** 整批还原里每一行的状态徽标 / 文案。 */
function batchTone(r: { ok: boolean; message: string }) {
  if (r.ok) return 'dh-badge-run'
  if (r.message === '待还原' || r.message === '还原中…') return 'dh-badge-plain'
  return 'dh-badge-err'
}

function batchLabel(r: { ok: boolean; message: string }) {
  if (r.ok) return '成功'
  if (r.message === '待还原') return '待还原'
  if (r.message === '还原中…') return '进行中'
  return '失败'
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
    toast.success(`清理了 ${res.removed} 份备份`, `释放 ${formatBytes(res.freedBytes)}`)
    await load()
  } catch (e) {
    toast.error('清理失败', e instanceof Error ? e.message : String(e))
  } finally {
    busy.value = ''
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
    void app.loadSettings()
    toast.success('保留策略已保存', '每天自动清理与「清理过期」都会按它执行')
  } catch (e) {
    toast.error('保存失败', e instanceof Error ? e.message : String(e))
  } finally {
    savingPolicy.value = false
  }
}

// ---------- compose 项目 ----------

async function backupProject(p: ProjectInfo) {
  busy.value = 'proj:' + p.project
  try {
    const item = await api.post<ProjectBackupItem>('/api/backups/projects/snapshot', { project: p.project })
    toast.success(`已备份项目 ${p.project}`, `${(item.files ?? []).join('、')} · ${formatBytes(item.size)}`)
    await loadProjects()
    await load()
  } catch (e) {
    toast.error('项目备份失败', e instanceof Error ? e.message : String(e))
  } finally {
    busy.value = ''
  }
}

async function backupAllProjects() {
  busy.value = 'proj-all'
  try {
    const res = await api.post<{ items: ProjectBackupItem[]; failed: { container: string; error: string }[] }>(
      '/api/backups/projects/snapshot-all',
      {},
    )
    if (res.failed?.length) {
      toast.info(`已备份 ${res.items.length} 个项目`, `${res.failed.length} 个跳过（文件看不见）`)
    } else {
      toast.success(`已备份全部 ${res.items.length} 个项目`)
    }
    await loadProjects()
    await load()
  } catch (e) {
    toast.error('批量备份项目失败', e instanceof Error ? e.message : String(e))
  } finally {
    busy.value = ''
  }
}

async function downloadProjectBackup(p: ProjectInfo, b?: ProjectBackupItem) {
  const target = b ?? p.backups?.[0]
  if (!target) return
  const files = target.files ?? []
  const single = files.length === 1
  const only = files[0]
  const q = new URLSearchParams({ project: p.project, ts: target.ts })
  if (single && only) q.set('file', only)
  const name = single && only ? only : `${p.project}-${target.ts}.zip`
  try {
    await downloadFile(`/api/backups/projects/download?${q.toString()}`, name)
    toast.success(single ? '已下载 yaml' : '已下载项目备份包', name)
  } catch (e) {
    toast.error('下载失败', e instanceof Error ? e.message : String(e))
  }
}

function openProjectRestore(p: ProjectInfo) {
  projectRestoreTarget.value = p
  projectChosen.value = p.backups?.[0]?.ts ?? ''
  projectResult.value = null
}

async function doProjectRestore() {
  const p = projectRestoreTarget.value
  if (!p || !projectChosen.value) return
  projectRestoring.value = true
  try {
    const res = await api.post<ProjectRestoreResult>('/api/backups/projects/restore', {
      project: p.project,
      ts: projectChosen.value,
      preSnapshot: true,
    })
    projectResult.value = res
    if (res.ok) toast.success('项目文件已还原', res.message)
    else toast.error('还原未完成', res.message)
    await loadProjects()
    await load()
  } catch (e) {
    const p2 = (e as { payload?: ProjectRestoreResult }).payload
    if (p2 && typeof p2 === 'object' && 'files' in p2) {
      projectResult.value = p2
      toast.error('还原未完成', p2.message)
    } else {
      toast.error('还原失败', e instanceof Error ? e.message : String(e))
    }
  } finally {
    projectRestoring.value = false
  }
}

async function deleteProjectBackup(p: ProjectInfo, b: ProjectBackupItem) {
  try {
    await api.del(
      `/api/backups/projects/${encodeURIComponent(p.project)}/${encodeURIComponent(b.ts)}`,
    )
    toast.success('已删除这份备份')
    await loadProjects()
    await load()
  } catch (e) {
    toast.error('删除失败', e instanceof Error ? e.message : String(e))
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

/** 给一个「照抄就能用」的挂载行：拿项目自己报出来的宿主路径反推目录。 */
function suggestedMount(p: ProjectInfo): string {
  const host = p.unreadable?.[0] ?? p.configFiles?.[0] ?? p.workDir ?? ''
  const dir = host.includes('/') ? host.slice(0, host.lastIndexOf('/')) : ''
  if (!dir) return `- /volume1/docker/${p.project}:/host/docker/${p.project}`
  return `- ${dir}:/host/docker/${p.project}`
}

async function copyText(s: string) {
  try {
    await navigator.clipboard.writeText(s)
    toast.success('已复制')
  } catch {
    toast.info('复制失败', '请手动选中复制')
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
    case 'restore-pre':
      return '还原前'
    default:
      return '未知来源'
  }
}

function reasonTone(reason: string) {
  switch (reason) {
    case 'pre_update':
      return 'dh-badge-warn'
    case 'restore-pre':
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
    <!-- 页头：主按钮跟着标签页走 -->
    <div class="dh-phead">
      <div class="min-w-0 flex-1">
        <div class="dh-h1">备份与恢复</div>
        <div class="dh-sub">{{ pageSub }}</div>
      </div>
      <div class="flex flex-wrap items-center gap-2">
        <template v-if="tab === 'snapshots'">
          <select
            v-model="snapshotTarget"
            class="dh-select !mb-0 !w-auto !min-w-[170px] !py-[5px] !text-[12px]"
            :disabled="busy === 'snapshot'"
            @change="snapshotOne"
          >
            <option value="">备份单个容器…</option>
            <option v-for="c in containers" :key="c.id" :value="c.name">{{ c.name }}</option>
          </select>
          <button class="dh-btn dh-btn-primary" :disabled="busy === 'snapshot-all'" @click="snapshotAll">
            <Loader2 v-if="busy === 'snapshot-all'" class="h-3.5 w-3.5 dh-spin" />
            <Archive v-else class="h-3.5 w-3.5" />立即备份全部容器
          </button>
        </template>
        <button
          v-else
          class="dh-btn dh-btn-primary"
          :disabled="busy === 'proj-all'"
          @click="backupAllProjects"
        >
          <Loader2 v-if="busy === 'proj-all'" class="h-3.5 w-3.5 dh-spin" />
          <Archive v-else class="h-3.5 w-3.5" />备份全部项目
        </button>
      </div>
    </div>

    <div class="dh-seg w-fit">
      <button :data-on="tab === 'snapshots'" @click="tab = 'snapshots'">容器配置快照</button>
      <button :data-on="tab === 'projects'" @click="tab = 'projects'">compose 项目</button>
    </div>

    <!-- ==================== 容器配置快照 ==================== -->
    <template v-if="tab === 'snapshots'">
      <div class="grid grid-cols-2 gap-3 lg:grid-cols-3">
        <div class="dh-card p-3.5">
          <div class="text-[12px] text-text-4">快照份数</div>
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
          <div class="text-[12px] text-text-4">项目配置备份</div>
          <div class="mt-1.5 text-[20px] font-semibold leading-none">{{ stats?.projectSnapshots ?? 0 }}</div>
          <div class="mt-1 text-[11px] text-text-5">
            compose 项目 yaml · {{ formatBytes(stats?.projectSizeBytes) }}
          </div>
        </div>
      </div>

      <div class="dh-card">
        <div class="dh-card-head">
          <Archive class="h-3.5 w-3.5 text-text-4" />
          <span>快照</span>
          <span class="text-[11.5px] font-normal text-text-5">
            {{ view === 'batch' ? '点击「全量」行展开该批的容器明细' : '一份快照一行' }}
          </span>
          <span class="ml-auto flex flex-wrap items-center gap-2 font-normal">
            <div class="dh-seg">
              <button :data-on="view === 'batch'" @click="view = 'batch'">按批次</button>
              <button :data-on="view === 'container'" @click="view = 'container'">按容器</button>
            </div>
            <input
              v-if="view === 'container'"
              v-model="filter"
              class="dh-input !mb-0 !w-[150px] !py-[5px] !text-[11.5px]"
              placeholder="筛选容器…"
            />
          </span>
        </div>
        <div class="flex flex-wrap items-center gap-2 border-b border-line-1 px-3.5 py-2">
          <button class="dh-btn dh-btn-sm dh-btn-ghost" @click="openImport">
            <Upload class="h-3 w-3" />导入备份包
          </button>
          <button class="dh-btn dh-btn-sm dh-btn-ghost" :disabled="busy === 'prune'" @click="confirmPrune = true">
            <Trash2 class="h-3 w-3" />清理过期
          </button>
          <button class="dh-btn dh-btn-sm dh-btn-ghost ml-auto" :disabled="loading" @click="load">
            <RefreshCw class="h-3 w-3" :class="loading ? 'dh-spin' : ''" />刷新
          </button>
        </div>

        <EmptyState
          v-if="!backups.length"
          :icon="Archive"
          :title="loading ? '正在载入…' : '还没有任何快照'"
          description="更新容器时 Dockhelm 会自动写一份快照（用于失败回滚），也可以点右上角「立即备份全部容器」把当前所有容器的配置各存一份。"
        />

        <!-- 按批次 -->
        <div v-else-if="view === 'batch'" class="overflow-x-auto">
          <table class="dh-table">
            <thead>
              <tr>
                <th>时间</th>
                <th>包含</th>
                <th class="w-[90px]">来源</th>
                <th class="w-[90px]">大小</th>
                <th class="w-[210px]" />
              </tr>
            </thead>
            <tbody>
              <template v-for="b in batches" :key="b.ts">
                <tr class="cursor-pointer" @click="toggleBatch(b.ts)">
                  <td>
                    <div class="flex items-center gap-1.5 text-[12.5px] font-medium text-text-1">
                      <component
                        :is="expanded[b.ts] ? ChevronDown : ChevronRight"
                        class="h-3.5 w-3.5 flex-none text-text-5"
                      />
                      {{ formatStamp(b.ts) }}
                    </div>
                    <div class="ml-5 text-[10.5px] text-text-6">{{ relativeTime(b.created) }}</div>
                  </td>
                  <td>
                    <span class="dh-badge" :class="isFullBatch(b) ? 'dh-badge-accent' : 'dh-badge-plain'">
                      {{ isFullBatch(b) ? '全量' : '单个' }}
                    </span>
                    <span class="ml-1.5 text-[11.5px] text-text-4">
                      {{ isFullBatch(b) ? `${b.items.length} 个容器` : firstOf(b).container }}
                    </span>
                  </td>
                  <td>
                    <span class="dh-badge" :class="reasonTone(b.reason)">{{ reasonLabel(b.reason) }}</span>
                  </td>
                  <td class="font-mono text-[11.5px] text-text-4">{{ formatBytes(b.size) }}</td>
                  <td @click.stop>
                    <div class="flex justify-end gap-1.5">
                      <button
                        v-if="isFullBatch(b)"
                        class="dh-btn dh-btn-sm"
                        @click="openBatchRestore(b)"
                      >
                        <Layers class="h-3 w-3" />还原本批
                      </button>
                      <button v-else class="dh-btn dh-btn-sm" @click="openRestore(firstOf(b))">
                        <RotateCcw class="h-3 w-3" />还原…
                      </button>
                      <button class="dh-btn dh-btn-sm" @click="batchExport(b)">
                        <Download class="h-3 w-3" />导出
                      </button>
                      <button
                        class="dh-btn dh-btn-sm dh-btn-danger"
                        @click="isFullBatch(b) ? (removeBatch = b) : (removeTarget = firstOf(b))"
                      >
                        <Trash2 class="h-3 w-3" />
                      </button>
                    </div>
                  </td>
                </tr>
                <tr v-if="expanded[b.ts]">
                  <td colspan="5" class="!py-2.5">
                    <div class="flex flex-col gap-1.5">
                      <div class="text-[11px] text-text-5">
                        这一批里 {{ b.items.length }} 个容器 · 每份都是完整配置，可单独还原
                      </div>
                      <div class="grid gap-1.5 sm:grid-cols-2 lg:grid-cols-3">
                        <div
                          v-for="it in batchItems(b)"
                          :key="it.container"
                          class="flex items-center gap-2 rounded-[9px] border border-line-1 bg-ink-800 px-2.5 py-1.5"
                        >
                          <span class="min-w-0 flex-1 truncate text-[11.5px] text-text-3" :title="it.container">
                            {{ it.container }}
                          </span>
                          <span class="dh-badge dh-badge-plain">{{ formatBytes(it.size) }}</span>
                          <button class="dh-btn dh-btn-sm" @click="openRestore(it)">
                            <RotateCcw class="h-3 w-3" />还原
                          </button>
                          <button
                            class="dh-btn dh-btn-sm dh-btn-danger"
                            title="只删这一份"
                            @click="removeTarget = it"
                          >
                            <Trash2 class="h-3 w-3" />
                          </button>
                        </div>
                      </div>
                      <button
                        v-if="b.items.length > 6 && !batchFull[b.ts]"
                        class="dh-btn dh-btn-sm dh-btn-ghost w-fit"
                        @click="batchFull[b.ts] = true"
                      >
                        展开全部（还有 {{ b.items.length - 6 }} 个）
                      </button>
                    </div>
                  </td>
                </tr>
              </template>
            </tbody>
          </table>
        </div>

        <!-- 按容器 -->
        <div v-else class="overflow-x-auto">
          <table class="dh-table">
            <thead>
              <tr>
                <th>容器</th>
                <th class="w-[140px]">快照时间</th>
                <th class="w-[80px]">来源</th>
                <th class="w-[80px]">大小</th>
                <th class="w-[100px]">快照时状态</th>
                <th class="w-[220px]" />
              </tr>
            </thead>
            <tbody>
              <tr v-for="item in flatRows" :key="item.container + '/' + item.ts">
                <td>
                  <div class="text-[12.5px] font-medium text-text-1">{{ item.container }}</div>
                  <div class="max-w-[220px] truncate font-mono text-[10.5px] text-text-6" :title="item.image">
                    {{ item.image || '—' }}
                  </div>
                </td>
                <td>
                  <div class="text-[12px] text-text-3">{{ formatStamp(item.ts) }}</div>
                  <div class="text-[10.5px] text-text-6">{{ relativeTime(item.created) }}</div>
                </td>
                <td>
                  <span class="dh-badge" :class="reasonTone(item.reason)">{{ reasonLabel(item.reason) }}</span>
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
                    <button class="dh-btn dh-btn-sm" @click="exportSnapshot(item)">
                      <Download class="h-3 w-3" />导出
                    </button>
                    <button class="dh-btn dh-btn-sm dh-btn-danger" @click="removeTarget = item">
                      <Trash2 class="h-3 w-3" />
                    </button>
                  </div>
                </td>
              </tr>
              <tr v-if="!flatRows.length">
                <td colspan="6" class="text-[12px] text-text-4">没有匹配的快照。</td>
              </tr>
            </tbody>
          </table>
        </div>
      </div>

      <div class="grid gap-3 lg:grid-cols-2">
        <!-- 保留策略 -->
        <div v-if="policy" class="dh-card">
          <div class="dh-card-head">
            <Trash2 class="h-3.5 w-3.5 text-text-4" />
            <span>保留策略</span>
            <button class="dh-btn dh-btn-sm dh-btn-primary ml-auto" :disabled="savingPolicy" @click="savePolicy">
              <Loader2 v-if="savingPolicy" class="h-3 w-3 dh-spin" />
              <Save v-else class="h-3 w-3" />保存
            </button>
          </div>
          <div class="flex flex-col gap-3 p-3.5">
            <SettingRow title="每容器保留最近份数" sub="超出后按时间滚动覆盖（更新前快照不计入这个额度）">
              <select
                v-model.number="policy.backupKeepPerContainer"
                class="dh-select !mb-0 !w-[110px] !py-[5px] !text-[11.5px]"
              >
                <option :value="0">不限制</option>
                <option :value="5">5 份</option>
                <option :value="10">10 份</option>
                <option :value="20">20 份</option>
                <option :value="50">50 份</option>
              </select>
            </SettingRow>
            <SettingRow title="快照保留期" sub="到期自动清理">
              <select
                v-model.number="policy.backupMaxAgeDays"
                class="dh-select !mb-0 !w-[110px] !py-[5px] !text-[11.5px]"
              >
                <option :value="0">不限制</option>
                <option :value="7">7 天</option>
                <option :value="30">30 天</option>
                <option :value="90">90 天</option>
                <option :value="365">365 天</option>
              </select>
            </SettingRow>
            <SettingRow title="快照总容量上限" sub="达到上限后先清最旧的快照">
              <select
                v-model.number="policy.backupMaxTotalMB"
                class="dh-select !mb-0 !w-[110px] !py-[5px] !text-[11.5px]"
              >
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

        <!-- 存储位置（路径映射保留在这里：项目备份能不能读到 yaml 全看它） -->
        <div class="dh-card">
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
            <div class="border-t border-line-1 pt-3">
              <div class="mb-2 flex items-center gap-2">
                <FolderTree class="h-3.5 w-3.5 text-text-4" />
                <span class="text-[12px] text-text-4">宿主机路径映射（启动时从自身容器自动识别）</span>
              </div>
              <div v-if="!pathMappings.length" class="text-[11.5px] leading-relaxed text-text-5">
                没有识别到任何挂载映射 —— Dockhelm 读不到宿主机的 docker 目录，
                compose 项目那份 yaml 也就无从备份。在 compose 里挂一行即可，例如
                <code class="font-mono text-accent">/volume1/docker:/host/docker</code>
                （右边叫什么名字都行，启动时会自动识别）。
              </div>
              <div v-else class="flex flex-col gap-1.5">
                <div
                  v-for="(m, i) in pathMappings"
                  :key="i"
                  class="flex flex-wrap items-center gap-2 text-[11.5px]"
                >
                  <span class="dh-badge" :class="m.visible ? 'dh-badge-run' : 'dh-badge-warn'">
                    {{ m.visible ? '可见' : '不可见' }}
                  </span>
                  <code class="font-mono text-text-3">{{ m.host }}</code>
                  <span class="text-text-6">→</span>
                  <code class="font-mono text-text-3">{{ m.container }}</code>
                  <span v-if="m.source === 'env'" class="text-[11px] text-text-6">来自 DOCKHELM_HOST_ROOTS</span>
                </div>
              </div>
            </div>
          </div>
        </div>
      </div>
    </template>

    <!-- ==================== compose 项目 ==================== -->
    <template v-else>
      <div class="dh-card">
        <div class="dh-card-head">
          <FolderTree class="h-3.5 w-3.5 text-text-4" />
          <span>项目</span>
          <span class="ml-auto flex items-center gap-2 font-normal">
            <span class="text-[11.5px] text-text-5">按项目名排序</span>
            <button class="dh-btn dh-btn-sm dh-btn-ghost" :disabled="loading" @click="loadProjects">
              <RefreshCw class="h-3 w-3" :class="loading ? 'dh-spin' : ''" />刷新
            </button>
          </span>
        </div>

        <EmptyState
          v-if="!projects.length"
          :icon="FolderTree"
          :title="loading ? '正在载入…' : '没有发现 compose 项目'"
          description="只有由 docker compose 创建的容器才会带上项目标签。"
        />

        <div v-else class="overflow-x-auto">
          <table class="dh-table">
            <thead>
              <tr>
                <th>项目</th>
                <th>文件</th>
                <th class="w-[150px]">最近备份</th>
                <th class="w-[280px]" />
              </tr>
            </thead>
            <tbody>
              <tr v-for="p in projects" :key="p.project">
                <td>
                  <div class="text-[12.5px] font-medium text-text-1">{{ p.project }}</div>
                  <div class="text-[10.5px] text-text-6">{{ p.containers?.length ?? 0 }} 个容器</div>
                </td>
                <td>
                  <div class="flex flex-wrap gap-1.5">
                    <span
                      v-for="f in p.readable"
                      :key="f.path"
                      class="inline-flex items-center gap-1 rounded-md border border-line-1 bg-ink-800 px-1.5 py-[2px] font-mono text-[10.5px] text-text-4"
                      :title="f.hostPath"
                    >
                      {{ f.name }} ✓
                    </span>
                    <span
                      v-for="u in p.unreadable"
                      :key="u"
                      class="inline-flex items-center gap-1 rounded-md border border-line-warn bg-soft-warn px-1.5 py-[2px] font-mono text-[10.5px] text-warn-text"
                      :title="u"
                    >
                      {{ u.split('/').pop() }} 不可见
                    </span>
                    <span v-if="!p.readable.length && !p.unreadable.length" class="text-[11px] text-text-5">
                      没读到任何 yaml
                    </span>
                  </div>
                </td>
                <td>
                  <template v-if="p.backups?.length">
                    <div class="text-[12px] text-text-3">{{ formatStamp(latestOf(p)?.ts ?? '') }}</div>
                    <div class="text-[10.5px] text-text-6">{{ p.backups.length }} 份历史</div>
                  </template>
                  <span v-else class="text-[11.5px] text-text-5">从未备份</span>
                </td>
                <td>
                  <div class="flex justify-end gap-1.5">
                    <button
                      class="dh-btn dh-btn-sm"
                      :disabled="!p.readable.length || busy === 'proj:' + p.project"
                      @click="backupProject(p)"
                    >
                      <Loader2 v-if="busy === 'proj:' + p.project" class="h-3 w-3 dh-spin" />
                      <Archive v-else class="h-3 w-3" />立即备份
                    </button>
                    <template v-if="p.backups?.length">
                      <button class="dh-btn dh-btn-sm" @click="downloadProjectBackup(p)">
                        <Download class="h-3 w-3" />{{
                          (latestOf(p)?.files?.length ?? 0) === 1 ? '下载 yaml' : '下载 zip'
                        }}
                      </button>
                      <button class="dh-btn dh-btn-sm" @click="openProjectRestore(p)">
                        <RotateCcw class="h-3 w-3" />还原
                      </button>
                    </template>
                    <button v-else-if="!p.readable.length" class="dh-btn dh-btn-sm" @click="mountHelpTarget = p">
                      看挂载方法
                    </button>
                  </div>
                </td>
              </tr>
            </tbody>
          </table>
        </div>
      </div>

      <div class="grid gap-3 lg:grid-cols-2">
        <div class="dh-card">
          <div class="dh-card-head">
            <FileCode2 class="h-3.5 w-3.5 text-text-4" />
            <span>一个项目的备份里有什么</span>
          </div>
          <div class="flex flex-col gap-2.5 p-3.5">
            <div class="flex items-start gap-2 text-[11.5px] leading-relaxed text-text-4">
              <span class="dh-badge dh-badge-plain">1</span>
              <span>项目目录下的 yaml 文本（compose / docker-compose / override）</span>
            </div>
            <div class="flex items-start gap-2 text-[11.5px] leading-relaxed text-text-4">
              <span class="dh-badge dh-badge-plain">2</span>
              <span><code class="font-mono">.env</code>（容器密码、端口这些变量都在这里）</span>
            </div>
            <div class="flex items-start gap-2 text-[11.5px] leading-relaxed text-text-4">
              <span class="dh-badge dh-badge-plain">3</span>
              <span>每个文件的原路径 + 校验值，还原时逐个核对</span>
            </div>
            <div class="border-t border-line-1 pt-2.5 text-[11.5px] leading-relaxed text-text-5">
              不打包项目目录里的其它文件（数据库、缓存、上传的文件）—— 那些属于<b class="text-text-4">数据</b>，
              仍然要你自己备份。
            </div>
            <div v-if="p0Readable.length" class="border-t border-line-1 pt-2.5">
              <div class="mb-1.5 text-[11.5px] text-text-5">随便看一个项目的文件内容：</div>
              <div class="flex flex-wrap gap-1.5">
                <button
                  v-for="f in p0Readable.slice(0, 4)"
                  :key="f.path"
                  class="dh-btn dh-btn-sm"
                  @click="openFile(f)"
                >
                  <Eye class="h-3 w-3" />{{ f.name }}
                </button>
              </div>
            </div>
          </div>
        </div>

        <div class="dh-card">
          <div class="dh-card-head">
            <AlertTriangle class="h-3.5 w-3.5 text-text-4" />
            <span>「不可见」是什么意思</span>
          </div>
          <div class="flex flex-col gap-2.5 p-3.5">
            <div class="text-[11.5px] leading-relaxed text-text-4">
              意思是 Dockhelm 在<b class="text-text-3">自己的容器里</b>读不到这个文件。
              把项目的宿主目录挂进 Dockhelm 容器（右边随便叫什么），重建容器后即可备份：
            </div>
            <div class="flex items-center gap-2 rounded-[9px] border border-line-1 bg-ink-800 px-2.5 py-2">
              <code class="min-w-0 flex-1 truncate font-mono text-[11px] text-text-3">
                - /volume1/docker/&lt;项目&gt;:/host/docker/&lt;项目&gt;
              </code>
              <button
                class="dh-btn dh-btn-sm flex-none"
                @click="copyText('- /volume1/docker/<项目>:/host/docker/<项目>')"
              >
                <Copy class="h-3 w-3" />复制
              </button>
            </div>
            <div class="text-[11.5px] leading-relaxed text-text-5">
              宿主机的 docker 目录整体挂进来最省事，映射会在启动时自动识别。
              它<b class="text-text-4">只影响 compose 项目页</b>能不能读到 yaml ——
              容器快照走的是 Docker 接口，一个目录都不用挂。
            </div>
          </div>
        </div>
      </div>
    </template>

    <!-- ==================== 单容器还原 ==================== -->
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
        <div
          class="flex items-start gap-2.5 rounded-[10px] border border-line-warn bg-soft-warn px-3 py-2.5 text-[12px] leading-relaxed text-warn-text"
        >
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

    <!-- ==================== 整批还原 ==================== -->
    <Modal
      :open="!!batchTarget"
      title="还原本批"
      :subtitle="batchTarget ? `${batchTarget.items.length} 个容器 · 快照 ${batchTarget.ts}` : ''"
      width="620px"
      :busy="batchRunning"
      @close="batchTarget = null"
    >
      <div class="flex flex-col gap-3">
        <div
          class="flex items-start gap-2.5 rounded-[10px] border border-line-warn bg-soft-warn px-3 py-2.5 text-[12px] leading-relaxed text-warn-text"
        >
          <AlertTriangle class="mt-[2px] h-4 w-4 flex-none" />
          <div>
            会把这个批次里的容器<b>逐个</b>按快照重建（串行执行，不并发）。每个容器还原前都会先为它
            <b>当前的状态</b>存一份快照，所以单个失败也能退回去。
          </div>
        </div>

        <div v-if="batchResults.length" class="flex flex-col gap-1.5">
          <div class="flex items-center gap-2 text-[11.5px] text-text-4">
            <span>进度 {{ batchProgress }} / {{ batchTarget?.items.length ?? 0 }}</span>
            <div class="h-[5px] flex-1 overflow-hidden rounded-full bg-ink-800">
              <div
                class="h-full rounded-full bg-accent transition-all"
                :style="{
                  width: (batchProgress / Math.max(1, batchTarget?.items.length ?? 1)) * 100 + '%',
                }"
              />
            </div>
          </div>
          <div class="dh-scroll max-h-[300px] overflow-auto rounded-[10px] border border-line-1">
            <div
              v-for="r in batchResults"
              :key="r.container"
              class="flex items-center gap-2 border-b border-line-row px-3 py-2 last:border-b-0"
            >
              <span class="dh-badge" :class="r.ok ? 'dh-badge-run' : batchTone(r)">
                {{ batchLabel(r) }}
              </span>
              <span class="min-w-0 flex-1 truncate text-[12px] text-text-2">{{ r.container }}</span>
              <span class="max-w-[220px] truncate text-[11px] text-text-5" :title="r.message">{{ r.message }}</span>
            </div>
          </div>
        </div>
      </div>
      <template #footer>
        <button class="dh-btn" :disabled="batchRunning" @click="batchTarget = null">
          {{ batchResults.length ? '关闭' : '取消' }}
        </button>
        <button
          v-if="!batchResults.length || batchProgress < (batchTarget?.items.length ?? 0)"
          class="dh-btn dh-btn-danger"
          :disabled="batchRunning"
          @click="runBatchRestore"
        >
          <Layers class="h-3.5 w-3.5" />{{ batchRunning ? '还原中…' : '开始还原' }}
        </button>
      </template>
    </Modal>

    <!-- ==================== 项目还原 ==================== -->
    <Modal
      :open="!!projectRestoreTarget"
      title="还原项目文件"
      :subtitle="projectRestoreTarget?.project"
      width="660px"
      :busy="projectRestoring"
      @close="projectRestoreTarget = null"
    >
      <div v-if="projectResult" class="flex flex-col gap-3">
        <div
          class="flex items-center gap-2 rounded-[10px] border px-3 py-2.5 text-[12.5px]"
          :class="
            projectResult.ok
              ? 'border-line-ok bg-soft-ok text-run-text'
              : 'border-line-err bg-soft-err text-err-text'
          "
        >
          <component :is="projectResult.ok ? CheckCircle2 : AlertTriangle" class="h-4 w-4 flex-none" />
          {{ projectResult.message }}
        </div>
        <div class="overflow-hidden rounded-[10px] border border-line-1">
          <table class="dh-table">
            <thead>
              <tr>
                <th class="w-[160px]">文件</th>
                <th>写回路径</th>
                <th class="w-[180px]">结果</th>
              </tr>
            </thead>
            <tbody>
              <tr v-for="f in projectResult.files" :key="f.name">
                <td class="font-mono text-[11.5px] text-text-2">{{ f.name }}</td>
                <td class="max-w-[240px] truncate font-mono text-[11px] text-text-5" :title="f.hostPath">
                  {{ f.hostPath }}
                </td>
                <td class="text-[11.5px]" :class="f.written ? 'text-run-text' : 'text-err-text'">{{ f.note }}</td>
              </tr>
            </tbody>
          </table>
        </div>
      </div>

      <div v-else class="flex flex-col gap-3">
        <div
          class="flex items-start gap-2.5 rounded-[10px] border border-line-warn bg-soft-warn px-3 py-2.5 text-[12px] leading-relaxed text-warn-text"
        >
          <AlertTriangle class="mt-[2px] h-4 w-4 flex-none" />
          <div>
            还原会把这些文件<b>写回它们原来的宿主路径</b>，覆盖现在的同名文件。动手前 Dockhelm 会先给
            <b>当前状态</b>存一份项目备份，写错了还能退回来。
          </div>
        </div>

        <div>
          <div class="mb-2 text-[12px] text-text-3">选择要还原的那一份</div>
          <div class="flex flex-col gap-1.5">
            <label
              v-for="b in projectRestoreTarget?.backups ?? []"
              :key="b.ts"
              class="flex cursor-pointer items-center gap-2.5 rounded-[9px] border px-3 py-2"
              :class="
                projectChosen === b.ts ? 'border-line-accent-soft bg-soft-accent' : 'border-line-1 bg-ink-800'
              "
            >
              <input v-model="projectChosen" type="radio" :value="b.ts" class="h-[13px] w-[13px] accent-accent" />
              <span class="min-w-0 flex-1">
                <span class="text-[12.5px] text-text-2">{{ formatStamp(b.ts) }}</span>
                <span class="ml-2 text-[11px] text-text-5">{{ relativeTime(b.created) }}</span>
                <span class="block truncate font-mono text-[10.5px] text-text-6">{{ (b.files ?? []).join('、') }}</span>
              </span>
              <span class="dh-badge dh-badge-plain">{{ formatBytes(b.size) }}</span>
              <button
                class="dh-btn dh-btn-sm dh-btn-danger flex-none"
                :disabled="projectRestoring"
                title="删除这一份备份"
                @click.prevent="projectRestoreTarget && deleteProjectBackup(projectRestoreTarget, b)"
              >
                <Trash2 class="h-3 w-3" />
              </button>
            </label>
          </div>
        </div>

        <div v-if="projectRestoreTarget?.backups?.[0] === undefined" class="text-[12px] text-text-5">
          这个项目还没有任何备份，先去上面点「立即备份」。
        </div>
      </div>

      <template #footer>
        <button class="dh-btn" @click="projectRestoreTarget = null">
          {{ projectResult ? '关闭' : '取消' }}
        </button>
        <button
          v-if="!projectResult"
          class="dh-btn dh-btn-danger"
          :disabled="projectRestoring || !projectChosen"
          @click="doProjectRestore"
        >
          <RotateCcw class="h-3.5 w-3.5" />{{ projectRestoring ? '还原中…' : '确认还原' }}
        </button>
      </template>
    </Modal>

    <!-- ==================== 看挂载方法 ==================== -->
    <Modal
      :open="!!mountHelpTarget"
      title="让 Dockhelm 读到这个项目的文件"
      :subtitle="mountHelpTarget?.project"
      width="640px"
      @close="mountHelpTarget = null"
    >
      <div v-if="mountHelpTarget" class="flex flex-col gap-3">
        <div class="text-[12.5px] leading-relaxed text-text-3">
          这个项目的文件在 Dockhelm 容器里看不见，所以备份不了。在 Dockhelm 的 compose 里加一行挂载即可 ——
          <b class="text-text-2">冒号右边叫什么名字都行</b>，启动时会自动识别。
        </div>
        <div class="flex items-center gap-2 rounded-[10px] border border-line-1 bg-ink-800 px-3 py-2.5">
          <code class="min-w-0 flex-1 break-all font-mono text-[11.5px] text-accent-text">
            {{ suggestedMount(mountHelpTarget) }}
          </code>
          <button class="dh-btn dh-btn-sm flex-none" @click="copyText(suggestedMount(mountHelpTarget))">
            <Copy class="h-3 w-3" />复制
          </button>
        </div>
        <div class="text-[11.5px] leading-relaxed text-text-5">
          改完 <code class="font-mono">docker compose up -d</code> 重建 Dockhelm 容器，
          回到这一页就会变成可备份。宿主机的整个 docker 目录挂进来最省事。
        </div>
      </div>
      <template #footer>
        <button class="dh-btn" @click="mountHelpTarget = null">关闭</button>
      </template>
    </Modal>

    <!-- ==================== 导入备份包 ==================== -->
    <Modal
      :open="showImport"
      title="导入备份包"
      subtitle="粘贴一份快照 JSON"
      width="640px"
      :busy="importing"
      @close="showImport = false"
    >
      <div class="flex flex-col gap-3">
        <div
          class="flex items-start gap-2.5 rounded-[10px] border border-line-1 bg-ink-800 px-3 py-2.5 text-[11.5px] leading-relaxed text-text-4"
        >
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
            <input v-model="importContainer" class="dh-input" placeholder="redis" @change="guessContainer" />
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

    <!-- ==================== 删除快照 ==================== -->
    <Modal
      :open="!!removeTarget"
      title="删除这一份快照"
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

    <!-- 删除整批：一轮备份动作产生的所有快照一起删 -->
    <Modal
      :open="!!removeBatch"
      title="删除这一整批快照"
      :subtitle="removeBatch ? `${removeBatch.items.length} 个容器 · ${removeBatch.ts}` : ''"
      width="440px"
      :busy="removing"
      @close="removeBatch = null"
    >
      <div class="flex flex-col gap-2 text-[12.5px] text-text-3">
        <div>
          会删掉这一批里的 <b class="text-text-2">{{ removeBatch?.items.length ?? 0 }} 份快照</b>（每个容器各一份），
          删掉后它们都无法再用于还原。已经运行中的容器不受影响。
        </div>
        <div class="text-[11.5px] text-text-5">
          只想删其中某一个容器的话，先把这一行展开，在它那一行点「还原」右边的删除按钮。
        </div>
      </div>
      <template #footer>
        <button class="dh-btn" :disabled="removing" @click="removeBatch = null">取消</button>
        <button class="dh-btn dh-btn-danger" :disabled="removing" @click="confirmRemoveBatch">
          <Loader2 v-if="removing" class="h-3.5 w-3.5 dh-spin" />
          <Trash2 v-else class="h-3.5 w-3.5" />{{ removing ? '删除中…' : '删除整批' }}
        </button>
      </template>
    </Modal>

    <!-- 清理过期：会真删文件，必须二次确认 -->
    <Modal :open="confirmPrune" title="清理过期备份" width="460px" :busy="busy === 'prune'" @close="confirmPrune = false">
      <div class="flex flex-col gap-3">
        <div
          class="flex items-start gap-2.5 rounded-[10px] border border-line-warn bg-soft-warn px-3 py-2.5 text-[12.5px] leading-relaxed text-warn-text"
        >
          <AlertTriangle class="mt-[2px] h-4 w-4 flex-none" />
          <div>
            会按下面的保留策略<b>直接删除备份文件</b>，<b>不可恢复</b>。容器快照与项目备份都会被清到，
            「更新前快照永不自动清理」打开时，那部分不会被碰到。
          </div>
        </div>
        <div class="rounded-[10px] border border-line-1 bg-ink-800 px-3 py-2.5 text-[12px] leading-relaxed text-text-3">
          当前策略：每容器保留最近
          <b>{{ policy?.backupKeepPerContainer ? policy.backupKeepPerContainer + ' 份' : '不限' }}</b>
          · 保留 {{ policy?.backupMaxAgeDays ? policy.backupMaxAgeDays + ' 天' : '不限' }}
          · 总容量上限 {{ policy?.backupMaxTotalMB ? policy.backupMaxTotalMB + ' MB' : '不限' }}
          <div class="mt-1 text-[11.5px] text-text-5">要改策略请在上面的「保留策略」卡里改并保存。</div>
        </div>
        <div class="text-[12px] text-text-4">
          当前共 {{ backups.length }} 份容器快照（{{ formatBytes(stats?.sizeBytes) }}）+ {{ stats?.projectSnapshots ?? 0 }} 份项目备份（{{ formatBytes(stats?.projectSizeBytes) }}）。
        </div>
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
    <Modal :open="!!showFile" title="文件内容" :subtitle="showFile?.path" width="720px" @close="showFile = null">
      <pre
        class="dh-scroll max-h-[520px] overflow-auto whitespace-pre-wrap break-all rounded-[10px] border border-line-1 bg-ink-800 p-3 font-mono text-[11.5px] leading-[1.7] text-text-2"
      >{{ showFile?.content }}</pre>
    </Modal>
  </div>
</template>
