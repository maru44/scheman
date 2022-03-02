package core

import (
	"context"
	"strings"

	"github.com/friendsofgo/errors"
	"github.com/maru44/scheman/definition"
	"github.com/maru44/scheman/file"
	"github.com/maru44/scheman/notion"
	"github.com/spf13/viper"
	"github.com/volatiletech/sqlboiler/v4/boilingcore"
	"github.com/volatiletech/sqlboiler/v4/drivers"
)

type (
	SchemanState struct {
		*boilingcore.State
		IgnoreAttributes []string
		Defs             map[Service]definition.Definition
	}

	Service string
)

const (
	ServiceNotion = Service("NOTION")
	ServiceFile   = Service("FILE")
)

var Services = []Service{
	ServiceNotion, ServiceFile,
}

func (s *SchemanState) Run() error {
	ctx := context.Background()
	for _, def := range s.Defs {
		if def.IsDefinition() {
			if err := def.Upsert(ctx); err != nil {
				return err
			}
		}
		if def.IsMermaid() {
			if err := def.Mermaid(ctx); err != nil {
				return err
			}
		}
	}

	return nil
}

func New(config *boilingcore.Config) (*SchemanState, error) {
	s := &SchemanState{
		State: &boilingcore.State{
			Config: config,
		},
		Defs: map[Service]definition.Definition{},
	}

	s.Driver = drivers.GetDriver(config.DriverName)
	if err := s.initDBInfo(config.DriverConfig); err != nil {
		return nil, errors.Wrap(err, "unable to initialize tables")
	}

	ignoreAttrs := viper.GetStringSlice("attr-ignore")
	ignores := make(map[string]int, len(ignoreAttrs))
	for _, a := range ignoreAttrs {
		ignores[a]++
	}
	isIgnoreView := viper.GetBool("disable-views")
	commonInfo := &definition.CommonInfo{
		TablesByConnection: s.Tables,
		DriverName:         config.DriverName,
		IgnoreAttributes:   ignores,
		IsIgnoreView:       isIgnoreView,
		RawMermaid:         s.genMermaid(isIgnoreView),
	}

	return s.setDefinition(commonInfo)
}

func (s *SchemanState) setDefinition(info *definition.CommonInfo) (*SchemanState, error) {
	services := viper.GetStringSlice("services")
	for _, service := range services {
		switch service {
		case string(ServiceNotion):
			pageID := viper.GetString("notion-page-id")
			token := viper.GetString("notion-token")
			n, err := notion.NewNotion(
				pageID,
				viper.GetString("notion-table-index"),
				token,
				info,
			)
			if err != nil {
				return nil, err
			}
			s.Defs[ServiceNotion] = n
		case string(ServiceFile):
			definitionFile := viper.GetString("def-file")
			erdFile := viper.GetString("erd-file")
			s.Defs[ServiceFile] = file.NewFile(definitionFile, erdFile, info)
		default:
			return nil, errors.Errorf("The service have not been supported yet: %s", service)
		}
	}

	mermaidOutputs := viper.GetStringSlice("erd-outputs")
	for _, m := range mermaidOutputs {
		if d, ok := s.Defs[Service(m)]; ok {
			d.EnableMermaid()
			continue
		}

		switch m {
		case string(ServiceNotion):
			pageID := viper.GetString("notion-page-id")
			token := viper.GetString("notion-token")
			n, err := notion.NewNotionOnlyMermaid(
				pageID,
				viper.GetString("notion-table-index"),
				token,
				info,
			)
			if err != nil {
				return nil, err
			}
			s.Defs[ServiceNotion] = n
		case string(ServiceFile):
			erdFile := viper.GetString("erd-file")
			s.Defs[ServiceFile] = file.NewFileOnlyMermaid(erdFile, info)
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
