package chaininfo

import (
	"embed"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"os"
	"strings"
	"github.com/johnsaigle/findaccount/types"
)


var (
	//go:embed chain-registry/*/*.json
	chainsFs embed.FS

	// //go:embed static/*
	// StaticFs embed.FS

	Infos = make(map[string]*types.ChainInfo)
)

// TODO I would like to remove the init function because I don't know if there is a good way to do error handling
// within it. Temporary compromise is to panic but this is risky if other projects use this package.
func init() {
	log.SetOutput(os.Stderr)
	log.SetFlags(log.Lshortfile)

	// The chain-registry directory is a submodule to https://github.com/cosmos/chain-registry/
	registryFiles, err := chainsFs.ReadDir("chain-registry")
	if err != nil {
		panic("Could not read chain-registry directory. No way to recover")
	}

	for _, entry := range registryFiles {
		// We want directories that do not start with an underscore or period 
		if !entry.IsDir() {
			continue
		} 
		name := entry.Name()
		// NOTE: embedFS regex pattern above might take care of this for us..
		if strings.HasPrefix(name, "_") || strings.HasPrefix(name, ".") {
			continue
		}

		f, e := chainsFs.Open(fmt.Sprintf("chain-registry/%s/chain.json", name))
		if e != nil {
			log.Println(e)
			continue
		}
		b, e := io.ReadAll(f)
		if e != nil {
			log.Println(e)
			continue
		}
		defer f.Close()
		chainInfo := &types.ChainInfo{}
		e = json.Unmarshal(b, chainInfo)
		if e != nil {
			log.Println(e)
			continue
		}
		if chainInfo != nil && len(chainInfo.Apis.Rpc) > 0 {
			Infos[name] = chainInfo
		}
	}

	// add extra known-good RPC servers....
	for name, prefix := range additional {
		if Infos[name] == nil {
			log.Println(name, "is not defined skipping addition of RPC")
			continue
		}
		if Infos[name].Apis.Rpc == nil {
			Infos[name].Apis.Rpc = make([]types.Rpc, 0)
		}
		for _, node := range prefix {
			Infos[name].Apis.Rpc = append(Infos[name].Apis.Rpc, types.Rpc{Address: node})
		}
	}

}

// Additional known nodes, not all nodes from cosmos repo are live....

var additional = map[string][]string{
	"secretnetwork": {"tcp://scrt-rpc.blockpane.com:26657"},
	"osmosis":       {"https://osmosis-rpc.polkachu.com:443"},
	"chihuahua":     {"https://chihuahua-rpc.mercury-nodes.net:443"},
	"emoney":        {"https://emoney.validator.network:443"},
	"kava":          {"https://rpc.data.kava.io:443"},
	"stargaze":      {"https://rpc.stargaze-apis.com:443"},
	"juno":          {"https://juno-rpc.polkachu.com:443"},
}

// Prefixes maps the chain name to the bech32 address prefix.
// TODO: validate these are all correct!!!
// TODO: delete anything that is in the chainlist. move any remaining to 'additional' for now
var Prefixes = map[string]string{
	"agoric":         "agoric",
	"akash":          "akash",
	"assetmantle":    "mantle",
	"axelar":         "axelar",
	"bandchain":      "band",
	"bitcanna":       "bcn",
	"bitsong":        "bitsong",
	"bostrom":        "bostrom",
	"carbon":         "carbon",
	"cerberus":       "cerberus",
	"cheqd":          "cheqd",
	"chihuahua":      "chihuahua",
	"comdex":         "comdex",
	"cosmoshub":      "cosmos",
	"crescent":       "cre",
	"cronos":         "cro",
	"cryptoorgchain": "cro",
	"decentr":        "decentr",
	"desmos":         "desmos",
	"dig":            "dig",
	"emoney":         "emoney",
	"evmos":          "evmos",
	"fetchhub":       "fetchhub",
	"firmachain":     "firma",
	"galaxy":         "galaxy",
	"genesisl1":      "genesis",
	"gravitybridge":  "gravity",
	"impacthub":      "impact",
	"injective":      "inj",
	"irisnet":        "i",
	"juno":           "juno",
	"kava":           "kava",
	"kichain":        "ki",
	"konstellation":  "darc",
	"likecoin":       "like",
	"meme":           "meme",
	"microtick":      "micro",
	"nomic":          "nomic",
	"octa":           "octa",
	"odin":           "odin",
	"oraichain":      "orai",
	"osmosis":        "osmo",
	"panacea":        "panacea",
	"persistence":    "persistence",
	"provenance":     "pb",
	"regen":          "regen",
	"rizon":          "rizon",
	"secretnetwork":  "secret",
	"sentinel":       "sent",
	"shentu":         "certic",
	"sifchain":       "sif",
	"sommelier":      "somm",
	"stargaze":       "stars",
	"starname":       "star",
	"terra":          "terra",
	"umee":           "umee",
	"vidulum":        "vdl",
	// not working!!!
	//"arkh":           "arkh", //this is a testnet?
	//"echelon":        "echelon",
	//"logos":          "logos",
	//"lumnetwork":     "lum",
	//"mythos":         "mythos",
	//"thorchain":      "thor",
}
