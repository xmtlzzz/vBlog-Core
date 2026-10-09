/**
 * Parse JWT payload without verifying signature (client-side decode).
 * @param {string} token
 * @returns {object|null}
 */
export function parseJwt(token) {
  if (!token || typeof token !== 'string') return null
  try {
    const parts = token.split('.')
    if (parts.length !== 3) return null
    const base64 = parts[1].replace(/-/g, '+').replace(/_/g, '/')
    const jsonStr = decodeURIComponent(
      atob(base64)
        .split('')
        .map(c => '%' + ('00' + c.charCodeAt(0).toString(16)).slice(-2))
        .join('')
    )
    return JSON.parse(jsonStr)
  } catch {
    return null
  }
}

/**
 * Checks if a JWT token is expired, invalid, or missing.
 * @param {string} token
 * @returns {boolean} true if token is missing, malformed, or expired
 */
export function isTokenExpired(token) {
  if (!token) return true
  const payload = parseJwt(token)
  if (!payload) return true
  if (!payload.exp) return false
  // 5 seconds buffer before true expiration
  return Date.now() >= (payload.exp * 1000 - 5000)
}
