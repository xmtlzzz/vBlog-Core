# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

vBlog Core is a customizable, lightweight blog system for geeks and vibe coders. Users can customize blog components. The system consists of a React frontend, Go backend (RESTful + gRPC), PostgreSQL database, and a Wails-based desktop client.

## Tech Stack

- **Frontend**: React, clean minimal UI with light/dark theme (CSS variables)
- **Backend**: Go — RESTful API server + gRPC server (for desktop client communication)
- **Database**: PostgreSQL — each module gets its own table, indexes built on demand
- **Desktop Client**: Wails (Go + web frontend), communicates with blog server via gRPC
- **Content**: Markdown-based blog posts with syntax highlighting

## Architecture

```
vBlog Core
├── web/          # React frontend (blog + admin)
├── server/       # Go backend (REST + gRPC)
├── client/       # Wails desktop client
├── docs/         # PRD and documentation
└── hdx/          # HTML prototypes (reference for UI implementation)
```

### Two Separate Frontend Apps

1. **Blog (public)** — index, post detail, archives, tags, about pages
2. **Admin (management)** — sidebar layout with: dashboard, posts CRUD, tags, comments, component customization, settings, gRPC status

### Go Backend Endpoints

- RESTful API for blog CRUD, tags, comments, settings
- gRPC service for desktop client: blog list change detection, follow/notification sync

### Database Design

- One table per module (posts, tags, comments, settings, etc.)
- Indexes added per actual query patterns, not preemptively

## UI Design System (from prototypes)

All prototypes in `hdx/` define the exact design tokens. Use these CSS variables:

```css
/* Light */
--bg: #fafafa; --surface: #ffffff; --fg: #171717; --muted: #737373;
--border: #e5e5e5; --accent: #2563eb; --accent-soft: #eff6ff;
--success: #16a34a; --warning: #f59e0b; --error: #dc2626;

/* Dark */
--bg: #0a0a0a; --surface: #141414; --fg: #ededed; --muted: #a3a3a3;
--border: #262626; --accent: #3b82f6; --accent-soft: #172554;
--success: #22c55e; --warning: #eab308; --error: #ef4444;

/* Fonts */
--font-sans: -apple-system, BlinkMacSystemFont, 'SF Pro Text', 'Segoe UI', system-ui, sans-serif;
--font-display: -apple-system, BlinkMacSystemFont, 'SF Pro Display', 'Segoe UI', system-ui, sans-serif;
--font-mono: 'JetBrains Mono', 'IBM Plex Mono', ui-monospace, Menlo, monospace;
```

- Theme toggle via `data-theme="dark"` on `<html>`, persist to localStorage
- Nav: sticky top bar with blur backdrop, 56px height
- Admin: 220px sidebar + main content area
- Animations: fade-in on scroll, smooth CSS transitions (0.3s for theme, 0.15s for interactions)
- Status colors: success (green), warning (yellow), error (red) with soft background variants

## Key Pages

| Page | Route | Source |
|------|-------|--------|
| Homepage | `/` | `hdx/index.html` — hero, stats bar, tag filter, post list, pagination |
| Post Detail | `/post/:id` | `hdx/post.html` — markdown content, prev/next nav |
| Archives | `/archives` | `hdx/archives.html` |
| Tags | `/tags` | `hdx/tags.html` |
| About | `/about` | `hdx/about.html` — avatar, bio, tech stack |
| Admin Dashboard | `/admin` | `hdx/admin.html` — overview stats |
| Post Management | `/admin/posts` | `hdx/admin-posts.html` |
| Tag Management | `/admin/tags` | `hdx/admin-tags.html` |
| Comment Management | `/admin/comments` | `hdx/admin-comments.html` |
| Component Custom | `/admin/custom` | `hdx/admin-custom.html` |
| Settings | `/admin/settings` | `hdx/admin-settings.html` |
| gRPC Status | `/admin/grpc` | `hdx/admin-grpc.html` |

## Implementation Principles

- Minimize code — simplest working implementation first
- No over-engineering or premature abstractions
- Each database module is an independent table
- Blog content is Markdown; admin editor outputs Markdown
- REST for web frontend, gRPC only for desktop client communication
