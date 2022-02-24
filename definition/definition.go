package definition

import (
	"context"
)

type (
	Definition interface {
		Upsert(context.Context) error
	}

	Service string
)

const (
	ServiceNotion = Service("NOTION")
)
