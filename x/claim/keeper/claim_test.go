package keeper_test

import (
	"time"

	"github.com/cosmos/cosmos-sdk/crypto/keys/secp256k1"
	sdk "github.com/cosmos/cosmos-sdk/types"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	sdkstakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	stakingkeeper "github.com/iqlusioninc/liquidity-staking-module/x/staking/keeper"
	stakingtypes "github.com/iqlusioninc/liquidity-staking-module/x/staking/types"

	"github.com/eve-network/eve/x/claim/types"
)

func (suite *KeeperTestSuite) TestHookOfUnclaimableAccount() {
	pub1 := secp256k1.GenPrivKey().PubKey()
	addr1 := sdk.AccAddress(pub1.Address())
	suite.app.AccountKeeper.SetAccount(suite.ctx, authtypes.NewBaseAccount(addr1, nil, 0, 0))

	claim, err := suite.app.ClaimKeeper.GetClaimRecord(suite.ctx, addr1)
	suite.NoError(err)
	suite.Equal(types.ClaimRecord{}, claim)

	suite.app.ClaimKeeper.AfterDelegationModified(suite.ctx, addr1, sdk.ValAddress(addr1))

	balances := suite.app.BankKeeper.GetAllBalances(suite.ctx, addr1)
	suite.Equal(sdk.Coins{}, balances)
}

func (suite *KeeperTestSuite) TestHookBeforeAirdropStart() {
	suite.SetupTest()

	airdropStartTime := time.Now().Add(time.Hour)

	suite.app.ClaimKeeper.SetParams(suite.ctx, types.Params{
		AirdropEnabled:     true,
		ClaimDenom:         types.DefaultClaimDenom,
		AirdropStartTime:   airdropStartTime,
		DurationUntilDecay: time.Hour,
		DurationOfDecay:    time.Hour * 4,
	})

	pub1 := secp256k1.GenPrivKey().PubKey()
	addr1 := sdk.AccAddress(pub1.Address())

	claimRecords := []types.ClaimRecord{
		{
			Address:                addr1.String(),
			InitialClaimableAmount: sdk.NewCoins(sdk.NewInt64Coin(types.DefaultClaimDenom, 1000)),
			ActionCompleted:        []bool{false, false, false, false, false},
		},
	}
	suite.app.AccountKeeper.SetAccount(suite.ctx, authtypes.NewBaseAccount(addr1, nil, 0, 0))

	err := suite.app.ClaimKeeper.SetClaimRecords(suite.ctx, claimRecords)
	suite.Require().NoError(err)

	coins, err := suite.app.ClaimKeeper.GetUserTotalClaimable(suite.ctx, addr1)
	suite.NoError(err)
	// Now, it is before starting air drop, so this value should return the empty coins
	suite.True(coins.Empty())

	coins, err = suite.app.ClaimKeeper.GetClaimableAmountForAction(suite.ctx, addr1, types.ActionDelegateStake)
	suite.NoError(err)
	// Now, it is before starting air drop, so this value should return the empty coins
	suite.True(coins.Empty())

	suite.app.ClaimKeeper.AfterDelegationModified(suite.ctx, addr1, sdk.ValAddress(addr1))
	balances := suite.app.BankKeeper.GetAllBalances(suite.ctx, addr1)
	// Now, it is before starting air drop, so claim module should not send the balances to the user after delegate.
	suite.True(balances.Empty())

	suite.app.ClaimKeeper.AfterDelegationModified(suite.ctx.WithBlockTime(airdropStartTime), addr1, sdk.ValAddress(addr1))
	balances = suite.app.BankKeeper.GetAllBalances(suite.ctx, addr1)
	// Now, it is the time for air drop, so claim module should send the balances to the user after delegate.
	suite.Equal(claimRecords[0].InitialClaimableAmount.AmountOf(types.DefaultClaimDenom).Quo(sdk.NewInt(int64(len(types.Action_value)))), balances.AmountOf(types.DefaultClaimDenom))
}

func (suite *KeeperTestSuite) TestAirdropDisabled() {
	suite.SetupTest()

	airdropStartTime := time.Now().Add(time.Hour)

	suite.app.ClaimKeeper.SetParams(suite.ctx, types.Params{
		AirdropEnabled:     false,
		ClaimDenom:         types.DefaultClaimDenom,
		DurationUntilDecay: time.Hour,
		DurationOfDecay:    time.Hour * 4,
	})

	pub1 := secp256k1.GenPrivKey().PubKey()
	addr1 := sdk.AccAddress(pub1.Address())

	claimRecords := []types.ClaimRecord{
		{
			Address:                addr1.String(),
			InitialClaimableAmount: sdk.NewCoins(sdk.NewInt64Coin(types.DefaultClaimDenom, 1000)),
			ActionCompleted:        []bool{false, false, false, false, false},
		},
	}
	suite.app.AccountKeeper.SetAccount(suite.ctx, authtypes.NewBaseAccount(addr1, nil, 0, 0))

	err := suite.app.ClaimKeeper.SetClaimRecords(suite.ctx, claimRecords)
	suite.Require().NoError(err)

	coins, err := suite.app.ClaimKeeper.GetUserTotalClaimable(suite.ctx, addr1)
	suite.NoError(err)
	// Now, it is before starting air drop, so this value should return the empty coins
	suite.True(coins.Empty())

	coins, err = suite.app.ClaimKeeper.GetClaimableAmountForAction(suite.ctx, addr1, types.ActionDelegateStake)
	suite.NoError(err)
	// Now, it is before starting air drop, so this value should return the empty coins
	suite.True(coins.Empty())

	suite.app.ClaimKeeper.AfterDelegationModified(suite.ctx, addr1, sdk.ValAddress(addr1))
	balances := suite.app.BankKeeper.GetAllBalances(suite.ctx, addr1)
	// Now, it is before starting air drop, so claim module should not send the balances to the user after delegate.
	suite.True(balances.Empty())

	suite.app.ClaimKeeper.AfterProposalVote(suite.ctx, 1, addr1)
	balances = suite.app.BankKeeper.GetAllBalances(suite.ctx, addr1)
	// Now, it is before starting air drop, so claim module should not send the balances to the user after vote.
	suite.True(balances.Empty())

	// set airdrop enabled but with invalid date
	suite.app.ClaimKeeper.SetParams(suite.ctx, types.Params{
		AirdropEnabled:     true,
		ClaimDenom:         types.DefaultClaimDenom,
		DurationUntilDecay: time.Hour,
		DurationOfDecay:    time.Hour * 4,
	})

	suite.app.ClaimKeeper.AfterDelegationModified(suite.ctx, addr1, sdk.ValAddress(addr1))
	balances = suite.app.BankKeeper.GetAllBalances(suite.ctx, addr1)
	// Now airdrop is enabled but a potential misconfiguraion on start time
	suite.True(balances.Empty())

	suite.app.ClaimKeeper.AfterProposalVote(suite.ctx, 1, addr1)
	balances = suite.app.BankKeeper.GetAllBalances(suite.ctx, addr1)
	// Now airdrop is enabled but a potential misconfiguraion on start time, so claim module should not send the balances to the user after vote.
	suite.True(balances.Empty())

	// set airdrop enabled but with date in the future
	suite.app.ClaimKeeper.SetParams(suite.ctx, types.Params{
		AirdropEnabled:     true,
		AirdropStartTime:   airdropStartTime.Add(time.Hour),
		ClaimDenom:         types.DefaultClaimDenom,
		DurationUntilDecay: time.Hour,
		DurationOfDecay:    time.Hour * 4,
	})

	suite.app.ClaimKeeper.AfterDelegationModified(suite.ctx, addr1, sdk.ValAddress(addr1))
	balances = suite.app.BankKeeper.GetAllBalances(suite.ctx, addr1)
	// Now airdrop is enabled  and date is not empty but block time still behid
	suite.True(balances.Empty())

	suite.app.ClaimKeeper.AfterProposalVote(suite.ctx, 1, addr1)
	balances = suite.app.BankKeeper.GetAllBalances(suite.ctx, addr1)
	// Now airdrop is enabled  and date is not empty but block time still behid
	suite.True(balances.Empty())

	// add extra 2 hours
	suite.app.ClaimKeeper.AfterDelegationModified(suite.ctx.WithBlockTime(airdropStartTime.Add(time.Hour*2)), addr1, sdk.ValAddress(addr1))
	balances = suite.app.BankKeeper.GetAllBalances(suite.ctx, addr1)
	// Now, it is the time for air drop, so claim module should send the balances to the user after delegate.
	suite.Equal(claimRecords[0].InitialClaimableAmount.AmountOf(types.DefaultClaimDenom).Quo(sdk.NewInt(int64(len(types.Action_value)))), balances.AmountOf(types.DefaultClaimDenom))
}

func (suite *KeeperTestSuite) TestDuplicatedActionNotWithdrawRepeatedly() {
	suite.SetupTest()

	pub1 := secp256k1.GenPrivKey().PubKey()
	addr1 := sdk.AccAddress(pub1.Address())

	claimRecords := []types.ClaimRecord{
		{
			Address:                addr1.String(),
			InitialClaimableAmount: sdk.NewCoins(sdk.NewInt64Coin(types.DefaultClaimDenom, 2000)),
			ActionCompleted:        []bool{false, false, false, false, false},
		},
	}
	suite.app.AccountKeeper.SetAccount(suite.ctx, authtypes.NewBaseAccount(addr1, nil, 0, 0))

	err := suite.app.ClaimKeeper.SetClaimRecords(suite.ctx, claimRecords)
	suite.Require().NoError(err)

	coins1, err := suite.app.ClaimKeeper.GetUserTotalClaimable(suite.ctx, addr1)
	suite.Require().NoError(err)
	suite.Require().Equal(coins1, claimRecords[0].InitialClaimableAmount)

	suite.app.ClaimKeeper.AfterDelegationModified(suite.ctx, addr1, sdk.ValAddress(addr1))
	claim, err := suite.app.ClaimKeeper.GetClaimRecord(suite.ctx, addr1)
	suite.NoError(err)
	suite.True(claim.ActionCompleted[types.ActionDelegateStake])

	claimedCoins := suite.app.BankKeeper.GetAllBalances(suite.ctx, addr1)
	suite.Require().Equal(claimedCoins.AmountOf(types.DefaultClaimDenom), claimRecords[0].InitialClaimableAmount.AmountOf(types.DefaultClaimDenom).Quo(sdk.NewInt(5)))

	suite.app.ClaimKeeper.AfterDelegationModified(suite.ctx, addr1, sdk.ValAddress(addr1))
	claim, err = suite.app.ClaimKeeper.GetClaimRecord(suite.ctx, addr1)
	suite.NoError(err)
	suite.True(claim.ActionCompleted[types.ActionDelegateStake])

	claimedCoins = suite.app.BankKeeper.GetAllBalances(suite.ctx, addr1)
	suite.Require().Equal(claimedCoins.AmountOf(types.DefaultClaimDenom), claimRecords[0].InitialClaimableAmount.AmountOf(types.DefaultClaimDenom).Quo(sdk.NewInt(5)))
}

func (suite *KeeperTestSuite) TestNotRunningGenesisBlock() {
	suite.ctx = suite.ctx.WithBlockHeight(1)
	suite.app.ClaimKeeper.CreateModuleAccount(suite.ctx, sdk.NewCoin(sdk.DefaultBondDenom, sdk.NewInt(10000000)))
	// set airdrop enabled but with date in the future
	suite.app.ClaimKeeper.SetParams(suite.ctx, types.Params{
		AirdropEnabled:     true,
		AirdropStartTime:   time.Now().Add(time.Hour * -1),
		ClaimDenom:         sdk.DefaultBondDenom,
		DurationUntilDecay: time.Hour,
		DurationOfDecay:    time.Hour * 4,
		AllowedClaimers:    make([]types.ClaimAuthorization, 0),
	})

	pub1 := secp256k1.GenPrivKey().PubKey()
	addr1 := sdk.AccAddress(pub1.Address())

	claimRecords := []types.ClaimRecord{
		{
			Address:                addr1.String(),
			InitialClaimableAmount: sdk.NewCoins(sdk.NewInt64Coin(types.DefaultClaimDenom, 2000)),
			ActionCompleted:        []bool{false, false, false, false, false},
		},
	}
	suite.app.AccountKeeper.SetAccount(suite.ctx, authtypes.NewBaseAccount(addr1, nil, 0, 0))

	err := suite.app.ClaimKeeper.SetClaimRecords(suite.ctx, claimRecords)
	suite.Require().NoError(err)

	coins1, err := suite.app.ClaimKeeper.GetUserTotalClaimable(suite.ctx, addr1)
	suite.Require().NoError(err)
	suite.Require().Equal(coins1, claimRecords[0].InitialClaimableAmount)

	suite.app.ClaimKeeper.AfterDelegationModified(suite.ctx, addr1, sdk.ValAddress(addr1))
	claim, err := suite.app.ClaimKeeper.GetClaimRecord(suite.ctx, addr1)
	suite.NoError(err)
	suite.False(claim.ActionCompleted[types.ActionDelegateStake])

	coins1, err = suite.app.ClaimKeeper.GetUserTotalClaimable(suite.ctx, addr1)
	suite.Require().NoError(err)
	suite.Require().Equal(coins1, claimRecords[0].InitialClaimableAmount)
}

func (suite *KeeperTestSuite) TestDelegationAutoWithdrawAndDelegateMore() {
	suite.SetupTest()
	suite.app.ClaimKeeper.CreateModuleAccount(suite.ctx, sdk.NewCoin(sdk.DefaultBondDenom, sdk.NewInt(10000000)))
	// set airdrop enabled but with date in the future
	suite.app.ClaimKeeper.SetParams(suite.ctx, types.Params{
		AirdropEnabled:     true,
		AirdropStartTime:   time.Now().Add(time.Hour * -1),
		ClaimDenom:         sdk.DefaultBondDenom,
		DurationUntilDecay: time.Hour,
		DurationOfDecay:    time.Hour * 4,
		AllowedClaimers:    make([]types.ClaimAuthorization, 0),
	})

	pub1 := secp256k1.GenPrivKey().PubKey()
	pub2 := secp256k1.GenPrivKey().PubKey()
	addr1 := sdk.AccAddress(pub1.Address())
	addr2 := sdk.AccAddress(pub2.Address())

	claimRecords := []types.ClaimRecord{
		{
			Address:                addr1.String(),
			InitialClaimableAmount: sdk.NewCoins(sdk.NewInt64Coin(sdk.DefaultBondDenom, 1000)),
			ActionCompleted:        []bool{false, false, false, false, false},
		},
		{
			Address:                addr2.String(),
			InitialClaimableAmount: sdk.NewCoins(sdk.NewInt64Coin(sdk.DefaultBondDenom, 1000)),
			ActionCompleted:        []bool{false, false, false, false, false},
		},
	}

	suite.app.AccountKeeper.SetAccount(suite.ctx, authtypes.NewBaseAccount(addr1, nil, 0, 0))
	suite.app.AccountKeeper.SetAccount(suite.ctx, authtypes.NewBaseAccount(addr2, nil, 0, 0))

	err := suite.app.ClaimKeeper.SetClaimRecords(suite.ctx, claimRecords)
	suite.Require().NoError(err)

	cr, err := suite.app.ClaimKeeper.GetClaimRecord(suite.ctx, addr1)
	suite.Require().NoError(err)
	suite.Require().Equal(cr, claimRecords[0])
	coins1, err := suite.app.ClaimKeeper.GetUserTotalClaimable(suite.ctx, addr1)
	suite.Require().NoError(err)
	suite.Require().Equal(claimRecords[1].InitialClaimableAmount.String(), coins1.String())

	coins2, err := suite.app.ClaimKeeper.GetUserTotalClaimable(suite.ctx, addr2)
	suite.Require().NoError(err)
	suite.Require().Equal(coins2, claimRecords[1].InitialClaimableAmount)

	// addr1 becomes validator
	validator, err := stakingtypes.NewValidator(sdk.ValAddress(addr1), pub1, stakingtypes.Description{})
	suite.Require().NoError(err)
	validator = stakingkeeper.TestingUpdateValidator(suite.app.StakingKeeper, suite.ctx, validator, true)
	suite.app.StakingKeeper.AfterValidatorCreated(suite.ctx, validator.GetOperator())

	validator, _ = validator.AddTokensFromDel(sdk.TokensFromConsensusPower(1, sdk.DefaultPowerReduction))
	delAmount := sdk.TokensFromConsensusPower(1, sdk.DefaultPowerReduction)
	err = FundAccount(suite.app.BankKeeper, suite.ctx, addr2, sdk.NewCoins(sdk.NewCoin(sdk.DefaultBondDenom, delAmount)))
	suite.NoError(err)

	balance := suite.app.BankKeeper.GetAllBalances(suite.ctx, addr2)
	suite.Require().Equal(
		balance.AmountOf(sdk.DefaultBondDenom).String(),
		delAmount.String())

	_, err = suite.app.StakingKeeper.Delegate(suite.ctx, addr2, delAmount, sdkstakingtypes.Unbonded, validator, true)
	suite.NoError(err)

	// delegation should automatically call claim and withdraw balance
	claimedCoins := suite.app.BankKeeper.GetAllBalances(suite.ctx, addr2)
	suite.Require().Equal(
		claimedCoins.AmountOf(sdk.DefaultBondDenom).String(),
		claimRecords[1].InitialClaimableAmount.AmountOf(sdk.DefaultBondDenom).Quo(sdk.NewInt(int64(len(claimRecords[1].ActionCompleted)))).String())

	_, err = suite.app.StakingKeeper.Delegate(suite.ctx, addr2, claimedCoins.AmountOf(sdk.DefaultBondDenom), sdkstakingtypes.Unbonded, validator, true)
	suite.NoError(err)
}

func (suite *KeeperTestSuite) TestEndAirdrop() {
	// set airdrop enabled but with date in the future
	suite.app.ClaimKeeper.SetParams(suite.ctx, types.Params{
		AirdropEnabled:     true,
		AirdropStartTime:   time.Now().Add(time.Hour * -1),
		ClaimDenom:         types.DefaultClaimDenom,
		DurationUntilDecay: time.Hour,
		DurationOfDecay:    time.Hour * 4,
		AllowedClaimers:    make([]types.ClaimAuthorization, 0),
	})

	pub1 := secp256k1.GenPrivKey().PubKey()
	pub2 := secp256k1.GenPrivKey().PubKey()
	addr1 := sdk.AccAddress(pub1.Address())
	addr2 := sdk.AccAddress(pub2.Address())

	claimRecords := []types.ClaimRecord{
		{
			Address:                addr1.String(),
			InitialClaimableAmount: sdk.NewCoins(sdk.NewInt64Coin(types.DefaultClaimDenom, 1000)),
			ActionCompleted:        []bool{false, false, false, false, false},
		},
		{
			Address:                addr2.String(),
			InitialClaimableAmount: sdk.NewCoins(sdk.NewInt64Coin(types.DefaultClaimDenom, 1000)),
			ActionCompleted:        []bool{false, false, false, false, false},
		},
	}

	suite.app.AccountKeeper.SetAccount(suite.ctx, authtypes.NewBaseAccount(addr1, nil, 0, 0))
	suite.app.AccountKeeper.SetAccount(suite.ctx, authtypes.NewBaseAccount(addr2, nil, 0, 0))

	err := suite.app.ClaimKeeper.SetClaimRecords(suite.ctx, claimRecords)
	suite.Require().NoError(err)

	err = suite.app.ClaimKeeper.EndAirdrop(suite.ctx)
	suite.Require().NoError(err)

	moduleAccAddr := suite.app.AccountKeeper.GetModuleAddress(types.ModuleName)
	coins := suite.app.BankKeeper.GetBalance(suite.ctx, moduleAccAddr, types.DefaultClaimDenom)
	suite.Require().Equal(sdk.NewInt64Coin(types.DefaultClaimDenom, 0).String(), coins.String())
}
