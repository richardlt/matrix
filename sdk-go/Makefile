# Regenerate the Go bindings from the .proto files.
#
# Requires protoc (https://github.com/protocolbuffers/protobuf/releases) and the two
# plugins, which are plain Go binaries:
#   go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
#   go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest
#
# The proto include root is this directory, so imports read "common/types.proto". Output
# paths come from each file's `option go_package`, with the module prefix stripped by
# --go_opt=module=..., which lands every generated file next to its .proto.
MODULE := github.com/richardlt/matrix/sdk-go

proto:
	protoc -I . \
		--go_out=. --go_opt=module=$(MODULE) \
		--go-grpc_out=. --go-grpc_opt=module=$(MODULE) \
		common/*.proto display/*.proto player/*.proto software/*.proto
