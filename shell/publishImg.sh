
docker build . -t cloud.canister.io:5000/walker088/chatgpt_bot_dongato
docker image prune
docker image ls
docker push cloud.canister.io:5000/walker088/chatgpt_bot_dongato:latest
