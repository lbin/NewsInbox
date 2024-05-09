# NewsInbox

## Pre-requisites

- Python 3.10+
- RSSHub

## Installation

```bash
docker-compose -f docker-compose.middleware.yaml -p newsinbox up -d
nohup celery -A app.celery beat --loglevel INFO &
nohup celery -A app.celery worker --loglevel INFO -P gevent &
```

## Supported News

- [x] 36kr
- [x] InfoQ
- [x] 差评 not updating
- [x] ReadHub
- [x] 东西智库
- [x] DeepMind
- [x] 量子位/tag
- [x] 品玩
- [x] 人人都是产品经理
- [x] 少数派
- [x] X
  - [ ] Elon Musk
