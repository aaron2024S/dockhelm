<script setup lang="ts">
import { computed, onMounted, ref } from 'vue'
import { useRouter } from 'vue-router'
import { Loader2 } from 'lucide-vue-next'
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
  if (!errorMsg.value) return ''
  if (!locked.value && remaining.value !== null && remaining.value > 0) {
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
      errorMsg.value = e.message
      if (e.code === 'locked') locked.value = true
      const payload = e.payload as { remaining?: number } | undefined
      if (typeof payload?.remaining === 'number') remaining.value = payload.remaining
      if (e.code === 'bad_password') password.value = ''
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
      <div class="grid h-[27px] w-[27px] flex-none place-items-center rounded-[9px] bg-accent text-[14px] font-bold text-accent-ink">
        D
      </div>
      <div class="text-[14px] font-semibold tracking-[.2px]">Dockhelm</div>
      <div class="text-[12px] text-text-4">登录与账户</div>
      <div class="ml-auto">
        <div class="dh-avatar" title="登录后可在这里进入账户设置">A</div>
      </div>
    </header>

    <div class="relative flex flex-1 flex-col items-center justify-center gap-[14px] px-5 py-10">
      <!-- 顶部柔光，对应设计稿 .lgpane.soft -->
      <div
        class="pointer-events-none absolute inset-x-0 top-0 h-[420px]"
        style="background: radial-gradient(560px 260px at 50% 0%, rgba(45, 212, 191, 0.1), transparent 72%)"
      />

      <div class="relative flex w-[312px] flex-col gap-[13px] rounded-[16px] border border-line-1 bg-ink-700 p-[22px]">
        <div class="dh-lglogo">D</div>

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
                <button type="button" class="dh-eye" @click="showPw = !showPw">
                  {{ showPw ? '隐藏' : '显示' }}
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
              <div v-if="confirm && confirm !== password" class="mt-[5px] text-[11.5px] text-[#fca5a5]">
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
              <button type="button" class="dh-eye" @click="showPw = !showPw">
                {{ showPw ? '隐藏' : '显示' }}
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
            <input v-model="keep" type="checkbox" class="h-[14px] w-[14px] accent-[#2dd4bf]" />
            保持登录（7 天，勾选延长到 30 天）
          </label>

          <button
            type="button"
            class="w-full rounded-[10px] bg-accent px-3 py-2.5 text-center text-[13px] font-semibold text-accent-ink transition hover:brightness-110 disabled:cursor-not-allowed disabled:opacity-45"
            :disabled="!canSubmit"
            @click="submit"
          >
            {{ busy ? '请稍候…' : isSetup ? '完成并进入' : '登 录' }}
          </button>

          <div v-if="!isSetup" class="text-center text-[11.5px] text-text-6">
            连续错 {{ app.maxFailures }} 次锁定 5 分钟 · 每次失败延迟 1 秒
          </div>
        </template>
      </div>

      <div class="relative w-[312px] text-center text-[11.5px] leading-[1.7] text-text-6">
        <template v-if="isSetup">
          密码只存 <span class="font-mono">bcrypt</span> 哈希，明文不落盘、也不写日志。
          <br />
          Dockhelm 持有 Docker 套接字，等于拥有宿主机 root 权限，请务必设置密码。
        </template>
        <template v-else>
          会话用 HttpOnly Cookie，有效期 7 天；勾「保持登录」延长到 30 天。
        </template>
      </div>
    </div>
  </div>
</template>
