package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

var _ sdk.Msg = &MsgMintAndSendTokens{}

func NewMsgMintAndSendTokens(owner string, denom string, amount int32, recipient string) *MsgMintAndSendTokens {
	return &MsgMintAndSendTokens{
		Owner:     owner,
		Denom:     denom,
		Amount:    amount,
		Recipient: recipient,
	}
}

func (msg *MsgMintAndSendTokens) Route() string {
	return RouterKey
}

func (msg *MsgMintAndSendTokens) Type() string {
	return "MintAndSendTokens"
}

func (msg *MsgMintAndSendTokens) GetSigners() []sdk.AccAddress {
	owner, err := sdk.AccAddressFromBech32(msg.Owner)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{owner}
}

func (msg *MsgMintAndSendTokens) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
}

func (msg *MsgMintAndSendTokens) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(msg.Owner)
	if err != nil {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidAddress, "invalid owner address (%s)", err)
	}
	return nil
}
