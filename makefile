go:
	protoc --proto_path=schema schema/*.proto --go_out=server/schema --go-grpc_out=server/schema --grpc-gateway_out=server/schema --go_opt=paths=source_relative --go-grpc_opt=paths=source_relative --grpc-gateway_opt=paths=source_relative
dart:
	protoc --dart_out=grpc:client/lib/ui/service/schema --grpc-gateway_out=logtostderr=true:client/lib/ui/service/schema -I ./schema ./schema/*.proto
run:
	go run server/main.go
test:
	protoc --proto_path=schema schema/*.proto --dart_out=server/schema --dart-grpc_out=server/schema --grpc-gateway_out=server/schema --dart_opt=paths=source_relative --dart-grpc_opt=paths=source_relative --grpc-gateway_opt=paths=source_relative
activate:
	dart pub global activate protoc_plugin