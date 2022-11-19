package client

import (
	govclient "github.com/cosmos/cosmos-sdk/x/gov/client"
	"github.com/eve-network/eve/x/claim/client/cli"
)

var ProposalHandlers = govclient.NewProposalHandler(cli.ProposalAirdropInitialCmd)
