package core

import (
	"testing"

	"github.com/maru44/scheman/definition"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
)

func Test_NewSchemanState(t *testing.T) {
	tests := []struct {
		name            string
		driverName      string
		envs            map[string]interface{}
		wantIsNotionNil bool
		wantIsFileNil   bool
		wantIsNotionDef bool
		wantIsNotionERD bool
		wantIsFileDef   bool
		wantIsFileERD   bool
	}{
		{
			name: "notion def + erd",
			envs: map[string]interface{}{
				"services":       []string{string(ServiceNotion)},
				"erd-outputs":    []string{string(ServiceNotion)},
				"notion-page-id": "a",
				"notion-token":   "a",
				"erd-file":       "a",
			},
			wantIsNotionNil: false,
			wantIsFileNil:   true,
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

			schemanState := &SchemanState{
				Defs: map[Service]definition.Definition{},
			}
			ss, err := schemanState.setDefinition(&definition.CommonInfo{})
			assert.NoError(t, err)
			for _, s := range Services {
				d := ss.Defs[s]
				if s == ServiceNotion {
					assert.Equal(t, tt.wantIsNotionNil, d == nil)
					if d != nil {
						assert.Equal(t, tt.wantIsNotionDef, d.IsDefinition())
						assert.Equal(t, tt.wantIsNotionERD, d.IsMermaid())
					}
				}

				if s == ServiceFile {
					assert.Equal(t, tt.wantIsFileNil, d == nil)
					if d != nil {
						assert.Equal(t, tt.wantIsFileDef, d.IsDefinition())
						assert.Equal(t, tt.wantIsFileERD, d.IsMermaid())
					}
				}
			}
		})
	}
}
