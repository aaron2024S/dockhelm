import { defineStore } from 'pinia'
import { api, setUnauthorizedHandler } from '@/api/client'
import type { SessionInfo, Settings } from '@/api/types'

/** 全局应用状态：登录态、设置、在线状态。 */
export const useAppStore = defineStore('app', {
  state: () => ({
    ready: false,
    initialized: true,
    loggedIn: false,
    sessionCount: 0,
    failures: 0,
    maxFailures: 5,
    minPasswordLength: 6,
    settings: null as Settings | null,
    dockerOnline: true,
    sidebarCollapsed: false,
    /** bootstrap 的进行中 promise（非响应式用途，仅用于去重）。 */
    _bootstrapping: null as Promise<void> | null,
  }),
  getters: {
    authReady: (s) => s.ready,
  },
  actions: {
    async bootstrap() {
      // 冷启动时 App 与登录页都会调 bootstrap（登录页是子组件，先触发），
      // 不去重就会并发发两条 /api/session、跑两遍 loadSettings。
      // 用同一个 promise 兜住，谁先来谁真正干活。
      if (this._bootstrapping) return this._bootstrapping
      this._bootstrapping = this._bootstrap()
      try {
        await this._bootstrapping
      } finally {
        this._bootstrapping = null
      }
    },
    async _bootstrap() {
      setUnauthorizedHandler(() => {
        this.loggedIn = false
        this.ready = true
        if (location.pathname !== '/login') {
          location.replace('/login')
        }
      })
      try {
        const s = await api.get<SessionInfo>('/api/session')
        this.initialized = s.initialized
        this.loggedIn = s.loggedIn
        this.sessionCount = s.sessionCount
        this.failures = s.failures
        this.maxFailures = s.maxFailures
        this.minPasswordLength = s.minPassword || 6
      } catch {
        this.loggedIn = false
      } finally {
        this.ready = true
      }
      if (this.loggedIn) {
        void this.loadSettings()
      }
    },
    async login(password: string, keep: boolean) {
      await api.post('/api/login', { password, keep })
      this.loggedIn = true
      await this.loadSettings()
    },
    async setup(password: string, confirm: string) {
      const res = await api.post<{ autoLogin: boolean }>('/api/setup', { password, confirm })
      this.initialized = true
      if (res.autoLogin) this.loggedIn = true
      return res
    },
    async logout() {
      try {
        await api.post('/api/logout')
      } finally {
        this.loggedIn = false
        this.settings = null
      }
    },
    async loadSettings() {
      try {
        this.settings = await api.get<Settings>('/api/settings')
      } catch {
        /* 忽略：设置页会再拉一次 */
      }
    },
    async saveSettings(patch: Partial<Settings>) {
      // 局部更新：只提交 patch 里出现过的键，不会把别处的改动覆盖掉
      this.settings = await api.patch<Settings>('/api/settings', patch)
      return this.settings
    },
  },
})
