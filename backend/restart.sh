ps -ef | grep learnerai | grep -v grep | awk '{print $2}' | xargs kill -9
rm nohup.out
go build -v  -o learnerai
nohup ./learnerai &