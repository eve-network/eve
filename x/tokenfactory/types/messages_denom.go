package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

// MsgCreateDenom
var _ sdk.Msg = &MsgCreateDenom{}

func NewMsgCreateDenom(
	owner string,
	name string,
	denom string,
	precision int32,
	maxSupply int32,
	canChangeMaxSupply bool,

) *MsgCreateDenom {
	return &MsgCreateDenom{
		Owner:              owner,
		Name:               name,
		Denom:              "factory/" + owner + "/" + denom,
		Precision:          precision,
		MaxSupply:          maxSupply,
		CanChangeMaxSupply: canChangeMaxSupply,
	}
}

func (msg *MsgCreateDenom) Route() string {
	return RouterKey
}

func (msg *MsgCreateDenom) Type() string {
	return "CreateDenom"
}

func (msg *MsgCreateDenom) GetSigners() []sdk.AccAddress {
	owner, err := sdk.AccAddressFromBech32(msg.Owner)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{owner}
}

func (msg *MsgCreateDenom) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
}

func (msg *MsgCreateDenom) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(msg.Owner)
	if err != nil {
		// return sdkerrors.Wrapf(sdkerrors.ErrInvalidAddress, "invalid owner address (%s)", err)
		return sdkerrors.Wrap(sdkerrors.ErrInvalidAddress, "invalid address")
	}

	// tickerLength := len(msg.Denom) // ! we do factory/{addr}/{ticker} now
	// if tickerLength < 3 {
	// 	return sdkerrors.Wrapf(sdkerrors.ErrInvalidRequest, "Denom length must be at least 3 chars long")
	// }
	// if tickerLength > 10 {
	// 	return sdkerrors.Wrapf(sdkerrors.ErrInvalidRequest, "Denom length must be 10 chars long maximum")
	// }
	if msg.MaxSupply == 0 {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidRequest, "Max Supply must be greater than 0")
	}

	return nil
}

// MsgUpdateDenom
var _ sdk.Msg = &MsgUpdateDenom{}

func NewMsgUpdateDenom(
	owner string,
	denom string,
	maxSupply int32,
	canChangeMaxSupply bool,
	description string, // TODO: can we make a metadata struct & use somehow?
	token_image string,
	website string,
) *MsgUpdateDenom {
	return &MsgUpdateDenom{
		Owner:              owner,
		Denom:              denom,
		MaxSupply:          maxSupply,
		CanChangeMaxSupply: canChangeMaxSupply,
		Description:        description,
		TokenImage:         token_image,
		Website:            website,
	}
}

func (msg *MsgUpdateDenom) Route() string {
	return RouterKey
}

func (msg *MsgUpdateDenom) Type() string {
	return "UpdateDenom"
}

func (msg *MsgUpdateDenom) GetSigners() []sdk.AccAddress {
	owner, err := sdk.AccAddressFromBech32(msg.Owner)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{owner}
}

func (msg *MsgUpdateDenom) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
}

func (msg *MsgUpdateDenom) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(msg.Owner)
	if err != nil {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidAddress, "invalid owner address (%s)", err)
	}
	if msg.MaxSupply == 0 {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidRequest, "Max Supply must be greater than 0")
	}
	return nil
}
