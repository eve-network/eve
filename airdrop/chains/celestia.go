package chains

// error max size response
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

func Celestia() ([]banktypes.Balance, []config.Reward, int, error) {
	err := godotenv.Load()
	if err != nil {
		return nil, nil, 0, fmt.Errorf("failed to load env: %w", err)
	}

	blockHeight, err := utils.GetLatestHeight(config.GetCelestiaConfig().RPC + "/status")
	if err != nil {
		return nil, nil, 0, fmt.Errorf("failed to get latest height for Celestia: %w", err)
	}

	grpcAddr := config.GetCelestiaConfig().GRPCAddr
	grpcConn, err := utils.SetupGRPCConnection(grpcAddr)
	if err != nil {
		return nil, nil, 0, fmt.Errorf("failed to connect to gRPC Celestia: %w", err)
	}
	defer grpcConn.Close()
	stakingClient := stakingtypes.NewQueryClient(grpcConn)

	delegators := []stakingtypes.DelegationResponse{}

	validators, err := utils.GetValidators(stakingClient, blockHeight)
	if err != nil {
		return nil, nil, 0, fmt.Errorf("failed to get Celestia validators: %w", err)
	}

	log.Println("Validators: ", len(validators))
	for validatorIndex, validator := range validators {
		url := config.GetCelestiaConfig().API + "/cosmos/staking/v1beta1/validators/" + validator.OperatorAddress + "/delegations?pagination.limit=" + strconv.Itoa(config.LimitPerPage) + "&pagination.count_total=true"
		delegations, total, err := utils.FetchDelegations(url)
		if err != nil {
			return nil, nil, 0, fmt.Errorf("failed to fetch delegations for Celestia: %w", err)
		}
		log.Println(validator.OperatorAddress)
		log.Println("Response ", len(delegations))
		log.Println("Celestia validator "+strconv.Itoa(validatorIndex)+" ", total)
		delegators = append(delegators, delegations...)
	}

	usd := sdkmath.LegacyMustNewDecFromStr("20")

	apiURL := config.APICoingecko + config.GetCelestiaConfig().CoinID + "&vs_currencies=usd"
	tokenInUsd, err := utils.FetchTokenPrice(apiURL, config.GetCelestiaConfig().CoinID)
	if err != nil {
		return nil, nil, 0, fmt.Errorf("failed to fetch Celestia token price: %w", err)
	}
	tokenIn20Usd := usd.Quo(tokenInUsd)

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
		eveAirdrop := (eveAirdrop.MulInt64(int64(config.GetCelestiaConfig().Percent))).QuoInt64(100).Mul(token).QuoTruncate(totalTokenDelegate)
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
			ChainID:         config.GetCelestiaConfig().ChainID,
		})
		testAmount = eveAirdrop.Add(testAmount)
		balanceInfo = append(balanceInfo, banktypes.Balance{
			Address: eveBech32Address,
			Coins:   sdk.NewCoins(sdk.NewCoin("eve", eveAirdrop.TruncateInt())),
		})
	}
	log.Println("Celestia balance: ", testAmount)
	// Write delegations to file
	// fileForDebug, _ := json.MarshalIndent(rewardInfo, "", " ")
	// _ = os.WriteFile("rewards.json", fileForDebug, 0644)

	// fileBalance, _ := json.MarshalIndent(balanceInfo, "", " ")
	// _ = os.WriteFile("balance.json", fileBalance, 0644)
	return balanceInfo, rewardInfo, len(balanceInfo), nil
}
