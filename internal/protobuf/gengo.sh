#!/usr/bin/env bash

set -e

docker run -v `pwd`:/defs qianuuu.cn/protoc-go

mv ./pb-go/*.go ./

rm -r ./pb-go

