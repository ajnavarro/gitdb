package model

import "fmt"

// refs/gitdb/[DBName]/[TableName]/[rootElementHash]
const rowReference = "refs/gitdb/%s/%s/%s"

func RowReference(DBName, tableName, rootElementHash string) string {
	return fmt.Sprintf(rowReference, DBName, tableName, rootElementHash)
}
