package keeper_test

import (
	"strconv"
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/stretchr/testify/require"

	keepertest "github.com/notional-labs/eve/testutil/keeper"
	"github.com/notional-labs/eve/x/tokenfactory/keeper"
	"github.com/notional-labs/eve/x/tokenfactory/types"
)

// Prevent strconv unused error
var _ = strconv.IntSize

func TestDenomMsgServerCreate(t *testing.T) {
	k, ctx := keepertest.TokenfactoryKeeper(t)
	srv := keeper.NewMsgServerImpl(*k)
	wctx := sdk.WrapSDKContext(ctx)
	owner := "A"
	for i := 0; i < 5; i++ {
		expected := &types.MsgCreateDenom{Owner: owner,
			Denom: strconv.Itoa(i),
		}
		_, err := srv.CreateDenom(wctx, expected)
		require.NoError(t, err)
		rst, found := k.GetDenom(ctx,
			expected.Denom,
		)
		require.True(t, found)
		require.Equal(t, expected.Owner, rst.Owner)
	}
}

func TestDenomMsgServerUpdate(t *testing.T) {
	owner := "A"

	for _, tc := range []struct {
		desc    string
		request *types.MsgUpdateDenom
		err     error
	}{
		{
			desc: "Completed",
			request: &types.MsgUpdateDenom{Owner: owner,
				Denom: strconv.Itoa(0),
			},
		},
		{
			desc: "Unauthorized",
			request: &types.MsgUpdateDenom{Owner: "B",
				Denom: strconv.Itoa(0),
			},
			err: sdkerrors.ErrUnauthorized,
		},
		{
			desc: "KeyNotFound",
			request: &types.MsgUpdateDenom{Owner: owner,
				Denom: strconv.Itoa(100000),
			},
			err: sdkerrors.ErrKeyNotFound,
		},
	} {
		t.Run(tc.desc, func(t *testing.T) {
			k, ctx := keepertest.TokenfactoryKeeper(t)
			srv := keeper.NewMsgServerImpl(*k)
			wctx := sdk.WrapSDKContext(ctx)
			expected := &types.MsgCreateDenom{Owner: owner,
				Denom: strconv.Itoa(0),
			}
			_, err := srv.CreateDenom(wctx, expected)
			require.NoError(t, err)

			_, err = srv.UpdateDenom(wctx, tc.request)
			if tc.err != nil {
				require.ErrorIs(t, err, tc.err)
			} else {
				require.NoError(t, err)
				rst, found := k.GetDenom(ctx,
					expected.Denom,
				)
				require.True(t, found)
				require.Equal(t, expected.Owner, rst.Owner)
			}
		})
	}
}
