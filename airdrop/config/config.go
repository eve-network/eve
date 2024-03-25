package config

func GetCosmosHubConfig() *ChainClientConfig {
	return &ChainClientConfig{
		Key:           "default",
		ChainID:       "cosmoshub-4",
		GRPCAddr:      "cosmos-grpc.publicnode.com:443",
		AccountPrefix: "cosmos",
		CoinId:        "cosmos",
		Percent:       10,
		RPC:           "https://cosmos-rpc.polkachu.com",
		API:           "https://cosmos-rest.publicnode.com",
	}
}

func GetComposableConfig() *ChainClientConfig {
	return &ChainClientConfig{
		Key:           "default",
		ChainID:       "centauri-1",
		GRPCAddr:      "composable-grpc.polkachu.com:22290",
		AccountPrefix: "centauri",
		Percent:       10,
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
		Percent:       10,
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
		Percent:       10,
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
		Percent:       10,
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
		Percent:       10,
		CoinId:        "stargaze",
		RPC:           "https://stargaze-rpc.polkachu.com",
		API:           "https://rest.stargaze-apis.com",
	}
}

func GetNeutronConfig() *ChainClientConfig {
	return &ChainClientConfig{
		Key:           "default",
		ChainID:       "neutron-1",
		GRPCAddr:      "neutron-grpc.polkachu.com:19190",
		AccountPrefix: "neutron",
		Percent:       10,
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
		Percent:       10,
		CoinId:        "terra-luna-2",
		RPC:           "https://terra-rpc.polkachu.com",
		API:           "https://terra-rest.publicnode.com",
	}
}

func GetBostromConfig() *ChainClientConfig {
	return &ChainClientConfig{
		Key:           "default",
		ChainID:       "bostrom",
		GRPCAddr:      "grpc.cyber.bronbro.io:443",
		AccountPrefix: "bostrom",
		Percent:       10,
		CoinId:        "bostrom",
		API:           "https://lcd.bostrom.cybernode.ai",
	}
}

func GetTerracConfig() *ChainClientConfig {
	return &ChainClientConfig{
		Key:           "default",
		ChainID:       "columbus-5",
		GRPCAddr:      "terra-classic-grpc.publicnode.com:443",
		AccountPrefix: "terra",
		Percent:       10,
		CoinId:        "terra-luna",
		API:           "https://terra-classic-lcd.publicnode.com",
	}
}
