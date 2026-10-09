<script setup lang="ts">
import { computed, onMounted, ref } from 'vue'
import { RouterLink } from 'vue-router'
import { Network, PlugZap, RefreshCw } from 'lucide-vue-next'
import { api } from '@/api/client'
import type { ContainerView, DockerNetwork } from '@/api/types'
import { useToastStore } from '@/stores/toast'
import EmptyState from '@/components/EmptyState.vue'

const toast = useToastStore()

const containers = ref<ContainerView[]>([])
const networks = ref<DockerNetwork[]>([])
const loading = ref(true)

interface PortRow {
  container: string
  state: string
  /** 宿主机侧（可能是 IP:端口，也可能是「仅容器内」）。 */
  host: string
  target: string
  proto: string
}

/**
 * 端口映射来自 /api/containers 里已经格式化好的字符串：
 * 「8080→80/tcp」= 映射到宿主机任意地址；「192.168.3.62:8080→80/tcp」= 只绑定某个网卡；
 * 「80/tcp」= 没有对外发布，只能在容器网络里访问。
 */
const portRows = computed<PortRow[]>(() => {
  const out: PortRow[] = []
  for (const c of containers.value) {
    for (const p of c.ports) {
      const arrow = p.includes('→')
      const halves = p.split('→')
      const left = halves[0] ?? ''
      const right = halves[1] ?? ''
      const tail = arrow ? right : left
      const seg = tail.split('/')
      out.push({
        container: c.name,
        state: c.state,
        host: arrow ? left : '仅容器内',
        target: seg[0] ?? '',
        proto: seg[1] ?? 'tcp',
      })
    }
  }
  // 发布到宿主机的排前面，其次按端口号排序
  return out.sort((a, b) => {
    const pa = a.host === '仅容器内' ? 1 : 0
    const pb = b.host === '仅容器内' ? 1 : 0
    if (pa !== pb) return pa - pb
    return Number(a.target) - Number(b.target)
  })
})

const published = computed(() => portRows.value.filter((r) => r.host !== '仅容器内').length)

function subnetOf(n: DockerNetwork) {
  return n.IPAM?.Config?.[0]?.Subnet || '—'
}
function gatewayOf(n: DockerNetwork) {
  return n.IPAM?.Config?.[0]?.Gateway || '—'
}
function membersOf(n: DockerNetwork) {
  return Object.values(n.Containers ?? {})
    .map((c) => c.Name)
    .sort()
}

async function load() {
  loading.value = true
  try {
    const [c, n] = await Promise.all([
      api.get<{ containers: ContainerView[] }>('/api/containers'),
      api.get<{ networks: DockerNetwork[] }>('/api/networks'),
    ])
    containers.value = c.containers ?? []
    // 只保留自定义网络与内置 bridge/host/none，按名称排序
    networks.value = (n.networks ?? []).slice().sort((a, b) => a.Name.localeCompare(b.Name))
  } catch (e) {
    toast.error('读取网络信息失败', e instanceof Error ? e.message : String(e))
  } finally {
    loading.value = false
  }
}

onMounted(() => void load())
</script>

<template>
  <div class="flex flex-col gap-3.5 p-[18px]">
    <div class="dh-phead">
      <div class="dh-h1">网络与端口</div>
      <div class="dh-sub">
        {{ networks.length }} 个网络 · {{ portRows.length }} 条端口映射（其中 {{ published }} 条对宿主机发布）
      </div>
      <div class="ml-auto flex flex-wrap gap-2">
        <button class="dh-btn" :disabled="loading" @click="load">
          <RefreshCw class="h-3.5 w-3.5" :class="loading ? 'dh-spin' : ''" />刷新
        </button>
      </div>
    </div>

    <div class="dh-banner dh-banner-info items-start">
      <span class="mt-[3px] h-[13px] w-[13px] flex-none rounded-[4px] bg-current opacity-50" />
      <span>
        这里的网络与端口来自 Docker 守护进程的实时状态，<b>只读</b>。Dockhelm 不会替你改
        <code class="text-text-3">daemon.json</code>，也不会给容器增删端口 —— 那需要重建容器，请走 compose 或
        <code class="text-text-3">docker run</code>。
      </span>
    </div>

    <!-- 端口映射 -->
    <div class="dh-card">
      <div class="dh-card-head">
        <PlugZap class="h-3.5 w-3.5 text-text-4" />
        <span>端口映射</span>
        <span class="ml-auto text-[11.5px] font-normal text-text-5">
          {{ published }} / {{ portRows.length }} 条对宿主机发布
        </span>
      </div>

      <div v-if="loading && !portRows.length" class="grid h-[120px] place-items-center">
        <RefreshCw class="h-5 w-5 dh-spin text-text-5" />
      </div>
      <div v-else-if="!portRows.length" class="dh-card-body text-[12.5px] text-text-4">
        现有容器都没有暴露端口。
      </div>
      <div v-else class="dh-scroll max-h-[520px] overflow-auto">
        <table class="dh-table">
          <thead>
            <tr>
              <th class="w-[220px]">容器</th>
              <th class="w-[240px]">宿主地址 → 容器</th>
              <th class="w-[80px]">协议</th>
              <th class="w-[92px]">状态</th>
            </tr>
          </thead>
          <tbody>
            <tr v-for="(r, i) in portRows" :key="i">
              <td>
                <RouterLink
                  :to="`/containers/${encodeURIComponent(r.container)}`"
                  class="text-[12.5px] font-medium hover:text-accent"
                >
                  {{ r.container }}
                </RouterLink>
              </td>
              <td class="font-mono text-[11.5px]">
                <span :class="r.host === '仅容器内' ? 'text-text-5' : 'text-text-2'">{{ r.host }}</span>
                <span class="mx-1.5 text-text-5">→</span>
                <span class="text-text-2">{{ r.target }}/{{ r.proto }}</span>
              </td>
              <td class="font-mono text-[11.5px] text-text-4">{{ r.proto }}</td>
              <td>
                <span class="dh-badge" :class="r.state === 'running' ? 'dh-badge-run' : 'dh-badge-stop'">
                  <span
                    class="h-[7px] w-[7px] rounded-full"
                    :class="r.state === 'running' ? 'bg-run' : 'bg-stop'"
                  />
                  {{ r.state === 'running' ? '监听中' : '未监听' }}
                </span>
              </td>
            </tr>
          </tbody>
        </table>
      </div>
    </div>

    <!-- 网络 -->
    <div class="dh-card">
      <div class="dh-card-head">
        <Network class="h-3.5 w-3.5 text-text-4" />
        <span>Docker 网络</span>
        <span class="ml-auto text-[11.5px] font-normal text-text-5">{{ networks.length }} 个</span>
      </div>

      <EmptyState
        v-if="!loading && !networks.length"
        :icon="Network"
        title="没有读到网络"
        description="确认 Docker 守护进程是否正常。"
      />
      <div v-else class="dh-scroll max-h-[560px] overflow-auto">
        <table class="dh-table">
          <thead>
            <tr>
              <th class="w-[190px]">名称</th>
              <th class="w-[110px]">驱动</th>
              <th class="w-[90px]">作用域</th>
              <th class="w-[150px]">子网</th>
              <th class="w-[140px]">网关</th>
              <th class="w-[80px]">容器</th>
              <th>成员</th>
            </tr>
          </thead>
          <tbody>
            <tr v-for="n in networks" :key="n.Id">
              <td>
                <div class="text-[12.5px] font-medium">{{ n.Name }}</div>
                <div class="flex flex-wrap gap-1.5 pt-1">
                  <span v-if="n.Internal" class="dh-badge dh-badge-plain">内部</span>
                  <span v-if="n.Attachable" class="dh-badge dh-badge-plain">可附加</span>
                </div>
              </td>
              <td class="font-mono text-[11.5px] text-text-3">{{ n.Driver }}</td>
              <td class="text-[11.5px] text-text-4">{{ n.Scope }}</td>
              <td class="font-mono text-[11.5px] text-text-3">{{ subnetOf(n) }}</td>
              <td class="font-mono text-[11.5px] text-text-3">{{ gatewayOf(n) }}</td>
              <td class="text-[12px] text-text-2">{{ membersOf(n).length }}</td>
              <td class="text-[11.5px] text-text-4">
                <span v-if="!membersOf(n).length" class="text-text-5">—</span>
                <template v-else>
                  <template v-for="(m, i) in membersOf(n).slice(0, 6)" :key="m">
                    <span v-if="i"> · </span>
                    <RouterLink :to="`/containers/${encodeURIComponent(m)}`" class="hover:text-accent">
                      {{ m }}
                    </RouterLink>
                  </template>
                  <span v-if="membersOf(n).length > 6" class="text-text-5">
                    …等 {{ membersOf(n).length }} 个
                  </span>
                </template>
              </td>
            </tr>
          </tbody>
        </table>
      </div>
    </div>
  </div>
</template>
