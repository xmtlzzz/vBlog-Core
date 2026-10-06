import MarkdownIt from 'markdown-it'

const markdown = new MarkdownIt({ html: false })

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
  const text = markdown.parse(String(source), {})
    .filter(token => token.type === 'inline')
    .map(token => inlineText(token.children)).join(' ')
    .replace(/\s+/g, ' ').trim()
  const characters = Array.from(text)
  return characters.length > limit ? characters.slice(0, limit).join('') + '...' : text
}

// Token maps identify real headings; fenced code is never rewritten.
export function articleMarkdown(source = '') {
  if (!source) return ''
  const lines = source.split('\n')
  const tokens = markdown.parse(source, {})
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
