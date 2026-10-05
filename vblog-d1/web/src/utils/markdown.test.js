import { describe, expect, it } from 'vitest'
import { articleMarkdown, plainExcerpt } from './markdown'

describe('article Markdown', () => {
  it('moves quoted and list headings below the article H1 while preserving their containers', () => {
    expect(articleMarkdown('> # 引用标题\n\n- # 列表标题')).toBe('> ## 引用标题\n\n- ## 列表标题')
  })
  it('normalizes setext headings without rewriting fenced code', () => {
    expect(articleMarkdown('标题\n===\n\n```md\n# 原样保留\n```')).toBe('## 标题\n\n```md\n# 原样保留\n```')
  })
  it('uses readable prose and Unicode character boundaries for excerpts', () => {
    expect(plainExcerpt('## [音乐](https://example.com) **练习**\n\n```js\nignore()\n```')).toBe('音乐 练习')
    expect(plainExcerpt('😀音乐🎸', 2)).toBe('😀音...')
  })
})
