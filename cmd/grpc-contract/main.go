// This program is free software: you can redistribute it and/or modify
// it under the terms of the GNU Lesser General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.

// This program is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// GNU Lesser General Public License for more details.

// You should have received a copy of the GNU Lesser General Public License
// along with this program.  If not, see <http://www.gnu.org/licenses/>.

package main

import (
	fmt "fmt"
	"io/ioutil"
	"os"
	"path"
	"strings"

	"github.com/getamis/grpc-contract/internal/impl"
	"github.com/getamis/sol2proto/util"
	flag "github.com/spf13/pflag"
	parser "github.com/zpatrick/go-parser"
)

var (
	filepath string
	goTypes  []string
)

func init() {
	flag.StringArrayVar(&goTypes, "types", []string{}, "the go-binding files")
	flag.StringVar(&filepath, "path", ".", "path")
}

func main() {
	flag.Parse()

	// Find all proto generated files
	filesInfos, err := ioutil.ReadDir(filepath)
	if err != nil {
		fmt.Printf("Failed to list files: %v\n", err)
		os.Exit(-1)
	}
	var files []string
	for _, f := range filesInfos {
		if !f.IsDir() && strings.HasSuffix(f.Name(), ".pb.go") {
			files = append(files, path.Join(filepath, f.Name()))
		}
	}
	if len(files) == 0 {
		fmt.Printf("Cannot find the generated proto files")
		os.Exit(-1)
	}
	goFiles, err := parser.ParseFiles(files)
	if err != nil {
		fmt.Printf("Failed to parse file: %v\n", err)
		os.Exit(-1)
	}
	if len(goFiles) == 0 {
		fmt.Printf("Cannot get the go files")
		os.Exit(-1)
	}

	// Save the common util functions
	grpcUtils := impl.Utils{
		Package: goFiles[0].Package,
	}
	grpcUtils.Write(filepath, "grpc_utils.go")

	for _, goType := range goTypes {
		file := path.Join(filepath, goType+".go")
		goBindingFile, err := parser.ParseSingleFile(file)
		if err != nil {
			fmt.Printf("Failed to parse file: %v\n", err)
			os.Exit(-1)
		}

		contract := impl.NewContract(goFiles[0].Package, util.ToCamelCase(goType), append(files, file))

		// Try to find the grpc server intreface
		for _, goFile := range goFiles {
			for _, i := range goFile.Interfaces {
				if !contract.IsServerInterface(i.Name) {
					continue
				}
				for _, m := range i.Methods {
					// Find request struct
					requestStructName := m.Params[1].Type[1:]
					request := findGoStruct(requestStructName, goFiles)
					if request == nil {
						fmt.Printf("Failed to corresponding request struct in method %v\n", m.Name)
						os.Exit(-1)
					}

					// Find response struct
					responseStructName := m.Results[0].Type[1:]
					response := findGoStruct(responseStructName, goFiles)
					if response == nil {
						fmt.Printf("Failed to corresponding response struct in method %v\n", m.Name)
						os.Exit(-1)
					}

					contract.Methods = append(contract.Methods, impl.NewMethod(m, request, response, goBindingFile, contract.StructName))
				}
				break
			}
		}
		contract.Write(filepath, goType+"_server.go")
	}
}

func findGoStruct(name string, goFiles []*parser.GoFile) *parser.GoStruct {
	for _, g := range goFiles {
		for _, s := range g.Structs {
			if name == s.Name {
				return s
			}
		}
	}
	return nil
}
