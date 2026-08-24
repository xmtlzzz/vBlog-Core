// vBlog D1 —— Cloudflare Worker：博客 + 后台 REST API
// 复刻「vBlog Core」Go 后端（go-restful + GORM）的对外 JSON 契约，
// 前端 web/ 无需任何改动。
//
// 组成：
//   · API 路由（/api/*）→ 本 Worker + D1（SQLite）
//   · 静态资源（SPA）   → Worker Static Assets（env.ASSETS）
//   · 图片上传          → R2 对象存储
//   · 每日统计快照      → Cron Trigger（scheduled）
// 零 npm 依赖（JWT/PBKDF2 用 WebCrypto，见 util.js）。
//
// 部署见 ../README.md；wrangler.toml 中需配置 D1/R2 绑定与 JWT_SECRET。

import {
  json, fail, signJWT, verifyJWT, hashPassword, checkPassword,
  fmtDate, fmtDateTime, calcReadTime, buildExcerpt, nowISO,
} from './util.js';

export default {
  async fetch(request, env, ctx) {
    const url = new URL(request.url);
    const path = url.pathname;
    const method = request.method.toUpperCase();

    if (method === 'OPTIONS') {
      return new Response(null, { status: 204, headers: { 'Access-Control-Allow-Origin': '*', 'Access-Control-Allow-Methods': 'GET,POST,PUT,PATCH,DELETE,OPTIONS', 'Access-Control-Allow-Headers': 'Content-Type,Authorization' } });
    }

    // 静态资源走 ASSETS（SPA fallback 由 not_found_handling 处理）
    if (path.startsWith('/uploads/') && method === 'GET') {
      return serveImage(env, path);
    }
    if (!path.startsWith('/api')) {
      // 仅记录页面导航（无扩展名或 .html），静态资源不计 PV，避免统计虚高与 D1 写配额浪费
      const isPage = method === 'GET' && (!/\.[\w]+$/.test(path) || path.endsWith('.html'));
      if (isPage) ctx.waitUntil(recordPageView(env, request, path));
      return env.ASSETS.fetch(request);
    }

    // API 请求不记 page_view（减少 D1 写入压力）
    return route(request, env, url, path, method);
  },

  // 每日 00:05(UTC) 快照统计（wrangler.toml [triggers]）
  async scheduled(_event, env, _ctx) {
    await dailySnapshot(env);
  },
};

// ══════════════════════════ 路由分发 ══════════════════════════
async function route(request, env, url, path, method) {
  // 路径形如 /api/settings → ['api','settings']，跳过首段 'api'
  const seg = path.split('/').filter(Boolean);
  const [, first, second, third, fourth] = seg;

  try {
    // ── 认证 ────────────────────────────────────────────────
    if (first === 'auth') {
      if (method === 'POST' && second === 'login') return login(request, env);
      if (method === 'POST' && second === 'register') return register(request, env);
    }

    // ── 站点设置 ─────────────────────────────────────────────
    if (first === 'settings') {
      if (method === 'GET') return listSettings(env);
      if (method === 'PUT') return requireAuth(request, env, () => saveSettings(request, env));
      if (method === 'POST' && second === 'reset') return requireAuth(request, env, () => resetSettings(env));
    }

    // ── 仪表盘统计 ───────────────────────────────────────────
    if (first === 'dashboard' && second === 'stats' && method === 'GET') {
      return dashboardStats(env);
    }

    // ── 标签 ─────────────────────────────────────────────────
    if (first === 'tags') {
      if (!second) {
        if (method === 'GET') return listTags(env);
        if (method === 'POST') return requireAuth(request, env, () => createTag(request, env));
      } else if (/^\d+$/.test(second)) {
        if (method === 'PUT') return requireAuth(request, env, () => updateTag(request, env, second));
        if (method === 'DELETE') return requireAuth(request, env, () => deleteTag(env, second));
      }
    }

    // ── 文章 ─────────────────────────────────────────────────
    if (first === 'posts') {
      // /api/posts/trash（必须先于 {id} 匹配）
      if (second === 'trash' && method === 'GET') {
        return requireAuth(request, env, () => listTrash(env));
      }
      if (second === 'trash' && method === 'POST' && third === 'empty') {
        return requireAuth(request, env, () => emptyTrash(env));
      }
      // /api/posts/{postId}/comments
      if (third === 'comments') {
        if (method === 'GET') return listCommentsByPost(env, second);
        if (method === 'POST') return createPublicComment(request, env, second);
      }
      if (!second) {
        if (method === 'GET') return listPosts(env, url);
        if (method === 'POST') return requireAuth(request, env, (claims) => createPost(request, env, claims));
      } else if (/^\d+$/.test(second)) {
        if (third === 'restore' && method === 'POST') return requireAuth(request, env, () => restorePost(env, second));
        if (third === 'permanent' && method === 'DELETE') return requireAuth(request, env, () => permanentDeletePost(env, second));
        if (!third) {
          if (method === 'GET') return getPost(env, second);
          if (method === 'PUT') return requireAuth(request, env, () => updatePost(request, env, second));
          if (method === 'DELETE') return requireAuth(request, env, () => deletePost(env, second));
        }
      }
    }

    // ── 评论（管理端）────────────────────────────────────────
    if (first === 'comments') {
      if (!second) {
        if (method === 'GET') return requireAuth(request, env, () => listComments(env, url));
        if (method === 'POST') return requireAuth(request, env, () => createComment(request, env));
      } else if (/^\d+$/.test(second)) {
        if (third === 'approve' && method === 'PATCH') return requireAuth(request, env, () => setCommentStatus(env, second, 'approved'));
        if (third === 'spam' && method === 'PATCH') return requireAuth(request, env, () => setCommentStatus(env, second, 'spam'));
        if (!third && method === 'DELETE') return requireAuth(request, env, () => deleteComment(env, second));
      }
    }

    // ── 组件 ─────────────────────────────────────────────────
    if (first === 'components') {
      if (second === 'active' && method === 'GET') return listComponents(env, 'active', { cf: { cacheTtl: 60 } });
      if (!second) {
        if (method === 'GET') return requireAuth(request, env, () => listComponents(env, ''));
        if (method === 'POST') return requireAuth(request, env, () => createComponent(request, env));
      } else if (/^\d+$/.test(second)) {
        if (third === 'toggle' && method === 'PATCH') return requireAuth(request, env, () => toggleComponent(env, second));
        if (!third) {
          if (method === 'PUT') return requireAuth(request, env, () => updateComponent(request, env, second));
          if (method === 'DELETE') return requireAuth(request, env, () => deleteComponent(env, second));
        }
      }
    }

    // ── 上传（→ R2）──────────────────────────────────────────
    if (first === 'upload' && method === 'POST') {
      return requireAuth(request, env, () => uploadImage(request, env));
    }

    // ── RSS ──────────────────────────────────────────────────
    if (first === 'rss' && method === 'GET') return rssFeed(env);

    return fail('not found', 404);
  } catch (e) {
    console.error('route error:', path, e);
    return fail('internal server error', 500);
  }
}

// ══════════════════════════ 鉴权 ══════════════════════════
// 基于 Cache API 的按 IP 固定窗口限流（Worker 无状态、零依赖的折衷方案；
// Cache 是最终一致的全局副本，极端并发下限流略松，但足以挡住脚本爆破）
async function rateLimit(request, windowSecs, limit) {
  const ip = request.headers.get('CF-Connecting-IP') || 'unknown';
  const key = new Request('https://ratelimit.local/' + ip, { method: 'GET' });
  const cache = caches.default;
  let count = 0;
  const hit = await cache.match(key);
  if (hit) {
    count = parseInt(hit.headers.get('X-Count') || '0', 10);
    const resetAt = parseInt(hit.headers.get('X-Reset') || '0', 10);
    if (Date.now() > resetAt) count = 0;
  }
  if (count >= limit) return false;
  const res = new Response(null, {
    headers: {
      'X-Count': String(count + 1),
      'X-Reset': String(Date.now() + windowSecs * 1000),
      'Cache-Control': `public, max-age=${windowSecs}`,
    },
  });
  await cache.put(key, res);
  return true;
}

async function requireAuth(request, env, handler) {
  const auth = request.headers.get('Authorization') || '';
  const token = auth.startsWith('Bearer ') ? auth.slice(7) : '';
  if (!token) return fail('missing token', 401);
  const claims = await verifyJWT(env.JWT_SECRET || '', token);
  if (!claims) return fail('invalid token', 401);
  return handler(claims);
}

async function issueTokens(env, user) {
  const access = await signJWT(env.JWT_SECRET || '', { user_id: user.id, username: user.username }, 24 * 3600);
  const refresh = await signJWT(env.JWT_SECRET || '', { user_id: user.id, username: user.username }, 7 * 24 * 3600);
  return json({ access_token: access, refresh_token: refresh });
}

// ══════════════════════════ 认证 ══════════════════════════
async function login(request, env) {
  if (!(await rateLimit(request, 60, 10))) return fail('请求过于频繁，请稍后再试', 429);
  const body = await readJson(request);
  if (!body) return fail('invalid request body');
  const user = await env.DB.prepare('SELECT * FROM users WHERE username = ?').bind(String(body.username || '')).first();
  if (!user || !(await checkPassword(user.password, String(body.password || '')))) {
    return fail('invalid credentials', 401);
  }
  return issueTokens(env, user);
}

async function register(request, env) {
  if (!(await rateLimit(request, 60, 5))) return fail('请求过于频繁，请稍后再试', 429);
  // 首号闸门：已存在用户则关闭公开注册（首个注册者即管理员）
  const anyUser = await env.DB.prepare('SELECT id FROM users LIMIT 1').first();
  if (anyUser) return fail('注册已关闭：管理员已存在', 403);
  const body = await readJson(request);
  const username = String((body && body.username) || '');
  const password = String((body && body.password) || '');
  const email = String((body && body.email) || '');
  if (!username || !password) return fail('注册失败，用户名可能已存在');
  const exists = await env.DB.prepare('SELECT id FROM users WHERE username = ?').bind(username).first();
  if (exists) return fail('注册失败，用户名可能已存在');
  const hash = await hashPassword(password);
  const res = await env.DB.prepare(
    'INSERT INTO users (username, password, email, role, created_at, updated_at) VALUES (?, ?, ?, ?, ?, ?)'
  ).bind(username, hash, email, 'admin', nowISO(), nowISO()).run();
  const id = res.meta.last_row_id;
  return issueTokens(env, { id, username });
}

// ══════════════════════════ 设置 ══════════════════════════
const DEFAULT_SETTINGS = {
  site_title: 'vBlog',
  site_description: 'A lightweight blog for geeks',
  site_url: '',
  posts_per_page: '10',
  theme: 'default',
  footer_text: 'Powered by vBlog Core',
  grpc_api_key: '',
  grpc_port: '50051',
};

async function listSettings(env) {
  const rows = await env.DB.prepare('SELECT key, value FROM settings').all();
  const out = {};
  for (const r of rows.results) out[r.key] = r.value;
  // 公开端点不泄露机密项（gRPC 连接密钥）；管理端保存时留空表示不变更
  delete out.grpc_api_key;
  return json(out, 200, { cf: { cacheTtl: 60 } });
}

async function saveSettings(request, env) {
  const body = await readJson(request);
  if (!body || typeof body !== 'object') return fail('invalid request body');
  const entries = Object.entries(body).filter(([k, v]) => !(k === 'grpc_api_key' && !v));
  if (entries.length === 0) return fail('invalid request body');
  // 全部 UPSERT 合并为一次 batch，原子写入
  const stmt = env.DB.prepare(
    `INSERT INTO settings (key, value, updated_at) VALUES (?, ?, ?)
     ON CONFLICT(key) DO UPDATE SET value = excluded.value, updated_at = excluded.updated_at`
  );
  await env.DB.batch(entries.map(([k, v]) => stmt.bind(k, String(v ?? ''), nowISO())));
  return json({ message: 'ok' });
}

async function resetSettings(env) {
  // 与 Go 版语义一致：重置默认值但保留现有 gRPC API Key（默认表中该键为空串，跳过即保留）
  const entries = Object.entries(DEFAULT_SETTINGS).filter(([k]) => k !== 'grpc_api_key');
  await env.DB.batch(entries.map(([k, v]) =>
    env.DB.prepare(
      `INSERT INTO settings (key, value, updated_at) VALUES (?, ?, ?)
       ON CONFLICT(key) DO UPDATE SET value = excluded.value, updated_at = excluded.updated_at`
    ).bind(k, v, nowISO())
  ));
  return json({ message: 'ok' });
}

// ══════════════════════════ 统计 ══════════════════════════
async function dashboardStats(env) {
  const stats = await env.DB.prepare(`
    SELECT
      (SELECT COUNT(*) FROM posts WHERE status = 'published' AND deleted_at IS NULL) AS total_posts,
      (SELECT COALESCE(SUM(views), 0) FROM posts WHERE deleted_at IS NULL) AS total_views,
      (SELECT COUNT(*) FROM comments) AS total_comments,
      (SELECT COUNT(*) FROM tags) AS total_tags
  `).first();
  return json({
    total_posts: stats.total_posts || 0,
    total_views: stats.total_views || 0,
    total_comments: stats.total_comments || 0,
    total_tags: stats.total_tags || 0,
  }, 200, { cf: { cacheTtl: 60 } });
}

// ══════════════════════════ 标签 ══════════════════════════
async function listTags(env) {
  const rows = await env.DB.prepare('SELECT id, name, description, created_at FROM tags ORDER BY created_at DESC').all();
  return json({ data: rows.results }, 200, { cf: { cacheTtl: 60 } });
}

async function createTag(request, env) {
  const body = await readJson(request);
  if (!body || !body.name) return fail('invalid request body');
  try {
    const res = await env.DB.prepare('INSERT INTO tags (name, description, created_at) VALUES (?, ?, ?)')
      .bind(String(body.name), String(body.description || ''), nowISO()).run();
    const tag = await env.DB.prepare('SELECT id, name, description, created_at FROM tags WHERE id = ?').bind(res.meta.last_row_id).first();
    return json(tag, 201);
  } catch {
    return fail('tag already exists', 400);
  }
}

async function updateTag(request, env, id) {
  const body = await readJson(request);
  if (!body) return fail('invalid request body');
  const name = body.name !== undefined ? String(body.name) : undefined;
  const desc = body.description !== undefined ? String(body.description) : undefined;
  const sets = [];
  const binds = [];
  if (name !== undefined) { sets.push('name = ?'); binds.push(name); }
  if (desc !== undefined) { sets.push('description = ?'); binds.push(desc); }
  if (sets.length === 0) return fail('invalid request body');
  binds.push(id);
  await env.DB.prepare(`UPDATE tags SET ${sets.join(', ')} WHERE id = ?`).bind(...binds).run();
  const tag = await env.DB.prepare('SELECT id, name, description, created_at FROM tags WHERE id = ?').bind(id).first();
  return tag ? json(tag) : fail('not found', 404);
}

async function deleteTag(env, id) {
  await env.DB.prepare('DELETE FROM post_tags WHERE tag_id = ?').bind(id).run();
  await env.DB.prepare('DELETE FROM tags WHERE id = ?').bind(id).run();
  return json({ message: 'ok' });
}

// ══════════════════════════ 文章 ══════════════════════════
// 标签子查询：一条 SQL 聚合出 JSON，消除「每篇文章一次查询」的 N+1 延迟
const TAGS_JSON_SQL = `(
  SELECT json_group_array(json_object(
    'id', t.id, 'name', t.name, 'description', t.description, 'created_at', t.created_at)
    ORDER BY pt.rowid)
  FROM post_tags pt JOIN tags t ON t.id = pt.tag_id
  WHERE pt.post_id = p.id
) AS tags_json`;

function parseTagsJson(tagsJson) {
  if (!tagsJson) return [];
  try {
    const arr = JSON.parse(tagsJson);
    return Array.isArray(arr) ? arr : [];
  } catch {
    return [];
  }
}

function toPostResp(row) {
  return {
    id: row.id, title: row.title, content: row.content, excerpt: row.excerpt,
    status: row.status, pinned: !!row.pinned, views: row.views, read_time: row.read_time,
    author_id: row.author_id, tags: parseTagsJson(row.tags_json),
    created_at: fmtDate(row.created_at), updated_at: fmtDateTime(row.updated_at),
  };
}

async function postById(env, id) {
  const row = await env.DB.prepare(
    `SELECT p.*, ${TAGS_JSON_SQL} FROM posts p WHERE p.id = ? AND p.deleted_at IS NULL`
  ).bind(id).first();
  if (!row) return null;
  return toPostResp(row);
}

const POST_SELECT = `SELECT p.*, ${TAGS_JSON_SQL} FROM posts p`;

async function listPosts(env, url) {
  const page = Math.max(1, parseInt(url.searchParams.get('page') || '1', 10) || 1);
  const perPage = Math.max(1, parseInt(url.searchParams.get('per_page') || '5', 10) || 5);
  const tag = url.searchParams.get('tag');
  const status = url.searchParams.get('status');
  const search = url.searchParams.get('search');

  const where = ['p.deleted_at IS NULL'];
  const binds = [];
  // 注意：SQL 里 JOIN 的占位符先于 WHERE，故 tag 参数必须先入 binds
  let join = '';
  if (tag) {
    join = 'JOIN post_tags ptf ON ptf.post_id = p.id JOIN tags tf ON tf.id = ptf.tag_id AND tf.name = ?';
    binds.push(tag);
  }
  if (status) { where.push('p.status = ?'); binds.push(status); }
  if (search) { where.push('p.title LIKE ?'); binds.push('%' + search + '%'); }
  const whereSql = 'WHERE ' + where.join(' AND ');

  const total = await env.DB.prepare(`SELECT COUNT(DISTINCT p.id) AS n FROM posts p ${join} ${whereSql}`).bind(...binds).first();
  const rows = await env.DB.prepare(
    `${POST_SELECT} ${join} ${whereSql} ORDER BY p.pinned DESC, p.created_at DESC LIMIT ? OFFSET ?`
  ).bind(...binds, perPage, (page - 1) * perPage).all();

  const data = rows.results.map(toPostResp);
  return json({ data, total: total.n || 0, page }, 200, { cf: { cacheTtl: 30 } });
}

async function getPost(env, id) {
  const row = await env.DB.prepare(
    `${POST_SELECT} WHERE p.id = ? AND p.deleted_at IS NULL`
  ).bind(id).first();
  if (!row) return fail('not found', 404);
  // 阅读量 +1（对齐 Go 端 GetByID；不缓存详情保证计数准确）
  await env.DB.prepare('UPDATE posts SET views = views + 1 WHERE id = ?').bind(id).run();
  row.views = (row.views || 0) + 1;
  return json(toPostResp(row));
}

async function resolveTags(env, tags) {
  const names = [];
  for (const t of tags || []) {
    if (t && t.name) names.push(String(t.name));
  }
  if (names.length === 0) return [];
  const out = [];
  for (const name of names) {
    let tag = await env.DB.prepare('SELECT id, name, description, created_at FROM tags WHERE name = ?').bind(name).first();
    if (!tag) {
      try {
        const res = await env.DB.prepare('INSERT INTO tags (name, description, created_at) VALUES (?, ?, ?)')
          .bind(name, '', nowISO()).run();
        tag = { id: res.meta.last_row_id, name, description: '', created_at: nowISO() };
      } catch {
        tag = await env.DB.prepare('SELECT id, name, description, created_at FROM tags WHERE name = ?').bind(name).first();
      }
    }
    out.push(tag);
  }
  return out;
}

// D1 batch 原子执行：删除旧关联 + 写入新关联，中途失败整体回滚
async function replacePostTags(env, postId, tags) {
  const del = env.DB.prepare('DELETE FROM post_tags WHERE post_id = ?').bind(postId);
  const ins = tags.map((t) =>
    env.DB.prepare('INSERT OR IGNORE INTO post_tags (post_id, tag_id) VALUES (?, ?)').bind(postId, t.id)
  );
  await env.DB.batch([del, ...ins]);
}

async function createPost(request, env, claims) {
  const body = await readJson(request);
  if (!body) return fail('invalid request body');
  const title = String(body.title || '').trim();
  if (!title) return fail('title required');
  const tags = await resolveTags(env, body.tags);
  const pinned = body.pinned ? 1 : 0;
  const res = await env.DB.prepare(
    `INSERT INTO posts (title, content, excerpt, status, pinned, views, read_time, author_id, created_at, updated_at)
     VALUES (?, ?, ?, ?, ?, 0, ?, ?, ?, ?)`
  ).bind(
    title, String(body.content || ''), String(body.excerpt || buildExcerpt(body.content)),
    String(body.status || 'draft'), pinned, calcReadTime(body.content), (claims && claims.user_id) || 0, nowISO(), nowISO()
  ).run();
  const id = res.meta.last_row_id;
  // 两步执行：INSERT 需先拿 last_row_id 作标签关联的外键，无法并入同一 batch；
  // 第二步失败会留下无标签文章（可重新编辑补标签），属已知取舍
  await replacePostTags(env, id, tags);
  return json(await postById(env, id), 201);
}

async function updatePost(request, env, id) {
  const body = await readJson(request);
  if (!body) return fail('invalid request body');
  const exists = await env.DB.prepare('SELECT id FROM posts WHERE id = ? AND deleted_at IS NULL').bind(id).first();
  if (!exists) return fail('not found', 404);
  const tags = await resolveTags(env, body.tags);
  const content = String(body.content || '');
  const pinned = body.pinned ? 1 : 0;
  // UPDATE 与标签关联替换合并为一次 batch，保证原子性
  await env.DB.batch([
    env.DB.prepare(
      `UPDATE posts SET title = ?, content = ?, excerpt = ?, status = ?, pinned = ?, read_time = ?, updated_at = ?
       WHERE id = ?`
    ).bind(
      String(body.title || '').trim(), content,
      String(body.excerpt || buildExcerpt(content)),
      String(body.status || 'draft'), pinned, calcReadTime(content), nowISO(), id
    ),
    env.DB.prepare('DELETE FROM post_tags WHERE post_id = ?').bind(id),
    ...tags.map((t) =>
      env.DB.prepare('INSERT OR IGNORE INTO post_tags (post_id, tag_id) VALUES (?, ?)').bind(id, t.id)
    ),
  ]);
  return json(await postById(env, id));
}

async function deletePost(env, id) {
  const res = await env.DB.prepare('UPDATE posts SET deleted_at = ? WHERE id = ? AND deleted_at IS NULL').bind(nowISO(), id).run();
  if (res.meta.changes === 0) return fail('not found', 404);
  return json({ message: 'ok' });
}

async function listTrash(env) {
  const rows = await env.DB.prepare(
    `${POST_SELECT} WHERE p.deleted_at IS NOT NULL ORDER BY p.deleted_at DESC`
  ).all();
  return json({ data: rows.results.map(toPostResp) });
}

// 逐篇串行删除（permanentDeletePost 内部多表清理，量小可接受）
async function emptyTrash(env) {
  const rows = await env.DB.prepare('SELECT id FROM posts WHERE deleted_at IS NOT NULL').all();
  for (const row of rows.results) {
    await permanentDeletePost(env, String(row.id));
  }
  return json({ message: 'ok' });
}

async function restorePost(env, id) {
  const res = await env.DB.prepare('UPDATE posts SET deleted_at = NULL WHERE id = ? AND deleted_at IS NOT NULL').bind(id).run();
  if (res.meta.changes === 0) return fail('not found', 404);
  return json({ message: 'ok' });
}

async function permanentDeletePost(env, id) {
  await env.DB.prepare('DELETE FROM post_tags WHERE post_id = ?').bind(id).run();
  await env.DB.prepare('DELETE FROM comments WHERE post_id = ?').bind(id).run();
  await env.DB.prepare('DELETE FROM posts WHERE id = ?').bind(id).run();
  return json({ message: 'ok' });
}

// ══════════════════════════ 评论 ══════════════════════════
async function listComments(env, url) {
  const page = Math.max(1, parseInt(url.searchParams.get('page') || '1', 10) || 1);
  const perPage = Math.max(1, parseInt(url.searchParams.get('per_page') || '10', 10) || 10);
  const status = url.searchParams.get('status');
  const search = url.searchParams.get('search');
  const where = [];
  const binds = [];
  if (status) { where.push('status = ?'); binds.push(status); }
  if (search) { where.push('body LIKE ?'); binds.push('%' + search + '%'); }
  const whereSql = where.length ? 'WHERE ' + where.join(' AND ') : '';
  const total = await env.DB.prepare(`SELECT COUNT(*) AS n FROM comments ${whereSql}`).bind(...binds).first();
  const rows = await env.DB.prepare(
    `SELECT * FROM comments ${whereSql} ORDER BY created_at DESC LIMIT ? OFFSET ?`
  ).bind(...binds, perPage, (page - 1) * perPage).all();
  return json({ data: rows.results, total: total.n || 0, page });
}

async function createComment(request, env) {
  const body = await readJson(request);
  if (!body || !body.post_id || !body.body) return fail('invalid request body');
  const res = await env.DB.prepare(
    `INSERT INTO comments (post_id, author_name, author_email, body, status, created_at)
     VALUES (?, ?, ?, ?, 'pending', ?)`
  ).bind(body.post_id, String(body.author_name || ''), String(body.author_email || ''), String(body.body), nowISO()).run();
  const comment = await env.DB.prepare('SELECT * FROM comments WHERE id = ?').bind(res.meta.last_row_id).first();
  return json(comment, 201);
}

async function listCommentsByPost(env, postId) {
  const rows = await env.DB.prepare(
    `SELECT * FROM comments WHERE post_id = ? AND status = 'approved' ORDER BY created_at DESC`
  ).bind(postId).all();
  return json({ data: rows.results, total: rows.results.length, page: 1 });
}

async function createPublicComment(request, env, postId) {
  const body = await readJson(request);
  if (!body || !body.body) return fail('invalid request body');
  // post_id 必须是已存在的未删除文章，拒绝脏数据
  if (!/^\d+$/.test(String(postId))) return fail('invalid post id', 400);
  const post = await env.DB.prepare('SELECT id FROM posts WHERE id = ? AND deleted_at IS NULL').bind(postId).first();
  if (!post) return fail('post not found', 404);
  // 设置里 enable_comments = 'false' 时禁止评论
  const setting = await env.DB.prepare("SELECT value FROM settings WHERE key = 'enable_comments'").first();
  if (setting && setting.value === 'false') return fail('comments disabled', 400);
  await env.DB.prepare(
    `INSERT INTO comments (post_id, author_name, author_email, body, status, created_at)
     VALUES (?, ?, ?, ?, 'pending', ?)`
  ).bind(postId, String(body.author_name || '').slice(0, 100), String(body.author_email || '').slice(0, 255), String(body.body), nowISO()).run();
  return json({ message: '评论已提交，等待审核' }, 201);
}

async function setCommentStatus(env, id, status) {
  const res = await env.DB.prepare('UPDATE comments SET status = ? WHERE id = ?').bind(status, id).run();
  if (res.meta.changes === 0) return fail('not found', 404);
  return json({ message: 'ok' });
}

async function deleteComment(env, id) {
  await env.DB.prepare('DELETE FROM comments WHERE id = ?').bind(id).run();
  return json({ message: 'ok' });
}

// ══════════════════════════ 组件 ══════════════════════════
async function listComponents(env, status, opts = {}) {
  const sql = status
    ? 'SELECT * FROM components WHERE status = ? ORDER BY id ASC'
    : 'SELECT * FROM components ORDER BY id ASC';
  const rows = status
    ? await env.DB.prepare(sql).bind(status).all()
    : await env.DB.prepare(sql).all();
  return json({ data: rows.results }, 200, opts);
}

async function createComponent(request, env) {
  const body = await readJson(request);
  if (!body || !body.name) return fail('invalid request body');
  const res = await env.DB.prepare(
    `INSERT INTO components (name, description, version, code, category, origin, status, created_at, updated_at)
     VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)`
  ).bind(
    String(body.name), String(body.description || ''), String(body.version || ''),
    String(body.code || ''), String(body.category || ''), String(body.origin || 'built-in'),
    String(body.status || 'active'), nowISO(), nowISO()
  ).run();
  const comp = await env.DB.prepare('SELECT * FROM components WHERE id = ?').bind(res.meta.last_row_id).first();
  return json(comp, 201);
}

async function updateComponent(request, env, id) {
  const body = await readJson(request);
  if (!body) return fail('invalid request body');
  const fields = ['name', 'description', 'version', 'code', 'category', 'origin', 'status'];
  const sets = [];
  const binds = [];
  for (const f of fields) {
    if (body[f] !== undefined) { sets.push(`${f} = ?`); binds.push(String(body[f])); }
  }
  if (sets.length === 0) return fail('invalid request body');
  sets.push('updated_at = ?');
  binds.push(nowISO(), id);
  const res = await env.DB.prepare(`UPDATE components SET ${sets.join(', ')} WHERE id = ?`).bind(...binds).run();
  if (res.meta.changes === 0) return fail('not found', 404);
  const comp = await env.DB.prepare('SELECT * FROM components WHERE id = ?').bind(id).first();
  return json(comp);
}

async function deleteComponent(env, id) {
  await env.DB.prepare('DELETE FROM components WHERE id = ?').bind(id).run();
  return json({ message: 'ok' });
}

async function toggleComponent(env, id) {
  const comp = await env.DB.prepare('SELECT id, status FROM components WHERE id = ?').bind(id).first();
  if (!comp) return fail('not found', 404);
  const next = comp.status === 'active' ? 'inactive' : 'active';
  await env.DB.prepare('UPDATE components SET status = ?, updated_at = ? WHERE id = ?').bind(next, nowISO(), id).run();
  return json({ message: 'ok' });
}

// ══════════════════════════ 上传 → R2 / KV 回退 ══════════════════════════
// 图片存储策略：
// · 配置了 R2（R2_PUBLIC_URL 非空）→ 上传到 R2 桶，返回 CDN URL
// · 未配置 R2 但有 KV 命名空间（IMG 绑定）→ 存 KV，返回本站 /uploads/<name>（Worker 读取）
const UPLOAD_MAX_SIZE = 10 * 1024 * 1024; // 10MB，对齐 Go 版
const ALLOWED_IMAGE_TYPES = new Set(['image/png', 'image/jpeg', 'image/gif', 'image/webp']);

async function uploadImage(request, env) {
  const form = await request.formData().catch(() => null);
  const file = form && form.get('file');
  // 类型白名单（客户端 file.type，收窄到位图四类；svg 一律拒绝）；大小上限 10MB 对齐 Go 版
  if (!file || !ALLOWED_IMAGE_TYPES.has(file.type)) return fail('only png/jpeg/gif/webp images allowed', 400);
  if (file.size > UPLOAD_MAX_SIZE) return fail('文件超过 10MB 限制', 400);

  const ext = (file.name || '').includes('.') ? '.' + file.name.split('.').pop().toLowerCase() : '.png';
  const key = `uploads/${Date.now()}${Math.floor(Math.random() * 1e6)}${ext}`; // 对齐 Go 端纳秒时间戳命名
  const buf = await file.arrayBuffer();
  const contentType = file.type;

  // ── R2 路径（首选）─────────────────────────────────────────────
  if (env.R2 && env.R2_PUBLIC_URL) {
    await env.R2.put(key, buf, { httpMetadata: { contentType } });
    return json({ url: String(env.R2_PUBLIC_URL).replace(/\/+$/, '') + '/' + key });
  }

  // ── KV 回退路径（无 R2 / 未配置 R2_PUBLIC_URL）─────────────────
  if (!env.IMG) return fail('upload storage not configured', 500);
  await env.IMG.put(key, buf, { metadata: { contentType } });
  const origin = new URL(request.url).origin;
  return json({ url: origin + '/' + key });
}

// 从 KV 读取图片（无 R2 时的回退图床）
const IMAGE_MIME = {
  png: 'image/png', jpg: 'image/jpeg', jpeg: 'image/jpeg', gif: 'image/gif',
  webp: 'image/webp', avif: 'image/avif', bmp: 'image/bmp',
};

async function serveImage(env, path) {
  const key = path.replace(/^\//, '');
  if (!env.IMG) return fail('not found', 404);
  const obj = await env.IMG.get(key, { type: 'arrayBuffer', cacheTtl: 86400 });
  if (!obj) return fail('not found', 404);
  // 兼容两种 KV 返回形态：有的运行时返回 {value,...} 包装，有的直接返回值
  const buf = obj.value !== undefined ? obj.value : obj;
  const ext = key.split('.').pop().toLowerCase();
  return new Response(buf, {
    headers: {
      'Content-Type': IMAGE_MIME[ext] || 'application/octet-stream',
      'Cache-Control': 'public, max-age=31536000, immutable',
      'Content-Length': buf.byteLength,
    },
  });
}

// ══════════════════════════ RSS ══════════════════════════
async function rssFeed(env) {
  const setting = await env.DB.prepare("SELECT value FROM settings WHERE key = 'site_title'").first();
  const siteTitle = (setting && setting.value) || 'vBlog';
  const rows = await env.DB.prepare(
    `SELECT id, title, excerpt, created_at FROM posts WHERE status = 'published' AND deleted_at IS NULL
     ORDER BY created_at DESC LIMIT 20`
  ).all();

  const rfc1123 = (iso) => new Date(iso).toUTCString().replace('GMT', '+0000');
  const items = rows.results.map((p) =>
    `      <item>\n` +
    `        <title>${escapeXml(p.title)}</title>\n` +
    `        <link>/post/${p.id}</link>\n` +
    `        <description>${escapeXml(p.excerpt)}</description>\n` +
    `        <pubDate>${rfc1123(p.created_at)}</pubDate>\n` +
    `        <guid>post-${p.id}</guid>\n` +
    `      </item>`
  ).join('\n');

  const xml =
    `<?xml version="1.0" encoding="UTF-8"?>\n` +
    `<rss version="2.0">\n` +
    `  <channel>\n` +
    `    <title>${escapeXml(siteTitle)}</title>\n` +
    `    <link>/</link>\n` +
    `    <description>RSS Feed for ${escapeXml(siteTitle)}</description>\n` +
    `    <language>zh-CN</language>\n` +
    `    <lastBuildDate>${rfc1123(nowISO())}</lastBuildDate>\n` +
    items + '\n' +
    `  </channel>\n` +
    `</rss>\n`;

  return new Response(xml, {
    headers: { 'Content-Type': 'application/xml; charset=utf-8' },
    cf: { cacheTtl: 300 },
  });
}

const escapeXml = (s) => String(s || '')
  .replace(/&/g, '&amp;').replace(/</g, '&lt;').replace(/>/g, '&gt;')
  .replace(/"/g, '&quot;').replace(/'/g, '&apos;');

// ══════════════════════════ 统计落库 ══════════════════════════
async function recordPageView(env, request, path) {
  try {
    const ip = request.headers.get('CF-Connecting-IP') || '';
    const ua = (request.headers.get('User-Agent') || '').slice(0, 500);
    await env.DB.prepare('INSERT INTO page_views (ip, path, user_agent, created_at) VALUES (?, ?, ?, ?)')
      .bind(ip, path.slice(0, 500), ua, nowISO()).run();
  } catch (e) {
    // 统计失败不影响主流程
  }
}

async function dailySnapshot(env) {
  // 补偿机制：从「已有最新快照的次日」补到昨天，漏跑的日子自动回填（回看上限 90 天，与明细保留期一致）
  const yesterday = new Date(Date.now() - 86400000).toISOString().slice(0, 10); // 昨天(UTC)
  const earliest = new Date(Date.now() - 90 * 86400000).toISOString().slice(0, 10);
  const latest = await env.DB.prepare('SELECT MAX(stat_date) AS d FROM daily_stats').first();
  let day = (latest && latest.d && latest.d > earliest) ? latest.d : earliest;
  while (day < yesterday) {
    day = new Date(new Date(day + 'T00:00:00Z').getTime() + 86400000).toISOString().slice(0, 10);
    await snapshotDay(env, day); // 从缺口次日起逐日补齐
  }
  await snapshotDay(env, yesterday);
  // 清理 90 天前的明细，控制 D1 写入量
  await env.DB.prepare("DELETE FROM page_views WHERE date(created_at) < date('now', '-90 days')").run();
}

// 统计某一天的 PV/UV 并 UPSERT 进 daily_stats（其余计数为当前总量快照）
async function snapshotDay(env, day) {
  const pvRow = await env.DB.prepare('SELECT COUNT(*) AS n FROM page_views WHERE date(created_at) = ?').bind(day).first();
  const uvRow = await env.DB.prepare("SELECT COUNT(DISTINCT ip) AS n FROM page_views WHERE date(created_at) = ? AND ip != ''").bind(day).first();
  const s = await env.DB.prepare(`
    SELECT
      (SELECT COUNT(*) FROM posts WHERE status = 'published' AND deleted_at IS NULL) AS post_count,
      (SELECT COALESCE(SUM(views), 0) FROM posts WHERE deleted_at IS NULL) AS view_total,
      (SELECT COUNT(*) FROM comments) AS comment_count,
      (SELECT COUNT(*) FROM tags) AS tag_count
  `).first();
  await env.DB.prepare(
    `INSERT INTO daily_stats (stat_date, pv, uv, post_count, view_total, comment_count, tag_count, created_at)
     VALUES (?, ?, ?, ?, ?, ?, ?, ?)
     ON CONFLICT(stat_date) DO UPDATE SET
       pv = excluded.pv, uv = excluded.uv, post_count = excluded.post_count,
       view_total = excluded.view_total, comment_count = excluded.comment_count, tag_count = excluded.tag_count`
  ).bind(day, pvRow.n || 0, uvRow.n || 0, s.post_count || 0, s.view_total || 0, s.comment_count || 0, s.tag_count || 0, nowISO()).run();
}

// ══════════════════════════ 工具 ══════════════════════════
async function readJson(request) {
  try {
    const text = await request.text();
    if (!text) return null;
    return JSON.parse(text);
  } catch {
    return null;
  }
}