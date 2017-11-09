package main

import (
	"fmt"
	"log"
	"net"
	"os"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/getamis/grpc-contract/examples/account"
	"github.com/getamis/grpc-contract/examples/contracts/name_service"
	flag "github.com/spf13/pflag"
	"github.com/spf13/viper"
	"google.golang.org/grpc"
)

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
		account := account.New(conn, privateKey)

		addr, _, _, err = name_service.DeployNameService(account.TransactOpts(), conn)
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
	name_service.RegisterNameServiceServer(s, name_service.NewServer(addr, conn))

	lis, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	s.Serve(lis)
}
