TOOLS_PROTO_BUILDER_IMAGE := cvdb-tools-proto-builder
TOOLS_PROTO_BUILDER_IMAGE_DOCKER_CONTEXT := tools/proto-builder
TOOLS_PROTO_BUILDER_IMAGE_STAMP := .$(TOOLS_PROTO_BUILDER_IMAGE).stamp

.PHONY: build-proto

PROTOS := $(wildcard ./shared/proto/src/*.proto)
PROTOS_IN_CONTAINER := $(patsubst ./shared/proto/src/%, /home/protos/import/%,$(PROTOS))

$(TOOLS_PROTO_BUILDER_IMAGE_STAMP): $(TOOLS_PROTO_BUILDER_IMAGE_DOCKER_CONTEXT)/Dockerfile
	@echo "Building proto compiler image..."
	@docker build -t $(TOOLS_PROTO_BUILDER_IMAGE) --load $(TOOLS_PROTO_BUILDER_IMAGE_DOCKER_CONTEXT)
	@touch $@

build-proto: $(TOOLS_PROTO_BUILDER_IMAGE_STAMP)
	@echo "Building proto files..."
	@docker run \
		-it \
		--rm \
		-v ./shared/proto/src:/home/protos/import \
		-v ./shared/proto/build:/home/protos/export \
		$(TOOLS_PROTO_BUILDER_IMAGE) \
		protoc \
			--proto_path=/home/protos/import/ \
			--go_out=/home/protos/export/ \
			--go_opt=paths=import \
			--go-grpc_out=/home/protos/export/ \
			--go-grpc_opt=paths=import \
			$(PROTOS_IN_CONTAINER)
