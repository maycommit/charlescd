CIRCLE_SOURCE_PATH=./internal/manager/circle

OUT_DIR=./pkg/grpc

protoc --go_out=./pkg/grpc --go_opt=paths=import \
 --go-grpc_out=./pkg/grpc --go-grpc_opt=paths=import \
 ./internal/manager/circle/circle.proto --experimental_allow_proto3_optional

# protoc -I=$CIRCLE_SOURCE_PATH --go_out=plugins=grpc:$OUT_DIR $CIRCLE_SOURCE_PATH/circle.proto --experimental_allow_proto3_optional