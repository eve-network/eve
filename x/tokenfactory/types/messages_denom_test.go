package types

import (
	"testing"

	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/notional-labs/eve/testutil/sample"
	"github.com/stretchr/testify/require"
)

func TestMsgCreateDenom_ValidateBasic(t *testing.T) {
	tests := []struct {
		name string
		msg  MsgCreateDenom
		err  error
	}{
		{
			name: "invalid address",
			msg: MsgCreateDenom{
				Owner: "invalid_address",
			},
			err: sdkerrors.ErrInvalidAddress,
		}, {
			name: "valid address",
			msg: MsgCreateDenom{
				Owner: sample.AccAddress(),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.msg.ValidateBasic()
			if tt.err != nil {
				require.ErrorIs(t, err, tt.err)
				return
			}
			require.NoError(t, err)
		})
	}
}

func TestMsgUpdateDenom_ValidateBasic(t *testing.T) {
	tests := []struct {
		name string
		msg  MsgUpdateDenom
		err  error
	}{
		{
			name: "invalid address",
			msg: MsgUpdateDenom{
				Owner: "invalid_address",
			},
			err: sdkerrors.ErrInvalidAddress,
		}, {
			name: "valid address",
			msg: MsgUpdateDenom{
				Owner: sample.AccAddress(),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.msg.ValidateBasic()
			if tt.err != nil {
				require.ErrorIs(t, err, tt.err)
				return
			}
			require.NoError(t, err)
		})
	}
}
