/** 与后端一致的接口类型定义。 */

export type CheckStatus = 'up_to_date' | 'update_available' | 'unknown' | 'no_upstream'
export type ResultStatus =
  | 'up_to_date'
  | 'updated'
  | 'failed'
  | 'broken'
  | 'skipped'
  | 'pulling'
  | 'stopped'
  | 'created'

export interface AboutInfo {
  name: string
  nameCN: string
  version: string
  tagline: string
  description: string
  author: string
  authorURL: string
  repoURL: string
  license: string
  commit: string
  commitShort: string
  buildTime: string
  goVersion: string
  platform: string
  startedAt: string
}

export interface AboutResponse {
  about: AboutInfo
  runtime: Record<string, unknown>
}

export interface SessionInfo {
  initialized: boolean
  loggedIn: boolean
  failures: number
  maxFailures: number
  /** 因连续失败被锁定时，还要等多少秒才能再试；未锁定为 0。 */
  lockedFor: number
  /** 锁定提示文案（含「约 N 分钟」），未锁定时为空串。措辞由服务端统一生成。 */
  lockedHint: string
  minPassword: number
  sessionCount: number
}

/** 一条端口映射。published=false 表示只在容器网络内可见，没发布到宿主机。 */
export interface PortView {
  hostIp: string
  hostPort: number
  innerPort: number
  proto: string
  published: boolean
}

export interface ContainerView {
  name: string
  id: string
  shortId: string
  image: string
  imageId: string
  state: string
  status: string
  created: number
  ports: string[]
  portList: PortView[]
  project: string
  labels: Record<string, string>
  hasUpdate: boolean
  updateKnown: boolean
  self: boolean
  excluded: boolean
}

export interface OverviewResponse {
  version: { version: string; name: string; commit: string }
  containers: { total: number; running: number; stopped: number; paused: number; unhealthy: number }
  images: { total: number; sizeBytes: number; reclaimable: number }
  updates: {
    available: number
    unknown: number
    checkedAt: string
    items: { container: string; image: string; reason: string; localDigest: string; remoteDigest: string }[]
  }
  /** 所有运行中容器的合计占用（容器数为 0 时后端不返回这个字段）。 */
  usage?: { cpuPercent: number; memUsed: number; memTotal: number }
  disk: { free: number; total: number }
  notify: { sentToday: number }
  recent: RunLog[]
  excluded: string[]
  self: string
  docker?: {
    /** 宿主机主机名，顶栏的「Docker 控制台 · xxx」用它。 */
    name: string
    version: string
    os: string
    arch: string
    kernel: string
    cpus: number
    memTotal: number
    rootDir: string
    mirrors: string[]
  }
  dockerError?: string
}

export interface RunLog {
  id: number
  ts: string
  kind: string
  ref: string
  status: string
  message: string
  detail: string
}

export interface CheckResult {
  container: string
  id: string
  image: string
  localDigest: string
  remoteDigest: string
  status: CheckStatus
  reason: string
  checkedAt: string
}

export interface UpdatesResponse {
  results: CheckResult[]
  checkedAt: string
  summary: { updateAvailable: number; upToDate: number; unknown: number; noUpstream: number }
  excluded: string[]
  selfName: string
}

export interface Schedule {
  id: number
  name: string
  cron: string
  action: string
  targets: string[]
  enabled: boolean
  lastRun: string
  lastStatus: string
  lastMessage: string
  createdAt: string
  nextRun?: string
}

export interface ScheduleAction {
  key: string
  label: string
  description: string
  needsTargets: boolean
}

export interface MirrorConfig {
  url: string
  note: string
  enabled: boolean
  lastTested?: string
  latencyMs?: number
  ok?: boolean
  err?: string
  builtin?: boolean
}

export interface RegistrySettings {
  mirrors: MirrorConfig[]
  pullMirror: string
  insecure: boolean
}

export interface RegistriesResponse {
  settings: RegistrySettings
  /** 预置的常用加速源清单。已随首次启动写进 settings.mirrors，
   *  这里只用于「列表被删空后一键找回」，不再渲染成独立栏位。 */
  presets: MirrorConfig[]
  daemonMirrors: string[]
  daemonError: string
  snippet: string
  explain: string
}

export interface SnapshotItem {
  container: string
  ts: string
  path: string
  size: number
  image: string
  running: boolean
  created: string
  /** 快照来源：manual / pre_update / scheduled / restore-pre。 */
  reason: string
}

/** 一次批量备份（「立即备份全部容器」）的结果。同一批的快照共用 ts。 */
export interface BatchResult {
  ts: string
  total: number
  items: SnapshotItem[]
  failed: { container: string; error: string }[]
}

export interface DiffEntry {
  field: string
  snapshot: string
  current: string
}

export interface BackupStats {
  snapshots: number
  sizeBytes: number
  containers: number
  dir: string
  /** 项目（compose yaml）备份的份数与占用，与容器快照分开计。 */
  projectSnapshots: number
  projectSizeBytes: number
  volumeRootMounted: boolean
  dockerRootVisible: boolean
  dockerRoot: string
  pathMappings: PathMapping[]
}

/** 一条「宿主机路径 → 容器内路径」映射。 */
export interface PathMapping {
  host: string
  container: string
  /** auto = 启动时从自身容器挂载自动识别；env = DOCKHELM_HOST_ROOTS 显式声明 */
  source: 'auto' | 'env'
  /** 容器内这个路径当前是否真的存在 */
  visible: boolean
}

export interface ProjectFile {
  name: string
  hostPath: string
  path: string
  size: number
}

/** 一份项目（compose yaml）备份。 */
export interface ProjectBackupItem {
  project: string
  ts: string
  created: string
  size: number
  files: string[]
  reason: string
}

export interface ProjectInfo {
  project: string
  containers: string[]
  configFiles: string[]
  readable: ProjectFile[]
  unreadable: string[]
  workDir: string
  /** 这个项目的 yaml 备份历史（新 → 旧）。 */
  backups: ProjectBackupItem[]
}

/** 项目还原的结果（逐文件）。 */
export interface ProjectRestoreResult {
  project: string
  ts: string
  files: { name: string; hostPath: string; written: boolean; note: string }[]
  ok: boolean
  message: string
}

export interface RestoreResult {
  container: string
  snapshot: string
  image: string
  steps: string[]
  ok: boolean
  message: string
}

export interface Channel {
  id: number
  name: string
  type: string
  enabled: boolean
  config: Record<string, unknown>
  createdAt: string
}

export interface ChannelField {
  key: string
  label: string
  type: string
  required: boolean
  placeholder: string
  help: string
}

export interface ChannelPreset {
  type: string
  label: string
  description: string
  fields: ChannelField[]
  method: string
  allowCustomHeaders: boolean
}

export interface EventDef {
  event: string
  label: string
  group: string
  level: 'urgent' | 'normal'
  default: boolean
  description: string
}

export interface NotifyEventRow {
  event: string
  enabled: boolean
  level: string
}

export interface NotifySettings {
  enabled: boolean
  quietEnabled: boolean
  quietStart: string
  quietEnd: string
  quietNormalMode: 'digest' | 'drop'
  quietUrgentSend: boolean
  dedupeWindow: number
  dailyLimit: number
  panelURL: string
  sentToday: number
  inQuietHours: boolean
}

export interface NotifyRecord {
  id: number
  ts: string
  event: string
  level: string
  title: string
  body: string
  ok: boolean
  errmsg: string
  channel: string
}

export interface Settings {
  exclude: string[]
  panelURL: string
  concurrency: number
  logRetention: number
  checkOnStart: boolean

  // —— 检测 ——
  /** 周期性自动巡检间隔（小时），0 = 关闭。 */
  checkIntervalHours: number
  /** 巡检发现新版本时推一条通知。 */
  notifyOnCheck: boolean

  // —— 更新策略 ——
  /** 自动更新总开关。关闭时定时任务只检测、绝不动容器。 */
  autoApply: boolean
  /** 同一轮批量更新里，同一个镜像只下载一次。 */
  pullOnce: boolean
  /** 重建容器前先写一份配置快照。 */
  backupBefore: boolean
  /** 更新成功后清理没有任何容器引用的旧镜像。 */
  cleanupAfter: boolean
  /** 显式写了域名的镜像（ghcr.io 等）直连，不套加速源。 */
  directFirst: boolean

  // —— 备份保留策略 ——
  backupKeepPerContainer: number
  backupMaxAgeDays: number
  backupMaxTotalMB: number
  /** 更新前快照永不自动清理。 */
  backupKeepPreUpdate: boolean
}

/** 自动更新：单个容器在本轮里的去向。 */
export interface AutoRunItem {
  name: string
  image: string
  status: 'updated' | 'up_to_date' | 'failed' | 'skipped' | 'pending'
  reason: string
}

/** 自动更新：一轮「巡检 + 可选执行」的结果。 */
export interface AutoRunSummary {
  startedAt: string
  finishedAt: string
  trigger: string
  dryRun: boolean
  checked: number
  available: number
  updated: number
  failed: number
  reclaimedMB: number
  durationMs: number
  containers: AutoRunItem[]
  error?: string
}

/** 自动更新：候选筛选结论（会被更新 / 被跳过 / 为什么）。 */
export interface AutoCandidate {
  name: string
  image: string
  running: boolean
  hasUpdate: boolean
  willUpdate: boolean
  protected: boolean
  excluded: boolean
  reason: string
}

export interface AutoUpdateInfo {
  enabled: boolean
  checkIntervalHours: number
  notifyOnCheck: boolean
  intervalChoices: number[]
  policy: {
    concurrency: number
    pullOnce: boolean
    backupBefore: boolean
    cleanupAfter: boolean
    directFirst: boolean
  }
  running: boolean
  nextCheckAt: string
  lastCheckAt: string
  candidates: AutoCandidate[]
  willUpdate: number
  exclude: string[]
  selfName: string
  lastRun: AutoRunSummary | null
}

export interface ImageView {
  id: string
  shortId: string
  tags: string[]
  /** 从 RepoDigests 推导的仓库名 —— 无 tag 镜像（未使用）靠它显示可读名称。 */
  repo?: string
  size: number
  created: number
  containers: number
  dangling: boolean
}

export interface BusEvent {
  id: number
  topic: string
  kind: string
  time: string
  data?: Record<string, unknown>
  status?: string
}

/** Docker 网络（GET /api/networks 原样透传，字段按需取用）。 */
export interface DockerNetwork {
  Id: string
  Name: string
  Driver: string
  Scope: string
  Internal: boolean
  Attachable: boolean
  Created?: string
  IPAM?: { Config?: { Subnet?: string; Gateway?: string }[] }
  Containers?: Record<string, { Name: string; IPv4Address?: string }>
  Labels?: Record<string, string>
}
