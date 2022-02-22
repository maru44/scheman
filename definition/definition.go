package definition

import (
	"context"
)

type (
	Definition interface {
		Upsert(context.Context) error
	}

	Platform string

	CreateDefinitionHandler func() Definition
)

const (
	PlatformNotion = "NOTION"
	// PlatformSpreadSheat = "SpreadSheat"
)

var (
	DefByPlatform = map[Platform]func() Definition{}
)
