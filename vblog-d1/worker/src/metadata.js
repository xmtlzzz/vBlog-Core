import { buildExcerpt } from './util.js';

const titles = { '/archives': '归档', '/modules': '模块', '/tags': '标签', '/friends': '友情链接', '/about': '关于' };
const escape = value => String(value).replace(/[&<>"']/g, character => ({ '&': '&amp;', '<': '&lt;', '>': '&gt;', '"': '&quot;', "'": '&#39;' })[character]);

export async function pageMetadata(env, url) {
  const rows = await env.DB.prepare("SELECT key, value FROM settings WHERE key IN ('site_title','site_description','description','site_url')").all();
  const settings = Object.fromEntries((rows.results || []).map(row => [row.key, row.value]));
  let origin = url.origin;
  try {
    const configured = new URL(settings.site_url || env.SITE_URL);
    if (['http:', 'https:'].includes(configured.protocol)) origin = configured.origin;
  } catch { /* Local preview uses its own origin. */ }
  const page = { site: settings.site_title || 'vBlog', title: titles[url.pathname] || '', description: buildExcerpt(settings.description || settings.site_description || '记录开发、技术与日常。', 160), url: new URL(url.pathname, origin).href, article: false, missing: false };
  if (url.pathname.startsWith('/post/')) {
    const match = url.pathname.match(/^\/post\/(\d+)$/);
    const post = match ? await env.DB.prepare("SELECT title, excerpt, content FROM posts WHERE id = ? AND status = 'published' AND deleted_at IS NULL").bind(Number(match[1])).first() : null;
    if (post) {
      page.title = post.title;
      page.description = buildExcerpt(post.excerpt || post.content, 160);
      page.article = true;
    } else { page.title = '文章不存在'; page.missing = true; }
  }
  return page;
}

export function rewriteMetadata(response, page) {
  if (!response.headers.get('content-type')?.includes('text/html') || response.status !== 200) return response;
  const title = page.title ? `${page.title} · ${page.site}` : page.site;
  const tags = `<title>${escape(title)}</title><meta name="description" content="${escape(page.description)}"><meta property="og:title" content="${escape(page.title || page.site)}"><meta property="og:description" content="${escape(page.description)}"><meta property="og:site_name" content="${escape(page.site)}"><meta property="og:type" content="${page.article ? 'article' : 'website'}"><meta property="og:url" content="${escape(page.url)}"><link rel="canonical" href="${escape(page.url)}">${page.missing ? '<meta name="robots" content="noindex">' : ''}`;
  const headers = new Headers(response.headers);
  headers.delete('etag');
  headers.delete('content-length');
  headers.set('cache-control', 'no-cache');
  const rewritten = new Response(response.body, { status: page.missing ? 404 : response.status, headers });
  return new HTMLRewriter()
    .on('title, meta[name="description"], meta[property^="og:"], link[rel="canonical"]', { element(element) { element.remove(); } })
    .on('head', { element(element) { element.append(tags, { html: true }); } })
    .transform(rewritten);
}
