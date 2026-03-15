package internal

import (
	"encoding/json"
	"errors"

	"github.com/ethereum/go-ethereum/common"
)

var (
	ErrTransaction = errors.New("transaction")
)

type InnerTransaction struct {
	From    *common.Address
	Gas     uint64
	Value   uint64
	Input   []byte
	Nonce   uint64
	ChainID uint64
}

//nolint:tagliatelle
type innerJSON struct {
	From    *common.Address `json:"from"`
	Gas     string          `json:"gas"`
	Value   string          `json:"value"`
	Input   string          `json:"input"`
	Nonce   string          `json:"nonce"`
	ChainID string          `json:"chainId"`
}

func (i *InnerTransaction) MarshalJSON() ([]byte, error) {
	bytes, err := json.Marshal(innerJSON{
		From:    i.From,
		Gas:     Uint64ToHexString(i.Gas),
		Value:   Uint64ToHexString(i.Value),
		Input:   BytesToHexString(i.Input),
		Nonce:   Uint64ToHexString(i.Nonce),
		ChainID: Uint64ToHexString(i.ChainID),
	})
	if err != nil {
		return nil, transactionError("marshaling inner", err)
	}
	return bytes, nil
}

func (i *InnerTransaction) UnmarshalJSON(data []byte) error {
	inner := innerJSON{}
	err := json.Unmarshal(data, &inner)
	if err != nil {
		return transactionError("unmarshaling inner", err)
	}

	gas, err := ParseUint64FromHexString(inner.Gas)
	if err != nil {
		return transactionError("parsing gas", err)
	}

	value, err := ParseUint64FromHexString(inner.Value)
	if err != nil {
		return transactionError("parsing value", err)
	}

	input, err := ParseBytesFromHexString(inner.Input)
	if err != nil {
		return transactionError("parsing input", err)
	}

	nonce, err := ParseUint64FromHexString(inner.Nonce)
	if err != nil {
		return transactionError("parsing nonce", err)
	}

	chainID, err := ParseUint64FromHexString(inner.ChainID)
	if err != nil {
		return transactionError("parsing chain ID", err)
	}

	i.From = inner.From
	i.Gas = gas
	i.Value = value
	i.Input = input
	i.Nonce = nonce
	i.ChainID = chainID
	return nil
}

type Transaction struct {
	Hash            []byte
	TransactionType string
	ContractName    string
	ContractAddress *common.Address
	Inner           *InnerTransaction
	IsFixedGasLimit bool
}

//nolint:tagliatelle
type transactionJSON struct {
	Hash            string            `json:"hash"`
	TransactionType string            `json:"transactionType"`
	ContractName    string            `json:"contractName"`
	ContractAddress *common.Address   `json:"contractAddress"`
	Inner           *InnerTransaction `json:"transaction"`
	IsFixedGasLimit bool              `json:"isFixedGasLimit"`
}

func (t *Transaction) MarshalJSON() ([]byte, error) {
	bytes, err := json.Marshal(transactionJSON{
		Hash:            BytesToHexString(t.Hash),
		TransactionType: t.TransactionType,
		ContractName:    t.ContractName,
		ContractAddress: t.ContractAddress,
		Inner:           t.Inner,
		IsFixedGasLimit: t.IsFixedGasLimit,
	})
	if err != nil {
		return nil, transactionError("marshaling transaction", err)
	}
	return bytes, nil
}

func (t *Transaction) UnmarshalJSON(data []byte) error {
	transaction := transactionJSON{}
	err := json.Unmarshal(data, &transaction)
	if err != nil {
		return transactionError("unmarshaling transaction", err)
	}

	hash, err := ParseBytesFromHexString(transaction.Hash)
	if err != nil {
		return transactionError("parsing hash", err)
	}

	t.Hash = hash
	t.TransactionType = transaction.TransactionType
	t.ContractName = transaction.ContractName
	t.ContractAddress = transaction.ContractAddress
	t.Inner = transaction.Inner
	t.IsFixedGasLimit = transaction.IsFixedGasLimit
	return nil
}
