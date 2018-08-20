package model

type OperationType int

const (
	OpInvalid OperationType = iota
	OpOverride
	OpAppend
	OpDelete
)

type OperationBlockType int

const (
	OpBlockInvalid OperationBlockType = iota
	OpBlockAdd
	//TODO	OpBlockCheckpoint
	// TODO OpBlockDelete
)
