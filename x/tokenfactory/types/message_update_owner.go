package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

var _ sdk.Msg = &MsgUpdateOwner{}

func NewMsgUpdateOwner(owner string, denom string, newOwner string) *MsgUpdateOwner {
	return &MsgUpdateOwner{
		Owner:    owner,
		Denom:    denom,
		NewOwner: newOwner,
	}
}

func (msg *MsgUpdateOwner) Route() string {
	return RouterKey
}

func (msg *MsgUpdateOwner) Type() string {
	return "UpdateOwner"
}

func (msg *MsgUpdateOwner) GetSigners() []sdk.AccAddress {
	owner, err := sdk.AccAddressFromBech32(msg.Owner)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{owner}
}

func (msg *MsgUpdateOwner) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
}

func (msg *MsgUpdateOwner) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(msg.Owner)
	if err != nil {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidAddress, "invalid owner address (%s)", err)
	}
	return nil
}
