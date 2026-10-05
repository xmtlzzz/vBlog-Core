INSERT OR REPLACE INTO settings (key,value) VALUES ('site_title','vBlog'), ('author_name','xmtlz'), ('site_description','关于代码、工具与开发过程的记录。'), ('enable_comments','false');
INSERT OR REPLACE INTO tags (id,name) VALUES (1,'开发');
INSERT OR REPLACE INTO post_tags (post_id,tag_id) VALUES (3,1);
INSERT OR REPLACE INTO posts (id,title,content,excerpt,status,read_time,created_at) VALUES
(3,'项目导航','# 项目导航

记录网站与工具的开发过程，让作品和知识有一个共同入口。

## 开源项目

- [vMaker](https://vmaker.xmtlz.dev)：浏览开源作品与最近进展。
- [vBlog Core](https://github.com/xmtlzzz/vBlog-Core)：Markdown 博客与内容管理。

## 开发笔记

```md
# 代码里的一级标题保持不变
```

把想法落成工具，也记录做出这些选择的原因。','# 项目 [导航](https://example.com) **介绍**','published',2,'2026-08-31'),
(5,'WanFunControlToVideo','# WanFunControlToVideo

学习视频工作流时的配置记录。

## 工作流笔记

保持参考与实验结果可追溯。','开发工具与工作流的学习笔记。','published',3,'2026-09-01'),
(4,'本地草稿测试','# 测试内容','用于检查草稿是否进入公开列表。','draft',1,'2026-09-02');
