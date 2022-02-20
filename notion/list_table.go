package notion

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"

	gn "github.com/dstotijn/go-notion"
)

type (
	listTable struct {
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
			if errors.Is(err, gn.ErrValidation) {
				fmt.Println("validation")
				// @TODO impl
				// create database for list table
			}
			fmt.Println(err) // @TODO delete
			return nil, err
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
			ls = append(ls, l)
		}
	}
	return ls, nil
}
