# NewsInbox

NewsInbox 用于聚合 RSS/公众号等资讯来源，支持内容抓取、摘要生成与推送。项目包含 Flask API、Celery 任务以及核心抓取与处理逻辑。

## 目录结构

- `api/`: Flask API 服务与任务入口。
- `newsinbox/`: 核心逻辑（抓取、摘要、推送等）。
- `backend/`: Go 后端（如有需要可单独部署）。
- `docker/`: 中间件与服务的 Docker 配置。
- `dev/`: 开发工具脚本。
- `tests/`: 测试代码。

## 环境依赖

- Python 3.10+
- PostgreSQL
- Redis
- RSSHub（需要 RSS 源时启用）
- Docker（可选，用于启动中间件）

## 快速开始

### 1. 安装依赖

```bash
python -m venv .venv
source .venv/bin/activate
pip install -r requirements.txt
```

### 2. 配置环境变量

复制并修改配置文件：

```bash
cp api/.env.example api/.env
```

至少需要配置：

- `SECRET_KEY`
- `DB_USERNAME` / `DB_PASSWORD` / `DB_HOST` / `DB_PORT` / `DB_DATABASE`
- `REDIS_HOST` / `REDIS_PORT` / `REDIS_PASSWORD`
- `CELERY_BROKER_URL`

### 3. 启动依赖服务（可选 Docker）

```bash
docker compose -f docker/docker-compose.middleware.yaml -p newsinbox up -d
```

### 4. 启动服务

使用脚本启动（推荐）：

```bash
./api/start.sh
```

或手动启动：

```bash
cd api
celery -A app.celery worker --loglevel INFO -P gevent
celery -A app.celery beat --loglevel INFO

gunicorn --bind 0.0.0.0:6001 --workers 2 --worker-class gevent --timeout 200 --preload app:app
```

停止服务：

```bash
./api/stop.sh
```

### 5. 定时任务示例

```cron
30 00 * * * cd /home/ubuntu/NewsInbox/api && ./stop.sh
35 00 * * * cd /home/ubuntu/NewsInbox/api && ./start.sh
```

## 开发工具

- 代码格式化与 Lint：

```bash
./dev/reformat
```

## 支持的资讯来源

### RSSHub

- [x] 36kr
- [x] InfoQ
- [x] ReadHub
- [x] 东西智库
- [x] DeepMind
- [x] 量子位/tag
- [x] 品玩
- [x] 人人都是产品经理
- [x] 少数派
- [x] Latepost
- [ ] 差评（not updating）
- [x] Techcrunch
- [x] theverge
- [x] producthunt
- [ ] oshwhub.json（解析失败，需要 VPN）
- [ ] fastcompany（解析失败，需要 VPN）
- [ ] X

### 公众号

- [ ] ~~AI产品榜/Kaixinhanguoyu~~
- [x] Founder Park/Founder-Park
- [x] 量子位/QbitAI
- [x] 甲子光年/jazzyear
- [x] DeepTech深科技/deeptechchina
- [x] 新智元/AI_era
- [x] 机器之心/almosthuman2014
- [x] 讯飞AIEd/XFjyjsyjy
- [ ] 世界人工智能大会/gh_00d68db4a358

## Features

- [ ] readhub.cn 相关功能
