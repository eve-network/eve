package bindings

import "cosmossdk.io/math"

type EveMsg struct {
	/// Contracts can create denoms, namespaced under the contract's address.
	/// A contract may create any number of independent sub-denoms.
	CreateDenom *CreateDenom `json:"create_denom,omitempty"`
	/// Contracts can change the admin of a denom that they are the admin of.
	ChangeOwner *ChangeOwner `json:"change_admin,omitempty"`
	/// Contracts can mint native tokens for an existing factory denom
	/// that they are the admin of.
	MintTokens *MintTokens `json:"mint_tokens,omitempty"`
}

// / CreateDenom creates a new factory denom, of denomination:
// / factory/{creating contract address}/{Subdenom}
// / Subdenom can be of length at most 44 characters, in [0-9a-zA-Z./]
// / The (creating contract address, subdenom) pair must be unique.
// / The created denom's admin is the creating contract address,
// / but this admin can be changed using the ChangeAdmin binding.
type CreateDenom struct {
	Name               string `json:"name"`
	Subdenom           string `json:"subdenom"`
	Precision          int32  `json:"precision"`
	MaxSupply          int32  `json:"max_supply"`
	Supply             int32  `json:"supply"`
	CanChangeMaxSupply bool   `json:"can_change_max_supply"`
	Owner              bool   `json:"owner"`
	// TODO: metadata in another Call? Require all?
}

// / ChangeAdmin changes the admin for a factory denom.
// / If the NewAdminAddress is empty, the denom has no admin.
type ChangeOwner struct {
	Denom           string `json:"denom"` // factory/{contract address}/{subdenom}
	NewOwnerAddress string `json:"new_owner"`
}

type MintTokens struct {
	Denom     string   `json:"denom"`
	Amount    math.Int `json:"amount"`
	Recipient string   `json:"recipient"`
}
