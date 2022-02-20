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
	tablesByName := map[string]definition.Table{}
	tablesDefinedInNotion := make([]definition.Table, len(ls))
	for i, l := range ls {
		t, err := n.getDefTable(ctx, l.PageID, l.TableName)
		if err != nil {
			return err
		}
		tablesDefinedInNotion[i] = *t
		tablesByName[l.TableName] = *t
	}

	tablesByConnection := n.TablesByConnection
	for _, tc := range tablesByConnection {
		if tn, ok := tablesByName[tc.Name]; ok {
			// update all table rows
			for _, col := range tc.Columns {
				updated := false

				// update
				for _, colN := range tn.Columns {
					if col.Name == colN.Name {
						n.updateDefRow(ctx, colN)
						updated = true
						break
					}
				}

				if !updated {
					c := definition.ConvertCol(col, tc.PKey)
					n.createDefRow(ctx, tn.PageID, c)
				}
			}
			continue
		}
		// insert table
	}

	// for _, tn := range tablesDefinedInNotion {
	// }

	return nil
}

func (n *Notion) Upsert(ctx context.Context) {

}
