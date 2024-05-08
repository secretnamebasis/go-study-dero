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

func (c *RPCConfig) SetEndpoint(m string) {
	if strings.Contains(m, "DERO.") {
		c.Endpoint = c.NodeEndpoint
	} else {
		c.Endpoint = c.WalletEndpoint
	}
}

func (c *RPCConfig) NewClient(m string) jsonrpc.RPCClient {
	c.SetEndpoint(m)
	return jsonrpc.NewClientWithOpts(c.Endpoint, c.Opts)
}

func (c *RPCConfig) Call(m string, r interface{}, params ...interface{}) error {
	c.Client = c.NewClient(m)
	if len(params) > 0 {
		return c.Client.CallFor(&r, m, params...)
	}
	return c.Client.CallFor(&r, m)
}

func (c *RPCConfig) BlockTemplate(p rpc.GetBlockTemplate_Params) (*rpc.GetBlockTemplate_Result, error) {
	var m = "DERO.GetBlockTemplate"
	var r = new(rpc.GetBlockTemplate_Result)
	c.Client = c.NewClient(m)
	err = c.Call(m, &r, p)
	if err != nil {
		return nil, err
	}
	return r, err
}

func (c *RPCConfig) GasEstimate(params rpc.GasEstimate_Params) (*rpc.GasEstimate_Result, error) {
	var m = "DERO.GetGasEstimate"
	var r = new(rpc.GasEstimate_Result)
	c.Client = c.NewClient(m)
	err = c.Call(m, &r, params)
	if err != nil {
		return nil, err
	}
	return r, err
}

func (c *RPCConfig) Address() (address *rpc.Address, err error) {
	var m = "GetAddress"
	var r = new(rpc.GetAddress_Result)
	err = c.Call(m, r)
	if err != nil {
		return address, err
	}
	return rpc.NewAddress(r.Address)
}

func (c *RPCConfig) Title() (string, error) {
	var m = "GetAddress"
	var r = new(rpc.GetAddress_Result)
	var err = c.Call(m, r)
	if err != nil {
		return "", err
	}

	return r.Address, nil
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
