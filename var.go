package main

import (
	"github.com/deroproject/derohe/rpc"
	"github.com/ybbus/jsonrpc"
)

var (
	rpcClient           jsonrpc.RPCClient
	nodeEndpoint        string = "http://secretnamebasis.site:10102/json_rpc"
	walletEndpoint      string = "http://127.0.0.1:10103/json_rpc"
	endpoint            string
	username            string = "secret"
	password            string = "pass"
	encodedAuth         string
	endpointAuth        string
	encodedEndpointAuth string
	opts                *jsonrpc.RPCClientOpts
)

var t = T{
	title:   "dero1qyvqpdftj8r6005xs20rnflakmwa5pdxg9vcjzdcuywq2t8skqhvwqglt6x0g",
	address: rpc.Address{},
	number:  42,
}

var p P = &t
