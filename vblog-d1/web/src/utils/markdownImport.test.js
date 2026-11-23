import { describe, expect, it, vi } from 'vitest'
import { scanLocalImages, matchLocalImage, rewriteImages, importImages } from './markdownImport'

function file(path) {
  const result = new File(['image'], path.split('/').pop(), { type: 'image/png' })
  Object.defineProperty(result, 'webkitRelativePath', { value: path })
  return result
}

describe('Markdown image import', () => {
  it('finds rendered local images, including references, while preserving code and remote URLs', () => {
    const content = '# 文章\n\n![截图](./assets/a(1).png "说明")\n![第二张][shot]\n\n[shot]: <./assets/中文 图.png> "引用说明"\n\n[下载](./assets/a(1).png)\n![远程](https://example.com/a.png)\n![本站](/uploads/old.png)\n`![示例](ignore.png)`\n\n```md\n![示例](ignore.png)\n```'
    const images = scanLocalImages(content)
    expect(images.map(image => image.source)).toEqual(['./assets/a(1).png', './assets/中文 图.png'])
    const rewritten = rewriteImages(content, images.map((image, i) => ({ ...image, url: `https://blog.example/uploads/${i}.png` })))
    expect(rewritten).toContain('![截图](https://blog.example/uploads/0.png "说明")')
    expect(rewritten).toContain('![第二张](https://blog.example/uploads/1.png "引用说明")')
    expect(rewritten).toContain('[下载](./assets/a(1).png)')
    expect(rewritten).toContain('[shot]: <./assets/中文 图.png> "引用说明"')
    expect(rewritten).toContain('`![示例](ignore.png)`')
    expect(rewritten).toContain('```md\n![示例](ignore.png)\n```')
  })

  it('preserves CRLF, quotes, lists and tables and rewrites each occurrence only', () => {
    const content = '> ![甲](./a.png)\r\n\r\n- ![乙](./a.png)\r\n\r\n| 图片 | 说明 |\r\n| --- | --- |\r\n| ![丙](./b.png) | 原文 |\r\n'
    const images = scanLocalImages(content)
    expect(images.map(image => image.source)).toEqual(['./a.png', './b.png'])
    expect(rewriteImages(content, images.map(image => ({ ...image, url: 'https://blog.example/x.png' })))).toBe(content.replaceAll('./a.png', 'https://blog.example/x.png').replaceAll('./b.png', 'https://blog.example/x.png'))
  })

  it('detects Windows and file URLs for explicit selection, and ignores escaped Markdown', () => {
    const images = scanLocalImages('![本地](C:\\notes\\a.png)\n![文件](file:///C:/notes/b.png)\n\\![不是图片](ignore.png)')
    expect(images.map(image => image.source)).toEqual(['C:\\notes\\a.png', 'file:///C:/notes/b.png'])
    expect(matchLocalImage(images[0].source, file('notes/article.md'), [file('notes/a.png')]).file).toBeNull()
  })

  it('resolves encoded relative paths from the Markdown directory and never guesses by basename', () => {
    const md = file('notes/posts/article.md')
    const image = file('notes/assets/中文 图.png')
    const duplicate = file('notes/other/中文 图.png')
    expect(matchLocalImage('../assets/%E4%B8%AD%E6%96%87%20%E5%9B%BE.png', md, [image, duplicate]).file).toBe(image)
    const ambiguous = matchLocalImage('./中文 图.png', md, [image, duplicate])
    expect(ambiguous.file).toBeNull()
    expect(ambiguous.candidates).toEqual([image, duplicate])
    expect(matchLocalImage('../../outside.png', md, [file('outside.png')]).file).toBeNull()
    expect(matchLocalImage('a.png', new File([''], 'article.md'), [file('notes/a.png')]).file).toBeNull()
  })

  it('uploads only referenced files, deduplicates a file and caches success across retries', async () => {
    const content = '![甲](./a.png)\n![乙](a.png)\n![丙](b.png)'
    const a = file('notes/a.png'), b = file('notes/b.png')
    const images = scanLocalImages(content).map(image => ({ ...image, file: image.source.includes('a.png') ? a : b }))
    const upload = vi.fn().mockResolvedValueOnce('https://blog.example/a.png').mockRejectedValueOnce(new Error('网络中断')).mockResolvedValueOnce('https://blog.example/b.png')
    const cache = new Map()
    await expect(importImages({ content, images, upload, cache })).rejects.toThrow('网络中断')
    expect(content).toContain('./a.png')
    expect(upload).toHaveBeenCalledTimes(2)
    const result = await importImages({ content, images, upload, cache })
    expect(upload).toHaveBeenCalledTimes(3)
    expect(result).toBe('![甲](https://blog.example/a.png)\n![乙](https://blog.example/a.png)\n![丙](https://blog.example/b.png)')
  })

  it('matches Windows relative paths when a separately selected Markdown is supplemented with a directory', () => {
    const md = new File([''], 'article.md')
    for (const path of ['notes/assets/images/shot.png', 'assets/images/shot.png', 'images/shot.png']) {
      const image = file(path)
      expect(matchLocalImage('assets\\images\\shot.png', md, [image]).file).toBe(image)
    }
    const image = file('images/中文 图.png')
    expect(matchLocalImage('./assets/images/%E4%B8%AD%E6%96%87%20%E5%9B%BE.png', md, [image]).file).toBe(image)
  })

  it('matches supplemental directories against the resolved Markdown path, with exact matches taking priority', () => {
    const md = file('notes/posts/article.md')
    const extra = file('images/shot.png')
    const exact = file('notes/assets/images/shot.png')
    expect(matchLocalImage('../assets/images/shot.png', md, [extra]).file).toBe(extra)
    expect(matchLocalImage('../assets/images/shot.png', md, [extra, exact]).file).toBe(exact)
    expect(matchLocalImage('../assets/images/shot.png', md, [file('notes/other/assets/images/shot.png')]).file).toBeNull()
  })

  it('leaves ambiguous directory paths, unmatched folders and absolute paths for explicit selection', () => {
    const md = new File([''], 'article.md')
    const first = file('one/assets/images/shot.png'), second = file('two/assets/images/shot.png')
    const ambiguous = matchLocalImage('assets/images/shot.png', md, [first, second])
    expect(ambiguous.file).toBeNull()
    expect(ambiguous.candidates).toEqual([first, second])
    expect(matchLocalImage('assets/images/shot.png', md, [file('other/shot.png')]).file).toBeNull()
    expect(matchLocalImage('C:\\assets\\images\\shot.png', md, [first]).file).toBeNull()
    expect(matchLocalImage('file:///C:/assets/images/shot.png', md, [first]).file).toBeNull()
  })

  it('blocks missing files before uploading anything', async () => {
    const upload = vi.fn()
    await expect(importImages({ content: '![图](missing.png)', images: scanLocalImages('![图](missing.png)'), upload })).rejects.toThrow('图片')
    expect(upload).not.toHaveBeenCalled()
  })
})
