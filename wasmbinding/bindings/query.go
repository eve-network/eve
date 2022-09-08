package bindings

// OsmosisQuery contains osmosis custom queries.
// See https://github.com/osmosis-labs/osmosis-bindings/blob/main/packages/bindings/src/query.rs
type EveQuery struct {
	/// Given a subdenom minted by a contract via `EveMsg::MintTokens`,
	/// returns the full denom as used by `BankMsg::Send`.
	FullDenom *FullDenom `json:"full_denom,omitempty"`
	/// Returns the admin of a denom, if the denom is a Token Factory denom.
	DenomAdmin *DenomAdmin `json:"denom_admin,omitempty"`
}

type FullDenom struct {
	CreatorAddr string `json:"creator_addr"`
	Subdenom    string `json:"subdenom"`
}

type DenomAdmin struct {
	Subdenom string `json:"subdenom"`
}

type DenomAdminResponse struct {
	Admin string `json:"admin"`
}

type FullDenomResponse struct {
	Name               string `json:"name"`
	Denom              string `json:"denom"`
	Precision          int32  `json:"precision"`
	MaxSupply          int32  `json:"max_supply"`
	Supply             int32  `json:"supply"`
	CanChangeMaxSupply bool   `json:"can_change_max_supply"`
	Owner              string `json:"owner"`

	// extra metadata some tokens may not have.
	Description string `json:"description"`
	TokenImage  string `json:"token_image"`
	Website     string `json:"website"`
}
