import assert from 'node:assert/strict';
import { buildExcerpt } from './util.js';
import { pageMetadata } from './metadata.js';

assert.equal(buildExcerpt('# 项目\n\n[导航](https://example.com) **介绍**'), '项目 导航 介绍');
assert.equal(buildExcerpt('简介\n\n```js\nconst secret = 1;\n```\n\n更多内容'), '简介 更多内容');
assert.equal(buildExcerpt('**你好世界**', 3), '你好世...');
assert.equal(buildExcerpt('😀音乐🎸', 2), '😀音...');
const env = { SITE_URL: 'https://vblog.xmtlz.dev', DB: { prepare(sql) {
  return { all: async () => ({ results: [{ key: 'site_title', value: 'vBlog' }] }), bind(id) { assert.ok(sql.includes("status = 'published'")); return { first: async () => id === 3 ? { title: '项目导航', content: '# 简介\n\n[项目](https://example.com)' } : null }; } };
} } };
const article = await pageMetadata(env, new URL('http://localhost/post/3'));
assert.equal(article.title, '项目导航');
assert.equal(article.description, '简介 项目');
assert.equal(article.url, 'https://vblog.xmtlz.dev/post/3');
assert.equal((await pageMetadata(env, new URL('http://localhost/post/4'))).missing, true);
console.log('Markdown excerpts and published-article metadata passed');
