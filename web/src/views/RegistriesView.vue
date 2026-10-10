<script setup lang="ts">
import { computed, onMounted, ref } from 'vue'
import {
  AlertTriangle,
  CheckCircle2,
  Copy,
  Gauge,
  GripVertical,
  Info,
  Loader2,
  Plus,
  Rocket,
  Save,
  Zap,
} from 'lucide-vue-next'
import { api } from '@/api/client'
import type { MirrorConfig, RegistriesResponse, RegistrySettings, Settings } from '@/api/types'
import { formatDayTime } from '@/utils/format'
import { useAppStore } from '@/stores/app'
import { useToastStore } from '@/stores/toast'
import EmptyState from '@/components/EmptyState.vue'
import SettingRow from '@/components/SettingRow.vue'
import ToggleSwitch from '@/components/ToggleSwitch.vue'
import UiSwitch from '@/components/UiSwitch.vue'

const toast = useToastStore()
const app = useAppStore()
const data = ref<RegistriesResponse | null>(null)
const loading = ref(true)
const testing = ref(false)
const testingUrl = ref('')
const newUrl = ref('')
const newNote = ref('')
const copied = ref(false)
const saving = ref(false)

const mirrors = computed<MirrorConfig[]>(() => data.value?.settings?.mirrors ?? [])
const presets = computed<MirrorConfig[]>(() => data.value?.presets ?? [])
const daemonMirrors = computed<string[]>(() => data.value?.daemonMirrors ?? [])

/** 与后端 normalizeMirror 同口径：忽略大小写与结尾斜杠。 */
function urlKey(u: string) {
  return u.trim().toLowerCase().replace(/\/+$/, '')
}

const pullMirror = computed({
  get: () => data.value?.settings?.pullMirror ?? '',
  set: (v: string) => {
    if (data.value?.settings) data.value.settings.pullMirror = v
  },
})

async function load() {
  loading.value = true
  try {
    data.value = await api.get<RegistriesResponse>('/api/registries')
    savedSig.value = sigOf(data.value.settings)
  } catch (e) {
    toast.error('读取加速源配置失败', e instanceof Error ? e.message : String(e))
  } finally {
    loading.value = false
  }
}

async function save() {
  if (!data.value) return
  saving.value = true
  let registrySaved = false
  try {
    await api.put('/api/registries', { settings: data.value.settings })
    registrySaved = true
    // directFirst 属于全局设置。用 PATCH 只提交这一个键 ——
    // 以前是整份 Settings PUT 回去，会把这个页面进页面之后在设置页改的东西一起还原。
    if (policy.value) {
      policy.value = await api.patch<Settings>('/api/settings', { directFirst: policy.value.directFirst })
    }
    void app.loadSettings()
    toast.success('已保存')
    await load()
  } catch (e) {
    // 两段保存是分开的：说清哪一段成了、哪一段没成，别只丢一句「保存失败」
    toast.error(
      registrySaved ? '加速源已保存，但 directFirst 策略没保存成功' : '保存失败',
      e instanceof Error ? e.message : String(e),
    )
  } finally {
    saving.value = false
  }
}

function addMirror() {
  const url = newUrl.value.trim()
  if (!url || !data.value) return
  // 去重：列表用 url 当 key，允许重复会让 Vue 报重复 key、按索引删还会删错行
  if (mirrors.value.some((m) => urlKey(m.url) === urlKey(url))) {
    toast.error('这个地址已经在列表里了')
    return
  }
  data.value.settings.mirrors.push({ url, note: newNote.value.trim(), enabled: true })
  newUrl.value = ''
  newNote.value = ''
}

/**
 * 把预置的常用加速源补回「我的加速源」，并直接保存。
 *
 * 正常情况下用不到 —— 预置清单在服务首次启动时就已经整份灌进来了。留这个入口只为兜底：
 * 用户把列表删空之后还能一键找回。删掉的条目本身不会自己回来（后端只灌一次）。
 */
async function addPresets() {
  if (!data.value) return
  let added = 0
  for (const p of presets.value) {
    if (mirrors.value.some((m) => urlKey(m.url) === urlKey(p.url))) continue
    data.value.settings.mirrors.push({ ...p })
    added++
  }
  if (added) await save()
}

function removeMirror(idx: number) {
  data.value?.settings.mirrors.splice(idx, 1)
}

/* ---------------------------------------------------------------------------
   拖动排序：顺序就是这个列表的优先级（daemonSnippet 按数组顺序生成），
   所以「拖一下就变顺序」比抄来抄去地址直观得多。

   用 pointer 事件而不是 HTML5 dragstart：后者在触屏上根本不触发，
   而这一页用户多半是拿手机开的。手柄上写了 touch-action:none，
   拖动时不会连带把页面滚起来。
--------------------------------------------------------------------------- */
const dragFrom = ref(-1)
const dragOver = ref(-1)
let dragMoved = false

function moveMirror(from: number, to: number) {
  const list = data.value?.settings.mirrors
  if (!list) return
  if (from === to || from < 0 || to < 0 || from >= list.length || to >= list.length) return
  const [row] = list.splice(from, 1)
  if (!row) return
  list.splice(to, 0, row)
  dragFrom.value = to
  dragOver.value = to
  dragMoved = true
}

function startDrag(i: number, ev: PointerEvent) {
  if (mirrors.value.length < 2) return
  ev.preventDefault()
  dragFrom.value = i
  dragOver.value = i
  dragMoved = false
  window.addEventListener('pointermove', onDragMove)
  window.addEventListener('pointerup', endDrag)
  window.addEventListener('pointercancel', endDrag)
}

/** 指针越过某行中线就把它换过去（实时换位，拖动时看到的就是最终顺序）。 */
function onDragMove(ev: PointerEvent) {
  const from = dragFrom.value
  if (from < 0) return
  const rows = Array.from(document.querySelectorAll<HTMLElement>('[data-mirror-row]'))
  for (let k = 0; k < rows.length; k++) {
    if (k === from) continue
    const el = rows[k]
    if (!el) continue
    const rect = el.getBoundingClientRect()
    const mid = rect.top + rect.height / 2
    if (k < from && ev.clientY < mid) return moveMirror(from, k)
    if (k > from && ev.clientY > mid) return moveMirror(from, k)
  }
}

function endDrag() {
  dragFrom.value = -1
  dragOver.value = -1
  window.removeEventListener('pointermove', onDragMove)
  window.removeEventListener('pointerup', endDrag)
  window.removeEventListener('pointercancel', endDrag)
  if (dragMoved) toast.success('顺序已调整', '点「保存」后生效')
  dragMoved = false
}

/** 键盘等价操作：手柄聚焦后用 ↑ ↓ 换位（列表里有 5、6 条时很顺手）。 */
function onGripKey(i: number, ev: KeyboardEvent) {
  if (ev.key !== 'ArrowUp' && ev.key !== 'ArrowDown') return
  ev.preventDefault()
  const to = ev.key === 'ArrowUp' ? i - 1 : i + 1
  if (to < 0 || to >= mirrors.value.length) return
  moveMirror(i, to)
  toast.success('顺序已调整', '点「保存」后生效')
  dragMoved = false
}

/**
 * 保存基线。只比对「用户能改的字段」——
 * 测速结果（延迟 / 报错）是点一下就会写回列表的，不能因此把页面标成「有未保存的改动」。
 */
const savedSig = ref('')
function sigOf(s: RegistrySettings | undefined): string {
  if (!s) return ''
  return JSON.stringify({
    mirrors: s.mirrors.map((m) => [m.url, m.note, m.enabled]),
    pullMirror: s.pullMirror,
    insecure: s.insecure,
  })
}
const dirty = computed(() => !!data.value && sigOf(data.value.settings) !== savedSig.value)

/** 把测速结果按地址写回列表里对应的那一行。 */
function applyResults(list: MirrorConfig[], results: MirrorConfig[]) {
  for (const r of results) {
    const hit = list.find((x) => urlKey(x.url) === urlKey(r.url))
    if (!hit) continue
    hit.latencyMs = r.latencyMs
    hit.ok = r.ok
    hit.err = r.err
    hit.lastTested = r.lastTested
  }
}

async function testOne(url: string) {
  testingUrl.value = url
  try {
    const res = await api.post<MirrorConfig>('/api/registries/test', { url })
    // 只把结果写回本地，不整页重载 —— 重载会把还没点「保存」的添加/删除一起丢掉
    applyResults(mirrors.value, [res])
    if (res.ok) {
      toast.success(`${res.url} 可用`, `延迟 ${res.latencyMs} ms`)
    } else {
      toast.error(`${res.url} 不可用`, res.err || '连接失败')
    }
  } catch (e) {
    toast.error('测速失败', e instanceof Error ? e.message : String(e))
  } finally {
    testingUrl.value = ''
  }
}

async function testAll() {
  testing.value = true
  try {
    const urls = [...new Set(mirrors.value.map((m) => m.url))]
    const res = await api.post<{ results: MirrorConfig[] }>('/api/registries/test-all', { urls })
    applyResults(mirrors.value, res.results ?? [])
    toast.success('测速完成', '结果已写回各自那一行')
  } catch (e) {
    toast.error('批量测速失败', e instanceof Error ? e.message : String(e))
  } finally {
    testing.value = false
  }
}

/** 页头的「添加加速源」：把新增加速源的输入框滚进视野。 */
function focusAdd() {
  document.getElementById('mirror-add')?.scrollIntoView({ behavior: 'smooth', block: 'center' })
  document.getElementById('mirror-add')?.querySelector('input')?.focus()
}

async function copySnippet() {
  const text = data.value?.snippet ?? ''
  try {
    await navigator.clipboard.writeText(text)
    copied.value = true
    window.setTimeout(() => (copied.value = false), 1800)
    toast.success('已复制到剪贴板')
  } catch {
    toast.error('复制失败', '请手动选中代码块复制')
  }
}

/** 徽标配色：以「测没测过」为准（lastTested 为空 = 从未测过）。 */
const latencyTone = (m: MirrorConfig) => {
  if (!m.lastTested) return 'dh-badge-plain'
  if (!m.ok) return 'dh-badge-err'
  if ((m.latencyMs ?? 9999) < 300) return 'dh-badge-run'
  if ((m.latencyMs ?? 9999) < 1200) return 'dh-badge-warn'
  return 'dh-badge-plain'
}

/** 连通性徽标文案：正常时给延迟，失败时给人话。 */
function connText(m: MirrorConfig): string {
  if (!m.lastTested) return '未测速'
  if (m.ok) return `正常 · ${m.latencyMs ?? 0} ms`
  return shortErr(m.err)
}

/**
 * 把探测失败的原因压成一句话。
 *
 * 后端直出的是 Go 的网络错误原文（可能带完整 URL、几十个字），
 * 塞进徽标里既难看也读不懂 —— 这里归成「域名解析失败 / 403 拒绝访问」这类结论。
 */
function shortErr(err?: string): string {
  const s = (err ?? '').trim()
  if (!s) return '连接失败'
  if (/no such host|server misbehaving|lookup .* on /i.test(s)) return '域名解析失败'
  if (/deadline exceeded|timed out|timeout/i.test(s)) return '连接超时'
  if (/connection refused/i.test(s)) return '拒绝连接'
  if (/certificate|x509|tls:/i.test(s)) return '证书错误'
  if (/connection reset|unexpected EOF/i.test(s)) return '连接被重置'
  const http = s.match(/HTTP\s*(\d{3})/)
  if (http) {
    const code = Number(http[1])
    if (code === 401 || code === 403) return `${code} 拒绝访问`
    if (code === 404) return '404 未找到'
    if (code >= 500) return `${code} 服务异常`
    return `HTTP ${code}`
  }
  return s.length > 18 ? `${s.slice(0, 18)}…` : s
}

/**
 * 本页只借 /api/settings 里的一个开关：directFirst（显式域名优先）。
 *
 * 「并发度」「检测频率」「检测后通知」这些旋钮已经完整存在于「设置 → 更新与检测」，
 * 这里再放一份就是同一份配置的第二、第三个入口 —— 改了一处不知道另一处也会变，
 * 所以本页不再重复提供，只保留与「加速源怎么用」直接相关的那一项。
 */
const policy = ref<Settings | null>(null)

async function loadPolicy() {
  try {
    policy.value = await api.get<Settings>('/api/settings')
  } catch {
    policy.value = null
  }
}

onMounted(() => {
  void load()
  void loadPolicy()
})
</script>

<template>
  <div class="flex flex-col gap-3.5 p-[18px]">
    <div class="dh-phead">
      <div class="dh-h1">镜像加速源</div>
      <div class="dh-sub">按顺序优先使用，失败自动顺延下一个</div>
      <div class="ml-auto flex gap-2">
        <button class="dh-btn" :disabled="testing" @click="testAll">
          <Gauge class="h-3.5 w-3.5" :class="testing ? 'dh-spin' : ''" />测试全部
        </button>
        <button class="dh-btn dh-btn-primary" @click="focusAdd">添加加速源</button>
      </div>
    </div>

    <!-- 现状：守护进程真正生效的加速源 -->
    <div class="dh-card">
      <div class="dh-card-head">
        <Info class="h-3.5 w-3.5 text-text-4" />
        <span>Docker 守护进程当前生效的加速源</span>
        <span class="ml-auto text-[11.5px] font-normal text-text-5">来源：daemon.json</span>
      </div>
      <div class="dh-card-body flex flex-col gap-2.5">
        <div class="text-[12px] leading-relaxed text-text-4">{{ data?.explain }}</div>
        <div v-if="data?.daemonError" class="text-[12px] text-err-text">
          无法读取守护进程信息：{{ data.daemonError }}
        </div>
        <div v-else-if="daemonMirrors.length" class="flex flex-wrap gap-1.5">
          <span
            v-for="m in daemonMirrors"
            :key="m"
            class="rounded-md bg-ink-800 px-2 py-[3px] font-mono text-[11.5px] text-accent"
          >
            {{ m }}
          </span>
        </div>
        <div v-else class="flex items-center gap-2 text-[12px] text-text-5">
          <AlertTriangle class="h-3.5 w-3.5" />
          守护进程没有配置任何加速源，拉取镜像会直连 Docker Hub。
        </div>
      </div>
    </div>

    <div class="grid grid-cols-1 gap-3.5 xl:grid-cols-2">
      <!-- 我的加速源列表：整行铺满 —— 表格化之后列需要宽度，挤在半栏里会被压扁 -->
      <div class="dh-card xl:col-span-2">
        <div class="dh-card-head">
          <Rocket class="h-3.5 w-3.5 text-text-4" />
          <span>我的加速源</span>
          <div class="ml-auto flex items-center gap-2">
            <span v-if="dirty" class="text-[11.5px] text-warn-text">有未保存的改动</span>
            <button class="dh-btn dh-btn-sm" :disabled="testing" @click="testAll">
              <Gauge class="h-3 w-3" :class="testing ? 'dh-spin' : ''" />批量测速
            </button>
            <button class="dh-btn dh-btn-sm dh-btn-primary" :disabled="saving" @click="save">保存</button>
          </div>
        </div>

        <div id="mirror-add" class="flex flex-col gap-2 border-b border-line-1 p-3">
          <div class="flex flex-wrap gap-2">
            <input v-model="newUrl" class="dh-input flex-1 !min-w-[180px]" placeholder="https://你的加速站地址" />
            <input v-model="newNote" class="dh-input !w-[110px]" placeholder="备注" />
            <button class="dh-btn" :disabled="!newUrl.trim()" @click="addMirror">
              <Plus class="h-3.5 w-3.5" />添加
            </button>
          </div>
        </div>

        <EmptyState
          v-if="!mirrors.length"
          :icon="Rocket"
          :title="loading ? '正在载入…' : '还没有添加加速源'"
          :description="
            loading
              ? ''
              : '预置的常用加速站随首次启动就已经写进来了，不想要哪条直接删；这里被删空之后可以一键找回，也可以手动填写自建加速站地址。'
          "
          :action-label="!loading && presets.length ? '加入预置的常用加速站' : ''"
          @action="addPresets"
        />
        <div v-else class="overflow-x-auto">
          <table class="dh-table dh-table-fixed">
            <thead>
              <tr>
                <th class="w-[38px]" />
                <th>加速源地址</th>
                <th class="w-[175px]">连通性</th>
                <th class="w-[130px]">上次测速</th>
                <th class="w-[70px] text-right">启用</th>
                <th class="w-[130px] text-right">操作</th>
              </tr>
            </thead>
            <tbody>
              <tr
                v-for="(m, i) in mirrors"
                :key="m.url"
                data-mirror-row
                :class="{
                  'dh-drag-row': dragFrom === i,
                  'dh-drag-over': dragOver === i && dragFrom >= 0 && dragFrom !== i,
                }"
              >
                <td class="!px-2">
                  <button
                    type="button"
                    class="dh-grip"
                    :aria-label="`拖动调整 ${m.url} 的优先级，也可用上下方向键`"
                    title="拖动调整优先级（也可用 ↑ ↓）"
                    @pointerdown="startDrag(i, $event)"
                    @keydown="onGripKey(i, $event)"
                  >
                    <GripVertical class="h-3.5 w-3.5" />
                  </button>
                </td>
                <td>
                  <div class="truncate font-mono text-[12px] text-text-2" :title="m.url">{{ m.url }}</div>
                  <div class="mt-0.5 flex items-center gap-1.5 text-[11px] text-text-5">
                    <!-- 「预置」只是个来源标记：表示这条是随应用自带的，删掉与删自建条目没区别 -->
                    <span v-if="m.builtin" class="dh-badge dh-badge-plain flex-none">预置</span>
                    <span v-if="m.note" class="truncate">{{ m.note }}</span>
                  </div>
                </td>
                <td>
                  <span class="dh-badge" :class="latencyTone(m)">{{ connText(m) }}</span>
                </td>
                <td class="text-[11.5px] text-text-5">
                  {{ m.lastTested ? formatDayTime(m.lastTested) : '—' }}
                </td>
                <td class="text-right">
                  <UiSwitch
                    :model-value="m.enabled"
                    :label="`启用 ${m.url}`"
                    @update:model-value="(v: boolean) => (m.enabled = v)"
                  />
                </td>
                <td class="text-right whitespace-nowrap">
                  <button class="dh-link" :disabled="testingUrl === m.url" @click="testOne(m.url)">
                    {{ testingUrl === m.url ? '测速中…' : '测速' }}
                  </button>
                  <span class="mx-1.5 text-text-6">·</span>
                  <button class="dh-link dh-link-danger" @click="removeMirror(i)">删除</button>
                </td>
              </tr>
            </tbody>
          </table>
          <div class="border-t border-line-row px-4 py-2 text-[11px] text-text-6">
            拖动左侧把手调整优先级：从上到下依次尝试，前面失败会自动顺延下一个。
          </div>
        </div>
      </div>

      <!-- 拉取策略 -->
      <div class="dh-card">
        <div class="dh-card-head">
          <Zap class="h-3.5 w-3.5 text-text-4" />
          <span>Dockhelm 自己的拉取策略</span>
        </div>
        <div class="dh-card-body flex flex-col gap-3">
          <div>
            <label class="dh-label">拉取加速源</label>
            <select v-model="pullMirror" class="dh-select">
              <option value="">不指定（完全交给 Docker 守护进程）</option>
              <option v-for="m in mirrors" :key="m.url" :value="m.url">{{ m.url }}</option>
            </select>
            <div class="mt-1.5 text-[11.5px] leading-relaxed text-text-5">
              保持「不指定」时，拉取行为与手动执行 <code class="text-text-3">docker pull</code> 完全一致。
              如果你没法改 daemon.json，可以在这里指定一个加速源：Dockhelm 会用
              <code class="text-text-3">&lt;加速站&gt;/&lt;仓库&gt;:&lt;标签&gt;</code> 拉取，
              拉完再打回原始标签，compose 与其它工具仍然按原来的名字找得到镜像。
              <br />
              <b class="text-text-3">注意</b>：只对来自 Docker Hub 的镜像生效。
            </div>
          </div>
          <div v-if="policy" class="border-t border-line-1 pt-3">
            <SettingRow title="显式域名优先，不套用加速" sub="ghcr.io、私有仓库等已经写明域名的镜像直连 —— 加速站通常只镜像 Docker Hub">
              <ToggleSwitch v-model="policy.directFirst" label="显式域名优先" />
            </SettingRow>
          </div>
          <button class="dh-btn dh-btn-primary" :disabled="saving" @click="save">
            <Loader2 v-if="saving" class="h-3.5 w-3.5 dh-spin" />
            <Save v-else class="h-3.5 w-3.5" />保存设置
          </button>
        </div>
        </div>

        <!-- daemon.json 片段 -->
        <div class="dh-card">
          <div class="dh-card-head">
            <Copy class="h-3.5 w-3.5 text-text-4" />
            <span>daemon.json 片段</span>
            <button class="dh-btn dh-btn-sm ml-auto" @click="copySnippet">
              <component :is="copied ? CheckCircle2 : Copy" class="h-3 w-3" />
              {{ copied ? '已复制' : '复制' }}
            </button>
          </div>
          <div class="dh-card-body">
            <div class="text-[11.5px] leading-relaxed text-text-5">
              把这段贴进 NAS 上的 <code class="text-text-3">/etc/docker/daemon.json</code>（群晖在
              Docker 套件的设置里，或 Container Manager 的「注册表镜像」），重启 Docker 服务后生效。
              内容是上表里<b class="text-text-3">已启用</b>的加速源，顺序与上表一致。
            </div>
            <pre class="mt-2.5 overflow-x-auto rounded-[9px] border border-line-3 bg-ink-800 px-3 py-2.5 font-mono text-[11.5px] leading-[1.8] text-text-2">{{ data?.snippet }}</pre>
          </div>
        </div>
    </div>

    <!--
      这里原本有一张「推荐加速源」卡：一份独立的内置清单，用户要逐条点「添加」才进「我的加速源」。
      同一批数据在两个地方出现、还要手动搬运，已经废掉 —— 预置清单现在只在服务首次启动时
      整份写进「我的加速源」（默认启用），不想要就逐条删，删空后可用列表空态里的按钮一键找回。
    -->
  </div>
</template>
