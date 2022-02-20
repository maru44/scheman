package notion

import (
	gn "github.com/dstotijn/go-notion"
)

type (
	listTable struct {
		TableName string
		PageID    string
	}

	listTableResponse struct {
		TableName gn.PageTitle `json:"Table Name"`
		PageID    struct {
			R []gn.RichText `json:"rich_text,omitempty"`
		} `json:"Page ID"`
	}
)

func (l *listTableResponse) toList() *listTable {
	if len(l.TableName.Title) == 0 || l.TableName.Title[0].PlainText == "" {
		return nil
	}
	if len(l.PageID.R) == 0 {
		return &listTable{
			TableName: l.TableName.Title[0].PlainText,
		}
	}

	return &listTable{
		TableName: l.TableName.Title[0].PlainText,
		PageID:    l.PageID.R[0].PlainText,
	}
}
