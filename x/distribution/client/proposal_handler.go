package client

import (
	govclient "github.com/cosmos/cosmos-sdk/x/gov/client"
	"github.com/notional-labs/eve/x/distribution/client/cli"
)

// ProposalHandler is the community spend proposal handler.
var (
	ProposalHandler = govclient.NewProposalHandler(cli.GetCmdSubmitProposal)
)
