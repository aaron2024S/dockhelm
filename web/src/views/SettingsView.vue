<script setup lang="ts">
import { computed, onMounted, ref } from 'vue'
import { RouterLink } from 'vue-router'
import {
  Ban,
  KeyRound,
  Loader2,
  LogOut,
  RefreshCw,
  Save,
  ScrollText,
  ShieldAlert,
  Trash2,
  User,
  Users,
} from 'lucide-vue-next'
import { api } from '@/api/client'
import type { ContainerView, RunLog, Settings } from '@/api/types'
import { formatDateTime, relativeTime } from '@/utils/format'
import { useToastStore } from '@/stores/toast'
import { useAppStore } from '@/stores/app'

const toast = useToastStore()
const app = useAppStore()

const settings = ref<Settings>({
  exclude: [],
  panelURL: '',
  concurrency: 2,
  deepCheckCron: '',
  logRetention: 500,
  checkOnStart: true,
})
const containers = ref<ContainerView[]>([])
const logs = ref<RunLog[]>([])
const account = ref<Record<string, any> | null>(null)
const loading = ref(true)
const saving = ref(false)
const newExclude = ref('')

const oldPw = ref('')
const newPw = ref('')
const confirmPw = ref('')
const changing = ref(false)

/** 模板里安全渲染「双花括号」变量写法（直接写字面量会被 Vue 的插值分隔符截断）。 */
const tplVar = (name: string) => `{{${name}}}`

const recentLogins = computed<{ ts: string; ip: string; ok: boolean }[]>(
  () => (account.value?.recentLogins as { ts: string; ip: string; ok: boolean }[]) ?? [],
)

const pwError = computed(() => {
  if (!newPw.value) return ''
  if (newPw.value.length < app.minPasswordLength) return `新密码至少 ${app.minPasswordLength} 位`
  if (newPw.value !== confirmPw.value) return '两次输入的新密码不一致'
  return ''
})

async function load() {
  loading.value = true
  try {
    const [s, c, l, a] = await Promise.all([
      api.get<Settings>('/api/settings'),
      api.get<{ containers: ContainerView[] }>('/api/containers'),
      api.get<{ logs: RunLog[] }>('/api/logs', { limit: 120 }),
      api.get<Record<string, any>>('/api/account'),
    ])
    settings.value = s
    containers.value = c.containers ?? []
    logs.value = l.logs ?? []
    account.value = a
  } catch (e) {
    toast.error('读取设置失败', e instanceof Error ? e.message : String(e))
  } finally {
    loading.value = false
  }
}

async function save() {
  saving.value = true
  try {
    settings.value = await api.put<Settings>('/api/settings', settings.value)
    await app.loadSettings()
    toast.success('设置已保存', '并发度改动会在下一次批量更新时生效')
  } catch (e) {
    toast.error('保存失败', e instanceof Error ? e.message : String(e))
  } finally {
    saving.value = false
  }
}

function addExclude() {
  const n = newExclude.value.trim()
  if (!n || settings.value.exclude.includes(n)) return
  settings.value.exclude.push(n)
  newExclude.value = ''
}

function removeExclude(name: string) {
  settings.value.exclude = settings.value.exclude.filter((x) => x !== name)
}

async function changePassword() {
  if (pwError.value || !oldPw.value || !newPw.value) return
  changing.value = true
  try {
    await api.post('/api/account/password', {
      oldPassword: oldPw.value,
      newPassword: newPw.value,
      confirm: confirmPw.value,
    })
    toast.success('密码已修改', '所有设备上的登录会话都已失效，请重新登录')
    oldPw.value = ''
    newPw.value = ''
    confirmPw.value = ''
    window.setTimeout(() => void doLogout(), 1200)
  } catch (e) {
    toast.error('修改失败', e instanceof Error ? e.message : String(e))
  } finally {
    changing.value = false
  }
}

/** 退出登录并把浏览器带回登录页。 */
async function doLogout() {
  try {
    await app.logout()
  } finally {
    window.location.replace('/login')
  }
}

async function clearLogs() {
  try {
    await api.del('/api/logs')
    logs.value = []
    toast.success('运行记录已清空')
  } catch (e) {
    toast.error('清空失败', e instanceof Error ? e.message : String(e))
  }
}

onMounted(() => void load())
</script>

<template>
  <div class="flex flex-col gap-3.5 p-[18px]">
    <!-- 更新行为 -->
    <div class="dh-card">
      <div class="dh-card-head">
        <RefreshCw class="h-3.5 w-3.5 text-text-4" />
        <span>更新行为</span>
        <button class="dh-btn dh-btn-sm dh-btn-primary ml-auto" :disabled="saving" @click="save">
          <Loader2 v-if="saving" class="h-3 w-3 dh-spin" />
          <Save v-else class="h-3 w-3" />保存
        </button>
      </div>
      <div class="grid grid-cols-1 gap-x-8 gap-y-4 p-3.5 lg:grid-cols-2">
        <div>
          <label class="dh-label">批量更新并发度</label>
          <select v-model.number="settings.concurrency" class="dh-select">
            <option :value="1">1（串行，最稳）</option>
            <option :value="2">2（推荐）</option>
            <option :value="3">3</option>
            <option :value="4">4</option>
            <option :value="6">6（激进）</option>
          </select>
          <div class="mt-1 text-[11px] leading-relaxed text-text-5">
            同时更新的容器数量。并发越高越快，但更容易触发镜像仓库限流。
          </div>
        </div>

        <div>
          <label class="dh-label">面板地址（通知模板里的 <code>{{ tplVar('url') }}</code>）</label>
          <input v-model="settings.panelURL" class="dh-input" placeholder="http://192.168.1.10:8080" />
        </div>

        <label class="flex items-start gap-2.5">
          <input v-model="settings.checkOnStart" type="checkbox" class="mt-[3px] h-[14px] w-[14px] accent-[#2dd4bf]" />
          <span class="text-[12.5px]">
            启动后自动巡检一次
            <div class="text-[11px] text-text-5">延迟 8 秒执行，只读检查，不会停止任何容器。</div>
          </span>
        </label>

        <div>
          <label class="dh-label">运行记录保留条数</label>
          <input v-model.number="settings.logRetention" type="number" min="50" class="dh-input" />
        </div>

        <div class="lg:col-span-2">
          <div class="flex items-start gap-2.5 rounded-[10px] border border-line-1 bg-ink-800 px-3 py-2.5 text-[11.5px] leading-relaxed text-text-4">
            <ShieldAlert class="mt-[1px] h-3.5 w-3.5 flex-none text-[#fbbf24]" />
            <span>
              这两类容器永远不会被自动更新：<b class="text-text-3">Dockhelm 自己</b>，以及下面的排除列表。
              计划任务里的「全部容器」也会自动跳过它们。
            </span>
          </div>
        </div>
      </div>
    </div>

    <!-- 排除列表 -->
    <div class="dh-card">
      <div class="dh-card-head">
        <Ban class="h-3.5 w-3.5 text-text-4" />
        <span>排除列表</span>
        <span class="ml-auto text-[11.5px] font-normal text-text-5">
          {{ settings.exclude.length }} 个容器永不自动更新
        </span>
      </div>
      <div class="flex flex-col gap-2.5 p-3.5">
        <div class="flex flex-wrap gap-2">
          <select v-model="newExclude" class="dh-select !w-auto !min-w-[190px]">
            <option value="">选择容器…</option>
            <option
              v-for="c in containers.filter((x) => !settings.exclude.includes(x.name))"
              :key="c.id"
              :value="c.name"
            >
              {{ c.name }}
            </option>
          </select>
          <button class="dh-btn" :disabled="!newExclude" @click="addExclude">加入排除</button>
        </div>
        <div v-if="app.settings === null && loading" class="text-[12px] text-text-5">正在载入…</div>
        <div v-if="!settings.exclude.length" class="text-[12px] text-text-5">排除列表为空。</div>
        <div v-else class="flex flex-wrap gap-1.5">
          <span
            v-for="n in settings.exclude"
            :key="n"
            class="inline-flex items-center gap-1.5 rounded-full border border-line-3 px-2.5 py-[3px] font-mono text-[11.5px] text-text-3"
          >
            {{ n }}
            <button class="text-text-5 hover:text-[#fca5a5]" @click="removeExclude(n)">×</button>
          </span>
        </div>
        <div class="text-[11px] text-text-5">
          计划任务里的「全部容器」会自动跳过排除项与 Dockhelm 自身。
        </div>
      </div>
    </div>

    <div class="grid grid-cols-1 gap-3.5 xl:grid-cols-2">
      <!-- 账户 -->
      <div class="dh-card">
        <div class="dh-card-head">
          <User class="h-3.5 w-3.5 text-text-4" />
          <span>账户与安全</span>
          <RouterLink to="/about" class="ml-auto text-[11.5px] font-normal text-text-5 hover:text-accent">
            关于 Dockhelm
          </RouterLink>
        </div>
        <div class="flex flex-col gap-3.5 p-3.5">
          <div class="grid grid-cols-2 gap-3 text-[12px]">
            <div class="text-text-5">当前会话数</div>
            <div class="text-text-2">
              <Users class="mr-1 inline h-3 w-3 text-text-5" />{{ account?.sessionCount ?? '—' }}
            </div>
            <div class="text-text-5">当前会话创建于</div>
            <div class="text-text-2">
              {{ account?.currentSession ? formatDateTime(account.currentSession.createdAt) : '—' }}
            </div>
            <div class="text-text-5">当前会话过期于</div>
            <div class="text-text-2">
              {{ account?.currentSession ? formatDateTime(account.currentSession.expiresAt) : '—' }}
            </div>
          </div>

          <div class="border-t border-line-1 pt-3.5">
            <div class="mb-2.5 flex items-center gap-2 text-[12.5px] font-medium">
              <KeyRound class="h-3.5 w-3.5 text-text-4" />修改登录密码
            </div>
            <div class="flex flex-col gap-2.5">
              <input v-model="oldPw" type="password" class="dh-input" placeholder="当前密码" autocomplete="current-password" />
              <input v-model="newPw" type="password" class="dh-input" placeholder="新密码" autocomplete="new-password" />
              <input v-model="confirmPw" type="password" class="dh-input" placeholder="再输一次新密码" autocomplete="new-password" />
              <div v-if="pwError" class="text-[11.5px] text-[#fca5a5]">{{ pwError }}</div>
              <button
                class="dh-btn dh-btn-primary"
                :disabled="changing || !oldPw || !!pwError || !newPw"
                @click="changePassword"
              >
                <Loader2 v-if="changing" class="h-3.5 w-3.5 dh-spin" />
                <KeyRound v-else class="h-3.5 w-3.5" />
                修改密码
              </button>
              <div class="text-[11px] leading-relaxed text-text-5">
                修改后所有设备上的会话立即失效，需要重新登录。
                忘记密码时：在 NAS 面板里删除 <code class="text-text-3">data/auth.json</code> 后重启容器，
                或在 compose 里临时加 <code class="text-text-3">DOCKHELM_PASSWORD=新密码</code>。
              </div>
            </div>
          </div>

          <div class="border-t border-line-1 pt-3.5">
            <div class="mb-2 text-[12.5px] font-medium">最近登录记录</div>
            <div v-if="!recentLogins.length" class="text-[11.5px] text-text-5">暂无记录</div>
            <div v-else class="flex flex-col gap-1">
              <div
                v-for="(r, i) in recentLogins"
                :key="i"
                class="flex items-center gap-2 text-[11.5px]"
              >
                <span
                  class="h-[6px] w-[6px] flex-none rounded-full"
                  :class="r.ok ? 'bg-run' : 'bg-err'"
                />
                <span class="w-[86px] flex-none text-text-5">{{ relativeTime(r.ts) }}</span>
                <span class="font-mono text-text-3">{{ r.ip }}</span>
                <span class="ml-auto text-text-5">{{ r.ok ? '成功' : '失败' }}</span>
              </div>
            </div>
          </div>

          <button class="dh-btn" @click="doLogout">
            <LogOut class="h-3.5 w-3.5" />退出登录
          </button>
        </div>
      </div>

      <!-- 运行记录 -->
      <div class="dh-card">
        <div class="dh-card-head">
          <ScrollText class="h-3.5 w-3.5 text-text-4" />
          <span>运行记录</span>
          <div class="ml-auto flex gap-2">
            <button class="dh-btn dh-btn-sm" :disabled="loading" @click="load">
              <RefreshCw class="h-3 w-3" :class="loading ? 'dh-spin' : ''" />
            </button>
            <button class="dh-btn dh-btn-sm dh-btn-danger" :disabled="!logs.length" @click="clearLogs">
              <Trash2 class="h-3 w-3" />清空
            </button>
          </div>
        </div>
        <div v-if="!logs.length" class="dh-card-body text-[12.5px] text-text-4">暂无记录</div>
        <div v-else class="dh-scroll max-h-[560px] overflow-auto">
          <table class="dh-table">
            <thead>
              <tr>
                <th class="w-[120px]">时间</th>
                <th class="w-[86px]">类型</th>
                <th>说明</th>
              </tr>
            </thead>
            <tbody>
              <tr v-for="l in logs" :key="l.id">
                <td class="whitespace-nowrap text-[11.5px] text-text-5">{{ relativeTime(l.ts) }}</td>
                <td>
                  <span
                    class="dh-badge"
                    :class="l.status === 'failed' ? 'dh-badge-err' : l.status === 'success' || l.status === 'up_to_date' ? 'dh-badge-accent' : 'dh-badge-plain'"
                  >
                    {{ l.kind }}
                  </span>
                </td>
                <td>
                  <div class="text-[12px] text-text-2">{{ l.message }}</div>
                  <div v-if="l.ref" class="font-mono text-[10.5px] text-text-6">{{ l.ref }}</div>
                </td>
              </tr>
            </tbody>
          </table>
        </div>
      </div>
    </div>
  </div>
</template>
