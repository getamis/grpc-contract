An example to show how to crate DAPP using [sol2proto](https://github.com/getamis/sol2proto) and grpc-contract.

## Build
This example requires a Go (version 1.7 or later) compiler. You can install them using your favorite package manager. Once the dependencies are installed,

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

## Trouble shootings
1. Old export format no longer supported
- Got error message like this:
```
could not import google.golang.org/grpc (/path/pkg/darwin_amd64/google.golang.org/grpc.a: import "google.golang.org/grpc": old export format no longer supported (recompile library))
```
- How to resolve?

Some libraries installed in your `$GOPATH` are in the old format, which are built go1.6 or before. Make sure all libraries under your `$GOPATH` are recompiled with your current go compiler.
```
go install google.golang.org/grpc
```
2. Inconsistent context library path
- Got error message like this:
```
cannot use server literal (type *server) as type NameServiceServer in return argument:
    *server does not implement NameServiceServer (wrong type for GetName method)
        have GetName("context".Context, *Empty) (*GetNameResp, error)
        want GetName("golang.org/x/net/context".Context, *Empty) (*GetNameResp, error)
```
- How to resolve?

The `context` is in the standard library Go 1.7 already. Make sure the latest version of grpc and protoc plugin are installed.
```
go get -u google.golang.org/grpc
go get -u github.com/golang/protobuf/protoc-gen-go
```

3. Wrong imports in generated file
- You may find you have wrong imports in the generated file, such as:

```
github.com/markya0616/go-ethereum/core/types
```

The correct one should be:
```
github.com/ethereum/go-ethereum/core/types
```

The main reason is that [goImports](https://godoc.org/golang.org/x/tools/cmd/goimports) finds the wrong imports. The solution is to create a configuration file at `$GOPATH/src/.goimportsignore` and put in the ignore wrong path, i.e., `github.com/markya0616/go-ethereum/`. After that, generate files again, it should work.