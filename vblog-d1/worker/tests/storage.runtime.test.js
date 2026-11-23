import assert from 'node:assert/strict';
import { after, test } from 'node:test';
import * as miniflare from 'miniflare';
import { signJWT } from '../src/util.js';

// Isolated in-memory resources. Never uses the production Wrangler bindings.
const secret = 'local-r2-runtime-only';
const options = {
  modules: true, scriptPath: '.wrangler/test-bundle/index.js',
  compatibilityDate: '2025-08-01',
  bindings: { JWT_SECRET: secret, R2_PUBLIC_URL: '' },
  r2Buckets: ['R2'], kvNamespaces: ['IMG'], d1Databases: ['DB'],
};
const mf = new miniflare.Miniflare(miniflare.convertV4MiniflareOptions ? miniflare.convertV4MiniflareOptions(options) : options);
after(() => mf.dispose());

test('workerd uploads into R2 and serves the same bytes, HEAD, cache validation and legacy KV', async () => {
  const token = await signJWT(secret, { user_id: 1, username: 'local-test' }, 300);
  const png = Buffer.from('iVBORw0KGgoAAAANSUhEUgAAAAEAAAABCAQAAAC1HAwCAAAAC0lEQVR42mP8/x8AAwMCAO+/l9sAAAAASUVORK5CYII=', 'base64');
  const form = new FormData();
  form.append('file', new File([png], 'test.png', { type: 'image/png' }));
  const payload = new Request('https://local.example/api/upload', { method: 'POST', body: form });
  const uploaded = await mf.dispatchFetch(payload.url, { method: 'POST', headers: { Authorization: `Bearer ${token}`, 'Content-Type': payload.headers.get('Content-Type') }, body: await payload.arrayBuffer() });
  assert.equal(uploaded.status, 200, await uploaded.clone().text());
  const { url } = await uploaded.json();
  const key = new URL(url).pathname.slice(1);
  const r2 = await mf.getR2Bucket('R2');
  assert.deepEqual(Buffer.from(await (await r2.get(key)).arrayBuffer()), png);
  const response = await mf.dispatchFetch(url);
  assert.equal(response.status, 200);
  assert.equal(response.headers.get('Content-Type'), 'image/png');
  assert.deepEqual(Buffer.from(await response.arrayBuffer()), png);
  const etag = response.headers.get('ETag');
  assert.ok(etag);
  const head = await mf.dispatchFetch(url, { method: 'HEAD' });
  assert.equal(head.status, 200);
  assert.equal(head.headers.get('Content-Length'), String(png.length));
  assert.equal(await head.text(), '');
  assert.equal((await mf.dispatchFetch(url, { headers: { 'If-None-Match': etag } })).status, 304);
  const kv = await mf.getKVNamespace('IMG');
  assert.equal(await kv.get(key), null);
  await kv.put('uploads/legacy.png', png);
  const legacy = await mf.dispatchFetch('https://local.example/uploads/legacy.png');
  assert.deepEqual(Buffer.from(await legacy.arrayBuffer()), png);
  const health = await mf.dispatchFetch('https://local.example/api/health');
  assert.equal(health.status, 200);
  assert.equal((await health.json()).components.r2.status, 'UP');
});
