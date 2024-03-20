package config

func GetCosmosHubConfig() *ChainClientConfig {
	return &ChainClientConfig{
		Key:           "default",
		ChainID:       "cosmoshub-4",
		GRPCAddr:      "grpc-cosmoshub-ia.cosmosia.notional.ventures:443",
		AccountPrefix: "cosmos",
		CoinId:        "cosmos",
		Percent:       0,
		NodeStatusUrl: "https://cosmos-rpc.polkachu.com/status",
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
		NodeStatusUrl: "https://composable-rpc.polkachu.com/status",
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
		NodeStatusUrl: "https://celestia-rpc.polkachu.com/status",
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
		NodeStatusUrl: "https://sentinel-rpc.polkachu.com/status",
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
		NodeStatusUrl: "https://akash-rpc.polkachu.com/status",
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
		NodeStatusUrl: "https://stargaze-rpc.polkachu.com/status",
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
		NodeStatusUrl: "https://neutron-rpc.polkachu.com/status",
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
		NodeStatusUrl: "https://terra-rpc.polkachu.com/status",
	}
}

func GetBostromConfig() *ChainClientConfig {
	return &ChainClientConfig{
		Key:           "default",
		ChainID:       "bostrom",
		GRPCAddr:      "grpc-cyber-ia.cosmosia.notional.ventures:443",
		AccountPrefix: "bostrom",
		Percent:       13,
		CoinId:        "bostrom",
		NodeStatusUrl: "https://rpc-cyber-ia.cosmosia.notional.ventures/status",
	}
}
