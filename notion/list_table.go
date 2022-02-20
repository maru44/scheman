package notion

import (
	gn "github.com/dstotijn/go-notion"
)

type (
	listTable struct {
		TableName       string
		TableIDInNotion string
	}

	listTableResponse struct {
		TableName       gn.PageTitle `json:"Table Name"`
		TableIDInNotion struct {
			R []gn.RichText `json:"rich_text,omitempty"`
		} `json:"DB ID defining Table"`
	}
)

func (l *listTableResponse) toList() *listTable {
	if len(l.TableName.Title) == 0 || l.TableName.Title[0].PlainText == "" {
		return nil
	}
	if len(l.TableIDInNotion.R) == 0 {
		return &listTable{
			TableName: l.TableName.Title[0].PlainText,
		}
	}

	return &listTable{
		TableName:       l.TableName.Title[0].PlainText,
		TableIDInNotion: l.TableIDInNotion.R[0].PlainText,
	}
}
