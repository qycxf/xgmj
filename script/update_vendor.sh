#!/usr/bin/env bash

set -e

DIR=${PWD##*/}

cd ..
echo $PWD

rm -rf vendor
govendor init && govendor add +external
rm -rf vendor/qianuuu.com/poker
rm -rf vendor/qianuuu.com/mahjong

cd $DIR
