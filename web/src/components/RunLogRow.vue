<script setup lang="ts">
/**
 * 运行记录的一行 —— 总览「最近执行记录」与设置页「运行记录」**共用同一个组件**。
 *
 * 用户 2026-10-13：「运行记录其实就是最近执行记录更多的展示，最近执行记录都是一行，
 * 为什么运行记录要搞 2 行，和最近执行记录的显示内容保持一致就好了」。所以两处渲染
 * 就该是同一份代码 —— 设置页以前是表格（说明一行、对象名一行），列和徽标都对不上。
 *
 * 列：时间 | 类型徽标 | 结果徽标 | 「对象名 + 说明」一行（超长截断，title 出全文）。
 * 对象名（ref）只在一次操作涉及多个对象时才有（清理镜像的名单、全量备份的失败名单），
 * 单个对象的记录把名字写进说明里，见 utils/format.ts 的 runRefShown。
 */
import type { RunLog } from '@/api/types'
import { formatDayTime, runKindLabel, runRefShown } from '@/utils/format'

const props = defineProps<{ log: RunLog }>()

/** 时间列：今天只出 HH:mm，昨天带前缀，更早出日期。 */
function logTime(ts: string) {
  const t = Date.parse(ts)
  if (Number.isNaN(t)) return '—'
  const d = new Date(t)
  const p = (n: number) => String(n).padStart(2, '0')
  const hm = `${p(d.getHours())}:${p(d.getMinutes())}`
  const today0 = new Date()
  today0.setHours(0, 0, 0, 0)
  const day0 = new Date(t)
  day0.setHours(0, 0, 0, 0)
  const diff = Math.round((day0.getTime() - today0.getTime()) / 86400000)
  if (diff === 0) return hm
  if (diff === -1) return `昨天 ${hm}`
  return formatDayTime(ts, 'date')
}

/** 类型徽标配色：失败红、计划任务紫、其余中性。 */
function kindClass(kind: string, status: string) {
  if (status === 'failed') return 'dh-badge-err'
  if (kind === 'schedule') return 'dh-badge-accent'
  return 'dh-badge-plain'
}

function resultBadge(l: RunLog) {
  if (l.status === 'success') return { text: '成功', cls: 'dh-badge-run' }
  if (l.status === 'up_to_date') return { text: '跳过', cls: 'dh-badge-accent' }
  if (l.status === 'failed') return { text: '失败', cls: 'dh-badge-err' }
  if (l.status === 'done') return { text: '完成', cls: 'dh-badge-run' }
  return { text: l.status || '—', cls: 'dh-badge-plain' }
}

function refShown(): boolean {
  return runRefShown(props.log.kind, props.log.ref)
}
</script>

<template>
  <div class="dh-tl">
    <div class="w-[86px] flex-none text-[11.5px] text-text-5 max-md:order-1">
      {{ logTime(log.ts) }}
    </div>
    <!-- 徽标定宽居中：行与行之间徽标列起点对齐，
         不因「更新/镜像」两字与「自动巡检/计划任务」四字之差而参差。 -->
    <span
      class="dh-badge min-w-[64px] flex-none justify-center max-md:order-2"
      :class="kindClass(log.kind, log.status)"
    >
      {{ runKindLabel(log.kind) }}
    </span>
    <span
      class="dh-badge min-w-[42px] flex-none justify-center max-md:order-3"
      :class="resultBadge(log).cls"
    >
      {{ resultBadge(log).text }}
    </span>
    <!-- 对象名 + 说明是同一列的整体：对象名作说明格内的等宽前缀，中间不留空档。 -->
    <span
      class="min-w-0 flex-1 truncate text-[12.5px] text-text-2 max-md:order-4"
      :title="refShown() ? `${log.ref} ${log.message}` : log.message"
    >
      <span v-if="refShown()" class="mr-1.5 font-mono text-[11.5px] text-text-3">{{ log.ref }}</span>{{ log.message }}
    </span>
  </div>
</template>
