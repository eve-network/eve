package config

func GetCosmosHubConfig() *ChainClientConfig {
	return &ChainClientConfig{
		Key:           "default",
		ChainID:       "cosmoshub-4",
		GRPCAddr:      "grpc-cosmoshub-ia.cosmosia.notional.ventures:443",
		AccountPrefix: "cosmos",
		CoinId:        "cosmos",
		Percent:       0,
		RPC:           "https://cosmos-rpc.polkachu.com",
	}
}

func GetComposableConfig() *ChainClientConfig {
	return &ChainClientConfig{
		Key:           "default",
		ChainID:       "centauri-1",
		GRPCAddr:      "composable-grpc.polkachu.com:22290",
		AccountPrefix: "centauri",
		Percent:       13,
		CoinId:        "picasso",
		RPC:           "https://composable-rpc.polkachu.com",
	}
}

func GetCelestiaConfig() *ChainClientConfig {
	return &ChainClientConfig{
		Key:           "default",
		ChainID:       "celestia",
		GRPCAddr:      "celestia-grpc.polkachu.com:11690",
		AccountPrefix: "celestia",
		Percent:       13,
		CoinId:        "celestia",
		RPC:           "https://celestia-rpc.polkachu.com",
		API:           "https://celestia-api.polkachu.com",
	}
}

func GetSentinelConfig() *ChainClientConfig {
	return &ChainClientConfig{
		Key:           "default",
		ChainID:       "sentinelhub-2",
		GRPCAddr:      "sentinel-grpc.polkachu.com:23990",
		AccountPrefix: "sent",
		Percent:       13,
		CoinId:        "sentinel",
		RPC:           "https://sentinel-rpc.polkachu.com",
	}
}

func GetAkashConfig() *ChainClientConfig {
	return &ChainClientConfig{
		Key:           "default",
		ChainID:       "akashnet-2",
		GRPCAddr:      "akash-grpc.polkachu.com:12890",
		AccountPrefix: "akash",
		Percent:       13,
		CoinId:        "akash-network",
		RPC:           "https://akash-rpc.polkachu.com",
	}
}

func GetStargazeConfig() *ChainClientConfig {
	return &ChainClientConfig{
		Key:           "default",
		ChainID:       "stargaze-1",
		GRPCAddr:      "stargaze-grpc.polkachu.com:13790",
		AccountPrefix: "stars",
		Percent:       13,
		CoinId:        "stargaze",
		RPC:           "https://stargaze-rpc.polkachu.com",
	}
}

func GetNeutronConfig() *ChainClientConfig {
	return &ChainClientConfig{
		Key:           "default",
		ChainID:       "neutron-1",
		GRPCAddr:      "neutron-grpc.polkachu.com:19190",
		AccountPrefix: "neutron",
		Percent:       13,
		CoinId:        "neutron-3",
		RPC:           "https://neutron-rpc.polkachu.com",
	}
}

func GetTerraConfig() *ChainClientConfig {
	return &ChainClientConfig{
		Key:           "default",
		ChainID:       "phoenix-1",
		GRPCAddr:      "terra-grpc.polkachu.com:11790",
		AccountPrefix: "terra",
		Percent:       13,
		CoinId:        "terra-luna-2",
		RPC:           "https://terra-rpc.polkachu.com",
	}
}

func GetBostromConfig() *ChainClientConfig {
	return &ChainClientConfig{
		Key:           "default",
		ChainID:       "bostrom",
		GRPCAddr:      "grpc.cyber.bronbro.io:443",
		AccountPrefix: "bostrom",
		Percent:       13,
		CoinId:        "bostrom",
		API:           "https://lcd.bostrom.cybernode.ai",
	}
}
