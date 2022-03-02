package definition

import (
	"context"

	"github.com/volatiletech/sqlboiler/v4/drivers"
)

type (
	Definition interface {
		EnableMermaid()
		IsDefinition() bool
		IsMermaid() bool

		Upsert(context.Context) error
		Mermaid(context.Context) error
	}

	CommonInfo struct {
		TablesByConnection []drivers.Table
		DriverName         string
		IgnoreAttributes   map[string]int
		IsIgnoreView       bool
		RawMermaid         string
	}
)

func (c *CommonInfo) IsShownAttr(attr string) bool {
	_, ok := c.IgnoreAttributes[attr]
	return ok
}
