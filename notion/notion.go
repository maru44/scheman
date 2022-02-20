package notion

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"

	gn "github.com/dstotijn/go-notion"
	"github.com/maru44/scheman/definition"
)

type (
	Notion struct {
		PageID        string
		TableListDBID string
		cli           *gn.Client
	}
)

func NewNotion(pageID, tableListDBID, token string) definition.Definition {
	return &Notion{
		PageID:        pageID,
		TableListDBID: tableListDBID,
		cli:           gn.NewClient(token),
	}
}

func (n *Notion) GetCurrent(ctx context.Context) {
	ls, err := n.getListTable(ctx)
	if err != nil {
		fmt.Println(err)
	}
	// fmt.Println(ls)
	for _, l := range ls {
		fmt.Println(l)
	}
}

func (n *Notion) Upsert() {}

func (n *Notion) getListTable(ctx context.Context) ([]*listTable, error) {
	hasNext := true
	startCursor := ""
	var res []gn.Page
	for hasNext {
		query := &gn.DatabaseQuery{
			StartCursor: startCursor,
		}
		q, err := n.cli.QueryDatabase(ctx, n.TableListDBID, query)
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
