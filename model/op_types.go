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
	// OpBlockInvalid is an invalid OperationBlockType.
	OpBlockInvalid OperationBlockType = iota
	// OpBlockAdd type means add everithing is in that block to the previous one.
	OpBlockAdd
	// OpBlockDelete is a block that is erasing all previous information.
	// If a operation block has this type, the content must be ignored.
	OpBlockDelete
	// OpBlockCheckpoint is a block that has all information from previous block
	// squashed into one. You don't need to read ahead this block in the history
	OpBlockCheckpoint
)
