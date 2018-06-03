
.PHONY: gen-grpc
gen-grpc:
	@mkdir -p client/gen
	#@mkdir -p server/gen
	protoc --proto_path=proto chat.proto --go_out=plugins=grpc:client/gen/	
	#protoc --proto_path=proto chat.proto --java_out=server/src/main/proto/ # java is generated by gradle!

.PHONY: build-client
build-client:
	cd client && go build client.go

.PHONY: run-client
run-client:
	cd client && ./client

.PHONY: build-server
build-server:
	cd server && ./gradlew installDist

.PHONY: run-server
run-server:
	server/build/install/server/bin/chat-server

.PHONY: clean
clean:
	rm -rf server/build
	rm -rf server/out
	rm -rf client/client