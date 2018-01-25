proto:
	protoc -I ${GOPATH}/src -I . common/*.proto --go_out=plugins=grpc:.
	protoc -I ${GOPATH}/src -I . software/*.proto --go_out=plugins=grpc:.
	protoc -I ${GOPATH}/src -I . display/*.proto --go_out=plugins=grpc:.
	protoc -I ${GOPATH}/src -I . player/*.proto --go_out=plugins=grpc:.