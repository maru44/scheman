package notion

import (
	"context"
	"fmt"
)

func (n *Notion) getPage(ctx context.Context) error {
	page, err := n.cli.FindPageByID(ctx, n.PageID)
	if err != nil {
		return err
	}

	fmt.Println(page)
	return nil
}
