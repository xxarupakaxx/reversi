.PHONY: proto
proto:
	@protoc -I . proto/*.proto --go_out=plugins=grpc:gen/pb --go_opt=paths=source_relative