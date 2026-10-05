// @vitest-environment jsdom
import { afterEach, beforeEach, describe, expect, it, vi } from 'vitest'
import { mount, flushPromises } from '@vue/test-utils'
import { createPinia } from 'pinia'
import { createMemoryHistory, createRouter } from 'vue-router'
import Home from './Home.vue'
import Post from './Post.vue'
import BlogNav from '../shared/BlogNav.vue'
import PostCard from '../shared/PostCard.vue'
import api from '../api/request'

vi.mock('../api/request', () => ({ default: { get: vi.fn() } }))
const wrappers = []
const post = { id: 3, title: '项目导航', content: '# 展示页\n\n```md\n# keep code\n```', excerpt: '# 项目 [导航](https://example.com) **介绍**', tags: [], created_at: '2026-08-31' }

async function render(component, path = '/', props = {}) {
  const router = createRouter({ history: createMemoryHistory(), routes: [{ path: '/:pathMatch(.*)*', component: { template: '<div />' } }, { path: '/post/:id', component: { template: '<div />' } }] })
  await router.push(path)
  await router.isReady()
  const wrapper = mount(component, { props, global: { plugins: [createPinia(), router], stubs: { BlogFooter: true, CustomWidgets: true, CommentSection: true, MdCatalog: true, MdPreview: { props: ['modelValue'], template: '<div class="preview-source">{{ modelValue }}</div>' }, ElPagination: true } } })
  wrappers.push(wrapper)
  await flushPromises()
  return { wrapper, router }
}

beforeEach(() => {
  vi.useFakeTimers()
  document.title = 'vBlog'
  window.scrollTo = vi.fn()
  Element.prototype.scrollIntoView = vi.fn()
  window.matchMedia = vi.fn(() => ({ matches: true, addEventListener: vi.fn(), removeEventListener: vi.fn() }))
  api.get.mockImplementation(async (url, options) => {
    if (url === '/settings') return { site_title: 'vBlog', author_name: 'xmtlz' }
    if (url === '/dashboard/stats') return { total_posts: 1, total_views: 1, total_tags: 0 }
    if (url === '/posts') return { data: options?.params?.search ? [] : [post], total: options?.params?.search ? 0 : 1 }
    return post
  })
})
afterEach(() => { wrappers.splice(0).forEach(w => w.unmount()); vi.clearAllTimers(); vi.useRealTimers(); vi.clearAllMocks() })

describe('public blog experience', () => {
  it('shows readable excerpts without Markdown link or heading syntax', async () => {
    const { wrapper } = await render(PostCard, '/', { post })
    expect(wrapper.get('.post-excerpt').text()).toBe('项目 导航 介绍')
  })
  it('distinguishes no search matches and clears back to the article list', async () => {
    const { wrapper } = await render(Home)
    await wrapper.get('input').setValue('missing')
    await vi.advanceTimersByTimeAsync(350)
    await flushPromises()
    expect(wrapper.text()).toContain('没有找到匹配文章')
    await wrapper.get('.empty-state button').trigger('click')
    await flushPromises()
    expect(wrapper.get('input').element.value).toBe('')
    expect(wrapper.text()).toContain('项目导航')
  })
  it('does not intercept the browser Find shortcut', async () => {
    await render(Home)
    const event = new KeyboardEvent('keydown', { key: 'f', ctrlKey: true, cancelable: true })
    window.dispatchEvent(event)
    expect(event.defaultPrevented).toBe(false)
  })
  it('exposes a menu disclosure and closes it after navigation', async () => {
    const { wrapper, router } = await render(BlogNav)
    const menu = wrapper.find('button[aria-controls="blog-navigation"]')
    expect(menu.exists()).toBe(true)
    await menu.trigger('click')
    expect(menu.attributes('aria-expanded')).toBe('true')
    await router.push('/about')
    await flushPromises()
    expect(menu.attributes('aria-expanded')).toBe('false')
  })
  it('uses article metadata and reserves H1 for the article title', async () => {
    const { wrapper } = await render(Post, '/post/3')
    expect(document.title).toBe('项目导航 · vBlog')
    expect(wrapper.get('.author-name').text()).toBe('xmtlz')
    expect(document.querySelector('meta[property="og:title"]')?.content).toBe('项目导航')
    expect(wrapper.get('.preview-source').text()).toContain('## 展示页')
    expect(wrapper.get('.preview-source').text()).toContain('# keep code')
  })
})
