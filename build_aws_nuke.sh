docker build -t aws_nuke_build .
docker run -v $(pwd):/go/src/github.com/rebuy-de/aws-nuke -e UID=$(id -u) -e GID=$(id -g) -it aws_nuke_build
