package chains

import (
	"fmt"
	"log"
	"strconv"

	"github.com/eve-network/eve/airdrop/config"
	"github.com/eve-network/eve/airdrop/utils"
	"github.com/joho/godotenv"

	sdkmath "cosmossdk.io/math"

	sdk "github.com/cosmos/cosmos-sdk/types"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
)

func Composable() ([]banktypes.Balance, []config.Reward, int, error) {
	err := godotenv.Load()
	if err != nil {
		return nil, nil, 0, fmt.Errorf("failed to load env: %w", err)
	}

	blockHeight, err := utils.GetLatestHeight(config.GetComposableConfig().RPC + "/status")
	if err != nil {
		return nil, nil, 0, fmt.Errorf("failed to get latest height for Composable: %w", err)
	}

	grpcAddr := config.GetComposableConfig().GRPCAddr
	grpcConn, err := utils.SetupGRPCConnection(grpcAddr)
	if err != nil {
		return nil, nil, 0, fmt.Errorf("failed to connect to gRPC Composable: %w", err)
	}
	defer grpcConn.Close()
	stakingClient := stakingtypes.NewQueryClient(grpcConn)

	delegators := []stakingtypes.DelegationResponse{}

	validators, err := utils.GetValidators(stakingClient, blockHeight)
	if err != nil {
		return nil, nil, 0, fmt.Errorf("failed to get Composable validators: %w", err)
	}

	log.Println("Validators: ", len(validators))
	for validatorIndex, validator := range validators {
		delegationsResponse, err := utils.GetValidatorDelegations(stakingClient, validator.OperatorAddress, blockHeight)
		if err != nil {
			return nil, nil, 0, fmt.Errorf("failed to query delegate info for Composable validator: %w", err)
		}
		total := delegationsResponse.Pagination.Total
		log.Println("Response ", len(delegationsResponse.DelegationResponses))
		log.Println("Composable validator "+strconv.Itoa(validatorIndex)+" ", total)
		delegators = append(delegators, delegationsResponse.DelegationResponses...)
	}

	usd := sdkmath.LegacyMustNewDecFromStr("20")

	apiURL := config.APICoingecko + config.GetComposableConfig().CoinID + "&vs_currencies=usd"
	tokenInUsd, err := utils.FetchTokenPrice(apiURL, config.GetComposableConfig().CoinID)
	if err != nil {
		return nil, nil, 0, fmt.Errorf("failed to fetch Composable token price: %w", err)
	}
	tokenIn20Usd := usd.QuoTruncate(tokenInUsd)

	rewardInfo := []config.Reward{}
	balanceInfo := []banktypes.Balance{}

	totalTokenDelegate := sdkmath.LegacyMustNewDecFromStr("0")
	for _, delegator := range delegators {
		validatorIndex := utils.FindValidatorInfo(validators, delegator.Delegation.ValidatorAddress)
		validatorInfo := validators[validatorIndex]
		token := (delegator.Delegation.Shares.MulInt(validatorInfo.Tokens)).QuoTruncate(validatorInfo.DelegatorShares)
		totalTokenDelegate = totalTokenDelegate.Add(token)
	}
	eveAirdrop := sdkmath.LegacyMustNewDecFromStr(config.EveAirdrop)
	testAmount, _ := sdkmath.LegacyNewDecFromStr("0")
	for _, delegator := range delegators {
		validatorIndex := utils.FindValidatorInfo(validators, delegator.Delegation.ValidatorAddress)
		validatorInfo := validators[validatorIndex]
		token := (delegator.Delegation.Shares.MulInt(validatorInfo.Tokens)).QuoTruncate(validatorInfo.DelegatorShares)
		if token.LT(tokenIn20Usd) {
			continue
		}
		eveAirdrop := (eveAirdrop.MulInt64(int64(config.GetComposableConfig().Percent))).QuoInt64(100).Mul(token).QuoTruncate(totalTokenDelegate)
		eveBech32Address, err := utils.ConvertBech32Address(delegator.Delegation.DelegatorAddress)
		if err != nil {
			return nil, nil, 0, fmt.Errorf("failed to convert Bech32Address: %w", err)
		}
		rewardInfo = append(rewardInfo, config.Reward{
			Address:         delegator.Delegation.DelegatorAddress,
			EveAddress:      eveBech32Address,
			Shares:          delegator.Delegation.Shares,
			Token:           token,
			EveAirdropToken: eveAirdrop,
			ChainID:         config.GetComposableConfig().ChainID,
		})
		testAmount = eveAirdrop.Add(testAmount)
		balanceInfo = append(balanceInfo, banktypes.Balance{
			Address: eveBech32Address,
			Coins:   sdk.NewCoins(sdk.NewCoin("eve", eveAirdrop.TruncateInt())),
		})
	}
	log.Println("Composable balance: ", testAmount)
	// Write delegations to file
	// fileForDebug, _ := json.MarshalIndent(rewardInfo, "", " ")
	// _ = os.WriteFile("rewards.json", fileForDebug, 0644)

	// fileBalance, _ := json.MarshalIndent(balanceInfo, "", " ")
	// _ = os.WriteFile("balance.json", fileBalance, 0644)
	return balanceInfo, rewardInfo, len(balanceInfo), nil
}
