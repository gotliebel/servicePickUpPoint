start:
	docker-compose up --build -d zookeeper kafka1 kafka2 kafka3 postgres jaeger
migrateUp:
	./migrations.sh up
migrateStatus:
	./migrations.sh status
down:
	docker-compose down
test:
	go test -v ./...
integration_test:
	go test --tags=integration ./...
generate-mocks:
	mockgen -source ./internal/service/storage.go -destination=./internal/service/mocks/storage.go -package=mocks


.PHONY:start migrateUp migrateStatus down test generate-mocks

LOCAL_BIN:=$(CURDIR)/bin

PROTOC = PATH="$$PATH:$(LOCAL_BIN)" protoc

ORDER_PROTO_PATH:="api/proto/order/v1"

.PHONY: .bin-deps
.bin-deps:
	GOBIN=$(LOCAL_BIN) go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
	GOBIN=$(LOCAL_BIN) go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest
	GOBIN=$(LOCAL_BIN) go install github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-grpc-gateway@latest
	GOBIN=$(LOCAL_BIN) go install github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-openapiv2@latest
	GOBIN=$(LOCAL_BIN) go install github.com/envoyproxy/protoc-gen-validate@latest

.vendor-proto: vendor-proto/google/protobuf vendor-proto/google/api vendor-proto/protoc-gen-openapiv2/options vendor-proto/validate

vendor-proto/protoc-gen-openapiv2/options:
	git clone -b main --single-branch -n --depth=1 --filter=tree:0 \
 		https://github.com/grpc-ecosystem/grpc-gateway vendor.proto/grpc-ecosystem && \
 	cd vendor.proto/grpc-ecosystem && \
	git sparse-checkout set --no-cone protoc-gen-openapiv2/options && \
	git checkout
	mkdir -p vendor.proto/protoc-gen-openapiv2
	mv vendor.proto/grpc-ecosystem/protoc-gen-openapiv2/options vendor.proto/protoc-gen-openapiv2
	rm -rf vendor.proto/grpc-ecosystem


vendor-proto/google/protobuf:
	git clone -b main --single-branch -n --depth=1 --filter=tree:0 \
		https://github.com/protocolbuffers/protobuf vendor.proto/protobuf &&\
	cd vendor.proto/protobuf &&\
	git sparse-checkout set --no-cone src/google/protobuf &&\
	git checkout
	mkdir -p vendor.proto/google
	mv vendor.proto/protobuf/src/google/protobuf vendor.proto/google
	rm -rf vendor.proto/protobuf

vendor-proto/google/api:
	git clone -b master --single-branch -n --depth=1 --filter=tree:0 \
 		https://github.com/googleapis/googleapis vendor.proto/googleapis && \
 	cd vendor.proto/googleapis && \
	git sparse-checkout set --no-cone google/api && \
	git checkout
	mkdir -p  vendor.proto/google
	mv vendor.proto/googleapis/google/api vendor.proto/google
	rm -rf vendor.proto/googleapis

vendor-proto/validate:
	git clone -b main --single-branch --depth=2 --filter=tree:0 \
		https://github.com/bufbuild/protoc-gen-validate vendor.proto/tmp && \
		cd vendor.proto/tmp && \
		git sparse-checkout set --no-cone validate &&\
		git checkout
		mkdir -p vendor.proto/validate
		mv vendor.proto/tmp/validate vendor.proto/
		rm -rf vendor.proto/tmp


.PHONY: generate
generate: .bin-deps .vendor-proto
	mkdir -p pkg/${ORDER_PROTO_PATH}
	protoc -I api/proto \
		-I vendor.proto \
		${ORDER_PROTO_PATH}/order.proto \
		--plugin=protoc-gen-go=$(LOCAL_BIN)/protoc-gen-go --go_out=./pkg/${ORDER_PROTO_PATH} --go_opt=paths=source_relative\
		--plugin=protoc-gen-go-grpc=$(LOCAL_BIN)/protoc-gen-go-grpc --go-grpc_out=./pkg/${ORDER_PROTO_PATH} --go-grpc_opt=paths=source_relative \
		--plugin=protoc-gen-grpc-gateway=$(LOCAL_BIN)/protoc-gen-grpc-gateway --grpc-gateway_out ./pkg/api/proto/order/v1  --grpc-gateway_opt  paths=source_relative --grpc-gateway_opt generate_unbound_methods=true \
		--plugin=protoc-gen-openapiv2=$(LOCAL_BIN)/protoc-gen-openapiv2 --openapiv2_out=./pkg/api/proto/order/v1 \
		--plugin=protoc-gen-validate=$(LOCAL_BIN)/protoc-gen-validate --validate_out="lang=go,paths=source_relative:pkg/api/proto/order/v1"

