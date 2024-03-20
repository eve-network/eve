package config

import (
	"encoding/json"

	"cosmossdk.io/math"
)

type ChainClientConfig struct {
	Key           string `json:"key" yaml:"key"`
	ChainID       string `json:"chain-id" yaml:"chain-id"`
	GRPCAddr      string `json:"grpc-addr" yaml:"grpc-addr"`
	AccountPrefix string `json:"account-prefix" yaml:"account-prefix"`
	Percent       int    `json:"percent" yaml:"percent"`
	CoinId        string `json:"coin-id" yaml:"coin-id"`
	NodeStatusUrl string `json:"node-status" yaml:"node-status"`
}

type Reward struct {
	Address         string         `json:"address"`
	EveAddress      string         `json:"eve_address"`
	Shares          math.LegacyDec `json:"shares"`
	Token           math.LegacyDec `json:"tokens"`
	EveAirdropToken math.LegacyDec `json:"eve"`
	ChainId         string         `json:"chain"`
}

type ComposablePrice struct {
	Token Price `json:"picasso"`
}

type AkashPrice struct {
	Token Price `json:"akash-network"`
}

type CelestiaPrice struct {
	Token Price `json:"celestia"`
}

type CosmosPrice struct {
	Token Price `json:"cosmos"`
}

type NeutronPrice struct {
	Token Price `json:"neutron-3"`
}

type SentinelPrice struct {
	Token Price `json:"sentinel"`
}

type StargazePrice struct {
	Token Price `json:"stargaze"`
}

type TerraPrice struct {
	Token Price `json:"terra-luna-2"`
}

type BostromPrice struct {
	Token Price `json:"bostrom"`
}

type Price struct {
	USD json.Number `json:"usd"`
}

type NodeResponse struct {
	Id      json.Number `json:"id"`
	JsonRPC string      `json:"jsonrpc"`
	Result  Result      `json:"result"`
}

type Result struct {
	NodeInfo      NodeInfo      `json:"node_info"`
	SyncInfo      SyncInfo      `json:"sync_info"`
	ValidatorInfo ValidatorInfo `json:"validator_info"`
}

type NodeInfo struct{}
type SyncInfo struct {
	CatchingUp           bool   `json:"catching_up"`
	EarlieastAppHash     string `json:"earliest_app_hash"`
	EarlieastBlockHash   string `json:"earliest_block_hash"`
	EarlieastBlockHeight string `json:"earliest_block_height"`
	EarlieastBlockTime   string `json:"earliest_block_time"`
	LatestAppHash        string `json:"latest_app_hash"`
	LatestBlockHash      string `json:"latest_block_hash"`
	LatestBlockHeight    string `json:"latest_block_height"`
	LatestBlockTime      string `json:"latest_block_time"`
}
type ValidatorInfo struct{}
