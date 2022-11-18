package types

import (
	"strings"
	"time"

	"github.com/cosmos/cosmos-sdk/types/errors"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"
	govv1beta1 "github.com/cosmos/cosmos-sdk/x/gov/types/v1beta1"
)

const (
	ProposalTypeInitialAirdrop = "InitialAirdrop"
)

var _ govv1beta1.Content = &AirdropProposal{}

func init() {
	govv1beta1.RegisterProposalType(ProposalTypeInitialAirdrop)
}

func NewInitialAirdropProposal(
	title string,
	description string,
	startTime time.Time,
	durationOfAirdrop time.Duration,
	denom string,
	airdropList []Claimable,
) *AirdropProposal {
	return &AirdropProposal{title, description, startTime, durationOfAirdrop, denom, airdropList}
}

func (p AirdropProposal) ProposalRoute() string { return RouterKey }

func (p AirdropProposal) ProposalType() string { return ProposalTypeInitialAirdrop }

func (p AirdropProposal) ValidateBasic() error {
	if err := validateProposalCommons(p.Title, p.Description); err != nil {
		return err
	}
	if err := validateDuration(p.DurationOfAirdrop); err != nil {
		return err
	}
	if len(p.AirdropList) == 0 {
		return errors.ErrInvalidRequest
	}
	return nil
}

func validateProposalCommons(title, description string) error {
	if strings.TrimSpace(title) != title {
		return sdkerrors.Wrap(govtypes.ErrInvalidProposalContent, "proposal title must not start/end with white spaces")
	}
	if len(title) == 0 {
		return sdkerrors.Wrap(govtypes.ErrInvalidProposalContent, "proposal title cannot be blank")
	}
	if len(title) > govv1beta1.MaxTitleLength {
		return sdkerrors.Wrapf(govtypes.ErrInvalidProposalContent, "proposal title is longer than max length of %d", govv1beta1.MaxTitleLength)
	}
	if strings.TrimSpace(description) != description {
		return sdkerrors.Wrap(govtypes.ErrInvalidProposalContent, "proposal description must not start/end with white spaces")
	}
	if len(description) == 0 {
		return sdkerrors.Wrap(govtypes.ErrInvalidProposalContent, "proposal description cannot be blank")
	}
	if len(description) > govv1beta1.MaxDescriptionLength {
		return sdkerrors.Wrapf(govtypes.ErrInvalidProposalContent, "proposal description is longer than max length of %d", govv1beta1.MaxDescriptionLength)
	}
	return nil
}
