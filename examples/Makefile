setup:
	# install sol2proto
	go install github.com/getamis/sol2proto
	# install grpc-contract
	go install github.com/getamis/grpc-contract/cmd/grpc-contract
	# install abigen
	go install github.com/ethereum/go-ethereum/cmd/abigen

run:
	mkdir -p contracts/$(pkg)
	sol2proto --pkg $(pkg) --abi $(name).abi > contracts/$(pkg)/$(pkg).proto
	protoc --go_out=plugins=grpc:. contracts/$(pkg)/$(pkg).proto
	abigen --type $(name) --abi $(name).abi --pkg $(pkg) --out ./contracts/$(pkg)/$(pkg).go --bin $(name).bin
	grpc-contract -type $(pkg) -path ./contracts/$(pkg) > ./contracts/$(pkg)/$(pkg)_server.go

server:
	go build -v -o ./build/bin/server ./cmd/server
	@echo "Done building."
	@echo "Run \"$(GOBIN)/server\" to launch server."

client:
	go build -v -o ./build/bin/client ./cmd/client
	@echo "Done building."
	@echo "Run \"$(GOBIN)/client\" to launch client."
