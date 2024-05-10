package main

import (
	"fmt"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/deroproject/derohe/cryptography/crypto"
	"github.com/deroproject/derohe/rpc"
)

func TestF(t *testing.T) {
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
	testCases := []struct {
		name  string
		check func(*testing.T, *rpc.Address)
	}{
		{
			name: "is a DERO address",
			check: func(t *testing.T, given *rpc.Address) {
				if !given.IsDERONetwork() {
					t.Error()
				}
			},
		},
		{
			name: "is mainnet",
			check: func(t *testing.T, given *rpc.Address) {
				if !given.Mainnet {
					t.Error()
				}
			},
		},
		{
			name: "does not contain proof",
			check: func(t *testing.T, given *rpc.Address) {
				if given.Proof {
					t.Error()
				}
			},
		},
		{
			name: "network is 0",
			check: func(t *testing.T, given *rpc.Address) {
				if given.Network != 0 {
					t.Error()
				}
			},
		},
		{
			name: "the string and encoded compressed PublicKey are the same",
			check: func(t *testing.T, given *rpc.Address) {
				if given.PublicKey.String() != string(given.PublicKey.EncodeCompressed()) {
					t.Error()
				}
			},
		},
		{
			name: "test PublicKey length",
			check: func(t *testing.T, given *rpc.Address) {
				expectedLength := 33
				if len(given.Compressed()) != expectedLength {
					t.Errorf(
						"Expected length of %d bytes, got %s bytes", // we use %s to stringify the bytes
						expectedLength,
						given.Compressed(),
					)
				}
			},
		},
	}

	given, err := c.Address()
	if err != nil {
		t.Error()
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			tc.check(t, given)
		})
	}
}

func TestGasEstimate(t *testing.T) {
	var s string = "dero1qyvqpdftj8r6005xs20rnflakmwa5pdxg9vcjzdcuywq2t8skqhvwqglt6x0g"
	var scid string = "0000000000000000000000000000000000000000000000000000000000000001"
	var given *rpc.GasEstimate_Result

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
				Name:     "name",
				DataType: rpc.DataString, // "S"
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
			if given.GasCompute < 1 {
				t.Error()
			}
			if given.GasStorage < 1 {
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

		   if r.GasCompute,
		       r.GasStorage, err = s.RunSC(
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
	p = rpc.GasEstimate_Params(tp)
	t.Run(
		"for dero",
		func(t *testing.T) {
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

func TestBlockTemplate(t *testing.T) {
	var p rpc.GetBlockTemplate_Params = rpc.GetBlockTemplate_Params{
		Miner:          c.Username,
		Wallet_Address: "dero1qyvqpdftj8r6005xs20rnflakmwa5pdxg9vcjzdcuywq2t8skqhvwqglt6x0g",
		Block:          false,
	}
	var given *rpc.GetBlockTemplate_Result
	given, err = c.BlockTemplate(p)
	fmt.Printf("%+v\n", given)
	if err != nil {
		t.Errorf("%s", err)
	}
	// Ensure given is not nil before accessing its fields
	if given == nil {
		t.Errorf("Received nil result")
		return // Exiting early if the result is nil
	}
	t.Run(
		"test JobID",
		func(t *testing.T) {
			// EX: JOB ID: 1715200867796.0.secret

			// Validate JobID format
			when := strings.Split(given.JobID, ".")
			if len(when) != 3 {
				t.Errorf("Invalid JobID format: %s", given.JobID)
			}
			t.Run(
				"JobID contains network",
				func(t *testing.T) {
					if !strings.EqualFold(when[1], "0") {
						t.Errorf("JobID does not contain newtwork")
					}
				},
			)
			t.Run(
				"JobID contains username",
				func(t *testing.T) {
					if !strings.EqualFold(when[2], c.Username) {
						t.Errorf("JobID does not contain username")
					}
				},
			)

			t.Run(
				"JobId contains timestamp",
				func(t *testing.T) {
					var timestamp, err = strconv.ParseInt(when[0], 10, 64)
					if err != nil {
						t.Errorf("Invalid timestamp in JobID: %s", when[0])
					}
					seconds := timestamp / 1000 // Convert milliseconds to seconds

					var timeObj = time.Unix(seconds, 0)
					if timeObj.IsZero() {
						t.Errorf("Invalid timestamp value: %d", seconds)
					}
					var test = []struct {
						test    *testing.T
						method  string
						x, y, z int
					}{
						{
							test:   t,
							method: "year",
							x:      timeObj.Year(),
							y:      2000,
							z:      2100,
						},
						{
							test:   t,
							method: "month",
							x:      int(timeObj.Month()),
							y:      1,
							z:      12,
						},
						{
							test:   t,
							method: "day",
							x:      timeObj.Day(),
							y:      1,
							z:      31,
						},
						{
							test:   t,
							method: "hour",
							x:      timeObj.Hour(),
							y:      0,
							z:      23,
						},
						{
							test:   t,
							method: "minute",
							x:      timeObj.Minute(),
							y:      0,
							z:      59,
						},
						{
							test:   t,
							method: "second",
							x:      timeObj.Second(),
							y:      0,
							z:      59,
						},
					}

					var validate = func(t *testing.T, method string, x, y, z int) {
						t.Helper() // t.Helper() is a method that marks the calling function as a helper function and not a test
						if x < y || x > z {
							t.Errorf("Invalid %s in timestamp: %d", method, x)
						}
					}

					for _, tc := range test {
						t.Run(
							"validate time components",
							func(t *testing.T) {
								validate(tc.test, tc.method, tc.x, tc.y, tc.z)
							},
						)
					}
				},
			)
		},
	)
}
