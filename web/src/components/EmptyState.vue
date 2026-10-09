<script setup lang="ts">
import type { Component } from 'vue'

withDefaults(
  defineProps<{
    icon?: Component
    title: string
    description?: string
    actionLabel?: string
  }>(),
  { icon: undefined, description: '', actionLabel: '' },
)
const emit = defineEmits<{ action: [] }>()
</script>

<template>
  <div class="flex flex-col items-center justify-center gap-2 px-6 py-12 text-center">
    <div v-if="icon" class="grid h-10 w-10 place-items-center rounded-xl border border-line-1 bg-ink-750 text-text-5">
      <component :is="icon" class="h-[18px] w-[18px]" />
    </div>
    <div class="text-[13px] font-medium text-text-2">{{ title }}</div>
    <div v-if="description" class="max-w-[420px] text-[12px] leading-relaxed text-text-5">
      {{ description }}
    </div>
    <button v-if="actionLabel" type="button" class="dh-btn mt-1.5" @click="emit('action')">
      {{ actionLabel }}
    </button>
  </div>
</template>
