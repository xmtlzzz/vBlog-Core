# gRPC Analytics + Wails Desktop Client Design

Date: 2026-05-10

## Overview

Add gRPC server to vBlog Core for real-time analytics streaming, and build a Wails desktop client that displays blog change notifications and trend comparisons (day/week/month).

## Current State

- REST API (go-restful) fully implemented
- Dashboard stats: 4 aggregate counts via `/api/dashboard/stats`, no time-series
- Post view counting exists (`IncrementViews`)
- No visitor (PV/UV) tracking
- No gRPC code, no `.proto` files, no `client/` directory

## Architecture

**Approach: Change Log + gRPC Streaming Push**

- `daily_stats` table: daily snapshots for trend comparison
- `change_log` table: records each data change for streaming push and history
- `page_views` table: PV/UV tracking
- gRPC server with unary RPCs (trends) + server-streaming RPC (change push)
- Wails client with notification/change center UI

## Database Design

### daily_stats

```sql
CREATE TABLE daily_stats (
    id            SERIAL PRIMARY KEY,
    stat_date     DATE NOT NULL UNIQUE,
    pv            BIGINT DEFAULT 0,
    uv            BIGINT DEFAULT 0,
    post_count    INT DEFAULT 0,
    view_total    BIGINT DEFAULT 0,
    comment_count INT DEFAULT 0,
    tag_count     INT DEFAULT 0,
    created_at    TIMESTAMP DEFAULT NOW()
);
```

### change_log

```sql
CREATE TABLE change_log (
    id          SERIAL PRIMARY KEY,
    change_type VARCHAR(50) NOT NULL,
    target_id   INT,
    title       VARCHAR(200),
    detail      TEXT,
    created_at  TIMESTAMP DEFAULT NOW()
);
CREATE INDEX idx_change_log_created ON change_log(created_at);
```

### page_views

```sql
CREATE TABLE page_views (
    id         SERIAL PRIMARY KEY,
    ip         VARCHAR(45),
    path       VARCHAR(500),
    user_agent TEXT,
    created_at TIMESTAMP DEFAULT NOW()
);
CREATE INDEX idx_pv_created ON page_views(created_at);
CREATE INDEX idx_pv_ip_date ON page_views(ip, created_at);
```

## gRPC Service

### Proto Definition

```protobuf
syntax = "proto3";
package vblog;
option go_package = "vblog-core/proto";

service BlogAnalytics {
  rpc GetTrends(GetTrendsRequest) returns (GetTrendsResponse);
  rpc GetLatestStats(Empty) returns (LatestStats);
  rpc WatchChanges(WatchRequest) returns (stream ChangeEvent);
  rpc Ping(Empty) returns (Empty);
}

message GetTrendsRequest {
  string granularity = 1;  // "day" | "week" | "month"
  int32 count = 2;
}

message TrendPoint {
  string label = 1;
  int64 pv = 2;
  int64 uv = 3;
  int64 view_total = 4;
  int64 comment_count = 5;
  int64 post_count = 6;
}

message GetTrendsResponse {
  repeated TrendPoint points = 1;
}

message LatestStats {
  int64 pv_today = 1;
  int64 uv_today = 2;
  int64 total_posts = 3;
  int64 total_views = 4;
  int64 total_comments = 5;
  int64 total_tags = 6;
  int64 pv_yesterday = 7;
  int64 uv_yesterday = 8;
  int64 views_today_delta = 9;
  int64 comments_today_delta = 10;
}

message WatchRequest {
  string api_key = 1;
  int64 since_id = 2;
}

message ChangeEvent {
  int64 id = 1;
  string type = 2;
  string title = 3;
  string detail = 4;
  string timestamp = 5;
}

message Empty {}
```

### Authentication

- Settings table adds `grpc_api_key` config
- All RPCs except `Ping` require valid api_key in request
- Client stores server address + api_key locally

### Streaming Flow

```
Client connects WatchChanges(api_key, since_id)
  -> Server validates api_key
  -> Reads change_log entries after since_id (catch-up)
  -> Enters watch loop: checks change_log every second
  -> New change found -> push ChangeEvent to client
  -> Client receives -> updates UI notification card
```

## Server Changes

### New Files

| File | Purpose |
|------|---------|
| `server/proto/analytics.proto` | Proto definition |
| `server/proto/analytics.pb.go` | Generated |
| `server/proto/analytics_grpc.pb.go` | Generated |
| `server/grpc/server.go` | gRPC server setup, listener on :50051 |
| `server/grpc/analytics.go` | BlogAnalytics service implementation |
| `server/service/daily_stats.go` | DailyStats CRUD |
| `server/service/change_log.go` | ChangeLog CRUD + write helpers |
| `server/service/page_view.go` | PV/UV recording + stats |
| `server/model/daily_stats.go` | DailyStats model |
| `server/model/change_log.go` | ChangeLog model |
| `server/model/page_view.go` | PageView model |

### Modified Files

| File | Change |
|------|--------|
| `server/cmd/main.go` | Start gRPC server alongside HTTP, add PV middleware |
| `server/service/post.go` | Write change_log on new post |
| `server/service/comment.go` | Write change_log on new comment |
| `server/middleware/pv.go` | New: record page views, deduplicate UV by IP+day |

### Change Log Triggers

Write to `change_log` when:
- New post published: `type=new_post`, title=post title
- New comment created: `type=new_comment`, title="评论 on {post title}"
- Post views hit milestone (100/500/1000/5000): `type=view_milestone`
- Site PV hits milestone: `type=pv_milestone`
- New tag created: `type=tag_added`

### Daily Snapshot

- Cron job at 00:05 daily: aggregate current stats into `daily_stats`
- On first request if today's snapshot missing: compute on-the-fly and cache

## Wails Client

### Project Structure

```
client/
├── main.go
├── app.go
├── proto/                  # Generated gRPC code
├── frontend/
│   ├── src/
│   │   ├── App.vue
│   │   ├── components/
│   │   │   ├── ChangeCard.vue
│   │   │   ├── TrendPanel.vue
│   │   │   ├── StatsBar.vue
│   │   │   └── Settings.vue
│   │   └── main.js
│   └── index.html
└── wails.json
```

### UI Layout

```
┌─────────────────────────────────────┐
│  vBlog Monitor          [Settings]  │
├─────────────────────────────────────┤
│  Stats Bar: PV / UV / Views / Δ%   │
├─────────────────────────────────────┤
│  Change Feed                        │
│  [ChangeCard] new_post: "xxx"       │
│  [ChangeCard] new_comment: "xxx"    │
│  [ChangeCard] view_milestone: 1000  │
├─────────────────────────────────────┤
│  Trend Panel (day/week/month)       │
│  [Mini chart: PV, UV, Views]       │
└─────────────────────────────────────┘
```

### Change Types

| Type | Icon | Trigger |
|------|------|---------|
| `new_post` | New post published |
| `new_comment` | New comment received |
| `view_milestone` | Post views hit 100/500/1000/5000 |
| `pv_milestone` | Site PV hits milestone |
| `tag_added` | New tag created |

### Core Logic (app.go)

- `Connect(addr, apiKey)` — establish gRPC connection
- `WatchChanges()` — goroutine listens stream, emits Wails Events to frontend
- `GetTrends(granularity, count)` — query trend data
- `GetLatestStats()` — get current stats
- `Disconnect()` — close connection

### Frontend

- Vue 3, minimal UI matching blog design system (CSS variables)
- Listen to Wails Events for real-time change cards
- Periodic refresh of stats and trends
- Settings page: server address + API key input
- Trend panel: day/week/month toggle with mini chart

## Implementation Order

1. **Proto + gRPC skeleton** — define proto, generate code, start gRPC server
2. **Database models + migrations** — daily_stats, change_log, page_views tables
3. **PV/UV tracking middleware** — record page views, deduplicate UV
4. **Change log service** — write helpers, integrate into post/comment services
5. **Daily snapshot cron** — aggregate stats daily
6. **gRPC Analytics service** — implement GetTrends, GetLatestStats, WatchChanges
7. **API Key auth** — settings config + gRPC interceptor
8. **Wails client scaffolding** — init project, basic window
9. **Client Go backend** — gRPC client methods, Wails bindings
10. **Client frontend** — ChangeCard, StatsBar, TrendPanel, Settings
