/**
 * 量「由 dockhelm 按顺序启动」卡片头里标题 / 说明 / 右对齐「预计」三段的基线位置。
 * 基线用零尺寸 inline-block 探针取（height:0 ⇒ rect.top 就是该行基线 y）。
 * 全程零网络：route 全接管，读真 web/dist。
 */
import { chromium } from 'playwright-core'
import { readFileSync, existsSync } from 'node:fs'
import { join, extname } from 'node:path'

const DIST = 'D:/workbuddy/dockhelm/web/dist'
const MIME = { '.html': 'text/html; charset=utf-8', '.js': 'text/javascript', '.css': 'text/css', '.svg': 'image/svg+xml', '.png': 'image/png', '.woff2': 'font/woff2' }
const names = ['alpha', 'beta', 'gamma']

function mkStub() {
  const containers = [
    ...names.map((n) => ({ name: n, state: 'exited', restartPolicy: 'on-failure', managed: true, project: '', auto: true, delaySec: 3 })),
    { name: 'web', state: 'running', restartPolicy: 'always', managed: false, project: '', auto: false, delaySec: 0 },
  ]
  const plan = { enabled: true, initialDelaySec: 30, waitTimeoutSec: 60, defaultDelaySec: 3, items: names.map((n) => ({ name: n, auto: true, delaySec: 3 })) }
  const status = { running: false, total: 0, index: 0, success: 0, skipped: 0, failed: 0, phase: 'idle' }
  return (p) => {
    if (p.includes('/api/boot/status')) return { status }
    if (p.includes('/api/boot')) return { plan, containers, self: { name: 'dockhelm', restartPolicy: 'on-failure' }, lastRun: null, status }
    if (p.includes('/api/session')) return { initialized: true, loggedIn: true }
    if (p.includes('/api/settings')) return { exclude: [] }
    if (p.includes('/api/overview')) return { containers: { total: 4, running: 1, stopped: 3, paused: 0, unhealthy: 0 }, images: { total: 1, sizeBytes: 1, reclaimable: 0 }, updates: { available: 0, unknown: 0, checkedAt: '', items: [] }, disk: { free: 1, total: 2 }, notify: { sentToday: 0 }, recent: [], excluded: [], self: 'dockhelm' }
    return {}
  }
}

async function routeAll(page) {
  const stub = mkStub()
  await page.route('**/*', async (route) => {
    const req = route.request()
    const url = new URL(req.url())
    if (url.hostname !== 'dockhelm.local') return route.fulfill({ status: 204, body: '' })
    if (url.pathname.startsWith('/api/')) {
      return route.fulfill({ status: 200, contentType: 'application/json', body: JSON.stringify(stub(url.pathname)) })
    }
    let p = join(DIST, url.pathname)
    if (!existsSync(p) || url.pathname === '/' || !extname(p)) p = join(DIST, 'index.html')
    return route.fulfill({ status: 200, contentType: MIME[extname(p)] || 'application/octet-stream', body: readFileSync(p) })
  })
}

// 在元素末尾插零尺寸探针，返回其 rect.top（= 基线）
const MEASURE = `(el) => {
  const probe = document.createElement('span')
  probe.style.cssText = 'display:inline-block;width:0;height:0;overflow:hidden'
  el.appendChild(probe)
  const y = probe.getBoundingClientRect().top
  probe.remove()
  return y
}`

const b = await chromium.launch({ channel: 'msedge', headless: true })
const page = await b.newPage({ viewport: { width: 1280, height: 900 }, deviceScaleFactor: 2 })
await routeAll(page)
await page.goto('https://dockhelm.local/boot', { waitUntil: 'networkidle' })
await page.waitForTimeout(700)

const info = await page.evaluate((measureSrc) => {
  const measure = eval(measureSrc)
  const head = [...document.querySelectorAll('.dh-card-head')].find((h) => h.textContent.includes('由 dockhelm 按顺序启动'))
  if (!head) return { error: 'head not found' }
  const spans = [...head.children]
  const out = {
    align: getComputedStyle(head).alignItems,
    items: spans.map((s) => {
      const cs = getComputedStyle(s)
      const r = s.getBoundingClientRect()
      return { text: s.textContent.trim().slice(0, 14), fontSize: cs.fontSize, lineHeight: cs.lineHeight, top: +r.top.toFixed(2), h: +r.height.toFixed(2), baseline: +measure(s).toFixed(2) }
    }),
  }
  return out
}, MEASURE)

console.log('align-items =', info.align)
info.items.forEach((i) => console.log(`  ${i.text.padEnd(16)} font=${i.fontSize} lh=${i.lineHeight} top=${i.top} h=${i.h} baseline=${i.baseline}`))

const bases = info.items.map((i) => i.baseline)
const spread = +(Math.max(...bases) - Math.min(...bases)).toFixed(2)
console.log(`基线极差 = ${spread}px → ${spread <= 0.6 ? '✓ 已对齐' : '✗ 仍未对齐'}`)

// 还原成改动前的写法，量出原来的偏差，作为对照
const before = await page.evaluate((measureSrc) => {
  const measure = eval(measureSrc)
  const head = [...document.querySelectorAll('.dh-card-head')].find((h) => h.textContent.includes('由 dockhelm 按顺序启动'))
  const spans = [...head.children]
  head.style.alignItems = 'flex-start'
  spans[0].style.paddingTop = '2px'
  spans[2].style.paddingTop = '1px'
  const ys = spans.map((s) => +measure(s).toFixed(2))
  return { ys, spread: +(Math.max(...ys) - Math.min(...ys)).toFixed(2) }
}, MEASURE)
console.log(`改动前（items-start + pt-[2px]/pt-[1px]）基线: ${before.ys.join(' / ')}  极差=${before.spread}px`)

// 截图：桌面 + 手机宽度
const head = page.locator('.dh-card-head', { hasText: '由 dockhelm 按顺序启动' }).first()
await head.scrollIntoViewIfNeeded()
await page.evaluate(() => {
  // 还原现场（上面改过 inline 样式）
  const h = [...document.querySelectorAll('.dh-card-head')].find((x) => x.textContent.includes('由 dockhelm 按顺序启动'))
  h.style.alignItems = ''
  ;[...h.children].forEach((s) => (s.style.paddingTop = ''))
})
await page.waitForTimeout(200)
await head.screenshot({ path: 'D:/workbuddy/dockhelm/.tmp/head-desktop.png' })

await page.setViewportSize({ width: 390, height: 844 })
await page.waitForTimeout(400)
await head.screenshot({ path: 'D:/workbuddy/dockhelm/.tmp/head-mobile.png' })

await b.close()
