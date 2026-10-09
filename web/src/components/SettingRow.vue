<script setup lang="ts">
/**
 * 一行设置项：左侧标题 + 说明，右侧放控件（开关 / 下拉 / 数值 / 徽标）。
 * 设计稿里这类「crow」重复了十几处，抽一个组件保证间距与排版完全一致。
 */
defineProps<{
  title: string
  sub?: string
}>()
</script>

<template>
  <div class="flex items-center gap-3">
    <div class="min-w-0 flex-1">
      <div class="text-[12.5px] leading-snug text-text-1">{{ title }}</div>
      <div v-if="sub" class="mt-[2px] text-[11px] leading-relaxed text-text-5">{{ sub }}</div>
    </div>
    <div class="flex flex-none items-center gap-2">
      <slot />
    </div>
  </div>
</template>
