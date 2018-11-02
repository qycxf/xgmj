#!/usr/bin/env bash

set -e

echo "go build app ..."

# 设置 appname 分支名称，git 提交记录等信息
PROJECT_NAME=$(echo $(dirname $(pwd)) | rev | awk -F \/ '{print $1}' | rev)
BRNAME=$(git symbolic-ref --short HEAD)
BRNAME=${BRNAME/\//-}
VERSION=$(git log --pretty=format:"%h-%cd" --date=short | head -1  | awk '{print $1}')
IMAGE_NAME=$PROJECT_NAME
if [ "$BRNAME"x = "master"x ]; then
	IMAGE_NAME=$PROJECT_NAME
else
	IMAGE_NAME=$PROJECT_NAME-$BRNAME
fi
APP_NAME=$IMAGE_NAME-$VERSION

# 编译 go 程序
echo "docker build $APP_NAME ..."
# CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o app ..

PACKAGE_NAME=qianuuu.com/$PROJECT_NAME
LC_PATH=$PWD/../..
PJ_PATH=/go/src/qianuuu.com
docker run --rm -e CGO_ENABLED=0 \
	   -v $LC_PATH:$PJ_PATH \
	   -w $PJ_PATH/$PROJECT_NAME/docker qianuuu.cn/golang go build -o app ..

# 创建 Dockerfile 文件
rm -f ./Dockerfile
touch ./Dockerfile
echo 'FROM qianuuu.cn/alpine' > Dockerfile
echo ADD app /$APP_NAME >> Dockerfile
echo CMD [\"/${APP_NAME}\"] >> Dockerfile

# 构建 docker image
REGISTRY=qianuuu.cn
echo go build $REGISTRY/$IMAGE_NAME
if [ -n "$1" ]; then
	docker build --force-rm=true -t $REGISTRY/$IMAGE_NAME:$1 .
else
	docker build --force-rm=true -t $REGISTRY/$IMAGE_NAME .
fi

rm ./app
rm ./Dockerfile

echo ''
echo "docker build $REGISTRY/$IMAGE_NAME succeed"

