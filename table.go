package gitdb

import (
	"github.com/ajnavarro/gitdb/git"
	"github.com/ajnavarro/gitdb/model"
)

type Table struct {
	DBName string
	Name   string
	repo   *git.Repository
}

func (t *Table) NewRow(fields []*model.Field, author *model.Author) (string, error) {
	opBlock := model.OperationBlock{}
	for _, f := range fields {
		opBlock.Ops = append(opBlock.Ops, &model.Operation{
			Field: f,
			Type:  model.OpOverride,
		})
	}

	data, err := opBlock.Marshal(model.OpBlockAdd)
	if err != nil {
		return "", err
	}

	return t.repo.NewRow(t.DBName, t.Name, data, author)
}

func (t *Table) UpdateRow(rowID string, opBlock *model.OperationBlock, author *model.Author) error {
	data, err := opBlock.Marshal(model.OpBlockAdd)
	if err != nil {
		return err
	}

	return t.repo.UpdateRow(rowID, t.DBName, t.Name, data, author)
}

func (t *Table) GetRow(rowID string) *model.Row {
	return nil
}

func (t *Table) GetRows() RowIterator {
	return nil
}

type RowIterator interface {
	Next() (*model.Row, error)
	Close() error
}
