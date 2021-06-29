#!/bin/bash
set -e


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

# turn on cli
# cli is already on

# 설치
# CORE_PEER_ADDRESS
docker exec cli peer chaincode install -n $CC_NAME -v $VERSION -p "$CC_SRC_PATH" 

# # 배포
docker exec cli peer chaincode instantiate -n $CC_NAME -v $VERSION -C "$CHANNEL_NAME" -c '{"Args":[]}' -P "OR ('Org1MSP.member')"

sleep 3

# test
# docker exec cli peer chaincode invoke -n $CC_NAME -C $CHANNEL_NAME -c '{"Args":["issue","00001","MagnetoCorp","MagentCorp","1 Jan 2020","30 Dec 2020","5000000"]}'

# docker exec cli peer chaincode invoke -n $CC_NAME -C $CHANNEL_NAME -c '{"Args":["registerUser","user1"]}'
# docker exec cli peer chaincode invoke -n $CC_NAME -C $CHANNEL_NAME -c '{"Args":["registerUser","user2"]}'

# docker exec cli peer chaincode invoke -n $CC_NAME -C $CHANNEL_NAME -c '{"Args":["registerProject","proj1","user1"]}'
# docker exec cli peer chaincode invoke -n $CC_NAME -C $CHANNEL_NAME -c '{"Args":["joinProject","proj1","user1"]}'
# docker exec cli peer chaincode invoke -n $CC_NAME -C $CHANNEL_NAME -c '{"Args":["completeProject","proj1"]}'
# docker exec cli peer chaincode invoke -n $CC_NAME -C $CHANNEL_NAME -c '{"Args":["recordScore","proj1","user1","30"]}'


# fabcar
# docker exec cli peer chaincode install -n papercontract -v $VERSION -p "$CC_SRC_PATH" -l "$CC_RUNTIME_LANGUAGE"
# docker exec cli peer chaincode instantiate -o orderer.example.com:7050 -C mychannel -n fabcar -l "$CC_RUNTIME_LANGUAGE" -v 1.0 -c '{"Args":[]}' -P "OR ('Org1MSP.member','Org2MSP.member')"
# sleep 10
# docker exec cli peer chaincode invoke -o orderer.example.com:7050 -C mychannel -n fabcar -c '{"function":"initLedger","Args":[]}'
