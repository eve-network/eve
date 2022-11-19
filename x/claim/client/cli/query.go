package cli

import (
	"context"
	"fmt"
	"strings"

	"github.com/spf13/cobra"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/version"

	// "github.com/cosmos/cosmos-sdk/client/flags"
	// sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/eve-network/eve/x/claim/types"
)

// GetQueryCmd returns the cli query commands for this module
func GetQueryCmd() *cobra.Command {
	// Group claim queries under a subcommand
	cmd := &cobra.Command{
		Use:                        types.ModuleName,
		Short:                      fmt.Sprintf("Querying commands for the %s module", types.ModuleName),
		DisableFlagParsing:         true,
		SuggestionsMinimumDistance: 2,
		RunE:                       client.ValidateCmd,
	}
	cmd.AddCommand(
		GetCmdQueryParams(),
		GetCmdQueryClaimRecord(),
		GetCmdQueryClaimable(),
	)
	// this line is used by starport scaffolding # 1

	return cmd
}

// GetCmdQueryParams implements a command to return the current claim
// parameters.
func GetCmdQueryParams() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "params",
		Short: "Query the current claims parameters",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}
			queryClient := types.NewQueryClient(clientCtx)

			params := &types.QueryParamsRequest{}
			res, err := queryClient.Params(context.Background(), params)
			if err != nil {
				return err
			}

			return clientCtx.PrintProto(&res.Params)
		},
	}

	flags.AddQueryFlagsToCmd(cmd)

	return cmd
}

// GetCmdQueryClaimRecord implements the query claim-records command.
func GetCmdQueryClaimRecord() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "claim-record [address]",
		Args:  cobra.ExactArgs(1),
		Short: "Query the claim record for an account.",
		Long: strings.TrimSpace(
			fmt.Sprintf(`Query the claim record for an account.
This contains an address' initial claimable amounts, and the completed actions.

Example:
$ %s query claim claim-record <address>
`,
				version.AppName,
			),
		),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}
			queryClient := types.NewQueryClient(clientCtx)
			// Query store
			res, err := queryClient.ClaimRecord(context.Background(), &types.QueryClaimRecordRequest{Address: args[0]})
			if err != nil {
				return err
			}
			return clientCtx.PrintObjectLegacy(res)
		},
	}
	flags.AddQueryFlagsToCmd(cmd)
	return cmd
}

// GetCmdQueryClaimable implements the query claimables command.
func GetCmdQueryClaimable() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "total-claimable [address]",
		Args:  cobra.ExactArgs(1),
		Short: "Query the total claimable amount remaining for an account.",
		Long: strings.TrimSpace(
			fmt.Sprintf(`Query the total claimable amount remaining for an account.
Example:
$ %s query claim total-claimable eve1258nuq58cz9tfcge5e5egeq69fdvdy7rxmjksa
`,
				version.AppName,
			),
		),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}
			queryClient := types.NewQueryClient(clientCtx)
			// Query store
			res, err := queryClient.Claimable(context.Background(), &types.QueryTotalClaimableRequest{
				Address: args[0],
			})
			if err != nil {
				return err
			}
			return clientCtx.PrintObjectLegacy(res)
		},
	}
	flags.AddQueryFlagsToCmd(cmd)
	return cmd
}
