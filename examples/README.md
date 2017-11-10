An example to show how to crate DAPP using [sol2proto](https://github.com/getamis/sol2proto) and grpc-contract.

## Build

1. Install [proto](https://github.com/google/protobuf/releases/) and [solc](http://solidity.readthedocs.io/en/develop/installing-solidity.html)

2. Install other tools

```
$ make setup
```

3. Put the solidity file into this folder and build.

```
# name=${contract_name} path=${package_path}
$ make run name=NameService path=github.com/getamis/grpc-contract/examples
```

4. Build server binary

```
make server
```

5. Build client binary

```
make client
```

## Run
1. Run a geth with web socket and 8546 port

2. Run grpc server
- Deploy new contract with deployer private key
```
build/bin/server --private_key $deployer_key
```
- Use existing contract address
```
build/bin/server --contract_address $contract_address
```

3. Run grpc client to connect grpc server
