package main

import (
	"testing"

	"github.com/deroproject/derohe/cryptography/crypto"
	"github.com/deroproject/derohe/rpc"
)

func TestF(t *testing.T) {
	var c *RPCConfig = NewRPCConfig(nodeEndpoint, walletEndpoint, username, password)
	address, _ := c.Address()
	this := address.String()
	given := T{
		title:  this,
		number: len(this),
	}

	when := F(&given)
	t.Run(
		"returns bool",
		func(t *testing.T) {
			// Check the outcome
			if !when {
				t.Errorf("Expected F(when) to return true, but got false")
			}
		},
	)
	t.Run(
		"returns values",
		func(t *testing.T) {
			when := T{
				title:  "dero1",
				number: 33,
			}
			t.Run(
				"title is not empty",
				func(t *testing.T) {
					// Check the outcome
					if given.title == "" {
						t.Errorf("expected: %s , got: %s", when.title, given.title)
					}
				},
			)
			t.Run(
				"title as expected",
				func(t *testing.T) {
					// Check the outcome
					if given.title[:5] != when.title[:5] {
						t.Errorf("expected: %s , got: %s", when.title, given.title)
					}
				},
			)
			t.Run(
				"number as expected",
				func(t *testing.T) {
					// Check the outcome
					if given.number/2 != when.number {
						t.Errorf("expected: %s , got: %s", when.title, given.title)
					}
				},
			)
		},
	)
}
func TestTitle(t *testing.T) {
	var c *RPCConfig = NewRPCConfig(nodeEndpoint, walletEndpoint, username, password)

	given, err := c.Title()
	if err != nil {
		t.Error()
	}

	when := len(given)
	expect66chars := 66

	if when != expect66chars {
		t.Errorf("Expected length of %d characters, got %d characters", expect66chars, when)
	}
}

func TestAddress(t *testing.T) {
	var c *RPCConfig = NewRPCConfig(nodeEndpoint, walletEndpoint, username, password)

	given, err := c.Address()
	if err != nil {
		t.Error()
	}
	t.Run(
		"test rpc.Address",
		func(t *testing.T) {
			t.Run(
				"is a DERO address",
				func(t *testing.T) {
					if !given.IsDERONetwork() {
						t.Error()
					}
				},
			)
			t.Run(
				"is mainnet",
				func(t *testing.T) {
					if !given.Mainnet {
						t.Error()
					}
				},
			)
			t.Run(
				"does not contain proof",
				func(t *testing.T) {
					if given.Proof {
						t.Error()
					}
				},
			)
			t.Run(
				"network is 0",
				func(t *testing.T) {
					if given.Network != 0 {
						t.Error()
					}
				},
			)
			t.Run(
				"the string and encoded compressed PulicKey are the same",
				func(t *testing.T) {
					if given.PublicKey.String() != string(given.PublicKey.EncodeCompressed()) {
						t.Error()
					}
				},
			)
			t.Run(
				"test PublicKey length",
				func(t *testing.T) {
					then := 33
					if len(given.Compressed()) != then {
						t.Errorf(
							"Expected length of %d bytes, got %s bytes", // we use %s to stringify the bytes
							then,
							given.Compressed(),
						)
					}
				},
			)
		},
	)
}

func TestGasEstimate(t *testing.T) {
	var c *RPCConfig = NewRPCConfig(nodeEndpoint, walletEndpoint, username, password)
	var s string = "dero1qyvqpdftj8r6005xs20rnflakmwa5pdxg9vcjzdcuywq2t8skqhvwqglt6x0g"
	var scid string = "0000000000000000000000000000000000000000000000000000000000000001"
	var given *rpc.GasEstimate_Result
	var err error

	var (
		a = rpc.Arguments{
			{
				Name:     rpc.SCACTION,   // "SC_ACTION""
				DataType: rpc.DataUint64, // "U"
				Value:    uint64(rpc.SC_CALL),
			},
			{
				Name:     rpc.SCID,     // "SC_ID"
				DataType: rpc.DataHash, // "H"
				Value:    crypto.HashHexToHash(scid),
			},
			{
				Name:     "entrypoint",
				DataType: rpc.DataString, // "S"
				Value:    "Register",
			},
			{
				Name:     "name",         // "SC_ID"
				DataType: rpc.DataString, // "H"
				Value:    "secretnamebasis",
			},
		}

		tp = rpc.Transfer_Params{
			Transfers: nil,
			SC_Value:  0,
			SC_ID:     scid,
			SC_RPC:    a,
			Ringsize:  uint64(128),
			Fees:      uint64(0),
			Signer:    s,
		}
		p rpc.GasEstimate_Params = rpc.GasEstimate_Params(
			tp,
		)
	)

	t.Run(
		"for name service contract", func(t *testing.T) {
			given, err = c.GasEstimate(p)
			if given.Status != "OK" {
				t.Errorf("%s\n", given.Status)
			}

			if err != nil {
				t.Errorf("%s\n", err)
			}
			if given.GasCompute < 0 {
				t.Error()
			}
			if given.GasStorage < 0 {
				t.Error()
			}
		},
	)
	a = rpc.Arguments{
		{
			Name:     rpc.RPC_COMMENT, // "C"
			DataType: rpc.DataString,  // "S"
			Value:    "Hello World",
		},
	}
	scid = "0000000000000000000000000000000000000000000000000000000000000000"
	var (
		d  string       = "dero1qynmz4tgkmtmmspqmywvjjmtl0x8vn5ahz4xwaldw0hu6r5500hryqgptvnj8"
		tx rpc.Transfer = rpc.Transfer{
			Amount:      uint64(1),
			SCID:        crypto.ZEROHASH, // what "ledger"
			Destination: d,
			Payload_RPC: a,
			Burn:        uint64(0),
		}
		txs []rpc.Transfer = []rpc.Transfer{
			tx,
		}
	)
	tp = rpc.Transfer_Params{
		Transfers: txs,
		SC_Value:  0,
		SC_ID:     scid,
		SC_RPC:    nil, // we are getting a stack trace error  because this is nil

		/* cmd/derod/rpc/rpc_dero_estimategas.go

		   if result.GasCompute,
		       result.GasStorage, err = s.RunSC(
		           incoming_values,
		           p.SC_RPC, // because this value doesn't exists it goes *blegh*
		           signer,
		           0,
		       ); err != nil {
		           return
		   }

		*/

		Ringsize: uint64(128),
		Fees:     uint64(0),
		Signer:   s,
	}
	t.Run(
		"for dero",
		func(t *testing.T) {

			p = rpc.GasEstimate_Params(tp)
			given, err = c.GasEstimate(p)
			t.Run(
				"gives an error",
				func(t *testing.T) {
					if err == nil || given != nil {
						t.Errorf("%s\n", given.Status)
					}
				},
			)
		},
	)
}
