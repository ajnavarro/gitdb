package ops

import (
	"fmt"
	"io"

	"github.com/ajnavarro/gitdb/model"
)

// TODO use go-errors
// TODO take into account commit timestamps :S
// TODO add more block operations
// TODO try to eval operations in reverse order, to be able to generate fields traversing commit tree just one time
type Evaluator struct {
}

func (e *Evaluator) Eval(opsIter model.OperationBlockIter) ([]*model.Field, error) {
	fieldByName := make(map[string]*model.Field)
	for {
		block, err := opsIter.Next()
		if err == io.EOF {
			return getFields(fieldByName), nil
		}

		if err != nil {
			return nil, err
		}

		switch bt := block.Type; bt {
		case model.OpBlockAdd:
			for _, op := range block.Ops {
				switch t := op.Type; t {
				case model.OpOverride:
					fieldByName[op.Field.Key] = op.Field
				case model.OpAppend:
					f, ok := fieldByName[op.Field.Key]
					if !ok {
						fieldByName[op.Field.Key] = op.Field
						break
					}

					newField := &model.Field{
						Key:   op.Field.Key,
						Value: append(f.Value, op.Field.Value...),
					}

					fieldByName[op.Field.Key] = newField
				case model.OpDelete:
					delete(fieldByName, op.Field.Key)
				default:
					return nil, fmt.Errorf("unsupported operation type: %d", t)
				}
			}
		default:
			return nil, fmt.Errorf("unsupported block type: %d", bt)
		}
	}
}

func getFields(fieldsMap map[string]*model.Field) []*model.Field {
	out := make([]*model.Field, len(fieldsMap))
	for _, v := range fieldsMap {
		out = append(out, v)
	}

	return out
}
