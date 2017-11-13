package main

import (
	"context"
	"fmt"
	"os"

	"github.com/getamis/grpc-contract/examples/contracts/name_service"
	flag "github.com/spf13/pflag"
	"github.com/spf13/viper"
	"google.golang.org/grpc"
)

const (
	serverName = "server"
)

var (
	serverFlag = flag.String(serverName, "127.0.0.1:5555", "server port")
)

func main() {
	flag.Parse()

	viper.BindPFlags(flag.CommandLine)
	viper.AutomaticEnv() // read in environment variables that match

	server := viper.GetString(serverName)
	if server == "" {
		fmt.Printf("No server specified\n")
		os.Exit(-1)
	}

	// Set up a connection to the server.
	conn, err := grpc.Dial(server, grpc.WithInsecure())
	if err != nil {
		fmt.Printf("Did not connect: %v\n", err)
		os.Exit(-1)
	}
	defer conn.Close()
	c := name_service.NewNameServiceClient(conn)

	// Contact the server and print out its response.
	tx, err := c.SetName(context.Background(), &name_service.SetNameReq{
		Opts: &name_service.TransactOpts{
			PrivateKey: "9ad3ea7650babad5d1976b75b3141278942cebbe423e84d7a6800ae67a0a74b5",
			Nonce:      -1,
			Value:      0,
			GasLimit:   4712388,
		},
		Name: "Mark",
	})
	if err != nil {
		fmt.Printf("Failed to set name: %v\n", err)
		os.Exit(1)
	}
	fmt.Printf("Set name in tx: %v\n", tx.Hash)

	r, err := c.GetName(context.Background(), &name_service.EmptyReq{})
	if err != nil {
		fmt.Printf("Failed to get name: %v\n", err)
		os.Exit(1)
	}
	fmt.Printf("Name: %v\n", r.Arg)
}
