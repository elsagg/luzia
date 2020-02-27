package domain

import (
	"os"

	"github.com/elsagg/schemaless-data-go"
)

// Domain is a domain model
type Domain struct {
	datasource *schemaless.DataSource
}

// GetCellLatest will return the last version of a data cell in a specified datasource
func (d *Domain) GetCellLatest(RowKey string, ColumnKey string) (dataCell *schemaless.DataCell, err error) {
	dataCell, err = d.datasource.GetCellLatest(RowKey, ColumnKey)
	return
}

// GetCell will return a data cell in a specified datasource
func (d *Domain) GetCell(RowKey string, ColumnKey string, RefKey int64) (dataCell *schemaless.DataCell, err error) {
	dataCell, err = d.datasource.GetCell(RowKey, ColumnKey, RefKey)
	return
}

// PutCell will save and return a data cell in a specified datasource
func (d *Domain) PutCell(RowKey string, ColumnKey string, Body interface{}) (dataCell *schemaless.DataCell, err error) {
	dataCell, err = d.datasource.PutCell(RowKey, ColumnKey, Body)
	return
}

// NewDomain will create a new domain model
func NewDomain(dataSourceName string) *Domain {
	conn := schemaless.NewConnection(&schemaless.ConnectionOptions{
		Host:       os.Getenv("DATABASE_HOST"),
		Port:       os.Getenv("DATABASE_PORT"),
		Username:   os.Getenv("DATABASE_USER"),
		Password:   os.Getenv("DATABASE_PASS"),
		Parameters: map[string]string{"parseTime": "true"},
	})

	ds := schemaless.NewDataSource(dataSourceName, conn)

	return &Domain{ds}
}
