// Copyright (c) 2017 Arista Networks, Inc.
// Use of this source code is governed by the Apache License 2.0
// that can be found in the COPYING file.

package gnmi

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"io/ioutil"
	"strings"

	pb "github.com/openconfig/gnmi/proto/gnmi"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/metadata"
)

const (
	defaultPort = "6030"
)

// Config is the gnmi.Client config
type Config struct {
	Addr     string
	CAFile   string
	CertFile string
	KeyFile  string
	Password string
	Username string
	TLS      bool
}

// Dial connects to a gnmi service and returns a client
func Dial(cfg *Config) (pb.GNMIClient, error) {
	var opts []grpc.DialOption
	if cfg.TLS || cfg.CAFile != "" || cfg.CertFile != "" {
		tlsConfig := &tls.Config{}
		if cfg.CAFile != "" {
			b, err := ioutil.ReadFile(cfg.CAFile)
			if err != nil {
				return nil, err
			}
			cp := x509.NewCertPool()
			if !cp.AppendCertsFromPEM(b) {
				return nil, fmt.Errorf("credentials: failed to append certificates")
			}
			tlsConfig.RootCAs = cp
		} else {
			tlsConfig.InsecureSkipVerify = true
		}
		if cfg.CertFile != "" {
			if cfg.KeyFile == "" {
				return nil, fmt.Errorf("please provide both -certfile and -keyfile")
			}
			cert, err := tls.LoadX509KeyPair(cfg.CertFile, cfg.KeyFile)
			if err != nil {
				return nil, err
			}
			tlsConfig.Certificates = []tls.Certificate{cert}
		}
		opts = append(opts, grpc.WithTransportCredentials(credentials.NewTLS(tlsConfig)))
	} else {
		opts = append(opts, grpc.WithInsecure())
	}

	if !strings.ContainsRune(cfg.Addr, ':') {
		cfg.Addr += ":" + defaultPort
	}
	conn, err := grpc.Dial(cfg.Addr, opts...)
	if err != nil {
		return nil, fmt.Errorf("failed to dial: %s", err)
	}

	return pb.NewGNMIClient(conn), nil
}

// NewContext returns a new context with username and password
// metadata if they are set in cfg.
func NewContext(ctx context.Context, cfg *Config) context.Context {
	if cfg.Username != "" {
		ctx = metadata.NewOutgoingContext(ctx, metadata.Pairs(
			"username", cfg.Username,
			"password", cfg.Password))
	}
	return ctx
}

// NewGetRequest returns a GetRequest for the given paths
func NewGetRequest(paths [][]string) (*pb.GetRequest, error) {
	req := &pb.GetRequest{
		Path: make([]*pb.Path, len(paths)),
	}
	for i, p := range paths {
		gnmiPath, err := ParseGNMIElements(p)
		if err != nil {
			return nil, err
		}
		req.Path[i] = gnmiPath
	}
	return req, nil
}

// NewSubscribeRequest returns a SubscribeRequest for the given paths
func NewSubscribeRequest(paths [][]string) (*pb.SubscribeRequest, error) {
	subList := &pb.SubscriptionList{
		Subscription: make([]*pb.Subscription, len(paths)),
	}
	for i, p := range paths {
		gnmiPath, err := ParseGNMIElements(p)
		if err != nil {
			return nil, err
		}
		subList.Subscription[i] = &pb.Subscription{Path: gnmiPath}
	}
	return &pb.SubscribeRequest{
		Request: &pb.SubscribeRequest_Subscribe{Subscribe: subList}}, nil
}
