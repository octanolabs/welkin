package store

// State ...
// Main instance state
type State struct {
	Address  string `bson:"address" json:"address"`
	Mnemonic string `bson:"mnemonic" json:"mnemonic"`
	HdPath   string `bson:"hdPath" json:"hdPath"`
	ChainID  string `bson:"chainID" json:"chainID"`
	Node     string `bson:"node" json:"node"`
}
