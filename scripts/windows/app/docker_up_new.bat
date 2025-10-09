cd ..\..\..\

env $(cat ../env/local.env | xargs) docker compose -f deploy/docker/docker-compose.yml -f deploy/docker/docker-compose.local.yml up -d --build

pause
