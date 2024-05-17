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

### RSSHub

- [x] 36kr
- [x] InfoQ
- [ ] 差评 not updating
- [x] ReadHub
- [x] 东西智库
- [x] DeepMind
- [x] 量子位/tag
- [x] 品玩
- [x] 人人都是产品经理
- [x] 少数派
- [x] Latepost
- [x] Techcrunch
- [ ] X

### 公众号

- [ ] ~~AI产品榜/Kaixinhanguoyu~~
- [x] Founder Park/Founder-Park
- [x] 量子位/QbitAI
- [ ] 世界人工智能大会/gh_00d68db4a358
- [x] 甲子光年/jazzyear
- [x] DeepTech深科技/deeptechchina
- [x] 新智元/AI_era
- [x] 机器之心/almosthuman2014
- [x] 讯飞AIEd/XFjyjsyjy

## Features

- [ ] readhub.cn相关功能
  