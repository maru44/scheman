package definition

import (
	"context"
)

type (
	Definition interface {
		Upsert(context.Context) error
	}

	Service string

	CreateDefinitionHandler func() Definition
)

const (
	ServiceNotion = Service("NOTION")
	// ServiceSpreadSheat = "SpreadSheat"
)

var (
	DefByService = map[Service]func() Definition{}
)
