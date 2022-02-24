package definition

import (
	"context"
)

type (
	Definition interface {
		Upsert(context.Context) error
	}
)
