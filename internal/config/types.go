package config

type Config struct {
	PasswordSalt string       `json:"password_salt"`
	Chains       []ChainConfig `json:"chains"`
}

type ChainConfig struct {
	Name           string `json:"name"`
	Executable      string `json:"executable"`
	WalletName      string `json:"wallet_name"`
	WalletPassword  string `json:"wallet_password"`
	GasPrices       string `json:"gas_prices"`
	ChainID         string `json:"chain_id"`
	RPCNode         string `json:"rpc_node"`
}
