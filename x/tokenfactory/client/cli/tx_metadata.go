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

// TODO: a lot of boilerplate code, need to move query -> its own function
// to return the Denom if no error. How to handle error OR value? Struct?
// TODO: Way to unpack values from a struct -> NewMsgUpdateDenom automatically?

func CmdUpdateDenomDescription() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "update-denom-desc [denom] [description]",
		Short: "Add/Update a Denom's description",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) (err error) {
			argDenom := args[0]
			argDescription := args[1]

			// query default denom data that we are not setting.
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}
			queryClient := types.NewQueryClient(clientCtx)
			params := &types.QueryGetDenomRequest{
				Denom: argDenom,
			}
			res, err := queryClient.Denom(context.Background(), params)
			if err != nil {
				return err
			}

			msg := types.NewMsgUpdateDenom(
				clientCtx.GetFromAddress().String(),
				res.GetDenom().Denom,
				res.GetDenom().MaxSupply,
				res.GetDenom().CanChangeMaxSupply,
				argDescription,
				res.GetDenom().TokenImage,
				res.GetDenom().Website,
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

func CmdUpdateDenomTokenImage() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "update-denom-image [denom] [image_url]",
		Short: "Update a Denom's image via a URL",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) (err error) {
			argDenom := args[0]
			argImageUrl := args[1]

			// query default denom data.
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}
			queryClient := types.NewQueryClient(clientCtx)
			params := &types.QueryGetDenomRequest{
				Denom: argDenom,
			}
			res, err := queryClient.Denom(context.Background(), params)
			if err != nil {
				return err
			}

			imageurl, err := cast.ToStringE(argImageUrl)
			if err != nil {
				return err
			}

			msg := types.NewMsgUpdateDenom(
				clientCtx.GetFromAddress().String(),
				res.GetDenom().Denom,
				res.GetDenom().MaxSupply,
				res.GetDenom().CanChangeMaxSupply,
				res.GetDenom().Description,
				imageurl,
				res.GetDenom().Website,
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

func CmdUpdateDenomWebsite() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "update-denom-website [denom] [website_link]",
		Short: "Update a Denom's image via a URL",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) (err error) {
			argDenom := args[0]
			argWebsiteLink := args[1]

			// query default denom data.
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}
			queryClient := types.NewQueryClient(clientCtx)
			params := &types.QueryGetDenomRequest{
				Denom: argDenom,
			}
			res, err := queryClient.Denom(context.Background(), params)
			if err != nil {
				return err
			}

			website, err := cast.ToStringE(argWebsiteLink)
			if err != nil {
				return err
			}

			msg := types.NewMsgUpdateDenom(
				clientCtx.GetFromAddress().String(),
				res.GetDenom().Denom,
				res.GetDenom().MaxSupply,
				res.GetDenom().CanChangeMaxSupply,
				res.GetDenom().Description,
				res.GetDenom().TokenImage,
				website,
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
