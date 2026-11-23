import { afterEach, beforeEach, describe, expect, it, vi } from 'vitest'
import { flushPromises, mount } from '@vue/test-utils'
import MarkdownImport from './MarkdownImport.vue'
import { uploadImageFile } from '../utils/imageUpload'

vi.mock('../utils/imageUpload', () => ({ uploadImageFile: vi.fn() }))
const wrappers = []
function file(path, text = 'image') {
  const result = new File([text], path.split('/').pop(), { type: path.endsWith('.md') ? 'text/markdown' : 'image/png' })
  Object.defineProperty(result, 'webkitRelativePath', { value: path.includes('/') ? path : '' })
  return result
}
function render() {
  const wrapper = mount(MarkdownImport, { global: { stubs: {
    ElDialog: { props: ['modelValue'], template: '<div v-if="modelValue"><slot /><slot name="footer" /></div>' },
    ElButton: { props: ['disabled', 'loading'], template: '<button :disabled="disabled || loading"><slot /></button>' },
    ElAlert: { props: ['title'], template: '<p role="alert">{{ title }}</p>' },
    ElProgress: true,
  } } })
  wrappers.push(wrapper)
  return wrapper
}
async function select(wrapper, index, files) {
  const input = wrapper.findAll('input')[index]
  Object.defineProperty(input.element, 'files', { configurable: true, value: files })
  await input.trigger('change')
  await flushPromises()
  await vi.waitFor(() => expect(wrapper.text()).not.toContain('正在读取文章'))
}
function submit(wrapper) { return wrapper.findAll('button').find(button => button.text().includes('进入编辑')) }
beforeEach(() => { uploadImageFile.mockReset(); uploadImageFile.mockResolvedValue('https://blog.example/uploads/a.webp') })
afterEach(() => wrappers.splice(0).forEach(wrapper => wrapper.unmount()))

describe('Markdown import dialog', () => {
  it('uploads the referenced file before emitting rewritten content for the editor', async () => {
    const wrapper = render()
    const image = file('notes/assets/a.png')
    await select(wrapper, 1, [file('notes/article.md', '# 标题\n![图](assets/a.png)'), image, file('notes/unused.png')])
    expect(uploadImageFile).not.toHaveBeenCalled()
    await submit(wrapper).trigger('click')
    await flushPromises()
    expect(uploadImageFile).toHaveBeenCalledExactlyOnceWith(image)
    expect(wrapper.emitted('imported')[0][0]).toEqual({ title: '标题', content: '# 标题\n![图](https://blog.example/uploads/a.webp)' })
  })

  it('blocks missing images and requires explicit confirmation for absolute paths', async () => {
    const wrapper = render()
    await select(wrapper, 1, [file('notes/article.md', '![图](C:\\notes\\a.png)'), file('notes/a.png')])
    expect(submit(wrapper).attributes('disabled')).toBeDefined()
    expect(wrapper.text()).toContain('1 张待补选或确认')
    await wrapper.get('.image-match select').setValue('1')
    expect(submit(wrapper).attributes('disabled')).toBeUndefined()
    await submit(wrapper).trigger('click')
    await flushPromises()
    expect(wrapper.emitted('imported')).toHaveLength(1)
  })

  it('offers a Markdown chooser instead of silently picking one article', async () => {
    const wrapper = render()
    await select(wrapper, 1, [file('notes/one.md', '# 一'), file('notes/two.md', '# 二')])
    expect(submit(wrapper).attributes('disabled')).toBeDefined()
    await wrapper.get('.article-choice select').setValue('notes/two.md')
    await vi.waitFor(() => expect(submit(wrapper).attributes('disabled')).toBeUndefined())
    await submit(wrapper).trigger('click')
    await flushPromises()
    expect(wrapper.emitted('imported')[0][0].title).toBe('二')
  })

  it('retains successful uploads when a later upload fails and retries only the failure', async () => {
    uploadImageFile.mockReset().mockResolvedValueOnce('https://blog.example/a.png').mockRejectedValueOnce(new Error('连接中断')).mockResolvedValueOnce('https://blog.example/b.png')
    const wrapper = render()
    await select(wrapper, 1, [file('notes/article.md', '![甲](a.png)\n![乙](b.png)'), file('notes/a.png'), file('notes/b.png')])
    await submit(wrapper).trigger('click')
    await flushPromises()
    expect(wrapper.emitted('imported')).toBeUndefined()
    expect(wrapper.text()).toContain('连接中断')
    await submit(wrapper).trigger('click')
    await flushPromises()
    expect(uploadImageFile).toHaveBeenCalledTimes(3)
    expect(wrapper.emitted('imported')[0][0].content).toBe('![甲](https://blog.example/a.png)\n![乙](https://blog.example/b.png)')
  })

  it('keeps the single Markdown workflow for articles without local images', async () => {
    const wrapper = render()
    const content = '# 纯文字\n`![示例](code.png)`'
    await select(wrapper, 0, [file('article.md', content)])
    expect(wrapper.emitted('imported')[0][0]).toEqual({ title: '纯文字', content })
    expect(uploadImageFile).not.toHaveBeenCalled()
  })

  it('automatically selects matching relative paths after supplementing a single Markdown with an image directory', async () => {
    const wrapper = render()
    await select(wrapper, 0, [file('article.md', '# 相对路径\n![图](assets\\images\\shot.png)')])
    expect(submit(wrapper).attributes('disabled')).toBeDefined()
    const image = file('images/shot.png')
    await select(wrapper, 2, [image, file('images/unused.png')])
    expect(wrapper.get('.image-match select').element.value).toBe('1')
    expect(submit(wrapper).attributes('disabled')).toBeUndefined()
    expect(uploadImageFile).not.toHaveBeenCalled()
    await submit(wrapper).trigger('click')
    await flushPromises()
    expect(uploadImageFile).toHaveBeenCalledExactlyOnceWith(image)
    expect(wrapper.emitted('imported')[0][0].content).toBe('# 相对路径\n![图](https://blog.example/uploads/a.webp)')
  })

  it('asks for a choice when relative paths are ambiguous and keeps a manual choice when more files are added', async () => {
    const wrapper = render()
    await select(wrapper, 0, [file('article.md', '![图](assets/images/shot.png)')])
    const first = file('one/assets/images/shot.png'), second = file('two/assets/images/shot.png')
    await select(wrapper, 2, [first, second])
    expect(submit(wrapper).attributes('disabled')).toBeDefined()
    expect(wrapper.get('.image-match select').element.value).toBe('')
    await wrapper.get('.image-match select').setValue('2')
    await select(wrapper, 2, [file('images/shot.png')])
    expect(wrapper.get('.image-match select').element.value).toBe('2')
    await submit(wrapper).trigger('click')
    await flushPromises()
    expect(uploadImageFile).toHaveBeenCalledExactlyOnceWith(second)
  })
})
