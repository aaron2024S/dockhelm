<script setup lang="ts">
/**
 * 一行设置项：左侧标题 + 说明，右侧放控件（开关 / 下拉 / 数值 / 徽标）。
 * 设计稿里这类「crow」重复了十几处，抽一个组件保证间距与排版完全一致。
 * stack：控件较宽（如 URL 输入框）时在手机端改为上下堆叠，
 * 否则固定宽度的控件会把左侧文字挤成一个字一行的竖排。
 */
defineProps<{
  title: string
  sub?: string
  stack?: boolean
}>()
</script>

<template>
  <div class="flex items-center gap-3" :class="stack ? 'max-md:flex-col max-md:items-stretch' : ''">
    <div class="min-w-0 flex-1">
      <div class="text-[12.5px] leading-snug text-text-1">{{ title }}</div>
      <div v-if="sub" class="mt-[2px] text-[11px] leading-relaxed text-text-5">{{ sub }}</div>
    </div>
    <div class="flex flex-none items-center gap-2" :class="stack ? 'max-md:w-full' : ''">
      <slot />
    </div>
  </div>
</template>
