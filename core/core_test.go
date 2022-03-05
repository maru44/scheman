package core

import (
	"testing"

	"github.com/fatih/color"
	"github.com/maru44/scheman/definition"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
)

func Test_SetDefinitionToSchemanState(t *testing.T) {
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
				"services":       []string{"Notion"},
				"erd-outputs":    []string{"Notion"},
				"notion-page-id": "a",
				"notion-token":   "a",
			},
			wantIsNotionNil: false,
			wantIsFileNil:   true,
			wantIsNotionDef: true,
			wantIsNotionERD: true,
			wantIsFileDef:   false,
			wantIsFileERD:   false,
		},
		{
			name: "notion only erd",
			envs: map[string]interface{}{
				"services":       []string{},
				"erd-outputs":    []string{"notioN"},
				"notion-page-id": "a",
				"notion-token":   "a",
			},
			wantIsNotionNil: false,
			wantIsFileNil:   true,
			wantIsNotionDef: false,
			wantIsNotionERD: true,
			wantIsFileDef:   false,
			wantIsFileERD:   false,
		},
		{
			name: "file only def",
			envs: map[string]interface{}{
				"services":    []string{},
				"erd-outputs": []string{},
				"def-file":    "a",
			},
			wantIsNotionNil: true,
			wantIsFileNil:   false,
			wantIsNotionDef: false,
			wantIsNotionERD: false,
			wantIsFileDef:   true,
			wantIsFileERD:   false,
		},
		{
			name: "file only erd",
			envs: map[string]interface{}{
				"services":    []string{},
				"erd-outputs": []string{},
				"def-file":    "a",
				"erd-file":    "a",
			},
			wantIsNotionNil: true,
			wantIsFileNil:   false,
			wantIsNotionDef: false,
			wantIsNotionERD: false,
			wantIsFileDef:   true,
			wantIsFileERD:   true,
		},
		{
			name: "all",
			envs: map[string]interface{}{
				"services":       []string{"notion"},
				"erd-outputs":    []string{"NOTION"},
				"notion-page-id": "a",
				"notion-token":   "a",
				"def-file":       "a",
				"erd-file":       "a",
			},
			wantIsNotionNil: false,
			wantIsFileNil:   false,
			wantIsNotionDef: true,
			wantIsNotionERD: true,
			wantIsFileDef:   true,
			wantIsFileERD:   true,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			color.Red(tt.name)
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
