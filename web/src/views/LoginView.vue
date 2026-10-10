<script setup lang="ts">
import { computed, onMounted, onUnmounted, ref } from 'vue'
import { useRouter } from 'vue-router'
import { Eye, EyeOff, Loader2 } from 'lucide-vue-next'
import { useAppStore } from '@/stores/app'
import { ApiError } from '@/api/client'

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
const title = computed(() => (isSetup.value ? '首次启动' : 'Dockhelm'))
const subtitle = computed(() => (isSetup.value ? '设置访问密码后才能使用' : '请输入密码以继续'))

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

const canSubmit = computed(() => {
  if (busy.value || locked.value) return false
  if (isSetup.value) return password.value.length >= app.minPasswordLength && password.value === confirm.value
  return password.value.length > 0
})

/** 错误条语气：锁定 / 首次设置出错用红，登录失败可重试用琥珀（与设计稿一致）。 */
const bannerTone = computed(() =>
  locked.value || isSetup.value ? 'dh-banner-err' : 'dh-banner-warn',
)

const bannerText = computed(() => {
  if (locked.value) {
    // 光说「请稍后再试」等于没说：把服务端的「约 N 分钟」和本地倒计时一起摆出来
    return lockLeft.value > 0
      ? `${errorMsg.value}（剩余 ${mmss(lockLeft.value)}）`
      : errorMsg.value
  }
  if (!errorMsg.value) return ''
  if (remaining.value !== null && remaining.value > 0) {
    return `${errorMsg.value}，还可尝试 ${remaining.value} 次`
  }
  return errorMsg.value
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

function onEnter() {
  void submit()
}
</script>

<template>
  <div class="flex min-h-screen flex-col bg-ink-900">
    <!-- 全宽顶栏（与其它页面同一套外壳） -->
    <header class="flex flex-none items-center gap-2.5 border-b border-line-2 bg-ink-850 px-[18px] py-3">
      <!--
        顶栏标必须与 AppShell 完全一致（同一套罗盘标）。
        这里原来落着一个字母「D」—— 2026-10-09 换标时漏了这一处（当时只算了 5 个落点）。
      -->
      <div
        class="grid h-[27px] w-[27px] flex-none place-items-center rounded-[9px] bg-accent text-accent-ink"
      >
        <svg viewBox="0 0 32 32" class="h-[15px] w-[15px]">
          <circle cx="16" cy="16" r="12.3" fill="none" stroke="currentColor" stroke-width="2.6" />
          <path
            fill-rule="evenodd"
            fill="currentColor"
            d="M16 6.6 18.5 13.5 25.4 16 18.5 18.5 16 25.4 13.5 18.5 6.6 16 13.5 13.5ZM18 16A2 2 0 1 0 14 16A2 2 0 1 0 18 16Z"
          />
        </svg>
      </div>
      <div class="text-[14px] font-semibold tracking-[.2px]">Dockhelm</div>
      <div class="text-[12px] text-text-4">登录与账户</div>
    </header>

    <div class="relative flex flex-1 flex-col items-center justify-center gap-[14px] px-5 py-10">
      <!-- 顶部柔光，对应设计稿 .lgpane.soft -->
      <div
        class="pointer-events-none absolute inset-x-0 top-0 h-[420px]"
        style="background: radial-gradient(560px 260px at 50% 0%, var(--color-glow), transparent 72%)"
      />

      <div class="relative flex w-[312px] flex-col gap-[13px] rounded-[16px] border border-line-1 bg-ink-700 p-[22px]">
        <div class="dh-lglogo">
          <svg viewBox="0 0 32 32" class="h-[23px] w-[23px]">
            <circle cx="16" cy="16" r="12.3" fill="none" stroke="currentColor" stroke-width="2.5" />
            <path
              fill-rule="evenodd"
              fill="currentColor"
              d="M16 6.6 18.5 13.5 25.4 16 18.5 18.5 16 25.4 13.5 18.5 6.6 16 13.5 13.5ZM18 16A2 2 0 1 0 14 16A2 2 0 1 0 18 16Z"
            />
          </svg>
        </div>

        <div class="text-center">
          <div class="text-[16px] font-semibold">{{ title }}</div>
          <div class="mt-[3px] text-[12px] text-text-4">{{ subtitle }}</div>
        </div>

        <div v-if="mode === 'loading'" class="grid h-[110px] place-items-center text-text-5">
          <Loader2 class="h-5 w-5 dh-spin" />
        </div>

        <template v-else>
          <!-- 首次启动 -->
          <template v-if="isSetup">
            <div>
              <label class="mb-[5px] block text-[12px] text-text-4">新密码</label>
              <div class="dh-field">
                <input
                  v-model="password"
                  :type="showPw ? 'text' : 'password'"
                  :placeholder="`至少 ${app.minPasswordLength} 位`"
                  autocomplete="new-password"
                  autofocus
                  @keyup.enter="onEnter"
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
              <div class="mt-[7px] flex gap-1">
                <i
                  v-for="i in 4"
                  :key="i"
                  class="h-[4px] flex-1 rounded-[2px]"
                  :class="i <= strength ? 'bg-run' : 'bg-line-1'"
                />
              </div>
              <div class="mt-[5px] text-[11.5px] text-text-5">
                强度：{{ strengthLabel }} · 至少 {{ app.minPasswordLength }} 位，建议含数字与符号
              </div>
            </div>

            <div>
              <label class="mb-[5px] block text-[12px] text-text-4">确认密码</label>
              <div class="dh-field">
                <input
                  v-model="confirm"
                  :type="showPw ? 'text' : 'password'"
                  placeholder="再输一次"
                  autocomplete="new-password"
                  @keyup.enter="onEnter"
                />
              </div>
              <div v-if="confirm && confirm !== password" class="mt-[5px] text-[11.5px] text-err-text">
                两次输入不一致
              </div>
            </div>
          </template>

          <!-- 正常登录 -->
          <template v-else>
            <div class="dh-field">
              <input
                v-model="password"
                :type="showPw ? 'text' : 'password'"
                placeholder="请输入密码"
                autocomplete="current-password"
                autofocus
                @keyup.enter="onEnter"
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
          </template>

          <div
            v-if="bannerText"
            class="dh-banner !gap-2.5 !py-2 !pl-[11px] !pr-[11px] !text-[12px]"
            :class="bannerTone"
          >
            <span class="h-[13px] w-[13px] flex-none rounded-[4px] bg-current opacity-50" />
            {{ bannerText }}
          </div>

          <label
            v-if="!isSetup"
            class="flex cursor-pointer items-center gap-2 text-[12px] text-text-3"
          >
            <input v-model="keep" type="checkbox" class="h-[14px] w-[14px] accent-accent" />
            保持登录（7 天）
          </label>

          <button
            type="button"
            class="w-full rounded-[10px] bg-accent px-3 py-2.5 text-center text-[13px] font-semibold text-accent-ink transition hover:brightness-110 disabled:cursor-not-allowed disabled:opacity-45"
            :disabled="!canSubmit"
            @click="submit"
          >
            {{ busy ? '请稍候…' : locked ? `已锁定 ${mmss(lockLeft)}` : isSetup ? '完成并进入' : '登 录' }}
          </button>

          <div v-if="!isSetup" class="text-center text-[11.5px] text-text-6">
            连续错 {{ app.maxFailures }} 次锁定 5 分钟
          </div>
        </template>
      </div>

      <div v-if="isSetup" class="relative w-[312px] text-center text-[11.5px] leading-[1.7] text-text-6">
        密码只存 <span class="font-mono">bcrypt</span> 哈希，明文不落盘、也不写日志。
        <br />
        Dockhelm 持有 Docker 套接字，等于拥有宿主机 root 权限，请务必设置密码。
      </div>
    </div>
  </div>
</template>
