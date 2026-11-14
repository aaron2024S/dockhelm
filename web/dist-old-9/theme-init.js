/**
 * 主题预设：必须在样式表加载之前「同步」执行。
 *
 * 这段代码原来内联在 index.html 里。挪出来是因为全站上了 CSP
 * （script-src 'self'），内联脚本会被浏览器直接拒绝 —— 而给 script-src
 * 放开 'unsafe-inline' 等于把 CSP 这层防线整个废掉。
 *
 * 注意：它是普通 <script>（不是 module），在 <head> 里同步执行、阻塞渲染，
 * 所以不会出现「先按深色画一帧再变白」的闪烁。
 *
 * ⚠ 逻辑与 src/stores/theme.ts 的 applyTheme() 必须保持一致，改一处要改两处。
 */
;(function () {
  try {
    var mode = localStorage.getItem('dockhelm.theme') || 'dark'
    var dark =
      mode === 'dark' ||
      (mode === 'system' &&
        !!window.matchMedia &&
        window.matchMedia('(prefers-color-scheme: dark)').matches)
    var root = document.documentElement
    root.classList.toggle('dark', dark)
    root.classList.toggle('light', !dark)
    root.style.colorScheme = dark ? 'dark' : 'light'
    var meta = document.querySelector('meta[name="theme-color"]')
    if (meta) meta.setAttribute('content', dark ? '#05080b' : '#eceff5')
  } catch (e) {
    /* 读不到偏好（隐私模式）时保持 HTML 上的默认深色 */
  }
})()
