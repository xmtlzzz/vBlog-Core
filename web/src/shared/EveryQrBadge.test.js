import { afterEach, beforeEach, describe, expect, it, vi } from 'vitest'
import { flushPromises, mount } from '@vue/test-utils'
import { createPinia } from 'pinia'
import EveryQrBadge from './EveryQrBadge.vue'
import { useThemeStore } from '../stores/theme'
import { createEveryQRCodeIdentity } from '@every-qrcode/core'
import { mountSeed } from '@every-qrcode/renderer-webgpu'

vi.mock('@every-qrcode/core', () => ({
  CURRENT_GENERATOR_VERSION: 1,
  createEveryQRCodeIdentity: vi.fn(async () => ({ qr: {} })),
  createQRSvgPath: vi.fn(() => ({ size: 33, path: 'M4 4h7v7H4z' }))
}))
vi.mock('@every-qrcode/renderer-webgpu', () => ({
  createSeedModel: vi.fn(async () => ({})),
  mountSeed: vi.fn(() => ({ setFlat: vi.fn(), setZoom: vi.fn(), setScene: vi.fn(), resize: vi.fn(), dispose: vi.fn() }))
}))
const wrappers = []
beforeEach(() => {
  vi.clearAllMocks()
  localStorage.clear()
  window.matchMedia = vi.fn(() => ({ matches: false }))
  vi.stubGlobal('ResizeObserver', class { observe() {} disconnect() {} })
})
afterEach(() => { wrappers.splice(0).forEach(wrapper => wrapper.unmount()); vi.unstubAllGlobals() })
async function render(props = {}) {
  const pinia = createPinia()
  const wrapper = mount(EveryQrBadge, { props: { url: 'https://example.com', ...props }, global: { plugins: [pinia] } })
  wrappers.push(wrapper)
  await flushPromises()
  return { wrapper, theme: useThemeStore(pinia) }
}

describe('QR theme and scan appearance', () => {
  it('keeps the native colored renderer visible after unfolding and switching themes', async () => {
    const { wrapper, theme } = await render()
    const renderer = mountSeed.mock.results[0].value
    expect(wrapper.classes()).not.toContain('is-flat')
    await wrapper.trigger('click')
    expect(wrapper.classes()).toContain('is-flat')
    expect(renderer.setFlat).toHaveBeenLastCalledWith(true)
    expect(renderer.setZoom).toHaveBeenLastCalledWith(1)
    expect(wrapper.find('canvas').isVisible()).toBe(true)
    expect(wrapper.find('svg').exists()).toBe(false)
    theme.toggle()
    await flushPromises()
    expect(createEveryQRCodeIdentity).toHaveBeenCalledExactlyOnceWith('https://example.com', { identityScope: 'url' })
    expect(mountSeed).toHaveBeenCalledTimes(1)
    expect(wrapper.find('svg').exists()).toBe(false)
    await wrapper.trigger('click')
    expect(wrapper.classes()).not.toContain('is-flat')
    expect(renderer.setFlat).toHaveBeenLastCalledWith(false)
    expect(wrapper.emitted('viewchange')).toEqual([[true], [false]])
  })
  it('uses theme paper for default white and rejects dark or translucent scan backgrounds', async () => {
    const { wrapper } = await render({ background: '#ffffff' })
    expect(wrapper.element.style.getPropertyValue('--qr-scan-bg')).toBe('')
    await wrapper.setProps({ background: '#111111' })
    expect(wrapper.element.style.getPropertyValue('--qr-scan-bg')).toBe('')
    await wrapper.setProps({ background: '#ffffff80' })
    expect(wrapper.element.style.getPropertyValue('--qr-scan-bg')).toBe('')
    await wrapper.setProps({ background: '#fff8ee' })
    expect(wrapper.element.style.getPropertyValue('--qr-scan-bg')).toBe('#fff8ee')
  })
  it('keeps a readable static QR surface when WebGPU fails', async () => {
    const { wrapper } = await render()
    mountSeed.mock.calls[0][4].onError(new Error('WebGPU unavailable'))
    await flushPromises()
    expect(wrapper.classes()).toContain('is-flat')
    expect(wrapper.get('svg path').attributes('d')).toBe('M4 4h7v7H4z')
    expect(wrapper.get('svg rect').attributes('fill')).toBe('var(--qr-scan-bg)')
    expect(wrapper.emitted('viewchange')).toEqual([[true]])
  })
})
