package wasmbinding

import (
	"github.com/CosmWasm/wasmd/x/wasm"
	wasmkeeper "github.com/CosmWasm/wasmd/x/wasm/keeper"

	tokenfactorykeeper "github.com/notional-labs/eve/x/tokenfactory/keeper"
)

func RegisterCustomPlugins(

	tokenFactory *tokenfactorykeeper.Keeper,
) []wasmkeeper.Option {
	wasmQueryPlugin := NewQueryPlugin(tokenFactory)

	queryPluginOpt := wasmkeeper.WithQueryPlugins(&wasmkeeper.QueryPlugins{
		Custom: CustomQuerier(wasmQueryPlugin),
	})
	messengerDecoratorOpt := wasmkeeper.WithMessageHandlerDecorator(
		CustomMessageDecorator(tokenFactory),
	)

	return []wasm.Option{
		queryPluginOpt,
		messengerDecoratorOpt,
	}
}
