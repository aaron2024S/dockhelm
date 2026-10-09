<script setup lang="ts">
import { computed, onMounted, ref } from 'vue'
import { useRouter } from 'vue-router'
import { Eye, EyeOff, KeyRound, Loader2, ShieldCheck } from 'lucide-vue-next'
import { useAppStore } from '@/stores/app'
import { ApiError } from '@/api/client'
import { version as brandVersion } from '@/config'

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
const title = computed(() => (isSetup.value ? '设置访问密码' : '登录 Dockhelm'))
const subtitle = computed(() =>
  isSetup.value
    ? 'Dockhelm 持有 Docker 套接字，等于拥有宿主机 root 权限。首次使用必须先设置密码，否则局域网内任何人都能删除你的容器。'
    : '输入访问密码以继续。',
)

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
const strengthLabel = computed(() => ['', '偏弱', '一般', '不错', '很强'][strength.value])

const canSubmit = computed(() => {
  if (busy.value || locked.value) return false
  if (isSetup.value) return password.value.length >= app.minPasswordLength && password.value === confirm.value
  return password.value.length > 0
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
  <div class="grid min-h-screen place-items-center bg-ink-900 px-5 py-10">
    <!-- 背景光晕，和设计稿的 soft 面板一致 -->
    <div
      class="pointer-events-none absolute inset-x-0 top-0 h-[420px]"
      style="background: radial-gradient(680px 300px at 50% 0%, rgba(45, 212, 191, 0.1), transparent 72%)"
    />

    <div class="relative w-full max-w-[380px]">
      <div class="dh-card flex flex-col gap-3.5 p-6">
        <div class="grid h-[42px] w-[42px] place-items-center self-center rounded-[13px] bg-accent text-accent-ink">
          <svg viewBox="0 0 32 32" class="h-5 w-5" fill="none" stroke="currentColor" stroke-width="2.6">
            <path d="M16 7l7 4v10l-7 4-7-4V11z" stroke-linejoin="round" />
            <path d="M16 15v10M9 11l7 4 7-4" stroke-linejoin="round" />
          </svg>
        </div>

        <div class="text-center">
          <div class="text-[16px] font-semibold">{{ title }}</div>
          <div class="mt-1 whitespace-pre-line text-[12px] leading-relaxed text-text-4">
            {{ subtitle }}
          </div>
        </div>

        <div v-if="mode === 'loading'" class="grid h-24 place-items-center text-text-5">
          <Loader2 class="h-5 w-5 dh-spin" />
        </div>

        <template v-else>
          <div>
            <label class="dh-label">访问密码</label>
            <div class="relative">
              <input
                v-model="password"
                :type="showPw ? 'text' : 'password'"
                class="dh-input !py-[10px] !pr-10 !tracking-[2px]"
                :placeholder="isSetup ? `至少 ${app.minPasswordLength} 位` : '请输入密码'"
                autocomplete="current-password"
                autofocus
                @keyup.enter="onEnter"
              />
              <button
                type="button"
                class="absolute right-2 top-1/2 grid h-6 w-6 -translate-y-1/2 place-items-center rounded-md text-text-5 hover:text-text-1"
                @click="showPw = !showPw"
              >
                <component :is="showPw ? EyeOff : Eye" class="h-3.5 w-3.5" />
              </button>
            </div>
            <div v-if="isSetup" class="mt-[7px] flex gap-1">
              <i
                v-for="i in 4"
                :key="i"
                class="h-[3px] flex-1 rounded-full"
                :class="i <= strength ? 'bg-[#34d399]' : 'bg-line-1'"
              />
            </div>
            <div v-if="isSetup && password" class="mt-1 text-[11px] text-text-5">
              强度：{{ strengthLabel }}
            </div>
          </div>

          <div v-if="isSetup">
            <label class="dh-label">再输一次</label>
            <input
              v-model="confirm"
              :type="showPw ? 'text' : 'password'"
              class="dh-input !py-[10px] !tracking-[2px]"
              placeholder="重复上面的密码"
              autocomplete="new-password"
              @keyup.enter="onEnter"
            />
            <div
              v-if="confirm && confirm !== password"
              class="mt-1 text-[11px] text-[#fca5a5]"
            >
              两次输入不一致
            </div>
          </div>

          <label v-else class="flex cursor-pointer items-center gap-2 text-[12px] text-text-3">
            <input v-model="keep" type="checkbox" class="h-[14px] w-[14px] accent-[#2dd4bf]" />
            保持登录（30 天内免登录）
          </label>

          <div
            v-if="errorMsg"
            class="rounded-[9px] border border-[rgba(248,113,113,.35)] bg-[rgba(248,113,113,.08)] px-3 py-2 text-[12px] leading-relaxed text-[#fca5a5]"
          >
            {{ errorMsg }}
            <span v-if="!locked && remaining !== null && remaining > 0" class="text-text-4">
              （还可尝试 {{ remaining }} 次）
            </span>
          </div>

          <button
            type="button"
            class="dh-btn dh-btn-primary !py-[10px] !text-[13px]"
            :disabled="!canSubmit"
            @click="submit"
          >
            <Loader2 v-if="busy" class="h-3.5 w-3.5 dh-spin" />
            <KeyRound v-else class="h-3.5 w-3.5" />
            {{ isSetup ? '设置并进入' : '登录' }}
          </button>

          <div v-if="isSetup" class="flex items-start gap-2 rounded-[9px] border border-line-1 bg-ink-800 px-3 py-2.5">
            <ShieldCheck class="mt-[1px] h-3.5 w-3.5 flex-none text-text-5" />
            <div class="text-[11.5px] leading-relaxed text-text-5">
              密码只保存 bcrypt 哈希，明文不落盘、不写日志。
              <br />
              忘了密码：在 NAS 面板里删除 <code class="text-text-3">data/auth.json</code> 后重启容器即可重设；
              也可以在 compose 里加一行 <code class="text-text-3">DOCKHELM_PASSWORD=新密码</code> 强制覆盖。
            </div>
          </div>
        </template>
      </div>

      <div class="mt-3 text-center text-[11px] text-text-6">
        Dockhelm v{{ brandVersion }} · 自托管 Docker 容器管理面板
      </div>
    </div>
  </div>
</template>
