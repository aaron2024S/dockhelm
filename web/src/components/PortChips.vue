<script setup lang="ts">
/**
 * 端口徽标组。
 *
 * 两种端口在视觉上必须能一眼分开：
 *   - 已发布到宿主机：实底 + 描边，宿主端口加粗（这是用户真正要访问的入口）
 *   - 只在容器网络里：虚线框（没有对外暴露，别让人误以为能连上）
 * 超过 max 条收成 +N，鼠标悬停列出全部。
 */
import { computed } from 'vue'
import type { PortView } from '@/api/types'

const props = withDefaults(
  defineProps<{
    ports: PortView[]
    /** 最多显示几条，超出收成 +N。传 0 表示全部显示。 */
    max?: number
  }>(),
  { max: 3 },
)

const visible = computed(() => (props.max > 0 ? props.ports.slice(0, props.max) : props.ports))
const hidden = computed(() => (props.max > 0 ? props.ports.length - props.max : 0))

function keyOf(p: PortView) {
  return `${p.hostIp}:${p.hostPort}>${p.innerPort}/${p.proto}`
}

function tipOf(p: PortView) {
  return p.published
    ? `已发布到宿主机 ${p.hostIp || '0.0.0.0'}:${p.hostPort} → 容器内 ${p.innerPort}/${p.proto}`
    : `仅容器内可见（未发布到宿主机）：${p.innerPort}/${p.proto}`
}

const allTip = computed(() =>
  props.ports
    .map((p) => (p.published ? `${p.hostIp || '0.0.0.0'}:${p.hostPort} → ${p.innerPort}/${p.proto}` : `${p.innerPort}/${p.proto}（仅容器内）`))
    .join('\n'),
)
</script>

<template>
  <div v-if="ports.length" class="flex flex-wrap items-center gap-1">
    <span
      v-for="p in visible"
      :key="keyOf(p)"
      class="inline-flex items-center gap-1 rounded-md px-1.5 py-[2px] font-mono text-[10.5px]"
      :class="
        p.published
          ? 'bg-ink-850 text-text-2 ring-1 ring-inset ring-line-3'
          : 'border border-dashed border-line-3 text-text-5'
      "
      :title="tipOf(p)"
    >
      <template v-if="p.published">
        <span class="font-semibold text-text-1">{{ p.hostPort }}</span>
        <span class="text-text-6">→</span>
        <span>{{ p.innerPort }}/{{ p.proto }}</span>
      </template>
      <template v-else>
        <span>{{ p.innerPort }}/{{ p.proto }}</span>
      </template>
    </span>
    <span v-if="hidden > 0" class="px-1 text-[10.5px] text-text-5" :title="allTip">
      +{{ hidden }}
    </span>
  </div>
</template>
