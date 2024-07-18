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

func main() {
	startTime := time.Now()

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
			return chains.Cosmosnft(config.BadkidsContractAddress, int64(config.GetBadKidsConfig().Percent), config.GetStargazeConfig().API)
		},
		"cryptonium": func() ([]banktypes.Balance, []config.Reward, int, error) {
			return chains.Cosmosnft(config.CryptoniumContractAddress, int64(config.GetCryptoniumConfig().Percent), config.GetStargazeConfig().API)
		},
		"milady": chains.Ethereumnft,
	}

	wg := sync.WaitGroup{}
	balanceInfoCh := make(chan []banktypes.Balance, len(balanceFunctions))
	lengthBalanceInfoCh := make(chan int, len(balanceFunctions))
	errFuncCh := make(chan string, len(balanceFunctions))

	for name, fn := range balanceFunctions {
		wg.Add(1)
		go func(name string, fn utils.BalanceFunction) {
			defer wg.Done()
			log.Printf("Fetching balance info: %s\n", name)
			info, _, len, err := fn()
			if err != nil {
				log.Printf("Error executing balanceFunction %s: %v\n", name, err)
				errFuncCh <- name
				return
			}
			balanceInfoCh <- info
			lengthBalanceInfoCh <- len
		}(name, fn)
	}

	go func() {
		wg.Wait()
		close(balanceInfoCh)
		close(lengthBalanceInfoCh)
		close(errFuncCh)
	}()

	total := 0
	balanceAkashInfo := []banktypes.Balance{}

	for lenCh := range lengthBalanceInfoCh {
		total += lenCh
	}

	for infoCh := range balanceInfoCh {
		balanceAkashInfo = append(balanceAkashInfo, infoCh...)
	}

	for funcCh := range errFuncCh {
		errFuncName := funcCh
		log.Printf("Retrying failed balance function: %s\n", errFuncName)
		info, _, len, err := utils.RetryableBalanceFunc(balanceFunctions[errFuncName])()
		if err != nil {
			panic(fmt.Sprintf("Error executing balanceFunction %s: %v", errFuncName, err))
		}
		total += len
		balanceAkashInfo = append(balanceAkashInfo, info...)
	}

	log.Printf("Total: %d\n", total)
	log.Printf("Number of balances: %d\n", len(balanceAkashInfo))

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

	log.Printf("Check balance: %d\n", checkBalance)

	// Write balance info to file
	fileBalance, err := json.MarshalIndent(balanceInfo, "", " ")
	if err != nil {
		log.Fatal("Failed to marshal balance info:", err)
	}

	err = os.WriteFile("balance.json", fileBalance, 0o600)
	if err != nil {
		log.Fatal("Failed to write balance info to file:", err)
	}

	duration := time.Since(startTime)
	log.Printf("Total time taken: %v\n", duration)
}
