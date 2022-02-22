package core

import (
	"context"
	"strings"

	"github.com/friendsofgo/errors"
	"github.com/maru44/scheman/definition"
	"github.com/maru44/scheman/notion"
	"github.com/spf13/viper"
	"github.com/volatiletech/sqlboiler/v4/boilingcore"
	"github.com/volatiletech/sqlboiler/v4/drivers"
)

type (
	SchemanState struct {
		*boilingcore.State
		Defs map[definition.Service]definition.Definition
	}
)

func New(config *boilingcore.Config) (*SchemanState, error) {
	s := &SchemanState{
		State: &boilingcore.State{
			Config: config,
		},
		Defs: map[definition.Service]definition.Definition{},
	}

	s.Driver = drivers.GetDriver(config.DriverName)
	if err := s.initDBInfo(config.DriverConfig); err != nil {
		return nil, errors.Wrap(err, "unable to initialize tables")
	}

	services := viper.GetStringSlice("services")
	for _, service := range services {
		switch service {
		case string(definition.ServiceNotion):
			pageID := viper.GetString("notion-page-id")
			token := viper.GetString("notion-token")
			if pageID == "" {
				return nil, errors.New("notion-page-id is not set")
			}
			if token == "" {
				return nil, errors.New("notion-token is not set")
			}
			s.Defs[definition.ServiceNotion] = notion.NewNotion(
				pageID,
				viper.GetString("notion-table-list-id"),
				token,
				s.Tables,
				config.DriverName,
			)
		default:
			return nil, errors.New("Service have not been supported yet")
		}
	}

	return s, nil
}

func (s *SchemanState) initDBInfo(config map[string]interface{}) error {
	dbInfo, err := s.Driver.Assemble(config)
	if err != nil {
		return errors.Wrap(err, "unable to fetch table data")
	}

	if len(dbInfo.Tables) == 0 {
		return errors.New("no tables found in database")
	}

	if err := checkPKeys(dbInfo.Tables); err != nil {
		return err
	}

	s.Schema = dbInfo.Schema
	s.Tables = dbInfo.Tables
	s.Dialect = dbInfo.Dialect

	return nil
}

func checkPKeys(tables []drivers.Table) error {
	var missingPKey []string
	for _, t := range tables {
		if !t.IsView && t.PKey == nil {
			missingPKey = append(missingPKey, t.Name)
		}
	}

	if len(missingPKey) != 0 {
		return errors.Errorf("primary key missing in tables (%s)", strings.Join(missingPKey, ", "))
	}
	return nil
}

func (s *SchemanState) Run() error {
	ctx := context.Background()
	for _, def := range s.Defs {
		if err := def.Upsert(ctx); err != nil {
			return err
		}
	}

	return nil
}

func (s *SchemanState) AddDef(key definition.Service, def definition.Definition) {
	s.Defs[key] = def
}
