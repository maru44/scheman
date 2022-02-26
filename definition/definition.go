package definition

import (
	"context"

	"github.com/volatiletech/sqlboiler/v4/drivers"
)

type (
	Definition interface {
		Upsert(context.Context) error
		SetMermaid(string)
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
