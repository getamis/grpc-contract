package main

import (
	"context"
	"fmt"
	"log"
	"math/big"
	"net"
	"os"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
	contracts "github.com/getamis/grpc-contract/examples/pb"
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
	ethereumFlag        = flag.String(ethereumName, "http://192.168.99.100:8545", "the ethereum client address")
	portFlag            = flag.String(portName, "127.0.0.1:5555", "server port")
	privateKeyFlag      = flag.String(privateKeyName, "9ad3ea7650babad5d1976b75b3141278942cebbe423e84d7a6800ae67a0a74b5", "deployer's private key")
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

		addr, _, _, err = contracts.DeployNameService(auth, conn)
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
	contracts.RegisterNameServiceServer(s, contracts.NewNameServiceServer(addr, conn, nil))

	lis, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	s.Serve(lis)
}
