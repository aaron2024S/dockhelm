import { defineStore } from 'pinia'

export type ToastKind = 'success' | 'error' | 'info'

export interface Toast {
  id: number
  kind: ToastKind
  message: string
  detail?: string
}

let seq = 0

/** 轻量吐司提示。 */
export const useToastStore = defineStore('toast', {
  state: () => ({
    items: [] as Toast[],
  }),
  actions: {
    push(kind: ToastKind, message: string, detail?: string) {
      const id = ++seq
      this.items.push({ id, kind, message, detail })
      const ttl = kind === 'error' ? 6500 : 3800
      window.setTimeout(() => this.dismiss(id), ttl)
      return id
    },
    success(message: string, detail?: string) {
      return this.push('success', message, detail)
    },
    error(message: string, detail?: string) {
      return this.push('error', message, detail)
    },
    info(message: string, detail?: string) {
      return this.push('info', message, detail)
    },
    dismiss(id: number) {
      this.items = this.items.filter((t) => t.id !== id)
    },
  },
})
