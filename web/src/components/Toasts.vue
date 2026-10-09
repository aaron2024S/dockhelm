<script setup lang="ts">
import { CheckCircle2, Info, X, XCircle } from 'lucide-vue-next'
import { useToastStore } from '@/stores/toast'

const toast = useToastStore()

const icons = {
  success: CheckCircle2,
  error: XCircle,
  info: Info,
} as const

const tones = {
  success: 'border-line-ok text-run-text',
  error: 'border-line-err text-err-text',
  info: 'border-accent-line text-accent',
} as const</script>

<template>
  <div class="pointer-events-none fixed bottom-5 right-5 z-[80] flex w-[330px] flex-col gap-2">
    <TransitionGroup
      enter-active-class="transition duration-200 ease-out"
      enter-from-class="translate-y-2 opacity-0"
      leave-active-class="transition duration-150 ease-in"
      leave-to-class="translate-y-1 opacity-0"
    >
      <div
        v-for="item in toast.items"
        :key="item.id"
        class="pointer-events-auto flex items-start gap-2.5 rounded-xl border bg-ink-700/97 px-3.5 py-3 shadow-[var(--shadow-toast)] backdrop-blur"
        :class="tones[item.kind]"
      >
        <component :is="icons[item.kind]" class="mt-[1px] h-4 w-4 flex-none" />
        <div class="min-w-0 flex-1">
          <div class="text-[12.5px] font-medium text-text-1">{{ item.message }}</div>
          <div v-if="item.detail" class="mt-0.5 break-words text-[11.5px] leading-relaxed text-text-4">
            {{ item.detail }}
          </div>
        </div>
        <button
          type="button"
          class="dh-tap -m-1.5 flex-none rounded-md p-1.5 text-text-5 transition-colors hover:bg-ink-650 hover:text-text-1"
          aria-label="关闭提示"
          @click="toast.dismiss(item.id)"
        >
          <X class="h-3.5 w-3.5" />
        </button>
      </div>
    </TransitionGroup>
  </div>
</template>
