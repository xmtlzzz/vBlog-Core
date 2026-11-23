import { describe, expect, it } from 'vitest'
import { articleMarkdown, cleanTitle, plainExcerpt } from './markdown'

describe('article Markdown', () => {
  it('moves quoted and list headings below the article H1 while preserving their containers', () => {
    expect(articleMarkdown('> # 引用标题\n\n- # 列表标题')).toBe('> ## 引用标题\n\n- ## 列表标题')
  })
  it('normalizes setext headings without rewriting fenced code', () => {
    expect(articleMarkdown('标题\n===\n\n```md\n# 原样保留\n```')).toBe('## 标题\n\n```md\n# 原样保留\n```')
  })
  it('converts :::info containers, normalizes Windows image slashes, and removes empty images', () => {
    expect(articleMarkdown(':::info\n第一行提示\n:::\n\n![图](assets\\sub\\a.png)\n![]()')).toBe('> **Info**\n>\n> 第一行提示\n\n![图](assets/sub/a.png)\n')
  })
  it('unescapes Markdown punctuation in titles', () => {
    expect(cleanTitle('whatis\\&whereis')).toBe('whatis&whereis')
  })
  it('uses readable prose and Unicode character boundaries for excerpts', () => {
    expect(plainExcerpt('## [音乐](https://example.com) **练习**\n\n```js\nignore()\n```')).toBe('音乐 练习')
    expect(plainExcerpt(':::info\n提示内容\n:::\n\nSee [this commit](https://example.com/trun')).toBe('提示内容 See this commit')
    expect(plainExcerpt('😀音乐🎸', 2)).toBe('😀音...')
  })
})
