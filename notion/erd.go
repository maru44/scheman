package notion

import (
	"context"

	gn "github.com/dstotijn/go-notion"
)

func (n *Notion) createERD(ctx context.Context) (*string, error) {
	typ := "mermaid"
	p, err := n.cli.CreatePage(ctx, gn.CreatePageParams{
		ParentType: gn.ParentTypePage,
		ParentID:   n.PageID,
		Title: []gn.RichText{
			{
				Text: &gn.Text{
					Content: "ERD",
				},
			},
		},
		Children: []gn.Block{
			{
				Type: gn.BlockTypeCode,
				Code: &gn.Code{
					Language: &typ,
					RichTextBlock: gn.RichTextBlock{
						Text: []gn.RichText{
							{
								Text: &gn.Text{
									Content: n.RawMermaid,
								},
							},
						},
					},
				},
			},
		},
	})
	if err != nil {
		return nil, err
	}
	return &p.ID, err
}
