import MarkdownIt from 'markdown-it'

// Parse the editor's Markdown grammar, then replace source spans without reformatting it.
export function scanLocalImages(content) {
  const md = new MarkdownIt({ html: false })
  const destinations = new Map()
  const normalize = md.normalizeLink
  md.normalizeLink = value => {
    const normalized = normalize(value)
    if (!destinations.has(normalized)) destinations.set(normalized, value)
    return normalized
  }
  const validate = md.validateLink
  md.validateLink = href => /^file:/i.test(href) || validate(href)
  const imageRule = md.inline.ruler.__rules__.find(rule => rule.name === 'image').fn
  let insideImage = false
  md.inline.ruler.at('image', (state, silent) => {
    if (insideImage || silent) return imageRule(state, silent)
    const start = state.pos
    const labelEnd = md.helpers.parseLinkLabel(state, start + 1, false)
    insideImage = true
    let matched
    try { matched = imageRule(state, silent) } finally { insideImage = false }
    if (!matched) return false
    const token = state.tokens.at(-1)
    let position = labelEnd + 1
    if (state.src[position] === '(') {
      position++
      while (position < state.posMax && /\s/.test(state.src[position])) position++
      const destination = md.helpers.parseLinkDestination(state.src, position, state.posMax)
      const angle = state.src[position] === '<'
      token.meta = { source: destination.str, start: position + (angle ? 1 : 0), end: destination.pos - (angle ? 1 : 0) }
    } else {
      token.meta = { source: destinations.get(token.attrGet('src')) || token.attrGet('src'), start, end: state.pos, labelEnd, title: token.attrGet('title') || '', reference: true }
    }
    return true
  })

  const offsets = [0]
  for (let i = 0; i < content.length; i++) if (content[i] === '\n') offsets.push(i + 1)
  offsets.push(content.length)
  const images = new Map()
  let range = [0, offsets.length - 1], cursor = 0
  for (const token of md.parse(content, {})) {
    if (token.map) range = token.map
    if (token.type !== 'inline') continue
    const local = (token.children || []).filter(child => child.type === 'image' && isLocalImage(child.meta?.source))
    if (!local.length) continue
    const mapping = []
    let search = Math.max(offsets[range[0]], cursor)
    for (const part of token.content.split('\n')) {
      const found = content.indexOf(part, search)
      if (found < 0 || found + part.length > offsets[range[1]]) throw new Error('无法定位图片路径，请检查 Markdown 格式')
      for (let i = 0; i < part.length; i++) mapping.push(found + i)
      search = found + part.length
      mapping.push(content.indexOf('\n', search))
      if (content[search] === '\r') search++
      if (content[search] === '\n') search++
    }
    cursor = search
    for (const image of local) {
      const span = image.meta
      const source = span.source
      const occurrence = { start: mapping[span.start], end: mapping[span.end - 1] + 1 }
      if (span.reference) {
        occurrence.prefix = content.slice(mapping[span.start], mapping[span.labelEnd] + 1)
        occurrence.title = span.title
      }
      if (!images.has(source)) images.set(source, { source, occurrences: [] })
      images.get(source).occurrences.push(occurrence)
    }
  }
  return [...images.values()]
}

function isLocalImage(source = '') {
  if (!source || /^(?:https?:|data:|blob:|\/\/)/i.test(source) || source.startsWith('/uploads/')) return false
  return /^file:/i.test(source) || /^[a-z]:[\\/]/i.test(source) || !/^[a-z][a-z\d+.-]*:/i.test(source)
}

function imagePath(source) {
  let path = source.replace(/[?#].*$/, '')
  try { path = decodeURIComponent(path) } catch { /* retain literal percent characters */ }
  return path.replace(/\\/g, '/')
}

export function matchLocalImage(source, markdownFile, files) {
  const path = imagePath(source)
  const candidates = files.filter(file => file.name === path.split('/').pop())
  const result = { file: null, candidates }
  const mdPath = markdownFile?.webkitRelativePath?.replace(/\\/g, '/')
  // Absolute paths still require an explicit choice, even when their tail matches.
  if (/^(?:file:|[a-z]:|\/)/i.test(path)) return result
  const parts = mdPath ? mdPath.split('/').slice(0, -1) : []
  for (const segment of path.split('/')) {
    if (!segment || segment === '.') continue
    if (segment === '..') {
      if (parts.length <= (mdPath ? 1 : 0)) return result
      parts.pop()
    } else parts.push(segment)
  }
  if (mdPath) {
    const exact = files.filter(file => file.webkitRelativePath === parts.join('/'))
    if (exact.length) {
      if (exact.length === 1) result.file = exact[0]
      return result
    }
  }
  const relative = (mdPath ? parts.slice(1) : parts).join('/')
  // Directory pickers include their selected root in webkitRelativePath. Align
  // that root with the article's relative path, retaining at least a folder name.
  if (!relative.includes('/')) return result
  const matches = candidates.filter(file => {
    const selectedPath = file.webkitRelativePath?.replace(/\\/g, '/') || ''
    if (!selectedPath.includes('/')) return false
    // Files already under the Markdown's known root must match the exact path.
    if (mdPath && selectedPath.split('/')[0] === mdPath.split('/')[0]) return false
    return selectedPath === relative || selectedPath.endsWith('/' + relative) || relative.endsWith('/' + selectedPath)
  })
  if (matches.length === 1) result.file = matches[0]
  return result
}

export function rewriteImages(content, images) {
  const replacements = images.flatMap(image => image.occurrences.map(occurrence => {
    let value = image.url
    if (occurrence.prefix) {
      const title = occurrence.title ? ` "${occurrence.title.replace(/\\/g, '\\\\').replace(/"/g, '\\"')}"` : ''
      value = `${occurrence.prefix}(${image.url}${title})`
    }
    return { ...occurrence, value }
  })).sort((a, b) => b.start - a.start)
  for (const replacement of replacements) content = content.slice(0, replacement.start) + replacement.value + content.slice(replacement.end)
  return content
}

export async function importImages({ content, images, upload, cache = new Map(), onProgress = () => {} }) {
  if (images.some(image => !image.file)) throw new Error('请先补齐或确认所有本地图片')
  for (const image of images) {
    image.error = ''
    try {
      image.status = 'uploading'
      onProgress()
      if (!cache.has(image.file)) {
        const url = await upload(image.file)
        if (!/^https?:\/\//i.test(url)) throw new Error('图片上传未返回有效地址')
        cache.set(image.file, url)
      }
      image.url = cache.get(image.file)
      image.status = 'uploaded'
      onProgress()
    } catch (error) {
      image.status = 'error'
      image.error = error.response?.data?.error || error.message || '图片上传失败，请重试'
      onProgress()
      throw error
    }
  }
  return rewriteImages(content, images)
}
