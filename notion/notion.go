package notion

import (
	"context"

	gn "github.com/dstotijn/go-notion"
	"github.com/maru44/scheman/definition"
	"github.com/volatiletech/sqlboiler/v4/drivers"
)

type (
	Notion struct {
		PageID             string
		TableListDBID      string
		cli                *gn.Client
		TablesByConnection []drivers.Table
	}
)

func NewNotion(pageID, tableListDBID, token string, tables []drivers.Table) definition.Definition {
	return &Notion{
		PageID:             pageID,
		TableListDBID:      tableListDBID,
		cli:                gn.NewClient(token),
		TablesByConnection: tables,
	}
}

func (n *Notion) GetCurrent(ctx context.Context) error {
	ls, err := n.getListTable(ctx)
	if err != nil {
		return err
	}
	tablesInNotionByName := map[string]definition.Table{}
	tablesDefinedInNotion := make([]definition.Table, len(ls))
	for i, l := range ls {
		t, err := n.getDefTable(ctx, l.PageID, l.TableName)
		if err != nil {
			return err
		}
		tablesDefinedInNotion[i] = *t
		tablesInNotionByName[l.TableName] = *t
	}
	tablesByConnection := n.TablesByConnection

	for _, tn := range tablesDefinedInNotion {
		// judge if the table exists in connection.
		existsInConnection := false
		for _, tc := range tablesByConnection {
			if tc.Name == tn.Name {
				existsInConnection = true
				break
			}
		}
		if !existsInConnection {
			// drop def table in notion
		}
	}

	for _, tc := range tablesByConnection {
		if tn, ok := tablesInNotionByName[tc.Name]; ok {
			columnNamesByConnection := map[string]int{}
			columnInNotionByColumnName := map[string]definition.Column{}

			for _, col := range tn.Columns {
				columnInNotionByColumnName[col.Name] = col
			}

			// loop for connection columns
			// update or create column in notion.
			for _, col := range tc.Columns {
				columnNamesByConnection[col.Name]++
				// If column name already exists in notion, update the row in notion.
				if _, ok := columnInNotionByColumnName[col.Name]; ok {
					currentColumn := definition.ConvertCol(col, tc.PKey)
					if err := n.updateDefRow(ctx, currentColumn); err != nil {
						return err
					}
					continue
				}

				// If column name does not exists in notion, insert row in notion.
				c := definition.ConvertCol(col, tc.PKey)
				if err := n.createDefRow(ctx, tn.PageID, c); err != nil {
					return err
				}
				continue
			}

			// loop for notion columns
			// If column name does not exists in notion,
			// delete the column in notion.
			for columnNameN, col := range columnInNotionByColumnName {
				if _, ok := columnNamesByConnection[columnNameN]; !ok {
					if err := n.deleteDefTable(ctx, col.RowID); err != nil {
						return err
					}
				}
			}
			continue
		}

		// If table does not exists,
		// insert table and insert columns as row.
		dbID, err := n.createDefTable(ctx, tc.Name)
		if err != nil {
			return err
		}
		for _, col := range tc.Columns {
			c := definition.ConvertCol(col, tc.PKey)
			if err := n.createDefRow(ctx, *dbID, c); err != nil {
				return err
			}
		}
	}

	return nil
}

func (n *Notion) Upsert(ctx context.Context) {

}
