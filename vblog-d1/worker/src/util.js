// vBlog D1 · util.js —— 零 npm 依赖工具
// JWT(HS256) 与密码哈希(PBKDF2-SHA256) 全部基于 WebCrypto，可在 Worker 直接运行。

const enc = new TextEncoder();
const dec = new TextDecoder();

// ── base64url ─────────────────────────────────────────────────────
function b64url(bytes) {
  let s = '';
  for (const b of new Uint8Array(bytes)) s += String.fromCharCode(b);
  return btoa(s).replace(/\+/g, '-').replace(/\//g, '_').replace(/=+$/, '');
}

function b64urlToBytes(str) {
  let s = str.replace(/-/g, '+').replace(/_/g, '/');
  while (s.length % 4) s += '=';
  const bin = atob(s);
  const out = new Uint8Array(bin.length);
  for (let i = 0; i < bin.length; i++) out[i] = bin.charCodeAt(i);
  return out;
}

// ── JWT HS256 ─────────────────────────────────────────────────────
async function hmacSha256(secret, data) {
  const key = await crypto.subtle.importKey(
    'raw', enc.encode(secret), { name: 'HMAC', hash: 'SHA-256' }, false, ['sign']
  );
  return new Uint8Array(await crypto.subtle.sign('HMAC', key, data));
}

export async function signJWT(secret, claims, expiresSecs) {
  const now = Math.floor(Date.now() / 1000);
  const payload = Object.assign({}, claims, { iat: now, exp: now + expiresSecs });
  const h = b64url(enc.encode(JSON.stringify({ alg: 'HS256', typ: 'JWT' })));
  const p = b64url(enc.encode(JSON.stringify(payload)));
  const sig = b64url(await hmacSha256(secret, enc.encode(h + '.' + p)));
  return h + '.' + p + '.' + sig;
}

export async function verifyJWT(secret, token) {
  const parts = String(token || '').split('.');
  if (parts.length !== 3) return null;
  const [h, p, sig] = parts;
  const expected = b64url(await hmacSha256(secret, enc.encode(h + '.' + p)));
  const a = b64urlToBytes(sig);
  const b = b64urlToBytes(expected);
  if (a.length !== b.length) return null;
  let diff = 0;
  for (let i = 0; i < a.length; i++) diff |= a[i] ^ b[i];
  if (diff !== 0) return null;
  try {
    const claims = JSON.parse(dec.decode(b64urlToBytes(p)));
    const now = Math.floor(Date.now() / 1000);
    if (!claims || !claims.exp || claims.exp < now) return null;
    return claims;
  } catch {
    return null;
  }
}

// ── 密码哈希：pbkdf2$<iter>$<saltB64>$<hashB64> ───────────────────
const PBKDF2_ITER = 100000;

async function derive(password, salt, iter, len) {
  const base = await crypto.subtle.importKey('raw', enc.encode(password), 'PBKDF2', false, ['deriveBits']);
  const bits = await crypto.subtle.deriveBits(
    { name: 'PBKDF2', hash: 'SHA-256', salt, iterations: iter }, base, len * 8
  );
  return new Uint8Array(bits);
}

export async function hashPassword(password) {
  const salt = crypto.getRandomValues(new Uint8Array(16));
  const key = await derive(password, salt, PBKDF2_ITER, 32);
  return `pbkdf2$${PBKDF2_ITER}$${b64url(salt)}$${b64url(key)}`;
}

export async function checkPassword(hash, password) {
  const parts = String(hash || '').split('$');
  if (parts.length !== 4 || parts[0] !== 'pbkdf2') return false;
  const iter = parseInt(parts[1], 10) || PBKDF2_ITER;
  const salt = b64urlToBytes(parts[2]);
  const expected = b64urlToBytes(parts[3]);
  const got = await derive(password, salt, iter, expected.length);
  if (expected.length !== got.length) return false;
  let diff = 0;
  for (let i = 0; i < expected.length; i++) diff |= expected[i] ^ got[i];
  return diff === 0;
}

// ── HTTP 响应工具 ─────────────────────────────────────────────────
// opts 可传 ResponseInit 扩展字段，如 { cf: { cacheTtl: 60 } } 做边缘缓存
export function json(data, status = 200, opts = {}) {
  return new Response(JSON.stringify(data), {
    status,
    headers: { 'Content-Type': 'application/json; charset=utf-8' },
    ...opts,
  });
}

export function fail(msg, status = 400) {
  return json({ error: msg }, status);
}

// ── 日期/文本工具（对齐 Go 端契约：created_at=日期，updated_at=日期时间）─
const pad = (n) => String(n).padStart(2, '0');

export function fmtDate(iso) {
  if (!iso) return '';
  const d = new Date(iso);
  if (isNaN(d.getTime())) return '';
  return `${d.getFullYear()}-${pad(d.getMonth() + 1)}-${pad(d.getDate())}`;
}

export function fmtDateTime(iso) {
  if (!iso) return '';
  const d = new Date(iso);
  if (isNaN(d.getTime())) return '';
  return `${d.getFullYear()}-${pad(d.getMonth() + 1)}-${pad(d.getDate())} ${pad(d.getHours())}:${pad(d.getMinutes())}`;
}

// 阅读时间：~500 字符/分钟（对齐 Go 端 CalcReadTime）
export function calcReadTime(content) {
  const minutes = Math.ceil((content || '').length / 500);
  return minutes < 1 ? 1 : minutes;
}

// 摘要：截断到 maxLen 字符并加省略号（对齐 Go 端 BuildExcerpt）
export function buildExcerpt(content, maxLen = 200) {
  const s = content || '';
  if (s.length <= maxLen) return s;
  return s.slice(0, maxLen) + '...';
}

export function nowISO() {
  return new Date().toISOString();
}