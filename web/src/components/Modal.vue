<script setup lang="ts">
import { X } from 'lucide-vue-next'

withDefaults(
  defineProps<{
    open: boolean
    title?: string
    subtitle?: string
    width?: string
    closeOnBackdrop?: boolean
    /**
     * 提交中：屏蔽「点遮罩关闭」和右上角 × 关闭。
     * 以前删除/还原/清理这类会真动数据的弹窗，请求还没回来也能被关掉 ——
     * 关掉后结果无从得知，用户会以为没执行而再点一次。
     */
    busy?: boolean
  }>(),
  { title: '', subtitle: '', width: '520px', closeOnBackdrop: true, busy: false },
)

const emit = defineEmits<{ close: [] }>()
</script>

<template>
  <Teleport to="body">
    <Transition
      enter-active-class="transition duration-150 ease-out"
      enter-from-class="opacity-0"
      leave-active-class="transition duration-120 ease-in"
      leave-to-class="opacity-0"
    >
      <div
        v-if="open"
        class="fixed inset-0 z-[70] grid place-items-center overflow-y-auto bg-black/60 p-4 backdrop-blur-[2px]"
        @click.self="closeOnBackdrop && !busy && emit('close')"
      >
        <div
          class="dh-card w-full shadow-[var(--shadow-modal)]"
          :style="{ maxWidth: width }"
          role="dialog"
          aria-modal="true"
        >
          <div class="flex items-start gap-3 border-b border-line-1 px-4 py-3">
            <div class="min-w-0 flex-1">
              <div class="truncate text-[13.5px] font-semibold">{{ title }}</div>
              <div v-if="subtitle" class="mt-0.5 text-[11.5px] leading-relaxed text-text-4">
                {{ subtitle }}
              </div>
            </div>
            <button
              type="button"
              :disabled="busy"
              class="dh-tap grid h-6 w-6 flex-none place-items-center rounded-md text-text-5 transition-colors hover:bg-ink-650 hover:text-text-1 disabled:cursor-not-allowed disabled:opacity-40 disabled:hover:bg-transparent"
              @click="emit('close')"
            >
              <X class="h-3.5 w-3.5" />
            </button>
          </div>
          <div class="px-4 py-4">
            <slot />
          </div>
          <div v-if="$slots.footer" class="flex items-center justify-end gap-2 border-t border-line-1 px-4 py-3">
            <slot name="footer" />
          </div>
        </div>
      </div>
    </Transition>
  </Teleport>
</template>
