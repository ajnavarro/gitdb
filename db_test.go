package gitdb

import (
	"testing"

	"github.com/ajnavarro/gitdb/model"
)

func TestSimple(t *testing.T) {
	db, err := NewDB("/path/to/repo", "TESTDB")
	checkErr(err)

	table, err := db.Table("TESTTABLE")
	checkErr(err)

	cols := []*model.Field{
		{
			Key:   "c1",
			Value: []byte("column one content"),
		},
		{
			Key:   "c2",
			Value: []byte("column two content"),
		},
	}

	author := &model.Author{
		Email: "EMAIL",
		Name:  "NAME",
	}

	id, err := table.NewRow(cols, author)
	checkErr(err)

	ops := &model.OperationBlock{
		Type: model.OpBlockAdd,
		Ops: []*model.Operation{
			{
				Type: model.OpOverride,
				Field: &model.Field{
					Key:   "c1",
					Value: []byte("other value for column 1"),
				},
			},
		},
	}

	err = table.UpdateRow(id, ops, author)
	checkErr(err)
}

func checkErr(err error) {
	if err != nil {
		panic(err)
	}
}
