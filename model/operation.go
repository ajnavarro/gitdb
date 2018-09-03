package model

import (
	"encoding/json"
)

type OperationBlock struct {
	Type OperationBlockType `json:"t"`
	Ops  []*Operation       `json:"ops"`
}

func (pb *OperationBlock) Marshal() ([]byte, error) {
	return json.Marshal(pb)
}

func UnmarshalOperationBlock(data []byte) (*OperationBlock, error) {
	op := &OperationBlock{}
	if err := json.Unmarshal(data, op); err != nil {
		return nil, err
	}

	return op, nil
}

type Operation struct {
	Type  OperationType `json:"t"`
	Field *Field        `json:"f"`
}
