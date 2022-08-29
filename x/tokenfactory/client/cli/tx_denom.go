package cli

import (
	"context"

	"github.com/spf13/cast"
	"github.com/spf13/cobra"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/client/tx"
	"github.com/notional-labs/eve/x/tokenfactory/types"
)

func CmdCreateDenom() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "create-denom [denom] [description] [precision] [max-supply] [can-change-max-supply]",
		Short: "Create a new Denom",
		Args:  cobra.ExactArgs(5),
		RunE: func(cmd *cobra.Command, args []string) (err error) {
			// Get indexes
			indexDenom := args[0]

			// Get value arguments
			argDescription := args[1]
			argPrecision, err := cast.ToInt32E(args[2])
			if err != nil {
				return err
			}
			argMaxSupply, err := cast.ToInt32E(args[3])
			if err != nil {
				return err
			}
			argCanChangeMaxSupply, err := cast.ToBoolE(args[4])
			if err != nil {
				return err
			}

			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			msg := types.NewMsgCreateDenom(
				clientCtx.GetFromAddress().String(),
				indexDenom,
				argDescription,
				argPrecision,
				argMaxSupply,
				argCanChangeMaxSupply,
			)
			if err := msg.ValidateBasic(); err != nil {
				return err
			}
			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}

	flags.AddTxFlagsToCmd(cmd)

	return cmd
}

func CmdUpdateDenomDescription() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "update-denom-desc [denom] [description]",
		Short: "Update a Denom's description",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) (err error) {
			indexDenom := args[0]
			argDescription := args[1]

			// query default denom data that we are not setting.
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			queryClient := types.NewQueryClient(clientCtx)
			argDenom := args[0]

			params := &types.QueryGetDenomRequest{
				Denom: argDenom,
			}

			res, err := queryClient.Denom(context.Background(), params)
			if err != nil {
				return err
			}

			// argMaxSupply, err := cast.ToInt32E(args[2])
			// if err != nil {
			// 	return err
			// }
			// argCanChangeMaxSupply, err := cast.ToBoolE(args[3])
			// if err != nil {
			// 	return err
			// }

			// query maxSupply from the denom

			msg := types.NewMsgUpdateDenom(
				clientCtx.GetFromAddress().String(),
				indexDenom,
				argDescription,
				res.GetDenom().MaxSupply,
				res.GetDenom().CanChangeMaxSupply,
			)
			if err := msg.ValidateBasic(); err != nil {
				return err
			}
			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}

	flags.AddTxFlagsToCmd(cmd)

	return cmd
}
