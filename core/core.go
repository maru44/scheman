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
		Defs map[definition.Platform]definition.Definition
	}
)

func New(config *boilingcore.Config) (*SchemanState, error) {
	s := &SchemanState{
		State: &boilingcore.State{
			Config: config,
		},
		Defs: map[definition.Platform]definition.Definition{},
	}

	s.Driver = drivers.GetDriver(config.DriverName)
	if err := s.initDBInfo(config.DriverConfig); err != nil {
		return nil, errors.Wrap(err, "unable to initialize tables")
	}

	platform := viper.GetString("platform")
	if platform == definition.PlatformNotion {
		s.Defs[definition.PlatformNotion] = notion.NewNotion(
			viper.GetString("notion_page_id"),
			viper.GetString("notion_table_list_id"),
			viper.GetString("notion_token"),
			s.Tables,
			config.DriverName,
		)
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

func (s *SchemanState) AddDef(key definition.Platform, def definition.Definition) {
	s.Defs[key] = def
}
