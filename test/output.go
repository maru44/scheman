package test

import (
	"testing"

	"github.com/maru44/scheman/core"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
	"github.com/volatiletech/sqlboiler/v4/boilingcore"
)

func Test_NewSchemanState(t *testing.T) {
	tests := []struct {
		name            string
		envs            map[string]interface{}
		wantIsNotionDef bool
		wantIsNotionERD bool
		wantIsFileDef   bool
		wantIsFileERD   bool
	}{
		{
			name: "notion def + erd",
			envs: map[string]interface{}{
				"services":       []string{string(core.ServiceNotion)},
				"erd-outputs":    []string{string(core.ServiceNotion)},
				"notion-page-id": "a",
				"notion-token":   "a",
				"erd-file":       "a",
			},
			wantIsNotionDef: true,
			wantIsNotionERD: true,
			wantIsFileDef:   false,
			wantIsFileERD:   false,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			for k, v := range tt.envs {
				viper.Set(k, v)
			}
			schemanState, err := core.New(&boilingcore.Config{})
			assert.NoError(t, err)
			// for _, d := range schemanState.Defs {

			// }
		})
	}
}
