// SiteStats - 博客站点统计卡片
// 后台 → 组件定制 → 上传组件 → 导入此文件

function SiteStats(container) {
  var style = document.createElement('style')
  style.textContent =
    '.site-stats{' +
    'display:grid;grid-template-columns:repeat(auto-fit,minmax(140px,1fr));gap:12px;' +
    'padding:20px;background:var(--surface,#fff);border:1px solid var(--border,#e5e5e5);' +
    'border-radius:var(--radius,8px);font-family:var(--font-sans,system-ui,sans-serif)' +
    '}' +
    '.site-stats .stat{' +
    'text-align:center;padding:12px 8px;border-radius:6px;transition:background .2s' +
    '}' +
    '.site-stats .stat:hover{background:var(--accent-soft,#eff6ff)}' +
    '.site-stats .stat-val{' +
    'display:block;font-family:var(--font-display,system-ui);font-size:28px;font-weight:700;' +
    'letter-spacing:-.02em;color:var(--fg,#171717);font-variant-numeric:tabular-nums' +
    '}' +
    '.site-stats .stat-label{' +
    'display:block;font-size:12px;color:var(--muted,#737373);margin-top:4px;' +
    'text-transform:uppercase;letter-spacing:.04em' +
    '}' +
    '.site-stats .stat-bar{' +
    'height:3px;border-radius:2px;background:var(--border,#e5e5e5);margin-top:8px;overflow:hidden' +
    '}' +
    '.site-stats .stat-bar-fill{' +
    'height:100%;border-radius:2px;background:var(--accent,#2563eb);transition:width 1.2s ease' +
    '}'
  document.head.appendChild(style)

  var d = window.__BLOG_DATA__ || {}
  var stats = [
    { label: '文章', value: d.posts || 0, max: Math.max((d.posts || 0) * 2, 10) },
    { label: '阅读', value: d.views || 0, max: Math.max((d.views || 0) * 2, 100) },
    { label: '标签', value: d.tags || 0, max: Math.max((d.tags || 0) * 2, 10) },
    { label: '访客', value: d.visitors || 0, max: Math.max((d.visitors || 0) * 2, 100) }
  ]

  var el = document.createElement('div')
  el.className = 'site-stats'

  stats.forEach(function(s) {
    var div = document.createElement('div')
    div.className = 'stat'
    var pct = Math.round((s.value / s.max) * 100)
    div.innerHTML =
      '<span class="stat-val" data-target="' + s.value + '">0</span>' +
      '<span class="stat-label">' + s.label + '</span>' +
      '<div class="stat-bar"><div class="stat-bar-fill" style="width:0%"></div></div>'
    el.appendChild(div)
  })

  container.appendChild(el)

  // 数字滚动动画
  var vals = el.querySelectorAll('.stat-val')
  var fills = el.querySelectorAll('.stat-bar-fill')
  var duration = 1200
  var start = null

  function animate(ts) {
    if (!start) start = ts
    var progress = Math.min((ts - start) / duration, 1)
    var ease = 1 - Math.pow(1 - progress, 3)
    vals.forEach(function(v, i) {
      var target = parseInt(v.getAttribute('data-target'))
      v.textContent = Math.round(target * ease).toLocaleString()
    })
    fills.forEach(function(f, i) {
      var pct = Math.round((stats[i].value / stats[i].max) * 100)
      f.style.width = (pct * ease) + '%'
    })
    if (progress < 1) requestAnimationFrame(animate)
  }

  requestAnimationFrame(animate)

  return function cleanup() {
    style.remove()
    el.remove()
  }
}
