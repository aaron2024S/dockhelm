<script setup lang="ts">
import { computed, onMounted, ref } from 'vue'
import {
  Ban,
  KeyRound,
  Loader2,
  LogOut,
  RefreshCw,
  Save,
  ScrollText,
  SearchCheck,
  ShieldAlert,
  SlidersHorizontal,
  Trash2,
  User,
  Users,
  Zap,
} from 'lucide-vue-next'
import { api } from '@/api/client'
import type { ContainerView, RunLog, Settings } from '@/api/types'
import { formatDateTime, relativeTime } from '@/utils/format'
import { useToastStore } from '@/stores/toast'
import { useAppStore } from '@/stores/app'
import Modal from '@/components/Modal.vue'
import SettingRow from '@/components/SettingRow.vue'
import ToggleSwitch from '@/components/ToggleSwitch.vue'

const toast = useToastStore()
const app = useAppStore()

const settings = ref<Settings>({
  exclude: [],
  panelURL: '',
  concurrency: 2,
  deepCheckCron: '',
  logRetention: 500,
  checkOnStart: true,
  checkIntervalHours: 6,
  notifyOnCheck: false,
  autoApply: false,
  pullOnce: true,
  backupBefore: true,
  cleanupAfter: true,
  directFirst: true,
  backupKeepPerContainer: 10,
  backupMaxAgeDays: 30,
  backupMaxTotalMB: 2048,
  backupKeepPreUpdate: true,
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
const showOld = ref(false)
const showNew = ref(false)
const changing = ref(false)
/** 「清空运行记录」的二次确认 —— 其余破坏性操作都有确认框，这个以前点了就清。 */
const confirmClearLogs = ref(false)
/** 清空请求进行中 —— 防连点。 */
const clearingLogs = ref(false)

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
    // 只提交本页真正在编辑的字段（PATCH）。
    //
    // 全局设置被拆在三个页面里：本页、镜像加速源页（directFirst）、备份与恢复页
    // （backupKeep*/backupMax*）。以前这里把整份 Settings PUT 回去，那份快照是
    // 进页面时拉的 —— 期间在备份页改过的保留策略会被静默还原。PATCH 后后端只改
    // 出现的键，其余保持原值。
    const s = settings.value
    settings.value = await api.patch<Settings>('/api/settings', {
      exclude: s.exclude,
      panelURL: s.panelURL,
      concurrency: s.concurrency,
      logRetention: s.logRetention,
      checkOnStart: s.checkOnStart,
      checkIntervalHours: s.checkIntervalHours,
      notifyOnCheck: s.notifyOnCheck,
      autoApply: s.autoApply,
      pullOnce: s.pullOnce,
      backupBefore: s.backupBefore,
      cleanupAfter: s.cleanupAfter,
    })
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
  if (clearingLogs.value) return
  clearingLogs.value = true
  confirmClearLogs.value = false
  try {
    await api.del('/api/logs')
    logs.value = []
    toast.success('运行记录已清空')
  } catch (e) {
    toast.error('清空失败', e instanceof Error ? e.message : String(e))
  } finally {
    clearingLogs.value = false
  }
}

onMounted(() => void load())
</script>

<template>
  <div class="flex flex-col gap-3.5 p-[18px]">
    <div class="dh-phead">
      <div class="dh-h1">设置</div>
      <div class="dh-sub">更新与检测、排除列表与账户安全</div>
    </div>

    <!-- 更新与检测 -->
    <div class="dh-card">
      <div class="dh-card-head">
        <RefreshCw class="h-3.5 w-3.5 text-text-4" />
        <span>更新与检测</span>
        <span class="ml-2 text-[11.5px] font-normal text-text-5">
          {{ settings.checkIntervalHours > 0 ? `每 ${settings.checkIntervalHours} 小时巡检一次` : '未开启周期巡检' }}
          · 自动更新{{ settings.autoApply ? '已开启' : '已关闭' }}
        </span>
        <button class="dh-btn dh-btn-sm dh-btn-primary ml-auto" :disabled="saving" @click="save">
          <Loader2 v-if="saving" class="h-3 w-3 dh-spin" />
          <Save v-else class="h-3 w-3" />保存
        </button>
      </div>

      <div class="flex flex-col gap-3 p-3.5">
        <!-- 检测：这是「多久检查一次」的落点 -->
        <div class="rounded-[10px] border border-line-1 bg-ink-800 p-3">
          <div class="mb-2.5 flex items-center gap-2 text-[12px] font-medium text-text-3">
            <SearchCheck class="h-3.5 w-3.5 text-text-4" />检测
          </div>
          <div class="flex flex-col gap-2.5">
            <SettingRow title="检测频率" :sub="settings.checkIntervalHours > 0 ? `每 ${settings.checkIntervalHours} 小时自动扫描一次镜像仓库` : '关闭后只在手动点「重新检测」时才检查'">
              <select v-model.number="settings.checkIntervalHours" class="dh-select !w-[120px] !py-[5px] !text-[11.5px]">
                <option :value="0">关闭</option>
                <option :value="1">每 1 小时</option>
                <option :value="3">每 3 小时</option>
                <option :value="6">每 6 小时</option>
                <option :value="12">每 12 小时</option>
                <option :value="24">每 24 小时</option>
              </select>
            </SettingRow>
            <SettingRow title="摘要来源" sub="由本机 Docker 守护进程解析，与 docker pull 走同一数据源">
              <span class="dh-badge dh-badge-plain">推荐</span>
            </SettingRow>
            <SettingRow title="检测完成后通知" sub="巡检发现问题时推一条通知到已配置的渠道">
              <ToggleSwitch v-model="settings.notifyOnCheck" label="检测完成后通知" />
            </SettingRow>
            <SettingRow title="启动后自动巡检一次" sub="延迟 8 秒执行，只读检查，不会停止任何容器">
              <ToggleSwitch v-model="settings.checkOnStart" label="启动后自动巡检一次" />
            </SettingRow>
          </div>
        </div>

        <!-- 自动更新：回答「检测到有更新会不会自己动手」 -->
        <div class="rounded-[10px] border p-3" :class="settings.autoApply ? 'border-line-warn bg-soft-warn' : 'border-line-1 bg-ink-800'">
          <div class="mb-2.5 flex items-center gap-2 text-[12px] font-medium text-text-3">
            <Zap class="h-3.5 w-3.5" :class="settings.autoApply ? 'text-warn-text' : 'text-text-4'" />自动更新
          </div>
          <div class="flex flex-col gap-2.5">
            <SettingRow
              title="检测到新版本后自动更新"
              :sub="
                settings.autoApply
                  ? '每轮巡检结束后，把有更新的容器（排除列表与自己除外）自动重建到新镜像'
                  : '当前只检测、不动手 —— 发现更新后要你在更新中心手动点'
              "
            >
              <ToggleSwitch v-model="settings.autoApply" label="自动更新" />
            </SettingRow>
            <div
              v-if="settings.autoApply"
              class="flex items-start gap-2 rounded-[8px] border border-line-warn px-2.5 py-2 text-[11.5px] leading-relaxed text-warn-text"
            >
              <ShieldAlert class="mt-[1px] h-3.5 w-3.5 flex-none" />
              <span>
                自动更新会<b>真的重启容器</b>。建议保持「更新前自动备份容器配置」开启，并确认排除列表里
                放好了数据库、反向代理这类不能随便重启的服务。更新中心页可以先看「本轮会更新哪几个」再决定。
              </span>
            </div>
          </div>
        </div>

        <!-- 执行策略 -->
        <div class="rounded-[10px] border border-line-1 bg-ink-800 p-3">
          <div class="mb-2.5 flex items-center gap-2 text-[12px] font-medium text-text-3">
            <SlidersHorizontal class="h-3.5 w-3.5 text-text-4" />执行策略
          </div>
          <div class="flex flex-col gap-2.5">
            <SettingRow title="仅在镜像真正变化时重启容器" sub="比对更新前后的镜像 ID，一致就完全不动它">
              <span class="dh-badge dh-badge-plain">始终开启</span>
            </SettingRow>
            <SettingRow title="同一镜像的多个容器只拉取一次" sub="4 个容器共用 nginx:alpine 时，下载 1 次、重建 4 个">
              <ToggleSwitch v-model="settings.pullOnce" label="同一镜像只拉取一次" />
            </SettingRow>
            <SettingRow title="更新前自动备份容器配置" sub="失败可一键回滚到更新前">
              <ToggleSwitch v-model="settings.backupBefore" label="更新前自动备份容器配置" />
            </SettingRow>
            <SettingRow title="更新后清理旧镜像" sub="确认没有任何容器再引用后才删除">
              <ToggleSwitch v-model="settings.cleanupAfter" label="更新后清理旧镜像" />
            </SettingRow>
            <SettingRow title="批量更新并发度" sub="并发越高越快，但更容易触发镜像仓库限流">
              <select v-model.number="settings.concurrency" class="dh-select !w-[150px] !py-[5px] !text-[11.5px]">
                <option :value="1">1（串行，最稳）</option>
                <option :value="2">2（推荐）</option>
                <option :value="3">3</option>
                <option :value="4">4</option>
                <option :value="6">6（激进）</option>
              </select>
            </SettingRow>
          </div>
        </div>

        <!-- 面板与记录 -->
        <div class="rounded-[10px] border border-line-1 bg-ink-800 p-3">
          <div class="mb-2.5 flex items-center gap-2 text-[12px] font-medium text-text-3">
            <ScrollText class="h-3.5 w-3.5 text-text-4" />面板与记录
          </div>
          <div class="flex flex-col gap-2.5">
            <SettingRow title="面板地址" sub="通知模板里的 {{url}} 用它拼可点击的链接">
              <input v-model="settings.panelURL" class="dh-input !w-[260px]" placeholder="http://192.168.1.10:5923" />
            </SettingRow>
            <SettingRow title="运行记录保留条数" sub="超出后按时间滚动覆盖，只影响面板里的历史列表">
              <input v-model.number="settings.logRetention" type="number" min="50" class="dh-input !w-[110px]" />
            </SettingRow>
          </div>
        </div>

        <div class="flex items-start gap-2.5 rounded-[10px] border border-line-1 bg-ink-800 px-3 py-2.5 text-[11.5px] leading-relaxed text-text-4">
          <ShieldAlert class="mt-[1px] h-3.5 w-3.5 flex-none text-warn-text" />
          <span>
            这两类容器永远不会被自动更新：<b class="text-text-3">Dockhelm 自己</b>，以及下面的排除列表。
            计划任务也只会作用于你在任务里<b class="text-text-3">明确勾选</b>的容器（「不选」不代表全部）。
          </span>
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
            <button class="dh-tap text-text-5 hover:text-err-text" @click="removeExclude(n)">×</button>
          </span>
        </div>
        <div class="text-[11px] text-text-5">
          加进这里的容器不会被自动更新，也不会被任何计划任务作用到（计划任务必须逐个勾选容器）。
        </div>
      </div>
    </div>

    <div class="grid grid-cols-1 gap-3.5 xl:grid-cols-2">
      <!-- 账户 -->
      <div class="dh-card">
        <div class="dh-card-head">
          <User class="h-3.5 w-3.5 text-text-4" />
          <span>账户</span>
          <span class="ml-auto text-[11.5px] font-mono font-normal text-text-5">右上角头像 → 账户</span>
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
            <div class="flex flex-col gap-[11px]">
              <div>
                <label class="mb-[5px] block text-[12px] text-text-4">当前密码</label>
                <div class="dh-field">
                  <input
                    v-model="oldPw"
                    :type="showOld ? 'text' : 'password'"
                    placeholder="当前使用的密码"
                    autocomplete="current-password"
                  />
                  <button type="button" class="dh-eye" @click="showOld = !showOld">
                    {{ showOld ? '隐藏' : '显示' }}
                  </button>
                </div>
              </div>
              <div>
                <label class="mb-[5px] block text-[12px] text-text-4">新密码</label>
                <div class="dh-field">
                  <input
                    v-model="newPw"
                    :type="showNew ? 'text' : 'password'"
                    :placeholder="`至少 ${app.minPasswordLength} 位`"
                    autocomplete="new-password"
                  />
                  <button type="button" class="dh-eye" @click="showNew = !showNew">
                    {{ showNew ? '隐藏' : '显示' }}
                  </button>
                </div>
              </div>
              <div>
                <label class="mb-[5px] block text-[12px] text-text-4">确认新密码</label>
                <div class="dh-field">
                  <input
                    v-model="confirmPw"
                    :type="showNew ? 'text' : 'password'"
                    placeholder="再输一次新密码"
                    autocomplete="new-password"
                  />
                </div>
              </div>
              <div v-if="pwError" class="text-[11.5px] text-err-text">{{ pwError }}</div>
              <div class="dh-banner dh-banner-info !gap-2.5 !py-[9px] !pl-3 !pr-3 !text-[12px]">
                <span class="h-[13px] w-[13px] flex-none rounded-[4px] bg-current opacity-50" />
                <span class="min-w-0 flex-1">改密成功后其他设备上的登录会立即失效，需重新登录。</span>
              </div>
              <div class="flex justify-end">
                <button
                  class="dh-btn dh-btn-primary"
                  :disabled="changing || !oldPw || !!pwError || !newPw"
                  @click="changePassword"
                >
                  <Loader2 v-if="changing" class="h-3.5 w-3.5 dh-spin" />
                  <KeyRound v-else class="h-3.5 w-3.5" />
                  保存新密码
                </button>
              </div>
              <div class="text-[11px] leading-relaxed text-text-5">
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
            <button
              class="dh-btn dh-btn-sm dh-btn-danger"
              :disabled="!logs.length || clearingLogs"
              @click="confirmClearLogs = true"
            >
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
    <!-- 清空确认：会真删记录，必须二次确认 -->
    <Modal
      :open="confirmClearLogs"
      title="清空运行记录"
      width="430px"
      :busy="clearingLogs"
      @close="confirmClearLogs = false"
    >
      <div class="text-[12.5px] leading-relaxed text-text-3">清空后全部更新 / 巡检 / 备份的<b class="text-err-text">执行记录</b>会立刻消失，不可恢复。容器、镜像与快照本身都不受影响。</div>
      <template #footer>
        <button class="dh-btn" :disabled="clearingLogs" @click="confirmClearLogs = false">取消</button>
        <button class="dh-btn dh-btn-danger" :disabled="clearingLogs" @click="clearLogs">
          <Trash2 class="h-3.5 w-3.5" />{{ clearingLogs ? '清空中…' : '确认清空' }}
        </button>
      </template>
    </Modal>

  </div>
</template>
