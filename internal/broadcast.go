package internal

import (
	"errors"
	"fmt"

	"github.com/ethereum/go-ethereum/common"
)

var (
	ErrBroadcast = errors.New("broadcast")
)

type Broadcast struct {
	Transactions []*Transaction `json:"transactions"`
	Receipts     []*Receipt     `json:"receipts"`
	Timestamp    uint64         `json:"timestamp"`
	Chain        uint64         `json:"chain"`
	Commit       string         `json:"commit"`
}

func (b *Broadcast) GetContractAddress(name string) (*common.Address, error) {
	for _, tx := range b.Transactions {
		if tx.ContractName == name {
			return tx.ContractAddress, nil
		}
	}
	msg := fmt.Sprintf("getting contract address for '%s'", name)
	return nil, broadcastError(msg, nil)
}
