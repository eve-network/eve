package app

import (
	"flag"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/address"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"

	"github.com/cosmos/cosmos-sdk/types/simulation"
)

const (
	HumanCoinUnit = "eve"
	BaseCoinUnit  = "ueve"
	UeveExponent  = 6

	DefaultBondDenom = BaseCoinUnit

	// Bech32PrefixAccAddr defines the Bech32 prefix of an account's address.
	AccountAddressPrefix = "eve"
)

var (
	// PREFIXES
	// Bech32PrefixAccAddr defines the Bech32 prefix of an account's address.
	Bech32PrefixAccAddr = AccountAddressPrefix
	// Bech32PrefixAccPub defines the Bech32 prefix of an account's public key.
	Bech32PrefixAccPub = AccountAddressPrefix + sdk.PrefixPublic
	// Bech32PrefixValAddr defines the Bech32 prefix of a validator's operator address.
	Bech32PrefixValAddr = AccountAddressPrefix + sdk.PrefixValidator + sdk.PrefixOperator
	// Bech32PrefixValPub defines the Bech32 prefix of a validator's operator public key.
	Bech32PrefixValPub = AccountAddressPrefix + sdk.PrefixValidator + sdk.PrefixOperator + sdk.PrefixPublic
	// Bech32PrefixConsAddr defines the Bech32 prefix of a consensus node address.
	Bech32PrefixConsAddr = AccountAddressPrefix + sdk.PrefixValidator + sdk.PrefixConsensus
	// Bech32PrefixConsPub defines the Bech32 prefix of a consensus node public key.
	Bech32PrefixConsPub = AccountAddressPrefix + sdk.PrefixValidator + sdk.PrefixConsensus + sdk.PrefixPublic
)

// List of available flags for the simulator
var (
	FlagGenesisFileValue        string
	FlagParamsFileValue         string
	FlagExportParamsPathValue   string
	FlagExportParamsHeightValue int
	FlagExportStatePathValue    string
	FlagExportStatsPathValue    string
	FlagSeedValue               int64
	FlagInitialBlockHeightValue int
	FlagNumBlocksValue          int
	FlagBlockSizeValue          int
	FlagLeanValue               bool
	FlagCommitValue             bool
	FlagOnOperationValue        bool // TODO Remove in favor of binary search for invariant violation
	FlagAllInvariantsValue      bool
	FlagDBBackendValue          string

	FlagEnabledValue     bool
	FlagVerboseValue     bool
	FlagPeriodValue      uint
	FlagGenesisTimeValue int64
)

// GetSimulatorFlags gets the values of all the available simulation flags
func GetSimulatorFlags() {
	// config fields
	flag.StringVar(&FlagGenesisFileValue, "Genesis", "", "custom simulation genesis file; cannot be used with params file")
	flag.StringVar(&FlagParamsFileValue, "Params", "", "custom simulation params file which overrides any random params; cannot be used with genesis")
	flag.StringVar(&FlagExportParamsPathValue, "ExportParamsPath", "", "custom file path to save the exported params JSON")
	flag.IntVar(&FlagExportParamsHeightValue, "ExportParamsHeight", 0, "height to which export the randomly generated params")
	flag.StringVar(&FlagExportStatePathValue, "ExportStatePath", "", "custom file path to save the exported app state JSON")
	flag.StringVar(&FlagExportStatsPathValue, "ExportStatsPath", "", "custom file path to save the exported simulation statistics JSON")
	flag.Int64Var(&FlagSeedValue, "Seed", 42, "simulation random seed")
	flag.IntVar(&FlagInitialBlockHeightValue, "InitialBlockHeight", 1, "initial block to start the simulation")
	flag.IntVar(&FlagNumBlocksValue, "NumBlocks", 500, "number of new blocks to simulate from the initial block height")
	flag.IntVar(&FlagBlockSizeValue, "BlockSize", 200, "operations per block")
	flag.BoolVar(&FlagLeanValue, "Lean", false, "lean simulation log output")
	flag.BoolVar(&FlagCommitValue, "Commit", false, "have the simulation commit")
	flag.BoolVar(&FlagOnOperationValue, "SimulateEveryOperation", false, "run slow invariants every operation")
	flag.BoolVar(&FlagAllInvariantsValue, "PrintAllInvariants", false, "print all invariants if a broken invariant is found")
	flag.StringVar(&FlagDBBackendValue, "DBBackend", "pebbledb", "custom db backend type")

	// simulation flags
	flag.BoolVar(&FlagEnabledValue, "Enabled", false, "enable the simulation")
	flag.BoolVar(&FlagVerboseValue, "Verbose", false, "verbose log output")
	flag.UintVar(&FlagPeriodValue, "Period", 0, "run slow invariants only once every period assertions")
	flag.Int64Var(&FlagGenesisTimeValue, "GenesisTime", 0, "override genesis UNIX time instead of using a random UNIX time")
}

// NewConfigFromFlags creates a simulation from the retrieved values of the flags.
func NewConfigFromFlags() simulation.Config {
	return simulation.Config{
		GenesisFile:        FlagGenesisFileValue,
		ParamsFile:         FlagParamsFileValue,
		ExportParamsPath:   FlagExportParamsPathValue,
		ExportParamsHeight: FlagExportParamsHeightValue,
		ExportStatePath:    FlagExportStatePathValue,
		ExportStatsPath:    FlagExportStatsPathValue,
		Seed:               FlagSeedValue,
		InitialBlockHeight: FlagInitialBlockHeightValue,
		NumBlocks:          FlagNumBlocksValue,
		BlockSize:          FlagBlockSizeValue,
		Lean:               FlagLeanValue,
		Commit:             FlagCommitValue,
		OnOperation:        FlagOnOperationValue,
		AllInvariants:      FlagAllInvariantsValue,
		DBBackend:          FlagDBBackendValue,
	}
}

// https://github.com/sei-protocol/sei-chain/blob/b52d2f472e520bcf8bd8024d7b6090b4f14dbc64/app/params/config.go#L53
func SetAddressPrefixes() {
	config := sdk.GetConfig()
	config.SetBech32PrefixForAccount(Bech32PrefixAccAddr, Bech32PrefixAccPub)
	config.SetBech32PrefixForValidator(Bech32PrefixValAddr, Bech32PrefixValPub)
	config.SetBech32PrefixForConsensusNode(Bech32PrefixConsAddr, Bech32PrefixConsPub)

	// This is copied from the cosmos sdk v0.43.0-beta1
	// source: https://github.com/cosmos/cosmos-sdk/blob/v0.43.0-beta1/types/address.go#L141
	config.SetAddressVerifier(func(bytes []byte) error {
		if len(bytes) == 0 {
			return sdkerrors.Wrap(sdkerrors.ErrUnknownAddress, "addresses cannot be empty")
		}

		if len(bytes) > address.MaxAddrLen {
			return sdkerrors.Wrapf(sdkerrors.ErrUnknownAddress, "address max length is %d, got %d", address.MaxAddrLen, len(bytes))
		}

		// TODO: Do we want to allow addresses of lengths other than 20 and 32 bytes?
		if len(bytes) != 20 && len(bytes) != 32 {
			return sdkerrors.Wrapf(sdkerrors.ErrUnknownAddress, "address length must be 20 or 32 bytes, got %d", len(bytes))
		}

		return nil
	})
}
