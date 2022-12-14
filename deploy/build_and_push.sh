docker build --platform linux/amd64 -f deploy/Dockerfile -t go-tron-tiny-wallet . &&
docker tag go-tron-tiny-wallet docker.io/youracc/go-tron-tiny-wallet:latest &&
docker push docker.io/youracc/go-tron-tiny-wallet:latest


