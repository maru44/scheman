package notion

import (
	"context"
	"encoding/json"

	gn "github.com/dstotijn/go-notion"
	"github.com/maru44/scheman/definition"
	"github.com/volatiletech/sqlboiler/v4/drivers"
)

var (
	initialRichText = gn.DatabaseProperty{
		RichText: &gn.EmptyMetadata{},
		Type:     gn.DBPropTypeRichText,
	}

	initialCheckbox = gn.DatabaseProperty{
		Checkbox: &gn.EmptyMetadata{},
		Type:     gn.DBPropTypeCheckbox,
	}

	initialDataTypeProperty = gn.DatabaseProperty{
		Type: gn.DBPropTypeSelect,
		Select: &gn.SelectMetadata{
			Options: []gn.SelectOptions{
				{Name: "uuid", Color: gn.ColorDefault},
				{Name: "integer", Color: gn.ColorBlue},
				{Name: "int", Color: gn.ColorBlue},
				{Name: "bigint", Color: gn.ColorBlue},
				{Name: "smallint", Color: gn.ColorBlue},
				{Name: "tinyint", Color: gn.ColorBlue},
				{Name: "float", Color: gn.ColorBlue},
				{Name: "money", Color: gn.ColorBlue},
				{Name: "double precision", Color: gn.ColorBlue},
				{Name: "numeric", Color: gn.ColorBlue},
				{Name: "pg_lsn", Color: gn.ColorBlue},
				{Name: "enum", Color: gn.ColorBrown},
				{Name: "character", Color: gn.ColorOrange},
				{Name: "char", Color: gn.ColorOrange},
				{Name: "varchar", Color: gn.ColorOrange},
				{Name: "inet", Color: gn.ColorOrange},
				{Name: "cidr", Color: gn.ColorOrange},
				{Name: "macaddr", Color: gn.ColorOrange},
				{Name: "tsquery", Color: gn.ColorOrange},
				{Name: "tsvector", Color: gn.ColorOrange},
				{Name: "tinytext", Color: gn.ColorRed},
				{Name: "text", Color: gn.ColorRed},
				{Name: "mediumtext", Color: gn.ColorRed},
				{Name: "longtext", Color: gn.ColorRed},
				{Name: "date", Color: gn.ColorPurple},
				{Name: "datetime", Color: gn.ColorPurple},
				{Name: "time", Color: gn.ColorPurple},
				{Name: "timestamp", Color: gn.ColorPurple},
				{Name: "timestamp with time zone", Color: gn.ColorPurple},
				{Name: "timestamp without time zone", Color: gn.ColorPurple},
				{Name: "interval", Color: gn.ColorPurple},
				{Name: "boolean", Color: gn.ColorGreen},
				{Name: "json", Color: gn.ColorPink},
				{Name: "jsonb", Color: gn.ColorPink},
				{Name: "bytea", Color: gn.ColorYellow},
				{Name: "binary", Color: gn.ColorYellow},
				{Name: "varbinary", Color: gn.ColorYellow},
				{Name: "tinyblob", Color: gn.ColorYellow},
				{Name: "blob", Color: gn.ColorYellow},
				{Name: "mediumblob", Color: gn.ColorYellow},
				{Name: "longblob", Color: gn.ColorYellow},
				{Name: "circle", Color: gn.ColorYellow},
				{Name: "line", Color: gn.ColorYellow},
				{Name: "lseg", Color: gn.ColorYellow},
				{Name: "path", Color: gn.ColorYellow},
				{Name: "point", Color: gn.ColorYellow},
				{Name: "box", Color: gn.ColorYellow},
				{Name: "polygon", Color: gn.ColorYellow},
				{Name: "txid_snapshot", Color: gn.ColorYellow},
				{Name: "xml", Color: gn.ColorYellow},
				{Name: "USER-DEFINED", Color: gn.ColorYellow},
				{Name: "ARRAYinteger", Color: gn.ColorGray},
				{Name: "ARRAYboolean", Color: gn.ColorGray},
				{Name: "ARRAYnumeric", Color: gn.ColorGray},
				{Name: "ARRAYbytea", Color: gn.ColorGray},
				{Name: "ARRAYjson", Color: gn.ColorGray},
				{Name: "ARRAYjsonb", Color: gn.ColorGray},
				{Name: "ARRAYcharacter varying", Color: gn.ColorGray},
			},
		},
	}

	allAttrs = []string{
		"Column Name", "Data Type", "Default", "PK", "Auto Generate",
		"Unique", "Null", "Comment", "Enum", "Free Entry",
	}
)

func (n *Notion) createDefTable(ctx context.Context, tableNameForNotion string) (*string, error) {
	params := gn.CreateDatabaseParams{
		ParentPageID: n.PageID,
		Title: []gn.RichText{
			{
				Text: &gn.Text{
					Content: tableNameForNotion,
				},
			},
		},
		Properties: gn.DatabaseProperties{
			"Column Name": gn.DatabaseProperty{
				Type:  gn.DBPropTypeTitle,
				Title: &gn.EmptyMetadata{},
			},
		},
	}
	if _, ok := n.IgnoreAttributes["Data Type"]; !ok {
		params.Properties["Data Type"] = initialDataTypeProperty
	}
	if _, ok := n.IgnoreAttributes["Default"]; !ok {
		params.Properties["Default"] = initialRichText
	}
	if _, ok := n.IgnoreAttributes["PK"]; !ok {
		params.Properties["PK"] = initialCheckbox
	}
	if _, ok := n.IgnoreAttributes["Auto Generate"]; !ok {
		params.Properties["Auto Generate"] = initialCheckbox
	}
	if _, ok := n.IgnoreAttributes["Unique"]; !ok {
		params.Properties["Unique"] = initialCheckbox
	}
	if _, ok := n.IgnoreAttributes["Null"]; !ok {
		params.Properties["Null"] = initialCheckbox
	}
	if _, ok := n.IgnoreAttributes["Comment"]; !ok {
		params.Properties["Comment"] = initialRichText
	}
	if _, ok := n.IgnoreAttributes["Free Entry"]; !ok {
		params.Properties["Free Entry"] = initialRichText
	}
	if _, ok := n.IgnoreAttributes["Enum"]; !ok {
		params.Properties["Enum"] = gn.DatabaseProperty{
			Type: gn.DBPropTypeMultiSelect,
			MultiSelect: &gn.SelectMetadata{
				Options: []gn.SelectOptions{},
			},
		}
	}
	db, err := n.cli.CreateDatabase(ctx, params)
	if err != nil {
		return nil, err
	}

	return &db.ID, nil
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

	var cols []definition.Column
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

		if len(columnProps.ColumnName.Title) == 0 || columnProps.DataType.Select.Name == "" {
			continue
		}

		col := definition.Column{
			RowID: r.ID,
			Name:  columnProps.ColumnName.Title[0].PlainText,
		}

		if _, ok := n.IgnoreAttributes["Data Type"]; !ok {
			col.DBType = columnProps.DataType.Select.Name
		}
		if _, ok := n.IgnoreAttributes["Unique"]; !ok && columnProps.Unique.Checkbox != nil && *columnProps.Unique.Checkbox {
			col.PK = true
		}
		if _, ok := n.IgnoreAttributes["Null"]; !ok && columnProps.Nullable.Checkbox != nil && *columnProps.Nullable.Checkbox {
			col.Nullable = true
		}
		if _, ok := n.IgnoreAttributes["Auto Generate"]; !ok && columnProps.AutoGen.Checkbox != nil && *columnProps.AutoGen.Checkbox {
			col.AutoGenerated = true
		}
		if _, ok := n.IgnoreAttributes["PK"]; !ok && columnProps.PK.Checkbox != nil && *columnProps.PK.Checkbox {
			col.PK = true
		}
		if _, ok := n.IgnoreAttributes["Default"]; !ok && len(columnProps.Default.RichText) != 0 {
			col.Default = columnProps.Default.RichText[0].PlainText
		}
		if _, ok := n.IgnoreAttributes["Comment"]; !ok && len(columnProps.Comment.RichText) != 0 {
			col.Comment = columnProps.Comment.RichText[0].PlainText
		}
		if _, ok := n.IgnoreAttributes["Enum"]; !ok && len(columnProps.Enum.MultiSelect) != 0 {
			ms := columnProps.Enum.MultiSelect
			enums := make([]string, len(ms))
			for i, e := range ms {
				enums[i] = e.Name
			}
			col.Enum = enums
		}
		if _, ok := n.IgnoreAttributes["Free Entry"]; !ok && len(columnProps.FreeText.RichText) != 0 {
			col.FreeText = columnProps.FreeText.RichText[0].PlainText
		}

		cols = append(cols, col)
	}
	table.Columns = cols

	return table, nil
}

func (n *Notion) tableNameForNotion(table drivers.Table) string {
	if table.IsView {
		if n.IsIgnoreView {
			return ""
		}
		return table.Name + " (View)"
	}
	return table.Name
}
