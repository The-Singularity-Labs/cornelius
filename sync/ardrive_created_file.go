package sync

import (
	"fmt"
	"strconv"
)

type TxData struct {
	Created []File            `json:"created"`
	Tips    []Tip             `json:"tips"`
	Fees    map[string]string `json:"fees"`
}

type File struct {
	Type         string `json:"type"`
	EntityName   string `json:"entityName"`
	EntityId     string `json:"entityId"`
	DataTxId     string `json:"data_tx_id,omitempty"`
	MetadataTxId string `json:"metadata_tx_id,omitempty"`
	BundledIn    string `json:"bundledIn,omitempty"`
	SourceUri    string `json:"sourceUri,omitempty"`
}

type Tip struct {
	TxId      string `json:"txId"`
	Recipient string `json:"recipient"`
	Winston   string `json:"winston"`
}

func (tx TxData) EntityIds() []string {
	ids := []string{}
	for _, file := range tx.Created {
		if file.EntityId != "" {
			ids = append(ids, file.EntityId)
		}
	}

	return ids
}

func (tx TxData) TotalFees() (int64, error) {
	var total int64
	for _, feeString := range tx.Fees {
		fee, err := strconv.ParseInt(feeString, 10, 64)
		if err != nil {
			return 0, fmt.Errorf("unable to parse %q as int64: %w", feeString, err)
		}

		total += fee
	}

	return total, nil
}

func (tx TxData) EntityId() string {
	for _, f := range tx.Created {
		if f.EntityId != "" {
			return f.EntityId
		}
	}
	return ""
}
