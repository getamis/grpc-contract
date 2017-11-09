grpc-contract:
	go build -v -o ./build/bin/grpc-contract ./cmd/grpc-contract
	@echo "Done building."
	@echo "Run \"$(GOBIN)/grpc-contract\" to launch grpc-contract."
