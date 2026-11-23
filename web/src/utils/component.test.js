import { describe, expect, it, vi } from 'vitest'
import { flushPromises, mount } from '@vue/test-utils'
import { createPinia } from 'pinia'
import Modules from '../blog/Modules.vue'
import { useThemeStore } from '../stores/theme'
import { buildSrcdoc } from './component'

vi.mock('../api/request', () => ({ default: { get: vi.fn(async url => url === '/components/active' ? [{ id: 1, name: '交互组件', code: 'function render(root) { root.innerHTML = "<input value=kept>" }' }] : {}) } }))

describe('embedded component themes', () => {
  it('applies the selected theme while retaining component code and blog data', () => {
    const code = 'function render(root) { root.textContent = window.__BLOG_DATA__.posts }'
    const source = buildSrcdoc(code, { theme: 'dark', data: { posts: 42 } })
    expect(source).toContain('<html data-theme="dark">')
    expect(source).toContain('color-scheme: dark')
    expect(source).toContain(code)
    expect(source).toContain('window.__BLOG_DATA__={"posts":42}')
  })
  it('defaults to light for an invalid theme and keeps empty previews empty', () => {
    expect(buildSrcdoc('function render() {}', { theme: '"><script>bad()</script>' })).toContain('<html data-theme="light">')
    expect(buildSrcdoc('')).toBe('')
  })
  it('updates embedded colors without reloading the component document or losing its state', async () => {
    localStorage.clear()
    window.matchMedia = vi.fn(() => ({ matches: false }))
    const pinia = createPinia()
    const wrapper = mount(Modules, { global: { plugins: [pinia], stubs: { BlogNav: true, BlogFooter: true, QrLinksSection: true } } })
    try {
      await flushPromises()
      const frame = wrapper.get('iframe')
      const source = frame.attributes('srcdoc')
      const message = vi.fn()
      Object.defineProperty(frame.element, 'contentWindow', { value: { postMessage: message }, configurable: true })
      useThemeStore(pinia).toggle()
      await flushPromises()
      expect(wrapper.get('iframe').attributes('srcdoc')).toBe(source)
      expect(message).toHaveBeenCalledWith({ type: 'vblog-theme', theme: 'dark' }, '*')
    } finally { wrapper.unmount() }
  })
})
