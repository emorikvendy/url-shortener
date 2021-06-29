// Code generated by gopkg.in/reform.v1. DO NOT EDIT.

package postgres

import (
	"fmt"
	"strings"

	"gopkg.in/reform.v1"
	"gopkg.in/reform.v1/parse"
)

type urlModelTableType struct {
	s parse.StructInfo
	z []interface{}
}

// Schema returns a schema name in SQL database ("").
func (v *urlModelTableType) Schema() string {
	return v.s.SQLSchema
}

// Name returns a view or table name in SQL database ("url").
func (v *urlModelTableType) Name() string {
	return v.s.SQLName
}

// Columns returns a new slice of column names for that view or table in SQL database.
func (v *urlModelTableType) Columns() []string {
	return []string{
		"id",
		"name",
		"link",
		"hash",
		"created_at",
		"updated_at",
	}
}

// NewStruct makes a new struct for that view or table.
func (v *urlModelTableType) NewStruct() reform.Struct {
	return new(urlModel)
}

// NewRecord makes a new record for that table.
func (v *urlModelTableType) NewRecord() reform.Record {
	return new(urlModel)
}

// PKColumnIndex returns an index of primary key column for that table in SQL database.
func (v *urlModelTableType) PKColumnIndex() uint {
	return uint(v.s.PKFieldIndex)
}

// urlModelTable represents url view or table in SQL database.
var urlModelTable = &urlModelTableType{
	s: parse.StructInfo{
		Type:    "urlModel",
		SQLName: "url",
		Fields: []parse.FieldInfo{
			{Name: "ID", Type: "int64", Column: "id"},
			{Name: "Name", Type: "string", Column: "name"},
			{Name: "Link", Type: "string", Column: "link"},
			{Name: "Hash", Type: "string", Column: "hash"},
			{Name: "CreatedAt", Type: "time.Time", Column: "created_at"},
			{Name: "UpdatedAt", Type: "time.Time", Column: "updated_at"},
		},
		PKFieldIndex: 0,
	},
	z: new(urlModel).Values(),
}

// String returns a string representation of this struct or record.
func (model urlModel) String() string {
	res := make([]string, 6)
	res[0] = "ID: " + reform.Inspect(model.ID, true)
	res[1] = "Name: " + reform.Inspect(model.Name, true)
	res[2] = "Link: " + reform.Inspect(model.Link, true)
	res[3] = "Hash: " + reform.Inspect(model.Hash, true)
	res[4] = "CreatedAt: " + reform.Inspect(model.CreatedAt, true)
	res[5] = "UpdatedAt: " + reform.Inspect(model.UpdatedAt, true)
	return strings.Join(res, ", ")
}

// Values returns a slice of struct or record field values.
// Returned interface{} values are never untyped nils.
func (model *urlModel) Values() []interface{} {
	return []interface{}{
		model.ID,
		model.Name,
		model.Link,
		model.Hash,
		model.CreatedAt,
		model.UpdatedAt,
	}
}

// Pointers returns a slice of pointers to struct or record fields.
// Returned interface{} values are never untyped nils.
func (model *urlModel) Pointers() []interface{} {
	return []interface{}{
		&model.ID,
		&model.Name,
		&model.Link,
		&model.Hash,
		&model.CreatedAt,
		&model.UpdatedAt,
	}
}

// View returns View object for that struct.
func (model *urlModel) View() reform.View {
	return urlModelTable
}

// Table returns Table object for that record.
func (model *urlModel) Table() reform.Table {
	return urlModelTable
}

// PKValue returns a value of primary key for that record.
// Returned interface{} value is never untyped nil.
func (model *urlModel) PKValue() interface{} {
	return model.ID
}

// PKPointer returns a pointer to primary key field for that record.
// Returned interface{} value is never untyped nil.
func (model *urlModel) PKPointer() interface{} {
	return &model.ID
}

// HasPK returns true if record has non-zero primary key set, false otherwise.
func (model *urlModel) HasPK() bool {
	return model.ID != urlModelTable.z[urlModelTable.s.PKFieldIndex]
}

// SetPK sets record primary key, if possible.
//
// Deprecated: prefer direct field assignment where possible: s.ID = pk.
func (model *urlModel) SetPK(pk interface{}) {
	reform.SetPK(model, pk)
}

// check interfaces
var (
	_ reform.View   = urlModelTable
	_ reform.Struct = (*urlModel)(nil)
	_ reform.Table  = urlModelTable
	_ reform.Record = (*urlModel)(nil)
	_ fmt.Stringer  = (*urlModel)(nil)
)

func init() {
	parse.AssertUpToDate(&urlModelTable.s, new(urlModel))
}
