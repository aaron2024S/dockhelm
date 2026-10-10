<script setup lang="ts">
import { computed, onMounted, onUnmounted, ref } from 'vue'
import { RouterLink } from 'vue-router'
import {
  AlertTriangle,
  CheckCircle2,
  Clock,
  Download,
  HelpCircle,
  Info,
  Loader2,
  Play,
  RefreshCw,
  ShieldQuestion,
  Zap,
} from 'lucide-vue-next'
import { api, openStream } from '@/api/client'
import type { AutoUpdateInfo, CheckResult, UpdatesResponse } from '@/api/types'
import { checkLabel, formatDateTime, relativeTime, shortImage } from '@/utils/format'
import { useToastStore } from '@/stores/toast'
import Modal from '@/components/Modal.vue'

const toast = useToastStore()
const results = ref<CheckResult[]>([])
const checkedAt = ref('')
const loading = ref(false)
const deepLoading = ref(false)
const progress = ref<{ name: string; message: string; status: string }[]>([])
const showInfo = ref(false)
let closeStream: (() => void) | null = null
/** 组件已卸载 —— 用来打断 runAutoCycle 里那个最多 60 秒的轮询循环。 */
let disposed = false

// —— 自动更新 ——
// 注意：本页**不**维护自动更新开关与执行策略（pullOnce / backupBefore / cleanupAfter /
// concurrency）——那些都是 /api/settings 里的同一批字段，唯一入口在「设置 → 更新与检测」。
// 这里只读 /api/updates/auto 的状态（enabled / 下次巡检 / 本轮清单 / 上一轮结果）。
const auto = ref<AutoUpdateInfo | null>(null)
const autoBusy = ref(false)

const available = computed(() => results.value.filter((r) => r.status === 'update_available'))

/** 按容器名索引巡检结果 —— 候选表里给「未检测到更新」的行补上具体判定（已是最新/无法判定/本地镜像）。 */
const resultByName = computed(() => {
  const m = new Map<string, CheckResult>()
  for (const r of results.value) m.set(r.container, r)
  return m
})

/** 会被自动更新的那几个容器（预览）。 */
const willUpdate = computed(() => (auto.value?.candidates ?? []).filter((c) => c.willUpdate))
/** 被跳过的（含排除列表、自身、以及未检测到更新的）。 */
const skipped = computed(() => (auto.value?.candidates ?? []).filter((c) => !c.willUpdate && (c.protected || c.excluded)))

async function load() {
  try {
    const res = await api.get<UpdatesResponse>('/api/updates')
    results.value = res.results ?? []
    checkedAt.value = res.checkedAt
  } catch (e) {
    toast.error('读取巡检结果失败', e instanceof Error ? e.message : String(e))
  }
}

async function check(deep: boolean) {
  if (deep) deepLoading.value = true
  else loading.value = true
  try {
    const res = await api.post<{ results: CheckResult[] }>('/api/updates/check', { deep })
    results.value = res.results ?? []
    checkedAt.value = new Date().toISOString()
    const n = (res.results ?? []).filter((r) => r.status === 'update_available').length
    toast.success(deep ? '深度检测完成' : '巡检完成', `发现 ${n} 个有可用更新`)
  } catch (e) {
    toast.error('巡检失败', e instanceof Error ? e.message : String(e))
  } finally {
    deepLoading.value = false
    loading.value = false
  }
}

const toneOf = (status: string) => checkLabel(status).tone

// ---------- 自动更新 ----------

async function loadAuto() {
  try {
    auto.value = await api.get<AutoUpdateInfo>('/api/updates/auto')
  } catch {
    // 自动更新信息读不到不影响主流程（例如守护进程暂时不可达）
    auto.value = null
  }
}

/** 立即跑一轮；dryRun 为真时只巡检不动容器。 */
async function runAutoCycle(dryRun: boolean) {
  autoBusy.value = true
  try {
    await api.post('/api/updates/auto-run', { dryRun })
    toast.info(dryRun ? '已开始巡检（不动容器）' : '已开始自动更新', '完成后本页会自动刷新')
    // 后端在后台跑，轮询几次把结果拉回来。
    // 循环必须能被卸载打断：否则用户点完就离开页面，这个循环还会持续
    // 最多一分钟对着已卸载的组件写 state、继续发请求。
    for (let i = 0; i < 40 && !disposed; i++) {
      await new Promise((r) => window.setTimeout(r, 1500))
      if (disposed) return
      const info = await api.get<AutoUpdateInfo>('/api/updates/auto').catch(() => null)
      if (disposed) return
      if (info) auto.value = info
      if (info && !info.running) break
    }
    if (disposed) return
    await Promise.all([load(), loadAuto()])
    const last = auto.value?.lastRun
    if (last && !last.dryRun) {
      toast.success('自动更新完成', `更新 ${last.updated} 个，失败 ${last.failed} 个`)
    }
  } catch (e) {
    if (disposed) return
    toast.error('启动失败', e instanceof Error ? e.message : String(e))
  } finally {
    if (!disposed) autoBusy.value = false
  }
}

onMounted(() => {
  void load()
  void loadAuto()
  closeStream = openStream('/api/events/stream', (topic, ev) => {
    if (topic !== 'update') return
    const d = ev.data ?? {}
    const name = String(d.container ?? '')
    if (ev.kind === 'step' || ev.kind === 'container_status' || ev.kind === 'pull_progress') {
      const msg = String(d.message ?? d.status ?? '')
      if (!msg) return
      const idx = progress.value.findIndex((p) => p.name === name)
      const entry = { name, message: msg, status: String(ev.status ?? 'running') }
      if (idx >= 0) progress.value[idx] = entry
      else progress.value.unshift(entry)
      if (progress.value.length > 12) progress.value.pop()
    }
    if (ev.kind === 'container_status') {
      const st = String(d.status ?? '')
      if (['updated', 'failed', 'broken', 'up_to_date'].includes(st)) {
        const entry = progress.value.find((p) => p.name === name)
        if (entry) entry.status = st === 'updated' || st === 'up_to_date' ? 'success' : 'failed'
      }
    }
    if (ev.kind === 'batch_done') {
      // 手动批量更新的入口已从本页移除（更新容器走容器页或自动更新），
      // 但自动更新跑完也会发这个事件 —— 借它刷新一次巡检结果。
      void load()
    }
    // 自动更新：巡检结束 / 整轮结束都要把这块刷新一下
    if (ev.kind === 'auto_check_done' || ev.kind === 'auto_done') {
      void loadAuto()
      if (ev.kind === 'auto_done') void load()
    }
  })
})

onUnmounted(() => {
  disposed = true
  closeStream?.()
})
</script>

<template>
  <div class="flex flex-col gap-3.5 p-[18px]">
    <div class="dh-phead">
      <div class="dh-h1">更新中心</div>
      <button
        class="dh-tap inline-flex items-center gap-1 rounded-md px-1.5 py-[3px] text-[11.5px] text-text-5 transition-colors hover:bg-ink-750 hover:text-accent"
        title="为什么 Dockhelm 不会误停容器"
        @click="showInfo = true"
      >
        <Info class="h-3.5 w-3.5" />了解更多
      </button>
      <div class="dh-sub">
        {{ checkedAt ? `上次检测 ${relativeTime(checkedAt)} · 共比对 ${results.length} 个容器` : '还没有检测过' }}
      </div>
      <div class="ml-auto flex gap-2">
        <button class="dh-btn" :disabled="loading || deepLoading" @click="check(true)">
          <Zap class="h-3.5 w-3.5" :class="deepLoading ? 'dh-spin' : ''" />深度检测
        </button>
        <button class="dh-btn" :disabled="loading || deepLoading" @click="check(false)">
          <RefreshCw class="h-3.5 w-3.5" :class="loading ? 'dh-spin' : ''" />重新检测
        </button>
      </div>
    </div>

    <!-- 汇总 -->
    <div class="grid grid-cols-2 gap-3 lg:grid-cols-4">
      <div class="dh-card p-3.5">
        <div class="flex items-center gap-2 text-[12px] text-text-4"><Download class="h-3.5 w-3.5" />有可用更新</div>
        <div class="mt-1.5 text-[22px] font-semibold leading-none" :class="available.length ? 'text-warn-text' : ''">
          {{ available.length }}
        </div>
      </div>
      <div class="dh-card p-3.5">
        <div class="flex items-center gap-2 text-[12px] text-text-4"><CheckCircle2 class="h-3.5 w-3.5" />已是最新</div>
        <div class="mt-1.5 text-[22px] font-semibold leading-none">
          {{ results.filter((r) => r.status === 'up_to_date').length }}
        </div>
      </div>
      <div class="dh-card p-3.5">
        <div class="flex items-center gap-2 text-[12px] text-text-4"><ShieldQuestion class="h-3.5 w-3.5" />无法判定</div>
        <div class="mt-1.5 text-[22px] font-semibold leading-none">
          {{ results.filter((r) => r.status === 'unknown').length }}
        </div>
        <div class="mt-1 text-[11px] text-text-6">网络/认证问题导致，绝不当作有更新</div>
      </div>
      <div class="dh-card p-3.5">
        <div class="flex items-center gap-2 text-[12px] text-text-4"><HelpCircle class="h-3.5 w-3.5" />本地镜像</div>
        <div class="mt-1.5 text-[22px] font-semibold leading-none">
          {{ results.filter((r) => r.status === 'no_upstream').length }}
        </div>
        <div class="mt-1 text-[11px] text-text-6">本地构建，无远端可比对</div>
      </div>
    </div>

    <!-- 进度：自动更新（以及容器页发起的更新）的实时步骤都走这条 SSE -->
    <div v-if="progress.length" class="dh-card">
      <div class="dh-card-head">
        <Loader2 v-if="autoBusy" class="h-3.5 w-3.5 dh-spin text-accent" />
        <Zap v-else class="h-3.5 w-3.5 text-text-4" />
        <span>{{ autoBusy ? '更新进行中' : '最近一次更新进度' }}</span>
        <span class="ml-auto text-[11.5px] font-normal text-text-5">{{ progress.length }} 条</span>
      </div>
      <div class="flex flex-col">
        <div
          v-for="(p, i) in progress"
          :key="i"
          class="flex items-center gap-2.5 border-b border-line-row px-3.5 py-2 text-[12px] last:border-b-0"
        >
          <span
            class="h-[6px] w-[6px] flex-none rounded-full"
            :class="p.status === 'failed' ? 'bg-err' : p.status === 'success' ? 'bg-run' : 'bg-accent'"
          />
          <span class="w-[170px] flex-none truncate font-mono text-[11.5px] text-text-3">{{ p.name }}</span>
          <span class="min-w-0 flex-1 truncate text-text-2">{{ p.message }}</span>
        </div>
      </div>
    </div>

    <!--
      「有可用更新的容器」卡已删（它带全选/勾选/「更新选中的」一整套手动批量更新操作）：
      同一批容器的判定信息已经在下面自动更新卡的候选表里（candidates 覆盖全部容器），
      再列一遍纯属重复；手动更新走容器页（候选表的容器名可直接点过去），本页只做预览与自动更新。
    -->

    <!-- 自动更新 -->
    <!--
      卡头一行说清三件事：开关状态（点徽标直达设置）、下次巡检时间、手动触发按钮。
      原先这里有一个两栏的「总开关 / 下次巡检」说明区，信息与本行完全重复，已删。
    -->
    <div class="dh-card">
      <div class="dh-card-head">
        <Zap class="h-3.5 w-3.5" :class="auto?.enabled ? 'text-warn-text' : 'text-text-4'" />
        <span>自动更新</span>
        <RouterLink
          to="/settings"
          class="dh-tap dh-badge"
          :class="auto?.enabled ? 'dh-badge-warn' : 'dh-badge-plain'"
          :title="auto?.enabled ? '已开启 · 去「设置 → 更新与检测」调整' : '已关闭 · 去「设置 → 更新与检测」打开'"
        >
          {{ auto?.enabled ? '已开启' : '已关闭' }}
        </RouterLink>
        <div class="ml-auto flex flex-wrap items-center gap-x-3 gap-y-1">
          <span
            class="flex items-center gap-1.5 text-[11.5px] text-text-5"
            :title="auto?.lastCheckAt ? `上次巡检 ${formatDateTime(auto.lastCheckAt)}` : '还没有巡检记录'"
          >
            <Clock class="h-3 w-3" />
            {{ auto?.nextCheckAt ? `下次巡检 ${formatDateTime(auto.nextCheckAt)}` : '未开启周期巡检' }}
            <template v-if="auto?.checkIntervalHours">· 每 {{ auto.checkIntervalHours }} 小时一次</template>
            <template v-if="auto?.lastCheckAt">· 上次 {{ relativeTime(auto.lastCheckAt) }}</template>
          </span>
          <button class="dh-btn dh-btn-sm" :disabled="autoBusy" @click="runAutoCycle(true)">
            <RefreshCw class="h-3 w-3" :class="autoBusy ? 'dh-spin' : ''" />立即巡检一轮
          </button>
          <button class="dh-btn dh-btn-sm dh-btn-primary" :disabled="autoBusy || !auto?.enabled" @click="runAutoCycle(false)">
            <Play class="h-3 w-3" />立即执行自动更新
          </button>
        </div>
      </div>

      <div class="flex flex-col gap-3 p-3.5">
        <!-- 本轮会发生什么 -->
        <div class="rounded-[10px] border border-line-1 bg-ink-800">
          <div class="flex flex-wrap items-center gap-2 border-b border-line-1 px-3 py-2">
            <span class="text-[12px] font-medium text-text-3">本轮会发生什么</span>
            <!-- 自动更新关闭时绝不能说「将被更新」——那一轮什么都不会发生，只是标记出有更新的容器 -->
            <span class="dh-badge dh-badge-warn">
              {{ auto?.enabled ? `${willUpdate.length} 个容器将被更新` : `${willUpdate.length} 个容器有可用更新` }}
            </span>
            <span v-if="auto?.enabled === false" class="dh-badge dh-badge-plain">自动更新已关闭 · 手动更新请进容器页</span>
            <span class="dh-badge dh-badge-plain">{{ skipped.length }} 个容器被保护/排除</span>
            <RouterLink to="/settings" class="dh-tap-txt ml-auto text-[11.5px] text-text-5 hover:text-accent">
              去设置里调整检测频率与排除列表 →
            </RouterLink>
          </div>
          <div v-if="!auto?.candidates?.length" class="px-3 py-3 text-[12px] text-text-5">
            还没有巡检结果，点右上角「立即巡检一轮」先跑一次。
          </div>
          <div v-else class="overflow-x-auto">
            <table class="dh-table">
              <thead>
                <tr>
                  <th class="w-[200px]">容器</th>
                  <th class="w-[220px]">镜像</th>
                  <th class="w-[110px]">状态</th>
                  <th>去向</th>
                </tr>
              </thead>
              <tbody>
                <tr v-for="c in auto.candidates" :key="c.name">
                  <td>
                    <!-- 容器名直接链到容器详情：手动更新就从那里发起（本页不再提供批量更新操作） -->
                    <RouterLink
                      :to="`/containers/${encodeURIComponent(c.name)}`"
                      class="dh-tap-txt font-mono text-[11.5px] text-text-3 hover:text-accent"
                    >
                      {{ c.name }}
                    </RouterLink>
                  </td>
                  <td class="max-w-[220px] truncate font-mono text-[11px] text-text-5">{{ shortImage(c.image) }}</td>
                  <td>
                    <span class="dh-badge" :class="c.running ? 'dh-badge-run' : 'dh-badge-stop'">
                      {{ c.running ? '运行中' : '已停止' }}
                    </span>
                  </td>
                  <td class="text-[11.5px]">
                    <!--
                      「去向」吸收了原「其余容器」表的全部信息：未命中将更新/保护/排除的行，
                      不再笼统写「无需更新」，而是带上巡检的具体判定（已是最新/无法判定/本地镜像）与原因。
                      自动更新关闭时不说「将更新」——那一轮什么都不会发生，只标记「有可用更新」。
                    -->
                    <span v-if="c.willUpdate" class="dh-badge dh-badge-warn">
                      {{ auto?.enabled ? '将更新' : '有可用更新' }}
                    </span>
                    <span v-else-if="c.protected" class="dh-badge dh-badge-accent">受保护</span>
                    <span v-else-if="c.excluded" class="dh-badge dh-badge-plain">已排除</span>
                    <span
                      v-else-if="resultByName.get(c.name)"
                      class="dh-badge"
                      :class="`dh-badge-${toneOf(resultByName.get(c.name)!.status)}`"
                    >
                      {{ checkLabel(resultByName.get(c.name)!.status).text }}
                    </span>
                    <span v-else class="dh-badge dh-badge-plain">未检测到更新</span>
                    <span class="ml-2 text-text-5">{{ resultByName.get(c.name)?.reason || c.reason }}</span>
                  </td>
                </tr>
              </tbody>
            </table>
          </div>
        </div>

        <!--
          「上一轮结果」条与「自动更新默认关闭」长提示已删：
          上一轮的更新/失败数在提交后的 toast、顶部汇总卡与「设置 → 运行记录」里都有；
          开关状态与去设置的入口已在卡头的状态徽标上，不再重复一段说明文字。
        -->
      </div>
    </div>

    <!-- 更新策略 -->
    <!--
      这里原本有一张「更新策略」卡（同一镜像只拉一次 / 更新前备份 / 更新后清理 + 保存按钮），
      每一行都是「设置 → 更新与检测 → 执行策略」里的同一批开关，属于同一份配置的第二个入口，
      已整体删除；本页只用上面那句指路。
    -->

    <!--
      「其余容器」卡已删：它与上面候选表的容器清单完全重复（candidates 覆盖全部容器），
      判定信息（已是最新/无法判定/本地镜像 + 原因）已并入候选表的「去向」列。
    -->

    <!--
      「确认执行更新」弹窗已随手动批量更新入口一起删除
      （那一整套勾选 / 全选 / 更新选中的 / 强制重建都在说明同一件事：这里原本是手动更新的入口）。
      手动更新走容器详情页；本页只保留自动更新与「立即执行自动更新」一个会动容器的按钮。
    -->

    <!-- 原理说明 -->
    <Modal :open="showInfo" title="为什么 Dockhelm 不会误停容器" width="620px" @close="showInfo = false">
      <div class="flex flex-col gap-3 text-[12.5px] leading-relaxed text-text-3">
        <div class="flex items-start gap-2.5">
          <AlertTriangle class="mt-[2px] h-4 w-4 flex-none text-warn-text" />
          <div>
            <b class="text-text-1">dockerCopilot 的两个缺陷。</b>
            一是它的检测路径自拼 registry 请求、用一份硬编码的加速站列表，从不读取守护进程
            <code class="text-text-2">daemon.json</code> 里真正生效的 <code class="text-text-2">registry-mirrors</code>，
            于是「检测说有新版本、拉取说已是最新」会永久互相矛盾。
            二是它的更新流程无条件执行 <code class="text-text-2">pull → stop → rename → create → start</code>，
            <b class="text-text-1">算出了新镜像 ID 却只用来决定要不要删旧镜像，从不用于决定要不要停容器</b>。
            两者叠加，就出现了「一次更新十几个，连没更新的容器也全被停掉」。
          </div>
        </div>
        <div class="flex items-start gap-2.5">
          <CheckCircle2 class="mt-[2px] h-4 w-4 flex-none text-run-text" />
          <div>
            <b class="text-text-1">Dockhelm 的做法。</b>
            检测走守护进程的 <code class="text-text-2">/distribution</code> 接口（与 pull 同一套仓库端点解析），
            检测失败一律标记为「无法判定」而不是「有新版本」。
            更新时先拉取，再比对容器使用的镜像 ID：<b class="text-text-1">一致就直接返回，容器一个字节都不碰</b>；
            不一致才停止旧容器、改名保留（<code class="text-text-2">&lt;名字&gt;__bak_&lt;时间&gt;</code>）、
            用原配置创建同名新容器、启动并做健康检查，失败自动把旧容器改回原名并启动。
          </div>
        </div>
        <div class="flex items-start gap-2.5">
          <Info class="mt-[2px] h-4 w-4 flex-none text-accent" />
          <div>
            <b class="text-text-1">「深度检测」是什么。</b>
            它真的拉一次镜像再比对镜像 ID —— 这是唯一 100% 同源的判定。
            镜像已最新时守护进程只会下载 manifest 与 config（几 KB），不会下载层文件。
            如果你怀疑某个容器被误报，用它确认即可。
          </div>
        </div>
      </div>
      <template #footer>
        <button class="dh-btn dh-btn-primary" @click="showInfo = false">明白了</button>
      </template>
    </Modal>
  </div>
</template>
