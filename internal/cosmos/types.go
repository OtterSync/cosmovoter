package cosmos

type ChainConfig struct {
	ChainID       string `json:"chain_id"`
	Executable    string `json:"executable"`
	WalletName    string `json:"wallet_name"`
	WalletAddress string `json:"wallet_address"`
	Password      string `json:"password"`
	Salt          string `json:"salt"`
	GasPrices     string `json:"gas_prices"`
	RPCNode       string `json:"rpc_node"`
}

type Proposal struct {
	ID          string `json:"id"`
	Title       string `json:"title"`
	Description string `json:"description"`
	Status      string `json:"status"`
}

type Vote struct {
	ProposalID string `json:"proposal_id"`
	Voter      string `json:"voter"`
	Option     string `json:"option"`
}
