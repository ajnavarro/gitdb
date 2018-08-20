package gitdb

import (
	"encoding/json"

	"github.com/ajnavarro/gitdb/model"

	"github.com/ajnavarro/gitdb/git"
)

type Table struct {
	DBName string
	Name   string
	repo   *git.Repository
}

func (t *Table) NewRow(cols []*Column, author *model.Author) (string, error) {
	ops := Operations{}
	for _, c := range cols {
		ops = append(ops, &Operation{
			Column: c,
			Type:   model.OpOverride,
		})
	}

	data, err := t.opsToBytes(ops, model.OpBlockCheckpoint)
	if err != nil {
		return "", err
	}

	return t.repo.NewRow(t.DBName, t.Name, data, author)
}

func (t *Table) UpdateRow(rowID string, ops Operations, author *model.Author) error {
	data, err := t.opsToBytes(ops, model.OpBlockAdd)
	if err != nil {
		return err
	}

	return t.repo.UpdateRow(rowID, t.DBName, t.Name, data, author)
}

func (t *Table) GetRow(rowID string) *Row {
	return nil
}

func (t *Table) GetRows() RowIterator {
	return nil
}

type RowIterator interface {
	Next() (*Row, error)
	Close() error
}

type Column struct {
	Key   string
	Value []byte
}

type Row struct {
	ID      string
	Columns []*Column
}

type Operations []*Operation

type Operation struct {
	Type   model.OperationType
	Column *Column
}

func (t *Table) opsToBytes(ops Operations, opType model.OperationBlockType) ([]byte, error) {
	jOps := []*model.Operation{}
	for _, o := range ops {
		jOps = append(jOps, &model.Operation{
			Type:  int(o.Type),
			Key:   o.Column.Key,
			Value: o.Column.Value,
		})
	}

	jop := &model.OperationBlock{
		Type: int(opType),
		Ops:  jOps,
	}

	return json.Marshal(jop)
}
