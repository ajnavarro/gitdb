package model

type Operation struct {
	Type  int    `json:"t"`
	Key   string `json:"k"`
	Value []byte `json:"v"`
}

type OperationBlock struct {
	Type int          `json:"t"`
	Ops  []*Operation `json:"ops"`
}
