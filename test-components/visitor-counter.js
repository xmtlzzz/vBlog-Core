// VisitorCounter - 在线访客计数器
// 后台 → 组件定制 → 上传组件 → 导入此文件

function VisitorCounter(container) {
  var style = document.createElement('style')
  style.textContent =
    '.visitor-counter{' +
    'display:flex;align-items:center;gap:10px;padding:14px 18px;' +
    'background:var(--surface,#fff);border:1px solid var(--border,#e5e5e5);' +
    'border-radius:var(--radius,8px);font-family:var(--font-sans,system-ui,sans-serif);' +
    'font-size:14px;color:var(--muted,#737373)' +
    '}' +
    '.visitor-counter .pulse-dot{' +
    'width:10px;height:10px;border-radius:50%;background:var(--success,#16a34a);' +
    'flex-shrink:0;animation:vc-pulse 2s infinite' +
    '}' +
    '.visitor-counter .count{' +
    'font-weight:700;font-size:18px;color:var(--fg,#171717);' +
    'font-variant-numeric:tabular-nums;min-width:28px;text-align:center' +
    '}' +
    '.visitor-counter .label{color:var(--muted,#737373);font-size:13px}' +
    '@keyframes vc-pulse{0%,100%{opacity:1;transform:scale(1)}50%{opacity:.4;transform:scale(.85)}}'
  document.head.appendChild(style)

  var el = document.createElement('div')
  el.className = 'visitor-counter'

  var count = Math.floor(Math.random() * 50) + 10
  el.innerHTML =
    '<span class="pulse-dot"></span>' +
    '<span class="count">' + count + '</span>' +
    '<span class="label">人正在浏览</span>'
  container.appendChild(el)

  var countEl = el.querySelector('.count')
  var timer = setInterval(function() {
    count += Math.random() > 0.5 ? 1 : -1
    if (count < 1) count = 1
    countEl.textContent = count
  }, 5000)

  return function cleanup() {
    clearInterval(timer)
    style.remove()
    el.remove()
  }
}
