package cli

import (
	"context"
	"fmt"
	"time"

	"github.com/spf13/cast"
	"github.com/spf13/cobra"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/client/tx"

	"github.com/notional-labs/eve/x/tokenfactory/types"
)

var (
	DefaultRelativePacketTimeoutTimestamp = uint64((time.Duration(10) * time.Minute).Nanoseconds())
)

const (
// flagPacketTimeoutTimestamp = "packet-timeout-timestamp"
)

// GetTxCmd returns the transaction commands for this module
func GetTxCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:                        types.ModuleName,
		Short:                      fmt.Sprintf("%s transactions subcommands", types.ModuleName),
		DisableFlagParsing:         true,
		SuggestionsMinimumDistance: 2,
		RunE:                       client.ValidateCmd,
	}

	// core
	cmd.AddCommand(CmdCreateDenom())
	cmd.AddCommand(CmdMintAndSendTokens())
	cmd.AddCommand(CmdUpdateOwner())
	cmd.AddCommand(CmdUpdateDenomCanChangeMaxSupply())

	// metadata
	cmd.AddCommand(CmdUpdateDenomDescription())
	cmd.AddCommand(CmdUpdateDenomTokenImage())
	cmd.AddCommand(CmdUpdateDenomWebsite())

	return cmd
}

func CmdCreateDenom() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "create-denom [name] [denom] [precision] [max-supply] [can-change-max-supply]",
		Short: "Create a new Denom",
		Args:  cobra.ExactArgs(5),
		RunE: func(cmd *cobra.Command, args []string) (err error) {
			name := args[0]
			indexDenom := args[1]
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
				name,
				indexDenom,
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

func CmdUpdateOwner() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "update-owner [denom] [new-owner]",
		Short: "Broadcast message UpdateOwner",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) (err error) {
			argDenom := args[0]
			argNewOwner := args[1]

			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			msg := types.NewMsgUpdateOwner(
				clientCtx.GetFromAddress().String(),
				argDenom,
				argNewOwner,
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

func CmdMintAndSendTokens() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "mint-and-send-tokens [denom] [amount] [recipient]",
		Short: "Broadcast message MintAndSendTokens",
		Args:  cobra.ExactArgs(3),
		RunE: func(cmd *cobra.Command, args []string) (err error) {
			argDenom := args[0]
			argAmount, err := cast.ToInt32E(args[1])
			if err != nil {
				return err
			}
			argRecipient := args[2]

			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			msg := types.NewMsgMintAndSendTokens(
				clientCtx.GetFromAddress().String(),
				argDenom,
				argAmount,
				argRecipient,
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

// is this even needed after init of a token?
func CmdUpdateDenomCanChangeMaxSupply() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "update-denom-change-supply [denom] [true/false]",
		Short: "Update a Denom's image via a URL",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) (err error) {
			argDenom := args[0]
			argCanChangeMaxSupply := args[1]

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

			canChange, err := cast.ToBoolE(argCanChangeMaxSupply)
			if err != nil {
				return err
			}

			msg := types.NewMsgUpdateDenom(
				clientCtx.GetFromAddress().String(),
				res.GetDenom().Denom,
				res.GetDenom().MaxSupply,
				canChange,
				res.GetDenom().Description,
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
