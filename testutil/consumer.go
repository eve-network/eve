package testutil

import (
	abci "github.com/cometbft/cometbft/abci/types"
	tmtypes "github.com/cometbft/cometbft/types"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	"github.com/cosmos/interchain-security/v4/testutil/crypto"
	ccv "github.com/cosmos/interchain-security/v4/x/ccv/types"
	"time"

	clienttypes "github.com/cosmos/ibc-go/v8/modules/core/02-client/types"
	commitmenttypes "github.com/cosmos/ibc-go/v8/modules/core/23-commitment/types"

	ibctmtypes "github.com/cosmos/ibc-go/v8/modules/light-clients/07-tendermint"
	consumertypes "github.com/cosmos/interchain-security/v4/x/ccv/consumer/types"
)

func CreateMinimalConsumerTestGenesis() *consumertypes.GenesisState {
	// create validator set
	cId := crypto.NewCryptoIdentityFromIntSeed(234234)
	pubKey := cId.TMCryptoPubKey()
	validator := tmtypes.NewValidator(pubKey, 1)
	valset := []abci.ValidatorUpdate{tmtypes.TM2PB.ValidatorUpdate(validator)}

	// create ibc client and last consensus states
	provConsState := ibctmtypes.NewConsensusState(
		time.Now(),
		commitmenttypes.NewMerkleRoot([]byte("apphash")),
		tmtypes.NewValidatorSet([]*tmtypes.Validator{validator}).Hash(),
	)

	provClientState := ibctmtypes.NewClientState(
		"provider",
		ibctmtypes.DefaultTrustLevel,
		1,
		stakingtypes.DefaultUnbondingTime,
		time.Second*10,
		clienttypes.Height{RevisionNumber: 0, RevisionHeight: 1},
		commitmenttypes.GetSDKSpecs(),
		[]string{"upgrade", "upgradedIBCState"},
	)

	// create default parameters for a new chain
	params := ccv.DefaultParams()
	params.Enabled = true
	state := consumertypes.NewInitialGenesisState(
		provClientState,
		provConsState,
		valset,
		params,
	)
	state.ProviderConsensusState = provConsState
	state.ProviderClientState = provClientState

	return state
}
