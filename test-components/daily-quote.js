// DailyQuote - 每日格言卡片
// 后台 → 组件定制 → 上传组件 → 导入此文件

function DailyQuote(container) {
  var style = document.createElement('style')
  style.textContent =
    '.daily-quote{' +
    'padding:24px;background:var(--surface,#fff);border:1px solid var(--border,#e5e5e5);' +
    'border-radius:var(--radius,8px);font-family:var(--font-sans,system-ui,sans-serif);' +
    'position:relative;overflow:hidden' +
    '}' +
    '.daily-quote::before{' +
    'content:"“";position:absolute;top:-8px;left:12px;font-size:80px;' +
    'color:var(--accent,#2563eb);opacity:.12;font-family:Georgia,serif;line-height:1' +
    '}' +
    '.daily-quote .text{' +
    'font-size:16px;line-height:1.7;color:var(--fg,#171717);font-style:italic;' +
    'margin-bottom:12px;position:relative;z-index:1' +
    '}' +
    '.daily-quote .author{' +
    'font-size:13px;color:var(--muted,#737373);text-align:right' +
    '}' +
    '.daily-quote .refresh{' +
    'position:absolute;top:12px;right:12px;background:none;border:none;' +
    'color:var(--muted,#737373);cursor:pointer;font-size:16px;padding:4px;' +
    'border-radius:4px;transition:all .15s;line-height:1' +
    '}' +
    '.daily-quote .refresh:hover{color:var(--accent,#2563eb);background:var(--accent-soft,#eff6ff)}'
  document.head.appendChild(style)

  var quotes = [
    { text: '代码是写给人看的，顺便让机器执行。', author: 'Harold Abelson' },
    { text: '过早优化是万恶之源。', author: 'Donald Knuth' },
    { text: '简单是可靠的前提。', author: 'Edsger W. Dijkstra' },
    { text: '任何足够先进的技术，都与魔法无异。', author: 'Arthur C. Clarke' },
    { text: '程序必须是写给人阅读的，只是顺便让机器执行。', author: 'Abelson & Sussman' },
    { text: '没有银弹。', author: 'Fred Brooks' },
    { text: '好的工具让复杂的事情变简单。', author: '未知' },
    { text: 'Talk is cheap. Show me the code.', author: 'Linus Torvalds' },
    { text: '先让它能用，再让它正确，最后让它快。', author: 'Kent Beck' },
    { text: '软件设计就像做菜，少即是多。', author: '未知' }
  ]

  var el = document.createElement('div')
  el.className = 'daily-quote'
  container.appendChild(el)

  function showRandom() {
    var q = quotes[Math.floor(Math.random() * quotes.length)]
    el.innerHTML =
      '<button class="refresh" title="换一条">&#x21bb;</button>' +
      '<div class="text">' + q.text + '</div>' +
      '<div class="author">— ' + q.author + '</div>'
    el.querySelector('.refresh').onclick = showRandom
  }

  showRandom()

  return function cleanup() {
    style.remove()
    el.remove()
  }
}
