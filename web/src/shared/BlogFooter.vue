<template>
  <footer class="blog-footer">
    <span>© 2026 vBlog Core · 用代码写作，用文字思考</span>
    <div class="footer-links">
      <a href="#">GitHub</a>
      <a href="#">RSS</a>
      <a href="#">gRPC API</a>
    </div>
  </footer>
</template>

<style scoped>
.blog-footer {
  max-width: 1080px;
  margin: 0 auto;
  padding: 24px 24px 48px;
  border-top: 1px solid var(--border);
  display: flex;
  justify-content: space-between;
  align-items: center;
  font-size: 13px;
  color: var(--muted);
}
.blog-footer a {
  color: var(--muted);
  text-decoration: none;
  transition: color 0.15s;
}
.blog-footer a:hover {
  color: var(--fg);
}
.footer-links {
  display: flex;
  gap: 16px;
}

@media (max-width: 640px) {
  .blog-footer {
    padding: 24px 16px 32px;
    flex-direction: column;
    gap: 12px;
  }
}
</style>
