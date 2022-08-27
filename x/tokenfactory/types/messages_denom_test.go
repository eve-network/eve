package types

import (
	"testing"
)

// TODO: write test

// func TestMsgCreateDenom_ValidateBasic(t *testing.T) {
// 	tests := []struct {
// 		name string
// 		msg  MsgCreateDenom
// 		err  error
// 	}{
// 		{
// 			name: "invalid address",
// 			msg: MsgCreateDenom{
// 				Owner: "invalid_address",
// 			},
// 			err: sdkerrors.ErrInvalidAddress,
// 		}, {
// 			name: "valid address",
// 			msg: MsgCreateDenom{
// 				Owner: sample.AccAddress(),
// 			},
// 		},
// 	}
// 	for _, tt := range tests {
// 		t.Run(tt.name, func(t *testing.T) {
// 			// err := tt.msg.ValidateBasic()
// 			// if tt.err != nil {
// 			// 	fmt.Printf("err: %v\n\n", err)
// 			// 	require.ErrorIs(t, err, tt.err)
// 			// 	return
// 			// }
// 			// require.NoError(t, err)

// 			err := tt.msg.ValidateBasic()
// 			if tt.err != nil {
// 				fmt.Printf("err: %v->%v\n\n--", err, tt.err)
// 				require.ErrorIs(t, err, tt.err)
// 				return
// 			}
// 			require.NoError(t, err)

// 		})
// 	}
// }

func TestMsgUpdateDenom_ValidateBasic(t *testing.T) {
	// TODO:
}
