package definition

import "context"

type (
	Definition interface {
		GetCurrent(context.Context)
		Upsert()
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