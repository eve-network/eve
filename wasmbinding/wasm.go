package wasmbinding

import (
	"github.com/CosmWasm/wasmd/x/wasm"
	wasmkeeper "github.com/CosmWasm/wasmd/x/wasm/keeper"
	bankkeeper "github.com/cosmos/cosmos-sdk/x/bank/keeper"
	"github.com/eve-network/eve/wasmbinding/encoder"

	tokenfactorykeeper "github.com/eve-network/eve/x/tokenfactory/keeper"
)

func RegisterCustomPlugins(
	bank *bankkeeper.BaseKeeper,
	tokenFactory *tokenfactorykeeper.Keeper,
	registry *encoder.EncoderRegistry,
) []wasmkeeper.Option {
	wasmQueryPlugin := NewQueryPlugin(tokenFactory)

	queryPluginOpt := wasmkeeper.WithQueryPlugins(&wasmkeeper.QueryPlugins{
		Custom: CustomQuerier(wasmQueryPlugin),
	})
	messengerDecoratorOpt := wasmkeeper.WithMessageHandlerDecorator(
		CustomMessageDecorator(bank, tokenFactory),
	)
	encoderPluginOpt := wasmkeeper.WithMessageEncoders(
		encoder.MessageEncoders(registry),
	)

	return []wasm.Option{
		queryPluginOpt,
		messengerDecoratorOpt,
		encoderPluginOpt,
	}
}
