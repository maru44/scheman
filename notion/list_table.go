package notion

import (
	"context"
	"encoding/json"

	gn "github.com/dstotijn/go-notion"
	"github.com/friendsofgo/errors"
)

type (
	listTable struct {
		ID        string
		TableName string
		PageID    string
	}

	listTableResponse struct {
		// title
		TableName gn.DatabasePageProperty `json:"Table Name"`
		// rich_text
		PageID gn.DatabasePageProperty `json:"Page ID"`
	}
)

func (l *listTableResponse) toList() *listTable {
	if len(l.TableName.Title) == 0 || l.TableName.Title[0].PlainText == "" {
		return nil
	}
	if len(l.PageID.RichText) == 0 {
		return &listTable{
			TableName: l.TableName.Title[0].PlainText,
		}
	}

	return &listTable{
		TableName: l.TableName.Title[0].PlainText,
		PageID:    l.PageID.RichText[0].PlainText,
	}
}

func (n *Notion) createListRow(ctx context.Context, tableName, tableID string) error {
	if _, err := n.cli.CreatePage(ctx, gn.CreatePageParams{
		ParentType: gn.ParentTypeDatabase,
		ParentID:   n.TableListDBID,
		DatabasePageProperties: &gn.DatabasePageProperties{
			"Table Name": gn.DatabasePageProperty{
				Title: []gn.RichText{
					{
						Text: &gn.Text{
							Content: tableName,
						},
					},
				},
			},
			"Page ID": gn.DatabasePageProperty{
				RichText: []gn.RichText{
					{
						Text: &gn.Text{
							Content: tableID,
						},
					},
				},
			},
		},
	}); err != nil {
		return err
	}
	return nil
}

func (n *Notion) getListTable(ctx context.Context) ([]*listTable, error) {
	hasNext := true
	startCursor := ""
	var res []gn.Page
	for hasNext {
		q, err := n.cli.QueryDatabase(ctx, n.TableListDBID, &gn.DatabaseQuery{
			StartCursor: startCursor,
		})
		if err != nil {
			return nil, errors.Wrap(err, "failed to query")
		}
		res = append(res, q.Results...)

		hasNext = q.HasMore
		if q.NextCursor != nil {
			startCursor = *q.NextCursor
		}
	}

	var ls []*listTable
	for _, r := range res {
		j, err := json.Marshal(r.Properties)
		if err != nil {
			return nil, err
		}

		var props listTableResponse
		if err := json.Unmarshal(j, &props); err != nil {
			return nil, err
		}
		if l := props.toList(); l != nil {
			l.ID = r.ID
			ls = append(ls, l)
		}
	}
	return ls, nil
}

func (n *Notion) createListTable(ctx context.Context) (*string, error) {
	db, err := n.cli.CreateDatabase(ctx, gn.CreateDatabaseParams{
		ParentPageID: n.PageID,
		Title: []gn.RichText{
			{
				Text: &gn.Text{
					Content: "Table Index",
				},
			},
		},
		Properties: gn.DatabaseProperties{
			"Table Name": gn.DatabaseProperty{
				Type:  gn.DBPropTypeTitle,
				Title: &gn.EmptyMetadata{},
			},
			"Page ID": initialRichText,
		},
	})
	if err != nil {
		return nil, err
	}
	return &db.ID, nil
}
