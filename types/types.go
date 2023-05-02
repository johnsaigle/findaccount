package types

type ChainInfo struct {
	Apis struct {
		Rpc []Rpc `json:"rpc"`
	} `json:"apis"`
	Bech32Prefix string `json:"bech32_prefix"`
	Explorers []Explorer `json:"explorers"`
}

type Rpc struct {
	Address string `json:"address"`
}

type Explorer struct {
	Url string `json:"url"`
}

