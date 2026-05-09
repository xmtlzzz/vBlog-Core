package main

import (
	"fmt"
	"log"
	"math/rand"

	"vblog-core/config"
	"vblog-core/model"
	"vblog-core/service"
)

var titles = []string{
	"Go 语言并发编程最佳实践", "深入理解 PostgreSQL 索引优化", "Vue 3 Composition API 完全指南",
	"Docker 容器化部署实战", "Redis 缓存策略与设计模式", "微服务架构设计原则",
	"Kubernetes 入门到精通", "GraphQL vs REST API 对比", "Linux 性能调优技巧",
	"Git 工作流与团队协作", "TypeScript 高级类型体操", "Rust 所有权系统深入解析",
	"系统设计面试常见题型", "CI/CD 流水线搭建指南", "WebAssembly 前端新范式",
	"OAuth 2.0 认证授权详解", "消息队列选型与实践", "数据库事务隔离级别",
	"HTTP/3 与 QUIC 协议", "WebSocket 实时通信实现", "CSS Grid 布局完全指南",
	"Python 异步编程深入", "Java 虚拟机调优实战", "React Server Components",
	"Go 泛型编程入门", "PostgreSQL 全文搜索", "Nginx 反向代理配置",
	"Elasticsearch 搜索引擎", "Terraform 基础设施即代码", "Prometheus 监控体系",
	"GraphQL Schema 设计", "gRPC 流式通信实战", "MongoDB 聚合管道",
	"Swift 并发编程模型", "Kotlin 协程深入理解", "Flutter 跨平台开发",
	"机器学习工程化实践", "深度学习模型部署", "数据湖架构设计",
	"流式计算框架对比", "分布式系统 CAP 定理", "一致性哈希算法",
	"设计模式之观察者模式", "SOLID 原则实践", "领域驱动设计入门",
	"函数式编程思想", "响应式编程范式", "事件溯源架构",
	"单元测试最佳实践", "代码审查的艺术",
}

var contents = []string{
	"## 引言\n\n在现代软件开发中，%s 是一个非常重要的主题。本文将深入探讨其核心概念和实际应用。\n\n## 基础概念\n\n首先我们需要理解几个关键概念。这些概念是后续学习的基础。\n\n```\n// 示例代码\nfunc main() {\n    fmt.Println(\"Hello, World!\")\n}\n```\n\n## 实践经验\n\n经过多年的工程实践，我总结了以下几点经验：\n\n- 保持代码简洁\n- 编写可测试的代码\n- 持续重构\n\n## 总结\n\n%s 是每个开发者都应该掌握的技能。希望本文对你有所帮助。",
	"## 背景\n\n最近在项目中遇到了一个关于 %s 的问题，花了不少时间解决，在此记录一下。\n\n## 问题描述\n\n在生产环境中，我们发现系统性能出现了明显下降。经过排查，问题出在 %s 上。\n\n## 解决方案\n\n经过多次尝试，我们最终采用了以下方案：\n\n1. 优化数据结构\n2. 引入缓存机制\n3. 异步处理非关键路径\n\n```go\n// 优化后的代码\nfunc process(data []byte) error {\n    // 处理逻辑\n    return nil\n}\n```\n\n## 效果\n\n优化后，系统性能提升了 3 倍，响应时间从 200ms 降低到 60ms。",
	"## 前言\n\n今天来聊聊 %s 这个话题。作为一个有 5 年经验的开发者，我想分享一些自己的见解。\n\n## 为什么要学习它？\n\n在当今快速发展的技术领域，掌握 %s 可以让你：\n\n- 提高开发效率\n- 写出更优雅的代码\n- 解决复杂的技术问题\n\n## 核心要点\n\n### 第一点\n\n理解底层原理比记住 API 更重要。\n\n### 第二点\n\n实践出真知，多写代码多踩坑。\n\n### 第三点\n\n阅读优秀项目的源码是提升最快的方式。\n\n## 推荐资源\n\n- 官方文档\n- GitHub 开源项目\n- 技术博客\n\n## 结语\n\n学习是一个持续的过程，保持好奇心最重要。",
}

var tags = []string{
	"Go", "Vue", "Docker", "PostgreSQL", "Redis", "TypeScript", "React",
	"Kubernetes", "Linux", "Git", "Python", "Rust", "系统设计", "前端",
	"后端", "数据库", "DevOps", "架构", "性能优化", "安全",
}

func main() {
	cfg := config.Load()
	db, err := cfg.DB.Connect()
	if err != nil {
		log.Fatalf("db connect failed: %v", err)
	}

	postSvc := service.NewPostService(db)

	// Ensure tags exist
	var tagModels []model.Tag
	for _, name := range tags {
		var t model.Tag
		if err := db.Where("name = ?", name).First(&t).Error; err != nil {
			t = model.Tag{Name: name}
			db.Create(&t)
		}
		tagModels = append(tagModels, t)
	}

	statuses := []string{"published", "published", "published", "draft"}

	for i := 0; i < 50; i++ {
		title := titles[i%len(titles)]
		if i >= len(titles) {
			title = fmt.Sprintf("%s（续篇 %d）", titles[i%len(titles)], i/len(titles)+1)
		}
		content := fmt.Sprintf(contents[i%len(contents)], title, title)

		// Pick 1-3 random tags
		n := rand.Intn(3) + 1
		shuffled := make([]model.Tag, len(tagModels))
		copy(shuffled, tagModels)
		rand.Shuffle(len(shuffled), func(i, j int) { shuffled[i], shuffled[j] = shuffled[j], shuffled[i] })
		postTags := shuffled[:n]

		post := &model.Post{
			Title:   title,
			Content: content,
			Status:  statuses[rand.Intn(len(statuses))],
			Tags:    postTags,
		}

		if err := postSvc.Create(post); err != nil {
			log.Printf("failed to create post %d: %v", i+1, err)
		} else {
			log.Printf("created post %d: %s", i+1, post.Title)
		}
	}

	log.Println("done! 50 test posts created.")
}
