package model

import (
	"encoding/json"
)

type OperationBlock struct {
	Type OperationBlockType `json:"t"`
	Ops  []*Operation       `json:"ops"`
}

func (pb *OperationBlock) Marshal(opType OperationBlockType) ([]byte, error) {
	return json.Marshal(pb)
}

type Operation struct {
	Type  OperationType `json:"t"`
	Field *Field        `json:"f"`
}

type OperationBlockIter interface {
	Next() (OperationBlock, error)
	Close() error
}
