package ops

import (
	"fmt"
	"io"

	"github.com/ajnavarro/gitdb/model"
)

type OperationBlockIter interface {
	Next() (*model.OperationBlock, error)
	Close() error
}

// TODO use go-errors
// TODO take into account commit timestamps :S
type Evaluator struct {
	fields          map[string]*model.Field
	processedFields map[string]bool
	appending       map[string]bool
}

func (e *Evaluator) Resolve(iter OperationBlockIter) ([]*model.Field, error) {
	e.fields = make(map[string]*model.Field)
	e.processedFields = make(map[string]bool)
	e.appending = make(map[string]bool)

For:
	for {
		ob, err := iter.Next()
		if err == io.EOF {
			break
		}

		if err != nil {
			return nil, err
		}

		switch t := ob.Type; t {
		case model.OpBlockAdd:
			if err := e.resolveAdd(ob); err != nil {
				return nil, err
			}
		case model.OpBlockDelete:
			e.fields = make(map[string]*model.Field)
			break For
		default:
			return nil, fmt.Errorf("invalid operation block type: %d", t)
		}
	}

	return e.fieldsSlice(), nil
}

func (e *Evaluator) fieldsSlice() []*model.Field {
	var fields = make([]*model.Field, 0)
	for _, v := range e.fields {

		fields = append(fields, v)
	}

	return fields
}

func (e *Evaluator) resolveAdd(actual *model.OperationBlock) error {
	for _, o := range actual.Ops {
		if processed := e.processedFields[o.Field.Key]; processed {
			continue
		}

		switch t := o.Type; t {
		case model.OpOverride:
			f := o.Field
			if ap := e.appending[o.Field.Key]; ap {
				f = e.merge(e.fields[o.Field.Key], f)
				e.appending[o.Field.Key] = false
			}

			e.fields[o.Field.Key] = f
			e.processedFields[o.Field.Key] = true
		case model.OpAppend:
			oldField, ok := e.fields[o.Field.Key]
			if !ok {
				e.appending[o.Field.Key] = true
				e.fields[o.Field.Key] = o.Field
				continue
			}

			e.fields[o.Field.Key] = e.merge(oldField, o.Field)
		case model.OpDelete:
			e.processedFields[o.Field.Key] = true
			e.appending[o.Field.Key] = false
			delete(e.fields, o.Field.Key)
		default:
			return fmt.Errorf("invalid operation type. operation type ID: %d", t)
		}
	}

	return nil
}

func (e *Evaluator) merge(old, new *model.Field) *model.Field {
	return &model.Field{
		Key:   old.Key,
		Value: append(new.Value, old.Value...),
	}
}
