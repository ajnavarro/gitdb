package gitdb

import (
	"github.com/ajnavarro/gitdb/git"
	"github.com/ajnavarro/gitdb/model"
	"github.com/ajnavarro/gitdb/ops"
)

type Table struct {
	DBName string
	Name   string
	repo   *git.Repository
	ev     *ops.Evaluator
}

func NewTable(name, DBName string, r *git.Repository) *Table {
	return &Table{
		Name:   name,
		DBName: DBName,
		repo:   r,
		ev:     &ops.Evaluator{},
	}
}

func (t *Table) NewRow(fields []*model.Field, author *model.Author) (string, error) {
	opBlock := model.OperationBlock{Type: model.OpBlockAdd}
	for _, f := range fields {
		opBlock.Ops = append(opBlock.Ops, &model.Operation{
			Field: f,
			Type:  model.OpOverride,
		})
	}

	data, err := opBlock.Marshal()
	if err != nil {
		return "", err
	}

	return t.repo.NewRow(t.DBName, t.Name, data, author)
}

func (t *Table) UpdateRow(rowID string, opBlock *model.OperationBlock, author *model.Author) error {
	data, err := opBlock.Marshal()
	if err != nil {
		return err
	}

	return t.repo.UpdateRow(rowID, t.DBName, t.Name, data, author)
}

func (t *Table) GetRow(rowID string) (*model.Row, error) {
	opIter, err := t.repo.GetOperationBlocks(rowID, t.DBName, t.Name)
	if err != nil {
		return nil, err
	}

	fields, err := t.ev.Resolve(opIter)
	if err != nil {
		return nil, err
	}

	return &model.Row{ID: rowID, Fields: fields}, nil
}

func (t *Table) GetRows() (RowIterator, error) {
	return nil, nil
}

type RowIterator interface {
	Next() (*model.Row, error)
	Close() error
}
