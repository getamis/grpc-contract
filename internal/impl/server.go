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

package impl

import (
	"bytes"
	"fmt"
	"html"
	"html/template"
	"os"

	"github.com/getamis/grpc-contract/internal/util"
)

type Server struct {
	ContractName   string
	ProjectPackage string
}

var ServerTemplate string = `package main

import "context"
import "fmt"
import "log"
import "math/big"
import "net"
import "os"

import "github.com/ethereum/go-ethereum/accounts/abi/bind"
import "github.com/ethereum/go-ethereum/common"
import "github.com/ethereum/go-ethereum/crypto"
import "github.com/ethereum/go-ethereum/ethclient"
import {{ .ContractName }} "{{ .ProjectPackage }}"
import flag "github.com/spf13/pflag"
import "github.com/spf13/viper"
import "google.golang.org/grpc"

const (
	ethereumName        = "ethereum"
	portName            = "port"
	privateKeyName      = "private_key"
	contractAddressName = "contract_address"
)

var (
	ethereumFlag        = flag.String(ethereumName, "ws://127.0.0.1:8546", "the ethereum client address")
	portFlag            = flag.String(portName, "127.0.0.1:5555", "server port")
	privateKeyFlag      = flag.String(privateKeyName, "", "deployer's private key")
	contractAddressFlag = flag.String(contractAddressName, "", "contract address")
)

func main() {
	flag.Parse()

	viper.BindPFlags(flag.CommandLine)
	viper.AutomaticEnv() // read in environment variables that match

	ethereum := viper.GetString(ethereumName)
	if ethereum == "" {
		fmt.Printf("No ethereum client specified\n")
		os.Exit(-1)
	}

	port := viper.GetString(portName)
	if ethereum == "" {
		fmt.Printf("No listen port specified\n")
		os.Exit(-1)
	}

	// connect to ethereum client
	conn, err := ethclient.Dial(ethereum)
	if err != nil {
		fmt.Printf("Failed to connect ethereum: %v\n", err)
		os.Exit(-1)
	}

	privateKey := viper.GetString(privateKeyName)

	// Deploy contracts
	var addr common.Address
	if privateKey != "" {
		// set up auth
		key, err := crypto.ToECDSA(common.Hex2Bytes(privateKey))
		if err != nil {
			fmt.Printf("Failed to get private key: %v\n", err)
			os.Exit(-1)
		}
		auth := bind.NewKeyedTransactor(key)
		auth.GasLimit = big.NewInt(int64(4712388))
		// get nonce
		nonce, err := conn.NonceAt(context.Background(), auth.From, nil)
		if err != nil {
			fmt.Printf("Failed to get nonce: %v\n", err)
			os.Exit(-1)
		}
		auth.Nonce = big.NewInt(int64(nonce))
		auth.GasPrice = big.NewInt(20000000000)

		addr, _, _, err = {{ .ContractName }}.Deploy{{ .ContractName }}(auth, conn)
		if err != nil {
			fmt.Printf("Failed to deploy contract: %v\n", err)
			os.Exit(-1)
		}

		fmt.Printf("Deployed contract: %v\n", addr.Hex())
	} else {
		address := viper.GetString(contractAddressName)
		if address == "" {
			fmt.Printf("No contract address specified\n")
			os.Exit(-1)
		}
		addr = common.HexToAddress(address)
	}

	s := grpc.NewServer()
	{{ .ContractName }}.Register{{ .ContractName }}Server(s, {{ .ContractName }}.NewServer(addr, conn))

	lis, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	s.Serve(lis)
}
`

func (s *Server) Write(filepath, filename string) {
	implTemplate, err := template.New("server").Parse(ServerTemplate)
	if err != nil {
		fmt.Printf("Failed to parse template: %v\n", err)
		os.Exit(-1)
	}
	result := new(bytes.Buffer)
	err = implTemplate.Execute(result, s)
	if err != nil {
		fmt.Printf("Failed to render template: %v\n", err)
		os.Exit(-1)
	}
	content := html.UnescapeString(html.UnescapeString(result.String()))
	util.WriteFile(content, filepath, filename)
}
