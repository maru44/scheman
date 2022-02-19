package db

import (
	"fmt"
	"strings"

	"github.com/friendsofgo/errors"
	"github.com/volatiletech/sqlboiler/v4/boilingcore"
	"github.com/volatiletech/sqlboiler/v4/drivers"
)

type SchemanState struct {
	*boilingcore.State
}

func New(config *boilingcore.Config) (*SchemanState, error) {
	s := &SchemanState{
		State: &boilingcore.State{
			Config: config,
		},
	}

	s.Driver = drivers.GetDriver(config.DriverName)
	if err := s.initDBInfo(config.DriverConfig); err != nil {
		return nil, errors.Wrap(err, "unable to initialize tables")
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
	fmt.Println(s.Tables)
	return nil
}
