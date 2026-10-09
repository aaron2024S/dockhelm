/**
 * 后端接口封装。
 *
 * 约定：所有接口都在 /api 下；未登录时后端返回 401，这里统一抛 UnauthorizedError，
 * 由路由守卫或全局处理跳到登录页。
 */

import type { BusEvent } from './types'

export class ApiError extends Error {
  status: number
  code?: string
  payload?: unknown
  constructor(message: string, status: number, code?: string, payload?: unknown) {
    super(message)
    this.name = 'ApiError'
    this.status = status
    this.code = code
    this.payload = payload
  }
}

export class UnauthorizedError extends ApiError {
  constructor(message = '未登录或会话已过期') {
    super(message, 401, 'unauthorized')
    this.name = 'UnauthorizedError'
  }
}

/** 登录失效时的回调（由 store 注册，避免这里直接依赖 router）。 */
let onUnauthorized: (() => void) | null = null
export function setUnauthorizedHandler(fn: () => void) {
  onUnauthorized = fn
}

type Options = {
  method?: string
  body?: unknown
  query?: Record<string, string | number | boolean | undefined | null>
}

function buildUrl(path: string, query?: Options['query']): string {
  let url = path
  if (query) {
    const params = new URLSearchParams()
    for (const [k, v] of Object.entries(query)) {
      if (v === undefined || v === null || v === '') continue
      params.set(k, String(v))
    }
    const qs = params.toString()
    if (qs) url += (url.includes('?') ? '&' : '?') + qs
  }
  return url
}

export async function request<T>(path: string, opts: Options = {}): Promise<T> {
  const init: RequestInit = {
    method: opts.method ?? 'GET',
    credentials: 'same-origin',
    headers: {},
  }
  if (opts.body !== undefined) {
    init.headers = { 'Content-Type': 'application/json' }
    init.body = JSON.stringify(opts.body)
  }
  let res: Response
  try {
    res = await fetch(buildUrl(path, opts.query), init)
  } catch (e) {
    throw new ApiError('网络请求失败：' + (e instanceof Error ? e.message : String(e)), 0)
  }

  if (res.status === 401) {
    onUnauthorized?.()
    throw new UnauthorizedError()
  }

  const text = await res.text()
  let data: unknown = null
  if (text) {
    try {
      data = JSON.parse(text)
    } catch {
      data = text
    }
  }

  if (!res.ok) {
    const obj = (data ?? {}) as Record<string, unknown>
    const msg = typeof obj.error === 'string' ? obj.error : `请求失败（HTTP ${res.status}）`
    throw new ApiError(msg, res.status, typeof obj.code === 'string' ? obj.code : undefined, data)
  }
  return data as T
}

export const api = {
  get: <T>(path: string, query?: Options['query']) => request<T>(path, { query }),
  post: <T>(path: string, body?: unknown, query?: Options['query']) =>
    request<T>(path, { method: 'POST', body, query }),
  put: <T>(path: string, body?: unknown) => request<T>(path, { method: 'PUT', body }),
  del: <T>(path: string, query?: Options['query']) => request<T>(path, { method: 'DELETE', query }),
  /** 原始文本（容器日志用） */
  text: async (path: string, query?: Options['query']): Promise<string> => {
    const res = await fetch(buildUrl(path, query), { credentials: 'same-origin' })
    if (res.status === 401) {
      onUnauthorized?.()
      throw new UnauthorizedError()
    }
    if (!res.ok) throw new ApiError(await res.text(), res.status)
    return res.text()
  },
}

/** 建立 SSE 连接，返回关闭函数。 */
export function openStream(
  path: string,
  onEvent: (event: string, ev: BusEvent) => void,
  onError?: () => void,
): () => void {
  let closed = false
  let es: EventSource | null = null
  let retry = 0
  let timer: number | undefined

  const connect = () => {
    if (closed) return
    es = new EventSource(path, { withCredentials: true })
    es.onmessage = (e) => handle(e)
    // 后端把 topic 写进 event 字段，这里对几个已知 topic 显式监听
    for (const topic of ['update', 'schedule', 'notify', 'container', 'system', 'message']) {
      es.addEventListener(topic, (e) => handle(e as MessageEvent))
    }
    es.onerror = () => {
      es?.close()
      es = null
      if (closed) return
      retry = Math.min(retry + 1, 6)
      onError?.()
      // 指数退避重连，最长 15 秒
      timer = window.setTimeout(connect, Math.min(1000 * 2 ** retry, 15000))
    }
    es.onopen = () => {
      retry = 0
    }
  }

  const handle = (e: MessageEvent) => {
    if (!e.data) return
    try {
      const parsed = JSON.parse(e.data) as BusEvent
      const topic = typeof parsed.topic === 'string' ? parsed.topic : 'message'
      onEvent(topic, parsed)
    } catch {
      /* 忽略 ping 之类的非 JSON 帧 */
    }
  }

  connect()
  return () => {
    closed = true
    if (timer) window.clearTimeout(timer)
    es?.close()
  }
}
