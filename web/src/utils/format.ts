/** 展示用的格式化工具（保持纯函数，方便单测）。 */

export function formatBytes(n: number | undefined | null, digits = 1): string {
  if (n === undefined || n === null || Number.isNaN(n) || n <= 0) return '0 B'
  const units = ['B', 'KB', 'MB', 'GB', 'TB', 'PB']
  let i = 0
  let v = n
  while (v >= 1024 && i < units.length - 1) {
    v /= 1024
    i++
  }
  const fixed = i === 0 ? 0 : digits
  return `${v.toFixed(fixed)} ${units[i]}`
}

/** 相对时间：刚刚 / 3 分钟前 / 2 小时前 / 3 天前 / 具体日期。 */
export function relativeTime(input: string | number | undefined | null): string {
  if (!input) return '—'
  const t = typeof input === 'number' ? input * 1000 : Date.parse(input)
  if (Number.isNaN(t)) return '—'
  const diff = Date.now() - t
  const abs = Math.abs(diff)
  const future = diff < 0
  const sec = Math.round(abs / 1000)
  const fmt = (n: number, unit: string) => (future ? `${n} ${unit}后` : `${n} ${unit}前`)
  if (sec < 45) return future ? '即将' : '刚刚'
  const min = Math.round(sec / 60)
  if (min < 60) return fmt(min, '分钟')
  const hr = Math.round(min / 60)
  if (hr < 24) return fmt(hr, '小时')
  const day = Math.round(hr / 24)
  if (day < 30) return fmt(day, '天')
  return formatDateTime(t)
}

export function formatDateTime(input: string | number | undefined | null): string {
  if (!input) return '—'
  const t = typeof input === 'number' ? (input > 1e12 ? input : input * 1000) : Date.parse(input)
  if (Number.isNaN(t)) return '—'
  const d = new Date(t)
  const p = (n: number) => String(n).padStart(2, '0')
  return `${d.getFullYear()}-${p(d.getMonth() + 1)}-${p(d.getDate())} ${p(d.getHours())}:${p(d.getMinutes())}:${p(d.getSeconds())}`
}

export function formatDateShort(input: string | number | undefined | null): string {
  if (!input) return '—'
  const t = typeof input === 'number' ? input * 1000 : Date.parse(input)
  if (Number.isNaN(t)) return '—'
  const d = new Date(t)
  const p = (n: number) => String(n).padStart(2, '0')
  return `${p(d.getMonth() + 1)}-${p(d.getDate())} ${p(d.getHours())}:${p(d.getMinutes())}`
}

/** 把秒数转成「3 天 4 小时」这种可读形式。 */
export function formatDuration(seconds: number | undefined | null): string {
  if (!seconds || seconds < 0) return '—'
  const d = Math.floor(seconds / 86400)
  const h = Math.floor((seconds % 86400) / 3600)
  const m = Math.floor((seconds % 3600) / 60)
  const parts: string[] = []
  if (d) parts.push(`${d} 天`)
  if (h) parts.push(`${h} 小时`)
  if (!d && m) parts.push(`${m} 分钟`)
  if (!parts.length) parts.push(`${Math.floor(seconds)} 秒`)
  return parts.join(' ')
}

const STATUS_LABEL: Record<string, string> = {
  running: '运行中',
  exited: '已停止',
  created: '已创建',
  paused: '已暂停',
  restarting: '重启中',
  dead: '异常',
  removing: '删除中',
}

export function containerStateLabel(state: string): string {
  return STATUS_LABEL[state] ?? state
}

const CHECK_LABEL: Record<string, { text: string; tone: 'accent' | 'warn' | 'plain' | 'err' }> = {
  up_to_date: { text: '已是最新', tone: 'accent' },
  update_available: { text: '有新版本', tone: 'warn' },
  unknown: { text: '无法判定', tone: 'plain' },
  no_upstream: { text: '本地镜像', tone: 'plain' },
}

export function checkLabel(status: string) {
  return CHECK_LABEL[status] ?? { text: status, tone: 'plain' as const }
}

const RESULT_LABEL: Record<string, { text: string; tone: 'accent' | 'warn' | 'plain' | 'err' }> = {
  up_to_date: { text: '已是最新·已跳过', tone: 'accent' },
  updated: { text: '更新成功', tone: 'accent' },
  failed: { text: '更新失败·已回滚', tone: 'warn' },
  broken: { text: '失败且回滚失败', tone: 'err' },
  skipped: { text: '已跳过', tone: 'plain' },
  pulling: { text: '拉取中', tone: 'plain' },
  stopped: { text: '已停止旧容器', tone: 'plain' },
  created: { text: '已创建新容器', tone: 'plain' },
}

export function resultLabel(status: string) {
  return RESULT_LABEL[status] ?? { text: status, tone: 'plain' as const }
}

/** 截断长文本。 */
export function truncate(s: string, n: number): string {
  if (!s) return ''
  return s.length <= n ? s : s.slice(0, n) + '…'
}

/** 简化镜像引用用于展示（去掉 library/ 前缀）。 */
export function shortImage(image: string): string {
  if (!image) return '—'
  return image.replace(/^docker\.io\/library\//, '').replace(/^library\//, '')
}

/**
 * 带日期的短时间：今天 04:00 / 明天 04:00 / 周三 08:00 / 10 月 10 日 02:00。
 * mode 为 'date' 时只出日期部分（10 月 3 日）。
 */
export function formatDayTime(
  input: string | number | undefined | null,
  mode: 'datetime' | 'date' = 'datetime',
): string {
  if (!input) return '—'
  const t = typeof input === 'number' ? (input > 1e12 ? input : input * 1000) : Date.parse(input)
  if (Number.isNaN(t)) return '—'
  const d = new Date(t)
  const p = (n: number) => String(n).padStart(2, '0')
  const hm = `${p(d.getHours())}:${p(d.getMinutes())}`

  const day0 = new Date(t)
  day0.setHours(0, 0, 0, 0)
  const today0 = new Date()
  today0.setHours(0, 0, 0, 0)
  const diffDays = Math.round((day0.getTime() - today0.getTime()) / 86400000)

  let day: string
  if (diffDays === 0) day = '今天'
  else if (diffDays === 1) day = '明天'
  else if (diffDays === 2) day = '后天'
  else if (diffDays > 2 && diffDays < 7) day = `周${'日一二三四五六'[d.getDay()]}`
  else day = `${d.getMonth() + 1} 月 ${d.getDate()} 日`

  return mode === 'date' ? day : `${day} ${hm}`
}

/** 把 cron 表达式翻译成中文说明（覆盖常见写法，认不出就原样返回）。 */
export function explainCron(expr: string): string {
  const e = expr.trim()
  if (!e) return ''
  if (e === '@daily' || e === '@midnight') return '每天 00:00'
  if (e === '@hourly') return '每小时'
  if (e === '@weekly') return '每周'
  if (e === '@monthly') return '每月'
  if (e === '@yearly' || e === '@annually') return '每年'
  const every = /^@every\s+(.+)$/.exec(e)
  if (every) return `每 ${every[1]}`

  const parts = e.split(/\s+/)
  if (parts.length !== 5) return e
  // parts.length 已确认等于 5，这里显式取值（tsconfig 开了 noUncheckedIndexedAccess）
  const min = parts[0] as string
  const hour = parts[1] as string
  const dom = parts[2] as string
  const mon = parts[3] as string
  const dow = parts[4] as string
  const isNum = (s: string) => /^\d+$/.test(s)
  const pad = (s: string) => s.padStart(2, '0')

  if (isNum(min) && isNum(hour)) {
    const hm = `${pad(hour)}:${pad(min)}`
    if (dom === '*' && mon === '*' && dow === '*') return `每天 ${hm}`
    if (dom === '*' && mon === '*' && dow !== '*') {
      const names: Record<string, string> = {
        '0': '周日', '1': '周一', '2': '周二', '3': '周三',
        '4': '周四', '5': '周五', '6': '周六', '7': '周日',
      }
      const list = dow.split(',').map((d) => names[d] ?? `周${d}`).join('、')
      return `每${list} ${hm}`
    }
    if (mon === '*' && dow === '*' && isNum(dom)) return `每月 ${dom} 号 ${hm}`
  }
  if (min === '*' && isNum(hour)) return `每小时的第 ${hour} 分（${pad(hour)}:${min.replace('*', '00')} 起每分钟）`
  if (min.startsWith('*/')) return `每 ${min.slice(2)} 分钟`
  if (hour.startsWith('*/')) return `每 ${hour.slice(2)} 小时`
  return e
}
