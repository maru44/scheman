package notion

import (
	"context"
	"errors"

	gn "github.com/dstotijn/go-notion"
	"github.com/fatih/color"
	"github.com/maru44/scheman/definition"
)

type (
	Notion struct {
		*definition.CommonInfo
		PageID       string
		TableIndexID string
		MermaidERDID string
		cli          *gn.Client
	}
)

func NewNotion(pageID, tableIndexID, token string, info *definition.CommonInfo) (definition.Definition, error) {
	if pageID == "" {
		return nil, errors.New("notion-page-id is not set")
	}
	if token == "" {
		return nil, errors.New("notion-token is not set")
	}
	return &Notion{
		PageID:       pageID,
		TableIndexID: tableIndexID,
		cli:          gn.NewClient(token),
		CommonInfo:   info,
	}, nil
}

func (n *Notion) SetMermaid(m string) {
	n.RawMermaid = m
}

func (n *Notion) Upsert(ctx context.Context) error {
	if n.TablesByConnection == nil {
		return nil
	}
	color.Green("Getting tables in Notion ...")

	newListTableID := ""
	if n.TableIndexID == "" {
		id, err := n.createListTable(ctx)
		if err != nil {
			return err
		}
		n.TableIndexID = *id
		newListTableID = *id
	}

	ls, err := n.getListTable(ctx)
	if err != nil {
		return err
	}
	color.Green("Success to get tables in Notion!")
	listTableIDByTableName := map[string]string{}
	for _, t := range ls {
		listTableIDByTableName[t.TableName] = t.ID
	}

	tablesInNotionByName := map[string]definition.Table{}
	var tablesDefinedInNotion []definition.Table
	for _, l := range ls {
		if l.TableName == "Mermaid ERD" {
			continue
		}
		t, err := n.getDefTable(ctx, l.PageID, l.TableName)
		if err != nil {
			return err
		}
		tablesDefinedInNotion = append(tablesDefinedInNotion, *t)
		tablesInNotionByName[l.TableName] = *t
	}
	tablesByConnection := n.TablesByConnection

	// for delete notion definition table
	// and update table attributes
	for _, tn := range tablesDefinedInNotion {
		// judge if the table exists in connection.
		existsInConnection := false
		for _, tc := range tablesByConnection {
			if n.tableNameForNotion(tc) == tn.Name {
				existsInConnection = true
				break
			}
		}
		if !existsInConnection {
			// drop def table in notion
			if err := n.deleteRowOrTable(ctx, tn.PageID); err != nil {
				return err
			}
			// drop from list table
			if err := n.deleteRowOrTable(ctx, listTableIDByTableName[tn.Name]); err != nil {
				return err
			}
			continue
		}

		// if exists in connection
		// update attributes if required
		updateAttrProps, err := n.updateAttrProps(ctx, tn.PageID)
		if err != nil {
			return err
		}
		if len(updateAttrProps) != 0 {
			if _, err := n.cli.UpdateDatabase(ctx, tn.PageID, gn.UpdateDatabaseParams{
				Properties: updateAttrProps,
			}); err != nil {
				return err
			}
		}
	}

	for _, tc := range tablesByConnection {
		tableNameForNotion := n.tableNameForNotion(tc)
		if tableNameForNotion == "" {
			continue
		}

		color.Green("Writing Notion Table: %s", tableNameForNotion)
		if tn, ok := tablesInNotionByName[tableNameForNotion]; ok {
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
				if colNotion, ok := columnInNotionByColumnName[col.Name]; ok {
					currentColumn := definition.ConvertCol(col, tc.PKey, n.DriverName)
					currentColumn.RowID = colNotion.RowID
					currentColumn.FreeText = colNotion.FreeText
					if err := n.updateDefRow(ctx, currentColumn); err != nil {
						return err
					}
					continue
				}

				// If column name does not exists in notion, insert row in notion.
				c := definition.ConvertCol(col, tc.PKey, n.DriverName)
				if err := n.createDefRow(ctx, tn.PageID, c); err != nil {
					return err
				}
			}

			// loop for notion columns
			// If column name does not exists in notion,
			// delete the column in notion.
			for columnNameN, col := range columnInNotionByColumnName {
				if _, ok := columnNamesByConnection[columnNameN]; !ok {
					if err := n.deleteRowOrTable(ctx, col.RowID); err != nil {
						return err
					}
				}
			}
			continue
		}

		// If table does not exists,
		// insert table and insert columns as row.
		dbID, err := n.createDefTable(ctx, tableNameForNotion)
		if err != nil {
			return err
		}
		for i := range tc.Columns {
			// for reverse
			col := tc.Columns[len(tc.Columns)-1-i]

			c := definition.ConvertCol(col, tc.PKey, n.DriverName)
			if err := n.createDefRow(ctx, *dbID, c); err != nil {
				return err
			}
		}
		if err := n.createListRow(ctx, tableNameForNotion, *dbID); err != nil {
			return err
		}
	}

	if newListTableID != "" {
		color.Yellow(
			"We created new Table Index Database.\nYou have to set following config.\n\nkey: notion-table-index\nvalue: %s",
			newListTableID,
		)
	}

	return nil
}

func (n *Notion) Mermaid(ctx context.Context) error {
	if n.RawMermaid == "" {
		return nil
	}

	newListTableID := ""
	if n.TableIndexID == "" {
		id, err := n.createListTable(ctx)
		if err != nil {
			return err
		}
		n.TableIndexID = *id
		newListTableID = *id
		color.Green("Success to get tables in Notion!")
	}

	ls, err := n.getListTable(ctx)
	if err != nil {
		return err
	}
	// if already exists
	for _, l := range ls {
		if l.TableName == "Mermaid ERD" {
			// drop mermaid ERD
			if err := n.deleteRowOrTable(ctx, l.PageID); err != nil {
				return err
			}
			// drop from list table
			if err := n.deleteRowOrTable(ctx, l.ID); err != nil {
				return err
			}
			break
		}
	}

	color.Green("Writing Notion Mermaid ERD: ERD")
	pID, err := n.createERD(ctx)
	if err != nil {
		return err
	}
	if err := n.createListRow(ctx, "Mermaid ERD", *pID); err != nil {
		return err
	}

	if newListTableID != "" {
		color.Yellow(
			"We created new Table Index Database.\nYou have to set following config.\n\nkey: notion-table-index\nvalue: %s",
			newListTableID,
		)
	}

	return nil
}

func (n *Notion) deleteRowOrTable(ctx context.Context, id string) error {
	if _, err := n.cli.DeleteBlock(ctx, id); err != nil {
		return err
	}
	return nil
}
