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
- [x] 差评
- [ ] X
  - [ ] Elon Musk

