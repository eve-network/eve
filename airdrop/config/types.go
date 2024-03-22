package config

import (
	"encoding/json"

	"cosmossdk.io/math"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
)

type ChainClientConfig struct {
	Key           string `json:"key" yaml:"key"`
	ChainID       string `json:"chain-id" yaml:"chain-id"`
	GRPCAddr      string `json:"grpc-addr" yaml:"grpc-addr"`
	AccountPrefix string `json:"account-prefix" yaml:"account-prefix"`
	Percent       int    `json:"percent" yaml:"percent"`
	CoinId        string `json:"coin-id" yaml:"coin-id"`
	RPC           string `json:"rpc" yaml:"rpc"`
	API           string `json:"api" yaml:"api"`
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

type ValidatorResponse struct {
	Validators []Validator `json:"validators"`
	Pagination Pagination  `json:"pagination"`
}

type Pagination struct {
	// next_key is the key to be passed to PageRequest.key to
	// query the next page most efficiently
	NextKey []byte `protobuf:"bytes,1,opt,name=next_key,json=nextKey,proto3" json:"next_key,omitempty"`
	// total is total number of results available if PageRequest.count_total
	// was set, its value is undefined otherwise
	Total string `protobuf:"varint,2,opt,name=total,proto3" json:"total,omitempty"`
}

type Validator struct {
	// operator_address defines the address of the validator's operator; bech encoded in JSON.
	OperatorAddress string `protobuf:"bytes,1,opt,name=operator_address,json=operatorAddress,proto3" json:"operator_address,omitempty" yaml:"operator_address"`
	// tokens define the delegated tokens (incl. self-delegation).
	Tokens math.Int `protobuf:"bytes,5,opt,name=tokens,proto3,customtype=github.com/cosmos/cosmos-sdk/types.Int" json:"tokens"`
	// delegator_shares defines total shares issued to a validator's delegators.
	DelegatorShares math.LegacyDec `protobuf:"bytes,6,opt,name=delegator_shares,json=delegatorShares,proto3,customtype=github.com/cosmos/cosmos-sdk/types.Dec" json:"delegator_shares" yaml:"delegator_shares"`
}

type QueryValidatorDelegationsResponse struct {
	DelegationResponses stakingtypes.DelegationResponses `protobuf:"bytes,1,rep,name=delegation_responses,json=delegationResponses,proto3,castrepeated=DelegationResponses" json:"delegation_responses"`
	// pagination defines the pagination in the response.
	Pagination Pagination `protobuf:"bytes,2,opt,name=pagination,proto3" json:"pagination,omitempty"`
}
