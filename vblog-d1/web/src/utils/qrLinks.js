// 可用配色预设（与 EveryQrBadge 内 QR_PALETTES 对齐）
export const QR_PALETTE_NAMES = ['sakura', 'forest', 'ocean', 'sunset']

// 解析后台「项目二维码导航」配置：
// 每行「名称|URL」或「名称|URL|配色」或只填 URL（名称自动取域名，配色按行轮换）
export function parseQrLinks(raw) {
  const text = (raw || '').trim()
  if (!text) return []
  const list = []
  for (const line of text.split('\n')) {
    const t = line.trim()
    if (!t) continue
    let name = ''
    let url = ''
    let palette = ''
    if (t.includes('|')) {
      const parts = t.split('|')
      name = (parts[0] || '').trim()
      url = (parts[1] || '').trim() || name
      const p = (parts[2] || '').trim().toLowerCase()
      if (QR_PALETTE_NAMES.includes(p)) palette = p
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
    list.push({ name, url, palette })
  }
  return list
}
