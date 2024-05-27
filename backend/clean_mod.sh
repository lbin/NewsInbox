export GIT_TERMINAL_PROMPT=1
go env -w GOPROXY=https://goproxy.cn,direct
go clean --modcache
git config --global credential.helper store
