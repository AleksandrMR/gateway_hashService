

get-deps:
	go get -u github.com/brianvoe/gofakeit
	go get -u github.com/grpc-ecosystem/grpc-gateway/v2/runtime
	go get -u google.golang.org/grpc
	go get -u google.golang.org/grpc/credentials/insecure
	go get -u google.golang.org/grpc/reflection
	go get -u google.golang.org/protobuf/types/known/timestamppb
	go get -u github.com/ilyakaznacheev/cleanenv

run-local:
	go run ./cmd/main.go --config=./config/local.yaml
