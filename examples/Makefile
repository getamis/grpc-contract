filename:=$(shell echo '$(name)' | perl -pe 's/([A-Z])/_\L\1/g' | sed 's/^_//')

setup:
	# install proto, protoc-gen-go
	go get github.com/golang/protobuf/proto
	go get github.com/golang/protobuf/protoc-gen-go
	# install sol2proto
	go get github.com/getamis/sol2proto
	# install grpc-contract
	go get github.com/getamis/grpc-contract/cmd/grpc-contract
	# install abigen
	go get github.com/ethereum/go-ethereum/cmd/abigen

run:
	solc --bin --abi --optimize --overwrite -o . $(name).sol
	mkdir -p pb
	sol2proto --pkg pb --abi $(name).abi
	mv *.proto pb
	protoc --go_out=plugins=grpc:./pb -I pb pb/*.proto -I ../vendor pb/*.proto
	abigen --type $(name) --abi $(name).abi --pkg pb --out ./pb/$(filename).go --bin $(name).bin
	grpc-contract --types $(filename) --path ./pb --pb-path ./pb

server:
	go build -v -o ./build/bin/server ./cmd/server
	@echo "Done building."
	@echo "Run \"$(GOBIN)/server\" to launch server."

client:
	go build -v -o ./build/bin/client ./cmd/client
	@echo "Done building."
	@echo "Run \"$(GOBIN)/client\" to launch client."
