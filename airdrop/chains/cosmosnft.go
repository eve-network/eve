package chains

import (
	"fmt"
	"log"
	"sync"

	"github.com/eve-network/eve/airdrop/config"
	"github.com/eve-network/eve/airdrop/utils"

	sdkmath "cosmossdk.io/math"

	sdk "github.com/cosmos/cosmos-sdk/types"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
)

func Cosmosnft(contract string, percent int64) ([]banktypes.Balance, []config.Reward, int, error) {
	tokenIds, err := utils.FetchTokenIds(contract, config.GetStargazeConfig().API)
	if err != nil {
		return nil, nil, 0, fmt.Errorf("failed to fetch token ids: %w", err)
	}

	// Create channels to receive balance info and rewards
	balanceCh := make(chan banktypes.Balance)
	rewardCh := make(chan config.Reward)
	doneCh := make(chan struct{})

	// Use a buffered channel as a semaphore to limit concurrency
	semaphore := make(chan struct{}, 10) // Limit to 10 concurrent goroutines

	wg := &sync.WaitGroup{}
	wg.Add(len(tokenIds))

	for _, token := range tokenIds {
		semaphore <- struct{}{} // Acquire semaphore
		go func(token string) {
			defer func() { <-semaphore }() // Release semaphore when done
			defer wg.Done()

			nftHolders, err := utils.FetchTokenInfo(token, contract, config.GetStargazeConfig().API)
			if err != nil {
				log.Printf("Error fetching token info for %s: %v\n", token, err)
				return
			}

			if nftHolders.Address == "" {
				log.Printf("Token id: %s is not NFT\n", token)
				return
			}

			eveBech32Address, err := utils.ConvertBech32Address(nftHolders.Address)
			if err != nil {
				log.Printf("Error converting Bech32Address for %s: %v\n", nftHolders.Address, err)
				return
			}

			allEveAirdrop := sdkmath.LegacyMustNewDecFromStr(config.EveAirdrop)
			eveAirdrop := (allEveAirdrop.MulInt64(percent)).QuoInt64(100).QuoInt(sdkmath.NewInt(int64(len(tokenIds))))
			rewardCh <- config.Reward{
				Address:         nftHolders.Address,
				EveAddress:      eveBech32Address,
				EveAirdropToken: eveAirdrop,
				ChainID:         config.GetBadKidsConfig().ChainID,
			}

			balanceCh <- banktypes.Balance{
				Address: eveBech32Address,
				Coins:   sdk.NewCoins(sdk.NewCoin("eve", eveAirdrop.TruncateInt())),
			}
		}(token)
	}

	// Close channels once all goroutines are done
	go func() {
		wg.Wait()
		close(balanceCh)
		close(rewardCh)
		close(doneCh)
	}()

	var balanceInfo []banktypes.Balance
	var rewardInfo []config.Reward

	for {
		select {
		case balance, ok := <-balanceCh:
			if !ok {
				balanceCh = nil
			} else {
				balanceInfo = append(balanceInfo, balance)
			}
		case reward, ok := <-rewardCh:
			if !ok {
				rewardCh = nil
			} else {
				rewardInfo = append(rewardInfo, reward)
			}
		case <-doneCh:
			return balanceInfo, rewardInfo, len(balanceInfo), nil
		}
	}
}
