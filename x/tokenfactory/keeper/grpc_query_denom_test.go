package keeper_test

import (
	"strconv"
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/query"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	keepertest "github.com/notional-labs/eve/testutil/keeper"
	"github.com/notional-labs/eve/x/tokenfactory/types"
)

// Prevent strconv unused error
var _ = strconv.IntSize

func TestDenomQuerySingle(t *testing.T) {
	keeper, ctx := keepertest.TokenfactoryKeeper(t)
	wctx := sdk.WrapSDKContext(ctx)
	msgs := createNDenom(keeper, ctx, 2)
	for _, tc := range []struct {
		desc     string
		request  *types.QueryGetDenomRequest
		response *types.QueryGetDenomResponse
		err      error
	}{
		{
			desc: "First",
			request: &types.QueryGetDenomRequest{
				Denom: msgs[0].Denom,
			},
			response: &types.QueryGetDenomResponse{Denom: msgs[0]},
		},
		{
			desc: "Second",
			request: &types.QueryGetDenomRequest{
				Denom: msgs[1].Denom,
			},
			response: &types.QueryGetDenomResponse{Denom: msgs[1]},
		},
		{
			desc: "KeyNotFound",
			request: &types.QueryGetDenomRequest{
				Denom: strconv.Itoa(100000),
			},
			err: status.Error(codes.InvalidArgument, "not found"),
		},
		{
			desc: "InvalidRequest",
			err:  status.Error(codes.InvalidArgument, "invalid request"),
		},
	} {
		t.Run(tc.desc, func(t *testing.T) {
			response, err := keeper.Denom(wctx, tc.request)
			if tc.err != nil {
				require.ErrorIs(t, err, tc.err)
			} else {
				require.Equal(t, tc.response, response)
			}
		})
	}
}

func TestDenomQueryPaginated(t *testing.T) {
	keeper, ctx := keepertest.TokenfactoryKeeper(t)
	wctx := sdk.WrapSDKContext(ctx)
	msgs := createNDenom(keeper, ctx, 5)

	request := func(next []byte, offset, limit uint64, total bool) *types.QueryAllDenomRequest {
		return &types.QueryAllDenomRequest{
			Pagination: &query.PageRequest{
				Key:        next,
				Offset:     offset,
				Limit:      limit,
				CountTotal: total,
			},
		}
	}
	t.Run("ByOffset", func(t *testing.T) {
		step := 2
		for i := 0; i < len(msgs); i += step {
			resp, err := keeper.DenomAll(wctx, request(nil, uint64(i), uint64(step), false))
			require.NoError(t, err)
			require.LessOrEqual(t, len(resp.Denom), step)
			require.Subset(t, msgs, resp.Denom)
		}
	})
	t.Run("ByKey", func(t *testing.T) {
		step := 2
		var next []byte
		for i := 0; i < len(msgs); i += step {
			resp, err := keeper.DenomAll(wctx, request(next, 0, uint64(step), false))
			require.NoError(t, err)
			require.LessOrEqual(t, len(resp.Denom), step)
			require.Subset(t, msgs, resp.Denom)
			next = resp.Pagination.NextKey
		}
	})
	t.Run("Total", func(t *testing.T) {
		resp, err := keeper.DenomAll(wctx, request(nil, 0, 0, true))
		require.NoError(t, err)
		require.Equal(t, len(msgs), int(resp.Pagination.Total))
		require.ElementsMatch(t, msgs, resp.Denom)
	})
	t.Run("InvalidRequest", func(t *testing.T) {
		_, err := keeper.DenomAll(wctx, nil)
		require.ErrorIs(t, err, status.Error(codes.InvalidArgument, "invalid request"))
	})
}
