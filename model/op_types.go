package model

type OperationType int

const (
	OpOverride OperationType = iota
	OpAdd
	OpDelete
)

type OperationBlockType int

const (
	OpBlockCheckpoint OperationBlockType = iota
	OpBlockAdd
	OpBlockDelete
)
