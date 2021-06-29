#!/bin/bash
set -e
./start.sh
./deploy_cc.sh

# 변수설정
# CC_SRC_PATH=github.com/paper-contract/
# docker base_path : /opt/gopath/src/
CC_SRC_PATH=chaincode/
CHANNEL_NAME=mychannel
CC_NAME=teamate
VERSION=0.9
CC_RUNTIME_LANGUAGE=go

# maybe in the future
# makeCommand(){
#     cmd=`{"Args":["$1",`
#     return $cmd
# }

# test

docker exec cli peer chaincode invoke -n $CC_NAME -C $CHANNEL_NAME -c '{"Args":["registerUser","user1"]}'
docker exec cli peer chaincode invoke -n $CC_NAME -C $CHANNEL_NAME -c '{"Args":["registerUser","user2"]}'

sleep 3

docker exec cli peer chaincode invoke -n $CC_NAME -C $CHANNEL_NAME -c '{"Args":["registerProject","proj1","user1"]}'
sleep 3
docker exec cli peer chaincode invoke -n $CC_NAME -C $CHANNEL_NAME -c '{"Args":["joinProject","proj1","user1"]}'
sleep 3
docker exec cli peer chaincode invoke -n $CC_NAME -C $CHANNEL_NAME -c '{"Args":["completeProject","proj1"]}'
sleep 3
docker exec cli peer chaincode invoke -n $CC_NAME -C $CHANNEL_NAME -c '{"Args":["recordScore","proj1","user1","30"]}'

