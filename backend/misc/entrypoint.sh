#!/bin/sh

set -e

#echo "path is:$PATH before"
#export PATH=$PATH:/learnerai
#echo "path is:$PATH after"

DEFAULT_API=api.yaml
DEFAULT_CRAWLER=learnerai.yaml


# 如果不为空, 复制环境变量
if [ -n "$CONFIG_API" ]; then
  if [ "$CONFIG_API" != "$DEFAULT_API" ]; then
    echo "cp $CONFIG_API to conf/api.yaml;"
    cp -f ./conf/$CONFIG_API ./conf/api.yaml
  fi
fi

if [ -n "$CONFIG_CRAWLER" ]; then
  if [ "$CONFIG_CRAWLER" != "$DEFAULT_CRAWLER" ]; then
    echo "cp $CONFIG_CRAWLER to conf/learnerai.yaml;"
    cp -f ./conf/$CONFIG_CRAWLER ./conf/learnerai.yaml
  fi
fi


# 最终运行肯定是api.yaml和learnerai.yaml
CONFIG_API=$DEFAULT_API
CONFIG_CRAWLER=$DEFAULT_CRAWLER

echo "after"
echo "final CONFIG_API is: $CONFIG_API"
echo "final CONFIG_CRAWLER is: $CONFIG_CRAWLER"

exec "$@" #./learnerai
