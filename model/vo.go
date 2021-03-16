package model

////////////////
type Welcome struct {
	Contact string `json:"contact"`
	Job     string `json:"job"`
	Github  string `json:"github"`
}

////////////////

type BlockchainInfoResp struct {
	Chain         string `json:"chain"`
	Blocks        int    `json:"blocks"`
	Headers       int    `json:"headers"`
	BestBlockHash string `json:"bestBlockHash"`
	Difficulty    string `json:"difficulty"`
	MedianTime    int    `json:"medianTime"`
	Chainwork     string `json:"chainwork"`
}

type BlockInfoResp struct {
	Height         int    `json:"height"`
	BlockIdHex     string `json:"id"`
	PrevBlockIdHex string `json:"prev"`
	NextBlockIdHex string `json:"next"`
	MerkleRootHex  string `json:"merkle"`
	TxCount        int    `json:"ntx"`
	InSatoshi      int    `json:"inSatoshi"`
	OutSatoshi     int    `json:"outSatoshi"`
	CoinbaseOut    int    `json:"coinbaseOut"`
	BlockTime      int    `json:"timestamp"`
	Bits           int    `json:"bits"`
	BlockSize      int    `json:"size"`
}

type TxInfoResp struct {
	TxIdHex    string `json:"txid"`
	InCount    int    `json:"nIn"`
	OutCount   int    `json:"nOut"`
	TxSize     int    `json:"size"`
	LockTime   int    `json:"locktime"`
	InSatoshi  int    `json:"inSatoshi"`
	OutSatoshi int    `json:"outSatoshi"`
	BlockTime  int    `json:"timestamp"`
	Height     int    `json:"height"`
	BlockIdHex string `json:"blkid"`
	Idx        int    `json:"idx"`
}

type TxInSpentResp struct {
	Height   int    `json:"height"`
	TxIdHex  string `json:"txid"`
	Idx      int    `json:"idx"`
	UtxIdHex string `json:"utxid"`
	Vout     int    `json:"vout"`
}

type TxInResp struct {
	Height       int    `json:"height"`
	TxIdHex      string `json:"txid"`
	Idx          int    `json:"idx"`
	ScriptSigHex string `json:"scriptSig"`
	Sequence     int    `json:"sequence"`

	HeightTxo     int    `json:"heightTxo"`
	UtxIdHex      string `json:"utxid"`
	Vout          int    `json:"vout"`
	Address       string `json:"address"`
	GenesisHex    string `json:"genesis"`
	Satoshi       int    `json:"satoshi"`
	ScriptTypeHex string `json:"scriptType"`
	ScriptPkHex   string `json:"scriptPk"`
}

type TxOutResp struct {
	TxIdHex       string `json:"txid"`
	Vout          int    `json:"vout"`
	Address       string `json:"address"`
	GenesisHex    string `json:"genesis"`
	Satoshi       int    `json:"satoshi"`
	ScriptTypeHex string `json:"scriptType"`
	ScriptPkHex   string `json:"scriptPk"`
	Height        int    `json:"height"`
	Idx           int    `json:"idx"`
}

type TxOutHistoryResp struct {
	TxIdHex       string `json:"txid"`
	Vout          int    `json:"vout"`
	Address       string `json:"address"`
	GenesisHex    string `json:"genesis"`
	Satoshi       int    `json:"satoshi"`
	ScriptTypeHex string `json:"scriptType"`
	Height        int    `json:"height"`
	Idx           int    `json:"idx"`
	IOType        int    `json:"ioType"`
}

type TxOutStatusResp struct {
	TxIdHex       string `json:"txid"`
	Vout          int    `json:"vout"`
	Address       string `json:"address"`
	GenesisHex    string `json:"genesis"`
	Satoshi       int    `json:"satoshi"`
	ScriptTypeHex string `json:"scriptType"`
	ScriptPkHex   string `json:"scriptPk"`
	Height        int    `json:"height"`

	TxIdSpentHex string `json:"txidSpent"`
	HeightSpent  int    `json:"heightSpent"`
}
