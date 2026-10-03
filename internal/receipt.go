package internal

import (
	"encoding/json"
	"errors"

	"github.com/ethereum/go-ethereum/common"
)

var (
	ErrReceipt = errors.New("receipt")
)

type Receipt struct {
	Status            uint64
	CumulativeGasUsed uint64
	Logs              []*Log
	LogsBloom         []byte
	Type              uint64
	TransactionHash   []byte
	TransactionIndex  uint64
	BlockHash         []byte
	BlockNumber       uint64
	GasUsed           uint64
	EffectiveGasPrice uint64
	BlobGasPrice      uint64
	From              *common.Address
	To                *common.Address
	ContractAddress   *common.Address
}

//nolint:tagliatelle
type receiptJSON struct {
	Status            string          `json:"status"`
	CumulativeGasUsed string          `json:"cumulativeGasUsed"`
	Logs              []*Log          `json:"logs"`
	LogsBloom         string          `json:"logsBloom"`
	Type              string          `json:"type"`
	TransactionHash   string          `json:"transactionHash"`
	TransactionIndex  string          `json:"transactionIndex"`
	BlockHash         string          `json:"blockHash"`
	BlockNumber       string          `json:"blockNumber"`
	GasUsed           string          `json:"gasUsed"`
	EffectiveGasPrice string          `json:"effectiveGasPrice"`
	BlobGasPrice      string          `json:"blobGasPrice"`
	From              *common.Address `json:"from"`
	To                *common.Address `json:"to"`
	ContractAddress   *common.Address `json:"contractAddress"`
}

func (r *Receipt) MarshalJSON() ([]byte, error) {
	bytes, err := json.Marshal(receiptJSON{
		Status:            Uint64ToHexString(r.Status),
		CumulativeGasUsed: Uint64ToHexString(r.CumulativeGasUsed),
		Logs:              r.Logs,
		LogsBloom:         BytesToHexString(r.LogsBloom),
		Type:              Uint64ToHexString(r.Type),
		TransactionHash:   BytesToHexString(r.TransactionHash),
		TransactionIndex:  Uint64ToHexString(r.TransactionIndex),
		BlockHash:         BytesToHexString(r.BlockHash),
		BlockNumber:       Uint64ToHexString(r.BlockNumber),
		GasUsed:           Uint64ToHexString(r.GasUsed),
		EffectiveGasPrice: Uint64ToHexString(r.EffectiveGasPrice),
		BlobGasPrice:      Uint64ToHexString(r.BlobGasPrice),
		From:              r.From,
		To:                r.To,
		ContractAddress:   r.ContractAddress,
	})
	if err != nil {
		return nil, receiptError("marshaling receipt", err)
	}
	return bytes, nil
}

func (r *Receipt) UnmarshalJSON(data []byte) error {
	receipt := receiptJSON{}
	err := json.Unmarshal(data, &receipt)
	if err != nil {
		return receiptError("unmarshaling receipt", err)
	}

	status, err := ParseUint64FromHexString(receipt.Status)
	if err != nil {
		return receiptError("parsing status", err)
	}

	cumulativeGasUsed, err := ParseUint64FromHexString(receipt.CumulativeGasUsed)
	if err != nil {
		return receiptError("parsing cumulative gas used", err)
	}

	logsBloom, err := ParseBytesFromHexString(receipt.LogsBloom)
	if err != nil {
		return receiptError("parsing logs bloom", err)
	}

	rType, err := ParseUint64FromHexString(receipt.Type)
	if err != nil {
		return receiptError("parsing type", err)
	}

	transactionHash, err := ParseBytesFromHexString(receipt.TransactionHash)
	if err != nil {
		return receiptError("parsing transaction hash", err)
	}

	transactionIndex, err := ParseUint64FromHexString(receipt.TransactionIndex)
	if err != nil {
		return receiptError("parsing transaction index", err)
	}

	blockHash, err := ParseBytesFromHexString(receipt.BlockHash)
	if err != nil {
		return receiptError("parsing block hash", err)
	}

	blockNumber, err := ParseUint64FromHexString(receipt.BlockNumber)
	if err != nil {
		return receiptError("parsing block number", err)
	}

	gasUsed, err := ParseUint64FromHexString(receipt.GasUsed)
	if err != nil {
		return receiptError("parsing gas used", err)
	}

	effectiveGasPrice, err := ParseUint64FromHexString(receipt.EffectiveGasPrice)
	if err != nil {
		return receiptError("parsing effective gas price", err)
	}

	var blobGasPrice uint64
	if receipt.BlobGasPrice != "" {
		blobGasPrice, err = ParseUint64FromHexString(receipt.BlobGasPrice)
		if err != nil {
			return receiptError("parsing blob gas price", err)
		}
	}

	r.Status = status
	r.CumulativeGasUsed = cumulativeGasUsed
	r.Logs = receipt.Logs
	r.LogsBloom = logsBloom
	r.Type = rType
	r.TransactionHash = transactionHash
	r.TransactionIndex = transactionIndex
	r.BlockHash = blockHash
	r.BlockNumber = blockNumber
	r.GasUsed = gasUsed
	r.EffectiveGasPrice = effectiveGasPrice
	r.BlobGasPrice = blobGasPrice
	r.From = receipt.From
	r.To = receipt.To
	r.ContractAddress = receipt.ContractAddress
	return nil
}
