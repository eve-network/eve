package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"sync"
	"time"

	"github.com/eve-network/eve/airdrop/chains"
	"github.com/eve-network/eve/airdrop/config"
	"github.com/eve-network/eve/airdrop/utils"

	sdkmath "cosmossdk.io/math"

	sdk "github.com/cosmos/cosmos-sdk/types"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
)

// got to export genesis state from neutron and bostrom chain

const (
	Badkids    = "stars19jq6mj84cnt9p7sagjxqf8hxtczwc8wlpuwe4sh62w45aheseues57n420"
	Cryptonium = "stars1g2ptrqnky5pu70r3g584zpk76cwqplyc63e8apwayau6l3jr8c0sp9q45u"
)

func main() {
	// Capture start time
	startTime := time.Now()

	// Define balance functions with their associated names
	balanceFunctions := map[string]utils.BalanceFunction{
		"akash":      chains.Akash,
		"bostrom":    chains.Bostrom,
		"celestia":   chains.Celestia,
		"composable": chains.Composable,
		"cosmos":     chains.Cosmos,
		"neutron":    chains.Neutron,
		"sentinel":   chains.Sentinel,
		"stargaze":   chains.Stargaze,
		"terra":      chains.Terra,
		"terrac":     chains.Terrac,
		"badkids": func() ([]banktypes.Balance, []config.Reward, int, error) {
			return chains.Cosmosnft(Badkids, int64(config.GetBadKidsConfig().Percent))
		},
		"cryptonium": func() ([]banktypes.Balance, []config.Reward, int, error) {
			return chains.Cosmosnft(Cryptonium, int64(config.GetCryptoniumConfig().Percent))
		},
		// need set coin type on Eve
		"milady": chains.Ethereumnft,
	}

	lenBalanceFunctions := len(balanceFunctions)
	wg := &sync.WaitGroup{}
	wg.Add(lenBalanceFunctions)

	// Channel to collect balance info from goroutines
	balanceInfoCh := make(chan []banktypes.Balance, lenBalanceFunctions)

	// Channel to collect length of balance info from goroutines
	lengthBalanceInfoCh := make(chan int, lenBalanceFunctions)

	// Channel to collect name of error balanceFunction from goroutines
	errFuncCh := make(chan string, lenBalanceFunctions)

	// Iterate over the balanceFunctions map and run each function in a goroutine
	for name, fn := range balanceFunctions {
		go func(name string, fn utils.BalanceFunction) {
			defer wg.Done()
			log.Println("Fetching balance info: ", name)
			info, _, len, err := fn() // Call the function
			if err != nil {
				log.Printf("Error executing balanceFunction %s: %v\n", name, err)
				errFuncCh <- name // Send the error function's name to channel
				return
			}
			balanceInfoCh <- info      // Send balance info to channel
			lengthBalanceInfoCh <- len // Send length of balance info to channel
		}(name, fn)
	}

	go func() {
		// Wait for all goroutines to complete
		wg.Wait()
		// Close channels
		close(balanceInfoCh)
		close(lengthBalanceInfoCh)
		close(errFuncCh)
	}()

	total := 0
	balanceAkashInfo := []banktypes.Balance{}

	// Collect results from channels
	for lenCh := range lengthBalanceInfoCh {
		total += lenCh
	}

	for infoCh := range balanceInfoCh {
		balanceAkashInfo = append(balanceAkashInfo, infoCh...)
	}

	for funcCh := range errFuncCh {
		// Retrieve the error function's name from the channel
		errFuncName := funcCh
		// Retry the failed balance function
		log.Println("Retry the failed balance function: ", errFuncName)
		info, _, len, err := utils.RetryableBalanceFunc(balanceFunctions[errFuncName])()
		if err != nil {
			panic(fmt.Sprintf("error executing balanceFunction %s: %v", errFuncName, err))
		}
		total += len
		balanceAkashInfo = append(balanceAkashInfo, info...)
	}

	log.Println("Total: ", total)
	log.Println(len(balanceAkashInfo))

	airdropMap := make(map[string]int)
	for _, info := range balanceAkashInfo {
		amount := airdropMap[info.Address]
		airdropMap[info.Address] = amount + int(info.Coins.AmountOf("eve").Int64())
	}

	balanceInfo := []banktypes.Balance{}
	checkBalance := 0
	for address, amount := range airdropMap {
		if amount == 0 {
			continue
		}
		checkBalance += amount
		balanceInfo = append(balanceInfo, banktypes.Balance{
			Address: address,
			Coins:   sdk.NewCoins(sdk.NewCoin("eve", sdkmath.NewInt(int64(amount)))),
		})
	}

	log.Println("Check balance: ", checkBalance)

	// // Write delegations to file
	// fileForDebug, _ := json.MarshalIndent(rewardComposableInfo, "", " ")
	// _ = os.WriteFile("rewards.json", fileForDebug, 0644)

	fileBalance, _ := json.MarshalIndent(balanceInfo, "", " ")
	_ = os.WriteFile("balance.json", fileBalance, 0o600)

	// Calculate and print total time duration
	duration := time.Since(startTime)
	log.Printf("Total time taken: %v\n", duration)
}
