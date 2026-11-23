import MarkdownIt from 'markdown-it'

const markdown = new MarkdownIt({ html: false })

export function cleanTitle(title = '') {
  return String(title || '').replace(/\\([&_*~[\]()!#])/g, '$1').trim()
}

function stripContainerFences(source = '') {
  return String(source || '')
    .replace(/^[ \t]*:::[ \t]*[a-zA-Z0-9_-]*[^\n]*$/gm, '')
    .replace(/(?:^|\s):::(?:[a-zA-Z0-9_-]+)?(?=\s|$)/g, ' ')
}

function stripTruncatedLinks(text = '') {
  return String(text || '')
    .replace(/!\[[^\]]*\]\([^)]*$/g, '')
    .replace(/\[([^\]]+)\]\([^)]*$/g, '$1')
}

function inlineText(tokens = []) {
  return tokens.map(token => {
    if (token.type === 'image') return token.content
    if (token.children) return inlineText(token.children)
    if (['text', 'code_inline'].includes(token.type)) return token.content
    if (['softbreak', 'hardbreak'].includes(token.type)) return ' '
    return ''
  }).join('')
}

export function plainExcerpt(source = '', limit = 200) {
  const cleaned = stripTruncatedLinks(stripContainerFences(source))
  const text = markdown.parse(cleaned, {})
    .filter(token => token.type === 'inline')
    .map(token => inlineText(token.children)).join(' ')
    .replace(/\s+/g, ' ').trim()
  const characters = Array.from(text)
  return characters.length > limit ? characters.slice(0, limit).join('') + '...' : text
}

const CONTAINER_LABELS = {
  info: 'Info',
  tip: 'Tip',
  note: 'Note',
  warning: 'Warning',
  danger: 'Danger',
}

// Normalize container blocks (`:::info ... :::`), Windows image backslashes, and
// empty image tags outside fenced code blocks before shifting heading levels.
function preprocessArticleSource(source = '') {
  const lines = String(source).split('\n')
  const out = []
  let fenceChar = ''
  let fenceLen = 0
  let inContainer = false

  for (const rawLine of lines) {
    const fenceMatch = rawLine.match(/^[ \t]*(`{3,}|~{3,})/)
    if (fenceMatch) {
      const marker = fenceMatch[1]
      if (!fenceChar) {
        fenceChar = marker[0]
        fenceLen = marker.length
      } else if (marker[0] === fenceChar && marker.length >= fenceLen && /^[ \t]*$/.test(rawLine.slice(fenceMatch[0].length))) {
        fenceChar = ''
        fenceLen = 0
      }
      out.push(inContainer ? `> ${rawLine}` : rawLine)
      continue
    }

    if (fenceChar) {
      out.push(inContainer ? `> ${rawLine}` : rawLine)
      continue
    }

    const openContainer = rawLine.match(/^[ \t]*:::[ \t]*([a-zA-Z0-9_-]+)(?:[ \t]+(.+))?[ \t]*$/)
    if (!inContainer && openContainer) {
      inContainer = true
      const kind = openContainer[1].toLowerCase()
      const heading = (openContainer[2] || CONTAINER_LABELS[kind] || openContainer[1]).trim()
      out.push(`> **${heading}**`)
      out.push('>')
      continue
    }

    if (inContainer && /^[ \t]*:::[ \t]*$/.test(rawLine)) {
      inContainer = false
      continue
    }

    // Remove empty Markdown images `![]()` and normalize Windows `\` in image URLs.
    const normalized = rawLine
      .replace(/!\[\s*\]\(\s*\)/g, '')
      .replace(/!\[([^\]]*)\]\(([^)\s]+?)(\s+["'][^"']*["'])?\)/g, (_match, alt, url, titlePart = '') => {
        const cleanUrl = url.includes('\\') && !/^https?:\/\//i.test(url)
          ? url.replace(/\\+/g, '/')
          : url
        return `![${alt}](${cleanUrl}${titlePart})`
      })

    out.push(inContainer ? (normalized.trim() ? `> ${normalized}` : '>') : normalized)
  }

  return out.join('\n')
}

// Token maps identify real headings; fenced code is never rewritten.
export function articleMarkdown(source = '') {
  if (!source) return ''
  const preprocessed = preprocessArticleSource(source)
  const lines = preprocessed.split('\n')
  const tokens = markdown.parse(preprocessed, {})
  const headings = tokens.flatMap((token, index) =>
    token.type === 'heading_open' && token.map && tokens[index + 1]
      ? [{ token, text: tokens[index + 1].content }]
      : []
  )
  for (const { token, text } of headings.reverse()) {
    if (!token.map) continue
    const [start, end] = token.map
    if (start >= lines.length) continue
    const level = Math.min(Number(token.tag.slice(1)) + 1, 6)
    const match = lines[start].match(/^([ \t]*(?:(?:>[ \t]*|(?:[-+*]|\d+[.)])[ \t]+)[ \t]*)*)/)
    const prefix = match ? match[1] : ''
    lines.splice(start, end - start, `${prefix}${'#'.repeat(level)} ${text}`)
  }
  return lines.join('\n')
}
