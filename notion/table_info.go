package notion

import (
	"context"
	"encoding/json"

	gn "github.com/dstotijn/go-notion"
	"github.com/maru44/scheman/definition"
)

func (n *Notion) createDefTable(ctx context.Context, tableName string) (*string, error) {
	db, err := n.cli.CreateDatabase(ctx, gn.CreateDatabaseParams{
		ParentPageID: n.PageID,
		Title: []gn.RichText{
			{
				Text: &gn.Text{
					Content: tableName,
				},
			},
		},
		Properties: gn.DatabaseProperties{
			"Column Name": gn.DatabaseProperty{
				Title: &gn.EmptyMetadata{},
			},
			"Data Type": gn.DatabaseProperty{
				Select: &gn.SelectMetadata{
					Options: []gn.SelectOptions{
						{Name: "uuid", Color: gn.ColorDefault},
						{Name: "int", Color: gn.ColorBlue},
						{Name: "bigint", Color: gn.ColorBlue},
						{Name: "smallint", Color: gn.ColorBlue},
						{Name: "tinyint", Color: gn.ColorBlue},
						{Name: "float", Color: gn.ColorBlue},
						{Name: "numeric", Color: gn.ColorBlue},
						{Name: "char", Color: gn.ColorOrange},
						{Name: "varchar", Color: gn.ColorOrange},
						{Name: "tinytext", Color: gn.ColorOrange},
						{Name: "text", Color: gn.ColorOrange},
						{Name: "mediumtext", Color: gn.ColorOrange},
						{Name: "longtext", Color: gn.ColorOrange},
						{Name: "date", Color: gn.ColorPurple},
						{Name: "datetime", Color: gn.ColorPurple},
						{Name: "time", Color: gn.ColorPurple},
						{Name: "timestamp", Color: gn.ColorPurple},
						{Name: "boolean", Color: gn.ColorGreen},
						{Name: "json", Color: gn.ColorPink},
						{Name: "jsonb", Color: gn.ColorPink},
					},
				},
			},
			"Default": gn.DatabaseProperty{
				RichText: &gn.EmptyMetadata{},
			},
			"PK": gn.DatabaseProperty{
				Checkbox: &gn.EmptyMetadata{},
			},
			"Auto Generate": gn.DatabaseProperty{
				Checkbox: &gn.EmptyMetadata{},
			},
			"Unique": gn.DatabaseProperty{
				Checkbox: &gn.EmptyMetadata{},
			},
			"Null": gn.DatabaseProperty{
				Checkbox: &gn.EmptyMetadata{},
			},
			"Comment": gn.DatabaseProperty{
				RichText: &gn.EmptyMetadata{},
			},
			"Free Entry": gn.DatabaseProperty{
				RichText: &gn.EmptyMetadata{},
			},
		},
	})
	if err != nil {
		return nil, err
	}

	return &db.ID, nil
}

func (n *Notion) deleteDefTable(ctx context.Context, tableID string) error {
	if _, err := n.cli.DeleteBlock(ctx, tableID); err != nil {
		return err
	}
	return nil
}

func (n *Notion) getDefTable(ctx context.Context, tableID, tableName string) (*definition.Table, error) {
	// request to notion api
	hasNext := true
	startCursor := ""
	var res []gn.Page
	for hasNext {
		q, err := n.cli.QueryDatabase(ctx, tableID, &gn.DatabaseQuery{
			StartCursor: startCursor,
		})
		if err != nil {
			return nil, err
		}
		res = append(res, q.Results...)

		hasNext = q.HasMore
		if q.NextCursor != nil {
			startCursor = *q.NextCursor
		}
	}

	table := &definition.Table{
		PageID: tableID,
		Name:   tableName,
	}

	// convert notion response to []sqlboiler.drivers.table
	for _, r := range res {
		j, err := json.Marshal(r.Properties)
		if err != nil {
			return nil, err
		}

		var columnProps columnProps
		if err := json.Unmarshal(j, &columnProps); err != nil {
			return nil, err
		}

		if len(columnProps.ColumnName.Title) == 0 || len(columnProps.DataType.RichText) == 0 {
			continue
		}

		col := definition.Column{
			RowID:         r.ID,
			Name:          columnProps.ColumnName.Title[0].PlainText,
			DBType:        columnProps.DataType.Select.Name,
			Unique:        *columnProps.Unique.Checkbox,
			Nullable:      *columnProps.Nullable.Checkbox,
			AutoGenerated: *columnProps.AutoGen.Checkbox,
		}

		if columnProps.PK.Checkbox != nil && *columnProps.PK.Checkbox {
			col.PK = true
		}
		if len(columnProps.Default.RichText) != 0 {
			col.Default = columnProps.Default.RichText[0].PlainText
		}
		if len(columnProps.Comment.RichText) != 0 {
			col.Comment = columnProps.Comment.RichText[0].PlainText
		}
		if len(columnProps.FreeText.RichText) != 0 {
			col.FreeText = columnProps.FreeText.RichText[0].PlainText
		}

		table.Columns = append(table.Columns, col)
	}

	return table, nil
}