import { defineStore } from 'pinia'

/**
 * 主题三态。数组顺序 = 顶栏按钮的循环顺序（深色 → 浅色 → 跟随系统 → 深色）。
 *
 * 两套配色只差 CSS 变量，见 style.css 顶部的 `@theme static` 与 `html.light` 覆盖；
 * 组件里写的都是 bg-ink-700 / text-text-3 这类引用变量的类，因此不需要条件类。
 */
export type ThemeMode = 'dark' | 'light' | 'system'

export const THEME_ORDER: readonly ThemeMode[] = ['dark', 'light', 'system'] as const

/** localStorage 键名。index.html 里的防闪烁脚本读的是同一个键。 */
const KEY = 'dockhelm.theme'

export const THEME_LABEL: Record<ThemeMode, string> = {
  dark: '深色',
  light: '浅色',
  system: '跟随系统',
}

/** 地址栏 / 浏览器 UI 配色，与 style.css 的 --color-ink-950 保持一致。 */
const BAR_COLOR: Record<'dark' | 'light', string> = { dark: '#05080b', light: '#eceff5' }

function prefersDark(): boolean {
  return window.matchMedia?.('(prefers-color-scheme: dark)').matches ?? true
}

function isMode(v: string | null): v is ThemeMode {
  return v === 'dark' || v === 'light' || v === 'system'
}

/** 读持久化偏好。默认 dark：不因为新增了主题功能就把老用户突然变白。 */
export function readStoredMode(): ThemeMode {
  try {
    const v = localStorage.getItem(KEY)
    return isMode(v) ? v : 'dark'
  } catch {
    return 'dark'
  }
}

/** 三态 → 实际生效的两态。 */
export function resolveMode(mode: ThemeMode): 'dark' | 'light' {
  if (mode === 'light') return 'light'
  if (mode === 'dark') return 'dark'
  return prefersDark() ? 'dark' : 'light'
}

/** 循环切换的下一档。 */
export function nextMode(mode: ThemeMode): ThemeMode {
  const i = THEME_ORDER.indexOf(mode)
  return THEME_ORDER[(i + 1) % THEME_ORDER.length] ?? 'dark'
}

/**
 * 把模式落到 <html> 上。三件事：
 * ① class dark/light —— 决定走哪套 CSS 变量；
 * ② color-scheme —— 让滚动条、原生下拉、日期选择器这些浏览器绘制的东西跟着换；
 * ③ theme-color —— 移动端地址栏配色。
 *
 * ⚠ 这段逻辑在 index.html 里有一份同步版内联脚本（首屏防闪烁），改这里要同步改那边。
 */
export function applyTheme(mode: ThemeMode) {
  const eff = resolveMode(mode)
  const root = document.documentElement
  root.classList.toggle('dark', eff === 'dark')
  root.classList.toggle('light', eff === 'light')
  root.style.colorScheme = eff
  document.querySelector('meta[name="theme-color"]')?.setAttribute('content', BAR_COLOR[eff])
  return eff
}

export const useThemeStore = defineStore('theme', {
  state: () => ({
    /** 用户选的模式（三态） */
    mode: 'dark' as ThemeMode,
    /** 实际生效的两态，供 UI 展示「跟随系统」时当前到底是哪个 */
    effective: 'dark' as 'dark' | 'light',
  }),
  getters: {
    label: (s) => THEME_LABEL[s.mode],
    /** 点一下会切到什么，用于按钮的 title 提示。 */
    nextLabel: (s) => THEME_LABEL[nextMode(s.mode)],
  },
  actions: {
    /**
     * 首次挂载时调用：把 <html> 上的 class 与 store 对齐，
     * 并在「跟随系统」模式下监听系统主题变化（macOS/Windows 的自动日夜切换）。
     */
    init() {
      this.mode = readStoredMode()
      this.effective = applyTheme(this.mode)
      const mq = window.matchMedia?.('(prefers-color-scheme: dark)')
      // addEventListener 在老 Safari 上不存在，带上 addListener 兜底
      const onChange = () => {
        if (this.mode === 'system') this.effective = applyTheme(this.mode)
      }
      if (mq?.addEventListener) mq.addEventListener('change', onChange)
      else mq?.addListener?.(onChange)
    },
    set(mode: ThemeMode) {
      this.mode = mode
      try {
        localStorage.setItem(KEY, mode)
      } catch {
        /* 隐私模式下写不进去：本次会话仍然生效，只是记不住 */
      }
      this.effective = applyTheme(mode)
    },
    /** 顶栏按钮：深色 → 浅色 → 跟随系统 → 深色 循环。 */
    cycle() {
      this.set(nextMode(this.mode))
    },
  },
})
