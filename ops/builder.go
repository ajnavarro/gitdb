package ops

import "github.com/ajnavarro/gitdb/model"

type Builder struct {
	ops []*model.Operation
}

func NewBuilder() *Builder {
	return &Builder{}
}

func (b *Builder) add(t model.OperationType, key string, value []byte) {
	f := &model.Field{Key: key, Value: value}
	b.ops = append(b.ops, &model.Operation{Field: f, Type: t})
}

func (b *Builder) Append(key string, value []byte) *Builder {
	b.add(model.OpAppend, key, value)
	return b
}

func (b *Builder) Override(key string, value []byte) *Builder {
	b.add(model.OpOverride, key, value)
	return b
}

func (b *Builder) Delete(key string) *Builder {
	b.add(model.OpDelete, key, nil)
	return b
}

func (b *Builder) Build() *model.OperationBlock {
	return &model.OperationBlock{
		Ops:  b.ops,
		Type: model.OpBlockAdd,
	}
}

func (b *Builder) DeleteRow() *model.OperationBlock {
	return &model.OperationBlock{
		Type: model.OpBlockDelete,
	}
}
