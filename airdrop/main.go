package main

import (
	"encoding/json"
	"fmt"
	"math"
	"os"
	"sync"
	"time"

	"github.com/eve-network/eve/airdrop/config"
	"github.com/eve-network/eve/airdrop/utils"

	sdkmath "cosmossdk.io/math"

	sdk "github.com/cosmos/cosmos-sdk/types"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
)

// got to export genesis state from neutron and bostrom chain

const (
	EveAirdrop   = "1000000000" // 1,000,000,000
	Badkids      = "stars19jq6mj84cnt9p7sagjxqf8hxtczwc8wlpuwe4sh62w45aheseues57n420"
	Cryptonium   = "stars1g2ptrqnky5pu70r3g584zpk76cwqplyc63e8apwayau6l3jr8c0sp9q45u"
	APICoingecko = "https://api.coingecko.com/api/v3/simple/price?ids="
	MaxRetries   = 5
	BackOff      = 200 * time.Millisecond
)

// Define a function type that returns balance info, reward info and length
type balanceFunction func() ([]banktypes.Balance, []config.Reward, int, error)

// Retryable function to wrap balanceFunction with retry logic
func retryable(fn balanceFunction) balanceFunction {
	return func() ([]banktypes.Balance, []config.Reward, int, error) {
		for attempt := 1; attempt <= MaxRetries; attempt++ {
			balances, rewards, length, err := fn()
			if err == nil {
				return balances, rewards, length, nil
			}
			fmt.Printf("Failed attempt %d for function %s: %v\n", attempt, utils.GetFunctionName(fn), err)
		}
		return nil, nil, 0, fmt.Errorf("maximum retries reached for function %s", utils.GetFunctionName(fn))
	}
}

func main() {
	// Capture start time
	startTime := time.Now()

	// Define balance functions with their associated names
	balanceFunctions := map[string]balanceFunction{
		"akash":      akash,
		"bostrom":    bostrom,
		"celestia":   celestia,
		"composable": composable,
		"cosmos":     cosmos,
		"neutron":    neutron,
		"sentinel":   sentinel,
		"stargaze":   stargaze,
		"terra":      terra,
		"terrac":     terrac,
		"badkids": func() ([]banktypes.Balance, []config.Reward, int, error) {
			return cosmosnft(Badkids, int64(config.GetBadKidsConfig().Percent))
		},
		"cryptonium": func() ([]banktypes.Balance, []config.Reward, int, error) {
			return cosmosnft(Cryptonium, int64(config.GetCryptoniumConfig().Percent))
		},
		// need set coin type on Eve
		"milady": ethereumnft,
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
		go func(name string, fn balanceFunction) {
			defer wg.Done()
			fmt.Println("Fetching balance info: ", name)
			info, _, len, err := fn() // Call the function
			if err != nil {
				fmt.Printf("Error executing balanceFunction %s: %v\n", name, err)
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
		fmt.Println("Retry the failed balance function: ", errFuncName)
		info, _, len, err := retryable(balanceFunctions[errFuncName])()
		if err != nil {
			panic(fmt.Sprintf("error executing balanceFunction %s: %v", errFuncName, err))
		}
		total += len
		balanceAkashInfo = append(balanceAkashInfo, info...)
	}

	fmt.Println("Total: ", total)
	fmt.Println(len(balanceAkashInfo))

	airdropMap := make(map[string]int)
	for _, info := range balanceAkashInfo {
		amount := airdropMap[info.Address]
		airdropMap[info.Address] = amount + int(info.Coins.AmountOf("eve").Int64())
	}

	balanceInfo := []banktypes.Balance{}
	checkBalance := 0
	for address, amount := range airdropMap {
		checkBalance += amount
		balanceInfo = append(balanceInfo, banktypes.Balance{
			Address: address,
			Coins:   sdk.NewCoins(sdk.NewCoin("eve", sdkmath.NewInt(int64(amount)))),
		})
	}

	fmt.Println("Check balance: ", checkBalance)

	// // Write delegations to file
	// fileForDebug, _ := json.MarshalIndent(rewardComposableInfo, "", " ")
	// _ = os.WriteFile("rewards.json", fileForDebug, 0644)

	fileBalance, _ := json.MarshalIndent(balanceInfo, "", " ")
	_ = os.WriteFile("balance.json", fileBalance, 0o600)

	// Calculate and print total time duration
	duration := time.Since(startTime)
	fmt.Printf("Total time taken: %v\n", duration)
}

// Define a function type that returns token price from a price source
type tokenPriceFunction func(apiURL string) (sdkmath.LegacyDec, error)

func fetchTokenPriceWithRetry(fn tokenPriceFunction) tokenPriceFunction {
	return func(apiURL string) (sdkmath.LegacyDec, error) {
		for attempt := 1; attempt <= MaxRetries; attempt++ {
			data, err := fn(apiURL)
			if err == nil {
				return data, nil
			}

			fmt.Printf("Failed attempt %d for function %s: %v\n", attempt, utils.GetFunctionName(fn), err)

			if attempt < MaxRetries {
				// Calculate backoff duration using exponential backoff strategy
				backoffDuration := time.Duration(BackOff.Seconds() * math.Pow(2, float64(attempt)))
				fmt.Printf("Retrying after %s...\n", backoffDuration)
				time.Sleep(backoffDuration)
			}
		}
		return sdkmath.LegacyDec{}, fmt.Errorf("maximum retries reached for function %s", utils.GetFunctionName(fn))
	}
}
