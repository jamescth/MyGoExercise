#!/bin/sh
# https://github.com/grpc-ecosystem/grpc-gateway
#
# install protocolBuffers
# 	mkdir tmp
# 	cd tmp
# 	git clone https://github.com/google/protobuf
# 	cd protobuf
# 	./autogen.sh
# 	./configure
# 	make
# 	make check
# 	sudo make install
#
# go get
#	go get -u github.com/grpc-ecosystem/grpc-gateway/protoc-gen-grpc-gateway
#	go get -u github.com/grpc-ecosystem/grpc-gateway/protoc-gen-swagger
#	go get -u github.com/golang/protobuf/protoc-gen-go


protoc -I/usr/local/include -I. \
	-I$GOPATH/src \
	-I$GOPATH/src/github.com/grpc-ecosystem/grpc-gateway/third_party/googleapis \
	--go_out=Mgoogle/api/annotations.proto=github.com/grpc-ecosystem/grpc-gateway/third_party/googleapis/google/api,plugins=grpc:. \
	./hello.proto

# generate reverse-proxy
protoc -I/usr/local/include -I. \
	-I$GOPATH/src \
	-I$GOPATH/src/github.com/grpc-ecosystem/grpc-gateway/third_party/googleapis \
	--grpc-gateway_out=logtostderr=true:. \
	./hello.proto
