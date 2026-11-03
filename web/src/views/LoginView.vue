<script setup lang="ts">
/**
 * 登录 / 首次设置密码。
 *
 * 布局：左边品牌区（介绍面板能干什么 + 鉴权说明），右边表单。
 * 单用户密码模式 —— mailnest 那种「用户名 + 密码」在这里没有意义，
 * 所以首次部署只让设一个密码，之后也只要求输密码。
 *
 * 交互上几件容易被忽略的事，都在下面处理了：
 * ① 用真 <form> + submit：回车能提交，浏览器密码管理器也认得出这是登录框；
 * ② 刷新页面时若仍在锁定期，倒计时必须接上（否则用户以为能试、一直失败）；
 * ③ 服务端返回的剩余次数要用服务端的值，不能自己减（多标签页各错一次就错乱了）。
 */
import { computed, onMounted, onUnmounted, ref } from 'vue'
import { useRouter } from 'vue-router'
import {
  ArchiveRestore,
  ArrowRight,
  Bell,
  Boxes,
  CalendarClock,
  Eye,
  EyeOff,
  Layers,
  ListOrdered,
  Loader2,
  Lock,
  RefreshCw,
  TriangleAlert,
} from 'lucide-vue-next'
import { useAppStore } from '@/stores/app'
import { ApiError } from '@/api/client'
import { version } from '@/config'

const app = useAppStore()
const router = useRouter()

const mode = ref<'loading' | 'setup' | 'login'>('loading')
const password = ref('')
const confirm = ref('')
const showPw = ref(false)
const keep = ref(true)
const busy = ref(false)
const errorMsg = ref('')
const remaining = ref<number | null>(null)
const locked = ref(false)
/** 锁定解除的时刻（毫秒时间戳）。用相对时限而不是服务端绝对时间，省得依赖两端时钟一致。 */
const lockDeadline = ref(0)
/** 每秒推一下的「现在」，只为驱动倒计时重算。 */
const nowTick = ref(Date.now())
let lockTimer: number | undefined

/** 左侧品牌区条目：只写用户真正能感知的能力，不写实现细节。 */
const features = [
  { icon: Boxes, name: '容器', desc: '启停重启、日志跟随、实时资源占用，环境变量自动脱敏' },
  { icon: Layers, name: '镜像与加速源', desc: '体积与引用数、清理未使用；实测生效的加速站与 daemon.json' },
  { icon: RefreshCw, name: '更新与自动更新', desc: '巡检后标出可更新，按周期自动更新；镜像没变就不动容器' },
  { icon: CalendarClock, name: '计划任务', desc: '5 段 cron：启动 / 停止 / 重启 / 更新 / 备份 / 清理' },
  { icon: ListOrdered, name: '顺序启动', desc: '开机后按你排的顺序串行拉起：拖拽排序、逐行延时、逐行开关' },
  { icon: ArchiveRestore, name: '备份与恢复', desc: '配置快照、差异预览、一键还原；compose / .env 一并备份' },
  { icon: Bell, name: '通知', desc: '11 类渠道、按事件订阅；静默时段、去重、失败重试与推送历史' },
]

/** 锁定还要等多少秒（向上取整，0 表示已解锁）。 */
const lockLeft = computed(() => {
  const ms = lockDeadline.value - nowTick.value
  return ms > 0 ? Math.ceil(ms / 1000) : 0
})

function mmss(total: number) {
  const m = Math.floor(total / 60)
  const s = total % 60
  return `${m}:${String(s).padStart(2, '0')}`
}

function ensureLockTicker() {
  if (lockTimer !== undefined) return
  lockTimer = window.setInterval(() => {
    nowTick.value = Date.now()
    if (lockDeadline.value <= Date.now()) clearLock()
  }, 1000)
}

/** 进入锁定态：msg 是服务端给的「约 N 分钟」，seconds > 0 时顺带起倒计时。 */
function applyLock(msg: string, seconds: number) {
  errorMsg.value = msg
  locked.value = true
  if (seconds > 0) {
    lockDeadline.value = Date.now() + seconds * 1000
    nowTick.value = Date.now()
    ensureLockTicker()
  } else {
    lockDeadline.value = 0
  }
}

/** 锁定结束（或重载后已过期）：清掉倒计时与错误条，按钮重新可用。 */
function clearLock() {
  locked.value = false
  lockDeadline.value = 0
  errorMsg.value = ''
  remaining.value = app.maxFailures
  if (lockTimer !== undefined) {
    window.clearInterval(lockTimer)
    lockTimer = undefined
  }
}

onUnmounted(() => {
  if (lockTimer !== undefined) window.clearInterval(lockTimer)
})

const isSetup = computed(() => mode.value === 'setup')

/** 密码强度：0~4 */
const strength = computed(() => {
  const p = password.value
  if (!p) return 0
  let s = 0
  if (p.length >= 6) s++
  if (p.length >= 10) s++
  if (/[A-Za-z]/.test(p) && /\d/.test(p)) s++
  if (/[^A-Za-z0-9]/.test(p)) s++
  return Math.min(s, 4)
})
const strengthLabel = computed(() => ['', '偏弱', '一般', '较强', '很强'][strength.value])

/** 强度条下的说明：没输入时别显示「强度：」这个空壳。 */
const strengthHint = computed(() =>
  password.value
    ? `强度：${strengthLabel.value} · 至少 ${app.minPasswordLength} 位，建议含数字与符号`
    : `至少 ${app.minPasswordLength} 位，建议含数字与符号`,
)

const canSubmit = computed(() => {
  if (busy.value || locked.value) return false
  if (isSetup.value)
    return password.value.length >= app.minPasswordLength && password.value === confirm.value
  return password.value.length > 0
})

/** 错误条语气：锁定 / 首次设置出错用红，登录失败可重试用琥珀（与设计稿一致）。 */
const bannerTone = computed(() =>
  locked.value || isSetup.value ? 'dh-banner-err' : 'dh-banner-warn',
)

const bannerText = computed(() => {
  if (locked.value) {
    // 光说「请稍后再试」等于没说：把服务端的「约 N 分钟」和本地倒计时一起摆出来
    return lockLeft.value > 0 ? `${errorMsg.value}（剩余 ${mmss(lockLeft.value)}）` : errorMsg.value
  }
  if (!errorMsg.value) return ''
  if (remaining.value !== null && remaining.value > 0) {
    return `${errorMsg.value}，还可尝试 ${remaining.value} 次`
  }
  return errorMsg.value
})

const submitLabel = computed(() => {
  if (busy.value) return '请稍候…'
  if (locked.value) return lockLeft.value > 0 ? `已锁定 ${mmss(lockLeft.value)}` : '已锁定'
  return isSetup.value ? '设置并进入' : '登 录'
})

onMounted(async () => {
  if (!app.ready) await app.bootstrap()
  if (!app.initialized) {
    mode.value = 'setup'
    return
  }
  if (app.loggedIn) {
    router.replace('/overview')
    return
  }
  mode.value = 'login'
  remaining.value = app.maxFailures - app.failures
  // 刷新页面后仍在锁定期：提示与倒计时都要接上，别让用户以为能试却一直失败
  if (app.lockedFor > 0) applyLock(app.lockedHint, app.lockedFor)
})

async function submit() {
  if (!canSubmit.value) return
  busy.value = true
  errorMsg.value = ''
  try {
    if (isSetup.value) {
      const res = await app.setup(password.value, confirm.value)
      if (res.autoLogin) {
        router.replace('/overview')
      } else {
        mode.value = 'login'
        password.value = ''
        confirm.value = ''
      }
    } else {
      await app.login(password.value, keep.value)
      router.replace('/overview')
    }
  } catch (e) {
    if (e instanceof ApiError) {
      const payload = e.payload as { remaining?: number; retryAfter?: number } | undefined
      if (e.code === 'locked') {
        applyLock(e.message, typeof payload?.retryAfter === 'number' ? payload.retryAfter : 0)
      } else {
        errorMsg.value = e.message
        if (typeof payload?.remaining === 'number') remaining.value = payload.remaining
        if (e.code === 'bad_password') password.value = ''
      }
    } else {
      errorMsg.value = e instanceof Error ? e.message : String(e)
    }
  } finally {
    busy.value = false
  }
}
</script>

<template>
  <div class="dh-auth">
    <!-- ============ 左：品牌区 ============ -->
    <aside class="dh-auth-brand">
      <div class="dh-auth-top">
        <div class="dh-auth-mark">
          <svg viewBox="0 0 32 32" class="h-[21px] w-[21px]">
            <circle cx="16" cy="16" r="12.3" fill="none" stroke="currentColor" stroke-width="2.6" />
            <path
              fill-rule="evenodd"
              fill="currentColor"
              d="M16 5.9 18.76 13.24 26.1 16 18.76 18.76 16 26.1 13.24 18.76 5.9 16 13.24 13.24ZM18.2 16A2.2 2.2 0 1 0 13.8 16A2.2 2.2 0 1 0 18.2 16Z"
            />
          </svg>
        </div>
        <div class="text-[16px] font-semibold tracking-[-.2px]">Dockhelm</div>
      </div>

      <div class="dh-auth-mid">
        <h1>容器、镜像、任务，<br />一个面板全收住</h1>
        <p class="dh-auth-lead">
          打开就是总览：容器、镜像、更新、任务都在这儿。<br />
          不用再记一堆 docker 命令，也不用 SSH 翻 compose。
        </p>
        <ul class="dh-auth-feats">
          <li v-for="f in features" :key="f.name">
            <span class="dh-auth-ico">
              <component :is="f.icon" class="h-[15px] w-[15px]" />
            </span>
            <div>
              <b>{{ f.name }}</b>
              <span class="d">{{ f.desc }}</span>
            </div>
          </li>
        </ul>
        <p class="dh-auth-extra">
          登录只要一个密码：<b>bcrypt</b> 哈希落盘、HttpOnly 会话、连续失败锁定。
        </p>
      </div>

      <div class="dh-auth-brandfoot">© 2026 Dockhelm · MIT 许可</div>
    </aside>

    <!-- ============ 右：表单 ============ -->
    <main class="dh-auth-form">
      <form class="dh-auth-box" @submit.prevent="submit">
        <template v-if="isSetup">
          <span class="dh-auth-eyebrow"><i />首次启动</span>
          <h2>设置访问密码</h2>
          <p class="dh-auth-sub">这台 Dockhelm 还没有密码。先设一个，之后用它进入面板。</p>
        </template>
        <template v-else>
          <h2>登录 Dockhelm</h2>
          <p class="dh-auth-sub">输入访问密码继续。</p>
        </template>

        <div v-if="mode === 'loading'" class="grid h-[120px] place-items-center text-text-5">
          <Loader2 class="h-5 w-5 dh-spin" />
        </div>

        <template v-else>
          <!-- 密码（首次要输两遍） -->
          <div class="mb-[15px]">
            <label class="mb-[7px] block text-[12.5px] font-semibold text-text-3" for="dh-pw">
              {{ isSetup ? '密码' : '访问密码' }}
            </label>
            <div class="dh-field !h-[42px]">
              <Lock class="h-[15px] w-[15px] flex-none text-text-6" />
              <input
                id="dh-pw"
                v-model="password"
                :type="showPw ? 'text' : 'password'"
                :placeholder="isSetup ? `至少 ${app.minPasswordLength} 位` : '请输入密码'"
                :autocomplete="isSetup ? 'new-password' : 'current-password'"
                autofocus
              />
              <button
                type="button"
                class="dh-eye"
                :title="showPw ? '隐藏密码' : '显示密码'"
                :aria-label="showPw ? '隐藏密码' : '显示密码'"
                @click="showPw = !showPw"
              >
                <EyeOff v-if="showPw" class="h-[15px] w-[15px]" />
                <Eye v-else class="h-[15px] w-[15px]" />
              </button>
            </div>

            <template v-if="isSetup">
              <div class="mt-[8px] flex gap-1">
                <i
                  v-for="i in 4"
                  :key="i"
                  class="h-[4px] flex-1 rounded-[2px]"
                  :class="i <= strength ? 'bg-run' : 'bg-line-1'"
                />
              </div>
              <p class="mt-[6px] text-[11.5px] text-text-6">{{ strengthHint }}</p>
            </template>
          </div>

          <div v-if="isSetup" class="mb-[15px]">
            <label class="mb-[7px] block text-[12.5px] font-semibold text-text-3" for="dh-pw2">
              确认密码
            </label>
            <div class="dh-field !h-[42px]">
              <Lock class="h-[15px] w-[15px] flex-none text-text-6" />
              <input
                id="dh-pw2"
                v-model="confirm"
                :type="showPw ? 'text' : 'password'"
                placeholder="再输一次"
                autocomplete="new-password"
              />
            </div>
            <p v-if="confirm && confirm !== password" class="mt-[6px] text-[11.5px] text-err-text">
              两次输入不一致
            </p>
          </div>

          <div
            v-if="bannerText"
            class="dh-banner mb-[14px] !items-start !gap-2 !py-[9px] !pl-[11px] !pr-[11px] !text-[12.5px]"
            :class="bannerTone"
            role="alert"
          >
            <TriangleAlert class="h-[14px] w-[14px] flex-none translate-y-[1px]" />
            <span>{{ bannerText }}</span>
          </div>

          <label
            v-if="!isSetup"
            class="flex cursor-pointer items-center gap-2 text-[12px] text-text-3"
          >
            <input v-model="keep" type="checkbox" class="h-[14px] w-[14px] accent-accent" />
            保持登录（7 天）
          </label>

          <button
            type="submit"
            class="dh-auth-submit"
            :disabled="!canSubmit"
            :aria-busy="busy"
          >
            <Loader2 v-if="busy" class="h-[15px] w-[15px] dh-spin" />
            <span>{{ submitLabel }}</span>
            <ArrowRight v-if="!busy" class="arrow" />
          </button>

          <p v-if="isSetup" class="dh-auth-note">
            密码只存 <code class="font-mono">bcrypt</code> 哈希，明文不落盘、也不写日志。<br />
            Dockhelm 持有 Docker 套接字，等于宿主机 root 权限，请务必设置密码。
          </p>
          <p v-else class="dh-auth-note">连续错 {{ app.maxFailures }} 次锁定 5 分钟</p>
        </template>
      </form>

      <div class="dh-auth-formfoot">
        Dockhelm {{ version }}<span class="sep">·</span>单用户密码
      </div>
    </main>
  </div>
</template>
