<script setup lang="ts">
import { computed, onMounted } from 'vue'
import { RouterView, useRoute } from 'vue-router'
import AppShell from '@/components/AppShell.vue'
import Toasts from '@/components/Toasts.vue'
import { useAppStore } from '@/stores/app'

const app = useAppStore()
const route = useRoute()

const isPublic = computed(() => route.meta.public === true)
const showShell = computed(() => !isPublic.value && app.loggedIn)

onMounted(() => {
  void app.bootstrap()
})
</script>

<template>
  <div class="min-h-full">
    <!-- 启动时先显示一个极简的加载态，避免登录页闪一下再跳走 -->
    <div v-if="!app.ready" class="fixed inset-0 grid place-items-center bg-ink-950">
      <div class="flex flex-col items-center gap-3">
        <div class="grid h-11 w-11 place-items-center rounded-[13px] bg-accent text-accent-ink">
          <svg viewBox="0 0 32 32" class="h-6 w-6" fill="none" stroke="currentColor" stroke-width="2.4">
            <path d="M16 7l7 4v10l-7 4-7-4V11z" stroke-linejoin="round" />
            <path d="M16 15v10M9 11l7 4 7-4" stroke-linejoin="round" />
          </svg>
        </div>
        <div class="text-[12px] text-text-5">正在载入 Dockhelm…</div>
      </div>
    </div>

    <AppShell v-else-if="showShell">
      <RouterView v-slot="{ Component }">
        <component :is="Component" />
      </RouterView>
    </AppShell>

    <RouterView v-else />

    <Toasts />
  </div>
</template>
