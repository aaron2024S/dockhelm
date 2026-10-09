<script setup lang="ts">
import { computed, onMounted, ref } from 'vue'
import {
  Boxes,
  CalendarClock,
  Cpu,
  ExternalLink,
  Github,
  HardDrive,
  Heart,
  Info,
  Layers,
  Rocket,
  Server,
  ShieldCheck,
  Zap,
} from 'lucide-vue-next'
import { api } from '@/api/client'
import type { AboutResponse } from '@/api/types'
import { formatBytes, formatDuration } from '@/utils/format'
import { useToastStore } from '@/stores/toast'

const toast = useToastStore()
const data = ref<AboutResponse | null>(null)
const loading = ref(true)

const about = computed(() => data.value?.about)
const runtime = computed(() => (data.value?.runtime ?? {}) as Record<string, any>)

const uptimeSeconds = computed(() => Number(runtime.value.uptimeSeconds ?? 0))

async function load() {
  loading.value = true
  try {
    data.value = await api.get<AboutResponse>('/api/about')
  } catch (e) {
    toast.error('读取版本信息失败', e instanceof Error ? e.message : String(e))
  } finally {
    loading.value = false
  }
}

async function copyCommit() {
  const c = about.value?.commit
  if (!c) return
  try {
    await navigator.clipboard.writeText(c)
    toast.success('已复制完整提交号')
  } catch {
    toast.error('复制失败')
  }
}

/** 功能清单：关于页要能在没连 Docker 的时候也把「这软件能干什么」讲清楚。 */
const features = [
  { icon: Boxes, title: '容器管理', desc: '启停、重启、重命名、日志、实时资源占用、环境变量与挂载一览' },
  { icon: Layers, title: '镜像管理', desc: '列表、清理悬空镜像、一键删除；显示体积与被引用情况' },
  { icon: Zap, title: '更新中心', desc: '先拉取再比对镜像 ID，镜像没变就完全不动容器；失败自动回滚' },
  { icon: CalendarClock, title: '计划任务', desc: '标准 cron 表达式，定时启停 / 重启 / 更新 / 备份容器' },
  { icon: Rocket, title: '加速源', desc: '显示守护进程真实生效的镜像站、批量测速、生成 daemon.json 片段' },
  { icon: HardDrive, title: '备份与恢复', desc: '容器配置快照、差异预览、一键还原；compose 文件可读可看' },
  { icon: ShieldCheck, title: '登录鉴权', desc: '单用户密码 + bcrypt 哈希 + HttpOnly 会话 + 连续失败锁定' },
  { icon: Heart, title: '事件通知', desc: '10 类渠道预设 + 自定义 Webhook，事件订阅、静默时段、防轰炸' },
]

const stack = [
  { label: '后端', value: 'Go + 标准库 net/http（无 Web 框架）' },
  { label: '容器对接', value: '直接使用 Docker Engine HTTP API' },
  { label: '调度', value: 'robfig/cron v3' },
  { label: '密码哈希', value: 'golang.org/x/crypto/bcrypt' },
  { label: '持久化', value: '原子写入的 JSON 文档（无 CGO、无外部数据库）' },
  { label: '前端', value: 'Vue 3 + TypeScript + Vite + Tailwind CSS + Pinia + Vue Router' },
  { label: '图标', value: 'lucide' },
  { label: '分发', value: '单二进制 + //go:embed 内嵌前端，单容器部署' },
]

onMounted(() => void load())
</script>

<template>
  <div class="flex flex-col gap-3.5 p-[18px]">
    <!-- 头部 -->
    <div class="dh-card overflow-hidden">
      <div
        class="flex flex-wrap items-center gap-4 p-5"
        style="background: radial-gradient(620px 200px at 12% 0%, rgba(45, 212, 191, 0.1), transparent 70%)"
      >
        <div class="grid h-[58px] w-[58px] flex-none place-items-center rounded-[17px] bg-accent text-accent-ink">
          <svg viewBox="0 0 32 32" class="h-7 w-7" fill="none" stroke="currentColor" stroke-width="2.4">
            <path d="M16 7l7 4v10l-7 4-7-4V11z" stroke-linejoin="round" />
            <path d="M16 15v10M9 11l7 4 7-4" stroke-linejoin="round" />
          </svg>
        </div>
        <div class="min-w-0 flex-1">
          <div class="flex flex-wrap items-center gap-2.5">
            <span class="text-[20px] font-semibold tracking-[-0.2px]">{{ about?.name ?? 'Dockhelm' }}</span>
            <span class="dh-badge dh-badge-accent">v{{ about?.version ?? '—' }}</span>
            <span class="dh-badge dh-badge-plain">{{ about?.license }}</span>
            <button
              v-if="about?.commitShort"
              type="button"
              class="dh-badge dh-badge-plain cursor-pointer"
              title="点击复制完整提交号"
              @click="copyCommit"
            >
              <Github class="h-3 w-3" />{{ about.commitShort }}
            </button>
          </div>
          <div class="mt-1.5 text-[13px] text-text-2">{{ about?.tagline }}</div>
          <div class="mt-1 text-[12px] leading-relaxed text-text-4">{{ about?.description }}</div>
        </div>
        <div class="flex flex-none flex-col gap-2">
          <a
            :href="about?.repoURL"
            target="_blank"
            rel="noreferrer noopener"
            class="dh-btn dh-btn-primary"
          >
            <Github class="h-3.5 w-3.5" />项目主页
            <ExternalLink class="h-3 w-3" />
          </a>
          <a
            :href="about?.authorURL"
            target="_blank"
            rel="noreferrer noopener"
            class="dh-btn"
          >
            <Heart class="h-3.5 w-3.5" />作者主页
            <ExternalLink class="h-3 w-3" />
          </a>
        </div>
      </div>

      <div class="grid grid-cols-2 gap-x-6 gap-y-3 border-t border-line-1 p-4 lg:grid-cols-4">
        <div>
          <div class="text-[11.5px] text-text-5">版本</div>
          <div class="mt-0.5 text-[13px] font-medium">{{ about?.version ?? '—' }}</div>
        </div>
        <div>
          <div class="text-[11.5px] text-text-5">作者</div>
          <div class="mt-0.5 text-[13px] font-medium">{{ about?.author ?? '—' }}</div>
        </div>
        <div>
          <div class="text-[11.5px] text-text-5">构建时间</div>
          <div class="mt-0.5 text-[13px] font-medium">{{ about?.buildTime || '本地构建未注入' }}</div>
        </div>
        <div>
          <div class="text-[11.5px] text-text-5">许可协议</div>
          <div class="mt-0.5 text-[13px] font-medium">{{ about?.license }}</div>
        </div>
        <div>
          <div class="text-[11.5px] text-text-5">提交</div>
          <div class="mt-0.5 break-all font-mono text-[12px]">{{ about?.commit || '未知' }}</div>
        </div>
        <div>
          <div class="text-[11.5px] text-text-5">Go 版本</div>
          <div class="mt-0.5 font-mono text-[12px]">{{ about?.goVersion || '—' }}</div>
        </div>
        <div>
          <div class="text-[11.5px] text-text-5">运行平台</div>
          <div class="mt-0.5 font-mono text-[12px]">{{ about?.platform || '—' }}</div>
        </div>
        <div>
          <div class="text-[11.5px] text-text-5">已运行</div>
          <div class="mt-0.5 text-[13px] font-medium">
            {{ loading ? '载入中…' : formatDuration(uptimeSeconds) }}
          </div>
        </div>
      </div>
    </div>

    <!-- 功能 -->
    <div class="dh-card">
      <div class="dh-card-head">
        <Info class="h-3.5 w-3.5 text-text-4" />
        <span>它能做什么</span>
      </div>
      <div class="grid grid-cols-1 gap-2.5 p-3.5 md:grid-cols-2 xl:grid-cols-4">
        <div
          v-for="f in features"
          :key="f.title"
          class="rounded-[10px] border border-line-1 bg-ink-800 p-3"
        >
          <div class="flex items-center gap-2">
            <component :is="f.icon" class="h-3.5 w-3.5 text-accent" />
            <span class="text-[12.5px] font-medium">{{ f.title }}</span>
          </div>
          <div class="mt-1.5 text-[11.5px] leading-relaxed text-text-4">{{ f.desc }}</div>
        </div>
      </div>
    </div>

    <div class="grid grid-cols-1 gap-3.5 xl:grid-cols-2">
      <!-- 运行环境 -->
      <div class="dh-card">
        <div class="dh-card-head">
          <Server class="h-3.5 w-3.5 text-text-4" />
          <span>运行环境</span>
        </div>
        <div class="grid grid-cols-2 gap-x-5 gap-y-2.5 p-3.5 text-[12px]">
          <div class="text-text-5">数据目录</div>
          <div class="truncate font-mono text-[11.5px]" :title="String(runtime.dataDir ?? '')">
            {{ runtime.dataDir ?? '—' }}
          </div>
          <div class="text-text-5">Docker 地址</div>
          <div class="truncate font-mono text-[11.5px]">{{ runtime.dockerHost ?? '—' }}</div>
          <div class="text-text-5">Docker 版本</div>
          <div class="text-text-2">
            {{ runtime.dockerVersion ?? '未连接' }}
            <span v-if="runtime.dockerApiVersion" class="text-text-5">（API v{{ runtime.dockerApiVersion }}）</span>
          </div>
          <div class="text-text-5">Docker 系统</div>
          <div class="truncate text-text-2">
            {{ [runtime.dockerOs, runtime.dockerArch, runtime.dockerKernel].filter(Boolean).join(' · ') || '—' }}
          </div>
          <div class="text-text-5">Docker 数据根</div>
          <div class="truncate font-mono text-[11.5px]" :title="String(runtime.dockerRoot ?? '')">
            {{ runtime.dockerRoot ?? '—' }}
          </div>
          <div class="text-text-5">容器 / 镜像总数</div>
          <div class="text-text-2">
            {{ runtime.containersTotal ?? '—' }} 个容器 · {{ runtime.imagesTotal ?? '—' }} 个镜像
          </div>
        </div>
      </div>

      <!-- 技术栈 -->
      <div class="dh-card">
        <div class="dh-card-head">
          <Cpu class="h-3.5 w-3.5 text-text-4" />
          <span>技术栈</span>
        </div>
        <div class="flex flex-col p-3.5">
          <div
            v-for="s in stack"
            :key="s.label"
            class="flex gap-4 border-b border-[#171f2a] py-2 text-[12px] last:border-b-0"
          >
            <div class="w-[76px] flex-none text-text-5">{{ s.label }}</div>
            <div class="min-w-0 flex-1 text-text-2">{{ s.value }}</div>
          </div>
        </div>
      </div>
    </div>

    <!-- 设计取向 -->
    <div class="dh-card">
      <div class="dh-card-head">
        <ShieldCheck class="h-3.5 w-3.5 text-text-4" />
        <span>设计取向</span>
      </div>
      <div class="grid grid-cols-1 gap-3 p-3.5 lg:grid-cols-2">
        <div class="rounded-[10px] border border-line-1 bg-ink-800 p-3">
          <div class="text-[12.5px] font-medium text-text-1">检测与拉取必须同源</div>
          <div class="mt-1.5 text-[11.5px] leading-relaxed text-text-4">
            检测走 Docker 守护进程自己解析的仓库端点，拉取走守护进程的 ImagePull。
            检测失败一律标记为「无法判定」，绝不退化成「有新版本」——
            这正是同类工具出现「检测说有更新、拉取说已是最新」永久误报的根因。
          </div>
        </div>
        <div class="rounded-[10px] border border-line-1 bg-ink-800 p-3">
          <div class="text-[12.5px] font-medium text-text-1">镜像没变就不碰容器</div>
          <div class="mt-1.5 text-[11.5px] leading-relaxed text-text-4">
            更新流程是先拉取、再比对容器使用的镜像 ID。一致就直接返回，
            <b class="text-text-3">容器不会被停止、不会被重建</b>。只有镜像真的变化时才会走
            「停旧 → 改名保留 → 建新 → 健康检查 → 失败自动回滚」。
          </div>
        </div>
        <div class="rounded-[10px] border border-line-1 bg-ink-800 p-3">
          <div class="text-[12.5px] font-medium text-text-1">看得见才叫能备份</div>
          <div class="mt-1.5 text-[11.5px] leading-relaxed text-text-4">
            绑定挂载的数据在宿主机目录上，Dockhelm 默认看不见它。界面上会明确标注
            「知道路径但看不见」，绝不假装备份成功 —— 只记路径，并告诉你去挂载哪个目录。
          </div>
        </div>
        <div class="rounded-[10px] border border-line-1 bg-ink-800 p-3">
          <div class="text-[12.5px] font-medium text-text-1">少依赖，容易自托管</div>
          <div class="mt-1.5 text-[11.5px] leading-relaxed text-text-4">
            整个后端只有两个第三方依赖（cron 与 bcrypt），HTTP 用标准库，持久化用原子写入的 JSON。
            因此可以 CGO_ENABLED=0 交叉编译 amd64 / arm64，也不需要外部数据库。
          </div>
        </div>
      </div>
      <div class="border-t border-line-1 px-4 py-3 text-[11.5px] leading-relaxed text-text-5">
        Dockhelm 会挂载 <code class="text-text-3">/var/run/docker.sock</code>，这等价于宿主机的 root 权限。
        请务必设置登录密码，并且只在可信的内网中暴露端口。配置存储占用：
        {{ formatBytes(1024) }} 量级的小文件，随时可以直接拷贝。
      </div>
    </div>
  </div>
</template>
