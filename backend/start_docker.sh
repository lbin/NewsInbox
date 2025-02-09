sudo docker build -t learnerai . -f local.learnerai.Dockerfile
sudo docker ps -a  | grep learnerai | awk -F " " '{print $1}' |xargs docker container stop|xargs docker container rm
sudo docker run -tid  --restart=always \
 	    --network=host \
        -v /home/ubuntu/NewsInbox/backend/conf:/home/ubuntu/NewsInbox/backend/conf \
        -v /etc/localtime:/etc/localtime \
        --name learnerai \
        learnerai