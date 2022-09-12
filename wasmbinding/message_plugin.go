package wasmbinding

import (
	"encoding/json"

	wasmvmtypes "github.com/CosmWasm/wasmvm/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	bankkeeper "github.com/cosmos/cosmos-sdk/x/bank/keeper"
	wasmkeeper "github.com/notional-labs/eve/x/wasm/keeper"

	"github.com/notional-labs/eve/wasmbinding/bindings"

	tokenfactorykeeper "github.com/notional-labs/eve/x/tokenfactory/keeper"
	tokenfactorytypes "github.com/notional-labs/eve/x/tokenfactory/types"
)

// CustomMessageDecorator returns decorator for custom CosmWasm bindings messages
func CustomMessageDecorator(bank *bankkeeper.BaseKeeper, tokenFactory *tokenfactorykeeper.Keeper) func(wasmkeeper.Messenger) wasmkeeper.Messenger {
	return func(old wasmkeeper.Messenger) wasmkeeper.Messenger {
		return &CustomMessenger{
			wrapped:      old,
			bank:         bank,
			tokenFactory: tokenFactory,
		}
	}
}

type CustomMessenger struct {
	wrapped      wasmkeeper.Messenger
	bank         *bankkeeper.BaseKeeper
	tokenFactory *tokenfactorykeeper.Keeper
}

var _ wasmkeeper.Messenger = (*CustomMessenger)(nil)

// DispatchMsg executes on the contractMsg.
func (m *CustomMessenger) DispatchMsg(ctx sdk.Context, contractAddr sdk.AccAddress, contractIBCPortID string, msg wasmvmtypes.CosmosMsg) ([]sdk.Event, [][]byte, error) {
	if msg.Custom != nil {
		// only handle the happy path where this is really creating / minting / swapping ...
		// leave everything else for the wrapped version
		var contractMsg bindings.EveMsg
		if err := json.Unmarshal(msg.Custom, &contractMsg); err != nil {
			return nil, nil, sdkerrors.Wrap(err, "eve msg")
		}
		if contractMsg.CreateDenom != nil {
			return m.createDenom(ctx, contractAddr, contractMsg.CreateDenom)
		}
		if contractMsg.MintTokens != nil {
			return m.mintTokens(ctx, contractAddr, contractMsg.MintTokens)
		}
		if contractMsg.ChangeOwner != nil {
			return m.changeAdmin(ctx, contractAddr, contractMsg.ChangeOwner)
		}
	}
	return m.wrapped.DispatchMsg(ctx, contractAddr, contractIBCPortID, msg)
}

// createDenom creates a new token denom
func (m *CustomMessenger) createDenom(ctx sdk.Context, contractAddr sdk.AccAddress, createDenom *bindings.CreateDenom) ([]sdk.Event, [][]byte, error) {
	err := PerformCreateDenom(m.tokenFactory, m.bank, ctx, contractAddr, createDenom)
	if err != nil {
		return nil, nil, sdkerrors.Wrap(err, "perform create denom")
	}
	return nil, nil, nil
}

// PerformCreateDenom is used with createDenom to create a token denom; validates the msgCreateDenom.
func PerformCreateDenom(f *tokenfactorykeeper.Keeper, b *bankkeeper.BaseKeeper, ctx sdk.Context, contractAddr sdk.AccAddress, createDenom *bindings.CreateDenom) error {
	if createDenom == nil {
		return wasmvmtypes.InvalidRequest{Err: "create denom null create denom"}
	}

	msgServer := tokenfactorykeeper.NewMsgServerImpl(*f)

	msgCreateDenom := tokenfactorytypes.NewMsgCreateDenom(
		contractAddr.String(),
		createDenom.Name,
		createDenom.Subdenom,
		createDenom.Precision,
		createDenom.MaxSupply,
		createDenom.CanChangeMaxSupply,
	)

	if err := msgCreateDenom.ValidateBasic(); err != nil {
		return sdkerrors.Wrap(err, "failed validating MsgCreateDenom")
	}

	// Create denom
	_, err := msgServer.CreateDenom(
		sdk.WrapSDKContext(ctx),
		msgCreateDenom,
	)
	if err != nil {
		return sdkerrors.Wrap(err, "creating denom")
	}
	return nil
}

// mintTokens mints tokens of a specified denom to an address.
func (m *CustomMessenger) mintTokens(ctx sdk.Context, contractAddr sdk.AccAddress, mint *bindings.MintTokens) ([]sdk.Event, [][]byte, error) {
	err := PerformMint(m.tokenFactory, m.bank, ctx, contractAddr, mint)
	if err != nil {
		return nil, nil, sdkerrors.Wrap(err, "perform mint")
	}
	return nil, nil, nil
}

// PerformMint used with mintTokens to validate the mint message and mint through token factory.
func PerformMint(f *tokenfactorykeeper.Keeper, b *bankkeeper.BaseKeeper, ctx sdk.Context, contractAddr sdk.AccAddress, mint *bindings.MintTokens) error {
	if mint == nil {
		return wasmvmtypes.InvalidRequest{Err: "mint token null mint"}
	}
	rcpt, err := parseAddress(mint.Recipient)
	if err != nil {
		return err
	}

	coin := sdk.Coin{Denom: mint.Denom, Amount: mint.Amount}
	sdkMsg := tokenfactorytypes.NewMsgMintAndSendTokens(
		contractAddr.String(),
		mint.Denom,
		int32(mint.Amount.Int64()),
		rcpt.String(),
	)
	if err = sdkMsg.ValidateBasic(); err != nil {
		return err
	}

	// Mint through token factory / message server
	msgServer := tokenfactorykeeper.NewMsgServerImpl(*f)
	_, err = msgServer.MintAndSendTokens(sdk.WrapSDKContext(ctx), sdkMsg)
	if err != nil {
		return sdkerrors.Wrap(err, "minting coins from message")
	}
	err = b.SendCoins(ctx, contractAddr, rcpt, sdk.NewCoins(coin))
	if err != nil {
		return sdkerrors.Wrap(err, "sending newly minted coins from message")
	}
	return nil
}

// changeAdmin changes the admin.
func (m *CustomMessenger) changeAdmin(ctx sdk.Context, contractAddr sdk.AccAddress, changeAdmin *bindings.ChangeOwner) ([]sdk.Event, [][]byte, error) {
	err := ChangeAdmin(m.tokenFactory, ctx, contractAddr, changeAdmin)
	if err != nil {
		return nil, nil, sdkerrors.Wrap(err, "failed to change admin")
	}
	return nil, nil, nil
}

// ChangeAdmin is used with changeAdmin to validate changeAdmin messages and to dispatch.
func ChangeAdmin(f *tokenfactorykeeper.Keeper, ctx sdk.Context, contractAddr sdk.AccAddress, changeAdmin *bindings.ChangeOwner) error {
	if changeAdmin == nil {
		return wasmvmtypes.InvalidRequest{Err: "changeAdmin is nil"}
	}
	newAdminAddr, err := parseAddress(changeAdmin.NewOwnerAddress)
	if err != nil {
		return err
	}

	changeAdminMsg := tokenfactorytypes.NewMsgUpdateOwner(contractAddr.String(), changeAdmin.Denom, newAdminAddr.String())
	if err := changeAdminMsg.ValidateBasic(); err != nil {
		return err
	}

	msgServer := tokenfactorykeeper.NewMsgServerImpl(*f)
	_, err = msgServer.UpdateOwner(sdk.WrapSDKContext(ctx), changeAdminMsg)
	if err != nil {
		return sdkerrors.Wrap(err, "failed changing admin from message")
	}
	return nil
}

// TODO: always pass through the full address
// GetFullDenom is a function, not method, so the message_plugin can use it
// func GetDenomData(contract string, subDenom string) (string, error) {
// 	// Address validation
// 	if _, err := parseAddress(contract); err != nil {
// 		return "", err
// 	}
// 	fullDenom, err := tokenfactorytypes.GetTokenDenom(contract, subDenom)
// 	if err != nil {
// 		return "", sdkerrors.Wrap(err, "validate sub-denom")
// 	}

// 	return fullDenom, nil
// }

// parseAddress parses address from bech32 string and verifies its format.
func parseAddress(addr string) (sdk.AccAddress, error) {
	parsed, err := sdk.AccAddressFromBech32(addr)
	if err != nil {
		return nil, sdkerrors.Wrap(err, "address from bech32")
	}
	err = sdk.VerifyAddressFormat(parsed)
	if err != nil {
		return nil, sdkerrors.Wrap(err, "verify address format")
	}
	return parsed, nil
}
