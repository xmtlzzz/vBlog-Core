import assert from 'node:assert/strict';
import { test } from 'node:test';
import worker from '../src/index.js';

function environment(failed) {
  const calls = [];
  return {
    calls,
    DB: { prepare(sql) {
      calls.push(sql);
      return { async first() {
        if (failed === 'd1') throw new Error('private database detail');
        return { '1': 1 };
      } };
    } },
    IMG: { async get(key) {
      calls.push(key);
      if (failed === 'kv') throw new Error('private storage detail');
      return null;
    } },
    ASSETS: { fetch() { throw new Error('Health must not enter SPA'); } },
  };
}

for (const failed of [undefined, 'd1', 'kv', 'missing-kv', 'missing-db']) {
  test(`health probes actual dependencies: ${failed || 'healthy'}`, async () => {
    const env = environment(failed);
    if (failed === 'missing-kv') delete env.IMG;
    if (failed === 'missing-db') delete env.DB;
    const response = await worker.fetch(new Request('https://test.example/api/health'), env, {
      waitUntil() { throw new Error('Health must not write page statistics'); },
    });
    assert.equal(response.status, failed ? 503 : 200);
    assert.equal(response.headers.get('Cache-Control'), 'no-store');
    assert.match(response.headers.get('Content-Type'), /application\/json/);
    const body = await response.json();
    assert.equal(body.status, failed ? 'DOWN' : 'UP');
    for (const name of ['d1', 'kv']) {
      assert.equal(body.components[name].status,
        failed === name || failed === `missing-${name === 'd1' ? 'db' : 'kv'}` ? 'DOWN' : 'UP');
      assert.ok(Number.isFinite(body.components[name].latency_ms));
    }
    if (env.DB) assert.ok(env.calls.includes('SELECT 1'));
    if (env.IMG) assert.ok(env.calls.includes('__vpulse_health__'));
    assert.ok(!JSON.stringify(body).includes('private'));
  });
}

test('health matches only the public GET endpoint', async () => {
  for (const [path, method] of [['/api/health/extra', 'GET'], ['/api/health', 'POST']]) {
    const env = environment();
    assert.equal((await worker.fetch(new Request(`https://test.example${path}`, { method }), env, {})).status, 404);
    assert.deepEqual(env.calls, []);
  }
});
