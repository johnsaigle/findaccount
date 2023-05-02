package findaccount

import (
	"fmt"
	"log"
	"sort"
	"sync"

	"github.com/cosmos/cosmos-sdk/types/bech32"
	"github.com/johnsaigle/findaccount/pkg/chaininfo"
	"github.com/johnsaigle/findaccount/pkg/client"
)

var accountsMux sync.Mutex
var infos = chaininfo.Infos //populated by init code when the script gets run

type ChainResult struct {
	Chain      string `json:"chain"`
	Address    string `json:"address"`
	Validator  string `json:"is_validator"`
	HasBalance bool   `json:"hasBalance"`
	Coins      string `json:"coins"`
	Error      string `json:"error"`
	Link       string `json:"link"`
}

func (r ChainResult) CsvHeader() string {
	return "chain,address,validator,has balance,coins,error"
}

func (r ChainResult) ToCsv() string {
	return fmt.Sprintf("%s,%s,%q,%v,%s,%s", r.Chain, r.Address, r.Validator, r.HasBalance, r.Coins, r.Error)
}

// SearchAccounts is the entrypoint for performing a search
func SearchAccounts(account, rpc, prefix string) ([]ChainResult, error) {
	// TODO : validate rpc and prefix
	// i.e. if prefix is not alphanumeric
	// i.e. if rpc is not well-formed (may need a URL-parsing library
	results := make([]ChainResult, 0)

	// override chain-registry infos if custom rpc provided
	if rpc != "" && prefix != "" {
		panic("Not implemented!")
	}

	addrMap, err := ConvertToAccounts(account)
	if err != nil {
		return results, err
	}

	wg := &sync.WaitGroup{}
	wg.Add(len(infos))
	for k, v := range infos {
		var link string

		accountsMux.Lock()
		chain, rpcs := k, v
		addr := addrMap[k]
		if len(infos[k].Explorers) > 0 {
			link = infos[k].Explorers[0].Url
		}
		accountsMux.Unlock()

		go func() {
			bal, coins, e := client.QueryAccount(rpcs, chain, addr)
			errStr := "ok"
			if e != nil {
				errStr = e.Error()
			}
			val, _ := client.IsValidator(rpcs, chain, addr)
			results = append(results, ChainResult{
				Chain:      chain,
				Address:    addr,
				Validator:  val,
				HasBalance: bal,
				Coins:      coins,
				Error:      errStr,
				Link:       link,
			})
			wg.Done()
		}()
	}
	wg.Wait()

	sort.Slice(results, func(i, j int) bool {
		return sort.StringsAreSorted([]string{results[i].Chain, results[j].Chain})
	})

	return results, err
}

// Takes a string that should be a bech32 address. Returns error if it isn't 
// Extracts the bytes that represent the actual address (without HRP and checksum)
// Iterates over the ChainInfo struct to obtain all bech32 prefixes extract from the chain-registry.
// Encode the address bytes using all bech32 prefixes
// Returns a mapping of chain names to generated addresses
func ConvertToAccounts(s string) (map[string]string, error) {
	accounts := make(map[string]string)
	_, b64, err := bech32.DecodeAndConvert(s)

	if err != nil {
		return nil, err
	}

	for name, chainInfo := range infos {
		addr, e := bech32.ConvertAndEncode(chainInfo.Bech32Prefix, b64)
		if e != nil {
			log.Println(name, e)
		}
		accounts[name] = addr
	}

	return accounts, nil
}
