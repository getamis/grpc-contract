package main

import (
	"flag"
	fmt "fmt"
	"os"
	"path"

	"github.com/getamis/grpc-contract/internal/impl"
	"github.com/getamis/sol2proto/util"
	parser "github.com/zpatrick/go-parser"
)

var (
	filepath    string
	goType      string
	packagePath string
)

func init() {
	flag.StringVar(&goType, "type", "", "the go file from proto")
	flag.StringVar(&filepath, "path", ".", "path")
	flag.StringVar(&packagePath, "package", ".", "package path")
}

func main() {
	flag.Parse()

	goFile, err := parser.ParseFile(path.Join(filepath, goType+".pb.go"))
	if err != nil {
		fmt.Printf("Failed to parse file: %v\n", err)
		os.Exit(-1)
	}

	goBindingFile, err := parser.ParseFile(path.Join(filepath, goType+".go"))
	if err != nil {
		fmt.Printf("Failed to parse file: %v\n", err)
		os.Exit(-1)
	}

	contract := impl.Contract{
		Imports: []string{
			"context",
			"math/big",
			"os",

			"github.com/ethereum/go-ethereum/accounts/abi/bind",
			"github.com/ethereum/go-ethereum/common",
			"github.com/ethereum/go-ethereum/crypto",
		},
		Package: goFile.Package,
		Name:    util.ToCamelCase(goType),
	}

	// Try to find the grpc server intreface
	for _, i := range goFile.Interfaces {
		if !contract.IsServerInterface(i.Name) {
			continue
		}
		for _, m := range i.Methods {
			// Find request struct
			requestStructName := m.Params[1].Type[1:]
			var request *parser.GoStruct
			for _, s := range goFile.Structs {
				if requestStructName == s.Name {
					request = s
					break
				}
			}
			if request == nil {
				fmt.Printf("Failed to corresponding request struct in method %v\n", m.Name)
				os.Exit(-1)
			}

			// Find response struct
			responseStructName := m.Results[0].Type[1:]
			var response *parser.GoStruct
			for _, s := range goFile.Structs {
				if responseStructName == s.Name {
					response = s
					break
				}
			}
			if response == nil {
				fmt.Printf("Failed to corresponding response struct in method %v\n", m.Name)
				os.Exit(-1)
			}

			contract.Methods = append(contract.Methods, impl.NewMethod(m, request, response, goBindingFile))
		}
		break
	}
	contract.Write(filepath, goType+"_server.go")

	server := &impl.Server{
		ContractName:    goType,
		ContractPackage: path.Join(packagePath, filepath),
	}
	server.Write("cmd/server", "main.go")
}
