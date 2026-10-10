<script setup lang="ts">
import { computed, onMounted, ref } from 'vue'
import {
  AlertTriangle,
  CheckCircle2,
  Copy,
  Gauge,
  Info,
  Loader2,
  Plus,
  Rocket,
  Save,
  Trash2,
  XCircle,
  Zap,
} from 'lucide-vue-next'
import { api } from '@/api/client'
import type { MirrorConfig, RegistriesResponse, Settings } from '@/api/types'
import { useAppStore } from '@/stores/app'
import { useToastStore } from '@/stores/toast'
import EmptyState from '@/components/EmptyState.vue'
import SettingRow from '@/components/SettingRow.vue'
import ToggleSwitch from '@/components/ToggleSwitch.vue'

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
      <!-- 我的加速源列表 -->
      <div class="dh-card">
        <div class="dh-card-head">
          <Rocket class="h-3.5 w-3.5 text-text-4" />
          <span>我的加速源</span>
          <div class="ml-auto flex gap-2">
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
        <div v-else class="flex flex-col">
          <div
            v-for="(m, i) in mirrors"
            :key="m.url"
            class="flex flex-wrap items-center gap-2 border-b border-line-row px-3 py-2.5 last:border-b-0"
          >
            <input
              type="checkbox"
              class="h-[14px] w-[14px] accent-accent"
              :checked="m.enabled"
              @change="m.enabled = !m.enabled"
              title="启用这个加速源"
            />
            <div class="min-w-0 flex-1">
              <div class="truncate font-mono text-[11.5px] text-text-2">{{ m.url }}</div>
              <div class="flex items-center gap-1.5 text-[11px] text-text-5">
                <!-- 「预置」只是个来源标记：表示这条是随应用自带的，删掉与删自建条目没区别 -->
                <span v-if="m.builtin" class="dh-badge dh-badge-plain flex-none">预置</span>
                <span v-if="m.note" class="truncate">{{ m.note }}</span>
              </div>
            </div>
            <span v-if="m.lastTested" class="dh-badge" :class="latencyTone(m)">
              <CheckCircle2 v-if="m.ok" class="h-3 w-3" />
              <XCircle v-else class="h-3 w-3" />
              {{ m.ok ? m.latencyMs + ' ms' : '不可用' }}
            </span>
            <button class="dh-btn dh-btn-sm" :disabled="testingUrl === m.url" @click="testOne(m.url)">
              <Zap class="h-3 w-3" :class="testingUrl === m.url ? 'dh-spin' : ''" />测速
            </button>
            <button class="dh-btn dh-btn-sm dh-btn-danger" @click="removeMirror(i)">
              <Trash2 class="h-3 w-3" />
            </button>
          </div>
        </div>
      </div>

      <div class="flex flex-col gap-3.5">
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
              内容是上表里<b class="text-text-3">勾选启用</b>的加速源。
            </div>
            <pre class="mt-2.5 overflow-x-auto rounded-[9px] border border-line-3 bg-ink-800 px-3 py-2.5 font-mono text-[11.5px] leading-[1.8] text-text-2">{{ data?.snippet }}</pre>
          </div>
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
