package types

import (
	"github.com/cosmos/cosmos-sdk/codec"
	cdctypes "github.com/cosmos/cosmos-sdk/codec/types"
	cryptocodec "github.com/cosmos/cosmos-sdk/crypto/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/msgservice"
	govv1beta1 "github.com/cosmos/cosmos-sdk/x/gov/types/v1beta1"
)

func RegisterLegacyAminoCodec(cdc *codec.LegacyAmino) {
	cdc.RegisterConcrete(&MsgInitialClaim{}, "claim/InitialClaim", nil)
	cdc.RegisterConcrete(&MsgClaim{}, "claim/Claim", nil)
	// this line is used by starport scaffolding # 2

	cdc.RegisterConcrete(&AirdropProposal{}, "claim/AirdropProposal", nil)
}

func RegisterInterfaces(registry cdctypes.InterfaceRegistry) {
	registry.RegisterImplementations((*sdk.Msg)(nil),
		&MsgInitialClaim{},
		&MsgClaim{},
	)
	registry.RegisterImplementations(
		(*govv1beta1.Content)(nil),
		&AirdropProposal{},
	)

	// this line is used by starport scaffolding # 3
	msgservice.RegisterMsgServiceDesc(registry, &_Msg_serviceDesc)
}

var (
	amino = codec.NewLegacyAmino()
	// ModuleCdc = codec.NewProtoCodec(cdctypes.NewInterfaceRegistry())
	ModuleCdc = codec.NewAminoCodec(amino)
)

func init() {
	RegisterLegacyAminoCodec(amino)
	cryptocodec.RegisterCrypto(amino)
	amino.Seal()
}
