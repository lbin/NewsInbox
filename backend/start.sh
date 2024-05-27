#sudo docker build -t learnerai . -f local.learnerai.Dockerfile
sudo docker ps -a  | grep learnerai | awk -F " " '{print $1}' |xargs docker container stop|xargs docker container rm
sudo docker run -tid  --restart=always \
 	    --network=host \
            -v /data/logs/learnerai/runtime:/data/logs/learnerai/runtime \
            -v /home/ubuntu/NewsInbox/backend/opt:/home/ubuntu/NewsInbox/backend/opt \
            -v /etc/ssl/certs:/etc/ssl/certs \
            -v /etc/localtime:/etc/localtime \
            -v /mnt:/mnt \
            --name learnerai \
            learnerai
