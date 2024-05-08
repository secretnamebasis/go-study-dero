package main

import (
	"encoding/base64"
	"strings"

	"github.com/deroproject/derohe/rpc"
	"github.com/ybbus/jsonrpc"
)

func NewRPCConfig(nodeEndpoint, walletEndpoint, username, password string) *RPCConfig {
	endpointAuth = username + ":" + password
	encodedAuth = base64.StdEncoding.EncodeToString([]byte(endpointAuth))
	opts = &jsonrpc.RPCClientOpts{
		CustomHeaders: map[string]string{
			"Authorization": "Basic " + encodedAuth,
		},
	}

	return &RPCConfig{
		NodeEndpoint:   nodeEndpoint,
		WalletEndpoint: walletEndpoint,
		Username:       username,
		Password:       password,
		EndpointAuth:   endpointAuth,
		EncodedAuth:    encodedAuth,
		Opts:           opts,
	}
}

func (c *RPCConfig) SetEndpoint(method string) {
	if strings.Contains(method, "DERO.") {
		c.Endpoint = c.NodeEndpoint
	} else {
		c.Endpoint = c.WalletEndpoint
	}
}

func (c *RPCConfig) NewClient(method string) jsonrpc.RPCClient {
	c.SetEndpoint(method)
	return jsonrpc.NewClientWithOpts(c.Endpoint, c.Opts)
}

func (c *RPCConfig) Call(method string, result interface{}, params ...interface{}) error {
	c.Client = c.NewClient(method)
	if len(params) > 0 {
		return c.Client.CallFor(&result, method, params...)
	}
	return c.Client.CallFor(&result, method)
}

func (c *RPCConfig) GasEstimate(params rpc.GasEstimate_Params) (*rpc.GasEstimate_Result, error) {
	var method = "DERO.GetGasEstimate"
	var result = new(rpc.GasEstimate_Result)
	c.Client = c.NewClient(method)
	var err = c.Call(method, &result, params)
	if err != nil {
		return nil, err
	}
	return result, err
}

func (c *RPCConfig) Address() (address *rpc.Address, err error) {
	var method = "GetAddress"
	var result = new(rpc.GetAddress_Result)
	err = c.Call(method, result)
	if err != nil {
		return address, err
	}
	return rpc.NewAddress(result.Address)
}

func (c *RPCConfig) Title() (string, error) {
	var method = "GetAddress"
	var result = new(rpc.GetAddress_Result)
	var err = c.Call(method, result)
	if err != nil {
		return "", err
	}

	return result.Address, nil
}

func F(x P) bool {
	if x.hasString() && x.hasInt() {
		return true
	}
	return false
}

func (t T) hasString() bool {
	_, ok := any(t.title).(string)
	return ok
}

func (t T) hasInt() bool {
	_, ok := any(t.number).(int)
	return ok
}
