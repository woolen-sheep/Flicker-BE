# Flicker-BE

## 简介

Flicker闪卡的后端，同时也是华科CS2019的软件工程Project。

Flicker是一个基于 Web 平台的知识分享、学习与记忆平台，且支持更高的可定制化程度、更广的学习内容覆盖面。

## 项目结构

```
├─config - 读取并解析配置文件
├─configs - 放置配置文件
├─constant - 定义全局常量
├─controller - controller层，负责处理请求
│  └─param - 定义前后端交互的参数类型
├─docs - 放置文档
├─middleware - 自定义中间件
├─model - model层，负责与数据库交互
├─router - 路由
└─util - 工具
    ├─context - 全局上下文工具
    └─log - logger
```

## 部署

依赖：

- docker
- docker-compose

```bash
git clone https://github.com/woolen-sheep/Flicker-BE.git
cd Flicker-BE
sudo docker-compose up --build
```

