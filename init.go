package findaccount

import (
	"embed"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"os"
	"strings"
)

var (
	//go:embed chain-registry/*/*.json
	chainsFs embed.FS

	//go:embed static/*
	StaticFs embed.FS

	infos = make(map[string]*ChainInfo)
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
		chainInfo := &ChainInfo{}
		e = json.Unmarshal(b, chainInfo)
		if e != nil {
			log.Println(e)
			continue
		}
		if chainInfo != nil && len(chainInfo.Apis.Rpc) > 0 {
			infos[name] = chainInfo
		}
	}

	// add extra known-good RPC servers....
	for k, v := range additional {
		if infos[k] == nil {
			log.Println(k, "is not defined skipping addition of RPC")
			continue
		}
		if infos[k].Apis.Rpc == nil {
			infos[k].Apis.Rpc = make([]Rpc, 0)
		}
		for _, node := range v {
			infos[k].Apis.Rpc = append(infos[k].Apis.Rpc, Rpc{Address: node})
		}
	}

}
