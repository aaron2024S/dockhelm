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
  minPassword: number
  sessionCount: number
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
  images: { total: number; sizeBytes: number }
  updates: {
    available: number
    unknown: number
    checkedAt: string
    items: { container: string; image: string; reason: string }[]
  }
  disk: { free: number; total: number }
  notify: { sentToday: number }
  recent: RunLog[]
  excluded: string[]
  self: string
  docker?: {
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
  suggestions: MirrorConfig[]
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

export interface ProjectInfo {
  project: string
  containers: string[]
  configFiles: string[]
  readable: ProjectFile[]
  unreadable: string[]
  workDir: string
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
  deepCheckCron: string
  logRetention: number
  checkOnStart: boolean
}

export interface ImageView {
  id: string
  shortId: string
  tags: string[]
  digests: string[]
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
