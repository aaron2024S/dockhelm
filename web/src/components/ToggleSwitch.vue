<script setup lang="ts">
/**
 * 开关。用 design token 里的 .dh-switch，配 v-model 直接用。
 * 之所以不用原生 checkbox：设计稿里开关是「轨道 + 圆点」，
 * 而原生 checkbox 在暗色下各浏览器渲染差异太大。
 */
const model = defineModel<boolean>({ required: true })

defineProps<{
  disabled?: boolean
  /** 无障碍标签：开关自身没有文字，读屏需要知道它控制什么。 */
  label?: string
}>()
</script>

<template>
  <button
    type="button"
    class="dh-switch"
    role="switch"
    :aria-checked="model"
    :aria-label="label"
    :data-on="model"
    :disabled="disabled"
    @click="model = !model"
  />
</template>
