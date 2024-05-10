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
- [ ] 差评 not updating
- [x] ReadHub
- [x] 东西智库
- [x] DeepMind Can't find RSS
- [ ] 量子位/tag
- [x] 品玩
- [x] 人人都是产品经理
- [x] 少数派
- [x] Latepost
- [ ] X
  - [ ] Elon Musk

- [ ] ~~AI产品榜/Kaixinhanguoyu~~
- [ ] Founder Park/Founder-Park
- [ ] 量子位/QbitAI
- [ ] 世界人工智能大会/gh_00d68db4a358
- [ ] 甲子光年/jazzyear
- [ ] DeepTech深科技/deeptechchina
- [ ] 新智元/AI_era
- [ ] 机器之心/almosthuman2014
- [ ] 讯飞AIEd/XFjyjsyjy
- [ ] AgeClub/AgeClub
- [ ] AgeTech新视野/AgeTech2030
- [ ] ITH康养家/huikaolayanglao
- [ ] 创意老龄/chuangyilaoling
- [ ] 开心果carefree/carefreegame
- [ ] 系龄人/Xilingren2020
- [ ] 新老年洞察/NewagingProX
- [ ] 周燕珉工作室/ZYMstudio
- [ ] AgeTech Collaborative/https://agetechcollaborative.org/
- [ ] Aging and Health Technology Watch/https://www.ageinplacetech.com/
- [ ] 美国退休人员协会网站/https://www.aarp.org/
- [ ] 国外老龄化产品分类图/https://thegerontechnologist.com/
- [ ] AI寒武纪/gh_7e5d9d010744
- [ ] HyperAI超神经/HyperAI

## Features

- [ ] readhub.cn相关功能
  