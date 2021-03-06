
.PHONY: install-protoc
install-protoc:
	curl -OL https://github.com/google/protobuf/releases/download/v3.3.0/protoc-3.3.0-linux-x86_64.zip
	unzip protoc-3.3.0-linux-x86_64.zip -d protoc3
	sudo mv protoc3/bin/* /usr/local/bin/
	sudo mv protoc3/include/* /usr/local/include/
	sudo chown $USER /usr/local/bin/protoc\nsudo chown -R $USER /usr/local/include/google
	go get -u github.com/golang/protobuf/protoc-gen-go

.PHONY: gen-grpc
gen-grpc:
	@mkdir -p client/go/gen
	@mkdir -p client/python/gen
	protoc proto/chat.proto -I proto  --go_out=plugins=grpc:client/go/gen
	protoc proto/chat.proto -I proto --python_out=plugins=grpc:client/python/gen
	cp proto/chat.proto server/src/main/proto/ #java is generated by gradle!

.PHONY: build-client
build-client:
	cd client/go && go build client.go

.PHONY: run-client
run-client-go:
	cd client/go/ && ./client

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
	rm -rf client/go/client
