import assert from 'node:assert/strict';
import { test } from 'node:test';
import worker from '../src/index.js';
import { signJWT } from '../src/util.js';

const png = Buffer.from('iVBORw0KGgoAAAANSUhEUgAAAAEAAAABCAQAAAC1HAwCAAAAC0lEQVR42mP8/x8AAwMCAO+/l9sAAAAASUVORK5CYII=', 'base64');
const secret = 'local-storage-regression-only';
const token = await signJWT(secret, { user_id: 1, username: 'test' }, 300);
const auth = { Authorization: `Bearer ${token}` };
function uploadRequest(data = png, name = 'photo.png', type = 'image/png', headers = auth) {
  const form = new FormData();
  form.append('file', new File([data], name, { type }));
  return new Request('https://blog.example/api/upload', { method: 'POST', headers, body: form });
}
function environment() {
  const objects = new Map();
  return {
    JWT_SECRET: secret,
    objects,
    R2: {
      async put(key, body, options) { objects.set(key, { body, options }); },
      async get(key) {
        const obj = objects.get(key);
        return obj ? { body: new Response(obj.body).body, size: obj.body.byteLength, httpEtag: '"test-etag"', writeHttpMetadata(headers) { headers.set('Content-Type', obj.options.httpMetadata.contentType); } } : null;
      },
    },
    IMG: { async put() { throw new Error('New images must not enter KV'); }, async get() { return null; } },
  };
}

test('private R2 works without a public URL and serves image bytes, MIME, ETag and HEAD', async () => {
  const env = environment();
  const uploaded = await worker.fetch(uploadRequest(png, 'wrong.html'), env, {});
  assert.equal(uploaded.status, 200);
  const { url } = await uploaded.json();
  assert.match(url, /^https:\/\/blog.example\/uploads\/.+\.png$/);
  const response = await worker.fetch(new Request(url), env, {});
  assert.equal(response.status, 200);
  assert.equal(response.headers.get('Content-Type'), 'image/png');
  assert.equal(response.headers.get('ETag'), '"test-etag"');
  assert.deepEqual(Buffer.from(await response.arrayBuffer()), png);
  assert.equal((await worker.fetch(new Request(url, { method: 'HEAD' }), env, {})).headers.get('Content-Length'), String(png.length));
  const cached = await worker.fetch(new Request(url, { headers: { 'If-None-Match': '"test-etag"' } }), env, {});
  assert.equal(cached.status, 304);
});

test('R2 public URLs remain supported and old KV images remain readable', async () => {
  const env = environment();
  env.R2_PUBLIC_URL = 'https://img.example/';
  const response = await worker.fetch(uploadRequest(), env, {});
  assert.match((await response.json()).url, /^https:\/\/img.example\/uploads\//);
  env.IMG.get = async key => key === 'uploads/old.png' ? png : null;
  const legacy = await worker.fetch(new Request('https://blog.example/uploads/old.png'), env, {});
  assert.deepEqual(Buffer.from(await legacy.arrayBuffer()), png);
  assert.equal((await worker.fetch(new Request('https://blog.example/uploads/missing.png'), env, {})).status, 404);
});

test('upload requires authentication, rejects fake images and oversized files', async () => {
  const env = environment();
  assert.equal((await worker.fetch(uploadRequest(png, 'a.png', 'image/png', {}), env, {})).status, 401);
  assert.equal((await worker.fetch(uploadRequest('<script>x</script>'), env, {})).status, 400);
  assert.equal((await worker.fetch(uploadRequest(new Uint8Array(10 * 1024 * 1024 + 1)), env, {})).status, 400);
  assert.equal(env.objects.size, 0);
});

test('KV-only installations still upload successfully', async () => {
  const env = environment();
  delete env.R2;
  let stored;
  env.IMG.put = async (key, body) => { stored = { key, body }; };
  const response = await worker.fetch(uploadRequest(), env, {});
  assert.equal(response.status, 200);
  assert.ok(stored.key.startsWith('uploads/'));
  assert.deepEqual(Buffer.from(stored.body), png);
});

test('R2 failure returns an error without silently writing new images to KV', async () => {
  const env = environment();
  env.R2.put = async () => { throw new Error('storage unavailable'); };
  assert.equal((await worker.fetch(uploadRequest(), env, {})).status, 500);
  env.R2.get = async () => { throw new Error('storage unavailable'); };
  assert.equal((await worker.fetch(new Request('https://blog.example/uploads/a.png'), env, {})).status, 503);
  assert.equal(env.objects.size, 0);
});
