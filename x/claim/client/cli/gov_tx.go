package cli

import (
	"os"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/tx"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/gov/client/cli"
	govv1beta1 "github.com/cosmos/cosmos-sdk/x/gov/types/v1beta1"
	"github.com/eve-network/eve/x/claim/types"
	"github.com/spf13/cobra"
)

func ProposalAirdropInitialCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "initial-airdrop",
		Short: "Submit a proposal to initialize airdrop",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, _, _, deposit, err := getProposalInfo(cmd)
			if err != nil {
				return err
			}

			proposal, err := ParseAirdopInitialProposal(clientCtx.Codec, args[0])
			if err != nil {
				return err
			}

			content := types.NewInitialAirdropProposal(
				proposal.Title,
				proposal.Description,
				proposal.AirdropStartTime,
				proposal.DurationOfAirdrop,
				proposal.ClaimDenom,
				proposal.AirdropList,
			)

			if err = content.ValidateBasic(); err != nil {
				return err
			}

			msg, err := govv1beta1.NewMsgSubmitProposal(content, deposit, clientCtx.GetFromAddress())
			if err != nil {
				return err
			}
			if err = msg.ValidateBasic(); err != nil {
				return err
			}
			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}

	// proposal flags
	cmd.Flags().String(cli.FlagTitle, "", "Title of proposal")
	cmd.Flags().String(cli.FlagDescription, "", "Description of proposal")
	cmd.Flags().String(cli.FlagDeposit, "", "Deposit of proposal")
	return cmd
}

func getProposalInfo(cmd *cobra.Command) (client.Context, string, string, sdk.Coins, error) {
	clientCtx, err := client.GetClientTxContext(cmd)
	if err != nil {
		return client.Context{}, "", "", nil, err
	}

	proposalTitle, err := cmd.Flags().GetString(cli.FlagTitle)
	if err != nil {
		return clientCtx, proposalTitle, "", nil, err
	}

	proposalDescr, err := cmd.Flags().GetString(cli.FlagDescription)
	if err != nil {
		return client.Context{}, proposalTitle, proposalDescr, nil, err
	}

	depositArg, err := cmd.Flags().GetString(cli.FlagDeposit)
	if err != nil {
		return client.Context{}, proposalTitle, proposalDescr, nil, err
	}

	deposit, err := sdk.ParseCoinsNormalized(depositArg)
	if err != nil {
		return client.Context{}, proposalTitle, proposalDescr, deposit, err
	}

	return clientCtx, proposalTitle, proposalDescr, deposit, nil
}

func ParseAirdopInitialProposal(cdc codec.JSONCodec, proposalFile string) (types.AirdropProposal, error) {
	var proposal types.AirdropProposal

	contents, err := os.ReadFile(proposalFile)
	if err != nil {
		return proposal, err
	}

	err = cdc.UnmarshalJSON(contents, &proposal)
	if err != nil {
		return proposal, err
	}

	return proposal, nil
}
