import { plainExcerpt } from './markdown'

let cachedSettings = null

function meta(attribute, key, content) {
  let element = document.head.querySelector(`meta[${attribute}="${key}"]`)
  if (!element) {
    element = document.createElement('meta')
    element.setAttribute(attribute, key)
    document.head.append(element)
  }
  element.content = content
}

export function updateMetadata(settings = {}, { title, description, article = false, path = window.location.pathname } = {}) {
  if (settings && Object.keys(settings).length > 0) {
    cachedSettings = settings
  }
  const effectiveSettings = cachedSettings || settings || {}
  const site = effectiveSettings.site_title || 'vBlog'
  const summary = plainExcerpt(description || effectiveSettings.description || effectiveSettings.site_description || '记录开发、技术与日常。', 160)
  document.title = title ? `${title} · ${site}` : site
  meta('name', 'description', summary)
  meta('property', 'og:title', title || site)
  meta('property', 'og:description', summary)
  meta('property', 'og:site_name', site)
  meta('property', 'og:type', article ? 'article' : 'website')
  meta('property', 'og:url', new URL(path, window.location.origin).href)
  let canonical = document.head.querySelector('link[rel="canonical"]')
  if (!canonical) { canonical = document.createElement('link'); canonical.rel = 'canonical'; document.head.append(canonical) }
  canonical.href = new URL(path, window.location.origin).href
}
