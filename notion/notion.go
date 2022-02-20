package notion

import (
	"context"
	"errors"
	"fmt"

	gn "github.com/dstotijn/go-notion"
	"github.com/maru44/scheman/definition"
)

type Notion struct {
	PageID        string
	TableListDBID string
	Token         string
}

func NewNotion(pageID, tableListDBID, token string) definition.Definition {
	return &Notion{
		PageID:        pageID,
		TableListDBID: tableListDBID,
		Token:         token,
	}
}

func (n *Notion) GetCurrent(ctx context.Context) {
	cli := gn.NewClient(n.Token)
	n.getListTable(ctx, cli)
}

func (n *Notion) Upsert() {}

func (n *Notion) getListTable(ctx context.Context, cli *gn.Client) error {
	db, err := cli.FindDatabaseByID(ctx, n.TableListDBID)
	if err != nil {
		if errors.Is(err, gn.ErrValidation) {
			fmt.Println("validation")
			// @TODO impl
			// create database for list table
		}
		fmt.Println(err)
		return err
	}

	fmt.Println("db", db)

	fmt.Println(db.Properties)
	return nil
}
