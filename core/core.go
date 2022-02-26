package core

import (
	"context"
	"fmt"
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
		IgnoreAttributes []string
		Mermaid          string
		Defs             map[Service]definition.Definition
	}

	Service string
)

const (
	ServiceNotion = Service("NOTION")
)

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

	services := viper.GetStringSlice("services")
	for _, service := range services {
		switch service {
		case string(ServiceNotion):
			pageID := viper.GetString("notion-page-id")
			token := viper.GetString("notion-token")
			if pageID == "" {
				return nil, errors.New("notion-page-id is not set")
			}
			if token == "" {
				return nil, errors.New("notion-token is not set")
			}
			s.Defs[ServiceNotion] = notion.NewNotion(
				pageID,
				viper.GetString("notion-table-index"),
				token,
				s.Tables,
				config.DriverName,
				ignores,
				isIgnoreView,
			)
		default:
			return nil, errors.Errorf("The service have not been supported yet: %s", service)
		}
	}

	mermaidOutputs := viper.GetStringSlice("mermaid-outputs")
	if len(mermaidOutputs) != 0 {
		s.Mermaid = s.genMermaid()
	}
	for _, m := range mermaidOutputs {
		if d, ok := s.Defs[Service(m)]; ok {
			d.SetMermaid(s.Mermaid)
			continue
		}

		switch m {
		case string(ServiceNotion):
			pageID := viper.GetString("notion-page-id")
			token := viper.GetString("notion-token")
			if pageID == "" {
				return nil, errors.New("notion-page-id is not set")
			}
			if token == "" {
				return nil, errors.New("notion-token is not set")
			}
			s.Defs[ServiceNotion] = notion.NewNotion(
				pageID,
				viper.GetString("notion-table-index"),
				token,
				nil,
				config.DriverName,
				ignores,
				isIgnoreView,
			)
			s.Defs[ServiceNotion].SetMermaid(s.Mermaid)
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
		if err := def.Mermaid(ctx); err != nil {
			return err
		}
	}

	return nil
}

type (
	rel struct {
		TableName    string
		RelTableName string
		Nullable     bool
	}
)

func (s *SchemanState) genMermaid() string {
	var (
		fKeys, ones, manys []rel
		cols               string
		rels               string
	)
	for _, t := range s.Tables {
		for _, r := range t.FKeys {
			fKeys = append(fKeys, rel{
				TableName:    r.Table,
				RelTableName: r.ForeignTable,
				Nullable:     r.Nullable,
			})
		}
		for _, r := range t.ToOneRelationships {
			ones = append(ones, rel{
				TableName:    r.Table,
				RelTableName: r.ForeignTable,
				Nullable:     r.Nullable,
			})
		}
		for _, r := range t.ToManyRelationships {
			manys = append(manys, rel{
				TableName:    r.Table,
				RelTableName: r.ForeignTable,
				Nullable:     r.Nullable,
			})
		}

		cols += fmt.Sprintf("%s {\n", t.Name)
		for _, col := range t.Columns {
			c := definition.ConvertCol(col, t.PKey, s.Config.DriverName)
			cols += fmt.Sprintf("  %s %s\n", c.DBType, c.Name)
		}
		cols += "}\n\n"
	}

	for _, f := range fKeys {
		isOne := false
		for _, r := range ones {
			if f.TableName == r.RelTableName && f.RelTableName == r.TableName {
				switch f.Nullable {
				case true:
					if r.Nullable {
						rels += fmt.Sprintf("%s |o--o| %s : own\n", r.TableName, f.TableName)
					} else {
						rels += fmt.Sprintf("%s |o--|| %s : own\n", r.TableName, f.TableName)
					}
				case false:
					if r.Nullable {
						rels += fmt.Sprintf("%s ||--o| %s : own\n", r.TableName, f.TableName)
					} else {
						rels += fmt.Sprintf("%s ||--|| %s : own\n", r.TableName, f.TableName)
					}
				}
				isOne = true
				break
			}
		}
		if isOne {
			continue
		}
		for _, r := range manys {
			if f.TableName == r.RelTableName && f.RelTableName == r.TableName {
				switch f.Nullable {
				case true:
					if r.Nullable {
						rels += fmt.Sprintf("%s |o--o{ %s : has \n", r.TableName, f.TableName)
					} else {
						rels += fmt.Sprintf("%s |o--|{ %s : has \n", r.TableName, f.TableName)
					}
				case false:
					if r.Nullable {
						rels += fmt.Sprintf("%s ||--o{ %s : has \n", r.TableName, f.TableName)
					} else {
						rels += fmt.Sprintf("%s ||--|{ %s : has \n", r.TableName, f.TableName)
					}
				}
				isOne = true
				break
			}
		}
	}
	return "erDiagram\n\n" + rels + "\n\n" + cols
}
