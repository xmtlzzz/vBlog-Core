// 解析后台「项目二维码导航」配置：每行「名称|URL」或只填 URL（名称自动取域名）
export function parseQrLinks(raw) {
  const text = (raw || '').trim()
  if (!text) return []
  const list = []
  for (const line of text.split('\n')) {
    const t = line.trim()
    if (!t) continue
    let name = ''
    let url = ''
    if (t.includes('|')) {
      const parts = t.split('|')
      name = (parts[0] || '').trim()
      url = (parts[1] || '').trim() || name
    } else {
      url = t
    }
    if (!/^https?:\/\//i.test(url)) continue
    if (!name) {
      try {
        name = new URL(url).host.replace(/^www\./, '')
      } catch {
        continue
      }
    }
    list.push({ name, url })
  }
  return list
}
