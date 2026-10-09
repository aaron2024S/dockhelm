<script setup lang="ts">
import { computed, onMounted, ref } from 'vue'
import {
  AlertTriangle,
  CheckCircle2,
  Copy,
  Gauge,
  Info,
  Plus,
  Rocket,
  Trash2,
  XCircle,
  Zap,
} from 'lucide-vue-next'
import { api } from '@/api/client'
import type { MirrorConfig, RegistriesResponse } from '@/api/types'
import { useToastStore } from '@/stores/toast'
import EmptyState from '@/components/EmptyState.vue'

const toast = useToastStore()
const data = ref<RegistriesResponse | null>(null)
const loading = ref(true)
const testing = ref(false)
const testingUrl = ref('')
const newUrl = ref('')
const newNote = ref('')
const copied = ref(false)
const saving = ref(false)

const mirrors = computed<MirrorConfig[]>(() => data.value?.settings.mirrors ?? [])
const suggestions = computed<MirrorConfig[]>(() => data.value?.suggestions ?? [])
const daemonMirrors = computed<string[]>(() => data.value?.daemonMirrors ?? [])
const pullMirror = computed({
  get: () => data.value?.settings.pullMirror ?? '',
  set: (v: string) => {
    if (data.value) data.value.settings.pullMirror = v
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
  try {
    await api.put('/api/registries', { settings: data.value.settings })
    toast.success('已保存')
    await load()
  } catch (e) {
    toast.error('保存失败', e instanceof Error ? e.message : String(e))
  } finally {
    saving.value = false
  }
}

function addMirror() {
  const url = newUrl.value.trim()
  if (!url || !data.value) return
  data.value.settings.mirrors.push({ url, note: newNote.value.trim(), enabled: true })
  newUrl.value = ''
  newNote.value = ''
}

function addSuggestion(s: MirrorConfig) {
  if (!data.value) return
  data.value.settings.mirrors.push({ ...s, enabled: true })
  data.value.suggestions = data.value.suggestions.filter((x) => x.url !== s.url)
}

function removeMirror(idx: number) {
  data.value?.settings.mirrors.splice(idx, 1)
}

async function testOne(url: string) {
  testingUrl.value = url
  try {
    const res = await api.post<MirrorConfig>('/api/registries/test', { url })
    if (res.ok) {
      toast.success(`${res.url} 可用`, `延迟 ${res.latencyMs} ms`)
    } else {
      toast.error(`${res.url} 不可用`, res.err || '连接失败')
    }
    await load()
  } catch (e) {
    toast.error('测速失败', e instanceof Error ? e.message : String(e))
  } finally {
    testingUrl.value = ''
  }
}

async function testAll() {
  testing.value = true
  try {
    const urls = [
      ...mirrors.value.map((m) => m.url),
      ...suggestions.value.map((s) => s.url),
    ]
    await api.post('/api/registries/test-all', { urls })
    toast.success('测速完成', '结果已按延迟排序并保存')
    await load()
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

const latencyTone = (m: MirrorConfig) => {
  if (m.ok === undefined) return 'dh-badge-plain'
  if (!m.ok) return 'dh-badge-err'
  if ((m.latencyMs ?? 9999) < 300) return 'dh-badge-run'
  if ((m.latencyMs ?? 9999) < 1200) return 'dh-badge-warn'
  return 'dh-badge-plain'
}

onMounted(() => void load())
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
        <div v-if="data?.daemonError" class="text-[12px] text-[#fca5a5]">
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

    <div class="grid grid-cols-1 gap-3.5 xl:grid-cols-[1.3fr_1fr]">
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
          description="可以从下方的推荐列表一键添加，也可以手动填写自建加速站地址。"
        />
        <div v-else class="flex flex-col">
          <div
            v-for="(m, i) in mirrors"
            :key="m.url"
            class="flex flex-wrap items-center gap-2 border-b border-[#171f2a] px-3 py-2.5 last:border-b-0"
          >
            <input
              type="checkbox"
              class="h-[14px] w-[14px] accent-[#2dd4bf]"
              :checked="m.enabled"
              @change="m.enabled = !m.enabled"
              title="启用这个加速源"
            />
            <div class="min-w-0 flex-1">
              <div class="truncate font-mono text-[11.5px] text-text-2">{{ m.url }}</div>
              <div v-if="m.note" class="text-[11px] text-text-5">{{ m.note }}</div>
            </div>
            <span v-if="m.ok !== undefined" class="dh-badge" :class="latencyTone(m)">
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
            <button class="dh-btn dh-btn-primary" :disabled="saving" @click="save">保存设置</button>
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

    <!-- 推荐 -->
    <div v-if="suggestions.length" class="dh-card">
      <div class="dh-card-head">
        <Rocket class="h-3.5 w-3.5 text-text-4" />
        <span>推荐加速源</span>
        <span class="ml-auto text-[11.5px] font-normal text-text-5">
          点「添加」加入列表；这些只是起点，可用性请用测速确认
        </span>
      </div>
      <div class="grid grid-cols-1 gap-2 p-3 md:grid-cols-2 xl:grid-cols-3">
        <div
          v-for="s in suggestions"
          :key="s.url"
          class="flex items-center gap-2 rounded-[10px] border border-line-1 bg-ink-800 px-3 py-2"
        >
          <div class="min-w-0 flex-1">
            <div class="truncate font-mono text-[11.5px] text-text-2">{{ s.url }}</div>
            <div class="text-[11px] text-text-5">{{ s.note }}</div>
          </div>
          <span v-if="s.ok !== undefined" class="dh-badge" :class="latencyTone(s)">
            {{ s.ok ? s.latencyMs + ' ms' : '不可用' }}
          </span>
          <button class="dh-btn dh-btn-sm" @click="testOne(s.url)">
            <Zap class="h-3 w-3" />
          </button>
          <button class="dh-btn dh-btn-sm" @click="addSuggestion(s)">
            <Plus class="h-3 w-3" />添加
          </button>
        </div>
      </div>
    </div>
  </div>
</template>
