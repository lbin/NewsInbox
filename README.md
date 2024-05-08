# NewsInbox

## Pre-requisites

- RSSHub

## Installation

```bash
docker-compose -f docker-compose.middleware.yaml -p newsinbox up -d
celery -A app.celery beat --loglevel INFO
```

## Supported News

- [x] 36kr
- [x] 差评
- [ ] X
  - [ ] Elon Musk

