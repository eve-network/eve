package wasmbinding

import (
	"github.com/CosmWasm/wasmd/x/wasm"
	wasmkeeper "github.com/CosmWasm/wasmd/x/wasm/keeper"
	"github.com/eve-network/eve/wasmbinding/encoder"
)

func RegisterCustomPlugins(
	registry *encoder.EncoderRegistry,
) []wasmkeeper.Option {
	encoderPluginOpt := wasmkeeper.WithMessageEncoders(
		encoder.MessageEncoders(registry),
	)

	return []wasm.Option{
		encoderPluginOpt,
	}
}
