package proto

//go:generate protoc ./auth.proto --go_out=./ --go-grpc_out=./
//go:generate protoc ./communications.proto --go_out=./ --go-grpc_out=./
//go:generate protoc ./personalities.proto --go_out=./ --go-grpc_out=./
//go:generate protoc ./message.proto --go_out=./ --go-grpc_out=./
//go:generate protoc ./survey.proto --go_out=./ --go-grpc_out=./
//go:generate protoc ./payments.proto --go_out=./ --go-grpc_out=./
