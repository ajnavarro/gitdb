package gitdb

import (
	"fmt"
	"io/ioutil"
	"testing"

	"github.com/ajnavarro/gitdb/model"
	"github.com/ajnavarro/gitdb/ops"
	"github.com/stretchr/testify/require"
	gogit "gopkg.in/src-d/go-git.v4"
)

var testCols = []*model.Field{
	{
		Key:   "c1",
		Value: []byte("column one content"),
	},
	{
		Key:   "c2",
		Value: []byte("column two content"),
	},
}

var testColsUpdate = []*model.Field{
	{
		Key:   "c1",
		Value: []byte("override content 1"),
	},
	{
		Key:   "c2",
		Value: []byte("column two contentdata added"),
	},
}

var testColsDelete = []*model.Field{
	{
		Key:   "c2",
		Value: []byte("column two contentdata added"),
	},
}

var testAuthor = &model.Author{
	Email: "EMAIL",
	Name:  "NAME",
}

func TestCreate(t *testing.T) {
	req := require.New(t)
	create(req)
}

func TestUpdate(t *testing.T) {
	req := require.New(t)
	path, rowID := create(req)
	println(rowID)
	read(req, path, rowID, testCols)

	ob := ops.NewBuilder().
		Override("c1", []byte("override content 1")).
		Append("c2", []byte("data added")).
		Build()

	update(req, path, rowID, ob)
	read(req, path, rowID, testColsUpdate)
}

func TestDelete(t *testing.T) {
	req := require.New(t)
	path, rowID := create(req)

	read(req, path, rowID, testCols)

	ob := ops.NewBuilder().
		Override("c1", []byte("override content 1")).
		Append("c2", []byte("data added")).
		Build()

	update(req, path, rowID, ob)
	read(req, path, rowID, testColsUpdate)

	ob = ops.NewBuilder().
		Delete("c1").
		Build()

	update(req, path, rowID, ob)
	read(req, path, rowID, testColsDelete)

	ob = ops.NewBuilder().DeleteRow()

	update(req, path, rowID, ob)
	read(req, path, rowID, []*model.Field{})
}

func BenchmarkAppendMultipleCols(b *testing.B) {
	req := require.New(b)
	path, rowID := create(req)
	for i := 0; i < b.N; i++ {
		builder := ops.NewBuilder()
		for j := 0; j < b.N; j++ {
			builder.Append(fmt.Sprintf("c%d ", i), []byte(fmt.Sprintf("%d ", i)))
		}

		update(req, path, rowID, builder.Build())
	}

	b.ResetTimer()

	db := createDatabase(req, path)
	table, err := db.Table("TABLE")
	req.NoError(err)

	row, err := table.GetRow(rowID)
	req.NoError(err)

	req.Equal(rowID, row.ID)
}

func BenchmarkAppendOneCol(b *testing.B) {
	req := require.New(b)
	path, rowID := create(req)
	for i := 0; i < b.N; i++ {
		builder := ops.NewBuilder()
		builder.Append("c1", []byte(fmt.Sprintf("%d ", i)))
		update(req, path, rowID, builder.Build())
	}

	b.ResetTimer()

	db := createDatabase(req, path)
	table, err := db.Table("TABLE")
	req.NoError(err)

	row, err := table.GetRow(rowID)
	req.NoError(err)

	req.Equal(rowID, row.ID)
}

func BenchmarkOverrideOneCol(b *testing.B) {
	req := require.New(b)
	path, rowID := create(req)
	for i := 0; i < b.N; i++ {
		builder := ops.NewBuilder()
		builder.Override("c1", []byte(fmt.Sprintf("%d ", i)))
		update(req, path, rowID, builder.Build())
	}

	b.ResetTimer()

	db := createDatabase(req, path)
	table, err := db.Table("TABLE")
	req.NoError(err)

	row, err := table.GetRow(rowID)
	req.NoError(err)

	req.Equal(rowID, row.ID)
}

func BenchmarkOverrideMultipleCols(b *testing.B) {
	req := require.New(b)
	path, rowID := create(req)
	for i := 0; i < b.N; i++ {
		builder := ops.NewBuilder()
		for j := 0; j < b.N; j++ {
			builder.Override(fmt.Sprintf("c%d ", i), []byte(fmt.Sprintf("%d ", i)))
		}

		update(req, path, rowID, builder.Build())
	}

	b.ResetTimer()

	db := createDatabase(req, path)
	table, err := db.Table("TABLE")
	req.NoError(err)

	row, err := table.GetRow(rowID)
	req.NoError(err)

	req.Equal(rowID, row.ID)
}

func create(req *require.Assertions) (path, rowID string) {
	path = createRepo(req)
	db := createDatabase(req, path)

	table, err := db.Table("TABLE")
	req.NoError(err)

	rowID, err = table.NewRow(testCols, testAuthor)
	req.NoError(err)
	req.NotEqual("", rowID)

	return
}

func read(req *require.Assertions, path, rowID string, cols []*model.Field) {
	db := createDatabase(req, path)
	table, err := db.Table("TABLE")
	req.NoError(err)

	row, err := table.GetRow(rowID)
	req.NoError(err)

	req.Equal(rowID, row.ID)
	req.Equal(cols, row.Fields)
}

func update(req *require.Assertions, path, rowID string, ob *model.OperationBlock) {
	db := createDatabase(req, path)
	table, err := db.Table("TABLE")
	req.NoError(err)

	err = table.UpdateRow(rowID, ob, testAuthor)
	req.NoError(err)
}

func createDatabase(req *require.Assertions, path string) *DB {
	db, err := NewDB(path, "TESTDB")
	req.NoError(err)

	return db
}

func createRepo(req *require.Assertions) string {
	path, err := ioutil.TempDir("", "repo-")
	req.NoError(err)

	_, err = gogit.PlainInit(path, false)
	req.NoError(err)

	return path
}
