import assert from 'node:assert/strict';
import { signJWT } from '../src/util.js';
const base = process.env.VBLOG_TEST_URL || 'http://127.0.0.1:8787';
assert.ok(['localhost', '127.0.0.1'].includes(new URL(base).hostname), 'Runtime regression must target a local test server');
const token = await signJWT('local-ui-regression-only', { user_id: 1, username: 'local-test' }, 300);
const headers = { Authorization: `Bearer ${token}`, 'Content-Type': 'application/json' };
const publicPosts = await (await fetch(`${base}/api/posts`)).json();
assert.ok(publicPosts.data.every(post => post.status === 'published'), 'anonymous lists must exclude drafts without requiring a status filter');
assert.equal((await fetch(`${base}/api/posts/4`)).status, 404, 'anonymous visitors cannot open a draft by ID');
assert.equal((await fetch(`${base}/api/posts/4`, { headers })).status, 200, 'authenticated editors can still open drafts');
for (const requestHeaders of [{}, { 'Sec-Fetch-Mode': 'navigate', Accept: 'text/html' }]) {
  const articleHTML = await (await fetch(`${base}/post/3`, { headers: requestHeaders })).text();
  assert.ok(articleHTML.includes('<title>项目导航 · vBlog</title>'));
  assert.ok(articleHTML.includes('content="article"'));
  assert.equal((articleHTML.match(/<title>/g) || []).length, 1);
}
const original = await (await fetch(`${base}/api/posts/3`, { headers })).json();
try {
  const save = await fetch(`${base}/api/posts/3`, { method: 'PUT', headers, body: JSON.stringify({ ...original, status: 'draft' }) });
  assert.equal(save.status, 200);
  const draft = await save.json();
  assert.equal(draft.title, original.title);
  assert.equal(draft.content, original.content);
  assert.deepEqual(draft.tags, original.tags);
  assert.equal((await fetch(`${base}/api/posts/3`)).status, 404);
  assert.equal((await fetch(`${base}/post/3`)).status, 404);
  assert.equal((await fetch(`${base}/post/3`, { headers: { 'Sec-Fetch-Mode': 'navigate', Accept: 'text/html' } })).status, 404);
} finally {
  const restore = await fetch(`${base}/api/posts/3`, { method: 'PUT', headers, body: JSON.stringify(original) });
  assert.equal(restore.status, 200);
}
console.log('Local D1 runtime: public/draft separation, editor access, reversible draft conversion, and server HTML metadata passed.');
