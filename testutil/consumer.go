package testutil

import (
	"time"

	ibcclienttypes "github.com/cosmos/ibc-go/v8/modules/core/02-client/types" //nolint:staticcheck
	ibccommitmenttypes "github.com/cosmos/ibc-go/v8/modules/core/23-commitment/types"
	ibctmtypes "github.com/cosmos/ibc-go/v8/modules/light-clients/07-tendermint"
	ccvconsumertypes "github.com/cosmos/interchain-security/v5/x/ccv/consumer/types"
	ccvprovidertypes "github.com/cosmos/interchain-security/v5/x/ccv/provider/types"
	"github.com/cosmos/interchain-security/v5/x/ccv/types"
)

func CreateMinimalConsumerTestGenesis() *ccvconsumertypes.GenesisState {
	genesisState := ccvconsumertypes.DefaultGenesisState()
	genesisState.Params.Enabled = true
	genesisState.NewChain = true
	genesisState.Provider.ClientState = ccvprovidertypes.DefaultParams().TemplateClient
	genesisState.Provider.ClientState.ChainId = "eve"
	genesisState.Provider.ClientState.LatestHeight = ibcclienttypes.Height{RevisionNumber: 0, RevisionHeight: 1}
	trustPeriod, err := types.CalculateTrustPeriod(genesisState.Params.UnbondingPeriod, ccvprovidertypes.DefaultTrustingPeriodFraction)
	if err != nil {
		panic("provider client trusting period error")
	}
	genesisState.Provider.ClientState.TrustingPeriod = trustPeriod
	genesisState.Provider.ClientState.UnbondingPeriod = genesisState.Params.UnbondingPeriod
	genesisState.Provider.ClientState.MaxClockDrift = ccvprovidertypes.DefaultMaxClockDrift
	genesisState.Provider.ConsensusState = &ibctmtypes.ConsensusState{
		Timestamp: time.Now().UTC(),
		Root:      ibccommitmenttypes.MerkleRoot{Hash: []byte("dummy")},
	}

	return genesisState
}
