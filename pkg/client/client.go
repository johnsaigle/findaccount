package client

import (
	"context"
	"errors"
	"fmt"
	"regexp"
	"strings"

	"github.com/cosmos/cosmos-sdk/types/bech32"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	staketypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	"github.com/johnsaigle/findaccount/types"
	findaccounttypes "github.com/johnsaigle/findaccount/types"
	rpchttp "github.com/tendermint/tendermint/rpc/client/http"
)

var portRex = regexp.MustCompile(`.*:\d+$`)
var protoRex = regexp.MustCompile(`^\w+://`)

// TODO adding REST API support would be nice for nodes that do not have RPC enabled
func NewClient(rpcaddress string) (*rpchttp.HTTP, error) {
	client := &rpchttp.HTTP{}
	var err error
	rpcaddress= strings.TrimRight(rpcaddress, "/")
	var unknown bool

	if !portRex.MatchString(rpcaddress) {
		switch protoRex.FindString(rpcaddress) {
		case "https://":
			rpcaddress = rpcaddress + ":443"
		case "http://":
			rpcaddress = rpcaddress + ":80"
		case "tcp://":
			rpcaddress = rpcaddress + ":26657"
		default:
			unknown = true
		}
	}
	if unknown {
		err = errors.New("Unknown protocol")
		return nil, err
	}
	client, err = rpchttp.NewWithTimeout(rpcaddress, "/websocket", 10)
	if err != nil {
		return nil, err
	}
	status, err := client.Status(context.Background())
	if err != nil || status.SyncInfo.CatchingUp {
		return nil, err
	}
	return client, err
}

func NewClientFromChainInfo(rpcs []types.Rpc, chain string) (*rpchttp.HTTP, error) {
	client := &rpchttp.HTTP{}
	var err error
	ok := false
	for i := range rpcs {
		endpoint := rpcs[len(rpcs)-1-i]
		endpoint.Address = strings.TrimRight(endpoint.Address, "/")
		var unknown bool

		if !portRex.MatchString(endpoint.Address) {
			switch protoRex.FindString(endpoint.Address) {
			case "https://":
				endpoint.Address = endpoint.Address + ":443"
			case "http://":
				endpoint.Address = endpoint.Address + ":80"
			case "tcp://":
				endpoint.Address = endpoint.Address + ":26657"
			default:
				unknown = true
			}
		}
		if unknown {
			err = errors.New("Unknown protocol")
			continue
		}
		client, err = rpchttp.NewWithTimeout(endpoint.Address, "/websocket", 10)
		if err != nil {
			continue
		}
		status, err := client.Status(context.Background())
		if err != nil || status.SyncInfo.CatchingUp {
			continue
		}
		ok = true
		break
	}
	if !ok {
		err = fmt.Errorf("could not connect to any endpoints for %s: %w", chain, err)
	}
	return client, err
}

// TODO change to accept a client as parameter rather than build one. this function queries a single
// RPC endpoint anyway; it doesn't need to build the client.
func IsValidator(client rpchttp.HTTP, account, prefix string) (validator string, err error) {
	// client, err := NewClientFromChainInfo(info, chain)
	// if err != nil {
	// 	return
	// }
	// Check if the account is also a validator
	_, b64, err := bech32.DecodeAndConvert(account)
	if err != nil {
		return
	}
	// accountsMux.Lock()
	// // FIXME remove Prefixes and replace with chainInfo
	// prefix := Prefixes[chain]
	// accountsMux.Unlock()
	addr, _ := bech32.ConvertAndEncode(prefix+"valoper", b64)
	valQ := staketypes.QueryValidatorRequest{ValidatorAddr: addr}
	valQuery, err := valQ.Marshal()
	if err != nil {
		return
	}
	valResult, err := client.ABCIQuery(context.Background(), "/cosmos.staking.v1beta1.Query/Validator", valQuery)
	if err != nil {
		return
	}
	if len(valResult.Response.Value) > 0 {
		valResp := staketypes.QueryValidatorResponse{}
		err = valResp.Unmarshal(valResult.Response.Value)
		if err != nil {
			return
		}
		validator = valResp.Validator.GetMoniker()
		//fmt.Println(valResp)

	}
	return
}

func QueryAccountFromChainInfo(client rpchttp.HTTP, info *findaccounttypes.ChainInfo, account string) (hasBalance bool, balances string, err error) {
	q := banktypes.QueryBalanceRequest{Address: account}
	var query []byte
	query, err = q.Marshal()
	if err != nil {
		err = fmt.Errorf("Could not marshal QueryBalanceRequest: %w", err)
		return
	}
	result, err := client.ABCIQuery(context.Background(), "/cosmos.bank.v1beta1.Query/AllBalances", query)
	if err != nil {
		err = fmt.Errorf("Could not complete ABCIQuery: %w", err)
		return
	}

	if len(result.Response.Value) > 0 {
		balResp := banktypes.QueryBalanceResponse{}
		err = balResp.Unmarshal(result.Response.Value)
		if err != nil {
			err = fmt.Errorf("Could not unmarshal QueryBalanceResponse: %w", err)
			return
		}
		if balResp.Balance != nil {
			balances = balResp.Balance.String()
			hasBalance = true
		}
	}

	return
}

// TODO change to accept a client as parameter rather than build one. this function queries a single
// RPC endpoint anyway; it doesn't need to build the client.
func QueryAccount(client rpchttp.HTTP, account string) (hasBalance bool, balances string, err error) {

	q := banktypes.QueryBalanceRequest{Address: account}
	var query []byte
	query, err = q.Marshal()
	if err != nil {
		err = fmt.Errorf("Could not marshal QueryBalanceRequest: %w", err)
		return
	}
	result, err := client.ABCIQuery(context.Background(), "/cosmos.bank.v1beta1.Query/AllBalances", query)
	if err != nil {
		err = fmt.Errorf("Could not complete ABCIQuery: %w", err)
		return
	}

	if len(result.Response.Value) > 0 {
		balResp := banktypes.QueryBalanceResponse{}
		err = balResp.Unmarshal(result.Response.Value)
		if err != nil {
			err = fmt.Errorf("Could not unmarshal QueryBalanceResponse: %w", err)
			return
		}
		if balResp.Balance != nil {
			balances = balResp.Balance.String()
			hasBalance = true
		}
	}

	return
}
