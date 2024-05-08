package main

import (
	"github.com/deroproject/derohe/rpc"
	"github.com/ybbus/jsonrpc"
)

// Define the T struct
type T struct {
	title   string
	address rpc.Address
	number  int
}

type P interface {
	hasString() bool
	hasInt() bool
}

type RPCConfig struct {
	Client         jsonrpc.RPCClient
	NodeEndpoint   string
	WalletEndpoint string
	Endpoint       string
	Username       string
	Password       string
	EndpointAuth   string
	EncodedAuth    string
	Opts           *jsonrpc.RPCClientOpts
}
